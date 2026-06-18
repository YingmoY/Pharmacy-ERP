// Package ai 是 AI 发票识别业务模块，负责路由注册和依赖组装。
package ai

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/ai/handler"
	"github.com/YingmoY/PharmacyERP/internal/ai/repository"
	"github.com/YingmoY/PharmacyERP/internal/ai/service"
	inboundRepo "github.com/YingmoY/PharmacyERP/internal/inbound/repository"
	inboundService "github.com/YingmoY/PharmacyERP/internal/inbound/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 AI 发票业务模块的 HTTP 适配层。
type Module struct {
	h         *handler.Handler
	jwtSecret string
}

// NewModule 组装 AI 发票模块依赖。
func NewModule(db *gorm.DB, logger *zap.Logger, mqClient *mq.Client, jwtSecret string, aiServiceBaseURL string, aiServiceTimeout time.Duration) *Module {
	// 构建入库单仓储与服务（AI 转换时需要创建入库单）。
	ibRepo := inboundRepo.NewInboundRepo(db)
	ibSvc := inboundService.NewInboundService(db, ibRepo, mqClient, logger)

	invoiceRepo := repository.NewAIInvoiceRepo(db)
	svc := service.NewAIInvoiceService(db, invoiceRepo, ibRepo, ibSvc, logger, aiServiceBaseURL, aiServiceTimeout)
	h := handler.New(svc)
	return &Module{h: h, jwtSecret: jwtSecret}
}

// RegisterRoutes 将 AI 发票相关路由挂载到 /api/v1 下。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/ai/invoices", middleware.JWTAuth(m.jwtSecret))
	g.POST("/recognize", m.h.RecognizeInvoice)
	g.GET("", m.h.ListInvoices)
	g.GET("/:id", m.h.GetInvoice)
	g.POST("/:id/convert-to-inbound", m.h.ConvertToInbound)
}
