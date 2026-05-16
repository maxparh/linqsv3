<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="visible"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 backdrop-blur-sm"
        @click.self="handleOverlayClick"
      >
        <!-- Модалка с анимацией всплытия -->
        <Transition
          enter-active-class="transition-all duration-300 ease-out"
          enter-from-class="opacity-0 translate-y-4 scale-95"
          enter-to-class="opacity-100 translate-y-0 scale-100"
          leave-active-class="transition-all duration-200 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-show="visible"
            class="bg-white rounded-card border border-card-border p-8 max-w-[320px] w-full text-center relative shadow-2xl"
          >
            <!-- Иконка: динамическая -->
            <div
              v-if="showIcon"
              class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center"
              :class="iconBgClass"
            >
              <!-- Иконка ошибки / подтверждения (красный крестик) -->
              <svg
                v-if="type === 'error' || type === 'confirm'"
                width="32"
                height="32"
                viewBox="0 0 32 32"
                fill="none"
              >
                <path
                  d="M16 8V16M16 20H16.01M8 16C8 11.5817 11.5817 8 16 8C20.4183 8 24 11.5817 24 16C24 20.4183 20.4183 24 16 24C11.5817 24 8 20.4183 8 16Z"
                  :stroke="iconColor"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>

              <!-- Иконка успеха (зелёная галочка) -->
              <svg
                v-else-if="type === 'success'"
                width="32"
                height="32"
                viewBox="0 0 32 32"
                fill="none"
              >
                <path
                  d="M6 16L13 23L26 9"
                  stroke="#10B981"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>

              <!-- Иконка инфо (синяя "i") -->
              <svg
                v-else-if="type === 'info'"
                width="32"
                height="32"
                viewBox="0 0 32 32"
                fill="none"
              >
                <path
                  d="M16 8V16M16 20H16.01M8 16C8 11.5817 11.5817 8 16 8C20.4183 8 24 11.5817 24 16C24 20.4183 20.4183 24 16 24C11.5817 24 8 20.4183 8 16Z"
                  stroke="#3B82F6"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>

            <!-- Заголовок -->
            <h3 class="font-manrope font-bold text-[20px] text-text-primary mb-2">
              {{ title }}
            </h3>

            <!-- Описание -->
            <p class="font-inter text-[17px] text-text-secondary mb-6">
              {{ message }}
            </p>

            <!-- 🔥 Кнопки действий (если есть) -->
            <div
              v-if="showActions && buttons?.length"
              class="flex items-center justify-center gap-3"
            >
              <button
                v-for="(btn, idx) in buttons"
                :key="idx"
                @click="handleButtonClick(btn)"
                class="h-10 px-5 rounded-[10px] font-inter text-[15px] font-medium transition-colors hover:opacity-90 active:scale-95"
                :class="getButtonClass(btn.variant)"
              >
                {{ btn.text }}
              </button>
            </div>

            <!-- 🔥 Кнопка закрытия (если нет кнопок действий) -->
            <button
              v-else-if="showCloseButton"
              @click="close"
              class="absolute top-4 right-4 w-8 h-8 flex items-center justify-center text-text-secondary hover:text-text-primary transition-colors"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <path
                  d="M6 6L18 18M6 18L18 6"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
              </svg>
            </button>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

// 🔥 Тип для кнопки тоста
export interface ToastButton {
  text: string
  variant?: 'primary' | 'secondary' | 'danger'
  action?: () => void
  closeAfter?: boolean
}

// 🔥 ЕДИНЫЙ вызов defineProps + withDefaults (исправлено!)
const props = withDefaults(
  defineProps<{
    show: boolean
    title: string
    message: string
    duration?: number
    type?: 'success' | 'error' | 'confirm' | 'info'
    buttons?: ToastButton[]
    showActions?: boolean
    showCloseButton?: boolean
    showIcon?: boolean
  }>(),
  {
    duration: 3000,
    type: 'success',
    buttons: () => [],
    showActions: false,
    showCloseButton: true,
    showIcon: true,
  },
)

const emit = defineEmits<{
  close: []
  confirm: []
  cancel: []
  action: [index: number]
}>()

const visible = ref(false)
let timeoutId: ReturnType<typeof setTimeout> | null = null

// 🔥 Вычисляемые стили для иконки
const iconBgClass = computed(() => {
  switch (props.type) {
    case 'error':
    case 'confirm':
      return 'bg-red-100'
    case 'success':
      return 'bg-success/10'
    case 'info':
      return 'bg-primary/10'
    default:
      return 'bg-page-bg'
  }
})

const iconColor = computed(() => {
  return props.type === 'error' || props.type === 'confirm' ? '#EF4444' : '#10B981'
})

// 🔥 Стили для кнопок
const getButtonClass = (variant: ToastButton['variant'] = 'primary') => {
  const base = 'hover:opacity-90 active:scale-95 transition-all'
  switch (variant) {
    case 'danger':
      return `${base} bg-error text-white hover:bg-[#dc2626]`
    case 'secondary':
      return `${base} bg-page-bg text-text-primary border border-card-border hover:bg-card-border/30`
    case 'primary':
    default:
      return `${base} bg-primary text-white hover:bg-[#013d41]`
  }
}

// 🔥 Обработчики
const startTimer = () => {
  if (timeoutId) clearTimeout(timeoutId)
  // Не запускаем таймер, если есть кнопки действий
  if (props.showActions && props.buttons?.length) return

  timeoutId = setTimeout(() => {
    close()
  }, props.duration || 3000)
}

const close = () => {
  visible.value = false
  if (timeoutId) {
    clearTimeout(timeoutId)
    timeoutId = null
  }
  emit('close')
}

const showWithAnimation = () => {
  visible.value = true
  startTimer()
}

const handleOverlayClick = () => {
  // Если есть кнопки — не закрывать по клику на оверлей
  if (props.showActions && props.buttons?.length) return
  close()
}

const handleButtonClick = (btn: ToastButton) => {
  if (btn.action) btn.action()

  // Эмитим события для удобства
  if (btn.text.toLowerCase().includes('отмена') || btn.text.toLowerCase().includes('cancel')) {
    emit('cancel')
  } else if (
    btn.text.toLowerCase().includes('удалить') ||
    btn.text.toLowerCase().includes('confirm') ||
    btn.text.toLowerCase().includes('подтвердить')
  ) {
    emit('confirm')
  }
  emit('action', props.buttons?.indexOf(btn) || 0)

  // Закрываем тост, если не указано иное
  if (btn.closeAfter !== false) {
    close()
  }
}

watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      showWithAnimation()
    } else {
      close()
    }
  },
  { immediate: true },
)
</script>
