// 看板统计 API
import { get } from './request'
import type { DashboardOverview, SalesTrendData, TopDrugItem } from '@/types/report'

// 获取看板概览数据
export function getDashboardOverview() {
  return get<DashboardOverview>('/dashboard/overview')
}

// 获取销售趋势（指定日期范围）
export function getSalesTrend(params: { start_date: string; end_date: string; granularity?: 'DAY' | 'WEEK' | 'MONTH' }) {
  return get<SalesTrendData>('/dashboard/sales-trend', params as Record<string, unknown>)
}

// 获取热销药品排行
export function getTopSellingDrugs(params: { start_date: string; end_date: string; top_n?: number; sort_by?: 'QUANTITY' | 'AMOUNT' }) {
  return get<TopDrugItem[]>('/dashboard/top-drugs', params as Record<string, unknown>)
}

// 获取库存统计
export function getInventoryStats() {
  return get<Record<string, unknown>>('/dashboard/inventory-stats')
}
