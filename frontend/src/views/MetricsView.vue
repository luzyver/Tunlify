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
    <h1 class="cyber-heading text-xl lg:text-2xl">Metrics</h1>

    <div v-if="metrics" class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Total Requests</p>
        <p class="text-2xl text-neon text-glow mt-1 font-mono">{{ metrics.total_requests?.toLocaleString() || '0' }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Active Connections</p>
        <p class="text-2xl text-cyan mt-1 font-mono" style="text-shadow: 0 0 8px #00d4ff80">{{ metrics.active_connections || '0' }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Response Codes</p>
        <div class="mt-2 space-y-1">
          <div v-for="(count, code) in metrics.response_codes" :key="code" class="flex justify-between text-xs font-mono">
            <span class="text-muted">{{ code }}</span>
            <span :class="String(code).startsWith('2') ? 'text-neon' : String(code).startsWith('4') ? 'text-magenta' : 'text-danger'">{{ count }}</span>
          </div>
          <p v-if="!Object.keys(metrics.response_codes || {}).length" class="text-muted text-xs">—</p>
        </div>
      </div>
    </div>

    <div v-else class="text-muted text-sm tracking-wider">// LOADING METRICS<span class="animate-blink text-neon">_</span></div>
  </div>
</template>
