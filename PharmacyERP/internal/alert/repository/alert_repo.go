// Package repository 实现告警模块的数据访问层。
package repository

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/alert/model"
	"gorm.io/gorm"
)

// AlertFilter 告警列表查询过滤条件。
type AlertFilter struct {
	// Status 按状态筛选（nil 表示不限）。
	Status *int16
	// EventType 按告警类型筛选。
	EventType string
	// Severity 按严重等级筛选。
	Severity string
	// Page 页码（从 1 开始）。
	Page int
	// PageSize 每页条数。
	PageSize int
}

// AlertRepo 告警仓储接口。
type AlertRepo interface {
	// List 分页查询告警列表。
	List(ctx context.Context, filter AlertFilter) ([]*model.Alert, int64, error)
	// GetByID 按 ID 查询告警。
	GetByID(ctx context.Context, id int64) (*model.Alert, error)
	// UpdateResolve 将告警标记为已解决。
	UpdateResolve(ctx context.Context, id int64, resolvedBy int64, resolution string, closedAt time.Time) error
	// UpdateIgnore 将告警标记为已忽略。
	UpdateIgnore(ctx context.Context, id int64, ignoredBy int64, ignoredAt time.Time, reason string) error
	// CreateAlert 创建新告警记录。
	CreateAlert(ctx context.Context, alert *model.Alert) error
}

type alertRepo struct {
	db *gorm.DB
}

// NewAlertRepo 创建告警仓储实例。
func NewAlertRepo(db *gorm.DB) AlertRepo {
	return &alertRepo{db: db}
}

// List 分页查询告警列表。
func (r *alertRepo) List(ctx context.Context, filter AlertFilter) ([]*model.Alert, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Alert{})

	// 按状态筛选
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	// 按告警类型筛选
	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	// 按严重等级筛选
	if filter.Severity != "" {
		query = query.Where("severity = ?", filter.Severity)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Alert{}, 0, nil
	}

	// 分页，按创建时间降序
	offset := (filter.Page - 1) * filter.PageSize
	var alerts []*model.Alert
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// GetByID 按 ID 查询告警。
func (r *alertRepo) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	var alert model.Alert
	if err := r.db.WithContext(ctx).First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateResolve 将告警标记为已解决状态。
func (r *alertRepo) UpdateResolve(ctx context.Context, id int64, resolvedBy int64, resolution string, closedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Alert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      model.AlertStatusResolved,
			"resolved_by": resolvedBy,
			"resolution":  resolution,
			"closed_at":   closedAt,
		}).Error
}

// UpdateIgnore 将告警标记为已忽略状态。
func (r *alertRepo) UpdateIgnore(ctx context.Context, id int64, ignoredBy int64, ignoredAt time.Time, reason string) error {
	return r.db.WithContext(ctx).Model(&model.Alert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.AlertStatusIgnored,
			"ignored_by": ignoredBy,
			"ignored_at": ignoredAt,
			"resolution": reason,
			"closed_at":  ignoredAt,
		}).Error
}

// CreateAlert 创建新告警记录。
func (r *alertRepo) CreateAlert(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}
