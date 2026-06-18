package service

import (
	"context"
	"errors"
	"time"

	drugModel "github.com/YingmoY/PharmacyERP/internal/drug/model"
	inventoryModel "github.com/YingmoY/PharmacyERP/internal/inventory/model"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	salesModel "github.com/YingmoY/PharmacyERP/internal/sales/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// traceLogRow is an internal scan target for enriched drug_trace_log queries.
type traceLogRow struct {
	ID               int64     `gorm:"column:id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	TraceCode        string    `gorm:"column:trace_code"`
	ActionType       string    `gorm:"column:action_type"`
	FromStatus       *string   `gorm:"column:from_status"`
	ToStatus         *string   `gorm:"column:to_status"`
	OperatorID       int64     `gorm:"column:operator_id"`
	OperatorName     string    `gorm:"column:operator_name"`
	RelatedNo        *string   `gorm:"column:related_no"`
	Remark           *string   `gorm:"column:remark"`
	DrugID           *int64    `gorm:"column:drug_id"`
	DrugName         string    `gorm:"column:drug_name"`
	OrderID          *int64    `gorm:"column:order_id"`
	OrderItemID      *int64    `gorm:"column:order_item_id"`
	RequestID        *string   `gorm:"column:request_id"`
	FromLocationID   *int64    `gorm:"column:from_location_id"`
	ToLocationID     *int64    `gorm:"column:to_location_id"`
	FromLocationCode string    `gorm:"column:from_location_code"`
	ToLocationCode   string    `gorm:"column:to_location_code"`
}

func rowToLog(r traceLogRow) *inventoryModel.DrugTraceLog {
	l := &inventoryModel.DrugTraceLog{}
	l.ID = r.ID
	l.CreatedAt = r.CreatedAt
	l.UpdatedAt = r.UpdatedAt
	l.TraceCode = r.TraceCode
	l.ActionType = r.ActionType
	l.FromStatus = r.FromStatus
	l.ToStatus = r.ToStatus
	l.OperatorID = r.OperatorID
	l.OperatorName = r.OperatorName
	l.RelatedNo = r.RelatedNo
	l.Remark = r.Remark
	l.DrugID = r.DrugID
	l.DrugName = r.DrugName
	l.OrderID = r.OrderID
	l.OrderItemID = r.OrderItemID
	l.RequestID = r.RequestID
	l.FromLocationID = r.FromLocationID
	l.ToLocationID = r.ToLocationID
	l.FromLocationCode = r.FromLocationCode
	l.ToLocationCode = r.ToLocationCode
	return l
}

// TraceInfoResult 追溯码详情（含药品信息）
type TraceInfoResult struct {
	TraceCode   string  `json:"trace_code"`
	Status      string  `json:"status"`
	DrugID      int64   `json:"drug_id"`
	DrugName    string  `json:"drug_name"`
	DrugCode    string  `json:"drug_code"`
	BatchNumber string  `json:"batch_number"`
	ExpireDate  string  `json:"expire_date"`
	LocationID  *int64  `json:"location_id,omitempty"`
	LocationCode string `json:"location_code,omitempty"`
	IsReserved  bool    `json:"is_reserved"`
	ReservationNo *string `json:"reservation_no,omitempty"`
}

// ValidateResult 追溯码验证结果
type ValidateResult struct {
	TraceCode   string `json:"trace_code"`
	Exists      bool   `json:"exists"`
	Status      string `json:"status,omitempty"`
	DrugID      int64  `json:"drug_id,omitempty"`
	DrugName    string `json:"drug_name,omitempty"`
	BatchNumber string `json:"batch_number,omitempty"`
	ExpireDate  string `json:"expire_date,omitempty"`
	LocationCode string `json:"location_code,omitempty"`
	IsAvailable bool   `json:"is_available"`
	IsReserved  bool   `json:"is_reserved"`
	Reason      string `json:"reason,omitempty"`
}

// ValidateReq 追溯码验证请求
type ValidateReq struct {
	TraceCode string `json:"trace_code" binding:"required"`
}

// LogFilter 追溯日志查询过滤条件
type LogFilter struct {
	Page     int
	PageSize int
}

// TraceService 定义追溯码查询能力
type TraceService interface {
	// GetTraceInfo 查询追溯码库存状态及药品信息
	GetTraceInfo(ctx context.Context, traceCode string) (*TraceInfoResult, error)
	// GetFullChain 查询追溯码完整操作链（时间倒序）
	GetFullChain(ctx context.Context, traceCode string) ([]*inventoryModel.DrugTraceLog, error)
	// GetLogs 分页查询追溯日志
	GetLogs(ctx context.Context, traceCode string, filter LogFilter) ([]*inventoryModel.DrugTraceLog, int64, error)
	// Validate 验证追溯码是否可用于销售
	Validate(ctx context.Context, req ValidateReq) (*ValidateResult, error)
}

type traceService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewTraceService 创建追溯码服务实例
func NewTraceService(db *gorm.DB, logger *zap.Logger) TraceService {
	return &traceService{db: db, logger: logger}
}

// GetTraceInfo 查询追溯码详情（含药品信息和货位）
func (s *traceService) GetTraceInfo(ctx context.Context, traceCode string) (*TraceInfoResult, error) {
	// 1. 查询追溯库存
	var trace inventoryModel.TraceInventory
	if err := s.db.WithContext(ctx).Where("trace_code = ?", traceCode).First(&trace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrTraceCodeNotFound
		}
		return nil, err
	}

	result := &TraceInfoResult{
		TraceCode:   trace.TraceCode,
		Status:      trace.Status,
		DrugID:      trace.DrugID,
		BatchNumber: trace.BatchNumber,
		ExpireDate:  trace.ExpireDate.Format("2006-01-02"),
		LocationID:  trace.LocationID,
	}

	// 2. 查询药品信息
	var drug drugModel.DrugInfo
	if err := s.db.WithContext(ctx).Where("id = ?", trace.DrugID).First(&drug).Error; err == nil {
		result.DrugName = drug.CommonName
		result.DrugCode = drug.DrugCode
	}

	// 3. 查询货位信息
	if trace.LocationID != nil {
		var loc locationModel.LocationInfo
		if err := s.db.WithContext(ctx).Where("id = ?", *trace.LocationID).First(&loc).Error; err == nil {
			result.LocationCode = loc.LocationCode
		}
	}

	// 4. 检查是否被预留
	var rsv salesModel.TraceReservation
	if err := s.db.WithContext(ctx).
		Where("trace_code = ? AND status = ?", traceCode, salesModel.ReservationStatusReserved).
		First(&rsv).Error; err == nil {
		result.IsReserved = true
		result.ReservationNo = &rsv.ReservationNo
	}

	return result, nil
}

// enrichedLogQuery returns a base query that JOINs operator name, drug name, and location codes.
func (s *traceService) enrichedLogQuery(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("drug_trace_log dtl").
		Joins("LEFT JOIN sys_user u ON u.id = dtl.operator_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN drug_trace_inventory dti ON dti.trace_code = dtl.trace_code AND dti.deleted_at IS NULL").
		Joins("LEFT JOIN drug_info di ON di.id = COALESCE(dtl.drug_id, dti.drug_id) AND di.deleted_at IS NULL").
		Joins("LEFT JOIN location_info lf ON lf.id = dtl.from_location_id AND lf.deleted_at IS NULL").
		Joins("LEFT JOIN location_info lt ON lt.id = dtl.to_location_id AND lt.deleted_at IS NULL").
		Select(`dtl.id, dtl.created_at, dtl.updated_at,
			dtl.trace_code, dtl.action_type,
			dtl.from_status, dtl.to_status,
			dtl.operator_id, COALESCE(u.real_name, u.username, '') AS operator_name,
			dtl.related_no, dtl.remark,
			dtl.drug_id, COALESCE(di.common_name, '') AS drug_name,
			dtl.order_id, dtl.order_item_id, dtl.request_id,
			dtl.from_location_id, dtl.to_location_id,
			COALESCE(lf.location_code, '') AS from_location_code,
			COALESCE(lt.location_code, '') AS to_location_code`)
}

// GetFullChain 查询追溯码完整操作链（时间倒序）
func (s *traceService) GetFullChain(ctx context.Context, traceCode string) ([]*inventoryModel.DrugTraceLog, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&inventoryModel.TraceInventory{}).
		Where("trace_code = ?", traceCode).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ecode.ErrTraceCodeNotFound
	}

	var rows []traceLogRow
	if err := s.enrichedLogQuery(ctx).
		Where("dtl.trace_code = ? AND dtl.deleted_at IS NULL", traceCode).
		Order("dtl.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	logs := make([]*inventoryModel.DrugTraceLog, 0, len(rows))
	for _, r := range rows {
		logs = append(logs, rowToLog(r))
	}
	return logs, nil
}

// GetLogs 分页查询追溯日志
func (s *traceService) GetLogs(ctx context.Context, traceCode string, filter LogFilter) ([]*inventoryModel.DrugTraceLog, int64, error) {
	var traceCount int64
	if err := s.db.WithContext(ctx).Model(&inventoryModel.TraceInventory{}).
		Where("trace_code = ?", traceCode).Count(&traceCount).Error; err != nil {
		return nil, 0, err
	}
	if traceCount == 0 {
		return nil, 0, ecode.ErrTraceCodeNotFound
	}

	var total int64
	if err := s.db.WithContext(ctx).Table("drug_trace_log").
		Where("trace_code = ? AND deleted_at IS NULL", traceCode).
		Count(&total).Error; err != nil {
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

	var rows []traceLogRow
	if err := s.enrichedLogQuery(ctx).
		Where("dtl.trace_code = ? AND dtl.deleted_at IS NULL", traceCode).
		Order("dtl.created_at DESC").
		Offset(offset).Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	logs := make([]*inventoryModel.DrugTraceLog, 0, len(rows))
	for _, r := range rows {
		logs = append(logs, rowToLog(r))
	}
	return logs, total, nil
}

// Validate 验证追溯码是否可用于销售
func (s *traceService) Validate(ctx context.Context, req ValidateReq) (*ValidateResult, error) {
	result := &ValidateResult{
		TraceCode: req.TraceCode,
	}

	// 1. 查询追溯库存
	var trace inventoryModel.TraceInventory
	if err := s.db.WithContext(ctx).Where("trace_code = ?", req.TraceCode).First(&trace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Exists = false
			result.IsAvailable = false
			result.Reason = "trace code not found"
			return result, nil
		}
		return nil, err
	}

	result.Exists = true
	result.Status = trace.Status
	result.DrugID = trace.DrugID
	result.BatchNumber = trace.BatchNumber
	result.ExpireDate = trace.ExpireDate.Format("2006-01-02")

	// 2. 查询药品信息
	var drug drugModel.DrugInfo
	if err := s.db.WithContext(ctx).Where("id = ?", trace.DrugID).First(&drug).Error; err == nil {
		result.DrugName = drug.CommonName
	}

	// 3. 查询货位
	if trace.LocationID != nil {
		var loc locationModel.LocationInfo
		if err := s.db.WithContext(ctx).Where("id = ?", *trace.LocationID).First(&loc).Error; err == nil {
			result.LocationCode = loc.LocationCode
		}
	}

	// 4. 检查状态是否可用
	if trace.Status != inventoryModel.TraceInventoryStatusInStock {
		result.IsAvailable = false
		result.Reason = "trace code is not IN_STOCK (status: " + trace.Status + ")"
		return result, nil
	}

	// 5. 检查是否被预留
	var rsvCount int64
	if err := s.db.WithContext(ctx).Model(&salesModel.TraceReservation{}).
		Where("trace_code = ? AND status = ?", req.TraceCode, salesModel.ReservationStatusReserved).
		Count(&rsvCount).Error; err != nil {
		return nil, err
	}

	result.IsReserved = rsvCount > 0
	if rsvCount > 0 {
		result.IsAvailable = false
		result.Reason = "trace code is already reserved"
		return result, nil
	}

	result.IsAvailable = true
	return result, nil
}

