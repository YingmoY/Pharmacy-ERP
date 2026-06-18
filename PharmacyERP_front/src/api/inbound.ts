// 入库管理 API
import { get, post, put, del } from './request'
import type {
  InboundOrder,
  InboundOrderDetail,
  CreateInboundOrderRequest,
  UpdateInboundOrderRequest,
  CreateInboundDetailRequest,
  InboundOrderListQuery,
  ConfirmTraceCodeRequest,
  ConfirmTraceCodeResult,
} from '@/types/inbound'
import type { PageResponse } from '@/types/common'

// 获取入库单列表
export function getInboundOrderList(params: InboundOrderListQuery) {
  return get<PageResponse<InboundOrder>>('/inbound-orders', params as Record<string, unknown>)
}

// 获取入库单详情
export function getInboundOrderDetail(id: number) {
  return get<InboundOrder>(`/inbound-orders/${id}`)
}

// 创建入库单
export function createInboundOrder(data: CreateInboundOrderRequest) {
  return post<InboundOrder>('/inbound-orders', data)
}

// 更新入库单
export function updateInboundOrder(id: number, data: UpdateInboundOrderRequest) {
  return put<InboundOrder>(`/inbound-orders/${id}`, data)
}

// 提交入库单（DRAFT -> PENDING_CONFIRM）
export function submitInboundOrder(id: number) {
  return post<InboundOrder>(`/inbound-orders/${id}/submit`)
}

// 完成入库单（PENDING_CONFIRM -> COMPLETED）
export function completeInboundOrder(id: number) {
  return post<InboundOrder>(`/inbound-orders/${id}/complete`)
}

// 取消入库单
export function cancelInboundOrder(id: number, reason?: string) {
  return post<InboundOrder>(`/inbound-orders/${id}/cancel`, { reason })
}

// 获取入库明细列表
export function getInboundDetailList(orderId: number) {
  return get<InboundOrderDetail[]>(`/inbound-orders/${orderId}/details`)
}

// 新增入库明细
export function addInboundDetail(orderId: number, data: CreateInboundDetailRequest) {
  return post<InboundOrderDetail>(`/inbound-orders/${orderId}/details`, data)
}

// 修改入库明细
export function updateInboundDetail(orderId: number, detailId: number, data: Partial<CreateInboundDetailRequest>) {
  return put<InboundOrderDetail>(`/inbound-orders/${orderId}/details/${detailId}`, data)
}

// 删除入库明细
export function deleteInboundDetail(orderId: number, detailId: number) {
  return del<null>(`/inbound-orders/${orderId}/details/${detailId}`)
}

// 追溯码扫码确认
export function confirmTraceCode(orderId: number, data: ConfirmTraceCodeRequest) {
  return post<ConfirmTraceCodeResult>(`/inbound-orders/${orderId}/confirm-trace`, data)
}
