package service

import (
	"context"
	"errors"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/inventory/model"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InventoryService 定义库存查询与调整能力。
type InventoryService interface {
	// ListInventory 分页查询追溯库存列表，支持多条件过滤。
	ListInventory(ctx context.Context, req ListInventoryRequest) (*core.PageResult, error)

	// GetInventorySummary 统计各状态下的库存数量，以及近效期、预约数量。
	GetInventorySummary(ctx context.Context) (*InventorySummary, error)

	// ListPendingShelving 查询待上架追溯码（状态=PENDING，关联入库单已完成）。
	ListPendingShelving(ctx context.Context, page, pageSize int) (*core.PageResult, error)

	// ListNearExpire 查询近30天内到期的在库药品追溯码。
	ListNearExpire(ctx context.Context, page, pageSize int) (*core.PageResult, error)

	// ListRecommendSale 推荐销售顺序（FEFO：先到期先出），仅限在库未预约追溯码。
	ListRecommendSale(ctx context.Context, drugID int64, page, pageSize int) (*core.PageResult, error)

	// ListByDrug 查询某药品 ID 下的所有追溯码记录。
	ListByDrug(ctx context.Context, drugID int64, page, pageSize int) (*core.PageResult, error)

	// ListByLocation 查询某货位 ID 下的所有追溯码记录。
	ListByLocation(ctx context.Context, locationID int64, page, pageSize int) (*core.PageResult, error)

	// ManualStatusChange 手动变更追溯码状态（仅管理员），写 inventory_adjustment + trace_log。
	ManualStatusChange(ctx context.Context, req ManualStatusChangeRequest) error

	// ListAdjustments 分页查询库存调整记录。
	ListAdjustments(ctx context.Context, req ListAdjustmentsRequest) (*core.PageResult, error)

	// GetAdjustment 查询单条调整记录。
	GetAdjustment(ctx context.Context, id int64) (*model.InventoryAdjustment, error)

	// CreateAdjustment 手动创建调整记录（除ManualStatusChange外的通用入口）。
	CreateAdjustment(ctx context.Context, req CreateAdjustmentRequest) (*model.InventoryAdjustment, error)
}

// ListInventoryRequest 库存列表查询条件。
type ListInventoryRequest struct {
	DrugID          *int64
	LocationID      *int64
	Status          string
	BatchNumber     string
	ExpireDateStart *time.Time
	ExpireDateEnd   *time.Time
	Page            int
	PageSize        int
}

// InventorySummary 库存统计汇总。
type InventorySummary struct {
	TotalInStock       int64 `json:"total_in_stock"`
	TotalPending       int64 `json:"total_pending"`
	TotalSold          int64 `json:"total_sold"`
	TotalMisplaced     int64 `json:"total_misplaced"`
	TotalLossCandidate int64 `json:"total_loss_candidate"`
	TotalLost          int64 `json:"total_lost"`
	NearExpireCount    int64 `json:"near_expire_count"`
}

// ManualStatusChangeRequest 手动状态变更请求。
type ManualStatusChangeRequest struct {
	TraceCode  string
	ToStatus   string
	Reason     string
	Remark     string
	OperatorID int64
}

// ListAdjustmentsRequest 调整记录查询条件。
type ListAdjustmentsRequest struct {
	TraceCode  string
	DrugID     *int64
	AdjustType string
	Page       int
	PageSize   int
}

// CreateAdjustmentRequest 手动创建调整请求。
type CreateAdjustmentRequest struct {
	TraceCode      string
	AdjustType     string
	FromLocationID *int64
	ToLocationID   *int64
	Reason         string
	Remark         string
	OperatorID     int64
}

// inventoryService 是 InventoryService 的默认实现。
type inventoryService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewInventoryService 创建库存查询与调整服务。
func NewInventoryService(db *gorm.DB, logger *zap.Logger) InventoryService {
	return &inventoryService{db: db, logger: logger}
}

func (s *inventoryService) ListInventory(ctx context.Context, req ListInventoryRequest) (*core.PageResult, error) {
	query := s.db.WithContext(ctx).Model(&model.TraceInventory{})

	if req.DrugID != nil {
		query = query.Where("drug_id = ?", *req.DrugID)
	}
	if req.LocationID != nil {
		query = query.Where("location_id = ?", *req.LocationID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.BatchNumber != "" {
		query = query.Where("batch_number = ?", req.BatchNumber)
	}
	if req.ExpireDateStart != nil {
		query = query.Where("expire_date >= ?", *req.ExpireDateStart)
	}
	if req.ExpireDateEnd != nil {
		query = query.Where("expire_date <= ?", *req.ExpireDateEnd)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.TraceInventory
	if err := query.
		Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	s.enrichTraceInventory(ctx, list)
	result := core.NewPageResult(total, req.Page, req.PageSize, list)
	return &result, nil
}

// enrichTraceInventory 批量填充追溯库存的药品名称、货位编码和入库单号。
func (s *inventoryService) enrichTraceInventory(ctx context.Context, list []model.TraceInventory) {
	if len(list) == 0 {
		return
	}
	drugIDs := make([]int64, 0, len(list))
	locationIDs := make([]int64, 0, len(list))
	orderIDs := make([]int64, 0, len(list))
	for i := range list {
		drugIDs = append(drugIDs, list[i].DrugID)
		if list[i].LocationID != nil {
			locationIDs = append(locationIDs, *list[i].LocationID)
		}
		orderIDs = append(orderIDs, list[i].InboundOrderID)
	}

	type drugRow struct {
		ID            int64
		CommonName    string
		Specification string
		Manufacturer  string
	}
	var drugs []drugRow
	s.db.WithContext(ctx).Table("drug_info").Select("id, common_name, specification, manufacturer").Where("id IN ?", drugIDs).Scan(&drugs)
	drugMap := make(map[int64]drugRow, len(drugs))
	for _, d := range drugs {
		drugMap[d.ID] = d
	}

	locMap := make(map[int64]string)
	if len(locationIDs) > 0 {
		type locRow struct {
			ID           int64
			LocationCode string
		}
		var locs []locRow
		s.db.WithContext(ctx).Table("location_info").Select("id, location_code").Where("id IN ?", locationIDs).Scan(&locs)
		for _, l := range locs {
			locMap[l.ID] = l.LocationCode
		}
	}

	type orderRow struct {
		ID      int64
		OrderNo string
	}
	var orders []orderRow
	s.db.WithContext(ctx).Table("inbound_order").Select("id, order_no").Where("id IN ?", orderIDs).Scan(&orders)
	orderMap := make(map[int64]string, len(orders))
	for _, o := range orders {
		orderMap[o.ID] = o.OrderNo
	}

	for i := range list {
		if d, ok := drugMap[list[i].DrugID]; ok {
			list[i].DrugName = d.CommonName
			list[i].Specification = d.Specification
			list[i].Manufacturer = d.Manufacturer
		}
		if list[i].LocationID != nil {
			code := locMap[*list[i].LocationID]
			list[i].LocationCode = code
			list[i].SystemLocationCode = code
		}
		list[i].InboundOrderNo = orderMap[list[i].InboundOrderID]
	}
}

func (s *inventoryService) GetInventorySummary(ctx context.Context) (*InventorySummary, error) {
	type statusCount struct {
		Status string
		Cnt    int64
	}

	var counts []statusCount
	err := s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Select("status, COUNT(*) as cnt").
		Group("status").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	summary := &InventorySummary{}
	for _, c := range counts {
		switch c.Status {
		case model.TraceInventoryStatusInStock:
			summary.TotalInStock = c.Cnt
		case model.TraceInventoryStatusPending:
			summary.TotalPending = c.Cnt
		case model.TraceInventoryStatusSold:
			summary.TotalSold = c.Cnt
		case model.TraceInventoryStatusMisplaced:
			summary.TotalMisplaced = c.Cnt
		case model.TraceInventoryStatusLossCandidate:
			summary.TotalLossCandidate = c.Cnt
		case model.TraceInventoryStatusLost:
			summary.TotalLost = c.Cnt
		}
	}

	// 统计近30天到期数量（仅统计在库药品）。
	nearExpireDeadline := time.Now().AddDate(0, 0, 30)
	var nearExpireCount int64
	err = s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("status = ? AND expire_date <= ?", model.TraceInventoryStatusInStock, nearExpireDeadline).
		Count(&nearExpireCount).Error
	if err != nil {
		return nil, err
	}
	summary.NearExpireCount = nearExpireCount

	return summary, nil
}

func (s *inventoryService) ListPendingShelving(ctx context.Context, page, pageSize int) (*core.PageResult, error) {
	// 查询状态为PENDING且关联入库单已完成的追溯码。
	query := s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Joins("JOIN inbound_order ON inbound_order.id = drug_trace_inventory.inbound_order_id AND inbound_order.deleted_at IS NULL").
		Where("drug_trace_inventory.status = ?", model.TraceInventoryStatusPending).
		Where("inbound_order.status = ?", model.InboundOrderStatusCompleted)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.TraceInventory
	if err := query.
		Order("drug_trace_inventory.created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	s.enrichTraceInventory(ctx, list)
	result := core.NewPageResult(total, page, pageSize, list)
	return &result, nil
}

// NearExpireItem 近效期药品汇总条目（按药品+批号+货位聚合）。
type NearExpireItem struct {
	DrugID        int64  `json:"drug_id"`
	DrugName      string `json:"drug_name"`
	Specification string `json:"specification"`
	BatchNumber   string `json:"batch_number"`
	ExpireDate    string `json:"expire_date"`
	RemainingDays int    `json:"remaining_days"`
	LocationCode  string `json:"location_code"`
	Count         int64  `json:"count"`
	AlertLevel    string `json:"alert_level"`
}

func (s *inventoryService) ListNearExpire(ctx context.Context, page, pageSize int) (*core.PageResult, error) {
	// 使用聚合查询，按药品+批号+货位分组，统计各组件数量及最近到期日。
	type rawRow struct {
		DrugID       int64
		CommonName   string
		Specification string
		BatchNumber  string
		ExpireDate   string
		LocationCode string
		Count        int64
	}

	countSQL := `
		SELECT COUNT(*) FROM (
			SELECT ti.drug_id, ti.batch_number, ti.expire_date, ti.location_id
			FROM drug_trace_inventory ti
			WHERE ti.status = 'IN_STOCK'
			  AND ti.expire_date <= CURRENT_DATE + INTERVAL '30 days'
			  AND ti.deleted_at IS NULL
			GROUP BY ti.drug_id, ti.batch_number, ti.expire_date, ti.location_id
		) sub`

	var total int64
	if err := s.db.WithContext(ctx).Raw(countSQL).Scan(&total).Error; err != nil {
		return nil, err
	}

	listSQL := `
		SELECT
			d.id        AS drug_id,
			d.common_name AS common_name,
			d.specification,
			ti.batch_number,
			TO_CHAR(ti.expire_date, 'YYYY-MM-DD') AS expire_date,
			COALESCE(l.location_code, '') AS location_code,
			COUNT(*) AS count
		FROM drug_trace_inventory ti
		JOIN drug_info d ON d.id = ti.drug_id AND d.deleted_at IS NULL
		LEFT JOIN location_info l ON l.id = ti.location_id
		WHERE ti.status = 'IN_STOCK'
		  AND ti.expire_date <= CURRENT_DATE + INTERVAL '30 days'
		  AND ti.deleted_at IS NULL
		GROUP BY d.id, d.common_name, d.specification, ti.batch_number, ti.expire_date, l.location_code
		ORDER BY ti.expire_date ASC
		LIMIT ? OFFSET ?`

	var rows []rawRow
	offset := (page - 1) * pageSize
	if err := s.db.WithContext(ctx).Raw(listSQL, pageSize, offset).Scan(&rows).Error; err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	items := make([]NearExpireItem, 0, len(rows))
	for _, r := range rows {
		expDate, _ := time.Parse("2006-01-02", r.ExpireDate)
		remaining := int(expDate.Sub(today).Hours() / 24)
		if remaining < 0 {
			remaining = 0
		}
		level := "LOW"
		if remaining <= 7 {
			level = "HIGH"
		} else if remaining <= 15 {
			level = "MEDIUM"
		}
		items = append(items, NearExpireItem{
			DrugID:        r.DrugID,
			DrugName:      r.CommonName,
			Specification: r.Specification,
			BatchNumber:   r.BatchNumber,
			ExpireDate:    r.ExpireDate,
			RemainingDays: remaining,
			LocationCode:  r.LocationCode,
			Count:         r.Count,
			AlertLevel:    level,
		})
	}

	result := core.NewPageResult(total, page, pageSize, items)
	return &result, nil
}

func (s *inventoryService) ListRecommendSale(ctx context.Context, drugID int64, page, pageSize int) (*core.PageResult, error) {
	if drugID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	// 查询在库、未被预约的追溯码，按到期日升序（FEFO先到期先出）。
	query := s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("drug_id = ? AND status = ?", drugID, model.TraceInventoryStatusInStock).
		Where("trace_code NOT IN (SELECT trace_code FROM trace_reservation WHERE status = 'RESERVED')")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.TraceInventory
	if err := query.
		Order("expire_date ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := core.NewPageResult(total, page, pageSize, list)
	return &result, nil
}

func (s *inventoryService) ListByDrug(ctx context.Context, drugID int64, page, pageSize int) (*core.PageResult, error) {
	query := s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("drug_id = ?", drugID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.TraceInventory
	if err := query.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	s.enrichTraceInventory(ctx, list)
	result := core.NewPageResult(total, page, pageSize, list)
	return &result, nil
}

func (s *inventoryService) ListByLocation(ctx context.Context, locationID int64, page, pageSize int) (*core.PageResult, error) {
	query := s.db.WithContext(ctx).
		Model(&model.TraceInventory{}).
		Where("location_id = ?", locationID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.TraceInventory
	if err := query.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := core.NewPageResult(total, page, pageSize, list)
	return &result, nil
}

func (s *inventoryService) ManualStatusChange(ctx context.Context, req ManualStatusChangeRequest) error {
	if req.TraceCode == "" || req.ToStatus == "" || req.Reason == "" || req.OperatorID <= 0 {
		return ecode.ErrParamInvalid
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 加行锁读取追溯码当前状态。
		var trace model.TraceInventory
		err := tx.Clauses().
			Where("trace_code = ?", req.TraceCode).
			First(&trace).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrTraceCodeNotFound
			}
			return err
		}

		fromStatus := trace.Status
		toStatus := req.ToStatus

		// 更新追溯码状态。
		if err := tx.Model(&model.TraceInventory{}).
			Where("trace_code = ?", req.TraceCode).
			Update("status", toStatus).Error; err != nil {
			return err
		}

		// 写入库存调整记录。
		adj := &model.InventoryAdjustment{
			TraceCode:    req.TraceCode,
			DrugID:       trace.DrugID,
			AdjustType:   model.AdjustTypeStatusChange,
			BeforeStatus: &fromStatus,
			AfterStatus:  &toStatus,
			Reason:       req.Reason,
			OperatorID:   req.OperatorID,
		}
		if err := tx.Create(adj).Error; err != nil {
			return err
		}

		// 写入追溯日志。
		relatedNo := ""
		actionType := "INVENTORY"
		remark := req.Remark
		log := &model.DrugTraceLog{
			TraceCode:  req.TraceCode,
			ActionType: actionType,
			FromStatus: &fromStatus,
			ToStatus:   &toStatus,
			OperatorID: req.OperatorID,
			RelatedNo:  &relatedNo,
			Remark:     &remark,
		}
		if trace.LocationID != nil {
			log.FromLocationID = trace.LocationID
		}
		return tx.Create(log).Error
	})
}

func (s *inventoryService) ListAdjustments(ctx context.Context, req ListAdjustmentsRequest) (*core.PageResult, error) {
	query := s.db.WithContext(ctx).Model(&model.InventoryAdjustment{})

	if req.TraceCode != "" {
		query = query.Where("trace_code = ?", req.TraceCode)
	}
	if req.DrugID != nil {
		query = query.Where("drug_id = ?", *req.DrugID)
	}
	if req.AdjustType != "" {
		query = query.Where("adjust_type = ?", req.AdjustType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.InventoryAdjustment
	if err := query.
		Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	s.enrichAdjustments(ctx, list)
	result := core.NewPageResult(total, req.Page, req.PageSize, list)
	return &result, nil
}

// enrichAdjustments 批量填充库存调整记录的药品名称、货位编码和操作人名称。
func (s *inventoryService) enrichAdjustments(ctx context.Context, list []model.InventoryAdjustment) {
	if len(list) == 0 {
		return
	}
	drugIDs := make([]int64, 0, len(list))
	locIDs := make([]int64, 0)
	opIDs := make([]int64, 0, len(list))
	for i := range list {
		drugIDs = append(drugIDs, list[i].DrugID)
		opIDs = append(opIDs, list[i].OperatorID)
		if list[i].FromLocationID != nil {
			locIDs = append(locIDs, *list[i].FromLocationID)
		}
		if list[i].ToLocationID != nil {
			locIDs = append(locIDs, *list[i].ToLocationID)
		}
	}

	type drugRow struct {
		ID         int64
		CommonName string
	}
	var drugs []drugRow
	s.db.WithContext(ctx).Table("drug_info").Select("id, common_name").Where("id IN ?", drugIDs).Scan(&drugs)
	drugMap := make(map[int64]string, len(drugs))
	for _, d := range drugs {
		drugMap[d.ID] = d.CommonName
	}

	locMap := make(map[int64]string)
	if len(locIDs) > 0 {
		type locRow struct {
			ID           int64
			LocationCode string
		}
		var locs []locRow
		s.db.WithContext(ctx).Table("location_info").Select("id, location_code").Where("id IN ?", locIDs).Scan(&locs)
		for _, l := range locs {
			locMap[l.ID] = l.LocationCode
		}
	}

	type userRow struct {
		ID       int64
		RealName string
	}
	var users []userRow
	s.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", opIDs).Scan(&users)
	userMap := make(map[int64]string, len(users))
	for _, u := range users {
		userMap[u.ID] = u.RealName
	}

	for i := range list {
		list[i].DrugName = drugMap[list[i].DrugID]
		list[i].OperatorName = userMap[list[i].OperatorID]
		if list[i].FromLocationID != nil {
			list[i].FromLocationCode = locMap[*list[i].FromLocationID]
		}
		if list[i].ToLocationID != nil {
			list[i].ToLocationCode = locMap[*list[i].ToLocationID]
		}
	}
}

func (s *inventoryService) GetAdjustment(ctx context.Context, id int64) (*model.InventoryAdjustment, error) {
	var adj model.InventoryAdjustment
	err := s.db.WithContext(ctx).First(&adj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return &adj, nil
}

func (s *inventoryService) CreateAdjustment(ctx context.Context, req CreateAdjustmentRequest) (*model.InventoryAdjustment, error) {
	if req.TraceCode == "" || req.AdjustType == "" || req.Reason == "" || req.OperatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	var trace model.TraceInventory
	err := s.db.WithContext(ctx).Where("trace_code = ?", req.TraceCode).First(&trace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrTraceCodeNotFound
		}
		return nil, err
	}

	fromStatus := trace.Status

	adj := &model.InventoryAdjustment{
		TraceCode:      req.TraceCode,
		DrugID:         trace.DrugID,
		AdjustType:     req.AdjustType,
		BeforeStatus:   &fromStatus,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Reason:         req.Reason,
		OperatorID:     req.OperatorID,
	}

	// 需要同时检查货位是否存在（如果提供了 ToLocationID）。
	if req.ToLocationID != nil {
		var loc locationModel.LocationInfo
		err := s.db.WithContext(ctx).
			Where("id = ? AND status = ?", *req.ToLocationID, locationModel.LocationStatusEnabled).
			First(&loc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ecode.ErrLocationNotFound
			}
			return nil, err
		}
	}

	if err := s.db.WithContext(ctx).Create(adj).Error; err != nil {
		return nil, err
	}

	return adj, nil
}
