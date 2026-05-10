<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <!-- Почта -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Email</label>
      <input
        v-model="form.email"
        type="email"
        placeholder="your@email.com"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
        required
      />
    </div>

    <!-- Имя -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Имя</label>
      <input
        v-model="form.firstName"
        type="text"
        placeholder="Иван"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
        required
      />
    </div>

    <!-- Фамилия -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Фамилия</label>
      <input
        v-model="form.lastName"
        type="text"
        placeholder="Иванов"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
        required
      />
    </div>

    <!-- Телефон (необязательно) -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Телефон</label>
      <input
        v-model="form.phone"
        type="tel"
        placeholder="+7 (999) 123-45-67"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
      />
    </div>

    <!-- Пароль -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Пароль</label>
      <input
        v-model="form.password"
        type="password"
        placeholder="••••••••"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
        required
      />
    </div>

    <!-- Подтвердите пароль -->
    <div>
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5"
        >Подтвердите пароль</label
      >
      <input
        v-model="form.confirmPassword"
        type="password"
        placeholder="••••••••"
        class="w-full h-10 px-4 bg-white border border-[#E2E8F0] rounded-[10px] text-[17px] text-[#0F172A] placeholder:text-[#64748B] focus:outline-none focus:border-[#014751] transition-colors"
        required
      />
    </div>

    <!-- Согласие -->
    <label class="flex items-start gap-2 cursor-pointer">
      <input
        v-model="form.agree"
        type="checkbox"
        class="w-4 h-4 mt-0.5 rounded border-[#E2E8F0] text-[#014751] focus:ring-[#014751]"
        required
      />
      <span class="text-[14px] font-['Inter'] text-[#475569] leading-relaxed">
        Я согласен с
        <a href="#" class="text-[#014751] hover:underline">условиями использования</a> и
        <a href="#" class="text-[#014751] hover:underline">политикой конфиденциальности</a>
      </span>
    </label>

    <!-- Ошибка -->
    <div v-if="error" class="text-red-600 text-sm">{{ error }}</div>

    <!-- Кнопка -->
    <button
      type="submit"
      class="w-full h-10 bg-[#014751] text-white rounded-[10px] text-[17px] font-medium font-['Inter'] hover:bg-[#013d41] transition-colors"
      :disabled="loading"
    >
      {{ loading ? 'Обработка...' : 'Зарегистрироваться' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const error = ref('')
const loading = ref(false)

// 🔥 Объявляем API_URL здесь, чтобы он был виден во всей функции
const API_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  password: '',
  confirmPassword: '',
  agree: false,
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  // Валидация
  if (form.password !== form.confirmPassword) {
    error.value = 'Пароли не совпадают'
    loading.value = false
    return
  }

  try {
    // 1. РЕГИСТРАЦИЯ
    const registerRes = await fetch(`${API_URL}/api/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: form.email,
        password: form.password,
        // firstName/lastName можно добавить в бэкенд позже
      }),
    })

    if (!registerRes.ok) {
      const text = await registerRes.text()
      error.value = text || 'Ошибка регистрации'
      loading.value = false
      return
    }

    // 2. АВТО-ЛОГИН (чтобы получить токен)
    const loginResponse = await fetch(`${API_URL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: form.email,
        password: form.password,
      }),
    })

    // 🔥 ПРОВЕРКА: loginResponse определен в этой области видимости
    if (loginResponse.ok) {
      const loginData = await loginResponse.json()

      // Отладка в консоль
      if (import.meta.env.DEV) {
        console.log('Login data:', loginData)
      }

      if (loginData.access_token) {
        // 🔥 Сохраняем токен с правильным ключом
        localStorage.setItem('access_token', loginData.access_token)
        router.push('/')
      } else {
        console.error('No access_token in response')
        error.value = 'Ошибка: токен не получен'
      }
    } else {
      // Если авто-логин не сработал
      alert('Регистрация успешна! Теперь войдите.')
      router.push('/auth')
    }
  } catch (err) {
    console.error('Registration error:', err)
    error.value = 'Ошибка сети или сервера'
  } finally {
    loading.value = false
  }
}
</script>
