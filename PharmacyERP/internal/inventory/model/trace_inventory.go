package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// TraceInventoryStatusPending 待上架（入库确认后）。
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

// TraceInventory 表示按唯一药品追溯码管理的核心库存表。
// 对应数据库表：drug_trace_inventory。
type TraceInventory struct {
	core.BaseModel
	TraceCode       string    `gorm:"column:trace_code;type:varchar(100);not null;uniqueIndex" json:"trace_code"`
	DrugID          int64     `gorm:"column:drug_id;type:bigint;not null;index" json:"drug_id"`
	BatchNumber     string    `gorm:"column:batch_number;type:varchar(50);not null" json:"batch_number"`
	ExpireDate      time.Time `gorm:"column:expire_date;type:date;not null" json:"expire_date"`
	LocationID      *int64    `gorm:"column:location_id;type:bigint" json:"location_id,omitempty"`
	Status          string    `gorm:"column:status;type:varchar(20);not null;index" json:"status"`
	InboundOrderID  int64     `gorm:"column:inbound_order_id;type:bigint;not null;index" json:"inbound_order_id"`
	InboundDetailID int64     `gorm:"column:inbound_detail_id;type:bigint;not null;index" json:"inbound_detail_id"`

	// SoldAt 售出时间。
	SoldAt *time.Time `gorm:"column:sold_at;type:timestamptz" json:"sold_at,omitempty"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	DrugName           string `gorm:"-" json:"drug_name,omitempty"`
	Specification      string `gorm:"-" json:"specification,omitempty"`
	Manufacturer       string `gorm:"-" json:"manufacturer,omitempty"`
	LocationCode       string `gorm:"-" json:"location_code,omitempty"`
	SystemLocationCode string `gorm:"-" json:"system_location_code,omitempty"`
	ScannedLocationCode string `gorm:"-" json:"scanned_location_code,omitempty"`
	InboundOrderNo     string `gorm:"-" json:"inbound_order_no,omitempty"`
}

func (TraceInventory) TableName() string {
	return "drug_trace_inventory"
}

// IsInStock 判断当前追溯码药品是否处于在库状态。
func (t TraceInventory) IsInStock() bool {
	return t.Status == TraceInventoryStatusInStock
}

// IsAbnormal 判断当前库存状态是否属于异常（错放或疑似丢失）。
func (t TraceInventory) IsAbnormal() bool {
	return t.Status == TraceInventoryStatusMisplaced || t.Status == TraceInventoryStatusLossCandidate || t.Status == TraceInventoryStatusLost
}
