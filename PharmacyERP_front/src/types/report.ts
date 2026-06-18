// 报表相关类型

// 看板概览（对应 DashboardOverview schema）
export interface DashboardOverview {
  today_sales_amount: number
  today_sales_count: number
  today_inbound_count: number
  in_stock_count: number
  near_expire_count: number
  loss_candidate_count: number
  pending_shelving_count: number
  active_alert_count: number
}

// 销售趋势数据（对应 SalesTrendData schema）
export interface SalesTrendPoint {
  date: string
  sales_amount: number
  sales_count: number
}

export interface SalesTrendData {
  start_date: string
  end_date: string
  granularity: 'DAY' | 'WEEK' | 'MONTH'
  points: SalesTrendPoint[]
}

// 热销药品（对应 TopDrugItem schema）
export interface TopDrugItem {
  rank: number
  drug_id: number
  drug_code: string
  common_name: string
  specification: string
  total_quantity: number
  total_amount: number
}

// 销售报表数据
export interface SalesReportData {
  start_date: string
  end_date: string
  total_orders: number
  total_amount: number
  items: SalesReportRow[]
}

export interface SalesReportRow {
  order_no: string
  order_date: string
  cashier_name: string
  drug_name: string
  specification: string
  batch_number: string
  trace_code: string
  quantity: number
  unit_price: number
  subtotal: number
  payment_method: string
}

// 入库报表数据
export interface InboundReportData {
  start_date: string
  end_date: string
  total_orders: number
  total_amount: number
  items: InboundReportRow[]
}

export interface InboundReportRow {
  order_no: string
  order_date: string
  supplier_id: number
  supplier_name: string
  invoice_no: string
  drug_name: string
  specification: string
  batch_number: string
  expire_date: string
  planned_qty: number
  confirmed_qty: number
  unit_price: number
  subtotal: number
}

// 库存报表数据
export interface InventoryReportData {
  generated_at: string
  total_count: number
  items: InventoryReportRow[]
}

export interface InventoryReportRow {
  trace_code: string
  drug_code: string
  common_name: string
  specification: string
  manufacturer: string
  batch_number: string
  expire_date: string
  status: string
  location_code: string | null
  unit_price: number
}

// 追溯日志报表数据（带分页，与后端 core.PageResult 对齐）
export interface TraceLogReportData {
  total: number
  page: number
  page_size: number
  total_pages: number
  list: TraceLogReportRow[]
}

export interface TraceLogReportRow {
  id: number
  trace_code: string
  drug_name: string
  action_type: string
  from_status: string
  to_status: string
  operator_name: string
  related_no: string
  remark: string
  created_at: string
}

// 导出任务（openapi.yaml ReportExportTask）
export interface ReportExportTask {
  task_id: string
  report_type: string
  status: 'PENDING' | 'RUNNING' | 'SUCCESS' | 'FAILED'
  file_id?: string | null
  message?: string | null
  created_at: string
}
