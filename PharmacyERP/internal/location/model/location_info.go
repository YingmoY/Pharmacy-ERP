package model

import "github.com/YingmoY/PharmacyERP/internal/pkg/core"

const (
	// LocationStatusDisabled 意味着该位置不可用于操作。
	LocationStatusDisabled int8 = 0
	// LocationStatusEnabled 意味着该位置可用于操作。
	LocationStatusEnabled int8 = 1
)

// LocationInfo 代表仓库中一个具体位置的信息，对应数据库表 location_info。
type LocationInfo struct {
	core.BaseModel
	LocationCode string  `gorm:"column:location_code;uniqueIndex;size:50;not null" json:"location_code"`
	LocationName string  `gorm:"column:location_name;size:100;not null" json:"location_name"`
	Area         string  `gorm:"column:area;size:20;not null" json:"area"`
	Shelf        string  `gorm:"column:shelf;size:20" json:"shelf"`
	Layer        string  `gorm:"column:layer;size:20" json:"layer"`
	Position     string  `gorm:"column:position;size:20" json:"position"`
	Capacity     int     `gorm:"column:capacity;default:0" json:"capacity"`
	Status       int8    `gorm:"column:status;default:1" json:"status"`
	Remark       *string `gorm:"column:remark;type:text" json:"remark,omitempty"`
}

// TableName overrides GORM pluralized naming.
func (LocationInfo) TableName() string {
	return "location_info"
}
