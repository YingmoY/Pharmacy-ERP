package shelving

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/shelving/handler"
	"github.com/YingmoY/PharmacyERP/internal/shelving/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 shelving 业务模块的 HTTP 适配层。
//
// 职责：
// 1) 依赖注入：组合服务层与处理器层；
// 2) 将上架相关接口挂载到路由组。
type Module struct {
	shelvingHdl *handler.ShelvingHandler
	jwtSecret   string
}

// NewModule 构建上架模块，注入数据库连接与日志。
func NewModule(db *gorm.DB, logger *zap.Logger, jwtSecret string) *Module {
	shelvingSvc := service.NewShelvingService(db, logger)
	return &Module{
		shelvingHdl: handler.NewShelvingHandler(shelvingSvc),
		jwtSecret:   jwtSecret,
	}
}

// RegisterRoutes 注册上架模块路由。
//
// 路由前缀：/api/v1/shelving
// 开放接口：
// - GET  /shelving/pending     待上架列表（分页）；
// - POST /shelving/scan        单条扫码上架；
// - POST /shelving/batch       批量上架；
// - POST /shelving/relocate    货位调整（移位）；
// - GET  /shelving/mix-check   货位混合陈列检查。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	shelvingGroup := group.Group("/shelving", middleware.JWTAuth(m.jwtSecret))
	shelvingGroup.GET("/pending", m.shelvingHdl.GetPending)
	shelvingGroup.POST("/scan", m.shelvingHdl.Scan)
	shelvingGroup.POST("/batch", m.shelvingHdl.Batch)
	shelvingGroup.POST("/relocate", m.shelvingHdl.Relocate)
	shelvingGroup.GET("/mix-check", m.shelvingHdl.MixCheck)
}
