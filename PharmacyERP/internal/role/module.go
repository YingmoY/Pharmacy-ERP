package role

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/role/handler"
	"github.com/YingmoY/PharmacyERP/internal/role/repository"
	"github.com/YingmoY/PharmacyERP/internal/role/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 角色权限管理模块
type Module struct {
	handler *handler.Handler
	cfg     *config.JWTConfig
	db      *gorm.DB
}

// NewModule 创建角色管理模块，初始化所有依赖
func NewModule(db *gorm.DB, log *zap.Logger, cfg *config.JWTConfig) *Module {
	// 初始化仓储层
	rRepo := repository.NewRoleRepository(db)

	// 初始化服务层
	svc := service.NewRoleService(rRepo, log)

	// 初始化处理器层
	h := handler.NewHandler(svc, log)

	return &Module{handler: h, cfg: cfg, db: db}
}

// RegisterRoutes 注册角色权限管理相关路由（JWT + RBAC 保护）
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	// JWT + RBAC 保护
	jwtMiddleware := middleware.JWTAuth(m.cfg.Secret)
	rbacMiddleware := middleware.CasbinRBACAuth(m.db)

	// 角色 CRUD 接口
	roles := group.Group("/roles")
	roles.Use(jwtMiddleware, rbacMiddleware)
	{
		roles.GET("", m.handler.ListRoles)
		roles.POST("", m.handler.CreateRole)
		roles.GET("/:id", m.handler.GetRole)
		roles.PUT("/:id", m.handler.UpdateRole)
		roles.DELETE("/:id", m.handler.DeleteRole)
		roles.GET("/:id/permissions", m.handler.GetRolePermissions)
		roles.PUT("/:id/permissions", m.handler.SetRolePermissions)
	}

	// 权限列表查询接口
	perms := group.Group("/permissions")
	perms.Use(jwtMiddleware, rbacMiddleware)
	{
		perms.GET("", m.handler.ListPermissions)
	}
}
