// Package repository 定义 AI 发票仓储层接口及 GORM 实现。
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/ai/model"
	"gorm.io/gorm"
)

// AIInvoiceListFilter AI 发票列表查询条件。
type AIInvoiceListFilter struct {
	// Status 按状态过滤（空则不过滤）。
	Status string
	// SupplierID 按已匹配供应商 ID 过滤（matched_supplier_id，0 则不过滤）。
	SupplierID int64
	// Page 页码（从 1 开始）。
	Page int
	// PageSize 每页条数（默认 20）。
	PageSize int
	// StartDate 创建时间起始。
	StartDate *time.Time
	// EndDate 创建时间结束。
	EndDate *time.Time
}

// AIInvoiceRepo 定义 AI 发票记录的读写能力。
type AIInvoiceRepo interface {
	// Create 创建 AI 发票记录。
	Create(ctx context.Context, record *model.AIInvoiceRecord) error
	// FindByID 按主键查询。
	FindByID(ctx context.Context, id int64) (*model.AIInvoiceRecord, error)
	// List 分页查询列表，返回数据和总条数。
	List(ctx context.Context, filter AIInvoiceListFilter) ([]*model.AIInvoiceRecord, int64, error)
	// UpdateStatus 更新状态及可选字段（识别结果等）。
	UpdateStatus(ctx context.Context, id int64, updates map[string]interface{}) error
	// SetInboundOrder 写回关联入库单 ID 及转换时间。
	SetInboundOrder(ctx context.Context, id, inboundOrderID int64) error
}

// aiInvoiceRepo 是 AIInvoiceRepo 的 GORM 实现。
type aiInvoiceRepo struct {
	db *gorm.DB
}

// NewAIInvoiceRepo 创建仓储实例。
func NewAIInvoiceRepo(db *gorm.DB) AIInvoiceRepo {
	return &aiInvoiceRepo{db: db}
}

// Create 插入一条 AI 发票记录。
func (r *aiInvoiceRepo) Create(ctx context.Context, record *model.AIInvoiceRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// FindByID 按主键查询 AI 发票记录。
func (r *aiInvoiceRepo) FindByID(ctx context.Context, id int64) (*model.AIInvoiceRecord, error) {
	var record model.AIInvoiceRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &record, nil
}

// List 分页查询 AI 发票列表。
func (r *aiInvoiceRepo) List(ctx context.Context, filter AIInvoiceListFilter) ([]*model.AIInvoiceRecord, int64, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.AIInvoiceRecord{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.SupplierID > 0 {
		query = query.Where("matched_supplier_id = ?", filter.SupplierID)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []*model.AIInvoiceRecord
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// UpdateStatus 更新 AI 发票状态及其他识别字段。
func (r *aiInvoiceRepo) UpdateStatus(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.AIInvoiceRecord{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// SetInboundOrder 将识别记录的 inbound_order_id 和 converted_at 写回。
func (r *aiInvoiceRepo) SetInboundOrder(ctx context.Context, id, inboundOrderID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.AIInvoiceRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"inbound_order_id": inboundOrderID,
			"converted_at":     now,
		}).Error
}
