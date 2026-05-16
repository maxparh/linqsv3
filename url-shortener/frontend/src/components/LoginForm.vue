<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Email</label>
      <input
        v-model="form.identifier"
        type="text"
        placeholder="Ваш email или телефон"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
      />
    </div>

    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Пароль</label>
      <input
        v-model="form.password"
        type="password"
        placeholder="••••••••"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
      />
    </div>

    <div class="flex items-center justify-between">
      <label class="flex items-center gap-2 cursor-pointer">
        <input
          v-model="form.remember"
          type="checkbox"
          class="w-4 h-4 rounded border-[#E2E8F0] text-[#014751] focus:ring-[#014751]"
        />
        <span class="text-[14px] font-['Inter'] text-[#475569]">Запомнить меня</span>
      </label>
      <a href="#" class="text-[14px] font-['Inter'] text-[#014751] hover:underline"
        >Забыли пароль?</a
      >
    </div>

    <button
      type="submit"
      class="w-full h-10 bg-[#014751] text-white rounded-[10px] text-[17px] font-medium font-['Inter'] hover:bg-[#013d41] transition-colors"
    >
      Войти
    </button>
  </form>

  <ToastNotification
    :show="toastShow"
    :title="toastTitle"
    :message="toastMessage"
    :duration="3000"
    type="error"
    @close="toastShow = false"
  />
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import ToastNotification from '@/components/ToastNotification.vue' // ✅ Проверь путь

// Состояние для модального уведомления
const toastShow = ref(false)
const toastTitle = ref('')
const toastMessage = ref('')

// Вспомогательная функция для показа
const showToast = (title: string, message: string) => {
  toastTitle.value = title
  toastMessage.value = message
  toastShow.value = true
  // Компонент сам закроется через duration (3 сек)
}

const API_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Router
const router = useRouter()

// Форма логина
const form = reactive({
  identifier: '',
  password: '',
  remember: false,
})

// Функция логина
const handleLogin = async () => {
  try {
    const response = await fetch(`${API_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        identifier: form.identifier,
        password: form.password,
      }),
    })

    if (!response.ok) {
      const errorText = await response.text()
      showToast('Ошибка входа', 'Неверный email/телефон или пароль')
      return
    }

    const data = await response.json()

    if (!data.access_token) {
      showToast('Ошибка', 'Токен не получен')
      return
    }

    // Принудительно ждем записи
    await new Promise((resolve) => setTimeout(resolve, 100))
    localStorage.setItem('access_token', data.access_token)
    // Проверяем что записалось
    const check = localStorage.getItem('access_token')
    if (import.meta.env.DEV) {
      console.log('Token saved?', !!check)
    }
    if (!check) {
      showToast('Ошибка', 'Не удалось сохранить токен')
      return
    }
    await router.push('/')
  } catch (err) {
    console.error(err)
    showToast('Ошибка сети', 'Проверь подключение к интернету')
  }
}

// Сабмит формы
const handleSubmit = () => {
  handleLogin()
}
</script>
