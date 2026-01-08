import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const Login = () => import('../views/LoginView.vue')
const Home = () => import('../views/HomeView.vue')
const Chat = () => import('../views/ChatView.vue')

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/login', component: Login },
    { path: '/home', component: Home, meta: { requiresAuth: true } },
    { path: '/chat/:type/:targetId', component: Chat, meta: { requiresAuth: true } }
  ]
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthed) {
    return '/login'
  }
  if (to.path === '/login' && auth.isAuthed) {
    return '/home'
  }
})
