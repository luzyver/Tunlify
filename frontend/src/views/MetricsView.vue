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

interface CodeRow {
  code: string
  count: number
  pct: number
  scale: number
  bar: string
  label: string
}

const SIGNAL: Record<string, { bar: string; label: string }> = {
  '2': { bar: 'success', label: 'Success' },
  '3': { bar: 'secondary', label: 'Redirect' },
  '4': { bar: 'warning', label: 'Client error' },
  '5': { bar: 'error', label: 'Server error' },
}

function meta(code: string) {
  return SIGNAL[code[0]] || { bar: 'grey', label: 'Other' }
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
        label: m.label,
      }
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
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Metrics</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Cloudflared traffic</h1>
      </div>
      <v-btn variant="tonal" :loading="loading" @click="fetchMetrics">
        <v-icon start>mdi-refresh</v-icon>
        Refresh
      </v-btn>
    </div>

    <v-row class="mb-6">
      <v-col cols="12" sm="6">
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Total requests</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ formatNumber(totalRequests) }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6">
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Distinct codes</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ rows.length }}</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-card>
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Response codes</span>
        <span class="text-caption text-medium-emphasis tabular-nums">
          {{ formatNumber(totalCounted) }} counted &middot; auto-refresh 10s
        </span>
      </v-card-title>

      <div v-if="rows.length" class="pa-4">
        <div
          v-for="row in rows"
          :key="row.code"
          class="d-grid align-center ga-4 px-2 rounded hover-bg-grey-lighten-4"
          style="grid-template-columns: 56px 130px 1fr 96px 76px; height: 36px;"
        >
          <span class="font-mono text-body-2 font-weight-medium" :class="`text-${row.bar}`">{{ row.code }}</span>
          <span class="text-caption text-medium-emphasis d-none d-md-block text-truncate">{{ row.label }}</span>
          <v-progress-linear
            :model-value="row.scale"
            :color="row.bar"
            height="8"
            rounded
            class="flex-grow-1"
          />
          <span class="tabular-nums text-body-2 text-on-surface text-right">{{ formatNumber(row.count) }}</span>
          <span class="tabular-nums text-caption text-medium-emphasis text-right">{{ formatPct(row.pct) }}</span>
        </div>
      </div>

      <div v-else class="pa-8 text-center text-body-2 text-medium-emphasis">
        {{ metrics ? 'No response codes recorded yet' : 'Loading metrics...' }}
      </div>

      <div v-if="rows.length" class="pa-4 d-flex flex-wrap ga-4 border-top" style="border-top: 1px solid #ded9ca;">
        <div v-for="cls in [
          { code: '2xx', label: 'Success', color: 'success' },
          { code: '3xx', label: 'Redirect', color: 'secondary' },
          { code: '4xx', label: 'Client error', color: 'warning' },
          { code: '5xx', label: 'Server error', color: 'error' },
        ]" :key="cls.code" class="d-flex align-center ga-2">
          <v-icon :color="cls.color" size="x-small">mdi-checkbox-blank-circle</v-icon>
          <span class="text-caption text-medium-emphasis font-mono">{{ cls.code }}</span>
          <span class="text-caption text-medium-emphasis">&middot; {{ cls.label }}</span>
        </div>
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.hover-bg-grey-lighten-4:hover {
  background: #f5f5f5;
}
.border-top {
  border-top: 1px solid #ded9ca;
}
</style>
