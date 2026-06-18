package repository

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"gorm.io/gorm"
)

// LocationFilter 货位列表查询过滤条件。
type LocationFilter struct {
	Keyword  string // 编码/名称/区域模糊搜索
	Status   *int8  // nil 表示不过滤
	Area     string // 按区域过滤
	Page     int
	PageSize int
}

// LocationRepository 定义货位仓储接口。
type LocationRepository interface {
	// FindByID 根据主键查询货位。
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.LocationInfo, error)
	// FindByCode 根据货位编码查询。
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.LocationInfo, error)
	// List 分页查询货位列表。
	List(ctx context.Context, db *gorm.DB, filter LocationFilter) ([]*model.LocationInfo, int64, error)
	// Create 新增货位记录。
	Create(ctx context.Context, db *gorm.DB, location *model.LocationInfo) error
	// Update 按字段 map 更新货位。
	Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error
	// UpdateStatus 仅更新状态字段。
	UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error
	// Delete 软删除货位记录。
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	// GetAreas 查询所有不重复的区域值。
	GetAreas(ctx context.Context, db *gorm.DB) ([]string, error)
	// HasActiveInventory 检查货位是否存在在库/错架/盘亏候选状态的追溯库存。
	HasActiveInventory(ctx context.Context, db *gorm.DB, id int64) (bool, error)
}

type locationRepo struct{}

// NewLocationRepo 创建货位仓储实现。
func NewLocationRepo() LocationRepository {
	return &locationRepo{}
}

// FindByID 根据主键查询货位。
func (r *locationRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.LocationInfo, error) {
	var loc model.LocationInfo
	if err := db.WithContext(ctx).Where("id = ?", id).First(&loc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrLocationNotFound
		}
		return nil, err
	}
	return &loc, nil
}

// FindByCode 根据货位编码查询。
func (r *locationRepo) FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.LocationInfo, error) {
	var loc model.LocationInfo
	if err := db.WithContext(ctx).Where("location_code = ?", code).First(&loc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrLocationNotFound
		}
		return nil, err
	}
	return &loc, nil
}

// List 分页查询货位列表。
func (r *locationRepo) List(ctx context.Context, db *gorm.DB, filter LocationFilter) ([]*model.LocationInfo, int64, error) {
	q := db.WithContext(ctx).Model(&model.LocationInfo{})

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("location_code ILIKE ? OR location_name ILIKE ? OR area ILIKE ?", like, like, like)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.Area != "" {
		q = q.Where("area = ?", filter.Area)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	var list []*model.LocationInfo
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Create 新增货位记录。
func (r *locationRepo) Create(ctx context.Context, db *gorm.DB, location *model.LocationInfo) error {
	return db.WithContext(ctx).Create(location).Error
}

// Update 按字段 map 更新货位。
func (r *locationRepo) Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.WithContext(ctx).Model(&model.LocationInfo{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStatus 仅更新状态字段。
func (r *locationRepo) UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error {
	return db.WithContext(ctx).Model(&model.LocationInfo{}).Where("id = ?", id).
		Update("status", status).Error
}

// Delete 软删除货位。
func (r *locationRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&model.LocationInfo{}, id).Error
}

// GetAreas 查询所有不重复的区域值（仅启用状态）。
func (r *locationRepo) GetAreas(ctx context.Context, db *gorm.DB) ([]string, error) {
	var areas []string
	if err := db.WithContext(ctx).Model(&model.LocationInfo{}).
		Distinct("area").
		Where("area != '' AND area IS NOT NULL").
		Order("area ASC").
		Pluck("area", &areas).Error; err != nil {
		return nil, err
	}
	return areas, nil
}

// HasActiveInventory 检查货位是否存在在库/错架/盘亏候选状态追溯库存。
func (r *locationRepo) HasActiveInventory(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Table("drug_trace_inventory").
		Where("location_id = ? AND status IN ? AND deleted_at IS NULL",
			id, []string{"IN_STOCK", "MISPLACED", "LOSS_CANDIDATE"}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
