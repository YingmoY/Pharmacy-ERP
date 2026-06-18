package supplier

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/supplier/handler"
	"github.com/YingmoY/PharmacyERP/internal/supplier/repository"
	"github.com/YingmoY/PharmacyERP/internal/supplier/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 supplier 业务模块的 HTTP 适配层。
type Module struct {
	db       *gorm.DB
	log      *zap.Logger
	mqClient *mq.Client
	handler  *handler.Handler
}

// NewModule 构建供应商模块，注入依赖。
func NewModule(db *gorm.DB, log *zap.Logger, mqClient *mq.Client) *Module {
	repo := repository.NewSupplierRepo()
	svc := service.NewSupplierService(db, repo, log, mqClient)
	h := handler.NewHandler(svc, log)
	return &Module{
		db:       db,
		log:      log,
		mqClient: mqClient,
		handler:  h,
	}
}

// RegisterRoutes 注册供应商模块路由。
// JWT 鉴权由上层路由组统一应用，此处不再重复添加。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/suppliers")

	g.GET("", m.handler.ListSuppliers)
	g.POST("", m.handler.CreateSupplier)
	g.GET("/:id", m.handler.GetSupplier)
	g.PUT("/:id", m.handler.UpdateSupplier)
	g.DELETE("/:id", m.handler.DeleteSupplier)
	g.PATCH("/:id/status", m.handler.UpdateSupplierStatus)
}
