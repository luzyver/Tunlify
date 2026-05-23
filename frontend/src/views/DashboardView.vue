<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue'
import { ExternalLink, RotateCcw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const status = ref<any>(null)
const projects = ref<any[]>([])
const actionLoading = ref('')
const error = ref('')

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
        <span class="account-tile-value">{{ status.uptime || '—' }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Ingress rules</span>
        <span class="account-tile-value">{{ status.ingress_count ?? 0 }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Tunnel</span>
        <span class="account-tile-value-sm truncate block">{{ status.tunnel_name || 'n/a' }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Memory</span>
        <span class="account-tile-meta block">{{ status.memory_usage || '—' }}</span>
      </div>
      <div class="account-tile">
        <span class="account-tile-label">Version</span>
        <span class="account-tile-meta block truncate">{{ versionLabel }}</span>
      </div>
    </section>

    <section v-if="status?.hostnames?.length" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Active hostnames</span>
        <span class="text-2xs text-text-dim tabular-nums">{{ status.hostnames.length }} total</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th class="w-8"></th>
            <th>Hostname</th>
            <th class="num">Open</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="h in status.hostnames" :key="h">
            <td class="w-8">
              <span class="dot bg-success" aria-hidden="true"></span>
              <span class="sr-only">Live</span>
            </td>
            <td class="font-mono text-text">{{ h }}</td>
            <td class="num">
              <a
                :href="`https://${h}`"
                target="_blank"
                rel="noopener"
                class="inline-flex items-center gap-1 text-text-muted hover:text-accent transition-colors"
              >
                <ExternalLink class="w-3.5 h-3.5" :stroke-width="1.75" />
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section v-if="projects.length" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Projects</span>
        <router-link to="/projects" class="text-xs text-text-muted hover:text-accent transition-colors">
          View all →
        </router-link>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th>Project</th>
            <th>Path</th>
            <th class="num">Created</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="p in projects"
            :key="p.id"
            class="cursor-pointer"
            @click="$router.push('/projects')"
          >
            <td class="font-medium text-text">{{ p.name }}</td>
            <td class="font-mono text-xs text-text-muted truncate max-w-[260px]">{{ p.path }}</td>
            <td class="num text-text-dim text-xs">{{ p.created_at?.split('T')[0] || '—' }}</td>
          </tr>
        </tbody>
      </table>
    </section>

    <p v-if="!status" class="text-sm text-text-muted">Loading status…</p>
  </div>
</template>
