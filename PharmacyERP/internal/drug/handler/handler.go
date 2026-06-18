package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/drug/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 药品模块 HTTP 处理器。
type Handler struct {
	svc service.DrugService
	log *zap.Logger
}

// NewHandler 创建药品处理器实例。
func NewHandler(svc service.DrugService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// ListDrugs 查询药品列表。
// GET /drugs?keyword=&status=&is_prescription=&page=1&page_size=20
func (h *Handler) ListDrugs(c *gin.Context) {
	var req service.DrugListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	list, total, err := h.svc.ListDrugs(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, list))
}

// CreateDrug 创建药品。
// POST /drugs
func (h *Handler) CreateDrug(c *gin.Context) {
	var req service.CreateDrugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.CreateDrug(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetDrug 根据 ID 获取药品详情。
// GET /drugs/:id
func (h *Handler) GetDrug(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	dto, err := h.svc.GetDrug(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// UpdateDrug 更新药品信息。
// PUT /drugs/:id
func (h *Handler) UpdateDrug(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var req service.UpdateDrugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.UpdateDrug(c.Request.Context(), id, req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// DeleteDrug 软删除药品。
// DELETE /drugs/:id
func (h *Handler) DeleteDrug(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var operatorID int64
	if uid, ok := middleware.GetCurrentUserID(c); ok {
		operatorID = uid
	}

	if err := h.svc.DeleteDrug(c.Request.Context(), id, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "deleted": true})
}

// UpdateDrugStatus 更新药品状态（启用/停用）。
// PATCH /drugs/:id/status
func (h *Handler) UpdateDrugStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var body struct {
		Status *int8 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	var operatorID int64
	if uid, ok := middleware.GetCurrentUserID(c); ok {
		operatorID = uid
	}

	if err := h.svc.UpdateDrugStatus(c.Request.Context(), id, *body.Status, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "status": *body.Status})
}

// GetInventorySummary 获取药品库存汇总统计。
// GET /drugs/:id/inventory-summary
func (h *Handler) GetInventorySummary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	dto, err := h.svc.GetInventorySummary(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetDrugSaleInfo 获取药品销售相关信息。
// GET /drugs/:id/sale-info
func (h *Handler) GetDrugSaleInfo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	dto, err := h.svc.GetDrugSaleInfo(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetDrugByCode 根据药品编码获取药品。
// GET /drugs/code/:drug_code
func (h *Handler) GetDrugByCode(c *gin.Context) {
	code := c.Param("drug_code")
	if code == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "drug_code is required")
		return
	}

	dto, err := h.svc.GetDrugByCode(c.Request.Context(), code)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// SearchDrugs 药品快速搜索（自动补全）。
// GET /drugs/search?q=&limit=10
func (h *Handler) SearchDrugs(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	list, err := h.svc.SearchDrugs(c.Request.Context(), q, limit)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, list)
}
