package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pharmacist/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReviewFilter 审核列表查询过滤条件
type ReviewFilter struct {
	Status   string
	OrderNo  string
	Page     int
	PageSize int
}

// ReviewRepo 定义药师审核表数据库操作能力
type ReviewRepo interface {
	// CreateReview 创建审核记录（事务内使用）
	CreateReview(ctx context.Context, tx *gorm.DB, review *model.AuditReview) error
	// GetReviewByID 根据 ID 查询审核记录
	GetReviewByID(ctx context.Context, db *gorm.DB, id int64) (*model.AuditReview, error)
	// GetReviewByIDForUpdate 根据 ID 加行锁查询审核记录（事务内使用）
	GetReviewByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.AuditReview, error)
	// GetPendingReviewByOrderID 查询订单下待审核记录
	GetPendingReviewByOrderID(ctx context.Context, db *gorm.DB, orderID int64) (*model.AuditReview, error)
	// GetLatestReviewByOrderID 查询订单最新审核记录
	GetLatestReviewByOrderID(ctx context.Context, db *gorm.DB, orderID int64) (*model.AuditReview, error)
	// UpdateReview 更新审核记录字段
	UpdateReview(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error
	// ListReviews 分页查询审核列表
	ListReviews(ctx context.Context, db *gorm.DB, filter ReviewFilter) ([]*model.AuditReview, int64, error)
	// GenReviewNo 生成唯一审核单号（REV-YYYYMMDD-XXXX）
	GenReviewNo(ctx context.Context, db *gorm.DB) (string, error)
}

// reviewRepo 是 ReviewRepo 的 GORM 实现
type reviewRepo struct{}

// NewReviewRepo 创建 ReviewRepo 实例
func NewReviewRepo() ReviewRepo {
	return &reviewRepo{}
}

func (r *reviewRepo) CreateReview(ctx context.Context, tx *gorm.DB, review *model.AuditReview) error {
	return tx.WithContext(ctx).Create(review).Error
}

func (r *reviewRepo) GetReviewByID(ctx context.Context, db *gorm.DB, id int64) (*model.AuditReview, error) {
	var review model.AuditReview
	err := db.WithContext(ctx).Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) GetReviewByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*model.AuditReview, error) {
	var review model.AuditReview
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) GetPendingReviewByOrderID(ctx context.Context, db *gorm.DB, orderID int64) (*model.AuditReview, error) {
	var review model.AuditReview
	err := db.WithContext(ctx).
		Where("order_id = ? AND status = ?", orderID, model.AuditReviewStatusPending).
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) GetLatestReviewByOrderID(ctx context.Context, db *gorm.DB, orderID int64) (*model.AuditReview, error) {
	var review model.AuditReview
	err := db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) UpdateReview(ctx context.Context, tx *gorm.DB, id int64, fields map[string]interface{}) error {
	return tx.WithContext(ctx).
		Model(&model.AuditReview{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *reviewRepo) ListReviews(ctx context.Context, db *gorm.DB, filter ReviewFilter) ([]*model.AuditReview, int64, error) {
	query := db.WithContext(ctx).Model(&model.AuditReview{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OrderNo != "" {
		query = query.Where("order_id IN (SELECT id FROM sales_order WHERE order_no = ?)", filter.OrderNo)
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

	var reviews []*model.AuditReview
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

// GenReviewNo 生成唯一审核单号，格式 REV-YYYYMMDD-XXXX
func (r *reviewRepo) GenReviewNo(ctx context.Context, db *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "REV-" + today + "-"

	var maxNo string
	err := db.WithContext(ctx).
		Model(&model.AuditReview{}).
		Where("review_no LIKE ?", prefix+"%").
		Order("review_no DESC").
		Limit(1).
		Pluck("review_no", &maxNo).Error
	if err != nil {
		return "", err
	}

	seq := 1
	if maxNo != "" && len(maxNo) > len(prefix) {
		suffix := maxNo[len(prefix):]
		var n int
		if _, scanErr := fmt.Sscanf(suffix, "%d", &n); scanErr == nil {
			seq = n + 1
		}
	}
	return fmt.Sprintf("%s%04d", prefix, seq), nil
}
