<template>
  <div class="chat">
    <aside class="sidebar">
      <div class="box">
        <div class="label">聊天对象</div>
        <div class="value">type={{ type }}, target={{ targetId }}</div>
      </div>

      <div class="box">
        <div class="label">发送者（你的 userId）</div>
        <input v-model="sendId" placeholder="send_id" />
        <button class="btn" @click="saveSendId">保存</button>
      </div>

      <div class="box">
        <button class="btn" @click="loadHistory" :disabled="loading">拉取历史</button>
        <button class="btn" @click="connectWs" :disabled="wsConnected">连接WS</button>
        <button class="btn" @click="disconnectWs" :disabled="!wsConnected">断开WS</button>
      </div>

      <p class="err" v-if="error">{{ error }}</p>
    </aside>

    <section class="main">
      <div class="messages">
        <div v-for="(m, idx) in messages" :key="idx" class="msg" :class="m.from === sendId ? 'mine' : 'other'">
          <div class="meta">
            <span>{{ m.from }}</span>
            <span class="muted">{{ m.uuid || '' }}</span>
          </div>
          <div class="bubble">{{ m.content }}</div>
        </div>
      </div>

      <div class="composer">
        <input v-model="text" placeholder="输入消息" @keydown.enter="onSend" />
        <button class="btn primary" @click="onSend" :disabled="!text">发送</button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRoute } from 'vue-router'
import { apiPost } from '@/lib/http'
import { useAuthStore } from '@/stores/auth'
import { WSClient } from '@/lib/ws'

const route = useRoute()
const auth = useAuthStore()

const type = computed(() => Number(route.params.type))
const targetId = computed(() => String(route.params.targetId))

const sendId = ref(auth.userId)
const text = ref('')
const loading = ref(false)
const error = ref('')

type Msg = { from: string; to: string; content: string; uuid?: string }
const messages = ref<Msg[]>([])

const ws = new WSClient(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api/v1/ws`)
const wsConnected = ref(false)

function saveSendId() {
  auth.setUserId(sendId.value)
}

async function loadHistory() {
  error.value = ''
  loading.value = true
  try {
    const list = await apiPost<any[]>('/chat/history', { target_id: targetId.value, type: type.value })
    // message model: { uuid, from_user_id, to_id, content, type }
    messages.value = list
      .slice()
      .reverse()
      .map((x) => ({ from: x.from_user_id, to: x.to_id, content: x.content, uuid: x.uuid }))
  } catch (e: any) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

function connectWs() {
  error.value = ''
  ws.connect((data) => {
    // consumer 推送的是 ChatMessageContent（不是 Message wrapper）: {send_id, receiver_id, type, content, uuid}
    if (data && data.send_id && data.receiver_id) {
      messages.value.push({ from: data.send_id, to: data.receiver_id, content: data.content, uuid: data.uuid })
    }
  })
  wsConnected.value = true
}

function disconnectWs() {
  ws.close()
  wsConnected.value = false
}

function onSend() {
  if (!text.value) return
  if (!wsConnected.value) {
    error.value = '请先连接 WebSocket'
    return
  }

  const payload = {
    action: 'chat_message',
    content: {
      send_id: sendId.value,
      receiver_id: targetId.value,
      type: 1,
      content: text.value
    }
  }

  ws.send(payload as any)
  text.value = ''
}

onBeforeUnmount(() => {
  ws.close()
})
</script>

<style scoped>
.chat {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 12px;
  height: calc(100vh - 80px);
}
@media (max-width: 900px) {
  .chat {
    grid-template-columns: 1fr;
    height: auto;
  }
}
.sidebar {
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 12px;
}
.box {
  margin-bottom: 12px;
}
.label {
  font-size: 12px;
  color: #666;
}
.value {
  font-weight: 600;
  margin-top: 6px;
}
input {
  width: 100%;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px;
  margin-top: 8px;
}
.btn {
  margin-top: 8px;
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
}
.btn.primary {
  background: #111;
  color: white;
  border-color: #111;
}
.main {
  border: 1px solid #eee;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.messages {
  flex: 1;
  padding: 12px;
  overflow: auto;
  background: #fafafa;
}
.msg {
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  max-width: 70%;
}
.msg.mine {
  margin-left: auto;
  align-items: flex-end;
}
.meta {
  font-size: 11px;
  color: #666;
  display: flex;
  gap: 8px;
}
.bubble {
  margin-top: 4px;
  background: white;
  border: 1px solid #eee;
  padding: 10px;
  border-radius: 10px;
  white-space: pre-wrap;
}
.composer {
  display: flex;
  gap: 10px;
  padding: 10px;
  border-top: 1px solid #eee;
  background: white;
}
.composer input {
  flex: 1;
  margin: 0;
}
.err {
  color: #b00020;
}
.muted {
  color: #999;
}
</style>
