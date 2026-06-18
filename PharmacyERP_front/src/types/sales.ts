// 销售相关类型

// 销售单状态
export type SalesOrderStatus =
  | 'PENDING'
  | 'PENDING_REVIEW'
  | 'APPROVED'
  | 'COMPLETED'
  | 'PARTIALLY_REFUNDED'
  | 'REFUNDED'
  | 'CANCELLED'

// 退货状态
export type RefundStatus = 'NONE' | 'REFUNDED'

// 支付方式
export type PaymentMethod = 'CASH' | 'ALIPAY' | 'WECHAT' | 'BANK_CARD' | 'MEDICARE'

// 退货模式
export type RefundMode = 'FULL' | 'PARTIAL'

// 销售单
export interface SalesOrder {
  id: number
  order_no: string
  status: SalesOrderStatus
  need_audit: boolean
  need_medicare: boolean
  is_prescription: boolean
  total_amount: string
  discount_amount: string
  actual_amount: string
  refund_amount: string
  medicare_amount: string
  personal_amount: string
  payment_method?: PaymentMethod
  medicare_transaction_id?: string
  mdtrt_id?: string
  cashier_id: number
  cashier_name: string
  created_at: string
  paid_at?: string
  cancelled_at?: string
  refunded_at?: string
  items?: SalesOrderItem[]
  review?: PharmacistReview
}

// 销售明细
export interface SalesOrderItem {
  id: number
  order_id: number
  drug_id: number
  drug_name: string
  specification: string
  manufacturer: string
  trace_code?: string
  batch_number?: string
  expire_date?: string
  location_code?: string
  quantity: number
  unit_price: string
  refund_status: RefundStatus
  refund_amount?: string
  refunded_at?: string
  refund_reason?: string
}

// 药师审核记录
export interface PharmacistReview {
  id: number
  review_no: string
  order_id: number
  status: 'PENDING' | 'APPROVED' | 'REJECTED' | 'CANCELLED'
  submitter_id: number
  submitter_name: string
  submitted_at: string
  pharmacist_id?: number
  pharmacist_name?: string
  reviewed_at?: string
  review_opinion?: string
}

// 创建销售单请求
export interface CreateSalesOrderRequest {
  is_prescription?: boolean
  items: CreateSalesItemRequest[]
}

// 创建销售明细
export interface CreateSalesItemRequest {
  drug_id: number
  trace_code?: string
}

// 结算请求
export interface SettleSalesOrderRequest {
  payment_method: PaymentMethod
  actual_amount: string
  discount_amount?: string
  // Medicare fields — only required when use_medicare is true
  use_medicare?: boolean
  med_type?: string
  insutype?: string
  acct_used_flag?: string
  // Patient identity — required when use_medicare is true
  mdtrt_cert_no?: string
  psn_no?: string
  psn_name?: string
}

// 医保预查询请求
export interface MedicarePreviewRequest {
  mdtrt_cert_no: string
  med_type?: string
  insutype?: string
  acct_used_flag?: string
}

// 医保预查询响应（费用分摊明细）
export interface MedicarePreviewResponse {
  psn_no: string
  psn_name: string
  total_amount: number
  fund_pay: number      // 统筹基金支付
  acct_pay: number      // 个人账户支付
  personal_cash: number // 个人现金支付
}

// 退货请求
export interface RefundSalesOrderRequest {
  refund_mode: RefundMode
  refund_reason: string
  detail_ids?: number[]
}

// 销售单列表查询
export interface SalesOrderListQuery {
  page?: number
  page_size?: number
  order_no?: string
  status?: SalesOrderStatus
  cashier_id?: number
  payment_method?: PaymentMethod
  created_start?: string
  created_end?: string
  need_audit?: boolean
  has_refund?: boolean
}

// 追溯码预占
export interface TraceReservation {
  id: number
  reservation_no: string
  order_id: number
  item_id: number
  trace_code: string
  status: 'RESERVED' | 'RELEASED' | 'CONSUMED' | 'EXPIRED'
  expire_at: string
  created_at: string
  released_at?: string
  confirmed_at?: string
}
