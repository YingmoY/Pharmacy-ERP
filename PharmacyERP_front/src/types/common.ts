// 通用类型定义

// 统一响应格式
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  request_id: string
}

// 分页响应（与后端 PageResult 结构对齐：list + 扁平分页字段）
export interface PageResponse<T> {
  total: number
  page: number
  page_size: number
  total_pages: number
  list: T[]
}

// 通用分页查询参数
export interface PageQuery {
  page?: number
  page_size?: number
}

// 状态枚举：通用启用/停用
export type EnableStatus = 0 | 1
