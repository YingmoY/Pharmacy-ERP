package entity

import (
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/drug/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
)

// Drug 是药品的领域实体，封装业务规则与状态流转逻辑，不直接依赖数据库层。
type Drug struct {
	ID               int64
	DrugCode         string
	CommonName       string
	TradeName        string
	Specification    string
	DosageForm       string
	Manufacturer     string
	ApprovalNumber   string
	IsPrescription   bool
	IsMedicare       bool
	Status           int8
	Barcode          *string
	Unit             string
	RetailPrice      *float64
	PurchasePrice    *float64
	StorageCondition string
	Remark           string
}

// FromModel 从 GORM 数据库模型构造药品领域实体。
func FromModel(m *model.DrugInfo) *Drug {
	return &Drug{
		ID:               m.ID,
		DrugCode:         m.DrugCode,
		CommonName:       m.CommonName,
		TradeName:        m.TradeName,
		Specification:    m.Specification,
		DosageForm:       m.DosageForm,
		Manufacturer:     m.Manufacturer,
		ApprovalNumber:   m.ApprovalNumber,
		IsPrescription:   m.IsPrescription,
		IsMedicare:       m.IsMedicare,
		Status:           m.Status,
		Barcode:          m.Barcode,
		Unit:             m.Unit,
		RetailPrice:      m.RetailPrice,
		PurchasePrice:    m.PurchasePrice,
		StorageCondition: m.StorageCondition,
		Remark:           m.Remark,
	}
}

// ToModel 将领域实体转换为 GORM 数据库模型，供仓储层写入使用。
// ID 为零时表示新建，GORM 将自动填充自增主键。
func (d *Drug) ToModel() *model.DrugInfo {
	return &model.DrugInfo{
		DrugCode:         d.DrugCode,
		CommonName:       d.CommonName,
		TradeName:        d.TradeName,
		Specification:    d.Specification,
		DosageForm:       d.DosageForm,
		Manufacturer:     d.Manufacturer,
		ApprovalNumber:   d.ApprovalNumber,
		IsPrescription:   d.IsPrescription,
		IsMedicare:       d.IsMedicare,
		Status:           d.Status,
		Barcode:          d.Barcode,
		Unit:             d.Unit,
		RetailPrice:      d.RetailPrice,
		PurchasePrice:    d.PurchasePrice,
		StorageCondition: d.StorageCondition,
		Remark:           d.Remark,
	}
}

// IsEnabled 返回药品是否处于启用状态，停用药品不允许参与入库、销售等业务流程。
func (d *Drug) IsEnabled() bool {
	return d.Status == model.DrugStatusEnabled
}

// NeedsAudit 返回该药品在销售时是否需要药师审核（即为处方药）。
func (d *Drug) NeedsAudit() bool {
	return d.IsPrescription
}

// ToggleStatus 在启用与停用之间切换药品状态。
func (d *Drug) ToggleStatus() {
	if d.Status == model.DrugStatusEnabled {
		d.Status = model.DrugStatusDisabled
	} else {
		d.Status = model.DrugStatusEnabled
	}
}

// Validate 校验药品领域实体的必填字段是否合法。
// 在执行创建或更新前调用，确保写入数据库的数据符合业务约束。
func (d *Drug) Validate() error {
	if strings.TrimSpace(d.DrugCode) == "" {
		return ecode.ErrParamInvalid
	}
	if strings.TrimSpace(d.CommonName) == "" {
		return ecode.ErrParamInvalid
	}
	if strings.TrimSpace(d.Specification) == "" {
		return ecode.ErrParamInvalid
	}
	if strings.TrimSpace(d.Manufacturer) == "" {
		return ecode.ErrParamInvalid
	}
	return nil
}
