package trace

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/trace/handler"
	"github.com/YingmoY/PharmacyERP/internal/trace/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 trace 追溯码查询模块的路由注册入口
type Module struct {
	handler *handler.Handler
	secret  string
}

// NewModule 构建 trace 模块，注入所有依赖
func NewModule(db *gorm.DB, logger *zap.Logger, jwtSecret string) *Module {
	svc := service.NewTraceService(db, logger)
	return &Module{
		handler: handler.New(svc),
		secret:  jwtSecret,
	}
}

// RegisterRoutes 注册追溯码查询模块路由
//
// 路由前缀：/api/v1
// 路由列表：
//   - GET  /trace/:trace_code           查询追溯码库存状态及药品信息
//   - GET  /trace/:trace_code/full-chain 查询追溯码完整操作链（时间倒序）
//   - GET  /trace/:trace_code/logs      分页查询追溯日志
//   - POST /trace/validate              验证追溯码是否可用于销售
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := middleware.JWTAuth(m.secret)

	g := group.Group("/trace", auth)
	{
		// 注意：/validate 必须在 /:trace_code 之前注册，否则 "validate" 会被当作 trace_code
		g.POST("/validate", m.handler.Validate)
		g.GET("/:trace_code", m.handler.GetTraceInfo)
		g.GET("/:trace_code/full-chain", m.handler.GetFullChain)
		g.GET("/:trace_code/logs", m.handler.GetLogs)
	}
}
