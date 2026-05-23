<script setup lang="ts">
import { ref } from 'vue'
import { Eye, RotateCcw, X } from 'lucide-vue-next'
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
  if (!confirm('Restore this backup? Current configuration will be overwritten.')) return
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
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Config</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Backups</h1>
      </div>
      <span class="text-2xs text-text-dim tabular-nums">{{ backups.length }} snapshots</span>
    </header>

    <div v-if="message" class="alert-success">{{ message }}</div>

    <DataTable
      :data="backups"
      :columns="columns"
      :page-size="20"
      :row-class="(row) => previewId === row.id ? 'is-selected' : undefined"
    >
      <template #cell-actions="{ row }">
        <div class="inline-flex items-center gap-2">
          <button @click="showPreview(row.id)" class="btn-secondary !py-1 !px-2.5 !text-xs">
            <Eye class="w-3.5 h-3.5" :stroke-width="1.75" />
            View
          </button>
          <button
            @click="restore(row.id)"
            :disabled="restoring === row.id"
            class="btn-primary !py-1 !px-2.5 !text-xs"
          >
            <RotateCcw class="w-3.5 h-3.5" :stroke-width="1.75" />
            {{ restoring === row.id ? 'Restoring…' : 'Restore' }}
          </button>
        </div>
      </template>
      <template #empty>No backups yet</template>
    </DataTable>

    <section v-if="previewId !== null" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Preview</span>
        <button @click="closePreview" class="btn-icon-ghost"><X class="w-4 h-4" :stroke-width="1.75" /></button>
      </div>
      <pre
        class="font-mono text-2xs leading-5 text-text-muted bg-bg-alt
               max-h-96 overflow-auto p-4 whitespace-pre scrollbar-thin"
      >{{ previewContent }}</pre>
    </section>
  </div>
</template>
