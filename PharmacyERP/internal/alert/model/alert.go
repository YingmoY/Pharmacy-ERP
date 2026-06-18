// Package model 定义告警模块的数据模型。
package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// AlertStatusActive 告警激活中（未处理）。
	AlertStatusActive = int16(0)
	// AlertStatusResolved 告警已解决。
	AlertStatusResolved = int16(1)
	// AlertStatusIgnored 告警已忽略。
	AlertStatusIgnored = int16(2)
)

const (
	// AlertEventTypeNearExpire 近效期告警。
	AlertEventTypeNearExpire = "NEAR_EXPIRE"
	// AlertEventTypeLowStock 低库存告警。
	AlertEventTypeLowStock = "LOW_STOCK"
	// AlertEventTypeLossCandidate 盘亏候选告警。
	AlertEventTypeLossCandidate = "LOSS_CANDIDATE"
	// AlertEventTypeMisplaced 错架告警。
	AlertEventTypeMisplaced = "MISPLACED"
)

const (
	// AlertSeverityHigh 高严重等级。
	AlertSeverityHigh = "HIGH"
	// AlertSeverityMedium 中严重等级。
	AlertSeverityMedium = "MEDIUM"
	// AlertSeverityLow 低严重等级。
	AlertSeverityLow = "LOW"
)

// Alert 映射 public.audit_event 表，存储系统告警事件。
type Alert struct {
	core.BaseModel
	// EventType 告警类型：NEAR_EXPIRE/LOW_STOCK/LOSS_CANDIDATE/MISPLACED 等。
	EventType string `gorm:"column:event_type;type:varchar(50);not null" json:"event_type"`
	// RelatedType 关联业务对象类型，如 drug_trace_inventory。
	RelatedType string `gorm:"column:related_type;type:varchar(50);not null" json:"related_type"`
	// RelatedID 关联业务对象 ID。
	RelatedID string `gorm:"column:related_id;type:varchar(100);not null" json:"related_id"`
	// Description 告警描述。
	Description *string `gorm:"column:description;type:text" json:"description,omitempty"`
	// AssignedTo 指派处理人 ID（关联 sys_user）。
	AssignedTo *int64 `gorm:"column:assigned_to;type:bigint" json:"assigned_to,omitempty"`
	// Status 状态：0=ACTIVE, 1=RESOLVED, 2=IGNORED。
	Status int16 `gorm:"column:status;type:smallint;default:0" json:"status"`
	// Resolution 处理说明。
	Resolution *string `gorm:"column:resolution;type:text" json:"resolution,omitempty"`
	// ClosedAt 关闭时间（解决或忽略时记录）。
	ClosedAt *time.Time `gorm:"column:closed_at;type:timestamptz" json:"closed_at,omitempty"`
	// Severity 严重等级：HIGH/MEDIUM/LOW。
	Severity string `gorm:"column:severity;type:varchar(20)" json:"severity"`
	// IgnoredAt 忽略时间。
	IgnoredAt *time.Time `gorm:"column:ignored_at;type:timestamptz" json:"ignored_at,omitempty"`
	// IgnoredBy 忽略操作人 ID。
	IgnoredBy *int64 `gorm:"column:ignored_by;type:bigint" json:"ignored_by,omitempty"`
	// ResolvedBy 解决操作人 ID。
	ResolvedBy *int64 `gorm:"column:resolved_by;type:bigint" json:"resolved_by,omitempty"`
}

// TableName 指定数据库表名。
func (Alert) TableName() string {
	return "audit_event"
}
