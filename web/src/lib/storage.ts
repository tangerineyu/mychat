const ACCESS_TOKEN_KEY = 'mychat_access_token'
const REFRESH_TOKEN_KEY = 'mychat_refresh_token'
const USER_ID_KEY = 'mychat_user_id'

export function getAccessToken(): string {
  return localStorage.getItem(ACCESS_TOKEN_KEY) || ''
}
export function setAccessToken(v: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, v)
}
export function getRefreshToken(): string {
  return localStorage.getItem(REFRESH_TOKEN_KEY) || ''
}
export function setRefreshToken(v: string) {
  localStorage.setItem(REFRESH_TOKEN_KEY, v)
}
export function getUserId(): string {
  return localStorage.getItem(USER_ID_KEY) || ''
}
export function setUserId(v: string) {
  localStorage.setItem(USER_ID_KEY, v)
}

export function clearAuth() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(USER_ID_KEY)
}

