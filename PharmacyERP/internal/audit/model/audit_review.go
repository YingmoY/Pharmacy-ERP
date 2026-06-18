package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// AuditReviewStatusPassed 审核通过。
	AuditReviewStatusPassed = "PASSED"
	// AuditReviewStatusRejected 审核驳回。
	AuditReviewStatusRejected = "REJECTED"
)

// AuditReview 映射 public.audit_review。
type AuditReview struct {
	core.BaseModel
	OrderID      int64     `gorm:"column:order_id;type:bigint;not null;index:idx_audit_review_order" json:"order_id"`
	PharmacistID int64     `gorm:"column:pharmacist_id;type:bigint;not null" json:"pharmacist_id"`
	Status       string    `gorm:"column:status;type:varchar(20);not null" json:"status"`
	Comment      *string   `gorm:"column:comment;type:text" json:"comment,omitempty"`
	ReviewedAt   time.Time `gorm:"column:reviewed_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"reviewed_at"`
}

func (AuditReview) TableName() string {
	return "audit_review"
}
