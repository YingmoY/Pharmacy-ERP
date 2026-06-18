# 智慧药店 ERP 业务规则说明（business.md）

> 本文档用于规范后端编码、AI 辅助编码与业务测试时应实现的核心业务逻辑。  
> 数据库结构以 `public.sql` 为准，接口契约以 `openapi.yaml` 为准，权限模型以 `casbin_model.conf` 为准。  
> 当数据库字段与 API DTO 存在轻微差异时，业务实现优先满足 API 暴露的 DTO 和本文档明确规则；API 未暴露且无法由已有数据推算的字段，当前版本不主动维护。

---

## 1. 系统定位与边界

### 1.1 系统定位

本系统是智慧药店 ERP 主系统，负责药店内部的基础资料、采购入库、追溯码库存、上架、销售、药师审核、退货退款、盘库、预警、报表、权限与审计等业务。

系统以“药品追溯码”为库存最小业务单位：

- 一个追溯码对应唯一一盒药。
- 一条销售明细对应一盒药、一个追溯码。
- `drug_trace_inventory` 只保存追溯码当前库存状态。
- 销售预占、释放、消耗以 `trace_reservation` 为唯一业务来源。
- 追溯轨迹以 `drug_trace_log` 记录全链路动作。

### 1.2 外部模块边界

主 ERP 不实现以下模块内部逻辑：

1. AI 发票识别子模块  
   主 ERP 只负责上传文件、调用外部 AI 服务、保存识别结果、展示给前端修正、将确认后的结果转入库单。

2. 医保网关子模块  
   主 ERP 可以保存销售单中的医保相关字段，但不在本文档中定义医保网关内部交互流程。医保网关接口、签名、交易重试等由独立模块负责。

---

## 2. 全局业务约定

### 2.1 当前用户与操作人

所有需要记录操作人的业务字段都由后端根据当前 JWT 登录用户生成，前端不得传入或覆盖。

典型字段包括：

- `operator_id`
- `creator_id`
- `cashier_id`
- `reserved_by`
- `submitter_id`
- `pharmacist_id`
- `refund_operator_id`
- `resolved_by`
- `ignored_by`
- `requested_by`

如果请求体中出现这些字段，后端应忽略或拒绝，具体以接口 DTO 为准。编码时不得信任前端传入的操作人身份。

### 2.2 软删除规则

业务主表通常包含 `deleted_at` 字段。默认删除行为为软删除。

软删除规则：

- 查询列表和详情默认过滤 `deleted_at IS NULL`。
- 被业务引用的数据，不允许直接删除。
- 删除药品、供应商、货位、角色等基础资料时，应先检查是否存在有效业务引用。
- 若存在有效引用，应返回业务冲突错误，不得强制删除。

### 2.3 启用/停用规则

基础资料状态字段通常使用：

- `1`：启用
- `0`：停用

停用后的业务限制：

- 停用药品不可新建入库单明细，不可新建销售明细。
- 停用供应商不可新建入库单。
- 停用货位不可执行上架、调拨、盘库扫码作为目标货位。
- 停用用户不可登录，不可作为新业务负责人或执行人。

历史业务数据不因基础资料停用而失效。

### 2.4 单号生成规则

业务单号由后端生成，前端不得传入。

建议格式：

| 业务 | 单号格式 |
|---|---|
| 入库单 | `IN-YYYYMMDD-XXXX` |
| 销售单 | `SO-YYYYMMDD-XXXX` |
| 药师审核单 | `REV-YYYYMMDD-XXXX` |
| 盘库任务 | `INV-YYYYMMDD-XXXX` |
| 扫码任务 | `SCAN-YYYYMMDD-XXXX` |
| 库存调整 | `ADJ-YYYYMMDD-XXXX` |
| 预占单 | `RSV-YYYYMMDD-XXXX` |
| 供应商编码 | `SUP-XXXX` 或业务指定编码 |
| 报表导出任务 | 全局唯一 `task_id` |

同一天内流水号应递增且避免并发重复。实现时应使用数据库唯一约束、序列、事务锁或专门的单号生成器保证唯一性。

### 2.5 金额计算规则

金额以数据库 `numeric` / 后端 decimal 类型计算，不得使用浮点数。

通用规则：

- 入库明细金额：`amount = planned_qty * unit_price`。
- 入库单总金额：所有未删除明细金额合计。
- 销售明细数量固定为 `1`。
- 销售明细单价以后端药品主数据 `drug_info.retail_price` 为准。
- 前端传入的 `unit_price` 只能作为兼容字段，不得决定最终售价。
- 如果前端传入 `unit_price` 与当前零售价不一致，后端应以后端价格覆盖，或直接返回价格不一致错误。推荐返回错误，提示前端刷新价格。
- 销售单 `total_amount` 为所有未删除、未退货明细价格合计。
- `discount_amount` 当前只作为销售单级优惠字段；没有专门改价权限前，不允许前端任意改变明细单价。
- `actual_amount = total_amount - discount_amount`，不得小于 `0`。
- 退货金额按被退明细的实际成交金额计算。若存在订单级优惠，当前版本可按明细原价占比分摊，或直接按明细价格退回；具体实现应保持前后一致。

### 2.6 事务边界

下列业务必须在同一数据库事务内完成：

- 创建销售单、创建销售明细、创建追溯码预占。
- 删除销售明细、释放对应追溯码预占、重算销售单金额。
- 销售结算、消耗预占、更新追溯码为 `SOLD`、写追溯日志。
- 销售取消、释放预占、更新销售单状态。
- 销售退货、更新销售明细退款状态、恢复追溯码库存、更新销售单退款状态、写追溯日志。
- 入库扫码确认、创建追溯码库存、更新入库明细确认数量、写追溯日志。
- 完成入库单时校验全部明细已确认，并更新入库单状态。
- 取消入库单时删除未上架追溯码、更新入库单状态。
- 上架、调拨、盘库错架处理、盘亏确认/驳回、库存调整及对应追溯日志。
- 药师审核通过/驳回与销售单状态变更。

事务内涉及追溯码、销售明细、预占记录时，应对相关行加锁，避免并发重复销售、重复退货、重复扫码。

### 2.7 幂等与重复操作

业务实现应尽量保证重复请求不会破坏数据一致性。

扫码相关规则：

- 入库扫码重复扫描同一追溯码，应返回重复扫码错误或幂等成功，但不得重复增加 `confirmed_qty`。
- 上架重复扫描已上架追溯码，应返回状态错误，不得再次写入上架。
- 盘库任务内重复扫描同一追溯码，应记录为重复或返回重复错误，不得重复计数。
- 销售预占同一追溯码时，只允许存在一条有效 `RESERVED` 记录。
- 退货同一销售明细只能退一次。

### 2.8 统一响应与错误处理

所有 JSON 响应使用统一包装：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "request_id": "..."
}
```

建议错误语义：

| 场景 | 建议 HTTP 状态 | 业务含义 |
|---|---:|---|
| 参数格式错误 | 400 | 请求体、路径参数、查询参数非法 |
| 未登录或 token 无效 | 401 | 认证失败 |
| 权限不足 | 403 | RBAC 校验不通过 |
| 资源不存在 | 404 | 业务对象不存在或已删除 |
| 状态冲突/重复编码/重复追溯码 | 409 | 当前状态不允许该操作 |
| 外部 AI/医保服务失败 | 502/503 | 外部依赖不可用或返回失败 |
| 服务器内部错误 | 500 | 未预期错误 |

---

## 3. 核心状态字典

### 3.1 入库单状态

| 状态 | 含义 | 允许操作 |
|---|---|---|
| `DRAFT` | 草稿 | 修改基本信息、增删改明细、提交、取消 |
| `PENDING_CONFIRM` | 待扫码确认 | 扫码确认、完成、取消 |
| `COMPLETED` | 已完成 | 查询、后续上架，不可取消 |
| `CANCELLED` | 已取消 | 仅查询 |

状态流转：

```text
DRAFT -> PENDING_CONFIRM -> COMPLETED
DRAFT -> CANCELLED
PENDING_CONFIRM -> CANCELLED
```

### 3.2 追溯码库存状态

| 状态 | 含义 |
|---|---|
| `PENDING` | 已入库确认，待上架 |
| `IN_STOCK` | 已上架，在库 |
| `SOLD` | 已销售出库 |
| `MISPLACED` | 盘库发现错架 |
| `LOSS_CANDIDATE` | 盘库未扫到，盘亏候选 |
| `LOST` | 已确认盘亏/丢失 |

状态流转示意：

```text
PENDING -> IN_STOCK -> SOLD
SOLD -> IN_STOCK              # 退货
IN_STOCK -> MISPLACED -> IN_STOCK
IN_STOCK -> LOSS_CANDIDATE -> LOST
LOSS_CANDIDATE -> IN_STOCK    # 驳回盘亏候选
```

注意：销售预占不改变 `drug_trace_inventory.status`。被预占的追溯码仍保持 `IN_STOCK`，但不可再次销售。

### 3.3 销售单状态

| 状态 | 含义 | 允许操作 |
|---|---|---|
| `PENDING` | 待结算 | 增删明细、结算、取消 |
| `PENDING_REVIEW` | 待药师审核 | 审核、取消；不允许修改明细 |
| `APPROVED` | 药师已审核通过 | 结算、取消 |
| `COMPLETED` | 已完成销售 | 退货 |
| `PARTIALLY_REFUNDED` | 部分退货 | 继续退剩余未退明细 |
| `REFUNDED` | 全单退货 | 仅查询 |
| `CANCELLED` | 已取消 | 仅查询 |

状态流转：

```text
PENDING -> COMPLETED
PENDING -> CANCELLED
PENDING_REVIEW -> APPROVED -> COMPLETED
PENDING_REVIEW -> CANCELLED
APPROVED -> CANCELLED
COMPLETED -> PARTIALLY_REFUNDED -> REFUNDED
COMPLETED -> REFUNDED
```

### 3.4 药师审核状态

| 状态 | 含义 |
|---|---|
| `PENDING` | 待审核 |
| `APPROVED` | 审核通过 |
| `REJECTED` | 审核驳回 |
| `CANCELLED` | 审核取消 |

审核驳回后，销售单直接变为 `CANCELLED`，并释放所有有效预占。

### 3.5 追溯码预占状态

| 状态 | 含义 |
|---|---|
| `RESERVED` | 已预占，尚未结算 |
| `RELEASED` | 已主动释放 |
| `CONSUMED` | 销售结算后已消耗 |
| `EXPIRED` | 超时自动释放 |

同一 `trace_code` 在 `RESERVED` 状态下只能存在一条有效记录。

### 3.6 盘库任务状态

| 状态 | 含义 |
|---|---|
| `PENDING` | 待开始 |
| `IN_PROGRESS` | 进行中 |
| `COMPLETED` | 已完成 |
| `CANCELLED` | 已取消 |

### 3.7 扫码任务状态

| 状态 | 含义 |
|---|---|
| `PENDING` | 待开始 |
| `IN_PROGRESS` | 进行中 |
| `COMPLETED` | 已完成 |
| `CANCELLED` | 已取消 |

扫码任务类型：

| 类型 | 含义 |
|---|---|
| `INBOUND` | 入库扫码 |
| `SHELVING` | 上架扫码 |
| `INVENTORY` | 盘库扫码 |

### 3.8 AI 发票识别状态

| 状态 | 含义 |
|---|---|
| `PENDING` | 待处理 |
| `PROCESSING` | 识别中 |
| `COMPLETED` | 识别完成 |
| `FAILED` | 识别失败 |

只有 `COMPLETED` 状态的识别记录允许转换为入库单。

### 3.9 预警状态映射

数据库 `audit_event.status` 与 API `AlertInfo.status` 映射：

| 数据库值 | API 状态 | 含义 |
|---:|---|---|
| `0` | `ACTIVE` | 待处理 |
| `1` | `RESOLVED` | 已处理 |
| `2` | `IGNORED` | 已忽略 |

---

## 4. 基础资料业务

## 4.1 药品资料

药品基础资料对应 `drug_info`。

核心规则：

- `drug_code` 全局唯一，不允许重复。
- `drug_code` 创建后不建议修改。
- `common_name`、`specification`、`manufacturer` 必填。
- `is_prescription` 表示是否处方药，是销售是否需要药师审核的重要依据。
- `retail_price` 是销售单价的后端来源。
- `status=0` 的药品不可用于新入库、新销售。
- 若药品存在有效库存、销售明细、入库明细等引用，不允许删除。

库存统计：

- 当前库存数量一般统计 `drug_trace_inventory.status = IN_STOCK` 的追溯码数量。
- 可售库存数量必须排除有效预占。
- 低库存阈值为 `3`，即可售库存数量小于等于 `3` 时触发低库存预警。

## 4.2 供应商资料

供应商基础资料对应 `supplier`。

核心规则：

- `supplier_code` 必填，且应唯一。
- `name` 必填。
- `status=0` 的供应商不可创建新入库单。
- 已被入库单引用的供应商不允许删除。
- AI 发票识别可能识别供应商名称，但转入库单时必须明确选择 `supplier_id`，不得仅使用供应商名称。

如果当前 OpenAPI 中 `CreateSupplierRequest` 尚未将 `supplier_code` 标记为 required，后端仍应按本文档校验 `supplier_code` 必填；后续建议同步修正 OpenAPI。

## 4.3 货位资料

货位基础资料对应 `location_info`。

核心规则：

- `location_code` 全局唯一。
- `location_name`、`area` 必填。
- 停用货位不可作为上架、调拨、盘库扫码目标货位。
- 若货位下存在 `IN_STOCK`、`MISPLACED`、`LOSS_CANDIDATE` 等未结束库存，不允许删除。
- 当前版本不做货位容量校验，`capacity` 字段暂不作为业务限制。

货位混放规则：

- 同一货位可以放多种药品。
- `/shelving/mix-check` 只返回提示，不强制禁止上架。
- 混放提示可基于该货位当前已有 `IN_STOCK` 药品种类、批号、有效期等信息生成。

## 4.4 文件资料

文件用于发票图片、PDF、报表导出文件等。

核心规则：

- 文件上传后生成 `file_id`。
- 文件业务归属通过业务表保存，如 `ai_invoice_record.file_id`、`report_export_task.file_id`。
- 文件删除不应破坏已存在业务记录；若需要清理，只做文件不可访问或后台清理。

---

## 5. AI 发票识别业务

### 5.1 业务流程

流程：

```text
上传发票文件 -> 调用外部 AI 子模块 -> 保存 ai_invoice_record -> 前端展示识别结果 -> 人工修正 -> 转入库单
```

主 ERP 的职责：

- 接收并保存发票文件。
- 调用外部 AI 发票识别服务。
- 保存识别状态、规范化结果 JSON、原始响应 JSON、错误信息。
- 提供识别记录列表与详情。
- 将人工确认后的识别结果转换为入库单草稿。

主 ERP 不负责：

- OCR 模型实现。
- LLM 提示词内部逻辑。
- AI 服务内部重试和调度策略。

### 5.2 识别记录

识别记录保存到 `ai_invoice_record`。

状态规则：

- 创建识别任务时状态为 `PENDING` 或 `PROCESSING`。
- 外部 AI 返回成功后状态为 `COMPLETED`。
- 外部 AI 返回失败、超时或解析失败后状态为 `FAILED`，记录 `error_message`。

识别结果：

- `result_json` 保存规范化后的 `InvoiceRecognizeResult`。
- `raw_response_json` 保存外部 AI 原始响应，便于排查。
- `recognized_supplier_name` 只作为参考。
- `matched_supplier_id` 如果无法唯一匹配供应商，可以为空。

### 5.3 转入库单

只有 `COMPLETED` 状态且尚未转换过的识别记录允许转入库单。

转换规则：

- 请求必须显式提交 `supplier_id`。
- `supplier_id` 必须指向启用状态的供应商。
- 请求必须提交人工确认后的 `items`。
- 每条 item 必须包含 `drug_id`、`batch_number`、`expire_date`、`planned_qty`、`unit_price`。
- `drug_id` 必须指向启用状态的药品。
- 转换成功后创建 `inbound_order(DRAFT)` 和对应 `inbound_order_detail`。
- 转换成功后回写 `ai_invoice_record.inbound_order_id` 和 `converted_at`。
- 同一张发票识别记录不允许多次转入库单。

如果用户需要重新转换，应取消或删除原入库单后，由后端提供专门的重置能力；当前版本不实现重复转换。

---

## 6. 入库业务

### 6.1 入库单创建

入库单对应 `inbound_order`，明细对应 `inbound_order_detail`。

创建规则：

- 创建后状态为 `DRAFT`。
- `supplier_id` 必填，且必须为启用供应商。
- `operator_id` / `creator_id` 由后端从 JWT 当前用户写入。
- 可在创建时同时提交明细，也可创建后单独添加明细。
- 明细中的药品必须启用。
- `planned_qty` 必须大于 `0`。
- `confirmed_qty` 初始为 `0`。
- `total_amount` 根据明细金额汇总。

### 6.2 入库单修改

只有 `DRAFT` 状态允许修改：

- 供应商
- 发票号
- 备注
- 新增明细
- 修改明细
- 删除明细

入库单进入 `PENDING_CONFIRM` 后，只允许：

- 扫码确认追溯码
- 完成入库
- 取消入库单

不得再修改供应商、发票号、明细、计划数量、批号、有效期、单价。

### 6.3 入库单提交

提交规则：

- 只有 `DRAFT` 状态允许提交。
- 提交前必须至少有一条有效明细。
- 每条明细必须满足药品、批号、有效期、计划数量、单价合法。
- 提交后状态变为 `PENDING_CONFIRM`。
- 写入提交时间 `submitted_at`。

### 6.4 追溯码扫码确认

扫码确认接口用于将实际收到的一盒药与入库单明细绑定。

确认规则：

- 只有 `PENDING_CONFIRM` 状态的入库单允许扫码确认。
- 请求必须指定 `detail_id` 和 `trace_code`。
- `detail_id` 必须属于当前入库单。
- 追溯码必须全局唯一，不得已存在有效 `drug_trace_inventory` 记录。
- 同一追溯码重复扫码不得重复增加确认数量。
- 明细 `confirmed_qty` 不得超过 `planned_qty`。
- 每扫码成功一个追溯码，创建一条 `drug_trace_inventory`：
  - `trace_code` 为扫码值
  - `drug_id`、`batch_number`、`expire_date` 来自入库明细
  - `status = PENDING`
  - `location_id = NULL`
  - `inbound_order_id`、`inbound_detail_id` 绑定当前入库单和明细
- 同时写入 `drug_trace_log(action_type=INBOUND)`。

### 6.5 完成入库

完成规则：

- 只有 `PENDING_CONFIRM` 状态允许完成。
- 必须所有明细全部扫码完成，即每条明细 `confirmed_qty == planned_qty`。
- 不允许少扫完成。
- 完成后入库单状态变为 `COMPLETED`。
- 写入 `completed_at`。
- 已确认追溯码保持 `PENDING`，等待上架。

### 6.6 取消入库单

取消规则：

- `DRAFT` 状态可以直接取消。
- `PENDING_CONFIRM` 状态可以取消，但前提是该入库单下没有追溯码已经上架为 `IN_STOCK` 或进入后续库存状态。
- 一旦任意追溯码已经上架为 `IN_STOCK`，不允许取消入库单。
- 取消 `PENDING_CONFIRM` 入库单时，应删除或软删除该入库单产生的 `PENDING` 追溯码库存记录。
- 取消后入库单状态变为 `CANCELLED`，写入 `cancelled_at`。
- 取消时应写操作日志。

建议取消校验：

```text
允许取消 = 入库单状态 in (DRAFT, PENDING_CONFIRM)
        AND 不存在该入库单关联的 drug_trace_inventory.status != PENDING
```

---

## 7. 上架与货位业务

### 7.1 待上架列表

待上架列表来自：

```text
drug_trace_inventory.status = PENDING
AND inbound_order.status = COMPLETED
```

只有所属入库单已完成的追溯码才允许上架。

### 7.2 单个上架

上架规则：

- 请求包含 `trace_code` 和 `location_code`。
- 追溯码必须存在，状态必须为 `PENDING`。
- 所属入库单必须为 `COMPLETED`。
- 货位必须存在且启用。
- 当前版本允许混放，只做提示，不阻止上架。
- 上架成功后：
  - `drug_trace_inventory.status = IN_STOCK`
  - `drug_trace_inventory.location_id = 目标货位 ID`
  - `last_action = SHELVING`
  - 写入 `drug_trace_log(action_type=SHELVING, from_status=PENDING, to_status=IN_STOCK, to_location_id=目标货位)`

### 7.3 批量上架

批量上架本质上是多次单个上架。

建议规则：

- 批量请求中每个追溯码独立校验。
- 如果要求强一致，任一失败则整批回滚。
- 如果要求移动端体验更好，可允许部分成功，并返回每条追溯码处理结果。
- 无论哪种策略，都必须在接口文档和实现中保持一致。

当前建议：批量上架允许部分成功，返回成功数、失败数和失败原因；每个追溯码单独事务或小批量事务。

### 7.4 货位调拨

调拨规则：

- 只有 `IN_STOCK` 或 `MISPLACED` 状态的追溯码允许调拨。
- 目标货位必须存在且启用。
- 调拨后：
  - 更新 `location_id`
  - 如果原状态为 `MISPLACED` 且是错架处理流程，则状态恢复为 `IN_STOCK`
  - 写入 `drug_trace_log(action_type=RELOCATION)`
  - 写入 `inventory_adjustment(adjust_type=RELOCATE)`，用于库存调整审计

普通调拨和盘库错架处理都应复用同一套内部服务。

### 7.5 货位混放检查

当前版本规则：

- 同一货位可以放多种药。
- 混放检查只做提示。
- 返回信息可包括当前货位已有药品种类、批号、近效期情况。
- 不因混放提示阻断上架或调拨。

---

## 8. 库存与追溯业务

### 8.1 库存当前状态

`drug_trace_inventory` 表只表示每个追溯码当前状态。

不得将销售预占直接写入库存状态。预占只看 `trace_reservation`。

常用库存口径：

1. 在库库存

```sql
status = 'IN_STOCK'
```

2. 可售库存

```sql
status = 'IN_STOCK'
AND NOT EXISTS (
  SELECT 1
  FROM trace_reservation r
  WHERE r.trace_code = drug_trace_inventory.trace_code
    AND r.status = 'RESERVED'
    AND r.deleted_at IS NULL
)
```

3. 待上架库存

```sql
status = 'PENDING'
```

4. 异常库存

```sql
status IN ('MISPLACED', 'LOSS_CANDIDATE', 'LOST')
```

### 8.2 推荐销售追溯码

推荐销售追溯码用于销售开单时自动选择库存。

推荐规则：

- 药品必须启用。
- 追溯码必须属于目标 `drug_id`。
- 追溯码状态必须为 `IN_STOCK`。
- 不存在有效 `RESERVED` 预占。
- 按近效期优先，即 `expire_date ASC`。
- 同一有效期下可按入库时间、追溯码 ID 排序。
- 不推荐 `MISPLACED`、`LOSS_CANDIDATE`、`LOST`、`SOLD`、`PENDING` 状态追溯码。

### 8.3 近效期库存

近效期阈值为 30 天。

近效期规则：

```text
expire_date <= 当前日期 + 30 天
AND status = IN_STOCK
```

近效期统计可以按药品、批号、货位聚合。

### 8.4 追溯码验证

追溯码验证接口用于销售、入库、盘库、上架前的快速校验。

应返回：

- 追溯码是否存在。
- 当前状态。
- 药品信息。
- 批号、有效期。
- 当前货位。
- 是否可售。
- 是否已被预占。
- 如果不可操作，返回原因。

### 8.5 追溯日志

所有改变追溯码状态或货位的动作必须写入 `drug_trace_log`。

动作类型：

| action_type | 触发业务 |
|---|---|
| `INBOUND` | 入库扫码确认 |
| `SHELVING` | 上架 |
| `SALE` | 销售结算出库 |
| `RETURN` | 销售退货 |
| `INVENTORY` | 盘库扫描或盘库状态变更 |
| `RELOCATION` | 货位调拨、错架处理 |
| `LOSS` | 确认盘亏 |

日志应尽量记录：

- `trace_code`
- `drug_id`
- `from_status`
- `to_status`
- `from_location_id`
- `to_location_id`
- `operator_id`
- `related_no`
- `order_id`
- `order_item_id`
- `request_id`
- `remark`

---

## 9. 销售与预占业务

### 9.1 创建销售单

销售单对应 `sales_order`，明细对应 `sales_order_item`。

创建规则：

- 创建销售单时必须至少包含一条明细。
- 每条明细表示一盒药，`quantity` 固定为 `1`。
- 如果购买多盒同一药品，应提交多条明细。
- 明细可以指定 `trace_code`，也可以不指定。
- 如果未指定 `trace_code`，后端按近效期优先自动选择一个可售追溯码。
- 销售单价以后端 `drug_info.retail_price` 为准。
- 创建销售单时应同时创建销售明细和追溯码预占。
- `cashier_id` 由当前用户写入。

创建时状态判断：

```text
need_audit = 请求 is_prescription = true
          OR 任意销售明细对应 drug_info.is_prescription = true
```

如果 `need_audit = true`：

- 销售单状态为 `PENDING_REVIEW`。
- 创建药师审核记录 `audit_review(PENDING)`。
- 进入审核后不允许修改明细。

如果 `need_audit = false`：

- 销售单状态为 `PENDING`。
- 可以直接结算。

### 9.2 添加销售明细

仅 `PENDING` 状态销售单允许添加明细。

添加规则：

- `PENDING_REVIEW` 后不允许修改明细。
- `APPROVED` 后不允许修改明细。
- `COMPLETED` 后不允许修改明细。
- 药品必须启用。
- 数量固定为 `1`。
- 后端自动确定最终销售价格。
- 后端自动选择或校验追溯码。
- 成功添加明细时，必须在同一事务中：
  - 创建 `sales_order_item`
  - 创建 `trace_reservation(RESERVED)`
  - 重算销售单金额

如果新增明细导致销售单需要药师审核，则应将销售单切换为 `PENDING_REVIEW` 并创建审核记录；但当前建议处方药销售在创建时一次性提交完整明细，避免中途切换复杂度。

### 9.3 删除销售明细

仅 `PENDING` 状态销售单允许删除明细。

删除规则：

- 删除明细时释放对应 `RESERVED` 预占。
- 重算销售单金额。
- 如果删除后销售单没有任何有效明细，应要求前端取消销售单，或后端自动取消销售单。推荐返回错误，要求至少保留一条明细或调用取消接口。

### 9.4 追溯码预占

预占规则：

- 预占只记录在 `trace_reservation`。
- 被预占追溯码的 `drug_trace_inventory.status` 仍为 `IN_STOCK`。
- 所有可售库存查询必须排除 `RESERVED` 预占。
- 同一追溯码同一时间只能有一条有效 `RESERVED`。
- 预占应设置 `expire_at`。
- 预占过期释放使用 RabbitMQ 死信队列实现。

### 9.5 RabbitMQ 死信队列释放预占

预占创建成功后，后端向 RabbitMQ 延迟/死信队列投递释放消息。

推荐消息内容：

```json
{
  "reservation_id": 123,
  "reservation_no": "RSV-20240101-0001",
  "sales_order_id": 456,
  "trace_code": "86901234567890123456",
  "expire_at": "2026-06-12T10:00:00+08:00"
}
```

释放消费者处理规则：

1. 根据 `reservation_id` 查询预占记录。
2. 如果记录不存在、已删除，直接确认消息。
3. 如果状态不是 `RESERVED`，说明已被结算、主动释放或取消，直接确认消息。
4. 如果状态为 `RESERVED` 且当前时间已到 `expire_at`：
   - 更新 `trace_reservation.status = EXPIRED`
   - 写入 `released_at`
   - 写入操作日志
5. 不修改 `drug_trace_inventory.status`，因为预占期间库存状态一直是 `IN_STOCK`。

并发规则：

- 支付结算和死信释放可能同时发生，必须对 `trace_reservation` 加行锁。
- 结算时只能消费状态仍为 `RESERVED` 的预占。
- 如果预占已过期，结算应失败，要求重新锁定追溯码。

### 9.6 销售结算

允许结算状态：

- `PENDING`：非处方销售，可直接结算。
- `APPROVED`：处方销售，药师审核通过后结算。

不允许结算状态：

- `PENDING_REVIEW`
- `COMPLETED`
- `PARTIALLY_REFUNDED`
- `REFUNDED`
- `CANCELLED`

结算规则：

- 所有销售明细必须存在有效 `RESERVED` 预占。
- 所有预占必须未过期。
- 每个追溯码当前库存状态必须仍为 `IN_STOCK`。
- 实付金额应与后端计算金额一致，或在允许误差范围内。
- 结算成功后：
  - 销售单状态变为 `COMPLETED`
  - 写入 `payment_method`、`actual_amount`、`paid_at`
  - 每条 `trace_reservation.status = CONSUMED`
  - 写入 `confirmed_at`
  - 每个追溯码 `drug_trace_inventory.status = SOLD`
  - 写入 `sold_at`
  - 写入 `drug_trace_log(action_type=SALE)`

### 9.7 销售取消

允许取消状态：

- `PENDING`
- `PENDING_REVIEW`
- `APPROVED`

取消规则：

- 取消时释放所有 `RESERVED` 预占。
- 销售单状态变为 `CANCELLED`。
- 写入 `cancelled_at`。
- 药师审核记录如仍为 `PENDING`，应变为 `CANCELLED`。
- 已 `COMPLETED` 的销售单不能取消，只能退货。

---

## 10. 药师审核业务

### 10.1 审核触发条件

后端根据以下规则自动计算 `need_audit`：

```text
need_audit = sales_order.is_prescription = true
          OR 任意销售明细药品 drug_info.is_prescription = true
```

前端不得直接传入或控制 `need_audit`。

### 10.2 提交审核

销售单进入 `PENDING_REVIEW` 时，应创建 `audit_review` 审核记录。

规则：

- 审核记录状态为 `PENDING`。
- `submitter_id` 为当前提交人。
- `submitted_at` 为提交时间。
- 待审核状态下 `pharmacist_id`、`reviewed_at`、`review_opinion` 允许为空。

如果创建销售单时已经判定需要审核，可以自动创建审核记录，无需前端再调用提交审核接口。若保留提交审核接口，则应保证重复提交不会创建多条有效审核记录。

### 10.3 审核通过

审核通过规则：

- 只有 `PENDING` 状态审核记录允许通过。
- 对应销售单必须为 `PENDING_REVIEW`。
- 审核通过后：
  - `audit_review.status = APPROVED`
  - 写入 `pharmacist_id = 当前用户`
  - 写入 `reviewed_at`
  - 写入 `review_opinion`
  - 销售单状态变为 `APPROVED`

### 10.4 审核驳回

审核驳回规则：

- 只有 `PENDING` 状态审核记录允许驳回。
- 对应销售单必须为 `PENDING_REVIEW`。
- 驳回后：
  - `audit_review.status = REJECTED`
  - 写入 `pharmacist_id = 当前用户`
  - 写入 `reviewed_at`
  - 写入 `review_opinion`
  - 销售单状态直接变为 `CANCELLED`
  - 释放所有 `RESERVED` 预占

驳回后的销售单不得再次提交审核。如需重新销售，应重新创建销售单。

---

## 11. 退货退款业务

### 11.1 退货范围

只有以下状态的销售单允许退货：

- `COMPLETED`
- `PARTIALLY_REFUNDED`

不允许退货状态：

- `PENDING`
- `PENDING_REVIEW`
- `APPROVED`
- `REFUNDED`
- `CANCELLED`

### 11.2 全单退货

`refund_mode = FULL`。

规则：

- 退还所有 `refund_status = NONE` 的销售明细。
- 已退明细不得重复退货。
- 所有被退明细的追溯码必须当前为 `SOLD`。
- 退货成功后，每条明细：
  - `refund_status = REFUNDED`
  - 写入 `refund_amount`
  - 写入 `refunded_at`
  - 写入 `refund_reason`
  - 写入 `refund_operator_id`
- 每个追溯码：
  - `status = IN_STOCK`
  - `location_id` 保持销售前原货位不变
  - `sold_at` 可保留历史销售时间，也可清空；建议保留，由追溯日志表达最新状态
  - 写入 `drug_trace_log(action_type=RETURN)`
- 销售单：
  - `status = REFUNDED`
  - 累加 `refund_amount`
  - 写入 `refunded_at`
  - 写入 `refund_reason`

### 11.3 部分退货

`refund_mode = PARTIAL`。

规则：

- `detail_ids` 必填且不能为空。
- `detail_ids` 必须全部属于当前销售单。
- 指定明细必须尚未退货。
- 指定明细对应追溯码必须当前为 `SOLD`。
- 每条明细独立退货。
- 如果退货后仍有未退明细，销售单状态变为 `PARTIALLY_REFUNDED`。
- 如果退货后所有明细均已退货，销售单状态变为 `REFUNDED`。

### 11.4 退货后的追溯码货位

退货后追溯码恢复为：

```text
status = IN_STOCK
location_id = 销售前原货位
```

当前退货接口不要求前端提交目标货位，不走重新上架流程。

如果未来希望退货后进入待验收/待上架，应新增状态或字段；当前版本不实现。

---

## 12. 盘库与库存调整业务

### 12.1 创建盘库任务

盘库任务对应 `inventory_task`。

创建规则：

- 创建后状态为 `PENDING`。
- `scope_type` 可为 `AREA`、`SHELF`、`LOCATION`。
- `scope_value` 表示盘点范围值。
- `creator_id` 由当前用户写入。
- 可选 `assignee_id` 指定执行人。

### 12.2 开始盘库任务

规则：

- 只有 `PENDING` 状态允许开始。
- 开始后状态变为 `IN_PROGRESS`。
- 写入 `start_time`。

### 12.3 盘库扫码

盘库扫码请求包含：

- `trace_code`
- `scanned_location_code`

规则：

- 只有 `IN_PROGRESS` 状态的盘库任务允许扫码。
- 扫描货位必须存在且启用。
- 同一任务中同一追溯码重复扫码不得重复计数。
- 后端根据当前系统库存判断差异类型。

差异判断：

1. `NORMAL`

```text
追溯码存在
AND status = IN_STOCK
AND system_location_id = scanned_location_id
AND 追溯码属于盘库任务范围
```

2. `MISPLACED_FOUND`

```text
追溯码存在
AND status = IN_STOCK
AND system_location_id != scanned_location_id
```

处理规则：

- 创建盘库明细，记录 `scanned_location_id` 和 `system_location_id`。
- 将追溯码状态更新为 `MISPLACED`。
- 写入 `drug_trace_log(action_type=INVENTORY, to_status=MISPLACED)`。
- 后续通过“处理错架并调整货位”接口恢复为 `IN_STOCK`。

3. `UNEXPECTED`

```text
当前系统中找不到该追溯码的 IN_STOCK 记录
```

包括：

- 追溯码不存在。
- 追溯码已售出。
- 追溯码待上架。
- 追溯码已盘亏。
- 追溯码处于其他非 `IN_STOCK` 状态。

处理规则：

- 只记录盘库明细。
- 不自动创建库存。
- 不自动修改非在库追溯码状态。
- 后续由人工处理。

### 12.4 完成盘库任务

完成规则：

- 只有 `IN_PROGRESS` 状态允许完成。
- 完成时根据盘库范围计算系统应盘点追溯码集合。
- 对范围内状态为 `IN_STOCK` 且未被本任务扫描到的追溯码，标记为 `LOSS_CANDIDATE`。
- 写入 `drug_trace_log(action_type=INVENTORY, to_status=LOSS_CANDIDATE)`。
- 盘库任务状态变为 `COMPLETED`。
- 写入 `end_time`。

注意：已经在扫码过程中被标记为 `MISPLACED` 的追溯码不应再重复标记为 `LOSS_CANDIDATE`。

### 12.5 取消盘库任务

规则：

- `PENDING` 状态可以直接取消。
- `IN_PROGRESS` 状态取消时，需要谨慎处理已产生的异常状态。
- 当前建议：若任务已产生 `MISPLACED` 或 `LOSS_CANDIDATE` 状态，不允许直接取消，应先处理异常或由管理员走手动调整。
- 若未产生库存状态变更，可取消并写入 `end_time`。

### 12.6 确认盘亏候选

确认规则：

- 追溯码必须属于该盘库任务产生的盘亏候选。
- 当前状态必须为 `LOSS_CANDIDATE`。
- 确认后：
  - `drug_trace_inventory.status = LOST`
  - 写入 `inventory_adjustment(adjust_type=LOSS)`
  - 写入 `drug_trace_log(action_type=LOSS)`

### 12.7 驳回盘亏候选

驳回规则：

- 当前状态必须为 `LOSS_CANDIDATE`。
- 驳回后：
  - `drug_trace_inventory.status = IN_STOCK`
  - `location_id` 保持原系统货位
  - 写入库存调整记录或操作日志
  - 写入 `drug_trace_log(action_type=INVENTORY, from_status=LOSS_CANDIDATE, to_status=IN_STOCK)`

### 12.8 处理错架并调整货位

规则：

- 当前状态必须为 `MISPLACED`。
- 请求指定目标货位。
- 目标货位必须存在且启用。
- 处理后：
  - 更新 `location_id = 目标货位`
  - `status = IN_STOCK`
  - 写入 `inventory_adjustment(adjust_type=RELOCATE)`
  - 写入 `drug_trace_log(action_type=RELOCATION, from_status=MISPLACED, to_status=IN_STOCK)`

### 12.9 手动库存调整

手动调整用于少量异常修正，不应绕过核心业务流程。

规则：

- 必须要求填写 `reason`。
- 必须记录操作人。
- 必须写入 `inventory_adjustment`。
- 涉及追溯码状态或货位变化时，必须写入 `drug_trace_log`。
- 不允许通过手动调整把已售追溯码直接改回在库，销售退货应走退货流程。
- 不允许通过手动调整创建重复追溯码。

---

## 13. 扫码任务业务

### 13.1 扫码任务定位

`scan_task` 是移动端批量扫码的包装层，不是唯一业务入口。

规则：

- 直接业务接口和扫码任务提交接口必须调用同一套内部业务服务。
- 不允许为扫码任务单独写一套绕过业务校验的逻辑。
- `/scan-tasks/{id}/submit` 只负责接收扫码结果，并按任务类型调用对应业务服务。
- `/scan-tasks/{id}/complete` 只关闭扫码任务本身，不得绕过入库、上架、盘库的状态校验。

### 13.2 创建扫码任务

规则：

- 创建后状态为 `PENDING`。
- `task_type` 为 `INBOUND`、`SHELVING`、`INVENTORY`。
- `related_id` 绑定对应业务对象：
  - `INBOUND`：入库单 ID
  - `SHELVING`：入库单 ID 或上架批次关联 ID，当前按入库单 ID 处理
  - `INVENTORY`：盘库任务 ID
- `operator_id` 由当前用户写入。

### 13.3 开始扫码任务

规则：

- 只有 `PENDING` 状态允许开始。
- 开始后状态为 `IN_PROGRESS`。
- 写入 `start_time`。

### 13.4 提交扫码结果

根据任务类型调用内部服务：

| task_type | 内部业务服务 |
|---|---|
| `INBOUND` | 入库追溯码确认服务 |
| `SHELVING` | 上架服务 |
| `INVENTORY` | 盘库扫码服务 |

每条扫码明细保存到 `scan_task_detail`，记录：

- `trace_code`
- `location_code`
- `scan_result`
- `error_msg`
- `scan_time`

扫码结果建议：

| scan_result | 含义 |
|---|---|
| `SUCCESS` | 成功 |
| `DUPLICATE` | 重复扫码 |
| `INVALID` | 追溯码或货位无效 |
| `STATUS_ERROR` | 当前状态不允许操作 |

### 13.5 完成与取消扫码任务

完成规则：

- 只有 `IN_PROGRESS` 状态允许完成。
- 完成只代表扫码任务结束，不代表入库单或盘库任务一定完成。
- 具体业务完成仍需调用对应业务接口，如完成入库、完成盘库。

取消规则：

- `PENDING` 或 `IN_PROGRESS` 可以取消。
- 已提交成功的业务扫码结果不因取消扫码任务自动回滚。
- 如需回滚，应走对应业务取消或调整流程。

---

## 14. 预警与通知业务

### 14.1 预警来源

预警通过 `audit_event` 表承载，对外映射为 `AlertInfo`。

预警类型：

| alert_type | 来源 |
|---|---|
| `NEAR_EXPIRE` | 30 天内近效期在库追溯码 |
| `LOW_STOCK` | 药品可售库存数量小于等于 3 |
| `LOSS_CANDIDATE` | 追溯码状态为 `LOSS_CANDIDATE` |
| `MISPLACED` | 追溯码状态为 `MISPLACED` |

### 14.2 近效期预警

规则：

```text
status = IN_STOCK
AND expire_date <= 当前日期 + 30 天
```

建议优先级：

- 7 天内：`HIGH`
- 8~15 天：`MEDIUM`
- 16~30 天：`LOW`

### 14.3 低库存预警

规则：

```text
可售库存数量 <= 3
```

可售库存必须排除有效预占。

### 14.4 盘亏与错架预警

规则：

- `drug_trace_inventory.status = LOSS_CANDIDATE` 产生盘亏候选预警。
- `drug_trace_inventory.status = MISPLACED` 产生错架预警。
- 预警处理不直接改变库存状态；库存状态必须通过盘亏确认、盘亏驳回、错架处理等业务接口改变。

### 14.5 处理预警

处理规则：

- `ACTIVE` 状态预警可以处理为 `RESOLVED`。
- 写入 `resolved_by`、`closed_at`、`resolution`。
- 已处理预警不可重复处理。

### 14.6 忽略预警

忽略规则：

- `ACTIVE` 状态预警可以忽略为 `IGNORED`。
- 写入 `ignored_by`、`ignored_at`、`resolution` 或忽略原因。
- 忽略只改变预警状态，不改变库存业务状态。

### 14.7 通知消息

通知对应 `notification`。

规则：

- 通知可以由预警、审核任务、报表导出等业务触发。
- `read_at IS NULL` 表示未读。
- 标记已读只更新 `read_at`。
- 全部标记已读只作用于当前用户。

---

## 15. 报表与统计业务

### 15.1 首页看板

首页看板数据来源：

| 字段 | 统计口径 |
|---|---|
| 今日销售额 | 今日 `COMPLETED`、`PARTIALLY_REFUNDED`、`REFUNDED` 销售单按实际销售额统计，需扣除退款 |
| 今日销售单数 | 今日创建或完成的销售单数量，推荐按 `paid_at` 统计已完成销售 |
| 今日入库单数 | 今日完成的入库单数量，推荐按 `completed_at` 统计 |
| 当前在库数量 | `drug_trace_inventory.status = IN_STOCK` 数量 |
| 近效期数量 | 30 天内近效期在库数量 |
| 盘亏候选数量 | `status = LOSS_CANDIDATE` 数量 |
| 待上架数量 | `status = PENDING` 数量 |
| 活跃预警数量 | `audit_event.status = 0` 数量 |

### 15.2 销售报表

销售报表以销售单和销售明细为数据源。

建议口径：

- 只统计未删除销售单。
- 默认统计 `COMPLETED`、`PARTIALLY_REFUNDED`、`REFUNDED`。
- 已取消销售单不计入销售额。
- 退款应从销售额中扣除。
- 支持按日期、药品、收银员、支付方式等维度过滤。

### 15.3 入库报表

入库报表以入库单和入库明细为数据源。

建议口径：

- 默认统计 `COMPLETED` 入库单。
- 草稿、待确认、已取消入库单不计入正式入库金额和数量。
- 支持按日期、供应商、药品、批号等维度过滤。

### 15.4 库存报表

库存报表以 `drug_trace_inventory` 当前状态为主。

建议口径：

- 按药品、批号、货位、状态聚合。
- 当前库存统计 `IN_STOCK`。
- 可售库存排除 `RESERVED` 预占。
- 异常库存单独统计 `MISPLACED`、`LOSS_CANDIDATE`、`LOST`。

### 15.5 追溯日志报表

追溯日志报表以 `drug_trace_log` 为数据源。

支持按以下条件过滤：

- 追溯码
- 药品
- 动作类型
- 操作人
- 时间范围
- 关联单号

### 15.6 报表导出任务

报表导出对应 `report_export_task`。

状态：

| 状态 | 含义 |
|---|---|
| `PENDING` | 待执行 |
| `RUNNING` | 执行中 |
| `SUCCESS` | 成功 |
| `FAILED` | 失败 |

规则：

- 创建导出任务时保存查询参数 `query_params`。
- 异步执行时状态从 `PENDING` 到 `RUNNING`。
- 成功后生成文件，写入 `file_id` 和 `finished_at`。
- 失败后写入 `message` 和 `finished_at`。
- 用户只能下载自己有权限访问的报表文件。

---

## 16. 权限与 RBAC 业务

### 16.1 权限来源

权限运行时来源为数据库：

- `sys_role`
- `sys_user_role`
- `sys_permission`
- `sys_role_permission`
- `sys_permission_api`

`casbin_policy.csv` 已弃用，不再作为权限来源。

`casbin_rule` 是由数据库函数 `refresh_casbin_rule_from_rbac()` 生成的 Casbin 标准规则表，不应作为人工维护的业务源表。

### 16.2 Casbin 模型

模型：

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

约定：

- `sub` 使用 `user:{用户ID}`，如 `user:12`。
- 用户角色关系通过 `g(user:{id}, ROLE_CODE)` 表示。
- 接口权限通过 `p(ROLE_CODE, obj, act)` 表示。
- OpenAPI 中的 `x-permission` 只声明 `code`、`resource`、`action`，不写死角色。

### 16.3 角色管理

规则：

- 角色编码 `code` 唯一。
- 内置角色 `built_in = true` 不允许删除。
- 内置角色可以限制修改编码和关键属性。
- 角色停用后，用户即使关联该角色，也不应获得该角色权限。
- 删除角色前必须确认没有用户关联，或先解除关联。

### 16.4 权限管理

规则：

- 权限表示业务能力，不一定等同于单个 API。
- 一个权限可以关联多个 API。
- 分配角色权限后，应刷新 Casbin 规则。
- 权限变更后应使后端权限缓存失效或重新加载。

### 16.5 用户角色管理

规则：

- 用户可拥有多个角色。
- 分配用户角色后，应刷新 Casbin 规则。
- 禁用用户后，该用户不得登录，也不得通过已有 token 继续访问敏感接口。

---

## 17. 审计与日志业务

### 17.1 登录日志

登录行为写入 `login_log`。

规则：

- 登录成功和失败都应记录。
- 记录用户名、用户 ID、IP、User-Agent、成功状态、失败原因。
- 登录成功后更新用户 `last_login_at` 和 `last_login_ip`。

### 17.2 操作日志

关键业务操作写入 `operation_log`。

建议记录：

- 创建、修改、删除基础资料。
- 入库提交、扫码确认、完成、取消。
- 上架、调拨。
- 销售创建、结算、取消、退货。
- 药师审核通过、驳回。
- 盘库开始、扫码、完成、盘亏确认/驳回、错架处理。
- 权限分配、角色变更、用户状态变更。
- 预警处理、忽略。
- 报表导出。

### 17.3 数据变更日志

关键表变更写入 `data_change_log`。

建议记录：

- `table_name`
- `record_id`
- `change_type`
- `operator_id`
- `before_data`
- `after_data`
- `changed_fields`
- `request_id`

适合记录的表：

- `drug_info`
- `supplier`
- `location_info`
- `sys_user`
- `sys_role`
- `sys_role_permission`
- `sys_user_role`
- `sales_order`
- `inbound_order`
- `drug_trace_inventory`

### 17.4 安全事件

安全事件写入 `security_event`。

建议场景：

- 多次登录失败。
- 无权限访问敏感接口。
- token 异常。
- 尝试操作不存在或越权业务数据。
- 高频请求或疑似攻击行为。

---

## 18. 并发控制重点

以下业务必须重点防并发问题。

### 18.1 追溯码重复入库

风险：两个请求同时确认同一个追溯码。

控制方式：

- `drug_trace_inventory.trace_code` 唯一约束。
- 事务内插入，冲突时返回重复追溯码错误。

### 18.2 追溯码重复销售

风险：两个销售单同时选择同一个可售追溯码。

控制方式：

- 查询可售追溯码时加锁或使用原子插入预占。
- `trace_reservation` 对 `status=RESERVED` 的 `trace_code` 建立唯一约束或通过事务锁保证。
- 创建销售明细和预占必须同事务。

### 18.3 预占过期与结算并发

风险：RabbitMQ 死信释放和收银结算同时发生。

控制方式：

- 消费释放消息时锁定预占行。
- 结算时锁定预占行。
- 状态机只允许：
  - `RESERVED -> CONSUMED`
  - `RESERVED -> RELEASED`
  - `RESERVED -> EXPIRED`
- 非 `RESERVED` 状态不得再次流转。

### 18.4 重复退货

风险：两个请求同时退同一明细。

控制方式：

- 退货时锁定销售明细行。
- 只允许 `refund_status = NONE` 的明细退货。
- 更新时带条件：`WHERE refund_status = 'NONE'`。

### 18.5 盘库状态并发

风险：盘库时库存同时被销售、调拨、退货。

当前建议：

- 盘库任务开始后，不全局锁库。
- 对扫码到的追溯码加行锁处理。
- 完成盘库计算盘亏候选时，以完成瞬间的库存状态为准。
- 对已不再是 `IN_STOCK` 的追溯码，不标记为盘亏候选。

---

## 19. API 与内部服务实现建议

### 19.1 不要让 Controller 写业务

建议分层：

```text
Controller -> Application Service -> Domain Service / Repository -> Database
```

Controller 只做：

- 参数绑定
- JWT 用户读取
- 权限注解/中间件
- 调用应用服务
- 返回统一响应

业务状态流转和事务必须放在 Service 层。

### 19.2 内部服务复用

必须复用的服务：

| 内部服务 | 被哪些接口复用 |
|---|---|
| 入库追溯码确认服务 | 直接入库扫码、扫码任务 INBOUND |
| 上架服务 | 直接上架、批量上架、扫码任务 SHELVING |
| 盘库扫码服务 | 盘库扫码、扫码任务 INVENTORY |
| 追溯码预占服务 | 创建销售单、添加销售明细、手动锁定追溯码 |
| 预占释放服务 | 删除明细、取消销售单、审核驳回、RabbitMQ 过期释放 |
| 销售结算服务 | 销售支付接口 |
| 退货服务 | 全单退货、部分退货 |
| 库存调整服务 | 手动调整、盘亏确认、错架处理、调拨 |

不得为不同入口复制粘贴业务逻辑。

### 19.3 字段维护原则

- API DTO 暴露的字段应优先维护。
- API 未暴露但可由已有数据推算的字段，可以查询时动态计算。
- API 未暴露、无法推算、也没有业务接口使用的字段，当前版本不维护。
- 任何涉及操作人、创建人、审核人、退款人、处理人的字段都由后端写入。

---

## 20. 当前版本明确不实现或弱实现的内容

当前版本不实现或仅弱实现：

- 多门店、多租户。
- 货位容量强校验。
- 复杂储存条件混放禁止规则。
- 医保网关内部交易协议。
- AI 子模块内部识别逻辑。
- 退货后重新验收入库流程。
- 复杂促销、会员价、阶梯价。
- 复杂库存冻结、锁库审批。
- 完整财务对账。

这些内容后续如需实现，应先更新数据库、OpenAPI 和本文档。

---

## 21. 关键业务规则速查

1. 入库必须全部扫码确认后才能完成，不允许少扫完成。
2. 入库单 `PENDING_CONFIRM` 后只能扫码确认、完成、取消，不能改明细。
3. 入库单下任意追溯码已经上架为 `IN_STOCK` 后，不允许取消入库单。
4. 追溯码必须在所属入库单 `COMPLETED` 后才能上架。
5. 销售明细一条就是一盒药，数量固定为 `1`。
6. 销售追溯码由后端按近效期优先自动选择，也允许前端指定后由后端校验。
7. 销售预占不改变库存状态，被预占追溯码仍为 `IN_STOCK`。
8. 所有可售库存查询必须排除有效 `RESERVED` 预占。
9. 预占过期释放使用 RabbitMQ 死信队列。
10. 销售价格以后端药品零售价为准。
11. `need_audit` 由后端根据销售单处方标记和药品属性计算。
12. `PENDING_REVIEW` 后不允许修改销售明细。
13. 审核驳回后销售单直接变 `CANCELLED`，并释放所有预占。
14. 已完成销售单不能取消，只能退货。
15. 退货后追溯码恢复为 `IN_STOCK`，货位保持销售前原货位。
16. 盘库找不到 `IN_STOCK` 追溯码时记为 `UNEXPECTED`。
17. 盘库发现错架时立即将追溯码状态改为 `MISPLACED`。
18. 驳回盘亏候选时追溯码恢复为 `IN_STOCK`。
19. 直接业务接口和扫码任务提交接口必须调用同一套内部服务。
20. 同一货位允许放多种药，混放只提示不阻止。
21. 近效期按 30 天计算。
22. 低库存阈值为可售库存小于等于 3。
23. 同一张 AI 发票识别记录不允许多次转入库单。
24. `supplier_code` 必填。
25. 货位容量当前版本不校验。

