<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { Plus, Trash2, Play, Square, RotateCcw, Rocket, Pencil, History } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()

interface Project { id: number; name: string; path: string; repo_url: string; git_username: string; git_token: string; created_at: string }

const projects = ref<Project[]>([])
const loading = ref(false)
const actionLoading = ref<Record<number, boolean>>({})
const showForm = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', path: '', repo_url: '', git_username: '', git_token: '' })
const deployRef = ref('')
const deployTarget = ref<number | null>(null)
const output = ref('')
const outputEl = ref<HTMLElement | null>(null)
const history = ref<any[]>([])
const historyTarget = ref<number | null>(null)

async function load() {
  loading.value = true
  try { projects.value = await apiFetch('/api/projects') } catch {}
  finally { loading.value = false }
}

function openAdd() {
  editingId.value = null
  form.value = { name: '', path: '', repo_url: '', git_username: '', git_token: '' }
  showForm.value = true
}

function openEdit(p: Project) {
  editingId.value = p.id
  form.value = { name: p.name, path: p.path, repo_url: p.repo_url, git_username: p.git_username, git_token: '' }
  showForm.value = true
}

async function save() {
  if (!form.value.name || !form.value.path) return
  try {
    if (editingId.value) {
      await apiFetch(`/api/projects/${editingId.value}`, { method: 'PUT', body: JSON.stringify(form.value) })
    } else {
      await apiFetch('/api/projects', { method: 'POST', body: JSON.stringify(form.value) })
    }
    showForm.value = false
    load()
  } catch {}
}

async function remove(id: number) {
  if (!confirm('Delete this project?')) return
  await apiFetch(`/api/projects/${id}`, { method: 'DELETE' })
  load()
}

async function action(id: number, act: string, body?: object) {
  if (actionLoading.value[id]) return
  actionLoading.value[id] = true
  output.value = ''
  try {
    await apiFetch(`/api/projects/${id}/${act}`, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
    // Poll output
    const poll = setInterval(async () => {
      try {
        const res: any = await apiFetch(`/api/projects/${id}/output`)
        output.value = (res.lines || []).join('\n')
        nextTick(() => { if (outputEl.value) outputEl.value.scrollTop = outputEl.value.scrollHeight })
        if (res.done) {
          clearInterval(poll)
          actionLoading.value[id] = false
          deployTarget.value = null
          deployRef.value = ''
        }
      } catch { clearInterval(poll); actionLoading.value[id] = false }
    }, 500)
  } catch (e: any) {
    output.value = e.message
    actionLoading.value[id] = false
    deployTarget.value = null
    deployRef.value = ''
  }
}

function startDeploy(p: Project) {
  if (!p.repo_url) {
    // No git repo — just pull latest images
    action(p.id, 'deploy', {})
    return
  }
  deployTarget.value = p.id
  deployRef.value = ''
}

function confirmDeploy() {
  if (deployTarget.value && deployRef.value) {
    action(deployTarget.value, 'deploy', { ref: deployRef.value })
  }
}

async function toggleHistory(id: number) {
  if (historyTarget.value === id) {
    historyTarget.value = null
    return
  }
  historyTarget.value = id
  try { history.value = await apiFetch(`/api/projects/${id}/history`) } catch { history.value = [] }
}

load()
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="cyber-heading text-xl lg:text-2xl">Projects</h1>
      <button @click="openAdd" class="cyber-btn-secondary flex items-center gap-2">
        <Plus class="w-4 h-4" :stroke-width="1.5" /> Add
      </button>
    </div>

    <!-- Add/Edit Form -->
    <div v-if="showForm" class="cyber-card space-y-3">
      <p class="text-sm text-neon uppercase tracking-wider">{{ editingId ? 'Edit' : 'New' }} Project</p>
      <input v-model="form.name" placeholder="Project name" class="cyber-input w-full" />
      <input v-model="form.path" placeholder="/path/to/project" class="cyber-input w-full" />
      <input v-model="form.repo_url" placeholder="https://github.com/user/repo.git (optional)" class="cyber-input w-full" />
      <div class="grid grid-cols-2 gap-2">
        <input v-model="form.git_username" placeholder="Git username (optional)" class="cyber-input" />
        <input v-model="form.git_token" type="password" placeholder="Git token (optional)" class="cyber-input" />
      </div>
      <div class="flex gap-2">
        <button @click="save" class="cyber-btn flex items-center gap-2">
          <Plus class="w-4 h-4" /> {{ editingId ? 'Save' : 'Create' }}
        </button>
        <button @click="showForm = false" class="cyber-btn-secondary">Cancel</button>
      </div>
    </div>

    <!-- Deploy Modal -->
    <div v-if="deployTarget" class="cyber-card space-y-3 border-neon">
      <p class="text-sm text-neon uppercase tracking-wider">Deploy — Enter branch or tag</p>
      <div class="flex gap-2">
        <input v-model="deployRef" placeholder="main / v1.0.0" class="cyber-input flex-1" @keyup.enter="confirmDeploy" :disabled="actionLoading[deployTarget]" />
        <button @click="confirmDeploy" :disabled="!deployRef || actionLoading[deployTarget]" class="cyber-btn flex items-center gap-2">
          <Rocket class="w-4 h-4" /> Deploy
        </button>
        <button @click="deployTarget = null" :disabled="actionLoading[deployTarget]" class="cyber-btn-secondary">Cancel</button>
      </div>
    </div>

    <!-- Output -->
    <div v-if="output || Object.values(actionLoading).some(v => v)" class="cyber-card">
      <p v-if="Object.values(actionLoading).some(v => v)" class="text-xs text-neon animate-pulse mb-2">// RUNNING...</p>
      <pre ref="outputEl" class="text-xs text-muted font-mono whitespace-pre-wrap max-h-96 overflow-y-auto">{{ output }}</pre>
    </div>

    <!-- Project List -->
    <div class="space-y-2">
      <div v-for="p in projects" :key="p.id" class="cyber-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-200 font-mono">{{ p.name }}</p>
            <p class="text-xs text-muted mt-0.5 font-mono">{{ p.path }}</p>
            <p v-if="p.repo_url" class="text-xs text-muted/60 mt-0.5 font-mono">{{ p.repo_url }}</p>
          </div>
          <div class="flex items-center gap-2">
            <button @click="action(p.id, 'up')" :disabled="actionLoading[p.id]" class="cyber-btn-sm" title="Up">
              <Play class="w-3.5 h-3.5" />
            </button>
            <button @click="action(p.id, 'down')" :disabled="actionLoading[p.id]" class="cyber-btn-sm" title="Down">
              <Square class="w-3.5 h-3.5" />
            </button>
            <button @click="action(p.id, 'restart')" :disabled="actionLoading[p.id]" class="cyber-btn-sm" title="Restart">
              <RotateCcw class="w-3.5 h-3.5" />
            </button>
            <button @click="startDeploy(p)" :disabled="actionLoading[p.id]" class="cyber-btn-sm text-neon" title="Deploy">
              <Rocket class="w-3.5 h-3.5" />
            </button>
            <button @click="toggleHistory(p.id)" class="cyber-btn-sm" title="History">
              <History class="w-3.5 h-3.5" />
            </button>
            <button @click="openEdit(p)" class="cyber-btn-sm" title="Edit">
              <Pencil class="w-3.5 h-3.5" />
            </button>
            <button @click="remove(p.id)" class="cyber-btn-sm text-danger" title="Delete">
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
        <!-- History -->
        <div v-if="historyTarget === p.id" class="mt-3 border-t border-border pt-3">
          <div v-for="h in history" :key="h.id" class="py-1">
            <div class="flex justify-between text-xs text-muted" :class="h.action === 'project_deploy' && 'cursor-pointer'" @click="h.action === 'project_deploy' && (h._open = !h._open)">
              <span class="font-mono">{{ h.action }}</span>
              <span>{{ h.created_at }}</span>
            </div>
            <pre v-if="h._open && h.action === 'project_deploy'" class="text-xs text-muted/70 font-mono whitespace-pre-wrap mt-1 ml-2 border-l border-border pl-2">{{ h.detail }}</pre>
          </div>
          <p v-if="!history.length" class="text-xs text-muted">// NO HISTORY</p>
        </div>
      </div>
      <div v-if="!projects.length && !loading" class="cyber-card text-center text-muted py-8 tracking-wider text-sm">
        // NO PROJECTS
      </div>
    </div>
  </div>
</template>
