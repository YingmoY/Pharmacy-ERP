package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/user/repository"
	"github.com/YingmoY/PharmacyERP/internal/user/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 用户管理 HTTP 处理器
type Handler struct {
	userSvc service.UserService
	log     *zap.Logger
}

// NewHandler 创建用户处理器实例
func NewHandler(userSvc service.UserService, log *zap.Logger) *Handler {
	return &Handler{userSvc: userSvc, log: log}
}

// ListUsers 分页查询用户列表
// GET /api/v1/users
func (h *Handler) ListUsers(c *gin.Context) {
	filter := repository.ListFilter{
		Username: c.Query("username"),
		RealName: c.Query("real_name"),
	}

	// 解析 status 过滤参数
	if statusStr := c.Query("status"); statusStr != "" {
		v, err := strconv.ParseInt(statusStr, 10, 16)
		if err == nil {
			s := int16(v)
			filter.Status = &s
		}
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	filter.Page = page
	filter.PageSize = pageSize

	result, err := h.userSvc.ListUsers(c.Request.Context(), filter)
	if err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 500, e.Code, e.Msg)
		return
	}

	core.Success(c, result)
}

// CreateUser 创建用户
// POST /api/v1/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req service.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	user, err := h.userSvc.CreateUser(c.Request.Context(), &req)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, user)
}

// GetUser 获取用户详情
// GET /api/v1/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	result, err := h.userSvc.GetUser(c.Request.Context(), id)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, result)
}

// UpdateUser 更新用户信息
// PUT /api/v1/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	var req service.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	user, err := h.userSvc.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, user)
}

// UpdateStatus 更新用户状态
// PATCH /api/v1/users/:id/status
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	var req struct {
		Status *int16 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	if err := h.userSvc.UpdateStatus(c.Request.Context(), id, *req.Status); err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, nil)
}

// ResetPassword 重置用户密码（管理员操作）
// POST /api/v1/users/:id/reset-password
func (h *Handler) ResetPassword(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	var req service.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	if err := h.userSvc.ResetPassword(c.Request.Context(), id, req.NewPassword); err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, nil)
}

// GetUserRoles 获取用户角色列表
// GET /api/v1/users/:id/roles
func (h *Handler) GetUserRoles(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	roles, err := h.userSvc.GetUserRoles(c.Request.Context(), id)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, roles)
}

// SetUserRoles 设置用户角色
// PUT /api/v1/users/:id/roles
func (h *Handler) SetUserRoles(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	var req struct {
		RoleCodes []string `json:"role_codes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	if err := h.userSvc.SetUserRoles(c.Request.Context(), id, req.RoleCodes); err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, gin.H{
		"user_id":    id,
		"role_codes": req.RoleCodes,
	})
}

// GetUserPermissions 获取用户权限列表
// GET /api/v1/users/:id/permissions
func (h *Handler) GetUserPermissions(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid user id")
		return
	}

	perms, err := h.userSvc.GetUserPermissions(c.Request.Context(), id)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, perms)
}

// parseIDParam 解析路由参数中的 ID
func parseIDParam(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}
