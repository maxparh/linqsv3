import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '@/views/AuthView.vue'
import HomeView from '@/views/HomeView.vue'
import ShortLinkView from '@/views/ShortLinkView.vue' // ← Добавь импорт

const isAuthenticated = true

const routes = [
  {
    path: '/',
    name: 'home',
    component: isAuthenticated ? HomeView : () => import('@/views/UnauthorizedView.vue'),
  },
  {
    path: '/auth',
    name: 'auth',
    component: AuthView,
  },
  
  // 🔥 ПОСЛЕДНИЙ: ловит короткие ссылки вида /mgnMyw
  {
    path: '/:shortCode',
    name: 'ShortLink',
    component: ShortLinkView,
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router