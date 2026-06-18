package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/supplier/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 供应商模块 HTTP 处理器。
type Handler struct {
	svc service.SupplierService
	log *zap.Logger
}

// NewHandler 创建供应商处理器实例。
func NewHandler(svc service.SupplierService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// ListSuppliers 查询供应商列表。
// GET /suppliers?keyword=&status=&page=1&page_size=20
func (h *Handler) ListSuppliers(c *gin.Context) {
	var req service.SupplierListRequest
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

	list, total, err := h.svc.ListSuppliers(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, list))
}

// CreateSupplier 创建供应商。
// POST /suppliers
func (h *Handler) CreateSupplier(c *gin.Context) {
	var req service.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.CreateSupplier(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetSupplier 根据 ID 获取供应商详情。
// GET /suppliers/:id
func (h *Handler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	dto, err := h.svc.GetSupplier(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// UpdateSupplier 更新供应商信息。
// PUT /suppliers/:id
func (h *Handler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var req service.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.UpdateSupplier(c.Request.Context(), id, req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// DeleteSupplier 软删除供应商。
// DELETE /suppliers/:id
func (h *Handler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var operatorID int64
	if uid, ok := middleware.GetCurrentUserID(c); ok {
		operatorID = uid
	}

	if err := h.svc.DeleteSupplier(c.Request.Context(), id, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "deleted": true})
}

// UpdateSupplierStatus 更新供应商状态。
// PATCH /suppliers/:id/status
func (h *Handler) UpdateSupplierStatus(c *gin.Context) {
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

	if err := h.svc.UpdateSupplierStatus(c.Request.Context(), id, *body.Status, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "status": *body.Status})
}
