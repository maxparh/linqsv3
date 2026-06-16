<template>
  <div class="min-h-screen bg-page-bg flex">
    <!-- Боковое меню -->
    <aside class="w-[224px] bg-white border-r border-card-border flex flex-col">
      <div class="p-6 flex items-center gap-3 pb-[80px]">
        <div class="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
          <img src="@/components/icons/linqs_logo.svg" alt="" />
        </div>
        <span class="font-inter font-semibold text-[17px] text-text-primary">Linqs</span>
      </div>

      <nav class="flex-1 px-4 space-y-1">
        <router-link
          to="/"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/home_nav_nactive.svg" alt="" />
          Главная
        </router-link>

        <router-link
          to="/links"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] bg-primary text-white font-inter text-[17px] font-bold"
        >
          <img src="@/components/icons/link_nav_active.svg" alt="" />
          Ссылки
        </router-link>

        <router-link
          to="/analytics"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/analit_nav_nactive.svg" alt="" />
          Аналитика
        </router-link>

        <router-link
          to="/settings"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/settings_nav_nactive.svg" alt="" />
          Настройки
        </router-link>
      </nav>

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
    <main class="flex-1 p-8" @click="showDropdownId = null">
      <div class="mb-8">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary mb-6">Ссылки</h1>

        <div class="flex items-center gap-4 mb-6">
          <!-- Поиск с белым фоном -->
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Начните вводить ссылку"
              class="w-full h-10 pl-10 pr-4 border border-card-border rounded-input font-inter text-[17px] text-text-primary placeholder:text-placeholder focus:outline-none focus:border-primary transition-colors bg-white"
            />
            <img
              src="@/components/icons/search.svg"
              alt=""
              class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none opacity-50"
            />
          </div>

          <button
            @click="showWIPPopup = true"
            class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white"
          >
            Фильтры
          </button>

          <button
            @click="showWIPPopup = true"
            class="h-10 px-6 bg-primary text-white rounded-input font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Новая ссылка
          </button>
        </div>
      </div>

      <!-- Таблица: скругление только сверху -->
      <div class="bg-white border border-card-border rounded-t-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-page-bg border-b border-card-border">
              <tr>
                <th class="w-12 px-4 py-3">
                  <input
                    type="checkbox"
                    class="w-4 h-4 rounded border-card-border"
                    @change="toggleSelectAll"
                    :checked="selectedLinks.length > 0 && selectedLinks.length === paginatedLinks.length"
                  />
                </th>
                <th class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary">
                  Ссылки
                </th>
                <th
                  class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary cursor-pointer hover:text-text-primary select-none"
                  @click.stop="toggleSort('createdAt')"
                >
                  <div class="flex items-center gap-1">
                    Дата создания
                    <div class="flex flex-col gap-0.1">
                      <img
                        :src="chevronUpIcon"
                        alt=""
                        class="w-5 h-5"
                        :class="sortField === 'createdAt' && sortDirection === 'asc' ? 'text-primary' : 'opacity-30'"
                      />
                      <img
                        :src="chevronDownIcon"
                        alt=""
                        class="w-5 h-5"
                        :class="sortField === 'createdAt' && sortDirection === 'desc' ? 'text-primary' : 'opacity-30'"
                      />
                    </div>
                  </div>
                </th>

                <th
                  class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary cursor-pointer hover:text-text-primary select-none"
                  @click.stop="toggleSort('expiresAt')"
                >
                  <div class="flex items-center gap-1">
                    Действует до
                    <div class="flex flex-col gap-0.1">
                      <img
                        :src="chevronUpIcon"
                        alt=""
                        class="w-5 h-5"
                        :class="sortField === 'expiresAt' && sortDirection === 'asc' ? 'text-primary' : 'opacity-30'"
                      />
                      <img
                        :src="chevronDownIcon"
                        alt=""
                        class="w-5 h-5"
                        :class="sortField === 'expiresAt' && sortDirection === 'desc' ? 'text-primary' : 'opacity-30'"
                      />
                    </div>
                  </div>
                </th>
                <th class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary">
                  Приватность
                </th>
                <th class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary">
                  Комментарий
                </th>
                <th class="text-left px-4 py-3 font-inter text-[17px] font-medium text-text-secondary">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="link in paginatedLinks"
                :key="link.id"
                class="border-b border-card-border last:border-0 hover:bg-page-bg/30 transition-colors"
              >
                <td class="px-4 py-3">
                  <input
                    type="checkbox"
                    class="w-4 h-4 rounded border-card-border"
                    :checked="selectedLinks.includes(link.id)"
                    @change.stop="toggleSelect(link.id)"
                  />
                </td>

                <!-- Ссылка + копировать -->
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <span class="font-inter text-[17px] font-medium text-text-primary">{{ link.shortUrl }}</span>
                    <button
                      @click.stop="copyLink(link.shortUrl)"
                      class="text-text-secondary hover:text-primary transition-colors p-1"
                    >
                      <img src="@/components/icons/copy.svg" alt="" />
                    </button>
                  </div>
                  <div class="font-inter text-[13px] text-text-secondary mt-0.5 truncate max-w-[250px]">
                    {{ link.originalUrl }}
                  </div>
                </td>

                <td class="px-4 py-3 font-inter text-[17px] text-text-primary whitespace-nowrap">
                  {{ formatDate(link.createdAt) }}
                </td>
                <td class="px-4 py-3 font-inter text-[17px] text-text-primary whitespace-nowrap">
                  {{ link.expiresAt ? formatDate(link.expiresAt) : '—' }}
                </td>

                <!-- Приватность -->
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <button
                      @click.stop="togglePrivacy(link.id)"
                      class="w-10 h-5 rounded-full transition-colors relative"
                      :class="link.isPrivate ? 'bg-primary' : 'bg-card-border'"
                    >
                      <div
                        class="w-4 h-4 bg-white rounded-full shadow-sm transition-transform absolute top-0.5"
                        :class="link.isPrivate ? 'translate-x-5' : 'translate-x-0.5'"
                      />
                    </button>
                    <img :src="link.isPrivate ? lockedIcon : lockIcon" alt="" />
                  </div>
                </td>

                <!-- Комментарий + редактировать -->
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <span class="font-inter text-[17px] text-text-primary truncate max-w-[180px]">{{
                      link.comment || '—'
                    }}</span>
                    <button
                      @click.stop="editComment(link.id)"
                      class="text-text-secondary hover:text-primary transition-colors p-1"
                    >
                      <img src="@/components/icons/edit.svg" alt="" />
                    </button>
                  </div>
                </td>

                <!-- Действия (три точки) -->
                <td class="px-4 py-3 relative">
                  <button
                    @click.stop="toggleDropdown(link.id)"
                    class="p-1 text-text-secondary hover:text-text-primary transition-colors rounded-lg hover:bg-page-bg"
                  >
                    <img src="@/components/icons/more_vert.svg" alt="" />
                  </button>

                  <Transition
                    enter-active-class="transition-all duration-200 ease-out"
                    enter-from-class="opacity-0 scale-95 -translate-y-1"
                    enter-to-class="opacity-100 scale-100 translate-y-0"
                    leave-active-class="transition-all duration-150 ease-in"
                    leave-from-class="opacity-100 scale-100 translate-y-0"
                    leave-to-class="opacity-0 scale-95 -translate-y-1"
                  >
                    <div
                      v-if="showDropdownId === link.id"
                      class="absolute right-2 top-10 bg-white border border-card-border rounded-[10px] shadow-xl py-2 w-44 z-30"
                      @click.stop
                    >
                      <button
                        @click="handleEditLink(link.id)"
                        class="w-full px-4 py-2.5 text-left font-inter text-[15px] text-text-primary hover:bg-page-bg transition-colors flex items-center gap-3"
                      >
                        <img src="@/components/icons/edit.svg" alt="" />
                        Редактировать
                      </button>
                      <div class="h-px bg-card-border mx-2 my-1"></div>
                      <button
                        @click="handleConfirmDelete(link.shortUrl)"
                        class="w-full px-4 py-2.5 text-left font-inter text-[15px] text-error hover:bg-red-50 transition-colors flex items-center gap-3"
                      >
                        <img src="@/components/icons/delete.svg" alt="" />
                        Удалить
                      </button>
                    </div>
                  </Transition>
                </td>
              </tr>
              <tr v-if="filteredLinks.length === 0">
                <td colspan="7" class="px-4 py-12 text-center">
                  <div class="font-inter text-[17px] text-text-secondary">
                    {{ searchQuery ? 'Ничего не найдено' : 'Список ссылок пуст' }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Пагинация: скругление только снизу -->
      <div
        v-if="filteredLinks.length > 0"
        class="bg-white border border-t-0 border-card-border rounded-b-card flex items-center justify-between px-4 py-3"
      >
        <div class="flex items-center gap-4">
          <select
            v-model.number="itemsPerPage"
            class="h-8 px-3 border border-card-border rounded-[8px] font-inter text-[17px] text-text-primary focus:outline-none focus:border-primary bg-white"
          >
            <option :value="1">1</option>
            <option :value="2">2</option>
            <option :value="3">3</option>
            <option :value="5">5</option>
            <option :value="10">10</option>
            <option :value="20">20</option>
          </select>
          <span class="font-inter text-[17px] text-text-secondary">элементов на страницу</span>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click.stop="prevPage"
            :disabled="currentPage === 1"
            class="flex items-center gap-1 px-3 py-1.5 text-[17px] text-text-secondary hover:text-text-primary disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            <img src="@/components/icons/arrow_left.svg" alt="" />
            Назад
          </button>
          <span
            class="px-3 py-1.5 font-inter text-[17px] text-text-primary bg-page-bg rounded-[8px] min-w-[60px] text-center"
          >
            {{ currentPage }}/{{ totalPages || 1 }}
          </span>
          <button
            @click.stop="nextPage"
            :disabled="currentPage === totalPages"
            class="flex items-center gap-1 px-3 py-1.5 text-[17px] text-text-secondary hover:text-text-primary disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            Вперёд
            <img src="@/components/icons/arrow_right.svg" alt="" />
          </button>
        </div>
      </div>
    </main>

    <!-- Попапы -->
    <ToastNotification
      v-if="deleteConfirmShow"
      :show="deleteConfirmShow"
      title="Удалить ссылку?"
      message="Это действие нельзя отменить"
      type="confirm"
      :show-actions="true"
      :show-icon="true"
      :buttons="deleteButtons"
      @close="deleteConfirmShow = false"
    />
    <ToastNotification
      v-if="toastShow && !deleteConfirmShow"
      :show="toastShow"
      :title="toastTitle"
      :message="toastMessage"
      :duration="3000"
      :type="toastType"
      :show-icon="true"
      @close="toastShow = false"
    />
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import ToastNotification from '@/components/ToastNotification.vue'

import lockIcon from '@/components/icons/lock.svg'
import lockedIcon from '@/components/icons/locked.svg'
import chevronUpIcon from '@/components/icons/chevron_up.svg'
import chevronDownIcon from '@/components/icons/chevron_down.svg'

// Константы
const API_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const SHORT_LINK_DOMAIN = import.meta.env.VITE_SHORT_LINK_DOMAIN || 'http://localhost'
const COOKIE_NAME = 'linqs_items_per_page'
const VALID_PAGE_SIZES = [1, 2, 3, 5, 10, 20] as const

// Хелперы для работы с куки
const setCookie = (name: string, value: string, days = 30): void => {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`
}

const getCookie = (name: string): string | null => {
  const matches = document.cookie.match(
    new RegExp(`(?:^|; )${name.replace(/([.$?*|{}()\[\]\\/+^])/g, '\\$1')}=([^;]*)`)
  )
  return matches && matches[1] ? decodeURIComponent(matches[1]) : null
}

const router = useRouter()

// Состояния UI
const showWIPPopup = ref(false)
const deleteConfirmShow = ref(false)
const linkCodeToDelete = ref<string | null>(null)
const toastShow = ref(false)
const toastTitle = ref('')
const toastMessage = ref('')
const toastType = ref<'success' | 'error' | 'info'>('success')

// Поиск и пагинация
const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref<(typeof VALID_PAGE_SIZES)[number]>(5)
const sortField = ref<'createdAt' | 'expiresAt' | null>(null)
const sortDirection = ref<'asc' | 'desc'>('asc')
const showDropdownId = ref<number | null>(null)

interface Link {
  id: number
  shortUrl: string
  originalUrl: string
  createdAt: string
  expiresAt: string | null
  comment: string
  isPrivate: boolean
}

const links = ref<Link[]>([])
const selectedLinks = ref<number[]>([])

// Вычисляемые свойства
const filteredLinks = computed(() => {
  if (!searchQuery.value) return links.value
  const query = searchQuery.value.toLowerCase()
  return links.value.filter(
    (link) =>
      link.originalUrl.toLowerCase().includes(query) ||
      link.shortUrl.toLowerCase().includes(query) ||
      link.comment.toLowerCase().includes(query)
  )
})

const sortedLinks = computed(() => {
  if (!sortField.value) return filteredLinks.value
  return [...filteredLinks.value].sort((a, b) => {
    const aVal = a[sortField.value!]
    const bVal = b[sortField.value!]
    if (!aVal && !bVal) return 0
    if (!aVal) return 1
    if (!bVal) return -1
    const comparison = aVal < bVal ? -1 : aVal > bVal ? 1 : 0
    return sortDirection.value === 'asc' ? comparison : -comparison
  })
})

const totalPages = computed(() => Math.ceil(sortedLinks.value.length / itemsPerPage.value))

const paginatedLinks = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return sortedLinks.value.slice(start, end)
})

// Watchers
const resetPagination = (): void => {
  currentPage.value = 1
}
watch([searchQuery, itemsPerPage], resetPagination)

watch(
  itemsPerPage,
  (newVal) => {
    if (newVal && VALID_PAGE_SIZES.includes(newVal)) {
      setCookie(COOKIE_NAME, String(newVal))
    }
  },
  { immediate: false }
)

// Методы
const toggleSort = (field: 'createdAt' | 'expiresAt'): void => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
}

const prevPage = (): void => {
  if (currentPage.value > 1) currentPage.value--
}

const nextPage = (): void => {
  if (currentPage.value < totalPages.value) currentPage.value++
}

const toggleSelect = (id: number): void => {
  const index = selectedLinks.value.indexOf(id)
  if (index > -1) {
    selectedLinks.value.splice(index, 1)
  } else {
    selectedLinks.value.push(id)
  }
}

const toggleSelectAll = (): void => {
  if (selectedLinks.value.length === paginatedLinks.value.length) {
    selectedLinks.value = []
  } else {
    selectedLinks.value = paginatedLinks.value.map((l) => l.id)
  }
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const copyLink = async (shortUrl: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(shortUrl)
    toastTitle.value = 'Готово!'
    toastMessage.value = 'Ссылка скопирована'
    toastType.value = 'success'
    toastShow.value = true
  } catch {
    toastTitle.value = 'Ошибка'
    toastMessage.value = 'Не удалось скопировать'
    toastType.value = 'error'
    toastShow.value = true
  }
}

const toggleDropdown = (id: number): void => {
  showDropdownId.value = showDropdownId.value === id ? null : id
}

const handleEditLink = (id: number): void => {
  showDropdownId.value = null
  editComment(id)
}

const handleConfirmDelete = (shortUrl: string): void => {
  showDropdownId.value = null
  confirmDelete(shortUrl)
}

const confirmDelete = (shortUrl: string): void => {
  const code = shortUrl.split('/').pop()
  if (!code) return
  linkCodeToDelete.value = code
  deleteConfirmShow.value = true
}

const handleDeleteConfirmed = async (): Promise<void> => {
  const code = linkCodeToDelete.value
  if (!code) return
  const token = localStorage.getItem('access_token')?.trim()
  if (!token) {
    router.push('/auth')
    return
  }
  try {
    const response = await fetch(`${API_URL}/links/${code}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) {
      toastTitle.value = 'Ошибка'
      toastMessage.value = 'Не удалось удалить'
      toastType.value = 'error'
      toastShow.value = true
      return
    }
    links.value = links.value.filter((l) => l.shortUrl.split('/').pop() !== code)
    toastTitle.value = 'Удалено'
    toastMessage.value = 'Ссылка удалена'
    toastType.value = 'success'
    toastShow.value = true
  } catch {
    toastTitle.value = 'Ошибка'
    toastMessage.value = 'Не удалось удалить'
    toastType.value = 'error'
    toastShow.value = true
  } finally {
    linkCodeToDelete.value = null
    deleteConfirmShow.value = false
  }
}

const togglePrivacy = (id: number): void => {
  const link = links.value.find((l) => l.id === id)
  if (link) link.isPrivate = !link.isPrivate
}

const editComment = (id: number): void => {
  showWIPPopup.value = true
}

const fetchLinks = async (): Promise<void> => {
  const token = localStorage.getItem('access_token')?.trim()
  if (!token) {
    router.push('/auth')
    return
  }
  try {
    const response = await fetch(`${API_URL}/links`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('access_token')
        router.push('/auth')
      }
      return
    }
    const data = await response.json()
    links.value = data?.map((item: any) => ({
      id: item.id,
      shortUrl: `${SHORT_LINK_DOMAIN}/${item.short_code}`,
      originalUrl: item.original_url?.replace(/^https?:\/\//, '') || '',
      createdAt: item.created_at,
      expiresAt: item.expires_at,
      comment: item.comment || '',
      isPrivate: item.is_private || false
    }))
  } catch (err) {
    console.error('💥 Network error:', err)
  }
}

// Кнопки для модалки удаления
const deleteButtons = computed(() => [
  {
    text: 'Отмена',
    variant: 'secondary' as const,
    action: cancelDelete
  },
  {
    text: 'Удалить',
    variant: 'danger' as const,
    action: handleDeleteConfirmed,
    closeAfter: true
  }
])

const cancelDelete = (): void => {
  deleteConfirmShow.value = false
  linkCodeToDelete.value = null
}

const handleLogout = (): void => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

// Инициализация
onMounted(() => {
  const saved = getCookie(COOKIE_NAME)
  if (saved) {
    const parsed = Number(saved)
    if (VALID_PAGE_SIZES.some((size) => size === parsed)) {
      itemsPerPage.value = parsed as (typeof VALID_PAGE_SIZES)[number]
    }
  }
  fetchLinks()
})
</script>