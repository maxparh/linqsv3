<template>
  <div
    class="bg-white rounded-card border border-card-border p-6 flex flex-col min-w-0 relative group"
  >
    <!-- Шапка виджета -->
    <div class="flex items-center justify-between gap-3 mb-4">
      <div class="min-w-0 flex items-center gap-2">
        <button v-if="draggable" type="button" class="widget-drag-handle w-8 h-8 shrink-0 flex items-center justify-center cursor-grab active:cursor-grabbing touch-none"
          :aria-label="`Переместить виджет: ${title}`" title="Переместить виджет (стрелки с клавиатуры)"
          @keydown.left.prevent="$emit('move', -1)" @keydown.up.prevent="$emit('move', -1)"
          @keydown.right.prevent="$emit('move', 1)" @keydown.down.prevent="$emit('move', 1)">
          <img src="@/components/icons/drag.svg" alt="" draggable="false" class="w-5 h-5 opacity-60" />
        </button>
        <span class="font-inter text-[17px] font-medium text-text-secondary">{{ title }}</span>
      </div>
      <button
        v-if="removable !== false"
        type="button"
        :aria-label="`Убрать виджет: ${title}`"
        title="Убрать виджет"
        @click="$emit('close')"
        class="w-8 h-8 shrink-0 flex items-center justify-center text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
      >
        <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
      </button>
    </div>

    <!-- Контент виджета -->
    <div class="min-w-0 flex-1 flex flex-col">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  removable?: boolean
  draggable?: boolean
}>(), {
  removable: true,
  draggable: false,
})
defineEmits<{
  close: []
  move: [direction: number]
}>()
</script>
