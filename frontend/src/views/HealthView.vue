<script setup lang="ts">
import { computed, ref } from 'vue'
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

const summary = computed(() => {
  const up = results.value.filter((r: any) => r.status === 'up').length
  return { up, total: results.value.length }
})
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Health</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Endpoint reachability</h1>
      </div>
      <button class="btn-secondary" :disabled="loading" @click="check">
        <RefreshCw class="w-4 h-4" :stroke-width="1.75" :class="loading && 'animate-spin'" />
        {{ loading ? 'Checking…' : 'Refresh' }}
      </button>
    </header>

    <section v-if="results.length" class="grid grid-cols-2 sm:grid-cols-3 gap-4">
      <div class="account-tile">
        <span class="account-tile-label">Total endpoints</span>
        <span class="account-tile-value">{{ summary.total }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Up</span>
        <span class="account-tile-value">{{ summary.up }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Down</span>
        <span class="account-tile-value">{{ summary.total - summary.up }}</span>
      </div>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Endpoints</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th class="w-8"></th>
            <th>Hostname</th>
            <th>Service</th>
            <th class="num">Latency</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in results" :key="r.hostname">
            <td>
              <span class="dot" :class="r.status === 'up' ? 'bg-success' : 'bg-danger'" aria-hidden="true"></span>
            </td>
            <td class="font-mono text-text">{{ r.hostname }}</td>
            <td class="font-mono text-xs text-text-muted truncate max-w-[260px]">{{ r.service }}</td>
            <td class="num text-text">{{ r.latency || '—' }}</td>
            <td>
              <span :class="r.status === 'up' ? 'badge-success' : 'badge-danger'">
                {{ r.status }}
              </span>
            </td>
          </tr>
          <tr v-if="!results.length && !loading">
            <td colspan="5" class="text-center text-text-muted py-8 text-sm">No endpoints to check</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>
