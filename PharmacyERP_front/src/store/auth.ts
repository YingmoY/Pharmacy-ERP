// 认证状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo } from '@/types/auth'
import { getCurrentUser, login as apiLogin, logout as apiLogout } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null)
  const accessToken = ref<string>(localStorage.getItem('access_token') || '')

  const isLoggedIn = computed(() => !!accessToken.value)

  // 当前用户权限点列表（来自 UserInfo.permission_codes）
  const permissions = computed(() => user.value?.permission_codes ?? [])

  function hasPermission(permCode: string): boolean {
    if (!user.value) return false
    return permissions.value.includes(permCode)
  }

  function hasAnyPermission(permCodes: string[]): boolean {
    return permCodes.some((code) => hasPermission(code))
  }

  // 登录（单 JWT，token 字段）
  async function login(username: string, password: string) {
    const res = await apiLogin({ username, password })
    accessToken.value = res.token
    user.value = res.user

    localStorage.setItem('access_token', res.token)
  }

  async function loadCurrentUser() {
    if (!accessToken.value) return
    try {
      user.value = await getCurrentUser()
    } catch {
      clearAuth()
    }
  }

  async function logout() {
    try {
      await apiLogout()
    } catch {
      // 忽略登出接口错误
    } finally {
      clearAuth()
    }
  }

  function clearAuth() {
    user.value = null
    accessToken.value = ''
    localStorage.removeItem('access_token')
  }

  return {
    user,
    accessToken,
    isLoggedIn,
    permissions,
    hasPermission,
    hasAnyPermission,
    login,
    logout,
    loadCurrentUser,
    clearAuth,
  }
})
