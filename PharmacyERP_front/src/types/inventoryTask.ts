// 盘库相关类型

export type InventoryTaskStatus = 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED'
export type InventoryTaskScopeType = 'AREA' | 'SHELF' | 'LOCATION'
export type ScanResultType = 'NORMAL' | 'MISPLACED_FOUND' | 'UNEXPECTED' | 'DUPLICATE'

// 盘库任务
export interface InventoryTask {
  id: number
  task_no: string
  scope_type: InventoryTaskScopeType
  scope_value: string
  status: InventoryTaskStatus
  creator_id: number
  creator_name: string
  assignee_id?: number
  assignee_name?: string
  remark?: string
  start_time?: string
  end_time?: string
  completed_at?: string
  created_at: string
  // 统计字段
  total_count?: number
  scanned_count?: number
  normal_count?: number
  misplaced_count?: number
  loss_candidate_count?: number
}

// 盘库明细
export interface InventoryTaskDetail {
  id: number
  task_id: number
  trace_code: string
  drug_id: number
  drug_name: string
  batch_number: string
  expire_date: string
  scanned_location_id?: number
  scanned_location_code?: string
  system_location_id?: number
  system_location_code?: string
  scan_result: ScanResultType
  scanned_at: string
}

// 创建盘库任务
export interface CreateInventoryTaskRequest {
  scope_type: InventoryTaskScopeType
  scope_value: string
  assignee_id?: number
  remark?: string
}

// 盘库扫码请求
export interface InventoryScanRequest {
  trace_code: string
  scanned_location_code: string
}

// 盘库扫码结果
export interface InventoryScanResult {
  success: boolean
  scan_result: ScanResultType
  trace_code: string
  drug_name?: string
  batch_number?: string
  system_location_code?: string
  scanned_location_code: string
  message: string
}

// 盘库任务列表查询
export interface InventoryTaskListQuery {
  page?: number
  page_size?: number
  task_no?: string
  scope_type?: InventoryTaskScopeType
  status?: InventoryTaskStatus
  assignee_id?: number
  start_time_start?: string
  start_time_end?: string
}

// 盘亏候选处理
export interface LossCandidateItem {
  id: number
  trace_code: string
  drug_id: number
  drug_name: string
  specification: string
  batch_number: string
  expire_date: string
  system_location_code: string
  status: 'LOSS_CANDIDATE'
}

// 错架记录
export interface MisplacedItem {
  id: number
  trace_code: string
  drug_id: number
  drug_name: string
  specification: string
  system_location_code: string
  scanned_location_code: string
  status: 'MISPLACED'
}
