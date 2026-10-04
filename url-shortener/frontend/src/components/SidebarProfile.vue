<template>
  <div class="flex items-center gap-3 px-4 py-3">
    <div class="w-8 h-8 shrink-0 overflow-hidden rounded-full bg-primary flex items-center justify-center">
      <img v-if="currentProfile?.avatar" :src="currentProfile.avatar" alt="" class="w-full h-full object-cover" />
      <span v-else class="text-white font-inter font-semibold text-[14px]">{{ initial }}</span>
    </div>
    <span class="min-w-0 break-words font-inter text-[17px] text-text-primary">{{ currentProfile?.first_name }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import type { Profile } from '@/types/profile'

const props = defineProps<{ profile?: Profile }>()
const loadedProfile = ref<Profile | null>(null)
const currentProfile = computed(() => props.profile ?? loadedProfile.value)
const initial = computed(() => Array.from(currentProfile.value?.first_name.trim() || '')[0]?.toUpperCase() || '')
const controller = new AbortController()

onMounted(async () => {
  if (props.profile) return
  const token = localStorage.getItem('access_token')?.trim()
  if (!token) return
  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || '/api'}/profile`, {
      headers: { Authorization: `Bearer ${token}` },
      signal: controller.signal,
    })
    if (response.ok) loadedProfile.value = await response.json()
  } catch {
    loadedProfile.value = null
  }
})
onBeforeUnmount(() => controller.abort())
</script>
