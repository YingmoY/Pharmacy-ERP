/*
 Navicat Premium Dump SQL

 Source Server         : PostgreSQL
 Source Server Type    : PostgreSQL
 Source Server Version : 180002 (180002)
 Source Host           : 127.0.0.1:5432
 Source Catalog        : pharmacy_erp
 Source Schema         : medicare

 Target Server Type    : PostgreSQL
 Target Server Version : 180002 (180002)
 File Encoding         : 65001

 Date: 12/06/2026 20:16:21
*/


-- ----------------------------
-- Sequence structure for audit_logs_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "medicare"."audit_logs_id_seq";
CREATE SEQUENCE "medicare"."audit_logs_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for medicare_api_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "medicare"."medicare_api_log_id_seq";
CREATE SEQUENCE "medicare"."medicare_api_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for medicare_transaction_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "medicare"."medicare_transaction_id_seq";
CREATE SEQUENCE "medicare"."medicare_transaction_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sign_sessions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "medicare"."sign_sessions_id_seq";
CREATE SEQUENCE "medicare"."sign_sessions_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for async_tasks
-- ----------------------------
DROP TABLE IF EXISTS "medicare"."async_tasks";
CREATE TABLE "medicare"."async_tasks" (
  "id" uuid NOT NULL,
  "infno" text COLLATE "pg_catalog"."default" NOT NULL,
  "operation" text COLLATE "pg_catalog"."default" NOT NULL,
  "erp_order_no" text COLLATE "pg_catalog"."default",
  "status" text COLLATE "pg_catalog"."default" NOT NULL,
  "attempts" int4 NOT NULL DEFAULT 0,
  "max_attempts" int4 NOT NULL DEFAULT 5,
  "request_payload" jsonb NOT NULL,
  "response_payload" jsonb,
  "error_message" text COLLATE "pg_catalog"."default",
  "next_run_at" timestamptz(6) NOT NULL DEFAULT now(),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for audit_logs
-- ----------------------------
DROP TABLE IF EXISTS "medicare"."audit_logs";
CREATE TABLE "medicare"."audit_logs" (
  "id" int8 NOT NULL DEFAULT nextval('"medicare".audit_logs_id_seq'::regclass),
  "trace_id" uuid NOT NULL,
  "msgid" text COLLATE "pg_catalog"."default" NOT NULL,
  "infno" text COLLATE "pg_catalog"."default" NOT NULL,
  "erp_order_no" text COLLATE "pg_catalog"."default",
  "operation" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" text COLLATE "pg_catalog"."default" NOT NULL,
  "request_payload" jsonb NOT NULL,
  "response_payload" jsonb,
  "error_message" text COLLATE "pg_catalog"."default",
  "http_status" int4,
  "elapsed_ms" int4,
  "sign_no" text COLLATE "pg_catalog"."default",
  "async_task_id" uuid,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for medicare_api_log
-- ----------------------------
DROP TABLE IF EXISTS "medicare"."medicare_api_log";
CREATE TABLE "medicare"."medicare_api_log" (
  "id" int8 NOT NULL DEFAULT nextval('"medicare".medicare_api_log_id_seq'::regclass),
  "transaction_no" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "api_endpoint" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "request_payload" text COLLATE "pg_catalog"."default",
  "response_payload" text COLLATE "pg_catalog"."default",
  "http_status" int4,
  "duration_ms" int4,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "medicare"."medicare_api_log" IS '医保外部接口调用日志表';

-- ----------------------------
-- Table structure for medicare_transaction
-- ----------------------------
DROP TABLE IF EXISTS "medicare"."medicare_transaction";
CREATE TABLE "medicare"."medicare_transaction" (
  "id" int8 NOT NULL DEFAULT nextval('"medicare".medicare_transaction_id_seq'::regclass),
  "transaction_no" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "erp_order_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "customer_id_no" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "trade_type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "total_amount" numeric(10,2) NOT NULL,
  "medicare_amount" numeric(10,2),
  "personal_amount" numeric(10,2),
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "fail_reason" varchar(255) COLLATE "pg_catalog"."default",
  "raw_response" jsonb,
  "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz(6)
)
;
COMMENT ON TABLE "medicare"."medicare_transaction" IS '医保交易流水表';

-- ----------------------------
-- Table structure for sign_sessions
-- ----------------------------
DROP TABLE IF EXISTS "medicare"."sign_sessions";
CREATE TABLE "medicare"."sign_sessions" (
  "id" int8 NOT NULL DEFAULT nextval('"medicare".sign_sessions_id_seq'::regclass),
  "sign_no" text COLLATE "pg_catalog"."default" NOT NULL,
  "operator_code" text COLLATE "pg_catalog"."default" NOT NULL,
  "operator_name" text COLLATE "pg_catalog"."default" NOT NULL,
  "request_msgid" text COLLATE "pg_catalog"."default" NOT NULL,
  "response" jsonb NOT NULL,
  "signed_at" timestamptz(6) NOT NULL DEFAULT now(),
  "expires_at" timestamptz(6),
  "active" bool NOT NULL DEFAULT true
)
;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "medicare"."audit_logs_id_seq"
OWNED BY "medicare"."audit_logs"."id";
SELECT setval('"medicare"."audit_logs_id_seq"', 25, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "medicare"."medicare_api_log_id_seq"
OWNED BY "medicare"."medicare_api_log"."id";
SELECT setval('"medicare"."medicare_api_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "medicare"."medicare_transaction_id_seq"
OWNED BY "medicare"."medicare_transaction"."id";
SELECT setval('"medicare"."medicare_transaction_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "medicare"."sign_sessions_id_seq"
OWNED BY "medicare"."sign_sessions"."id";
SELECT setval('"medicare"."sign_sessions_id_seq"', 9, true);

-- ----------------------------
-- Indexes structure for table async_tasks
-- ----------------------------
CREATE INDEX "idx_async_tasks_status_next_run" ON "medicare"."async_tasks" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "next_run_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table async_tasks
-- ----------------------------
ALTER TABLE "medicare"."async_tasks" ADD CONSTRAINT "async_tasks_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table audit_logs
-- ----------------------------
CREATE INDEX "idx_audit_logs_infno_created_at" ON "medicare"."audit_logs" USING btree (
  "infno" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_audit_logs_msgid" ON "medicare"."audit_logs" USING btree (
  "msgid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_audit_logs_order" ON "medicare"."audit_logs" USING btree (
  "erp_order_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_audit_logs_trace_id" ON "medicare"."audit_logs" USING btree (
  "trace_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table audit_logs
-- ----------------------------
ALTER TABLE "medicare"."audit_logs" ADD CONSTRAINT "audit_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table medicare_api_log
-- ----------------------------
CREATE INDEX "idx_medicare_api_log_tx" ON "medicare"."medicare_api_log" USING btree (
  "transaction_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table medicare_api_log
-- ----------------------------
CREATE TRIGGER "update_medicare_api_log_modtime" BEFORE UPDATE ON "medicare"."medicare_api_log"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Primary Key structure for table medicare_api_log
-- ----------------------------
ALTER TABLE "medicare"."medicare_api_log" ADD CONSTRAINT "medicare_api_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table medicare_transaction
-- ----------------------------
CREATE INDEX "idx_medicare_tx_erp_order" ON "medicare"."medicare_transaction" USING btree (
  "erp_order_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table medicare_transaction
-- ----------------------------
CREATE TRIGGER "update_medicare_transaction_modtime" BEFORE UPDATE ON "medicare"."medicare_transaction"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_modified_column"();

-- ----------------------------
-- Uniques structure for table medicare_transaction
-- ----------------------------
ALTER TABLE "medicare"."medicare_transaction" ADD CONSTRAINT "medicare_transaction_transaction_no_key" UNIQUE ("transaction_no");

-- ----------------------------
-- Primary Key structure for table medicare_transaction
-- ----------------------------
ALTER TABLE "medicare"."medicare_transaction" ADD CONSTRAINT "medicare_transaction_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sign_sessions
-- ----------------------------
CREATE INDEX "idx_sign_sessions_active_signed_at" ON "medicare"."sign_sessions" USING btree (
  "active" "pg_catalog"."bool_ops" ASC NULLS LAST,
  "signed_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table sign_sessions
-- ----------------------------
ALTER TABLE "medicare"."sign_sessions" ADD CONSTRAINT "sign_sessions_pkey" PRIMARY KEY ("id");
