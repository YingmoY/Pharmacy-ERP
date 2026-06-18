package model

import (
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
)

// User 对应 sys_user 表
type User struct {
	core.BaseModel
	Username     string     `gorm:"column:username;size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null"       json:"-"`
	RealName     string     `gorm:"column:real_name;size:64"                     json:"real_name"`
	Status       int16      `gorm:"column:status;default:1"                      json:"status"`
	Phone        string     `gorm:"column:phone;size:32"                         json:"phone"`
	Email        string     `gorm:"column:email;size:128"                        json:"email"`
	AvatarURL    string     `gorm:"column:avatar_url;size:512"                   json:"avatar_url"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"                         json:"last_login_at"`
	LastLoginIP  string     `gorm:"column:last_login_ip;size:64"                 json:"last_login_ip"`
	Remark       string     `gorm:"column:remark;size:512"                       json:"remark"`
}

func (User) TableName() string { return "sys_user" }

// UserRole 对应 sys_user_role 中间表（无 id 主键，使用复合主键）
type UserRole struct {
	UserID     int64     `gorm:"column:user_id;primaryKey"  json:"user_id"`
	RoleID     int64     `gorm:"column:role_id;primaryKey"  json:"role_id"`
	AssignedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"assigned_at"`
}

func (UserRole) TableName() string { return "sys_user_role" }

// LoginLog 用于插入登录日志记录（仅写入）
type LoginLog struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id"`
	Username  string    `gorm:"column:username;size:64"`
	IP        string    `gorm:"column:ip;size:64"`
	UserAgent string    `gorm:"column:user_agent;size:512"`
	Success   bool      `gorm:"column:success"`
	Message   string    `gorm:"column:message;size:255"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (LoginLog) TableName() string { return "login_log" }

// RoleSimple 角色简要信息（用于 UserDTO）
type RoleSimple struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// UserDTO 用于对外暴露的用户信息，不含敏感字段
type UserDTO struct {
	ID          int64        `json:"id"`
	Username    string       `json:"username"`
	RealName    string       `json:"real_name"`
	Phone       string       `json:"phone"`
	Email       string       `json:"email"`
	AvatarURL   string       `json:"avatar_url"`
	Status      int16        `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	LastLoginAt *time.Time   `json:"last_login_at"`
	LastLoginIP string       `json:"last_login_ip"`
	Roles       []RoleSimple `json:"roles"`
}

// ToDTO 将 User 转换为对外暴露的 DTO
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:          u.ID,
		Username:    u.Username,
		RealName:    u.RealName,
		Phone:       u.Phone,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
		LastLoginIP: u.LastLoginIP,
	}
}
