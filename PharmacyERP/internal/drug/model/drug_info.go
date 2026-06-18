package model

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// DrugStatusDisabled 代表药品不可用，可能是因为停产、缺货等原因。
	DrugStatusDisabled int8 = 0
	// DrugStatusEnabled 代表药品可用，正常供应。
	DrugStatusEnabled int8 = 1
)

// DrugInfo 代表药品的基本信息，包括药品编码、通用名称、商品名称、规格、剂型、生产厂家、批准文号、处方药标识、医保标识和状态等字段。
type DrugInfo struct {
	core.BaseModel
	DrugCode         string   `gorm:"column:drug_code;uniqueIndex;size:50;not null" json:"drug_code"`
	CommonName       string   `gorm:"column:common_name;size:100;not null" json:"common_name"`
	TradeName        string   `gorm:"column:trade_name;size:100" json:"trade_name"`
	Specification    string   `gorm:"column:specification;size:50;not null" json:"specification"`
	DosageForm       string   `gorm:"column:dosage_form;size:50" json:"dosage_form"`
	Manufacturer     string   `gorm:"column:manufacturer;size:100;not null" json:"manufacturer"`
	ApprovalNumber   string   `gorm:"column:approval_number;size:50" json:"approval_number"`
	IsPrescription      bool    `gorm:"column:is_prescription;default:false" json:"is_prescription"`
	IsMedicare          bool    `gorm:"column:is_medicare;default:false" json:"is_medicare"`
	MedListCodg         string  `gorm:"column:med_list_codg;size:50" json:"med_list_codg"`
	MedinsListCodg      string  `gorm:"column:medins_list_codg;size:50" json:"medins_list_codg"`
	FixmedinsHilistId   string  `gorm:"column:fixmedins_hilist_id;size:50" json:"fixmedins_hilist_id"`
	FixmedinsHilistName string  `gorm:"column:fixmedins_hilist_name;size:100" json:"fixmedins_hilist_name"`
	Status              int8    `gorm:"column:status;default:1" json:"status"`
	Barcode          *string  `gorm:"column:barcode;uniqueIndex;size:100" json:"barcode,omitempty"`
	Unit             string   `gorm:"column:unit;size:20" json:"unit"`
	RetailPrice      *float64 `gorm:"column:retail_price;type:numeric(10,2)" json:"retail_price,omitempty"`
	PurchasePrice    *float64 `gorm:"column:purchase_price;type:numeric(10,2)" json:"purchase_price,omitempty"`
	StorageCondition string   `gorm:"column:storage_condition;size:200" json:"storage_condition"`
	Remark           string   `gorm:"column:remark;type:text" json:"remark"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	InStockCount int64 `gorm:"-" json:"in_stock_count"`
}

// TableName 返回 DrugInfo 结构体对应的数据库表名 "drug_info"。这个方法是 GORM 框架用来映射结构体到数据库表的约定。
func (DrugInfo) TableName() string {
	return "drug_info"
}
