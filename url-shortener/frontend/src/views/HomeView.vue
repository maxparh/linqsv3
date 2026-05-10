<template>
  <div class="min-h-screen bg-page-bg flex">
    <!-- Боковое меню -->
    <aside class="w-[224px] bg-white border-r border-card-border flex flex-col">
      <!-- Логотип -->
      <div class="p-6 flex items-center gap-3 pb-[80px]">
        <div class="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
          <img src="@/components/icons/linqs_logo.svg" alt="" />
        </div>
        <span class="font-inter font-semibold text-[17px] text-text-primary">Linqs</span>
      </div>

      <!-- Навигация -->
      <nav class="flex-1 px-4 space-y-1">
        <!-- Активная: Главная -->
        <router-link
          to="/"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] bg-primary text-white font-inter text-[17px] font-bold"
        >
          <img src="@/components/icons/home_nav_active.svg" alt="" />
          Главная
        </router-link>

        <!-- Неактивные -->
        <button
          @click="showWIPPopup = true"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/link_nav_nactive.svg" alt="" />
          Ссылки
        </button>

        <button
          @click="showWIPPopup = true"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/analit_nav_nactive.svg" alt="" />
          Аналитика
        </button>

        <button
          @click="showWIPPopup = true"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/settings_nav_nactive.svg" alt="" />
          Настройки
        </button>
      </nav>

      <!-- Псевдопрофиль -->
      <div class="p-4 border-t border-card-border">
        <div class="flex items-center gap-3 px-4 py-3">
          <div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
            <span class="text-white font-inter font-semibold text-[14px]">A</span>
          </div>
          <span class="font-inter text-[17px] text-text-primary">admin</span>
        </div>
        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 px-4 py-7 text-error font-inter text-[17px] font-medium hover:bg-page-bg rounded-[10px] transition-colors"
        >
          <img src="@/components/icons/logout_btn.svg" alt="" class="w-[24px]" />
          Выход
        </button>
      </div>
    </aside>

    <!-- Основной контент -->
    <main class="flex-1 p-8">
      <!-- Хедер -->
      <div class="flex items-center justify-between mb-8">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Главная</h1>
      </div>

      <!-- Создание ссылки -->
      <div class="bg-white rounded-card border border-card-border p-6 mb-8">
        <div class="flex gap-4 mb-4">
          <input
            v-model="newLink.url"
            type="url"
            placeholder="https://example.com/very/long/url/that/needs/shortening"
            class="flex-1 h-10 px-4 border border-card-border rounded-[10px] font-inter text-[17px] text-text-primary placeholder:text-placeholder focus:outline-none focus:border-primary transition-colors"
          />
          <button
            @click="handleShorten"
            class="h-10 px-6 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Сократить
          </button>
        </div>
        <div class="flex items-center gap-4">
          <input
            v-model="newLink.comment"
            type="text"
            placeholder="Комментарий (необязательно)"
            class="flex-1 h-10 px-4 border border-card-border rounded-[10px] font-inter text-[17px] text-text-primary placeholder:text-placeholder focus:outline-none focus:border-primary transition-colors"
          />
        </div>
      </div>

      <!-- Список ссылок -->
      <div class="bg-white rounded-card border border-card-border p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="font-inter text-[17px] text-text-secondary">Все ссылки</h2>
        </div>

        <div class="space-y-4">
          <div
            v-for="link in links"
            :key="link.id"
            class="flex items-center justify-between py-3 border-b border-card-border last:border-0"
          >
            <div>
              <div class="font-inter text-[17px] font-medium text-text-primary mb-1">
                {{ link.shortUrl }}
              </div>
              <div class="font-inter text-[14px] text-text-secondary">
                {{ link.originalUrl }}
              </div>
            </div>

            <div class="flex items-center gap-6">
              <div class="text-right">
                <div class="font-inter text-[17px] text-text-primary">
                  {{ link.clicks }}
                </div>
                <div
                  class="font-inter text-[14px]"
                  :class="link.growth >= 0 ? 'text-success' : 'text-error'"
                >
                  {{ link.growth >= 0 ? '+' : '' }}{{ link.growth }}%
                </div>
              </div>

              <div class="flex items-center gap-2">
                <button
                  class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors"
                >
                  <img src="@/components/icons/link_secondary.svg" alt="" />
                </button>
                <button
                  @click="copyLink(link.shortUrl)"
                  class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors"
                >
                  <img src="@/components/icons/copy_secondary.svg" alt="" />
                </button>
                <button
                  class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors"
                >
                  <img src="@/components/icons/edit_secondary.svg" alt="" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Попап создания ссылки -->
    <LinkCreatedPopup
      v-if="showLinkPopup"
      :short-url="createdLink"
      @close="showLinkPopup = false"
    />

    <!-- Попап ошибки -->
    <LinkErrorPopup v-if="showErrorPopup" :url="errorUrl" @close="showErrorPopup = false" />

    <!-- Попап копирования -->
    <ToastNotification
      v-if="showToast"
      :show="showToast"
      title="Готово!"
      message="Ссылка скопирована в буфер обмена"
      @close="showToast = false"
    />

    <!-- В разработке попап -->
    <div
      v-if="showWIPPopup"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="showWIPPopup = false"
    >
      <div class="bg-white rounded-card border border-card-border p-8 max-w-[400px] text-center">
        <div
          class="w-16 h-16 mx-auto mb-4 rounded-full bg-page-bg flex items-center justify-center"
        >
          <img src="@/components/icons/dev.svg" alt="" />
        </div>
        <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-2">Упс!</h3>
        <p class="font-inter text-[17px] text-text-secondary mb-6">
          Данный функционал в разработке
        </p>
        <button
          @click="showWIPPopup = false"
          class="h-10 px-8 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
        >
          Понятно
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import LinkCreatedPopup from '@/components/LinkCreatedPopup.vue'
import LinkErrorPopup from '@/components/ErrorCreatedPopup.vue'
import ToastNotification from '@/components/ToastNotification.vue'

// 🔥 Используем переменную окружения вместо хардкода
const API_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
const router = useRouter()

// Состояние UI
const showWIPPopup = ref(false)
const showLinkPopup = ref(false)
const showToast = ref(false)
const showErrorPopup = ref(false)
const createdLink = ref('')
const errorUrl = ref('')

const newLink = reactive({
  url: '',
  comment: '',
})

interface Link {
  id: number
  shortUrl: string
  originalUrl: string
  clicks: string
  growth: number
}

const links = ref<Link[]>([])

// Копирование ссылки
const copyLink = async (shortUrl: string) => {
  try {
    await navigator.clipboard.writeText(shortUrl)
    showToast.value = true
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// Получение всех ссылок
const fetchLinks = async () => {
  const token = localStorage.getItem('access_token')?.trim()
  if (import.meta.env.DEV) {
    console.log('🔑 Token:', token ? 'EXISTS' : 'NOT FOUND')
  }
  if (!token) return router.push('/auth')

  try {
    if (import.meta.env.DEV) {
      console.log('📡 Sending request to /api/links...')
    }
    const response = await fetch(`${API_URL}/api/links`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    })
    if (import.meta.env.DEV) {
      console.log('📥 Response status:', response.status)
      console.log('📥 Response headers:', [...response.headers.entries()])
    }

    if (!response.ok) {
      const errorText = await response.text()
      console.error('❌ Error:', errorText)
      if (response.status === 401) {
        localStorage.removeItem('access_token')
        router.push('/auth')
      }
      return
    }

    const data = await response.json()
    if (import.meta.env.DEV) {
      console.log('✅ Success! Links:', data)
    }
    links.value = data?.map((item: any) => ({
      id: item.id,
      shortUrl: `${API_URL}/${item.short_code}`,
      originalUrl: item.original_url?.replace(/^https?:\/\//, '') || '',
      clicks: '0',
      growth: 0,
    }))
  } catch (err) {
    console.error('💥 Network error:', err)
  }
}

// Создание сокращенной ссылки
const handleShorten = async () => {
  if (!newLink.url) return

  const token = localStorage.getItem('access_token')?.trim()
  if (!token) return router.push('/auth')

  try {
    const response = await fetch(`${API_URL}/api/links`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      // 🔥 Бэкенд ожидает поле original_url
      body: JSON.stringify({ original_url: newLink.url }),
    })

    if (!response.ok) {
      errorUrl.value = newLink.url
      showErrorPopup.value = true
      return
    }

    // 🔥 Бэкенд возвращает JSON, а не текст
    const data = await response.json()
    const shortCode = data.short_code
    const shortUrl = `${API_URL}/${shortCode}`

    createdLink.value = shortUrl
    showLinkPopup.value = true

    // Добавляем в локальный список для мгновенного отображения
    links.value.unshift({
      id: data.id || Date.now(),
      shortUrl,
      originalUrl: newLink.url.replace(/^https?:\/\//, ''),
      clicks: '0',
      growth: 0,
    })

    newLink.url = ''
    newLink.comment = ''
  } catch (err) {
    console.error(err)
    errorUrl.value = newLink.url
    showErrorPopup.value = true
  }
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

onMounted(() => fetchLinks())
</script>
