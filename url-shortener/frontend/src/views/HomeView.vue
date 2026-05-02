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
            @click="createLink"
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
                <button class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors">
                  <img src="@/components/icons/link_secondary.svg" alt="" />
                </button>
                <button
                  @click="copyLink(link.shortUrl)"
                  class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors"
                >
                  <img src="@/components/icons/copy_secondary.svg" alt="" />
                </button>
                <button class="w-8 h-8 flex items-center justify-center text-text-secondary hover:text-primary transition-colors">
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
    <LinkErrorPopup
      v-if="showErrorPopup"
      :url="errorUrl"
      @close="showErrorPopup = false"
    />

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
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-page-bg flex items-center justify-center">
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

const router = useRouter()
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
    // используем короткий URL без изменений
    await navigator.clipboard.writeText(`${shortUrl}`)
    showToast.value = true
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// Получение всех ссылок
const fetchLinks = async () => {
  const token = localStorage.getItem('token')
  if (!token) return router.push('/auth')

  try {
    const response = await fetch('http://localhost:8080/api/links', {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        router.push('/auth')
      }
      return
    }

    const data = await response.json()
    links.value = data.map((item: any) => ({
      id: item.id,
      shortUrl: `http://localhost:8080/${item.short_link}`,
      originalUrl: item.original_link?.replace(/^https?:\/\//, '') || '',
      clicks: (item.clicks ?? 0).toString(),
      growth: 0,
    }))
  } catch (err) {
    console.error('Failed to fetch links:', err)
  }
}

// Создание сокращенной ссылки
const handleShorten = async () => {
  if (!newLink.url) return

  const token = localStorage.getItem('token')
  if (!token) return router.push('/auth')

  try {
    const response = await fetch('http://localhost:8080/api/shorten', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ url: newLink.url }),
    })

    if (!response.ok) {
      // Показываем негативный попап
      errorUrl.value = newLink.url
      showErrorPopup.value = true
      return
    }

    const shortCode = await response.text()
    const shortUrl = `http://localhost:8080/${shortCode}`
    createdLink.value = shortUrl

    links.value.unshift({
      id: Date.now(),
      shortUrl,
      originalUrl: newLink.url.replace(/^https?:\/\//, ''),
      clicks: '0',
      growth: 0,
    })

    newLink.url = ''
    newLink.comment = ''
    showLinkPopup.value = true
  } catch (err) {
    console.error(err)
    errorUrl.value = newLink.url
    showErrorPopup.value = true
  }
}

const createLink = () => handleShorten()

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/auth')
}

onMounted(() => fetchLinks())
</script>