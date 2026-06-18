package service

import (
	"context"
	"errors"
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/inventory/model"
	"github.com/YingmoY/PharmacyERP/internal/inventory/repository"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InboundService 定义入库核心业务能力。
type InboundService interface {
	// ConfirmInboundTraceCodes 执行入库确认事务：
	// 1) 校验入库单状态和入库明细数量
	// 2) 写入追溯库存（待上架）
	// 3) 回写明细已确认数量与入库单状态
	ConfirmInboundTraceCodes(ctx context.Context, req ConfirmInboundTraceCodesRequest) error

	// PutawayTraceCodes 执行上架事务：
	// 1) 校验追溯码集合参数
	// 2) 行锁读取对应追溯库存记录并校验状态
	// 3) 批量更新为在库状态并绑定货位
	PutawayTraceCodes(ctx context.Context, req PutawayTraceCodesRequest) error

	// IsValidOrderStatusTransition 校验入库单状态流转是否合法。
	IsValidOrderStatusTransition(fromStatus, toStatus string) bool

	// IsValidTraceStatusTransition 校验追溯库存状态流转是否合法。
	IsValidTraceStatusTransition(fromStatus, toStatus string) bool
}

// ConfirmInboundTraceCodesRequest 是入库确认请求参数。
type ConfirmInboundTraceCodesRequest struct {
	OrderID    int64    // 入库单 ID
	DetailID   int64    // 入库明细 ID
	TraceCodes []string // 本次确认的追溯码列表
	OperatorID int64    // 操作人 ID
}

// PutawayTraceCodesRequest 是上架请求参数。
type PutawayTraceCodesRequest struct {
	LocationID int64    // 目标货位 ID
	TraceCodes []string // 待上架追溯码列表
}

// inboundService 是 InboundService 的默认实现。
type inboundService struct {
	db        *gorm.DB
	orderRepo repository.InboundOrder
	traceRepo repository.TraceInventory
	logger    *zap.Logger
}

// NewInboundService 创建入库服务。
func NewInboundService(
	db *gorm.DB,
	orderRepo repository.InboundOrder,
	traceRepo repository.TraceInventory,
	logger *zap.Logger,
) InboundService {
	return &inboundService{
		db:        db,
		orderRepo: orderRepo,
		traceRepo: traceRepo,
		logger:    logger,
	}
}

func (s *inboundService) ConfirmInboundTraceCodes(ctx context.Context, req ConfirmInboundTraceCodesRequest) error {
	if req.OrderID <= 0 || req.DetailID <= 0 || req.OperatorID <= 0 {
		return ecode.ErrParamInvalid
	}
	if len(req.TraceCodes) == 0 {
		return ecode.ErrParamInvalid
	}

	cleanCodes, err := normalizeAndValidateTraceCodes(req.TraceCodes)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) 加锁读取入库单，保证并发下状态和数量判断一致。
		order, err := s.orderRepo.GetByIDForUpdate(ctx, tx, req.OrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrParamInvalid
			}
			return err
		}

		// 文档要求：入库确认必须在待确认状态下执行。
		if !strings.EqualFold(order.Status, model.InboundOrderStatusPendingConfirm) {
			return ecode.ErrStatusInvalid
		}

		// 2) 加锁读取入库明细，校验数量边界。
		detail, err := s.orderRepo.GetDetailByIDForUpdate(ctx, tx, req.OrderID, req.DetailID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrParamInvalid
			}
			return err
		}

		if detail.ConfirmedQty+int32(len(cleanCodes)) > detail.PlannedQty {
			return ecode.ErrParamInvalid
		}

		// 3) 检查追溯码是否已存在于库存表，防止重复入库。
		existing, err := s.traceRepo.FindExistingTraceCodes(ctx, tx, cleanCodes)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return ecode.ErrDuplicateScan
		}

		// 4) 批量写入追溯库存（初始状态：待上架）。
		inventories := make([]model.TraceInventory, 0, len(cleanCodes))
		for _, code := range cleanCodes {
			inventories = append(inventories, model.TraceInventory{
				TraceCode:       code,
				DrugID:          detail.DrugID,
				BatchNumber:     detail.BatchNumber,
				ExpireDate:      detail.ExpireDate,
				LocationID:      nil,
				InboundOrderID:  order.ID,
				InboundDetailID: detail.ID,
				Status:          model.TraceInventoryStatusPending,
			})
		}
		if err := s.traceRepo.BatchCreate(ctx, tx, inventories); err != nil {
			return err
		}

		// 5) 更新明细已确认数量。
		if err := s.orderRepo.IncreaseConfirmedQty(ctx, tx, detail.ID, int32(len(cleanCodes))); err != nil {
			return err
		}

		// 6) 仅当整单全部确认完成时，才推进到“已确认”状态。
		fullyConfirmed, err := s.orderRepo.IsOrderFullyConfirmed(ctx, tx, order.ID)
		if err != nil {
			return err
		}
		if fullyConfirmed {
			if err := s.orderRepo.UpdateStatus(ctx, tx, order.ID, model.InboundOrderStatusCompleted); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *inboundService) PutawayTraceCodes(ctx context.Context, req PutawayTraceCodesRequest) error {
	if req.LocationID <= 0 || len(req.TraceCodes) == 0 {
		return ecode.ErrParamInvalid
	}

	cleanCodes, err := normalizeAndValidateTraceCodes(req.TraceCodes)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 0) 校验货位是否存在且可用。
		var loc locationModel.LocationInfo
		err := tx.WithContext(ctx).
			Where("id = ? AND status = ?", req.LocationID, locationModel.LocationStatusEnabled).
			First(&loc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrLocationNotFound
			}
			return err
		}

		// 1) 行锁读取所有目标追溯码，保证并发上架安全。
		records, err := s.traceRepo.GetByTraceCodesForUpdate(ctx, tx, cleanCodes)
		if err != nil {
			return err
		}

		// 2) 数量必须一一对应，否则说明存在无效追溯码。
		if len(records) != len(cleanCodes) {
			return ecode.ErrTraceCodeNotFound
		}

		// 3) 仅允许 PENDING -> IN_STOCK 流转。
		for _, record := range records {
			if !s.IsValidTraceStatusTransition(record.Status, model.TraceInventoryStatusInStock) {
				return ecode.ErrTraceStatusLocked
			}
		}

		// 4) 批量更新货位与状态。
		if err := s.traceRepo.BatchPutaway(ctx, tx, cleanCodes, req.LocationID); err != nil {
			return err
		}

		return nil
	})
}

func (s *inboundService) IsValidOrderStatusTransition(fromStatus, toStatus string) bool {
	from := strings.ToUpper(strings.TrimSpace(fromStatus))
	to := strings.ToUpper(strings.TrimSpace(toStatus))

	allowed := map[string]map[string]struct{}{
		strings.ToUpper(model.InboundOrderStatusDraft): {
			strings.ToUpper(model.InboundOrderStatusPendingConfirm): {},
			strings.ToUpper(model.InboundOrderStatusCancelled):      {},
		},
		strings.ToUpper(model.InboundOrderStatusPendingConfirm): {
			strings.ToUpper(model.InboundOrderStatusCompleted): {},
			strings.ToUpper(model.InboundOrderStatusCancelled): {},
		},
	}

	next, ok := allowed[from]
	if !ok {
		return false
	}
	_, ok = next[to]
	return ok
}

func (s *inboundService) IsValidTraceStatusTransition(fromStatus, toStatus string) bool {
	from := strings.ToUpper(strings.TrimSpace(fromStatus))
	to := strings.ToUpper(strings.TrimSpace(toStatus))

	allowed := map[string]map[string]struct{}{
		model.TraceInventoryStatusPending: {
			model.TraceInventoryStatusInStock: {},
		},
		model.TraceInventoryStatusInStock: {
			model.TraceInventoryStatusSold:          {},
			model.TraceInventoryStatusMisplaced:     {},
			model.TraceInventoryStatusLossCandidate: {},
		},
		model.TraceInventoryStatusMisplaced: {
			model.TraceInventoryStatusInStock: {},
		},
		model.TraceInventoryStatusLossCandidate: {
			model.TraceInventoryStatusMisplaced: {},
			model.TraceInventoryStatusLost:      {},
		},
	}

	next, ok := allowed[from]
	if !ok {
		return false
	}
	_, ok = next[to]
	return ok
}

// normalizeAndValidateTraceCodes 负责清洗输入并校验任务内重复。
func normalizeAndValidateTraceCodes(codes []string) ([]string, error) {
	result := make([]string, 0, len(codes))
	seen := make(map[string]struct{}, len(codes))

	for _, raw := range codes {
		code := strings.TrimSpace(raw)
		if code == "" {
			return nil, ecode.ErrParamInvalid
		}
		if _, ok := seen[code]; ok {
			return nil, ecode.ErrDuplicateScan
		}
		seen[code] = struct{}{}
		result = append(result, code)
	}

	return result, nil
}
