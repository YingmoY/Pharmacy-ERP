package auth

import (
	"github.com/YingmoY/PharmacyERP/internal/auth/handler"
	"github.com/YingmoY/PharmacyERP/internal/auth/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	rolerepo "github.com/YingmoY/PharmacyERP/internal/role/repository"
	userrepo "github.com/YingmoY/PharmacyERP/internal/user/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 认证模块，聚合路由注册和依赖初始化
type Module struct {
	handler *handler.Handler
	cfg     *config.JWTConfig
	db      *gorm.DB
}

// NewModule 创建认证模块，初始化所有依赖
func NewModule(db *gorm.DB, log *zap.Logger, mqClient *mq.Client, cfg *config.JWTConfig) *Module {
	// 初始化仓储层
	uRepo := userrepo.NewUserRepository(db)
	rRepo := rolerepo.NewRoleRepository(db)

	// 初始化服务层
	svc := service.NewAuthService(db, uRepo, rRepo, mqClient, cfg, log)

	// 初始化处理器层
	h := handler.NewHandler(svc, log)

	return &Module{handler: h, cfg: cfg, db: db}
}

// RegisterRoutes 注册认证相关路由
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")

	// 公开接口：登录不需要 JWT
	auth.POST("/login", m.handler.Login)

	// 需要 JWT 鉴权的接口
	authRequired := auth.Group("")
	authRequired.Use(middleware.JWTAuth(m.cfg.Secret))
	{
		authRequired.POST("/logout", m.handler.Logout)
		authRequired.GET("/me", m.handler.Me)
		authRequired.PUT("/password", m.handler.ChangePassword)
	}
}
