package repository

import (
	"context"

	"github.com/YingmoY/PharmacyERP/internal/inventory/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TraceInventory 定义追溯库存表读写能力。
type TraceInventory interface {
	// BatchCreate 批量写入追溯库存记录（入库确认场景）。
	BatchCreate(ctx context.Context, tx *gorm.DB, records []model.TraceInventory) error
	// FindExistingTraceCodes 查询指定追溯码中已存在于库存表的数据（用于防重）。
	FindExistingTraceCodes(ctx context.Context, tx *gorm.DB, traceCodes []string) ([]string, error)
	// GetByTraceCodesForUpdate 对指定追溯码加行锁读取，保证上架事务并发安全。
	GetByTraceCodesForUpdate(ctx context.Context, tx *gorm.DB, traceCodes []string) ([]model.TraceInventory, error)
	// BatchPutaway 将一批追溯码从待上架更新为在库，并写入货位。
	BatchPutaway(ctx context.Context, tx *gorm.DB, traceCodes []string, locationID int64) error
}

// traceInventory 是 TraceInventory 的 GORM 实现。
type traceInventory struct{}

func NewTraceInventory() TraceInventory {
	return &traceInventory{}
}

func (r *traceInventory) BatchCreate(ctx context.Context, tx *gorm.DB, records []model.TraceInventory) error {
	if len(records) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&records).Error
}

func (r *traceInventory) FindExistingTraceCodes(ctx context.Context, tx *gorm.DB, traceCodes []string) ([]string, error) {
	if len(traceCodes) == 0 {
		return nil, nil
	}

	var existing []string
	err := tx.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("trace_code IN ?", traceCodes).
		Pluck("trace_code", &existing).Error
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *traceInventory) GetByTraceCodesForUpdate(ctx context.Context, tx *gorm.DB, traceCodes []string) ([]model.TraceInventory, error) {
	if len(traceCodes) == 0 {
		return nil, nil
	}

	var records []model.TraceInventory
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("trace_code IN ?", traceCodes).
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *traceInventory) BatchPutaway(ctx context.Context, tx *gorm.DB, traceCodes []string, locationID int64) error {
	if len(traceCodes) == 0 {
		return nil
	}

	return tx.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("trace_code IN ?", traceCodes).
		Updates(map[string]interface{}{
			"location_id": locationID,
			"status":      model.TraceInventoryStatusInStock,
		}).Error
}
