create schema if not exists medicare;

create table if not exists medicare.sign_sessions (
    id bigserial primary key,
    sign_no text not null,
    operator_code text not null,
    operator_name text not null,
    request_msgid text not null,
    response jsonb not null,
    signed_at timestamptz not null default now(),
    expires_at timestamptz,
    active boolean not null default true
);

create index if not exists idx_sign_sessions_active_signed_at
    on medicare.sign_sessions(active, signed_at desc);

create table if not exists medicare.audit_logs (
    id bigserial primary key,
    trace_id uuid not null,
    msgid text not null,
    infno text not null,
    erp_order_no text,
    operation text not null,
    status text not null,
    request_payload jsonb not null,
    response_payload jsonb,
    error_message text,
    http_status integer,
    elapsed_ms integer,
    sign_no text,
    async_task_id uuid,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_audit_logs_trace_id on medicare.audit_logs(trace_id);
create index if not exists idx_audit_logs_msgid on medicare.audit_logs(msgid);
create index if not exists idx_audit_logs_order on medicare.audit_logs(erp_order_no);
create index if not exists idx_audit_logs_infno_created_at on medicare.audit_logs(infno, created_at desc);

create table if not exists medicare.async_tasks (
    id uuid primary key,
    infno text not null,
    operation text not null,
    erp_order_no text,
    status text not null,
    attempts integer not null default 0,
    max_attempts integer not null default 5,
    request_payload jsonb not null,
    response_payload jsonb,
    error_message text,
    next_run_at timestamptz not null default now(),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_async_tasks_status_next_run
    on medicare.async_tasks(status, next_run_at);
