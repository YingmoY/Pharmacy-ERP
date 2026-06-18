// Package report 报表模块，提供销售、入库、库存和追溯日志报表查询与导出功能。
package report

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/report/handler"
	"github.com/YingmoY/PharmacyERP/internal/report/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 报表模块，整合路由、处理器与服务。
type Module struct {
	handler   *handler.Handler
	jwtSecret string
}

// NewModule 创建报表模块实例。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient interface{}, jwtSecret string) *Module {
	svc := service.NewReportService(db, log)
	h := handler.New(svc, log)
	return &Module{handler: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 注册报表模块路由至 /api/v1。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	reportGroup := group.Group("/reports", middleware.JWTAuth(m.jwtSecret))
	{
		// GET /reports/sales - 查询销售报表
		reportGroup.GET("/sales", m.handler.GetSalesReport)
		// POST /reports/sales/export - 导出销售报表
		reportGroup.POST("/sales/export", m.handler.ExportSalesReport)
		// GET /reports/inbound - 查询入库报表
		reportGroup.GET("/inbound", m.handler.GetInboundReport)
		// POST /reports/inbound/export - 导出入库报表
		reportGroup.POST("/inbound/export", m.handler.ExportInboundReport)
		// GET /reports/inventory - 查询库存报表
		reportGroup.GET("/inventory", m.handler.GetInventoryReport)
		// POST /reports/inventory/export - 导出库存报表
		reportGroup.POST("/inventory/export", m.handler.ExportInventoryReport)
		// GET /reports/trace-log - 查询追溯日志报表
		reportGroup.GET("/trace-log", m.handler.GetTraceLogReport)
		// POST /reports/trace-log/export - 导出追溯日志报表
		reportGroup.POST("/trace-log/export", m.handler.ExportTraceLogReport)
		// GET /reports/export-tasks/:task_id - 查询导出任务状态
		reportGroup.GET("/export-tasks/:task_id", m.handler.GetExportTask)
	}
}
