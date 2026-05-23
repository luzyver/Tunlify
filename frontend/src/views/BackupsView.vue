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
    <header class="border-b border-border pb-6">
      <p class="section-marker mb-2">RECOVERY</p>
      <h1 class="editorial-h1 !text-3xl">Config Backups</h1>
    </header>

    <div v-if="message" class="p-3 rounded-md text-sm border border-emerald-200 bg-emerald-50 text-emerald-700">{{ message }}</div>

    <div class="space-y-3">
      <div v-for="b in backups" :key="b.id" class="card-hover flex items-center justify-between">
        <span class="text-sm text-text-muted font-mono">{{ b.created_at }}</span>
        <div class="flex gap-2">
          <button @click="showPreview(b.id)" class="btn-secondary !px-3 !py-1.5 !text-xs">View</button>
          <button @click="restore(b.id)" class="btn-primary !px-3 !py-1.5 !text-xs">Restore</button>
        </div>
      </div>
      <div v-if="!backups.length" class="card text-center text-text-muted py-8 text-sm">No backups available</div>
    </div>

    <div v-if="preview" class="card">
      <div class="flex items-center justify-between mb-3">
        <p class="text-xs text-text-muted font-medium">Preview</p>
        <button @click="preview = ''" class="btn-secondary !px-3 !py-1 !text-xs">Close</button>
      </div>
      <pre class="font-mono text-xs text-text bg-surface p-4 rounded-md overflow-auto max-h-64">{{ preview }}</pre>
    </div>
  </div>
</template>
