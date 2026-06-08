<script setup lang="ts">
import { computed, ref } from 'vue'
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
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Health</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Endpoint reachability</h1>
      </div>
      <v-btn variant="tonal" :loading="loading" @click="check">
        <v-icon start>mdi-refresh</v-icon>
        {{ loading ? 'Checking...' : 'Refresh' }}
      </v-btn>
    </div>

    <v-row v-if="results.length" class="mb-6">
      <v-col cols="4" sm="3">
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Total endpoints</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ summary.total }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="4" sm="3">
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Up</p>
            <p class="text-h4 text-success tabular-nums">{{ summary.up }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="4" sm="3">
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Down</p>
            <p class="text-h4 text-error tabular-nums">{{ summary.total - summary.up }}</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <DataTable
      :data="results"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search hostname or service..."
      :page-size="25"
      :row-key="(row) => row.hostname"
    >
      <template #cell-status_dot="{ row }">
        <v-icon :color="row.status === 'up' ? 'success' : 'error'" size="small">mdi-checkbox-blank-circle</v-icon>
      </template>
      <template #cell-hostname="{ row }">
        <span class="font-mono text-on-surface">{{ row.hostname }}</span>
      </template>
      <template #cell-service="{ row }">
        <span class="font-mono text-caption text-medium-emphasis text-truncate" style="max-width: 260px;">{{ row.service }}</span>
      </template>
      <template #cell-latency_ms="{ row }">
        <span class="tabular-nums">{{ row.latency || '\u2014' }}</span>
      </template>
      <template #cell-status="{ row }">
        <v-chip :color="row.status === 'up' ? 'success' : 'error'" size="x-small" variant="tonal">{{ row.status }}</v-chip>
      </template>
      <template #empty>
        {{ loading ? 'Checking endpoints...' : 'No endpoints to check' }}
      </template>
    </DataTable>
  </div>
</template>
