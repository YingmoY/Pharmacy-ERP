// Package handler 实现审计模块的 HTTP 处理层。
package handler

import (
	"strconv"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/audit/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Handler 审计 HTTP 处理器。
type Handler struct {
	repo repository.AuditRepo
	log  *zap.Logger
}

// New 创建审计处理器实例。
func New(repo repository.AuditRepo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

// loginLogRequest 登录日志查询参数。
type loginLogRequest struct {
	UserID    *int64  `form:"user_id"`
	Username  string  `form:"username"`
	Success   *bool   `form:"success"`
	StartDate string  `form:"start_date"`
	EndDate   string  `form:"end_date"`
	Page      int     `form:"page,default=1" binding:"min=1"`
	PageSize  int     `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListLoginLogs 处理 GET /audit/login-logs 请求。
func (h *Handler) ListLoginLogs(c *gin.Context) {
	var req loginLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	filter := repository.LoginLogFilter{
		UserID:   req.UserID,
		Username: req.Username,
		Success:  req.Success,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		filter.StartTime = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		// 结束日期取当天末尾
		end := t.AddDate(0, 0, 1)
		filter.EndTime = &end
	}

	logs, total, err := h.repo.ListLoginLogs(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("查询登录日志失败", zap.Error(err))
		core.Fail(c, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, logs))
}

// operationLogRequest 操作日志查询参数。
type operationLogRequest struct {
	OperatorID   *int64 `form:"operator_id"`
	OperatorName string `form:"operator_name"`
	Module       string `form:"module"`
	Action       string `form:"action"`
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	Page         int    `form:"page,default=1" binding:"min=1"`
	PageSize     int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListOperationLogs 处理 GET /audit/operation-logs 请求。
func (h *Handler) ListOperationLogs(c *gin.Context) {
	var req operationLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	filter := repository.OperationLogFilter{
		OperatorID:   req.OperatorID,
		OperatorName: req.OperatorName,
		Module:       req.Module,
		Action:       req.Action,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		filter.StartTime = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		end := t.AddDate(0, 0, 1)
		filter.EndTime = &end
	}

	logs, total, err := h.repo.ListOperationLogs(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("查询操作日志失败", zap.Error(err))
		core.Fail(c, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, logs))
}

// GetOperationLog 处理 GET /audit/operation-logs/:id 请求，查询操作日志详情。
func (h *Handler) GetOperationLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的日志 ID")
		return
	}

	log, err := h.repo.GetOperationLog(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			core.Fail(c, ecode.ErrNotFound.Code, ecode.ErrNotFound.Msg)
			return
		}
		h.log.Error("查询操作日志详情失败", zap.Int64("id", id), zap.Error(err))
		core.Fail(c, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
		return
	}

	core.Success(c, log)
}

// dataChangeLogRequest 数据变更日志查询参数。
type dataChangeLogRequest struct {
	TableName  string `form:"table_name"`
	RecordID   string `form:"record_id"`
	ChangeType string `form:"change_type"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListDataChangeLogs 处理 GET /audit/data-change-logs 请求。
func (h *Handler) ListDataChangeLogs(c *gin.Context) {
	var req dataChangeLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	filter := repository.DataChangeLogFilter{
		TableName:  req.TableName,
		RecordID:   req.RecordID,
		ChangeType: req.ChangeType,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		filter.StartTime = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		end := t.AddDate(0, 0, 1)
		filter.EndTime = &end
	}

	logs, total, err := h.repo.ListDataChangeLogs(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("查询数据变更日志失败", zap.Error(err))
		core.Fail(c, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, logs))
}

// securityEventRequest 安全事件查询参数。
type securityEventRequest struct {
	UserID    *int64 `form:"user_id"`
	EventType string `form:"event_type"`
	Handled   *bool  `form:"handled"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListSecurityEvents 处理 GET /audit/security-events 请求。
func (h *Handler) ListSecurityEvents(c *gin.Context) {
	var req securityEventRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	filter := repository.SecurityEventFilter{
		UserID:    req.UserID,
		EventType: req.EventType,
		Handled:   req.Handled,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		filter.StartTime = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		end := t.AddDate(0, 0, 1)
		filter.EndTime = &end
	}

	events, total, err := h.repo.ListSecurityEvents(c.Request.Context(), filter)
	if err != nil {
		h.log.Error("查询安全事件失败", zap.Error(err))
		core.Fail(c, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, events))
}
