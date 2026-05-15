<!-- src/views/CookiePolicy.vue -->
<template>
  <div class="min-h-screen bg-page-bg flex items-start justify-center p-6">
    <div class="w-full max-w-[720px] bg-white rounded-card border border-card-border p-6 mt-8">
      
      <!-- Логотип -->
      <div class="p-6 flex items-center gap-3 pb-[80px]">
        <div class="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
          <img src="@/components/icons/linqs_logo.svg" alt="" />
        </div>
        <span class="font-inter font-semibold text-[17px] text-text-primary">Linqs</span>
      </div>

      <!-- Заголовок -->
      <h1 class="font-manrope font-bold text-[24px] text-text-primary text-center mb-6">
        Политика использования Cookies
      </h1>

      <!-- Контент -->
      <div class="space-y-5">
        <section>
          <h2 class="font-inter text-[17px] font-medium text-text-primary mb-2">1. Что такое cookies?</h2>
          <p class="font-inter text-[17px] text-text-secondary">
            Cookies — это небольшие текстовые файлы, которые сохраняются на вашем устройстве при посещении сайта. 
            Они помогают сайту запоминать ваши действия и предпочтения.
          </p>
        </section>

        <section>
          <h2 class="font-inter text-[17px] font-medium text-text-primary mb-2">2. Какие cookies мы используем?</h2>
          <ul class="space-y-3">
            <li class="flex gap-3">
              <span class="font-inter text-[17px] font-medium text-text-primary min-w-[140px]">Необходимые</span>
              <span class="font-inter text-[17px] text-text-secondary">
                Требуются для работы сайта: аутентификация, безопасность, сохранение сессии. Нельзя отключить.
              </span>
            </li>
            <li class="flex gap-3">
              <span class="font-inter text-[17px] font-medium text-text-primary min-w-[140px]">Аналитические</span>
              <span class="font-inter text-[17px] text-text-secondary">
                Помогают нам понимать, как пользователи взаимодействуют с сайтом. Можно отключить.
              </span>
            </li>
            <li class="flex gap-3">
              <span class="font-inter text-[17px] font-medium text-text-primary min-w-[140px]">Маркетинговые</span>
              <span class="font-inter text-[17px] text-text-secondary">
                Используются для персонализации рекламы. Можно отключить.
              </span>
            </li>
          </ul>
        </section>

        <section>
          <h2 class="font-inter text-[17px] font-medium text-text-primary mb-2">3. Как управлять cookies?</h2>
          <p class="font-inter text-[17px] text-text-secondary">
            Вы можете изменить настройки в любой момент через баннер внизу страницы. 
            Также большинство браузеров позволяют управлять cookies через свои настройки.
          </p>
        </section>

        <!-- Кнопка "Назад" -->
        <div class="pt-4 border-t border-card-border">
          <router-link
            to="/auth"
            class="inline-flex items-center gap-2 text-primary font-inter text-[17px] font-medium hover:underline"
          >
            <img src="@/components/icons/arrow_left.svg" alt=""/>
            Назад
          </router-link>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
// Мета-теги можно добавить при необходимости
</script>