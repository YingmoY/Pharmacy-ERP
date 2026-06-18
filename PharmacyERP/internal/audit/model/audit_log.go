// Package model 定义审计模块的数据模型（登录日志、操作日志、数据变更日志、安全事件）。
package model

import (
	"time"

	"gorm.io/datatypes"
)

// LoginLog 映射 public.login_log 表，记录用户登录日志。
type LoginLog struct {
	// ID 主键。
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// UserID 登录用户 ID（可为空，登录失败时可能无用户 ID）。
	UserID *int64 `gorm:"column:user_id;type:bigint" json:"user_id,omitempty"`
	// Username 登录用户名。
	Username string `gorm:"column:username;type:varchar(50)" json:"username"`
	// Success 是否登录成功。
	Success bool `gorm:"column:success;not null" json:"success"`
	// IP 登录来源 IP 地址。
	IP string `gorm:"column:ip;type:varchar(64)" json:"ip"`
	// UserAgent 客户端 User-Agent。
	UserAgent string `gorm:"column:user_agent;type:text" json:"user_agent"`
	// Message 附加消息（如失败原因）。
	Message string `gorm:"column:message;type:varchar(255)" json:"message"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

// TableName 指定数据库表名。
func (LoginLog) TableName() string {
	return "login_log"
}

// OperationLog 映射 public.operation_log 表，记录系统操作日志。
type OperationLog struct {
	// ID 主键。
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// Module 操作模块（对应 LogEvent.BusinessType）。
	Module string `gorm:"column:module;type:varchar(50)" json:"module"`
	// BusinessType 业务类型。
	BusinessType string `gorm:"column:business_type;type:varchar(50);not null" json:"business_type"`
	// BusinessID 业务对象 ID。
	BusinessID string `gorm:"column:business_id;type:varchar(100);not null" json:"business_id"`
	// Action 操作动作（如 CREATE/UPDATE/DELETE）。
	Action string `gorm:"column:action;type:varchar(50);not null" json:"action"`
	// OperatorID 操作人 ID（关联 sys_user）。
	OperatorID int64 `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	// ResourceType 资源类型。
	ResourceType string `gorm:"column:resource_type;type:varchar(50)" json:"resource_type"`
	// ResourceID 资源 ID。
	ResourceID string `gorm:"column:resource_id;type:varchar(100)" json:"resource_id"`
	// Detail 操作详情（JSONB）。
	Detail datatypes.JSON `gorm:"column:detail;type:jsonb" json:"detail"`
	// BeforeData 操作前数据（JSONB）。
	BeforeData datatypes.JSON `gorm:"column:before_data;type:jsonb" json:"before_data"`
	// AfterData 操作后数据（JSONB）。
	AfterData datatypes.JSON `gorm:"column:after_data;type:jsonb" json:"after_data"`
	// IP 操作来源 IP。
	IP string `gorm:"column:ip;type:varchar(64)" json:"ip"`
	// UserAgent 客户端 User-Agent。
	UserAgent string `gorm:"column:user_agent;type:text" json:"user_agent"`
	// RequestID 请求追踪 ID。
	RequestID string `gorm:"column:request_id;type:varchar(100)" json:"request_id"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	// UpdatedAt 更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// OperatorName 操作人名称，由仓储层批量填充，不存储于数据库。
	OperatorName string `gorm:"-" json:"operator_name,omitempty"`
}

// TableName 指定数据库表名。
func (OperationLog) TableName() string {
	return "operation_log"
}

// DataChangeLog 映射 public.data_change_log 表，记录数据变更日志。
type DataChangeLog struct {
	// ID 主键。
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// TargetTable 变更表名（GORM 列名 table_name）。
	TargetTable string `gorm:"column:table_name;type:varchar(100);not null" json:"table_name"`
	// RecordID 变更记录 ID。
	RecordID string `gorm:"column:record_id;type:varchar(100);not null" json:"record_id"`
	// ChangeType 变更类型：CREATE/UPDATE/DELETE。
	ChangeType string `gorm:"column:change_type;type:varchar(20);not null" json:"change_type"`
	// OperatorID 操作人 ID（关联 sys_user）。
	OperatorID *int64 `gorm:"column:operator_id;type:bigint" json:"operator_id,omitempty"`
	// OperatorName 操作人名称（冗余存储）。
	OperatorName string `gorm:"column:operator_name;type:varchar(100)" json:"operator_name"`
	// BeforeData 变更前数据（JSONB）。
	BeforeData datatypes.JSON `gorm:"column:before_data;type:jsonb" json:"before_data"`
	// AfterData 变更后数据（JSONB）。
	AfterData datatypes.JSON `gorm:"column:after_data;type:jsonb" json:"after_data"`
	// ChangedFields 变更字段列表（JSONB）。
	ChangedFields datatypes.JSON `gorm:"column:changed_fields;type:jsonb" json:"changed_fields"`
	// RequestID 请求追踪 ID。
	RequestID string `gorm:"column:request_id;type:varchar(100)" json:"request_id"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

// TableName 指定数据库表名。
func (DataChangeLog) TableName() string {
	return "data_change_log"
}

// SecurityEvent 映射 public.security_event 表，记录安全事件。
type SecurityEvent struct {
	// ID 主键。
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// EventType 安全事件类型。
	EventType string `gorm:"column:event_type;type:varchar(50);not null" json:"event_type"`
	// Severity 严重等级：HIGH/MEDIUM/LOW。
	Severity string `gorm:"column:severity;type:varchar(20);not null" json:"severity"`
	// UserID 关联用户 ID（可为空）。
	UserID *int64 `gorm:"column:user_id;type:bigint" json:"user_id,omitempty"`
	// Username 相关用户名。
	Username string `gorm:"column:username;type:varchar(50)" json:"username"`
	// IP 事件来源 IP。
	IP string `gorm:"column:ip;type:varchar(64)" json:"ip"`
	// Description 事件描述。
	Description string `gorm:"column:description;type:text;not null" json:"description"`
	// Detail 事件详情（JSONB）。
	Detail datatypes.JSON `gorm:"column:detail;type:jsonb" json:"detail"`
	// Handled 是否已处理。
	Handled bool `gorm:"column:handled;not null;default:false" json:"handled"`
	// HandledBy 处理人 ID（关联 sys_user）。
	HandledBy *int64 `gorm:"column:handled_by;type:bigint" json:"handled_by,omitempty"`
	// HandledAt 处理时间。
	HandledAt *time.Time `gorm:"column:handled_at;type:timestamptz" json:"handled_at,omitempty"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

// TableName 指定数据库表名。
func (SecurityEvent) TableName() string {
	return "security_event"
}
