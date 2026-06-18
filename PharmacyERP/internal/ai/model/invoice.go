// Package model 定义 AI 发票识别相关数据模型。
package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"gorm.io/datatypes"
)

// AI 发票识别状态常量。
const (
	// AIInvoiceStatusPending 待处理（刚创建，尚未开始识别）。
	AIInvoiceStatusPending = "PENDING"
	// AIInvoiceStatusProcessing 识别中。
	AIInvoiceStatusProcessing = "PROCESSING"
	// AIInvoiceStatusCompleted 识别完成。
	AIInvoiceStatusCompleted = "COMPLETED"
	// AIInvoiceStatusFailed 识别失败。
	AIInvoiceStatusFailed = "FAILED"
)

// AIInvoiceRecord AI 发票识别记录。
// 对应数据库表：ai_invoice_record。
type AIInvoiceRecord struct {
	core.BaseModel
	// FileID 上传文件 ID（外部文件服务）。
	FileID string `gorm:"column:file_id;type:varchar(200);not null" json:"file_id"`
	// FileName 原始文件名。
	FileName string `gorm:"column:file_name;type:varchar(200)" json:"file_name"`
	// Status 识别状态：PENDING/PROCESSING/COMPLETED/FAILED。
	Status string `gorm:"column:status;type:varchar(20);not null;default:'PENDING'" json:"status"`

	// RecognizedSupplierName AI 识别出的供应商名称。
	RecognizedSupplierName *string `gorm:"column:recognized_supplier_name;type:varchar(200)" json:"recognized_supplier_name,omitempty"`
	// MatchedSupplierID 系统匹配到的供应商 ID。
	MatchedSupplierID *int64 `gorm:"column:matched_supplier_id;type:bigint" json:"matched_supplier_id,omitempty"`

	// InvoiceNo 发票号。
	InvoiceNo *string `gorm:"column:invoice_no;type:varchar(100)" json:"invoice_no,omitempty"`
	// InvoiceDate 发票日期。
	InvoiceDate *time.Time `gorm:"column:invoice_date;type:date" json:"invoice_date,omitempty"`

	// ResultJSON AI 识别结果（结构化 JSON）；JSON 序列化为 "result" 与前端约定对齐。
	ResultJSON datatypes.JSON `gorm:"column:result_json;type:jsonb" json:"result,omitempty"`
	// RawResponseJSON AI 原始响应，仅供后端调试，不暴露给前端。
	RawResponseJSON datatypes.JSON `gorm:"column:raw_response_json;type:jsonb" json:"-"`

	// ErrorMessage 识别失败时的错误信息。
	ErrorMessage *string `gorm:"column:error_message;type:text" json:"error_message,omitempty"`

	// InboundOrderID 关联生成的入库单 ID（转换后写回）。
	InboundOrderID *int64 `gorm:"column:inbound_order_id;type:bigint" json:"inbound_order_id,omitempty"`
	// ConvertedAt 转换为入库单的时间。
	ConvertedAt *time.Time `gorm:"column:converted_at" json:"converted_at,omitempty"`

	// CreatorID 创建人 ID。
	CreatorID int64 `gorm:"column:creator_id;type:bigint;not null" json:"creator_id"`
	// Remark 备注。
	Remark *string `gorm:"column:remark;type:text" json:"remark,omitempty"`

	// 虚拟字段，由服务层批量填充，不存储于数据库。
	SupplierName string  `gorm:"-" json:"supplier_name,omitempty"`
	DrugCount    int     `gorm:"-" json:"drug_count"`
	TotalAmount  float64 `gorm:"-" json:"total_amount"`
	Confidence   float64 `gorm:"-" json:"confidence"`
}

// TableName 指定数据库表名。
func (AIInvoiceRecord) TableName() string {
	return "ai_invoice_record"
}
