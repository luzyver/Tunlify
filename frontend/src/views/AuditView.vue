<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'

const { apiFetch } = useApi()

interface AuditEntry {
  id: number
  action: string
  detail?: string
  ip_address?: string
  created_at: string
}

const entries = ref<AuditEntry[]>([])
const filterAction = ref('')
const limit = 500
const total = ref(0)

const actions = ['', 'login', 'logout', 'restart', 'config_edit', 'project_deploy']

async function fetchAudit() {
  const params = new URLSearchParams({ limit: String(limit), offset: '0' })
  if (filterAction.value) params.set('action', filterAction.value)
  const data = await apiFetch<{ entries: AuditEntry[]; total: number }>(`/api/audit?${params}`)
  entries.value = data.entries || []
  total.value = data.total
}

watch(filterAction, fetchAudit)
fetchAudit()

const columns: Column<AuditEntry>[] = [
  { key: 'created_at', label: 'Time', sortable: true, width: '200px', cellClass: 'num text-caption text-medium-emphasis' },
  { key: 'action', label: 'Action', sortable: true, width: '160px' },
  { key: 'detail', label: 'Detail', hideBelow: 'md' },
  { key: 'ip_address', label: 'IP', sortable: true, hideBelow: 'md', width: '140px', align: 'right', cellClass: 'num text-caption text-disabled', headerClass: 'num' },
]
</script>

<template>
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Audit</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Audit log</h1>
      </div>
      <span v-if="total > entries.length" class="text-caption text-medium-emphasis tabular-nums">
        showing newest {{ entries.length }} of {{ total }}
      </span>
    </div>

    <DataTable
      :data="entries"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search detail or IP..."
      :page-size="25"
    >
      <template #toolbar>
        <v-select
          v-model="filterAction"
          :items="actions.map(a => ({ title: a || 'All actions', value: a }))"
          hide-details
          density="compact"
          variant="outlined"
          style="min-width: 160px;"
        />
      </template>

      <template #cell-action="{ row }">
        <v-chip color="primary" size="x-small" variant="tonal" class="font-mono">{{ row.action }}</v-chip>
      </template>
      <template #cell-detail="{ row }">
        <span class="text-medium-emphasis text-caption text-truncate" style="max-width: 420px;">
          {{ row.detail?.split('\n')[0] || '\u2014' }}
        </span>
      </template>
      <template #cell-ip_address="{ row }">
        {{ row.ip_address || '\u2014' }}
      </template>
      <template #empty>
        {{ filterAction ? 'No entries for this action' : 'No audit entries' }}
      </template>
    </DataTable>
  </div>
</template>
