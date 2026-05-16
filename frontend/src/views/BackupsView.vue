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
  message.value = 'CONFIG RESTORED'
  preview.value = ''
  setTimeout(() => message.value = '', 3000)
}

load()
</script>

<template>
  <div class="space-y-6">
    <h1 class="cyber-heading text-xl lg:text-2xl">Config Backups</h1>

    <div v-if="message" class="border border-neon bg-neon/5 text-neon text-sm p-3 font-mono">[OK] {{ message }}</div>

    <div class="space-y-2">
      <div v-for="b in backups" :key="b.id" class="cyber-card flex items-center justify-between">
        <span class="text-xs text-muted font-mono">{{ b.created_at }}</span>
        <div class="flex gap-2">
          <button @click="showPreview(b.id)" class="cyber-btn-secondary !px-3 !py-1 !text-[10px]">VIEW</button>
          <button @click="restore(b.id)" class="cyber-btn !px-3 !py-1 !text-[10px]">RESTORE</button>
        </div>
      </div>
      <div v-if="!backups.length" class="cyber-card text-center text-muted py-8 text-sm tracking-wider">// NO BACKUPS</div>
    </div>

    <div v-if="preview" class="cyber-card">
      <p class="text-[10px] text-muted uppercase tracking-widest mb-2">// Preview</p>
      <pre class="font-mono text-xs text-neon/80 bg-void p-4 border border-border overflow-auto max-h-64">{{ preview }}</pre>
      <button @click="preview = ''" class="cyber-btn-secondary !px-3 !py-1 !text-[10px] mt-3">CLOSE</button>
    </div>
  </div>
</template>
