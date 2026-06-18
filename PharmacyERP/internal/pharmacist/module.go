package pharmacist

import (
	"github.com/YingmoY/PharmacyERP/internal/pharmacist/handler"
	"github.com/YingmoY/PharmacyERP/internal/pharmacist/repository"
	"github.com/YingmoY/PharmacyERP/internal/pharmacist/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 pharmacist 业务模块的路由注册入口
type Module struct {
	handler *handler.Handler
	secret  string
}

// NewModule 构建药师审核模块，注入所有依赖
func NewModule(db *gorm.DB, logger *zap.Logger, jwtSecret string) *Module {
	reviewRepo := repository.NewReviewRepo()
	svc := service.NewReviewService(db, reviewRepo, logger)

	return &Module{
		handler: handler.New(svc),
		secret:  jwtSecret,
	}
}

// RegisterRoutes 注册药师审核模块路由
//
// 路由前缀：/api/v1/pharmacist
// 路由列表：
//   - GET  /pharmacist/reviews         查询审核列表（支持 status 过滤）
//   - GET  /pharmacist/reviews/:id     查询审核记录详情
//   - POST /pharmacist/reviews/:id/approve 审核通过
//   - POST /pharmacist/reviews/:id/reject  审核驳回
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := middleware.JWTAuth(m.secret)

	g := group.Group("/pharmacist", auth)
	{
		reviews := g.Group("/reviews")
		{
			reviews.GET("", m.handler.ListReviews)
			reviews.GET("/:id", m.handler.GetReview)
			reviews.POST("/:id/approve", m.handler.ApproveReview)
			reviews.POST("/:id/reject", m.handler.RejectReview)
		}
	}
}
