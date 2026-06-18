// Package handler 提供入库单 HTTP 处理层。
package handler

import (
	"strconv"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/inbound/repository"
	"github.com/YingmoY/PharmacyERP/internal/inbound/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Handler 入库单 HTTP 处理器。
type Handler struct {
	svc service.InboundService
}

// New 创建 Handler 实例。
func New(svc service.InboundService) *Handler {
	return &Handler{svc: svc}
}

// ============================================================
// 请求/响应 DTO
// ============================================================

// createOrderRequest 创建入库单请求体。
type createOrderRequest struct {
	SupplierID int64                    `json:"supplier_id" binding:"required"`
	InvoiceNo  string                   `json:"invoice_no"`
	Remark     string                   `json:"remark"`
	Details    []service.CreateDetailReq `json:"details"`
}

// updateOrderRequest 更新入库单请求体。
type updateOrderRequest struct {
	SupplierID int64  `json:"supplier_id"`
	InvoiceNo  string `json:"invoice_no"`
	Remark     string `json:"remark"`
}

// addDetailRequest 新增明细请求体。
type addDetailRequest struct {
	DrugID      int64   `json:"drug_id" binding:"required"`
	BatchNumber string  `json:"batch_number" binding:"required"`
	ExpireDate  string  `json:"expire_date" binding:"required"`
	PlannedQty  int32   `json:"planned_qty" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"min=0"`
	Remark      string  `json:"remark"`
}

// updateDetailRequest 更新明细请求体。
type updateDetailRequest struct {
	BatchNumber string  `json:"batch_number"`
	ExpireDate  string  `json:"expire_date"`
	PlannedQty  int32   `json:"planned_qty" binding:"min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"min=0"`
	Remark      string  `json:"remark"`
}

// confirmTraceRequest 单条追溯码确认请求体。
type confirmTraceRequest struct {
	DetailID  int64  `json:"detail_id" binding:"required"`
	TraceCode string `json:"trace_code" binding:"required"`
}

// batchConfirmTraceRequest 批量追溯码确认请求体。
type batchConfirmTraceRequest struct {
	DetailID   int64    `json:"detail_id" binding:"required"`
	TraceCodes []string `json:"trace_codes" binding:"required,min=1"`
}

// ============================================================
// 辅助方法
// ============================================================

// getOrderID 从 URL 路径参数 ":id" 中解析 int64 ID。
func getOrderID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的入库单 ID")
		return 0, false
	}
	return id, true
}

// getDetailID 从 URL 路径参数 ":detail_id" 中解析 int64 ID。
func getDetailID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("detail_id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的明细 ID")
		return 0, false
	}
	return id, true
}

// requireUserID 从 JWT 上下文中获取当前用户 ID，不存在则响应 401。
func requireUserID(c *gin.Context) (int64, bool) {
	uid, ok := middleware.GetCurrentUserID(c)
	if !ok || uid <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return 0, false
	}
	return uid, true
}

// handleBizErr 统一处理业务错误响应。
func handleBizErr(c *gin.Context, err error) {
	bizErr := ecode.FromError(err)
	core.Fail(c, bizErr.Code, bizErr.Msg)
}

// ============================================================
// 入库单 CRUD
// ============================================================

// CreateOrder POST /inbound-orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}
	creatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	order, err := h.svc.CreateOrder(c.Request.Context(), service.CreateOrderReq{
		SupplierID: req.SupplierID,
		InvoiceNo:  req.InvoiceNo,
		Remark:     req.Remark,
		Details:    req.Details,
	}, creatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, order)
}

// GetOrder GET /inbound-orders/:id
func (h *Handler) GetOrder(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	order, err := h.svc.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, order)
}

// ListOrders GET /inbound-orders
func (h *Handler) ListOrders(c *gin.Context) {
	filter := repository.ListFilter{
		Status:   c.Query("status"),
		Keyword:  c.Query("order_no"),
		Page:     1,
		PageSize: 20,
	}
	if v, err := strconv.ParseInt(c.Query("supplier_id"), 10, 64); err == nil {
		filter.SupplierID = v
	}
	if v, err := strconv.Atoi(c.Query("page")); err == nil && v > 0 {
		filter.Page = v
	}
	if v, err := strconv.Atoi(c.Query("page_size")); err == nil && v > 0 {
		filter.PageSize = v
	}
	if v := c.Query("start_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			filter.StartDate = &t
		}
	}
	if v := c.Query("end_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			// 取当天结束时间。
			end := t.Add(24*time.Hour - time.Nanosecond)
			filter.EndDate = &end
		}
	}

	orders, total, err := h.svc.ListOrders(c.Request.Context(), filter)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, core.NewPageResult(total, filter.Page, filter.PageSize, orders))
}

// UpdateOrder PUT /inbound-orders/:id
func (h *Handler) UpdateOrder(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req updateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}
	err := h.svc.UpdateOrder(c.Request.Context(), orderID, service.UpdateOrderReq{
		SupplierID: req.SupplierID,
		InvoiceNo:  req.InvoiceNo,
		Remark:     req.Remark,
	}, operatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID})
}

// SubmitOrder POST /inbound-orders/:id/submit
func (h *Handler) SubmitOrder(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}
	if err := h.svc.SubmitOrder(c.Request.Context(), orderID, operatorID); err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID, "status": "PENDING_CONFIRM"})
}

// CompleteOrder POST /inbound-orders/:id/complete
func (h *Handler) CompleteOrder(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}
	if err := h.svc.CompleteOrder(c.Request.Context(), orderID, operatorID); err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID, "status": "COMPLETED"})
}

// CancelOrder POST /inbound-orders/:id/cancel
func (h *Handler) CancelOrder(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}
	if err := h.svc.CancelOrder(c.Request.Context(), orderID, operatorID); err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID, "status": "CANCELLED"})
}

// ============================================================
// 追溯码确认
// ============================================================

// ConfirmTrace POST /inbound-orders/:id/confirm-trace（单条）
func (h *Handler) ConfirmTrace(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req confirmTraceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	err := h.svc.ConfirmTraceCodes(c.Request.Context(), service.ConfirmTraceCodesReq{
		OrderID:    orderID,
		DetailID:   req.DetailID,
		TraceCodes: []string{req.TraceCode},
		OperatorID: operatorID,
	})
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID, "detail_id": req.DetailID, "confirmed_count": 1})
}

// BatchConfirmTrace POST /inbound-orders/:id/confirm-traces（批量）
func (h *Handler) BatchConfirmTrace(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req batchConfirmTraceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	err := h.svc.ConfirmTraceCodes(c.Request.Context(), service.ConfirmTraceCodesReq{
		OrderID:    orderID,
		DetailID:   req.DetailID,
		TraceCodes: req.TraceCodes,
		OperatorID: operatorID,
	})
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{
		"order_id":        orderID,
		"detail_id":       req.DetailID,
		"confirmed_count": len(req.TraceCodes),
	})
}

// GetInboundProgress GET /inbound-orders/:id/progress
func (h *Handler) GetInboundProgress(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	progress, err := h.svc.GetInboundProgress(c.Request.Context(), orderID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, progress)
}

// ============================================================
// 明细 CRUD
// ============================================================

// GetDetails GET /inbound-orders/:id/details
func (h *Handler) GetDetails(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	// 通过 GetOrder 同时加载明细。
	order, err := h.svc.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, order.Details)
}

// AddDetail POST /inbound-orders/:id/details
func (h *Handler) AddDetail(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req addDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}
	detail, err := h.svc.AddDetail(c.Request.Context(), orderID, service.CreateDetailReq{
		DrugID:      req.DrugID,
		BatchNumber: req.BatchNumber,
		ExpireDate:  req.ExpireDate,
		PlannedQty:  req.PlannedQty,
		UnitPrice:   req.UnitPrice,
		Remark:      req.Remark,
	}, operatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, detail)
}

// GetDetail GET /inbound-orders/:id/details/:detail_id
func (h *Handler) GetDetail(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	detailID, ok := getDetailID(c)
	if !ok {
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	for _, d := range order.Details {
		if d.ID == detailID {
			core.Success(c, d)
			return
		}
	}
	core.Fail(c, ecode.ErrNotFound.Code, ecode.ErrNotFound.Msg)
}

// UpdateDetail PUT /inbound-orders/:id/details/:detail_id
func (h *Handler) UpdateDetail(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	detailID, ok := getDetailID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req updateDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}
	updated, err := h.svc.UpdateDetail(c.Request.Context(), orderID, detailID, service.UpdateDetailReq{
		BatchNumber: req.BatchNumber,
		ExpireDate:  req.ExpireDate,
		PlannedQty:  req.PlannedQty,
		UnitPrice:   req.UnitPrice,
		Remark:      req.Remark,
	}, operatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, updated)
}

// DeleteDetail DELETE /inbound-orders/:id/details/:detail_id
func (h *Handler) DeleteDetail(c *gin.Context) {
	orderID, ok := getOrderID(c)
	if !ok {
		return
	}
	detailID, ok := getDetailID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	if err := h.svc.DeleteDetail(c.Request.Context(), orderID, detailID, operatorID); err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, gin.H{"order_id": orderID, "detail_id": detailID})
}
