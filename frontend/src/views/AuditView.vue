<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'


const { apiFetch } = useApi()

interface AuditEntry { id: number; action: string; detail?: string; ip_address?: string; created_at: string }

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
  { key: 'created_at', label: 'Time', sortable: true, width: '200px', cellClass: 'num font-mono text-caption' },
  { key: 'action', label: 'Action', sortable: true, width: '160px' },
  { key: 'detail', label: 'Detail', hideBelow: 'md' },
  { key: 'ip_address', label: 'IP', sortable: true, hideBelow: 'md', width: '140px', align: 'right', cellClass: 'num font-mono text-caption' },
]
</script>

<template>
  <div class="pa-6">
    <div class="d-flex align-center justify-space-between ga-4 mb-6 flex-wrap" style="gap: 16px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Audit</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Audit log</h1>
      </div>
      <span v-if="total > entries.length" class="font-mono tabular-nums" style="color: #94A3B8; font-size: 11px;">
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
        <div class="position-relative" style="min-width: 160px;">
          <select v-model="filterAction" style="width: 100%; background: rgba(0,0,0,0.5); border: 1px solid rgba(30, 41, 59, 0.8); border-radius: 8px; color: white; padding: 6px 12px; font-family: 'JetBrains Mono', monospace; font-size: 12px; height: 32px; outline: none; cursor: pointer; appearance: auto;">
            <option v-for="a in actions" :key="a" :value="a" style="background: #0F1115;">{{ a || 'All actions' }}</option>
          </select>
        </div>
      </template>

      <template #cell-action="{ row }">
        <span class="rounded-pill px-2 font-mono" style="background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 11px; line-height: 22px;">{{ row.action }}</span>
      </template>
      <template #cell-detail="{ row }">
        <span class="font-mono" style="color: #94A3B8; font-size: 12px; max-width: 420px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
          {{ row.detail?.split('\n')[0] || '\u2014' }}
        </span>
      </template>
      <template #cell-ip_address="{ row }">
        <span class="font-mono" style="color: #94A3B8; font-size: 12px;">{{ row.ip_address || '\u2014' }}</span>
      </template>
      <template #empty>
        {{ filterAction ? 'No entries for this action' : 'No audit entries' }}
      </template>
    </DataTable>
  </div>
</template>
