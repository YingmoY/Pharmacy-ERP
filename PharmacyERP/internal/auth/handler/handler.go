package handler

import (
	"github.com/YingmoY/PharmacyERP/internal/auth/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	rolemodel "github.com/YingmoY/PharmacyERP/internal/role/model"
	usermodel "github.com/YingmoY/PharmacyERP/internal/user/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 认证相关 HTTP 处理器
type Handler struct {
	authSvc service.AuthService
	log     *zap.Logger
}

// NewHandler 创建认证处理器实例
func NewHandler(authSvc service.AuthService, log *zap.Logger) *Handler {
	return &Handler{authSvc: authSvc, log: log}
}

// LoginReq 登录请求体
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordReq 修改密码请求体
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type userInfoResponse struct {
	ID              int64    `json:"id"`
	Username        string   `json:"username"`
	RealName        string   `json:"real_name"`
	Status          int16    `json:"status"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Roles           []string `json:"roles"`
	PermissionCodes []string `json:"permission_codes"`
}

func buildUserInfo(user usermodel.UserDTO, roles []rolemodel.Role, permissions []rolemodel.Permission) userInfoResponse {
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		permissionCodes = append(permissionCodes, permission.Code)
	}

	return userInfoResponse{
		ID:              user.ID,
		Username:        user.Username,
		RealName:        user.RealName,
		Status:          user.Status,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Roles:           roleCodes,
		PermissionCodes: permissionCodes,
	}
}

// Login 用户登录接口
// POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	result, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 400, e.Code, e.Msg)
		return
	}

	core.Success(c, gin.H{
		"token":      result.Token,
		"expires_in": result.ExpiresIn,
		"token_type": result.TokenType,
		"user":       buildUserInfo(result.User, result.Roles, result.Permissions),
	})
}

// Logout 用户登出接口
// POST /api/v1/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.FailWithStatus(c, 401, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	if err := h.authSvc.Logout(c.Request.Context(), userID); err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 500, e.Code, e.Msg)
		return
	}

	core.Success(c, nil)
}

// Me 获取当前登录用户信息
// GET /api/v1/auth/me
func (h *Handler) Me(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.FailWithStatus(c, 401, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	userDTO, roles, permissions, err := h.authSvc.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 400, e.Code, e.Msg)
		return
	}

	core.Success(c, buildUserInfo(*userDTO, roles, permissions))
}

// ChangePassword 修改当前用户密码
// PUT /api/v1/auth/password
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		core.FailWithStatus(c, 401, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	if err := h.authSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 400, e.Code, e.Msg)
		return
	}

	core.Success(c, nil)
}
