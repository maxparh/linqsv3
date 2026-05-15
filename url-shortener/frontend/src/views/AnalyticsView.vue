<template>
  <div class="min-h-screen bg-page-bg flex">
    <!-- Боковое меню -->
    <aside class="w-[224px] bg-white border-r border-card-border flex flex-col">
      <div class="p-6 flex items-center gap-3 pb-[80px]">
        <div class="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
          <img src="@/components/icons/linqs_logo.svg" alt="" />
        </div>
        <span class="font-inter font-semibold text-[17px] text-text-primary">Linqs</span>
      </div>

      <nav class="flex-1 px-4 space-y-1">
        <router-link
          to="/"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/home_nav_nactive.svg" alt="" />
          Главная
        </router-link>

        <router-link
          to="/links"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/link_nav_nactive.svg" alt="" />
          Ссылки
        </router-link>

        <router-link
          to="/analytics"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] bg-primary text-white font-inter text-[17px] font-bold"
        >
          <img src="@/components/icons/analit_nav_active.svg" alt="" />
          Аналитика
        </router-link>

        <button
          @click="showWIPPopup = true"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/settings_nav_nactive.svg" alt="" />
          Настройки
        </button>
      </nav>

      <div class="p-4 border-t border-card-border">
        <div class="flex items-center gap-3 px-4 py-3">
          <div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
            <span class="text-white font-inter font-semibold text-[14px]">A</span>
          </div>
          <span class="font-inter text-[17px] text-text-primary">admin</span>
        </div>
        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 px-4 py-7 text-error font-inter text-[17px] font-medium hover:bg-page-bg rounded-[10px] transition-colors"
        >
          <img src="@/components/icons/logout_btn.svg" alt="" class="w-[24px]" />
          Выход
        </button>
      </div>
    </aside>

    <!-- Основной контент -->
    <main class="flex-1 p-8">
      <!-- Шапка -->
      <div class="flex items-center justify-between mb-6">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Аналитика / Все</h1>
        <button
          @click="showWIPPopup = true"
          class="h-10 px-6 bg-primary text-white rounded-input font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
        >
          Добавить виджет
        </button>
      </div>

      <!-- Фильтры и поиск -->
      <div class="flex items-center gap-4 mb-8">
        <div class="flex-1 relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Начните вводить ссылку"
            class="w-full h-10 pl-10 pr-4 border border-card-border rounded-input font-inter text-[17px] text-text-primary placeholder:text-placeholder focus:outline-none focus:border-primary transition-colors bg-white"
          />
          <img src="@/components/icons/search.svg" alt="" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none opacity-50" />
        </div>

        <button @click="showWIPPopup = true" class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white flex items-center gap-2">
          Фильтры
          <img src="@/components/icons/chevron_down.svg" alt="" class="w-4 h-4" />
        </button>

        <button @click="showWIPPopup = true" class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white flex items-center gap-2">
          Период
          <img src="@/components/icons/chevron_down.svg" alt="" class="w-4 h-4" />
        </button>
      </div>

      <!-- Основной график -->
      <div class="bg-white rounded-card border border-card-border p-6 mb-6">
        <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">График количества переходов</h2>
        <div ref="lineChartRef" class="h-[300px] w-full"></div>
      </div>

      <!-- Сетка виджетов с drag-and-drop -->
      <draggable 
        v-model="widgets" 
        item-key="id" 
        class="grid grid-cols-3 gap-6"
        animation="200"
        ghost-class="opacity-50"
        @end="saveWidgetOrder"
      >
        <template #item="{ element }">
          <div class="bg-white rounded-card border border-card-border p-6 flex flex-col h-full relative group">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2 cursor-move">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">{{ element.title }}</span>
              </div>
              <button @click="removeWidget(element.id)" class="text-text-secondary hover:text-text-primary transition-colors opacity-0 group-hover:opacity-100">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            
            <div class="flex-1">
              <!-- 1. Переходы -->
              <div v-if="element.type === 'clicks'">
                <div class="font-manrope font-bold text-[28px] text-text-primary">144,4 тыс.</div>
                <div class="text-success font-inter text-[15px] mt-1">+4% относительно прошлой недели</div>
                <div class="text-text-secondary font-inter text-[15px] mt-1">135,3 тыс. на прошлой неделе</div>
              </div>

              <!-- 2. Уникальные -->
              <div v-else-if="element.type === 'unique'">
                <div class="font-manrope font-bold text-[28px] text-text-primary">65,4 тыс.</div>
                <div class="text-success font-inter text-[15px] mt-1">+4% относительно прошлой недели</div>
                <div class="text-text-secondary font-inter text-[15px] mt-1">62,9 тыс. на прошлой неделе</div>
              </div>

              <!-- 3. Отказы -->
              <div v-else-if="element.type === 'bounce'">
                <div class="font-manrope font-bold text-[28px] text-text-primary">6%</div>
                <div class="text-success font-inter text-[15px] mt-1">-2% относительно прошлой недели</div>
                <div class="mt-4 space-y-2">
                  <div class="text-text-secondary font-inter text-[15px] font-medium">Наиболее частые отказы</div>
                  <div class="flex items-center justify-between font-inter text-[17px] text-text-primary">
                    <span>linqs.ru/w4db7</span>
                    <span class="text-text-secondary">10%</span>
                  </div>
                  <div class="flex items-center justify-between font-inter text-[17px] text-text-primary">
                    <span>yandex.ru</span>
                    <span class="text-text-secondary">8%</span>
                  </div>
                </div>
              </div>

              <!-- 4. Локации -->
              <div v-else-if="element.type === 'locations'">
                <div class="space-y-3 mt-2">
                  <div v-for="loc in locations" :key="loc.country" class="flex items-center justify-between font-inter text-[17px]">
                    <div class="flex items-center gap-2">
                      <img 
                        :src="getFlagUrl(loc.code)" 
                        :alt="loc.country"
                        class="w-6 h-4 object-cover rounded-sm"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                      />
                      <span class="text-text-primary">{{ loc.country }}</span>
                    </div>
                    <span class="text-text-secondary font-medium">{{ loc.percent }}%</span>
                  </div>
                </div>
              </div>

              <!-- 5. Устройства -->
              <div v-else-if="element.type === 'devices'">
                <div class="flex items-center gap-4">
                  <div ref="donutChartRef" class="h-[160px] w-[160px]"></div>
                  <div class="space-y-2">
                    <div v-for="dev in devices" :key="dev.name" class="flex items-center gap-2 font-inter text-[15px]">
                      <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: dev.color }"></span>
                      <span class="text-text-primary">{{ dev.name }}</span>
                      <span class="text-text-secondary ml-auto">{{ dev.value }}%</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 6. Время -->
              <div v-else-if="element.type === 'time'">
                <div class="font-manrope font-bold text-[28px] text-text-primary">15 мин.</div>
                <div class="text-success font-inter text-[15px] mt-1">+6% относительно прошлой недели</div>
                <div class="text-text-secondary font-inter text-[15px] mt-1">14 мин. на прошлой неделе</div>
              </div>
            </div>
          </div>
        </template>
      </draggable>
    </main>

    <!-- Попап "В разработке" -->
    <div
      v-if="showWIPPopup"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="showWIPPopup = false"
    >
      <div class="bg-white rounded-card border border-card-border p-8 max-w-[400px] text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-page-bg flex items-center justify-center">
          <img src="@/components/icons/dev.svg" alt="" />
        </div>
        <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-2">Упс!</h3>
        <p class="font-inter text-[17px] text-text-secondary mb-6">Данный функционал в разработке</p>
        <button
          @click="showWIPPopup = false"
          class="h-10 px-8 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
        >
          Понятно
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import draggable from 'vuedraggable'

const router = useRouter()
const showWIPPopup = ref(false)
const searchQuery = ref('')

// 🔥 Рефы для графиков
const lineChartRef = ref<HTMLElement | null>(null)
const donutChartRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let donutChart: echarts.ECharts | null = null

// 🔥 Хелперы для куки
const setCookie = (name: string, value: string, days = 30) => {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`
}

const getCookie = (name: string): string | null => {
  const matches = document.cookie.match(new RegExp(`(?:^|; )${name.replace(/([.$?*|{}()\[\]\\/+^])/g, '\\$1')}=([^;]*)`))
  return matches && matches[1] ? decodeURIComponent(matches[1]) : null
}

// 🔥 Тип виджета
interface Widget {
  id: string
  type: 'clicks' | 'unique' | 'bounce' | 'locations' | 'devices' | 'time'
  title: string
}

// 🔥 Тип локации
interface Location {
  country: string
  code: string
  percent: number
}

// 🔥 Тип устройства
interface Device {
  name: string
  value: number
  color: string
}

// 🔥 Дефолтные виджеты
const defaultWidgets: Widget[] = [
  { id: '1', type: 'clicks', title: 'Переходы' },
  { id: '2', type: 'unique', title: 'Уникальные переходы' },
  { id: '3', type: 'bounce', title: 'Показатель отказа' },
  { id: '4', type: 'locations', title: 'Топ локаций' },
  { id: '5', type: 'devices', title: 'Наиболее распространенные устройства' },
  { id: '6', type: 'time', title: 'Среднее время на сайте' }
]

const widgets = ref<Widget[]>([])

// 🔥 Данные для локаций
const locations = ref<Location[]>([
  { country: 'Россия', code: 'Russia', percent: 50 },
  { country: 'Беларусь', code: 'Belarus', percent: 20 },
  { country: 'Германия', code: 'Germany', percent: 11 },
  { country: 'Финляндия', code: 'Finland', percent: 10 },
  { country: 'Франция', code: 'France', percent: 6 }
])

// 🔥 Данные для устройств
const devices = ref<Device[]>([
  { name: 'Компьютеры', value: 42, color: '#014751' },
  { name: 'Android', value: 28, color: '#14b8a6' },
  { name: 'iOS', value: 20, color: '#10b981' },
  { name: 'Другие', value: 10, color: '#3b82f6' }
])

// 🔥 Получить URL флага
const getFlagUrl = (countryCode: string) => {
  return `/flags/Flag ${countryCode}.png`
}

// 🔥 Инициализация графиков
const initCharts = () => {
  // Линейный график
  if (lineChartRef.value && !lineChart) {
    lineChart = echarts.init(lineChartRef.value)
    lineChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'],
        axisLine: { lineStyle: { color: '#e2e8f0' } },
        axisLabel: { color: '#475569', fontSize: 14 }
      },
      yAxis: {
        type: 'value',
        splitLine: { lineStyle: { color: '#f1f5f9' } },
        axisLabel: { color: '#475569', fontSize: 14 }
      },
      series: [{
        data: [25000, 60000, 20000, 65000, 110000, 85000, 130000],
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: { color: '#014751', width: 3 },
        itemStyle: { color: '#014751' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(1, 71, 81, 0.2)' },
              { offset: 1, color: 'rgba(1, 71, 81, 0)' }
            ]
          }
        }
      }]
    })
  }

  // Кольцевой график
  if (donutChartRef.value && !donutChart) {
    donutChart = echarts.init(donutChartRef.value)
    donutChart.setOption({
      tooltip: { trigger: 'item' },
      series: [{
        type: 'pie',
        radius: ['45%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        emphasis: { label: { show: true, fontSize: 16, fontWeight: 'bold' } },
        data: [
          { value: 42, name: 'Компьютеры', itemStyle: { color: '#014751' } },
          { value: 28, name: 'Android', itemStyle: { color: '#14b8a6' } },
          { value: 20, name: 'iOS', itemStyle: { color: '#10b981' } },
          { value: 10, name: 'Другие', itemStyle: { color: '#3b82f6' } }
        ]
      }]
    })
  }
}

// 🔥 Сохранение порядка виджетов
const saveWidgetOrder = () => {
  const ids = widgets.value.map((w: Widget) => w.id)
  setCookie('analytics_widgets_order', JSON.stringify(ids))
}

// 🔥 Удаление виджета
const removeWidget = (id: string) => {
  widgets.value = widgets.value.filter((w: Widget) => w.id !== id)
  saveWidgetOrder()
}

// 🔥 Загрузка порядка виджетов
const loadWidgetOrder = () => {
  const saved = getCookie('analytics_widgets_order')
  if (saved) {
    const savedIds: string[] = JSON.parse(saved)
    widgets.value = savedIds
      .map((id: string) => defaultWidgets.find((w: Widget) => w.id === id))
      .filter((w): w is Widget => w !== undefined)
    
    // Добавляем новые виджеты, если появились
    const existingIds = new Set(savedIds)
    defaultWidgets.forEach((w: Widget) => {
      if (!existingIds.has(w.id)) {
        widgets.value.push(w)
      }
    })
  } else {
    widgets.value = [...defaultWidgets]
  }
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

onMounted(() => {
  loadWidgetOrder()
  initCharts()
  
  // Ресайз графиков при изменении окна
  window.addEventListener('resize', () => {
    lineChart?.resize()
    donutChart?.resize()
  })
})
</script>