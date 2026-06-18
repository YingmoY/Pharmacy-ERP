// Package repository 定义入库单仓储层接口及其 GORM 实现。
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/inbound/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ListFilter 入库单列表查询条件。
type ListFilter struct {
	// Status 按状态过滤（可为空）。
	Status string
	// SupplierID 按供应商 ID 过滤（0 表示不过滤）。
	SupplierID int64
	// StartDate 创建时间起始（零值表示不过滤）。
	StartDate *time.Time
	// EndDate 创建时间结束（零值表示不过滤）。
	EndDate *time.Time
	// Keyword 关键字搜索（匹配单号、发票号）。
	Keyword string
	// Page 页码（从 1 开始）。
	Page int
	// PageSize 每页大小（默认 20）。
	PageSize int
}

// InboundRepo 定义入库单全量读写能力。
type InboundRepo interface {
	// GetByID 按 ID 查询入库单（含关联明细）。
	GetByID(ctx context.Context, id int64) (*model.InboundOrder, error)
	// GetByIDForUpdate 加行锁按 ID 查询入库单（用于事务内状态变更）。
	GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.InboundOrder, error)
	// GetByOrderNo 按单号查询入库单。
	GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error)

	// List 分页列表查询，返回数据和总条数。
	List(ctx context.Context, filter ListFilter) ([]*model.InboundOrder, int64, error)

	// Create 新建入库单（不含明细）。
	Create(ctx context.Context, tx *gorm.DB, order *model.InboundOrder) error
	// Update 更新可编辑字段（supplier_id、invoice_no、remark、operator_id）。
	Update(ctx context.Context, tx *gorm.DB, order *model.InboundOrder) error
	// UpdateStatus 更新入库单状态。
	UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string) error
	// UpdateTotalAmount 重新计算并写入 total_amount。
	UpdateTotalAmount(ctx context.Context, tx *gorm.DB, orderID int64) error

	// GetDetails 获取某入库单所有未删除明细。
	GetDetails(ctx context.Context, orderID int64) ([]*model.InboundOrderDetail, error)
	// GetDetailByID 按 ID 获取明细（校验 order_id 归属）。
	GetDetailByID(ctx context.Context, orderID, detailID int64) (*model.InboundOrderDetail, error)
	// GetDetailByIDForUpdate 加行锁按 ID 获取明细（事务内使用）。
	GetDetailByIDForUpdate(ctx context.Context, tx *gorm.DB, orderID, detailID int64) (*model.InboundOrderDetail, error)

	// AddDetail 新增明细行。
	AddDetail(ctx context.Context, tx *gorm.DB, detail *model.InboundOrderDetail) error
	// UpdateDetail 更新明细行可编辑字段。
	UpdateDetail(ctx context.Context, tx *gorm.DB, detail *model.InboundOrderDetail) error
	// DeleteDetail 软删除明细行。
	DeleteDetail(ctx context.Context, tx *gorm.DB, orderID, detailID int64) error

	// IncreaseConfirmedQty 原子累加明细已确认数量。
	IncreaseConfirmedQty(ctx context.Context, tx *gorm.DB, detailID int64, delta int32) error
	// IsOrderFullyConfirmed 判断入库单是否所有明细均已全部确认。
	IsOrderFullyConfirmed(ctx context.Context, tx *gorm.DB, orderID int64) (bool, error)

	// GetPendingTraceInventory 查询某入库单下状态为 PENDING 的追溯库存记录（取消前校验）。
	GetPendingTraceInventory(ctx context.Context, tx *gorm.DB, orderID int64) ([]*model.DrugTraceInventory, error)
	// SoftDeletePendingTraceInventory 软删除某入库单下所有 PENDING 状态的追溯库存记录。
	SoftDeletePendingTraceInventory(ctx context.Context, tx *gorm.DB, orderID int64) error

	// CreateTraceInventory 写入追溯库存记录。
	CreateTraceInventory(ctx context.Context, tx *gorm.DB, records []*model.DrugTraceInventory) error
	// ExistsTraceCode 判断追溯码是否已存在于 drug_trace_inventory。
	ExistsTraceCode(ctx context.Context, tx *gorm.DB, traceCode string) (bool, error)

	// CreateTraceLog 批量写入追溯日志。
	CreateTraceLog(ctx context.Context, tx *gorm.DB, logs []*model.DrugTraceLog) error

	// GenerateOrderNo 生成当日唯一入库单号，格式 IN-YYYYMMDD-XXXX。
	// 调用方须在事务中调用，以保证并发安全。
	GenerateOrderNo(ctx context.Context, tx *gorm.DB) (string, error)
}

// inboundRepo 是 InboundRepo 的 GORM 实现。
type inboundRepo struct {
	db *gorm.DB
}

// NewInboundRepo 创建仓储实例。
func NewInboundRepo(db *gorm.DB) InboundRepo {
	return &inboundRepo{db: db}
}

// ============================================================
// 入库单主表操作
// ============================================================

// GetByID 按主键查询入库单（不含明细）。
func (r *inboundRepo) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	var order model.InboundOrder
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByIDForUpdate 在事务内加行锁查询入库单。
func (r *inboundRepo) GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.InboundOrder, error) {
	var order model.InboundOrder
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 按单号查询入库单。
func (r *inboundRepo) GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error) {
	var order model.InboundOrder
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// List 分页查询入库单列表。
func (r *inboundRepo) List(ctx context.Context, filter ListFilter) ([]*model.InboundOrder, int64, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.InboundOrder{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.SupplierID > 0 {
		query = query.Where("supplier_id = ?", filter.SupplierID)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("order_no LIKE ? OR invoice_no LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []*model.InboundOrder
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// Create 新建入库单（在事务内执行）。
func (r *inboundRepo) Create(ctx context.Context, tx *gorm.DB, order *model.InboundOrder) error {
	return tx.WithContext(ctx).Create(order).Error
}

// Update 更新入库单基础字段。
func (r *inboundRepo) Update(ctx context.Context, tx *gorm.DB, order *model.InboundOrder) error {
	return tx.WithContext(ctx).Model(order).
		Select("supplier_id", "invoice_no", "remark", "operator_id").
		Updates(order).Error
}

// UpdateStatus 更新入库单状态及对应时间戳。
func (r *inboundRepo) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string) error {
	now := time.Now()
	updates := map[string]interface{}{"status": status}
	switch status {
	case model.InboundOrderStatusPendingConfirm:
		updates["submitted_at"] = now
	case model.InboundOrderStatusCompleted:
		updates["completed_at"] = now
	case model.InboundOrderStatusCancelled:
		updates["cancelled_at"] = now
	}
	return tx.WithContext(ctx).
		Model(&model.InboundOrder{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateTotalAmount 重算并写入 total_amount。
func (r *inboundRepo) UpdateTotalAmount(ctx context.Context, tx *gorm.DB, orderID int64) error {
	// 使用子查询对所有未软删除明细的 amount 求和。
	sql := `UPDATE inbound_order SET total_amount = (
		SELECT COALESCE(SUM(amount), 0) FROM inbound_order_detail
		WHERE order_id = ? AND deleted_at IS NULL
	) WHERE id = ?`
	return tx.WithContext(ctx).Exec(sql, orderID, orderID).Error
}

// ============================================================
// 入库单明细操作
// ============================================================

// GetDetails 查询某入库单所有未删除明细。
func (r *inboundRepo) GetDetails(ctx context.Context, orderID int64) ([]*model.InboundOrderDetail, error) {
	var details []*model.InboundOrderDetail
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("id ASC").
		Find(&details).Error
	return details, err
}

// GetDetailByID 按明细 ID 查询（同时校验所属入库单）。
func (r *inboundRepo) GetDetailByID(ctx context.Context, orderID, detailID int64) (*model.InboundOrderDetail, error) {
	var detail model.InboundOrderDetail
	err := r.db.WithContext(ctx).
		Where("id = ? AND order_id = ?", detailID, orderID).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &detail, nil
}

// GetDetailByIDForUpdate 在事务内加行锁查询明细。
func (r *inboundRepo) GetDetailByIDForUpdate(ctx context.Context, tx *gorm.DB, orderID, detailID int64) (*model.InboundOrderDetail, error) {
	var detail model.InboundOrderDetail
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND order_id = ?", detailID, orderID).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &detail, nil
}

// AddDetail 新增明细行。
func (r *inboundRepo) AddDetail(ctx context.Context, tx *gorm.DB, detail *model.InboundOrderDetail) error {
	return tx.WithContext(ctx).Create(detail).Error
}

// UpdateDetail 更新明细可编辑字段。
func (r *inboundRepo) UpdateDetail(ctx context.Context, tx *gorm.DB, detail *model.InboundOrderDetail) error {
	return tx.WithContext(ctx).Model(detail).
		Select("drug_id", "batch_number", "expire_date", "planned_qty", "unit_price", "amount", "remark").
		Updates(detail).Error
}

// DeleteDetail 软删除明细行（GORM 软删除），同时重新计算入库单总金额。
func (r *inboundRepo) DeleteDetail(ctx context.Context, tx *gorm.DB, orderID, detailID int64) error {
	return tx.WithContext(ctx).
		Where("id = ? AND order_id = ?", detailID, orderID).
		Delete(&model.InboundOrderDetail{}).Error
}

// IncreaseConfirmedQty 原子累加明细已确认数量。
func (r *inboundRepo) IncreaseConfirmedQty(ctx context.Context, tx *gorm.DB, detailID int64, delta int32) error {
	return tx.WithContext(ctx).
		Model(&model.InboundOrderDetail{}).
		Where("id = ?", detailID).
		Update("confirmed_qty", gorm.Expr("confirmed_qty + ?", delta)).Error
}

// IsOrderFullyConfirmed 判断入库单所有明细是否均已全量确认。
func (r *inboundRepo) IsOrderFullyConfirmed(ctx context.Context, tx *gorm.DB, orderID int64) (bool, error) {
	type agg struct {
		Planned   int64
		Confirmed int64
	}
	var summary agg
	err := tx.WithContext(ctx).
		Model(&model.InboundOrderDetail{}).
		Select("COALESCE(SUM(planned_qty), 0) AS planned, COALESCE(SUM(confirmed_qty), 0) AS confirmed").
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Scan(&summary).Error
	if err != nil {
		return false, err
	}
	// 没有有效明细时不应触发完成逻辑。
	if summary.Planned == 0 {
		return false, nil
	}
	return summary.Confirmed >= summary.Planned, nil
}

// ============================================================
// 追溯库存操作
// ============================================================

// GetPendingTraceInventory 查询某入库单下 PENDING 状态的追溯库存。
func (r *inboundRepo) GetPendingTraceInventory(ctx context.Context, tx *gorm.DB, orderID int64) ([]*model.DrugTraceInventory, error) {
	var records []*model.DrugTraceInventory
	err := tx.WithContext(ctx).
		Where("inbound_order_id = ? AND status = ?", orderID, model.TraceInventoryStatusPending).
		Find(&records).Error
	return records, err
}

// SoftDeletePendingTraceInventory 软删除某入库单下所有 PENDING 状态的追溯库存记录。
func (r *inboundRepo) SoftDeletePendingTraceInventory(ctx context.Context, tx *gorm.DB, orderID int64) error {
	return tx.WithContext(ctx).
		Where("inbound_order_id = ? AND status = ?", orderID, model.TraceInventoryStatusPending).
		Delete(&model.DrugTraceInventory{}).Error
}

// CreateTraceInventory 批量写入追溯库存记录。
func (r *inboundRepo) CreateTraceInventory(ctx context.Context, tx *gorm.DB, records []*model.DrugTraceInventory) error {
	if len(records) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&records).Error
}

// ExistsTraceCode 判断追溯码是否已在 drug_trace_inventory 表中存在。
func (r *inboundRepo) ExistsTraceCode(ctx context.Context, tx *gorm.DB, traceCode string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).
		Model(&model.DrugTraceInventory{}).
		Where("trace_code = ?", traceCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============================================================
// 追溯日志操作
// ============================================================

// CreateTraceLog 批量写入追溯日志。
func (r *inboundRepo) CreateTraceLog(ctx context.Context, tx *gorm.DB, logs []*model.DrugTraceLog) error {
	if len(logs) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&logs).Error
}

// ============================================================
// 单号生成
// ============================================================

// GenerateOrderNo 生成当日唯一入库单号，格式：IN-YYYYMMDD-XXXX。
// 利用当天已有单号数量 + 1 作为序号，调用方须在事务内执行以避免并发重复。
func (r *inboundRepo) GenerateOrderNo(ctx context.Context, tx *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "IN-" + today + "-"

	var count int64
	err := tx.WithContext(ctx).
		Model(&model.InboundOrder{}).
		Unscoped(). // 包含软删除记录，避免序号跳空后重用
		Where("order_no LIKE ?", prefix+"%").
		Count(&count).Error
	if err != nil {
		return "", fmt.Errorf("生成入库单号失败: %w", err)
	}
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}
