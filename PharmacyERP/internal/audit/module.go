// Package audit 审计模块，提供登录日志、操作日志、数据变更日志及安全事件查询接口。
package audit

import (
	"github.com/YingmoY/PharmacyERP/internal/audit/handler"
	"github.com/YingmoY/PharmacyERP/internal/audit/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 审计模块，整合路由、处理器与仓储。
type Module struct {
	handler   *handler.Handler
	jwtSecret string
}

// NewModule 创建审计模块实例。
func NewModule(db *gorm.DB, log *zap.Logger, jwtSecret string) *Module {
	repo := repository.NewAuditRepo(db)
	h := handler.New(repo, log)
	return &Module{handler: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 注册审计模块路由至 /api/v1/audit。
// 所有审计接口均需 JWT 鉴权。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auditGroup := group.Group("/audit", middleware.JWTAuth(m.jwtSecret))
	{
		// GET /audit/login-logs - 分页查询登录日志
		auditGroup.GET("/login-logs", m.handler.ListLoginLogs)
		// GET /audit/operation-logs - 分页查询操作日志
		auditGroup.GET("/operation-logs", m.handler.ListOperationLogs)
		// GET /audit/operation-logs/:id - 查询操作日志详情
		auditGroup.GET("/operation-logs/:id", m.handler.GetOperationLog)
		// GET /audit/data-change-logs - 分页查询数据变更日志
		auditGroup.GET("/data-change-logs", m.handler.ListDataChangeLogs)
		// GET /audit/security-events - 分页查询安全事件
		auditGroup.GET("/security-events", m.handler.ListSecurityEvents)
	}
}
