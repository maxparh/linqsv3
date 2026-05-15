import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '@/views/AuthView.vue'
import HomeView from '@/views/HomeView.vue'
import LinksView from '@/views/LinksView.vue'
import ShortLinkView from '@/views/ShortLinkView.vue'
import AnalyticsView from '@/views/AnalyticsView.vue'
import SettingsView from '@/views/SettingsView.vue'

const isAuthenticated = true

const routes = [
  {
    path: '/cookie-policy',
    name: 'CookiePolicy',
    component: () => import('@/views/CookiePolicyView.vue'),
    meta: { requiresAuth: false } // публичная страница
  },
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
  {
    path: '/links',
    name: 'links',
    component: isAuthenticated ? LinksView : () => import('@/views/UnauthorizedView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/analytics',
    name: 'analytics',
    component: isAuthenticated ? AnalyticsView : () => import('@/views/UnauthorizedView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: isAuthenticated ? SettingsView : () => import('@/views/UnauthorizedView.vue'),
    meta: { requiresAuth: true }
  },
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