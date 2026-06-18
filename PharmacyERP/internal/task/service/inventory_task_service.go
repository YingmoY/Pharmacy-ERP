package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inventoryModel "github.com/YingmoY/PharmacyERP/internal/inventory/model"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	taskModel "github.com/YingmoY/PharmacyERP/internal/task/model"
	"github.com/YingmoY/PharmacyERP/internal/task/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InventoryTaskService 定义盘点任务业务逻辑接口。
type InventoryTaskService interface {
	// CreateTask 创建盘点任务。
	CreateTask(ctx context.Context, req CreateInventoryTaskRequest) (*taskModel.InventoryTask, error)
	// StartTask 将任务从PENDING推进到IN_PROGRESS。
	StartTask(ctx context.Context, taskID int64) error
	// SubmitScan 提交单条追溯码扫描结果。
	SubmitScan(ctx context.Context, req SubmitInventoryScanRequest) (*taskModel.InventoryTaskDetail, error)
	// CompleteTask 完成任务：标记范围内未扫描在库药品为LOSS_CANDIDATE。
	CompleteTask(ctx context.Context, taskID int64, operatorID int64) error
	// CancelTask 取消任务（PENDING可直接取消；IN_PROGRESS须无MISPLACED/LOSS_CANDIDATE产生）。
	CancelTask(ctx context.Context, taskID int64) error
	// AssignTask 指定任务执行人。
	AssignTask(ctx context.Context, taskID int64, assigneeID int64) error
	// GetTask 查询单个任务。
	GetTask(ctx context.Context, taskID int64) (*taskModel.InventoryTask, error)
	// ListTasks 分页查询任务列表。
	ListTasks(ctx context.Context, page, pageSize int) ([]taskModel.InventoryTask, int64, error)
	// GetTaskDetails 查询任务的全部扫描明细。
	GetTaskDetails(ctx context.Context, taskID int64) ([]taskModel.InventoryTaskDetail, error)
	// GetTaskSummary 统计任务盘点结果摘要。
	GetTaskSummary(ctx context.Context, taskID int64) (*taskModel.InventoryTaskSummary, error)
	// ListLossCandidates 列出由本次任务产生的盘亏候选追溯码。
	ListLossCandidates(ctx context.Context, taskID int64) ([]inventoryModel.TraceInventory, error)
	// ConfirmLoss 确认盘亏：LOSS_CANDIDATE -> LOST。
	ConfirmLoss(ctx context.Context, taskID int64, traceCode string, operatorID int64) error
	// RejectLoss 拒绝盘亏：LOSS_CANDIDATE -> IN_STOCK。
	RejectLoss(ctx context.Context, taskID int64, traceCode string, operatorID int64) error
	// ListMisplaced 列出由本次任务产生的错架追溯码。
	ListMisplaced(ctx context.Context, taskID int64) ([]inventoryModel.TraceInventory, error)
	// RelocateMisplaced 移位错架追溯码：MISPLACED -> IN_STOCK，更新货位。
	RelocateMisplaced(ctx context.Context, taskID int64, traceCode string, locationCode string, operatorID int64) error
	// EnrichTaskCounts 批量填充任务列表的已扫码数量。
	EnrichTaskCounts(ctx context.Context, tasks []taskModel.InventoryTask)
}

// CreateInventoryTaskRequest 创建盘点任务请求。
type CreateInventoryTaskRequest struct {
	ScopeType  string
	ScopeValue string
	CreatorID  int64
	Remark     string
}

// SubmitInventoryScanRequest 提交盘点扫描请求。
type SubmitInventoryScanRequest struct {
	TaskID              int64
	TraceCode           string
	ScannedLocationCode string
	OperatorID          int64
}

// inventoryTaskService 是 InventoryTaskService 的默认实现。
type inventoryTaskService struct {
	db       *gorm.DB
	taskRepo repository.InventoryTask
	logger   *zap.Logger
}

// NewInventoryTaskService 创建盘点任务服务。
func NewInventoryTaskService(
	db *gorm.DB,
	taskRepo repository.InventoryTask,
	logger *zap.Logger,
) InventoryTaskService {
	return &inventoryTaskService{
		db:       db,
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// generateTaskNo 生成盘点任务编号，格式：INV-YYYYMMDD-XXXX（流水号用纳秒后4位保证唯一）。
func generateTaskNo() string {
	now := time.Now()
	seq := now.UnixNano() % 10000
	return fmt.Sprintf("INV-%s-%04d", now.Format("20060102"), seq)
}

func (s *inventoryTaskService) CreateTask(ctx context.Context, req CreateInventoryTaskRequest) (*taskModel.InventoryTask, error) {
	if req.ScopeType == "" || req.ScopeValue == "" || req.CreatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}
	scopeType := strings.ToUpper(req.ScopeType)
	if scopeType != taskModel.InventoryTaskScopeArea &&
		scopeType != taskModel.InventoryTaskScopeShelf &&
		scopeType != taskModel.InventoryTaskScopeLocation {
		return nil, ecode.ErrParamInvalid
	}

	var remark *string
	if req.Remark != "" {
		remark = &req.Remark
	}

	task := &taskModel.InventoryTask{
		TaskNo:     generateTaskNo(),
		ScopeType:  scopeType,
		ScopeValue: req.ScopeValue,
		CreatorID:  req.CreatorID,
		Status:     taskModel.InventoryTaskStatusPending,
		Remark:     remark,
	}

	if err := s.taskRepo.CreateTask(ctx, s.db, task); err != nil {
		s.logger.Error("创建盘点任务失败", zap.Error(err))
		return nil, err
	}
	return task, nil
}

func (s *inventoryTaskService) StartTask(ctx context.Context, taskID int64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrNotFound
		}
		return err
	}
	if task.Status != taskModel.InventoryTaskStatusPending {
		return ecode.ErrStatusInvalid
	}

	now := time.Now()
	task.Status = taskModel.InventoryTaskStatusInProgress
	task.StartTime = &now
	return s.taskRepo.UpdateTask(ctx, s.db, task)
}

func (s *inventoryTaskService) SubmitScan(ctx context.Context, req SubmitInventoryScanRequest) (*taskModel.InventoryTaskDetail, error) {
	if req.TraceCode == "" || req.ScannedLocationCode == "" || req.OperatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	var detail *taskModel.InventoryTaskDetail

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) 校验任务存在且处于进行中。
		task, err := s.taskRepo.GetTaskByID(ctx, tx, req.TaskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if task.Status != taskModel.InventoryTaskStatusInProgress {
			return ecode.ErrStatusInvalid
		}

		// 2) 检查任务内是否已扫过该追溯码（去重）。
		existing, err := s.taskRepo.GetDetailByTraceCode(ctx, tx, req.TaskID, req.TraceCode)
		if err != nil {
			return err
		}
		if existing != nil {
			return ecode.ErrDuplicateScan
		}

		// 3) 查询扫描货位。
		var scannedLoc locationModel.LocationInfo
		err = tx.WithContext(ctx).
			Where("location_code = ? AND status = ?", req.ScannedLocationCode, locationModel.LocationStatusEnabled).
			First(&scannedLoc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrLocationNotFound
			}
			return err
		}

		// 4) 查询追溯码的系统库存记录。
		var trace inventoryModel.TraceInventory
		err = tx.WithContext(ctx).
			Where("trace_code = ?", req.TraceCode).
			First(&trace).Error

		diffType := taskModel.InventoryDiscrepancyNormal
		var systemLocationID *int64

		if err != nil {
			// 追溯码不存在，属于意外扫描。
			diffType = taskModel.InventoryDiscrepancyUnexpected
		} else {
			systemLocationID = trace.LocationID

			if trace.Status != inventoryModel.TraceInventoryStatusInStock {
				// 非在库状态，属于意外扫描。
				diffType = taskModel.InventoryDiscrepancyUnexpected
			} else {
				// 检查追溯码是否在任务范围内，以及货位是否匹配。
				inScope := s.isTraceInScope(ctx, tx, task, &trace, scannedLoc.ID)
				locationMatch := trace.LocationID != nil && *trace.LocationID == scannedLoc.ID

				if inScope && locationMatch {
					diffType = taskModel.InventoryDiscrepancyNormal
				} else {
					// 货位不匹配，标记为错架找回，立即更新状态。
					diffType = taskModel.InventoryDiscrepancyMisplacedFound

					fromStatus := trace.Status
					toStatus := inventoryModel.TraceInventoryStatusMisplaced
					relatedNo := task.TaskNo
					remark := "盘点发现追溯码位置不符，标记为错架"

					// 更新追溯码状态为MISPLACED（不更新货位，记录在扫描位置）。
					if err := tx.Model(&inventoryModel.TraceInventory{}).
						Where("trace_code = ?", req.TraceCode).
						Update("status", toStatus).Error; err != nil {
						return err
					}

					// 写入追溯日志。
					log := &inventoryModel.DrugTraceLog{
						TraceCode:      req.TraceCode,
						ActionType:     "INVENTORY",
						FromStatus:     &fromStatus,
						ToStatus:       &toStatus,
						FromLocationID: &scannedLoc.ID,
						OperatorID:     req.OperatorID,
						RelatedNo:      &relatedNo,
						Remark:         &remark,
					}
					if err := tx.Create(log).Error; err != nil {
						return err
					}
				}
			}
		}

		// 5) 写入扫描明细。
		now := time.Now()
		detail = &taskModel.InventoryTaskDetail{
			TaskID:            req.TaskID,
			TraceCode:         req.TraceCode,
			LocationID:        scannedLoc.ID,
			ScannedLocationID: &scannedLoc.ID,
			SystemLocationID:  systemLocationID,
			DiffType:          diffType,
			OperatorID:        &req.OperatorID,
			ScanTime:          &now,
		}
		return s.taskRepo.CreateDetail(ctx, tx, detail)
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}

// isTraceInScope 判断追溯码所在货位是否在任务的盘点范围内。
func (s *inventoryTaskService) isTraceInScope(
	ctx context.Context,
	tx *gorm.DB,
	task *taskModel.InventoryTask,
	trace *inventoryModel.TraceInventory,
	scannedLocationID int64,
) bool {
	if trace.LocationID == nil {
		return false
	}

	switch task.ScopeType {
	case taskModel.InventoryTaskScopeLocation:
		// 范围值是货位 ID 字符串。
		return fmt.Sprintf("%d", *trace.LocationID) == task.ScopeValue
	case taskModel.InventoryTaskScopeArea:
		// 范围值是区域名称，通过 location_info.area 匹配。
		var loc locationModel.LocationInfo
		err := tx.WithContext(ctx).
			Where("id = ? AND area = ?", *trace.LocationID, task.ScopeValue).
			First(&loc).Error
		return err == nil
	case taskModel.InventoryTaskScopeShelf:
		// 范围值是货架前缀，通过 location_code LIKE 'shelf%' 匹配。
		var loc locationModel.LocationInfo
		err := tx.WithContext(ctx).
			Where("id = ? AND location_code LIKE ?", *trace.LocationID, task.ScopeValue+"%").
			First(&loc).Error
		return err == nil
	}
	return false
}

func (s *inventoryTaskService) CompleteTask(ctx context.Context, taskID int64, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		task, err := s.taskRepo.GetTaskByID(ctx, tx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if task.Status != taskModel.InventoryTaskStatusInProgress {
			return ecode.ErrStatusInvalid
		}

		// 查询本任务已扫描的所有追溯码。
		details, err := s.taskRepo.ListDetails(ctx, tx, taskID)
		if err != nil {
			return err
		}
		scannedCodes := make(map[string]struct{}, len(details))
		for _, d := range details {
			scannedCodes[d.TraceCode] = struct{}{}
		}

		// 查询范围内所有在库追溯码（未扫描到的将成为盘亏候选）。
		var inScopeTraces []inventoryModel.TraceInventory
		query := tx.WithContext(ctx).
			Model(&inventoryModel.TraceInventory{}).
			Joins("JOIN location_info ON location_info.id = drug_trace_inventory.location_id AND location_info.deleted_at IS NULL").
			Where("drug_trace_inventory.status = ?", inventoryModel.TraceInventoryStatusInStock)

		switch task.ScopeType {
		case taskModel.InventoryTaskScopeLocation:
			query = query.Where("drug_trace_inventory.location_id = ?", task.ScopeValue)
		case taskModel.InventoryTaskScopeArea:
			query = query.Where("location_info.area = ?", task.ScopeValue)
		case taskModel.InventoryTaskScopeShelf:
			query = query.Where("location_info.location_code LIKE ?", task.ScopeValue+"%")
		}

		if err := query.Find(&inScopeTraces).Error; err != nil {
			return err
		}

		// 找出未被扫描的在库追溯码，标记为 LOSS_CANDIDATE。
		var lossCandidateCodes []string
		for _, trace := range inScopeTraces {
			if _, scanned := scannedCodes[trace.TraceCode]; !scanned {
				lossCandidateCodes = append(lossCandidateCodes, trace.TraceCode)
			}
		}

		if len(lossCandidateCodes) > 0 {
			if err := tx.Model(&inventoryModel.TraceInventory{}).
				Where("trace_code IN ?", lossCandidateCodes).
				Update("status", inventoryModel.TraceInventoryStatusLossCandidate).Error; err != nil {
				return err
			}

			// 逐条写入追溯日志。
			toStatus := inventoryModel.TraceInventoryStatusLossCandidate
			relatedNo := task.TaskNo
			remark := "盘点完成，追溯码未被扫描，标记为盘亏候选"
			for _, trace := range inScopeTraces {
				if _, scanned := scannedCodes[trace.TraceCode]; scanned {
					continue
				}
				fromStatus := trace.Status
				log := &inventoryModel.DrugTraceLog{
					TraceCode:  trace.TraceCode,
					ActionType: "INVENTORY",
					FromStatus: &fromStatus,
					ToStatus:   &toStatus,
					OperatorID: operatorID,
					RelatedNo:  &relatedNo,
					Remark:     &remark,
				}
				if trace.LocationID != nil {
					log.FromLocationID = trace.LocationID
				}
				if err := tx.Create(log).Error; err != nil {
					return err
				}
			}
		}

		// 更新任务状态为已完成。
		now := time.Now()
		task.Status = taskModel.InventoryTaskStatusCompleted
		task.EndTime = &now
		return s.taskRepo.UpdateTask(ctx, tx, task)
	})
}

func (s *inventoryTaskService) CancelTask(ctx context.Context, taskID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		task, err := s.taskRepo.GetTaskByID(ctx, tx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		switch task.Status {
		case taskModel.InventoryTaskStatusPending:
			// PENDING 可以直接取消。
		case taskModel.InventoryTaskStatusInProgress:
			// IN_PROGRESS 须检查是否已产生 MISPLACED 或 LOSS_CANDIDATE 记录。
			details, err := s.taskRepo.ListDetails(ctx, tx, taskID)
			if err != nil {
				return err
			}
			for _, d := range details {
				if d.DiffType == taskModel.InventoryDiscrepancyMisplacedFound {
					return ecode.New(20007, "任务已产生错架记录，无法取消")
				}
			}
			// 检查是否有 LOSS_CANDIDATE 状态（CompleteTask 阶段产生，此处理论上不存在）
		default:
			return ecode.ErrStatusInvalid
		}

		task.Status = taskModel.InventoryTaskStatusCancelled
		return s.taskRepo.UpdateTask(ctx, tx, task)
	})
}

func (s *inventoryTaskService) AssignTask(ctx context.Context, taskID int64, assigneeID int64) error {
	if assigneeID <= 0 {
		return ecode.ErrParamInvalid
	}
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrNotFound
		}
		return err
	}
	task.AssigneeID = &assigneeID
	return s.taskRepo.UpdateTask(ctx, s.db, task)
}

func (s *inventoryTaskService) GetTask(ctx context.Context, taskID int64) (*taskModel.InventoryTask, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

func (s *inventoryTaskService) ListTasks(ctx context.Context, page, pageSize int) ([]taskModel.InventoryTask, int64, error) {
	return s.taskRepo.ListTasks(ctx, s.db, page, pageSize)
}

func (s *inventoryTaskService) GetTaskDetails(ctx context.Context, taskID int64) ([]taskModel.InventoryTaskDetail, error) {
	if _, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return s.taskRepo.ListDetails(ctx, s.db, taskID)
}

func (s *inventoryTaskService) GetTaskSummary(ctx context.Context, taskID int64) (*taskModel.InventoryTaskSummary, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	details, err := s.taskRepo.ListDetails(ctx, s.db, taskID)
	if err != nil {
		return nil, err
	}

	summary := &taskModel.InventoryTaskSummary{
		TaskID:     task.ID,
		TaskNo:     task.TaskNo,
		ScopeType:  task.ScopeType,
		ScopeValue: task.ScopeValue,
		Status:     task.Status,
	}
	summary.TotalScanned = int64(len(details))
	for _, d := range details {
		switch d.DiffType {
		case taskModel.InventoryDiscrepancyNormal:
			summary.NormalCount++
		case taskModel.InventoryDiscrepancyMisplacedFound:
			summary.MisplacedCount++
		case taskModel.InventoryDiscrepancyUnexpected:
			summary.UnexpectedCount++
		}
	}

	// 统计本任务产生的盘亏候选数量（通过关联追溯日志）。
	var lossCandidateCount int64
	s.db.WithContext(ctx).
		Table("drug_trace_log").
		Where("related_no = ? AND to_status = ?", task.TaskNo, inventoryModel.TraceInventoryStatusLossCandidate).
		Count(&lossCandidateCount)
	summary.LossCandidateCount = lossCandidateCount

	return summary, nil
}

func (s *inventoryTaskService) ListLossCandidates(ctx context.Context, taskID int64) ([]inventoryModel.TraceInventory, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	// 从追溯日志中找出本任务产生的盘亏候选追溯码。
	var traceCodes []string
	err = s.db.WithContext(ctx).
		Table("drug_trace_log").
		Select("trace_code").
		Where("related_no = ? AND to_status = ? AND deleted_at IS NULL", task.TaskNo, inventoryModel.TraceInventoryStatusLossCandidate).
		Pluck("trace_code", &traceCodes).Error
	if err != nil {
		return nil, err
	}
	if len(traceCodes) == 0 {
		return []inventoryModel.TraceInventory{}, nil
	}

	// 查询当前仍为 LOSS_CANDIDATE 的记录。
	var traces []inventoryModel.TraceInventory
	err = s.db.WithContext(ctx).
		Where("trace_code IN ? AND status = ?", traceCodes, inventoryModel.TraceInventoryStatusLossCandidate).
		Find(&traces).Error
	if err == nil {
		s.enrichTraceInventory(ctx, traces)
	}
	return traces, err
}

func (s *inventoryTaskService) ConfirmLoss(ctx context.Context, taskID int64, traceCode string, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		task, err := s.taskRepo.GetTaskByID(ctx, tx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		var trace inventoryModel.TraceInventory
		if err := tx.Where("trace_code = ?", traceCode).First(&trace).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}
		if trace.Status != inventoryModel.TraceInventoryStatusLossCandidate {
			return ecode.ErrStatusInvalid
		}

		fromStatus := trace.Status
		toStatus := inventoryModel.TraceInventoryStatusLost

		// 更新状态。
		if err := tx.Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", traceCode).
			Update("status", toStatus).Error; err != nil {
			return err
		}

		// 写入库存调整记录。
		reason := "盘点确认盘亏"
		adj := &inventoryModel.InventoryAdjustment{
			TraceCode:    traceCode,
			DrugID:       trace.DrugID,
			AdjustType:   inventoryModel.AdjustTypeLoss,
			BeforeStatus: &fromStatus,
			AfterStatus:  &toStatus,
			Reason:       reason,
			OperatorID:   operatorID,
		}
		if err := tx.Create(adj).Error; err != nil {
			return err
		}

		// 写入追溯日志。
		relatedNo := task.TaskNo
		remark := "盘点确认盘亏"
		log := &inventoryModel.DrugTraceLog{
			TraceCode:  traceCode,
			ActionType: "LOSS",
			FromStatus: &fromStatus,
			ToStatus:   &toStatus,
			OperatorID: operatorID,
			RelatedNo:  &relatedNo,
			Remark:     &remark,
		}
		return tx.Create(log).Error
	})
}

func (s *inventoryTaskService) RejectLoss(ctx context.Context, taskID int64, traceCode string, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.taskRepo.GetTaskByID(ctx, tx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		var trace inventoryModel.TraceInventory
		if err := tx.Where("trace_code = ?", traceCode).First(&trace).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}
		if trace.Status != inventoryModel.TraceInventoryStatusLossCandidate {
			return ecode.ErrStatusInvalid
		}

		fromStatus := trace.Status
		toStatus := inventoryModel.TraceInventoryStatusInStock

		// 恢复为在库状态。
		if err := tx.Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", traceCode).
			Update("status", toStatus).Error; err != nil {
			return err
		}

		// 写入追溯日志。
		remark := "拒绝盘亏，恢复在库状态"
		log := &inventoryModel.DrugTraceLog{
			TraceCode:  traceCode,
			ActionType: "INVENTORY",
			FromStatus: &fromStatus,
			ToStatus:   &toStatus,
			OperatorID: operatorID,
			Remark:     &remark,
		}
		return tx.Create(log).Error
	})
}

func (s *inventoryTaskService) ListMisplaced(ctx context.Context, taskID int64) ([]inventoryModel.TraceInventory, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, s.db, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	// 从追溯日志中找出本任务产生的错架追溯码。
	var traceCodes []string
	err = s.db.WithContext(ctx).
		Table("drug_trace_log").
		Select("trace_code").
		Where("related_no = ? AND to_status = ? AND deleted_at IS NULL", task.TaskNo, inventoryModel.TraceInventoryStatusMisplaced).
		Pluck("trace_code", &traceCodes).Error
	if err != nil {
		return nil, err
	}
	if len(traceCodes) == 0 {
		return []inventoryModel.TraceInventory{}, nil
	}

	// 查询当前仍为 MISPLACED 的记录。
	var traces []inventoryModel.TraceInventory
	err = s.db.WithContext(ctx).
		Where("trace_code IN ? AND status = ?", traceCodes, inventoryModel.TraceInventoryStatusMisplaced).
		Find(&traces).Error
	if err != nil {
		return nil, err
	}
	s.enrichTraceInventory(ctx, traces)
	s.enrichScannedLocation(ctx, taskID, traces)
	return traces, nil
}

// enrichScannedLocation 从盘点明细表批量填充错架记录的实际扫描货位编码。
func (s *inventoryTaskService) enrichScannedLocation(ctx context.Context, taskID int64, traces []inventoryModel.TraceInventory) {
	if len(traces) == 0 {
		return
	}
	codes := make([]string, 0, len(traces))
	for _, t := range traces {
		codes = append(codes, t.TraceCode)
	}

	// 从盘点明细取扫描货位ID（每个追溯码可能有多次扫描，取最新一次）。
	type detailRow struct {
		TraceCode         string
		ScannedLocationID int64
	}
	var rows []detailRow
	s.db.WithContext(ctx).
		Table("inventory_task_detail").
		Select("trace_code, scanned_location_id").
		Where("task_id = ? AND trace_code IN ? AND diff_type = ?", taskID, codes, taskModel.InventoryDiscrepancyMisplacedFound).
		Order("scan_time DESC").
		Scan(&rows)

	// 每个 trace_code 取第一条（最新）。
	scannedLocIDMap := make(map[string]int64, len(rows))
	for _, r := range rows {
		if _, exists := scannedLocIDMap[r.TraceCode]; !exists {
			scannedLocIDMap[r.TraceCode] = r.ScannedLocationID
		}
	}

	// 收集货位ID并批量查询货位编码。
	locIDs := make([]int64, 0, len(scannedLocIDMap))
	for _, id := range scannedLocIDMap {
		locIDs = append(locIDs, id)
	}
	locCodeMap := make(map[int64]string)
	if len(locIDs) > 0 {
		type locRow struct {
			ID           int64
			LocationCode string
		}
		var locs []locRow
		s.db.WithContext(ctx).Table("location_info").Select("id, location_code").Where("id IN ?", locIDs).Scan(&locs)
		for _, l := range locs {
			locCodeMap[l.ID] = l.LocationCode
		}
	}

	for i := range traces {
		if locID, ok := scannedLocIDMap[traces[i].TraceCode]; ok {
			traces[i].ScannedLocationCode = locCodeMap[locID]
		}
	}
}

// enrichTraceInventory 批量填充追溯库存的药品名称、规格和货位编码。
func (s *inventoryTaskService) enrichTraceInventory(ctx context.Context, traces []inventoryModel.TraceInventory) {
	if len(traces) == 0 {
		return
	}
	drugIDs := make([]int64, 0, len(traces))
	locIDs := make([]int64, 0)
	for i := range traces {
		drugIDs = append(drugIDs, traces[i].DrugID)
		if traces[i].LocationID != nil {
			locIDs = append(locIDs, *traces[i].LocationID)
		}
	}

	type drugRow struct {
		ID            int64
		CommonName    string
		Specification string
	}
	var drugs []drugRow
	s.db.WithContext(ctx).Table("drug_info").Select("id, common_name, specification").Where("id IN ?", drugIDs).Scan(&drugs)
	drugMap := make(map[int64]drugRow, len(drugs))
	for _, d := range drugs {
		drugMap[d.ID] = d
	}

	locMap := make(map[int64]string)
	if len(locIDs) > 0 {
		type locRow struct {
			ID           int64
			LocationCode string
		}
		var locs []locRow
		s.db.WithContext(ctx).Table("location_info").Select("id, location_code").Where("id IN ?", locIDs).Scan(&locs)
		for _, l := range locs {
			locMap[l.ID] = l.LocationCode
		}
	}

	for i := range traces {
		if d, ok := drugMap[traces[i].DrugID]; ok {
			traces[i].DrugName = d.CommonName
			traces[i].Specification = d.Specification
		}
		if traces[i].LocationID != nil {
			code := locMap[*traces[i].LocationID]
			traces[i].LocationCode = code
			traces[i].SystemLocationCode = code
		}
	}
}

func (s *inventoryTaskService) RelocateMisplaced(ctx context.Context, taskID int64, traceCode string, locationCode string, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.taskRepo.GetTaskByID(ctx, tx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		// 查询目标货位。
		var loc locationModel.LocationInfo
		err = tx.WithContext(ctx).
			Where("location_code = ? AND status = ?", locationCode, locationModel.LocationStatusEnabled).
			First(&loc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrLocationNotFound
			}
			return err
		}

		// 查询追溯码当前状态。
		var trace inventoryModel.TraceInventory
		if err := tx.Where("trace_code = ?", traceCode).First(&trace).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}
		if trace.Status != inventoryModel.TraceInventoryStatusMisplaced {
			return ecode.ErrStatusInvalid
		}

		fromStatus := trace.Status
		toStatus := inventoryModel.TraceInventoryStatusInStock
		fromLocationID := trace.LocationID
		toLocationID := loc.ID

		// 更新状态和货位。
		if err := tx.Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", traceCode).
			Updates(map[string]interface{}{
				"status":      toStatus,
				"location_id": toLocationID,
			}).Error; err != nil {
			return err
		}

		// 写入库存调整记录。
		reason := "错架药品移位归位"
		adj := &inventoryModel.InventoryAdjustment{
			TraceCode:      traceCode,
			DrugID:         trace.DrugID,
			AdjustType:     inventoryModel.AdjustTypeRelocate,
			BeforeStatus:   &fromStatus,
			AfterStatus:    &toStatus,
			FromLocationID: fromLocationID,
			ToLocationID:   &toLocationID,
			Reason:         reason,
			OperatorID:     operatorID,
		}
		if err := tx.Create(adj).Error; err != nil {
			return err
		}

		// 写入追溯日志。
		remark := "错架移位归位"
		log := &inventoryModel.DrugTraceLog{
			TraceCode:    traceCode,
			ActionType:   "RELOCATION",
			FromStatus:   &fromStatus,
			ToStatus:     &toStatus,
			ToLocationID: &toLocationID,
			OperatorID:   operatorID,
			Remark:       &remark,
		}
		return tx.Create(log).Error
	})
}

func (s *inventoryTaskService) EnrichTaskCounts(ctx context.Context, tasks []taskModel.InventoryTask) {
	if len(tasks) == 0 {
		return
	}
	ids := make([]int64, 0, len(tasks))
	for _, t := range tasks {
		ids = append(ids, t.ID)
	}

	type countRow struct {
		TaskID int64
		Cnt    int64
	}
	var rows []countRow
	s.db.WithContext(ctx).
		Table("inventory_task_detail").
		Select("task_id, COUNT(*) AS cnt").
		Where("task_id IN ? AND deleted_at IS NULL", ids).
		Group("task_id").
		Scan(&rows)

	countMap := make(map[int64]int64, len(rows))
	for _, r := range rows {
		countMap[r.TaskID] = r.Cnt
	}
	for i := range tasks {
		tasks[i].ScannedCount = countMap[tasks[i].ID]
	}
}
