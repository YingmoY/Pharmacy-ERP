// 库存与追溯 API
import { get, post } from './request'
import type {
  DrugTraceInventory,
  DrugTraceLog,
  TraceCodeValidation,
  InventoryListQuery,
  NearExpireInventory,
  InventoryAdjustment,
  ShelvingRequest,
  BatchShelvingRequest,
  ShelvingResult,
  RelocateRequest,
} from '@/types/inventory'
import type { PageResponse } from '@/types/common'

// 获取库存列表
export function getInventoryList(params: InventoryListQuery) {
  return get<PageResponse<DrugTraceInventory>>('/inventory', params as Record<string, unknown>)
}

// 获取追溯码详情（/trace/{trace_code}）
export function getTraceCodeDetail(traceCode: string) {
  return get<DrugTraceInventory>(`/trace/${traceCode}`)
}

// 验证追溯码（POST /trace/validate）
export function validateTraceCode(data: { trace_code: string; inbound_order_id?: number }) {
  return post<TraceCodeValidation>('/trace/validate', data)
}

// 获取追溯日志（/trace/{trace_code}/logs）
export function getTraceLog(traceCode: string, params?: { page?: number; page_size?: number }) {
  return get<PageResponse<DrugTraceLog>>(`/trace/${traceCode}/logs`, params as Record<string, unknown>)
}

// 获取待上架列表
export function getPendingShelvingList(params: { page?: number; page_size?: number; drug_id?: number; inbound_order_id?: number }) {
  return get<PageResponse<DrugTraceInventory>>('/shelving/pending', params as Record<string, unknown>)
}

// 单个上架（/shelving/scan）
export function shelveTraceCode(data: ShelvingRequest) {
  return post<ShelvingResult>('/shelving/scan', data)
}

// 批量上架
export function batchShelveTraceCodes(data: BatchShelvingRequest) {
  return post<{ success_count: number; fail_count: number; results: ShelvingResult[] }>('/shelving/batch', data)
}

// 货位调拨
export function relocateTraceCode(data: RelocateRequest) {
  return post<{ success: boolean; message: string }>('/shelving/relocate', data)
}

// 货位混放检查（GET /shelving/mix-check?location_code=...）
export function checkMixPlacement(location_code: string) {
  return get<import('@/types/location').MixCheckResult>('/shelving/mix-check', { location_code } as Record<string, unknown>)
}

// 获取近效期库存
export function getNearExpireInventory(params: { page?: number; page_size?: number; days?: number; drug_id?: number }) {
  return get<PageResponse<NearExpireInventory>>('/inventory/near-expire', params as Record<string, unknown>)
}

// 获取库存调整记录（/inventory-adjustments）
export function getInventoryAdjustments(params: {
  page?: number
  page_size?: number
  trace_code?: string
  adjust_type?: string
  start?: string
  end?: string
}) {
  return get<PageResponse<InventoryAdjustment>>('/inventory-adjustments', params as Record<string, unknown>)
}
