// 扫码任务 API
import { get, post } from './request'
import type { PageResponse } from '@/types/common'

export type ScanTaskType = 'INBOUND' | 'SHELVING' | 'INVENTORY'
export type ScanTaskStatus = 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED'

export interface ScanTask {
  id: number
  task_no: string
  task_type: ScanTaskType
  status: ScanTaskStatus
  related_id: number
  related_no: string
  operator_id: number
  operator_name: string
  remark?: string
  start_time?: string
  end_time?: string
  created_at: string
}

export interface ScanTaskDetail {
  id: number
  task_id: number
  trace_code: string
  location_code?: string
  scan_result: 'SUCCESS' | 'DUPLICATE' | 'INVALID' | 'STATUS_ERROR'
  error_msg?: string
  scan_time: string
}

export interface CreateScanTaskRequest {
  task_type: ScanTaskType
  related_id: number
  remark?: string
}

export interface SubmitScanRequest {
  trace_code: string
  location_code?: string
  detail_id?: number
}

export interface ScanTaskListQuery {
  page?: number
  page_size?: number
  task_no?: string
  task_type?: ScanTaskType
  status?: ScanTaskStatus
  operator_id?: number
}

// 获取扫码任务列表
export function getScanTaskList(params: ScanTaskListQuery) {
  return get<PageResponse<ScanTask>>('/scan-tasks', params as Record<string, unknown>)
}

// 获取扫码任务详情
export function getScanTaskDetail(id: number) {
  return get<ScanTask>(`/scan-tasks/${id}`)
}

// 创建扫码任务
export function createScanTask(data: CreateScanTaskRequest) {
  return post<ScanTask>('/scan-tasks', data)
}

// 开始扫码任务
export function startScanTask(id: number) {
  return post<ScanTask>(`/scan-tasks/${id}/start`)
}

// 提交扫码
export function submitScan(id: number, data: SubmitScanRequest) {
  return post<ScanTaskDetail>(`/scan-tasks/${id}/submit`, data)
}

// 完成扫码任务
export function completeScanTask(id: number) {
  return post<ScanTask>(`/scan-tasks/${id}/complete`)
}

// 取消扫码任务
export function cancelScanTask(id: number) {
  return post<ScanTask>(`/scan-tasks/${id}/cancel`)
}

// 获取扫码明细（全量）
export function getScanTaskDetails(id: number) {
  return get<ScanTaskDetail[]>(`/scan-tasks/${id}/details`)
}

// 获取扫码明细（分页，页面使用）
export function getScanTaskItems(id: number, params?: { page?: number; page_size?: number }) {
  return get<PageResponse<ScanTaskDetail>>(`/scan-tasks/${id}/details`, params as Record<string, unknown>)
}

// 提交扫码（页面使用）
export function submitScanResult(id: number, data: SubmitScanRequest) {
  return post<{ success: boolean; message: string }>(`/scan-tasks/${id}/submit`, data)
}

// 指派任务
export function assignScanTask(id: number, operator_id: number) {
  return post<ScanTask>(`/scan-tasks/${id}/assign`, { operator_id })
}
