// Package inbound 是入库单业务模块，负责路由注册和依赖组装。
package inbound

import (
	"github.com/YingmoY/PharmacyERP/internal/inbound/handler"
	"github.com/YingmoY/PharmacyERP/internal/inbound/repository"
	"github.com/YingmoY/PharmacyERP/internal/inbound/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 inbound 业务模块的 HTTP 适配层。
type Module struct {
	h         *handler.Handler
	jwtSecret string
}

// NewModule 组装 inbound 模块依赖。
func NewModule(db *gorm.DB, logger *zap.Logger, mqClient *mq.Client, jwtSecret string) *Module {
	repo := repository.NewInboundRepo(db)
	svc := service.NewInboundService(db, repo, mqClient, logger)
	h := handler.New(svc)
	return &Module{h: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 将入库单相关路由挂载到 /api/v1 下。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/inbound-orders", middleware.JWTAuth(m.jwtSecret))

	// 入库单主体操作
	g.GET("", m.h.ListOrders)
	g.POST("", m.h.CreateOrder)
	g.GET("/:id", m.h.GetOrder)
	g.PUT("/:id", m.h.UpdateOrder)
	g.POST("/:id/submit", m.h.SubmitOrder)
	g.POST("/:id/complete", m.h.CompleteOrder)
	g.POST("/:id/cancel", m.h.CancelOrder)

	// 追溯码确认
	g.POST("/:id/confirm-trace", m.h.ConfirmTrace)
	g.POST("/:id/confirm-traces", m.h.BatchConfirmTrace)

	// 入库进度
	g.GET("/:id/progress", m.h.GetInboundProgress)

	// 明细操作
	g.GET("/:id/details", m.h.GetDetails)
	g.POST("/:id/details", m.h.AddDetail)
	g.GET("/:id/details/:detail_id", m.h.GetDetail)
	g.PUT("/:id/details/:detail_id", m.h.UpdateDetail)
	g.DELETE("/:id/details/:detail_id", m.h.DeleteDetail)
}
