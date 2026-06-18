// Package model 定义入库单相关的数据模型。
package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// 入库单状态常量。
const (
	// InboundOrderStatusDraft 草稿状态，允许编辑明细。
	InboundOrderStatusDraft = "DRAFT"
	// InboundOrderStatusPendingConfirm 待确认状态，等待扫码录入追溯码。
	InboundOrderStatusPendingConfirm = "PENDING_CONFIRM"
	// InboundOrderStatusCompleted 已完成，所有明细均已全部确认到位。
	InboundOrderStatusCompleted = "COMPLETED"
	// InboundOrderStatusCancelled 已取消。
	InboundOrderStatusCancelled = "CANCELLED"
)

// 追溯库存状态常量。
const (
	// TraceInventoryStatusPending 待上架（入库确认后尚未上架）。
	TraceInventoryStatusPending = "PENDING"
	// TraceInventoryStatusInStock 在库（已上架）。
	TraceInventoryStatusInStock = "IN_STOCK"
	// TraceInventoryStatusSold 已售出。
	TraceInventoryStatusSold = "SOLD"
	// TraceInventoryStatusMisplaced 错架。
	TraceInventoryStatusMisplaced = "MISPLACED"
	// TraceInventoryStatusLossCandidate 盘亏候选。
	TraceInventoryStatusLossCandidate = "LOSS_CANDIDATE"
	// TraceInventoryStatusLost 已盘亏。
	TraceInventoryStatusLost = "LOST"
)

// 追溯日志动作类型常量。
const (
	TraceLogActionInbound    = "INBOUND"
	TraceLogActionShelving   = "SHELVING"
	TraceLogActionSale       = "SALE"
	TraceLogActionReturn     = "RETURN"
	TraceLogActionInventory  = "INVENTORY"
	TraceLogActionRelocation = "RELOCATION"
	TraceLogActionLoss       = "LOSS"
)

// InboundOrder 入库单主表。
// 对应数据库表：inbound_order。
type InboundOrder struct {
	core.BaseModel
	// OrderNo 入库单号，格式：IN-YYYYMMDD-XXXX。
	OrderNo     string  `gorm:"column:order_no;type:varchar(50);not null;uniqueIndex" json:"order_no"`
	// InvoiceNo 发票号，可为空。
	InvoiceNo   *string `gorm:"column:invoice_no;type:varchar(100)" json:"invoice_no,omitempty"`
	// OperatorID 操作人（经办人）ID，关联 sys_user。
	OperatorID  int64   `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	// CreatorID 创建人 ID，关联 sys_user。
	CreatorID   int64   `gorm:"column:creator_id;type:bigint;not null" json:"creator_id"`
	// SupplierID 供应商 ID，关联 supplier。
	SupplierID  int64   `gorm:"column:supplier_id;type:bigint;not null" json:"supplier_id"`
	// Status 单据状态：DRAFT/PENDING_CONFIRM/COMPLETED/CANCELLED。
	Status      string  `gorm:"column:status;type:varchar(20);not null;default:'DRAFT'" json:"status"`
	// TotalAmount 总金额，等于所有未删除明细金额之和。
	TotalAmount float64 `gorm:"column:total_amount;type:numeric(12,2);not null;default:0" json:"total_amount"`
	// Remark 备注。
	Remark      *string `gorm:"column:remark;type:text" json:"remark,omitempty"`

	// SubmittedAt 提交时间。
	SubmittedAt *time.Time `gorm:"column:submitted_at" json:"submitted_at,omitempty"`
	// CompletedAt 完成时间。
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	// CancelledAt 取消时间。
	CancelledAt *time.Time `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`

	// Details 入库明细，GORM 关联加载。
	Details []InboundOrderDetail `gorm:"foreignKey:OrderID;references:ID" json:"details,omitempty"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	SupplierName      string `gorm:"-" json:"supplier_name,omitempty"`
	CreatorName       string `gorm:"-" json:"creator_name,omitempty"`
	TotalPlannedQty   int32  `gorm:"-" json:"total_planned_qty"`
	TotalConfirmedQty int32  `gorm:"-" json:"total_confirmed_qty"`
}

// TableName 指定数据库表名。
func (InboundOrder) TableName() string {
	return "inbound_order"
}

// InboundOrderDetail 入库单明细行。
// 对应数据库表：inbound_order_detail。
type InboundOrderDetail struct {
	core.BaseModel
	// OrderID 所属入库单 ID。
	OrderID int64 `gorm:"column:order_id;type:bigint;not null;index" json:"order_id"`
	// DrugID 药品 ID，关联 drug_info。
	DrugID int64 `gorm:"column:drug_id;type:bigint;not null;index" json:"drug_id"`
	// BatchNumber 批号。
	BatchNumber string `gorm:"column:batch_number;type:varchar(50);not null" json:"batch_number"`
	// ExpireDate 有效期（date 类型）。
	ExpireDate time.Time `gorm:"column:expire_date;type:date;not null" json:"expire_date"`
	// PlannedQty 计划数量。
	PlannedQty int32 `gorm:"column:planned_qty;type:int;not null" json:"planned_qty"`
	// ConfirmedQty 已确认数量（扫码录入后累加）。
	ConfirmedQty int32 `gorm:"column:confirmed_qty;type:int;not null;default:0" json:"confirmed_qty"`
	// UnitPrice 单价。
	UnitPrice float64 `gorm:"column:unit_price;type:numeric(10,2);not null" json:"unit_price"`
	// Amount 金额 = PlannedQty * UnitPrice。
	Amount float64 `gorm:"column:amount;type:numeric(12,2);not null;default:0" json:"amount"`
	// Remark 备注。
	Remark *string `gorm:"column:remark;type:text" json:"remark,omitempty"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	DrugName      string `gorm:"-" json:"drug_name,omitempty"`
	Specification string `gorm:"-" json:"specification,omitempty"`
	Manufacturer  string `gorm:"-" json:"manufacturer,omitempty"`
}

// TableName 指定数据库表名。
func (InboundOrderDetail) TableName() string {
	return "inbound_order_detail"
}

// DrugTraceInventory 药品追溯库存记录。
// 对应数据库表：drug_trace_inventory。
type DrugTraceInventory struct {
	core.BaseModel
	// TraceCode 全局唯一追溯码。
	TraceCode string `gorm:"column:trace_code;type:varchar(100);not null;uniqueIndex" json:"trace_code"`
	// DrugID 药品 ID。
	DrugID int64 `gorm:"column:drug_id;type:bigint;not null;index" json:"drug_id"`
	// BatchNumber 批号。
	BatchNumber string `gorm:"column:batch_number;type:varchar(50);not null" json:"batch_number"`
	// ExpireDate 有效期。
	ExpireDate time.Time `gorm:"column:expire_date;type:date;not null" json:"expire_date"`
	// LocationID 当前货位 ID，可为 null（尚未上架）。
	LocationID *int64 `gorm:"column:location_id;type:bigint" json:"location_id,omitempty"`
	// Status 状态：PENDING/IN_STOCK/SOLD/MISPLACED/LOSS_CANDIDATE/LOST。
	Status string `gorm:"column:status;type:varchar(20);not null;index" json:"status"`
	// InboundOrderID 关联入库单 ID。
	InboundOrderID int64 `gorm:"column:inbound_order_id;type:bigint;not null;index" json:"inbound_order_id"`
	// InboundDetailID 关联入库明细 ID。
	InboundDetailID int64 `gorm:"column:inbound_detail_id;type:bigint;not null;index" json:"inbound_detail_id"`
	// SoldAt 售出时间。
	SoldAt *time.Time `gorm:"column:sold_at" json:"sold_at,omitempty"`
	// LastAction 最后操作记录。
	LastAction *string `gorm:"column:last_action;type:varchar(50)" json:"last_action,omitempty"`
}

// TableName 指定数据库表名。
func (DrugTraceInventory) TableName() string {
	return "drug_trace_inventory"
}

// DrugTraceLog 药品追溯日志。
// 对应数据库表：drug_trace_log。
type DrugTraceLog struct {
	core.BaseModel
	// TraceCode 追溯码。
	TraceCode string `gorm:"column:trace_code;type:varchar(100);not null;index" json:"trace_code"`
	// ActionType 动作类型：INBOUND/SHELVING/SALE/RETURN/INVENTORY/RELOCATION/LOSS。
	ActionType string `gorm:"column:action_type;type:varchar(20);not null" json:"action_type"`
	// FromStatus 流转前状态。
	FromStatus *string `gorm:"column:from_status;type:varchar(20)" json:"from_status,omitempty"`
	// ToStatus 流转后状态。
	ToStatus string `gorm:"column:to_status;type:varchar(20);not null" json:"to_status"`
	// OperatorID 操作人 ID。
	OperatorID int64 `gorm:"column:operator_id;type:bigint;not null" json:"operator_id"`
	// RelatedNo 关联单号（如入库单号）。
	RelatedNo *string `gorm:"column:related_no;type:varchar(100)" json:"related_no,omitempty"`
	// Remark 备注。
	Remark *string `gorm:"column:remark;type:text" json:"remark,omitempty"`
	// DrugID 药品 ID。
	DrugID int64 `gorm:"column:drug_id;type:bigint;not null" json:"drug_id"`
	// OrderID 关联入库单 ID。
	OrderID *int64 `gorm:"column:order_id;type:bigint" json:"order_id,omitempty"`
	// OrderItemID 关联入库明细 ID。
	OrderItemID *int64 `gorm:"column:order_item_id;type:bigint" json:"order_item_id,omitempty"`
	// RequestID 请求 ID，用于链路追踪。
	RequestID *string `gorm:"column:request_id;type:varchar(100)" json:"request_id,omitempty"`
	// FromLocationID 流转前货位 ID。
	FromLocationID *int64 `gorm:"column:from_location_id;type:bigint" json:"from_location_id,omitempty"`
	// ToLocationID 流转后货位 ID。
	ToLocationID *int64 `gorm:"column:to_location_id;type:bigint" json:"to_location_id,omitempty"`
}

// TableName 指定数据库表名。
func (DrugTraceLog) TableName() string {
	return "drug_trace_log"
}
