package service

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/role/model"
	"github.com/YingmoY/PharmacyERP/internal/role/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	Code        string `json:"code"        binding:"required,min=2,max=50"`
	Name        string `json:"name"        binding:"required,min=1,max=50"`
	Description string `json:"description"`
	Status      *int16 `json:"status" binding:"omitempty,oneof=0 1"`
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	Name        string `json:"name"        binding:"omitempty,min=1,max=50"`
	Description string `json:"description"`
	Status      *int16 `json:"status" binding:"omitempty,oneof=0 1"`
}

// RoleListResult 分页角色列表
type RoleListResult struct {
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
	List       []model.Role `json:"list"`
}

// RoleService 角色业务逻辑接口
type RoleService interface {
	GetRole(ctx context.Context, id int64) (*model.Role, error)
	ListRoles(ctx context.Context, filter repository.RoleListFilter) (*RoleListResult, error)
	CreateRole(ctx context.Context, req *CreateRoleReq) (*model.Role, error)
	UpdateRole(ctx context.Context, id int64, req *UpdateRoleReq) (*model.Role, error)
	DeleteRole(ctx context.Context, id int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionCodes []string) error
	ListPermissions(ctx context.Context) ([]model.Permission, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	log      *zap.Logger
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository, log *zap.Logger) RoleService {
	return &roleService{roleRepo: roleRepo, log: log}
}

// GetRole 获取角色详情
func (s *roleService) GetRole(ctx context.Context, id int64) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrRoleNotFound
		}
		return nil, ecode.ErrSystem
	}
	return role, nil
}

// ListRoles 分页查询角色列表
func (s *roleService) ListRoles(ctx context.Context, filter repository.RoleListFilter) (*RoleListResult, error) {
	roles, total, err := s.roleRepo.List(ctx, filter)
	if err != nil {
		return nil, ecode.ErrSystem
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	return &RoleListResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		List:       roles,
	}, nil
}

// CreateRole 创建新角色
func (s *roleService) CreateRole(ctx context.Context, req *CreateRoleReq) (*model.Role, error) {
	// 检查 code 是否已存在
	_, err := s.roleRepo.FindByCode(ctx, req.Code)
	if err == nil {
		return nil, ecode.ErrConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ecode.ErrSystem
	}

	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}

	role := &model.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		BuiltIn:     false,
		Status:      status,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		s.log.Error("创建角色失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}
	return role, nil
}

// UpdateRole 更新角色（内置角色不允许修改）
func (s *roleService) UpdateRole(ctx context.Context, id int64, req *UpdateRoleReq) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrRoleNotFound
		}
		return nil, ecode.ErrSystem
	}

	if role.BuiltIn {
		return nil, ecode.ErrBuiltInRole
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	role.Description = req.Description
	if req.Status != nil {
		role.Status = *req.Status
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		s.log.Error("更新角色失败", zap.Int64("role_id", id), zap.Error(err))
		return nil, ecode.ErrSystem
	}
	return role, nil
}

// DeleteRole 删除角色（内置角色和被使用的角色不允许删除）
func (s *roleService) DeleteRole(ctx context.Context, id int64) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrRoleNotFound
		}
		return ecode.ErrSystem
	}

	if role.BuiltIn {
		return ecode.ErrBuiltInRole
	}

	// 检查是否有用户使用该角色
	inUse, err := s.roleRepo.IsRoleInUse(ctx, id)
	if err != nil {
		return ecode.ErrSystem
	}
	if inUse {
		return ecode.ErrRoleInUse
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		s.log.Error("删除角色失败", zap.Int64("role_id", id), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// GetRolePermissions 获取角色权限列表
func (s *roleService) GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrRoleNotFound
		}
		return nil, ecode.ErrSystem
	}

	perms, err := s.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, ecode.ErrSystem
	}
	return perms, nil
}

// SetRolePermissions 设置角色权限（全量替换）
func (s *roleService) SetRolePermissions(ctx context.Context, roleID int64, permissionCodes []string) error {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrRoleNotFound
		}
		return ecode.ErrSystem
	}

	if role.BuiltIn {
		return ecode.ErrBuiltInRole
	}

	permIDs := make([]int64, 0, len(permissionCodes))
	for _, code := range permissionCodes {
		perm, err := s.roleRepo.FindPermissionByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrParamInvalid
			}
			return ecode.ErrSystem
		}
		permIDs = append(permIDs, perm.ID)
	}

	if err := s.roleRepo.SetRolePermissions(ctx, roleID, permIDs); err != nil {
		s.log.Error("设置角色权限失败", zap.Int64("role_id", roleID), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// ListPermissions 查询所有权限
func (s *roleService) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	perms, err := s.roleRepo.ListPermissions(ctx)
	if err != nil {
		return nil, ecode.ErrSystem
	}
	return perms, nil
}
