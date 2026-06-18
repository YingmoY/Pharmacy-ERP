// 预警与通知 API
import { get, post } from './request'
import type { AlertInfo, AlertListQuery, Notification, NotificationListQuery } from '@/types/alert'
import type { PageResponse } from '@/types/common'

// 获取预警列表
export function getAlertList(params: AlertListQuery) {
  return get<PageResponse<AlertInfo>>('/alerts', params as Record<string, unknown>)
}

// 获取预警详情
export function getAlertDetail(id: number) {
  return get<AlertInfo>(`/alerts/${id}`)
}

// 处理预警（remark 为处理备注）
export function resolveAlert(id: number, data: { remark: string }) {
  return post<null>(`/alerts/${id}/resolve`, data)
}

// 忽略预警（无请求体）
export function ignoreAlert(id: number) {
  return post<null>(`/alerts/${id}/ignore`)
}

// ===== 通知消息 =====

// 获取通知列表
export function getNotificationList(params: NotificationListQuery) {
  return get<PageResponse<Notification>>('/notifications', params as Record<string, unknown>)
}

// 标记已读
export function markNotificationRead(id: number) {
  return post<null>(`/notifications/${id}/read`)
}

// 全部标记已读
export function markAllNotificationsRead() {
  return post<null>('/notifications/read-all')
}

// 获取未读数量
export function getUnreadNotificationCount() {
  return get<{ count: number }>('/notifications/unread-count')
}
