package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	taskModel "github.com/YingmoY/PharmacyERP/internal/task/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InboundScanProvider 定义入库扫描上架调用接口（避免循环依赖）。
type InboundScanProvider interface {
	// ShelveByLocationCode 根据货位编码对单个追溯码执行上架操作。
	ShelveByLocationCode(ctx context.Context, traceCode, locationCode string, operatorID int64) error
}

// ShelvingScanProvider 定义上架类型扫描调用接口（避免循环依赖）。
type ShelvingScanProvider interface {
	// ShelveTrace 执行追溯码上架操作。
	ShelveTrace(ctx context.Context, traceCode, locationCode string, operatorID int64) error
}

// InventoryScanProvider 定义盘点类型扫描调用接口（避免循环依赖）。
type InventoryScanProvider interface {
	// ScanForInventory 提交盘点扫描，返回 diff_type（NORMAL/MISPLACED_FOUND/UNEXPECTED）。
	ScanForInventory(ctx context.Context, taskID int64, traceCode, locationCode string, operatorID int64) error
}

// ScanTaskService 定义扫码任务业务逻辑接口。
type ScanTaskService interface {
	// CreateScanTask 创建扫码任务。
	CreateScanTask(ctx context.Context, req CreateScanTaskRequest) (*taskModel.ScanTask, error)
	// StartScanTask 将任务从PENDING推进到IN_PROGRESS。
	StartScanTask(ctx context.Context, taskID int64) error
	// SubmitScan 提交单条扫描结果，根据任务类型路由到对应业务服务。
	SubmitScan(ctx context.Context, req SubmitScanTaskRequest) (*taskModel.ScanTaskDetail, error)
	// CompleteScanTask 完成扫码任务。
	CompleteScanTask(ctx context.Context, taskID int64) error
	// CancelScanTask 取消扫码任务（已提交结果不回滚）。
	CancelScanTask(ctx context.Context, taskID int64) error
	// GetScanTask 查询单个扫码任务。
	GetScanTask(ctx context.Context, taskID int64) (*taskModel.ScanTask, error)
	// ListScanTasks 分页查询扫码任务列表。
	ListScanTasks(ctx context.Context, page, pageSize int) ([]taskModel.ScanTask, int64, error)
	// GetScanTaskDetails 查询扫码任务的所有扫描明细。
	GetScanTaskDetails(ctx context.Context, taskID int64) ([]taskModel.ScanTaskDetail, error)
}

// CreateScanTaskRequest 创建扫码任务请求。
type CreateScanTaskRequest struct {
	TaskType   string
	RelatedID  int64
	OperatorID int64
	Remark     string
}

// SubmitScanTaskRequest 提交扫码请求。
type SubmitScanTaskRequest struct {
	TaskID       int64
	TraceCode    string
	LocationCode string
	DetailID     *int64 // 仅 INBOUND 类型使用
}

// scanTaskService 是 ScanTaskService 的默认实现。
type scanTaskService struct {
	db               *gorm.DB
	inboundProvider  InboundScanProvider
	shelvingProvider ShelvingScanProvider
	logger           *zap.Logger
}

// NewScanTaskService 创建扫码任务服务。
// inboundProvider 和 shelvingProvider 用于解耦循环依赖。
func NewScanTaskService(
	db *gorm.DB,
	inboundProvider InboundScanProvider,
	shelvingProvider ShelvingScanProvider,
	logger *zap.Logger,
) ScanTaskService {
	return &scanTaskService{
		db:               db,
		inboundProvider:  inboundProvider,
		shelvingProvider: shelvingProvider,
		logger:           logger,
	}
}

// generateScanTaskNo 生成扫码任务编号，格式：SCAN-YYYYMMDD-XXXX。
func generateScanTaskNo() string {
	now := time.Now()
	seq := now.UnixNano() % 10000
	return fmt.Sprintf("SCAN-%s-%04d", now.Format("20060102"), seq)
}

func (s *scanTaskService) CreateScanTask(ctx context.Context, req CreateScanTaskRequest) (*taskModel.ScanTask, error) {
	taskType := strings.ToUpper(strings.TrimSpace(req.TaskType))
	if taskType != taskModel.ScanTaskTypeInbound &&
		taskType != taskModel.ScanTaskTypeShelving &&
		taskType != taskModel.ScanTaskTypeInventory {
		return nil, ecode.ErrParamInvalid
	}
	if req.OperatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	var remark *string
	if req.Remark != "" {
		remark = &req.Remark
	}

	task := &taskModel.ScanTask{
		TaskNo:     generateScanTaskNo(),
		TaskType:   taskType,
		RelatedID:  req.RelatedID,
		OperatorID: req.OperatorID,
		Status:     taskModel.ScanTaskStatusPending,
		Remark:     remark,
	}

	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		s.logger.Error("创建扫码任务失败", zap.Error(err))
		return nil, err
	}
	return task, nil
}

func (s *scanTaskService) StartScanTask(ctx context.Context, taskID int64) error {
	task, err := s.getScanTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Status != taskModel.ScanTaskStatusPending {
		return ecode.ErrStatusInvalid
	}

	now := time.Now()
	return s.db.WithContext(ctx).
		Model(task).
		Updates(map[string]interface{}{
			"status":     taskModel.ScanTaskStatusInProgress,
			"start_time": now,
		}).Error
}

func (s *scanTaskService) SubmitScan(ctx context.Context, req SubmitScanTaskRequest) (*taskModel.ScanTaskDetail, error) {
	if req.TraceCode == "" || req.LocationCode == "" {
		return nil, ecode.ErrParamInvalid
	}

	task, err := s.getScanTask(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}
	if task.Status != taskModel.ScanTaskStatusInProgress {
		return nil, ecode.ErrStatusInvalid
	}

	// 检查任务内是否重复扫描。
	var existCount int64
	s.db.WithContext(ctx).
		Model(&taskModel.ScanTaskDetail{}).
		Where("task_id = ? AND trace_code = ?", task.ID, req.TraceCode).
		Count(&existCount)
	if existCount > 0 {
		detail := s.saveDetail(ctx, task.ID, req.TraceCode, req.LocationCode, taskModel.ScanResultDuplicate, "任务内重复扫描")
		return detail, nil
	}

	// 查询货位（供各类型业务服务调用）。
	var loc locationModel.LocationInfo
	locErr := s.db.WithContext(ctx).
		Where("location_code = ? AND status = ?", req.LocationCode, locationModel.LocationStatusEnabled).
		First(&loc).Error

	scanResult := taskModel.ScanResultSuccess
	errMsg := ""

	if locErr != nil {
		scanResult = taskModel.ScanResultInvalid
		errMsg = ecode.ErrLocationNotFound.Msg
	} else {
		// 根据任务类型路由到对应业务服务。
		var bizErr error
		switch task.TaskType {
		case taskModel.ScanTaskTypeInbound:
			if s.inboundProvider != nil {
				bizErr = s.inboundProvider.ShelveByLocationCode(ctx, req.TraceCode, req.LocationCode, task.OperatorID)
			}
		case taskModel.ScanTaskTypeShelving:
			if s.shelvingProvider != nil {
				bizErr = s.shelvingProvider.ShelveTrace(ctx, req.TraceCode, req.LocationCode, task.OperatorID)
			}
		case taskModel.ScanTaskTypeInventory:
			// 盘点类型的扫描结果由调用方自行处理，此处仅记录扫描结果。
			// 实际盘点逻辑在 InventoryTaskService.SubmitScan 中执行。
		}

		if bizErr != nil {
			scanResult = taskModel.ScanResultStatusError
			bizErrTyped := ecode.FromError(bizErr)
			errMsg = bizErrTyped.Msg
		}
	}

	detail := s.saveDetail(ctx, task.ID, req.TraceCode, req.LocationCode, scanResult, errMsg)
	return detail, nil
}

// saveDetail 创建并保存扫描明细记录。
func (s *scanTaskService) saveDetail(ctx context.Context, taskID int64, traceCode, locationCode, result, errMsg string) *taskModel.ScanTaskDetail {
	var locCode *string
	if locationCode != "" {
		locCode = &locationCode
	}
	var errMsgPtr *string
	if errMsg != "" {
		errMsgPtr = &errMsg
	}

	detail := &taskModel.ScanTaskDetail{
		TaskID:       taskID,
		TraceCode:    traceCode,
		LocationCode: locCode,
		ScanResult:   result,
		ErrorMsg:     errMsgPtr,
		ScanTime:     time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(detail).Error; err != nil {
		s.logger.Error("保存扫描明细失败", zap.Error(err))
	}
	return detail
}

func (s *scanTaskService) CompleteScanTask(ctx context.Context, taskID int64) error {
	task, err := s.getScanTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Status != taskModel.ScanTaskStatusInProgress {
		return ecode.ErrStatusInvalid
	}

	now := time.Now()
	return s.db.WithContext(ctx).
		Model(task).
		Updates(map[string]interface{}{
			"status":   taskModel.ScanTaskStatusCompleted,
			"end_time": now,
		}).Error
}

func (s *scanTaskService) CancelScanTask(ctx context.Context, taskID int64) error {
	task, err := s.getScanTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Status != taskModel.ScanTaskStatusPending && task.Status != taskModel.ScanTaskStatusInProgress {
		return ecode.ErrStatusInvalid
	}

	// 已提交扫描结果不回滚，直接标记取消。
	return s.db.WithContext(ctx).
		Model(task).
		Update("status", taskModel.ScanTaskStatusCancelled).Error
}

func (s *scanTaskService) GetScanTask(ctx context.Context, taskID int64) (*taskModel.ScanTask, error) {
	return s.getScanTask(ctx, taskID)
}

func (s *scanTaskService) ListScanTasks(ctx context.Context, page, pageSize int) ([]taskModel.ScanTask, int64, error) {
	var tasks []taskModel.ScanTask
	var total int64

	query := s.db.WithContext(ctx).Model(&taskModel.ScanTask{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	s.enrichScanTasks(ctx, tasks)
	return tasks, total, nil
}

// enrichScanTasks 批量填充扫码任务的操作人名称和关联单号。
func (s *scanTaskService) enrichScanTasks(ctx context.Context, tasks []taskModel.ScanTask) {
	if len(tasks) == 0 {
		return
	}

	// 收集操作人 ID。
	opIDs := make([]int64, 0, len(tasks))
	inboundIDs := make([]int64, 0)
	inventoryIDs := make([]int64, 0)
	for _, t := range tasks {
		opIDs = append(opIDs, t.OperatorID)
		switch t.TaskType {
		case taskModel.ScanTaskTypeInbound, taskModel.ScanTaskTypeShelving:
			inboundIDs = append(inboundIDs, t.RelatedID)
		case taskModel.ScanTaskTypeInventory:
			inventoryIDs = append(inventoryIDs, t.RelatedID)
		}
	}

	// 查询操作人名称。
	type userRow struct {
		ID       int64
		RealName string
	}
	var users []userRow
	s.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", opIDs).Scan(&users)
	userMap := make(map[int64]string, len(users))
	for _, u := range users {
		userMap[u.ID] = u.RealName
	}

	// 查询入库单号（INBOUND/SHELVING 类型）。
	inboundOrderNoMap := make(map[int64]string)
	if len(inboundIDs) > 0 {
		type orderRow struct {
			ID      int64
			OrderNo string
		}
		var orders []orderRow
		s.db.WithContext(ctx).Table("inbound_order").Select("id, order_no").Where("id IN ?", inboundIDs).Scan(&orders)
		for _, o := range orders {
			inboundOrderNoMap[o.ID] = o.OrderNo
		}
	}

	// 查询盘库任务号（INVENTORY 类型）。
	inventoryTaskNoMap := make(map[int64]string)
	if len(inventoryIDs) > 0 {
		type taskRow struct {
			ID     int64
			TaskNo string
		}
		var itasks []taskRow
		s.db.WithContext(ctx).Table("inventory_task").Select("id, task_no").Where("id IN ?", inventoryIDs).Scan(&itasks)
		for _, t := range itasks {
			inventoryTaskNoMap[t.ID] = t.TaskNo
		}
	}

	for i := range tasks {
		tasks[i].AssignedToName = userMap[tasks[i].OperatorID]
		switch tasks[i].TaskType {
		case taskModel.ScanTaskTypeInbound, taskModel.ScanTaskTypeShelving:
			tasks[i].RelatedOrderNo = inboundOrderNoMap[tasks[i].RelatedID]
		case taskModel.ScanTaskTypeInventory:
			tasks[i].RelatedOrderNo = inventoryTaskNoMap[tasks[i].RelatedID]
		}
	}
}

func (s *scanTaskService) GetScanTaskDetails(ctx context.Context, taskID int64) ([]taskModel.ScanTaskDetail, error) {
	if _, err := s.getScanTask(ctx, taskID); err != nil {
		return nil, err
	}

	var details []taskModel.ScanTaskDetail
	err := s.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Order("scan_time ASC").
		Find(&details).Error
	return details, err
}

// getScanTask 从数据库查询单个扫码任务（内部辅助方法）。
func (s *scanTaskService) getScanTask(ctx context.Context, taskID int64) (*taskModel.ScanTask, error) {
	var task taskModel.ScanTask
	err := s.db.WithContext(ctx).First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}
