<template>
  <Teleport to="body">
    <dialog ref="dialog" aria-labelledby="profile-edit-title" :aria-busy="saving || processingAvatar"
      class="profile-dialog m-auto w-[calc(100%-2rem)] max-w-[480px] max-h-[calc(100dvh-2rem)] overflow-y-auto rounded-card border border-card-border bg-white p-6 sm:p-8 shadow-2xl backdrop:bg-black/50 backdrop:backdrop-blur-sm"
      @cancel.prevent="close" @click="handleBackdropClick">
      <button type="button" aria-label="Закрыть" title="Закрыть" :disabled="busy" @click="close"
        class="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-input hover:bg-page-bg disabled:opacity-50">
        <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
      </button>
      <h2 id="profile-edit-title" class="pr-8 mb-6 font-manrope font-bold text-[20px] text-text-primary">Редактировать профиль</h2>
      <form @submit.prevent="submit" class="space-y-5">
        <fieldset :disabled="busy" class="min-w-0 space-y-5">
          <div class="flex items-center gap-4">
            <div class="w-20 h-20 shrink-0 overflow-hidden rounded-full bg-primary flex items-center justify-center">
              <img v-if="draft.avatar" :src="draft.avatar" alt="Аватар профиля" class="w-full h-full object-cover" />
              <span v-else class="text-white font-inter font-semibold text-[28px]">{{ avatarInitial }}</span>
            </div>
            <div class="min-w-0 space-y-2">
              <input ref="avatarInput" type="file" accept="image/jpeg,image/png,image/webp" class="hidden" aria-label="Выбрать аватар" @change="selectAvatar" />
              <button type="button" @click="avatarInput?.click()" class="flex items-center gap-2 font-inter text-[15px] font-medium text-primary hover:underline">
                <img src="@/components/icons/edit.svg" alt="" class="w-4 h-4" />
                {{ processingAvatar ? 'Обработка...' : 'Изменить фото' }}
              </button>
              <button v-if="draft.avatar" type="button" @click="draft.avatar = ''; avatarError = ''" class="flex items-center gap-2 font-inter text-[14px] text-error hover:underline">
                <img src="@/components/icons/delete.svg" alt="" class="w-4 h-4" />Удалить фото
              </button>
            </div>
          </div>
          <p v-if="avatarError" role="alert" class="text-error font-inter text-[14px]">{{ avatarError }}</p>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <label class="block min-w-0 font-inter text-[14px] text-text-secondary">
              Имя
              <input ref="firstNameInput" v-model="draft.first_name" name="first_name" autocomplete="given-name" required maxlength="100" class="profile-input" />
            </label>
            <label class="block min-w-0 font-inter text-[14px] text-text-secondary">
              Фамилия
              <input v-model="draft.last_name" name="last_name" autocomplete="family-name" required maxlength="100" class="profile-input" />
            </label>
          </div>
          <label class="block font-inter text-[14px] text-text-secondary">
            Электронная почта
            <input v-model="draft.email" name="email" type="email" autocomplete="email" required maxlength="255" class="profile-input" />
          </label>
          <label class="block font-inter text-[14px] text-text-secondary">
            Номер телефона
            <input v-model="draft.phone" name="phone" type="tel" autocomplete="tel" maxlength="40" class="profile-input" />
          </label>
        </fieldset>
        <p v-if="error" role="alert" class="font-inter text-[14px] text-error">{{ error }}</p>
        <div class="grid grid-cols-2 gap-3 pt-1">
          <button type="button" :disabled="busy" @click="close" class="h-10 px-3 rounded-input border border-card-border bg-page-bg font-inter text-[15px] font-medium text-text-primary hover:bg-card-border/30 disabled:opacity-50">Отмена</button>
          <button type="submit" :disabled="busy || !!avatarError" class="h-10 px-3 rounded-input bg-primary font-inter text-[15px] font-medium text-white hover:bg-[#013d41] disabled:opacity-50">{{ saving ? 'Сохранение...' : 'Сохранить' }}</button>
        </div>
      </form>
    </dialog>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, reactive, ref } from 'vue'
import type { Profile } from '@/types/profile'

const props = defineProps<{ profile: Profile; saving: boolean; error: string }>()
const emit = defineEmits<{ close: []; save: [profile: Profile] }>()
const draft = reactive<Profile>({ ...props.profile })
const avatarInitial = computed(() => Array.from(draft.first_name.trim())[0]?.toUpperCase() || '')
const dialog = ref<HTMLDialogElement | null>(null)
const firstNameInput = ref<HTMLInputElement | null>(null)
const avatarInput = ref<HTMLInputElement | null>(null)
const processingAvatar = ref(false)
const avatarError = ref('')
const busy = computed(() => props.saving || processingAvatar.value)
let previousOverflow = ''

onMounted(() => {
  previousOverflow = document.body.style.overflow
  document.body.style.overflow = 'hidden'
  dialog.value?.showModal()
  firstNameInput.value?.focus()
})
onBeforeUnmount(() => {
  dialog.value?.close()
  document.body.style.overflow = previousOverflow
})
const close = () => { if (!busy.value) emit('close') }
const handleBackdropClick = (event: MouseEvent) => {
  const element = dialog.value
  if (!element || event.target !== element) return
  const rect = element.getBoundingClientRect()
  if (event.clientX < rect.left || event.clientX > rect.right || event.clientY < rect.top || event.clientY > rect.bottom) close()
}
const submit = () => {
  if (!busy.value && !avatarError.value) emit('save', { ...draft })
}
const selectAvatar = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  avatarError.value = ''
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type) || file.size > 5 * 1024 * 1024) {
    avatarError.value = 'Выберите JPEG, PNG или WebP размером до 5 МБ'
    return
  }
  processingAvatar.value = true
  let bitmap: ImageBitmap | undefined
  try {
    bitmap = await createImageBitmap(file)
    if (bitmap.width > 8192 || bitmap.height > 8192) throw new Error('Image too large')
    const canvas = document.createElement('canvas')
    canvas.width = canvas.height = 256
    const context = canvas.getContext('2d')
    if (!context) throw new Error('Canvas unavailable')
    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, 256, 256)
    const side = Math.min(bitmap.width, bitmap.height)
    context.drawImage(bitmap, (bitmap.width - side) / 2, (bitmap.height - side) / 2, side, side, 0, 0, 256, 256)
    draft.avatar = canvas.toDataURL('image/jpeg', 0.85)
  } catch {
    avatarError.value = 'Не удалось прочитать изображение. Выберите другой файл'
  } finally {
    bitmap?.close()
    processingAvatar.value = false
  }
}
</script>

<style scoped>
.profile-dialog[open] { animation: profile-appear 200ms ease-out; }
.profile-input {
  display: block;
  width: 100%;
  min-width: 0;
  height: 40px;
  margin-top: 6px;
  padding: 0 12px;
  border: 1px solid var(--color-card-border);
  border-radius: var(--radius-input);
  color: var(--color-text-primary);
  background: white;
  font-size: 17px;
}
.profile-input:focus { outline: 2px solid var(--color-primary); outline-offset: 1px; }
@keyframes profile-appear { from { opacity: 0; transform: translateY(12px) scale(.97); } to { opacity: 1; transform: none; } }
@media (prefers-reduced-motion: reduce) { .profile-dialog[open] { animation: none; } }
</style>
