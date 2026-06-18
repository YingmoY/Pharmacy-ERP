package inventory

import (
	"github.com/YingmoY/PharmacyERP/internal/inventory/handler"
	"github.com/YingmoY/PharmacyERP/internal/inventory/repository"
	"github.com/YingmoY/PharmacyERP/internal/inventory/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 inventory 业务模块的 HTTP 适配层。
type Module struct {
	inventoryHdl *handler.InventoryHandler
	jwtSecret    string
}

// NewModule 构建 inventory 模块。
func NewModule(db *gorm.DB, logger *zap.Logger, jwtSecret string) *Module {
	orderRepo := repository.NewInboundOrder()
	traceRepo := repository.NewTraceInventory()
	_ = service.NewInboundService(db, orderRepo, traceRepo, logger)
	inventorySvc := service.NewInventoryService(db, logger)

	return &Module{
		inventoryHdl: handler.NewInventoryHandler(inventorySvc),
		jwtSecret:    jwtSecret,
	}
}

// RegisterRoutes 注册 inventory 模块路由。
//
// 路由前缀：/api/v1/inventory、/api/v1/inventory-adjustments
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := middleware.JWTAuth(m.jwtSecret)

	// ─── 库存查询接口 ──────────────────────────────────────────────────
	inventoryGroup := group.Group("/inventory", auth)
	inventoryGroup.GET("", m.inventoryHdl.ListInventory)
	inventoryGroup.GET("/summary", m.inventoryHdl.GetSummary)
	inventoryGroup.GET("/pending-shelving", m.inventoryHdl.ListPendingShelving)
	inventoryGroup.GET("/near-expire", m.inventoryHdl.ListNearExpire)
	inventoryGroup.GET("/recommend-sale", m.inventoryHdl.ListRecommendSale)
	inventoryGroup.GET("/drugs/:drug_id", m.inventoryHdl.ListByDrug)
	inventoryGroup.GET("/locations/:location_id", m.inventoryHdl.ListByLocation)

	// ─── 手动状态变更接口（仅管理员）─────────────────────────────────────
	inventoryGroup.PATCH("/:trace_code/status", m.inventoryHdl.ManualStatusChange)

	// ─── 库存调整记录接口 ──────────────────────────────────────────────
	adjustGroup := group.Group("/inventory-adjustments", auth)
	adjustGroup.GET("", m.inventoryHdl.ListAdjustments)
	adjustGroup.POST("", m.inventoryHdl.CreateAdjustment)
	adjustGroup.GET("/:id", m.inventoryHdl.GetAdjustment)
}
