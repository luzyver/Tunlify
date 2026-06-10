<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useApi, reportError } from '../composables/useApi'


const { apiFetch } = useApi()
const metrics = ref<any>(null)
const loading = ref(false)

async function fetchMetrics() {
  loading.value = true
  try { metrics.value = await apiFetch('/api/metrics') } catch (e) { reportError('failed to load metrics', e) }
  finally { loading.value = false }
}
fetchMetrics()
const interval = setInterval(fetchMetrics, 10_000)
onUnmounted(() => clearInterval(interval))

interface CodeRow { code: string; count: number; pct: number; scale: number; color: string; label: string }

const SIGNAL: Record<string, { color: string; label: string }> = {
  '2': { color: '#10B981', label: 'Success' },
  '3': { color: '#64748B', label: 'Redirect' },
  '4': { color: '#F7931A', label: 'Client error' },
  '5': { color: '#EF4444', label: 'Server error' },
}

function meta(code: string) { return SIGNAL[code[0]] || { color: '#475569', label: 'Other' } }

const rows = computed<CodeRow[]>(() => {
  const codes = metrics.value?.response_codes || {}
  const entries = Object.entries(codes).map(([code, count]) => ({ code, count: Number(count) }))
  if (!entries.length) return []
  const max = Math.max(...entries.map((e) => e.count))
  const total = entries.reduce((s, e) => s + e.count, 0)
  return entries.sort((a, b) => b.count - a.count).map((e) => {
    const m = meta(e.code)
    return { code: e.code, count: e.count, pct: total ? (e.count / total) * 100 : 0, scale: max ? (e.count / max) * 100 : 0, color: m.color, label: m.label }
  })
})

const totalRequests = computed(() => Number(metrics.value?.total_requests ?? 0))
const totalCounted = computed(() => rows.value.reduce((s, r) => s + r.count, 0))

function formatNumber(n: number) {
  if (n === undefined || n === null || Number.isNaN(n)) return '\u2014'
  return n.toLocaleString('en-US')
}
function formatPct(p: number) {
  if (p >= 10) return p.toFixed(1) + '%'
  return p.toFixed(2) + '%'
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <span class="page-badge">Console &middot; Metrics</span>
        <h1 class="page-title">Cloudflared traffic</h1>
      </div>
      <button class="btn btn-secondary" :disabled="loading" @click="fetchMetrics">
        <v-icon size="14" :class="{ 'spin': loading }">mdi-refresh</v-icon>
        Refresh
      </button>
    </div>

    <div class="mb-4">
      <div class="stat-grid" style="grid-template-columns: 1fr 1fr;">
        <div class="card stat-card card-hover">
          <p class="stat-label">Total requests</p>
          <p class="stat-value tabular-nums">{{ formatNumber(totalRequests) }}</p>
        </div>
        <div class="card stat-card card-hover">
          <p class="stat-label">Distinct codes</p>
          <p class="stat-value tabular-nums">{{ rows.length }}</p>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <span class="card-title">Response codes</span>
        <span class="stat-label" style="margin-bottom: 0; color: #71717A;">{{ formatNumber(totalCounted) }} counted &middot; auto-refresh 10s</span>
      </div>

      <div v-if="rows.length" class="pa-4">
        <div
          v-for="row in rows"
          :key="row.code"
          class="d-grid align-center ga-4 px-3 rounded-lg metric-row"
          style="grid-template-columns: 56px 130px 1fr 96px 76px; height: 36px;"
        >
          <span class="font-mono font-weight-medium" :style="{ color: row.color, fontSize: '14px' }">{{ row.code }}</span>
          <span class="font-mono d-none d-md-block text-truncate" style="color: #94A3B8; font-size: 12px;">{{ row.label }}</span>
          <div class="metric-bar">
            <div :style="{ width: row.scale + '%', background: row.color }"></div>
          </div>
          <span class="tabular-nums text-right font-mono" style="color: white; font-size: 14px;">{{ formatNumber(row.count) }}</span>
          <span class="tabular-nums text-right font-mono" style="color: #94A3B8; font-size: 12px;">{{ formatPct(row.pct) }}</span>
        </div>
      </div>

      <div v-else class="empty">
        {{ metrics ? 'No response codes recorded yet' : 'Loading metrics...' }}
      </div>

      <div v-if="rows.length" class="d-flex flex-wrap ga-4 px-6 py-3" style="border-top: 1px solid rgba(30, 41, 59, 0.6);">
        <div v-for="cls in [
          { code: '2xx', label: 'Success', color: '#FFD600' },
          { code: '3xx', label: 'Redirect', color: '#94A3B8' },
          { code: '4xx', label: 'Client error', color: '#F7931A' },
          { code: '5xx', label: 'Server error', color: '#EF4444' },
        ]" :key="cls.code" class="d-flex align-center ga-2">
          <span class="rounded-full d-inline-block" :style="{ width: '10px', height: '10px', background: cls.color }"></span>
          <span class="font-mono" style="color: #94A3B8; font-size: 11px;">{{ cls.code }}</span>
          <span class="font-mono" style="color: #94A3B8; font-size: 11px;">&middot; {{ cls.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.spin { animation: pg-spin 1s linear infinite; }
.metric-row { transition: background 0.2s; border-radius: 8px; }
.metric-row:hover { background: rgba(247, 147, 26, 0.03); }
.metric-bar { height: 8px; border-radius: 999px; overflow: hidden; background: rgba(255, 255, 255, 0.03); }
.metric-bar > div { height: 100%; border-radius: 999px; transition: width 0.5s ease-out; }
</style>
