package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/sales/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SalesRepo 定义销售模块所有数据库操作能力
type SalesRepo interface {
	// ==================== 销售订单 ====================

	// CreateOrder 创建销售订单（事务内使用）
	CreateOrder(ctx context.Context, tx *gorm.DB, order *model.SalesOrder) error
	// GetOrderByID 根据 ID 查询订单（不加锁）
	GetOrderByID(ctx context.Context, db *gorm.DB, id int64) (*model.SalesOrder, error)
	// GetOrderByIDForUpdate 根据 ID 加行锁查询订单（事务内使用）
	GetOrderByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.SalesOrder, error)
	// UpdateOrderStatus 更新订单状态
	UpdateOrderStatus(ctx context.Context, tx *gorm.DB, id int64, status string, extraFields map[string]interface{}) error
	// UpdateOrderAmounts 更新订单金额相关字段
	UpdateOrderAmounts(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error
	// ListOrders 分页查询订单列表
	ListOrders(ctx context.Context, db *gorm.DB, filter OrderFilter) ([]*model.SalesOrder, int64, error)
	// GenOrderNo 生成唯一订单号（SO-YYYYMMDD-XXXX）
	GenOrderNo(ctx context.Context, db *gorm.DB) (string, error)

	// ==================== 销售订单明细 ====================

	// CreateItem 创建订单明细（事务内使用）
	CreateItem(ctx context.Context, tx *gorm.DB, item *model.SalesOrderItem) error
	// GetItemByID 根据 ID 查询明细
	GetItemByID(ctx context.Context, db *gorm.DB, id int64) (*model.SalesOrderItem, error)
	// ListItemsByOrderID 查询订单下所有未删除明细
	ListItemsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.SalesOrderItem, error)
	// DeleteItem 软删除订单明细
	DeleteItem(ctx context.Context, tx *gorm.DB, id int64) error
	// UpdateItem 更新明细字段
	UpdateItem(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error
	// SumItemTotalByOrderID 汇总订单明细总金额
	SumItemTotalByOrderID(ctx context.Context, tx *gorm.DB, orderID int64) (float64, error)

	// ==================== 追溯预留 ====================

	// CreateReservation 创建追溯码预留（事务内使用）
	CreateReservation(ctx context.Context, tx *gorm.DB, rsv *model.TraceReservation) error
	// GetActiveReservationByTraceCode 查询追溯码当前有效预留（status=RESERVED）
	GetActiveReservationByTraceCode(ctx context.Context, db *gorm.DB, traceCode string) (*model.TraceReservation, error)
	// GetReservationsByOrderID 查询订单下所有预留记录
	GetReservationsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.TraceReservation, error)
	// GetReservationByOrderItemID 查询订单明细对应的预留记录
	GetReservationByOrderItemID(ctx context.Context, db *gorm.DB, orderID, itemID int64) (*model.TraceReservation, error)
	// UpdateReservationStatus 更新预留状态
	UpdateReservationStatus(ctx context.Context, tx *gorm.DB, id int64, status string, extraFields map[string]interface{}) error
	// GenReservationNo 生成唯一预留单号（RSV-YYYYMMDD-XXXX）
	GenReservationNo(ctx context.Context, db *gorm.DB) (string, error)
	// GetActiveReservationsByOrderID 查询订单下所有状态为 RESERVED 的预留
	GetActiveReservationsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.TraceReservation, error)
}

// OrderFilter 订单列表查询过滤条件
type OrderFilter struct {
	CashierID      *int64
	Status         string
	OrderNo        string
	StartDate      string
	EndDate        string
	IsPrescription *bool
	Page           int
	PageSize       int
}

// salesRepo 是 SalesRepo 的 GORM 实现
type salesRepo struct{}

// NewSalesRepo 创建销售 Repository 实例
func NewSalesRepo() SalesRepo {
	return &salesRepo{}
}

// ==================== 销售订单实现 ====================

func (r *salesRepo) CreateOrder(ctx context.Context, tx *gorm.DB, order *model.SalesOrder) error {
	return tx.WithContext(ctx).Create(order).Error
}

func (r *salesRepo) GetOrderByID(ctx context.Context, db *gorm.DB, id int64) (*model.SalesOrder, error) {
	var order model.SalesOrder
	err := db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *salesRepo) GetOrderByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.SalesOrder, error) {
	var order model.SalesOrder
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *salesRepo) UpdateOrderStatus(ctx context.Context, tx *gorm.DB, id int64, status string, extraFields map[string]interface{}) error {
	fields := map[string]interface{}{
		"status": status,
	}
	for k, v := range extraFields {
		fields[k] = v
	}
	return tx.WithContext(ctx).
		Model(&model.SalesOrder{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *salesRepo) UpdateOrderAmounts(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error {
	return tx.WithContext(ctx).
		Model(&model.SalesOrder{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *salesRepo) ListOrders(ctx context.Context, db *gorm.DB, filter OrderFilter) ([]*model.SalesOrder, int64, error) {
	query := db.WithContext(ctx).Model(&model.SalesOrder{})
	if filter.CashierID != nil && *filter.CashierID > 0 {
		query = query.Where("cashier_id = ?", *filter.CashierID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OrderNo != "" {
		query = query.Where("order_no = ?", filter.OrderNo)
	}
	if filter.StartDate != "" {
		query = query.Where("DATE(created_at) >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("DATE(created_at) <= ?", filter.EndDate)
	}
	if filter.IsPrescription != nil {
		query = query.Where("is_prescription = ?", *filter.IsPrescription)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
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
	offset := (page - 1) * pageSize

	var orders []*model.SalesOrder
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GenOrderNo 生成唯一订单号，格式 SO-YYYYMMDD-XXXX
// 使用序列自增方式：当天最大序号 +1
func (r *salesRepo) GenOrderNo(ctx context.Context, db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "SO-" + today + "-"

	var maxNo string
	err := db.WithContext(ctx).
		Model(&model.SalesOrder{}).
		Where("order_no LIKE ?", prefix+"%").
		Order("order_no DESC").
		Limit(1).
		Pluck("order_no", &maxNo).Error
	if err != nil {
		return "", err
	}

	seq := 1
	if maxNo != "" && len(maxNo) > len(prefix) {
		suffix := maxNo[len(prefix):]
		var n int
		if _, err := fmt.Sscanf(suffix, "%d", &n); err == nil {
			seq = n + 1
		}
	}
	return fmt.Sprintf("%s%04d", prefix, seq), nil
}

// ==================== 销售订单明细实现 ====================

func (r *salesRepo) CreateItem(ctx context.Context, tx *gorm.DB, item *model.SalesOrderItem) error {
	return tx.WithContext(ctx).Create(item).Error
}

func (r *salesRepo) GetItemByID(ctx context.Context, db *gorm.DB, id int64) (*model.SalesOrderItem, error) {
	var item model.SalesOrderItem
	err := db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *salesRepo) ListItemsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.SalesOrderItem, error) {
	var items []*model.SalesOrderItem
	err := db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *salesRepo) DeleteItem(ctx context.Context, tx *gorm.DB, id int64) error {
	return tx.WithContext(ctx).Delete(&model.SalesOrderItem{}, id).Error
}

func (r *salesRepo) UpdateItem(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error {
	return tx.WithContext(ctx).
		Model(&model.SalesOrderItem{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *salesRepo) SumItemTotalByOrderID(ctx context.Context, tx *gorm.DB, orderID int64) (float64, error) {
	var total float64
	err := tx.WithContext(ctx).
		Model(&model.SalesOrderItem{}).
		Where("order_id = ?", orderID).
		Select("COALESCE(SUM(subtotal_amount), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ==================== 追溯预留实现 ====================

func (r *salesRepo) CreateReservation(ctx context.Context, tx *gorm.DB, rsv *model.TraceReservation) error {
	return tx.WithContext(ctx).Create(rsv).Error
}

func (r *salesRepo) GetActiveReservationByTraceCode(ctx context.Context, db *gorm.DB, traceCode string) (*model.TraceReservation, error) {
	var rsv model.TraceReservation
	err := db.WithContext(ctx).
		Where("trace_code = ? AND status = ?", traceCode, model.ReservationStatusReserved).
		First(&rsv).Error
	if err != nil {
		return nil, err
	}
	return &rsv, nil
}

func (r *salesRepo) GetReservationsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.TraceReservation, error) {
	var rsvs []*model.TraceReservation
	err := db.WithContext(ctx).
		Where("sales_order_id = ?", orderID).
		Order("id ASC").
		Find(&rsvs).Error
	if err != nil {
		return nil, err
	}
	return rsvs, nil
}

func (r *salesRepo) GetReservationByOrderItemID(ctx context.Context, db *gorm.DB, orderID, itemID int64) (*model.TraceReservation, error) {
	var rsv model.TraceReservation
	err := db.WithContext(ctx).
		Where("sales_order_id = ? AND sales_order_item_id = ?", orderID, itemID).
		First(&rsv).Error
	if err != nil {
		return nil, err
	}
	return &rsv, nil
}

func (r *salesRepo) UpdateReservationStatus(ctx context.Context, tx *gorm.DB, id int64, status string, extraFields map[string]interface{}) error {
	fields := map[string]interface{}{
		"status": status,
	}
	for k, v := range extraFields {
		fields[k] = v
	}
	return tx.WithContext(ctx).
		Model(&model.TraceReservation{}).
		Where("id = ?", id).
		Updates(fields).Error
}

// GenReservationNo 生成唯一预留单号，格式 RSV-YYYYMMDD-XXXX
func (r *salesRepo) GenReservationNo(ctx context.Context, db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "RSV-" + today + "-"

	var maxNo string
	err := db.WithContext(ctx).
		Model(&model.TraceReservation{}).
		Where("reservation_no LIKE ?", prefix+"%").
		Order("reservation_no DESC").
		Limit(1).
		Pluck("reservation_no", &maxNo).Error
	if err != nil {
		return "", err
	}

	seq := 1
	if maxNo != "" && len(maxNo) > len(prefix) {
		suffix := maxNo[len(prefix):]
		var n int
		if _, err := fmt.Sscanf(suffix, "%d", &n); err == nil {
			seq = n + 1
		}
	}
	return fmt.Sprintf("%s%04d", prefix, seq), nil
}

func (r *salesRepo) GetActiveReservationsByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]*model.TraceReservation, error) {
	var rsvs []*model.TraceReservation
	err := db.WithContext(ctx).
		Where("sales_order_id = ? AND status = ?", orderID, model.ReservationStatusReserved).
		Order("id ASC").
		Find(&rsvs).Error
	if err != nil {
		return nil, err
	}
	return rsvs, nil
}
