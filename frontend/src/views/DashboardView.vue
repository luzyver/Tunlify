<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
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
  return raw.split(' ')[2] || raw || '\u2014'
})

const memory = computed(() => {
  const raw = (status.value?.memory_usage || '').trim()
  if (!raw) return { used: '\u2014', total: '' }
  const [used, total] = raw.split('/').map((s: string) => s.trim())
  return { used: used || '\u2014', total: total || '' }
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
  { key: 'created_at', label: 'Created', sortable: true, align: 'right', width: '140px', cellClass: 'num text-caption text-medium-emphasis', headerClass: 'num' },
]

fetchStatus()
fetchProjects()
const interval = setInterval(fetchStatus, 5000)
onUnmounted(() => clearInterval(interval))
</script>

<template>
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-8">
      <div>
        <p class="eyebrow mb-2">Console &middot; Status</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Tunnel</h1>
      </div>
      <v-btn
        color="primary"
        :loading="actionLoading === 'restart'"
        :disabled="!!actionLoading"
        @click="control('restart')"
      >
        <v-icon start>mdi-restart</v-icon>
        {{ actionLoading === 'restart' ? 'Restarting...' : 'Restart tunnel' }}
      </v-btn>
    </div>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <v-row v-if="status" class="mb-8">
      <v-col cols="6" lg>
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Uptime</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ status.uptime || '\u2014' }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="6" lg>
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Ingress rules</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ status.ingress_count ?? 0 }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="6" lg>
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Tunnel</p>
            <p class="text-h4 text-on-surface font-mono text-truncate">{{ status.tunnel_name || 'n/a' }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="6" lg>
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Memory</p>
            <p class="text-h4 text-on-surface tabular-nums">{{ memory.used }}</p>
            <p v-if="memory.total" class="text-caption text-medium-emphasis font-mono tabular-nums">of {{ memory.total }}</p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="6" lg>
        <v-card>
          <v-card-text>
            <p class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-1">Version</p>
            <p class="text-h4 text-on-surface font-mono text-truncate">{{ versionLabel }}</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <div v-if="hostnameRows.length" class="mb-8">
      <div class="d-flex align-center justify-space-between mb-3">
        <p class="eyebrow">Active hostnames</p>
        <span class="text-caption text-medium-emphasis tabular-nums">{{ hostnameRows.length }} total</span>
      </div>
      <DataTable
        :data="hostnameRows"
        :columns="hostnameCols"
        :show-pagination="false"
        :row-key="(row) => row.hostname"
      >
        <template #cell-status>
          <v-icon color="success" size="small">mdi-checkbox-blank-circle</v-icon>
        </template>
        <template #cell-hostname="{ row }">
          <span class="font-mono text-on-surface">{{ row.hostname }}</span>
        </template>
        <template #cell-open="{ row }">
          <a
            :href="`https://${row.hostname}`"
            target="_blank"
            rel="noopener"
            class="text-medium-emphasis text-decoration-none"
          >
            <v-icon size="small">mdi-open-in-new</v-icon>
          </a>
        </template>
      </DataTable>
    </div>

    <div v-if="projects.length" class="mb-8">
      <div class="d-flex align-center justify-space-between mb-3">
        <p class="eyebrow">Projects</p>
        <router-link to="/projects" class="text-caption text-medium-emphasis text-decoration-none">
          View all &rarr;
        </router-link>
      </div>
      <DataTable
        :data="projects"
        :columns="projectCols"
        :show-pagination="false"
      >
        <template #cell-name="{ row }">
          <span class="font-weight-medium text-on-surface" @click="router.push('/projects')">{{ row.name }}</span>
        </template>
        <template #cell-path="{ row }">
          <span class="font-mono text-caption text-medium-emphasis text-truncate" style="max-width: 260px;">{{ row.path }}</span>
        </template>
        <template #cell-created_at="{ row }">
          {{ row.created_at?.split('T')[0] || '\u2014' }}
        </template>
      </DataTable>
    </div>

    <p v-if="!status" class="text-medium-emphasis">Loading status...</p>
  </div>
</template>
