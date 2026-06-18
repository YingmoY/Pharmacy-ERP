package service

import (
	"context"
	"errors"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pharmacist/model"
	"github.com/YingmoY/PharmacyERP/internal/pharmacist/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	salesModel "github.com/YingmoY/PharmacyERP/internal/sales/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ReviewFilter 审核列表过滤条件
type ReviewFilter struct {
	Status   string
	OrderNo  string
	Page     int
	PageSize int
}

// ReviewDetail 审核详情响应，包含关联的完整销售单信息。
type ReviewDetail struct {
	*model.AuditReview
	Order *salesModel.SalesOrder `json:"order,omitempty"`
}

// ReviewService 定义药师审核核心业务能力
type ReviewService interface {
	// ListReviews 分页查询审核记录
	ListReviews(ctx context.Context, filter ReviewFilter) ([]*model.AuditReview, int64, error)
	// GetReview 根据 ID 查询审核记录（含完整销售单明细）
	GetReview(ctx context.Context, id int64) (*ReviewDetail, error)
	// ApproveReview 药师审核通过
	ApproveReview(ctx context.Context, id, pharmacistID int64, opinion string) error
	// RejectReview 药师审核驳回
	RejectReview(ctx context.Context, id, pharmacistID int64, opinion string) error
}

type reviewService struct {
	db         *gorm.DB
	reviewRepo repository.ReviewRepo
	logger     *zap.Logger
}

// NewReviewService 创建药师审核服务实例
func NewReviewService(
	db *gorm.DB,
	reviewRepo repository.ReviewRepo,
	logger *zap.Logger,
) ReviewService {
	return &reviewService{
		db:         db,
		reviewRepo: reviewRepo,
		logger:     logger,
	}
}

// ListReviews 分页查询审核记录列表
func (s *reviewService) ListReviews(ctx context.Context, filter ReviewFilter) ([]*model.AuditReview, int64, error) {
	reviews, total, err := s.reviewRepo.ListReviews(ctx, s.db, repository.ReviewFilter{
		Status:   filter.Status,
		OrderNo:  filter.OrderNo,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	s.enrichReviews(ctx, reviews)
	return reviews, total, nil
}

// GetReview 根据 ID 查询审核记录（含完整销售单明细）
func (s *reviewService) GetReview(ctx context.Context, id int64) (*ReviewDetail, error) {
	review, err := s.reviewRepo.GetReviewByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	s.enrichReviews(ctx, []*model.AuditReview{review})

	// 拉取完整销售单（含明细）
	var order salesModel.SalesOrder
	if err2 := s.db.WithContext(ctx).Preload("Items").First(&order, review.OrderID).Error; err2 == nil {
		s.enrichOrderItems(ctx, order.Items)
		return &ReviewDetail{AuditReview: review, Order: &order}, nil
	}
	return &ReviewDetail{AuditReview: review}, nil
}

// enrichOrderItems 批量填充销售单明细的药品名称、规格、处方药标识、批号、有效期。
func (s *reviewService) enrichOrderItems(ctx context.Context, items []salesModel.SalesOrderItem) {
	if len(items) == 0 {
		return
	}

	drugIDSet := make(map[int64]struct{})
	traceCodes := make([]string, 0, len(items))
	for i := range items {
		drugIDSet[items[i].DrugID] = struct{}{}
		traceCodes = append(traceCodes, items[i].TraceCode)
		items[i].UnitPrice = items[i].Price
	}

	drugIDs := make([]int64, 0, len(drugIDSet))
	for id := range drugIDSet {
		drugIDs = append(drugIDs, id)
	}

	type drugRow struct {
		ID             int64
		CommonName     string
		Specification  string
		IsPrescription bool
	}
	var drugs []drugRow
	s.db.WithContext(ctx).
		Table("drug_info").
		Select("id, common_name, specification, is_prescription").
		Where("id IN ?", drugIDs).
		Scan(&drugs)
	drugMap := make(map[int64]drugRow, len(drugs))
	for _, d := range drugs {
		drugMap[d.ID] = d
	}

	type traceRow struct {
		TraceCode   string
		BatchNumber string
		ExpireDate  string
	}
	var traces []traceRow
	s.db.WithContext(ctx).
		Table("drug_trace_inventory").
		Select("trace_code, batch_number, TO_CHAR(expire_date, 'YYYY-MM-DD') AS expire_date").
		Where("trace_code IN ?", traceCodes).
		Scan(&traces)
	traceMap := make(map[string]traceRow, len(traces))
	for _, t := range traces {
		traceMap[t.TraceCode] = t
	}

	for i := range items {
		if d, ok := drugMap[items[i].DrugID]; ok {
			items[i].DrugName = d.CommonName
			items[i].Specification = d.Specification
			items[i].IsPrescription = d.IsPrescription
		}
		if t, ok := traceMap[items[i].TraceCode]; ok {
			items[i].BatchNumber = t.BatchNumber
			items[i].ExpireDate = t.ExpireDate
		}
	}
}

// enrichReviews 批量为审核记录填充提交人名、药师名、销售单号
func (s *reviewService) enrichReviews(ctx context.Context, reviews []*model.AuditReview) {
	if len(reviews) == 0 {
		return
	}

	// 收集需要查询的用户 ID 和订单 ID
	userIDSet := make(map[int64]string)
	orderIDToNo := make(map[int64]string)
	var userIDs []int64
	var orderIDs []int64

	for _, r := range reviews {
		if r.SubmitterID != nil && *r.SubmitterID > 0 {
			if _, exists := userIDSet[*r.SubmitterID]; !exists {
				userIDSet[*r.SubmitterID] = ""
				userIDs = append(userIDs, *r.SubmitterID)
			}
		}
		if r.PharmacistID != nil && *r.PharmacistID > 0 {
			if _, exists := userIDSet[*r.PharmacistID]; !exists {
				userIDSet[*r.PharmacistID] = ""
				userIDs = append(userIDs, *r.PharmacistID)
			}
		}
		orderIDs = append(orderIDs, r.OrderID)
	}

	// 批量查询用户名
	if len(userIDs) > 0 {
		type userRow struct {
			ID       int64
			RealName string
		}
		var users []userRow
		s.db.WithContext(ctx).Table("sys_user").Select("id, real_name").Where("id IN ?", userIDs).Scan(&users)
		for _, u := range users {
			userIDSet[u.ID] = u.RealName
		}
	}

	// 批量查询订单号
	if len(orderIDs) > 0 {
		type orderRow struct {
			ID      int64
			OrderNo string
		}
		var orders []orderRow
		s.db.WithContext(ctx).Table("sales_order").Select("id, order_no").Where("id IN ?", orderIDs).Scan(&orders)
		for _, o := range orders {
			orderIDToNo[o.ID] = o.OrderNo
		}
	}

	// 填充虚拟字段
	for _, r := range reviews {
		if r.SubmitterID != nil {
			r.SubmitterName = userIDSet[*r.SubmitterID]
		}
		if r.PharmacistID != nil {
			r.PharmacistName = userIDSet[*r.PharmacistID]
		}
		r.OrderNo = orderIDToNo[r.OrderID]
	}
}

// ApproveReview 药师审核通过：
// 1. 查询审核记录，验证状态为 PENDING
// 2. 更新审核记录为 APPROVED
// 3. 更新销售订单状态为 APPROVED
func (s *reviewService) ApproveReview(ctx context.Context, id, pharmacistID int64, opinion string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询审核记录
		review, err := s.reviewRepo.GetReviewByIDForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if review.Status != model.AuditReviewStatusPending {
			return ecode.New(ecode.ErrStatusInvalid.Code, "review is not in PENDING status")
		}

		now := time.Now()

		// 2. 更新审核记录为 APPROVED
		updateFields := map[string]interface{}{
			"status":          model.AuditReviewStatusApproved,
			"pharmacist_id":   pharmacistID,
			"reviewed_at":     now,
		}
		if opinion != "" {
			updateFields["review_opinion"] = opinion
		}
		if err := s.reviewRepo.UpdateReview(ctx, tx, id, updateFields); err != nil {
			return err
		}

		// 3. 更新销售订单状态为 APPROVED
		if err := tx.WithContext(ctx).
			Model(&salesModel.SalesOrder{}).
			Where("id = ? AND status = ?", review.OrderID, salesModel.SalesOrderStatusPendingReview).
			Updates(map[string]interface{}{
				"status": salesModel.SalesOrderStatusApproved,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// RejectReview 药师审核驳回：
// 1. 查询审核记录，验证状态为 PENDING
// 2. 更新审核记录为 REJECTED
// 3. 释放订单所有 RESERVED 预留
// 4. 更新销售订单状态为 CANCELLED
func (s *reviewService) RejectReview(ctx context.Context, id, pharmacistID int64, opinion string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 加锁查询审核记录
		review, err := s.reviewRepo.GetReviewByIDForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrNotFound
			}
			return err
		}
		if review.Status != model.AuditReviewStatusPending {
			return ecode.New(ecode.ErrStatusInvalid.Code, "review is not in PENDING status")
		}

		now := time.Now()

		// 2. 更新审核记录为 REJECTED
		updateFields := map[string]interface{}{
			"status":        model.AuditReviewStatusRejected,
			"pharmacist_id": pharmacistID,
			"reviewed_at":   now,
		}
		if opinion != "" {
			updateFields["review_opinion"] = opinion
		}
		if err := s.reviewRepo.UpdateReview(ctx, tx, id, updateFields); err != nil {
			return err
		}

		// 3. 释放订单下所有 RESERVED 预留
		var rsvs []salesModel.TraceReservation
		if err := tx.WithContext(ctx).
			Where("sales_order_id = ? AND status = ?", review.OrderID, salesModel.ReservationStatusReserved).
			Find(&rsvs).Error; err != nil {
			return err
		}
		for _, rsv := range rsvs {
			if err := tx.WithContext(ctx).
				Model(&salesModel.TraceReservation{}).
				Where("id = ?", rsv.ID).
				Updates(map[string]interface{}{
					"status":      salesModel.ReservationStatusReleased,
					"released_at": now,
				}).Error; err != nil {
				return err
			}
		}

		// 4. 更新销售订单状态为 CANCELLED
		if err := tx.WithContext(ctx).
			Model(&salesModel.SalesOrder{}).
			Where("id = ?", review.OrderID).
			Updates(map[string]interface{}{
				"status":       salesModel.SalesOrderStatusCancelled,
				"cancelled_at": now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}
