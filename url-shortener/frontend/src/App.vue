<!-- App.vue -->
<template>
  <router-view />
  <CookieConsentBanner ref="cookieBanner" />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import CookieConsentBanner from '@/components/CookieConsentBanner.vue'

const router = useRouter()
const cookieBanner = ref<InstanceType<typeof CookieConsentBanner> | null>(null)

watch(
  () => router.currentRoute.value.path,
  (newPath, oldPath) => {
    // Если перешли с страницы авторизации на защищённую страницу
    if ((oldPath === '/auth' || oldPath === '/register') && 
        newPath !== '/auth' && newPath !== '/register') {
      // 🔥 Исправлено: вызываем правильный метод
      setTimeout(() => {
        cookieBanner.value?.checkAndShow()
      }, 500)
    }
  }
)
</script>