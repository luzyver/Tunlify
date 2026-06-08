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
  <div class="pa-6">
    <div class="d-flex align-center justify-space-between ga-4 mb-6 flex-wrap" style="gap: 16px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Config</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Backups</h1>
      </div>
      <span class="font-mono tabular-nums" style="color: #94A3B8; font-size: 11px;">{{ backups.length }} snapshots</span>
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
          <button
            class="d-inline-flex align-center ga-1 rounded-pill px-3 font-mono"
            style="height: 28px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 11px; cursor: pointer; transition: all 0.2s;"
            @click="showPreview(row.id)"
            @mouseenter="$event.target.style.background = 'rgba(247, 147, 26, 0.2)'"
            @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'"
          >
            <v-icon size="12">mdi-eye</v-icon>
            View
          </button>
          <button
            class="d-inline-flex align-center ga-1 rounded-pill px-3 font-mono"
            :style="{ height: '28px', background: 'linear-gradient(to right, #EA580C, #F7931A)', border: 'none', color: 'white', fontSize: '11px', cursor: restoring === row.id ? 'not-allowed' : 'pointer', boxShadow: '0 0 20px -5px rgba(234, 88, 12, 0.5)', transition: 'all 0.3s', opacity: restoring === row.id ? 0.6 : 1 }"
            :disabled="restoring === row.id"
            @click="restore(row.id)"
            @mouseenter="if(restoring !== row.id) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)' }"
            @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
          >
            <v-icon size="12">mdi-refresh</v-icon>
            {{ restoring === row.id ? 'Restoring...' : 'Restore' }}
          </button>
        </div>
      </template>
      <template #empty>No backups yet</template>
    </DataTable>

    <div v-if="previewId !== null" class="rounded-2xl mt-6" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Preview</span>
        <button @click="closePreview" style="background: none; border: none; color: #94A3B8; cursor: pointer; padding: 4px; transition: color 0.2s;" @mouseenter="$event.target.style.color = '#F7931A'" @mouseleave="$event.target.style.color = '#94A3B8'">
          <v-icon size="16">mdi-close</v-icon>
        </button>
      </div>
      <pre class="font-mono pa-4 overflow-auto scrollbar-thin" style="color: #94A3B8; font-size: 12px; line-height: 1.25; max-height: 384px; white-space: pre;">{{ previewContent }}</pre>
    </div>
  </div>
</template>
