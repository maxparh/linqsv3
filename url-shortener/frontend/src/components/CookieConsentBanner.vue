<!-- src/components/CookieConsentBanner.vue -->
<template>
  <Transition
    enter-active-class="transition-transform duration-300 ease-out"
    enter-from-class="translate-y-full"
    enter-to-class="translate-y-0"
    leave-active-class="transition-transform duration-200 ease-in"
    leave-from-class="translate-y-0"
    leave-to-class="translate-y-full"
  >
    <div
      v-if="show"
      class="fixed bottom-0 left-0 right-0 bg-white border-t border-card-border p-4 md:p-6 z-50 shadow-lg"
    >
      <div class="max-w-6xl mx-auto flex flex-col md:flex-row items-start md:items-center gap-4">
        
        <div class="flex-1">
          <p class="font-inter text-[17px] text-text-primary mb-2">
            Мы используем cookies 🔍
          </p>
          <p class="font-inter text-[14px] text-text-secondary">
            Это помогает нам улучшать работу сайта. Вы можете 
            <router-link to="/cookie-policy" class="text-primary hover:underline" target="_blank">
              ознакомиться подробно
            </router-link>.
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
          <button
            @click="handleAcceptMinimal"
            class="h-10 px-5 border border-card-border rounded-[10px] font-inter text-[17px] text-text-primary hover:bg-page-bg transition-colors"
          >
            Принять необходимые
          </button>
          <button
            @click="handleAcceptAll"
            class="h-10 px-5 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Принять
          </button>
        </div>

      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const API_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const CONSENT_COOKIE = 'cookie_consent'
const show = ref(false)

// Проверка: есть ли уже согласие в cookies
const hasConsent = () => {
  return document.cookie.split(';').some(c => c.trim().startsWith(CONSENT_COOKIE))
}

// Показать баннер, если согласия нет
const checkAndShow = () => {
  if (!hasConsent()) {
    setTimeout(() => show.value = true, 800)
  }
}

// Отправить выбор на бэкенд
const sendConsent = async (accepted: boolean, prefs: Record<string, boolean> = {}) => {
  try {
    const response = await fetch(`${API_URL}/cookies/consent`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        accepted,
        preferences: { 
          necessary: true, 
          analytics: prefs.analytics || false, 
          marketing: prefs.marketing || false 
        }
      }),
      credentials: 'include',
    })
    
    if (response.ok) {
      show.value = false
      return true
    }
  } catch (e) {
    console.error('Cookie consent error:', e)
  }
  return false
}

// Принять только необходимые
const handleAcceptMinimal = async () => {
  await sendConsent(true, { analytics: false, marketing: false })
}

// Принять все
const handleAcceptAll = async () => {
  await sendConsent(true, { analytics: true, marketing: true })
}

// Проверяем при загрузке компонента
onMounted(() => {
  checkAndShow()
})

// Экспортируем метод, чтобы можно было вызвать извне (после логина)
defineExpose({ checkAndShow })
</script>