import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '@/views/AuthView.vue'
import HomeView from '@/views/HomeView.vue'

const isAuthenticated = true // Поменяй на false чтобы проверить редирект

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
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
