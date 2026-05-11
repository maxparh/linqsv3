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
            <!-- Иконка: динамическая -->
            <div
              class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center"
              :class="type === 'error' ? 'bg-red-100' : 'bg-success/10'"
            >
              <!-- Иконка ошибки (красный крестик) -->
              <svg v-if="type === 'error'" width="32" height="32" viewBox="0 0 32 32" fill="none">
                <path
                  d="M16 8V16M16 20H16.01M8 16C8 11.5817 11.5817 8 16 8C20.4183 8 24 11.5817 24 16C24 20.4183 20.4183 24 16 24C11.5817 24 8 20.4183 8 16Z"
                  stroke="#EF4444"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              
              <!-- Иконка успеха (зелёная галочка) -->
              <svg v-else width="32" height="32" viewBox="0 0 32 32" fill="none">
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
  type?: 'success' | 'error'
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

const showWithAnimation = () => {
  visible.value = true
  startTimer()
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
  { immediate: true }
)
</script>