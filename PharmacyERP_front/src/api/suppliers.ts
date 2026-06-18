// 供应商 API
import { get, post, put, del } from './request'
import type { Supplier, CreateSupplierRequest, UpdateSupplierRequest, SupplierListQuery } from '@/types/supplier'
import type { PageResponse } from '@/types/common'

export function getSupplierList(params: SupplierListQuery) {
  return get<PageResponse<Supplier>>('/suppliers', params as Record<string, unknown>)
}

export function getSupplierDetail(id: number) {
  return get<Supplier>(`/suppliers/${id}`)
}

export function createSupplier(data: CreateSupplierRequest) {
  return post<Supplier>('/suppliers', data)
}

export function updateSupplier(id: number, data: UpdateSupplierRequest) {
  return put<Supplier>(`/suppliers/${id}`, data)
}

export function deleteSupplier(id: number) {
  return del<null>(`/suppliers/${id}`)
}

export function toggleSupplierStatus(id: number, status: 0 | 1) {
  return put<Supplier>(`/suppliers/${id}/status`, { status })
}
