<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const backups = ref<any[]>([])
const preview = ref('')
const message = ref('')

async function load() {
  try { backups.value = await apiFetch('/api/config/backups') } catch {}
}

async function showPreview(id: number) {
  const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
  preview.value = b.content
}

async function restore(id: number) {
  const b = await apiFetch<{ content: string }>(`/api/config/backups/${id}`)
  await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify({ content: b.content }) })
  message.value = 'Configuration restored'
  preview.value = ''
  setTimeout(() => message.value = '', 3000)
}

load()
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold tracking-tight">Config Backups</h1>

    <div v-if="message" class="p-3 rounded text-sm bg-success/10 text-emerald-800 border border-success/20">{{ message }}</div>

    <div class="card !p-0 overflow-hidden">
      <table class="table-tight">
        <thead>
          <tr>
            <th>Created</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in backups" :key="b.id">
            <td class="font-mono tabular-nums text-sm">{{ b.created_at }}</td>
            <td class="text-right">
              <div class="flex gap-2 justify-end">
                <button @click="showPreview(b.id)" class="btn-secondary !px-3 !py-1 !text-xs">View</button>
                <button @click="restore(b.id)" class="btn-primary !px-3 !py-1 !text-xs">Restore</button>
              </div>
            </td>
          </tr>
          <tr v-if="!backups.length">
            <td colspan="2" class="text-center text-text-muted py-8">No backups</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="preview" class="card">
      <div class="card-header flex items-center justify-between">
        <span>Preview</span>
        <button @click="preview = ''" class="btn-secondary !px-3 !py-1 !text-xs">Close</button>
      </div>
      <div class="card-body">
        <pre class="font-mono text-xs text-text bg-surface p-4 rounded overflow-auto max-h-64">{{ preview }}</pre>
      </div>
    </div>
  </div>
</template>
