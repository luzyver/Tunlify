<script setup lang="ts">
import { computed, ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'
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
    results.value = raw.map((r) => ({
      ...r,
      latency_ms: parseLatency(r.latency),
    }))
  } catch {}
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

    <DataTable
      :data="results"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search hostname or service…"
      :page-size="25"
      :row-key="(row) => row.hostname"
    >
      <template #cell-status_dot="{ row }">
        <span class="dot" :class="row.status === 'up' ? 'bg-success' : 'bg-danger'" aria-hidden="true"></span>
      </template>
      <template #cell-hostname="{ row }">
        <span class="font-mono text-text">{{ row.hostname }}</span>
      </template>
      <template #cell-service="{ row }">
        <span class="font-mono text-xs text-text-muted truncate max-w-[260px] block">{{ row.service }}</span>
      </template>
      <template #cell-latency_ms="{ row }">
        <span class="num">{{ row.latency || '—' }}</span>
      </template>
      <template #cell-status="{ row }">
        <span :class="row.status === 'up' ? 'badge-success' : 'badge-danger'">{{ row.status }}</span>
      </template>
      <template #empty>
        {{ loading ? 'Checking endpoints…' : 'No endpoints to check' }}
      </template>
    </DataTable>
  </div>
</template>
