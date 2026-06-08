<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'

const { apiFetch } = useApi()

interface Backup { id: number; created_at: string }

const backups = ref<Backup[]>([])
const previewId = ref<number | null>(null)
const previewContent = ref('')
const message = ref('')
const restoring = ref<number | null>(null)

async function load() {
  try { backups.value = await apiFetch('/api/config/backups') } catch {}
}

async function showPreview(id: number) {
  previewId.value = id
  const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
  previewContent.value = b.content
}

function closePreview() {
  previewId.value = null
  previewContent.value = ''
}

async function restore(id: number) {
  const confirmed = await confirm('Restore this backup? Current configuration will be overwritten.')
  if (!confirmed) return
  restoring.value = id
  try {
    const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
    await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify({ content: b.content }) })
    message.value = 'Configuration restored'
    closePreview()
    setTimeout(() => (message.value = ''), 3000)
  } finally {
    restoring.value = null
  }
}

load()

const columns: Column<Backup>[] = [
  { key: 'created_at', label: 'Created', sortable: true, cellClass: 'font-mono num' },
  { key: 'actions', label: '', align: 'right', width: '220px' },
]
</script>

<template>
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Config</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Backups</h1>
      </div>
      <span class="text-caption text-medium-emphasis tabular-nums">{{ backups.length }} snapshots</span>
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
          <v-btn variant="tonal" size="small" @click="showPreview(row.id)">
            <v-icon start size="x-small">mdi-eye</v-icon>
            View
          </v-btn>
          <v-btn
            color="primary"
            size="small"
            :loading="restoring === row.id"
            :disabled="restoring === row.id"
            @click="restore(row.id)"
          >
            <v-icon start size="x-small">mdi-restore</v-icon>
            {{ restoring === row.id ? 'Restoring...' : 'Restore' }}
          </v-btn>
        </div>
      </template>
      <template #empty>No backups yet</template>
    </DataTable>

    <v-card v-if="previewId !== null" class="mt-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Preview</span>
        <v-btn icon="mdi-close" variant="text" size="small" @click="closePreview" />
      </v-card-title>
      <pre
        class="font-mono text-caption text-medium-emphasis bg-grey-lighten-3 pa-4 overflow-auto scrollbar-thin"
        style="max-height: 384px; white-space: pre; line-height: 1.25;"
      >{{ previewContent }}</pre>
    </v-card>
  </div>
</template>
