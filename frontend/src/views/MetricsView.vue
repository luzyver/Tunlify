<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const metrics = ref<any>(null)
const loading = ref(false)

async function fetchMetrics() {
  loading.value = true
  try { metrics.value = await apiFetch('/api/metrics') } catch {}
  finally { loading.value = false }
}
fetchMetrics()
const interval = setInterval(fetchMetrics, 10_000)
onUnmounted(() => clearInterval(interval))

interface CodeRow {
  code: string
  count: number
  pct: number
  scale: number
  bar: string
  klass: string
  label: string
}

const SIGNAL: Record<string, { bar: string; klass: string; label: string }> = {
  '2': { bar: 'bg-success', klass: 'text-success', label: 'Success' },
  '3': { bar: 'bg-text-muted', klass: 'text-text-muted', label: 'Redirect' },
  '4': { bar: 'bg-warning', klass: 'text-warning', label: 'Client error' },
  '5': { bar: 'bg-danger', klass: 'text-danger', label: 'Server error' },
}

function meta(code: string) {
  return SIGNAL[code[0]] || { bar: 'bg-text-dim', klass: 'text-text-dim', label: 'Other' }
}

const rows = computed<CodeRow[]>(() => {
  const codes = metrics.value?.response_codes || {}
  const entries = Object.entries(codes).map(([code, count]) => ({ code, count: Number(count) }))
  if (!entries.length) return []
  const max = Math.max(...entries.map((e) => e.count))
  const total = entries.reduce((s, e) => s + e.count, 0)
  return entries
    .sort((a, b) => b.count - a.count)
    .map((e) => {
      const m = meta(e.code)
      return {
        code: e.code,
        count: e.count,
        pct: total ? (e.count / total) * 100 : 0,
        scale: max ? (e.count / max) * 100 : 0,
        bar: m.bar,
        klass: m.klass,
        label: m.label,
      }
    })
})

const totalRequests = computed(() => Number(metrics.value?.total_requests ?? 0))
const totalCounted = computed(() => rows.value.reduce((s, r) => s + r.count, 0))

function formatNumber(n: number) {
  if (n === undefined || n === null || Number.isNaN(n)) return '—'
  return n.toLocaleString('en-US')
}

function formatPct(p: number) {
  if (p >= 10) return p.toFixed(1) + '%'
  return p.toFixed(2) + '%'
}

const legend = [
  { code: '2xx', label: 'Success', dot: 'bg-success' },
  { code: '3xx', label: 'Redirect', dot: 'bg-text-muted' },
  { code: '4xx', label: 'Client error', dot: 'bg-warning' },
  { code: '5xx', label: 'Server error', dot: 'bg-danger' },
]
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

    <section class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div class="account-tile">
        <span class="account-tile-label">Total requests</span>
        <span class="account-tile-value">{{ formatNumber(totalRequests) }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Distinct codes</span>
        <span class="account-tile-value">{{ rows.length }}</span>
      </div>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Response codes</span>
        <span class="text-2xs text-text-dim tabular-nums">
          {{ formatNumber(totalCounted) }} counted · auto-refresh 10s
        </span>
      </div>

      <div v-if="rows.length" class="p-4 space-y-1">
        <div
          v-for="row in rows"
          :key="row.code"
          class="grid items-center gap-4 px-2 h-9 rounded-md hover:bg-surface-alt transition-colors"
          style="grid-template-columns: 56px 130px 1fr 96px 76px;"
        >
          <span class="font-mono text-sm font-medium" :class="row.klass">{{ row.code }}</span>
          <span class="text-xs text-text-muted hidden md:block truncate">{{ row.label }}</span>
          <div class="relative h-2 bg-bg-alt rounded-full overflow-hidden">
            <div
              class="absolute inset-y-0 left-0 rounded-full transition-[width] duration-500 ease-out"
              :class="row.bar"
              :style="{ width: row.scale + '%' }"
            ></div>
          </div>
          <span class="num text-sm text-text text-right">{{ formatNumber(row.count) }}</span>
          <span class="num text-xs text-text-dim text-right">{{ formatPct(row.pct) }}</span>
        </div>
      </div>

      <div v-else class="p-8 text-center text-sm text-text-muted">
        {{ metrics ? 'No response codes recorded yet' : 'Loading metrics…' }}
      </div>

      <div
        v-if="rows.length"
        class="px-4 py-3 flex flex-wrap items-center gap-x-5 gap-y-2 border-t border-border"
      >
        <div v-for="cls in legend" :key="cls.code" class="flex items-center gap-2">
          <span class="w-2.5 h-2.5 rounded-sm" :class="cls.dot"></span>
          <span class="text-2xs text-text-muted font-mono">{{ cls.code }}</span>
          <span class="text-2xs text-text-muted">· {{ cls.label }}</span>
        </div>
      </div>
    </section>
  </div>
</template>
