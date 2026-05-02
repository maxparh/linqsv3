<template>
  <div class="min-h-screen bg-[#F6F9FF] flex items-center justify-center p-6">
    <div class="w-full max-w-[480px] bg-white rounded-[24px] border border-[#E2E8F0] p-6">
      <!-- Logo -->
      <div class="flex items-center justify-center gap-3 mb-8">
        <div class="w-8 h-8 rounded-lg bg-[#014751] flex items-center justify-center">
          <img src="@/components/icons/linqs_logo.svg" alt="" />
        </div>
        <span class="font-['Inter'] font-semibold text-[17px] text-[#0F172A]">Linqs</span>
      </div>

      <!-- Tabs -->
      <div class="flex gap-2 p-1 bg-[#F6F9FF] rounded-[10px] mb-6">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'flex-1 py-2.5 px-4 rounded-[10px] text-[17px] font-medium font-[\'Inter\'] transition-colors duration-200',
            activeTab === tab.id
              ? 'bg-[#014751] text-white'
              : 'text-[#475569] hover:text-[#0F172A]',
          ]"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Forms -->
      <LoginForm v-if="activeTab === 'login'" />
      <RegisterForm v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'

const tabs = [
  { id: 'login', label: 'Вход' },
  { id: 'register', label: 'Регистрация' },
]

const activeTab = ref('login')
</script>
