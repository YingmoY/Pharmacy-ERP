package handler

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/trace/service"
	"github.com/gin-gonic/gin"
)

// Handler 是追溯码查询模块 HTTP 适配层
type Handler struct {
	svc service.TraceService
}

// New 创建追溯码查询 Handler
func New(svc service.TraceService) *Handler {
	return &Handler{svc: svc}
}

// GetTraceInfo GET /trace/:trace_code
func (h *Handler) GetTraceInfo(c *gin.Context) {
	traceCode := c.Param("trace_code")
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	info, err := h.svc.GetTraceInfo(c.Request.Context(), traceCode)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, info)
}

// GetFullChain GET /trace/:trace_code/full-chain
func (h *Handler) GetFullChain(c *gin.Context) {
	traceCode := c.Param("trace_code")
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	logs, err := h.svc.GetFullChain(c.Request.Context(), traceCode)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, logs)
}

// getLogsRequest 追溯日志分页查询参数
type getLogsRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// GetLogs GET /trace/:trace_code/logs
func (h *Handler) GetLogs(c *gin.Context) {
	traceCode := c.Param("trace_code")
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	var req getLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	logs, total, err := h.svc.GetLogs(c.Request.Context(), traceCode, service.LogFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, logs))
}

// Validate POST /trace/validate
func (h *Handler) Validate(c *gin.Context) {
	var req service.ValidateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	result, err := h.svc.Validate(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, result)
}
