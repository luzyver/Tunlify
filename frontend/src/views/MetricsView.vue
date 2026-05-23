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
const interval = setInterval(fetchMetrics, 10000)
onUnmounted(() => clearInterval(interval))

const responseRows = computed(() => {
  const codes = metrics.value?.response_codes || {}
  return Object.entries(codes).map(([code, count]) => ({ code, count: Number(count) }))
})

function codeClass(code: string) {
  if (code.startsWith('2')) return 'text-success'
  if (code.startsWith('3')) return 'text-text-muted'
  if (code.startsWith('4')) return 'text-warning'
  return 'text-danger'
}

function formatNumber(n: number | undefined) {
  if (n === undefined || n === null) return '—'
  return n.toLocaleString('en-US')
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

    <section class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div class="account-tile">
        <span class="account-tile-label">Total requests</span>
        <span class="account-tile-value">{{ formatNumber(metrics?.total_requests) }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Active connections</span>
        <span class="account-tile-value">{{ formatNumber(metrics?.active_connections) }}</span>
      </div>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Response codes</span>
        <span class="text-2xs text-text-dim">Auto-refresh · 10s</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th class="w-[120px]">Code</th>
            <th>Class</th>
            <th class="num">Count</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in responseRows" :key="row.code">
            <td class="font-mono" :class="codeClass(row.code)">{{ row.code }}</td>
            <td class="text-text-muted text-xs">
              <template v-if="row.code.startsWith('2')">Success</template>
              <template v-else-if="row.code.startsWith('3')">Redirect</template>
              <template v-else-if="row.code.startsWith('4')">Client error</template>
              <template v-else-if="row.code.startsWith('5')">Server error</template>
              <template v-else>Other</template>
            </td>
            <td class="num text-text">{{ formatNumber(row.count) }}</td>
          </tr>
          <tr v-if="!responseRows.length">
            <td colspan="3" class="text-center text-text-muted py-8 text-sm">
              {{ metrics ? 'No response codes recorded yet' : 'Loading metrics…' }}
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>
