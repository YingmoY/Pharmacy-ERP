package model

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// AdjustTypeRelocate 货位调整（移位）。
	AdjustTypeRelocate = "RELOCATE"
	// AdjustTypeLoss 确认盘亏（LOSS_CANDIDATE -> LOST）。
	AdjustTypeLoss = "LOSS"
	// AdjustTypeStatusChange 手动状态变更。
	AdjustTypeStatusChange = "STATUS_CHANGE"
	// AdjustTypeOther 其他调整原因。
	AdjustTypeOther = "OTHER"
)

// InventoryAdjustment 记录每一次人工或系统触发的库存调整操作。
// 对应数据库表：inventory_adjustment。
type InventoryAdjustment struct {
	core.BaseModel
	AdjustNo       string  `gorm:"column:adjust_no;type:varchar(50);not null;uniqueIndex" json:"adjust_no"`
	TraceCode      string  `gorm:"column:trace_code;type:varchar(100);not null;index"     json:"trace_code"`
	DrugID         int64   `gorm:"column:drug_id;type:bigint"                            json:"drug_id"`
	AdjustType     string  `gorm:"column:adjust_type;type:varchar(30);not null"          json:"adjust_type"`
	BeforeStatus   *string `gorm:"column:before_status;type:varchar(20)"                 json:"before_status,omitempty"`
	AfterStatus    *string `gorm:"column:after_status;type:varchar(20)"                  json:"after_status,omitempty"`
	FromLocationID *int64  `gorm:"column:from_location_id;type:bigint"                   json:"from_location_id,omitempty"`
	ToLocationID   *int64  `gorm:"column:to_location_id;type:bigint"                     json:"to_location_id,omitempty"`
	Reason         string  `gorm:"column:reason;type:text;not null"                      json:"reason"`
	OperatorID     int64   `gorm:"column:operator_id;type:bigint;not null"               json:"operator_id"`
	Status         string  `gorm:"column:status;type:varchar(20);not null;default:COMPLETED" json:"status"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	DrugName         string `gorm:"-" json:"drug_name,omitempty"`
	FromLocationCode string `gorm:"-" json:"from_location_code,omitempty"`
	ToLocationCode   string `gorm:"-" json:"to_location_code,omitempty"`
	OperatorName     string `gorm:"-" json:"operator_name,omitempty"`
}

func (InventoryAdjustment) TableName() string {
	return "inventory_adjustment"
}

// BeforeCreate 在插入前自动生成 adjust_no 和设置默认 status。
func (a *InventoryAdjustment) BeforeCreate(_ *gorm.DB) error {
	if a.AdjustNo == "" {
		a.AdjustNo = "ADJ-" + uuid.New().String()
	}
	if a.Status == "" {
		a.Status = "COMPLETED"
	}
	return nil
}
