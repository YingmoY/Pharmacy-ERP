package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// 审核状态常量
const (
	// AuditReviewStatusPending 待审核
	AuditReviewStatusPending = "PENDING"
	// AuditReviewStatusApproved 审核通过
	AuditReviewStatusApproved = "APPROVED"
	// AuditReviewStatusRejected 审核驳回
	AuditReviewStatusRejected = "REJECTED"
	// AuditReviewStatusCancelled 已取消
	AuditReviewStatusCancelled = "CANCELLED"
)

// AuditReview 映射 public.audit_review 表，记录药师审核信息
type AuditReview struct {
	core.BaseModel
	OrderID       int64      `gorm:"column:order_id;type:bigint;not null;index"             json:"order_id"`
	ReviewNo      string     `gorm:"column:review_no;type:varchar(50);not null;uniqueIndex" json:"review_no"`
	PharmacistID  *int64     `gorm:"column:pharmacist_id;type:bigint"                       json:"pharmacist_id,omitempty"`
	Status        string     `gorm:"column:status;type:varchar(20);not null;index"          json:"status"`
	Comment       *string    `gorm:"column:comment;type:text"                               json:"comment,omitempty"`
	ReviewedAt    *time.Time `gorm:"column:reviewed_at;type:timestamptz"                    json:"reviewed_at,omitempty"`
	SubmitterID   *int64     `gorm:"column:submitter_id;type:bigint"                        json:"submitter_id,omitempty"`
	SubmittedAt   *time.Time `gorm:"column:submitted_at;type:timestamptz"                   json:"submitted_at,omitempty"`
	ReviewOpinion *string    `gorm:"column:review_opinion;type:text"                        json:"review_opinion,omitempty"`
	// 虚拟字段：通过关联查询填充，不存储于数据库
	OrderNo        string `gorm:"-" json:"order_no,omitempty"`
	SubmitterName  string `gorm:"-" json:"submitter_name,omitempty"`
	PharmacistName string `gorm:"-" json:"pharmacist_name,omitempty"`
}

// TableName 返回数据库表名
func (AuditReview) TableName() string {
	return "audit_review"
}
