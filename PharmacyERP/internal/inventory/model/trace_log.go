package model

import "github.com/YingmoY/PharmacyERP/internal/pkg/core"

// DrugTraceLog maps to public.drug_trace_log.
type DrugTraceLog struct {
	core.BaseModel
	TraceCode      string  `gorm:"column:trace_code;type:varchar(100);not null;index:idx_trace_log_trace_code" json:"trace_code"`
	ActionType     string  `gorm:"column:action_type;type:varchar(50);not null" json:"action_type"`
	FromStatus     *string `gorm:"column:from_status;type:varchar(20)" json:"from_status,omitempty"`
	ToStatus       *string `gorm:"column:to_status;type:varchar(20)" json:"to_status,omitempty"`
	OperatorID     int64   `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	RelatedNo      *string `gorm:"column:related_no;type:varchar(100)" json:"related_no,omitempty"`
	Remark         *string `gorm:"column:remark;type:text" json:"remark,omitempty"`
	DrugID         *int64  `gorm:"column:drug_id;type:bigint" json:"drug_id,omitempty"`
	OrderID        *int64  `gorm:"column:order_id;type:bigint" json:"order_id,omitempty"`
	OrderItemID    *int64  `gorm:"column:order_item_id;type:bigint" json:"order_item_id,omitempty"`
	RequestID      *string `gorm:"column:request_id;type:varchar(100)" json:"request_id,omitempty"`
	FromLocationID *int64  `gorm:"column:from_location_id;type:bigint" json:"from_location_id,omitempty"`
	ToLocationID   *int64  `gorm:"column:to_location_id;type:bigint" json:"to_location_id,omitempty"`

	// Virtual display fields — populated via JOIN queries, not stored in DB.
	DrugName         string `gorm:"-" json:"drug_name,omitempty"`
	OperatorName     string `gorm:"-" json:"operator_name,omitempty"`
	FromLocationCode string `gorm:"-" json:"from_location_code,omitempty"`
	ToLocationCode   string `gorm:"-" json:"to_location_code,omitempty"`
}

func (DrugTraceLog) TableName() string {
	return "drug_trace_log"
}
