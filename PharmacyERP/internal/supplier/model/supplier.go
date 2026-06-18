package model

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// SupplierStatusDisabled 供应商停用状态。
	SupplierStatusDisabled int8 = 0
	// SupplierStatusEnabled 供应商启用状态。
	SupplierStatusEnabled int8 = 1
)

// Supplier 代表供应商基本信息，对应数据库表 supplier。
// 字段按实际表结构映射（contact_name / contact_phone / license_no）。
type Supplier struct {
	core.BaseModel
	SupplierCode string  `gorm:"column:supplier_code;uniqueIndex;size:50;not null" json:"supplier_code"`
	Name         string  `gorm:"column:name;size:100;not null" json:"name"`
	ContactName  string  `gorm:"column:contact_name;size:50" json:"contact_name"`
	ContactPhone string  `gorm:"column:contact_phone;size:30" json:"contact_phone"`
	LicenseNo    string  `gorm:"column:license_no;size:100" json:"license_no"`
	Address      string  `gorm:"column:address;size:200" json:"address"`
	Status       int8    `gorm:"column:status;default:1" json:"status"`
	Remark       *string `gorm:"column:remark;type:text" json:"remark,omitempty"`
}

// TableName 返回供应商表名。
func (Supplier) TableName() string {
	return "supplier"
}
