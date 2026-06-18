package service

import (
	"context"
	"errors"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	rolemodel "github.com/YingmoY/PharmacyERP/internal/role/model"
	rolerepo "github.com/YingmoY/PharmacyERP/internal/role/repository"
	"github.com/YingmoY/PharmacyERP/internal/user/model"
	"github.com/YingmoY/PharmacyERP/internal/user/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username  string   `json:"username" binding:"required,min=3,max=50"`
	Password  string   `json:"password" binding:"required,min=8"`
	RealName  string   `json:"real_name" binding:"required,max=50"`
	RoleCodes []string `json:"role_codes" binding:"required,min=1"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Status    *int16   `json:"status" binding:"omitempty,oneof=0 1"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	RealName  string `json:"real_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Remark    string `json:"remark"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserWithRoles 用户信息及角色
type UserWithRoles struct {
	User  model.UserDTO    `json:"user"`
	Roles []rolemodel.Role `json:"roles"`
}

// UserListResult 分页用户列表
type UserListResult struct {
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
	List       []model.UserDTO `json:"list"`
}

// UserService 用户业务逻辑接口
type UserService interface {
	GetUser(ctx context.Context, id int64) (*UserWithRoles, error)
	ListUsers(ctx context.Context, filter repository.ListFilter) (*UserListResult, error)
	CreateUser(ctx context.Context, req *CreateUserReq) (*model.UserDTO, error)
	UpdateUser(ctx context.Context, id int64, req *UpdateUserReq) (*model.UserDTO, error)
	UpdateStatus(ctx context.Context, id int64, status int16) error
	ResetPassword(ctx context.Context, id int64, newPassword string) error
	GetUserRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error)
	SetUserRoles(ctx context.Context, userID int64, roleCodes []string) error
	GetUserPermissions(ctx context.Context, userID int64) ([]rolemodel.Permission, error)
}

type userService struct {
	userRepo repository.UserRepository
	roleRepo rolerepo.RoleRepository
	log      *zap.Logger
}

// NewUserService 创建用户服务实例
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo rolerepo.RoleRepository,
	log *zap.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		log:      log,
	}
}

// GetUser 获取用户详情（含角色）
func (s *userService) GetUser(ctx context.Context, id int64) (*UserWithRoles, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrUserNotFound
		}
		return nil, ecode.ErrSystem
	}

	roles, err := s.userRepo.GetUserRoles(ctx, id)
	if err != nil {
		s.log.Error("获取用户角色失败", zap.Int64("user_id", id), zap.Error(err))
		roles = []rolemodel.Role{}
	}

	return &UserWithRoles{
		User:  user.ToDTO(),
		Roles: roles,
	}, nil
}

// ListUsers 分页查询用户列表
func (s *userService) ListUsers(ctx context.Context, filter repository.ListFilter) (*UserListResult, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, ecode.ErrSystem
	}

	userIDs := make([]int64, 0, len(users))
	for i := range users {
		userIDs = append(userIDs, users[i].ID)
	}
	rolesMap, err := s.userRepo.GetUserRolesBatch(ctx, userIDs)
	if err != nil {
		s.log.Error("批量获取用户角色失败", zap.Error(err))
		rolesMap = map[int64][]model.RoleSimple{}
	}

	dtos := make([]model.UserDTO, 0, len(users))
	for i := range users {
		dto := users[i].ToDTO()
		if roles, ok := rolesMap[users[i].ID]; ok {
			dto.Roles = roles
		} else {
			dto.Roles = []model.RoleSimple{}
		}
		dtos = append(dtos, dto)
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	return &UserListResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		List:       dtos,
	}, nil
}

// CreateUser 创建新用户
func (s *userService) CreateUser(ctx context.Context, req *CreateUserReq) (*model.UserDTO, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, ecode.ErrDuplicateUsername
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ecode.ErrSystem
	}

	// 生成密码哈希，cost=10
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		s.log.Error("生成密码哈希失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       status,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("创建用户失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	roleIDs := make([]int64, 0, len(req.RoleCodes))
	for _, code := range req.RoleCodes {
		role, err := s.roleRepo.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ecode.ErrParamInvalid
			}
			return nil, ecode.ErrSystem
		}
		roleIDs = append(roleIDs, role.ID)
	}
	if err := s.userRepo.SetUserRoles(ctx, user.ID, roleIDs); err != nil {
		s.log.Error("set initial user roles failed", zap.Int64("user_id", user.ID), zap.Error(err))
		return nil, ecode.ErrSystem
	}

	dto := user.ToDTO()
	return &dto, nil
}

// UpdateUser 更新用户基本信息
func (s *userService) UpdateUser(ctx context.Context, id int64, req *UpdateUserReq) (*model.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrUserNotFound
		}
		return nil, ecode.ErrSystem
	}

	// 更新可修改的字段
	user.RealName = req.RealName
	user.Phone = req.Phone
	user.Email = req.Email
	user.AvatarURL = req.AvatarURL
	user.Remark = req.Remark

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Error("更新用户失败", zap.Int64("user_id", id), zap.Error(err))
		return nil, ecode.ErrSystem
	}

	dto := user.ToDTO()
	return &dto, nil
}

// UpdateStatus 更新用户状态（启用/禁用）
func (s *userService) UpdateStatus(ctx context.Context, id int64, status int16) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrUserNotFound
		}
		return ecode.ErrSystem
	}

	if err := s.userRepo.UpdateStatus(ctx, id, status); err != nil {
		s.log.Error("更新用户状态失败", zap.Int64("user_id", id), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// ResetPassword 重置用户密码（管理员操作，不需要旧密码）
func (s *userService) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrUserNotFound
		}
		return ecode.ErrSystem
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		s.log.Error("生成密码哈希失败", zap.Error(err))
		return ecode.ErrSystem
	}

	if err := s.userRepo.UpdatePassword(ctx, id, string(hash)); err != nil {
		s.log.Error("重置密码失败", zap.Int64("user_id", id), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// GetUserRoles 获取用户角色列表
func (s *userService) GetUserRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrUserNotFound
		}
		return nil, ecode.ErrSystem
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, ecode.ErrSystem
	}
	return roles, nil
}

// SetUserRoles 设置用户角色（全量替换）
func (s *userService) SetUserRoles(ctx context.Context, userID int64, roleCodes []string) error {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrUserNotFound
		}
		return ecode.ErrSystem
	}

	roleIDs := make([]int64, 0, len(roleCodes))
	for _, code := range roleCodes {
		role, err := s.roleRepo.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecode.ErrParamInvalid
			}
			return ecode.ErrSystem
		}
		roleIDs = append(roleIDs, role.ID)
	}

	if err := s.userRepo.SetUserRoles(ctx, userID, roleIDs); err != nil {
		s.log.Error("设置用户角色失败", zap.Int64("user_id", userID), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// GetUserPermissions 获取用户所有权限
func (s *userService) GetUserPermissions(ctx context.Context, userID int64) ([]rolemodel.Permission, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrUserNotFound
		}
		return nil, ecode.ErrSystem
	}

	perms, err := s.userRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, ecode.ErrSystem
	}
	return perms, nil
}
