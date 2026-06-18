package repository

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/supplier/model"
	"gorm.io/gorm"
)

// SupplierFilter 供应商列表查询过滤条件。
type SupplierFilter struct {
	Keyword  string // 名称/编码模糊搜索
	Status   *int8  // nil 表示不过滤
	Page     int
	PageSize int
}

// SupplierRepository 定义供应商仓储接口。
type SupplierRepository interface {
	// FindByID 根据主键查询，未找到返回错误。
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.Supplier, error)
	// FindByCode 根据供应商编码查询。
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.Supplier, error)
	// List 分页查询供应商列表。
	List(ctx context.Context, db *gorm.DB, filter SupplierFilter) ([]*model.Supplier, int64, error)
	// Create 新增供应商记录。
	Create(ctx context.Context, db *gorm.DB, supplier *model.Supplier) error
	// Update 按字段 map 更新供应商。
	Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error
	// UpdateStatus 仅更新状态字段。
	UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error
	// Delete 软删除供应商记录。
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	// HasInboundOrders 检查供应商是否有关联的入库单记录。
	HasInboundOrders(ctx context.Context, db *gorm.DB, id int64) (bool, error)
}

type supplierRepo struct{}

// NewSupplierRepo 创建供应商仓储实现。
func NewSupplierRepo() SupplierRepository {
	return &supplierRepo{}
}

// FindByID 根据主键查询供应商。
func (r *supplierRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*model.Supplier, error) {
	var s model.Supplier
	if err := db.WithContext(ctx).Where("id = ?", id).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.New(40402, "supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

// FindByCode 根据供应商编码查询。
func (r *supplierRepo) FindByCode(ctx context.Context, db *gorm.DB, code string) (*model.Supplier, error) {
	var s model.Supplier
	if err := db.WithContext(ctx).Where("supplier_code = ?", code).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.New(40402, "supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

// List 分页查询供应商列表。
func (r *supplierRepo) List(ctx context.Context, db *gorm.DB, filter SupplierFilter) ([]*model.Supplier, int64, error) {
	q := db.WithContext(ctx).Model(&model.Supplier{})

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ? OR supplier_code ILIKE ? OR contact_name ILIKE ?", like, like, like)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
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

	var list []*model.Supplier
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Create 新增供应商记录。
func (r *supplierRepo) Create(ctx context.Context, db *gorm.DB, supplier *model.Supplier) error {
	return db.WithContext(ctx).Create(supplier).Error
}

// Update 按字段 map 更新供应商。
func (r *supplierRepo) Update(ctx context.Context, db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.WithContext(ctx).Model(&model.Supplier{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStatus 仅更新状态字段。
func (r *supplierRepo) UpdateStatus(ctx context.Context, db *gorm.DB, id int64, status int8) error {
	return db.WithContext(ctx).Model(&model.Supplier{}).Where("id = ?", id).
		Update("status", status).Error
}

// Delete 软删除供应商。
func (r *supplierRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&model.Supplier{}, id).Error
}

// HasInboundOrders 检查供应商是否有关联入库单。
// 当前 inbound_order 表使用 supplier 文本字段而非外键，
// 因此此处通过 supplier_code 做文本匹配检索。
// 如果后续改为外键关联，修改此处逻辑即可。
func (r *supplierRepo) HasInboundOrders(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	// 先获取供应商编码
	var s model.Supplier
	if err := db.WithContext(ctx).Select("supplier_code").Where("id = ?", id).First(&s).Error; err != nil {
		return false, err
	}

	var count int64
	if err := db.WithContext(ctx).Table("inbound_order").
		Where("supplier = ? AND deleted_at IS NULL", s.SupplierCode).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
