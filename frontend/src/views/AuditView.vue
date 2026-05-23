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
    <header class="border-b border-border pb-6">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div>
          <p class="section-marker mb-2">CHANGELOG</p>
          <h1 class="editorial-h1 !text-3xl">Audit Log</h1>
        </div>
        <select v-model="filterAction" class="input !w-auto !py-1.5 !text-xs">
          <option v-for="a in actions" :key="a" :value="a">{{ a || 'All actions' }}</option>
        </select>
      </div>
    </header>

    <div class="card !p-0 overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-border bg-surface">
            <th class="text-left px-4 py-3 text-xs font-medium text-text-muted">Time</th>
            <th class="text-left px-4 py-3 text-xs font-medium text-text-muted">Action</th>
            <th class="text-left px-4 py-3 text-xs font-medium text-text-muted hidden md:table-cell">Detail</th>
            <th class="text-left px-4 py-3 text-xs font-medium text-text-muted hidden md:table-cell">IP</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in entries" :key="e.id" class="border-b border-border/50 hover:bg-surface">
            <td class="px-4 py-3 text-xs text-text-muted font-mono">{{ e.created_at }}</td>
            <td class="px-4 py-3"><span class="badge-muted font-mono">{{ e.action }}</span></td>
            <td class="px-4 py-3 text-xs text-text-muted hidden md:table-cell">{{ e.detail?.split('\n')[0] || '—' }}</td>
            <td class="px-4 py-3 text-xs text-text-dim font-mono hidden md:table-cell">{{ e.ip_address }}</td>
          </tr>
          <tr v-if="!entries.length">
            <td colspan="4" class="px-4 py-12 text-center text-text-muted text-sm">No entries</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex items-center justify-between">
      <span class="text-xs text-text-muted">{{ total }} records</span>
      <div class="flex gap-2">
        <button @click="page--" :disabled="page === 0" class="btn-secondary !px-3 !py-1.5 !text-xs">Prev</button>
        <button @click="page++" :disabled="(page + 1) * limit >= total" class="btn-secondary !px-3 !py-1.5 !text-xs">Next</button>
      </div>
    </div>
  </div>
</template>
