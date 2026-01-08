import { defineStore } from 'pinia'
import { apiPost } from '@/lib/http'
import { clearAuth, getAccessToken, getRefreshToken, getUserId, setAccessToken, setRefreshToken, setUserId } from '@/lib/storage'
import { router } from '@/router'

type LoginResp = {
  token: string
  refresh_token: string
  nickname?: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: getAccessToken(),
    refreshToken: getRefreshToken(),
    userId: getUserId()
  }),
  getters: {
    isAuthed: (s) => !!s.accessToken
  },
  actions: {
    async login(telephone: string, password: string) {
      const data = await apiPost<LoginResp>('/login', { telephone, password })
      setAccessToken(data.token)
      setRefreshToken(data.refresh_token)
      this.accessToken = data.token
      this.refreshToken = data.refresh_token

      // 后端 login 返回里没有 userId；为了 WS / 列表请求，先从 token 里解码也行。
      // 这里用一个简单约定：让用户手动输入 userId 或从会话列表里再获得。
      // 更好的方式：后端在 login data 中返回 user uuid。
      await router.push('/home')
    },

    async register(telephone: string, password: string, nickname: string) {
      await apiPost<any>('/register', { telephone, password, nickname })
    },

    setUserId(uid: string) {
      setUserId(uid)
      this.userId = uid
    },

    logout() {
      clearAuth()
      this.accessToken = ''
      this.refreshToken = ''
      this.userId = ''
      router.push('/login')
    }
  }
})

