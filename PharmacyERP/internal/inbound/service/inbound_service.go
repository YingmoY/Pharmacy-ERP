// Package service 实现入库单核心业务逻辑。
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/inbound/model"
	"github.com/YingmoY/PharmacyERP/internal/inbound/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// =====================================================
// DTO 定义
// =====================================================

// CreateOrderReq 创建入库单请求。
type CreateOrderReq struct {
	SupplierID int64           `json:"supplier_id"`
	InvoiceNo  string          `json:"invoice_no"`
	Remark     string          `json:"remark"`
	Details    []CreateDetailReq `json:"details"`
}

// CreateDetailReq 创建明细行请求。
type CreateDetailReq struct {
	DrugID      int64   `json:"drug_id"`
	BatchNumber string  `json:"batch_number"`
	ExpireDate  string  `json:"expire_date"` // YYYY-MM-DD
	PlannedQty  int32   `json:"planned_qty"`
	UnitPrice   float64 `json:"unit_price"`
	Remark      string  `json:"remark"`
}

// UpdateOrderReq 更新入库单请求（仅 DRAFT 状态下可用）。
type UpdateOrderReq struct {
	SupplierID int64   `json:"supplier_id"`
	InvoiceNo  string  `json:"invoice_no"`
	Remark     string  `json:"remark"`
}

// UpdateDetailReq 更新明细行请求。
type UpdateDetailReq struct {
	BatchNumber string  `json:"batch_number"`
	ExpireDate  string  `json:"expire_date"`
	PlannedQty  int32   `json:"planned_qty"`
	UnitPrice   float64 `json:"unit_price"`
	Remark      string  `json:"remark"`
}

// ConfirmTraceCodesReq 批量确认追溯码请求。
type ConfirmTraceCodesReq struct {
	OrderID    int64
	DetailID   int64
	TraceCodes []string
	OperatorID int64
}

// ProgressDetail 入库进度明细行。
type ProgressDetail struct {
	DetailID     int64  `json:"detail_id"`
	DrugID       int64  `json:"drug_id"`
	DrugName     string `json:"drug_name"`
	BatchNumber  string `json:"batch_number"`
	PlannedQty   int32  `json:"planned_qty"`
	ConfirmedQty int32  `json:"confirmed_qty"`
}

// InboundProgressDTO 入库进度响应 DTO。
type InboundProgressDTO struct {
	OrderID        int64            `json:"order_id"`
	OrderNo        string           `json:"order_no"`
	Status         string           `json:"status"`
	TotalPlanned   int32            `json:"total_planned"`
	TotalConfirmed int32            `json:"total_confirmed"`
	Details        []ProgressDetail `json:"details"`
}

// =====================================================
// Service 接口
// =====================================================

// InboundService 入库单核心业务接口。
type InboundService interface {
	// CreateOrder 创建入库单，可同时携带明细行。
	CreateOrder(ctx context.Context, req CreateOrderReq, creatorID int64) (*model.InboundOrder, error)
	// SubmitOrder 提交入库单：DRAFT → PENDING_CONFIRM。
	SubmitOrder(ctx context.Context, orderID, operatorID int64) error
	// GetOrder 查询入库单详情（含明细）。
	GetOrder(ctx context.Context, orderID int64) (*model.InboundOrder, error)
	// ListOrders 分页查询入库单列表。
	ListOrders(ctx context.Context, filter repository.ListFilter) ([]*model.InboundOrder, int64, error)
	// UpdateOrder 更新草稿入库单基础信息。
	UpdateOrder(ctx context.Context, orderID int64, req UpdateOrderReq, operatorID int64) error
	// AddDetail 新增明细行（仅 DRAFT），返回已创建的明细。
	AddDetail(ctx context.Context, orderID int64, req CreateDetailReq, operatorID int64) (*model.InboundOrderDetail, error)
	// UpdateDetail 更新明细行（仅 DRAFT），返回更新后的明细。
	UpdateDetail(ctx context.Context, orderID, detailID int64, req UpdateDetailReq, operatorID int64) (*model.InboundOrderDetail, error)
	// DeleteDetail 删除明细行（仅 DRAFT）。
	DeleteDetail(ctx context.Context, orderID, detailID int64, operatorID int64) error
	// ConfirmTraceCodes 扫码确认追溯码（PENDING_CONFIRM 状态，支持批量）。
	ConfirmTraceCodes(ctx context.Context, req ConfirmTraceCodesReq) error
	// CompleteOrder 完成入库单：PENDING_CONFIRM → COMPLETED（需所有明细均已全量确认）。
	CompleteOrder(ctx context.Context, orderID, operatorID int64) error
	// CancelOrder 取消入库单：DRAFT → CANCELLED 或 PENDING_CONFIRM → CANCELLED（无 IN_STOCK 追溯码）。
	CancelOrder(ctx context.Context, orderID, operatorID int64) error
	// GetInboundProgress 查询入库确认进度。
	GetInboundProgress(ctx context.Context, orderID int64) (*InboundProgressDTO, error)
	// CreateOrderWithDetails 由 AI 发票转换等内部调用，在现有事务内创建入库单+明细。
	CreateOrderWithDetails(ctx context.Context, tx *gorm.DB, req CreateOrderReq, creatorID int64) (*model.InboundOrder, error)
}

// =====================================================
// Service 实现
// =====================================================

// inboundService 是 InboundService 的默认实现。
type inboundService struct {
	db       *gorm.DB
	repo     repository.InboundRepo
	mqClient *mq.Client
	logger   *zap.Logger
}

// NewInboundService 创建入库服务实例。
func NewInboundService(
	db *gorm.DB,
	repo repository.InboundRepo,
	mqClient *mq.Client,
	logger *zap.Logger,
) InboundService {
	return &inboundService{
		db:       db,
		repo:     repo,
		mqClient: mqClient,
		logger:   logger,
	}
}

// =====================================================
// 辅助方法
// =====================================================

// parseExpireDate 解析 YYYY-MM-DD 格式日期。
func parseExpireDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("expire_date 不能为空")
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("expire_date 格式无效，须为 YYYY-MM-DD")
	}
	return t, nil
}

// validateCreateDetail 校验明细创建请求字段。
func validateCreateDetail(req CreateDetailReq) error {
	if req.DrugID <= 0 {
		return ecode.ErrParamInvalid
	}
	if req.BatchNumber == "" {
		return ecode.ErrParamInvalid
	}
	if req.PlannedQty <= 0 {
		return ecode.ErrParamInvalid
	}
	if req.UnitPrice < 0 {
		return ecode.ErrParamInvalid
	}
	return nil
}

// checkSupplierActive 校验供应商是否存在且处于启用状态（status=1）。
func checkSupplierActive(ctx context.Context, tx *gorm.DB, supplierID int64) error {
	var status int8
	err := tx.WithContext(ctx).
		Table("supplier").
		Select("status").
		Where("id = ? AND deleted_at IS NULL", supplierID).
		Scan(&status).Error
	if err != nil {
		return err
	}
	if status != 1 {
		return ecode.New(10007, "供应商不存在或已停用")
	}
	return nil
}

// checkDrugActive 校验药品是否存在且处于启用状态（status=1）。
func checkDrugActive(ctx context.Context, tx *gorm.DB, drugID int64) error {
	var status int8
	err := tx.WithContext(ctx).
		Table("drug_info").
		Select("status").
		Where("id = ? AND deleted_at IS NULL", drugID).
		Scan(&status).Error
	if err != nil {
		return err
	}
	if status != 1 {
		return ecode.New(10008, "药品不存在或已停用")
	}
	return nil
}

// normalizeTraceCodes 清洗并去重追溯码。
func normalizeTraceCodes(codes []string) ([]string, error) {
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

// publishMQAsync 异步发送 MQ 操作日志，失败仅记录日志不影响主流程。
func (s *inboundService) publishMQAsync(event mq.LogEvent) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("MQ 发送 panic", zap.Any("recover", r))
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := s.mqClient.PublishLogEvent(ctx, event); err != nil {
			s.logger.Warn("MQ 发送失败", zap.String("action", event.Action), zap.Error(err))
		}
	}()
}

// =====================================================
// CreateOrder
// =====================================================

// CreateOrder 创建入库单，可同时创建明细行。
func (s *inboundService) CreateOrder(ctx context.Context, req CreateOrderReq, creatorID int64) (*model.InboundOrder, error) {
	if req.SupplierID <= 0 {
		return nil, ecode.ErrParamInvalid
	}
	if creatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	var created *model.InboundOrder
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		created, err = s.CreateOrderWithDetails(ctx, tx, req, creatorID)
		return err
	})
	if err != nil {
		return nil, err
	}

	s.publishMQAsync(mq.LogEvent{
		BusinessType: "inbound_order",
		BusinessID:   fmt.Sprintf("%d", created.ID),
		Action:       "CREATE",
		OperatorID:   creatorID,
		Detail:       map[string]interface{}{"order_no": created.OrderNo},
	})
	return created, nil
}

// CreateOrderWithDetails 在给定事务中创建入库单及明细，供内部（如 AI 发票转换）调用。
func (s *inboundService) CreateOrderWithDetails(ctx context.Context, tx *gorm.DB, req CreateOrderReq, creatorID int64) (*model.InboundOrder, error) {
	// 校验供应商状态。
	if err := checkSupplierActive(ctx, tx, req.SupplierID); err != nil {
		return nil, err
	}

	// 生成单号。
	orderNo, err := s.repo.GenerateOrderNo(ctx, tx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	order := &model.InboundOrder{
		OrderNo:    orderNo,
		SupplierID: req.SupplierID,
		CreatorID:  creatorID,
		OperatorID: creatorID,
		Status:     model.InboundOrderStatusDraft,
	}
	if req.InvoiceNo != "" {
		order.InvoiceNo = &req.InvoiceNo
	}
	if req.Remark != "" {
		order.Remark = &req.Remark
	}
	_ = now // suppress unused warning

	if err := s.repo.Create(ctx, tx, order); err != nil {
		return nil, err
	}

	// 批量创建明细行。
	totalAmount := 0.0
	for _, d := range req.Details {
		if err := validateCreateDetail(d); err != nil {
			return nil, err
		}
		if err := checkDrugActive(ctx, tx, d.DrugID); err != nil {
			return nil, err
		}
		expDate, err := parseExpireDate(d.ExpireDate)
		if err != nil {
			return nil, ecode.ErrParamInvalid
		}
		amount := float64(d.PlannedQty) * d.UnitPrice
		detail := &model.InboundOrderDetail{
			OrderID:     order.ID,
			DrugID:      d.DrugID,
			BatchNumber: d.BatchNumber,
			ExpireDate:  expDate,
			PlannedQty:  d.PlannedQty,
			UnitPrice:   d.UnitPrice,
			Amount:      amount,
		}
		if d.Remark != "" {
			detail.Remark = &d.Remark
		}
		if err := s.repo.AddDetail(ctx, tx, detail); err != nil {
			return nil, err
		}
		totalAmount += amount
	}

	// 写入 total_amount。
	if len(req.Details) > 0 {
		if err := s.repo.UpdateTotalAmount(ctx, tx, order.ID); err != nil {
			return nil, err
		}
		order.TotalAmount = totalAmount
	}

	return order, nil
}

// =====================================================
// SubmitOrder
// =====================================================

// SubmitOrder 提交入库单：DRAFT → PENDING_CONFIRM，须至少有一条明细。
func (s *inboundService) SubmitOrder(ctx context.Context, orderID, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusDraft {
			return ecode.ErrStatusInvalid
		}

		// 校验至少有一条未删除明细。
		details, err := s.repo.GetDetails(ctx, orderID)
		if err != nil {
			return err
		}
		if len(details) == 0 {
			return ecode.New(20010, "入库单必须至少有一条明细才能提交")
		}

		if err := s.repo.UpdateStatus(ctx, tx, orderID, model.InboundOrderStatusPendingConfirm); err != nil {
			return err
		}

		s.publishMQAsync(mq.LogEvent{
			BusinessType: "inbound_order",
			BusinessID:   fmt.Sprintf("%d", orderID),
			Action:       "SUBMIT",
			OperatorID:   operatorID,
			Detail:       map[string]interface{}{"order_no": order.OrderNo},
		})
		return nil
	})
}

// =====================================================
// GetOrder / ListOrders
// =====================================================

// GetOrder 查询入库单详情，附带明细行。
func (s *inboundService) GetOrder(ctx context.Context, orderID int64) (*model.InboundOrder, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	// 预加载明细。
	details, err := s.repo.GetDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, d := range details {
		order.Details = append(order.Details, *d)
	}
	// 填充供应商名称和创建人名称。
	s.enrichOrders(ctx, []*model.InboundOrder{order})
	// 填充明细行的药品信息。
	s.enrichDetails(ctx, order.Details)
	return order, nil
}

// enrichDetails 批量填充入库明细行的药品名称、规格和生产厂家。
func (s *inboundService) enrichDetails(ctx context.Context, details []model.InboundOrderDetail) {
	if len(details) == 0 {
		return
	}
	drugIDs := make([]int64, 0, len(details))
	for _, d := range details {
		drugIDs = append(drugIDs, d.DrugID)
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
	for i := range details {
		if d, ok := drugMap[details[i].DrugID]; ok {
			details[i].DrugName = d.CommonName
			details[i].Specification = d.Specification
			details[i].Manufacturer = d.Manufacturer
		}
	}
}

// ListOrders 分页查询入库单列表。
func (s *inboundService) ListOrders(ctx context.Context, filter repository.ListFilter) ([]*model.InboundOrder, int64, error) {
	orders, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	s.enrichOrders(ctx, orders)
	return orders, total, nil
}

// enrichOrders 批量填充入库单的供应商名称和创建人名称。
func (s *inboundService) enrichOrders(ctx context.Context, orders []*model.InboundOrder) {
	if len(orders) == 0 {
		return
	}
	supplierIDs := make([]int64, 0, len(orders))
	creatorIDs := make([]int64, 0, len(orders))
	for _, o := range orders {
		supplierIDs = append(supplierIDs, o.SupplierID)
		creatorIDs = append(creatorIDs, o.CreatorID)
	}

	type supplierRow struct {
		ID   int64
		Name string
	}
	var suppliers []supplierRow
	s.db.WithContext(ctx).Table("supplier").Select("id, name").Where("id IN ?", supplierIDs).Scan(&suppliers)
	supplierMap := make(map[int64]string, len(suppliers))
	for _, sup := range suppliers {
		supplierMap[sup.ID] = sup.Name
	}

	type userRow struct {
		ID       int64
		RealName string
	}
	var users []userRow
	s.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", creatorIDs).Scan(&users)
	userMap := make(map[int64]string, len(users))
	for _, u := range users {
		userMap[u.ID] = u.RealName
	}

	// Aggregate planned/confirmed quantities per order in one query.
	orderIDs := make([]int64, 0, len(orders))
	for _, o := range orders {
		orderIDs = append(orderIDs, o.ID)
	}
	type progressRow struct {
		OrderID        int64
		TotalPlanned   int32
		TotalConfirmed int32
	}
	var progresses []progressRow
	s.db.WithContext(ctx).
		Table("inbound_order_detail").
		Select("order_id, COALESCE(SUM(planned_qty),0) AS total_planned, COALESCE(SUM(confirmed_qty),0) AS total_confirmed").
		Where("order_id IN ?", orderIDs).
		Group("order_id").
		Scan(&progresses)
	progressMap := make(map[int64]progressRow, len(progresses))
	for _, p := range progresses {
		progressMap[p.OrderID] = p
	}

	for _, o := range orders {
		o.SupplierName = supplierMap[o.SupplierID]
		o.CreatorName = userMap[o.CreatorID]
		if p, ok := progressMap[o.ID]; ok {
			o.TotalPlannedQty = p.TotalPlanned
			o.TotalConfirmedQty = p.TotalConfirmed
		}
	}
}

// =====================================================
// UpdateOrder
// =====================================================

// UpdateOrder 更新草稿入库单基础信息（仅 DRAFT 状态）。
func (s *inboundService) UpdateOrder(ctx context.Context, orderID int64, req UpdateOrderReq, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusDraft {
			return ecode.ErrStatusInvalid
		}

		if req.SupplierID > 0 && req.SupplierID != order.SupplierID {
			if err := checkSupplierActive(ctx, tx, req.SupplierID); err != nil {
				return err
			}
			order.SupplierID = req.SupplierID
		}
		if req.InvoiceNo != "" {
			order.InvoiceNo = &req.InvoiceNo
		}
		if req.Remark != "" {
			order.Remark = &req.Remark
		}
		order.OperatorID = operatorID

		return s.repo.Update(ctx, tx, order)
	})
}

// =====================================================
// 明细 CRUD
// =====================================================

// AddDetail 新增明细行（仅 DRAFT 状态），返回已创建的明细。
func (s *inboundService) AddDetail(ctx context.Context, orderID int64, req CreateDetailReq, operatorID int64) (*model.InboundOrderDetail, error) {
	if err := validateCreateDetail(req); err != nil {
		return nil, err
	}
	expDate, err := parseExpireDate(req.ExpireDate)
	if err != nil {
		return nil, ecode.ErrParamInvalid
	}

	var result *model.InboundOrderDetail
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusDraft {
			return ecode.ErrStatusInvalid
		}
		if err := checkDrugActive(ctx, tx, req.DrugID); err != nil {
			return err
		}

		amount := float64(req.PlannedQty) * req.UnitPrice
		detail := &model.InboundOrderDetail{
			OrderID:     orderID,
			DrugID:      req.DrugID,
			BatchNumber: req.BatchNumber,
			ExpireDate:  expDate,
			PlannedQty:  req.PlannedQty,
			UnitPrice:   req.UnitPrice,
			Amount:      amount,
		}
		if req.Remark != "" {
			detail.Remark = &req.Remark
		}
		if err := s.repo.AddDetail(ctx, tx, detail); err != nil {
			return err
		}
		result = detail
		return s.repo.UpdateTotalAmount(ctx, tx, orderID)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateDetail 更新明细行（仅 DRAFT 状态），返回更新后的明细。
func (s *inboundService) UpdateDetail(ctx context.Context, orderID, detailID int64, req UpdateDetailReq, operatorID int64) (*model.InboundOrderDetail, error) {
	expDate, err := parseExpireDate(req.ExpireDate)
	if err != nil {
		return nil, ecode.ErrParamInvalid
	}

	var result *model.InboundOrderDetail
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusDraft {
			return ecode.ErrStatusInvalid
		}

		detail, err := s.repo.GetDetailByIDForUpdate(ctx, tx, orderID, detailID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		if req.PlannedQty > 0 {
			detail.PlannedQty = req.PlannedQty
		}
		if req.UnitPrice >= 0 {
			detail.UnitPrice = req.UnitPrice
		}
		if req.BatchNumber != "" {
			detail.BatchNumber = req.BatchNumber
		}
		detail.ExpireDate = expDate
		detail.Amount = float64(detail.PlannedQty) * detail.UnitPrice
		if req.Remark != "" {
			detail.Remark = &req.Remark
		}

		if err := s.repo.UpdateDetail(ctx, tx, detail); err != nil {
			return err
		}
		result = detail
		return s.repo.UpdateTotalAmount(ctx, tx, orderID)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteDetail 软删除明细行（仅 DRAFT 状态）。
func (s *inboundService) DeleteDetail(ctx context.Context, orderID, detailID int64, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusDraft {
			return ecode.ErrStatusInvalid
		}

		if err := s.repo.DeleteDetail(ctx, tx, orderID, detailID); err != nil {
			return err
		}
		return s.repo.UpdateTotalAmount(ctx, tx, orderID)
	})
}

// =====================================================
// ConfirmTraceCodes
// =====================================================

// ConfirmTraceCodes 扫码录入追溯码（PENDING_CONFIRM 状态，批量）。
func (s *inboundService) ConfirmTraceCodes(ctx context.Context, req ConfirmTraceCodesReq) error {
	if req.OrderID <= 0 || req.DetailID <= 0 || req.OperatorID <= 0 {
		return ecode.ErrParamInvalid
	}
	if len(req.TraceCodes) == 0 {
		return ecode.ErrParamInvalid
	}

	// 清洗并去重追溯码。
	cleanCodes, err := normalizeTraceCodes(req.TraceCodes)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 加锁读取入库单，保证并发安全。
		order, err := s.repo.GetByIDForUpdate(ctx, tx, req.OrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusPendingConfirm {
			return ecode.ErrStatusInvalid
		}

		// 加锁读取明细，校验 detail_id 归属。
		detail, err := s.repo.GetDetailByIDForUpdate(ctx, tx, req.OrderID, req.DetailID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		// 校验数量上限：confirmed + delta <= planned。
		if detail.ConfirmedQty+int32(len(cleanCodes)) > detail.PlannedQty {
			return ecode.New(20006, "确认数量超过计划数量")
		}

		// 逐一校验追溯码全局唯一性（避免重复入库）。
		for _, code := range cleanCodes {
			exists, err := s.repo.ExistsTraceCode(ctx, tx, code)
			if err != nil {
				return err
			}
			if exists {
				return ecode.ErrDuplicateScan
			}
		}

		// 批量写入追溯库存（初始状态：PENDING）。
		inventories := make([]*model.DrugTraceInventory, 0, len(cleanCodes))
		for _, code := range cleanCodes {
			inventories = append(inventories, &model.DrugTraceInventory{
				TraceCode:       code,
				DrugID:          detail.DrugID,
				BatchNumber:     detail.BatchNumber,
				ExpireDate:      detail.ExpireDate,
				InboundOrderID:  order.ID,
				InboundDetailID: detail.ID,
				Status:          model.TraceInventoryStatusPending,
			})
		}
		if err := s.repo.CreateTraceInventory(ctx, tx, inventories); err != nil {
			return err
		}

		// 原子累加已确认数量。
		if err := s.repo.IncreaseConfirmedQty(ctx, tx, detail.ID, int32(len(cleanCodes))); err != nil {
			return err
		}

		// 写入追溯日志（在同一事务内）。
		logs := make([]*model.DrugTraceLog, 0, len(cleanCodes))
		for _, code := range cleanCodes {
			fromStatus := ""
			toStatus := model.TraceInventoryStatusPending
			logs = append(logs, &model.DrugTraceLog{
				TraceCode:   code,
				ActionType:  model.TraceLogActionInbound,
				FromStatus:  &fromStatus,
				ToStatus:    toStatus,
				OperatorID:  req.OperatorID,
				RelatedNo:   &order.OrderNo,
				DrugID:      detail.DrugID,
				OrderID:     &order.ID,
				OrderItemID: &detail.ID,
			})
		}
		if err := s.repo.CreateTraceLog(ctx, tx, logs); err != nil {
			return err
		}

		return nil
	})
}

// =====================================================
// CompleteOrder
// =====================================================

// CompleteOrder 完成入库单：PENDING_CONFIRM → COMPLETED（须所有明细均已全量确认）。
func (s *inboundService) CompleteOrder(ctx context.Context, orderID, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.InboundOrderStatusPendingConfirm {
			return ecode.ErrStatusInvalid
		}

		// 校验所有明细均已全量确认。
		fullyConfirmed, err := s.repo.IsOrderFullyConfirmed(ctx, tx, orderID)
		if err != nil {
			return err
		}
		if !fullyConfirmed {
			return ecode.New(20007, "入库单存在未完全确认的明细，无法完成")
		}

		if err := s.repo.UpdateStatus(ctx, tx, orderID, model.InboundOrderStatusCompleted); err != nil {
			return err
		}

		s.publishMQAsync(mq.LogEvent{
			BusinessType: "inbound_order",
			BusinessID:   fmt.Sprintf("%d", orderID),
			Action:       "COMPLETE",
			OperatorID:   operatorID,
			Detail:       map[string]interface{}{"order_no": order.OrderNo},
		})
		return nil
	})
}

// =====================================================
// CancelOrder
// =====================================================

// CancelOrder 取消入库单。
// DRAFT → CANCELLED 直接取消；
// PENDING_CONFIRM → CANCELLED 须确认无 IN_STOCK 追溯码，并软删除 PENDING 追溯库存。
func (s *inboundService) CancelOrder(ctx context.Context, orderID, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}

		switch order.Status {
		case model.InboundOrderStatusDraft:
			// 草稿直接取消，无需其他处理。
		case model.InboundOrderStatusPendingConfirm:
			// 检查是否存在已上架（IN_STOCK）的追溯库存，如有则不允许取消。
			var inStockCount int64
			if err := tx.WithContext(ctx).
				Model(&model.DrugTraceInventory{}).
				Where("inbound_order_id = ? AND status = ?", orderID, model.TraceInventoryStatusInStock).
				Count(&inStockCount).Error; err != nil {
				return err
			}
			if inStockCount > 0 {
				return ecode.New(20008, "入库单存在已上架追溯码，无法取消")
			}
			// 软删除所有 PENDING 追溯库存。
			if err := s.repo.SoftDeletePendingTraceInventory(ctx, tx, orderID); err != nil {
				return err
			}
		default:
			return ecode.ErrStatusInvalid
		}

		if err := s.repo.UpdateStatus(ctx, tx, orderID, model.InboundOrderStatusCancelled); err != nil {
			return err
		}

		s.publishMQAsync(mq.LogEvent{
			BusinessType: "inbound_order",
			BusinessID:   fmt.Sprintf("%d", orderID),
			Action:       "CANCEL",
			OperatorID:   operatorID,
			Detail:       map[string]interface{}{"order_no": order.OrderNo},
		})
		return nil
	})
}

// =====================================================
// GetInboundProgress
// =====================================================

// GetInboundProgress 查询入库确认进度 DTO。
func (s *inboundService) GetInboundProgress(ctx context.Context, orderID int64) (*InboundProgressDTO, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	details, err := s.repo.GetDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var totalPlanned, totalConfirmed int32
	progressDetails := make([]ProgressDetail, 0, len(details))
	for _, d := range details {
		totalPlanned += d.PlannedQty
		totalConfirmed += d.ConfirmedQty

		// 查询药品名称（common_name）。
		var drugName string
		_ = s.db.WithContext(ctx).
			Table("drug_info").
			Select("common_name").
			Where("id = ?", d.DrugID).
			Scan(&drugName).Error

		progressDetails = append(progressDetails, ProgressDetail{
			DetailID:     d.ID,
			DrugID:       d.DrugID,
			DrugName:     drugName,
			BatchNumber:  d.BatchNumber,
			PlannedQty:   d.PlannedQty,
			ConfirmedQty: d.ConfirmedQty,
		})
	}

	return &InboundProgressDTO{
		OrderID:        order.ID,
		OrderNo:        order.OrderNo,
		Status:         order.Status,
		TotalPlanned:   totalPlanned,
		TotalConfirmed: totalConfirmed,
		Details:        progressDetails,
	}, nil
}
