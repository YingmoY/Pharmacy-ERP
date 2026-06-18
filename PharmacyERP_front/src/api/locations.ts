// 货位 API
import { get, post, put, del } from './request'
import type { LocationInfo, CreateLocationRequest, UpdateLocationRequest, LocationListQuery, MixCheckResult } from '@/types/location'
import type { PageResponse } from '@/types/common'

export function getLocationList(params: LocationListQuery) {
  return get<PageResponse<LocationInfo>>('/locations', params as Record<string, unknown>)
}

export function getLocationDetail(id: number) {
  return get<LocationInfo>(`/locations/${id}`)
}

export function getLocationByCode(code: string) {
  return get<LocationInfo>(`/locations/code/${code}`)
}

export function createLocation(data: CreateLocationRequest) {
  return post<LocationInfo>('/locations', data)
}

export function updateLocation(id: number, data: UpdateLocationRequest) {
  return put<LocationInfo>(`/locations/${id}`, data)
}

export function deleteLocation(id: number) {
  return del<null>(`/locations/${id}`)
}

export function toggleLocationStatus(id: number, status: 0 | 1) {
  return put<LocationInfo>(`/locations/${id}/status`, { status })
}

// 货位混放检查（GET /shelving/mix-check?location_code=...）
export function checkMixPlacement(location_code: string) {
  return get<MixCheckResult>('/shelving/mix-check', { location_code } as Record<string, unknown>)
}
