// Package model 定义报表模块的数据模型。
package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"gorm.io/datatypes"
)

const (
	// ExportTaskStatusPending 导出任务待处理。
	ExportTaskStatusPending = "PENDING"
	// ExportTaskStatusRunning 导出任务执行中。
	ExportTaskStatusRunning = "RUNNING"
	// ExportTaskStatusSuccess 导出任务已成功。
	ExportTaskStatusSuccess = "SUCCESS"
	// ExportTaskStatusFailed 导出任务失败。
	ExportTaskStatusFailed = "FAILED"
)

// ReportExportTask 映射 public.report_export_task 表，记录异步报表导出任务。
type ReportExportTask struct {
	core.BaseModel
	// TaskID 任务唯一标识（UUID 格式）。
	TaskID string `gorm:"column:task_id;type:varchar(100);not null;uniqueIndex" json:"task_id"`
	// ReportType 报表类型：SALES/INBOUND/INVENTORY/TRACE_LOG。
	ReportType string `gorm:"column:report_type;type:varchar(50);not null" json:"report_type"`
	// ExportFormat 导出格式，默认 xlsx。
	ExportFormat string `gorm:"column:export_format;type:varchar(20);not null;default:'xlsx'" json:"export_format"`
	// QueryParams 查询参数（JSON 格式）。
	QueryParams datatypes.JSON `gorm:"column:query_params;type:jsonb" json:"query_params"`
	// Status 任务状态：PENDING/RUNNING/SUCCESS/FAILED。
	Status string `gorm:"column:status;type:varchar(20);not null;default:'PENDING'" json:"status"`
	// FileID 生成的文件 ID（关联 file_info 表）。
	FileID *string `gorm:"column:file_id;type:varchar(100)" json:"file_id,omitempty"`
	// Message 任务消息（失败原因或成功提示）。
	Message *string `gorm:"column:message;type:text" json:"message,omitempty"`
	// RequestedBy 请求导出的用户 ID（关联 sys_user）。
	RequestedBy int64 `gorm:"column:requested_by;type:bigint;not null" json:"requested_by"`
	// StartedAt 任务开始时间。
	StartedAt *time.Time `gorm:"column:started_at;type:timestamptz" json:"started_at,omitempty"`
	// FinishedAt 任务完成时间。
	FinishedAt *time.Time `gorm:"column:finished_at;type:timestamptz" json:"finished_at,omitempty"`
}

// TableName 指定数据库表名。
func (ReportExportTask) TableName() string {
	return "report_export_task"
}
