package sales

import (
	pharmacistRepo "github.com/YingmoY/PharmacyERP/internal/pharmacist/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/medicare"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/sales/handler"
	"github.com/YingmoY/PharmacyERP/internal/sales/repository"
	"github.com/YingmoY/PharmacyERP/internal/sales/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 sales 业务模块的路由注册入口
type Module struct {
	handler *handler.Handler
	secret  string
}

// NewModule 构建 sales 模块，注入所有依赖
func NewModule(db *gorm.DB, logger *zap.Logger, mqClient *mq.Client, medicareClient *medicare.Client, jwtSecret string) *Module {
	salesRepo := repository.NewSalesRepo()
	reviewRepo := pharmacistRepo.NewReviewRepo()

	svc := service.NewSalesService(db, salesRepo, reviewRepo, mqClient, medicareClient, logger)

	return &Module{
		handler: handler.New(svc),
		secret:  jwtSecret,
	}
}

// RegisterRoutes 注册 sales 模块路由
//
// 路由前缀：/api/v1
// 路由列表：
//   - GET    /sales-orders                        查询订单列表
//   - POST   /sales-orders                        创建销售订单
//   - GET    /sales-orders/:id                    查询订单详情
//   - GET    /sales-orders/:id/details            查询订单明细列表
//   - POST   /sales-orders/:id/details            新增订单明细
//   - DELETE /sales-orders/:id/details/:detail_id 删除订单明细
//   - POST   /sales-orders/:id/pay                订单结算
//   - POST   /sales-orders/:id/cancel             取消订单
//   - POST   /sales-orders/:id/refund             订单退款
//   - GET    /sales-orders/:id/reserved-traces    查询预留追溯码列表
//   - GET    /sales-orders/:id/review-record      查询审核记录
//   - POST   /sales-orders/:id/submit-review      手动提交审核
//   - POST   /sales-orders/:id/reserve-trace      手动锁定追溯码
//   - POST   /sales-orders/:id/release-reservation 手动释放预留
//   - POST   /sales-orders/:id/scan-verify        扫码验证
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := middleware.JWTAuth(m.secret)

	g := group.Group("/sales-orders", auth)
	{
		g.GET("", m.handler.ListOrders)
		g.POST("", m.handler.CreateOrder)
		g.GET("/:id", m.handler.GetOrder)

		// 订单明细
		g.GET("/:id/details", m.handler.GetItems)
		g.POST("/:id/details", m.handler.AddItem)
		g.DELETE("/:id/details/:detail_id", m.handler.DeleteItem)

		// 订单操作
		g.POST("/:id/pay", m.handler.Pay)
		g.POST("/:id/cancel", m.handler.Cancel)
		g.POST("/:id/refund", m.handler.Refund)

		// 预留/追溯
		g.GET("/:id/reserved-traces", m.handler.GetReservedTraces)
		g.GET("/:id/review-record", m.handler.GetReviewRecord)
		g.POST("/:id/submit-review", m.handler.SubmitReview)
		g.POST("/:id/reserve-trace", m.handler.ReserveTrace)
		g.POST("/:id/release-reservation", m.handler.ReleaseReservation)
		g.POST("/:id/scan-verify", m.handler.ScanVerify)
		g.POST("/:id/medicare-preview", m.handler.MedicarePreview)
	}
}
