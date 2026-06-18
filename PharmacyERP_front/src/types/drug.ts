// 药品相关类型

export interface DrugInfo {
  id: number
  drug_code: string
  common_name: string
  trade_name?: string
  specification: string
  manufacturer: string
  dosage_form?: string
  unit?: string
  is_prescription: boolean
  is_medicare: boolean
  retail_price: string
  status: number
  created_at: string
  updated_at: string
  in_stock_count?: number
}

export interface CreateDrugRequest {
  drug_code: string
  common_name: string
  trade_name?: string
  specification: string
  manufacturer: string
  dosage_form?: string
  unit?: string
  is_prescription: boolean
  is_medicare: boolean
  retail_price: string
}

export interface UpdateDrugRequest extends Partial<CreateDrugRequest> {
  status?: number
}

export interface DrugListQuery {
  page?: number
  page_size?: number
  keyword?: string
  drug_code?: string
  common_name?: string
  manufacturer?: string
  is_prescription?: boolean
  status?: number
}

// AI 药品搜索请求（ai-openapi.yaml DrugSearchRequest）
export interface DrugSearchRequest {
  query: string
  search_mode?: 'KEYWORD' | 'FUZZY' | 'SEMANTIC' | 'HYBRID'
  limit?: number
  offset?: number
  filters?: {
    only_available?: boolean
    is_prescription?: boolean | null
    is_medicare?: boolean | null
    manufacturer?: string
    dosage_form?: string
    storage_condition?: string
    near_expire_only?: boolean
  }
  context?: {
    scene?: 'SALE' | 'INBOUND' | 'INVENTORY' | 'GENERAL'
    prefer_available?: boolean
  }
}

// 库存摘要（ai-openapi.yaml DrugInventorySummary）
export interface DrugInventorySummary {
  available_qty?: number
  in_stock_qty?: number
  reserved_qty?: number
  pending_qty?: number
  abnormal_qty?: number
  near_expire_available_qty?: number
  nearest_expire_date?: string | null
}

// AI 药品搜索结果条目（ai-openapi.yaml DrugSearchItem）
export interface DrugSearchResult {
  drug_id: number
  drug_code: string
  common_name: string
  trade_name?: string | null
  specification: string
  dosage_form?: string | null
  manufacturer: string
  approval_number?: string | null
  barcode?: string | null
  unit?: string | null
  retail_price?: string | null
  purchase_price?: string | null
  is_prescription: boolean
  is_medicare: boolean
  inventory?: DrugInventorySummary
  score: number
  match_reason?: string
  highlights?: string[]
}
