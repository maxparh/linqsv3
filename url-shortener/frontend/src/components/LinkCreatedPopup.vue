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
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path
                d="M15 5L5 15M5 5L15 15"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>

          <!-- Иконка успеха с анимацией -->
          <div
            class="w-16 h-16 mx-auto mb-4 rounded-full bg-success/10 flex items-center justify-center"
          >
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none" class="text-success">
              <path
                d="M6 16L13 23L26 9"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>

          <!-- Заголовок -->
          <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-6">
            Ссылка успешно создана!
          </h3>

          <!-- Сокращённая ссылка -->
          <div class="bg-page-bg rounded-[10px] p-4 mb-6">
            <div class="flex items-center justify-between gap-3">
              <span class="font-inter text-[17px] font-medium text-text-primary">
                {{ shortUrl }}
              </span>
              <button
                @click="copyLink"
                class="flex items-center gap-2 px-3 py-1.5 bg-white border border-card-border rounded-[8px] font-inter text-[14px] text-text-secondary hover:border-primary hover:text-primary transition-colors"
              >
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <path
                    d="M13.333 6H7.333C6.597 6 6 6.597 6 7.333V13.333C6 14.069 6.597 14.666 7.333 14.666H13.333C14.069 14.666 14.666 14.069 14.666 13.333V7.333C14.666 6.597 14.069 6 13.333 6Z"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                  <path
                    d="M2.667 10H2C1.264 10 0.667 9.403 0.667 8.667V2.667C0.667 1.931 1.264 1.334 2 1.334H8C8.736 1.334 9.333 1.931 9.333 2.667V3.334"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                {{ copied ? 'Скопировано!' : 'Копировать' }}
              </button>
            </div>
          </div>

          <!-- Кнопка закрыть -->
          <button
            @click="close"
            class="w-full h-10 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
          >
            Отлично!
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  shortUrl: string
}>()

const emit = defineEmits<{
  close: []
}>()

const visible = ref(false)
const animationClass = ref('')
const copied = ref(false)

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

const copyLink = async () => {
  try {
    await navigator.clipboard.writeText(props.shortUrl)
    copied.value = true
    setTimeout(() => (copied.value = false), 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// При появлении props.shortUrl показываем
watch(
  () => props.shortUrl,
  (newVal) => {
    if (newVal) {
      showWithAnimation()
    }
  },
  { immediate: true },
)
</script>
