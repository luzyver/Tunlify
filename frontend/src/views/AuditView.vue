<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()

const entries = ref<any[]>([])
const total = ref(0)
const page = ref(0)
const limit = 20
const filterAction = ref('')

const actions = ['', 'login', 'logout', 'restart', 'config_edit', 'project_deploy']

async function fetchAudit() {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(page.value * limit),
  })
  if (filterAction.value) params.set('action', filterAction.value)
  const data = await apiFetch<{ entries: any[]; total: number }>(`/api/audit?${params}`)
  entries.value = data.entries || []
  total.value = data.total
}

watch([page, filterAction], fetchAudit)
fetchAudit()

const rangeLabel = computed(() => {
  if (!total.value) return '0'
  const from = page.value * limit + 1
  const to = Math.min((page.value + 1) * limit, total.value)
  return `${from}–${to} of ${total.value}`
})

const hasNext = computed(() => (page.value + 1) * limit < total.value)
const hasPrev = computed(() => page.value > 0)
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Audit</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Audit log</h1>
      </div>
      <div class="flex items-center gap-2">
        <select v-model="filterAction" class="input !w-auto !py-1.5 !text-xs">
          <option v-for="a in actions" :key="a" :value="a">{{ a || 'All actions' }}</option>
        </select>
      </div>
    </header>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Entries</span>
        <span class="text-2xs text-text-dim tabular-nums">{{ rangeLabel }}</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th class="w-[180px]">Time</th>
            <th class="w-[160px]">Action</th>
            <th class="hidden md:table-cell">Detail</th>
            <th class="hidden md:table-cell w-[140px] num">IP</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in entries" :key="e.id">
            <td class="num text-xs text-text-muted">{{ e.created_at }}</td>
            <td>
              <span class="badge-accent font-mono">{{ e.action }}</span>
            </td>
            <td class="hidden md:table-cell text-text-muted text-xs truncate max-w-[420px]">
              {{ e.detail?.split('\n')[0] || '—' }}
            </td>
            <td class="hidden md:table-cell num text-xs text-text-dim">{{ e.ip_address || '—' }}</td>
          </tr>
          <tr v-if="!entries.length">
            <td colspan="4" class="text-center text-text-muted py-8 text-sm">No audit entries</td>
          </tr>
        </tbody>
      </table>
    </section>

    <div class="flex items-center justify-between gap-3">
      <span class="text-2xs text-text-dim tabular-nums">{{ rangeLabel }}</span>
      <div class="flex items-center gap-2">
        <button @click="page--" :disabled="!hasPrev" class="btn-secondary !py-1 !px-2.5 !text-xs">
          <ChevronLeft class="w-3.5 h-3.5" :stroke-width="1.75" />
          Previous
        </button>
        <button @click="page++" :disabled="!hasNext" class="btn-secondary !py-1 !px-2.5 !text-xs">
          Next
          <ChevronRight class="w-3.5 h-3.5" :stroke-width="1.75" />
        </button>
      </div>
    </div>
  </div>
</template>
