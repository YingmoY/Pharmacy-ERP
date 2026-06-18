package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/inventory/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// InventoryHandler 处理库存查询和调整相关 HTTP 请求。
type InventoryHandler struct {
	inventorySvc service.InventoryService
}

// NewInventoryHandler 创建库存 HTTP 处理器。
func NewInventoryHandler(inventorySvc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventorySvc: inventorySvc}
}

// ListInventory 分页查询追溯库存列表。
// GET /inventory
func (h *InventoryHandler) ListInventory(c *gin.Context) {
	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	req := service.ListInventoryRequest{
		Page:     pageQuery.Page,
		PageSize: pageQuery.PageSize,
	}

	if v := c.Query("status"); v != "" {
		req.Status = strings.ToUpper(v)
	}
	if v := c.Query("batch_number"); v != "" {
		req.BatchNumber = v
	}
	if v := c.Query("drug_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil && id > 0 {
			req.DrugID = &id
		}
	}
	if v := c.Query("location_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil && id > 0 {
			req.LocationID = &id
		}
	}
	if v := c.Query("expire_date_start"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			req.ExpireDateStart = &t
		}
	}
	if v := c.Query("expire_date_end"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			req.ExpireDateEnd = &t
		}
	}

	result, err := h.inventorySvc.ListInventory(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// GetSummary 获取库存统计汇总信息。
// GET /inventory/summary
func (h *InventoryHandler) GetSummary(c *gin.Context) {
	summary, err := h.inventorySvc.GetInventorySummary(c.Request.Context())
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, summary)
}

// ListPendingShelving 查询待上架追溯码列表。
// GET /inventory/pending-shelving
func (h *InventoryHandler) ListPendingShelving(c *gin.Context) {
	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.inventorySvc.ListPendingShelving(c.Request.Context(), pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// ListNearExpire 查询近效期在库药品列表。
// GET /inventory/near-expire
func (h *InventoryHandler) ListNearExpire(c *gin.Context) {
	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.inventorySvc.ListNearExpire(c.Request.Context(), pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// ListRecommendSale 推荐销售追溯码（FEFO先到期先出），drug_id 为必填查询参数。
// GET /inventory/recommend-sale
func (h *InventoryHandler) ListRecommendSale(c *gin.Context) {
	drugIDStr := c.Query("drug_id")
	if drugIDStr == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "drug_id is required")
		return
	}
	drugID, err := strconv.ParseInt(drugIDStr, 10, 64)
	if err != nil || drugID <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "drug_id must be a valid positive integer")
		return
	}

	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.inventorySvc.ListRecommendSale(c.Request.Context(), drugID, pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// ListByDrug 查询某药品的所有追溯码。
// GET /inventory/drugs/:drug_id
func (h *InventoryHandler) ListByDrug(c *gin.Context) {
	drugID, err := strconv.ParseInt(c.Param("drug_id"), 10, 64)
	if err != nil || drugID <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid drug_id")
		return
	}

	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.inventorySvc.ListByDrug(c.Request.Context(), drugID, pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// ListByLocation 查询某货位的所有追溯码。
// GET /inventory/locations/:location_id
func (h *InventoryHandler) ListByLocation(c *gin.Context) {
	locationID, err := strconv.ParseInt(c.Param("location_id"), 10, 64)
	if err != nil || locationID <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid location_id")
		return
	}

	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.inventorySvc.ListByLocation(c.Request.Context(), locationID, pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// ManualStatusChangeRequest 手动状态变更请求体。
type ManualStatusChangeRequest struct {
	ToStatus string `json:"to_status" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
	Remark   string `json:"remark"`
}

// ManualStatusChange 手动变更追溯码状态（仅管理员）。
// PATCH /inventory/:trace_code/status
func (h *InventoryHandler) ManualStatusChange(c *gin.Context) {
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	var req ManualStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	err := h.inventorySvc.ManualStatusChange(c.Request.Context(), service.ManualStatusChangeRequest{
		TraceCode:  traceCode,
		ToStatus:   strings.ToUpper(req.ToStatus),
		Reason:     req.Reason,
		Remark:     req.Remark,
		OperatorID: operatorID,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, gin.H{"trace_code": traceCode, "to_status": req.ToStatus})
}

// ListAdjustmentsRequest 调整记录查询请求体。
type listAdjustmentsQuery struct {
	core.PageQuery
	TraceCode  string `form:"trace_code"`
	DrugID     *int64 `form:"drug_id"`
	AdjustType string `form:"adjust_type"`
}

// ListAdjustments 分页查询库存调整记录。
// GET /inventory-adjustments
func (h *InventoryHandler) ListAdjustments(c *gin.Context) {
	var q listAdjustmentsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page = 1
		q.PageSize = 20
	}

	result, err := h.inventorySvc.ListAdjustments(c.Request.Context(), service.ListAdjustmentsRequest{
		TraceCode:  q.TraceCode,
		DrugID:     q.DrugID,
		AdjustType: q.AdjustType,
		Page:       q.Page,
		PageSize:   q.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// GetAdjustment 查询单条调整记录。
// GET /inventory-adjustments/:id
func (h *InventoryHandler) GetAdjustment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	adj, err := h.inventorySvc.GetAdjustment(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, adj)
}

// createAdjustmentRequest 手动创建调整记录请求体。
type createAdjustmentRequest struct {
	TraceCode      string `json:"trace_code" binding:"required"`
	AdjustType     string `json:"adjust_type" binding:"required"`
	FromLocationID *int64 `json:"from_location_id"`
	ToLocationID   *int64 `json:"to_location_id"`
	Reason         string `json:"reason" binding:"required"`
	Remark         string `json:"remark"`
}

// CreateAdjustment 手动创建调整记录。
// POST /inventory-adjustments
func (h *InventoryHandler) CreateAdjustment(c *gin.Context) {
	var req createAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID, ok := middleware.GetCurrentUserID(c)
	if !ok || operatorID <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	adj, err := h.inventorySvc.CreateAdjustment(c.Request.Context(), service.CreateAdjustmentRequest{
		TraceCode:      req.TraceCode,
		AdjustType:     strings.ToUpper(req.AdjustType),
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Reason:         req.Reason,
		Remark:         req.Remark,
		OperatorID:     operatorID,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, adj)
}
