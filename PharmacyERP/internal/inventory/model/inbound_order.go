package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// InboundOrderStatusDraft 草稿状态。
	InboundOrderStatusDraft = "DRAFT"
	// InboundOrderStatusPendingConfirm 待确认状态。
	InboundOrderStatusPendingConfirm = "PENDING_CONFIRM"
	// InboundOrderStatusCompleted 已完成状态。
	InboundOrderStatusCompleted = "COMPLETED"
	// 入库单状态：已取消
	InboundOrderStatusCancelled = "CANCELLED"
)

// InboundOrder 表示采购入库单主表。
// 对应数据库表：inbound_order。
type InboundOrder struct {
	core.BaseModel
	OrderNo    string               `gorm:"column:order_no;type:varchar(50);not null;uniqueIndex" json:"order_no"`
	Supplier   *string              `gorm:"column:supplier;type:varchar(100)" json:"supplier,omitempty"`
	InvoiceNo  *string              `gorm:"column:invoice_no;type:varchar(100)" json:"invoice_no,omitempty"`
	OperatorID int64                `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	Status     string               `gorm:"column:status;type:varchar(20);not null" json:"status"`
	Remark     *string              `gorm:"column:remark;type:text" json:"remark,omitempty"`
	Details    []InboundOrderDetail `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"details,omitempty"`
}

func (InboundOrder) TableName() string {
	return "inbound_order"
}

// InboundOrderDetail 表示入库单明细行。
// 对应数据库表：inbound_order_detail。
type InboundOrderDetail struct {
	core.BaseModel
	OrderID      int64     `gorm:"column:order_id;type:bigint;not null;index" json:"order_id"`
	DrugID       int64     `gorm:"column:drug_id;type:bigint;not null;index" json:"drug_id"`
	BatchNumber  string    `gorm:"column:batch_number;type:varchar(50);not null" json:"batch_number"`
	ExpireDate   time.Time `gorm:"column:expire_date;type:date;not null" json:"expire_date"`
	PlannedQty   int32     `gorm:"column:planned_qty;type:int;not null" json:"planned_qty"`
	ConfirmedQty int32     `gorm:"column:confirmed_qty;type:int;not null;default:0" json:"confirmed_qty"`
	UnitPrice    *float64  `gorm:"column:unit_price;type:decimal(10,2)" json:"unit_price,omitempty"`
}

func (InboundOrderDetail) TableName() string {
	return "inbound_order_detail"
}
