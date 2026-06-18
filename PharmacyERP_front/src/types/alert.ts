// 预警与通知类型

export type AlertType = 'NEAR_EXPIRE' | 'LOW_STOCK' | 'LOSS_CANDIDATE' | 'MISPLACED'
export type AlertStatus = 'ACTIVE' | 'RESOLVED' | 'IGNORED'
export type AlertPriority = 'HIGH' | 'MEDIUM' | 'LOW'

// 预警信息（对应 AlertInfo schema）
export interface AlertInfo {
  id: number
  alert_type: AlertType
  status: AlertStatus
  priority: AlertPriority
  title: string
  content: string
  ref_id?: number | null
  ref_type?: string | null
  resolver_id?: number | null
  resolver_name?: string | null
  resolved_at?: string | null
  created_at: string
}

// 预警查询参数
export interface AlertListQuery {
  page?: number
  page_size?: number
  alert_type?: AlertType
  priority?: AlertPriority
  status?: AlertStatus
}

// 通知消息（对应 NotificationInfo schema）
export interface Notification {
  id: number
  title: string
  content: string
  notification_type: 'ALERT' | 'SYSTEM' | 'TASK'
  is_read: boolean
  ref_id?: number | null
  ref_type?: string | null
  created_at: string
}

// 通知查询
export interface NotificationListQuery {
  page?: number
  page_size?: number
  is_read?: boolean
  notification_type?: string
}
