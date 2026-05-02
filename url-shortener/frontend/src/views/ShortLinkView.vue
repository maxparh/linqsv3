<template>
  <div class="flex justify-center items-center h-screen">
    <p class="text-text-secondary">Переходим по ссылке…</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router' // <- вот здесь

const route = useRoute()
const router = useRouter()

onMounted(async () => {
  const shortCode = route.params.shortCode
  try {
    const res = await fetch(`/api/resolve/${shortCode}`)
    if (!res.ok) {
      alert('Ссылка недоступна или не существует')
      return
    }
    const data = await res.json()
    const originalUrl = data.original_link

    // Редирект на оригинальный URL
    window.location.href = originalUrl
  } catch (err) {
    console.error(err)
    alert('Ошибка при переходе по ссылке')
  }
})
</script>