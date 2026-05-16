<template>
  <div
    class="bg-white rounded-card border border-card-border p-6 flex flex-col h-full relative group"
  >
    <!-- Шапка виджета -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <img
          src="@/components/icons/drag.svg"
          alt=""
          class="w-5 h-5 text-text-secondary opacity-40 cursor-grab active:cursor-grabbing"
        />
        <span class="font-inter text-[17px] font-medium text-text-secondary">{{ title }}</span>
      </div>
      <button
        @click="$emit('close')"
        class="text-text-secondary hover:text-text-primary transition-colors opacity-0 group-hover:opacity-100"
      >
        <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
      </button>
    </div>

    <!-- Контент виджета -->
    <div class="flex-1">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
}>()
defineEmits<{
  close: []
}>()
</script>
