// Package notification 通知模块，提供用户消息通知管理功能。
package notification

import (
	"github.com/YingmoY/PharmacyERP/internal/notification/handler"
	"github.com/YingmoY/PharmacyERP/internal/notification/repository"
	"github.com/YingmoY/PharmacyERP/internal/notification/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 通知模块，整合路由、处理器、服务与仓储。
type Module struct {
	handler   *handler.Handler
	jwtSecret string
}

// NewModule 创建通知模块实例。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient interface{}, jwtSecret string) *Module {
	repo := repository.NewNotificationRepo(db)
	svc := service.NewNotificationService(repo, log)
	h := handler.New(svc, log)
	return &Module{handler: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 注册通知模块路由至 /api/v1。
// 所有通知接口均需要 JWT 鉴权，用户只能查看和操作自己的通知。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	notifyGroup := group.Group("/notifications", middleware.JWTAuth(m.jwtSecret))
	{
		// GET /notifications - 查询当前用户通知列表（支持按已读状态过滤）
		notifyGroup.GET("", m.handler.ListNotifications)
		// GET /notifications/unread-count - 获取未读通知数量
		notifyGroup.GET("/unread-count", m.handler.GetUnreadCount)
		// POST /notifications/:id/read - 将指定通知标记为已读
		notifyGroup.POST("/:id/read", m.handler.MarkRead)
		// POST /notifications/read-all - 将当前用户所有通知标记为已读
		notifyGroup.POST("/read-all", m.handler.MarkAllRead)
	}
}
