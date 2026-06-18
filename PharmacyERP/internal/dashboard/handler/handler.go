// Package handler 实现仪表盘模块的 HTTP 处理层。
package handler

import (
	"github.com/YingmoY/PharmacyERP/internal/dashboard/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 仪表盘 HTTP 处理器。
type Handler struct {
	svc service.DashboardService
	log *zap.Logger
}

// New 创建仪表盘处理器实例。
func New(svc service.DashboardService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// GetOverview 处理 GET /dashboard/overview 请求，返回概览数据。
func (h *Handler) GetOverview(c *gin.Context) {
	data, err := h.svc.GetOverview(c.Request.Context())
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, data)
}

// salesTrendRequest 销售趋势查询参数。
type salesTrendRequest struct {
	// Days 查询最近天数，默认 7 天。
	Days int `form:"days,default=7" binding:"min=1,max=365"`
}

// GetSalesTrend 处理 GET /dashboard/sales-trend 请求，返回销售趋势数据。
func (h *Handler) GetSalesTrend(c *gin.Context) {
	var req salesTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	data, err := h.svc.GetSalesTrend(c.Request.Context(), req.Days)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, data)
}

// topDrugsRequest 销售排行查询参数。
type topDrugsRequest struct {
	// Limit 返回排行条数，默认 10。
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
	// Days 查询最近天数，默认 30 天。
	Days int `form:"days,default=30" binding:"min=1,max=365"`
}

// GetTopDrugs 处理 GET /dashboard/top-drugs 请求，返回药品销售排行。
func (h *Handler) GetTopDrugs(c *gin.Context) {
	var req topDrugsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	data, err := h.svc.GetTopDrugs(c.Request.Context(), req.Limit, req.Days)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, data)
}

// inboundStatsRequest 入库统计查询参数。
type inboundStatsRequest struct {
	// Days 查询最近天数，默认 30 天。
	Days int `form:"days,default=30" binding:"min=1,max=365"`
}

// GetInboundStats 处理 GET /dashboard/inbound-stats 请求，返回入库统计。
func (h *Handler) GetInboundStats(c *gin.Context) {
	var req inboundStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	data, err := h.svc.GetInboundStats(c.Request.Context(), req.Days)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, data)
}

// GetInventoryStats 处理 GET /dashboard/inventory-stats 请求，返回库存状态统计。
func (h *Handler) GetInventoryStats(c *gin.Context) {
	data, err := h.svc.GetInventoryStats(c.Request.Context())
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, data)
}
