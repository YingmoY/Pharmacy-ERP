// Package model 定义通知模块的数据模型。
package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// Notification 映射 public.notification 表，存储用户通知消息。
type Notification struct {
	core.BaseModel
	// UserID 接收通知的用户 ID（关联 sys_user）。
	UserID int64 `gorm:"column:user_id;type:bigint;not null;index:idx_notification_user_read,priority:1" json:"user_id"`
	// Title 通知标题。
	Title string `gorm:"column:title;type:varchar(100);not null" json:"title"`
	// Content 通知内容。
	Content string `gorm:"column:content;type:text;not null" json:"content"`
	// NotificationType 通知类型，如 ALERT/SYSTEM/ORDER 等。
	NotificationType string `gorm:"column:notification_type;type:varchar(50)" json:"notification_type"`
	// BusinessType 关联业务类型。
	BusinessType string `gorm:"column:business_type;type:varchar(50);index:idx_notification_business,priority:1" json:"business_type"`
	// BusinessID 关联业务对象 ID。
	BusinessID string `gorm:"column:business_id;type:varchar(100);index:idx_notification_business,priority:2" json:"business_id"`
	// ReadAt 已读时间，nil 表示未读。
	ReadAt *time.Time `gorm:"column:read_at;type:timestamptz;index:idx_notification_user_read,priority:2" json:"read_at,omitempty"`
}

// TableName 指定数据库表名。
func (Notification) TableName() string {
	return "notification"
}

// IsRead 判断通知是否已读。
func (n *Notification) IsRead() bool {
	return n.ReadAt != nil
}
