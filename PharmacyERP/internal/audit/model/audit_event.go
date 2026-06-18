package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// AuditEventStatusPending 待处理。
	AuditEventStatusPending int8 = 0
	// AuditEventStatusClosed 已关闭。
	AuditEventStatusClosed int8 = 1
)

// AuditEvent 映射 public.audit_event。
type AuditEvent struct {
	core.BaseModel
	EventType   string     `gorm:"column:event_type;type:varchar(50);not null" json:"event_type"`
	RelatedType string     `gorm:"column:related_type;type:varchar(50);not null" json:"related_type"`
	RelatedID   string     `gorm:"column:related_id;type:varchar(100);not null" json:"related_id"`
	Description *string    `gorm:"column:description;type:text" json:"description,omitempty"`
	AssignedTo  *int64     `gorm:"column:assigned_to;type:bigint" json:"assigned_to,omitempty"`
	Status      int8       `gorm:"column:status;type:smallint;default:0" json:"status"`
	Resolution  *string    `gorm:"column:resolution;type:text" json:"resolution,omitempty"`
	ClosedAt    *time.Time `gorm:"column:closed_at;type:timestamptz" json:"closed_at,omitempty"`
}

func (AuditEvent) TableName() string {
	return "audit_event"
}
