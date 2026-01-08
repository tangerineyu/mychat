<template>
  <div class="grid">
    <section class="panel">
      <h3>我的信息</h3>
      <div class="row">
        <label class="inline">
          userId
          <input v-model="uid" placeholder="填写自己的 userId（uuid）" />
        </label>
        <button class="btn" @click="saveUid">保存</button>
      </div>

      <div class="row">
        <button class="btn" @click="loadAll" :disabled="loading">刷新列表</button>
        <span class="muted" v-if="loading">加载中...</span>
      </div>

      <p class="err" v-if="error">{{ error }}</p>
    </section>

    <section class="panel">
      <h3>会话列表</h3>
      <ul class="list">
        <li v-for="s in sessions" :key="s.id || s.target_id" class="item" @click="openSession(s)">
          <div class="title">{{ titleForSession(s) }}</div>
          <div class="sub">{{ s.last_msg }} · {{ s.unread_cnt ?? 0 }}</div>
        </li>
      </ul>
    </section>

    <section class="panel">
      <h3>联系人</h3>
      <ul class="list">
        <li v-for="c in contacts" :key="c.user_id" class="item" @click="goChat(1, c.user_id)">
          <div class="title">{{ c.nickname }} ({{ c.user_id }})</div>
          <div class="sub">{{ c.desc || '' }}</div>
        </li>
      </ul>
    </section>

    <section class="panel">
      <h3>我的群组</h3>
      <ul class="list">
        <li v-for="g in groups" :key="g.uuid" class="item" @click="goChat(2, g.uuid)">
          <div class="title">{{ g.name }} ({{ g.uuid }})</div>
          <div class="sub">owner: {{ g.owner_id }}</div>
        </li>
      </ul>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiPost } from '@/lib/http'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const auth = useAuthStore()

const uid = ref(auth.userId)
const loading = ref(false)
const error = ref('')

type Session = any
type Contact = { user_id: string; nickname: string; avatar?: string; desc?: string }
type Group = any

const sessions = ref<Session[]>([])
const contacts = ref<Contact[]>([])
const groups = ref<Group[]>([])

function saveUid() {
  auth.setUserId(uid.value)
}

function titleForSession(s: any) {
  const t = s.type
  if (t === 1) return `私聊: ${s.target_id}`
  if (t === 2) return `群聊: ${s.target_id}`
  return s.target_id || 'unknown'
}

function openSession(s: any) {
  // 这里假设 session 返回有 type(1/2) target_id
  goChat(s.type ?? 1, s.target_id)
}

function goChat(type: number, targetId: string) {
  router.push(`/chat/${type}/${targetId}`)
}

async function loadAll() {
  error.value = ''
  loading.value = true
  try {
    // contacts
    contacts.value = await apiPost<Contact[]>('/contact/list', {})
    // sessions
    sessions.value = await apiPost<Session[]>('/session/list', {})
    // groups
    // 后端 group/loadMyGroup 返回结构不太确定（当前服务层实现看上去可能有问题），先直接调用
    groups.value = (await apiPost<any>('/group/loadMyGroup', {})) as any
  } catch (e: any) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAll()
})
</script>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
@media (max-width: 900px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
.panel {
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 12px;
}
.row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-top: 10px;
}
.inline {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
input {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px;
}
.btn {
  padding: 10px 12px;
  border: 1px solid #ddd;
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
}
.list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.item {
  padding: 10px;
  border: 1px solid #f0f0f0;
  border-radius: 10px;
  margin-top: 10px;
  cursor: pointer;
}
.title {
  font-weight: 600;
}
.sub {
  margin-top: 4px;
  font-size: 12px;
  color: #666;
}
.muted {
  color: #666;
  font-size: 12px;
}
.err {
  color: #b00020;
}
</style>

