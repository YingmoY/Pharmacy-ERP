// Package service 实现通知模块的业务逻辑层。
package service

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/notification/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
)

// NotificationDTO 通知数据传输对象。
type NotificationDTO struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	NotificationType string     `json:"notification_type"`
	BusinessType     string     `json:"business_type"`
	BusinessID       string     `json:"business_id"`
	ReadAt           *time.Time `json:"read_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// NotificationFilter 通知列表查询过滤条件。
type NotificationFilter struct {
	// Read 已读状态筛选：nil=全部，true=已读，false=未读。
	Read     *bool
	// Page 页码（从 1 开始）。
	Page     int
	// PageSize 每页条数。
	PageSize int
}

// NotificationService 通知服务接口。
type NotificationService interface {
	// ListNotifications 分页查询当前用户通知列表。
	ListNotifications(ctx context.Context, userID int64, filter NotificationFilter) ([]*NotificationDTO, int64, error)
	// GetUnreadCount 获取当前用户未读通知数量。
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	// MarkRead 将指定通知标记为已读。
	MarkRead(ctx context.Context, id int64, userID int64) error
	// MarkAllRead 将当前用户所有通知标记为已读。
	MarkAllRead(ctx context.Context, userID int64) error
}

type notificationService struct {
	repo repository.NotificationRepo
	log  *zap.Logger
}

// NewNotificationService 创建通知服务实例。
func NewNotificationService(repo repository.NotificationRepo, log *zap.Logger) NotificationService {
	return &notificationService{repo: repo, log: log}
}

// ListNotifications 分页查询当前用户通知列表。
func (s *notificationService) ListNotifications(ctx context.Context, userID int64, filter NotificationFilter) ([]*NotificationDTO, int64, error) {
	repoFilter := repository.NotificationFilter{
		Read:     filter.Read,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	notifications, total, err := s.repo.List(ctx, userID, repoFilter)
	if err != nil {
		s.log.Error("查询通知列表失败", zap.Int64("userID", userID), zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	result := make([]*NotificationDTO, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, &NotificationDTO{
			ID:               n.ID,
			Title:            n.Title,
			Content:          n.Content,
			NotificationType: n.NotificationType,
			BusinessType:     n.BusinessType,
			BusinessID:       n.BusinessID,
			ReadAt:           n.ReadAt,
			CreatedAt:        n.CreatedAt,
		})
	}
	return result, total, nil
}

// GetUnreadCount 获取当前用户未读通知数量。
func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	count, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		s.log.Error("查询未读通知数量失败", zap.Int64("userID", userID), zap.Error(err))
		return 0, ecode.ErrSystem
	}
	return count, nil
}

// MarkRead 将指定通知标记为已读。
func (s *notificationService) MarkRead(ctx context.Context, id int64, userID int64) error {
	if err := s.repo.MarkRead(ctx, id, userID); err != nil {
		s.log.Error("标记通知已读失败", zap.Int64("id", id), zap.Int64("userID", userID), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// MarkAllRead 将当前用户所有未读通知标记为已读。
func (s *notificationService) MarkAllRead(ctx context.Context, userID int64) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		s.log.Error("批量标记通知已读失败", zap.Int64("userID", userID), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}
