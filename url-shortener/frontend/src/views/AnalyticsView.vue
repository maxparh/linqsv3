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
        <router-link to="/" class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors">
          <img src="@/components/icons/home_nav_nactive.svg" alt="" /> Главная
        </router-link>
        <router-link to="/links" class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors">
          <img src="@/components/icons/link_nav_nactive.svg" alt="" /> Ссылки
        </router-link>
        <router-link to="/analytics" class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] bg-primary text-white font-inter text-[17px] font-bold">
          <img src="@/components/icons/analit_nav_active.svg" alt="" /> Аналитика
        </router-link>
        <router-link
          to="/settings"
          class="flex items-center gap-3 w-[188px] h-12 px-4 rounded-[10px] text-text-secondary font-inter text-[17px] font-medium hover:text-text-primary transition-colors"
        >
          <img src="@/components/icons/settings_nav_nactive.svg" alt="" />
          Настройки
        </router-link>
      </nav>

      <div class="p-4 border-t border-card-border">
        <div class="flex items-center gap-3 px-4 py-3">
          <div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
            <span class="text-white font-inter font-semibold text-[14px]">A</span>
          </div>
          <span class="font-inter text-[17px] text-text-primary">admin</span>
        </div>
        <button @click="handleLogout" class="w-full flex items-center gap-3 px-4 py-7 text-error font-inter text-[17px] font-medium hover:bg-page-bg rounded-[10px] transition-colors">
          <img src="@/components/icons/logout_btn.svg" alt="" class="w-[24px]" /> Выход
        </button>
      </div>
    </aside>

    <!-- Основной контент -->
    <main class="flex-1 p-8">
      <!-- Шапка -->
      <div class="flex items-center justify-between mb-6">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Аналитика / Все</h1>
        <button @click="showWIPPopup = true" class="h-10 px-6 bg-primary text-white rounded-input font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors">
          Добавить виджет
        </button>
      </div>

      <!-- Фильтры и поиск -->
      <div class="flex items-center gap-4 mb-8">
        <div class="flex-1 relative">
          <input v-model="searchQuery" type="text" placeholder="Начните вводить ссылку" class="w-full h-10 pl-10 pr-4 border border-card-border rounded-input font-inter text-[17px] text-text-primary placeholder:text-placeholder focus:outline-none focus:border-primary transition-colors bg-white" />
          <img src="@/components/icons/search.svg" alt="" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none opacity-50" />
        </div>
        <button @click="showWIPPopup = true" class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white flex items-center gap-2">
          Фильтры <img src="@/components/icons/chevron_down.svg" alt="" class="w-4 h-4" />
        </button>
        <button @click="showWIPPopup = true" class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white flex items-center gap-2">
          Период <img src="@/components/icons/chevron_down.svg" alt="" class="w-4 h-4" />
        </button>
      </div>

      <!-- Загрузка -->
      <div v-if="loading" class="flex items-center justify-center h-64 text-text-secondary font-inter text-[17px]">
        Загрузка аналитики...
      </div>

      <template v-else>
        <!-- Основной график -->
        <div class="bg-white rounded-card border border-card-border p-6 mb-6">
          <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">График количества переходов</h2>
          <div ref="lineChartRef" class="h-[300px] w-full"></div>
        </div>

        <!-- Сетка виджетов -->
        <div class="grid grid-cols-3 gap-6">
          <!-- 1. Переходы -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Переходы</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors opacity-0 group-hover:opacity-100">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="font-manrope font-bold text-[28px] text-text-primary">{{ formatClicks(overview.total_clicks) }}</div>
            <div class="text-success font-inter text-[15px] mt-1">+{{ ((overview.total_clicks / Math.max(overview.unique_clicks, 1)) * 100 - 100).toFixed(1) }}% относительно прошлой недели</div>
            <div class="text-text-secondary font-inter text-[15px] mt-1">{{ formatClicks(overview.unique_clicks) }} уникальных</div>
          </div>

          <!-- 2. Уникальные переходы -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Уникальные переходы</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="font-manrope font-bold text-[28px] text-text-primary">{{ formatClicks(overview.unique_clicks) }}</div>
            <div class="text-text-secondary font-inter text-[15px] mt-1">Из {{ formatClicks(overview.total_clicks) }} всего</div>
          </div>

          <!-- 3. Показатель отказа -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Показатель отказа</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="font-manrope font-bold text-[28px] text-text-primary">{{ overview.bounce_rate.toFixed(1) }}%</div>
            <div class="mt-4 space-y-2">
              <div class="text-text-secondary font-inter text-[15px] font-medium">Наиболее частые отказы</div>
              <div v-if="bounceData.length > 0" class="flex items-center justify-between font-inter text-[17px] text-text-primary">
                <span class="truncate max-w-[150px]">{{ bounceData[0].link_url }}</span>
                <span class="text-text-secondary">{{ bounceData[0].bounce_rate.toFixed(1) }}%</span>
              </div>
              <div v-else class="text-text-secondary font-inter text-[15px]">Нет данных</div>
            </div>
          </div>

          <!-- 4. Топ локаций -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Топ локаций</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="space-y-3 mt-2">
              <div v-for="loc in locations" :key="loc.country_code" class="flex items-center justify-between font-inter text-[17px]">
                <div class="flex items-center gap-2">
                  <img :src="getFlagUrl(loc.country_code)" :alt="loc.country" class="w-6 h-4 object-cover rounded-sm" @error="($event.target as HTMLImageElement).style.display = 'none'" />
                  <span class="text-text-primary">{{ loc.country }}</span>
                </div>
                <span class="text-text-secondary font-medium">{{ loc.percent.toFixed(1) }}%</span>
              </div>
              <div v-if="locations.length === 0" class="text-text-secondary font-inter text-[15px]">Нет данных</div>
            </div>
          </div>

          <!-- 5. Устройства -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Наиболее распространенные устройства</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="flex items-center gap-4">
              <div ref="donutChartRef" class="h-[160px] w-[160px]"></div>
              <div class="space-y-2">
                <div v-for="dev in devices" :key="dev.name" class="flex items-center gap-2 font-inter text-[15px]">
                  <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: dev.color }"></span>
                  <span class="text-text-primary">{{ dev.name }}</span>
                  <span class="text-text-secondary ml-auto">{{ dev.percent.toFixed(1) }}%</span>
                </div>
                <div v-if="devices.length === 0" class="text-text-secondary font-inter text-[15px]">Нет данных</div>
              </div>
            </div>
          </div>

          <!-- 6. Среднее время на сайте -->
          <div class="bg-white rounded-card border border-card-border p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <img src="@/components/icons/drag.svg" alt="" class="w-5 h-5 text-text-secondary opacity-40" />
                <span class="font-inter text-[17px] font-medium text-text-secondary">Среднее время на сайте</span>
              </div>
              <button class="text-text-secondary hover:text-text-primary transition-colors">
                <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
              </button>
            </div>
            <div class="font-manrope font-bold text-[28px] text-text-primary">{{ formatTime(overview.avg_time_on_site) }}</div>
            <div class="text-text-secondary font-inter text-[15px] mt-1">За последние 7 дней</div>
          </div>
        </div>
      </template>
    </main>

    <!-- Попап "В разработке" -->
    <div v-if="showWIPPopup" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showWIPPopup = false">
      <div class="bg-white rounded-card border border-card-border p-8 max-w-[400px] text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-page-bg flex items-center justify-center">
          <img src="@/components/icons/dev.svg" alt="" />
        </div>
        <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-2">Упс!</h3>
        <p class="font-inter text-[17px] text-text-secondary mb-6">Данный функционал в разработке</p>
        <button @click="showWIPPopup = false" class="h-10 px-8 bg-primary text-white rounded-[10px] font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors">Понятно</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'

const router = useRouter()
const showWIPPopup = ref(false)
const searchQuery = ref('')
const loading = ref(true)

// Рефы для графиков
const lineChartRef = ref<HTMLElement | null>(null)
const donutChartRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let donutChart: echarts.ECharts | null = null

// Данные аналитики
const overview = ref({
  total_clicks: 0,
  unique_clicks: 0,
  bounce_rate: 0,
  avg_time_on_site: 0
})
const locations = ref<any[]>([])
const devices = ref<any[]>([])
const bounceData = ref<any[]>([])

// Форматирование чисел в стиль "144,4 тыс."
const formatClicks = (val: number) => {
  if (val >= 1000000) return (val / 1000000).toFixed(1).replace('.', ',') + ' млн'
  if (val >= 1000) return (val / 1000).toFixed(1).replace('.', ',') + ' тыс.'
  return val.toString()
}

// Форматирование времени
const formatTime = (seconds: number) => {
  if (seconds < 60) return `${Math.round(seconds)} сек.`
  const mins = Math.floor(seconds / 60)
  const secs = Math.round(seconds % 60)
  return `${mins} мин.`
}

// URL флага
const getFlagUrl = (code: string) => `/flags/Flag ${code}.png`

// Инициализация графиков
const initCharts = () => {
  if (lineChartRef.value && !lineChart) {
    lineChart = echarts.init(lineChartRef.value)
    lineChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: [],
        axisLine: { lineStyle: { color: '#e2e8f0' } },
        axisLabel: { color: '#475569', fontSize: 14 }
      },
      yAxis: {
        type: 'value',
        splitLine: { lineStyle: { color: '#f1f5f9' } },
        axisLabel: { color: '#475569', fontSize: 14 }
      },
      series: [{
        data: [],
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: { color: '#014751', width: 3 },
        itemStyle: { color: '#014751' },
        areaStyle: {
          color: {
            type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [{ offset: 0, color: 'rgba(1, 71, 81, 0.2)' }, { offset: 1, color: 'rgba(1, 71, 81, 0)' }]
          }
        }
      }]
    })
  }

  if (donutChartRef.value && !donutChart) {
    donutChart = echarts.init(donutChartRef.value)
    donutChart.setOption({
      tooltip: { 
        trigger: 'item',
        formatter: '{b}: {c}%',
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#e2e8f0',
        borderWidth: 1,
        textStyle: { color: '#0f172a', fontSize: 14 },
        padding: [8, 12]
      },
      series: [{
        type: 'pie',
        radius: ['45%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        emphasis: { label: { show: false }, scale: true, scaleSize: 10 },
        data: []
      }]
    })
  }
}

// Обновление данных графиков
const updateCharts = () => {
  if (lineChart) {
    // Заглушка для графика, пока бэк не вернёт реальные данные по дням
    // В продакшене подставим clicksData.labels / clicksData.values
    lineChart.setOption({
      xAxis: { data: ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'] },
      series: [{ data: [0, 0, 0, 0, 0, 0, 0] }]
    })
  }
  if (donutChart) {
    donutChart.setOption({
      series: [{ data: devices.value.map(d => ({ value: d.percent, name: d.name, itemStyle: { color: d.color } })) }]
    })
  }
}

// Загрузка аналитики
const fetchAnalytics = async () => {
  const token = localStorage.getItem('access_token')?.trim()
  if (!token) return router.push('/auth')

  try {
    const [ov, loc, dev, bounce] = await Promise.all([
      fetch('/api/analytics/overview?days=7', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()),
      fetch('/api/analytics/top-locations?days=7&limit=5', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()),
      fetch('/api/analytics/devices?days=7', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json()),
      fetch('/api/analytics/bounce-rate?days=7', { headers: { Authorization: `Bearer ${token}` } }).then(r => r.json())
    ])

    overview.value = ov || {}
    locations.value = loc || []
    devices.value = dev || []
    bounceData.value = bounce || []

    updateCharts()
  } catch (err) {
    console.error('Failed to fetch analytics:', err)
  } finally {
    loading.value = false
  }
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

onMounted(() => {
  initCharts()
  fetchAnalytics()
  
  window.addEventListener('resize', () => {
    lineChart?.resize()
    donutChart?.resize()
  })
})
</script>