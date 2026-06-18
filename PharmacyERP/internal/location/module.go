package location

import (
	"github.com/YingmoY/PharmacyERP/internal/location/handler"
	"github.com/YingmoY/PharmacyERP/internal/location/repository"
	"github.com/YingmoY/PharmacyERP/internal/location/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 location 业务模块的 HTTP 适配层。
type Module struct {
	db       *gorm.DB
	log      *zap.Logger
	mqClient *mq.Client
	handler  *handler.Handler
}

// NewModule 构建货位模块，注入依赖。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient *mq.Client) *Module {
	repo := repository.NewLocationRepo()
	svc := service.NewLocationService(db, repo, log, mqClient)
	h := handler.NewHandler(svc, log)
	return &Module{
		db:       db,
		log:      log,
		mqClient: mqClient,
		handler:  h,
	}
}

// RegisterRoutes 注册货位模块路由。
// 带静态路径的路由（/areas、/code/:location_code）必须先于动态参数路由（/:id）注册。
// JWT 鉴权由上层路由组统一应用，此处不再重复添加。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/locations")

	// 静态路径优先于动态路径参数
	g.GET("/areas", m.handler.GetAreas)
	g.GET("/code/:location_code", m.handler.GetLocationByCode)

	g.GET("", m.handler.ListLocations)
	g.POST("", m.handler.CreateLocation)
	g.GET("/:id", m.handler.GetLocation)
	g.PUT("/:id", m.handler.UpdateLocation)
	g.DELETE("/:id", m.handler.DeleteLocation)
	g.PATCH("/:id/status", m.handler.UpdateLocationStatus)
	g.GET("/:id/drugs", m.handler.GetLocationDrugs)
}
