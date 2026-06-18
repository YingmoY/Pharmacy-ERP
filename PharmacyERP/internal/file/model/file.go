package model

import (
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

const (
	// FileStatusActive 文件正常状态。
	FileStatusActive int8 = 1
	// FileStatusDeleted 文件逻辑删除状态（软删除外的额外业务状态）。
	FileStatusDeleted int8 = 0
)

// FileInfo 代表上传文件的元数据信息，对应数据库表 file_info。
// 字段按实际表结构映射：storage_name / storage_path / content_type / file_hash 等。
type FileInfo struct {
	core.BaseModel
	FileID       string  `gorm:"column:file_id;uniqueIndex;size:100;not null" json:"file_id"`
	OriginalName string  `gorm:"column:original_name;size:255;not null" json:"original_name"`
	StorageName  string  `gorm:"column:storage_name;size:255;not null" json:"storage_name"`
	StoragePath  string  `gorm:"column:storage_path;size:500;not null" json:"storage_path"`
	ContentType  string  `gorm:"column:content_type;size:100" json:"content_type"`
	FileSize     int64   `gorm:"column:file_size;not null" json:"file_size"`
	FileHash     string  `gorm:"column:file_hash;size:100" json:"file_hash"`
	BusinessType string  `gorm:"column:business_type;size:50" json:"business_type"`
	BusinessID   string  `gorm:"column:business_id;size:100" json:"business_id"`
	UploaderID   int64   `gorm:"column:uploader_id;not null" json:"uploader_id"`
	Status       int8    `gorm:"column:status;default:1" json:"status"`
	Remark       *string `gorm:"column:remark;type:text" json:"remark,omitempty"`
}

// TableName 返回文件信息表名。
func (FileInfo) TableName() string {
	return "file_info"
}
