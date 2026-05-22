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
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <h1 class="cyber-heading text-xl lg:text-2xl">System Status</h1>
      <div class="flex gap-2">
        <button @click="control('restart')" :disabled="!!actionLoading" class="cyber-btn flex items-center gap-2">
          <RotateCcw class="w-4 h-4" :stroke-width="1.5" /> Restart
        </button>
      </div>
    </div>

    <div v-if="status" class="grid grid-cols-2 lg:grid-cols-5 gap-3">
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Uptime</p>
        <p class="text-lg text-neon text-glow mt-1 font-mono">{{ status.uptime || '—' }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Tunnel</p>
        <p class="text-sm text-gray-200 mt-1 truncate">{{ status.tunnel_name || 'n/a' }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Ingress</p>
        <p class="text-lg text-cyan mt-1" style="text-shadow: 0 0 8px #00d4ff80">{{ status.ingress_count || 0 }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Memory</p>
        <p class="text-sm text-gray-200 mt-1 font-mono">{{ status.memory_usage || '—' }}</p>
      </div>
      <div class="cyber-card-glow">
        <p class="text-[10px] text-muted uppercase tracking-widest">Version</p>
        <p class="text-sm text-magenta mt-1" style="text-shadow: 0 0 8px #ff00ff80">{{ status.version?.split(' ')[2] || '—' }}</p>
      </div>
    </div>

    <div v-if="status?.hostnames?.length" class="cyber-card">
      <h2 class="text-xs text-muted uppercase tracking-widest mb-3">// Active Hostnames</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <div v-for="h in status.hostnames" :key="h" class="flex items-center gap-2 px-3 py-2 bg-void border border-border/50">
          <span class="w-1.5 h-1.5 bg-neon shadow-neon animate-pulse"></span>
          <span class="text-sm text-neon/80 font-mono truncate">{{ h }}</span>
        </div>
      </div>
    </div>

    <div v-if="actionLoading" class="text-xs text-muted tracking-wider">
      > executing {{ actionLoading }}...<span class="animate-blink text-neon">_</span>
    </div>

    <div v-if="projects.length" class="cyber-card">
      <h2 class="text-xs text-muted uppercase tracking-widest mb-3">// Projects ({{ projects.length }})</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <router-link v-for="p in projects" :key="p.id" to="/projects" class="flex items-center gap-2 px-3 py-2 bg-void border border-border/50 hover:border-neon/50 transition-colors">
          <span class="w-1.5 h-1.5 bg-cyan shadow-neon"></span>
          <span class="text-sm text-gray-200 font-mono truncate">{{ p.name }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>
