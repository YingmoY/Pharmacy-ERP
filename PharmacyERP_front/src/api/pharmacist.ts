// 药师审核 API
import { get, post } from './request'
import type { PharmacistReview } from '@/types/sales'
import type { PageResponse } from '@/types/common'

export interface ReviewListQuery {
  page?: number
  page_size?: number
  status?: string
  order_no?: string
  submitter_id?: number
  submitted_start?: string
  submitted_end?: string
}

// 获取审核列表
export function getReviewList(params: ReviewListQuery) {
  return get<PageResponse<PharmacistReview>>('/pharmacist/reviews', params as Record<string, unknown>)
}

// 获取审核详情
export function getReviewDetail(id: number) {
  return get<PharmacistReview & { order: import('@/types/sales').SalesOrder }>(`/pharmacist/reviews/${id}`)
}

// 审核通过
export function approveReview(id: number, opinion?: string) {
  return post<PharmacistReview>(`/pharmacist/reviews/${id}/approve`, { review_opinion: opinion })
}

// 审核驳回
export function rejectReview(id: number, opinion: string) {
  return post<PharmacistReview>(`/pharmacist/reviews/${id}/reject`, { review_opinion: opinion })
}
