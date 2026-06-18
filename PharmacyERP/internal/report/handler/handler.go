// Package handler 实现报表模块的 HTTP 处理层。
package handler

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/report/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 报表 HTTP 处理器。
type Handler struct {
	svc service.ReportService
	log *zap.Logger
}

// New 创建报表处理器实例。
func New(svc service.ReportService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// salesQueryRequest 销售报表查询参数（Query String）。
type salesQueryRequest struct {
	StartDate     string `form:"start_date"`
	EndDate       string `form:"end_date"`
	CashierID     int64  `form:"cashier_id"`
	PaymentMethod string `form:"payment_method"`
	DrugID        int64  `form:"drug_id"`
	Status        string `form:"status"`
}

// GetSalesReport 处理 GET /reports/sales 请求。
func (h *Handler) GetSalesReport(c *gin.Context) {
	var req salesQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	report, err := h.svc.GetSalesReport(c.Request.Context(), service.SalesFilter{
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		CashierID:     req.CashierID,
		PaymentMethod: req.PaymentMethod,
		DrugID:        req.DrugID,
		Status:        req.Status,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, report)
}

// salesExportRequest 销售报表导出请求体。
type salesExportRequest struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	CashierID     int64  `json:"cashier_id"`
	PaymentMethod string `json:"payment_method"`
	DrugID        int64  `json:"drug_id"`
	Status        string `json:"status"`
}

// ExportSalesReport 处理 POST /reports/sales/export 请求。
func (h *Handler) ExportSalesReport(c *gin.Context) {
	var req salesExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	task, err := h.svc.ExportSalesReport(c.Request.Context(), service.SalesFilter{
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		CashierID:     req.CashierID,
		PaymentMethod: req.PaymentMethod,
		DrugID:        req.DrugID,
		Status:        req.Status,
	}, userID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// inboundQueryRequest 入库报表查询参数。
type inboundQueryRequest struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	SupplierID int64  `form:"supplier_id"`
	DrugID     int64  `form:"drug_id"`
	Status     string `form:"status"`
}

// GetInboundReport 处理 GET /reports/inbound 请求。
func (h *Handler) GetInboundReport(c *gin.Context) {
	var req inboundQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	report, err := h.svc.GetInboundReport(c.Request.Context(), service.InboundFilter{
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		SupplierID: req.SupplierID,
		DrugID:     req.DrugID,
		Status:     req.Status,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, report)
}

// inboundExportRequest 入库报表导出请求体。
type inboundExportRequest struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SupplierID int64  `json:"supplier_id"`
	DrugID     int64  `json:"drug_id"`
	Status     string `json:"status"`
}

// ExportInboundReport 处理 POST /reports/inbound/export 请求。
func (h *Handler) ExportInboundReport(c *gin.Context) {
	var req inboundExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	task, err := h.svc.ExportInboundReport(c.Request.Context(), service.InboundFilter{
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		SupplierID: req.SupplierID,
		DrugID:     req.DrugID,
		Status:     req.Status,
	}, userID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// inventoryQueryRequest 库存报表查询参数。
type inventoryQueryRequest struct {
	DrugID      int64  `form:"drug_id"`
	LocationID  int64  `form:"location_id"`
	Status      string `form:"status"`
	BatchNumber string `form:"batch_number"`
}

// GetInventoryReport 处理 GET /reports/inventory 请求。
func (h *Handler) GetInventoryReport(c *gin.Context) {
	var req inventoryQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	report, err := h.svc.GetInventoryReport(c.Request.Context(), service.InventoryFilter{
		DrugID:      req.DrugID,
		LocationID:  req.LocationID,
		Status:      req.Status,
		BatchNumber: req.BatchNumber,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, report)
}

// inventoryExportRequest 库存报表导出请求体。
type inventoryExportRequest struct {
	DrugID      int64  `json:"drug_id"`
	LocationID  int64  `json:"location_id"`
	Status      string `json:"status"`
	BatchNumber string `json:"batch_number"`
}

// ExportInventoryReport 处理 POST /reports/inventory/export 请求。
func (h *Handler) ExportInventoryReport(c *gin.Context) {
	var req inventoryExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	task, err := h.svc.ExportInventoryReport(c.Request.Context(), service.InventoryFilter{
		DrugID:      req.DrugID,
		LocationID:  req.LocationID,
		Status:      req.Status,
		BatchNumber: req.BatchNumber,
	}, userID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// traceLogQueryRequest 追溯日志报表查询参数。
type traceLogQueryRequest struct {
	TraceCode  string `form:"trace_code"`
	DrugID     int64  `form:"drug_id"`
	ActionType string `form:"action_type"`
	OperatorID int64  `form:"operator_id"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// GetTraceLogReport 处理 GET /reports/trace-log 请求。
func (h *Handler) GetTraceLogReport(c *gin.Context) {
	var req traceLogQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	items, total, err := h.svc.GetTraceLogReport(c.Request.Context(), service.TraceLogFilter{
		TraceCode:  req.TraceCode,
		DrugID:     req.DrugID,
		ActionType: req.ActionType,
		OperatorID: req.OperatorID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, core.NewPageResult(total, req.Page, req.PageSize, items))
}

// traceLogExportRequest 追溯日志报表导出请求体。
type traceLogExportRequest struct {
	TraceCode  string `json:"trace_code"`
	DrugID     int64  `json:"drug_id"`
	ActionType string `json:"action_type"`
	OperatorID int64  `json:"operator_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

// ExportTraceLogReport 处理 POST /reports/trace-log/export 请求。
func (h *Handler) ExportTraceLogReport(c *gin.Context) {
	var req traceLogExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	task, err := h.svc.ExportTraceLogReport(c.Request.Context(), service.TraceLogFilter{
		TraceCode:  req.TraceCode,
		DrugID:     req.DrugID,
		ActionType: req.ActionType,
		OperatorID: req.OperatorID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}, userID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// GetExportTask 处理 GET /reports/export-tasks/:task_id 请求，查询导出任务状态。
func (h *Handler) GetExportTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "task_id 不能为空")
		return
	}

	task, err := h.svc.GetExportTask(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}
