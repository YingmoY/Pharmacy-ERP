package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	drugModel "github.com/YingmoY/PharmacyERP/internal/drug/model"
	inventoryModel "github.com/YingmoY/PharmacyERP/internal/inventory/model"
	pharmacistModel "github.com/YingmoY/PharmacyERP/internal/pharmacist/model"
	pharmacistRepo "github.com/YingmoY/PharmacyERP/internal/pharmacist/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/medicare"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/sales/model"
	"github.com/YingmoY/PharmacyERP/internal/sales/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ==================== 请求/响应 DTO ====================

// CreateOrderItemReq 创建订单时单个药品项的请求
type CreateOrderItemReq struct {
	DrugID    int64   `json:"drug_id" binding:"required"`
	TraceCode *string `json:"trace_code"` // 可选，不填则自动选择
}

// CreateOrderReq 创建销售订单请求
type CreateOrderReq struct {
	CustomerName   *string              `json:"customer_name"`
	IsPrescription bool                 `json:"is_prescription"`
	Remark         *string              `json:"remark"`
	Items          []CreateOrderItemReq `json:"items" binding:"required,min=1"`
}

// PayReq 结算请求
type PayReq struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
	// UseMedicare lets the customer opt in to medicare settlement at checkout.
	UseMedicare  bool   `json:"use_medicare"`
	// Medicare-specific fields; only required when UseMedicare = true.
	MedType      string `json:"med_type"`       // e.g. "41" (门诊)
	InsuType     string `json:"insutype"`       // e.g. "310" (职工医保)
	AcctUsedFlag string `json:"acct_used_flag"` // "1" to use balance account
	// Patient identity fields — required when UseMedicare = true.
	MdtrtCertNo string `json:"mdtrt_cert_no"` // 医保凭证号（身份证号）
	PsnNo       string `json:"psn_no"`        // 参保人员编号（由1101查询返回）
	PsnName     string `json:"psn_name"`      // 参保人姓名
}

// MedicarePreviewReq 医保预查询请求
type MedicarePreviewReq struct {
	MdtrtCertNo  string `json:"mdtrt_cert_no" binding:"required"`
	MedType      string `json:"med_type"`
	InsuType     string `json:"insutype"`
	AcctUsedFlag string `json:"acct_used_flag"`
}

// MedicarePreviewResp 医保预查询结果（含费用分摊明细）
type MedicarePreviewResp struct {
	PsnNo        string  `json:"psn_no"`
	PsnName      string  `json:"psn_name"`
	TotalAmount  float64 `json:"total_amount"`
	FundPay      float64 `json:"fund_pay"`      // 统筹基金支付
	AcctPay      float64 `json:"acct_pay"`      // 个人账户支付
	PersonalCash float64 `json:"personal_cash"` // 个人现金支付
}

// RefundReq 退款请求
type RefundReq struct {
	RefundMode   string  `json:"refund_mode" binding:"required"` // FULL 或 PARTIAL
	RefundReason string  `json:"refund_reason"`
	DetailIDs    []int64 `json:"detail_ids"` // PARTIAL 模式下指定明细 ID
}

// OrderFilter 订单列表过滤条件
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

// ReserveTraceReq 手动锁定追溯码请求
type ReserveTraceReq struct {
	TraceCode string `json:"trace_code" binding:"required"`
	DrugID    int64  `json:"drug_id" binding:"required"`
}

// ReleaseReservationReq 手动释放预留请求
type ReleaseReservationReq struct {
	TraceCode string `json:"trace_code" binding:"required"`
}

// ScanVerifyReq 扫码验证请求
type ScanVerifyReq struct {
	TraceCode string `json:"trace_code" binding:"required"`
	DrugID    *int64 `json:"drug_id"` // 可选，如果提供则验证药品归属
}

// ScanVerifyResult 扫码验证结果
type ScanVerifyResult struct {
	TraceCode   string `json:"trace_code"`
	IsValid     bool   `json:"is_valid"`
	Status      string `json:"status"`
	DrugID      int64  `json:"drug_id"`
	DrugName    string `json:"drug_name"`
	BatchNumber string `json:"batch_number"`
	Reason      string `json:"reason,omitempty"`
}

// ==================== Service 接口 ====================

// SalesService 定义销售核心业务能力
type SalesService interface {
	// CreateOrder 创建销售订单（含预留追溯码）
	CreateOrder(ctx context.Context, req CreateOrderReq, cashierID int64) (*model.SalesOrder, error)
	// GetOrder 根据 ID 查询销售订单
	GetOrder(ctx context.Context, id int64) (*model.SalesOrder, error)
	// ListOrders 分页查询销售订单
	ListOrders(ctx context.Context, filter OrderFilter) ([]*model.SalesOrder, int64, error)
	// AddItem 向 PENDING 订单添加明细
	AddItem(ctx context.Context, orderID int64, req CreateOrderItemReq, operatorID int64) (*model.SalesOrderItem, error)
	// DeleteItem 从 PENDING 订单删除明细并释放预留
	DeleteItem(ctx context.Context, orderID, itemID, operatorID int64) error
	// GetItems 查询订单明细列表
	GetItems(ctx context.Context, orderID int64) ([]*model.SalesOrderItem, error)
	// Pay 结算支��
	Pay(ctx context.Context, orderID int64, req PayReq, operatorID int64) error
	// MedicarePreview 查询参保人信息并预估医保费用分摊
	MedicarePreview(ctx context.Context, orderID int64, req MedicarePreviewReq) (*MedicarePreviewResp, error)
	// Cancel 取消订单
	Cancel(ctx context.Context, orderID, operatorID int64) error
	// Refund 退款处理
	Refund(ctx context.Context, orderID int64, req RefundReq, operatorID int64) error
	// GetReservedTraces 查询订单预留的追溯码列表
	GetReservedTraces(ctx context.Context, orderID int64) ([]*model.TraceReservation, error)
	// GetReviewRecord 查询订单审核记录
	GetReviewRecord(ctx context.Context, orderID int64) (*pharmacistModel.AuditReview, error)
	// SubmitReview 手动提交审核（如需）
	SubmitReview(ctx context.Context, orderID, submitterID int64) error
	// ReserveTrace 手动为订单锁定追溯码
	ReserveTrace(ctx context.Context, orderID int64, req ReserveTraceReq, operatorID int64) (*model.TraceReservation, error)
	// ReleaseReservation 手动释放预留
	ReleaseReservation(ctx context.Context, orderID int64, req ReleaseReservationReq, operatorID int64) error
	// ScanVerify 扫码验证追溯码是否可售
	ScanVerify(ctx context.Context, req ScanVerifyReq) (*ScanVerifyResult, error)
}

// ==================== Service 实现 ====================

type salesService struct {
	db             *gorm.DB
	salesRepo      repository.SalesRepo
	reviewRepo     pharmacistRepo.ReviewRepo
	mqClient       *mq.Client
	medicareClient *medicare.Client
	logger         *zap.Logger
}

// NewSalesService 创建销售服务实例
func NewSalesService(
	db *gorm.DB,
	salesRepo repository.SalesRepo,
	reviewRepo pharmacistRepo.ReviewRepo,
	mqClient *mq.Client,
	medicareClient *medicare.Client,
	logger *zap.Logger,
) SalesService {
	return &salesService{
		db:             db,
		salesRepo:      salesRepo,
		reviewRepo:     reviewRepo,
		mqClient:       mqClient,
		medicareClient: medicareClient,
		logger:         logger,
	}
}

// ==================== CreateOrder ====================

func (s *salesService) CreateOrder(ctx context.Context, req CreateOrderReq, cashierID int64) (*model.SalesOrder, error) {
	if cashierID <= 0 {
		return nil, ecode.ErrUnauthorized
	}
	if len(req.Items) == 0 {
		return nil, ecode.ErrParamInvalid
	}

	var createdOrder *model.SalesOrder
	var reservations []*model.TraceReservation

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 生成订单号
		orderNo, err := s.salesRepo.GenOrderNo(ctx, tx)
		if err != nil {
			return err
		}

		// 2. 遍历请求项，验证药品并解析追溯码
		type itemData struct {
			drug      drugModel.DrugInfo
			trace     inventoryModel.TraceInventory
			itemReq   CreateOrderItemReq
		}

		items := make([]itemData, 0, len(req.Items))
		needAudit := req.IsPrescription

		for _, itemReq := range req.Items {
			// 查询药品信息
			var drug drugModel.DrugInfo
			if err := tx.WithContext(ctx).Where("id = ? AND status = ?", itemReq.DrugID, drugModel.DrugStatusEnabled).First(&drug).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ecode.New(ecode.ErrNotFound.Code, "drug not found or disabled")
				}
				return err
			}
			if drug.RetailPrice == nil {
				return ecode.New(ecode.ErrParamInvalid.Code, "drug has no retail price")
			}
			if drug.IsPrescription {
				needAudit = true
			}

			// 解析/自动选取追溯码
			var trace inventoryModel.TraceInventory
			if itemReq.TraceCode != nil && *itemReq.TraceCode != "" {
				// 指定了追溯码：验证存在且 IN_STOCK
				if err := tx.WithContext(ctx).
					Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
					Where("trace_code = ? AND drug_id = ? AND status = ?",
						*itemReq.TraceCode, itemReq.DrugID, inventoryModel.TraceInventoryStatusInStock).
					First(&trace).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return ecode.New(ecode.ErrTraceCodeNotFound.Code, "trace code not found, not IN_STOCK, or belongs to different drug")
					}
					return err
				}
				// 检查是否已被预留
				var rsvCount int64
				if err := tx.WithContext(ctx).Model(&model.TraceReservation{}).
					Where("trace_code = ? AND status = ?", trace.TraceCode, model.ReservationStatusReserved).
					Count(&rsvCount).Error; err != nil {
					return err
				}
				if rsvCount > 0 {
					return ecode.New(ecode.ErrConflict.Code, "trace code is already reserved")
				}
			} else {
				// 自动选取：按过期日期 ASC（FIFO 近效期优先），未被预留的 IN_STOCK 追溯码
				subQuery := tx.WithContext(ctx).
					Model(&model.TraceReservation{}).
					Select("trace_code").
					Where("status = ?", model.ReservationStatusReserved)

				if err := tx.WithContext(ctx).
					Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
					Where("drug_id = ? AND status = ? AND trace_code NOT IN (?)",
						itemReq.DrugID, inventoryModel.TraceInventoryStatusInStock, subQuery).
					Order("expire_date ASC").
					Limit(1).
					First(&trace).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return ecode.New(ecode.ErrTraceCodeNotFound.Code, "no available trace code for drug")
					}
					return err
				}
			}

			items = append(items, itemData{drug: drug, trace: trace, itemReq: itemReq})
		}

		// 3. 确定订单状态
		status := model.SalesOrderStatusPending
		if needAudit {
			status = model.SalesOrderStatusPendingReview
		}

		// 4. 计算总金额
		var totalAmount float64
		for _, item := range items {
			totalAmount += *item.drug.RetailPrice
		}

		// 5. 创建销售订单
		now := time.Now()
		order := &model.SalesOrder{
			OrderNo:        orderNo,
			CashierID:      cashierID,
			TotalAmount:    totalAmount,
			ActualAmount:   totalAmount,
			NeedAudit:      needAudit,
			Status:         status,
			CustomerName:   req.CustomerName,
			IsPrescription: req.IsPrescription,
			Remark:         req.Remark,
		}
		if err := s.salesRepo.CreateOrder(ctx, tx, order); err != nil {
			return err
		}

		// 6. 创建订单明细 + 追溯预留
		for _, item := range items {
			price := *item.drug.RetailPrice
			orderItem := &model.SalesOrderItem{
				OrderID:        order.ID,
				DrugID:         item.drug.ID,
				TraceCode:      item.trace.TraceCode,
				Price:          price,
				Quantity:       1,
				SubtotalAmount: price,
				RefundStatus:   model.ItemRefundStatusNone,
			}
			if err := s.salesRepo.CreateItem(ctx, tx, orderItem); err != nil {
				return err
			}

			// 生成预留单号
			rsvNo, err := s.salesRepo.GenReservationNo(ctx, tx)
			if err != nil {
				return err
			}

			expireAt := now.Add(30 * time.Minute)
			rsv := &model.TraceReservation{
				ReservationNo:    rsvNo,
				SalesOrderID:     order.ID,
				SalesOrderItemID: &orderItem.ID,
				TraceCode:        item.trace.TraceCode,
				DrugID:           item.drug.ID,
				ReservedBy:       cashierID,
				Status:           model.ReservationStatusReserved,
				ReservedAt:       now,
				ExpireAt:         expireAt,
			}
			if err := s.salesRepo.CreateReservation(ctx, tx, rsv); err != nil {
				return err
			}
			reservations = append(reservations, rsv)
		}

		// 7. 如果需要审核，创建审核记录
		if needAudit {
			reviewNo, err := s.reviewRepo.GenReviewNo(ctx, tx)
			if err != nil {
				return err
			}
			submittedAt := now
			review := &pharmacistModel.AuditReview{
				OrderID:     order.ID,
				ReviewNo:    reviewNo,
				Status:      pharmacistModel.AuditReviewStatusPending,
				SubmitterID: &cashierID,
				SubmittedAt: &submittedAt,
			}
			if err := s.reviewRepo.CreateReview(ctx, tx, review); err != nil {
				return err
			}
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload order with items so the response includes the created line items.
	if fullOrder, loadErr := s.salesRepo.GetOrderByID(ctx, s.db, createdOrder.ID); loadErr == nil {
		createdOrder = fullOrder
		s.enrichOrder(ctx, createdOrder)
	}

	// 8. 事务提交后异步发布预留过期消息
	for _, rsv := range reservations {
		event := mq.ReservationExpireEvent{
			ReservationID: rsv.ID,
			ReservationNo: rsv.ReservationNo,
			SalesOrderID:  rsv.SalesOrderID,
			TraceCode:     rsv.TraceCode,
			ExpireAt:      rsv.ExpireAt,
		}
		if pubErr := s.mqClient.PublishReservationExpire(ctx, event); pubErr != nil {
			s.logger.Warn("发布预留过期消息失败",
				zap.String("reservation_no", rsv.ReservationNo),
				zap.Error(pubErr),
			)
		}
	}

	return createdOrder, nil
}

// ==================== GetOrder ====================

func (s *salesService) GetOrder(ctx context.Context, id int64) (*model.SalesOrder, error) {
	order, err := s.salesRepo.GetOrderByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	s.enrichOrder(ctx, order)
	return order, nil
}

// enrichOrder populates virtual fields on a SalesOrder and its Items.
func (s *salesService) enrichOrder(ctx context.Context, order *model.SalesOrder) {
	if order == nil {
		return
	}
	if order.CashierID > 0 {
		var cashier struct{ RealName string }
		s.db.WithContext(ctx).Table("sys_user").Select("real_name").Where("id = ?", order.CashierID).Scan(&cashier)
		order.CashierName = cashier.RealName
	}
	if len(order.Items) == 0 {
		return
	}

	drugIDs := make([]int64, 0, len(order.Items))
	traceCodes := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		drugIDs = append(drugIDs, item.DrugID)
		if item.TraceCode != "" {
			traceCodes = append(traceCodes, item.TraceCode)
		}
	}

	type drugRow struct {
		ID            int64
		CommonName    string
		Specification string
		Manufacturer  string
	}
	var drugs []drugRow
	s.db.WithContext(ctx).Table("drug_info").
		Select("id, common_name, specification, manufacturer").
		Where("id IN ?", drugIDs).Scan(&drugs)
	drugMap := make(map[int64]drugRow, len(drugs))
	for _, d := range drugs {
		drugMap[d.ID] = d
	}

	type traceRow struct {
		TraceCode   string
		BatchNumber string
		ExpireDate  string
		LocationID  *int64
	}
	var traces []traceRow
	if len(traceCodes) > 0 {
		s.db.WithContext(ctx).Table("drug_trace_inventory").
			Select("trace_code, batch_number, TO_CHAR(expire_date, 'YYYY-MM-DD') AS expire_date, location_id").
			Where("trace_code IN ?", traceCodes).Scan(&traces)
	}
	traceMap := make(map[string]traceRow, len(traces))
	locationIDs := make([]int64, 0)
	for _, t := range traces {
		traceMap[t.TraceCode] = t
		if t.LocationID != nil {
			locationIDs = append(locationIDs, *t.LocationID)
		}
	}

	locMap := make(map[int64]string)
	if len(locationIDs) > 0 {
		type locRow struct {
			ID           int64
			LocationCode string
		}
		var locs []locRow
		s.db.WithContext(ctx).Table("location_info").
			Select("id, location_code").
			Where("id IN ?", locationIDs).Scan(&locs)
		for _, l := range locs {
			locMap[l.ID] = l.LocationCode
		}
	}

	for i := range order.Items {
		item := &order.Items[i]
		item.UnitPrice = item.Price
		if d, ok := drugMap[item.DrugID]; ok {
			item.DrugName = d.CommonName
			item.Specification = d.Specification
			item.Manufacturer = d.Manufacturer
		}
		if t, ok := traceMap[item.TraceCode]; ok {
			item.BatchNumber = t.BatchNumber
			item.ExpireDate = t.ExpireDate
			if t.LocationID != nil {
				item.LocationCode = locMap[*t.LocationID]
			}
		}
	}
}

// ==================== ListOrders ====================

func (s *salesService) ListOrders(ctx context.Context, filter OrderFilter) ([]*model.SalesOrder, int64, error) {
	orders, total, err := s.salesRepo.ListOrders(ctx, s.db, repository.OrderFilter{
		CashierID:      filter.CashierID,
		Status:         filter.Status,
		OrderNo:        filter.OrderNo,
		StartDate:      filter.StartDate,
		EndDate:        filter.EndDate,
		IsPrescription: filter.IsPrescription,
		Page:           filter.Page,
		PageSize:       filter.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	// Batch-fetch cashier names for the list view
	if len(orders) > 0 {
		cashierIDs := make([]int64, 0, len(orders))
		for _, o := range orders {
			cashierIDs = append(cashierIDs, o.CashierID)
		}
		type userRow struct {
			ID       int64
			RealName string
		}
		var users []userRow
		s.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", cashierIDs).Scan(&users)
		nameMap := make(map[int64]string, len(users))
		for _, u := range users {
			nameMap[u.ID] = u.RealName
		}
		for _, o := range orders {
			o.CashierName = nameMap[o.CashierID]
		}
	}
	return orders, total, nil
}

// ==================== AddItem ====================

func (s *salesService) AddItem(ctx context.Context, orderID int64, req CreateOrderItemReq, operatorID int64) (*model.SalesOrderItem, error) {
	var createdItem *model.SalesOrderItem

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询订单
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		// 只有 PENDING 状态可以添加明细
		if order.Status != model.SalesOrderStatusPending {
			return ecode.New(ecode.ErrStatusInvalid.Code, "can only add items to PENDING order")
		}

		// 2. 查询药品
		var drug drugModel.DrugInfo
		if err := tx.WithContext(ctx).Where("id = ? AND status = ?", req.DrugID, drugModel.DrugStatusEnabled).First(&drug).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.New(ecode.ErrNotFound.Code, "drug not found or disabled")
			}
			return err
		}
		if drug.RetailPrice == nil {
			return ecode.New(ecode.ErrParamInvalid.Code, "drug has no retail price")
		}

		// 3. 选取追溯码
		var trace inventoryModel.TraceInventory
		if req.TraceCode != nil && *req.TraceCode != "" {
			if err := tx.WithContext(ctx).
				Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Where("trace_code = ? AND drug_id = ? AND status = ?",
					*req.TraceCode, req.DrugID, inventoryModel.TraceInventoryStatusInStock).
				First(&trace).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ecode.New(ecode.ErrTraceCodeNotFound.Code, "trace code not found or unavailable")
				}
				return err
			}
			var rsvCount int64
			if err := tx.WithContext(ctx).Model(&model.TraceReservation{}).
				Where("trace_code = ? AND status = ?", trace.TraceCode, model.ReservationStatusReserved).
				Count(&rsvCount).Error; err != nil {
				return err
			}
			if rsvCount > 0 {
				return ecode.New(ecode.ErrConflict.Code, "trace code is already reserved")
			}
		} else {
			subQuery := tx.WithContext(ctx).Model(&model.TraceReservation{}).
				Select("trace_code").
				Where("status = ?", model.ReservationStatusReserved)
			if err := tx.WithContext(ctx).
				Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Where("drug_id = ? AND status = ? AND trace_code NOT IN (?)",
					req.DrugID, inventoryModel.TraceInventoryStatusInStock, subQuery).
				Order("expire_date ASC").
				Limit(1).
				First(&trace).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ecode.New(ecode.ErrTraceCodeNotFound.Code, "no available trace code for drug")
				}
				return err
			}
		}

		price := *drug.RetailPrice
		orderItem := &model.SalesOrderItem{
			OrderID:        order.ID,
			DrugID:         drug.ID,
			TraceCode:      trace.TraceCode,
			Price:          price,
			Quantity:       1,
			SubtotalAmount: price,
			RefundStatus:   model.ItemRefundStatusNone,
		}
		if err := s.salesRepo.CreateItem(ctx, tx, orderItem); err != nil {
			return err
		}

		// 4. 创建预留
		rsvNo, err := s.salesRepo.GenReservationNo(ctx, tx)
		if err != nil {
			return err
		}
		now := time.Now()
		expireAt := now.Add(30 * time.Minute)
		rsv := &model.TraceReservation{
			ReservationNo:    rsvNo,
			SalesOrderID:     order.ID,
			SalesOrderItemID: &orderItem.ID,
			TraceCode:        trace.TraceCode,
			DrugID:           drug.ID,
			ReservedBy:       operatorID,
			Status:           model.ReservationStatusReserved,
			ReservedAt:       now,
			ExpireAt:         expireAt,
		}
		if err := s.salesRepo.CreateReservation(ctx, tx, rsv); err != nil {
			return err
		}

		// 5. 重新计算总金额
		newTotal, err := s.salesRepo.SumItemTotalByOrderID(ctx, tx, order.ID)
		if err != nil {
			return err
		}
		if err := s.salesRepo.UpdateOrderAmounts(ctx, tx, order.ID, map[string]interface{}{
			"total_amount":  newTotal,
			"actual_amount": newTotal - order.DiscountAmount,
		}); err != nil {
			return err
		}

		// 6. 如果新增了处方药，更新 need_audit 并切换状态为 PENDING_REVIEW
		if drug.IsPrescription && !order.NeedAudit {
			reviewNo, err := s.reviewRepo.GenReviewNo(ctx, tx)
			if err != nil {
				return err
			}
			submittedAt := now
			review := &pharmacistModel.AuditReview{
				OrderID:     order.ID,
				ReviewNo:    reviewNo,
				Status:      pharmacistModel.AuditReviewStatusPending,
				SubmitterID: &operatorID,
				SubmittedAt: &submittedAt,
			}
			if err := s.reviewRepo.CreateReview(ctx, tx, review); err != nil {
				return err
			}
			if err := s.salesRepo.UpdateOrderStatus(ctx, tx, order.ID, model.SalesOrderStatusPendingReview, map[string]interface{}{
				"need_audit": true,
			}); err != nil {
				return err
			}
		}

		// Populate virtual fields immediately from already-fetched data
		createdItem = orderItem
		createdItem.UnitPrice = *drug.RetailPrice
		createdItem.DrugName = drug.CommonName
		createdItem.Specification = drug.Specification
		createdItem.Manufacturer = drug.Manufacturer
		createdItem.BatchNumber = trace.BatchNumber
		createdItem.ExpireDate = trace.ExpireDate.Format("2006-01-02")
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdItem, nil
}

// ==================== DeleteItem ====================

func (s *salesService) DeleteItem(ctx context.Context, orderID, itemID, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询订单
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if order.Status != model.SalesOrderStatusPending {
			return ecode.New(ecode.ErrStatusInvalid.Code, "can only delete items from PENDING order")
		}

		// 2. 查询明细
		item, err := s.salesRepo.GetItemByID(ctx, tx, itemID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if item.OrderID != orderID {
			return ecode.ErrNotFound
		}

		// 3. 释放预留
		var rsv model.TraceReservation
		if err := tx.WithContext(ctx).
			Where("sales_order_id = ? AND sales_order_item_id = ? AND status = ?",
				orderID, itemID, model.ReservationStatusReserved).
			First(&rsv).Error; err == nil {
			// 找到预留记录，释放
			now := time.Now()
			if err := s.salesRepo.UpdateReservationStatus(ctx, tx, rsv.ID, model.ReservationStatusReleased, map[string]interface{}{
				"released_at": now,
			}); err != nil {
				return err
			}
		}

		// 4. 软删除明细
		if err := s.salesRepo.DeleteItem(ctx, tx, itemID); err != nil {
			return err
		}

		// 5. 重新计算总金额
		newTotal, err := s.salesRepo.SumItemTotalByOrderID(ctx, tx, orderID)
		if err != nil {
			return err
		}
		return s.salesRepo.UpdateOrderAmounts(ctx, tx, orderID, map[string]interface{}{
			"total_amount":  newTotal,
			"actual_amount": newTotal - order.DiscountAmount,
		})
	})
}

// ==================== GetItems ====================

func (s *salesService) GetItems(ctx context.Context, orderID int64) ([]*model.SalesOrderItem, error) {
	// 确认订单存在
	if _, err := s.salesRepo.GetOrderByID(ctx, s.db, orderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return s.salesRepo.ListItemsByOrderID(ctx, s.db, orderID)
}

// ==================== Pay ====================

func (s *salesService) Pay(ctx context.Context, orderID int64, req PayReq, operatorID int64) error {
	if req.PaymentMethod == "" {
		return ecode.ErrParamInvalid
	}

	// Pre-load order and items for the medicare gateway call (outside DB transaction).
	order, err := s.salesRepo.GetOrderByID(ctx, s.db, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrNotFound
		}
		return err
	}
	if order.Status != model.SalesOrderStatusPending && order.Status != model.SalesOrderStatusApproved {
		return ecode.New(ecode.ErrStatusInvalid.Code, "order cannot be paid in current status")
	}

	preloadedItems, err := s.salesRepo.ListItemsByOrderID(ctx, s.db, orderID)
	if err != nil {
		return err
	}
	if len(preloadedItems) == 0 {
		return ecode.New(ecode.ErrParamInvalid.Code, "order has no items")
	}

	// Validate Medicare patient info when Medicare payment is selected.
	if req.UseMedicare {
		if req.MdtrtCertNo == "" {
			return ecode.New(ecode.ErrParamInvalid.Code, "医保凭证号不能为空")
		}
		if req.PsnNo == "" {
			return ecode.New(ecode.ErrParamInvalid.Code, "参保人员编号不能为空，请先执行医保查询")
		}
	}

	// Medicare settlement: call 2101 → 2102 before the DB transaction.
	// Only proceed when the customer explicitly opted in AND the client is configured.
	var setlInfo *medicare.SetlInfo
	if req.UseMedicare && s.medicareClient != nil {
		info, err := s.doMedicareSettle(ctx, order, preloadedItems, req)
		if err != nil {
			return err
		}
		setlInfo = info
	}

	// DB transaction: validate reservations + mark items SOLD + complete order.
	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Re-lock order.
		lockedOrder, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if lockedOrder.Status != model.SalesOrderStatusPending && lockedOrder.Status != model.SalesOrderStatusApproved {
			return ecode.New(ecode.ErrStatusInvalid.Code, "order cannot be paid in current status")
		}

		now := time.Now()

		for _, item := range preloadedItems {
			var rsv model.TraceReservation
			if err := tx.WithContext(ctx).
				Where("sales_order_id = ? AND sales_order_item_id = ? AND status = ?",
					orderID, item.ID, model.ReservationStatusReserved).
				First(&rsv).Error; err != nil {
				return ecode.New(ecode.ErrConflict.Code, "reservation not found or expired for item")
			}
			if now.After(rsv.ExpireAt) {
				return ecode.New(ecode.ErrConflict.Code, "reservation expired for trace code: "+rsv.TraceCode)
			}

			var trace inventoryModel.TraceInventory
			if err := tx.WithContext(ctx).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("trace_code = ? AND status = ?", item.TraceCode, inventoryModel.TraceInventoryStatusInStock).
				First(&trace).Error; err != nil {
				return ecode.New(ecode.ErrTraceStatusLocked.Code, "trace code is no longer IN_STOCK: "+item.TraceCode)
			}

			soldAt := now
			if err := tx.WithContext(ctx).
				Model(&inventoryModel.TraceInventory{}).
				Where("trace_code = ?", item.TraceCode).
				Updates(map[string]interface{}{
					"status":  inventoryModel.TraceInventoryStatusSold,
					"sold_at": soldAt,
				}).Error; err != nil {
				return err
			}

			if err := s.salesRepo.UpdateReservationStatus(ctx, tx, rsv.ID, model.ReservationStatusConsumed, map[string]interface{}{
				"confirmed_at": now,
			}); err != nil {
				return err
			}

			fromStatus := inventoryModel.TraceInventoryStatusInStock
			toStatus := inventoryModel.TraceInventoryStatusSold
			relatedNo := lockedOrder.OrderNo
			if err := tx.WithContext(ctx).Create(&inventoryModel.DrugTraceLog{
				TraceCode:  item.TraceCode,
				ActionType: "SALE",
				FromStatus: &fromStatus,
				ToStatus:   &toStatus,
				OperatorID: operatorID,
				RelatedNo:  &relatedNo,
			}).Error; err != nil {
				return err
			}
		}

		extraFields := map[string]interface{}{
			"payment_method": req.PaymentMethod,
			"paid_at":        now,
			"actual_amount":  lockedOrder.TotalAmount - lockedOrder.DiscountAmount,
		}
		if setlInfo != nil {
			setlID := setlInfo.SetlID
			mdtrtID := setlInfo.MdtrtID
			extraFields["medicare_transaction_id"] = setlID
			extraFields["mdtrt_id"] = mdtrtID
			extraFields["medicare_amount"] = setlInfo.FundPaySumamt + setlInfo.AcctPay
			extraFields["personal_amount"] = lockedOrder.TotalAmount - setlInfo.FundPaySumamt - setlInfo.AcctPay
		}
		return s.salesRepo.UpdateOrderStatus(ctx, tx, orderID, model.SalesOrderStatusCompleted, extraFields)
	})

	if txErr != nil {
		return txErr
	}

	// After successful commit, upload sale details to 3505 (best-effort).
	if setlInfo != nil && s.medicareClient != nil {
		s.uploadSaleAsync(ctx, order, preloadedItems, req, setlInfo)
	}

	return nil
}

// doMedicareSettle calls 2101 (pre-settle) then 2102 (settle) and returns settlement identifiers.
func (s *salesService) doMedicareSettle(
	ctx context.Context,
	order *model.SalesOrder,
	items []*model.SalesOrderItem,
	req PayReq,
) (*medicare.SetlInfo, error) {
	medType := req.MedType
	if medType == "" {
		medType = "41"
	}
	insuType := req.InsuType
	if insuType == "" {
		insuType = "310"
	}
	acctUsedFlag := req.AcctUsedFlag
	if acctUsedFlag == "" {
		acctUsedFlag = "1"
	}

	begntime := time.Now().Format("2006-01-02 15:04:05")
	drugdetail := make([]map[string]any, 0, len(items))
	for i, item := range items {
		var drug drugModel.DrugInfo
		if err := s.db.WithContext(ctx).First(&drug, item.DrugID).Error; err != nil {
			return nil, fmt.Errorf("load drug %d: %w", item.DrugID, err)
		}
		drugdetail = append(drugdetail, map[string]any{
			"feedetl_sn":          fmt.Sprintf("FEE%04d", i+1),
			"rxno":                fmt.Sprintf("HIRX%04d", i+1),
			"rx_circ_flag":        "1",
			"fee_ocur_time":       begntime,
			"med_list_codg":       drug.MedListCodg,
			"medins_list_codg":    drug.MedinsListCodg,
			"det_item_fee_sumamt": item.SubtotalAmount,
			"cnt":                 item.Quantity,
			"pric":                item.Price,
		})
	}

	mdtrtCertNo := req.MdtrtCertNo
	if mdtrtCertNo == "" {
		mdtrtCertNo = "320500199001010011"
	}
	psnNo := req.PsnNo
	if psnNo == "" {
		psnNo = "PSN0001"
	}

	preSettleInput := map[string]any{
		"druginfo": map[string]any{
			"mdtrt_cert_type":         "02",
			"mdtrt_cert_no":           mdtrtCertNo,
			"psn_no":                  psnNo,
			"med_type":                medType,
			"insutype":                insuType,
			"acct_used_flag":          acctUsedFlag,
			"medfee_sumamt":           order.TotalAmount,
			"begntime":                begntime,
			"minpacunt_drug_trac_cnt": 0,
			"mcs_trac_cnt":            0,
		},
		"drugdetail": drugdetail,
	}

	_, err := s.medicareClient.PreSettle(ctx, medicare.BusinessRequest{
		ERPOrderNo: order.OrderNo,
		Input:      preSettleInput,
	})
	if err != nil {
		s.logger.Warn("medicare 2101 pre-settle failed", zap.String("order_no", order.OrderNo), zap.Error(err))
		return nil, ecode.New(ecode.ErrMedicareTrialFail.Code, "medicare pre-settlement failed: "+err.Error())
	}

	settleInput := map[string]any{
		"druginfo": map[string]any{
			"mdtrt_cert_type": "02",
			"mdtrt_cert_no":   mdtrtCertNo,
			"psn_no":          psnNo,
			"med_type":        medType,
			"insutype":        insuType,
			"acct_used_flag":  acctUsedFlag,
			"medfee_sumamt":   order.TotalAmount,
			"begntime":        begntime,
		},
		"drugdetail": drugdetail,
	}

	settleResp, err := s.medicareClient.Settle(ctx, medicare.BusinessRequest{
		ERPOrderNo: order.OrderNo,
		Input:      settleInput,
	})
	if err != nil {
		s.logger.Error("medicare 2102 settle failed", zap.String("order_no", order.OrderNo), zap.Error(err))
		return nil, ecode.New(ecode.ErrMedicareTrialFail.Code, "medicare settlement failed: "+err.Error())
	}

	setlInfo, err := settleResp.ExtractSetlInfo()
	if err != nil {
		return nil, ecode.New(ecode.ErrMedicareTrialFail.Code, "parse settlement result: "+err.Error())
	}
	return setlInfo, nil
}

// uploadSaleAsync calls 3505 for each item after a successful settlement; failures are logged only.
func (s *salesService) uploadSaleAsync(
	ctx context.Context,
	order *model.SalesOrder,
	items []*model.SalesOrderItem,
	req PayReq,
	setlInfo *medicare.SetlInfo,
) {
	certNo := req.MdtrtCertNo
	if certNo == "" {
		certNo = "320500199001010011"
	}
	psnName := req.PsnName
	if psnName == "" {
		psnName = "参保人"
	}

	for _, item := range items {
		var drug drugModel.DrugInfo
		if err := s.db.WithContext(ctx).First(&drug, item.DrugID).Error; err != nil {
			s.logger.Warn("3505 upload: load drug failed", zap.Int64("drug_id", item.DrugID), zap.Error(err))
			continue
		}
		_, err := s.medicareClient.UploadSale(ctx, medicare.BusinessRequest{
			ERPOrderNo: order.OrderNo,
			Input: map[string]any{
				"selinfo": map[string]any{
					"med_list_codg":         drug.MedListCodg,
					"fixmedins_hilist_id":   drug.FixmedinsHilistId,
					"fixmedins_hilist_name": drug.FixmedinsHilistName,
					"fixmedins_bchno":       item.TraceCode,
					"rtal_docno":            order.OrderNo,
					"setl_id":               setlInfo.SetlID,
					"psn_no":                setlInfo.PsnNo,
					"certno":                certNo,
					"psn_name":              psnName,
					"sel_retn_cnt":          item.Quantity,
					"finl_trns_pric":        item.Price,
				},
			},
		})
		if err != nil {
			s.logger.Warn("medicare 3505 upload failed",
				zap.String("order_no", order.OrderNo),
				zap.String("trace_code", item.TraceCode),
				zap.Error(err),
			)
		}
	}
}

// ==================== Cancel ====================

func (s *salesService) Cancel(ctx context.Context, orderID, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询订单
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		// 仅允许 PENDING/PENDING_REVIEW/APPROVED 取消
		allowedStatuses := map[string]bool{
			model.SalesOrderStatusPending:       true,
			model.SalesOrderStatusPendingReview: true,
			model.SalesOrderStatusApproved:      true,
		}
		if !allowedStatuses[order.Status] {
			return ecode.New(ecode.ErrStatusInvalid.Code, "order cannot be cancelled in current status")
		}

		now := time.Now()

		// 2. 释放所有 RESERVED 预留
		var rsvs []model.TraceReservation
		if err := tx.WithContext(ctx).
			Where("sales_order_id = ? AND status = ?", orderID, model.ReservationStatusReserved).
			Find(&rsvs).Error; err != nil {
			return err
		}
		for _, rsv := range rsvs {
			if err := s.salesRepo.UpdateReservationStatus(ctx, tx, rsv.ID, model.ReservationStatusReleased, map[string]interface{}{
				"released_at": now,
			}); err != nil {
				return err
			}
		}

		// 3. 如果存在待审核的审核记录，将其取消
		var review pharmacistModel.AuditReview
		if err := tx.WithContext(ctx).
			Where("order_id = ? AND status = ?", orderID, pharmacistModel.AuditReviewStatusPending).
			First(&review).Error; err == nil {
			if err := s.reviewRepo.UpdateReview(ctx, tx, review.ID, map[string]interface{}{
				"status": pharmacistModel.AuditReviewStatusCancelled,
			}); err != nil {
				return err
			}
		}

		// 4. 更新订单状态
		return s.salesRepo.UpdateOrderStatus(ctx, tx, orderID, model.SalesOrderStatusCancelled, map[string]interface{}{
			"cancelled_at": now,
		})
	})
}

// ==================== Refund ====================

func (s *salesService) Refund(ctx context.Context, orderID int64, req RefundReq, operatorID int64) error {
	if req.RefundMode != model.RefundModeFull && req.RefundMode != model.RefundModePartial {
		return ecode.ErrParamInvalid
	}
	if req.RefundMode == model.RefundModePartial && len(req.DetailIDs) == 0 {
		return ecode.ErrParamInvalid
	}

	// Pre-load order outside the transaction so we can reference it for the 2103 call.
	var order *model.SalesOrder

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询订单
		o, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		order = o
		// 仅允许 COMPLETED 或 PARTIALLY_REFUNDED 退款
		if order.Status != model.SalesOrderStatusCompleted && order.Status != model.SalesOrderStatusPartiallyRefunded {
			return ecode.New(ecode.ErrStatusInvalid.Code, "order cannot be refunded in current status")
		}

		// 2. 查询所有明细
		allItems, err := s.salesRepo.ListItemsByOrderID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		now := time.Now()

		// 3. 确定退款明细
		var refundItems []*model.SalesOrderItem
		if req.RefundMode == model.RefundModeFull {
			// 全额退款：所有未退款明细
			for _, item := range allItems {
				if item.RefundStatus == model.ItemRefundStatusNone {
					refundItems = append(refundItems, item)
				}
			}
		} else {
			// 部分退款：指定明细 ID
			idSet := make(map[int64]bool, len(req.DetailIDs))
			for _, id := range req.DetailIDs {
				idSet[id] = true
			}
			for _, item := range allItems {
				if idSet[item.ID] {
					if item.RefundStatus != model.ItemRefundStatusNone {
						return ecode.New(ecode.ErrStatusInvalid.Code, "item already refunded")
					}
					refundItems = append(refundItems, item)
				}
			}
			if len(refundItems) != len(req.DetailIDs) {
				return ecode.New(ecode.ErrParamInvalid.Code, "some detail ids not found in order")
			}
		}

		if len(refundItems) == 0 {
			return ecode.New(ecode.ErrParamInvalid.Code, "no items to refund")
		}

		// 4. 处理每个退款明细
		var newRefundTotal float64
		for _, item := range refundItems {
			refundReason := req.RefundReason

			// 更新明细退款状态
			if err := s.salesRepo.UpdateItem(ctx, tx, item.ID, map[string]interface{}{
				"refund_status":      model.ItemRefundStatusRefunded,
				"refund_amount":      item.SubtotalAmount,
				"refunded_at":        now,
				"refund_reason":      refundReason,
				"refund_operator_id": operatorID,
			}); err != nil {
				return err
			}

			newRefundTotal += item.SubtotalAmount

			// 将追溯码从 SOLD 恢复为 IN_STOCK（加行锁）
			var trace inventoryModel.TraceInventory
			if err := tx.WithContext(ctx).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("trace_code = ?", item.TraceCode).
				First(&trace).Error; err != nil {
				return err
			}

			fromStatus := trace.Status
			toStatus := inventoryModel.TraceInventoryStatusInStock
			// 退款：恢复为 IN_STOCK，清空 sold_at（使用原生 SQL 避免 GORM 忽略 nil 值）
			if err := tx.WithContext(ctx).
				Exec("UPDATE drug_trace_inventory SET status = ?, sold_at = NULL, updated_at = NOW() WHERE trace_code = ?",
					inventoryModel.TraceInventoryStatusInStock, item.TraceCode).Error; err != nil {
				return err
			}

			// 写入追溯日志
			relatedNo := order.OrderNo
			remarkStr := "退款退货: " + refundReason
			if err := tx.WithContext(ctx).Create(&inventoryModel.DrugTraceLog{
				TraceCode:  item.TraceCode,
				ActionType: "RETURN",
				FromStatus: &fromStatus,
				ToStatus:   &toStatus,
				OperatorID: operatorID,
				RelatedNo:  &relatedNo,
				Remark:     &remarkStr,
			}).Error; err != nil {
				return err
			}
		}

		// 5. 更新订单退款金额和状态
		updatedRefundAmount := order.RefundAmount + newRefundTotal

		// 判断是全额退款还是部分退款：step4 已将本次明细更新为 REFUNDED，
		// 此处 count 得到的是本次退款后剩余未退款明细数量，无需再减。
		var remainingNone int64
		if err := tx.WithContext(ctx).Model(&model.SalesOrderItem{}).
			Where("order_id = ? AND refund_status = ?", orderID, model.ItemRefundStatusNone).
			Count(&remainingNone).Error; err != nil {
			return err
		}

		newStatus := model.SalesOrderStatusPartiallyRefunded
		if remainingNone <= 0 {
			newStatus = model.SalesOrderStatusRefunded
		}

		extraFields := map[string]interface{}{
			"refund_amount": updatedRefundAmount,
		}
		if req.RefundReason != "" {
			extraFields["refund_reason"] = req.RefundReason
		}
		if newStatus == model.SalesOrderStatusRefunded {
			extraFields["refunded_at"] = now
		}

		return s.salesRepo.UpdateOrderStatus(ctx, tx, orderID, newStatus, extraFields)
	})
	if txErr != nil {
		return txErr
	}

	// For full refund of a medicare-settled order, cancel the medicare settlement (2103).
	if req.RefundMode == model.RefundModeFull &&
		order.NeedMedicare &&
		order.MedicareTransactionID != nil &&
		order.MdtrtID != nil &&
		s.medicareClient != nil {
		_, err := s.medicareClient.CancelSettle(ctx, medicare.BusinessRequest{
			ERPOrderNo: order.OrderNo,
			Input: map[string]any{
				"data": map[string]any{
					"setl_id":  *order.MedicareTransactionID,
					"mdtrt_id": *order.MdtrtID,
					"psn_no":   "PSN0001",
				},
			},
		})
		if err != nil {
			s.logger.Error("medicare 2103 cancel-settle failed",
				zap.String("order_no", order.OrderNo),
				zap.String("setl_id", *order.MedicareTransactionID),
				zap.Error(err),
			)
			// Return error so the caller knows the medicare cancellation failed.
			return ecode.New(ecode.ErrMedicareTrialFail.Code, "medicare settlement cancellation failed: "+err.Error())
		}
	}
	return nil
}

// ==================== GetReservedTraces ====================

func (s *salesService) GetReservedTraces(ctx context.Context, orderID int64) ([]*model.TraceReservation, error) {
	if _, err := s.salesRepo.GetOrderByID(ctx, s.db, orderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return s.salesRepo.GetReservationsByOrderID(ctx, s.db, orderID)
}

// ==================== GetReviewRecord ====================

func (s *salesService) GetReviewRecord(ctx context.Context, orderID int64) (*pharmacistModel.AuditReview, error) {
	review, err := s.reviewRepo.GetLatestReviewByOrderID(ctx, s.db, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	return review, nil
}

// ==================== SubmitReview ====================

func (s *salesService) SubmitReview(ctx context.Context, orderID, submitterID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		// 只有 PENDING 状态的订单可以手动提交审核
		if order.Status != model.SalesOrderStatusPending {
			return ecode.New(ecode.ErrStatusInvalid.Code, "only PENDING order can be submitted for review")
		}

		now := time.Now()
		reviewNo, err := s.reviewRepo.GenReviewNo(ctx, tx)
		if err != nil {
			return err
		}
		submittedAt := now
		review := &pharmacistModel.AuditReview{
			OrderID:     order.ID,
			ReviewNo:    reviewNo,
			Status:      pharmacistModel.AuditReviewStatusPending,
			SubmitterID: &submitterID,
			SubmittedAt: &submittedAt,
		}
		if err := s.reviewRepo.CreateReview(ctx, tx, review); err != nil {
			return err
		}

		return s.salesRepo.UpdateOrderStatus(ctx, tx, orderID, model.SalesOrderStatusPendingReview, map[string]interface{}{
			"need_audit": true,
		})
	})
}

// ==================== ReserveTrace ====================

func (s *salesService) ReserveTrace(ctx context.Context, orderID int64, req ReserveTraceReq, operatorID int64) (*model.TraceReservation, error) {
	var createdRsv *model.TraceReservation

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		// 仅 PENDING/PENDING_REVIEW/APPROVED 允许手动锁定
		allowedStatuses := map[string]bool{
			model.SalesOrderStatusPending:       true,
			model.SalesOrderStatusPendingReview: true,
			model.SalesOrderStatusApproved:      true,
		}
		if !allowedStatuses[order.Status] {
			return ecode.New(ecode.ErrStatusInvalid.Code, "cannot reserve trace in current order status")
		}

		// 验证追溯码
		var trace inventoryModel.TraceInventory
		if err := tx.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("trace_code = ? AND drug_id = ? AND status = ?",
				req.TraceCode, req.DrugID, inventoryModel.TraceInventoryStatusInStock).
			First(&trace).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.New(ecode.ErrTraceCodeNotFound.Code, "trace code not found or unavailable")
			}
			return err
		}

		// 检查是否已被预留
		var rsvCount int64
		if err := tx.WithContext(ctx).Model(&model.TraceReservation{}).
			Where("trace_code = ? AND status = ?", req.TraceCode, model.ReservationStatusReserved).
			Count(&rsvCount).Error; err != nil {
			return err
		}
		if rsvCount > 0 {
			return ecode.New(ecode.ErrConflict.Code, "trace code already reserved")
		}

		rsvNo, err := s.salesRepo.GenReservationNo(ctx, tx)
		if err != nil {
			return err
		}

		now := time.Now()
		rsv := &model.TraceReservation{
			ReservationNo: rsvNo,
			SalesOrderID:  orderID,
			TraceCode:     req.TraceCode,
			DrugID:        req.DrugID,
			ReservedBy:    operatorID,
			Status:        model.ReservationStatusReserved,
			ReservedAt:    now,
			ExpireAt:      now.Add(30 * time.Minute),
		}
		if err := s.salesRepo.CreateReservation(ctx, tx, rsv); err != nil {
			return err
		}
		createdRsv = rsv
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 发布过期消息
	if createdRsv != nil {
		event := mq.ReservationExpireEvent{
			ReservationID: createdRsv.ID,
			ReservationNo: createdRsv.ReservationNo,
			SalesOrderID:  createdRsv.SalesOrderID,
			TraceCode:     createdRsv.TraceCode,
			ExpireAt:      createdRsv.ExpireAt,
		}
		if pubErr := s.mqClient.PublishReservationExpire(ctx, event); pubErr != nil {
			s.logger.Warn("发布预留过期消息失败",
				zap.String("reservation_no", createdRsv.ReservationNo),
				zap.Error(pubErr),
			)
		}
	}

	return createdRsv, nil
}

// ==================== ReleaseReservation ====================

func (s *salesService) ReleaseReservation(ctx context.Context, orderID int64, req ReleaseReservationReq, operatorID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.salesRepo.GetOrderByIDForUpdate(ctx, tx, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		_ = order // 校验订单存在即可

		var rsv model.TraceReservation
		if err := tx.WithContext(ctx).
			Where("sales_order_id = ? AND trace_code = ? AND status = ?",
				orderID, req.TraceCode, model.ReservationStatusReserved).
			First(&rsv).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.New(ecode.ErrNotFound.Code, "no active reservation found for trace code")
			}
			return err
		}

		now := time.Now()
		return s.salesRepo.UpdateReservationStatus(ctx, tx, rsv.ID, model.ReservationStatusReleased, map[string]interface{}{
			"released_at": now,
		})
	})
}

// ==================== MedicarePreview ====================

func (s *salesService) MedicarePreview(ctx context.Context, orderID int64, req MedicarePreviewReq) (*MedicarePreviewResp, error) {
	order, err := s.salesRepo.GetOrderByID(ctx, s.db, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	medType := req.MedType
	if medType == "" {
		medType = "41"
	}
	insuType := req.InsuType
	if insuType == "" {
		insuType = "310"
	}
	acctUsedFlag := req.AcctUsedFlag
	if acctUsedFlag == "" {
		acctUsedFlag = "1"
	}

	// Query patient info via 1101.
	var psnNo, psnName string
	if s.medicareClient != nil {
		personResp, qErr := s.medicareClient.QueryPerson(ctx, medicare.BusinessRequest{
			ERPOrderNo: order.OrderNo,
			Input: map[string]any{
				"data": map[string]any{
					"mdtrt_cert_type": "02",
					"mdtrt_cert_no":   req.MdtrtCertNo,
					"card_sn":         "",
					"begntime":        time.Now().Format("2006-01-02 15:04:05"),
					"psn_cert_type":   "01",
					"certno":          req.MdtrtCertNo,
					"psn_name":        "",
				},
			},
		})
		if qErr != nil {
			return nil, ecode.New(ecode.ErrMedicareTrialFail.Code, "查询参保人信息失败: "+qErr.Error())
		}
		if info, pErr := personResp.ExtractPersonInfo(); pErr == nil {
			psnNo = info.PsnNo
			psnName = info.PsnName
		}
	}
	// Fallback for dev/simulation when client is not configured or response lacks info.
	if psnNo == "" {
		psnNo = "PSN0001"
	}
	if psnName == "" {
		psnName = "参保人"
	}

	total := order.TotalAmount
	fundPay, acctPay := calcMedicareCoverage(total, insuType, medType, acctUsedFlag == "1")
	personalCash := math.Round((total-fundPay-acctPay)*100) / 100
	if personalCash < 0 {
		personalCash = 0
	}

	return &MedicarePreviewResp{
		PsnNo:        psnNo,
		PsnName:      psnName,
		TotalAmount:  total,
		FundPay:      fundPay,
		AcctPay:      acctPay,
		PersonalCash: personalCash,
	}, nil
}

// calcMedicareCoverage returns estimated fund and account payment amounts.
// These are simplified rates used for preview; actual amounts come from the 2102 gateway response.
func calcMedicareCoverage(total float64, insuType, medType string, useAcct bool) (fundPay, acctPay float64) {
	var fundRate, acctRate float64
	switch insuType {
	case "310": // 职工医保
		switch medType {
		case "41": // 门诊
			fundRate = 0.60
			if useAcct {
				acctRate = 0.10
			}
		case "11": // 住院
			fundRate = 0.80
		}
	case "390": // 居民医保
		switch medType {
		case "41": // 门诊
			fundRate = 0.50
		case "11": // 住院
			fundRate = 0.70
		}
	}
	fundPay = math.Round(total*fundRate*100) / 100
	acctPay = math.Round(total*acctRate*100) / 100
	return
}

// ==================== ScanVerify ====================

func (s *salesService) ScanVerify(ctx context.Context, req ScanVerifyReq) (*ScanVerifyResult, error) {
	result := &ScanVerifyResult{
		TraceCode: req.TraceCode,
	}

	// 查询追溯库存
	var trace inventoryModel.TraceInventory
	if err := s.db.WithContext(ctx).Where("trace_code = ?", req.TraceCode).First(&trace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.IsValid = false
			result.Reason = "trace code not found"
			return result, nil
		}
		return nil, err
	}

	result.Status = trace.Status
	result.DrugID = trace.DrugID
	result.BatchNumber = trace.BatchNumber

	// 查询药品信息
	var drug drugModel.DrugInfo
	if err := s.db.WithContext(ctx).Where("id = ?", trace.DrugID).First(&drug).Error; err == nil {
		result.DrugName = drug.CommonName
	}

	// 如果提供了 drug_id，验证归属
	if req.DrugID != nil && *req.DrugID != trace.DrugID {
		result.IsValid = false
		result.Reason = "trace code belongs to a different drug"
		return result, nil
	}

	// 检查状态
	if trace.Status != inventoryModel.TraceInventoryStatusInStock {
		result.IsValid = false
		result.Reason = "trace code is not IN_STOCK (status: " + trace.Status + ")"
		return result, nil
	}

	// 检查是否被预留
	var rsvCount int64
	if err := s.db.WithContext(ctx).Model(&model.TraceReservation{}).
		Where("trace_code = ? AND status = ?", req.TraceCode, model.ReservationStatusReserved).
		Count(&rsvCount).Error; err != nil {
		return nil, err
	}

	result.IsValid = rsvCount == 0
	if rsvCount > 0 {
		result.Reason = "trace code is already reserved"
	}

	return result, nil
}
