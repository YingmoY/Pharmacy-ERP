// 药品基础资料 API
import { get, post, put, del } from './request'
import type { DrugInfo, CreateDrugRequest, UpdateDrugRequest, DrugListQuery } from '@/types/drug'
import type { PageResponse } from '@/types/common'

// 获取药品列表
export function getDrugList(params: DrugListQuery) {
  return get<PageResponse<DrugInfo>>('/drugs', params as Record<string, unknown>)
}

// 获取药品详情
export function getDrugDetail(id: number) {
  return get<DrugInfo>(`/drugs/${id}`)
}

// 创建药品
export function createDrug(data: CreateDrugRequest) {
  return post<DrugInfo>('/drugs', data)
}

// 更新药品
export function updateDrug(id: number, data: UpdateDrugRequest) {
  return put<DrugInfo>(`/drugs/${id}`, data)
}

// 删除药品
export function deleteDrug(id: number) {
  return del<null>(`/drugs/${id}`)
}

// 启用/停用药品
export function toggleDrugStatus(id: number, status: 0 | 1) {
  return put<DrugInfo>(`/drugs/${id}/status`, { status })
}

// 获取药品追溯码库存明细
export function getDrugInventoryList(drugId: number, params?: { page?: number; page_size?: number }) {
  return get<import('@/types/common').PageResponse<{
    id: number
    trace_code: string
    batch_number: string
    expire_date: string
    status: string
    location_code: string
    drug_name?: string
  }>>(`/inventory/drugs/${drugId}`, params as Record<string, unknown>)
}

// 获取药品库存汇总
export function getDrugInventorySummary(id: number) {
  return get<{
    in_stock: number
    available: number
    pending_shelving: number
    reserved: number
    near_expire: number
  }>(`/drugs/${id}/inventory-summary`)
}
