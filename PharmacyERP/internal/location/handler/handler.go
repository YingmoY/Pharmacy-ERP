package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/location/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 货位模块 HTTP 处理器。
type Handler struct {
	svc service.LocationService
	log *zap.Logger
}

// NewHandler 创建货位处理器实例。
func NewHandler(svc service.LocationService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// ListLocations 查询货位列表。
// GET /locations?keyword=&status=&area=&page=1&page_size=20
func (h *Handler) ListLocations(c *gin.Context) {
	var req service.LocationListRequest
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

	list, total, err := h.svc.ListLocations(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, list))
}

// CreateLocation 创建货位。
// POST /locations
func (h *Handler) CreateLocation(c *gin.Context) {
	var req service.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.CreateLocation(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetLocation 根据 ID 获取货位详情。
// GET /locations/:id
func (h *Handler) GetLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	dto, err := h.svc.GetLocation(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// UpdateLocation 更新货位信息。
// PUT /locations/:id
func (h *Handler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var req service.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if uid, ok := middleware.GetCurrentUserID(c); ok {
		req.OperatorID = uid
	}

	dto, err := h.svc.UpdateLocation(c.Request.Context(), id, req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// DeleteLocation 软删除货位。
// DELETE /locations/:id
func (h *Handler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	var operatorID int64
	if uid, ok := middleware.GetCurrentUserID(c); ok {
		operatorID = uid
	}

	if err := h.svc.DeleteLocation(c.Request.Context(), id, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "deleted": true})
}

// UpdateLocationStatus 更新货位状态。
// PATCH /locations/:id/status
func (h *Handler) UpdateLocationStatus(c *gin.Context) {
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

	if err := h.svc.UpdateLocationStatus(c.Request.Context(), id, *body.Status, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, gin.H{"id": id, "status": *body.Status})
}

// GetLocationDrugs 获取货位在库药品列表。
// GET /locations/:id/drugs
func (h *Handler) GetLocationDrugs(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid id")
		return
	}

	list, err := h.svc.GetLocationDrugs(c.Request.Context(), id)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, list)
}

// GetAreas 获取所有不重复区域列表。
// GET /locations/areas
func (h *Handler) GetAreas(c *gin.Context) {
	areas, err := h.svc.GetAreas(c.Request.Context())
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, areas)
}

// GetLocationByCode 根据货位编码获取货位。
// GET /locations/code/:location_code
func (h *Handler) GetLocationByCode(c *gin.Context) {
	code := c.Param("location_code")
	if code == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "location_code is required")
		return
	}

	dto, err := h.svc.GetLocationByCode(c.Request.Context(), code)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}
