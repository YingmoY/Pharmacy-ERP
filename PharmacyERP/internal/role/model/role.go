package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// Role 对应 sys_role 表
type Role struct {
	core.BaseModel
	Code        string `gorm:"column:code;size:64;not null;uniqueIndex" json:"code"`
	Name        string `gorm:"column:name;size:128;not null"            json:"name"`
	Description string `gorm:"column:description;size:512"             json:"description"`
	BuiltIn     bool   `gorm:"column:built_in;default:false"            json:"built_in"`
	Status      int16  `gorm:"column:status;default:1"                  json:"status"`
}

func (Role) TableName() string { return "sys_role" }

// Permission 对应 sys_permission 表
type Permission struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"column:code;size:128;not null;uniqueIndex" json:"code"`
	Name        string    `gorm:"column:name;size:128;not null"             json:"name"`
	Description string    `gorm:"column:description;size:512"              json:"description"`
	Resource    string    `gorm:"column:resource;size:128"                 json:"resource"`
	Action      string    `gorm:"column:action;size:64"                    json:"action"`
	Method      string    `gorm:"-" json:"method,omitempty"`
	Path        string    `gorm:"-" json:"path,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"         json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"         json:"updated_at"`

	// 关联的 API 列表（不存储到主表）
	APIs []PermissionAPI `gorm:"foreignKey:PermissionCode;references:Code" json:"apis,omitempty"`
}

func (Permission) TableName() string { return "sys_permission" }

// PermissionAPI 对应 sys_permission_api 表
type PermissionAPI struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement"         json:"id"`
	PermissionCode string    `gorm:"column:permission_code;not null;index"      json:"permission_code"`
	PathPattern    string    `gorm:"column:path_pattern;size:255;not null"      json:"path_pattern"`
	HttpMethod     string    `gorm:"column:http_method;size:20;not null"        json:"http_method"`
	Summary        string    `gorm:"column:summary;size:255"                    json:"summary,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"           json:"created_at"`
}

func (PermissionAPI) TableName() string { return "sys_permission_api" }

// RolePermission 表示角色的权限 DTO（用于响应）
type RolePermissionDTO struct {
	RoleID      int64        `json:"role_id"`
	Permissions []Permission `json:"permissions"`
}

type RolePermission struct {
	RoleID       int64 `gorm:"column:role_id;primaryKey"`
	PermissionID int64 `gorm:"column:permission_id;primaryKey"`
}

func (RolePermission) TableName() string { return "sys_role_permission" }
