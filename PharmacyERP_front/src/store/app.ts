// 应用全局状态
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUnreadNotificationCount, getAlertList } from '@/api/alert'
import { useAuthStore } from '@/store/auth'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const collapsed = ref(false)

  // 未读通知数
  const unreadNotificationCount = ref(0)

  // 活跃预警数
  const activeAlertCount = ref(0)

  // 切换侧边栏
  function toggleCollapsed() {
    collapsed.value = !collapsed.value
  }

  // 刷新未读数量（仅在已登录时执行，防止登录前触发 10002 错误提示）
  async function refreshCounts() {
    const authStore = useAuthStore()
    if (!authStore.isLoggedIn) return
    try {
      const [notifRes, alertRes] = await Promise.all([
        getUnreadNotificationCount(),
        getAlertList({ page: 1, page_size: 1, status: 'ACTIVE' }),
      ])
      unreadNotificationCount.value = notifRes.count
      activeAlertCount.value = alertRes.total
    } catch {
      // 忽略错误，不影响主流程
    }
  }

  return {
    collapsed,
    unreadNotificationCount,
    activeAlertCount,
    toggleCollapsed,
    refreshCounts,
  }
})
