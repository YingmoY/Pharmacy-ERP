// 货位相关类型

export interface LocationInfo {
  id: number
  location_code: string
  location_name: string
  area: string
  shelf?: string
  layer?: number
  position?: number
  capacity?: number
  status: number
  created_at: string
  updated_at: string
}

export interface CreateLocationRequest {
  location_code: string
  location_name: string
  area: string
  shelf?: string
  layer?: number
  position?: number
  capacity?: number
}

export interface UpdateLocationRequest extends Partial<CreateLocationRequest> {
  status?: number
}

export interface LocationListQuery {
  page?: number
  page_size?: number
  location_code?: string
  location_name?: string
  area?: string
  status?: number
}

// 混放检查结果
export interface MixCheckResult {
  location_code: string
  location_name: string
  has_mixed_drugs: boolean
  drugs?: Array<{
    drug_id: number
    drug_name: string
    count: number
    batch_numbers?: string[]
  }>
}
