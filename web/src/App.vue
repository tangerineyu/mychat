<template>
  <div class="layout">
    <header class="topbar">
      <div class="brand">My Chat</div>
      <div class="spacer"></div>
      <template v-if="auth.isAuthed">
        <span class="uid">{{ auth.userId }}</span>
        <button class="btn" @click="logout">退出</button>
      </template>
    </header>

    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

function logout() {
  auth.logout()
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto;
}
.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
}
.brand {
  font-weight: 700;
}
.spacer {
  flex: 1;
}
.uid {
  font-size: 12px;
  color: #666;
}
.btn {
  padding: 6px 10px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
}
.main {
  flex: 1;
  padding: 16px;
}
</style>

