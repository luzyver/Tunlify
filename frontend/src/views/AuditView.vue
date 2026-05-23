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
  { key: 'created_at', label: 'Time', sortable: true, width: '200px', cellClass: 'num text-xs text-text-muted' },
  { key: 'action', label: 'Action', sortable: true, width: '160px' },
  { key: 'detail', label: 'Detail', hideBelow: 'md' },
  { key: 'ip_address', label: 'IP', sortable: true, hideBelow: 'md', width: '140px', align: 'right', cellClass: 'num text-xs text-text-dim', headerClass: 'num' },
]
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Audit</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Audit log</h1>
      </div>
      <span v-if="total > entries.length" class="text-2xs text-text-dim tabular-nums">
        showing newest {{ entries.length }} of {{ total }}
      </span>
    </header>

    <DataTable
      :data="entries"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search detail or IP…"
      :page-size="25"
    >
      <template #toolbar>
        <select v-model="filterAction" class="input !w-auto !py-1.5 !text-xs">
          <option v-for="a in actions" :key="a" :value="a">{{ a || 'All actions' }}</option>
        </select>
      </template>

      <template #cell-action="{ row }">
        <span class="badge-accent font-mono">{{ row.action }}</span>
      </template>
      <template #cell-detail="{ row }">
        <span class="text-text-muted text-xs truncate max-w-[420px] block">
          {{ row.detail?.split('\n')[0] || '—' }}
        </span>
      </template>
      <template #cell-ip_address="{ row }">
        {{ row.ip_address || '—' }}
      </template>
      <template #empty>
        {{ filterAction ? 'No entries for this action' : 'No audit entries' }}
      </template>
    </DataTable>
  </div>
</template>
