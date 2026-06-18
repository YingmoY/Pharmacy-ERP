package gateway

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"medicare-gateway/internal/config"
	"medicare-gateway/internal/db"
	"medicare-gateway/internal/medicare"
)

type Queue interface {
	PublishTask(ctx context.Context, task TaskMessage) error
	Close() error
}

type TaskMessage struct {
	TaskID     string          `json:"task_id"`
	TraceID    string          `json:"trace_id"`
	Infno      string          `json:"infno"`
	Operation  string          `json:"operation"`
	ERPOrderNo string          `json:"erp_order_no,omitempty"`
	Payload    json.RawMessage `json:"payload"`
}

type Service struct {
	cfg    config.Config
	store  *db.Store
	client *medicare.Client
	queue  Queue
	logger *slog.Logger
}

func NewService(cfg config.Config, store *db.Store, client *medicare.Client, queue Queue, logger *slog.Logger) *Service {
	return &Service{cfg: cfg, store: store, client: client, queue: queue, logger: logger}
}

func (s *Service) SignIn(ctx context.Context) (Response, error) {
	traceID := NewTraceID()
	input := map[string]any{
		"signIn": map[string]any{
			"opter_no":   s.cfg.Operator,
			"opter_name": s.cfg.OperatorName,
			"mac":        "00-00-00-00-00-00",
			"ip":         "127.0.0.1",
		},
	}
	envelope := s.envelope("9001", "", mustJSON(input))
	return s.callAndAudit(ctx, traceID, "9001", "sign_in", "", envelope, true)
}

func (s *Service) PersonInfo(ctx context.Context) (Response, error) {
	input := mustJSON(map[string]any{
		"data": map[string]any{
			"mdtrt_cert_type": s.cfg.DefaultCertType,
			"mdtrt_cert_no":   s.cfg.DefaultCertNo,
			"card_sn":         "",
			"begntime":        TimeText(time.Now()),
			"psn_cert_type":   "01",
			"certno":          s.cfg.DefaultCertNo,
			"psn_name":        "",
		},
	})
	return s.Call(ctx, "1101", ERPRequest{Input: input})
}

func (s *Service) Call(ctx context.Context, infno string, req ERPRequest) (Response, error) {
	traceID := NewTraceID()
	input := req.BusinessInput()
	if len(input) == 0 {
		input = json.RawMessage(`{}`)
	}

	signNo := ""
	if infno != "9001" {
		session, err := s.store.LatestSignSession(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			sign, signErr := s.SignIn(ctx)
			if signErr != nil {
				return sign, signErr
			}
			var signBody struct {
				Data json.RawMessage `json:"data"`
			}
			_ = json.Unmarshal(mustJSON(sign), &signBody)
			session, err = s.store.LatestSignSession(ctx)
		}
		if err != nil {
			return Response{TraceID: traceID, Infno: infno, Status: "failed", Error: err.Error()}, err
		}
		signNo = session.SignNo
	}

	envelope := s.envelope(infno, signNo, input)
	if req.Async && s.queue != nil && isAsyncFriendly(infno) {
		taskID := NewTraceID()
		payload := mustJSON(envelope)
		msg := TaskMessage{TaskID: taskID, TraceID: traceID, Infno: infno, Operation: operationName(infno), ERPOrderNo: req.ERPOrderNo, Payload: payload}
		if err := s.queue.PublishTask(ctx, msg); err != nil {
			s.logger.Warn("async publish failed, fallback to sync call", "trace_id", traceID, "infno", infno, "error", err)
		} else {
			_ = s.store.InsertAsyncTask(ctx, db.AsyncTask{
				ID:             taskID,
				Infno:          infno,
				Operation:      operationName(infno),
				ERPOrderNo:     req.ERPOrderNo,
				Status:         "queued",
				RequestPayload: payload,
			})
			_ = s.store.InsertAudit(ctx, db.AuditLog{
				TraceID:        traceID,
				MsgID:          envelope.MsgID,
				Infno:          infno,
				ERPOrderNo:     req.ERPOrderNo,
				Operation:      operationName(infno),
				Status:         "queued",
				RequestPayload: payload,
				SignNo:         signNo,
				AsyncTaskID:    taskID,
			})
			return Response{TraceID: traceID, MsgID: envelope.MsgID, Infno: infno, Status: "queued"}, nil
		}
	}

	return s.callAndAudit(ctx, traceID, infno, operationName(infno), req.ERPOrderNo, envelope, infno == "9001")
}

func (s *Service) ProcessTask(ctx context.Context, task TaskMessage) error {
	var envelope MedicareEnvelope
	if err := json.Unmarshal(task.Payload, &envelope); err != nil {
		_ = s.store.UpdateAsyncTask(ctx, db.AsyncTask{ID: task.TaskID, Status: "failed", ErrorMessage: err.Error()})
		return err
	}
	resp, err := s.callAndAudit(ctx, task.TraceID, task.Infno, task.Operation, task.ERPOrderNo, envelope, task.Infno == "9001")
	status := "done"
	errText := ""
	if err != nil || resp.Status == "failed" {
		status = "failed"
		errText = resp.Error
		if errText == "" && err != nil {
			errText = err.Error()
		}
	}
	_ = s.store.UpdateAsyncTask(ctx, db.AsyncTask{
		ID:              task.TaskID,
		Status:          status,
		ResponsePayload: resp.Data,
		ErrorMessage:    errText,
	})
	return err
}

func (s *Service) callAndAudit(ctx context.Context, traceID, infno, operation, orderNo string, envelope MedicareEnvelope, saveSign bool) (Response, error) {
	payload := mustJSON(envelope)
	result, err := s.client.Call(ctx, infno, envelope)
	status := "success"
	errText := ""
	if err != nil {
		status = "failed"
		errText = err.Error()
	} else if result.StatusCode < 200 || result.StatusCode >= 300 || medicare.Infcode(result.Body) != "0" {
		status = "failed"
		errText = medicare.ErrMsg(result.Body)
	}

	if saveSign && status == "success" {
		signNo := extractSignNo(result.Body)
		if signNo != "" {
			_ = s.store.SaveSignSession(ctx, db.SignSession{
				SignNo:       signNo,
				OperatorCode: s.cfg.Operator,
				OperatorName: s.cfg.OperatorName,
				RequestMsgID: envelope.MsgID,
				Response:     result.Body,
			})
			envelope.SignNo = signNo
		}
	}

	auditErr := s.store.InsertAudit(ctx, db.AuditLog{
		TraceID:         traceID,
		MsgID:           envelope.MsgID,
		Infno:           infno,
		ERPOrderNo:      orderNo,
		Operation:       operation,
		Status:          status,
		RequestPayload:  payload,
		ResponsePayload: result.Body,
		ErrorMessage:    errText,
		HTTPStatus:      result.StatusCode,
		ElapsedMS:       int(result.Elapsed / time.Millisecond),
		SignNo:          envelope.SignNo,
	})
	if auditErr != nil {
		s.logger.Error("audit insert failed", "trace_id", traceID, "error", auditErr)
	}
	if err != nil {
		return Response{TraceID: traceID, MsgID: envelope.MsgID, Infno: infno, Status: status, Error: errText}, err
	}
	return Response{TraceID: traceID, MsgID: envelope.MsgID, Infno: infno, Status: status, Data: result.Body, Error: errText}, nil
}

func (s *Service) envelope(infno, signNo string, input json.RawMessage) MedicareEnvelope {
	return MedicareEnvelope{
		Infno:          infno,
		MsgID:          NewMsgID(infno),
		MdtrtareaAdmvs: s.cfg.AreaCode,
		InsuplcAdmdvs:  s.cfg.InsuplcCode,
		RecerSysCode:   "CSB",
		Infver:         "1.0.0",
		OpterType:      s.cfg.OperatorType,
		Opter:          s.cfg.Operator,
		OpterName:      s.cfg.OperatorName,
		InfTime:        TimeText(time.Now()),
		FixmedinsCode:  s.cfg.FixmedinsCode,
		FixmedinsName:  s.cfg.FixmedinsName,
		SignNo:         signNo,
		Input:          input,
	}
}

func extractSignNo(body json.RawMessage) string {
	var out struct {
		Output map[string]any `json:"output"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return ""
	}
	for _, key := range []string{"sign_no", "signNo"} {
		if value, ok := out.Output[key].(string); ok {
			return value
		}
	}
	raw, _ := json.Marshal(out.Output)
	var nested struct {
		Signinoutb struct {
			SignNo string `json:"sign_no"`
		} `json:"signinoutb"`
		SignIn struct {
			SignNo string `json:"sign_no"`
		} `json:"signIn"`
	}
	_ = json.Unmarshal(raw, &nested)
	if nested.Signinoutb.SignNo != "" {
		return nested.Signinoutb.SignNo
	}
	return nested.SignIn.SignNo
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func isAsyncFriendly(infno string) bool {
	switch infno {
	case "2102", "2103", "3505", "3201", "3202":
		return true
	default:
		return false
	}
}

func operationName(infno string) string {
	switch infno {
	case "9001":
		return "sign_in"
	case "1101":
		return "person_info"
	case "2101":
		return "drugstore_presettle"
	case "2102":
		return "drugstore_settle"
	case "2103":
		return "drugstore_settle_cancel"
	case "3505":
		return "goods_sale_upload"
	case "3201":
		return "settlement_list_query"
	case "3202":
		return "settlement_detail_query"
	default:
		return "medicare_call"
	}
}
