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
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Health Check</h1>
      <button @click="check" :disabled="loading" class="btn-secondary">
        <RefreshCw class="w-4 h-4" :class="loading && 'animate-spin'" /> Refresh
      </button>
    </div>

    <div class="card !p-0 overflow-hidden">
      <table class="table-tight">
        <thead>
          <tr>
            <th>Hostname</th>
            <th>Service</th>
            <th>Latency</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in results" :key="r.hostname">
            <td class="font-mono text-sm">{{ r.hostname }}</td>
            <td class="text-text-muted">{{ r.service }}</td>
            <td class="font-mono tabular-nums text-text-dim">{{ r.latency }}</td>
            <td><span :class="r.status === 'up' ? 'badge-success' : 'badge-danger'">{{ r.status }}</span></td>
          </tr>
          <tr v-if="!results.length && !loading">
            <td colspan="4" class="text-center text-text-muted py-8">No services</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
