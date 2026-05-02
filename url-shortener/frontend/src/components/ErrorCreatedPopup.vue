<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="visible"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 backdrop-blur-sm"
        @click.self="close"
      >
        <div
          class="bg-white rounded-card border border-card-border p-8 max-w-[420px] w-full text-center relative shadow-2xl"
          :class="animationClass"
        >
          <!-- Кнопка закрытия -->
          <button
            @click="close"
            class="absolute top-4 right-4 w-8 h-8 flex items-center justify-center text-text-secondary hover:text-text-primary transition-colors"
          >
            
          </button>

          <!-- Иконка ошибки -->
          <div
            class="w-16 h-16 mx-auto mb-4 rounded-full bg-red-100 flex items-center justify-center"
          >
            <svg width="32" height="32" viewBox="0 0 32 32" fill="white" class="text-red-600">
              
              <path d="M16 10V18" stroke="white" stroke-width="4" stroke-linecap="round" />
              <circle cx="16" cy="22" r="1.5" fill="white" />
            </svg>
          </div>

          <!-- Заголовок -->
          <h3 class="font-manrope font-bold text-[24px] text-red-600 mb-6">
            Сайт недоступен или не существует
          </h3>

          <!-- Сообщение или URL -->
          <div v-if="url" class="bg-page-bg rounded-[10px] p-4 mb-6">
            <div class="font-inter text-[17px] font-medium text-text-primary break-all">
              {{ url }}
            </div>
          </div>

          <!-- Кнопка закрыть -->
          <button
            @click="close"
            class="w-full h-10 bg-red-600 text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-red-700 transition-colors"
          >
            Понятно
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  url?: string
}>()

const emit = defineEmits<{
  close: []
}>()

const visible = ref(false)
const animationClass = ref('')

// Анимация входа
const showWithAnimation = () => {
  visible.value = true
  animationClass.value = 'animate-in zoom-in duration-300'
  setTimeout(() => {
    animationClass.value = ''
  }, 300)
}

// Закрытие с анимацией
const close = () => {
  animationClass.value = 'animate-out zoom-out duration-200'
  setTimeout(() => {
    visible.value = false
    animationClass.value = ''
    emit('close')
  }, 200)
}

// Показ при появлении props.url
watch(
  () => props.url,
  (newVal) => {
    if (newVal) {
      showWithAnimation()
    }
  },
  { immediate: true },
)
</script>
