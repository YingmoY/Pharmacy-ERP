package repository

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/inventory/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InboundOrder 定义入库单读写能力。
type InboundOrder interface {
	GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.InboundOrder, error)
	GetDetailByIDForUpdate(ctx context.Context, tx *gorm.DB, orderID, detailID int64) (*model.InboundOrderDetail, error)
	IncreaseConfirmedQty(ctx context.Context, tx *gorm.DB, detailID int64, delta int32) error
	IsOrderFullyConfirmed(ctx context.Context, tx *gorm.DB, orderID int64) (bool, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string) error
}

// inboundOrder 是 InboundOrder 的 GORM 实现。
type inboundOrder struct{}

func NewInboundOrder() InboundOrder {
	return &inboundOrder{}
}

func (r *inboundOrder) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string) error {
	return tx.WithContext(ctx).
		Model(&model.InboundOrder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *inboundOrder) GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.InboundOrder, error) {
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

func (r *inboundOrder) GetDetailByIDForUpdate(ctx context.Context, tx *gorm.DB, orderID, detailID int64) (*model.InboundOrderDetail, error) {
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

func (r *inboundOrder) IncreaseConfirmedQty(ctx context.Context, tx *gorm.DB, detailID int64, delta int32) error {
	return tx.WithContext(ctx).
		Model(&model.InboundOrderDetail{}).
		Where("id = ?", detailID).
		Update("confirmed_qty", gorm.Expr("confirmed_qty + ?", delta)).Error
}

func (r *inboundOrder) IsOrderFullyConfirmed(ctx context.Context, tx *gorm.DB, orderID int64) (bool, error) {
	type agg struct {
		Planned   int64
		Confirmed int64
	}

	var summary agg
	err := tx.WithContext(ctx).
		Model(&model.InboundOrderDetail{}).
		Select("COALESCE(SUM(planned_qty), 0) AS planned, COALESCE(SUM(confirmed_qty), 0) AS confirmed").
		Where("order_id = ?", orderID).
		Scan(&summary).Error
	if err != nil {
		return false, err
	}

	// 没有明细时不应进入完成态，避免误完成空单据。
	if summary.Planned == 0 {
		return false, nil
	}

	return summary.Confirmed >= summary.Planned, nil
}
