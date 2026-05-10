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
    <div class="relative z-10">
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Телефон</label>
      <VueTelInput
        v-model="phoneDisplay"
        :defaultCountry="'RU'"
        :autoFormat="true"
        :validCharactersOnly="true"
        :placeholder="'\u002B7 (999) 123-45-67'"
        @input="handleTelInput"
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
      <label class="block text-[14px] font-['Inter'] text-[#475569] mb-1.5">Подтвердите пароль</label>
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
import { VueTelInput } from 'vue-tel-input'
import 'vue-tel-input/vue-tel-input.css'

const router = useRouter()
const error = ref('')
const loading = ref(false)
const phoneDisplay = ref('')  // Форматированный номер для отображения
const phoneE164 = ref('')     // Чистый номер для API (+79991234567)

const API_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',  // Отправляем в БД
  password: '',
  confirmPassword: '',
  agree: false,
})

// 🔥 Безопасный парсер: принимает ЛЮБОЙ тип, возвращает string | null
const parsePhoneToE164 = (input: unknown): string | null => {
  // 1. Приводим к строке безопасно
  let value: string
  if (input == null) return null
  if (typeof input === 'string') {
    value = input
  } else if (typeof input === 'number') {
    value = String(input)
  } else if (input instanceof Event) {
    // Нативное событие — берем значение из target
    const target = input.target as HTMLInputElement
    value = target?.value || ''
  } else {
    // Попытка вызвать toString
    try {
      value = String(input)
    } catch {
      return null
    }
  }

  // 2. Пустая строка
  if (!value.trim()) return null

  // 3. Оставляем только цифры и +
  const digits = value.replace(/[^\d+]/g, '')
  if (!digits) return null

  // 4. Нормализация в E.164
  if (digits.startsWith('8') && digits.length === 11) {
    return '+7' + digits.slice(1)
  }
  if (digits.startsWith('+7')) {
    return digits
  }
  if (digits.startsWith('7') && digits.length === 11) {
    return '+' + digits
  }
  if (digits.length === 10) {
    return '+7' + digits
  }

  return null // Не валидный формат
}

// 🔥 Обработчик @input для vue-tel-input v9
// Может прийти: string, Event, или объект — парсер всё обработает
const handleTelInput = (payload: unknown) => {
  const result = parsePhoneToE164(payload)
  phoneE164.value = result || ''
  console.log('📱 handleTelInput:', { payload: typeof payload, result: phoneE164.value })
}

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  if (form.password !== form.confirmPassword) {
    error.value = 'Пароли не совпадают'
    loading.value = false
    return
  }

  try {
    // 🔥 Финальный парсинг прямо перед отправкой (двойная защита)
    const finalPhone = parsePhoneToE164(phoneDisplay.value) || phoneE164.value || null
    const phonePayload = finalPhone ? finalPhone : undefined

    console.log('🔍 phoneDisplay:', phoneDisplay.value)
    console.log('🔍 phoneE164:', phoneE164.value)
    console.log('🔍 finalPhone:', finalPhone)
    console.log('🔍 phonePayload:', phonePayload)

    // 1. РЕГИСТРАЦИЯ
    const registerRes = await fetch(`${API_URL}/api/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: form.email,
        password: form.password,
        phone: phonePayload,
      }),
    })

    if (!registerRes.ok) {
      const text = await registerRes.text()
      error.value = text || 'Ошибка регистрации'
      loading.value = false
      return
    }

    // 2. АВТО-ЛОГИН
    const loginResponse = await fetch(`${API_URL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        identifier: form.email,
        password: form.password,
      }),
    })

    if (loginResponse.ok) {
      const loginData = await loginResponse.json()
      if (import.meta.env.DEV) {
        console.log('Login ', loginData)
      }

      if (loginData.access_token) {
        localStorage.setItem('access_token', loginData.access_token)
        router.push('/')
      } else {
        console.error('No access_token in response')
        error.value = 'Ошибка: токен не получен'
      }
    } else {
      alert('Регистрация успешна! Теперь войдите.')
      router.push('/auth')
    }
  } catch (err) {
    console.error('Registration error:', err)
    error.value = 'Ошибка сети или сервера'
  } finally {
    loading.value = false
    phoneDisplay.value = ''
    phoneE164.value = ''
    form.phone = ''
  }
}
</script>

<style scoped>
:deep(.vue-tel-input) { border: none !important; box-shadow: none !important; }
:deep(.vue-tel-input input) {
  height: 40px !important; border: 1px solid #E2E8F0 !important;
  border-radius: 10px !important; background: white !important;
  font-family: 'Inter', sans-serif !important; font-size: 17px !important;
  color: #0F172A !important; padding: 0 16px !important;
  transition: border-color 0.2s !important;
}
:deep(.vue-tel-input input:focus) { outline: none !important; border-color: #014751 !important; }
:deep(.vue-tel-input input::placeholder) { color: #64748B !important; }
:deep(.vue-tel-input .dropdown) {
  z-index: 50 !important; border-radius: 10px !important;
  border: 1px solid #E2E8F0 !important; font-family: 'Inter', sans-serif !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
  max-height: 200px !important; overflow-y: auto !important;
}
:deep(.vue-tel-input .selection) {
  height: 40px !important; border-right: 1px solid #E2E8F0 !important;
  background: white !important;
}
</style>