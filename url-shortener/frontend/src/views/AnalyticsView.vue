<template>
  <div class="analytics-page min-h-screen bg-page-bg flex">
    <aside class="page-sidebar analytics-sidebar w-[224px] bg-white border-r border-card-border flex flex-col">
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
        <SidebarProfile />
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
    <main class="analytics-main flex-1 p-8">
      <!-- Шапка -->
      <div class="analytics-header flex items-center justify-between mb-6">
        <h1 class="font-manrope font-bold text-[32px] text-text-primary">Аналитика / Все</h1>
        <button
          @click="widgetDialog?.showModal()"
          class="h-10 px-6 bg-primary text-white rounded-input font-inter text-[17px] font-medium hover:bg-[#013d41] transition-colors"
        >
          Добавить виджет
        </button>
      </div>

      <!-- Фильтры и поиск -->
      <div class="analytics-filters flex items-center gap-4 mb-8">
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
          class="h-10 px-4 border border-card-border rounded-input font-inter text-[17px] text-text-secondary bg-white focus:outline-none focus:border-primary"
        >
          <option value="7">7 дней</option>
          <option value="30">30 дней</option>
          <option value="90">90 дней</option>
        </select>
      </div>

      <AnalyticsWidget title="График количества переходов" :removable="false" class="analytics-chart mb-6">
        <div class="relative">
          <div ref="lineChartRef" class="h-[300px] w-full"></div>
          <div v-if="loading" class="absolute inset-0 bg-white/80 flex items-center justify-center text-text-secondary">Загрузка...</div>
        </div>
      </AnalyticsWidget>

      <p v-if="loadError" role="alert" class="text-error mb-4">{{ loadError }}</p>
      <VueDraggableNext v-model="selectedWidgets" class="analytics-widgets grid grid-cols-3 gap-6"
        handle=".widget-drag-handle" :animation="180" :force-fallback="true"
        :fallback-tolerance="5" ghost-class="widget-ghost" chosen-class="widget-chosen">
        <AnalyticsWidget
          v-for="id in selectedWidgets"
          :key="id"
          draggable
          :title="widgetOptions.find(widget => widget.id === id)?.title || ''"
          @close="removeWidget(id)"
          @move="moveWidget(id, $event)"
        >
          <template v-if="id === '1' || id === '2'">
          <div class="font-manrope font-bold text-[28px] text-text-primary">
            {{ loading ? '…' : formatNumber(id === '1' ? overview.total_clicks : overview.unique_clicks) }}
          </div>
          <div class="mt-3 font-inter text-[15px] space-y-2">
            <p v-if="loading" class="text-text-secondary">Загрузка сравнения...</p>
            <template v-else-if="!loadError && previousClicks(id) !== undefined">
              <p :style="{ color: currentClicks(id) > previousClicks(id)! ? '#10B981' : currentClicks(id) < previousClicks(id)! ? '#EF4444' : '#475569' }">
                {{ comparisonText(id) }} {{ selectedDays === '7' ? 'относительно прошлой недели' : `относительно предыдущих ${selectedDays} дней` }}
              </p>
              <p class="text-text-secondary">{{ formatNumber(previousClicks(id)!) }} {{ selectedDays === '7' ? 'за прошлую неделю' : `за предыдущие ${selectedDays} дней` }}</p>
            </template>
            <p v-else class="text-text-secondary">Сравнение недоступно</p>
          </div>
          </template>
          <template v-else-if="id === '4'">
            <div class="space-y-3">
              <div v-for="loc in locations" :key="loc.country_code" class="flex items-center justify-between gap-3 font-inter text-[17px]">
                <div class="flex items-center gap-2 min-w-0">
                  <img :src="getFlagUrl(loc.country_code)" :alt="loc.country" class="w-6 h-4 shrink-0 object-cover rounded-sm" @error="($event.target as HTMLImageElement).style.display = 'none'" />
                  <span class="text-text-primary">{{ loc.country }}</span>
                </div>
                <span class="text-text-secondary shrink-0">{{ loc.percent.toFixed(1) }}%</span>
              </div>
              <p v-if="!locations.length" class="text-text-secondary">{{ loading ? 'Загрузка...' : 'Нет данных' }}</p>
            </div>
            <a href="https://db-ip.com" target="_blank" rel="noopener noreferrer" class="self-start mt-auto pt-4 font-inter text-[12px] text-text-secondary underline underline-offset-2 hover:text-primary">IP Geolocation by DB-IP</a>
          </template>
          <template v-else-if="id === '5'">
            <div class="device-breakdown flex items-center gap-4">
              <div :ref="(el) => { donutChartRef = el as HTMLElement | null }" class="shrink-0 h-[160px] w-[160px]"></div>
              <div class="device-legend min-w-0 flex-1 space-y-2">
                <div v-for="dev in devices" :key="dev.name" class="flex items-center gap-2 font-inter text-[15px]">
                  <span class="w-3 h-3 shrink-0 rounded-full" :style="{ backgroundColor: dev.color }"></span>
                  <span class="text-text-primary">{{ dev.name }}</span>
                  <span class="text-text-secondary ml-auto">{{ dev.percent.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
            <p v-if="!devices.length" class="text-text-secondary">{{ loading ? 'Загрузка...' : 'Нет данных' }}</p>
          </template>
        </AnalyticsWidget>
      </VueDraggableNext>
      <dialog ref="widgetDialog" class="widget-dialog bg-white text-text-primary rounded-card border border-card-border p-6" aria-labelledby="widget-dialog-title" @click="closeWidgetBackdrop">
        <div class="flex items-center justify-between gap-4 mb-4">
          <h2 id="widget-dialog-title" class="font-manrope font-bold text-[22px]">Добавить виджет</h2>
          <button type="button" autofocus aria-label="Закрыть" title="Закрыть" class="w-10 h-10 shrink-0 flex items-center justify-center" @click="widgetDialog?.close()">
            <img src="@/components/icons/close_x.svg" alt="" class="w-5 h-5" />
          </button>
        </div>
        <div class="divide-y divide-card-border">
          <button v-for="widget in widgetOptions" :key="widget.id" type="button"
            :disabled="!widget.available || selectedWidgets.includes(widget.id)"
            class="w-full flex items-center justify-between gap-4 py-4 text-left font-inter text-[16px] hover:text-primary disabled:opacity-50 disabled:cursor-default"
            @click="addWidget(widget.id)">
            <span>{{ widget.title }}</span>
            <span v-if="!widget.available" class="text-[13px] shrink-0">Скоро</span>
            <span v-else-if="selectedWidgets.includes(widget.id)" class="text-[13px] shrink-0">Добавлен</span>
            <span v-else aria-hidden="true" class="text-[24px] shrink-0">+</span>
          </button>
        </div>
      </dialog>
    </main>

    <!-- Попап "В разработке" -->
    <div
      v-if="showWIPPopup"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto"
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
import SidebarProfile from '@/components/SidebarProfile.vue'
import AnalyticsWidget from '@/components/AnalyticsWidget.vue'
import { VueDraggableNext } from 'vue-draggable-next'
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import {
  analyticsAPI,
  type AnalyticsOverview,
  type LocationStat,
  type DeviceStat,
} from '@/services/analytics'
import { countryNameMap } from '@/data/countries'

const router = useRouter()
const showWIPPopup = ref(false)
const searchQuery = ref('')
const loading = ref(false)
const selectedDays = ref('7')
const widgetDialog = ref<HTMLDialogElement | null>(null)
const loadError = ref('')
const widgetOptions = [
  { id: '1', title: 'Переходы', available: true },
  { id: '2', title: 'Уникальные переходы', available: true },
  { id: '4', title: 'Топ локаций', available: true },
  { id: '5', title: 'Устройства', available: true },
  { id: '3', title: 'Показатель отказа', available: false },
  { id: '6', title: 'Среднее время на сайте', available: false },
]
const storageKey = 'analytics.widgets.v1'
const selectedWidgets = ref<string[]>(['1'])
try {
  const saved: unknown = JSON.parse(localStorage.getItem(storageKey) || 'null')
  if (Array.isArray(saved)) {
    selectedWidgets.value = [...new Set(saved.filter((id): id is string =>
      typeof id === 'string' && widgetOptions.some(widget => widget.id === id && widget.available)))]
  }
} catch { /* Keep the default layout when storage is unavailable or invalid. */ }
watch(selectedWidgets, (ids) => {
  try { localStorage.setItem(storageKey, JSON.stringify(ids)) } catch { /* Layout remains usable without storage. */ }
})
const addWidget = (id: string) => {
  if (selectedWidgets.value.includes(id) || !widgetOptions.some(widget => widget.id === id && widget.available)) return
  selectedWidgets.value = [...selectedWidgets.value, id]
  widgetDialog.value?.close()
}
const closeWidgetBackdrop = (event: MouseEvent) => {
  const dialog = widgetDialog.value
  if (!dialog || event.target !== dialog) return
  const bounds = dialog.getBoundingClientRect()
  if (event.clientX < bounds.left || event.clientX > bounds.right || event.clientY < bounds.top || event.clientY > bounds.bottom) dialog.close()
}

const lineChartRef = ref<HTMLElement | null>(null)
const donutChartRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let donutChart: echarts.ECharts | null = null
let chartObserver: ResizeObserver | undefined

// Данные из API
const overview = ref<AnalyticsOverview>({
  total_clicks: 0,
  unique_clicks: 0,
  bounce_rate: 0,
  avg_time_on_site: 0,
})


const locations = ref<LocationStat[]>([])
const devices = ref<DeviceStat[]>([])
const clicksOverTime = ref<{ labels: string[]; values: number[] }>({ labels: [], values: [] })

// Загрузка всех данных
const loadAllData = async () => {
  loading.value = true
  loadError.value = ''
  const days = parseInt(selectedDays.value)

  try {
    console.log('📡 Fetching analytics for', days, 'days...')

    const [overviewData, clicksData, locationsData, devicesData] = await Promise.all([
      analyticsAPI.getOverview(days),
      analyticsAPI.getClicksOverTime(days),
      analyticsAPI.getTopLocations(days, 5),
      analyticsAPI.getDeviceStats(days),
    ])

    console.log('✅ Overview:', overviewData)
    console.log('✅ Clicks over time:', clicksData)
    console.log('✅ Locations:', locationsData)
    console.log('✅ Devices:', devicesData)

    overview.value = overviewData
    clicksOverTime.value = clicksData
    locations.value = locationsData || []
    devices.value = devicesData || []

    updateLineChart()
    updateDonutChart()
  } catch (error: any) {
    loadError.value = 'Не удалось загрузить аналитику'
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
    value: d.percent ?? 0,
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
        data: data,
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

const currentClicks = (id: string) => id === '1' ? overview.value.total_clicks : overview.value.unique_clicks
const previousClicks = (id: string) => id === '1' ? overview.value.previous_total_clicks : overview.value.previous_unique_clicks
const comparisonText = (id: string): string => {
  const previous = previousClicks(id)
  if (previous === undefined) return ''
  const difference = currentClicks(id) - previous
  if (difference === 0) return 'Без изменений'
  if (previous === 0) return `+${formatNumber(difference)} переходов`
  const percent = Math.abs(difference / previous * 100)
  return `${difference > 0 ? '+' : '-'}${percent < 0.1 ? '<0,1' : percent.toLocaleString('ru-RU', { maximumFractionDigits: 1 })}%`
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
      tooltip: { trigger: 'axis', confine: true },
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
  selectedWidgets.value = selectedWidgets.value.filter(widget => widget !== id)
}

const moveWidget = (id: string, direction: number) => {
  const index = selectedWidgets.value.indexOf(id)
  const target = index + direction
  if (index < 0 || target < 0 || target >= selectedWidgets.value.length) return
  const reordered = [...selectedWidgets.value]
  reordered.splice(index, 1)
  reordered.splice(target, 0, id)
  selectedWidgets.value = reordered
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  router.push('/auth')
}

onMounted(() => {
  initCharts()
  loadAllData()

  chartObserver = new ResizeObserver(() => {
    lineChart?.resize()
    donutChart?.resize()
  })
  if (lineChartRef.value) chartObserver.observe(lineChartRef.value)
  if (donutChartRef.value) chartObserver.observe(donutChartRef.value)
})

watch(donutChartRef, (element, previous) => {
  if (previous) chartObserver?.unobserve(previous)
  donutChart?.dispose()
  donutChart = null
  if (element) {
    initCharts()
    updateDonutChart()
    chartObserver?.observe(element)
  }
}, { flush: 'post' })

onBeforeUnmount(() => {
  chartObserver?.disconnect()
  lineChart?.dispose()
  donutChart?.dispose()
})

// Перезагрузка при смене периода
watch(selectedDays, () => {
  loadAllData()
})
</script>

<style scoped>
.analytics-widgets :deep(.widget-ghost) { opacity: 0.3; }
.analytics-widgets :deep(.widget-chosen) { outline: 2px solid var(--color-primary); }
.widget-dialog {
  margin: auto;
  width: min(460px, calc(100% - 32px));
  max-height: calc(100dvh - 32px);
  overflow-y: auto;
}
.widget-dialog::backdrop {
  background: rgb(0 0 0 / 50%);
  backdrop-filter: blur(4px);
}
.widget-dialog[open] { animation: widget-appear 200ms ease-out; }
@keyframes widget-appear {
  from { opacity: 0; transform: translateY(12px) scale(.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
@media (prefers-reduced-motion: reduce) {
  .widget-dialog[open] { animation: none; }
}

.analytics-main,
.analytics-widgets > *,
.analytics-widgets .flex > div {
  min-width: 0;
}

.analytics-sidebar {
  flex-shrink: 0;
}

.analytics-header {
  gap: 16px;
}

.analytics-header > button,
.analytics-widgets button,
.analytics-widgets img,
.device-legend span:first-child,
.device-legend span:last-child {
  flex-shrink: 0;
}

.analytics-widgets {
  overflow-wrap: anywhere;
}

.analytics-widgets .justify-between {
  gap: 12px;
}

@media (width < 1280px) {
  .analytics-main {
    padding: 24px;
  }

  .analytics-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .analytics-header h1 {
    font-size: 28px;
  }

  .analytics-filters {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 12px;
  }

  .analytics-filters > div {
    grid-column: 1 / -1;
  }

  .analytics-filters > button {
    justify-content: space-between;
  }

  .analytics-widgets {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
  }

  .analytics-widgets > *,
  .analytics-chart {
    padding: 20px;
  }

  .device-breakdown {
    flex-direction: column;
  }

  .device-legend {
    width: 100%;
  }
}

@media (width < 830px) {
  .analytics-page {
    flex-direction: column;
  }

  .analytics-sidebar {
    position: fixed;
    height: auto;
    overflow: visible;
    inset: auto 0 0;
    z-index: 40;
    width: 100%;
    border-right: 0;
    border-top: 1px solid var(--color-card-border);
    padding-bottom: env(safe-area-inset-bottom, 0px);
  }

  .analytics-sidebar > div:first-child,
  .analytics-sidebar > div:last-child {
    display: none;
  }

  .analytics-sidebar nav {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 4px;
    padding: 8px;
  }

  .analytics-sidebar nav > a {
    width: 100%;
    min-width: 0;
    height: 56px;
    margin: 0;
    padding: 4px 2px;
    flex-direction: column;
    justify-content: center;
    gap: 4px;
    font-size: 12px;
    line-height: 16px;
    background: transparent;
    color: var(--color-text-secondary);
    font-weight: 500;
  }

  .analytics-sidebar nav img {
    width: 24px;
    height: 24px;
    flex-shrink: 0;
    filter: grayscale(1);
    opacity: 0.65;
  }

  .analytics-sidebar nav > a[aria-current='page'] {
    color: var(--color-primary);
    font-weight: 600;
  }

  .analytics-sidebar nav > a[aria-current='page'] img {
    filter: brightness(0) saturate(100%) invert(20%) sepia(32%) saturate(1623%) hue-rotate(143deg) brightness(92%) contrast(99%);
    opacity: 1;
  }

  .analytics-main {
    padding: 24px 16px calc(96px + env(safe-area-inset-bottom, 0px));
  }

  .analytics-header h1 {
    font-size: 24px;
  }

  .analytics-header > button {
    width: 100%;
  }

  .analytics-widgets {
    grid-template-columns: minmax(0, 1fr);
  }

  .analytics-widgets > *,
  .analytics-chart {
    padding: 16px;
  }
}
</style>
