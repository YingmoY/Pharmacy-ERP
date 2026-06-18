// 库存与追溯相关类型

// 追溯码库存状态
export type TraceInventoryStatus =
  | 'PENDING'
  | 'IN_STOCK'
  | 'SOLD'
  | 'MISPLACED'
  | 'LOSS_CANDIDATE'
  | 'LOST'

// 追溯码库存
export interface DrugTraceInventory {
  id: number
  trace_code: string
  drug_id: number
  drug_name: string
  specification: string
  manufacturer: string
  batch_number: string
  expire_date: string
  status: TraceInventoryStatus
  location_id?: number
  location_code?: string
  location_name?: string
  is_reserved: boolean
  inbound_order_id: number
  inbound_order_no: string
  last_action: string
  created_at: string
  updated_at: string
}

// 追溯日志
export interface DrugTraceLog {
  id: number
  trace_code: string
  drug_id: number
  drug_name: string
  action_type:
    | 'INBOUND'
    | 'SHELVING'
    | 'SALE'
    | 'RETURN'
    | 'INVENTORY'
    | 'RELOCATION'
    | 'LOSS'
  from_status?: TraceInventoryStatus
  to_status: TraceInventoryStatus
  from_location_id?: number
  from_location_code?: string
  to_location_id?: number
  to_location_code?: string
  operator_id: number
  operator_name: string
  related_no?: string
  order_id?: number
  remark?: string
  created_at: string
}

// 追溯码验证结果
export interface TraceCodeValidation {
  trace_code: string
  exists: boolean
  status?: TraceInventoryStatus
  drug_id?: number
  drug_name?: string
  batch_number?: string
  expire_date?: string
  location_id?: number
  location_code?: string
  is_available: boolean
  is_reserved: boolean
  reason?: string
}

// 库存查询参数
export interface InventoryListQuery {
  page?: number
  page_size?: number
  trace_code?: string
  drug_id?: number
  status?: TraceInventoryStatus
  location_id?: number
  batch_number?: string
  expire_start?: string
  expire_end?: string
  is_reserved?: boolean
  near_expire?: boolean
}

// 近效期库存
export interface NearExpireInventory {
  drug_id: number
  drug_name: string
  specification: string
  batch_number: string
  expire_date: string
  remaining_days: number
  location_code: string
  count: number
  alert_level: 'HIGH' | 'MEDIUM' | 'LOW'
}

// 库存调整记录
export interface InventoryAdjustment {
  id: number
  adjustment_no: string
  trace_code: string
  drug_name: string
  adjust_type: 'RELOCATE' | 'LOSS' | 'RETURN' | 'MANUAL'
  from_status?: TraceInventoryStatus
  to_status: TraceInventoryStatus
  from_location_code?: string
  to_location_code?: string
  reason: string
  operator_id: number
  operator_name: string
  created_at: string
}

// 上架请求
export interface ShelvingRequest {
  trace_code: string
  location_code: string
}

// 批量上架请求
export interface BatchShelvingRequest {
  items: ShelvingRequest[]
}

// 上架结果
export interface ShelvingResult {
  trace_code: string
  success: boolean
  location_code?: string
  drug_name?: string
  message: string
}

// 调拨请求
export interface RelocateRequest {
  trace_code: string
  location_code: string
  remark?: string
}
