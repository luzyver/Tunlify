<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const metrics = ref<any>(null)

async function fetchMetrics() {
  try { metrics.value = await apiFetch('/api/metrics') } catch {}
}
fetchMetrics()
const interval = setInterval(fetchMetrics, 10000)
onUnmounted(() => clearInterval(interval))
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold tracking-tight">Metrics</h1>

    <div v-if="metrics" class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Total Requests</p>
        <p class="metric-value mt-1">{{ metrics.total_requests?.toLocaleString() || '0' }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Active Connections</p>
        <p class="metric-value mt-1 text-accent">{{ metrics.active_connections || '0' }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Response Codes</p>
        <div class="mt-2 space-y-1">
          <div v-for="(count, code) in metrics.response_codes" :key="code" class="flex justify-between text-xs font-mono tabular-nums">
            <span class="text-text-muted">{{ code }}</span>
            <span :class="String(code).startsWith('2') ? 'text-success' : String(code).startsWith('4') ? 'text-warning' : 'text-danger'">{{ count }}</span>
          </div>
          <p v-if="!Object.keys(metrics.response_codes || {}).length" class="text-text-dim text-xs">—</p>
        </div>
      </div>
    </div>

    <div v-else class="text-text-muted text-sm">Loading metrics…</div>
  </div>
</template>
