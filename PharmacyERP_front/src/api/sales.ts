// 销售管理 API
import { get, post, del } from './request'
import type {
  SalesOrder,
  SalesOrderItem,
  CreateSalesOrderRequest,
  CreateSalesItemRequest,
  SettleSalesOrderRequest,
  RefundSalesOrderRequest,
  SalesOrderListQuery,
  TraceReservation,
  MedicarePreviewRequest,
  MedicarePreviewResponse,
} from '@/types/sales'
import type { PageResponse } from '@/types/common'

// 获取销售单列表
export function getSalesOrderList(params: SalesOrderListQuery) {
  return get<PageResponse<SalesOrder>>('/sales-orders', params as Record<string, unknown>)
}

// 获取销售单详情
export function getSalesOrderDetail(id: number) {
  return get<SalesOrder>(`/sales-orders/${id}`)
}

// 创建销售单
export function createSalesOrder(data: CreateSalesOrderRequest) {
  return post<SalesOrder>('/sales-orders', data)
}

// 添加销售明细（仅 PENDING 状态可用）
export function addSalesOrderItem(orderId: number, data: CreateSalesItemRequest) {
  return post<SalesOrderItem>(`/sales-orders/${orderId}/details`, data)
}

// 删除销售明细（仅 PENDING 状态可用）
export function deleteSalesOrderItem(orderId: number, itemId: number) {
  return del<null>(`/sales-orders/${orderId}/details/${itemId}`)
}

// 销售结算（PENDING 或 APPROVED 状态）
export function settleSalesOrder(orderId: number, data: SettleSalesOrderRequest) {
  return post<SalesOrder>(`/sales-orders/${orderId}/pay`, data)
}

// 取消销售单（PENDING / PENDING_REVIEW / APPROVED）
export function cancelSalesOrder(orderId: number, reason?: string) {
  return post<SalesOrder>(`/sales-orders/${orderId}/cancel`, { reason })
}

// 退货（COMPLETED / PARTIALLY_REFUNDED）
export function refundSalesOrder(orderId: number, data: RefundSalesOrderRequest) {
  return post<SalesOrder>(`/sales-orders/${orderId}/refund`, data)
}

// 获取销售明细列表
export function getSalesOrderItems(orderId: number) {
  return get<SalesOrderItem[]>(`/sales-orders/${orderId}/details`)
}

// 获取预占记录
export function getSalesReservations(orderId: number) {
  return get<TraceReservation[]>(`/sales-orders/${orderId}/reserved-traces`)
}

// 医保预查询（1101 患者身份验证 + 费用分摊预估）
export function getMedicarePreview(orderId: number, data: MedicarePreviewRequest) {
  return post<MedicarePreviewResponse>(`/sales-orders/${orderId}/medicare-preview`, data)
}

// 推荐可售追溯码（按近效期优先）
export function recommendTraceCode(drugId: number, quantity?: number) {
  return get<{ trace_codes: { trace_code: string; batch_number: string; expire_date: string; location_code: string }[] }>(
    `/inventory/recommend-sale`,
    { drug_id: drugId, quantity: quantity ?? 1 } as Record<string, unknown>,
  )
}
