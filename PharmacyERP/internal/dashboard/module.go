// Package dashboard 仪表盘模块，提供运营数据概览功能。
package dashboard

import (
	"github.com/YingmoY/PharmacyERP/internal/dashboard/handler"
	"github.com/YingmoY/PharmacyERP/internal/dashboard/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 仪表盘模块，整合路由、处理器与服务。
type Module struct {
	handler   *handler.Handler
	jwtSecret string
}

// NewModule 创建仪表盘模块实例。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient interface{}, jwtSecret string) *Module {
	svc := service.NewDashboardService(db, log)
	h := handler.New(svc, log)
	return &Module{handler: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 注册仪表盘模块路由至 /api/v1。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	dashGroup := group.Group("/dashboard", middleware.JWTAuth(m.jwtSecret))
	{
		// GET /dashboard/overview - 仪表盘概览数据
		dashGroup.GET("/overview", m.handler.GetOverview)
		// GET /dashboard/sales-trend - 销售趋势（query: days，默认 7 天）
		dashGroup.GET("/sales-trend", m.handler.GetSalesTrend)
		// GET /dashboard/top-drugs - 药品销售排行（query: limit，days）
		dashGroup.GET("/top-drugs", m.handler.GetTopDrugs)
		// GET /dashboard/inbound-stats - 入库统计（query: days，默认 30 天）
		dashGroup.GET("/inbound-stats", m.handler.GetInboundStats)
		// GET /dashboard/inventory-stats - 库存状态统计
		dashGroup.GET("/inventory-stats", m.handler.GetInventoryStats)
	}
}
