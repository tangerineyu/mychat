import { getAccessToken } from '@/lib/storage'

export type WSChatContent = {
  send_id: string
  receiver_id: string
  // 注意：后端协议里这里的 type 表示消息类型（1 文本 / 2 图片）
  type: number
  content: string
  uuid?: string
}

export type WSMessage =
  | { action: 'heartbeat' }
  | { action: 'chat_message'; content: WSChatContent; trace_id?: string }

export class WSClient {
  private ws: WebSocket | null = null
  private heartbeatTimer: number | null = null

  constructor(private url: string) {}

  connect(onMessage: (data: any) => void) {
    if (this.ws) return

    // backend supports query token as fallback
    const token = getAccessToken()
    const u = new URL(this.url)
    if (token) u.searchParams.set('token', token)

    this.ws = new WebSocket(u.toString())

    this.ws.onopen = () => {
      this.startHeartbeat()
    }

    this.ws.onmessage = (evt) => {
      try {
        const data = JSON.parse(evt.data)
        onMessage(data)
      } catch {
        onMessage(evt.data)
      }
    }

    this.ws.onclose = () => {
      this.stopHeartbeat()
      this.ws = null
    }

    this.ws.onerror = () => {
      // let onclose handle cleanup
    }
  }

  send(msg: WSMessage) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return
    this.ws.send(JSON.stringify(msg))
  }

  close() {
    this.stopHeartbeat()
    this.ws?.close()
    this.ws = null
  }

  private startHeartbeat() {
    this.stopHeartbeat()
    this.heartbeatTimer = window.setInterval(() => {
      this.send({ action: 'heartbeat' })
    }, 25_000)
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      window.clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }
}
