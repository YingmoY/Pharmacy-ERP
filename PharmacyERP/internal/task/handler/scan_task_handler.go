package handler

import (
	"strconv"
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/task/service"
	"github.com/gin-gonic/gin"
)

// ScanTaskHandler 处理扫码任务相关 HTTP 请求。
type ScanTaskHandler struct {
	scanTaskSvc service.ScanTaskService
}

// NewScanTaskHandler 创建扫码任务处理器。
func NewScanTaskHandler(svc service.ScanTaskService) *ScanTaskHandler {
	return &ScanTaskHandler{scanTaskSvc: svc}
}

// ListScanTasks 分页查询扫码任务列表。
// GET /scan-tasks
func (h *ScanTaskHandler) ListScanTasks(c *gin.Context) {
	var q core.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page = 1
		q.PageSize = 20
	}

	tasks, total, err := h.scanTaskSvc.ListScanTasks(c.Request.Context(), q.Page, q.PageSize)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, core.NewPageResult(total, q.Page, q.PageSize, tasks))
}

// createScanTaskReq 创建扫码任务请求体。
type createScanTaskReq struct {
	TaskType  string `json:"task_type" binding:"required"`
	RelatedID int64  `json:"related_id"`
	Remark    string `json:"remark"`
}

// CreateScanTask 创建扫码任务。
// POST /scan-tasks
func (h *ScanTaskHandler) CreateScanTask(c *gin.Context) {
	var req createScanTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	operatorID := currentOperatorID(c)
	task, err := h.scanTaskSvc.CreateScanTask(c.Request.Context(), service.CreateScanTaskRequest{
		TaskType:   strings.ToUpper(req.TaskType),
		RelatedID:  req.RelatedID,
		OperatorID: operatorID,
		Remark:     req.Remark,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// GetScanTask 查询单个扫码任务。
// GET /scan-tasks/:id
func (h *ScanTaskHandler) GetScanTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	task, err := h.scanTaskSvc.GetScanTask(c.Request.Context(), taskID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, task)
}

// StartScanTask 启动扫码任务。
// POST /scan-tasks/:id/start
func (h *ScanTaskHandler) StartScanTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.scanTaskSvc.StartScanTask(c.Request.Context(), taskID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// submitScanReq 提交扫码结果请求体。
type submitScanReq struct {
	TraceCode    string `json:"trace_code" binding:"required"`
	LocationCode string `json:"location_code" binding:"required"`
	DetailID     *int64 `json:"detail_id"`
}

// SubmitScan 提交单条扫码结果。
// POST /scan-tasks/:id/submit
func (h *ScanTaskHandler) SubmitScan(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	var req submitScanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	detail, err := h.scanTaskSvc.SubmitScan(c.Request.Context(), service.SubmitScanTaskRequest{
		TaskID:       taskID,
		TraceCode:    strings.TrimSpace(req.TraceCode),
		LocationCode: strings.TrimSpace(req.LocationCode),
		DetailID:     req.DetailID,
	})
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, detail)
}

// CompleteScanTask 完成扫码任务。
// POST /scan-tasks/:id/complete
func (h *ScanTaskHandler) CompleteScanTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.scanTaskSvc.CompleteScanTask(c.Request.Context(), taskID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// CancelScanTask 取消扫码任务（已提交结果不回滚）。
// POST /scan-tasks/:id/cancel
func (h *ScanTaskHandler) CancelScanTask(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}

	if err := h.scanTaskSvc.CancelScanTask(c.Request.Context(), taskID); err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}
	core.Success(c, nil)
}

// GetScanTaskDetails 查询扫码任务的扫描明细列表（分页）。
// GET /scan-tasks/:id/details
func (h *ScanTaskHandler) GetScanTaskDetails(c *gin.Context) {
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

	details, err := h.scanTaskSvc.GetScanTaskDetails(c.Request.Context(), taskID)
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
