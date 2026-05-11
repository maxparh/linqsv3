<template>
  <div class="flex justify-center items-center h-screen">
    <p class="text-text-secondary">Переходим по ссылке…</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

onMounted(async () => {
  const shortCode = route.params.shortCode
  
  if (!shortCode || typeof shortCode !== 'string') {
    alert('Неверный формат ссылки')
    router.push('/')
    return
  }
  
  try {
    const res = await fetch(`/api/resolve/${shortCode}`)
    if (!res.ok) {
      alert('Ссылка недоступна или не существует')
      router.push('/')
      return
    }
    const data = await res.json()
    
    // 🔥 Исправлено: original_link → original_url
    const originalUrl = data.original_url
    
    if (!originalUrl) {
      alert('Неверный формат ответа сервера')
      router.push('/')
      return
    }
    
    window.location.href = originalUrl
  } catch (err) {
    console.error(err)
    alert('Ошибка при переходе по ссылке')
    router.push('/')
  }
})
</script>
