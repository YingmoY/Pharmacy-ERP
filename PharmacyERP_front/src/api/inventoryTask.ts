// 盘库管理 API
import { get, post } from './request'
import type {
  InventoryTask,
  InventoryTaskDetail,
  CreateInventoryTaskRequest,
  InventoryScanRequest,
  InventoryScanResult,
  InventoryTaskListQuery,
  LossCandidateItem,
  MisplacedItem,
} from '@/types/inventoryTask'
import type { PageResponse } from '@/types/common'

// 获取盘库任务列表
export function getInventoryTaskList(params: InventoryTaskListQuery) {
  return get<PageResponse<InventoryTask>>('/inventory-tasks', params as Record<string, unknown>)
}

// 获取盘库任务详情
export function getInventoryTaskDetail(id: number) {
  return get<InventoryTask>(`/inventory-tasks/${id}`)
}

// 创建盘库任务
export function createInventoryTask(data: CreateInventoryTaskRequest) {
  return post<InventoryTask>('/inventory-tasks', data)
}

// 开始盘库任务
export function startInventoryTask(id: number) {
  return post<InventoryTask>(`/inventory-tasks/${id}/start`)
}

// 完成盘库任务
export function completeInventoryTask(id: number) {
  return post<InventoryTask>(`/inventory-tasks/${id}/complete`)
}

// 取消盘库任务
export function cancelInventoryTask(id: number) {
  return post<InventoryTask>(`/inventory-tasks/${id}/cancel`)
}

// 盘库扫码
export function scanInventory(id: number, data: InventoryScanRequest) {
  return post<InventoryScanResult>(`/inventory-tasks/${id}/scan`, data)
}

// 获取盘库明细
export function getInventoryTaskDetails(id: number, params?: { scan_result?: string; page?: number; page_size?: number }) {
  return get<PageResponse<InventoryTaskDetail>>(`/inventory-tasks/${id}/details`, params as Record<string, unknown>)
}

// 获取盘亏候选列表
export function getLossCandidates(id: number, params?: { page?: number; page_size?: number }) {
  return get<PageResponse<LossCandidateItem>>(`/inventory-tasks/${id}/loss-candidates`, params as Record<string, unknown>)
}

// 确认盘亏候选（单条，需提供原因）
// POST /inventory-tasks/{id}/loss-candidates/{trace_code}/confirm
export function confirmLossCandidate(taskId: number, traceCode: string, reason: string) {
  return post<null>(`/inventory-tasks/${taskId}/loss-candidates/${encodeURIComponent(traceCode)}/confirm`, { reason })
}

// 驳回盘亏候选（单条，需提供原因）
// POST /inventory-tasks/{id}/loss-candidates/{trace_code}/reject
export function rejectLossCandidate(taskId: number, traceCode: string, reason: string) {
  return post<null>(`/inventory-tasks/${taskId}/loss-candidates/${encodeURIComponent(traceCode)}/reject`, { reason })
}

// 获取错架列表（分页）
export function getMisplacedItems(id: number, params?: { page?: number; page_size?: number }) {
  return get<PageResponse<MisplacedItem>>(`/inventory-tasks/${id}/misplaced`, params as Record<string, unknown>)
}

// 获取错架列表（全量，兼容旧接口）
export function getMisplacedList(id: number) {
  return get<MisplacedItem[]>(`/inventory-tasks/${id}/misplaced`)
}

// 处理错架并调整货位（使用 location_id，非 location_code）
// POST /inventory-tasks/{id}/misplaced/{trace_code}/relocate
export function fixMisplaced(taskId: number, traceCode: string, targetLocationId: number, reason?: string) {
  return post<null>(`/inventory-tasks/${taskId}/misplaced/${encodeURIComponent(traceCode)}/relocate`, {
    target_location_id: targetLocationId,
    reason,
  })
}
