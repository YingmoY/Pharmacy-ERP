package repository

import (
	"context"

	"github.com/YingmoY/PharmacyERP/internal/role/model"
	"gorm.io/gorm"
)

// RoleRepository 瑙掕壊鏁版嵁璁块棶鎺ュ彛
type RoleRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Role, error)
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	List(ctx context.Context, filter RoleListFilter) ([]model.Role, int64, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error)
	SetRolePermissions(ctx context.Context, roleID int64, permIDs []int64) error
	IsRoleInUse(ctx context.Context, roleID int64) (bool, error)

	ListPermissions(ctx context.Context) ([]model.Permission, error)
	FindPermissionByID(ctx context.Context, id int64) (*model.Permission, error)
	FindPermissionByCode(ctx context.Context, code string) (*model.Permission, error)
}

// RoleListFilter 瑙掕壊鍒楄〃鏌ヨ杩囨护鏉′欢
type RoleListFilter struct {
	Name     string
	Status   *int16
	Page     int
	PageSize int
}

type roleRepo struct {
	db *gorm.DB
}

// NewRoleRepository 鍒涘缓瑙掕壊浠撳偍瀹炰緥
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepo{db: db}
}

// FindByID 鏍规嵁 ID 鏌ヨ瑙掕壊
func (r *roleRepo) FindByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 鏍规嵁 code 鏌ヨ瑙掕壊
func (r *roleRepo) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 鍒嗛〉鏌ヨ瑙掕壊鍒楄〃
func (r *roleRepo) List(ctx context.Context, filter RoleListFilter) ([]model.Role, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Role{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
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

	var roles []model.Role
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&roles).Error
	return roles, total, err
}

// Create 鍒涘缓瑙掕壊
func (r *roleRepo) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Update 鏇存柊瑙掕壊
func (r *roleRepo) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete deletes a role by id.
func (r *roleRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// GetRolePermissions 鑾峰彇瑙掕壊鍏宠仈鐨勬潈闄愬垪琛紙閫氳繃 casbin_rule 琛ㄥ叧鑱旓級
func (r *roleRepo) GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	// 鍏堣幏鍙栬鑹?code
	var permissions []model.Permission
	err := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Preload("APIs").
		Joins("INNER JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id = ?", roleID).
		Find(&permissions).Error
	fillPermissionAPIFields(permissions)
	return permissions, err
}

// SetRolePermissions 鎵归噺璁剧疆瑙掕壊鏉冮檺锛堥€氳繃鏇存柊 casbin_rule 琛級
func (r *roleRepo) SetRolePermissions(ctx context.Context, roleID int64, permIDs []int64) error {
	// 鑾峰彇瑙掕壊 code
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 鍒犻櫎璇ヨ鑹茬幇鏈夌殑 casbin p 瑙勫垯
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		if len(permIDs) == 0 {
			return nil
		}

		rolePerms := make([]model.RolePermission, 0, len(permIDs))
		for _, permID := range permIDs {
			rolePerms = append(rolePerms, model.RolePermission{RoleID: roleID, PermissionID: permID})
		}
		return tx.Create(&rolePerms).Error
	})
}

// IsRoleInUse 妫€鏌ヨ鑹叉槸鍚﹁鐢ㄦ埛浣跨敤
func (r *roleRepo) IsRoleInUse(ctx context.Context, roleID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("sys_user_role").
		Where("role_id = ?", roleID).Count(&count).Error
	return count > 0, err
}

// ListPermissions returns all permissions.
func (r *roleRepo) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Preload("APIs").Find(&permissions).Error
	fillPermissionAPIFields(permissions)
	return permissions, err
}

// FindPermissionByID 鏍规嵁 ID 鏌ヨ鏉冮檺
func (r *roleRepo) FindPermissionByID(ctx context.Context, id int64) (*model.Permission, error) {
	var perm model.Permission
	err := r.db.WithContext(ctx).Preload("APIs").Where("id = ?", id).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *roleRepo) FindPermissionByCode(ctx context.Context, code string) (*model.Permission, error) {
	var perm model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func fillPermissionAPIFields(permissions []model.Permission) {
	for i := range permissions {
		if len(permissions[i].APIs) == 0 {
			continue
		}
		permissions[i].Method = permissions[i].APIs[0].HttpMethod
		permissions[i].Path = permissions[i].APIs[0].PathPattern
	}
}
