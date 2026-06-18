// 认证相关 API
import { get, post, put } from './request'
import type { LoginRequest, LoginResponse, UserInfo } from '@/types/auth'

// 登录
export function login(data: LoginRequest) {
  return post<LoginResponse>('/auth/login', data)
}

// 获取当前用户信息（含权限）
export function getCurrentUser() {
  return get<UserInfo>('/auth/me')
}

// 登出
export function logout() {
  return post<null>('/auth/logout')
}

// 修改密码（PUT /auth/password）
export function changePassword(data: { old_password: string; new_password: string }) {
  return put<null>('/auth/password', data)
}
