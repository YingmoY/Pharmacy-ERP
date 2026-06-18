package service

import (
	"context"
	"errors"
	"strings"
	"time"

	inventoryModel "github.com/YingmoY/PharmacyERP/internal/inventory/model"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ShelveItem 上架请求中的单条目。
type ShelveItem struct {
	TraceCode    string
	LocationCode string
}

// ShelveResult 上架操作的单条结果。
type ShelveResult struct {
	TraceCode string `json:"trace_code"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// DrugMixInfo 混合陈列药品信息。
type DrugMixInfo struct {
	DrugID         int64     `json:"drug_id"`
	DrugName       string    `json:"drug_name"`
	Count          int64     `json:"count"`
	BatchNumbers   []string  `json:"batch_numbers"`
	EarliestExpire time.Time `json:"earliest_expire"`
}

// MixCheckResult 货位混合陈列检查结果。
type MixCheckResult struct {
	LocationCode  string        `json:"location_code"`
	LocationName  string        `json:"location_name"`
	HasMixedDrugs bool          `json:"has_mixed_drugs"`
	Drugs         []DrugMixInfo `json:"drugs"`
}

// ShelvingService 定义上架业务能力接口。
type ShelvingService interface {
	// ShelveTrace 将单个追溯码上架到指定货位（PENDING -> IN_STOCK）。
	ShelveTrace(ctx context.Context, traceCode, locationCode string, operatorID int64) error

	// BatchShelve 批量上架，每条独立处理，允许部分成功。
	BatchShelve(ctx context.Context, items []ShelveItem, operatorID int64) []ShelveResult

	// GetPendingList 分页查询待上架追溯码列表。
	GetPendingList(ctx context.Context, page, pageSize int) (*core.PageResult, error)

	// Relocate 将在库或错架追溯码移动到新货位，写调整记录与追溯日志。
	Relocate(ctx context.Context, traceCode, locationCode string, operatorID int64) error

	// MixCheck 检查指定货位上当前陈列的药品是否混放。
	MixCheck(ctx context.Context, locationCode string) (*MixCheckResult, error)
}

// shelvingService 是 ShelvingService 的默认实现。
type shelvingService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewShelvingService 创建上架服务。
func NewShelvingService(db *gorm.DB, logger *zap.Logger) ShelvingService {
	return &shelvingService{db: db, logger: logger}
}

// findActiveLocation 根据货位编码查询激活状态的货位记录。
func (s *shelvingService) findActiveLocation(ctx context.Context, db *gorm.DB, locationCode string) (*locationModel.LocationInfo, error) {
	var loc locationModel.LocationInfo
	err := db.WithContext(ctx).
		Where("location_code = ? AND status = ?", strings.TrimSpace(locationCode), locationModel.LocationStatusEnabled).
		First(&loc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrLocationNotFound
		}
		return nil, err
	}
	return &loc, nil
}

func (s *shelvingService) ShelveTrace(ctx context.Context, traceCode, locationCode string, operatorID int64) error {
	traceCode = strings.TrimSpace(traceCode)
	locationCode = strings.TrimSpace(locationCode)
	if traceCode == "" || locationCode == "" || operatorID <= 0 {
		return ecode.ErrParamInvalid
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) 校验货位是否存在且激活。
		loc, err := s.findActiveLocation(ctx, tx, locationCode)
		if err != nil {
			return err
		}

		// 2) 加行锁读取追溯码当前状态。
		var trace inventoryModel.TraceInventory
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("trace_code = ?", traceCode).
			First(&trace).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}

		// 3) 追溯码必须处于 PENDING 状态。
		if trace.Status != inventoryModel.TraceInventoryStatusPending {
			return ecode.ErrTraceStatusLocked
		}

		// 4) 校验关联入库单已完成。
		var orderStatus string
		err = tx.Table("inbound_order").
			Select("status").
			Where("id = ? AND deleted_at IS NULL", trace.InboundOrderID).
			Pluck("status", &orderStatus).Error
		if err != nil {
			return err
		}
		if orderStatus != inventoryModel.InboundOrderStatusCompleted {
			return ecode.New(20006, "inbound order not completed yet")
		}

		// 5) 更新状态：PENDING -> IN_STOCK，绑定货位。
		if err := tx.Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", traceCode).
			Updates(map[string]interface{}{
				"status":      inventoryModel.TraceInventoryStatusInStock,
				"location_id": loc.ID,
			}).Error; err != nil {
			return err
		}

		// 6) 写入追溯日志。
		fromStatus := inventoryModel.TraceInventoryStatusPending
		toStatus := inventoryModel.TraceInventoryStatusInStock
		actionType := "SHELVING"
		log := &inventoryModel.DrugTraceLog{
			TraceCode:    traceCode,
			ActionType:   actionType,
			FromStatus:   &fromStatus,
			ToStatus:     &toStatus,
			ToLocationID: &loc.ID,
			OperatorID:   operatorID,
		}
		return tx.Create(log).Error
	})
}

func (s *shelvingService) BatchShelve(ctx context.Context, items []ShelveItem, operatorID int64) []ShelveResult {
	results := make([]ShelveResult, 0, len(items))
	for _, item := range items {
		err := s.ShelveTrace(ctx, item.TraceCode, item.LocationCode, operatorID)
		if err != nil {
			bizErr := ecode.FromError(err)
			results = append(results, ShelveResult{
				TraceCode: item.TraceCode,
				Success:   false,
				Error:     bizErr.Msg,
			})
		} else {
			results = append(results, ShelveResult{
				TraceCode: item.TraceCode,
				Success:   true,
			})
		}
	}
	return results
}

func (s *shelvingService) GetPendingList(ctx context.Context, page, pageSize int) (*core.PageResult, error) {
	query := s.db.WithContext(ctx).
		Model(&inventoryModel.TraceInventory{}).
		Joins("JOIN inbound_order ON inbound_order.id = drug_trace_inventory.inbound_order_id AND inbound_order.deleted_at IS NULL").
		Where("drug_trace_inventory.status = ?", inventoryModel.TraceInventoryStatusPending).
		Where("inbound_order.status = ?", inventoryModel.InboundOrderStatusCompleted)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []inventoryModel.TraceInventory
	if err := query.
		Order("drug_trace_inventory.created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := core.NewPageResult(total, page, pageSize, list)
	return &result, nil
}

func (s *shelvingService) Relocate(ctx context.Context, traceCode, locationCode string, operatorID int64) error {
	traceCode = strings.TrimSpace(traceCode)
	locationCode = strings.TrimSpace(locationCode)
	if traceCode == "" || locationCode == "" || operatorID <= 0 {
		return ecode.ErrParamInvalid
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) 校验目标货位。
		loc, err := s.findActiveLocation(ctx, tx, locationCode)
		if err != nil {
			return err
		}

		// 2) 加行锁读取追溯码。
		var trace inventoryModel.TraceInventory
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("trace_code = ?", traceCode).
			First(&trace).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}

		// 3) 仅允许 IN_STOCK 或 MISPLACED 状态的追溯码进行移位。
		if trace.Status != inventoryModel.TraceInventoryStatusInStock &&
			trace.Status != inventoryModel.TraceInventoryStatusMisplaced {
			return ecode.ErrTraceStatusLocked
		}

		fromLocationID := trace.LocationID
		toLocationID := loc.ID

		// 4) 更新货位，IN_STOCK 保持 IN_STOCK，MISPLACED 恢复为 IN_STOCK。
		if err := tx.Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", traceCode).
			Updates(map[string]interface{}{
				"location_id": toLocationID,
				"status":      inventoryModel.TraceInventoryStatusInStock,
			}).Error; err != nil {
			return err
		}

		// 5) 写入库存调整记录。
		fromStatus := trace.Status
		toStatus := inventoryModel.TraceInventoryStatusInStock
		reason := "货位调整"
		adj := &inventoryModel.InventoryAdjustment{
			TraceCode:      traceCode,
			DrugID:         trace.DrugID,
			AdjustType:     inventoryModel.AdjustTypeRelocate,
			BeforeStatus:   &fromStatus,
			AfterStatus:    &toStatus,
			FromLocationID: fromLocationID,
			ToLocationID:   &toLocationID,
			Reason:         reason,
			OperatorID:     operatorID,
		}
		if err := tx.Create(adj).Error; err != nil {
			return err
		}

		// 6) 写入追溯日志。
		actionType := "RELOCATION"
		log := &inventoryModel.DrugTraceLog{
			TraceCode:    traceCode,
			ActionType:   actionType,
			FromStatus:   &fromStatus,
			ToStatus:     &toStatus,
			ToLocationID: &toLocationID,
			OperatorID:   operatorID,
		}
		return tx.Create(log).Error
	})
}

// drugMixRow 用于查询货位药品汇总的中间结构。
type drugMixRow struct {
	DrugID         int64     `gorm:"column:drug_id"`
	CommonName     string    `gorm:"column:common_name"`
	Cnt            int64     `gorm:"column:cnt"`
	EarliestExpire time.Time `gorm:"column:earliest_expire"`
}

func (s *shelvingService) MixCheck(ctx context.Context, locationCode string) (*MixCheckResult, error) {
	locationCode = strings.TrimSpace(locationCode)
	if locationCode == "" {
		return nil, ecode.ErrParamInvalid
	}

	// 查询货位基本信息。
	var loc locationModel.LocationInfo
	err := s.db.WithContext(ctx).
		Where("location_code = ?", locationCode).
		First(&loc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrLocationNotFound
		}
		return nil, err
	}

	// 查询该货位在库药品按 drug_id 分组汇总。
	var rows []drugMixRow
	err = s.db.WithContext(ctx).
		Table("drug_trace_inventory t").
		Select("t.drug_id, d.common_name, COUNT(*) as cnt, MIN(t.expire_date) as earliest_expire").
		Joins("LEFT JOIN drug_info d ON d.id = t.drug_id AND d.deleted_at IS NULL").
		Where("t.location_id = ? AND t.status = ? AND t.deleted_at IS NULL", loc.ID, inventoryModel.TraceInventoryStatusInStock).
		Group("t.drug_id, d.common_name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	drugs := make([]DrugMixInfo, 0, len(rows))
	for _, row := range rows {
		// 查询该 drug_id 在此货位的所有批次号。
		var batchNumbers []string
		s.db.WithContext(ctx).
			Model(&inventoryModel.TraceInventory{}).
			Select("DISTINCT batch_number").
			Where("location_id = ? AND drug_id = ? AND status = ?", loc.ID, row.DrugID, inventoryModel.TraceInventoryStatusInStock).
			Pluck("batch_number", &batchNumbers)

		drugs = append(drugs, DrugMixInfo{
			DrugID:         row.DrugID,
			DrugName:       row.CommonName,
			Count:          row.Cnt,
			BatchNumbers:   batchNumbers,
			EarliestExpire: row.EarliestExpire,
		})
	}

	return &MixCheckResult{
		LocationCode:  loc.LocationCode,
		LocationName:  loc.LocationName,
		HasMixedDrugs: len(drugs) > 1,
		Drugs:         drugs,
	}, nil
}
