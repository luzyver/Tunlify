<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const entries = ref<any[]>([])
const total = ref(0)
const page = ref(0)
const limit = 20
const filterAction = ref('')

async function fetchAudit() {
  const params = new URLSearchParams({ limit: String(limit), offset: String(page.value * limit) })
  if (filterAction.value) params.set('action', filterAction.value)
  const data = await apiFetch<{ entries: any[]; total: number }>(`/api/audit?${params}`)
  entries.value = data.entries || []
  total.value = data.total
}

watch([page, filterAction], fetchAudit)
fetchAudit()

const actions = ['', 'login', 'logout', 'restart', 'config_edit']
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <h1 class="cyber-heading text-xl lg:text-2xl">Audit Log</h1>
      <select v-model="filterAction" class="cyber-input !w-auto !py-1.5 !text-xs">
        <option v-for="a in actions" :key="a" :value="a">{{ a || '// ALL' }}</option>
      </select>
    </div>

    <div class="cyber-card !p-0 overflow-hidden">
      <table class="w-full text-xs font-mono">
        <thead>
          <tr class="border-b border-border bg-chrome/50">
            <th class="text-left px-4 py-3 text-muted uppercase tracking-widest">Time</th>
            <th class="text-left px-4 py-3 text-muted uppercase tracking-widest">Action</th>
            <th class="text-left px-4 py-3 text-muted uppercase tracking-widest hidden md:table-cell">Detail</th>
            <th class="text-left px-4 py-3 text-muted uppercase tracking-widest hidden md:table-cell">Origin</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in entries" :key="e.id" class="border-b border-border/50 hover:bg-neon/5">
            <td class="px-4 py-3 text-muted">{{ e.created_at }}</td>
            <td class="px-4 py-3">
              <span class="cyber-badge border-cyan text-cyan">{{ e.action }}</span>
            </td>
            <td class="px-4 py-3 text-muted hidden md:table-cell">{{ e.detail || '—' }}</td>
            <td class="px-4 py-3 text-muted hidden md:table-cell">{{ e.ip_address }}</td>
          </tr>
          <tr v-if="!entries.length">
            <td colspan="4" class="px-4 py-12 text-center text-muted tracking-wider">// NO ENTRIES</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex items-center justify-between">
      <span class="text-xs text-muted tracking-wider">{{ total }} records</span>
      <div class="flex gap-2">
        <button @click="page--" :disabled="page === 0" class="cyber-btn-secondary !px-3 !py-1.5 !text-[10px] disabled:opacity-30">PREV</button>
        <button @click="page++" :disabled="(page + 1) * limit >= total" class="cyber-btn-secondary !px-3 !py-1.5 !text-[10px] disabled:opacity-30">NEXT</button>
      </div>
    </div>
  </div>
</template>
