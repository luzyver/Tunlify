<script setup lang="ts">
import { computed, ref } from 'vue'
import { useApi, reportError } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'


const { apiFetch } = useApi()

interface HealthRow {
  hostname: string
  service: string
  latency: string
  latency_ms?: number
  status: string
}

const results = ref<HealthRow[]>([])
const loading = ref(false)

async function check() {
  loading.value = true
  try {
    const raw = (await apiFetch<HealthRow[]>('/api/health')) || []
    results.value = raw.map((r) => ({ ...r, latency_ms: parseLatency(r.latency) }))
  } catch (e) { reportError('failed to run health check', e) }
  finally { loading.value = false }
}
check()

function parseLatency(v?: string): number {
  if (!v) return Number.POSITIVE_INFINITY
  const n = parseFloat(v)
  if (Number.isNaN(n)) return Number.POSITIVE_INFINITY
  if (v.includes('ms')) return n
  if (v.includes('s')) return n * 1000
  return n
}

const summary = computed(() => {
  const up = results.value.filter((r) => r.status === 'up').length
  return { up, total: results.value.length }
})

const columns: Column<HealthRow>[] = [
  { key: 'status_dot', label: '', width: '32px' },
  { key: 'hostname', label: 'Hostname', sortable: true },
  { key: 'service', label: 'Service', sortable: true, hideBelow: 'md' },
  { key: 'latency_ms', label: 'Latency', sortable: true, align: 'right', cellClass: 'num', headerClass: 'num' },
  { key: 'status', label: 'Status', sortable: true, width: '120px' },
]
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <span class="page-badge">Console &middot; Health</span>
        <h1 class="page-title">Endpoint reachability</h1>
      </div>
      <button class="btn btn-secondary" :disabled="loading" @click="check">
        <v-icon size="14" :class="{ 'spin': loading }">mdi-refresh</v-icon>
        {{ loading ? 'Checking...' : 'Refresh' }}
      </button>
    </div>

    <div v-if="results.length" class="mb-4">
      <div class="stat-grid" style="grid-template-columns: repeat(3, 1fr);">
        <div class="card stat-card card-hover">
          <p class="stat-label">Total endpoints</p>
          <p class="stat-value tabular-nums">{{ summary.total }}</p>
        </div>
        <div class="card stat-card card-hover">
          <p class="stat-label">Up</p>
          <p class="stat-value stat-value--gold tabular-nums">{{ summary.up }}</p>
        </div>
        <div class="card stat-card card-hover">
          <p class="stat-label">Down</p>
          <p class="stat-value stat-value--red tabular-nums">{{ summary.total - summary.up }}</p>
        </div>
      </div>
    </div>

    <DataTable
      :data="results"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search hostname or service..."
      :page-size="25"
      :row-key="(row) => row.hostname"
    >
      <template #cell-status_dot="{ row }">
        <v-icon v-if="row.status === 'up'" size="14" color="#F7931A">mdi-check-circle</v-icon>
        <v-icon v-else size="14" color="#EF4444">mdi-close-circle</v-icon>
      </template>
      <template #cell-hostname="{ row }">
        <span class="font-mono" style="color: white;">{{ row.hostname }}</span>
      </template>
      <template #cell-service="{ row }">
        <span class="font-mono" style="color: #94A3B8; font-size: 12px; max-width: 260px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ row.service }}</span>
      </template>
      <template #cell-latency_ms="{ row }">
        <span class="tabular-nums" style="color: #94A3B8;">{{ row.latency || '\u2014' }}</span>
      </template>
      <template #cell-status="{ row }">
        <span class="rounded-pill px-2 font-mono" :style="{ background: row.status === 'up' ? 'rgba(247, 147, 26, 0.1)' : 'rgba(239, 68, 68, 0.1)', border: '1px solid ' + (row.status === 'up' ? 'rgba(247, 147, 26, 0.2)' : 'rgba(239, 68, 68, 0.3)'), color: row.status === 'up' ? '#F7931A' : '#EF4444', fontSize: '11px' }">{{ row.status }}</span>
      </template>
      <template #empty>
        {{ loading ? 'Checking endpoints...' : 'No endpoints to check' }}
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.spin { animation: pg-spin 1s linear infinite; }
</style>
