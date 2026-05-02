<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="visible"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 backdrop-blur-sm"
        @click.self="close"
      >
        <div
          class="bg-white rounded-card border border-card-border p-8 max-w-[320px] w-full text-center relative shadow-2xl"
          :class="animationClass"
        >
          <!-- Иконка успеха с анимацией -->
          <div
            class="w-16 h-16 mx-auto mb-4 rounded-full bg-success/10 flex items-center justify-center"
          >
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none" class="animate-pulse">
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
const animationClass = ref('')

let timeoutId: ReturnType<typeof setTimeout> | null = null

const startTimer = () => {
  if (timeoutId) clearTimeout(timeoutId)
  timeoutId = setTimeout(() => {
    close()
  }, props.duration || 2000)
}

const close = () => {
  animationClass.value = 'animate-out zoom-out duration-200'
  setTimeout(() => {
    visible.value = false
    animationClass.value = ''
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
    emit('close')
  }, 200)
}

// При появлении добавляем анимацию входа
const showWithAnimation = () => {
  visible.value = true
  animationClass.value = 'animate-in zoom-in duration-300'
  setTimeout(() => {
    animationClass.value = ''
  }, 300)
  startTimer()
}

if (props.show) {
  showWithAnimation()
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
)
</script>
