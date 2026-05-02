<template>
  <div class="min-h-screen bg-page-bg flex items-center justify-center p-6">
    <div
      class="w-full max-w-[480px] bg-white rounded-card border border-card-border p-8 text-center"
    >
      <!-- Неудача -->
      <div class="w-16 h-16 mx-auto mb-6 rounded-full bg-page-bg flex items-center justify-center">
        <img src="@/components/icons/oops_lock.svg" alt="" />
      </div>

      <!-- Заголовок -->
      <h1 class="font-manrope font-bold text-[32px] text-text-primary mb-3">Упс!</h1>

      <!-- Описание -->
      <p class="font-inter text-[17px] text-text-secondary mb-8">Для начала нужно авторизоваться</p>

      <!-- Кнопка -->
      <router-link
        to="/auth"
        class="inline-flex items-center justify-center py-2.5 px-8 bg-primary text-white rounded-input text-[17px] font-medium font-inter hover:bg-[#013d41] transition-colors"
      >
        Перейти к авторизации
      </router-link>
    </div>
  </div>
</template>
