<script setup lang="ts">
import { ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const results = ref<any[]>([])
const loading = ref(false)

async function check() {
  loading.value = true
  try { results.value = await apiFetch('/api/health') } catch {}
  finally { loading.value = false }
}
check()
</script>

<template>
  <div class="space-y-6">
    <header class="border-b border-border pb-6">
      <div class="flex items-center justify-between">
        <div>
          <p class="section-marker mb-2">MONITORING</p>
          <h1 class="editorial-h1 !text-3xl">Health Check</h1>
        </div>
        <button @click="check" :disabled="loading" class="btn-secondary">
          <RefreshCw class="w-4 h-4" :stroke-width="1.5" :class="loading && 'animate-spin'" /> Refresh
        </button>
      </div>
    </header>

    <div class="space-y-3">
      <div v-for="r in results" :key="r.hostname" class="card-hover flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-text">{{ r.hostname }}</p>
          <p class="text-xs text-text-muted mt-0.5">{{ r.service }}</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs text-text-dim font-mono">{{ r.latency }}</span>
          <span :class="r.status === 'up' ? 'badge-success' : 'badge-danger'">{{ r.status }}</span>
        </div>
      </div>
      <div v-if="!results.length && !loading" class="card text-center text-text-muted py-8 text-sm">
        No services to check
      </div>
    </div>
  </div>
</template>
