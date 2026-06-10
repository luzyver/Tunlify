<script setup lang="ts">
import { ref } from 'vue'
import { useApi, reportError } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'


const { apiFetch } = useApi()

interface Backup { id: number; created_at: string }

const backups = ref<Backup[]>([])
const previewId = ref<number | null>(null)
const previewContent = ref('')
const message = ref('')
const restoring = ref<number | null>(null)

async function load() {
  try { backups.value = await apiFetch('/api/config/backups') } catch (e) { reportError('failed to load backups', e) }
}

async function showPreview(id: number) {
  previewId.value = id
  const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
  previewContent.value = b.content
}

function closePreview() { previewId.value = null; previewContent.value = '' }

async function restore(id: number) {
  if (!confirm('Restore this backup? Current configuration will be overwritten.')) return
  restoring.value = id
  try {
    const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
    await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify({ content: b.content }) })
    message.value = 'Configuration restored'
    closePreview()
    setTimeout(() => (message.value = ''), 3000)
  } finally { restoring.value = null }
}

load()

const columns: Column<Backup>[] = [
  { key: 'created_at', label: 'Created', sortable: true, cellClass: 'font-mono num' },
  { key: 'actions', label: '', align: 'right', width: '220px' },
]
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <span class="page-badge">Console &middot; Config</span>
        <h1 class="page-title">Backups</h1>
      </div>
      <span class="stat-label" style="margin-bottom: 0; color: #71717A;">{{ backups.length }} snapshots</span>
    </div>

    <v-alert v-if="message" type="success" class="mb-4" variant="tonal">{{ message }}</v-alert>

    <DataTable
      :data="backups"
      :columns="columns"
      :page-size="20"
      :row-class="(row) => previewId === row.id ? 'is-selected' : undefined"
    >
      <template #cell-actions="{ row }">
        <div class="d-inline-flex align-center ga-2">
          <button class="btn btn-secondary btn-sm" @click="showPreview(row.id)">
            <v-icon size="12">mdi-eye</v-icon>
            View
          </button>
          <button
            class="btn btn-primary btn-sm"
            :disabled="restoring === row.id"
            @click="restore(row.id)"
          >
            <v-icon size="12">mdi-refresh</v-icon>
            {{ restoring === row.id ? 'Restoring...' : 'Restore' }}
          </button>
        </div>
      </template>
      <template #empty>No backups yet</template>
    </DataTable>

    <div v-if="previewId !== null" class="card mt-6">
      <div class="card-header">
        <span class="card-title">Preview</span>
        <button class="btn-ghost" @click="closePreview">
          <v-icon size="16">mdi-close</v-icon>
        </button>
      </div>
      <pre class="font-mono pa-4 overflow-auto scrollbar-thin" style="color: #94A3B8; font-size: 12px; line-height: 1.25; max-height: 384px; white-space: pre;">{{ previewContent }}</pre>
    </div>
  </div>
</template>
