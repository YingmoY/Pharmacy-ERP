package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/sales/service"
	"github.com/gin-gonic/gin"
)

// Handler 是销售模块 HTTP 适配层
type Handler struct {
	svc service.SalesService
}

// New 创建销售 Handler
func New(svc service.SalesService) *Handler {
	return &Handler{svc: svc}
}

// ==================== 销售订单接口 ====================

// listOrdersRequest 订单列表查询参数
type listOrdersRequest struct {
	CashierID      *int64 `form:"cashier_id"`
	Status         string `form:"status"`
	OrderNo        string `form:"order_no"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	IsPrescription *bool  `form:"is_prescription"`
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListOrders GET /sales-orders
func (h *Handler) ListOrders(c *gin.Context) {
	var req listOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	orders, total, err := h.svc.ListOrders(c.Request.Context(), service.OrderFilter{
		CashierID:      req.CashierID,
		Status:         req.Status,
		OrderNo:        req.OrderNo,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsPrescription: req.IsPrescription,
		Page:           req.Page,
		PageSize:       req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, orders))
}

// CreateOrder POST /sales-orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	cashierID, ok := middleware.GetCurrentUserID(c)
	if !ok || cashierID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	order, err := h.svc.CreateOrder(c.Request.Context(), req, cashierID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, order)
}

// GetOrder GET /sales-orders/:id
func (h *Handler) GetOrder(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, order)
}

// ==================== 订单明细接口 ====================

// GetItems GET /sales-orders/:id/details
func (h *Handler) GetItems(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	items, err := h.svc.GetItems(c.Request.Context(), orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, items)
}

// AddItem POST /sales-orders/:id/details
func (h *Handler) AddItem(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.CreateOrderItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	item, err := h.svc.AddItem(c.Request.Context(), orderID, req, operatorID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, item)
}

// DeleteItem DELETE /sales-orders/:id/details/:detail_id
func (h *Handler) DeleteItem(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}
	detailID, err := parseID(c, "detail_id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid detail id")
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.DeleteItem(c.Request.Context(), orderID, detailID, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"deleted": true})
}

// ==================== 结算/取消/退款接口 ====================

// Pay POST /sales-orders/:id/pay
func (h *Handler) Pay(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.PayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.Pay(ctx, orderID, req, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	order, err := h.svc.GetOrder(ctx, orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, order)
}

// Cancel POST /sales-orders/:id/cancel
func (h *Handler) Cancel(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.Cancel(ctx, orderID, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	order, err := h.svc.GetOrder(ctx, orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, order)
}

// Refund POST /sales-orders/:id/refund
func (h *Handler) Refund(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.RefundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.Refund(ctx, orderID, req, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	order, err := h.svc.GetOrder(ctx, orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, order)
}

// ==================== 预留/追溯接口 ====================

// GetReservedTraces GET /sales-orders/:id/reserved-traces
func (h *Handler) GetReservedTraces(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	rsvs, err := h.svc.GetReservedTraces(c.Request.Context(), orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, rsvs)
}

// GetReviewRecord GET /sales-orders/:id/review-record
func (h *Handler) GetReviewRecord(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	review, err := h.svc.GetReviewRecord(c.Request.Context(), orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, review)
}

// SubmitReview POST /sales-orders/:id/submit-review
func (h *Handler) SubmitReview(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	submitterID, ok := middleware.GetCurrentUserID(c)
	if !ok || submitterID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	ctx := c.Request.Context()
	if err := h.svc.SubmitReview(ctx, orderID, submitterID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	review, err := h.svc.GetReviewRecord(ctx, orderID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, review)
}

// ReserveTrace POST /sales-orders/:id/reserve-trace
func (h *Handler) ReserveTrace(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.ReserveTraceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	rsv, err := h.svc.ReserveTrace(c.Request.Context(), orderID, req, operatorID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, rsv)
}

// ReleaseReservation POST /sales-orders/:id/release-reservation
func (h *Handler) ReleaseReservation(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.ReleaseReservationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.svc.ReleaseReservation(c.Request.Context(), orderID, req, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"released": true})
}

// ScanVerify POST /sales-orders/:id/scan-verify
func (h *Handler) ScanVerify(c *gin.Context) {
	var req service.ScanVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	result, err := h.svc.ScanVerify(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, result)
}

// MedicarePreview POST /sales-orders/:id/medicare-preview
func (h *Handler) MedicarePreview(c *gin.Context) {
	orderID, err := parseID(c, "id")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid order id")
		return
	}

	var req service.MedicarePreviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	resp, err := h.svc.MedicarePreview(c.Request.Context(), orderID, req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, resp)
}

// ==================== 工具函数 ====================

// parseID 从路径参数中解析 int64 ID
func parseID(c *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(c.Param(param), 10, 64)
}
