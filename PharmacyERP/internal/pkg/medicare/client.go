package medicare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Enabled bool          `mapstructure:"enabled"`
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// BusinessRequest is the unified body the ERP sends to the gateway.
type BusinessRequest struct {
	ERPOrderNo string         `json:"erp_order_no"`
	Async      bool           `json:"async,omitempty"`
	Input      map[string]any `json:"input"`
}

// GatewayResponse is the unified response from the gateway.
type GatewayResponse struct {
	TraceID string         `json:"trace_id"`
	MsgID   string         `json:"msgid"`
	Infno   string         `json:"infno"`
	Status  string         `json:"status"` // success | failed | queued
	Data    map[string]any `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// SetlInfo holds the settlement identifiers extracted from a 2102 response.
type SetlInfo struct {
	SetlID        string  `json:"setl_id"`
	MdtrtID       string  `json:"mdtrt_id"`
	PsnNo         string  `json:"psn_no"`
	FundPaySumamt float64 `json:"fund_pay_sumamt"`
	AcctPay       float64 `json:"acct_pay"`
}

// PersonInfo holds patient data extracted from a 1101 response.
type PersonInfo struct {
	PsnNo   string
	PsnName string
}

// ExtractSetlInfo parses a 2102 GatewayResponse and returns settlement identifiers.
func (r *GatewayResponse) ExtractSetlInfo() (*SetlInfo, error) {
	if r.Data == nil {
		return nil, fmt.Errorf("gateway response has no data")
	}
	output, _ := r.Data["output"].(map[string]any)
	if output == nil {
		return nil, fmt.Errorf("gateway response missing output")
	}
	raw, _ := output["setlinfo"].(map[string]any)
	if raw == nil {
		return nil, fmt.Errorf("gateway response missing setlinfo")
	}
	info := &SetlInfo{
		SetlID:  strVal(raw["setl_id"]),
		MdtrtID: strVal(raw["mdtrt_id"]),
		PsnNo:   strVal(raw["psn_no"]),
	}
	if v, ok := raw["fund_pay_sumamt"].(float64); ok {
		info.FundPaySumamt = v
	}
	if v, ok := raw["acct_pay"].(float64); ok {
		info.AcctPay = v
	}
	return info, nil
}

func strVal(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// Client calls the medicare gateway.
type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL: cfg.BaseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) call(ctx context.Context, path string, body any) (*GatewayResponse, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gateway unreachable: %w", err)
	}
	defer resp.Body.Close()

	var gwResp GatewayResponse
	if err := json.NewDecoder(resp.Body).Decode(&gwResp); err != nil {
		return nil, fmt.Errorf("decode gateway response: %w", err)
	}
	if gwResp.Status == "failed" {
		msg := gwResp.Error
		if msg == "" {
			msg = "unknown gateway error"
		}
		return &gwResp, fmt.Errorf("gateway %s failed: %s", gwResp.Infno, msg)
	}
	return &gwResp, nil
}

func (c *Client) PreSettle(ctx context.Context, req BusinessRequest) (*GatewayResponse, error) {
	return c.call(ctx, "/api/2101", req)
}

func (c *Client) Settle(ctx context.Context, req BusinessRequest) (*GatewayResponse, error) {
	return c.call(ctx, "/api/2102", req)
}

func (c *Client) UploadSale(ctx context.Context, req BusinessRequest) (*GatewayResponse, error) {
	return c.call(ctx, "/api/3505", req)
}

func (c *Client) CancelSettle(ctx context.Context, req BusinessRequest) (*GatewayResponse, error) {
	return c.call(ctx, "/api/2103", req)
}

func (c *Client) QueryPerson(ctx context.Context, req BusinessRequest) (*GatewayResponse, error) {
	return c.call(ctx, "/api/1101", req)
}

// ExtractPersonInfo parses a 1101 GatewayResponse and returns patient identifiers.
func (r *GatewayResponse) ExtractPersonInfo() (*PersonInfo, error) {
	if r.Data == nil {
		return nil, fmt.Errorf("gateway response has no data")
	}
	output, _ := r.Data["output"].(map[string]any)
	if output == nil {
		return nil, fmt.Errorf("gateway response missing output")
	}
	baseinfo, _ := output["baseinfo"].(map[string]any)
	if baseinfo == nil {
		return nil, fmt.Errorf("gateway response missing baseinfo")
	}
	info := &PersonInfo{
		PsnNo:   strVal(baseinfo["psn_no"]),
		PsnName: strVal(baseinfo["psn_name"]),
	}
	if info.PsnNo == "" {
		return nil, fmt.Errorf("psn_no not found in response")
	}
	return info, nil
}
