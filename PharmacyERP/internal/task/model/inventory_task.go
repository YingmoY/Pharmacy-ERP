package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// InventoryTaskStatusPending 任务待开始。
	InventoryTaskStatusPending = "PENDING"
	// InventoryTaskStatusInProgress 任务盘点中。
	InventoryTaskStatusInProgress = "IN_PROGRESS"
	// InventoryTaskStatusCompleted 任务已完成。
	InventoryTaskStatusCompleted = "COMPLETED"
	// InventoryTaskStatusCancelled 任务已取消。
	InventoryTaskStatusCancelled = "CANCELLED"
)

const (
	// InventoryTaskScopeArea 按区域盘点。
	InventoryTaskScopeArea = "AREA"
	// InventoryTaskScopeShelf 按货架盘点。
	InventoryTaskScopeShelf = "SHELF"
	// InventoryTaskScopeLocation 按具体货位盘点。
	InventoryTaskScopeLocation = "LOCATION"
)

const (
	// InventoryDiscrepancyNormal 盘点结果正常（位置匹配）。
	InventoryDiscrepancyNormal = "NORMAL"
	// InventoryDiscrepancyMisplacedFound 扫描到错架药品（位置不匹配）。
	InventoryDiscrepancyMisplacedFound = "MISPLACED_FOUND"
	// InventoryDiscrepancyUnexpected 扫描到非在库状态的追溯码。
	InventoryDiscrepancyUnexpected = "UNEXPECTED"
)

// InventoryTask 映射 public.inventory_task。
// 用于记录一次盘点任务主信息（范围、创建人、执行人、状态与时间）。
type InventoryTask struct {
	core.BaseModel
	TaskNo     string                `gorm:"column:task_no;type:varchar(50);not null;uniqueIndex" json:"task_no"`
	ScopeType  string                `gorm:"column:scope_type;type:varchar(20);not null"          json:"scope_type"`
	ScopeValue string                `gorm:"column:scope_value;type:varchar(50);not null"         json:"scope_value"`
	CreatorID  int64                 `gorm:"column:creator_id;type:bigint;not null"               json:"creator_id"`
	AssigneeID *int64                `gorm:"column:assignee_id;type:bigint"                       json:"assignee_id,omitempty"`
	Status     string                `gorm:"column:status;type:varchar(20);not null;default:PENDING" json:"status"`
	StartTime  *time.Time            `gorm:"column:start_time;type:timestamptz"                   json:"start_time,omitempty"`
	EndTime    *time.Time            `gorm:"column:end_time;type:timestamptz"                     json:"end_time,omitempty"`
	Remark     *string               `gorm:"column:remark;type:text"                              json:"remark,omitempty"`
	Details    []InventoryTaskDetail `gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"details,omitempty"`

	// 虚拟字段，由服务层批量填充。
	ScannedCount       int64 `gorm:"-" json:"scanned_count"`
	NormalCount        int64 `gorm:"-" json:"normal_count"`
	MisplacedCount     int64 `gorm:"-" json:"misplaced_count"`
	UnexpectedCount    int64 `gorm:"-" json:"unexpected_count"`
	LossCandidateCount int64 `gorm:"-" json:"loss_candidate_count"`
}

func (InventoryTask) TableName() string {
	return "inventory_task"
}

// InventoryTaskDetail 映射 public.inventory_task_detail。
// 记录盘点过程中的实际扫描结果与差异类型。
type InventoryTaskDetail struct {
	core.BaseModel
	TaskID            int64      `gorm:"column:task_id;type:bigint;not null;index"           json:"task_id"`
	TraceCode         string     `gorm:"column:trace_code;type:varchar(100);not null;index"  json:"trace_code"`
	LocationID        int64      `gorm:"column:location_id;type:bigint;not null"             json:"location_id"`
	ScannedLocationID *int64     `gorm:"column:scanned_location_id;type:bigint"              json:"scanned_location_id,omitempty"`
	SystemLocationID  *int64     `gorm:"column:system_location_id;type:bigint"               json:"system_location_id,omitempty"`
	DiffType          string     `gorm:"column:discrepancy_type;type:varchar(20);not null"   json:"scan_result"`
	OperatorID        *int64     `gorm:"column:operator_id;type:bigint"                      json:"operator_id,omitempty"`
	ScanTime          *time.Time `gorm:"column:scanned_at;type:timestamptz"                  json:"scanned_at,omitempty"`
}

func (InventoryTaskDetail) TableName() string {
	return "inventory_task_detail"
}

// InventoryTaskSummary 盘点任务统计汇总。
type InventoryTaskSummary struct {
	TaskID            int64  `json:"task_id"`
	TaskNo            string `json:"task_no"`
	ScopeType         string `json:"scope_type"`
	ScopeValue        string `json:"scope_value"`
	Status            string `json:"status"`
	TotalScanned      int64  `json:"total_scanned"`
	NormalCount       int64  `json:"normal_count"`
	MisplacedCount    int64  `json:"misplaced_count"`
	UnexpectedCount   int64  `json:"unexpected_count"`
	LossCandidateCount int64 `json:"loss_candidate_count"`
}
