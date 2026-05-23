<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ExternalLink, RotateCcw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'

const { apiFetch } = useApi()
const router = useRouter()
const status = ref<any>(null)
const projects = ref<any[]>([])
const actionLoading = ref('')
const error = ref('')

interface HostRow { hostname: string }
interface ProjectRow { id: number; name: string; path: string; created_at: string }

async function fetchStatus() {
  try { status.value = await apiFetch('/api/status') } catch {}
}
async function fetchProjects() {
  try { projects.value = await apiFetch('/api/projects') } catch {}
}
async function control(action: string) {
  if (actionLoading.value) return
  actionLoading.value = action
  error.value = ''
  try {
    await apiFetch(`/api/control/${action}`, { method: 'POST' })
    await new Promise((r) => setTimeout(r, 1500))
    await fetchStatus()
  } catch (e: any) {
    error.value = e.message
  } finally {
    actionLoading.value = ''
  }
}

const versionLabel = computed(() => {
  const raw = status.value?.version || ''
  return raw.split(' ')[2] || raw || '—'
})

const memory = computed(() => {
  const raw = (status.value?.memory_usage || '').trim()
  if (!raw) return { used: '—', total: '' }
  const [used, total] = raw.split('/').map((s: string) => s.trim())
  return { used: used || '—', total: total || '' }
})

const hostnameRows = computed<HostRow[]>(() =>
  (status.value?.hostnames || []).map((h: string) => ({ hostname: h }))
)

const hostnameCols: Column<HostRow>[] = [
  { key: 'status', label: '', width: '32px' },
  { key: 'hostname', label: 'Hostname', sortable: true },
  { key: 'open', label: '', align: 'right', width: '60px' },
]

const projectCols: Column<ProjectRow>[] = [
  { key: 'name', label: 'Project', sortable: true },
  { key: 'path', label: 'Path', sortable: true },
  { key: 'created_at', label: 'Created', sortable: true, align: 'right', width: '140px', cellClass: 'num text-xs text-text-dim', headerClass: 'num' },
]

fetchStatus()
fetchProjects()
const interval = setInterval(fetchStatus, 5000)
onUnmounted(() => clearInterval(interval))
</script>

<template>
  <div class="space-y-8">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Status</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Tunnel</h1>
      </div>
      <button
        type="button"
        class="btn-primary"
        :disabled="!!actionLoading"
        @click="control('restart')"
      >
        <RotateCcw class="w-4 h-4" :stroke-width="1.75" :class="actionLoading === 'restart' && 'animate-spin'" />
        {{ actionLoading === 'restart' ? 'Restarting…' : 'Restart tunnel' }}
      </button>
    </header>

    <div v-if="error" class="alert-danger">{{ error }}</div>

    <section v-if="status" class="grid grid-cols-2 lg:grid-cols-5 gap-4">
      <div class="account-tile">
        <span class="account-tile-label">Uptime</span>
        <span class="account-tile-value" :title="status.uptime || ''">{{ status.uptime || '—' }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Ingress rules</span>
        <span class="account-tile-value">{{ status.ingress_count ?? 0 }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Tunnel</span>
        <span class="account-tile-value" :title="status.tunnel_name || ''">{{ status.tunnel_name || 'n/a' }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Memory</span>
        <span class="account-tile-value" :title="status.memory_usage || ''">{{ memory.used }}</span>
        <span v-if="memory.total" class="account-tile-sub">of {{ memory.total }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Version</span>
        <span class="account-tile-value" :title="versionLabel">{{ versionLabel }}</span>
      </div>
    </section>

    <section v-if="hostnameRows.length">
      <div class="flex items-center justify-between mb-3">
        <p class="eyebrow">Active hostnames</p>
        <span class="text-2xs text-text-dim tabular-nums">{{ hostnameRows.length }} total</span>
      </div>
      <DataTable
        :data="hostnameRows"
        :columns="hostnameCols"
        :show-pagination="false"
        :row-key="(row) => row.hostname"
      >
        <template #cell-status>
          <span class="dot bg-success" aria-hidden="true"></span>
          <span class="sr-only">Live</span>
        </template>
        <template #cell-hostname="{ row }">
          <span class="font-mono text-text">{{ row.hostname }}</span>
        </template>
        <template #cell-open="{ row }">
          <a
            :href="`https://${row.hostname}`"
            target="_blank"
            rel="noopener"
            class="inline-flex items-center gap-1 text-text-muted hover:text-accent transition-colors"
            @click.stop
          >
            <ExternalLink class="w-3.5 h-3.5" :stroke-width="1.75" />
          </a>
        </template>
      </DataTable>
    </section>

    <section v-if="projects.length">
      <div class="flex items-center justify-between mb-3">
        <p class="eyebrow">Projects</p>
        <router-link to="/projects" class="text-xs text-text-muted hover:text-accent transition-colors">
          View all →
        </router-link>
      </div>
      <DataTable
        :data="projects"
        :columns="projectCols"
        :show-pagination="false"
        :row-class="() => 'cursor-pointer'"
      >
        <template #cell-name="{ row }">
          <span class="font-medium text-text" @click="router.push('/projects')">{{ row.name }}</span>
        </template>
        <template #cell-path="{ row }">
          <span class="font-mono text-xs text-text-muted truncate max-w-[260px] block">{{ row.path }}</span>
        </template>
        <template #cell-created_at="{ row }">
          {{ row.created_at?.split('T')[0] || '—' }}
        </template>
      </DataTable>
    </section>

    <p v-if="!status" class="text-sm text-text-muted">Loading status…</p>
  </div>
</template>
