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
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Audit Log</h1>
      <select v-model="filterAction" class="input !w-auto !py-1.5">
        <option v-for="a in actions" :key="a" :value="a">{{ a || 'All actions' }}</option>
      </select>
    </div>

    <div class="card !p-0 overflow-hidden">
      <table class="table-tight">
        <thead>
          <tr>
            <th>Time</th>
            <th>Action</th>
            <th class="hidden md:table-cell">Detail</th>
            <th class="hidden md:table-cell">IP</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in entries" :key="e.id">
            <td class="font-mono tabular-nums text-text-muted text-xs">{{ e.created_at }}</td>
            <td><span class="badge-accent">{{ e.action }}</span></td>
            <td class="text-text-muted text-xs hidden md:table-cell">{{ e.detail?.split('\n')[0] || '—' }}</td>
            <td class="font-mono text-text-dim text-xs hidden md:table-cell">{{ e.ip_address }}</td>
          </tr>
          <tr v-if="!entries.length">
            <td colspan="4" class="text-center text-text-muted py-8">No entries</td>
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
