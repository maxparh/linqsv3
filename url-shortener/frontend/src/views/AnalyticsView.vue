<template>
  <div class="min-h-screen bg-page-bg flex">
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
          <img
            src="@/components/icons/search.svg"
            alt=""
            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none opacity-50"
          />
        </div>

        <button
          @click="showWIPPopup = true"
          class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary hover:text-text-primary transition-colors bg-white flex items-center gap-2"
        >
          Фильтры
          <img src="@/components/icons/chevron_down.svg" alt="" class="w-4 h-4" />
        </button>

        <select
          v-model="selectedDays"
          @change="loadAllData"
          class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary bg-white focus:outline-none focus:border-primary"
        >
          <option value="7">7 дней</option>
          <option value="30">30 дней</option>
          <option value="90">90 дней</option>
        </select>
      </div>

      <!-- Основной график -->
      <div class="bg-white rounded-card border border-card-border p-6 mb-6">
        <h2 class="font-inter text-[17px] font-medium text-text-secondary mb-4">
          График количества переходов
        </h2>
        <div v-if="loading" class="h-[300px] flex items-center justify-center">
          <div class="text-text-secondary">Загрузка...</div>
        </div>
        <div ref="lineChartRef" class="h-[300px] w-full"></div>
      </div>

      <!-- Сетка виджетов -->
      <div class="grid grid-cols-3 gap-6">
        <!-- 1. Переходы -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary">Переходы</span>
            </div>
            <button
              @click="removeWidget('1')"
              class="text-text-secondary hover:text-text-primary transition-colors"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
            </button>
          </div>
          <div class="font-manrope font-bold text-[28px] text-text-primary">
            {{ formatNumber(overview.total_clicks) }}
          </div>
          <div class="text-success font-inter text-[15px] mt-1">
            +{{ calculateGrowth(overview.total_clicks, prevOverview.total_clicks) }}% относительно
            прошлой недели
          </div>
          <div class="text-text-secondary font-inter text-[15px] mt-1">
            {{ formatNumber(prevOverview.total_clicks) }} на прошлой неделе
          </div>
        </div>

        <!-- 2. Уникальные переходы -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary"
                >Уникальные переходы</span
              >
            </div>
            <button
              @click="removeWidget('2')"
              class="text-text-secondary hover:text-text-primary transition-colors"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
            </button>
          </div>
          <div class="font-manrope font-bold text-[28px] text-text-primary">
            {{ formatNumber(overview.unique_clicks) }}
          </div>
          <div class="text-success font-inter text-[15px] mt-1">
            +{{ calculateGrowth(overview.unique_clicks, prevOverview.unique_clicks) }}% относительно
            прошлой недели
          </div>
          <div class="text-text-secondary font-inter text-[15px] mt-1">
            {{ formatNumber(prevOverview.unique_clicks) }} на прошлой неделе
          </div>
        </div>

        <!-- 3. Показатель отказа -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary"
                >Показатель отказа</span
              >
            </div>
            <button
              @click="removeWidget('3')"
              class="text-text-secondary hover:text-text-primary transition-colors"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
            </button>
          </div>
          <div class="font-manrope font-bold text-[28px] text-text-primary">
            {{ overview.bounce_rate.toFixed(1) }}%
          </div>
          <div class="text-success font-inter text-[15px] mt-1">
            -2% относительно прошлой недели
          </div>
          <div class="mt-4 space-y-2">
            <div class="text-text-secondary font-inter text-[15px] font-medium">
              Наиболее частые отказы
            </div>
            <div
              v-for="stat in bounceRate.slice(0, 2)"
              :key="stat.link_url"
              class="flex items-center justify-between font-inter text-[17px] text-text-primary"
            >
              <span>{{ truncateUrl(stat.link_url) }}</span>
              <span class="text-text-secondary">{{ stat.bounce_rate.toFixed(0) }}%</span>
            </div>
          </div>
        </div>

        <!-- 4. Топ локаций -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary"
                >Топ локаций</span
              >
            </div>
            <button
              @click="removeWidget('4')"
              class="text-text-secondary hover:text-text-primary transition-colors"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
            </button>
          </div>
          <div class="space-y-3 mt-2">
            <div
              v-for="loc in locations"
              :key="loc.country_code"
              class="flex items-center justify-between font-inter text-[17px]"
            >
              <div class="flex items-center gap-2">
                <img
                  :src="getFlagUrl(loc.country_code)"
                  :alt="loc.country"
                  class="w-6 h-4 object-cover rounded-sm"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <span class="text-text-primary">{{ loc.country }}</span>
              </div>
              <span class="text-text-secondary font-medium">{{ loc.percent.toFixed(1) }}%</span>
            </div>
          </div>
        </div>

        <!-- 5. Устройства -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary"
                >Наиболее распространенные устройства</span
              >
            </div>
            <button
              @click="removeWidget('5')"
              class="text-text-secondary hover:text-text-primary transition-colors flex-shrink-0"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5 flex-shrink-0" />
            </button>
          </div>
          <div class="flex items-center gap-4">
            <div ref="donutChartRef" class="h-[160px] w-[160px]"></div>
            <div class="space-y-2">
              <div
                v-for="dev in devices"
                :key="dev.name"
                class="flex items-center gap-2 font-inter text-[15px]"
              >
                <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: dev.color }"></span>
                <span class="text-text-primary">{{ dev.name }}</span>
                <span class="text-text-secondary ml-auto">{{ dev.percent.toFixed(1) }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 6. Среднее время на сайте -->
        <div class="bg-white rounded-card border border-card-border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <img
                src="@/components/icons/drag.svg"
                alt=""
                class="w-5 h-5 text-text-secondary opacity-40"
              />
              <span class="font-inter text-[17px] font-medium text-text-secondary"
                >Среднее время на сайте</span
              >
            </div>
            <button
              @click="removeWidget('6')"
              class="text-text-secondary hover:text-text-primary transition-colors"
            >
              <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
            </button>
          </div>
          <div class="font-manrope font-bold text-[28px] text-text-primary">
            {{ formatTime(overview.avg_time_on_site) }}
          </div>
          <div class="text-success font-inter text-[15px] mt-1">
            +6% относительно прошлой недели
          </div>
          <div class="text-text-secondary font-inter text-[15px] mt-1">
            {{ formatTime(overview.avg_time_on_site * 0.93) }} на прошлой неделе
          </div>
        </div>
      </div>
    </main>

    <!-- Попап "В разработке" -->
    <div
      v-if="showWIPPopup"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="showWIPPopup = false"
    >
      <div class="bg-white rounded-card border border-card-border p-8 max-w-[400px] text-center">
        <div
          class="w-16 h-16 mx-auto mb-4 rounded-full bg-page-bg flex items-center justify-center"
        >
          <img src="@/components/icons/dev.svg" alt="" />
        </div>
        <h3 class="font-manrope font-bold text-[24px] text-text-primary mb-2">Упс!</h3>
        <p class="font-inter text-[17px] text-text-secondary mb-6">
          Данный функционал в разработке
        </p>
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
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import {
  analyticsAPI,
  type AnalyticsOverview,
  type LocationStat,
  type DeviceStat,
  type BounceRateStat,
} from '@/services/analytics'
import { countryNameMap } from '@/data/countries'

const router = useRouter()
const showWIPPopup = ref(false)
const searchQuery = ref('')
const loading = ref(false)
const selectedDays = ref('7')

const lineChartRef = ref<HTMLElement | null>(null)
const donutChartRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let donutChart: echarts.ECharts | null = null

// Данные из API
const overview = ref<AnalyticsOverview>({
  total_clicks: 0,
  unique_clicks: 0,
  bounce_rate: 0,
  avg_time_on_site: 0,
})

const prevOverview = ref<AnalyticsOverview>({
  total_clicks: 0,
  unique_clicks: 0,
  bounce_rate: 0,
  avg_time_on_site: 0,
})

const locations = ref<LocationStat[]>([])
const devices = ref<DeviceStat[]>([])
const bounceRate = ref<BounceRateStat[]>([])
const clicksOverTime = ref<{ labels: string[]; values: number[] }>({ labels: [], values: [] })

// Загрузка всех данных
const loadAllData = async () => {
  loading.value = true
  const days = parseInt(selectedDays.value)

  try {
    console.log('📡 Fetching analytics for', days, 'days...')

    const [overviewData, clicksData, locationsData, devicesData, bounceData] = await Promise.all([
      analyticsAPI.getOverview(days),
      analyticsAPI.getClicksOverTime(days),
      analyticsAPI.getTopLocations(days, 5),
      analyticsAPI.getDeviceStats(days),
      analyticsAPI.getBounceRate(days),
    ])

    console.log('✅ Overview:', overviewData)
    console.log('✅ Clicks over time:', clicksData)
    console.log('✅ Locations:', locationsData)
    console.log('✅ Devices:', devicesData)
    console.log('✅ Bounce rate:', bounceData)

    overview.value = overviewData
    clicksOverTime.value = clicksData
    locations.value = locationsData || []
    devices.value = devicesData || []
    bounceRate.value = bounceData || []

    updateLineChart()
    updateDonutChart()
  } catch (error: any) {
    console.error('❌ Error loading analytics:', error)
    if (error.response) {
      console.error('Status:', error.response.status)
      console.error('Data:', error.response.data)
    }
  } finally {
    loading.value = false
  }
}
// Обновление линейного графика
const updateLineChart = () => {
  if (!lineChart) return

  const labels = clicksOverTime.value?.labels || []
  const values = clicksOverTime.value?.values || []

  lineChart.setOption({
    xAxis: { data: labels },
    series: [{ data: values }],
  })
}

// Обновление кругового графика
const updateDonutChart = () => {
  if (!donutChart) return

  const data = (devices.value || []).map((d) => ({
    value: d.percent ?? 0, // ← ИЗМЕНИЛ: теперь проценты вместо d.value
    name: d.name ?? 'Unknown',
    itemStyle: { color: d.color ?? '#94a3b8' },
  }))

  donutChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}%',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#e2e8f0',
      borderWidth: 1,
      textStyle: { color: '#0f172a', fontSize: 14 },
      padding: [8, 12],
    },
    series: [
      {
        type: 'pie',
        radius: ['45%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: false },
          scale: true,
          scaleSize: 10,
        },
        data: data, // ← Теперь здесь проценты
      },
    ],
  })
}
// Хелперы
const formatNumber = (num: number): string => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + ' млн.'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + ' тыс.'
  }
  return num.toString()
}

const formatTime = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  if (mins > 0) {
    return `${mins} мин.${secs > 0 ? ` ${secs} сек.` : ''}`
  }
  return `${secs} сек.`
}

const calculateGrowth = (current: number, prev: number): number => {
  if (prev === 0) return 0
  return Math.round(((current - prev) / prev) * 100)
}

const truncateUrl = (url: string): string => {
  if (url.length > 25) {
    return url.substring(0, 22) + '...'
  }
  return url
}

const getFlagUrl = (countryCode: string) => {
  const countryName = countryNameMap[countryCode]
  if (!countryName) {
    console.warn(`No flag mapping for country code: ${countryCode}`)
    return '' // или заглушку
  }
  return `/flags/Flag ${countryName}.png`
}

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
        axisLabel: { color: '#475569', fontSize: 14 },
      },
      yAxis: {
        type: 'value',
        splitLine: { lineStyle: { color: '#f1f5f9' } },
        axisLabel: { color: '#475569', fontSize: 14 },
      },
      series: [
        {
          data: [],
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 8,
          lineStyle: { color: '#014751', width: 3 },
          itemStyle: { color: '#014751' },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(1, 71, 81, 0.2)' },
                { offset: 1, color: 'rgba(1, 71, 81, 0)' },
              ],
            },
          },
        },
      ],
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
        padding: [8, 12],
      },
      series: [
        {
          type: 'pie',
          radius: ['45%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: false },
            scale: true,
            scaleSize: 10,
          },
          data: [],
        },
      ],
    })
  }
}

const removeWidget = (id: string) => {
  // Твоя логика удаления виджетов
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

onMounted(() => {
  initCharts()
  loadAllData()

  window.addEventListener('resize', () => {
    lineChart?.resize()
    donutChart?.resize()
  })
})

// Перезагрузка при смене периода
watch(selectedDays, () => {
  loadAllData()
})
</script>
