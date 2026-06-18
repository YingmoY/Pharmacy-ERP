package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/pharmacist/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Handler 是药师审核模块 HTTP 适配层
type Handler struct {
	svc service.ReviewService
}

// New 创建药师审核 Handler
func New(svc service.ReviewService) *Handler {
	return &Handler{svc: svc}
}

// listReviewsRequest 审核列表查询参数
type listReviewsRequest struct {
	Status   string `form:"status"`
	OrderNo  string `form:"order_no"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListReviews GET /pharmacist/reviews
func (h *Handler) ListReviews(c *gin.Context) {
	var req listReviewsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	reviews, total, err := h.svc.ListReviews(c.Request.Context(), service.ReviewFilter{
		Status:   req.Status,
		OrderNo:  req.OrderNo,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, reviews))
}

// GetReview GET /pharmacist/reviews/:id
func (h *Handler) GetReview(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid review id")
		return
	}

	review, err := h.svc.GetReview(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, review)
}

// approveRejectRequest 审核通过/驳回请求体
type approveRejectRequest struct {
	ReviewOpinion string `json:"review_opinion"`
}

// ApproveReview POST /pharmacist/reviews/:id/approve
func (h *Handler) ApproveReview(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid review id")
		return
	}

	var req approveRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	pharmacistID, ok := middleware.GetCurrentUserID(c)
	if !ok || pharmacistID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.ApproveReview(ctx, id, pharmacistID, req.ReviewOpinion); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	review, err := h.svc.GetReview(ctx, id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, review)
}

// RejectReview POST /pharmacist/reviews/:id/reject

func (h *Handler) RejectReview(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid review id")
		return
	}

	var req approveRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	pharmacistID, ok := middleware.GetCurrentUserID(c)
	if !ok || pharmacistID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.RejectReview(ctx, id, pharmacistID, req.ReviewOpinion); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	review, err := h.svc.GetReview(ctx, id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, review)
}

// parseID 从路径参数中解析 int64 ID
func parseID(c *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(c.Param(param), 10, 64)
}
