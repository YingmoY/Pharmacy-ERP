package medicare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Result struct {
	StatusCode int
	Body       json.RawMessage
	Elapsed    time.Duration
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Call(ctx context.Context, infno string, payload any) (Result, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return Result{}, err
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, infno)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	start := time.Now()
	resp, err := c.http.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return Result{Elapsed: elapsed}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{StatusCode: resp.StatusCode, Elapsed: elapsed}, err
	}
	if !json.Valid(respBody) {
		respBody, _ = json.Marshal(map[string]string{"raw": string(respBody)})
	}
	return Result{StatusCode: resp.StatusCode, Body: respBody, Elapsed: elapsed}, nil
}

func Infcode(body json.RawMessage) string {
	var envelope struct {
		Infcode any `json:"infcode"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	switch value := envelope.Infcode.(type) {
	case string:
		return value
	case float64:
		if value == 0 {
			return "0"
		}
		return fmt.Sprintf("%.0f", value)
	default:
		return ""
	}
}

func ErrMsg(body json.RawMessage) string {
	var envelope struct {
		ErrMsg string `json:"err_msg"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	return envelope.ErrMsg
}
