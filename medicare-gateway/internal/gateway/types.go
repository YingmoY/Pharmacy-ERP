package gateway

import (
	"encoding/json"
	"time"
)

type ERPRequest struct {
	ERPOrderNo string          `json:"erp_order_no"`
	Async      bool            `json:"async"`
	Input      json.RawMessage `json:"input"`
	RawInput   json.RawMessage `json:"Input"`
}

func (r ERPRequest) BusinessInput() json.RawMessage {
	if len(r.Input) > 0 {
		return r.Input
	}
	return r.RawInput
}

type MedicareEnvelope struct {
	Infno          string          `json:"infno"`
	MsgID          string          `json:"msgid"`
	MdtrtareaAdmvs string          `json:"mdtrtarea_admvs"`
	InsuplcAdmdvs  string          `json:"insuplc_admdvs"`
	RecerSysCode   string          `json:"recer_sys_code,omitempty"`
	DevNo          string          `json:"dev_no,omitempty"`
	DevSafeInfo    string          `json:"dev_safe_info,omitempty"`
	Cainfo         string          `json:"cainfo,omitempty"`
	Signtype       string          `json:"signtype,omitempty"`
	Infver         string          `json:"infver"`
	OpterType      string          `json:"opter_type"`
	Opter          string          `json:"opter"`
	OpterName      string          `json:"opter_name"`
	InfTime        string          `json:"inf_time"`
	FixmedinsCode  string          `json:"fixmedins_code"`
	FixmedinsName  string          `json:"fixmedins_name"`
	SignNo         string          `json:"sign_no,omitempty"`
	Input          json.RawMessage `json:"Input"`
}

type Response struct {
	TraceID string          `json:"trace_id"`
	MsgID   string          `json:"msgid,omitempty"`
	Infno   string          `json:"infno,omitempty"`
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func TimeText(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
