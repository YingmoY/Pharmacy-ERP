package repository

import (
	"context"

	rolemodel "github.com/YingmoY/PharmacyERP/internal/role/model"
	"github.com/YingmoY/PharmacyERP/internal/user/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, filter ListFilter) ([]model.User, int64, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateStatus(ctx context.Context, id int64, status int16) error
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	GetUserRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error)
	GetUserRolesBatch(ctx context.Context, userIDs []int64) (map[int64][]model.RoleSimple, error)
	SetUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	GetUserPermissions(ctx context.Context, userID int64) ([]rolemodel.Permission, error)
}

// ListFilter 用户列表查询过滤条件
type ListFilter struct {
	Username string
	RealName string
	Status   *int16
	Page     int
	PageSize int
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// FindByID 根据 ID 查询用户（过滤软删除）
func (r *userRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户（过滤软删除）
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 分页查询用户列表
func (r *userRepo) List(ctx context.Context, filter ListFilter) ([]model.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.User{})

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.RealName != "" {
		query = query.Where("real_name LIKE ?", "%"+filter.RealName+"%")
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
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

	var users []model.User
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户信息
func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateStatus 更新用户状态
func (r *userRepo) UpdateStatus(ctx context.Context, id int64, status int16) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdatePassword 更新用户密码哈希
func (r *userRepo) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

// GetUserRoles 获取用户绑定的所有角色
func (r *userRepo) GetUserRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error) {
	var roles []rolemodel.Role
	err := r.db.WithContext(ctx).
		Table("sys_role").
		Joins("INNER JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ? AND sys_role.deleted_at IS NULL", userID).
		Find(&roles).Error
	return roles, err
}

// GetUserRolesBatch 批量获取多个用户的角色（返回 userID → []RoleSimple 映射）
func (r *userRepo) GetUserRolesBatch(ctx context.Context, userIDs []int64) (map[int64][]model.RoleSimple, error) {
	if len(userIDs) == 0 {
		return map[int64][]model.RoleSimple{}, nil
	}
	type row struct {
		UserID   int64  `gorm:"column:user_id"`
		RoleID   int64  `gorm:"column:role_id"`
		RoleCode string `gorm:"column:role_code"`
		RoleName string `gorm:"column:role_name"`
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("sys_user_role").
		Select("sys_user_role.user_id, sys_role.id AS role_id, sys_role.code AS role_code, sys_role.name AS role_name").
		Joins("INNER JOIN sys_role ON sys_role.id = sys_user_role.role_id AND sys_role.deleted_at IS NULL").
		Where("sys_user_role.user_id IN ?", userIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]model.RoleSimple)
	for _, rr := range rows {
		result[rr.UserID] = append(result[rr.UserID], model.RoleSimple{
			ID:   rr.RoleID,
			Code: rr.RoleCode,
			Name: rr.RoleName,
		})
	}
	return result, nil
}

// SetUserRoles 批量设置用户角色（先删后插）
func (r *userRepo) SetUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		// 批量插入新关联
		userRoles := make([]model.UserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			userRoles = append(userRoles, model.UserRole{UserID: userID, RoleID: rid})
		}
		return tx.Create(&userRoles).Error
	})
}

// GetUserPermissions 获取用户所有角色关联的权限（去重）
func (r *userRepo) GetUserPermissions(ctx context.Context, userID int64) ([]rolemodel.Permission, error) {
	var permissions []rolemodel.Permission
	err := r.db.WithContext(ctx).
		Table("sys_permission").
		Joins("INNER JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Joins("INNER JOIN sys_role ON sys_role.id = sys_role_permission.role_id AND sys_role.deleted_at IS NULL").
		Joins("INNER JOIN sys_user_role ON sys_user_role.role_id = sys_role.id AND sys_user_role.user_id = ?", userID).
		Distinct("sys_permission.*").
		Find(&permissions).Error
	return permissions, err
}
