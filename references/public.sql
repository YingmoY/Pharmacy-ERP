/*
 Navicat Premium Dump SQL

 Source Server         : PostgreSQL
 Source Server Type    : PostgreSQL
 Source Server Version : 180002 (180002)
 Source Host           : 127.0.0.1:5432
 Source Catalog        : pharmacy_erp
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 180002 (180002)
 File Encoding         : 65001

 Date: 10/06/2026 21:53:50
*/



-- ============================================================
-- Fixed by ChatGPT: make this Navicat dump safe to re-run.
-- 1) Drop dependent views/tables/functions/sequences in a dependency-safe way.
-- 2) Keep the original schema/data definitions below, with only syntax/order fixes.
-- ============================================================
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
CREATE SCHEMA IF NOT EXISTS "public";
SET search_path = "public", "pg_catalog";

BEGIN;

-- Drop views first, because original Navicat order drops tables before views.
DROP VIEW IF EXISTS "public"."v_user_with_roles", "public"."v_sales_order", "public"."v_pharmacist_review", "public"."v_operation_log", "public"."v_inbound_order" CASCADE;

-- Drop tables with CASCADE so repeated initialization is not blocked by FKs/triggers/views.
DROP TABLE IF EXISTS "public"."trace_reservation", "public"."sys_user_role", "public"."sys_user", "public"."sys_role_permission", "public"."sys_role", "public"."sys_permission_api", "public"."sys_permission", "public"."supplier", "public"."security_event", "public"."scan_task_detail", "public"."scan_task", "public"."sales_order_item", "public"."sales_order", "public"."report_export_task", "public"."operation_log", "public"."notification", "public"."login_log", "public"."location_info", "public"."inventory_task_detail", "public"."inventory_task", "public"."inventory_adjustment", "public"."inbound_order_detail", "public"."inbound_order", "public"."file_info", "public"."drug_trace_log", "public"."drug_trace_inventory", "public"."drug_info", "public"."data_change_log", "public"."casbin_rule", "public"."audit_review", "public"."audit_event", "public"."ai_invoice_record" CASCADE;

-- Drop functions after tables/triggers have gone away.
DROP FUNCTION IF EXISTS "public"."refresh_casbin_rule_from_rbac"() CASCADE;
DROP FUNCTION IF EXISTS "public"."update_modified_column"() CASCADE;

-- Drop remaining standalone sequences. Owned sequences may already be removed by DROP TABLE CASCADE.
DROP SEQUENCE IF EXISTS "public"."ai_invoice_record_id_seq", "public"."audit_event_id_seq", "public"."audit_review_id_seq", "public"."casbin_rule_id_seq", "public"."data_change_log_id_seq", "public"."drug_info_id_seq", "public"."drug_trace_inventory_id_seq", "public"."drug_trace_log_id_seq", "public"."file_info_id_seq", "public"."inbound_order_detail_id_seq", "public"."inbound_order_id_seq", "public"."inventory_adjustment_id_seq", "public"."inventory_task_detail_id_seq", "public"."inventory_task_id_seq", "public"."location_info_id_seq", "public"."login_log_id_seq", "public"."notification_id_seq", "public"."operation_log_id_seq", "public"."report_export_task_id_seq", "public"."sales_order_id_seq", "public"."sales_order_item_id_seq", "public"."scan_task_detail_id_seq", "public"."scan_task_id_seq", "public"."security_event_id_seq", "public"."supplier_id_seq", "public"."sys_permission_api_id_seq", "public"."sys_permission_id_seq", "public"."sys_role_id_seq", "public"."sys_user_id_seq", "public"."trace_reservation_id_seq" CASCADE;

-- ----------------------------
-- Sequence structure for ai_invoice_record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_invoice_record_id_seq";
CREATE SEQUENCE "public"."ai_invoice_record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for audit_event_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."audit_event_id_seq";
CREATE SEQUENCE "public"."audit_event_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for audit_review_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."audit_review_id_seq";
CREATE SEQUENCE "public"."audit_review_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for casbin_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."casbin_rule_id_seq";
CREATE SEQUENCE "public"."casbin_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for data_change_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."data_change_log_id_seq";
CREATE SEQUENCE "public"."data_change_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for drug_info_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."drug_info_id_seq";
CREATE SEQUENCE "public"."drug_info_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for drug_trace_inventory_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."drug_trace_inventory_id_seq";
CREATE SEQUENCE "public"."drug_trace_inventory_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for drug_trace_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."drug_trace_log_id_seq";
CREATE SEQUENCE "public"."drug_trace_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for file_info_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."file_info_id_seq";
CREATE SEQUENCE "public"."file_info_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for inbound_order_detail_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."inbound_order_detail_id_seq";
CREATE SEQUENCE "public"."inbound_order_detail_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for inbound_order_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."inbound_order_id_seq";
CREATE SEQUENCE "public"."inbound_order_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for inventory_adjustment_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."inventory_adjustment_id_seq";
CREATE SEQUENCE "public"."inventory_adjustment_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for inventory_task_detail_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."inventory_task_detail_id_seq";
CREATE SEQUENCE "public"."inventory_task_detail_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for inventory_task_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."inventory_task_id_seq";
CREATE SEQUENCE "public"."inventory_task_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for location_info_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."location_info_id_seq";
CREATE SEQUENCE "public"."location_info_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for login_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."login_log_id_seq";
CREATE SEQUENCE "public"."login_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for notification_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notification_id_seq";
CREATE SEQUENCE "public"."notification_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for operation_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."operation_log_id_seq";
CREATE SEQUENCE "public"."operation_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for report_export_task_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."report_export_task_id_seq";
CREATE SEQUENCE "public"."report_export_task_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sales_order_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sales_order_id_seq";
CREATE SEQUENCE "public"."sales_order_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sales_order_item_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sales_order_item_id_seq";
CREATE SEQUENCE "public"."sales_order_item_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for scan_task_detail_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."scan_task_detail_id_seq";
CREATE SEQUENCE "public"."scan_task_detail_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for scan_task_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."scan_task_id_seq";
CREATE SEQUENCE "public"."scan_task_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for security_event_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."security_event_id_seq";
CREATE SEQUENCE "public"."security_event_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for supplier_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."supplier_id_seq";
CREATE SEQUENCE "public"."supplier_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_permission_api_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_permission_api_id_seq";
CREATE SEQUENCE "public"."sys_permission_api_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_permission_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_permission_id_seq";
CREATE SEQUENCE "public"."sys_permission_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_role_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_role_id_seq";
CREATE SEQUENCE "public"."sys_role_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_user_id_seq";
CREATE SEQUENCE "public"."sys_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for trace_reservation_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."trace_reservation_id_seq";
CREATE SEQUENCE "public"."trace_reservation_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for ai_invoice_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_invoice_record";
CREATE TABLE "public"."ai_invoice_record" (
  "id" int8 NOT NULL DEFAULT nextval('ai_invoice_record_id_seq'::regclass),
  "file_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "file_name" varchar(255) COLLATE "pg_catalog"."default",
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'PENDING'::character varying,
  "recognized_supplier_name" varchar(255) COLLATE "pg_catalog"."default",
  "matched_supplier_id" int8,
  "invoice_no" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_date" date,
  "result_json" jsonb,
  "raw_response_json" jsonb,
  "error_message" text COLLATE "pg_catalog"."default",
  "inbound_order_id" int8,
  "converted_at" timestamptz(6),
  "creator_id" int8 NOT NULL,
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON COLUMN "public"."ai_invoice_record"."result_json" IS '规范化后的识别结果 JSON，对应 OpenAPI InvoiceRecognizeResult';
COMMENT ON COLUMN "public"."ai_invoice_record"."raw_response_json" IS 'AI 子模块原始响应，便于排查识别问题';
COMMENT ON TABLE "public"."ai_invoice_record" IS 'AI 发票识别记录表；主 ERP 调用独立 AI 子模块后保存识别状态、结果 JSON、错误信息和转入库单关系';

-- ----------------------------
-- Records of ai_invoice_record
-- ----------------------------

-- ----------------------------
-- Table structure for audit_event
-- ----------------------------
DROP TABLE IF EXISTS "public"."audit_event";
CREATE TABLE "public"."audit_event" (
  "id" int8 NOT NULL DEFAULT nextval('audit_event_id_seq'::regclass),
  "event_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "related_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "related_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "assigned_to" int8,
  "status" int2 DEFAULT 0,
  "resolution" text COLLATE "pg_catalog"."default",
  "closed_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "severity" varchar(20) COLLATE "pg_catalog"."default",
  "ignored_at" timestamptz(6),
  "ignored_by" int8,
  "resolved_by" int8
)
;
COMMENT ON TABLE "public"."audit_event" IS '审计事件表';

-- ----------------------------
-- Records of audit_event
-- ----------------------------

-- ----------------------------
-- Table structure for audit_review
-- ----------------------------
DROP TABLE IF EXISTS "public"."audit_review";
CREATE TABLE "public"."audit_review" (
  "id" int8 NOT NULL DEFAULT nextval('audit_review_id_seq'::regclass),
  "order_id" int8 NOT NULL,
  "pharmacist_id" int8,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "comment" text COLLATE "pg_catalog"."default",
  "reviewed_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "review_no" varchar(50) COLLATE "pg_catalog"."default",
  "submitter_id" int8,
  "submitted_at" timestamptz(6),
  "review_opinion" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."audit_review"."pharmacist_id" IS '实际审核药师 ID；待审核状态允许为空';
COMMENT ON COLUMN "public"."audit_review"."reviewed_at" IS '实际审核时间；待审核状态允许为空，不应默认等于提交时间';
COMMENT ON COLUMN "public"."audit_review"."submitted_at" IS '提交审核时间';
COMMENT ON COLUMN "public"."audit_review"."review_opinion" IS '药师审核意见；待审核状态允许为空';
COMMENT ON TABLE "public"."audit_review" IS '药师审核记录表';

-- ----------------------------
-- Records of audit_review
-- ----------------------------

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_rule";
CREATE TABLE "public"."casbin_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default",
  "v0" varchar(100) COLLATE "pg_catalog"."default",
  "v1" varchar(255) COLLATE "pg_catalog"."default",
  "v2" varchar(255) COLLATE "pg_catalog"."default",
  "v3" varchar(255) COLLATE "pg_catalog"."default",
  "v4" varchar(255) COLLATE "pg_catalog"."default",
  "v5" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."casbin_rule" IS 'Casbin 标准策略表；本系统中作为运行时缓存，由 refresh_casbin_rule_from_rbac() 从 RBAC 业务表生成，不再手写维护 casbin_policy.csv';

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO "public"."casbin_rule" VALUES (1, 'p', 'ADMIN', '/api/v1/ai/invoices', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (2, 'p', 'ADMIN', '/api/v1/ai/invoices/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (3, 'p', 'ADMIN', '/api/v1/ai/invoices/:id/convert-to-inbound', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (4, 'p', 'ADMIN', '/api/v1/ai/invoices/recognize', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (5, 'p', 'ADMIN', '/api/v1/alerts', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (6, 'p', 'ADMIN', '/api/v1/alerts/:id/ignore', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (7, 'p', 'ADMIN', '/api/v1/alerts/:id/resolve', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (8, 'p', 'ADMIN', '/api/v1/alerts/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (9, 'p', 'ADMIN', '/api/v1/alerts/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (10, 'p', 'ADMIN', '/api/v1/audit/data-change-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (11, 'p', 'ADMIN', '/api/v1/audit/login-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (12, 'p', 'ADMIN', '/api/v1/audit/operation-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (13, 'p', 'ADMIN', '/api/v1/audit/operation-logs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (14, 'p', 'ADMIN', '/api/v1/audit/security-events', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (15, 'p', 'ADMIN', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (16, 'p', 'ADMIN', '/api/v1/auth/me', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (17, 'p', 'ADMIN', '/api/v1/auth/password', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (18, 'p', 'ADMIN', '/api/v1/dashboard/inbound-stats', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (19, 'p', 'ADMIN', '/api/v1/dashboard/inventory-stats', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (20, 'p', 'ADMIN', '/api/v1/dashboard/overview', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (21, 'p', 'ADMIN', '/api/v1/dashboard/sales-trend', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (22, 'p', 'ADMIN', '/api/v1/dashboard/top-drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (23, 'p', 'ADMIN', '/api/v1/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (24, 'p', 'ADMIN', '/api/v1/drugs', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (25, 'p', 'ADMIN', '/api/v1/drugs/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (26, 'p', 'ADMIN', '/api/v1/drugs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (27, 'p', 'ADMIN', '/api/v1/drugs/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (28, 'p', 'ADMIN', '/api/v1/drugs/:id/inventory-summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (29, 'p', 'ADMIN', '/api/v1/drugs/:id/sale-info', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (30, 'p', 'ADMIN', '/api/v1/drugs/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (31, 'p', 'ADMIN', '/api/v1/drugs/code/:drug_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (32, 'p', 'ADMIN', '/api/v1/drugs/search', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (33, 'p', 'ADMIN', '/api/v1/files/:file_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (34, 'p', 'ADMIN', '/api/v1/files/:file_id/download', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (35, 'p', 'ADMIN', '/api/v1/files/upload', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (36, 'p', 'ADMIN', '/api/v1/inbound-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (37, 'p', 'ADMIN', '/api/v1/inbound-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (38, 'p', 'ADMIN', '/api/v1/inbound-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (39, 'p', 'ADMIN', '/api/v1/inbound-orders/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (40, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (41, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (42, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/confirm-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (43, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/confirm-traces', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (44, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (45, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (46, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (47, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/details/:detail_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (48, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/details/:detail_id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (49, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/progress', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (50, 'p', 'ADMIN', '/api/v1/inbound-orders/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (51, 'p', 'ADMIN', '/api/v1/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (52, 'p', 'ADMIN', '/api/v1/inventory-adjustments', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (53, 'p', 'ADMIN', '/api/v1/inventory-adjustments', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (54, 'p', 'ADMIN', '/api/v1/inventory-adjustments/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (55, 'p', 'ADMIN', '/api/v1/inventory-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (56, 'p', 'ADMIN', '/api/v1/inventory-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (57, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (58, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/assign', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (59, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (60, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (61, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (62, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (63, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/confirm', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (64, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/reject', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (65, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/misplaced', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (66, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/misplaced/:trace_code/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (67, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (68, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (69, 'p', 'ADMIN', '/api/v1/inventory-tasks/:id/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (70, 'p', 'ADMIN', '/api/v1/inventory/:trace_code/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (71, 'p', 'ADMIN', '/api/v1/inventory/drugs/:drug_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (72, 'p', 'ADMIN', '/api/v1/inventory/locations/:location_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (73, 'p', 'ADMIN', '/api/v1/inventory/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (74, 'p', 'ADMIN', '/api/v1/inventory/pending-shelving', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (75, 'p', 'ADMIN', '/api/v1/inventory/recommend-sale', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (76, 'p', 'ADMIN', '/api/v1/inventory/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (77, 'p', 'ADMIN', '/api/v1/locations', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (78, 'p', 'ADMIN', '/api/v1/locations', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (79, 'p', 'ADMIN', '/api/v1/locations/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (80, 'p', 'ADMIN', '/api/v1/locations/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (81, 'p', 'ADMIN', '/api/v1/locations/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (82, 'p', 'ADMIN', '/api/v1/locations/:id/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (83, 'p', 'ADMIN', '/api/v1/locations/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (84, 'p', 'ADMIN', '/api/v1/locations/areas', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (85, 'p', 'ADMIN', '/api/v1/locations/code/:location_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (86, 'p', 'ADMIN', '/api/v1/notifications', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (87, 'p', 'ADMIN', '/api/v1/notifications/:id/read', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (88, 'p', 'ADMIN', '/api/v1/notifications/read-all', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (89, 'p', 'ADMIN', '/api/v1/notifications/unread-count', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (90, 'p', 'ADMIN', '/api/v1/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (91, 'p', 'ADMIN', '/api/v1/pharmacist/reviews', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (92, 'p', 'ADMIN', '/api/v1/pharmacist/reviews/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (93, 'p', 'ADMIN', '/api/v1/pharmacist/reviews/:id/approve', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (94, 'p', 'ADMIN', '/api/v1/pharmacist/reviews/:id/reject', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (95, 'p', 'ADMIN', '/api/v1/reports/export-tasks/:task_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (96, 'p', 'ADMIN', '/api/v1/reports/inbound', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (97, 'p', 'ADMIN', '/api/v1/reports/inbound/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (98, 'p', 'ADMIN', '/api/v1/reports/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (99, 'p', 'ADMIN', '/api/v1/reports/inventory/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (100, 'p', 'ADMIN', '/api/v1/reports/sales', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (101, 'p', 'ADMIN', '/api/v1/reports/sales/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (102, 'p', 'ADMIN', '/api/v1/reports/trace-log', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (103, 'p', 'ADMIN', '/api/v1/reports/trace-log/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (104, 'p', 'ADMIN', '/api/v1/roles', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (105, 'p', 'ADMIN', '/api/v1/roles', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (106, 'p', 'ADMIN', '/api/v1/roles/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (107, 'p', 'ADMIN', '/api/v1/roles/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (108, 'p', 'ADMIN', '/api/v1/roles/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (109, 'p', 'ADMIN', '/api/v1/roles/:id/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (110, 'p', 'ADMIN', '/api/v1/roles/:id/permissions', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (111, 'p', 'ADMIN', '/api/v1/sales-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (112, 'p', 'ADMIN', '/api/v1/sales-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (113, 'p', 'ADMIN', '/api/v1/sales-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (114, 'p', 'ADMIN', '/api/v1/sales-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (115, 'p', 'ADMIN', '/api/v1/sales-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (116, 'p', 'ADMIN', '/api/v1/sales-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (117, 'p', 'ADMIN', '/api/v1/sales-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (118, 'p', 'ADMIN', '/api/v1/sales-orders/:id/pay', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (119, 'p', 'ADMIN', '/api/v1/sales-orders/:id/refund', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (120, 'p', 'ADMIN', '/api/v1/sales-orders/:id/release-reservation', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (121, 'p', 'ADMIN', '/api/v1/sales-orders/:id/reserve-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (122, 'p', 'ADMIN', '/api/v1/sales-orders/:id/reserved-traces', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (123, 'p', 'ADMIN', '/api/v1/sales-orders/:id/review-record', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (124, 'p', 'ADMIN', '/api/v1/sales-orders/:id/scan-verify', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (125, 'p', 'ADMIN', '/api/v1/sales-orders/:id/submit-review', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (126, 'p', 'ADMIN', '/api/v1/scan-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (127, 'p', 'ADMIN', '/api/v1/scan-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (128, 'p', 'ADMIN', '/api/v1/scan-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (129, 'p', 'ADMIN', '/api/v1/scan-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (130, 'p', 'ADMIN', '/api/v1/scan-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (131, 'p', 'ADMIN', '/api/v1/scan-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (132, 'p', 'ADMIN', '/api/v1/scan-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (133, 'p', 'ADMIN', '/api/v1/scan-tasks/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (134, 'p', 'ADMIN', '/api/v1/shelving/batch', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (135, 'p', 'ADMIN', '/api/v1/shelving/mix-check', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (136, 'p', 'ADMIN', '/api/v1/shelving/pending', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (137, 'p', 'ADMIN', '/api/v1/shelving/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (138, 'p', 'ADMIN', '/api/v1/shelving/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (139, 'p', 'ADMIN', '/api/v1/suppliers', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (140, 'p', 'ADMIN', '/api/v1/suppliers', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (141, 'p', 'ADMIN', '/api/v1/suppliers/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (142, 'p', 'ADMIN', '/api/v1/suppliers/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (143, 'p', 'ADMIN', '/api/v1/suppliers/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (144, 'p', 'ADMIN', '/api/v1/suppliers/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (145, 'p', 'ADMIN', '/api/v1/trace/:trace_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (146, 'p', 'ADMIN', '/api/v1/trace/:trace_code/full-chain', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (147, 'p', 'ADMIN', '/api/v1/trace/:trace_code/logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (148, 'p', 'ADMIN', '/api/v1/trace/validate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (149, 'p', 'ADMIN', '/api/v1/users', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (150, 'p', 'ADMIN', '/api/v1/users', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (151, 'p', 'ADMIN', '/api/v1/users/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (152, 'p', 'ADMIN', '/api/v1/users/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (153, 'p', 'ADMIN', '/api/v1/users/:id/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (154, 'p', 'ADMIN', '/api/v1/users/:id/reset-password', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (155, 'p', 'ADMIN', '/api/v1/users/:id/roles', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (156, 'p', 'ADMIN', '/api/v1/users/:id/roles', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (157, 'p', 'ADMIN', '/api/v1/users/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (158, 'p', 'CASHIER', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (159, 'p', 'CASHIER', '/api/v1/auth/me', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (160, 'p', 'CASHIER', '/api/v1/auth/password', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (161, 'p', 'CASHIER', '/api/v1/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (162, 'p', 'CASHIER', '/api/v1/drugs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (163, 'p', 'CASHIER', '/api/v1/drugs/:id/inventory-summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (164, 'p', 'CASHIER', '/api/v1/drugs/:id/sale-info', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (165, 'p', 'CASHIER', '/api/v1/drugs/code/:drug_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (166, 'p', 'CASHIER', '/api/v1/drugs/search', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (167, 'p', 'CASHIER', '/api/v1/files/:file_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (168, 'p', 'CASHIER', '/api/v1/files/:file_id/download', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (169, 'p', 'CASHIER', '/api/v1/files/upload', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (170, 'p', 'CASHIER', '/api/v1/inbound-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (171, 'p', 'CASHIER', '/api/v1/inventory-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (172, 'p', 'CASHIER', '/api/v1/sales-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (173, 'p', 'CASHIER', '/api/v1/sales-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (174, 'p', 'CASHIER', '/api/v1/sales-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (175, 'p', 'CASHIER', '/api/v1/sales-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (176, 'p', 'CASHIER', '/api/v1/sales-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (177, 'p', 'CASHIER', '/api/v1/sales-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (178, 'p', 'CASHIER', '/api/v1/sales-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (179, 'p', 'CASHIER', '/api/v1/sales-orders/:id/pay', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (180, 'p', 'CASHIER', '/api/v1/sales-orders/:id/release-reservation', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (181, 'p', 'CASHIER', '/api/v1/sales-orders/:id/reserve-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (182, 'p', 'CASHIER', '/api/v1/sales-orders/:id/reserved-traces', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (183, 'p', 'CASHIER', '/api/v1/sales-orders/:id/review-record', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (184, 'p', 'CASHIER', '/api/v1/sales-orders/:id/scan-verify', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (185, 'p', 'CASHIER', '/api/v1/sales-orders/:id/submit-review', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (186, 'p', 'CASHIER', '/api/v1/scan-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (187, 'p', 'PHARMACIST', '/api/v1/alerts', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (188, 'p', 'PHARMACIST', '/api/v1/alerts/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (189, 'p', 'PHARMACIST', '/api/v1/alerts/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (190, 'p', 'PHARMACIST', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (191, 'p', 'PHARMACIST', '/api/v1/auth/me', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (192, 'p', 'PHARMACIST', '/api/v1/auth/password', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (193, 'p', 'PHARMACIST', '/api/v1/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (194, 'p', 'PHARMACIST', '/api/v1/drugs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (195, 'p', 'PHARMACIST', '/api/v1/drugs/:id/inventory-summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (196, 'p', 'PHARMACIST', '/api/v1/drugs/:id/sale-info', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (197, 'p', 'PHARMACIST', '/api/v1/drugs/code/:drug_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (198, 'p', 'PHARMACIST', '/api/v1/drugs/search', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (199, 'p', 'PHARMACIST', '/api/v1/files/:file_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (200, 'p', 'PHARMACIST', '/api/v1/files/:file_id/download', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (201, 'p', 'PHARMACIST', '/api/v1/files/upload', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (202, 'p', 'PHARMACIST', '/api/v1/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (203, 'p', 'PHARMACIST', '/api/v1/inventory-adjustments', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (204, 'p', 'PHARMACIST', '/api/v1/inventory-adjustments/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (205, 'p', 'PHARMACIST', '/api/v1/inventory/drugs/:drug_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (206, 'p', 'PHARMACIST', '/api/v1/inventory/locations/:location_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (207, 'p', 'PHARMACIST', '/api/v1/inventory/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (208, 'p', 'PHARMACIST', '/api/v1/inventory/pending-shelving', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (209, 'p', 'PHARMACIST', '/api/v1/inventory/recommend-sale', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (210, 'p', 'PHARMACIST', '/api/v1/inventory/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (211, 'p', 'PHARMACIST', '/api/v1/locations', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (212, 'p', 'PHARMACIST', '/api/v1/locations/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (213, 'p', 'PHARMACIST', '/api/v1/locations/:id/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (214, 'p', 'PHARMACIST', '/api/v1/locations/areas', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (215, 'p', 'PHARMACIST', '/api/v1/locations/code/:location_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (216, 'p', 'PHARMACIST', '/api/v1/notifications', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (217, 'p', 'PHARMACIST', '/api/v1/notifications/unread-count', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (218, 'p', 'PHARMACIST', '/api/v1/pharmacist/reviews', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (219, 'p', 'PHARMACIST', '/api/v1/pharmacist/reviews/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (220, 'p', 'PHARMACIST', '/api/v1/pharmacist/reviews/:id/approve', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (221, 'p', 'PHARMACIST', '/api/v1/pharmacist/reviews/:id/reject', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (222, 'p', 'PHARMACIST', '/api/v1/sales-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (223, 'p', 'PHARMACIST', '/api/v1/sales-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (224, 'p', 'PHARMACIST', '/api/v1/sales-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (225, 'p', 'PHARMACIST', '/api/v1/sales-orders/:id/reserved-traces', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (226, 'p', 'PHARMACIST', '/api/v1/sales-orders/:id/review-record', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (227, 'p', 'PHARMACIST', '/api/v1/trace/:trace_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (228, 'p', 'PHARMACIST', '/api/v1/trace/:trace_code/full-chain', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (229, 'p', 'PHARMACIST', '/api/v1/trace/:trace_code/logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (230, 'p', 'STORE_MANAGER', '/api/v1/ai/invoices', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (231, 'p', 'STORE_MANAGER', '/api/v1/ai/invoices/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (232, 'p', 'STORE_MANAGER', '/api/v1/ai/invoices/:id/convert-to-inbound', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (233, 'p', 'STORE_MANAGER', '/api/v1/ai/invoices/recognize', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (234, 'p', 'STORE_MANAGER', '/api/v1/alerts', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (235, 'p', 'STORE_MANAGER', '/api/v1/alerts/:id/ignore', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (236, 'p', 'STORE_MANAGER', '/api/v1/alerts/:id/resolve', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (237, 'p', 'STORE_MANAGER', '/api/v1/alerts/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (238, 'p', 'STORE_MANAGER', '/api/v1/alerts/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (239, 'p', 'STORE_MANAGER', '/api/v1/audit/data-change-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (240, 'p', 'STORE_MANAGER', '/api/v1/audit/login-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (241, 'p', 'STORE_MANAGER', '/api/v1/audit/operation-logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (242, 'p', 'STORE_MANAGER', '/api/v1/audit/operation-logs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (243, 'p', 'STORE_MANAGER', '/api/v1/audit/security-events', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (244, 'p', 'STORE_MANAGER', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (245, 'p', 'STORE_MANAGER', '/api/v1/auth/me', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (246, 'p', 'STORE_MANAGER', '/api/v1/auth/password', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (247, 'p', 'STORE_MANAGER', '/api/v1/dashboard/inbound-stats', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (248, 'p', 'STORE_MANAGER', '/api/v1/dashboard/inventory-stats', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (249, 'p', 'STORE_MANAGER', '/api/v1/dashboard/overview', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (250, 'p', 'STORE_MANAGER', '/api/v1/dashboard/sales-trend', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (251, 'p', 'STORE_MANAGER', '/api/v1/dashboard/top-drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (252, 'p', 'STORE_MANAGER', '/api/v1/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (253, 'p', 'STORE_MANAGER', '/api/v1/drugs', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (254, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (255, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (256, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (257, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id/inventory-summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (258, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id/sale-info', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (259, 'p', 'STORE_MANAGER', '/api/v1/drugs/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (260, 'p', 'STORE_MANAGER', '/api/v1/drugs/code/:drug_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (261, 'p', 'STORE_MANAGER', '/api/v1/drugs/search', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (262, 'p', 'STORE_MANAGER', '/api/v1/files/:file_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (263, 'p', 'STORE_MANAGER', '/api/v1/files/:file_id/download', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (264, 'p', 'STORE_MANAGER', '/api/v1/files/upload', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (265, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (266, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (267, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (268, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (269, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (270, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (271, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/confirm-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (272, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/confirm-traces', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (273, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (274, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (275, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (276, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/details/:detail_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (277, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/details/:detail_id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (278, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/progress', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (279, 'p', 'STORE_MANAGER', '/api/v1/inbound-orders/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (280, 'p', 'STORE_MANAGER', '/api/v1/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (281, 'p', 'STORE_MANAGER', '/api/v1/inventory-adjustments', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (282, 'p', 'STORE_MANAGER', '/api/v1/inventory-adjustments', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (283, 'p', 'STORE_MANAGER', '/api/v1/inventory-adjustments/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (284, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (285, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (286, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (287, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/assign', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (288, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (289, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (290, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (291, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (292, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/confirm', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (293, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/reject', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (294, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/misplaced', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (295, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/misplaced/:trace_code/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (296, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (297, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (298, 'p', 'STORE_MANAGER', '/api/v1/inventory-tasks/:id/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (299, 'p', 'STORE_MANAGER', '/api/v1/inventory/:trace_code/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (300, 'p', 'STORE_MANAGER', '/api/v1/inventory/drugs/:drug_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (301, 'p', 'STORE_MANAGER', '/api/v1/inventory/locations/:location_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (302, 'p', 'STORE_MANAGER', '/api/v1/inventory/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (303, 'p', 'STORE_MANAGER', '/api/v1/inventory/pending-shelving', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (304, 'p', 'STORE_MANAGER', '/api/v1/inventory/recommend-sale', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (305, 'p', 'STORE_MANAGER', '/api/v1/inventory/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (306, 'p', 'STORE_MANAGER', '/api/v1/locations', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (307, 'p', 'STORE_MANAGER', '/api/v1/locations', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (308, 'p', 'STORE_MANAGER', '/api/v1/locations/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (309, 'p', 'STORE_MANAGER', '/api/v1/locations/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (310, 'p', 'STORE_MANAGER', '/api/v1/locations/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (311, 'p', 'STORE_MANAGER', '/api/v1/locations/:id/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (312, 'p', 'STORE_MANAGER', '/api/v1/locations/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (313, 'p', 'STORE_MANAGER', '/api/v1/locations/areas', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (314, 'p', 'STORE_MANAGER', '/api/v1/locations/code/:location_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (315, 'p', 'STORE_MANAGER', '/api/v1/notifications', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (316, 'p', 'STORE_MANAGER', '/api/v1/notifications/:id/read', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (317, 'p', 'STORE_MANAGER', '/api/v1/notifications/read-all', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (318, 'p', 'STORE_MANAGER', '/api/v1/notifications/unread-count', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (319, 'p', 'STORE_MANAGER', '/api/v1/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (320, 'p', 'STORE_MANAGER', '/api/v1/pharmacist/reviews', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (321, 'p', 'STORE_MANAGER', '/api/v1/pharmacist/reviews/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (322, 'p', 'STORE_MANAGER', '/api/v1/reports/export-tasks/:task_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (323, 'p', 'STORE_MANAGER', '/api/v1/reports/inbound', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (324, 'p', 'STORE_MANAGER', '/api/v1/reports/inbound/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (325, 'p', 'STORE_MANAGER', '/api/v1/reports/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (326, 'p', 'STORE_MANAGER', '/api/v1/reports/inventory/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (327, 'p', 'STORE_MANAGER', '/api/v1/reports/sales', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (328, 'p', 'STORE_MANAGER', '/api/v1/reports/sales/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (329, 'p', 'STORE_MANAGER', '/api/v1/reports/trace-log', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (330, 'p', 'STORE_MANAGER', '/api/v1/reports/trace-log/export', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (331, 'p', 'STORE_MANAGER', '/api/v1/roles', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (332, 'p', 'STORE_MANAGER', '/api/v1/roles/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (333, 'p', 'STORE_MANAGER', '/api/v1/roles/:id/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (334, 'p', 'STORE_MANAGER', '/api/v1/sales-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (335, 'p', 'STORE_MANAGER', '/api/v1/sales-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (336, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (337, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (338, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (339, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (340, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (341, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/pay', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (342, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/refund', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (343, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/release-reservation', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (344, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/reserve-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (345, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/reserved-traces', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (346, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/review-record', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (347, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/scan-verify', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (348, 'p', 'STORE_MANAGER', '/api/v1/sales-orders/:id/submit-review', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (349, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (350, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (351, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (352, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (353, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (354, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (355, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (356, 'p', 'STORE_MANAGER', '/api/v1/scan-tasks/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (357, 'p', 'STORE_MANAGER', '/api/v1/shelving/batch', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (358, 'p', 'STORE_MANAGER', '/api/v1/shelving/mix-check', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (359, 'p', 'STORE_MANAGER', '/api/v1/shelving/pending', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (360, 'p', 'STORE_MANAGER', '/api/v1/shelving/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (361, 'p', 'STORE_MANAGER', '/api/v1/shelving/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (362, 'p', 'STORE_MANAGER', '/api/v1/suppliers', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (363, 'p', 'STORE_MANAGER', '/api/v1/suppliers', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (364, 'p', 'STORE_MANAGER', '/api/v1/suppliers/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (365, 'p', 'STORE_MANAGER', '/api/v1/suppliers/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (366, 'p', 'STORE_MANAGER', '/api/v1/suppliers/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (367, 'p', 'STORE_MANAGER', '/api/v1/suppliers/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (368, 'p', 'STORE_MANAGER', '/api/v1/trace/:trace_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (369, 'p', 'STORE_MANAGER', '/api/v1/trace/:trace_code/full-chain', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (370, 'p', 'STORE_MANAGER', '/api/v1/trace/:trace_code/logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (371, 'p', 'STORE_MANAGER', '/api/v1/trace/validate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (372, 'p', 'STORE_MANAGER', '/api/v1/users', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (373, 'p', 'STORE_MANAGER', '/api/v1/users/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (374, 'p', 'STORE_MANAGER', '/api/v1/users/:id/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (375, 'p', 'STORE_MANAGER', '/api/v1/users/:id/roles', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (376, 'p', 'WAREHOUSE', '/api/v1/ai/invoices', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (377, 'p', 'WAREHOUSE', '/api/v1/ai/invoices/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (378, 'p', 'WAREHOUSE', '/api/v1/ai/invoices/:id/convert-to-inbound', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (379, 'p', 'WAREHOUSE', '/api/v1/ai/invoices/recognize', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (380, 'p', 'WAREHOUSE', '/api/v1/alerts', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (381, 'p', 'WAREHOUSE', '/api/v1/alerts/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (382, 'p', 'WAREHOUSE', '/api/v1/alerts/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (383, 'p', 'WAREHOUSE', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (384, 'p', 'WAREHOUSE', '/api/v1/auth/me', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (385, 'p', 'WAREHOUSE', '/api/v1/auth/password', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (386, 'p', 'WAREHOUSE', '/api/v1/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (387, 'p', 'WAREHOUSE', '/api/v1/drugs/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (388, 'p', 'WAREHOUSE', '/api/v1/drugs/:id/inventory-summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (389, 'p', 'WAREHOUSE', '/api/v1/drugs/:id/sale-info', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (390, 'p', 'WAREHOUSE', '/api/v1/drugs/code/:drug_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (391, 'p', 'WAREHOUSE', '/api/v1/drugs/search', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (392, 'p', 'WAREHOUSE', '/api/v1/files/:file_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (393, 'p', 'WAREHOUSE', '/api/v1/files/:file_id/download', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (394, 'p', 'WAREHOUSE', '/api/v1/files/upload', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (395, 'p', 'WAREHOUSE', '/api/v1/inbound-orders', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (396, 'p', 'WAREHOUSE', '/api/v1/inbound-orders', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (397, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (398, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (399, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (400, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (401, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/confirm-trace', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (402, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/confirm-traces', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (403, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (404, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/details', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (405, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/details/:detail_id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (406, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/details/:detail_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (407, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/details/:detail_id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (408, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/progress', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (409, 'p', 'WAREHOUSE', '/api/v1/inbound-orders/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (410, 'p', 'WAREHOUSE', '/api/v1/inventory', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (411, 'p', 'WAREHOUSE', '/api/v1/inventory-adjustments', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (412, 'p', 'WAREHOUSE', '/api/v1/inventory-adjustments', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (413, 'p', 'WAREHOUSE', '/api/v1/inventory-adjustments/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (414, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (415, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (416, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (417, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/assign', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (418, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (419, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (420, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (421, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/loss-candidates', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (422, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/confirm', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (423, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/reject', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (424, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/misplaced', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (425, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/misplaced/:trace_code/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (426, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (427, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (428, 'p', 'WAREHOUSE', '/api/v1/inventory-tasks/:id/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (429, 'p', 'WAREHOUSE', '/api/v1/inventory/:trace_code/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (430, 'p', 'WAREHOUSE', '/api/v1/inventory/drugs/:drug_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (431, 'p', 'WAREHOUSE', '/api/v1/inventory/locations/:location_id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (432, 'p', 'WAREHOUSE', '/api/v1/inventory/near-expire', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (433, 'p', 'WAREHOUSE', '/api/v1/inventory/pending-shelving', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (434, 'p', 'WAREHOUSE', '/api/v1/inventory/recommend-sale', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (435, 'p', 'WAREHOUSE', '/api/v1/inventory/summary', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (436, 'p', 'WAREHOUSE', '/api/v1/locations', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (437, 'p', 'WAREHOUSE', '/api/v1/locations', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (438, 'p', 'WAREHOUSE', '/api/v1/locations/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (439, 'p', 'WAREHOUSE', '/api/v1/locations/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (440, 'p', 'WAREHOUSE', '/api/v1/locations/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (441, 'p', 'WAREHOUSE', '/api/v1/locations/:id/drugs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (442, 'p', 'WAREHOUSE', '/api/v1/locations/:id/status', 'PATCH', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (443, 'p', 'WAREHOUSE', '/api/v1/locations/areas', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (444, 'p', 'WAREHOUSE', '/api/v1/locations/code/:location_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (445, 'p', 'WAREHOUSE', '/api/v1/notifications', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (446, 'p', 'WAREHOUSE', '/api/v1/notifications/unread-count', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (447, 'p', 'WAREHOUSE', '/api/v1/sales-orders/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (448, 'p', 'WAREHOUSE', '/api/v1/scan-tasks', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (449, 'p', 'WAREHOUSE', '/api/v1/scan-tasks', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (450, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (451, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id/cancel', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (452, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id/complete', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (453, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id/details', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (454, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id/start', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (455, 'p', 'WAREHOUSE', '/api/v1/scan-tasks/:id/submit', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (456, 'p', 'WAREHOUSE', '/api/v1/shelving/batch', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (457, 'p', 'WAREHOUSE', '/api/v1/shelving/mix-check', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (458, 'p', 'WAREHOUSE', '/api/v1/shelving/pending', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (459, 'p', 'WAREHOUSE', '/api/v1/shelving/relocate', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (460, 'p', 'WAREHOUSE', '/api/v1/shelving/scan', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (461, 'p', 'WAREHOUSE', '/api/v1/suppliers', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (462, 'p', 'WAREHOUSE', '/api/v1/suppliers', 'POST', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (463, 'p', 'WAREHOUSE', '/api/v1/suppliers/:id', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (464, 'p', 'WAREHOUSE', '/api/v1/suppliers/:id', 'PUT', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (465, 'p', 'WAREHOUSE', '/api/v1/trace/:trace_code', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (466, 'p', 'WAREHOUSE', '/api/v1/trace/:trace_code/full-chain', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (467, 'p', 'WAREHOUSE', '/api/v1/trace/:trace_code/logs', 'GET', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (468, 'p', 'WAREHOUSE', '/api/v1/trace/validate', 'POST', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for data_change_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."data_change_log";
CREATE TABLE "public"."data_change_log" (
  "id" int8 NOT NULL DEFAULT nextval('data_change_log_id_seq'::regclass),
  "table_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "record_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "change_type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "operator_id" int8,
  "operator_name" varchar(100) COLLATE "pg_catalog"."default",
  "before_data" jsonb,
  "after_data" jsonb,
  "changed_fields" jsonb,
  "request_id" varchar(100) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."data_change_log" IS '数据变更日志表';

-- ----------------------------
-- Records of data_change_log
-- ----------------------------

-- ----------------------------
-- Table structure for drug_info
-- ----------------------------
DROP TABLE IF EXISTS "public"."drug_info";
CREATE TABLE "public"."drug_info" (
  "id" int8 NOT NULL DEFAULT nextval('drug_info_id_seq'::regclass),
  "drug_code" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "common_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "trade_name" varchar(100) COLLATE "pg_catalog"."default",
  "specification" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "dosage_form" varchar(50) COLLATE "pg_catalog"."default",
  "manufacturer" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "approval_number" varchar(50) COLLATE "pg_catalog"."default",
  "is_prescription" bool DEFAULT false,
  "is_medicare" bool DEFAULT false,
  "status" int2 DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "barcode" varchar(100) COLLATE "pg_catalog"."default",
  "unit" varchar(20) COLLATE "pg_catalog"."default",
  "retail_price" numeric(10,2),
  "purchase_price" numeric(10,2),
  "storage_condition" varchar(100) COLLATE "pg_catalog"."default",
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."drug_info" IS '药品基础信息表';

-- ----------------------------
-- Records of drug_info
-- ----------------------------

-- ----------------------------
-- Table structure for drug_trace_inventory
-- ----------------------------
DROP TABLE IF EXISTS "public"."drug_trace_inventory";
CREATE TABLE "public"."drug_trace_inventory" (
  "id" int8 NOT NULL DEFAULT nextval('drug_trace_inventory_id_seq'::regclass),
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "drug_id" int8 NOT NULL,
  "batch_number" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "expire_date" date NOT NULL,
  "location_id" int8,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "inbound_order_id" int8 NOT NULL,
  "inbound_detail_id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "sold_at" timestamptz(6),
  "last_action" varchar(50) COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."drug_trace_inventory" IS '药品追溯码库存表；仅保存追溯码当前库存状态，销售预占/释放/消耗以 trace_reservation 为准';

-- ----------------------------
-- Records of drug_trace_inventory
-- ----------------------------

-- ----------------------------
-- Table structure for drug_trace_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."drug_trace_log";
CREATE TABLE "public"."drug_trace_log" (
  "id" int8 NOT NULL DEFAULT nextval('drug_trace_log_id_seq'::regclass),
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "action_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "from_status" varchar(20) COLLATE "pg_catalog"."default",
  "to_status" varchar(20) COLLATE "pg_catalog"."default",
  "operator_id" int8 NOT NULL,
  "related_no" varchar(100) COLLATE "pg_catalog"."default",
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "drug_id" int8,
  "order_id" int8,
  "order_item_id" int8,
  "request_id" varchar(100) COLLATE "pg_catalog"."default",
  "from_location_id" int8,
  "to_location_id" int8
)
;
COMMENT ON COLUMN "public"."drug_trace_log"."from_location_id" IS '动作发生前货位ID；用于移位、盘点纠错、退货等轨迹展示，可为空';
COMMENT ON COLUMN "public"."drug_trace_log"."to_location_id" IS '动作发生后货位ID；用于上架、移位、盘点纠错、退货等轨迹展示，可为空';
COMMENT ON TABLE "public"."drug_trace_log" IS '药品业务轨迹表';

-- ----------------------------
-- Records of drug_trace_log
-- ----------------------------

-- ----------------------------
-- Table structure for file_info
-- ----------------------------
DROP TABLE IF EXISTS "public"."file_info";
CREATE TABLE "public"."file_info" (
  "id" int8 NOT NULL DEFAULT nextval('file_info_id_seq'::regclass),
  "file_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "original_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "storage_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "storage_path" varchar(500) COLLATE "pg_catalog"."default" NOT NULL,
  "content_type" varchar(100) COLLATE "pg_catalog"."default",
  "file_size" int8,
  "file_hash" varchar(128) COLLATE "pg_catalog"."default",
  "business_type" varchar(50) COLLATE "pg_catalog"."default",
  "business_id" varchar(100) COLLATE "pg_catalog"."default",
  "uploader_id" int8,
  "status" int2 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."file_info" IS '文件信息表';

-- ----------------------------
-- Records of file_info
-- ----------------------------

-- ----------------------------
-- Table structure for inbound_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."inbound_order";
CREATE TABLE "public"."inbound_order" (
  "id" int8 NOT NULL DEFAULT nextval('inbound_order_id_seq'::regclass),
  "order_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "invoice_no" varchar(100) COLLATE "pg_catalog"."default",
  "operator_id" int8 NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "supplier_id" int8 NOT NULL,
  "creator_id" int8,
  "total_amount" numeric(12,2) DEFAULT 0,
  "submitted_at" timestamptz(6),
  "completed_at" timestamptz(6),
  "cancelled_at" timestamptz(6)
)
;
COMMENT ON COLUMN "public"."inbound_order"."supplier_id" IS '供应商 ID，必须来自供应商主数据 supplier.id；旧 supplier 文本字段已删除';
COMMENT ON TABLE "public"."inbound_order" IS '入库单主表';

-- ----------------------------
-- Records of inbound_order
-- ----------------------------

-- ----------------------------
-- Table structure for inbound_order_detail
-- ----------------------------
DROP TABLE IF EXISTS "public"."inbound_order_detail";
CREATE TABLE "public"."inbound_order_detail" (
  "id" int8 NOT NULL DEFAULT nextval('inbound_order_detail_id_seq'::regclass),
  "order_id" int8 NOT NULL,
  "drug_id" int8 NOT NULL,
  "batch_number" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "expire_date" date NOT NULL,
  "planned_qty" int4 NOT NULL,
  "confirmed_qty" int4 DEFAULT 0,
  "unit_price" numeric(10,2),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "amount" numeric(12,2),
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."inbound_order_detail" IS '入库单明细表';

-- ----------------------------
-- Records of inbound_order_detail
-- ----------------------------

-- ----------------------------
-- Table structure for inventory_adjustment
-- ----------------------------
DROP TABLE IF EXISTS "public"."inventory_adjustment";
CREATE TABLE "public"."inventory_adjustment" (
  "id" int8 NOT NULL DEFAULT nextval('inventory_adjustment_id_seq'::regclass),
  "adjust_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "adjust_type" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "drug_id" int8,
  "from_location_id" int8,
  "to_location_id" int8,
  "before_status" varchar(20) COLLATE "pg_catalog"."default",
  "after_status" varchar(20) COLLATE "pg_catalog"."default",
  "reason" text COLLATE "pg_catalog"."default" NOT NULL,
  "operator_id" int8 NOT NULL,
  "related_task_id" int8,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'COMPLETED'::character varying,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."inventory_adjustment" IS '库存调整记录表';

-- ----------------------------
-- Records of inventory_adjustment
-- ----------------------------

-- ----------------------------
-- Table structure for inventory_task
-- ----------------------------
DROP TABLE IF EXISTS "public"."inventory_task";
CREATE TABLE "public"."inventory_task" (
  "id" int8 NOT NULL DEFAULT nextval('inventory_task_id_seq'::regclass),
  "task_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "scope_type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "scope_value" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "creator_id" int8 NOT NULL,
  "assignee_id" int8,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "start_time" timestamptz(6),
  "end_time" timestamptz(6),
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."inventory_task" IS '盘库任务主表';

-- ----------------------------
-- Records of inventory_task
-- ----------------------------

-- ----------------------------
-- Table structure for inventory_task_detail
-- ----------------------------
DROP TABLE IF EXISTS "public"."inventory_task_detail";
CREATE TABLE "public"."inventory_task_detail" (
  "id" int8 NOT NULL DEFAULT nextval('inventory_task_detail_id_seq'::regclass),
  "task_id" int8 NOT NULL,
  "location_id" int8 NOT NULL,
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "discrepancy_type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "scanned_location_id" int8,
  "system_location_id" int8,
  "operator_id" int8,
  "scanned_at" timestamptz(6)
)
;
COMMENT ON COLUMN "public"."inventory_task_detail"."location_id" IS '兼容旧字段；建议新逻辑使用 scanned_location_id / system_location_id';
COMMENT ON COLUMN "public"."inventory_task_detail"."scanned_location_id" IS '实际扫描时所在货位 ID';
COMMENT ON COLUMN "public"."inventory_task_detail"."system_location_id" IS '系统记录的应在货位 ID；UNEXPECTED 情况可为空';
COMMENT ON COLUMN "public"."inventory_task_detail"."operator_id" IS '扫码操作员 ID；由后端根据当前 JWT 用户写入';
COMMENT ON COLUMN "public"."inventory_task_detail"."scanned_at" IS '扫码时间';
COMMENT ON TABLE "public"."inventory_task_detail" IS '盘库任务明细表';

-- ----------------------------
-- Records of inventory_task_detail
-- ----------------------------

-- ----------------------------
-- Table structure for location_info
-- ----------------------------
DROP TABLE IF EXISTS "public"."location_info";
CREATE TABLE "public"."location_info" (
  "id" int8 NOT NULL DEFAULT nextval('location_info_id_seq'::regclass),
  "location_code" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "location_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "area" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "shelf" varchar(20) COLLATE "pg_catalog"."default",
  "layer" varchar(20) COLLATE "pg_catalog"."default",
  "position" varchar(20) COLLATE "pg_catalog"."default",
  "status" int2 DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "capacity" int4,
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."location_info" IS '货位信息表';

-- ----------------------------
-- Records of location_info
-- ----------------------------

-- ----------------------------
-- Table structure for login_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."login_log";
CREATE TABLE "public"."login_log" (
  "id" int8 NOT NULL DEFAULT nextval('login_log_id_seq'::regclass),
  "user_id" int8,
  "username" varchar(50) COLLATE "pg_catalog"."default",
  "success" bool NOT NULL,
  "ip" varchar(64) COLLATE "pg_catalog"."default",
  "user_agent" text COLLATE "pg_catalog"."default",
  "message" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."login_log" IS '登录日志表';

-- ----------------------------
-- Records of login_log
-- ----------------------------

-- ----------------------------
-- Table structure for notification
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification";
CREATE TABLE "public"."notification" (
  "id" int8 NOT NULL DEFAULT nextval('notification_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "title" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "notification_type" varchar(50) COLLATE "pg_catalog"."default",
  "business_type" varchar(50) COLLATE "pg_catalog"."default",
  "business_id" varchar(100) COLLATE "pg_catalog"."default",
  "read_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."notification" IS '用户通知表';

-- ----------------------------
-- Records of notification
-- ----------------------------

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."operation_log";
CREATE TABLE "public"."operation_log" (
  "id" int8 NOT NULL DEFAULT nextval('operation_log_id_seq'::regclass),
  "business_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "business_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "action" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "operator_id" int8 NOT NULL,
  "detail" jsonb,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "module" varchar(50) COLLATE "pg_catalog"."default",
  "resource_type" varchar(50) COLLATE "pg_catalog"."default",
  "resource_id" varchar(100) COLLATE "pg_catalog"."default",
  "before_data" jsonb,
  "after_data" jsonb,
  "ip" varchar(64) COLLATE "pg_catalog"."default",
  "user_agent" text COLLATE "pg_catalog"."default",
  "request_id" varchar(100) COLLATE "pg_catalog"."default"
)
;
COMMENT ON TABLE "public"."operation_log" IS '操作日志表';

-- ----------------------------
-- Records of operation_log
-- ----------------------------

-- ----------------------------
-- Table structure for report_export_task
-- ----------------------------
DROP TABLE IF EXISTS "public"."report_export_task";
CREATE TABLE "public"."report_export_task" (
  "id" int8 NOT NULL DEFAULT nextval('report_export_task_id_seq'::regclass),
  "task_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "report_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "export_format" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'xlsx'::character varying,
  "query_params" jsonb,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'PENDING'::character varying,
  "file_id" varchar(100) COLLATE "pg_catalog"."default",
  "message" text COLLATE "pg_catalog"."default",
  "requested_by" int8 NOT NULL,
  "started_at" timestamptz(6),
  "finished_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."report_export_task" IS '报表导出任务表';

-- ----------------------------
-- Records of report_export_task
-- ----------------------------

-- ----------------------------
-- Table structure for sales_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."sales_order";
CREATE TABLE "public"."sales_order" (
  "id" int8 NOT NULL DEFAULT nextval('sales_order_id_seq'::regclass),
  "order_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "cashier_id" int8 NOT NULL,
  "total_amount" numeric(10,2) NOT NULL,
  "medicare_amount" numeric(10,2) DEFAULT 0,
  "personal_amount" numeric(10,2) DEFAULT 0,
  "need_audit" bool DEFAULT false,
  "need_medicare" bool DEFAULT false,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "medicare_transaction_id" varchar(100) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "customer_name" varchar(100) COLLATE "pg_catalog"."default",
  "is_prescription" bool NOT NULL DEFAULT false,
  "discount_amount" numeric(10,2) NOT NULL DEFAULT 0,
  "actual_amount" numeric(10,2),
  "payment_method" varchar(30) COLLATE "pg_catalog"."default",
  "paid_at" timestamptz(6),
  "cancelled_at" timestamptz(6),
  "refunded_at" timestamptz(6),
  "refund_amount" numeric(10,2) NOT NULL DEFAULT 0,
  "refund_reason" text COLLATE "pg_catalog"."default",
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."sales_order"."need_audit" IS '是否需要药师审核；内部派生字段，由后端根据 is_prescription 和药品属性计算，前端不得直接控制';
COMMENT ON COLUMN "public"."sales_order"."is_prescription" IS '是否处方销售；处方药销售必须进入药师审核流程';
COMMENT ON TABLE "public"."sales_order" IS '销售订单主表';

-- ----------------------------
-- Records of sales_order
-- ----------------------------

-- ----------------------------
-- Table structure for sales_order_item
-- ----------------------------
DROP TABLE IF EXISTS "public"."sales_order_item";
CREATE TABLE "public"."sales_order_item" (
  "id" int8 NOT NULL DEFAULT nextval('sales_order_item_id_seq'::regclass),
  "order_id" int8 NOT NULL,
  "drug_id" int8 NOT NULL,
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "quantity" int4 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "subtotal_amount" numeric(10,2),
  "remark" text COLLATE "pg_catalog"."default",
  "refund_status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'NONE'::character varying,
  "refund_amount" numeric(10,2) NOT NULL DEFAULT 0,
  "refunded_at" timestamptz(6),
  "refund_reason" text COLLATE "pg_catalog"."default",
  "refund_operator_id" int8
)
;
COMMENT ON COLUMN "public"."sales_order_item"."quantity" IS '固定为 1；每条销售明细唯一对应一个追溯码，即一盒药';
COMMENT ON COLUMN "public"."sales_order_item"."refund_status" IS '明细级退货状态：NONE 未退货，REFUNDED 已退货；用于支持部分退货';
COMMENT ON COLUMN "public"."sales_order_item"."refund_amount" IS '该明细已退金额，未退货为 0';
COMMENT ON COLUMN "public"."sales_order_item"."refunded_at" IS '该明细退货时间';
COMMENT ON COLUMN "public"."sales_order_item"."refund_reason" IS '该明细退货原因';
COMMENT ON COLUMN "public"."sales_order_item"."refund_operator_id" IS '执行该明细退货的操作员 ID';
COMMENT ON TABLE "public"."sales_order_item" IS '销售订单明细表';

-- ----------------------------
-- Records of sales_order_item
-- ----------------------------

-- ----------------------------
-- Table structure for scan_task
-- ----------------------------
DROP TABLE IF EXISTS "public"."scan_task";
CREATE TABLE "public"."scan_task" (
  "id" int8 NOT NULL DEFAULT nextval('scan_task_id_seq'::regclass),
  "task_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "task_type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "related_id" int8 NOT NULL,
  "operator_id" int8 NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "start_time" timestamptz(6),
  "end_time" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."scan_task"."operator_id" IS '操作员 ID；由后端根据当前 JWT 用户写入，前端不得传入';
COMMENT ON TABLE "public"."scan_task" IS '扫码作业任务表';

-- ----------------------------
-- Records of scan_task
-- ----------------------------

-- ----------------------------
-- Table structure for scan_task_detail
-- ----------------------------
DROP TABLE IF EXISTS "public"."scan_task_detail";
CREATE TABLE "public"."scan_task_detail" (
  "id" int8 NOT NULL DEFAULT nextval('scan_task_detail_id_seq'::regclass),
  "task_id" int8 NOT NULL,
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "location_code" varchar(50) COLLATE "pg_catalog"."default",
  "scan_result" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "error_msg" varchar(255) COLLATE "pg_catalog"."default",
  "scan_time" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."scan_task_detail" IS '扫码作业明细表';

-- ----------------------------
-- Records of scan_task_detail
-- ----------------------------

-- ----------------------------
-- Table structure for security_event
-- ----------------------------
DROP TABLE IF EXISTS "public"."security_event";
CREATE TABLE "public"."security_event" (
  "id" int8 NOT NULL DEFAULT nextval('security_event_id_seq'::regclass),
  "event_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "severity" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "user_id" int8,
  "username" varchar(50) COLLATE "pg_catalog"."default",
  "ip" varchar(64) COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default" NOT NULL,
  "detail" jsonb,
  "handled" bool NOT NULL DEFAULT false,
  "handled_by" int8,
  "handled_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."security_event" IS '安全事件表';

-- ----------------------------
-- Records of security_event
-- ----------------------------

-- ----------------------------
-- Table structure for supplier
-- ----------------------------
DROP TABLE IF EXISTS "public"."supplier";
CREATE TABLE "public"."supplier" (
  "id" int8 NOT NULL DEFAULT nextval('supplier_id_seq'::regclass),
  "supplier_code" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "contact_name" varchar(50) COLLATE "pg_catalog"."default",
  "contact_phone" varchar(30) COLLATE "pg_catalog"."default",
  "license_no" varchar(100) COLLATE "pg_catalog"."default",
  "address" varchar(255) COLLATE "pg_catalog"."default",
  "status" int2 NOT NULL DEFAULT 1,
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."supplier" IS '供应商表';

-- ----------------------------
-- Records of supplier
-- ----------------------------

-- ----------------------------
-- Table structure for sys_permission
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_permission";
CREATE TABLE "public"."sys_permission" (
  "id" int8 NOT NULL DEFAULT nextval('sys_permission_id_seq'::regclass),
  "code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(150) COLLATE "pg_catalog"."default" NOT NULL,
  "resource" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "action" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "status" int2 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."sys_permission" IS '系统权限表，表示业务能力，不直接等同于单个 API';

-- ----------------------------
-- Records of sys_permission
-- ----------------------------
INSERT INTO "public"."sys_permission" VALUES (1, 'alerts.create', '处理预警', 'alerts', 'create', '将预警状态更新为 RESOLVED，记录处理人和处理备注。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (2, 'alerts.read', '预警列表', 'alerts', 'read', '查询系统预警记录，支持按类型、状态、优先级过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (3, 'audit.data_change.read', '分页查询数据变更日志', 'audit-data-change-logs', 'read', '分页查询数据变更日志', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (4, 'audit.login.read', '分页查询登录日志', 'audit-login-logs', 'read', '分页查询登录日志', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (5, 'audit.operation.read', '分页查询操作审计日志', 'audit-operation-logs', 'read', '分页查询操作审计日志', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (6, 'audit.security.read', '分页查询安全事件', 'audit-security-events', 'read', '分页查询安全事件', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (7, 'auth.self.read', '获取当前用户信息', 'auth', 'self', '返回当前登录用户的基本信息和角色。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (8, 'auth.self.update', '用户登出', 'auth', 'self', '使当前 access_token 失效（服务端加入黑名单或清除 session）。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (9, 'dashboard.read', '首页看板概览', 'dashboard', 'read', '返回首页所需的核心 KPI 数据，包括今日销售额、库存总量、近效期预警等。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (10, 'drugs.create', '创建药品', 'drugs', 'create', '新增药品基础资料。drug_code 全局唯一，重复时返回 409。仅 ADMIN 和 PHARMACIST 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (11, 'drugs.delete', '删除药品', 'drugs', 'delete', '软删除药品基础资料。若该药品有关联的在库追溯码，则不允许删除（返回 409）。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (12, 'drugs.read', '药品列表', 'drugs', 'read', '分页查询药品基础信息，支持按通用名、商品名、厂家、状态过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (13, 'drugs.status.update', '启用/停用药品', 'drugs', 'status.update', '切换药品启用状态。停用后不可新建入库单和销售单。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (14, 'drugs.update', '更新药品信息', 'drugs', 'update', '修改药品资料，drug_code 不可修改。仅 ADMIN 和 PHARMACIST 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (15, 'files.create', '上传文件', 'files', 'create', '上传发票图片、PDF 等文件，返回文件 ID 和访问 URL，供后续 AI 识别接口使用。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (16, 'files.read', '获取文件信息', 'files', 'read', '获取文件信息', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (17, 'inbound.orders.cancel', '取消入库单', 'inbound-orders', 'cancel', '取消入库单，仅 DRAFT 或 PENDING_CONFIRM 状态可取消。取消后已确认的追溯码记录将被删除。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (18, 'inbound.orders.complete', '完成入库', 'inbound-orders', 'complete', '将 PENDING_CONFIRM 状态的入库单标记为完成（COMPLETED）。
完成后，所有 confirmed_qty > 0 的明细行对应的追溯码保持 PENDING 状态，
等待上架操作。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (19, 'inbound.orders.create', '创建入库单', 'inbound-orders', 'create', '创建采购入库单（草稿状态）。
可手动填写，也可由 AI 发票识别模块自动生成。
创建后状态为 DRAFT，需调用提交接口流转至 PENDING_CONFIRM。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (20, 'inbound.orders.delete', '删除入库明细行', 'inbound-orders', 'delete', '删除指定明细行。仅 DRAFT 状态可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (21, 'inbound.orders.read', '入库单列表', 'inbound-orders', 'read', '分页查询入库单，支持按状态、供应商、日期范围过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (22, 'inbound.orders.update', '更新入库单基本信息', 'inbound-orders', 'update', '仅允许在 DRAFT 状态下修改供应商、发票号、备注等基本信息。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (23, 'inventory.adjustment.create', '创建库存调整', 'inventory-adjustments', 'create', '创建库存调整', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (24, 'inventory.adjustment.read', '分页查询库存调整记录', 'inventory-adjustments', 'read', '分页查询库存调整记录', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (25, 'inventory.loss.confirm', '确认盘亏候选', 'inventory-loss', 'confirm', '确认盘亏候选', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (26, 'inventory.loss.reject', '驳回盘亏候选', 'inventory-loss', 'reject', '驳回盘亏候选', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (27, 'inventory.misplaced.relocate', '处理错架并调整货位', 'inventory-misplaced', 'relocate', '处理错架并调整货位', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (28, 'inventory.tasks.cancel', '取消盘库任务', 'inventory-tasks', 'cancel', '取消 PENDING 或 IN_PROGRESS 状态的盘库任务。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (29, 'inventory.tasks.complete', '完成盘库任务', 'inventory-tasks', 'complete', '完成盘库任务，触发以下操作：
1. 计算盘亏候选：系统范围内 IN_STOCK 记录 - 实际扫描集合 = 盘亏候选
2. 将盘亏候选追溯码状态更新为 LOSS_CANDIDATE
3. 将任务状态更新为 COMPLETED，记录 end_time
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (30, 'inventory.tasks.create', '创建盘库任务', 'inventory-tasks', 'create', '创建盘库任务，指定盘点范围（区域/货架/货位）。
创建后状态为 PENDING，需分配执行人后才可开始。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (31, 'inventory.tasks.read', '盘库任务列表', 'inventory-tasks', 'read', '分页查询盘库任务，支持按状态、范围类型、创建人过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (32, 'inventory.tasks.scan', '盘库扫码提交', 'inventory-tasks', 'scan', '盘库过程中提交单次扫码结果。
系统根据追溯码的系统记录与实际扫描位置判断差异类型：
- NORMAL：追溯码在系统记录的货位扫到
- MISPLACED_FOUND：追溯码在其他货位扫到（系统记录在别处）
- UNEXPECTED：追溯码在系统中无 IN_STOCK 记录
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (33, 'inventory.tasks.start', '开始盘库任务', 'inventory-tasks', 'start', '将任务状态从 PENDING 更新为 IN_PROGRESS，记录 start_time。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (34, 'locations.create', '创建货位', 'locations', 'create', '新增货位。location_code 全局唯一。仅 ADMIN 和 WAREHOUSE 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (35, 'locations.delete', '删除货位', 'locations', 'delete', '软删除货位。若货位有在库药品，则不允许删除（返回 409）。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (36, 'locations.read', '货位列表', 'locations', 'read', '分页查询货位信息，支持按区域、货架、状态过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (37, 'locations.status.update', '启用/停用货位', 'locations', 'status.update', '启用/停用货位', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (38, 'locations.update', '更新货位信息', 'locations', 'update', '修改货位名称、区域等信息，location_code 不可修改。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (39, 'pharmacist.review.approve', '药师审核通过', 'pharmacist-reviews', 'approve', '药师审核通过', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (40, 'pharmacist.review.read', '分页查询药师审核任务', 'pharmacist-reviews', 'read', '分页查询药师审核任务', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (41, 'pharmacist.review.reject', '药师审核驳回', 'pharmacist-reviews', 'reject', '药师审核驳回', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (42, 'report.export_task.read', '查询报表导出任务', 'report-export-tasks', 'read', '查询报表导出任务', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (43, 'report.inbound.export', '导出入库报表', 'report-exports', 'export', '导出入库报表', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (44, 'report.inventory.export', '导出库存报表', 'report-exports', 'export', '导出库存报表', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (45, 'report.sales.export', '导出销售报表', 'report-exports', 'export', '导出销售报表', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (46, 'report.trace-log.export', '导出追溯日志报表', 'report-exports', 'export', '导出追溯日志报表', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (47, 'reports.read', '销售报表', 'reports', 'read', '按日期范围、药品、收银员等维度查询销售报表数据，支持导出 Excel。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (48, 'sales.orders.cancel', '取消销售单', 'sales-orders', 'cancel', '取消 PENDING 状态的销售单，释放所有已预占的追溯码。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (49, 'sales.orders.create', '创建销售单', 'sales-orders', 'create', '创建销售单（PENDING 状态）。
创建时可指定药品列表，系统根据近效期优先原则自动匹配追溯码。
若需手动指定追溯码，可通过 items[].trace_code 字段传入。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (50, 'sales.orders.delete', '删除销售明细行', 'sales-orders', 'delete', '从 PENDING 状态的销售单中删除指定明细行，释放已预占的追溯码。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (51, 'sales.orders.pay', '销售结算', 'sales-orders', 'pay', '对 PENDING 状态的销售单进行结算，触发以下操作：
1. 将销售单状态更新为 COMPLETED
2. 将所有明细行关联的追溯码状态更新为 SOLD
3. 在 drug_trace_log 中写入 SALE 动作记录
4. 记录支付方式和实收金额
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (52, 'sales.orders.read', '销售单列表', 'sales-orders', 'read', '分页查询销售单，支持按状态、收银员、日期范围过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (53, 'sales.orders.refund', '销售退货', 'sales-orders', 'refund', '对 COMPLETED 状态的销售单发起退货，触发以下操作：
1. 将销售单状态更新为 REFUNDED
2. 将退货明细对应的追溯码状态从 SOLD 恢复为 IN_STOCK
3. 在 drug_trace_log 中写入 RETURN 动作记录
4. 支持部分退货（通过 detail_ids 指定退货明细）
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (54, 'sales.orders.scan', '销售追溯码扫码验证', 'sales-orders', 'scan', '销售出库时扫描追溯码进行验证，确认追溯码与销售明细匹配。
验证通过后不立即更新状态，需调用结算接口统一提交。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (55, 'sales.review.read', '查询销售单审核记录', 'sales-reviews', 'read', '查询销售单审核记录', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (56, 'sales.review.submit', '提交销售单药师审核', 'sales-reviews', 'submit', '提交销售单药师审核', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (57, 'sales.trace.release', '释放销售追溯码锁定', 'sales-reservations', 'release', '释放销售追溯码锁定', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (58, 'sales.trace.reserve', '锁定销售追溯码', 'sales-reservations', 'reserve', '锁定销售追溯码', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (59, 'sales.trace.reserve.read', '查询销售单已锁定追溯码', 'sales-reservations', 'read', '查询销售单已锁定追溯码', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (60, 'scan.tasks.read', '扫码任务列表', 'scan-tasks', 'read', '查询扫码作业任务，支持按类型、状态、操作员过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (62, 'shelving.create', '批量上架', 'shelving', 'create', '将多个追溯码批量上架至同一货位。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (63, 'shelving.read', '待上架追溯码列表', 'shelving', 'read', '查询所有状态为 PENDING 的追溯码，支持按药品、入库单过滤，用于生成上架任务。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (64, 'shelving.scan', '单个追溯码上架', 'shelving', 'scan', '先扫货位码，再扫追溯码，完成上架操作。
系统将追溯码状态从 PENDING 更新为 IN_STOCK，并记录 location_id。
同时在 drug_trace_log 中写入 SHELVING 动作记录。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (65, 'supplier.create', '创建供应商', 'suppliers', 'create', '创建供应商', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (66, 'supplier.delete', '删除供应商', 'suppliers', 'delete', '删除供应商', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (67, 'supplier.read', '分页查询供应商', 'suppliers', 'read', '分页查询供应商', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (68, 'supplier.status.update', '启用/禁用供应商', 'suppliers', 'status.update', '启用/禁用供应商', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (69, 'supplier.update', '更新供应商', 'suppliers', 'update', '更新供应商', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (70, 'system.file.download', '下载文件', 'files', 'download', '下载文件', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (71, 'system.permission.read', '查询权限字典', 'permissions', 'read', '查询权限字典', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (72, 'system.role.create', '创建角色', 'roles', 'create', '创建角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (73, 'system.role.delete', '删除角色', 'roles', 'delete', '删除角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (74, 'system.role.permission.assign', '分配角色权限', 'role-permissions', 'assign', '分配角色权限', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (75, 'system.role.permission.read', '查询角色权限', 'role-permissions', 'read', '查询角色权限', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (76, 'system.role.read', '分页查询角色', 'roles', 'read', '分页查询角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (77, 'system.role.update', '更新角色', 'roles', 'update', '更新角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (78, 'system.user.permission.read', '查询用户有效权限', 'user-permissions', 'read', '查询用户有效权限', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (79, 'system.user.role.assign', '分配用户角色', 'user-roles', 'assign', '分配用户角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (80, 'system.user.role.read', '查询用户角色', 'user-roles', 'read', '查询用户角色', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (81, 'trace.inventory.create', '追溯码有效性验证', 'trace-inventory', 'create', '验证追溯码是否存在于系统中，以及当前状态是否允许指定操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (82, 'trace.inventory.read', '追溯码库存列表', 'trace-inventory', 'read', '分页查询追溯码库存记录，支持多维度过滤。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (83, 'trace.inventory.status.update', '手动更新追溯码状态', 'trace-inventory', 'status.update', '管理员手动修正追溯码状态（如将 MISPLACED 归位为 IN_STOCK）。
状态流转必须符合合法路径，否则返回 409。
', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (84, 'users.create', '创建用户', 'users', 'create', '创建新系统用户。仅 ADMIN 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (85, 'users.password.reset', '重置用户密码', 'users', 'password.reset', '管理员重置指定用户的密码。仅 ADMIN 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (86, 'users.read', '用户列表', 'users', 'read', '分页查询系统用户，支持按角色、状态、姓名过滤。仅 ADMIN 可访问。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (87, 'users.status.update', '启用/禁用用户', 'users', 'status.update', '切换用户账号状态。仅 ADMIN 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (88, 'users.update', '更新用户信息', 'users', 'update', '修改用户真实姓名、角色等信息。仅 ADMIN 可操作。', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (89, 'scan.tasks.create', '创建扫码任务', 'scan-tasks', 'create', '为入库、上架、盘库等业务创建扫码作业任务。', 1, '2026-06-08 16:34:39.566604+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (90, 'scan.tasks.start', '开始扫码任务', 'scan-tasks', 'start', '将扫码任务从 PENDING 流转为 IN_PROGRESS。', 1, '2026-06-08 16:34:39.566604+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (91, 'scan.tasks.submit', '提交扫码结果', 'scan-tasks', 'submit', '提交扫码任务的扫码结果。', 1, '2026-06-08 16:34:39.566604+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (92, 'scan.tasks.complete', '完成扫码任务', 'scan-tasks', 'complete', '完成扫码任务并触发后续业务状态流转。', 1, '2026-06-08 16:34:39.566604+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (93, 'scan.tasks.cancel', '取消扫码任务', 'scan-tasks', 'cancel', '取消仍可取消的扫码任务。', 1, '2026-06-08 16:34:39.566604+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (61, 'scan.tasks.scan', '创建扫码任务', 'scan-tasks', 'scan', '已废弃：该权限码同时覆盖创建、开始、提交、完成、取消扫码任务，容易产生权限漂移；请使用 scan.tasks.create/start/submit/complete/cancel。', 0, '2026-06-07 14:15:04.180655+00', '2026-06-08 16:34:39.566604+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (94, 'ai.invoices.read', '查询发票识别记录', 'ai-invoices', 'read', '查询主系统保存的 AI 发票识别记录列表和详情。', 1, '2026-06-10 02:42:45.441234+00', '2026-06-10 02:42:45.441234+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (95, 'ai.invoices.recognize', '发票 AI 识别', 'ai-invoices', 'recognize', '由主 ERP 系统调用 AI 子模块完成发票 OCR 和字段结构化识别，并保存识别结果。', 1, '2026-06-10 02:42:45.441234+00', '2026-06-10 02:42:45.441234+00', NULL);
INSERT INTO "public"."sys_permission" VALUES (96, 'ai.invoices.convert', '发票识别结果转入库单', 'ai-invoices', 'convert', '将已确认的发票识别结果转换为采购入库单，前端必须指定 supplier_id，不再使用供应商名称作为入库单供应商来源。', 1, '2026-06-10 02:42:45.441234+00', '2026-06-10 02:42:45.441234+00', NULL);

-- ----------------------------
-- Table structure for sys_permission_api
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_permission_api";
CREATE TABLE "public"."sys_permission_api" (
  "id" int8 NOT NULL DEFAULT nextval('sys_permission_api_id_seq'::regclass),
  "permission_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "path_pattern" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "http_method" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "summary" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."sys_permission_api" IS '权限与 API 路径映射表，用于从数据库生成/校验 Casbin policy';

-- ----------------------------
-- Records of sys_permission_api
-- ----------------------------
INSERT INTO "public"."sys_permission_api" VALUES (1, 'alerts.create', '/api/v1/alerts/:id/ignore', 'POST', '忽略预警', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (2, 'alerts.create', '/api/v1/alerts/:id/resolve', 'POST', '处理预警', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (3, 'alerts.create', '/api/v1/notifications/:id/read', 'POST', '标记通知已读', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (4, 'alerts.create', '/api/v1/notifications/read-all', 'POST', '全部标记已读', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (5, 'alerts.read', '/api/v1/alerts', 'GET', '预警列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (6, 'alerts.read', '/api/v1/alerts/loss-candidates', 'GET', '盘亏候选预警列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (7, 'alerts.read', '/api/v1/alerts/near-expire', 'GET', '近效期预警列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (8, 'alerts.read', '/api/v1/notifications', 'GET', '通知消息列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (9, 'alerts.read', '/api/v1/notifications/unread-count', 'GET', '未读通知数量', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (10, 'audit.data_change.read', '/api/v1/audit/data-change-logs', 'GET', '分页查询数据变更日志', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (11, 'audit.login.read', '/api/v1/audit/login-logs', 'GET', '分页查询登录日志', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (12, 'audit.operation.read', '/api/v1/audit/operation-logs', 'GET', '分页查询操作审计日志', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (13, 'audit.operation.read', '/api/v1/audit/operation-logs/:id', 'GET', '获取操作审计日志详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (14, 'audit.security.read', '/api/v1/audit/security-events', 'GET', '分页查询安全事件', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (15, 'auth.self.read', '/api/v1/auth/me', 'GET', '获取当前用户信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (16, 'auth.self.update', '/api/v1/auth/logout', 'POST', '用户登出', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (17, 'auth.self.update', '/api/v1/auth/password', 'PUT', '修改当前用户密码', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (18, 'dashboard.read', '/api/v1/dashboard/inbound-stats', 'GET', '入库统计', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (19, 'dashboard.read', '/api/v1/dashboard/inventory-stats', 'GET', '库存统计', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (20, 'dashboard.read', '/api/v1/dashboard/overview', 'GET', '首页看板概览', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (21, 'dashboard.read', '/api/v1/dashboard/sales-trend', 'GET', '销售趋势统计', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (22, 'dashboard.read', '/api/v1/dashboard/top-drugs', 'GET', '热销药品排行', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (23, 'drugs.create', '/api/v1/drugs', 'POST', '创建药品', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (24, 'drugs.delete', '/api/v1/drugs/:id', 'DELETE', '删除药品', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (25, 'drugs.read', '/api/v1/drugs', 'GET', '药品列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (26, 'drugs.read', '/api/v1/drugs/:id', 'GET', '药品详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (27, 'drugs.read', '/api/v1/drugs/:id/inventory-summary', 'GET', '药品库存汇总', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (28, 'drugs.read', '/api/v1/drugs/:id/sale-info', 'GET', '获取药品销售信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (29, 'drugs.read', '/api/v1/drugs/code/:drug_code', 'GET', '按药品编码查询', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (30, 'drugs.read', '/api/v1/drugs/search', 'GET', '药品模糊搜索', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (31, 'drugs.status.update', '/api/v1/drugs/:id/status', 'PATCH', '启用/停用药品', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (32, 'drugs.update', '/api/v1/drugs/:id', 'PUT', '更新药品信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (33, 'files.create', '/api/v1/files/upload', 'POST', '上传文件', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (34, 'files.read', '/api/v1/files/:file_id', 'GET', '获取文件信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (35, 'inbound.orders.cancel', '/api/v1/inbound-orders/:id/cancel', 'POST', '取消入库单', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (36, 'inbound.orders.complete', '/api/v1/inbound-orders/:id/complete', 'POST', '完成入库', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (37, 'inbound.orders.create', '/api/v1/inbound-orders', 'POST', '创建入库单', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (38, 'inbound.orders.create', '/api/v1/inbound-orders/:id/confirm-trace', 'POST', '单个追溯码扫码确认', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (39, 'inbound.orders.create', '/api/v1/inbound-orders/:id/confirm-traces', 'POST', '批量追溯码扫码确认', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (40, 'inbound.orders.create', '/api/v1/inbound-orders/:id/details', 'POST', '添加入库明细行', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (41, 'inbound.orders.create', '/api/v1/inbound-orders/:id/submit', 'POST', '提交入库单', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (42, 'inbound.orders.delete', '/api/v1/inbound-orders/:id/details/:detail_id', 'DELETE', '删除入库明细行', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (43, 'inbound.orders.read', '/api/v1/inbound-orders', 'GET', '入库单列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (44, 'inbound.orders.read', '/api/v1/inbound-orders/:id', 'GET', '入库单详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (45, 'inbound.orders.read', '/api/v1/inbound-orders/:id/details', 'GET', '入库单明细列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (46, 'inbound.orders.read', '/api/v1/inbound-orders/:id/details/:detail_id', 'GET', '入库明细详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (47, 'inbound.orders.read', '/api/v1/inbound-orders/:id/progress', 'GET', '入库进度查询', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (48, 'inbound.orders.update', '/api/v1/inbound-orders/:id', 'PUT', '更新入库单基本信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (49, 'inbound.orders.update', '/api/v1/inbound-orders/:id/details/:detail_id', 'PUT', '更新入库明细', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (50, 'inventory.adjustment.create', '/api/v1/inventory-adjustments', 'POST', '创建库存调整', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (51, 'inventory.adjustment.read', '/api/v1/inventory-adjustments', 'GET', '分页查询库存调整记录', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (52, 'inventory.adjustment.read', '/api/v1/inventory-adjustments/:id', 'GET', '获取库存调整详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (53, 'inventory.loss.confirm', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/confirm', 'POST', '确认盘亏候选', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (54, 'inventory.loss.reject', '/api/v1/inventory-tasks/:id/loss-candidates/:trace_code/reject', 'POST', '驳回盘亏候选', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (55, 'inventory.misplaced.relocate', '/api/v1/inventory-tasks/:id/misplaced/:trace_code/relocate', 'POST', '处理错架并调整货位', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (56, 'inventory.tasks.cancel', '/api/v1/inventory-tasks/:id/cancel', 'POST', '取消盘库任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (57, 'inventory.tasks.complete', '/api/v1/inventory-tasks/:id/complete', 'POST', '完成盘库任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (58, 'inventory.tasks.create', '/api/v1/inventory-tasks', 'POST', '创建盘库任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (59, 'inventory.tasks.create', '/api/v1/inventory-tasks/:id/assign', 'POST', '分配盘库执行人', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (60, 'inventory.tasks.read', '/api/v1/inventory-tasks', 'GET', '盘库任务列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (61, 'inventory.tasks.read', '/api/v1/inventory-tasks/:id', 'GET', '盘库任务详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (62, 'inventory.tasks.read', '/api/v1/inventory-tasks/:id/details', 'GET', '盘库明细列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (63, 'inventory.tasks.read', '/api/v1/inventory-tasks/:id/loss-candidates', 'GET', '盘亏候选列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (64, 'inventory.tasks.read', '/api/v1/inventory-tasks/:id/misplaced', 'GET', '错架药品列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (65, 'inventory.tasks.read', '/api/v1/inventory-tasks/:id/summary', 'GET', '盘库任务汇总', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (66, 'inventory.tasks.scan', '/api/v1/inventory-tasks/:id/scan', 'POST', '盘库扫码提交', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (67, 'inventory.tasks.start', '/api/v1/inventory-tasks/:id/start', 'POST', '开始盘库任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (68, 'locations.create', '/api/v1/locations', 'POST', '创建货位', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (69, 'locations.delete', '/api/v1/locations/:id', 'DELETE', '删除货位', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (70, 'locations.read', '/api/v1/locations', 'GET', '货位列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (71, 'locations.read', '/api/v1/locations/:id', 'GET', '货位详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (72, 'locations.read', '/api/v1/locations/:id/drugs', 'GET', '查询货位当前在库药品', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (73, 'locations.read', '/api/v1/locations/areas', 'GET', '获取所有区域列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (74, 'locations.read', '/api/v1/locations/code/:location_code', 'GET', '按货位编码查询', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (75, 'locations.status.update', '/api/v1/locations/:id/status', 'PATCH', '启用/停用货位', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (76, 'locations.update', '/api/v1/locations/:id', 'PUT', '更新货位信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (77, 'pharmacist.review.approve', '/api/v1/pharmacist/reviews/:id/approve', 'POST', '药师审核通过', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (78, 'pharmacist.review.read', '/api/v1/pharmacist/reviews', 'GET', '分页查询药师审核任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (79, 'pharmacist.review.read', '/api/v1/pharmacist/reviews/:id', 'GET', '获取药师审核详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (80, 'pharmacist.review.reject', '/api/v1/pharmacist/reviews/:id/reject', 'POST', '药师审核驳回', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (81, 'report.export_task.read', '/api/v1/reports/export-tasks/:task_id', 'GET', '查询报表导出任务', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (82, 'report.inbound.export', '/api/v1/reports/inbound/export', 'POST', '导出入库报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (83, 'report.inventory.export', '/api/v1/reports/inventory/export', 'POST', '导出库存报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (84, 'report.sales.export', '/api/v1/reports/sales/export', 'POST', '导出销售报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (85, 'report.trace-log.export', '/api/v1/reports/trace-log/export', 'POST', '导出追溯日志报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (86, 'reports.read', '/api/v1/reports/inbound', 'GET', '入库报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (87, 'reports.read', '/api/v1/reports/inventory', 'GET', '库存报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (88, 'reports.read', '/api/v1/reports/sales', 'GET', '销售报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (89, 'reports.read', '/api/v1/reports/trace-log', 'GET', '追溯日志报表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (90, 'sales.orders.cancel', '/api/v1/sales-orders/:id/cancel', 'POST', '取消销售单', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (91, 'sales.orders.create', '/api/v1/sales-orders', 'POST', '创建销售单', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (92, 'sales.orders.create', '/api/v1/sales-orders/:id/details', 'POST', '添加销售明细行', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (93, 'sales.orders.delete', '/api/v1/sales-orders/:id/details/:detail_id', 'DELETE', '删除销售明细行', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (94, 'sales.orders.pay', '/api/v1/sales-orders/:id/pay', 'POST', '销售结算', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (95, 'sales.orders.read', '/api/v1/sales-orders', 'GET', '销售单列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (96, 'sales.orders.read', '/api/v1/sales-orders/:id', 'GET', '销售单详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (97, 'sales.orders.read', '/api/v1/sales-orders/:id/details', 'GET', '销售单明细列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (98, 'sales.orders.refund', '/api/v1/sales-orders/:id/refund', 'POST', '销售退货', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (99, 'sales.orders.scan', '/api/v1/sales-orders/:id/scan-verify', 'POST', '销售追溯码扫码验证', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (100, 'sales.review.read', '/api/v1/sales-orders/:id/review-record', 'GET', '查询销售单审核记录', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (101, 'sales.review.submit', '/api/v1/sales-orders/:id/submit-review', 'POST', '提交销售单药师审核', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (102, 'sales.trace.release', '/api/v1/sales-orders/:id/release-reservation', 'POST', '释放销售追溯码锁定', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (103, 'sales.trace.reserve', '/api/v1/sales-orders/:id/reserve-trace', 'POST', '锁定销售追溯码', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (104, 'sales.trace.reserve.read', '/api/v1/sales-orders/:id/reserved-traces', 'GET', '查询销售单已锁定追溯码', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (105, 'scan.tasks.read', '/api/v1/scan-tasks', 'GET', '扫码任务列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (106, 'scan.tasks.read', '/api/v1/scan-tasks/:id', 'GET', '扫码任务详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (107, 'scan.tasks.read', '/api/v1/scan-tasks/:id/details', 'GET', '扫码明细列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (113, 'shelving.create', '/api/v1/shelving/batch', 'POST', '批量上架', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (114, 'shelving.create', '/api/v1/shelving/relocate', 'POST', '货位调拨', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (115, 'shelving.read', '/api/v1/shelving/mix-check', 'GET', '货位混放检查', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (116, 'shelving.read', '/api/v1/shelving/pending', 'GET', '待上架追溯码列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (117, 'shelving.scan', '/api/v1/shelving/scan', 'POST', '单个追溯码上架', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (118, 'supplier.create', '/api/v1/suppliers', 'POST', '创建供应商', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (119, 'supplier.delete', '/api/v1/suppliers/:id', 'DELETE', '删除供应商', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (120, 'supplier.read', '/api/v1/suppliers', 'GET', '分页查询供应商', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (121, 'supplier.read', '/api/v1/suppliers/:id', 'GET', '获取供应商详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (122, 'supplier.status.update', '/api/v1/suppliers/:id/status', 'PATCH', '启用/禁用供应商', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (123, 'supplier.update', '/api/v1/suppliers/:id', 'PUT', '更新供应商', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (124, 'system.file.download', '/api/v1/files/:file_id/download', 'GET', '下载文件', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (125, 'system.permission.read', '/api/v1/permissions', 'GET', '查询权限字典', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (126, 'system.role.create', '/api/v1/roles', 'POST', '创建角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (127, 'system.role.delete', '/api/v1/roles/:id', 'DELETE', '删除角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (128, 'system.role.permission.assign', '/api/v1/roles/:id/permissions', 'PUT', '分配角色权限', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (129, 'system.role.permission.read', '/api/v1/roles/:id/permissions', 'GET', '查询角色权限', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (130, 'system.role.read', '/api/v1/roles', 'GET', '分页查询角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (131, 'system.role.read', '/api/v1/roles/:id', 'GET', '获取角色详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (132, 'system.role.update', '/api/v1/roles/:id', 'PUT', '更新角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (133, 'system.user.permission.read', '/api/v1/users/:id/permissions', 'GET', '查询用户有效权限', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (134, 'system.user.role.assign', '/api/v1/users/:id/roles', 'PUT', '分配用户角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (135, 'system.user.role.read', '/api/v1/users/:id/roles', 'GET', '查询用户角色', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (136, 'trace.inventory.create', '/api/v1/trace/validate', 'POST', '追溯码有效性验证', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (137, 'trace.inventory.read', '/api/v1/inventory', 'GET', '追溯码库存列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (138, 'trace.inventory.read', '/api/v1/inventory/drugs/:drug_id', 'GET', '指定药品的库存列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (139, 'trace.inventory.read', '/api/v1/inventory/locations/:location_id', 'GET', '指定货位的库存列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (140, 'trace.inventory.read', '/api/v1/inventory/near-expire', 'GET', '近效期药品列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (141, 'trace.inventory.read', '/api/v1/inventory/pending-shelving', 'GET', '待上架库存列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (142, 'trace.inventory.read', '/api/v1/inventory/recommend-sale', 'GET', '推荐出库追溯码', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (143, 'trace.inventory.read', '/api/v1/inventory/summary', 'GET', '库存总览统计', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (144, 'trace.inventory.read', '/api/v1/trace/:trace_code', 'GET', '追溯码当前状态查询', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (145, 'trace.inventory.read', '/api/v1/trace/:trace_code/full-chain', 'GET', '追溯码全链路追溯', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (146, 'trace.inventory.read', '/api/v1/trace/:trace_code/logs', 'GET', '追溯码操作日志', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (147, 'trace.inventory.status.update', '/api/v1/inventory/:trace_code/status', 'PATCH', '手动更新追溯码状态', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (148, 'users.create', '/api/v1/users', 'POST', '创建用户', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (149, 'users.password.reset', '/api/v1/users/:id/reset-password', 'POST', '重置用户密码', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (150, 'users.read', '/api/v1/users', 'GET', '用户列表', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (151, 'users.read', '/api/v1/users/:id', 'GET', '用户详情', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (152, 'users.status.update', '/api/v1/users/:id/status', 'PATCH', '启用/禁用用户', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (153, 'users.update', '/api/v1/users/:id', 'PUT', '更新用户信息', '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_permission_api" VALUES (154, 'scan.tasks.create', '/api/v1/scan-tasks', 'POST', '创建扫码任务', '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_permission_api" VALUES (155, 'scan.tasks.start', '/api/v1/scan-tasks/:id/start', 'POST', '开始扫码任务', '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_permission_api" VALUES (156, 'scan.tasks.submit', '/api/v1/scan-tasks/:id/submit', 'POST', '提交扫码结果', '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_permission_api" VALUES (157, 'scan.tasks.complete', '/api/v1/scan-tasks/:id/complete', 'POST', '完成扫码任务', '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_permission_api" VALUES (158, 'scan.tasks.cancel', '/api/v1/scan-tasks/:id/cancel', 'POST', '取消扫码任务', '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_permission_api" VALUES (159, 'ai.invoices.read', '/api/v1/ai/invoices', 'GET', '发票识别记录列表', '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_permission_api" VALUES (160, 'ai.invoices.recognize', '/api/v1/ai/invoices/recognize', 'POST', '调用 AI 子模块识别发票', '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_permission_api" VALUES (161, 'ai.invoices.read', '/api/v1/ai/invoices/:id', 'GET', '发票识别记录详情', '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_permission_api" VALUES (162, 'ai.invoices.convert', '/api/v1/ai/invoices/:id/convert-to-inbound', 'POST', '发票识别结果转入库单', '2026-06-10 02:42:45.441234+00');

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role";
CREATE TABLE "public"."sys_role" (
  "id" int8 NOT NULL DEFAULT nextval('sys_role_id_seq'::regclass),
  "code" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "built_in" bool NOT NULL DEFAULT false,
  "status" int2 NOT NULL DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON COLUMN "public"."sys_role"."code" IS '角色编码，例如 ADMIN、STORE_MANAGER、PHARMACIST、CASHIER、WAREHOUSE';
COMMENT ON COLUMN "public"."sys_role"."built_in" IS '是否系统内置角色';
COMMENT ON TABLE "public"."sys_role" IS '系统角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO "public"."sys_role" VALUES (2, 'STORE_MANAGER', '店长', '负责门店业务管理、人员权限、报表和基础资料维护', 'f', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_role" VALUES (3, 'PHARMACIST', '药师', '负责处方药审核、药品与库存查询', 'f', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_role" VALUES (4, 'CASHIER', '收银员', '负责销售开单、收银、顾客购药流程', 'f', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_role" VALUES (5, 'WAREHOUSE', '仓库员', '负责入库、上架、盘点、库存作业', 'f', 1, '2026-06-07 14:15:04.180655+00', '2026-06-07 14:15:04.180655+00', NULL);
INSERT INTO "public"."sys_role" VALUES (1, 'ADMIN', '系统管理员', '拥有系统全部管理权限', 't', 1, '2026-06-07 14:15:04.180655+00', '2026-06-10 08:10:49.912781+00', NULL);

-- ----------------------------
-- Table structure for sys_role_permission
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role_permission";
CREATE TABLE "public"."sys_role_permission" (
  "role_id" int8 NOT NULL,
  "permission_id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."sys_role_permission" IS '角色权限关联表';

-- ----------------------------
-- Records of sys_role_permission
-- ----------------------------
INSERT INTO "public"."sys_role_permission" VALUES (1, 88, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 87, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 86, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 85, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 84, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 83, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 82, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 81, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 80, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 79, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 78, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 77, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 76, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 75, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 74, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 73, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 72, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 71, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 70, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 69, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 68, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 67, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 66, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 65, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 64, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 63, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 62, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 60, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 59, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 58, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 57, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 56, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 55, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 54, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 53, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 52, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 51, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 50, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 49, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 48, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 47, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 46, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 45, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 44, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 43, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 42, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 41, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 40, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 39, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 38, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 37, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 36, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 35, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 34, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 33, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 32, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 31, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 30, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 29, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 28, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 27, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 26, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 25, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 24, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 23, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 22, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 21, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 20, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 19, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 18, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 17, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 16, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 15, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 14, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 13, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 12, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 11, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 10, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 9, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 8, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 7, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 6, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 5, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 4, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 3, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 2, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 1, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 86, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 83, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 82, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 81, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 80, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 78, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 76, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 75, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 71, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 70, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 69, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 68, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 67, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 66, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 65, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 64, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 63, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 62, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 60, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 59, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 58, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 57, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 56, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 55, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 54, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 53, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 52, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 51, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 50, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 49, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 48, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 47, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 46, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 45, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 44, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 43, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 42, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 40, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 38, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 37, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 36, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 35, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 34, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 33, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 32, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 31, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 30, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 29, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 28, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 27, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 26, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 25, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 24, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 23, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 22, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 21, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 20, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 19, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 18, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 17, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 16, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 15, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 14, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 13, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 12, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 11, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 10, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 9, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 8, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 7, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 6, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 5, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 4, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 3, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 2, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 1, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 82, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 70, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 59, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 55, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 52, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 41, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 40, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 39, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 36, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 24, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 16, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 15, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 12, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 8, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 7, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (3, 2, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 70, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 59, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 58, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 57, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 56, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 55, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 54, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 52, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 51, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 50, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 49, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 48, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 28, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 17, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 16, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 15, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 12, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 8, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 7, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 83, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 82, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 81, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 70, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 69, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 67, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 65, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 64, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 63, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 62, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 60, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 48, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 38, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 37, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 36, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 35, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 34, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 33, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 32, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 31, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 30, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 29, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 28, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 27, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 26, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 25, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 24, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 23, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 22, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 21, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 20, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 19, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 18, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 17, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 16, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 15, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 12, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 8, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 7, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 2, '2026-06-07 14:15:04.180655+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 93, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 92, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 91, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 90, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 89, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 93, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 92, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 91, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 90, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 89, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (4, 93, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 93, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 92, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 91, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 90, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 89, '2026-06-08 16:34:39.566604+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 94, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 95, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (1, 96, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 94, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 95, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (2, 96, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 94, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 95, '2026-06-10 02:42:45.441234+00');
INSERT INTO "public"."sys_role_permission" VALUES (5, 96, '2026-06-10 02:42:45.441234+00');

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_user";
CREATE TABLE "public"."sys_user" (
  "id" int8 NOT NULL DEFAULT nextval('sys_user_id_seq'::regclass),
  "username" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "password_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "real_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 DEFAULT 1,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6),
  "phone" varchar(30) COLLATE "pg_catalog"."default",
  "email" varchar(100) COLLATE "pg_catalog"."default",
  "avatar_url" varchar(500) COLLATE "pg_catalog"."default",
  "last_login_at" timestamptz(6),
  "last_login_ip" varchar(64) COLLATE "pg_catalog"."default",
  "remark" text COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."sys_user"."phone" IS '手机号';
COMMENT ON COLUMN "public"."sys_user"."email" IS '邮箱';
COMMENT ON COLUMN "public"."sys_user"."avatar_url" IS '头像 URL';
COMMENT ON COLUMN "public"."sys_user"."last_login_at" IS '最后登录时间';
COMMENT ON COLUMN "public"."sys_user"."last_login_ip" IS '最后登录 IP';
COMMENT ON COLUMN "public"."sys_user"."remark" IS '备注';
COMMENT ON TABLE "public"."sys_user" IS '系统用户表';

-- ----------------------------
-- Records of sys_user
-- ----------------------------

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_user_role";
CREATE TABLE "public"."sys_user_role" (
  "user_id" int8 NOT NULL,
  "role_id" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."sys_user_role" IS '用户角色关联表，一个用户可拥有多个角色';

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------

-- ----------------------------
-- Table structure for trace_reservation
-- ----------------------------
DROP TABLE IF EXISTS "public"."trace_reservation";
CREATE TABLE "public"."trace_reservation" (
  "id" int8 NOT NULL DEFAULT nextval('trace_reservation_id_seq'::regclass),
  "reservation_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "sales_order_id" int8 NOT NULL,
  "sales_order_item_id" int8,
  "trace_code" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "drug_id" int8 NOT NULL,
  "reserved_by" int8 NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'RESERVED'::character varying,
  "reserved_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "released_at" timestamptz(6),
  "confirmed_at" timestamptz(6),
  "expire_at" timestamptz(6),
  "remark" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "public"."trace_reservation" IS '销售订单追溯码预占表；预占的唯一业务来源，同一 trace_code 在 RESERVED 状态下只能存在一条有效记录';

-- ----------------------------
-- Records of trace_reservation
-- ----------------------------

-- ----------------------------
-- Function structure for refresh_casbin_rule_from_rbac
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."refresh_casbin_rule_from_rbac"();
CREATE OR REPLACE FUNCTION "public"."refresh_casbin_rule_from_rbac"()
  RETURNS "pg_catalog"."void" AS $BODY$
BEGIN
  -- casbin_rule 是生成表/缓存表，不作为业务权限源。
  TRUNCATE TABLE public.casbin_rule RESTART IDENTITY;

  -- p: role -> api path/method
  INSERT INTO public.casbin_rule (ptype, v0, v1, v2)
  SELECT DISTINCT
    'p' AS ptype,
    r.code AS v0,
    pa.path_pattern AS v1,
    UPPER(pa.http_method) AS v2
  FROM public.sys_role r
  JOIN public.sys_role_permission rp ON rp.role_id = r.id
  JOIN public.sys_permission p ON p.id = rp.permission_id
  JOIN public.sys_permission_api pa ON pa.permission_code = p.code
  WHERE r.status = 1
    AND r.deleted_at IS NULL
    AND p.status = 1
    AND p.deleted_at IS NULL;

  -- g: user:{id} -> role_code
  -- 后端 Enforce 时 sub 应传 user:{用户ID}，例如 user:12。
  INSERT INTO public.casbin_rule (ptype, v0, v1)
  SELECT DISTINCT
    'g' AS ptype,
    'user:' || u.id::text AS v0,
    r.code AS v1
  FROM public.sys_user u
  JOIN public.sys_user_role ur ON ur.user_id = u.id
  JOIN public.sys_role r ON r.id = ur.role_id
  WHERE u.status = 1
    AND u.deleted_at IS NULL
    AND r.status = 1
    AND r.deleted_at IS NULL;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
COMMENT ON FUNCTION "public"."refresh_casbin_rule_from_rbac"() IS '从 sys_role/sys_user_role/sys_permission/sys_role_permission/sys_permission_api 生成 Casbin p/g 策略。权限变更或用户角色变更后调用。';

-- ----------------------------
-- Function structure for update_modified_column
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."update_modified_column"();
CREATE OR REPLACE FUNCTION "public"."update_modified_column"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

-- ----------------------------
-- View structure for v_user_with_roles
-- ----------------------------
DROP VIEW IF EXISTS "public"."v_user_with_roles";
CREATE VIEW "public"."v_user_with_roles" AS  SELECT u.id,
    u.username,
    u.real_name,
    u.phone,
    u.email,
    u.avatar_url,
    u.status,
    u.last_login_at,
    u.last_login_ip,
    u.created_at,
    u.updated_at,
    COALESCE(jsonb_agg(jsonb_build_object('id', r.id, 'code', r.code, 'name', r.name) ORDER BY r.id) FILTER (WHERE r.id IS NOT NULL), '[]'::jsonb) AS roles
   FROM sys_user u
     LEFT JOIN sys_user_role ur ON ur.user_id = u.id
     LEFT JOIN sys_role r ON r.id = ur.role_id AND r.deleted_at IS NULL
  WHERE u.deleted_at IS NULL
  GROUP BY u.id,
    u.username,
    u.real_name,
    u.phone,
    u.email,
    u.avatar_url,
    u.status,
    u.last_login_at,
    u.last_login_ip,
    u.created_at,
    u.updated_at;

-- ----------------------------
-- View structure for v_sales_order
-- ----------------------------
DROP VIEW IF EXISTS "public"."v_sales_order";
CREATE VIEW "public"."v_sales_order" AS  SELECT so.id,
    so.order_no,
    so.cashier_id,
    so.total_amount,
    so.medicare_amount,
    so.personal_amount,
    so.need_audit,
    so.need_medicare,
    so.status,
    so.medicare_transaction_id,
    so.created_at,
    so.updated_at,
    so.deleted_at,
    so.customer_name,
    so.is_prescription,
    so.discount_amount,
    so.actual_amount,
    so.payment_method,
    so.paid_at,
    so.cancelled_at,
    so.refunded_at,
    so.refund_amount,
    so.refund_reason,
    so.remark,
    u.real_name AS cashier_name
   FROM sales_order so
     LEFT JOIN sys_user u ON u.id = so.cashier_id
  WHERE so.deleted_at IS NULL;

-- ----------------------------
-- View structure for v_pharmacist_review
-- ----------------------------
DROP VIEW IF EXISTS "public"."v_pharmacist_review";
CREATE VIEW "public"."v_pharmacist_review" AS  SELECT ar.id,
    ar.review_no,
    ar.order_id AS sales_order_id,
    so.order_no,
    so.customer_name,
    ar.submitter_id,
    su.real_name AS submitter_name,
    ar.pharmacist_id,
    pu.real_name AS pharmacist_name,
    ar.status,
    COALESCE(ar.review_opinion, ar.comment) AS review_opinion,
    ar.submitted_at,
    ar.reviewed_at,
    ar.created_at,
    ar.updated_at
   FROM audit_review ar
     JOIN sales_order so ON so.id = ar.order_id
     LEFT JOIN sys_user su ON su.id = ar.submitter_id
     LEFT JOIN sys_user pu ON pu.id = ar.pharmacist_id
  WHERE ar.deleted_at IS NULL;

-- ----------------------------
-- View structure for v_operation_log
-- ----------------------------
DROP VIEW IF EXISTS "public"."v_operation_log";
CREATE VIEW "public"."v_operation_log" AS  SELECT ol.id,
    ol.business_type,
    ol.business_id,
    ol.action,
    ol.operator_id,
    ol.detail,
    ol.created_at,
    ol.updated_at,
    ol.deleted_at,
    ol.module,
    ol.resource_type,
    ol.resource_id,
    ol.before_data,
    ol.after_data,
    ol.ip,
    ol.user_agent,
    ol.request_id,
    u.real_name AS operator_name
   FROM operation_log ol
     LEFT JOIN sys_user u ON u.id = ol.operator_id
  WHERE ol.deleted_at IS NULL;

-- ----------------------------
-- View structure for v_inbound_order
-- ----------------------------
DROP VIEW IF EXISTS "public"."v_inbound_order";
CREATE VIEW "public"."v_inbound_order" AS  SELECT io.id,
    io.order_no,
    io.invoice_no,
    io.operator_id,
    io.status,
    io.remark,
    io.created_at,
    io.updated_at,
    io.deleted_at,
    io.supplier_id,
    io.creator_id,
    io.total_amount,
    io.submitted_at,
    io.completed_at,
    io.cancelled_at,
    COALESCE(io.creator_id, io.operator_id) AS api_creator_id,
    u.real_name AS creator_name,
    s.name AS supplier_name
   FROM inbound_order io
     LEFT JOIN sys_user u ON u.id = COALESCE(io.creator_id, io.operator_id)
     LEFT JOIN supplier s ON s.id = io.supplier_id
  WHERE io.deleted_at IS NULL;
COMMENT ON VIEW "public"."v_inbound_order" IS '入库单查询视图；supplier_name 仅由 supplier_id 关联供应商主数据得到';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_invoice_record_id_seq"
OWNED BY "public"."ai_invoice_record"."id";
SELECT setval('"public"."ai_invoice_record_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."audit_event_id_seq"
OWNED BY "public"."audit_event"."id";
SELECT setval('"public"."audit_event_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."audit_review_id_seq"
OWNED BY "public"."audit_review"."id";
SELECT setval('"public"."audit_review_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."casbin_rule_id_seq"
OWNED BY "public"."casbin_rule"."id";
SELECT setval('"public"."casbin_rule_id_seq"', 468, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."data_change_log_id_seq"
OWNED BY "public"."data_change_log"."id";
SELECT setval('"public"."data_change_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."drug_info_id_seq"
OWNED BY "public"."drug_info"."id";
SELECT setval('"public"."drug_info_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."drug_trace_inventory_id_seq"
OWNED BY "public"."drug_trace_inventory"."id";
SELECT setval('"public"."drug_trace_inventory_id_seq"', 267, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."drug_trace_log_id_seq"
OWNED BY "public"."drug_trace_log"."id";
SELECT setval('"public"."drug_trace_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."file_info_id_seq"
OWNED BY "public"."file_info"."id";
SELECT setval('"public"."file_info_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."inbound_order_detail_id_seq"
OWNED BY "public"."inbound_order_detail"."id";
SELECT setval('"public"."inbound_order_detail_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."inbound_order_id_seq"
OWNED BY "public"."inbound_order"."id";
SELECT setval('"public"."inbound_order_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."inventory_adjustment_id_seq"
OWNED BY "public"."inventory_adjustment"."id";
SELECT setval('"public"."inventory_adjustment_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."inventory_task_detail_id_seq"
OWNED BY "public"."inventory_task_detail"."id";
SELECT setval('"public"."inventory_task_detail_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."inventory_task_id_seq"
OWNED BY "public"."inventory_task"."id";
SELECT setval('"public"."inventory_task_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."location_info_id_seq"
OWNED BY "public"."location_info"."id";
SELECT setval('"public"."location_info_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."login_log_id_seq"
OWNED BY "public"."login_log"."id";
SELECT setval('"public"."login_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."notification_id_seq"
OWNED BY "public"."notification"."id";
SELECT setval('"public"."notification_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."operation_log_id_seq"
OWNED BY "public"."operation_log"."id";
SELECT setval('"public"."operation_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."report_export_task_id_seq"
OWNED BY "public"."report_export_task"."id";
SELECT setval('"public"."report_export_task_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sales_order_id_seq"
OWNED BY "public"."sales_order"."id";
SELECT setval('"public"."sales_order_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sales_order_item_id_seq"
OWNED BY "public"."sales_order_item"."id";
SELECT setval('"public"."sales_order_item_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."scan_task_detail_id_seq"
OWNED BY "public"."scan_task_detail"."id";
SELECT setval('"public"."scan_task_detail_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."scan_task_id_seq"
OWNED BY "public"."scan_task"."id";
SELECT setval('"public"."scan_task_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."security_event_id_seq"
OWNED BY "public"."security_event"."id";
SELECT setval('"public"."security_event_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."supplier_id_seq"
OWNED BY "public"."supplier"."id";
SELECT setval('"public"."supplier_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_permission_api_id_seq"
OWNED BY "public"."sys_permission_api"."id";
SELECT setval('"public"."sys_permission_api_id_seq"', 162, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_permission_id_seq"
OWNED BY "public"."sys_permission"."id";
SELECT setval('"public"."sys_permission_id_seq"', 96, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_role_id_seq"
OWNED BY "public"."sys_role"."id";
SELECT setval('"public"."sys_role_id_seq"', 6, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_user_id_seq"
OWNED BY "public"."sys_user"."id";
SELECT setval('"public"."sys_user_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."trace_reservation_id_seq"
OWNED BY "public"."trace_reservation"."id";
SELECT setval('"public"."trace_reservation_id_seq"', 1, false);

-- ----------------------------
-- Indexes structure for table ai_invoice_record
-- ----------------------------
CREATE INDEX "idx_ai_invoice_record_created_at" ON "public"."ai_invoice_record" USING btree (
  "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_invoice_record_creator" ON "public"."ai_invoice_record" USING btree (
  "creator_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_invoice_record_inbound_order" ON "public"."ai_invoice_record" USING btree (
  "inbound_order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_invoice_record_status" ON "public"."ai_invoice_record" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table ai_invoice_record
-- ----------------------------
CREATE TRIGGER "update_ai_invoice_record_modtime" BEFORE UPDATE ON "public"."ai_invoice_record"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Checks structure for table ai_invoice_record
-- ----------------------------
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "chk_ai_invoice_record_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'PROCESSING'::character varying, 'COMPLETED'::character varying, 'FAILED'::character varying]::text[])) NOT VALID;

-- ----------------------------
-- Primary Key structure for table ai_invoice_record
-- ----------------------------
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "ai_invoice_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table audit_event
-- ----------------------------
CREATE INDEX "idx_audit_event_status_time" ON "public"."audit_event" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table audit_event
-- ----------------------------
CREATE TRIGGER "update_audit_event_modtime" BEFORE UPDATE ON "public"."audit_event"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Primary Key structure for table audit_event
-- ----------------------------
ALTER TABLE "public"."audit_event" ADD CONSTRAINT "audit_event_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table audit_review
-- ----------------------------
CREATE INDEX "idx_audit_review_order" ON "public"."audit_review" USING btree (
  "order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_audit_review_pending_order" ON "public"."audit_review" USING btree (
  "order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL AND status::text = 'PENDING'::text;
COMMENT ON INDEX "public"."uk_audit_review_pending_order" IS '同一销售单同时只能存在一条有效待审核记录';
CREATE UNIQUE INDEX "uk_audit_review_review_no" ON "public"."audit_review" USING btree (
  "review_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE review_no IS NOT NULL;

-- ----------------------------
-- Triggers structure for table audit_review
-- ----------------------------
CREATE TRIGGER "update_audit_review_modtime" BEFORE UPDATE ON "public"."audit_review"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Checks structure for table audit_review
-- ----------------------------
ALTER TABLE "public"."audit_review" ADD CONSTRAINT "chk_audit_review_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'APPROVED'::character varying, 'REJECTED'::character varying, 'CANCELLED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table audit_review
-- ----------------------------
ALTER TABLE "public"."audit_review" ADD CONSTRAINT "audit_review_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table casbin_rule
-- ----------------------------
CREATE UNIQUE INDEX "uk_casbin_rule" ON "public"."casbin_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_rule
-- ----------------------------
ALTER TABLE "public"."casbin_rule" ADD CONSTRAINT "casbin_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table data_change_log
-- ----------------------------
CREATE INDEX "idx_data_change_log_operator_time" ON "public"."data_change_log" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_data_change_log_record" ON "public"."data_change_log" USING btree (
  "table_name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "record_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table data_change_log
-- ----------------------------
ALTER TABLE "public"."data_change_log" ADD CONSTRAINT "data_change_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table drug_info
-- ----------------------------
CREATE INDEX "idx_drug_info_status" ON "public"."drug_info" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX "uk_drug_info_barcode" ON "public"."drug_info" USING btree (
  "barcode" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE barcode IS NOT NULL;

-- ----------------------------
-- Triggers structure for table drug_info
-- ----------------------------
CREATE TRIGGER "update_drug_info_modtime" BEFORE UPDATE ON "public"."drug_info"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table drug_info
-- ----------------------------
ALTER TABLE "public"."drug_info" ADD CONSTRAINT "drug_info_drug_code_key" UNIQUE ("drug_code");

-- ----------------------------
-- Primary Key structure for table drug_info
-- ----------------------------
ALTER TABLE "public"."drug_info" ADD CONSTRAINT "drug_info_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table drug_trace_inventory
-- ----------------------------
CREATE INDEX "idx_trace_inventory_drug_expire_status" ON "public"."drug_trace_inventory" USING btree (
  "drug_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "expire_date" "pg_catalog"."date_ops" ASC NULLS LAST,
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_inventory_expire" ON "public"."drug_trace_inventory" USING btree (
  "expire_date" "pg_catalog"."date_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_inventory_location" ON "public"."drug_trace_inventory" USING btree (
  "location_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_inventory_status" ON "public"."drug_trace_inventory" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table drug_trace_inventory
-- ----------------------------
CREATE TRIGGER "update_drug_trace_inventory_modtime" BEFORE UPDATE ON "public"."drug_trace_inventory"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table drug_trace_inventory
-- ----------------------------
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "drug_trace_inventory_trace_code_key" UNIQUE ("trace_code");

-- ----------------------------
-- Checks structure for table drug_trace_inventory
-- ----------------------------
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "chk_trace_inventory_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'IN_STOCK'::character varying, 'SOLD'::character varying, 'MISPLACED'::character varying, 'LOSS_CANDIDATE'::character varying, 'LOST'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table drug_trace_inventory
-- ----------------------------
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "drug_trace_inventory_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table drug_trace_log
-- ----------------------------
CREATE INDEX "idx_trace_log_from_location" ON "public"."drug_trace_log" USING btree (
  "from_location_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_log_operator_time" ON "public"."drug_trace_log" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_trace_log_to_location" ON "public"."drug_trace_log" USING btree (
  "to_location_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_log_trace_code" ON "public"."drug_trace_log" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table drug_trace_log
-- ----------------------------
CREATE TRIGGER "update_drug_trace_log_modtime" BEFORE UPDATE ON "public"."drug_trace_log"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Primary Key structure for table drug_trace_log
-- ----------------------------
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "drug_trace_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table file_info
-- ----------------------------
CREATE INDEX "idx_file_info_business" ON "public"."file_info" USING btree (
  "business_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "business_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_file_info_uploader" ON "public"."file_info" USING btree (
  "uploader_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table file_info
-- ----------------------------
CREATE TRIGGER "update_file_info_modtime" BEFORE UPDATE ON "public"."file_info"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table file_info
-- ----------------------------
ALTER TABLE "public"."file_info" ADD CONSTRAINT "file_info_file_id_key" UNIQUE ("file_id");

-- ----------------------------
-- Primary Key structure for table file_info
-- ----------------------------
ALTER TABLE "public"."file_info" ADD CONSTRAINT "file_info_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table inbound_order
-- ----------------------------
CREATE INDEX "idx_inbound_order_status_time" ON "public"."inbound_order" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_inbound_order_supplier" ON "public"."inbound_order" USING btree (
  "supplier_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table inbound_order
-- ----------------------------
CREATE TRIGGER "update_inbound_order_modtime" BEFORE UPDATE ON "public"."inbound_order"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table inbound_order
-- ----------------------------
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "inbound_order_order_no_key" UNIQUE ("order_no");

-- ----------------------------
-- Checks structure for table inbound_order
-- ----------------------------
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "chk_inbound_order_status" CHECK (status::text = ANY (ARRAY['DRAFT'::character varying, 'PENDING_CONFIRM'::character varying, 'COMPLETED'::character varying, 'CANCELLED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table inbound_order
-- ----------------------------
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "inbound_order_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table inbound_order_detail
-- ----------------------------
CREATE INDEX "idx_inbound_detail_drug" ON "public"."inbound_order_detail" USING btree (
  "drug_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_inbound_detail_order_id" ON "public"."inbound_order_detail" USING btree (
  "order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table inbound_order_detail
-- ----------------------------
CREATE TRIGGER "update_inbound_order_detail_modtime" BEFORE UPDATE ON "public"."inbound_order_detail"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Checks structure for table inbound_order_detail
-- ----------------------------
ALTER TABLE "public"."inbound_order_detail" ADD CONSTRAINT "chk_inbound_order_detail_qty" CHECK (planned_qty > 0 AND COALESCE(confirmed_qty, 0) >= 0 AND COALESCE(confirmed_qty, 0) <= planned_qty) NOT VALID;
ALTER TABLE "public"."inbound_order_detail" ADD CONSTRAINT "chk_inbound_order_detail_amount" CHECK ((unit_price IS NULL OR unit_price >= 0::numeric) AND (amount IS NULL OR amount >= 0::numeric)) NOT VALID;
COMMENT ON CONSTRAINT "chk_inbound_order_detail_qty" ON "public"."inbound_order_detail" IS '入库计划数量必须大于 0，确认数量不能为负且不能超过计划数量';
COMMENT ON CONSTRAINT "chk_inbound_order_detail_amount" ON "public"."inbound_order_detail" IS '入库单价和金额不得为负';

-- ----------------------------
-- Primary Key structure for table inbound_order_detail
-- ----------------------------
ALTER TABLE "public"."inbound_order_detail" ADD CONSTRAINT "inbound_order_detail_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table inventory_adjustment
-- ----------------------------
CREATE INDEX "idx_inventory_adjustment_operator_time" ON "public"."inventory_adjustment" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_inventory_adjustment_trace" ON "public"."inventory_adjustment" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table inventory_adjustment
-- ----------------------------
CREATE TRIGGER "update_inventory_adjustment_modtime" BEFORE UPDATE ON "public"."inventory_adjustment"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table inventory_adjustment
-- ----------------------------
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "inventory_adjustment_adjust_no_key" UNIQUE ("adjust_no");

-- ----------------------------
-- Checks structure for table inventory_adjustment
-- ----------------------------
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "chk_inventory_adjustment_type" CHECK (adjust_type::text = ANY (ARRAY['LOSS'::character varying, 'GAIN'::character varying, 'RELOCATE'::character varying, 'DAMAGE'::character varying, 'EXPIRE'::character varying, 'MANUAL'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table inventory_adjustment
-- ----------------------------
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "inventory_adjustment_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Triggers structure for table inventory_task
-- ----------------------------
CREATE TRIGGER "update_inventory_task_modtime" BEFORE UPDATE ON "public"."inventory_task"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table inventory_task
-- ----------------------------
ALTER TABLE "public"."inventory_task" ADD CONSTRAINT "inventory_task_task_no_key" UNIQUE ("task_no");

-- ----------------------------
-- Checks structure for table inventory_task
-- ----------------------------
ALTER TABLE "public"."inventory_task" ADD CONSTRAINT "chk_inventory_task_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'IN_PROGRESS'::character varying, 'COMPLETED'::character varying, 'CANCELLED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table inventory_task
-- ----------------------------
ALTER TABLE "public"."inventory_task" ADD CONSTRAINT "inventory_task_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table inventory_task_detail
-- ----------------------------
CREATE INDEX "idx_inventory_detail_task" ON "public"."inventory_task_detail" USING btree (
  "task_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_inventory_task_detail_operator" ON "public"."inventory_task_detail" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_inventory_task_detail_scanned_at" ON "public"."inventory_task_detail" USING btree (
  "scanned_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_inventory_task_detail_scanned_location" ON "public"."inventory_task_detail" USING btree (
  "scanned_location_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_inventory_task_detail_system_location" ON "public"."inventory_task_detail" USING btree (
  "system_location_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table inventory_task_detail
-- ----------------------------
CREATE TRIGGER "update_inventory_task_detail_modtime" BEFORE UPDATE ON "public"."inventory_task_detail"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Checks structure for table inventory_task_detail
-- ----------------------------
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "chk_inventory_task_detail_discrepancy" CHECK (discrepancy_type::text = ANY (ARRAY['NORMAL'::character varying, 'MISPLACED_FOUND'::character varying, 'UNEXPECTED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table inventory_task_detail
-- ----------------------------
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "inventory_task_detail_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table location_info
-- ----------------------------
CREATE INDEX "idx_location_info_status" ON "public"."location_info" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Triggers structure for table location_info
-- ----------------------------
CREATE TRIGGER "update_location_info_modtime" BEFORE UPDATE ON "public"."location_info"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table location_info
-- ----------------------------
ALTER TABLE "public"."location_info" ADD CONSTRAINT "location_info_location_code_key" UNIQUE ("location_code");

-- ----------------------------
-- Primary Key structure for table location_info
-- ----------------------------
ALTER TABLE "public"."location_info" ADD CONSTRAINT "location_info_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table login_log
-- ----------------------------
CREATE INDEX "idx_login_log_ip_time" ON "public"."login_log" USING btree (
  "ip" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_login_log_user_time" ON "public"."login_log" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table login_log
-- ----------------------------
ALTER TABLE "public"."login_log" ADD CONSTRAINT "login_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table notification
-- ----------------------------
CREATE INDEX "idx_notification_business" ON "public"."notification" USING btree (
  "business_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "business_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_notification_user_read" ON "public"."notification" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "read_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table notification
-- ----------------------------
ALTER TABLE "public"."notification" ADD CONSTRAINT "notification_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table operation_log
-- ----------------------------
CREATE INDEX "idx_op_log_business" ON "public"."operation_log" USING btree (
  "business_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "business_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_op_log_resource" ON "public"."operation_log" USING btree (
  "resource_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "resource_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_operation_log_operator_time" ON "public"."operation_log" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table operation_log
-- ----------------------------
CREATE TRIGGER "update_operation_log_modtime" BEFORE UPDATE ON "public"."operation_log"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Primary Key structure for table operation_log
-- ----------------------------
ALTER TABLE "public"."operation_log" ADD CONSTRAINT "operation_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table report_export_task
-- ----------------------------
CREATE INDEX "idx_report_export_task_user_time" ON "public"."report_export_task" USING btree (
  "requested_by" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table report_export_task
-- ----------------------------
CREATE TRIGGER "update_report_export_task_modtime" BEFORE UPDATE ON "public"."report_export_task"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table report_export_task
-- ----------------------------
ALTER TABLE "public"."report_export_task" ADD CONSTRAINT "report_export_task_task_id_key" UNIQUE ("task_id");

-- ----------------------------
-- Checks structure for table report_export_task
-- ----------------------------
ALTER TABLE "public"."report_export_task" ADD CONSTRAINT "chk_report_export_task_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'RUNNING'::character varying, 'SUCCESS'::character varying, 'FAILED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table report_export_task
-- ----------------------------
ALTER TABLE "public"."report_export_task" ADD CONSTRAINT "report_export_task_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sales_order
-- ----------------------------
CREATE INDEX "idx_sales_order_cashier_time" ON "public"."sales_order" USING btree (
  "cashier_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_sales_order_deleted" ON "public"."sales_order" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sales_order_status_time" ON "public"."sales_order" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table sales_order
-- ----------------------------
CREATE TRIGGER "update_sales_order_modtime" BEFORE UPDATE ON "public"."sales_order"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table sales_order
-- ----------------------------
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "sales_order_order_no_key" UNIQUE ("order_no");

-- ----------------------------
-- Checks structure for table sales_order
-- ----------------------------
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "chk_sales_order_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'PENDING_REVIEW'::character varying, 'APPROVED'::character varying, 'COMPLETED'::character varying, 'PARTIALLY_REFUNDED'::character varying, 'REFUNDED'::character varying, 'CANCELLED'::character varying]::text[]));
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "chk_sales_order_prescription_need_audit" CHECK (is_prescription = false OR need_audit = true) NOT VALID;
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "chk_sales_order_amount_range" CHECK (total_amount >= 0::numeric AND COALESCE(actual_amount, 0::numeric) >= 0::numeric AND COALESCE(discount_amount, 0::numeric) >= 0::numeric AND COALESCE(refund_amount, 0::numeric) >= 0::numeric AND COALESCE(refund_amount, 0::numeric) <= COALESCE(actual_amount, total_amount)) NOT VALID;

-- ----------------------------
-- Primary Key structure for table sales_order
-- ----------------------------
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "sales_order_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sales_order_item
-- ----------------------------
CREATE INDEX "idx_sales_item_drug" ON "public"."sales_order_item" USING btree (
  "drug_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sales_item_order" ON "public"."sales_order_item" USING btree (
  "order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sales_item_trace" ON "public"."sales_order_item" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sales_order_item_refund_status" ON "public"."sales_order_item" USING btree (
  "refund_status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sales_order_item_refunded_at" ON "public"."sales_order_item" USING btree (
  "refunded_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_sales_order_item_active_trace" ON "public"."sales_order_item" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL AND refund_status::text = 'NONE'::text;
COMMENT ON INDEX "public"."uk_sales_order_item_active_trace" IS '同一 trace_code 同时只能存在一条有效未退销售明细，避免一盒药重复销售';

-- ----------------------------
-- Triggers structure for table sales_order_item
-- ----------------------------
CREATE TRIGGER "update_sales_order_item_modtime" BEFORE UPDATE ON "public"."sales_order_item"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Checks structure for table sales_order_item
-- ----------------------------
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "chk_sales_order_item_quantity_one" CHECK (quantity = 1);
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "chk_sales_order_item_price_nonnegative" CHECK (price >= 0::numeric);
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "chk_sales_order_item_subtotal_nonnegative" CHECK (subtotal_amount IS NULL OR subtotal_amount >= 0::numeric);
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "chk_sales_order_item_refund_status" CHECK (refund_status::text = ANY (ARRAY['NONE'::character varying, 'REFUNDED'::character varying]::text[]));
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "chk_sales_order_item_refund_amount_range" CHECK (refund_amount >= 0::numeric AND refund_amount <= price) NOT VALID;

-- ----------------------------
-- Primary Key structure for table sales_order_item
-- ----------------------------
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "sales_order_item_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Triggers structure for table scan_task
-- ----------------------------
CREATE TRIGGER "update_scan_task_modtime" BEFORE UPDATE ON "public"."scan_task"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table scan_task
-- ----------------------------
ALTER TABLE "public"."scan_task" ADD CONSTRAINT "scan_task_task_no_key" UNIQUE ("task_no");

-- ----------------------------
-- Checks structure for table scan_task
-- ----------------------------
ALTER TABLE "public"."scan_task" ADD CONSTRAINT "chk_scan_task_type" CHECK (task_type::text = ANY (ARRAY['INBOUND'::character varying, 'SHELVING'::character varying, 'INVENTORY'::character varying]::text[])) NOT VALID;
ALTER TABLE "public"."scan_task" ADD CONSTRAINT "chk_scan_task_status" CHECK (status::text = ANY (ARRAY['PENDING'::character varying, 'IN_PROGRESS'::character varying, 'COMPLETED'::character varying, 'CANCELLED'::character varying]::text[])) NOT VALID;
COMMENT ON CONSTRAINT "chk_scan_task_type" ON "public"."scan_task" IS '扫码任务类型：INBOUND 入库扫码，SHELVING 上架扫码，INVENTORY 盘库扫码';
COMMENT ON CONSTRAINT "chk_scan_task_status" ON "public"."scan_task" IS '扫码任务状态：PENDING、IN_PROGRESS、COMPLETED、CANCELLED';

-- ----------------------------
-- Primary Key structure for table scan_task
-- ----------------------------
ALTER TABLE "public"."scan_task" ADD CONSTRAINT "scan_task_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Triggers structure for table scan_task_detail
-- ----------------------------
CREATE TRIGGER "update_scan_task_detail_modtime" BEFORE UPDATE ON "public"."scan_task_detail"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table scan_task_detail
-- ----------------------------
ALTER TABLE "public"."scan_task_detail" ADD CONSTRAINT "uk_scan_task_trace" UNIQUE ("task_id", "trace_code");

-- ----------------------------
-- Checks structure for table scan_task_detail
-- ----------------------------
ALTER TABLE "public"."scan_task_detail" ADD CONSTRAINT "chk_scan_task_detail_result" CHECK (scan_result::text = ANY (ARRAY['SUCCESS'::character varying, 'DUPLICATE'::character varying, 'INVALID'::character varying, 'STATUS_ERROR'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table scan_task_detail
-- ----------------------------
ALTER TABLE "public"."scan_task_detail" ADD CONSTRAINT "scan_task_detail_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table security_event
-- ----------------------------
CREATE INDEX "idx_security_event_type_time" ON "public"."security_event" USING btree (
  "event_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table security_event
-- ----------------------------
ALTER TABLE "public"."security_event" ADD CONSTRAINT "security_event_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table supplier
-- ----------------------------
CREATE INDEX "idx_supplier_status" ON "public"."supplier" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Triggers structure for table supplier
-- ----------------------------
CREATE TRIGGER "update_supplier_modtime" BEFORE UPDATE ON "public"."supplier"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table supplier
-- ----------------------------
ALTER TABLE "public"."supplier" ADD CONSTRAINT "supplier_supplier_code_key" UNIQUE ("supplier_code");

-- ----------------------------
-- Primary Key structure for table supplier
-- ----------------------------
ALTER TABLE "public"."supplier" ADD CONSTRAINT "supplier_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_permission
-- ----------------------------
CREATE INDEX "idx_sys_permission_resource_action" ON "public"."sys_permission" USING btree (
  "resource" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "action" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table sys_permission
-- ----------------------------
CREATE TRIGGER "update_sys_permission_modtime" BEFORE UPDATE ON "public"."sys_permission"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table sys_permission
-- ----------------------------
ALTER TABLE "public"."sys_permission" ADD CONSTRAINT "sys_permission_code_key" UNIQUE ("code");

-- ----------------------------
-- Primary Key structure for table sys_permission
-- ----------------------------
ALTER TABLE "public"."sys_permission" ADD CONSTRAINT "sys_permission_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_permission_api
-- ----------------------------
CREATE INDEX "idx_sys_permission_api_path_method" ON "public"."sys_permission_api" USING btree (
  "path_pattern" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "http_method" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table sys_permission_api
-- ----------------------------
ALTER TABLE "public"."sys_permission_api" ADD CONSTRAINT "uk_sys_permission_api" UNIQUE ("permission_code", "path_pattern", "http_method");

-- ----------------------------
-- Primary Key structure for table sys_permission_api
-- ----------------------------
ALTER TABLE "public"."sys_permission_api" ADD CONSTRAINT "sys_permission_api_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Triggers structure for table sys_role
-- ----------------------------
CREATE TRIGGER "update_sys_role_modtime" BEFORE UPDATE ON "public"."sys_role"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table sys_role
-- ----------------------------
ALTER TABLE "public"."sys_role" ADD CONSTRAINT "sys_role_code_key" UNIQUE ("code");

-- ----------------------------
-- Primary Key structure for table sys_role
-- ----------------------------
ALTER TABLE "public"."sys_role" ADD CONSTRAINT "sys_role_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_role_permission
-- ----------------------------
CREATE INDEX "idx_sys_role_permission_permission" ON "public"."sys_role_permission" USING btree (
  "permission_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_role_permission
-- ----------------------------
ALTER TABLE "public"."sys_role_permission" ADD CONSTRAINT "sys_role_permission_pkey" PRIMARY KEY ("role_id", "permission_id");

-- ----------------------------
-- Indexes structure for table sys_user
-- ----------------------------
CREATE INDEX "idx_sys_user_status" ON "public"."sys_user" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Triggers structure for table sys_user
-- ----------------------------
CREATE TRIGGER "update_sys_user_modtime" BEFORE UPDATE ON "public"."sys_user"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table sys_user
-- ----------------------------
ALTER TABLE "public"."sys_user" ADD CONSTRAINT "sys_user_username_key" UNIQUE ("username");

-- ----------------------------
-- Primary Key structure for table sys_user
-- ----------------------------
ALTER TABLE "public"."sys_user" ADD CONSTRAINT "sys_user_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_user_role
-- ----------------------------
CREATE INDEX "idx_sys_user_role_role" ON "public"."sys_user_role" USING btree (
  "role_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_user_role
-- ----------------------------
ALTER TABLE "public"."sys_user_role" ADD CONSTRAINT "sys_user_role_pkey" PRIMARY KEY ("user_id", "role_id");

-- ----------------------------
-- Indexes structure for table trace_reservation
-- ----------------------------
CREATE INDEX "idx_trace_reservation_order" ON "public"."trace_reservation" USING btree (
  "sales_order_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_trace_reservation_trace" ON "public"."trace_reservation" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_trace_reservation_active_trace" ON "public"."trace_reservation" USING btree (
  "trace_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE status::text = 'RESERVED'::text AND deleted_at IS NULL;

-- ----------------------------
-- Triggers structure for table trace_reservation
-- ----------------------------
CREATE TRIGGER "update_trace_reservation_modtime" BEFORE UPDATE ON "public"."trace_reservation"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table trace_reservation
-- ----------------------------
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "trace_reservation_reservation_no_key" UNIQUE ("reservation_no");

-- ----------------------------
-- Checks structure for table trace_reservation
-- ----------------------------
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "chk_trace_reservation_status" CHECK (status::text = ANY (ARRAY['RESERVED'::character varying, 'RELEASED'::character varying, 'CONSUMED'::character varying, 'EXPIRED'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table trace_reservation
-- ----------------------------
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "trace_reservation_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table ai_invoice_record
-- ----------------------------
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "fk_ai_invoice_record_creator" FOREIGN KEY ("creator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "fk_ai_invoice_record_file" FOREIGN KEY ("file_id") REFERENCES "public"."file_info" ("file_id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "fk_ai_invoice_record_inbound_order" FOREIGN KEY ("inbound_order_id") REFERENCES "public"."inbound_order" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."ai_invoice_record" ADD CONSTRAINT "fk_ai_invoice_record_matched_supplier" FOREIGN KEY ("matched_supplier_id") REFERENCES "public"."supplier" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table audit_event
-- ----------------------------
ALTER TABLE "public"."audit_event" ADD CONSTRAINT "fk_audit_event_assigned_to" FOREIGN KEY ("assigned_to") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."audit_event" ADD CONSTRAINT "fk_audit_event_ignored_by" FOREIGN KEY ("ignored_by") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."audit_event" ADD CONSTRAINT "fk_audit_event_resolved_by" FOREIGN KEY ("resolved_by") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table audit_review
-- ----------------------------
ALTER TABLE "public"."audit_review" ADD CONSTRAINT "fk_audit_review_order" FOREIGN KEY ("order_id") REFERENCES "public"."sales_order" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."audit_review" ADD CONSTRAINT "fk_audit_review_pharmacist" FOREIGN KEY ("pharmacist_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."audit_review" ADD CONSTRAINT "fk_audit_review_submitter" FOREIGN KEY ("submitter_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table data_change_log
-- ----------------------------
ALTER TABLE "public"."data_change_log" ADD CONSTRAINT "fk_data_change_log_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table drug_trace_inventory
-- ----------------------------
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "fk_trace_inventory_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "fk_trace_inventory_inbound_detail" FOREIGN KEY ("inbound_detail_id") REFERENCES "public"."inbound_order_detail" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "fk_trace_inventory_inbound_order" FOREIGN KEY ("inbound_order_id") REFERENCES "public"."inbound_order" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_inventory" ADD CONSTRAINT "fk_trace_inventory_location" FOREIGN KEY ("location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table drug_trace_log
-- ----------------------------
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_from_location" FOREIGN KEY ("from_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_order" FOREIGN KEY ("order_id") REFERENCES "public"."sales_order" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_order_item" FOREIGN KEY ("order_item_id") REFERENCES "public"."sales_order_item" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."drug_trace_log" ADD CONSTRAINT "fk_trace_log_to_location" FOREIGN KEY ("to_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table file_info
-- ----------------------------
ALTER TABLE "public"."file_info" ADD CONSTRAINT "fk_file_info_uploader" FOREIGN KEY ("uploader_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table inbound_order
-- ----------------------------
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "fk_inbound_order_creator" FOREIGN KEY ("creator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "fk_inbound_order_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inbound_order" ADD CONSTRAINT "fk_inbound_order_supplier" FOREIGN KEY ("supplier_id") REFERENCES "public"."supplier" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table inbound_order_detail
-- ----------------------------
ALTER TABLE "public"."inbound_order_detail" ADD CONSTRAINT "fk_inbound_order_detail_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inbound_order_detail" ADD CONSTRAINT "fk_inbound_order_detail_order" FOREIGN KEY ("order_id") REFERENCES "public"."inbound_order" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table inventory_adjustment
-- ----------------------------
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_from_location" FOREIGN KEY ("from_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_task" FOREIGN KEY ("related_task_id") REFERENCES "public"."inventory_task" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_to_location" FOREIGN KEY ("to_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_adjustment" ADD CONSTRAINT "fk_inventory_adjustment_trace" FOREIGN KEY ("trace_code") REFERENCES "public"."drug_trace_inventory" ("trace_code") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table inventory_task
-- ----------------------------
ALTER TABLE "public"."inventory_task" ADD CONSTRAINT "fk_inventory_task_assignee" FOREIGN KEY ("assignee_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_task" ADD CONSTRAINT "fk_inventory_task_creator" FOREIGN KEY ("creator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table inventory_task_detail
-- ----------------------------
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "fk_inventory_task_detail_location" FOREIGN KEY ("location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "fk_inventory_task_detail_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "fk_inventory_task_detail_scanned_location" FOREIGN KEY ("scanned_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "fk_inventory_task_detail_system_location" FOREIGN KEY ("system_location_id") REFERENCES "public"."location_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."inventory_task_detail" ADD CONSTRAINT "fk_inventory_task_detail_task" FOREIGN KEY ("task_id") REFERENCES "public"."inventory_task" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table login_log
-- ----------------------------
ALTER TABLE "public"."login_log" ADD CONSTRAINT "fk_login_log_user" FOREIGN KEY ("user_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table notification
-- ----------------------------
ALTER TABLE "public"."notification" ADD CONSTRAINT "fk_notification_user" FOREIGN KEY ("user_id") REFERENCES "public"."sys_user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table operation_log
-- ----------------------------
ALTER TABLE "public"."operation_log" ADD CONSTRAINT "fk_operation_log_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table report_export_task
-- ----------------------------
ALTER TABLE "public"."report_export_task" ADD CONSTRAINT "fk_report_export_task_file" FOREIGN KEY ("file_id") REFERENCES "public"."file_info" ("file_id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."report_export_task" ADD CONSTRAINT "fk_report_export_task_user" FOREIGN KEY ("requested_by") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sales_order
-- ----------------------------
ALTER TABLE "public"."sales_order" ADD CONSTRAINT "fk_sales_order_cashier" FOREIGN KEY ("cashier_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sales_order_item
-- ----------------------------
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "fk_sales_order_item_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "fk_sales_order_item_order" FOREIGN KEY ("order_id") REFERENCES "public"."sales_order" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "fk_sales_order_item_refund_operator" FOREIGN KEY ("refund_operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."sales_order_item" ADD CONSTRAINT "fk_sales_order_item_trace" FOREIGN KEY ("trace_code") REFERENCES "public"."drug_trace_inventory" ("trace_code") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table scan_task
-- ----------------------------
ALTER TABLE "public"."scan_task" ADD CONSTRAINT "fk_scan_task_operator" FOREIGN KEY ("operator_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table scan_task_detail
-- ----------------------------
ALTER TABLE "public"."scan_task_detail" ADD CONSTRAINT "fk_scan_task_detail_task" FOREIGN KEY ("task_id") REFERENCES "public"."scan_task" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table security_event
-- ----------------------------
ALTER TABLE "public"."security_event" ADD CONSTRAINT "fk_security_event_handler" FOREIGN KEY ("handled_by") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."security_event" ADD CONSTRAINT "fk_security_event_user" FOREIGN KEY ("user_id") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sys_permission_api
-- ----------------------------
ALTER TABLE "public"."sys_permission_api" ADD CONSTRAINT "fk_sys_permission_api_permission" FOREIGN KEY ("permission_code") REFERENCES "public"."sys_permission" ("code") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sys_role_permission
-- ----------------------------
ALTER TABLE "public"."sys_role_permission" ADD CONSTRAINT "fk_sys_role_permission_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."sys_permission" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."sys_role_permission" ADD CONSTRAINT "fk_sys_role_permission_role" FOREIGN KEY ("role_id") REFERENCES "public"."sys_role" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sys_user_role
-- ----------------------------
ALTER TABLE "public"."sys_user_role" ADD CONSTRAINT "fk_sys_user_role_role" FOREIGN KEY ("role_id") REFERENCES "public"."sys_role" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."sys_user_role" ADD CONSTRAINT "fk_sys_user_role_user" FOREIGN KEY ("user_id") REFERENCES "public"."sys_user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table trace_reservation
-- ----------------------------
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "fk_trace_reservation_drug" FOREIGN KEY ("drug_id") REFERENCES "public"."drug_info" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "fk_trace_reservation_item" FOREIGN KEY ("sales_order_item_id") REFERENCES "public"."sales_order_item" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "fk_trace_reservation_order" FOREIGN KEY ("sales_order_id") REFERENCES "public"."sales_order" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "fk_trace_reservation_trace_code" FOREIGN KEY ("trace_code") REFERENCES "public"."drug_trace_inventory" ("trace_code") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."trace_reservation" ADD CONSTRAINT "fk_trace_reservation_user" FOREIGN KEY ("reserved_by") REFERENCES "public"."sys_user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

COMMIT;
