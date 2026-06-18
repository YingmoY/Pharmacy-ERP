package drug

import (
	"github.com/YingmoY/PharmacyERP/internal/drug/handler"
	"github.com/YingmoY/PharmacyERP/internal/drug/repository"
	"github.com/YingmoY/PharmacyERP/internal/drug/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 drug 业务模块的 HTTP 适配层。
type Module struct {
	db       *gorm.DB
	log      *zap.Logger
	mqClient *mq.Client
	handler  *handler.Handler
}

// NewModule 构建药品模块，注入依赖。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient *mq.Client) *Module {
	repo := repository.NewDrugRepo()
	svc := service.NewDrugService(db, repo, log, mqClient)
	h := handler.NewHandler(svc, log)
	return &Module{
		db:       db,
		log:      log,
		mqClient: mqClient,
		handler:  h,
	}
}

// RegisterRoutes 注册药品模块路由到指定路由组。
// 路由前缀由上层调用者传入（如 /api/v1）。
// JWT 鉴权由上层路由组统一应用，此处不再重复添加。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/drugs")

	// 注意：带静态路径段的路由（/search、/code/:drug_code）必须在
	// 带动态参数段（/:id）的路由之前注册，避免参数路由捕获静态路径。
	g.GET("/search", m.handler.SearchDrugs)
	g.GET("/code/:drug_code", m.handler.GetDrugByCode)

	g.GET("", m.handler.ListDrugs)
	g.POST("", m.handler.CreateDrug)
	g.GET("/:id", m.handler.GetDrug)
	g.PUT("/:id", m.handler.UpdateDrug)
	g.DELETE("/:id", m.handler.DeleteDrug)
	g.PATCH("/:id/status", m.handler.UpdateDrugStatus)
	g.GET("/:id/inventory-summary", m.handler.GetInventorySummary)
	g.GET("/:id/sale-info", m.handler.GetDrugSaleInfo)
}
