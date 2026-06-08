<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
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

interface CodeRow { code: string; count: number; pct: number; scale: number; color: string; label: string }

const SIGNAL: Record<string, { color: string; label: string }> = {
  '2': { color: '#FFD600', label: 'Success' },
  '3': { color: '#94A3B8', label: 'Redirect' },
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
  <div class="pa-6">
    <div class="d-flex align-center justify-space-between ga-4 mb-6 flex-wrap" style="gap: 16px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Metrics</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Cloudflared traffic</h1>
      </div>
      <button
        class="d-inline-flex align-center ga-2 rounded-pill px-4 font-mono"
        style="height: 36px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 12px; cursor: pointer; transition: all 0.3s;"
        :disabled="loading"
        @click="fetchMetrics"
        @mouseenter="if(!loading) { $event.target.style.background = 'rgba(247, 147, 26, 0.2)'; $event.target.style.borderColor = '#F7931A' }"
        @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'; $event.target.style.borderColor = 'rgba(247, 147, 26, 0.2)'"
      >
        <v-icon size="14" :class="{ 'spin': loading }">mdi-refresh</v-icon>
        Refresh
      </button>
    </div>

    <div class="mb-6">
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
        <div class="rounded-2xl p-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
          <p class="eyebrow mb-2">Total requests</p>
          <p class="font-heading font-semibold tabular-nums" style="color: white; font-size: 24px;">{{ formatNumber(totalRequests) }}</p>
        </div>
        <div class="rounded-2xl p-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
          <p class="eyebrow mb-2">Distinct codes</p>
          <p class="font-heading font-semibold tabular-nums" style="color: white; font-size: 24px;">{{ rows.length }}</p>
        </div>
      </div>
    </div>

    <div class="rounded-2xl" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Response codes</span>
        <span class="font-mono tabular-nums" style="color: #94A3B8; font-size: 11px;">{{ formatNumber(totalCounted) }} counted &middot; auto-refresh 10s</span>
      </div>

      <div v-if="rows.length" class="pa-4">
        <div v-for="row in rows" :key="row.code" class="d-grid align-center ga-4 px-3 rounded-lg" style="grid-template-columns: 56px 130px 1fr 96px 76px; height: 36px; transition: background 0.2s;" @mouseenter="$event.currentTarget.style.background = 'rgba(247, 147, 26, 0.03)'" @mouseleave="$event.currentTarget.style.background = 'transparent'">
          <span class="font-mono font-weight-medium" :style="{ color: row.color, fontSize: '14px' }">{{ row.code }}</span>
          <span class="font-mono d-none d-md-block text-truncate" style="color: #94A3B8; font-size: 12px;">{{ row.label }}</span>
          <div style="height: 8px; border-radius: 9999px; overflow: hidden; background: rgba(30, 41, 59, 0.6);">
            <div :style="{ width: row.scale + '%', height: '100%', background: row.color, borderRadius: '9999px', transition: 'width 0.5s ease-out' }"></div>
          </div>
          <span class="tabular-nums text-right font-mono" style="color: white; font-size: 14px;">{{ formatNumber(row.count) }}</span>
          <span class="tabular-nums text-right font-mono" style="color: #94A3B8; font-size: 12px;">{{ formatPct(row.pct) }}</span>
        </div>
      </div>

      <div v-else class="pa-8 text-center font-mono" style="color: #94A3B8; font-size: 13px;">
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
@keyframes spin { to { transform: rotate(360deg); } }
.spin { animation: spin 1s linear infinite; }
</style>
