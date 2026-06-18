// Package repository 实现通知模块的数据访问层。
package repository

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/notification/model"
	"gorm.io/gorm"
)

// NotificationFilter 通知列表查询过滤条件。
type NotificationFilter struct {
	// Read 已读状态筛选：nil=全部，true=已读，false=未读。
	Read     *bool
	// Page 页码（从 1 开始）。
	Page     int
	// PageSize 每页条数。
	PageSize int
}

// NotificationRepo 通知仓储接口。
type NotificationRepo interface {
	// List 分页查询用户通知列表（未读优先）。
	List(ctx context.Context, userID int64, filter NotificationFilter) ([]*model.Notification, int64, error)
	// GetUnreadCount 查询用户未读通知数量。
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	// MarkRead 将指定通知标记为已读（仅限该用户自己的通知）。
	MarkRead(ctx context.Context, id int64, userID int64) error
	// MarkAllRead 将该用户所有未读通知标记为已读。
	MarkAllRead(ctx context.Context, userID int64) error
}

type notificationRepo struct {
	db *gorm.DB
}

// NewNotificationRepo 创建通知仓储实例。
func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

// List 分页查询用户通知列表，未读消息优先展示。
func (r *notificationRepo) List(ctx context.Context, userID int64, filter NotificationFilter) ([]*model.Notification, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)

	// 按已读状态筛选
	if filter.Read != nil {
		if *filter.Read {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Notification{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var notifications []*model.Notification
	// 未读优先（read_at IS NULL 排前），再按创建时间降序
	if err := query.
		Order("read_at IS NOT NULL ASC, created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetUnreadCount 查询用户未读通知数量。
func (r *notificationRepo) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// MarkRead 将指定 ID 的通知标记为已读。
func (r *notificationRepo) MarkRead(ctx context.Context, id int64, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now).Error
}

// MarkAllRead 将该用户所有未读通知标记为已读。
func (r *notificationRepo) MarkAllRead(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now).Error
}
