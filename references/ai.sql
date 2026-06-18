-- PharmacyERP AI 子模块数据库初始化脚本
-- 目标：
-- 1. 新建 ai schema，供 AI 子模块读取 ERP 主数据与保存 AI 侧运行数据。
-- 2. public schema 仍是业务主数据唯一来源；ai schema 通过 VIEW 暴露药品、供应商、库存等只读数据。
-- 3. AI 子模块自身产生的搜索日志、发票识别任务、别名/同义词等数据保存在 ai schema。
--
-- 注意：
-- - 本脚本不会修改 public schema 的表结构。
-- - 发票转入库单仍由 ERP 后端完成，AI 子模块只返回识别草稿和匹配建议。
-- - 前端直连的药品搜索接口不鉴权；建议在网关侧做限流。
-- - 发票识别接口由 ERP 后端本地/内网直连，不鉴权；建议仅监听 localhost 或通过防火墙限制来源。

BEGIN;

CREATE SCHEMA IF NOT EXISTS ai;

COMMENT ON SCHEMA ai IS '智慧药店 ERP AI 子模块 schema；通过视图读取 public 主数据，并保存 AI 搜索、识别、日志等子模块数据';

-- 可选扩展：用于药品名称、厂家、规格等模糊匹配。
-- 如果部署用户没有创建扩展权限，可注释本行，应用层仍可使用 ILIKE/相似度算法实现搜索。
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- =========================================================
-- 1. 通用更新时间函数
-- =========================================================

CREATE OR REPLACE FUNCTION ai.set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- 2. AI 子模块自有表
-- =========================================================

CREATE TABLE IF NOT EXISTS ai.ai_request_log (
  id bigserial PRIMARY KEY,
  request_id varchar(100) NOT NULL DEFAULT (
    'AI-' || to_char(clock_timestamp(), 'YYYYMMDDHH24MISSMS') || '-' ||
    lpad((floor(random() * 1000000))::int::text, 6, '0')
  ),
  module varchar(50) NOT NULL,
  endpoint varchar(255),
  http_method varchar(20),
  caller_type varchar(30) NOT NULL DEFAULT 'UNKNOWN',
  client_ip varchar(64),
  user_agent text,
  request_json jsonb,
  response_json jsonb,
  success boolean,
  error_code varchar(100),
  error_message text,
  duration_ms int8,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_ai_request_log_request_id UNIQUE (request_id),
  CONSTRAINT chk_ai_request_log_caller_type CHECK (caller_type IN ('FRONTEND', 'ERP_BACKEND', 'SCHEDULER', 'UNKNOWN'))
);

COMMENT ON TABLE ai.ai_request_log IS 'AI 子模块通用请求日志；用于排查前端搜索、ERP 本地调用发票识别等请求';
COMMENT ON COLUMN ai.ai_request_log.caller_type IS 'FRONTEND=前端直连；ERP_BACKEND=ERP 后端本地/内网调用；SCHEDULER=定时任务';

CREATE INDEX IF NOT EXISTS idx_ai_request_log_module_created ON ai.ai_request_log (module, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ai_request_log_success_created ON ai.ai_request_log (success, created_at DESC);

CREATE TABLE IF NOT EXISTS ai.drug_search_alias (
  id bigserial PRIMARY KEY,
  drug_id int8 NOT NULL,
  alias_text varchar(200) NOT NULL,
  alias_type varchar(30) NOT NULL DEFAULT 'CUSTOM',
  weight numeric(6,3) NOT NULL DEFAULT 1.000,
  source varchar(50) NOT NULL DEFAULT 'MANUAL',
  status int2 NOT NULL DEFAULT 1,
  remark text,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamptz,
  CONSTRAINT chk_drug_search_alias_type CHECK (alias_type IN ('COMMON_ALIAS', 'TRADE_ALIAS', 'PINYIN', 'SYMPTOM', 'OCR_NAME', 'CUSTOM')),
  CONSTRAINT chk_drug_search_alias_status CHECK (status IN (0, 1))
);

COMMENT ON TABLE ai.drug_search_alias IS 'AI 药品搜索别名/同义词表；用于补充 public.drug_info 以外的搜索词，如商品名别名、拼音、症状词、OCR 常见误识别名';
COMMENT ON COLUMN ai.drug_search_alias.drug_id IS '对应 public.drug_info.id；不设外键，避免 AI 子模块强绑定 public 写入流程';

CREATE UNIQUE INDEX IF NOT EXISTS uk_drug_search_alias_active
  ON ai.drug_search_alias (drug_id, alias_text, alias_type)
  WHERE deleted_at IS NULL AND status = 1;
CREATE INDEX IF NOT EXISTS idx_drug_search_alias_drug_id ON ai.drug_search_alias (drug_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_drug_search_alias_text_trgm ON ai.drug_search_alias USING gin (alias_text gin_trgm_ops) WHERE deleted_at IS NULL AND status = 1;

DROP TRIGGER IF EXISTS trg_drug_search_alias_updated_at ON ai.drug_search_alias;
CREATE TRIGGER trg_drug_search_alias_updated_at
BEFORE UPDATE ON ai.drug_search_alias
FOR EACH ROW EXECUTE FUNCTION ai.set_updated_at();

CREATE TABLE IF NOT EXISTS ai.drug_search_log (
  id bigserial PRIMARY KEY,
  request_id varchar(100),
  query_text text NOT NULL,
  search_mode varchar(30) NOT NULL DEFAULT 'HYBRID',
  filters jsonb,
  result_count int4 NOT NULL DEFAULT 0,
  top_result_drug_id int8,
  latency_ms int8,
  client_ip varchar(64),
  user_agent text,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_drug_search_log_mode CHECK (search_mode IN ('KEYWORD', 'FUZZY', 'SEMANTIC', 'HYBRID'))
);

COMMENT ON TABLE ai.drug_search_log IS '前端直连药品智能搜索日志；用于后续优化关键词、别名、排序和召回策略';

CREATE INDEX IF NOT EXISTS idx_drug_search_log_created ON ai.drug_search_log (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_drug_search_log_query_trgm ON ai.drug_search_log USING gin (query_text gin_trgm_ops);

CREATE TABLE IF NOT EXISTS ai.drug_search_feedback (
  id bigserial PRIMARY KEY,
  request_id varchar(100),
  query_text text NOT NULL,
  selected_drug_id int8,
  feedback_type varchar(30) NOT NULL,
  feedback_note text,
  client_ip varchar(64),
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_drug_search_feedback_type CHECK (feedback_type IN ('CLICK', 'SELECT', 'NO_RESULT', 'BAD_RESULT', 'MANUAL_CORRECTION'))
);

COMMENT ON TABLE ai.drug_search_feedback IS '药品搜索反馈；用于记录用户点击、选择、无结果、结果不准等行为，辅助改进搜索';

CREATE INDEX IF NOT EXISTS idx_drug_search_feedback_created ON ai.drug_search_feedback (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_drug_search_feedback_selected_drug ON ai.drug_search_feedback (selected_drug_id);

CREATE TABLE IF NOT EXISTS ai.invoice_recognition_job (
  id bigserial PRIMARY KEY,
  request_id varchar(100) NOT NULL DEFAULT (
    'INVAI-' || to_char(clock_timestamp(), 'YYYYMMDDHH24MISSMS') || '-' ||
    lpad((floor(random() * 1000000))::int::text, 6, '0')
  ),
  erp_request_id varchar(100),
  erp_file_id varchar(100),
  file_name varchar(255),
  content_type varchar(100),
  file_size int8,
  file_hash varchar(128),
  status varchar(20) NOT NULL DEFAULT 'PENDING',
  ai_provider varchar(50),
  ai_model varchar(100),
  recognized_supplier_name varchar(255),
  matched_supplier_id int8,
  invoice_no varchar(100),
  invoice_date date,
  confidence numeric(5,4),
  result_json jsonb,
  raw_response_json jsonb,
  warnings_json jsonb,
  error_code varchar(100),
  error_message text,
  started_at timestamptz,
  finished_at timestamptz,
  duration_ms int8,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamptz,
  CONSTRAINT uk_invoice_recognition_job_request_id UNIQUE (request_id),
  CONSTRAINT chk_invoice_recognition_job_status CHECK (status IN ('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED')),
  CONSTRAINT chk_invoice_recognition_job_confidence CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 1))
);

COMMENT ON TABLE ai.invoice_recognition_job IS 'AI 服务侧发票识别任务；ERP 主库的 ai_invoice_record 仍是业务记录，本表仅保存 AI 服务运行态、原始响应、模型信息和排查数据';
COMMENT ON COLUMN ai.invoice_recognition_job.result_json IS '规范化识别结果，结构对应 ai-openapi.yaml 的 InvoiceRecognizeResult';
COMMENT ON COLUMN ai.invoice_recognition_job.raw_response_json IS 'OCR/LLM/视觉模型原始响应，便于排查';
COMMENT ON COLUMN ai.invoice_recognition_job.warnings_json IS '金额不一致、低置信度、未匹配药品等质检提示';

CREATE INDEX IF NOT EXISTS idx_invoice_recognition_job_status_created ON ai.invoice_recognition_job (status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_invoice_recognition_job_file_hash ON ai.invoice_recognition_job (file_hash) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_invoice_recognition_job_invoice_no ON ai.invoice_recognition_job (invoice_no) WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS trg_invoice_recognition_job_updated_at ON ai.invoice_recognition_job;
CREATE TRIGGER trg_invoice_recognition_job_updated_at
BEFORE UPDATE ON ai.invoice_recognition_job
FOR EACH ROW EXECUTE FUNCTION ai.set_updated_at();

CREATE TABLE IF NOT EXISTS ai.invoice_recognition_item_cache (
  id bigserial PRIMARY KEY,
  job_id int8 NOT NULL REFERENCES ai.invoice_recognition_job(id) ON DELETE CASCADE,
  row_no int4 NOT NULL,
  drug_name varchar(200),
  specification varchar(100),
  manufacturer varchar(200),
  batch_number varchar(100),
  expire_date date,
  quantity numeric(12,3),
  unit_price numeric(12,4),
  amount numeric(12,2),
  matched_drug_id int8,
  confidence numeric(5,4),
  candidates_json jsonb,
  warnings_json jsonb,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_invoice_recognition_item_cache_job_row UNIQUE (job_id, row_no),
  CONSTRAINT chk_invoice_recognition_item_cache_confidence CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 1))
);

COMMENT ON TABLE ai.invoice_recognition_item_cache IS '发票识别明细缓存；便于后续人工排查、模型评估和药品匹配优化';

CREATE INDEX IF NOT EXISTS idx_invoice_recognition_item_cache_job ON ai.invoice_recognition_item_cache (job_id);
CREATE INDEX IF NOT EXISTS idx_invoice_recognition_item_cache_matched_drug ON ai.invoice_recognition_item_cache (matched_drug_id);
CREATE INDEX IF NOT EXISTS idx_invoice_recognition_item_cache_drug_name_trgm ON ai.invoice_recognition_item_cache USING gin (drug_name gin_trgm_ops);

CREATE TABLE IF NOT EXISTS ai.ai_model_config (
  id bigserial PRIMARY KEY,
  module varchar(50) NOT NULL,
  provider varchar(50) NOT NULL,
  model_name varchar(100) NOT NULL,
  config_json jsonb,
  enabled boolean NOT NULL DEFAULT true,
  priority int4 NOT NULL DEFAULT 100,
  remark text,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamptz,
  CONSTRAINT uk_ai_model_config_module_provider_model UNIQUE (module, provider, model_name)
);

COMMENT ON TABLE ai.ai_model_config IS 'AI 模型配置表；用于记录发票识别、药品搜索等模块使用的模型/供应商配置，不保存密钥明文';

CREATE INDEX IF NOT EXISTS idx_ai_model_config_module_enabled ON ai.ai_model_config (module, enabled, priority) WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS trg_ai_model_config_updated_at ON ai.ai_model_config;
CREATE TRIGGER trg_ai_model_config_updated_at
BEFORE UPDATE ON ai.ai_model_config
FOR EACH ROW EXECUTE FUNCTION ai.set_updated_at();

-- =========================================================
-- 3. public 主数据只读视图
-- =========================================================

CREATE OR REPLACE VIEW ai.v_supplier_source AS
SELECT
  s.id AS supplier_id,
  s.supplier_code,
  s.name,
  s.contact_name,
  s.contact_phone,
  s.license_no,
  s.address,
  s.status,
  s.updated_at
FROM public.supplier s
WHERE s.deleted_at IS NULL;

COMMENT ON VIEW ai.v_supplier_source IS 'AI 子模块供应商只读视图；用于发票供应商名称匹配，主数据仍以 public.supplier 为准';

CREATE OR REPLACE VIEW ai.v_active_supplier_source AS
SELECT *
FROM ai.v_supplier_source
WHERE status = 1;

COMMENT ON VIEW ai.v_active_supplier_source IS '启用供应商视图；发票转入库单候选 supplier_id 只能来自启用供应商';

CREATE OR REPLACE VIEW ai.v_location_source AS
SELECT
  l.id AS location_id,
  l.location_code,
  l.location_name,
  l.area,
  l.shelf,
  l.layer,
  l.position,
  l.status,
  l.updated_at
FROM public.location_info l
WHERE l.deleted_at IS NULL;

COMMENT ON VIEW ai.v_location_source IS 'AI 子模块货位只读视图；供后续库存问答、盘库异常解释等能力使用';

CREATE OR REPLACE VIEW ai.v_drug_base_source AS
SELECT
  d.id AS drug_id,
  d.drug_code,
  d.common_name,
  d.trade_name,
  d.specification,
  d.dosage_form,
  d.manufacturer,
  d.approval_number,
  d.barcode,
  d.unit,
  d.retail_price,
  d.purchase_price,
  d.storage_condition,
  COALESCE(d.is_prescription, false) AS is_prescription,
  COALESCE(d.is_medicare, false) AS is_medicare,
  d.status,
  d.updated_at,
  concat_ws(' ',
    d.drug_code,
    d.common_name,
    d.trade_name,
    d.specification,
    d.dosage_form,
    d.manufacturer,
    d.approval_number,
    d.barcode,
    d.storage_condition
  ) AS base_search_text
FROM public.drug_info d
WHERE d.deleted_at IS NULL;

COMMENT ON VIEW ai.v_drug_base_source IS 'AI 子模块药品基础只读视图；包含启用/停用药品，搜索接口默认只返回启用药品';

CREATE OR REPLACE VIEW ai.v_drug_inventory_summary AS
WITH inv AS (
  SELECT
    i.drug_id,
    i.trace_code,
    i.status,
    i.expire_date
  FROM public.drug_trace_inventory i
  WHERE i.deleted_at IS NULL
), reserved AS (
  SELECT DISTINCT r.trace_code
  FROM public.trace_reservation r
  WHERE r.status = 'RESERVED'
    AND r.deleted_at IS NULL
)
SELECT
  inv.drug_id,
  COUNT(*) FILTER (WHERE inv.status = 'PENDING')::int4 AS pending_qty,
  COUNT(*) FILTER (WHERE inv.status = 'IN_STOCK')::int4 AS in_stock_qty,
  COUNT(*) FILTER (WHERE inv.status = 'IN_STOCK' AND reserved.trace_code IS NOT NULL)::int4 AS reserved_qty,
  COUNT(*) FILTER (WHERE inv.status = 'IN_STOCK' AND reserved.trace_code IS NULL)::int4 AS available_qty,
  COUNT(*) FILTER (WHERE inv.status = 'SOLD')::int4 AS sold_qty,
  COUNT(*) FILTER (WHERE inv.status IN ('MISPLACED', 'LOSS_CANDIDATE', 'LOST'))::int4 AS abnormal_qty,
  COUNT(*) FILTER (
    WHERE inv.status = 'IN_STOCK'
      AND reserved.trace_code IS NULL
      AND inv.expire_date <= CURRENT_DATE + INTERVAL '30 days'
  )::int4 AS near_expire_available_qty,
  MIN(inv.expire_date) FILTER (WHERE inv.status = 'IN_STOCK' AND reserved.trace_code IS NULL) AS nearest_expire_date
FROM inv
LEFT JOIN reserved ON reserved.trace_code = inv.trace_code
GROUP BY inv.drug_id;

COMMENT ON VIEW ai.v_drug_inventory_summary IS '药品库存汇总视图；可售库存=IN_STOCK 且排除 RESERVED 预占';

CREATE OR REPLACE VIEW ai.v_drug_search_source AS
WITH alias_agg AS (
  SELECT
    a.drug_id,
    string_agg(a.alias_text, ' ' ORDER BY a.weight DESC, a.id ASC) AS alias_text
  FROM ai.drug_search_alias a
  WHERE a.deleted_at IS NULL
    AND a.status = 1
  GROUP BY a.drug_id
)
SELECT
  d.drug_id,
  d.drug_code,
  d.common_name,
  d.trade_name,
  d.specification,
  d.dosage_form,
  d.manufacturer,
  d.approval_number,
  d.barcode,
  d.unit,
  d.retail_price,
  d.purchase_price,
  d.storage_condition,
  d.is_prescription,
  d.is_medicare,
  d.status,
  COALESCE(s.pending_qty, 0) AS pending_qty,
  COALESCE(s.in_stock_qty, 0) AS in_stock_qty,
  COALESCE(s.reserved_qty, 0) AS reserved_qty,
  COALESCE(s.available_qty, 0) AS available_qty,
  COALESCE(s.sold_qty, 0) AS sold_qty,
  COALESCE(s.abnormal_qty, 0) AS abnormal_qty,
  COALESCE(s.near_expire_available_qty, 0) AS near_expire_available_qty,
  s.nearest_expire_date,
  COALESCE(a.alias_text, '') AS alias_text,
  concat_ws(' ', d.base_search_text, COALESCE(a.alias_text, '')) AS search_text,
  d.updated_at
FROM ai.v_drug_base_source d
LEFT JOIN ai.v_drug_inventory_summary s ON s.drug_id = d.drug_id
LEFT JOIN alias_agg a ON a.drug_id = d.drug_id
WHERE d.status = 1;

COMMENT ON VIEW ai.v_drug_search_source IS '前端直连药品智能搜索数据源视图；只返回启用药品，并附带可售库存、近效期、别名搜索文本';

CREATE OR REPLACE VIEW ai.v_invoice_drug_match_source AS
SELECT
  drug_id,
  drug_code,
  common_name,
  trade_name,
  specification,
  dosage_form,
  manufacturer,
  approval_number,
  barcode,
  unit,
  purchase_price,
  retail_price,
  is_prescription,
  is_medicare,
  alias_text,
  search_text,
  updated_at
FROM ai.v_drug_search_source;

COMMENT ON VIEW ai.v_invoice_drug_match_source IS '发票识别明细匹配药品的数据源视图；由 AI 服务读取候选药品，但最终 drug_id 仍需 ERP/前端确认';

CREATE OR REPLACE VIEW ai.v_drug_batch_inventory_source AS
SELECT
  i.drug_id,
  d.drug_code,
  d.common_name,
  d.specification,
  d.manufacturer,
  i.batch_number,
  i.expire_date,
  i.location_id,
  l.location_code,
  l.location_name,
  COUNT(*) FILTER (WHERE i.status = 'IN_STOCK')::int4 AS in_stock_qty,
  COUNT(*) FILTER (WHERE i.status = 'PENDING')::int4 AS pending_qty,
  COUNT(*) FILTER (WHERE i.status IN ('MISPLACED', 'LOSS_CANDIDATE', 'LOST'))::int4 AS abnormal_qty
FROM public.drug_trace_inventory i
JOIN public.drug_info d ON d.id = i.drug_id AND d.deleted_at IS NULL
LEFT JOIN public.location_info l ON l.id = i.location_id AND l.deleted_at IS NULL
WHERE i.deleted_at IS NULL
GROUP BY
  i.drug_id,
  d.drug_code,
  d.common_name,
  d.specification,
  d.manufacturer,
  i.batch_number,
  i.expire_date,
  i.location_id,
  l.location_code,
  l.location_name;

COMMENT ON VIEW ai.v_drug_batch_inventory_source IS '药品批号/有效期/货位库存视图；供后续库存问答、近效期建议、搜索结果详情使用';

-- =========================================================
-- 4. 初始模型配置占位数据，可按实际部署修改
-- =========================================================

INSERT INTO ai.ai_model_config (module, provider, model_name, config_json, enabled, priority, remark)
VALUES
  ('INVOICE_RECOGNITION', 'LOCAL_OCR', 'default-invoice-ocr', '{"type":"placeholder"}'::jsonb, true, 100, '占位配置：发票 OCR/结构化识别模型'),
  ('DRUG_SEARCH', 'LOCAL_SEARCH', 'hybrid-keyword-fuzzy', '{"type":"keyword+trgm"}'::jsonb, true, 100, '占位配置：药品智能搜索，本地关键词/模糊匹配')
ON CONFLICT (module, provider, model_name) DO NOTHING;

COMMIT;
