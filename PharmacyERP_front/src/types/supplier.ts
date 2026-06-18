// 供应商相关类型

export interface Supplier {
  id: number
  supplier_code: string
  name: string
  contact_name?: string
  contact_phone?: string
  address?: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateSupplierRequest {
  supplier_code: string
  name: string
  contact_name?: string
  contact_phone?: string
  address?: string
}

export interface UpdateSupplierRequest extends Partial<CreateSupplierRequest> {
  status?: number
}

export interface SupplierListQuery {
  page?: number
  page_size?: number
  name?: string
  supplier_code?: string
  status?: number
}
