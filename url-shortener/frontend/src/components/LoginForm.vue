<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Email</label>
      <input
        v-model="form.email"
        type="email"
        placeholder="your@email.com"
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
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const API_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Router
const router = useRouter()

// Форма логина
const form = reactive({
  email: '',
  password: '',
  remember: false,
})

// Функция логина
const handleLogin = async () => {
  try {
    // 🔥 Используем переменную из .env
    const API_URL = import.meta.env.VITE_API_BASE_URL

    const response = await fetch(`${API_URL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: form.email,
        password: form.password,
      }),
    })

    if (!response.ok) {
      const errorText = await response.text()
      alert('Ошибка: ' + errorText)
      return
    }

    const data = await response.json()

    if (!data.access_token) {
      alert('Токен не получен')
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
      alert('Не удалось сохранить токен! Проверь настройки браузера')
      return
    }
    await router.push('/')
  } catch (err) {
    console.error(err)
    alert('Ошибка подключения к серверу')
  }
}

// Сабмит формы
const handleSubmit = () => {
  handleLogin()
}
</script>
