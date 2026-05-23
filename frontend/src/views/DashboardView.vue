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
  <div class="space-y-8">
    <!-- Editorial title block -->
    <header class="border-b border-border pb-8">
      <p class="section-marker mb-3">STATUS</p>
      <h1 class="editorial-h1">System Overview</h1>
      <p class="text-lg text-text-muted mt-3 max-w-reading">Tunnel status, active hostnames, and resource metrics at a glance.</p>
      <div class="flex items-center gap-4 mt-6">
        <button @click="control('restart')" :disabled="!!actionLoading" class="btn-primary">
          <RotateCcw class="w-4 h-4" :stroke-width="1.5" /> Restart Tunnel
        </button>
      </div>
    </header>

    <!-- Metrics grid -->
    <div v-if="status" class="grid grid-cols-2 lg:grid-cols-5 gap-4">
      <div class="card-hover">
        <p class="text-xs text-text-dim uppercase tracking-wide">Uptime</p>
        <p class="text-xl font-mono font-medium mt-1">{{ status.uptime || '—' }}</p>
      </div>
      <div class="card-hover">
        <p class="text-xs text-text-dim uppercase tracking-wide">Tunnel</p>
        <p class="text-sm text-text mt-1 truncate">{{ status.tunnel_name || 'n/a' }}</p>
      </div>
      <div class="card-hover">
        <p class="text-xs text-text-dim uppercase tracking-wide">Ingress</p>
        <p class="text-xl font-mono font-medium mt-1 text-accent">{{ status.ingress_count || 0 }}</p>
      </div>
      <div class="card-hover">
        <p class="text-xs text-text-dim uppercase tracking-wide">Memory</p>
        <p class="text-sm font-mono mt-1">{{ status.memory_usage || '—' }}</p>
      </div>
      <div class="card-hover">
        <p class="text-xs text-text-dim uppercase tracking-wide">Version</p>
        <p class="text-sm font-mono mt-1">{{ status.version?.split(' ')[2] || '—' }}</p>
      </div>
    </div>

    <!-- Hostnames -->
    <section v-if="status?.hostnames?.length">
      <p class="section-marker mb-4">HOSTNAMES</p>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <div v-for="h in status.hostnames" :key="h" class="card-hover flex items-center gap-3">
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
          <span class="text-sm font-mono truncate">{{ h }}</span>
        </div>
      </div>
    </section>

    <!-- Loading indicator -->
    <p v-if="actionLoading" class="text-sm text-text-muted">
      Executing {{ actionLoading }}…
    </p>

    <!-- Projects -->
    <section v-if="projects.length">
      <p class="section-marker mb-4">PROJECTS</p>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <router-link v-for="p in projects" :key="p.id" to="/projects" class="card-hover flex items-center gap-3">
          <span class="w-2 h-2 rounded-full bg-accent"></span>
          <span class="text-sm font-medium truncate">{{ p.name }}</span>
        </router-link>
      </div>
    </section>
  </div>
</template>
