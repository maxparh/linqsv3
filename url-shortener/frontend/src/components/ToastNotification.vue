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
        @click.self="close"
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
            <!-- Иконка успеха -->
            <div
              class="w-16 h-16 mx-auto mb-4 rounded-full bg-success/10 flex items-center justify-center"
            >
              <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                <path
                  d="M6 16L13 23L26 9"
                  stroke="#10B981"
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
            <p class="font-inter text-[17px] text-text-secondary">
              {{ message }}
            </p>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  show: boolean
  title: string
  message: string
  duration?: number
}>()

const emit = defineEmits<{
  close: []
}>()

const visible = ref(false)
let timeoutId: ReturnType<typeof setTimeout> | null = null

const startTimer = () => {
  if (timeoutId) clearTimeout(timeoutId)
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

// Показ с запуском таймера
const showWithAnimation = () => {
  visible.value = true
  startTimer()
}

// Реакция на изменение props.show
watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      showWithAnimation()
    } else {
      close()
    }
  },
  { immediate: true } // ← Важно: сработает при первом монтировании
)
</script>