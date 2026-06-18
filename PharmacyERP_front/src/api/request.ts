import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { message } from 'ant-design-vue'
import type { ApiResponse } from '@/types/common'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：自动附加 JWT token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// 响应拦截器：统一处理错误
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data
    if (res.code !== 200) {
      message.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message))
    }
    return response
  },
  (error) => {
    // 401：token 失效，跳转登录页
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      redirectToLogin()
      return Promise.reject(error)
    }

    // 403：权限不足
    if (error.response?.status === 403) {
      message.error('权限不足，无法执行此操作')
      return Promise.reject(error)
    }

    // 404：资源不存在
    if (error.response?.status === 404) {
      message.error('资源不存在')
      return Promise.reject(error)
    }

    // 409：业务状态冲突
    if (error.response?.status === 409) {
      const msg = error.response.data?.message || '操作冲突，请刷新后重试'
      message.error(msg)
      return Promise.reject(error)
    }

    // 502/503：外部服务不可用
    if (error.response?.status === 502 || error.response?.status === 503) {
      message.error('外部服务暂时不可用，请稍后重试')
      return Promise.reject(error)
    }

    // 其他错误
    const errMsg = error.response?.data?.message || error.message || '网络错误，请稍后重试'
    message.error(errMsg)
    return Promise.reject(error)
  },
)

function redirectToLogin() {
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

export async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const res = await request.get<ApiResponse<T>>(url, { params })
  return res.data.data
}

export async function post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  const res = await request.post<ApiResponse<T>>(url, data, config)
  return res.data.data
}

export async function put<T>(url: string, data?: unknown): Promise<T> {
  const res = await request.put<ApiResponse<T>>(url, data)
  return res.data.data
}

export async function patch<T>(url: string, data?: unknown): Promise<T> {
  const res = await request.patch<ApiResponse<T>>(url, data)
  return res.data.data
}

export async function del<T>(url: string): Promise<T> {
  const res = await request.delete<ApiResponse<T>>(url)
  return res.data.data
}

export default request
