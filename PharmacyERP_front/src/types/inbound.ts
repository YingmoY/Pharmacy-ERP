// 入库相关类型

// 入库单状态
export type InboundOrderStatus = 'DRAFT' | 'PENDING_CONFIRM' | 'COMPLETED' | 'CANCELLED'

// 入库单
export interface InboundOrder {
  id: number
  order_no: string
  supplier_id: number
  supplier_name: string
  status: InboundOrderStatus
  invoice_no?: string
  total_amount: string
  remark?: string
  creator_id: number
  creator_name: string
  created_at: string
  submitted_at?: string
  completed_at?: string
  cancelled_at?: string
  details?: InboundOrderDetail[]
  total_planned_qty: number
  total_confirmed_qty: number
}

// 入库明细
export interface InboundOrderDetail {
  id: number
  order_id: number
  drug_id: number
  drug_name: string
  specification: string
  manufacturer: string
  batch_number: string
  expire_date: string
  planned_qty: number
  confirmed_qty: number
  unit_price: string
  amount: string
}

// 创建入库单请求
export interface CreateInboundOrderRequest {
  supplier_id: number
  invoice_no?: string
  remark?: string
  details?: CreateInboundDetailRequest[]
}

// 创建入库明细
export interface CreateInboundDetailRequest {
  drug_id: number
  batch_number: string
  expire_date: string
  planned_qty: number
  unit_price: number
}

// 更新入库单请求
export interface UpdateInboundOrderRequest {
  supplier_id?: number
  invoice_no?: string
  remark?: string
}

// 入库单列表查询
export interface InboundOrderListQuery {
  page?: number
  page_size?: number
  order_no?: string
  supplier_id?: number
  status?: InboundOrderStatus
  created_start?: string
  created_end?: string
  completed_start?: string
  completed_end?: string
}

// 追溯码扫码确认请求
export interface ConfirmTraceCodeRequest {
  detail_id: number
  trace_code: string
}

// 扫码确认结果
export interface ConfirmTraceCodeResult {
  success: boolean
  trace_code: string
  drug_name: string
  batch_number: string
  detail_id: number
  confirmed_qty: number
  planned_qty: number
  error_code?: string
  message: string
}
