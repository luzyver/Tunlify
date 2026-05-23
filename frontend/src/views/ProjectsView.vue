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
    if (editingId.value) await apiFetch(`/api/projects/${editingId.value}`, { method: 'PUT', body: JSON.stringify(form.value) })
    else await apiFetch('/api/projects', { method: 'POST', body: JSON.stringify(form.value) })
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
    await apiFetch(`/api/projects/${id}/${act}`, { method: 'POST', body: body ? JSON.stringify(body) : undefined })
    const poll = setInterval(async () => {
      try {
        const res: any = await apiFetch(`/api/projects/${id}/output`)
        output.value = (res.lines || []).join('\n')
        nextTick(() => { if (outputEl.value) outputEl.value.scrollTop = outputEl.value.scrollHeight })
        if (res.done) { clearInterval(poll); actionLoading.value[id] = false; deployTarget.value = null; deployRef.value = '' }
      } catch { clearInterval(poll); actionLoading.value[id] = false }
    }, 500)
  } catch (e: any) { output.value = e.message; actionLoading.value[id] = false; deployTarget.value = null; deployRef.value = '' }
}

function startDeploy(p: Project) {
  if (!p.repo_url) { action(p.id, 'deploy', {}); return }
  deployTarget.value = p.id; deployRef.value = ''
}

function confirmDeploy() {
  if (deployTarget.value && deployRef.value) action(deployTarget.value, 'deploy', { ref: deployRef.value })
}

async function toggleHistory(id: number) {
  if (historyTarget.value === id) { historyTarget.value = null; return }
  historyTarget.value = id
  try { history.value = await apiFetch(`/api/projects/${id}/history`) } catch { history.value = [] }
}

load()
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Projects</h1>
      <button @click="openAdd" class="btn-secondary"><Plus class="w-4 h-4" /> Add</button>
    </div>

    <!-- Form -->
    <div v-if="showForm" class="card card-body space-y-3">
      <p class="text-sm font-bold">{{ editingId ? 'Edit' : 'New' }} Project</p>
      <input v-model="form.name" placeholder="Project name" class="input" />
      <input v-model="form.path" placeholder="/path/to/project" class="input" />
      <input v-model="form.repo_url" placeholder="https://github.com/user/repo.git (optional)" class="input" />
      <div class="grid grid-cols-2 gap-3">
        <input v-model="form.git_username" placeholder="Git username" class="input" />
        <input v-model="form.git_token" type="password" placeholder="Git token" class="input" />
      </div>
      <div class="flex gap-2">
        <button @click="save" class="btn-primary">{{ editingId ? 'Save' : 'Create' }}</button>
        <button @click="showForm = false" class="btn-secondary">Cancel</button>
      </div>
    </div>

    <!-- Deploy -->
    <div v-if="deployTarget" class="card card-body border-accent space-y-3">
      <p class="text-sm font-bold">Deploy — enter branch or tag</p>
      <div class="flex gap-2">
        <input v-model="deployRef" placeholder="main / v1.0.0" class="input flex-1" @keyup.enter="confirmDeploy" />
        <button @click="confirmDeploy" :disabled="!deployRef || actionLoading[deployTarget]" class="btn-primary"><Rocket class="w-4 h-4" /> Deploy</button>
        <button @click="deployTarget = null" class="btn-secondary">Cancel</button>
      </div>
    </div>

    <!-- Output -->
    <div v-if="output || Object.values(actionLoading).some(v => v)" class="card card-body">
      <p v-if="Object.values(actionLoading).some(v => v)" class="text-xs text-accent font-medium mb-2">Running…</p>
      <pre ref="outputEl" class="text-xs text-text-muted font-mono whitespace-pre-wrap max-h-96 overflow-y-auto">{{ output }}</pre>
    </div>

    <!-- List -->
    <div class="space-y-3">
      <div v-for="p in projects" :key="p.id" class="card card-body">
        <div class="flex items-center justify-between gap-4">
          <div class="min-w-0">
            <p class="text-sm font-bold truncate">{{ p.name }}</p>
            <p class="text-xs text-text-muted font-mono mt-0.5 truncate">{{ p.path }}</p>
            <p v-if="p.repo_url" class="text-xs text-text-dim font-mono mt-0.5 truncate">{{ p.repo_url }}</p>
          </div>
          <div class="flex items-center gap-1 shrink-0">
            <button @click="action(p.id, 'up')" :disabled="actionLoading[p.id]" class="btn-icon" title="Up"><Play class="w-3.5 h-3.5" /></button>
            <button @click="action(p.id, 'down')" :disabled="actionLoading[p.id]" class="btn-icon" title="Down"><Square class="w-3.5 h-3.5" /></button>
            <button @click="action(p.id, 'restart')" :disabled="actionLoading[p.id]" class="btn-icon" title="Restart"><RotateCcw class="w-3.5 h-3.5" /></button>
            <button @click="startDeploy(p)" :disabled="actionLoading[p.id]" class="btn-icon !text-accent !border-accent/30" title="Deploy"><Rocket class="w-3.5 h-3.5" /></button>
            <button @click="toggleHistory(p.id)" class="btn-icon" title="History"><History class="w-3.5 h-3.5" /></button>
            <button @click="openEdit(p)" class="btn-icon" title="Edit"><Pencil class="w-3.5 h-3.5" /></button>
            <button @click="remove(p.id)" class="btn-icon !text-danger !border-danger/30" title="Delete"><Trash2 class="w-3.5 h-3.5" /></button>
          </div>
        </div>
        <div v-if="historyTarget === p.id" class="mt-3 border-t border-border pt-3">
          <div v-for="h in history" :key="h.id" class="py-1">
            <div class="flex justify-between text-xs text-text-muted" :class="h.action === 'project_deploy' && 'cursor-pointer'" @click="h.action === 'project_deploy' && (h._open = !h._open)">
              <span class="font-mono">{{ h.action }}</span>
              <span class="tabular-nums">{{ h.created_at }}</span>
            </div>
            <pre v-if="h._open" class="text-xs text-text-dim font-mono whitespace-pre-wrap mt-1 ml-3 border-l-2 border-border pl-3">{{ h.detail }}</pre>
          </div>
          <p v-if="!history.length" class="text-xs text-text-muted">No history</p>
        </div>
      </div>
      <div v-if="!projects.length && !loading" class="card card-body text-center text-text-muted py-8 text-sm">No projects</div>
    </div>
  </div>
</template>
