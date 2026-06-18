package db

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

//go:embed schema.sql
var initSQL string

type Store struct {
	db *sql.DB
}

type AuditLog struct {
	TraceID         string
	MsgID           string
	Infno           string
	ERPOrderNo      string
	Operation       string
	Status          string
	RequestPayload  json.RawMessage
	ResponsePayload json.RawMessage
	ErrorMessage    string
	HTTPStatus      int
	ElapsedMS       int
	SignNo          string
	AsyncTaskID     string
}

type SignSession struct {
	SignNo       string          `json:"sign_no"`
	OperatorCode string          `json:"operator_code"`
	OperatorName string          `json:"operator_name"`
	RequestMsgID string          `json:"request_msgid"`
	Response     json.RawMessage `json:"response"`
	SignedAt     time.Time       `json:"signed_at"`
}

type AsyncTask struct {
	ID              string
	Infno           string
	Operation       string
	ERPOrderNo      string
	Status          string
	RequestPayload  json.RawMessage
	ResponsePayload json.RawMessage
	ErrorMessage    string
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, initSQL)
	return err
}

func (s *Store) InsertAudit(ctx context.Context, log AuditLog) error {
	req := string(log.RequestPayload)
	resp := string(log.ResponsePayload)
	_, err := s.db.ExecContext(ctx, `
insert into medicare.audit_logs (
    trace_id, msgid, infno, erp_order_no, operation, status,
    request_payload, response_payload, error_message, http_status,
    elapsed_ms, sign_no, async_task_id
) values ($1::uuid, $2, $3, nullif($4, ''), $5, $6, $7, nullif($8, '')::jsonb,
          nullif($9, ''), nullif($10, 0), nullif($11, 0), nullif($12, ''), nullif($13, '')::uuid)`,
		log.TraceID, log.MsgID, log.Infno, log.ERPOrderNo, log.Operation, log.Status,
		req, resp, log.ErrorMessage, log.HTTPStatus,
		log.ElapsedMS, log.SignNo, log.AsyncTaskID)
	return err
}

func (s *Store) SaveSignSession(ctx context.Context, session SignSession) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `update medicare.sign_sessions set active = false where active = true`); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
insert into medicare.sign_sessions (sign_no, operator_code, operator_name, request_msgid, response, signed_at, active)
values ($1, $2, $3, $4, $5::jsonb, now(), true)`,
		session.SignNo, session.OperatorCode, session.OperatorName, session.RequestMsgID, string(session.Response))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) LatestSignSession(ctx context.Context) (SignSession, error) {
	var session SignSession
	err := s.db.QueryRowContext(ctx, `
select sign_no, operator_code, operator_name, request_msgid, response, signed_at
from medicare.sign_sessions
where active = true
order by signed_at desc
limit 1`).Scan(&session.SignNo, &session.OperatorCode, &session.OperatorName, &session.RequestMsgID, &session.Response, &session.SignedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return session, sql.ErrNoRows
	}
	return session, err
}

func (s *Store) InsertAsyncTask(ctx context.Context, task AsyncTask) error {
	req := string(task.RequestPayload)
	_, err := s.db.ExecContext(ctx, `
insert into medicare.async_tasks (id, infno, operation, erp_order_no, status, request_payload)
values ($1::uuid, $2, $3, nullif($4, ''), $5, $6::jsonb)`,
		task.ID, task.Infno, task.Operation, task.ERPOrderNo, task.Status, req)
	return err
}

func (s *Store) UpdateAsyncTask(ctx context.Context, task AsyncTask) error {
	resp := string(task.ResponsePayload)
	_, err := s.db.ExecContext(ctx, `
update medicare.async_tasks
set status = $2,
    attempts = attempts + 1,
    response_payload = nullif($3, '')::jsonb,
    error_message = nullif($4, ''),
    updated_at = now()
where id = $1::uuid`,
		task.ID, task.Status, resp, task.ErrorMessage)
	return err
}

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
