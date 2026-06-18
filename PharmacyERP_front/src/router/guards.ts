// 路由守卫
import type { Router } from 'vue-router'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'

export function setupRouterGuards(router: Router) {
  // 前置守卫：验证登录状态和权限
  router.beforeEach(async (to, _from, next) => {
    const authStore = useAuthStore()
    const requiresAuth = to.meta.requiresAuth !== false

    // 不需要认证的页面（登录页）直接放行
    if (!requiresAuth) {
      // 已登录时访问登录页，跳转到首页
      if (authStore.isLoggedIn && to.path === '/login') {
        next('/dashboard')
        return
      }
      next()
      return
    }

    // 未登录，跳转登录页
    if (!authStore.isLoggedIn) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }

    // 已登录但未加载用户信息，加载一次
    if (!authStore.user) {
      await authStore.loadCurrentUser()
      // 加载失败（token 已过期）
      if (!authStore.user) {
        next({ path: '/login', query: { redirect: to.fullPath } })
        return
      }
    }

    // 权限检查
    const permission = to.meta.permission as string | undefined
    if (permission && !authStore.hasPermission(permission)) {
      const pageTitle = to.meta.title as string | undefined
      const hint = pageTitle ? `你无权访问"${pageTitle}"，需要 ${permission} 权限` : `你无权访问此内容，需要 ${permission} 权限`
      message.error(hint, 3)
      next(false)
      // 取消本次导航后回到上一页（若上一页存在）
      if (_from && _from.name) {
        router.back()
      } else {
        router.replace('/dashboard')
      }
      return
    }

    next()
  })

  // 后置守卫：更新页面标题
  router.afterEach((to) => {
    const title = to.meta.title as string | undefined
    document.title = title ? `${title} - 智慧药店ERP` : '智慧药店ERP'

    // 仅在需要认证的页面刷新角标，避免登录页无 token 时发出请求触发 10002 提示
    if (to.meta.requiresAuth !== false) {
      const appStore = useAppStore()
      appStore.refreshCounts()
    }
  })
}
