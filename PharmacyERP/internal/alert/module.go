// Package alert 告警模块，提供告警查询、处理及实时预警视图。
package alert

import (
	"github.com/YingmoY/PharmacyERP/internal/alert/handler"
	"github.com/YingmoY/PharmacyERP/internal/alert/repository"
	"github.com/YingmoY/PharmacyERP/internal/alert/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 告警模块，整合路由、处理器、服务与仓储。
type Module struct {
	handler *handler.Handler
	jwtSecret string
}

// NewModule 创建告警模块实例。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient interface{}, jwtSecret string) *Module {
	repo := repository.NewAlertRepo(db)
	svc := service.NewAlertService(repo, db, log)
	h := handler.New(svc, log)
	return &Module{handler: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 注册告警模块路由至 /api/v1。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	// 所有告警接口均需 JWT 鉴权
	alertGroup := group.Group("/alerts", middleware.JWTAuth(m.jwtSecret))
	{
		// GET /alerts - 分页查询告警列表（可按 status/event_type/severity 过滤）
		alertGroup.GET("", m.handler.ListAlerts)
		// POST /alerts/:id/resolve - 解决告警
		alertGroup.POST("/:id/resolve", m.handler.ResolveAlert)
		// POST /alerts/:id/ignore - 忽略告警
		alertGroup.POST("/:id/ignore", m.handler.IgnoreAlert)
		// GET /alerts/near-expire - 实时近效期药品列表
		alertGroup.GET("/near-expire", m.handler.ListNearExpire)
		// GET /alerts/loss-candidates - 实时盘亏候选药品列表
		alertGroup.GET("/loss-candidates", m.handler.ListLossCandidates)
	}
}
