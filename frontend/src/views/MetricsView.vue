<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import VueApexCharts from 'vue3-apexcharts'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const metrics = ref<any>(null)
const loading = ref(false)

interface Point { x: number; y: number }
const series = ref<Point[]>([])
const lastTotal = ref<number | null>(null)
const lastTime = ref<number | null>(null)
const MAX_POINTS = 60
const POLL_MS = 10_000

async function fetchMetrics() {
  loading.value = true
  try {
    const data: any = await apiFetch('/api/metrics')
    metrics.value = data
    const total = Number(data?.total_requests ?? 0)
    const now = Date.now()
    if (lastTotal.value !== null && lastTime.value !== null) {
      const elapsedSec = Math.max(1, (now - lastTime.value) / 1000)
      const delta = Math.max(0, total - lastTotal.value)
      const rate = delta / elapsedSec
      series.value = [...series.value, { x: now, y: rate }].slice(-MAX_POINTS)
    }
    lastTotal.value = total
    lastTime.value = now
  } catch {}
  finally { loading.value = false }
}

fetchMetrics()
const interval = setInterval(fetchMetrics, POLL_MS)
onUnmounted(() => clearInterval(interval))

const chartSeries = computed(() => [{ name: 'Requests / sec', data: series.value }])

const chartOptions = computed(() => ({
  chart: {
    type: 'area',
    height: 360,
    fontFamily: 'Inter, ui-sans-serif, system-ui, sans-serif',
    background: 'transparent',
    foreColor: '#5a5548',
    toolbar: { show: false },
    zoom: { enabled: false },
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 400,
      animateGradually: { enabled: false },
      dynamicAnimation: { enabled: true, speed: 350 },
    },
    sparkline: { enabled: false },
  },
  colors: ['#5266eb'],
  stroke: {
    curve: 'smooth',
    width: 2,
    lineCap: 'round',
  },
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.18,
      opacityTo: 0,
      stops: [0, 100],
      colorStops: [
        { offset: 0, color: '#5266eb', opacity: 0.18 },
        { offset: 100, color: '#5266eb', opacity: 0 },
      ],
    },
  },
  dataLabels: { enabled: false },
  markers: {
    size: 0,
    colors: ['#5266eb'],
    strokeColors: '#ffffff',
    strokeWidth: 2,
    hover: { size: 5 },
  },
  grid: {
    borderColor: '#ded9ca',
    strokeDashArray: 0,
    xaxis: { lines: { show: false } },
    yaxis: { lines: { show: true } },
    padding: { left: 10, right: 10, top: 0, bottom: 0 },
  },
  xaxis: {
    type: 'datetime',
    labels: {
      style: { colors: '#8a8478', fontFamily: 'Inter', fontSize: '11px' },
      datetimeUTC: false,
      format: 'HH:mm:ss',
    },
    axisBorder: { color: '#ded9ca' },
    axisTicks: { color: '#ded9ca' },
    crosshairs: {
      stroke: { color: '#c9c3b3', width: 1, dashArray: 0 },
    },
  },
  yaxis: {
    min: 0,
    forceNiceScale: true,
    labels: {
      style: { colors: '#8a8478', fontFamily: 'Inter', fontSize: '11px' },
      formatter: (v: number) => {
        if (v >= 1000) return (v / 1000).toFixed(1) + 'k'
        return v.toFixed(2)
      },
    },
  },
  tooltip: {
    enabled: true,
    theme: 'light',
    x: { format: 'HH:mm:ss' },
    y: {
      formatter: (v: number) => v.toFixed(2) + ' req/s',
    },
    style: { fontFamily: 'Inter', fontSize: '12px' },
    marker: { show: false },
  },
  legend: { show: false },
}))

const totalRequests = computed(() => Number(metrics.value?.total_requests ?? 0))

const currentRate = computed(() => {
  const last = series.value[series.value.length - 1]
  return last ? last.y : 0
})

const peakRate = computed(() => {
  if (!series.value.length) return 0
  return Math.max(...series.value.map((p) => p.y))
})

function formatNumber(n: number) {
  if (n === undefined || n === null || Number.isNaN(n)) return '—'
  return n.toLocaleString('en-US')
}
function formatRate(n: number) {
  if (n >= 100) return n.toFixed(0)
  if (n >= 10) return n.toFixed(1)
  return n.toFixed(2)
}
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Metrics</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Cloudflared traffic</h1>
      </div>
      <button class="btn-secondary" :disabled="loading" @click="fetchMetrics">
        <RefreshCw class="w-4 h-4" :stroke-width="1.75" :class="loading && 'animate-spin'" />
        Refresh
      </button>
    </header>

    <section class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="account-tile">
        <span class="account-tile-label">Total requests</span>
        <span class="account-tile-value">{{ formatNumber(totalRequests) }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Current rate</span>
        <span class="account-tile-value">{{ formatRate(currentRate) }}</span>
        <span class="account-tile-sub">requests / second</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Peak rate</span>
        <span class="account-tile-value">{{ formatRate(peakRate) }}</span>
        <span class="account-tile-sub">since session start</span>
      </div>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Request rate</span>
        <span class="text-2xs text-text-dim tabular-nums">
          window {{ series.length }}/{{ MAX_POINTS }} · sample every {{ POLL_MS / 1000 }}s
        </span>
      </div>
      <div class="px-2 py-2">
        <VueApexCharts
          v-if="series.length"
          type="area"
          height="360"
          :options="chartOptions"
          :series="chartSeries"
        />
        <div v-else class="text-center text-sm text-text-muted py-24">
          Collecting samples… first datapoint after {{ POLL_MS / 1000 }}s.
        </div>
      </div>
    </section>
  </div>
</template>
