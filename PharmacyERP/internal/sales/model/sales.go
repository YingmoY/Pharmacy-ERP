package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// 销售订单状态常量
const (
	// SalesOrderStatusPending 待结算（无需审核的订单初始状态）
	SalesOrderStatusPending = "PENDING"
	// SalesOrderStatusPendingReview 待审核（处方药或需审核订单的初始状态）
	SalesOrderStatusPendingReview = "PENDING_REVIEW"
	// SalesOrderStatusApproved 审核通过，可结算
	SalesOrderStatusApproved = "APPROVED"
	// SalesOrderStatusCompleted 已完成（已付款）
	SalesOrderStatusCompleted = "COMPLETED"
	// SalesOrderStatusPartiallyRefunded 部分退款
	SalesOrderStatusPartiallyRefunded = "PARTIALLY_REFUNDED"
	// SalesOrderStatusRefunded 全额退款
	SalesOrderStatusRefunded = "REFUNDED"
	// SalesOrderStatusCancelled 已取消
	SalesOrderStatusCancelled = "CANCELLED"
)

// 退款模式常量
const (
	// RefundModeFull 全额退款
	RefundModeFull = "FULL"
	// RefundModePartial 部分退款
	RefundModePartial = "PARTIAL"
)

// 销售订单明细退款状态常量
const (
	// ItemRefundStatusNone 未退款
	ItemRefundStatusNone = "NONE"
	// ItemRefundStatusRefunded 已退款
	ItemRefundStatusRefunded = "REFUNDED"
)

// 预留状态常量
const (
	// ReservationStatusReserved 已预留
	ReservationStatusReserved = "RESERVED"
	// ReservationStatusReleased 已释放
	ReservationStatusReleased = "RELEASED"
	// ReservationStatusConsumed 已消费（结算完成）
	ReservationStatusConsumed = "CONSUMED"
	// ReservationStatusExpired 已过期
	ReservationStatusExpired = "EXPIRED"
)

// SalesOrder 映射 public.sales_order 表
type SalesOrder struct {
	core.BaseModel
	OrderNo               string           `gorm:"column:order_no;type:varchar(50);not null;uniqueIndex"                  json:"order_no"`
	CashierID             int64            `gorm:"column:cashier_id;type:bigint;not null;index"                           json:"cashier_id"`
	CashierName           string           `gorm:"-"                                                                      json:"cashier_name,omitempty"`
	TotalAmount           float64          `gorm:"column:total_amount;type:numeric(10,2);not null;default:0"              json:"total_amount"`
	MedicareAmount        float64          `gorm:"column:medicare_amount;type:numeric(10,2);not null;default:0"           json:"medicare_amount"`
	PersonalAmount        float64          `gorm:"column:personal_amount;type:numeric(10,2);not null;default:0"           json:"personal_amount"`
	NeedAudit             bool             `gorm:"column:need_audit;not null;default:false"                               json:"need_audit"`
	NeedMedicare          bool             `gorm:"column:need_medicare;not null;default:false"                            json:"need_medicare"`
	Status                string           `gorm:"column:status;type:varchar(30);not null;index"                          json:"status"`
	MedicareTransactionID *string          `gorm:"column:medicare_transaction_id;type:varchar(100)"                       json:"medicare_transaction_id,omitempty"`
	MdtrtID               *string          `gorm:"column:mdtrt_id;type:varchar(100)"                                      json:"mdtrt_id,omitempty"`
	CustomerName          *string          `gorm:"column:customer_name;type:varchar(100)"                                 json:"customer_name,omitempty"`
	IsPrescription        bool             `gorm:"column:is_prescription;not null;default:false"                          json:"is_prescription"`
	DiscountAmount        float64          `gorm:"column:discount_amount;type:numeric(10,2);not null;default:0"           json:"discount_amount"`
	ActualAmount          float64          `gorm:"column:actual_amount;type:numeric(10,2);not null;default:0"             json:"actual_amount"`
	PaymentMethod         *string          `gorm:"column:payment_method;type:varchar(30)"                                 json:"payment_method,omitempty"`
	PaidAt                *time.Time       `gorm:"column:paid_at;type:timestamptz"                                        json:"paid_at,omitempty"`
	CancelledAt           *time.Time       `gorm:"column:cancelled_at;type:timestamptz"                                   json:"cancelled_at,omitempty"`
	RefundedAt            *time.Time       `gorm:"column:refunded_at;type:timestamptz"                                    json:"refunded_at,omitempty"`
	RefundAmount          float64          `gorm:"column:refund_amount;type:numeric(10,2);not null;default:0"             json:"refund_amount"`
	RefundReason          *string          `gorm:"column:refund_reason;type:text"                                         json:"refund_reason,omitempty"`
	Remark                *string          `gorm:"column:remark;type:text"                                                json:"remark,omitempty"`
	Items                 []SalesOrderItem `gorm:"foreignKey:OrderID;references:ID"                                       json:"items,omitempty"`
}

// TableName 返回数据库表名
func (SalesOrder) TableName() string {
	return "sales_order"
}

// SalesOrderItem 映射 public.sales_order_item 表
type SalesOrderItem struct {
	core.BaseModel
	OrderID          int64      `gorm:"column:order_id;type:bigint;not null;index"                          json:"order_id"`
	DrugID           int64      `gorm:"column:drug_id;type:bigint;not null;index"                           json:"drug_id"`
	TraceCode        string     `gorm:"column:trace_code;type:varchar(100);not null"                        json:"trace_code"`
	Price            float64    `gorm:"column:price;type:numeric(10,2);not null"                            json:"price"`
	Quantity         int32      `gorm:"column:quantity;type:int;not null;default:1"                         json:"quantity"`
	SubtotalAmount   float64    `gorm:"column:subtotal_amount;type:numeric(10,2);not null"                  json:"subtotal_amount"`
	Remark           *string    `gorm:"column:remark;type:text"                                             json:"remark,omitempty"`
	RefundStatus     string     `gorm:"column:refund_status;type:varchar(20);not null;default:'NONE'"       json:"refund_status"`
	RefundAmount     float64    `gorm:"column:refund_amount;type:numeric(10,2);not null;default:0"          json:"refund_amount"`
	RefundedAt       *time.Time `gorm:"column:refunded_at;type:timestamptz"                                 json:"refunded_at,omitempty"`
	RefundReason     *string    `gorm:"column:refund_reason;type:text"                                      json:"refund_reason,omitempty"`
	RefundOperatorID *int64     `gorm:"column:refund_operator_id;type:bigint"                               json:"refund_operator_id,omitempty"`
	// 虚拟字段：通过关联查询填充
	UnitPrice      float64 `gorm:"-" json:"unit_price"`
	DrugName       string  `gorm:"-" json:"drug_name,omitempty"`
	Specification  string  `gorm:"-" json:"specification,omitempty"`
	Manufacturer   string  `gorm:"-" json:"manufacturer,omitempty"`
	BatchNumber    string  `gorm:"-" json:"batch_number,omitempty"`
	ExpireDate     string  `gorm:"-" json:"expire_date,omitempty"`
	LocationCode   string  `gorm:"-" json:"location_code,omitempty"`
	IsPrescription bool    `gorm:"-" json:"is_prescription"`
}

// TableName 返回数据库表名
func (SalesOrderItem) TableName() string {
	return "sales_order_item"
}

// TraceReservation 映射 public.trace_reservation 表
type TraceReservation struct {
	core.BaseModel
	ReservationNo    string     `gorm:"column:reservation_no;type:varchar(50);not null;uniqueIndex"         json:"reservation_no"`
	SalesOrderID     int64      `gorm:"column:sales_order_id;type:bigint;not null;index"                    json:"sales_order_id"`
	SalesOrderItemID *int64     `gorm:"column:sales_order_item_id;type:bigint"                              json:"sales_order_item_id,omitempty"`
	TraceCode        string     `gorm:"column:trace_code;type:varchar(100);not null;index"                  json:"trace_code"`
	DrugID           int64      `gorm:"column:drug_id;type:bigint;not null"                                 json:"drug_id"`
	ReservedBy       int64      `gorm:"column:reserved_by;type:bigint;not null"                             json:"reserved_by"`
	Status           string     `gorm:"column:status;type:varchar(20);not null;default:'RESERVED';index"    json:"status"`
	ReservedAt       time.Time  `gorm:"column:reserved_at;type:timestamptz;not null"                        json:"reserved_at"`
	ReleasedAt       *time.Time `gorm:"column:released_at;type:timestamptz"                                 json:"released_at,omitempty"`
	ConfirmedAt      *time.Time `gorm:"column:confirmed_at;type:timestamptz"                                json:"confirmed_at,omitempty"`
	ExpireAt         time.Time  `gorm:"column:expire_at;type:timestamptz;not null"                          json:"expire_at"`
	Remark           *string    `gorm:"column:remark;type:text"                                             json:"remark,omitempty"`
}

// TableName 返回数据库表名
func (TraceReservation) TableName() string {
	return "trace_reservation"
}
