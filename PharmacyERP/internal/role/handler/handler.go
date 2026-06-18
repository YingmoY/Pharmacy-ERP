package handler

import (
	"strconv"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/role/repository"
	"github.com/YingmoY/PharmacyERP/internal/role/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 角色管理 HTTP 处理器
type Handler struct {
	roleSvc service.RoleService
	log     *zap.Logger
}

// NewHandler 创建角色处理器实例
func NewHandler(roleSvc service.RoleService, log *zap.Logger) *Handler {
	return &Handler{roleSvc: roleSvc, log: log}
}

// ListRoles 分页查询角色列表
// GET /api/v1/roles
func (h *Handler) ListRoles(c *gin.Context) {
	filter := repository.RoleListFilter{
		Name: c.Query("name"),
	}

	if statusStr := c.Query("status"); statusStr != "" {
		v, err := strconv.ParseInt(statusStr, 10, 16)
		if err == nil {
			s := int16(v)
			filter.Status = &s
		}
	}

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

	result, err := h.roleSvc.ListRoles(c.Request.Context(), filter)
	if err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 500, e.Code, e.Msg)
		return
	}

	core.Success(c, result)
}

// CreateRole 创建角色
// POST /api/v1/roles
func (h *Handler) CreateRole(c *gin.Context) {
	var req service.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	role, err := h.roleSvc.CreateRole(c.Request.Context(), &req)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, role)
}

// GetRole 获取角色详情
// GET /api/v1/roles/:id
func (h *Handler) GetRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid role id")
		return
	}

	role, err := h.roleSvc.GetRole(c.Request.Context(), id)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, role)
}

// UpdateRole 更新角色
// PUT /api/v1/roles/:id
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid role id")
		return
	}

	var req service.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	role, err := h.roleSvc.UpdateRole(c.Request.Context(), id, &req)
	if err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, role)
}

// DeleteRole 删除角色
// DELETE /api/v1/roles/:id
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid role id")
		return
	}

	if err := h.roleSvc.DeleteRole(c.Request.Context(), id); err != nil {
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

// GetRolePermissions 获取角色权限列表
// GET /api/v1/roles/:id/permissions
func (h *Handler) GetRolePermissions(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid role id")
		return
	}

	perms, err := h.roleSvc.GetRolePermissions(c.Request.Context(), id)
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

// SetRolePermissions 设置角色权限
// PUT /api/v1/roles/:id/permissions
func (h *Handler) SetRolePermissions(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, "invalid role id")
		return
	}

	var req struct {
		PermissionCodes []string `json:"permission_codes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		core.FailWithStatus(c, 400, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	if err := h.roleSvc.SetRolePermissions(c.Request.Context(), id, req.PermissionCodes); err != nil {
		e := ecode.FromError(err)
		status := 400
		if e.Code == ecode.ErrSystem.Code {
			status = 500
		}
		core.FailWithStatus(c, status, e.Code, e.Msg)
		return
	}

	core.Success(c, gin.H{
		"role_id":          id,
		"permission_codes": req.PermissionCodes,
	})
}

// ListPermissions 查询所有权限
// GET /api/v1/permissions
func (h *Handler) ListPermissions(c *gin.Context) {
	perms, err := h.roleSvc.ListPermissions(c.Request.Context())
	if err != nil {
		e := ecode.FromError(err)
		core.FailWithStatus(c, 500, e.Code, e.Msg)
		return
	}

	core.Success(c, perms)
}

// parseIDParam 解析路由参数中的 ID
func parseIDParam(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}
