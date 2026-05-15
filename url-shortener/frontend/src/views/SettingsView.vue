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
      <div class="mb-8">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Настройки</h1>
      </div>

      <!-- Профиль и Тарифы -->
      <div class="grid grid-cols-2 gap-6 mb-6">
        <!-- Профиль -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">Профиль</h2>
          
          <div class="flex gap-4 mb-6">
            <!-- Аватар -->
            <div class="w-16 h-16 rounded-full bg-page-bg flex items-center justify-center overflow-hidden">
              <!-- ICON: avatar_placeholder -->
              <div class="w-16 h-16 bg-gray-200 rounded-full">
                <img src="@/components/icons/avatar.png" alt="">
              </div>
            </div>
            
            <div class="flex-1">
              <div class="flex gap-8 mb-2">
                <div>
                  <div class="font-inter text-[14px] text-text-secondary mb-1">Имя</div>
                  <div class="font-inter text-[17px] text-text-primary font-medium">Александра</div>
                </div>
                <div>
                  <div class="font-inter text-[14px] text-text-secondary mb-1">Фамилия</div>
                  <div class="font-inter text-[17px] text-text-primary font-medium">Ермакова</div>
                </div>
              </div>
              <div>
                <div class="font-inter text-[14px] text-text-secondary mb-1">Тариф</div>
                <div class="font-inter text-[17px] text-primary font-medium">Базовый</div>
              </div>
            </div>
          </div>

          <h3 class="font-inter text-[17px] font-medium text-text-primary mb-4">Профиль</h3>
          
          <div class="grid grid-cols-2 gap-4 mb-6">
            <div class="h-16 px-4 border border-card-border rounded-[10px] bg-page-bg/50 flex flex-col justify-center">
              <div class="font-inter text-[14px] text-text-secondary mb-1">Электронная почта</div>
              <div class="font-inter text-[17px] text-text-primary">yourperfectemail@gmail.com</div>
            </div>
            <div class="h-16 px-4 border border-card-border rounded-[10px] bg-page-bg/50 flex flex-col justify-center">
              <div class="font-inter text-[14px] text-text-secondary mb-1">Номер телефона</div>
              <div class="font-inter text-[17px] text-text-primary">+79529525252</div>
            </div>
          </div>

          <button
            @click="showEditProfile = true"
            class="w-full h-10 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Редактировать
          </button>
        </div>

        <!-- Тарифы -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">Тарифы</h2>
          
          <div class="grid grid-cols-3 gap-4">
            <!-- Базовый -->
            <div class="border border-card-border rounded-[16px] p-4">
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">Базовый</div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">0 ₽ / мес.</div>
              
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
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">Профессиональный</div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">700 ₽ / мес.</div>
              
              <ul class="space-y-2 mb-6">
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">До 100 ссылок</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Расширенная аналитика</span>
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
              <div class="font-inter text-[17px] font-medium text-text-primary mb-2">Корпоративный</div>
              <div class="font-manrope font-bold text-[24px] text-text-primary mb-4">5000 ₽ / мес.</div>
              
              <ul class="space-y-2 mb-6">
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Неограниченные ссылки</span>
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
                  <span class="font-inter text-[14px] text-text-secondary">До 20 пользователей</span>
                </li>
                <li class="flex items-center gap-2">
                  <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                  <span class="font-inter text-[14px] text-text-secondary">Брендирование ссылок</span>
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
              <input
                v-model="settings.darkTheme"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
            </label>
          </div>

          <!-- Уведомления -->
          <div class="flex items-center justify-between py-2">
            <span class="font-inter text-[17px] text-text-primary">Уведомления об истекающих ссылках</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="settings.notifications"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
            </label>
          </div>

          <!-- 2FA -->
          <div class="flex items-center justify-between py-2">
            <span class="font-inter text-[17px] text-text-primary">Двухфакторная аутентификация</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="settings.twoFA"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-card-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
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

    <!-- Toast уведомления -->
    <ToastNotification
      v-if="toastShow"
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
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-error/10 flex items-center justify-center">
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
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import ToastNotification from '@/components/ToastNotification.vue'

const router = useRouter()

// Состояния UI
const toastShow = ref(false)
const toastTitle = ref('')
const toastMessage = ref('')
const toastType = ref<'success' | 'error' | 'info'>('success')
const deleteAccountConfirmShow = ref(false)
const showEditProfile = ref(false)

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