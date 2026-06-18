package handler

import (
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/shelving/service"
	"github.com/gin-gonic/gin"
)

// ShelvingHandler 处理上架相关 HTTP 请求。
type ShelvingHandler struct {
	shelvingSvc service.ShelvingService
}

// NewShelvingHandler 创建上架 HTTP 处理器。
func NewShelvingHandler(shelvingSvc service.ShelvingService) *ShelvingHandler {
	return &ShelvingHandler{shelvingSvc: shelvingSvc}
}

// GetPending 分页查询待上架列表。
// GET /shelving/pending
func (h *ShelvingHandler) GetPending(c *gin.Context) {
	var pageQuery core.PageQuery
	if err := c.ShouldBindQuery(&pageQuery); err != nil {
		pageQuery.Page = 1
		pageQuery.PageSize = 20
	}

	result, err := h.shelvingSvc.GetPendingList(c.Request.Context(), pageQuery.Page, pageQuery.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// shelveRequest 单条上架请求体。
type shelveRequest struct {
	TraceCode    string `json:"trace_code" binding:"required"`
	LocationCode string `json:"location_code" binding:"required"`
}

// Scan 单条追溯码上架（扫码上架）。
// POST /shelving/scan
func (h *ShelvingHandler) Scan(c *gin.Context) {
	var req shelveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID := h.currentOperatorID(c)
	err := h.shelvingSvc.ShelveTrace(
		c.Request.Context(),
		strings.TrimSpace(req.TraceCode),
		strings.TrimSpace(req.LocationCode),
		operatorID,
	)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, gin.H{"trace_code": req.TraceCode, "location_code": req.LocationCode})
}

// batchShelveRequest 批量上架请求体。
type batchShelveRequest struct {
	Items []shelveRequest `json:"items" binding:"required,min=1"`
}

// Batch 批量上架，允许部分成功。
// POST /shelving/batch
func (h *ShelvingHandler) Batch(c *gin.Context) {
	var req batchShelveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID := h.currentOperatorID(c)
	items := make([]service.ShelveItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, service.ShelveItem{
			TraceCode:    strings.TrimSpace(item.TraceCode),
			LocationCode: strings.TrimSpace(item.LocationCode),
		})
	}

	results := h.shelvingSvc.BatchShelve(c.Request.Context(), items, operatorID)

	successCount := 0
	failCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			failCount++
		}
	}

	core.Success(c, gin.H{
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
	})
}

// relocateRequest 货位调整请求体。
type relocateRequest struct {
	TraceCode    string `json:"trace_code" binding:"required"`
	LocationCode string `json:"location_code" binding:"required"`
}

// Relocate 将追溯码移动到新货位。
// POST /shelving/relocate
func (h *ShelvingHandler) Relocate(c *gin.Context) {
	var req relocateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID := h.currentOperatorID(c)
	err := h.shelvingSvc.Relocate(
		c.Request.Context(),
		strings.TrimSpace(req.TraceCode),
		strings.TrimSpace(req.LocationCode),
		operatorID,
	)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, gin.H{"trace_code": req.TraceCode, "location_code": req.LocationCode})
}

// MixCheck 检查货位是否存在混合陈列。
// GET /shelving/mix-check
func (h *ShelvingHandler) MixCheck(c *gin.Context) {
	locationCode := strings.TrimSpace(c.Query("location_code"))
	if locationCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "location_code is required")
		return
	}

	result, err := h.shelvingSvc.MixCheck(c.Request.Context(), locationCode)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, result)
}

// currentOperatorID 从 JWT 上下文中提取当前操作人 ID，如果没有则默认 1。
func (h *ShelvingHandler) currentOperatorID(c *gin.Context) int64 {
	if operatorID, ok := middleware.GetCurrentUserID(c); ok && operatorID > 0 {
		return operatorID
	}
	return 1
}
