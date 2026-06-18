package handler

import (
	"strconv"
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/task/service"
	"github.com/gin-gonic/gin"
)

// InventoryTaskHandler 处理盘点任务相关 HTTP 请求。
type InventoryTaskHandler struct {
	inventoryTaskSvc service.InventoryTaskService
}

// NewInventoryTaskHandler 创建盘点任务处理器。
func NewInventoryTaskHandler(svc service.InventoryTaskService) *InventoryTaskHandler {
	return &InventoryTaskHandler{inventoryTaskSvc: svc}
}

// parseTaskID 从 URI 参数解析任务 ID。
func parseTaskID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "invalid task id")
		return 0, false
	}
	return id, true
}

// currentOperatorID 从 JWT 中提取当前操作人 ID。
func currentOperatorID(c *gin.Context) int64 {
	if id, ok := middleware.GetCurrentUserID(c); ok && id > 0 {
		return id
	}
	return 1
}

// ListInventoryTasks 分页查询盘点任务列表。
// GET /inventory-tasks
func (h *InventoryTaskHandler) ListInventoryTasks(c *gin.Context) {
	var q core.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page = 1
		q.PageSize = 20
	}

	tasks, total, err := h.inventoryTaskSvc.ListTasks(c.Request.Context(), q.Page, q.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	h.inventoryTaskSvc.EnrichTaskCounts(c.Request.Context(), tasks)
	core.Success(c, core.NewPageResult(total, q.Page, q.PageSize, tasks))
}

// createInventoryTaskReq 创建盘点任务请求体。
type createInventoryTaskReq struct {
	ScopeType  string `json:"scope_type" binding:"required"`
	ScopeValue string `json:"scope_value" binding:"required"`
	Remark     string `json:"remark"`
}

// CreateInventoryTask 创建盘点任务。
// POST /inventory-tasks
func (h *InventoryTaskHandler) CreateInventoryTask(c *gin.Context) {
	var req createInventoryTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	creatorID := currentOperatorID(c)
	task, err := h.inventoryTaskSvc.CreateTask(c.Request.Context(), service.CreateInventoryTaskRequest{
		ScopeType:  strings.ToUpper(req.ScopeType),
		ScopeValue: req.ScopeValue,
		CreatorID:  creatorID,
		Remark:     req.Remark,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// GetInventoryTask 查询单个盘点任务。
// GET /inventory-tasks/:id
func (h *InventoryTaskHandler) GetInventoryTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	task, err := h.inventoryTaskSvc.GetTask(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	if summary, err2 := h.inventoryTaskSvc.GetTaskSummary(c.Request.Context(), taskID); err2 == nil {
		task.ScannedCount = summary.TotalScanned
		task.NormalCount = summary.NormalCount
		task.MisplacedCount = summary.MisplacedCount
		task.UnexpectedCount = summary.UnexpectedCount
		task.LossCandidateCount = summary.LossCandidateCount
	}
	core.Success(c, task)
}

// StartInventoryTask 启动盘点任务。
// POST /inventory-tasks/:id/start
func (h *InventoryTaskHandler) StartInventoryTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.inventoryTaskSvc.StartTask(c.Request.Context(), taskID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// submitInventoryScanReq 盘点扫码请求体。
type submitInventoryScanReq struct {
	TraceCode           string `json:"trace_code" binding:"required"`
	ScannedLocationCode string `json:"scanned_location_code" binding:"required"`
}

// SubmitInventoryScan 提交盘点扫描记录。
// POST /inventory-tasks/:id/scan
func (h *InventoryTaskHandler) SubmitInventoryScan(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	var req submitInventoryScanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	detail, err := h.inventoryTaskSvc.SubmitScan(c.Request.Context(), service.SubmitInventoryScanRequest{
		TaskID:              taskID,
		TraceCode:           strings.TrimSpace(req.TraceCode),
		ScannedLocationCode: strings.TrimSpace(req.ScannedLocationCode),
		OperatorID:          currentOperatorID(c),
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	msg := "正常"
	switch detail.DiffType {
	case "MISPLACED_FOUND":
		msg = "货位不符，已标记为错架"
	case "UNEXPECTED":
		msg = "追溯码不在库或状态异常"
	}
	core.Success(c, gin.H{
		"success":               true,
		"scan_result":           detail.DiffType,
		"trace_code":            req.TraceCode,
		"scanned_location_code": req.ScannedLocationCode,
		"message":               msg,
	})
}

// CompleteInventoryTask 完成盘点任务，返回任务统计摘要。
// POST /inventory-tasks/:id/complete
func (h *InventoryTaskHandler) CompleteInventoryTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	operatorID := currentOperatorID(c)
	if err := h.inventoryTaskSvc.CompleteTask(c.Request.Context(), taskID, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	summary, err := h.inventoryTaskSvc.GetTaskSummary(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, summary)
}

// CancelInventoryTask 取消盘点任务。
// POST /inventory-tasks/:id/cancel
func (h *InventoryTaskHandler) CancelInventoryTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.inventoryTaskSvc.CancelTask(c.Request.Context(), taskID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// assignTaskReq 指定任务执行人请求体。
type assignTaskReq struct {
	AssigneeID int64 `json:"assignee_id" binding:"required"`
}

// AssignInventoryTask 指定盘点任务执行人。
// POST /inventory-tasks/:id/assign
func (h *InventoryTaskHandler) AssignInventoryTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	var req assignTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	if err := h.inventoryTaskSvc.AssignTask(c.Request.Context(), taskID, req.AssigneeID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// GetInventoryTaskDetails 查询盘点任务扫描明细（分页）。
// GET /inventory-tasks/:id/details
func (h *InventoryTaskHandler) GetInventoryTaskDetails(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	page := 1
	pageSize := 20
	if v, err := strconv.Atoi(c.Query("page")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(c.Query("page_size")); err == nil && v > 0 {
		pageSize = v
	}

	details, err := h.inventoryTaskSvc.GetTaskDetails(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	total := int64(len(details))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(details) {
		start = len(details)
	}
	if end > len(details) {
		end = len(details)
	}
	core.Success(c, core.NewPageResult(total, page, pageSize, details[start:end]))
}

// GetInventoryTaskSummary 获取盘点任务统计摘要。
// GET /inventory-tasks/:id/summary
func (h *InventoryTaskHandler) GetInventoryTaskSummary(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	summary, err := h.inventoryTaskSvc.GetTaskSummary(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, summary)
}

// GetLossCandidates 查询盘点任务的盘亏候选列表。
// GET /inventory-tasks/:id/loss-candidates
func (h *InventoryTaskHandler) GetLossCandidates(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	traces, err := h.inventoryTaskSvc.ListLossCandidates(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, gin.H{"total": len(traces), "list": traces})
}

// ConfirmLoss 确认盘亏（LOSS_CANDIDATE -> LOST）。
// POST /inventory-tasks/:id/loss-candidates/:trace_code/confirm
func (h *InventoryTaskHandler) ConfirmLoss(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	operatorID := currentOperatorID(c)
	if err := h.inventoryTaskSvc.ConfirmLoss(c.Request.Context(), taskID, traceCode, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// RejectLoss 拒绝盘亏（LOSS_CANDIDATE -> IN_STOCK）。
// POST /inventory-tasks/:id/loss-candidates/:trace_code/reject
func (h *InventoryTaskHandler) RejectLoss(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	operatorID := currentOperatorID(c)
	if err := h.inventoryTaskSvc.RejectLoss(c.Request.Context(), taskID, traceCode, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// GetMisplaced 查询盘点任务的错架列表。
// GET /inventory-tasks/:id/misplaced
func (h *InventoryTaskHandler) GetMisplaced(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	traces, err := h.inventoryTaskSvc.ListMisplaced(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, gin.H{"total": len(traces), "list": traces})
}

// relocateMisplacedReq 错架移位请求体。
type relocateMisplacedReq struct {
	LocationCode string `json:"location_code" binding:"required"`
}

// RelocateMisplaced 将错架药品移位到正确货位。
// POST /inventory-tasks/:id/misplaced/:trace_code/relocate
func (h *InventoryTaskHandler) RelocateMisplaced(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	var req relocateMisplacedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID := currentOperatorID(c)
	if err := h.inventoryTaskSvc.RelocateMisplaced(c.Request.Context(), taskID, traceCode, req.LocationCode, operatorID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}
