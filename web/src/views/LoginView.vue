<template>
  <div class="wrap">
    <div class="card">
      <h2>登录</h2>

      <label>
        手机号
        <input v-model="telephone" placeholder="telephone" />
      </label>

      <label>
        密码
        <input v-model="password" type="password" placeholder="password" />
      </label>

      <button class="btn" :disabled="loading" @click="onLogin">登录</button>

      <div class="divider"></div>

      <h3>注册</h3>
      <label>
        昵称
        <input v-model="nickname" placeholder="nickname" />
      </label>

      <button class="btn secondary" :disabled="loading" @click="onRegister">注册</button>

      <p class="err" v-if="error">{{ error }}</p>

      <p class="tip">
        注意：当前后端 login 的 data 里没有返回 userId/uuid，但 WS 发消息需要 send_id。
        你登录后到 Home/Chat 页面手动填写一次 userId。
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const telephone = ref('')
const password = ref('')
const nickname = ref('')

const loading = ref(false)
const error = ref('')

async function onLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(telephone.value, password.value)
  } catch (e: any) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(telephone.value, password.value, nickname.value)
    error.value = '注册成功，请直接登录'
  } catch (e: any) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 80px);
}
.card {
  width: 420px;
  max-width: calc(100vw - 32px);
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 16px;
}
label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 12px;
  font-size: 14px;
}
input {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px;
}
.btn {
  margin-top: 14px;
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 8px;
  background: #111;
  color: white;
  cursor: pointer;
}
.btn.secondary {
  background: #444;
}
.divider {
  margin: 16px 0;
  height: 1px;
  background: #eee;
}
.err {
  margin-top: 10px;
  color: #b00020;
}
.tip {
  margin-top: 10px;
  color: #666;
  font-size: 12px;
  line-height: 1.5;
}
</style>

