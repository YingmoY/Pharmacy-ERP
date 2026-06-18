package user

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	rolerepo "github.com/YingmoY/PharmacyERP/internal/role/repository"
	"github.com/YingmoY/PharmacyERP/internal/user/handler"
	"github.com/YingmoY/PharmacyERP/internal/user/repository"
	"github.com/YingmoY/PharmacyERP/internal/user/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 用户管理模块
type Module struct {
	handler *handler.Handler
	cfg     *config.JWTConfig
	db      *gorm.DB
}

// NewModule 创建用户管理模块，初始化所有依赖
func NewModule(db *gorm.DB, log *zap.Logger, cfg *config.JWTConfig) *Module {
	// 初始化仓储层
	uRepo := repository.NewUserRepository(db)
	rRepo := rolerepo.NewRoleRepository(db)

	// 初始化服务层
	svc := service.NewUserService(uRepo, rRepo, log)

	// 初始化处理器层
	h := handler.NewHandler(svc, log)

	return &Module{handler: h, cfg: cfg, db: db}
}

// RegisterRoutes 注册用户管理相关路由（JWT + RBAC 保护）
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	// 所有用户管理接口都需要 JWT 鉴权和 RBAC 权限控制
	users := group.Group("/users")
	users.Use(middleware.JWTAuth(m.cfg.Secret))
	users.Use(middleware.CasbinRBACAuth(m.db))
	{
		users.GET("", m.handler.ListUsers)
		users.POST("", m.handler.CreateUser)
		users.GET("/:id", m.handler.GetUser)
		users.PUT("/:id", m.handler.UpdateUser)
		users.PATCH("/:id/status", m.handler.UpdateStatus)
		users.POST("/:id/reset-password", m.handler.ResetPassword)
		users.GET("/:id/roles", m.handler.GetUserRoles)
		users.PUT("/:id/roles", m.handler.SetUserRoles)
		users.GET("/:id/permissions", m.handler.GetUserPermissions)
	}
}
