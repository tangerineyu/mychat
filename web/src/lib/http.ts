import axios, { AxiosError } from 'axios'
import { getAccessToken, getRefreshToken, setAccessToken, setRefreshToken } from '@/lib/storage'

export type ApiResp<T> = { code: number; message: string; data: T }

export const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

http.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let refreshing: Promise<void> | null = null

async function refreshTokenIfNeeded() {
  const rt = getRefreshToken()
  if (!rt) throw new Error('no refresh token')

  // align with backend: POST /refresh-token { refresh_token }
  const res = await axios.post<ApiResp<{ access_token?: string; refresh_token?: string }>>(
    '/api/v1/refresh-token',
    { refresh_token: rt }
  )

  // backend's RefreshToken handler returns {code,message,data:{access_token,refresh_token}}
  if (res.data.code !== 200) throw new Error(res.data.message || 'refresh failed')

  const access = res.data.data?.access_token || ''
  const refresh = res.data.data?.refresh_token || ''
  if (!access || !refresh) throw new Error('refresh response missing tokens')

  setAccessToken(access)
  setRefreshToken(refresh)
}

http.interceptors.response.use(
  (resp) => resp,
  async (err: AxiosError<ApiResp<any>>) => {
    const status = err.response?.status
    const code = err.response?.data?.code

    // backend returns HTTP 200 with business code typically, but some handlers return 401.
    const shouldTryRefresh = status === 401 || code === 401 || code === 40001

    if (!shouldTryRefresh) throw err

    if (!refreshing) {
      refreshing = refreshTokenIfNeeded().finally(() => {
        refreshing = null
      })
    }

    await refreshing

    // retry original request
    const config = err.config
    if (!config) throw err

    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${getAccessToken()}`

    return http.request(config)
  }
)

export async function apiPost<T>(url: string, body?: any): Promise<T> {
  const res = await http.post<ApiResp<T>>(url, body || {})
  if (res.data.code !== 200) {
    throw new Error(res.data.message || `api error code=${res.data.code}`)
  }
  return res.data.data
}

