package repository

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/task/model"
	"gorm.io/gorm"
)

// InventoryTask 定义盘点任务仓储接口。
type InventoryTask interface {
	// CreateTask 创建盘点任务。
	CreateTask(ctx context.Context, tx *gorm.DB, task *model.InventoryTask) error
	// GetTaskByID 根据主键查询任务。
	GetTaskByID(ctx context.Context, tx *gorm.DB, id int64) (*model.InventoryTask, error)
	// GetTaskByNo 根据任务编号查询任务。
	GetTaskByNo(ctx context.Context, tx *gorm.DB, taskNo string) (*model.InventoryTask, error)
	// UpdateTask 更新任务。
	UpdateTask(ctx context.Context, tx *gorm.DB, task *model.InventoryTask) error
	// ListTasks 分页查询任务列表。
	ListTasks(ctx context.Context, tx *gorm.DB, page, pageSize int) ([]model.InventoryTask, int64, error)

	// Detail 相关
	// CreateDetail 写入扫描明细。
	CreateDetail(ctx context.Context, tx *gorm.DB, detail *model.InventoryTaskDetail) error
	// GetDetailByTraceCode 根据追溯码查询该任务内是否已有扫描记录（用于去重）。
	GetDetailByTraceCode(ctx context.Context, tx *gorm.DB, taskID int64, traceCode string) (*model.InventoryTaskDetail, error)
	// ListDetails 查询任务下所有扫描明细。
	ListDetails(ctx context.Context, tx *gorm.DB, taskID int64) ([]model.InventoryTaskDetail, error)
}

type inventoryTask struct{}

func NewInventoryTask() InventoryTask {
	return &inventoryTask{}
}

func (r *inventoryTask) CreateTask(ctx context.Context, tx *gorm.DB, task *model.InventoryTask) error {
	return tx.WithContext(ctx).Create(task).Error
}

func (r *inventoryTask) GetTaskByID(ctx context.Context, tx *gorm.DB, id int64) (*model.InventoryTask, error) {
	var task model.InventoryTask
	err := tx.WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *inventoryTask) GetTaskByNo(ctx context.Context, tx *gorm.DB, taskNo string) (*model.InventoryTask, error) {
	var task model.InventoryTask
	err := tx.WithContext(ctx).Where("task_no = ?", taskNo).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *inventoryTask) UpdateTask(ctx context.Context, tx *gorm.DB, task *model.InventoryTask) error {
	return tx.WithContext(ctx).Save(task).Error
}

func (r *inventoryTask) ListTasks(ctx context.Context, tx *gorm.DB, page, pageSize int) ([]model.InventoryTask, int64, error) {
	var tasks []model.InventoryTask
	var total int64

	query := tx.WithContext(ctx).Model(&model.InventoryTask{})
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
	return tasks, total, nil
}

func (r *inventoryTask) CreateDetail(ctx context.Context, tx *gorm.DB, detail *model.InventoryTaskDetail) error {
	return tx.WithContext(ctx).Create(detail).Error
}

func (r *inventoryTask) GetDetailByTraceCode(ctx context.Context, tx *gorm.DB, taskID int64, traceCode string) (*model.InventoryTaskDetail, error) {
	var detail model.InventoryTaskDetail
	err := tx.WithContext(ctx).
		Where("task_id = ? AND trace_code = ?", taskID, traceCode).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

func (r *inventoryTask) ListDetails(ctx context.Context, tx *gorm.DB, taskID int64) ([]model.InventoryTaskDetail, error) {
	var details []model.InventoryTaskDetail
	err := tx.WithContext(ctx).Where("task_id = ?", taskID).Find(&details).Error
	return details, err
}
