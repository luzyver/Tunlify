<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { RotateCcw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const status = ref<any>(null)
const projects = ref<any[]>([])
const actionLoading = ref('')

async function fetchStatus() {
  try { status.value = await apiFetch('/api/status') } catch {}
}
async function fetchProjects() {
  try { projects.value = await apiFetch('/api/projects') } catch {}
}
async function control(action: string) {
  actionLoading.value = action
  try {
    await apiFetch(`/api/control/${action}`, { method: 'POST' })
    await new Promise(r => setTimeout(r, 1500))
    await fetchStatus()
  } catch (e: any) { alert(e.message) }
  finally { actionLoading.value = '' }
}

fetchStatus()
fetchProjects()
const interval = setInterval(fetchStatus, 5000)
onUnmounted(() => clearInterval(interval))
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Status</h1>
      <button @click="control('restart')" :disabled="!!actionLoading" class="btn-primary">
        <RotateCcw class="w-4 h-4" /> Restart
      </button>
    </div>

    <div v-if="status" class="grid grid-cols-2 lg:grid-cols-5 gap-3">
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Uptime</p>
        <p class="metric-value mt-1">{{ status.uptime || '—' }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Tunnel</p>
        <p class="text-sm font-medium mt-1 truncate">{{ status.tunnel_name || 'n/a' }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Ingress</p>
        <p class="metric-value mt-1 text-accent">{{ status.ingress_count || 0 }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Memory</p>
        <p class="text-sm font-mono mt-1 tabular-nums">{{ status.memory_usage || '—' }}</p>
      </div>
      <div class="card card-body">
        <p class="text-2xs text-text-dim uppercase tracking-wide">Version</p>
        <p class="text-sm font-mono mt-1">{{ status.version?.split(' ')[2] || '—' }}</p>
      </div>
    </div>

    <div v-if="status?.hostnames?.length" class="card">
      <div class="card-header">Active Hostnames</div>
      <div class="card-body grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <div v-for="h in status.hostnames" :key="h" class="flex items-center gap-2 px-3 py-2 bg-surface rounded">
          <span class="w-2 h-2 rounded-full bg-success"></span>
          <span class="text-sm font-mono truncate">{{ h }}</span>
        </div>
      </div>
    </div>

    <p v-if="actionLoading" class="text-sm text-text-muted">Running {{ actionLoading }}…</p>

    <div v-if="projects.length" class="card">
      <div class="card-header">Projects ({{ projects.length }})</div>
      <div class="card-body grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <router-link v-for="p in projects" :key="p.id" to="/projects" class="flex items-center gap-2 px-3 py-2 bg-surface rounded hover:bg-bg-alt transition-colors">
          <span class="w-2 h-2 rounded-full bg-chart-5"></span>
          <span class="text-sm font-medium truncate">{{ p.name }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>
