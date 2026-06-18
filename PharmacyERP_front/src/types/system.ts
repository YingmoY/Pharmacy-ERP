// 系统管理相关类型

// 用户
export interface SysUser {
  id: number
  username: string
  real_name: string
  phone?: string
  email?: string
  status: number
  last_login_at?: string
  last_login_ip?: string
  created_at: string
  roles: UserRole[]
}

export interface UserRole {
  id: number
  code: string
  name: string
}

// 角色
export interface SysRole {
  id: number
  code: string
  name: string
  description?: string
  built_in: boolean
  status: number
  created_at: string
  permissions?: SysPermission[]
}

// 权限
export interface SysPermission {
  id: number
  code: string
  name: string
  resource: string
  action: string
  description?: string
  module?: string
  apis?: PermissionApi[]
}

// 权限关联API
export interface PermissionApi {
  id: number
  permission_id: number
  method: string
  path: string
}

// 登录日志
export interface LoginLog {
  id: number
  username: string
  success: boolean
  ip: string
  user_agent?: string
  message?: string
  created_at: string
}

// 操作日志
export interface OperationLog {
  id: number
  operator_id?: number
  operator_name?: string
  module: string
  action: string
  resource_type: string
  resource_id?: string
  before_data?: Record<string, unknown>
  after_data?: Record<string, unknown>
  ip?: string
  user_agent?: string
  created_at: string
}

// 安全事件
export interface SecurityEvent {
  id: number
  event_type: string
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL'
  description: string
  ip?: string
  created_at: string
}

// 数据变更日志
export interface DataChangeLog {
  id: number
  table_name: string
  record_id: string
  change_type: 'INSERT' | 'UPDATE' | 'DELETE'
  before_data?: Record<string, unknown>
  after_data?: Record<string, unknown>
  operator_name?: string
  created_at: string
}

// 创建用户请求
export interface CreateUserRequest {
  username: string
  password: string
  real_name: string
  role_codes: string[]
}

// 更新用户请求（只更新基础信息，角色通过 PUT /users/{id}/roles 单独维护）
export interface UpdateUserRequest {
  real_name?: string
}

// 创建角色请求
export interface CreateRoleRequest {
  code: string
  name: string
  description?: string
}

// 用户列表查询
export interface UserListQuery {
  page?: number
  page_size?: number
  username?: string
  real_name?: string
  status?: number
  role_id?: number
}

// 角色列表查询
export interface RoleListQuery {
  page?: number
  page_size?: number
  code?: string
  name?: string
  status?: number
}
