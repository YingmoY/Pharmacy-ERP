package entity

import (
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
)

// Location 是货位的领域实体，封装货位的业务规则与状态流转逻辑，不直接依赖数据库层。
type Location struct {
	ID           int64
	LocationCode string
	LocationName string
	Area         string
	Shelf        string
	Layer        string
	Position     string
	Capacity     int
	Status       int8
	Remark       *string
}

// FromModel 从 GORM 数据库模型构造货位领域实体。
func FromModel(m *model.LocationInfo) *Location {
	return &Location{
		ID:           m.ID,
		LocationCode: m.LocationCode,
		LocationName: m.LocationName,
		Area:         m.Area,
		Shelf:        m.Shelf,
		Layer:        m.Layer,
		Position:     m.Position,
		Capacity:     m.Capacity,
		Status:       m.Status,
		Remark:       m.Remark,
	}
}

// ToModel 将领域实体转换为 GORM 数据库模型，供仓储层写入使用。
// ID 为零时表示新建，GORM 将自动填充自增主键。
func (l *Location) ToModel() *model.LocationInfo {
	return &model.LocationInfo{
		LocationCode: l.LocationCode,
		LocationName: l.LocationName,
		Area:         l.Area,
		Shelf:        l.Shelf,
		Layer:        l.Layer,
		Position:     l.Position,
		Capacity:     l.Capacity,
		Status:       l.Status,
		Remark:       l.Remark,
	}
}

// IsEnabled 返回货位是否处于启用状态，停用货位不允许参与上架、盘库等业务流程。
func (l *Location) IsEnabled() bool {
	return l.Status == model.LocationStatusEnabled
}

// ToggleStatus 在启用与停用之间切换货位状态。
func (l *Location) ToggleStatus() {
	if l.Status == model.LocationStatusEnabled {
		l.Status = model.LocationStatusDisabled
	} else {
		l.Status = model.LocationStatusEnabled
	}
}

// Validate 校验货位领域实体的必填字段是否合法。
// 在执行创建或更新前调用，确保写入数据库的数据符合业务约束。
func (l *Location) Validate() error {
	if strings.TrimSpace(l.LocationCode) == "" {
		return ecode.ErrParamInvalid
	}
	if strings.TrimSpace(l.LocationName) == "" {
		return ecode.ErrParamInvalid
	}
	if strings.TrimSpace(l.Area) == "" {
		return ecode.ErrParamInvalid
	}
	return nil
}
