// 报表 API
import { get, post } from './request'
import type { SalesReportData, InboundReportData, InventoryReportData, TraceLogReportData, ReportExportTask } from '@/types/report'

// ===== 销售报表 =====
export function getSalesReport(params: {
  start_date: string
  end_date: string
  drug_id?: number
  cashier_id?: number
}) {
  return get<SalesReportData>('/reports/sales', params as Record<string, unknown>)
}

// ===== 入库报表 =====
export function getInboundReport(params: {
  start_date: string
  end_date: string
  supplier_id?: number
}) {
  return get<InboundReportData>('/reports/inbound', params as Record<string, unknown>)
}

// ===== 库存报表 =====
export function getInventoryReport(params: {
  drug_id?: number
  location_id?: number
  status?: string
}) {
  return get<InventoryReportData>('/reports/inventory', params as Record<string, unknown>)
}

// ===== 追溯日志报表（支持分页）=====
export function getTraceLogReport(params: {
  trace_code?: string
  drug_id?: number
  action_type?: string
  operator_id?: number
  date_start?: string
  date_end?: string
  related_no?: string
  page?: number
  page_size?: number
}) {
  return get<TraceLogReportData>('/reports/trace-log', params as Record<string, unknown>)
}

// ===== 报表导出（POST 到专属导出接口，返回导出任务）=====

export interface ReportExportRequest {
  format: 'XLSX' | 'CSV' | 'PDF'
  filters?: Record<string, unknown>
  async?: boolean
}

export function exportSalesReport(filters: Record<string, unknown>, format: ReportExportRequest['format'] = 'XLSX') {
  return post<ReportExportTask>('/reports/sales/export', { format, filters, async: true })
}

export function exportInboundReport(filters: Record<string, unknown>, format: ReportExportRequest['format'] = 'XLSX') {
  return post<ReportExportTask>('/reports/inbound/export', { format, filters, async: true })
}

export function exportInventoryReport(filters: Record<string, unknown>, format: ReportExportRequest['format'] = 'XLSX') {
  return post<ReportExportTask>('/reports/inventory/export', { format, filters, async: true })
}

export function exportTraceLogReport(filters: Record<string, unknown>, format: ReportExportRequest['format'] = 'XLSX') {
  return post<ReportExportTask>('/reports/trace-log/export', { format, filters, async: true })
}

// ===== 导出任务查询（单条）=====
export function getExportTaskDetail(task_id: string) {
  return get<ReportExportTask>(`/reports/export-tasks/${task_id}`)
}
