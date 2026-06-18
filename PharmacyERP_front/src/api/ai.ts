// AI 辅助 API
import axios from 'axios'
import { get, post } from './request'
import type { DrugSearchRequest, DrugSearchResult } from '@/types/drug'
import type { PageResponse } from '@/types/common'

// AI 药品搜索（前端直连 /ai/api/v1，无需 JWT）
const aiRequest = axios.create({
  baseURL: '/ai/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

export async function searchDrugsWithAI(data: DrugSearchRequest): Promise<DrugSearchResult[]> {
  const res = await aiRequest.post('/drugs/search', data)
  // 响应结构: { code, message, request_id, data: { query, total, items } }
  return res.data.data?.items ?? []
}

export interface DrugRecommendRequest {
  query: string
  limit?: number
  filters?: { only_available?: boolean }
}

export interface DrugRecommendResponse {
  query: string
  explanation: string
  terms: string[]
  total: number
  items: DrugSearchResult[]
}

export async function recommendDrugsWithAI(data: DrugRecommendRequest): Promise<DrugRecommendResponse> {
  const res = await aiRequest.post('/drugs/recommend', data, { timeout: 600000 })
  return res.data.data as DrugRecommendResponse
}

// ===== 发票识别（通过主 ERP /api/v1/ai/invoices 接口操作）=====

// ai_invoice_record 表结构（openapi.yaml InvoiceRecord）
export interface InvoiceRecord {
  id: number
  file_id: string
  file_name: string
  status: 'PENDING' | 'PROCESSING' | 'COMPLETED' | 'FAILED'
  result?: InvoiceRecognizeResult | null
  recognized_supplier_name?: string | null
  matched_supplier_id?: number | null
  invoice_no?: string | null
  supplier_name?: string
  drug_count?: number
  total_amount?: number
  confidence?: number
  error_message?: string | null
  inbound_order_id?: number | null
  creator_id?: number
  creator_name?: string
  created_at: string
}

// openapi.yaml InvoiceRecognizeResult（也被 ai-openapi.yaml 引用）
export interface InvoiceRecognizeResult {
  recognized_supplier_name?: string | null
  supplier_candidates?: SupplierCandidate[]
  matched_supplier_id?: number | null
  invoice_no?: string | null
  invoice_date?: string | null
  total_amount?: string | null
  confidence?: number | null
  items: InvoiceRecognizeItem[]
  warnings?: QualityWarning[]
}

export interface InvoiceRecognizeItem {
  row_index: number
  drug_name?: string | null
  specification?: string | null
  manufacturer?: string | null
  approval_number?: string | null
  batch_number?: string | null
  expire_date?: string | null
  quantity?: string | null
  unit_price?: string | null
  amount?: string | null
  confidence?: number | null
  matched_drug_id?: number | null
  drug_candidates?: DrugCandidate[]
  warnings?: QualityWarning[]
}

export interface QualityWarning {
  level: 'LOW' | 'MEDIUM' | 'HIGH'
  code?: string
  field?: string | null
  message: string
  suggestion?: string | null
}

export interface SupplierCandidate {
  supplier_id: number
  supplier_code: string
  name: string
  license_no?: string | null
  contact_name?: string | null
  contact_phone?: string | null
  confidence: number
  match_reason?: string
}

export interface DrugCandidate extends DrugSearchResult {
  confidence?: number
  candidate_type?: 'EXACT' | 'FUZZY' | 'SEMANTIC' | 'ALIAS' | 'MANUAL'
}

// 发票列表查询参数（openapi.yaml GET /api/v1/ai/invoices）
export interface InvoiceListQuery {
  page?: number
  page_size?: number
  status?: 'PENDING' | 'PROCESSING' | 'COMPLETED' | 'FAILED'
  supplier_id?: number
}

// 上传发票并触发 AI 识别（POST /api/v1/ai/invoices/recognize）
// 必须用 bare axios（不带默认 Content-Type: application/json），让浏览器自动设置 multipart/form-data 边界。
export async function uploadAndRecognizeInvoice(
  file: File,
  options?: { supplier_id?: number; remark?: string },
): Promise<InvoiceRecord> {
  const formData = new FormData()
  formData.append('file', file)
  if (options?.supplier_id != null) formData.append('supplier_id', String(options.supplier_id))
  if (options?.remark) formData.append('remark', options.remark)

  const token = localStorage.getItem('access_token')
  const res = await axios.post<{ code: number; message: string; data: InvoiceRecord }>(
    '/api/v1/ai/invoices/recognize',
    formData,
    { headers: token ? { Authorization: `Bearer ${token}` } : undefined },
  )
  if (res.data.code !== 200) throw new Error(res.data.message || '识别失败')
  return res.data.data
}

// 获取发票识别记录列表（GET /api/v1/ai/invoices）
export function getInvoiceList(params: InvoiceListQuery) {
  return get<PageResponse<InvoiceRecord>>('/ai/invoices', params as Record<string, unknown>)
}

// 获取发票识别记录详情（GET /api/v1/ai/invoices/{id}）
export function getInvoiceDetail(id: number) {
  return get<InvoiceRecord>(`/ai/invoices/${id}`)
}

// 将识别结果转为入库单（POST /api/v1/ai/invoices/{id}/convert-to-inbound）
export interface ConvertToInboundRequest {
  supplier_id: number
  invoice_no?: string
  remark?: string
  items: Array<{
    drug_id: number
    batch_number: string
    expire_date: string
    planned_qty: number
    unit_price: number
  }>
}

export function convertInvoiceToInbound(id: number, data: ConvertToInboundRequest) {
  return post<Record<string, unknown>>(`/ai/invoices/${id}/convert-to-inbound`, data)
}

// 页面别名函数
export const getInvoiceRecordList = getInvoiceList
export const getInvoiceRecord = getInvoiceDetail
