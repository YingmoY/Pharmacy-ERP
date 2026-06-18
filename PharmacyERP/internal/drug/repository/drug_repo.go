package repository

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/drug/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"gorm.io/gorm"
)

// DrugFilter 药品列表查询过滤条件。
type DrugFilter struct {
	Keyword        string // 通用名/商品名/药品编码/条形码模糊搜索
	Status         *int8  // 状态过滤，nil 表示不过滤
	IsPrescription *bool  // 处方药过滤，nil 表示不过滤
	Page           int
	Manufacturer   string
	IsMedicare     *bool
	PageSize       int
}

// DrugRepository 定义药品仓储接口。
type DrugRepository interface {
	// FindByID 根据主键查询药品，未找到返回 ErrTraceCodeNotFound。
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.DrugInfo, error)
	// FindByCode 根据药品编码查询，未找到返回 ErrTraceCodeNotFound。
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.DrugInfo, error)
	// List 分页查询药品列表，返回数据与总数。
	List(ctx context.Context, db *gorm.DB, filter DrugFilter) ([]*model.DrugInfo, int64, error)
	// Create 新增药品记录。
	Create(ctx context.Context, db *gorm.DB, drug *model.DrugInfo) error
	// Update 按字段 map 更新指定药品（只更新有值的字段）。
	Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error
	// UpdateStatus 仅更新药品状态字段。
	UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error
	// Delete 软删除药品记录。
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	// HasInventory 检查药品是否存在 drug_trace_inventory 记录（任何状态）。
	HasInventory(ctx context.Context, db *gorm.DB, id int64) (bool, error)
	// HasSalesItems 检查药品是否存在 sales_order_item 记录。
	HasSalesItems(ctx context.Context, db *gorm.DB, id int64) (bool, error)
	// HasInboundDetails 检查药品是否存在 inbound_order_detail 记录。
	HasInboundDetails(ctx context.Context, db *gorm.DB, id int64) (bool, error)
}

type drugRepo struct{}

// NewDrugRepo 创建药品仓储实现。
func NewDrugRepo() DrugRepository {
	return &drugRepo{}
}

// FindByID 根据主键查询药品。
func (r *drugRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.DrugInfo, error) {
	var drug model.DrugInfo
	if err := db.WithContext(ctx).Where("id = ?", id).First(&drug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.New(40401, "drug not found")
		}
		return nil, err
	}
	return &drug, nil
}

// FindByCode 根据药品编码查询。
func (r *drugRepo) FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.DrugInfo, error) {
	var drug model.DrugInfo
	if err := db.WithContext(ctx).Where("drug_code = ?", code).First(&drug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.New(40401, "drug not found")
		}
		return nil, err
	}
	return &drug, nil
}

// List 分页查询药品列表。
func (r *drugRepo) List(ctx context.Context, db *gorm.DB, filter DrugFilter) ([]*model.DrugInfo, int64, error) {
	q := db.WithContext(ctx).Model(&model.DrugInfo{})

	// 关键词搜索：通用名、商品名、药品编码、条形码
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("common_name ILIKE ? OR trade_name ILIKE ? OR drug_code ILIKE ? OR barcode ILIKE ?",
			like, like, like, like)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.IsPrescription != nil {
		q = q.Where("is_prescription = ?", *filter.IsPrescription)
	}
	if filter.Manufacturer != "" {
		q = q.Where("manufacturer ILIKE ?", "%"+filter.Manufacturer+"%")
	}
	if filter.IsMedicare != nil {
		q = q.Where("is_medicare = ?", *filter.IsMedicare)
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

	var list []*model.DrugInfo
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Create 新增药品记录。
func (r *drugRepo) Create(ctx context.Context, db *gorm.DB, drug *model.DrugInfo) error {
	return db.WithContext(ctx).Create(drug).Error
}

// Update 按字段 map 更新药品。
func (r *drugRepo) Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.WithContext(ctx).Model(&model.DrugInfo{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStatus 仅更新状态字段。
func (r *drugRepo) UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error {
	return db.WithContext(ctx).Model(&model.DrugInfo{}).Where("id = ?", id).
		Update("status", status).Error
}

// Delete 软删除药品。
func (r *drugRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&model.DrugInfo{}, id).Error
}

// HasInventory 检查 drug_trace_inventory 中是否存在该药品的追溯记录。
func (r *drugRepo) HasInventory(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Table("drug_trace_inventory").
		Where("drug_id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasSalesItems 检查 sales_order_item 中是否存在该药品的销售记录。
func (r *drugRepo) HasSalesItems(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Table("sales_order_item").
		Where("drug_id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasInboundDetails 检查 inbound_order_detail 中是否存在该药品的入库明细。
func (r *drugRepo) HasInboundDetails(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Table("inbound_order_detail").
		Where("drug_id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
