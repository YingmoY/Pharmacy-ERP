// 认证相关类型

export interface LoginRequest {
  username: string
  password: string
  client_type?: 'WEB' | 'MOBILE'
}

// 登录成功返回的 data 字段（单 JWT，无 refresh_token）
export interface LoginResponse {
  token: string
  expires_in: number
  token_type: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  real_name: string
  status: number
  created_at?: string
  updated_at?: string
  roles?: string[]
  permission_codes?: string[]
}
