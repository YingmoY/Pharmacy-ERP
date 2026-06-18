// Package handler 实现通知模块的 HTTP 处理层。
package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/notification/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 通知 HTTP 处理器。
type Handler struct {
	svc service.NotificationService
	log *zap.Logger
}

// New 创建通知处理器实例。
func New(svc service.NotificationService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// listNotificationsRequest 通知列表查询参数。
type listNotificationsRequest struct {
	// Read 已读状态筛选：不传=全部，"true"=已读，"false"=未读。
	Read     *bool `form:"read"`
	Page     int   `form:"page,default=1" binding:"min=1"`
	PageSize int   `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListNotifications 处理 GET /notifications 请求，查询当前用户通知列表。
func (h *Handler) ListNotifications(c *gin.Context) {
	var req listNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	// 从 JWT 获取当前用户 ID
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	notifications, total, err := h.svc.ListNotifications(c.Request.Context(), userID, service.NotificationFilter{
		Read:     req.Read,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, notifications))
}

// GetUnreadCount 处理 GET /notifications/unread-count 请求，获取未读通知数量。
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	count, err := h.svc.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"unread_count": count})
}

// MarkRead 处理 POST /notifications/:id/read 请求，标记指定通知为已读。
func (h *Handler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的通知 ID")
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.MarkRead(c.Request.Context(), id, userID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id})
}

// MarkAllRead 处理 POST /notifications/read-all 请求，将当前用户所有通知标记为已读。
func (h *Handler) MarkAllRead(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok || userID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"message": "所有通知已标记为已读"})
}
