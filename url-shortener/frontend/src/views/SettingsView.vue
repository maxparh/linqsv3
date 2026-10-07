<template>
  <div class="min-h-screen bg-page-bg flex flex-col md:flex-row">
    <!-- Боковое меню -->
    <aside class="w-full md:w-[224px] shrink-0 bg-white border-r border-card-border flex flex-col">
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
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/home_nav_nactive.svg" alt="" />
          Главная
        </router-link>

        <router-link
          to="/links"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/link_nav_nactive.svg" alt="" />
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
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] bg-primary text-white font-inter text-[17px] font-bold"
        >
          <img src="@/components/icons/settings_nav_active.svg" alt="" />
          Настройки
        </router-link>
      </nav>

      <!-- Профиль -->
      <div class="p-4 border-t border-card-border">
        <SidebarProfile :profile="profile" />
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
    <main class="flex-1 min-w-0 p-4 md:p-8">
      <!-- Хедер -->
      <div class="mb-8">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Настройки</h1>
      </div>

      <!-- Профиль и Тарифы -->
      <div class="grid grid-cols-1 2xl:grid-cols-2 gap-6 mb-6">
        <!-- Профиль -->
        <div :aria-busy="profileLoading" class="min-w-0 bg-white rounded-card border border-card-border p-6">
          <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">Профиль</h2>

          <p v-if="profileLoading" role="status" class="text-text-secondary">Загрузка профиля...</p>
          <div v-else-if="profileLoadError" role="alert" class="text-error">
            <p>{{ profileLoadError }}</p>
            <button type="button" @click="loadProfile" class="mt-3 text-primary underline">Повторить</button>
          </div>
          <template v-else>

          <div class="flex flex-col sm:flex-row gap-4 mb-6">
            <!-- Аватар -->
            <div
              class="w-16 h-16 shrink-0 rounded-full bg-primary flex items-center justify-center overflow-hidden"
            >
              <img v-if="profile.avatar" :src="profile.avatar" alt="Аватар профиля" class="w-full h-full object-cover" />
              <span v-else class="text-white font-inter font-semibold text-[24px]">{{ profileInitial }}</span>
            </div>

            <div class="flex-1 min-w-0">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-2">
                <label class="min-w-0">
                  <span class="block font-inter text-[14px] text-text-secondary mb-1">Имя</span>
                  <span class="block break-words font-inter text-[17px] text-text-primary font-medium">{{ profile.first_name || 'Не указано' }}</span>
                </label>
                <label class="min-w-0">
                  <span class="block font-inter text-[14px] text-text-secondary mb-1">Фамилия</span>
                  <span class="block break-words font-inter text-[17px] text-text-primary font-medium">{{ profile.last_name || 'Не указано' }}</span>
                </label>
              </div>
              <div>
                <div class="font-inter text-[14px] text-text-secondary mb-1">Тариф</div>
                <div class="font-inter text-[17px] text-primary font-medium">Базовый</div>
              </div>
            </div>
          </div>

          <h3 class="font-inter text-[17px] font-medium text-text-primary mb-4">Профиль</h3>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
            <label
              class="min-w-0 min-h-16 px-4 py-3 border border-card-border rounded-[10px] bg-page-bg/50 flex flex-col justify-center"
            >
              <span class="font-inter text-[14px] text-text-secondary mb-1">Электронная почта</span>
              <span class="break-all font-inter text-[17px] text-text-primary">{{ profile.email }}</span>
            </label>
            <label
              class="min-w-0 min-h-16 px-4 py-3 border border-card-border rounded-[10px] bg-page-bg/50 flex flex-col justify-center"
            >
              <span class="font-inter text-[14px] text-text-secondary mb-1">Номер телефона</span>
              <span class="break-words font-inter text-[17px] text-text-primary">{{ profile.phone || 'Не указан' }}</span>
            </label>
          </div>

          <button
            type="button"
            ref="editProfileButton"
            @click="startEditProfile"
            class="w-full h-10 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Редактировать
          </button>
          </template>
        </div>

        <!-- Тарифы -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">Тарифы</h2>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <!-- Базовый -->
            <div class="border border-card-border rounded-[16px] p-4">
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">Базовый</div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">
                0 ₽ / мес.
              </div>

              <ul class="space-y-2 mb-6">
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">До 20 ссылок</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Базовая аналитика</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Генерация QR-кодов</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">1 пользователь</span>
                </li>
              </ul>

              <button
                disabled
                class="w-full h-10 border border-primary text-primary rounded-[10px] font-inter text-[17px] font-medium opacity-60 cursor-not-allowed"
              >
                Текущий
              </button>
            </div>

            <!-- Профессиональный -->
            <div class="border border-card-border rounded-[16px] p-4">
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">
                Профессиональный
              </div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">
                700 ₽ / мес.
              </div>

              <ul class="space-y-2 mb-6">
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">До 100 ссылок</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary"
                    >Расширенная аналитика</span
                  >
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Генерация QR-кодов</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">До 5 пользователей</span>
                </li>
              </ul>

              <button
                @click="selectTariff('professional')"
                class="w-full h-10 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
              >
                Выбрать
              </button>
            </div>

            <!-- Корпоративный -->
            <div class="border border-card-border rounded-[16px] p-4">
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">
                Корпоративный
              </div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">
                5000 ₽ / мес.
              </div>

              <ul class="space-y-2 mb-6">
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary"
                    >Неограниченные ссылки</span
                  >
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Базовая аналитика</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Генерация QR-кодов</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary"
                    >До 20 пользователей</span
                  >
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary"
                    >Брендирование ссылок</span
                  >
                </li>
              </ul>

              <button
                @click="selectTariff('corporate')"
                class="w-full h-10 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
              >
                Выбрать
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Другие настройки -->
      <div class="bg-white rounded-card border border-card-border p-6">
        <h2 class="font-inter text-[17px] font-medium text-text-primary mb-6">Другие настройки</h2>

        <div class="space-y-4">
          <!-- Тёмная тема -->
          <div class="flex items-center justify-between py-2">
            <span class="font-inter text-[17px] text-text-primary">Тёмная тема</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="settings.darkTheme" type="checkbox" class="sr-only peer" />
              <div
                class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"
              ></div>
            </label>
          </div>

          <!-- Уведомления -->
          <div class="flex items-center justify-between py-2">
            <span class="font-inter text-[17px] text-text-primary"
              >Уведомления об истекающих ссылках</span
            >
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="settings.notifications" type="checkbox" class="sr-only peer" />
              <div
                class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"
              ></div>
            </label>
          </div>

          <!-- 2FA -->
          <div class="flex items-center justify-between py-2">
            <span class="font-inter text-[17px] text-text-primary"
              >Двухфакторная аутентификация</span
            >
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="settings.twoFA" type="checkbox" class="sr-only peer" />
              <div
                class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"
              ></div>
            </label>
          </div>
        </div>

        <!-- Удалить аккаунт -->
        <div class="mt-8 flex justify-center">
          <button
            @click="confirmDeleteAccount"
            class="h-10 px-6 bg-error text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-red-600 transition-colors"
          >
            Удалить аккаунт
          </button>
        </div>
      </div>
    </main>

    <ProfileEditPopup v-if="showEditProfile" ref="profilePopup" :profile="profile" :saving="profileSaving" :error="profileSaveError" @close="cancelEditProfile" @save="saveProfile" />

    <!-- Toast уведомления -->
    <ToastNotification
      :show="toastShow"
      :title="toastTitle"
      :message="toastMessage"
      :duration="3000"
      :type="toastType"
      :show-icon="true"
      @close="toastShow = false"
    />

    <!-- Попап подтверждения удаления аккаунта -->
    <div
      v-if="deleteAccountConfirmShow"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="deleteAccountConfirmShow = false"
    >
      <div class="bg-white rounded-card border border-card-border p-8 max-w-[400px] text-center">
        <div
          class="w-16 h-16 mx-auto mb-4 rounded-full bg-error/10 flex items-center justify-center"
        >
          <!-- ICON: warning -->
          <div class="w-8 h-8 bg-error rounded">warning</div>
        </div>
        <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-2">Удалить аккаунт?</h3>
        <p class="font-inter text-[17px] text-text-secondary mb-6">
          Это действие нельзя отменить. Все ваши данные будут безвозвратно удалены.
        </p>
        <div class="flex gap-3 justify-center">
          <button
            @click="deleteAccountConfirmShow = false"
            class="h-10 px-6 border border-card-border text-text-primary rounded-[10px] font-inter text-[17px] font-medium hover:bg-page-bg transition-colors"
          >
            Отмена
          </button>
          <button
            @click="handleDeleteAccount"
            class="h-10 px-6 bg-error text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-red-600 transition-colors"
          >
            Удалить
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import ToastNotification from '@/components/ToastNotification.vue'
import ProfileEditPopup from '@/components/ProfileEditPopup.vue'
import SidebarProfile from '@/components/SidebarProfile.vue'
import type { Profile } from '@/types/profile'

const router = useRouter()

// Состояния UI
const toastShow = ref(false)
const toastTitle = ref('')
const toastMessage = ref('')
const toastType = ref<'success' | 'error' | 'info'>('success')
const deleteAccountConfirmShow = ref(false)
const showEditProfile = ref(false)
const profilePopup = ref<InstanceType<typeof ProfileEditPopup> | null>(null)

const API_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const profile = reactive<Profile>({ first_name: '', last_name: '', email: '', phone: '', avatar: '' })
const profileLoading = ref(true)
const profileSaving = ref(false)
const profileLoadError = ref('')
const profileSaveError = ref('')
const editProfileButton = ref<HTMLButtonElement | null>(null)
const profileInitial = computed(() => Array.from(profile.first_name.trim())[0]?.toUpperCase() || '')

const requestProfile = async (method: 'GET' | 'PUT', body?: Profile): Promise<Profile> => {
  const token = localStorage.getItem('access_token')?.trim()
  if (!token) {
    await router.push('/auth')
    throw new Error('Войдите в аккаунт')
  }
  const response = await fetch(`${API_URL}/profile`, {
    method,
    headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (response.status === 401) {
    localStorage.removeItem('access_token')
    await router.push('/auth')
    throw new Error('Сессия истекла. Войдите снова')
  }
  if (!response.ok) {
    const data = await response.json().catch(() => ({}))
    const messages: Record<string, string> = {
      'email already in use': 'Эта почта уже используется другим аккаунтом',
      'phone already in use': 'Этот телефон уже используется другим аккаунтом',
      'invalid profile: names must contain 1 to 100 characters': 'Имя и фамилия должны содержать от 1 до 100 символов',
      'invalid profile: invalid email': 'Укажите корректную электронную почту',
      'invalid profile: invalid phone': 'Укажите корректный телефон с кодом страны',
      'invalid profile: invalid avatar': 'Не удалось сохранить фото. Выберите другое изображение',
      'user not found': 'Профиль не найден',
    }
    throw new Error(messages[data.error] || 'Не удалось выполнить запрос. Попробуйте ещё раз')
  }
  const data = await response.json()
  return { first_name: data.first_name, last_name: data.last_name, email: data.email, phone: data.phone || '', avatar: data.avatar || '' }
}

const loadProfile = async () => {
  profileLoading.value = true
  profileLoadError.value = ''
  try {
    Object.assign(profile, await requestProfile('GET'))
  } catch (error) {
    profileLoadError.value = error instanceof Error ? error.message : 'Не удалось загрузить профиль'
  } finally {
    profileLoading.value = false
  }
}

const startEditProfile = () => {
  profileSaveError.value = ''
  showEditProfile.value = true
}

const cancelEditProfile = async () => {
  await profilePopup.value?.closeWithAnimation()
  showEditProfile.value = false
  profileSaveError.value = ''
  await nextTick()
  editProfileButton.value?.focus()
}

const saveProfile = async (profileDraft: Profile) => {
  if (!showEditProfile.value || profileSaving.value) return
  profileSaving.value = true
  profileSaveError.value = ''
  try {
    const updated = await requestProfile('PUT', {
      first_name: profileDraft.first_name.trim(),
      last_name: profileDraft.last_name.trim(),
      email: profileDraft.email.trim(),
      phone: profileDraft.phone.trim(),
      avatar: profileDraft.avatar,
    })
    Object.assign(profile, updated)
    await cancelEditProfile()
    toastTitle.value = 'Профиль обновлён'
    toastMessage.value = 'Изменения сохранены'
    toastType.value = 'success'
    toastShow.value = true
  } catch (error) {
    profileSaveError.value = error instanceof Error ? error.message : 'Не удалось сохранить профиль'
  } finally {
    profileSaving.value = false
  }
}

onMounted(loadProfile)

// Настройки
const settings = reactive({
  darkTheme: false,
  notifications: true,
  twoFA: true,
})

// Выбор тарифа
const selectTariff = (tariff: string) => {
  toastTitle.value = 'Тариф изменён'
  toastMessage.value = `Вы выбрали тариф "${tariff === 'professional' ? 'Профессиональный' : 'Корпоративный'}"`
  toastType.value = 'success'
  toastShow.value = true
}

// Подтверждение удаления аккаунта
const confirmDeleteAccount = () => {
  deleteAccountConfirmShow.value = true
}

// Удаление аккаунта
const handleDeleteAccount = async () => {
  // Здесь будет вызов API
  toastTitle.value = 'Аккаунт удалён'
  toastMessage.value = 'Ваш аккаунт был успешно удалён'
  toastType.value = 'success'
  toastShow.value = true
  deleteAccountConfirmShow.value = false

  setTimeout(() => {
    router.push('/auth')
  }, 1500)
}

// Выход
const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}
</script>
