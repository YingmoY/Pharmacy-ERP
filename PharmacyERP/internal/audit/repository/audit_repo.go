// Package repository 实现审计模块的数据访问层。
package repository

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/audit/model"
	"gorm.io/gorm"
)

// LoginLogFilter 登录日志查询过滤条件。
type LoginLogFilter struct {
	UserID    *int64
	Username  string
	Success   *bool
	StartTime *time.Time
	EndTime   *time.Time
	Page      int
	PageSize  int
}

// OperationLogFilter 操作日志查询过滤条件。
type OperationLogFilter struct {
	OperatorID   *int64
	OperatorName string
	Module       string
	Action       string
	StartTime    *time.Time
	EndTime      *time.Time
	Page         int
	PageSize     int
}

// DataChangeLogFilter 数据变更日志查询过滤条件。
type DataChangeLogFilter struct {
	TableName  string
	RecordID   string
	ChangeType string
	StartTime  *time.Time
	EndTime    *time.Time
	Page       int
	PageSize   int
}

// SecurityEventFilter 安全事件查询过滤条件。
type SecurityEventFilter struct {
	UserID    *int64
	EventType string
	Handled   *bool
	StartTime *time.Time
	EndTime   *time.Time
	Page      int
	PageSize  int
}

// AuditRepo 审计仓储接口。
type AuditRepo interface {
	// ListLoginLogs 分页查询登录日志。
	ListLoginLogs(ctx context.Context, filter LoginLogFilter) ([]*model.LoginLog, int64, error)
	// ListOperationLogs 分页查询操作日志。
	ListOperationLogs(ctx context.Context, filter OperationLogFilter) ([]*model.OperationLog, int64, error)
	// GetOperationLog 按 ID 查询操作日志详情。
	GetOperationLog(ctx context.Context, id int64) (*model.OperationLog, error)
	// ListDataChangeLogs 分页查询数据变更日志。
	ListDataChangeLogs(ctx context.Context, filter DataChangeLogFilter) ([]*model.DataChangeLog, int64, error)
	// ListSecurityEvents 分页查询安全事件。
	ListSecurityEvents(ctx context.Context, filter SecurityEventFilter) ([]*model.SecurityEvent, int64, error)
	// SaveOperationLog 保存操作日志（由 MQ 消费者调用）。
	SaveOperationLog(ctx context.Context, log *model.OperationLog) error
}

type auditRepo struct {
	db *gorm.DB
}

// NewAuditRepo 创建审计仓储实例。
func NewAuditRepo(db *gorm.DB) AuditRepo {
	return &auditRepo{db: db}
}

// ListLoginLogs 分页查询登录日志，按创建时间降序。
func (r *auditRepo) ListLoginLogs(ctx context.Context, filter LoginLogFilter) ([]*model.LoginLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.LoginLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Username != "" {
		query = query.Where("username ILIKE ?", "%"+filter.Username+"%")
	}
	if filter.Success != nil {
		query = query.Where("success = ?", *filter.Success)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at < ?", *filter.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*model.LoginLog{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var logs []*model.LoginLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ListOperationLogs 分页查询操作日志，按创建时间降序，结果中包含 operator_name。
func (r *auditRepo) ListOperationLogs(ctx context.Context, filter OperationLogFilter) ([]*model.OperationLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.OperationLog{})

	if filter.OperatorID != nil {
		query = query.Where("operator_id = ?", *filter.OperatorID)
	} else if filter.OperatorName != "" {
		// 按操作人姓名模糊查找用户 ID 列表，再过滤日志。
		var userIDs []int64
		r.db.WithContext(ctx).Table("sys_user").Select("id").
			Where("real_name ILIKE ? AND deleted_at IS NULL", "%"+filter.OperatorName+"%").
			Scan(&userIDs)
		if len(userIDs) == 0 {
			return []*model.OperationLog{}, 0, nil
		}
		query = query.Where("operator_id IN ?", userIDs)
	}
	if filter.Module != "" {
		query = query.Where("module = ? OR business_type = ?", filter.Module, filter.Module)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at < ?", *filter.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*model.OperationLog{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var logs []*model.OperationLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	// 批量填充 operator_name。
	operatorIDs := make([]int64, 0, len(logs))
	for _, l := range logs {
		if l.OperatorID > 0 {
			operatorIDs = append(operatorIDs, l.OperatorID)
		}
	}
	if len(operatorIDs) > 0 {
		type userRow struct {
			ID       int64
			RealName string
		}
		var users []userRow
		r.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", operatorIDs).Scan(&users)
		userMap := make(map[int64]string, len(users))
		for _, u := range users {
			userMap[u.ID] = u.RealName
		}
		for _, l := range logs {
			l.OperatorName = userMap[l.OperatorID]
		}
	}

	return logs, total, nil
}

// GetOperationLog 按 ID 查询操作日志详情（含 operator_name）。
func (r *auditRepo) GetOperationLog(ctx context.Context, id int64) (*model.OperationLog, error) {
	var log model.OperationLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		return nil, err
	}
	if log.OperatorID > 0 {
		var name string
		r.db.WithContext(ctx).Table("sys_user").Select("real_name").Where("id = ?", log.OperatorID).Scan(&name)
		log.OperatorName = name
	}
	return &log, nil
}

// ListDataChangeLogs 分页查询数据变更日志，按创建时间降序。
func (r *auditRepo) ListDataChangeLogs(ctx context.Context, filter DataChangeLogFilter) ([]*model.DataChangeLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.DataChangeLog{})

	if filter.TableName != "" {
		query = query.Where("table_name = ?", filter.TableName)
	}
	if filter.RecordID != "" {
		query = query.Where("record_id = ?", filter.RecordID)
	}
	if filter.ChangeType != "" {
		query = query.Where("change_type = ?", filter.ChangeType)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at < ?", *filter.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*model.DataChangeLog{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var logs []*model.DataChangeLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ListSecurityEvents 分页查询安全事件，按创建时间降序。
func (r *auditRepo) ListSecurityEvents(ctx context.Context, filter SecurityEventFilter) ([]*model.SecurityEvent, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.SecurityEvent{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.Handled != nil {
		query = query.Where("handled = ?", *filter.Handled)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at < ?", *filter.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*model.SecurityEvent{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var events []*model.SecurityEvent
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

// SaveOperationLog 保存操作日志（供 MQ 消费者调用）。
func (r *auditRepo) SaveOperationLog(ctx context.Context, log *model.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
