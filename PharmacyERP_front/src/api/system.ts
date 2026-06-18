// 系统管理 API
import { get, post, put, del } from './request'
import type {
  SysUser,
  SysRole,
  SysPermission,
  LoginLog,
  OperationLog,
  SecurityEvent,
  DataChangeLog,
  CreateUserRequest,
  UpdateUserRequest,
  CreateRoleRequest,
  UserListQuery,
  RoleListQuery,
} from '@/types/system'
import type { PageResponse } from '@/types/common'

// ===== 用户管理 =====
export function getUserList(params: UserListQuery) {
  return get<PageResponse<SysUser>>('/users', params as Record<string, unknown>)
}

export function getUserDetail(id: number) {
  return get<SysUser>(`/users/${id}`)
}

export function createUser(data: CreateUserRequest) {
  return post<SysUser>('/users', data)
}

export function updateUser(id: number, data: UpdateUserRequest) {
  return put<SysUser>(`/users/${id}`, data)
}

export function deleteUser(id: number) {
  return del<null>(`/users/${id}`)
}

export function toggleUserStatus(id: number, status: 0 | 1) {
  return put<SysUser>(`/users/${id}/status`, { status })
}

export function resetUserPassword(id: number, new_password: string) {
  return post<null>(`/users/${id}/reset-password`, { new_password })
}

// 分配用户角色（PUT /users/{id}/roles）：使用 role_codes（角色编码字符串），全量覆盖
export function assignUserRoles(id: number, role_codes: string[]) {
  return put<{ user_id: number; username: string; role_codes: string[] }>(`/users/${id}/roles`, { role_codes })
}

// ===== 角色管理 =====
export function getRoleList(params: RoleListQuery) {
  return get<PageResponse<SysRole>>('/roles', params as Record<string, unknown>)
}

export function getRoleDetail(id: number) {
  return get<SysRole>(`/roles/${id}`)
}

export function createRole(data: CreateRoleRequest) {
  return post<SysRole>('/roles', data)
}

export function updateRole(id: number, data: Partial<CreateRoleRequest>) {
  return put<SysRole>(`/roles/${id}`, data)
}

export function deleteRole(id: number) {
  return del<null>(`/roles/${id}`)
}

// 查询角色权限
export function getRolePermissions(roleId: number) {
  return get<SysPermission[]>(`/roles/${roleId}/permissions`)
}

// 分配角色权限（PUT /roles/{id}/permissions）：使用 permission_codes（权限编码字符串），全量覆盖
export function assignRolePermissions(roleId: number, permission_codes: string[]) {
  return put<null>(`/roles/${roleId}/permissions`, { permission_codes })
}

// ===== 权限管理（只读字典）=====
export function getPermissionList() {
  return get<SysPermission[]>('/permissions')
}

// ===== 审计日志 =====
// 登录日志（GET /audit/login-logs）支持 username, success
export function getLoginLogs(params: { page?: number; page_size?: number; username?: string; success?: boolean }) {
  return get<PageResponse<LoginLog>>('/audit/login-logs', params as Record<string, unknown>)
}

// 操作日志（GET /audit/operation-logs）支持 operator_name, module, action, start_time, end_time
export function getOperationLogs(params: { page?: number; page_size?: number; operator_name?: string; module?: string; action?: string; start_time?: string; end_time?: string }) {
  return get<PageResponse<OperationLog>>('/audit/operation-logs', params as Record<string, unknown>)
}

// 安全事件（GET /audit/security-events）支持 severity, event_type
export function getSecurityEvents(params: { page?: number; page_size?: number; event_type?: string; severity?: string }) {
  return get<PageResponse<SecurityEvent>>('/audit/security-events', params as Record<string, unknown>)
}

// 数据变更日志（GET /audit/data-change-logs）支持 table_name, record_id, change_type
export function getDataChangeLogs(params: { page?: number; page_size?: number; table_name?: string; record_id?: string; change_type?: string }) {
  return get<PageResponse<DataChangeLog>>('/audit/data-change-logs', params as Record<string, unknown>)
}
