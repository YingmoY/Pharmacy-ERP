// Package handler 实现告警模块的 HTTP 处理层。
package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/alert/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 告警 HTTP 处理器。
type Handler struct {
	svc service.AlertService
	log *zap.Logger
}

// New 创建告警处理器实例。
func New(svc service.AlertService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// listAlertsRequest 告警列表查询参数。
type listAlertsRequest struct {
	// Status 状态筛选：ACTIVE/RESOLVED/IGNORED，空表示不限。
	Status    string `form:"status"`
	EventType string `form:"event_type"`
	Severity  string `form:"severity"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListAlerts 处理 GET /alerts 请求，分页查询告警列表。
func (h *Handler) ListAlerts(c *gin.Context) {
	var req listAlertsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	alerts, total, err := h.svc.ListAlerts(c.Request.Context(), service.AlertFilter{
		Status:    req.Status,
		EventType: req.EventType,
		Severity:  req.Severity,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, alerts))
}

// resolveAlertRequest 解决告警请求体。
type resolveAlertRequest struct {
	Remark string `json:"remark"`
}

// ResolveAlert 处理 POST /alerts/:id/resolve 请求。
func (h *Handler) ResolveAlert(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的告警 ID")
		return
	}

	var req resolveAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	// 从 JWT 中获取当前操作人 ID
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.ResolveAlert(c.Request.Context(), id, userID, req.Remark); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "status": "RESOLVED"})
}

// ignoreAlertRequest 忽略告警请求体。
type ignoreAlertRequest struct {
	Reason string `json:"reason"`
}

// IgnoreAlert 处理 POST /alerts/:id/ignore 请求。
func (h *Handler) IgnoreAlert(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的告警 ID")
		return
	}

	var req ignoreAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	// 从 JWT 中获取当前操作人 ID
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.IgnoreAlert(c.Request.Context(), id, userID, req.Reason); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "status": "IGNORED"})
}

// nearExpireRequest 近效期药品列表查询参数。
type nearExpireRequest struct {
	Days     int `form:"days,default=30" binding:"min=1,max=365"`
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListNearExpire 处理 GET /alerts/near-expire 请求，实时查询近效期药品。
func (h *Handler) ListNearExpire(c *gin.Context) {
	var req nearExpireRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	items, total, err := h.svc.ListNearExpire(c.Request.Context(), req.Days, req.Page, req.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, items))
}

// lossCandidateRequest 盘亏候选列表查询参数。
type lossCandidateRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListLossCandidates 处理 GET /alerts/loss-candidates 请求，实时查询盘亏候选药品。
func (h *Handler) ListLossCandidates(c *gin.Context) {
	var req lossCandidateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	items, total, err := h.svc.ListLossCandidates(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, items))
}
