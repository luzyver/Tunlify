<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'
import { useActionLogStore } from '../stores/actionLog'
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
const actionLog = useActionLogStore()

async function control(action: string) {
  if (actionLoading.value) return
  actionLoading.value = action
  error.value = ''
  const logId = actionLog.start(`${action} tunnel`, 'cloudflared')
  actionLog.append(logId, `Executing ${action} on cloudflared container...`)
  try {
    await apiFetch(`/api/control/${action}`, { method: 'POST' })
    actionLog.append(logId, `Successfully executed ${action}`)
    await new Promise((r) => setTimeout(r, 1500))
    await fetchStatus()
    actionLog.end(logId, 'success')
  } catch (e: any) {
    actionLog.append(logId, `Error: ${e.message}`)
    actionLog.end(logId, 'error')
    error.value = e.message
  } finally {
    actionLoading.value = ''
  }
}

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
  { key: 'created_at', label: 'Created', sortable: true, align: 'right', width: '140px', cellClass: 'num text-caption', headerClass: 'num' },
]

fetchStatus()
fetchProjects()
const interval = setInterval(fetchStatus, 5000)
onUnmounted(() => clearInterval(interval))
</script>

<template>
  <div class="pa-6">
    <div class="d-flex align-center justify-space-between ga-4 mb-8 flex-wrap" style="gap: 16px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Status</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Tunnel</h1>
      </div>
      <button
        class="d-inline-flex align-center ga-2 rounded-pill px-5 font-mono"
        style="height: 44px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
        :disabled="!!actionLoading"
        @click="control('restart')"
        @mouseenter="$event.target.style.transform = 'scale(1.03)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)'"
        @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
      >
        <v-icon size="16" :class="{ 'spin': actionLoading === 'restart' }">mdi-refresh</v-icon>
        {{ actionLoading === 'restart' ? 'Restarting...' : 'Restart tunnel' }}
      </button>
    </div>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <div v-if="status" class="mb-8">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px;">
        <div class="rounded-2xl pa-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6); transition: all 0.3s;" @mouseenter="$event.currentTarget.style.borderColor = 'rgba(247, 147, 26, 0.3)'; $event.currentTarget.style.boxShadow = '0 0 30px -10px rgba(247, 147, 26, 0.15)'" @mouseleave="$event.currentTarget.style.borderColor = 'rgba(30, 41, 59, 0.6)'; $event.currentTarget.style.boxShadow = 'none'">
          <p class="eyebrow mb-2">Uptime</p>
          <p class="font-heading font-semibold tabular-nums" style="color: white; font-size: 24px;">{{ status.uptime || '\u2014' }}</p>
        </div>
        <div class="rounded-2xl pa-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6); transition: all 0.3s;" @mouseenter="$event.currentTarget.style.borderColor = 'rgba(247, 147, 26, 0.3)'; $event.currentTarget.style.boxShadow = '0 0 30px -10px rgba(247, 147, 26, 0.15)'" @mouseleave="$event.currentTarget.style.borderColor = 'rgba(30, 41, 59, 0.6)'; $event.currentTarget.style.boxShadow = 'none'">
          <p class="eyebrow mb-2">Ingress rules</p>
          <p class="font-heading font-semibold tabular-nums" style="color: white; font-size: 24px;">{{ status.ingress_count ?? 0 }}</p>
        </div>
        <div class="rounded-2xl pa-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6); transition: all 0.3s;" @mouseenter="$event.currentTarget.style.borderColor = 'rgba(247, 147, 26, 0.3)'; $event.currentTarget.style.boxShadow = '0 0 30px -10px rgba(247, 147, 26, 0.15)'" @mouseleave="$event.currentTarget.style.borderColor = 'rgba(30, 41, 59, 0.6)'; $event.currentTarget.style.boxShadow = 'none'">
          <p class="eyebrow mb-2">Tunnel</p>
          <p class="font-heading font-semibold font-mono text-truncate" style="color: white; font-size: 24px;">{{ status.tunnel_name || 'n/a' }}</p>
        </div>
        <div class="rounded-2xl pa-5" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6); transition: all 0.3s;" @mouseenter="$event.currentTarget.style.borderColor = 'rgba(247, 147, 26, 0.3)'; $event.currentTarget.style.boxShadow = '0 0 30px -10px rgba(247, 147, 26, 0.15)'" @mouseleave="$event.currentTarget.style.borderColor = 'rgba(30, 41, 59, 0.6)'; $event.currentTarget.style.boxShadow = 'none'">
          <p class="eyebrow mb-2">Memory</p>
          <p class="font-heading font-semibold tabular-nums" style="color: white; font-size: 24px;">{{ memory.used }}</p>
          <p v-if="memory.total" class="font-mono" style="color: #94A3B8; font-size: 12px; margin-top: 2px;">of {{ memory.total }}</p>
        </div>
      </div>
    </div>

    <div v-if="hostnameRows.length" class="mb-8">
      <div class="d-flex align-center justify-space-between mb-4">
        <p class="eyebrow">Active hostnames</p>
        <span class="font-mono" style="color: #94A3B8; font-size: 11px;">{{ hostnameRows.length }} total</span>
      </div>
      <DataTable
        :data="hostnameRows"
        :columns="hostnameCols"
        :show-pagination="false"
        :row-key="(row) => row.hostname"
      >
        <template #cell-status>
          <v-icon size="14" color="#FFD600">mdi-check-circle</v-icon>
        </template>
        <template #cell-hostname="{ row }">
          <span class="font-mono" style="color: white;">{{ row.hostname }}</span>
        </template>
        <template #cell-open="{ row }">
          <a
            :href="`https://${row.hostname}`"
            target="_blank"
            rel="noopener"
            style="color: #94A3B8; transition: color 0.2s;"
            @mouseenter="$event.target.style.color = '#F7931A'"
            @mouseleave="$event.target.style.color = '#94A3B8'"
          >
            <v-icon size="14">mdi-open-in-new</v-icon>
          </a>
        </template>
      </DataTable>
    </div>

    <div v-if="projects.length" class="mb-8">
      <div class="d-flex align-center justify-space-between mb-4">
        <p class="eyebrow">Projects</p>
        <router-link to="/projects" class="font-mono text-decoration-none" style="color: #94A3B8; font-size: 12px; transition: color 0.2s;" @mouseenter="$event.target.style.color = '#F7931A'" @mouseleave="$event.target.style.color = '#94A3B8'">
          View all &rarr;
        </router-link>
      </div>
      <DataTable
        :data="projects"
        :columns="projectCols"
        :show-pagination="false"
      >
        <template #cell-name="{ row }">
          <span class="font-weight-medium" style="color: white; cursor: pointer;" @click="router.push('/projects')">{{ row.name }}</span>
        </template>
        <template #cell-path="{ row }">
          <span class="font-mono" style="color: #94A3B8; font-size: 12px; max-width: 260px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ row.path }}</span>
        </template>
        <template #cell-created_at="{ row }">
          <span class="font-mono" style="color: #94A3B8; font-size: 12px;">{{ row.created_at?.split('T')[0] || '\u2014' }}</span>
        </template>
      </DataTable>
    </div>

    <p v-if="!status" class="font-mono" style="color: #94A3B8;">Loading status...</p>
  </div>
</template>

<style scoped>
@keyframes spin { to { transform: rotate(360deg); } }
.spin { animation: spin 1s linear infinite; }
</style>
