package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// ScanTaskStatusPending 扫码任务待开始。
	ScanTaskStatusPending = "PENDING"
	// ScanTaskStatusInProgress 扫码任务进行中。
	ScanTaskStatusInProgress = "IN_PROGRESS"
	// ScanTaskStatusCompleted 扫码任务已完成。
	ScanTaskStatusCompleted = "COMPLETED"
	// ScanTaskStatusCancelled 扫码任务已取消。
	ScanTaskStatusCancelled = "CANCELLED"
)

const (
	// ScanTaskTypeInbound 入库类型扫码任务。
	ScanTaskTypeInbound = "INBOUND"
	// ScanTaskTypeShelving 上架类型扫码任务。
	ScanTaskTypeShelving = "SHELVING"
	// ScanTaskTypeInventory 盘点类型扫码任务。
	ScanTaskTypeInventory = "INVENTORY"
)

const (
	// ScanResultSuccess 扫描成功。
	ScanResultSuccess = "SUCCESS"
	// ScanResultDuplicate 重复扫描。
	ScanResultDuplicate = "DUPLICATE"
	// ScanResultInvalid 无效追溯码。
	ScanResultInvalid = "INVALID"
	// ScanResultStatusError 状态错误（业务规则不满足）。
	ScanResultStatusError = "STATUS_ERROR"
)

// ScanTask 映射 public.scan_task。
// 记录一次扫码操作会话（可包含多次单码扫描）。
type ScanTask struct {
	core.BaseModel
	TaskNo     string     `gorm:"column:task_no;type:varchar(50);not null;uniqueIndex" json:"task_no"`
	TaskType   string     `gorm:"column:task_type;type:varchar(20);not null" json:"task_type"`
	RelatedID  int64      `gorm:"column:related_id;type:bigint;not null" json:"related_id"`
	OperatorID int64      `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	Status     string     `gorm:"column:status;type:varchar(20);not null" json:"status"`
	StartTime  *time.Time `gorm:"column:start_time;type:timestamptz" json:"start_time,omitempty"`
	EndTime    *time.Time `gorm:"column:end_time;type:timestamptz" json:"end_time,omitempty"`
	Remark     *string    `gorm:"column:remark;type:text" json:"remark,omitempty"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	AssignedToName string `gorm:"-" json:"assigned_to_name,omitempty"`
	RelatedOrderNo string `gorm:"-" json:"related_order_no,omitempty"`
}

func (ScanTask) TableName() string {
	return "scan_task"
}

// ScanTaskDetail 映射 public.scan_task_detail。
// 记录单次扫码的结果详情。
type ScanTaskDetail struct {
	core.BaseModel
	TaskID       int64     `gorm:"column:task_id;type:bigint;not null;index" json:"task_id"`
	TraceCode    string    `gorm:"column:trace_code;type:varchar(100);not null" json:"trace_code"`
	LocationCode *string   `gorm:"column:location_code;type:varchar(50)" json:"location_code,omitempty"`
	ScanResult   string    `gorm:"column:scan_result;type:varchar(20);not null" json:"scan_result"`
	ErrorMsg     *string   `gorm:"column:error_msg;type:varchar(255)" json:"error_msg,omitempty"`
	ScanTime     time.Time `gorm:"column:scan_time;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"scan_time"`
}

func (ScanTaskDetail) TableName() string {
	return "scan_task_detail"
}
