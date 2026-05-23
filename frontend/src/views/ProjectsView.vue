<script setup lang="ts">
import { nextTick, ref } from 'vue'
import {
  History,
  Pencil,
  Play,
  Plus,
  Rocket,
  RotateCcw,
  Square,
  Trash2,
  X,
} from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()

interface Project {
  id: number
  name: string
  path: string
  repo_url: string
  git_username: string
  git_token: string
  created_at: string
}

interface HistoryEntry {
  id: number
  action: string
  detail: string
  created_at: string
  _open?: boolean
}

const projects = ref<Project[]>([])
const loading = ref(false)
const actionLoading = ref<Record<number, boolean>>({})
const error = ref('')

const showForm = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', path: '', repo_url: '', git_username: '', git_token: '' })

const deployTarget = ref<number | null>(null)
const deployRef = ref('')

const output = ref('')
const outputEl = ref<HTMLElement | null>(null)

const history = ref<HistoryEntry[]>([])
const historyTarget = ref<number | null>(null)

async function load() {
  loading.value = true
  error.value = ''
  try { projects.value = await apiFetch('/api/projects') } catch (e: any) { error.value = e.message }
  finally { loading.value = false }
}

function openAdd() {
  editingId.value = null
  form.value = { name: '', path: '', repo_url: '', git_username: '', git_token: '' }
  showForm.value = true
}

function openEdit(p: Project) {
  editingId.value = p.id
  form.value = {
    name: p.name,
    path: p.path,
    repo_url: p.repo_url,
    git_username: p.git_username,
    git_token: '',
  }
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
  } catch (e: any) {
    error.value = e.message
  }
}

async function remove(id: number) {
  if (!confirm('Delete this project? This cannot be undone.')) return
  try {
    await apiFetch(`/api/projects/${id}`, { method: 'DELETE' })
    load()
  } catch (e: any) { error.value = e.message }
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
      } catch {
        clearInterval(poll)
        actionLoading.value[id] = false
      }
    }, 500)
  } catch (e: any) {
    output.value = e.message
    actionLoading.value[id] = false
    deployTarget.value = null
    deployRef.value = ''
  }
}

function startDeploy(p: Project) {
  if (!p.repo_url) { action(p.id, 'deploy', {}); return }
  deployTarget.value = p.id
  deployRef.value = ''
}

function confirmDeploy() {
  if (deployTarget.value && deployRef.value) action(deployTarget.value, 'deploy', { ref: deployRef.value })
}

async function toggleHistory(id: number) {
  if (historyTarget.value === id) { historyTarget.value = null; return }
  historyTarget.value = id
  try { history.value = await apiFetch(`/api/projects/${id}/history`) }
  catch { history.value = [] }
}

const anyRunning = () => Object.values(actionLoading.value).some(Boolean)

load()
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Compose</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Projects</h1>
      </div>
      <button class="btn-primary" @click="openAdd">
        <Plus class="w-4 h-4" :stroke-width="1.75" />
        New project
      </button>
    </header>

    <div v-if="error" class="alert-danger">{{ error }}</div>

    <section v-if="showForm" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">{{ editingId ? 'Edit project' : 'New project' }}</span>
        <button @click="showForm = false" class="btn-icon-ghost"><X class="w-4 h-4" :stroke-width="1.75" /></button>
      </div>
      <div class="card-body grid sm:grid-cols-2 gap-4">
        <div>
          <label class="field-label">Name</label>
          <input v-model="form.name" class="input" placeholder="my-app" />
        </div>
        <div>
          <label class="field-label">Path</label>
          <input v-model="form.path" class="input font-mono" placeholder="/srv/my-app" />
        </div>
        <div class="sm:col-span-2">
          <label class="field-label">Repository URL <span class="text-text-dim">(optional)</span></label>
          <input v-model="form.repo_url" class="input font-mono" placeholder="https://github.com/user/repo.git" />
        </div>
        <div>
          <label class="field-label">Git username</label>
          <input v-model="form.git_username" class="input" placeholder="git" />
        </div>
        <div>
          <label class="field-label">Git token</label>
          <input v-model="form.git_token" type="password" class="input font-mono" placeholder="ghp_…" />
        </div>
        <div class="sm:col-span-2 flex gap-3 pt-1">
          <button @click="save" class="btn-primary">{{ editingId ? 'Save changes' : 'Create project' }}</button>
          <button @click="showForm = false" class="btn-secondary">Cancel</button>
        </div>
      </div>
    </section>

    <section v-if="deployTarget" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title flex items-center gap-2">
          <Rocket class="w-4 h-4 text-accent" :stroke-width="1.75" />
          Deploy — choose ref
        </span>
        <button @click="deployTarget = null" class="btn-icon-ghost"><X class="w-4 h-4" :stroke-width="1.75" /></button>
      </div>
      <div class="card-body flex gap-3">
        <input
          v-model="deployRef"
          placeholder="main / v1.0.0"
          class="input flex-1 font-mono"
          @keyup.enter="confirmDeploy"
        />
        <button
          @click="confirmDeploy"
          :disabled="!deployRef || actionLoading[deployTarget]"
          class="btn-primary"
        >
          Deploy
        </button>
      </div>
    </section>

    <section v-if="output || anyRunning()" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Output</span>
        <span v-if="anyRunning()" class="badge-accent">
          <span class="dot bg-accent"></span> Running
        </span>
      </div>
      <pre
        ref="outputEl"
        class="font-mono text-2xs leading-5 text-text-muted bg-bg-alt
               max-h-96 overflow-auto p-4 whitespace-pre-wrap scrollbar-thin"
      >{{ output || '(waiting)' }}</pre>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">All projects</span>
        <span class="text-2xs text-text-dim tabular-nums">{{ projects.length }}</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th>Name</th>
            <th>Path</th>
            <th class="hidden lg:table-cell">Repository</th>
            <th class="num">Created</th>
            <th class="text-right pr-4">Actions</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="p in projects" :key="p.id">
            <tr :class="historyTarget === p.id && 'is-selected'">
              <td class="font-medium text-text">{{ p.name }}</td>
              <td class="font-mono text-xs text-text-muted truncate max-w-[260px]">{{ p.path }}</td>
              <td class="hidden lg:table-cell font-mono text-xs text-text-dim truncate max-w-[260px]">
                {{ p.repo_url || '—' }}
              </td>
              <td class="num text-xs text-text-dim">{{ p.created_at?.split('T')[0] || '—' }}</td>
              <td class="text-right pr-4">
                <div class="inline-flex items-center gap-1">
                  <button @click="action(p.id, 'up')" :disabled="actionLoading[p.id]" class="btn-icon-ghost" title="Up">
                    <Play class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="action(p.id, 'down')" :disabled="actionLoading[p.id]" class="btn-icon-ghost" title="Down">
                    <Square class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="action(p.id, 'restart')" :disabled="actionLoading[p.id]" class="btn-icon-ghost" title="Restart">
                    <RotateCcw class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="startDeploy(p)" :disabled="actionLoading[p.id]" class="btn-icon-ghost text-accent hover:bg-accent-soft" title="Deploy">
                    <Rocket class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="toggleHistory(p.id)" class="btn-icon-ghost" title="History">
                    <History class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="openEdit(p)" class="btn-icon-ghost" title="Edit">
                    <Pencil class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                  <button @click="remove(p.id)" class="btn-icon-ghost text-danger hover:bg-danger/5" title="Delete">
                    <Trash2 class="w-3.5 h-3.5" :stroke-width="1.75" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="historyTarget === p.id">
              <td colspan="5" class="!p-0 bg-bg-alt">
                <div class="px-6 py-4 space-y-2">
                  <p class="eyebrow-dim">History</p>
                  <div v-if="history.length" class="space-y-1">
                    <div v-for="h in history" :key="h.id" class="text-xs">
                      <button
                        type="button"
                        class="w-full flex items-center justify-between gap-3 py-1
                               text-text-muted hover:text-text transition-colors"
                        :disabled="h.action !== 'project_deploy'"
                        @click="h._open = !h._open"
                      >
                        <span class="font-mono">{{ h.action }}</span>
                        <span class="num text-text-dim">{{ h.created_at }}</span>
                      </button>
                      <pre
                        v-if="h._open"
                        class="font-mono text-2xs whitespace-pre-wrap text-text-muted
                               border-l-2 border-border pl-3 ml-1 my-1"
                      >{{ h.detail }}</pre>
                    </div>
                  </div>
                  <p v-else class="text-xs text-text-dim">No history yet.</p>
                </div>
              </td>
            </tr>
</template>
          <tr v-if="!projects.length && !loading">
            <td colspan="5" class="text-center text-text-muted py-12 text-sm">
              No projects yet — create one above.
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>
