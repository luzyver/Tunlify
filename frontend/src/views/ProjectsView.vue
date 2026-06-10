<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { useApi } from '../composables/useApi'
import { useActionLogStore } from '../stores/actionLog'
import DataTable, { type Column } from '../components/DataTable.vue'


const { apiFetch } = useApi()
const actionLog = useActionLogStore()

interface Project {
  id: number; name: string; path: string; repo_url: string
  git_username: string; git_token: string; created_at: string
}
interface HistoryEntry {
  id: number; action: string; detail: string; created_at: string; _open?: boolean
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
const historyName = ref('')

async function load() {
  loading.value = true; error.value = ''
  try { projects.value = await apiFetch('/api/projects') } catch (e: any) { error.value = e.message }
  finally { loading.value = false }
}

function openAdd() {
  editingId.value = null; form.value = { name: '', path: '', repo_url: '', git_username: '', git_token: '' }; showForm.value = true
}
function openEdit(p: Project) {
  editingId.value = p.id; form.value = { name: p.name, path: p.path, repo_url: p.repo_url, git_username: p.git_username, git_token: '' }; showForm.value = true
}

async function save() {
  if (!form.value.name || !form.value.path) return
  try {
    if (editingId.value) await apiFetch(`/api/projects/${editingId.value}`, { method: 'PUT', body: JSON.stringify(form.value) })
    else await apiFetch('/api/projects', { method: 'POST', body: JSON.stringify(form.value) })
    showForm.value = false; load()
  } catch (e: any) { error.value = e.message }
}

async function remove(id: number) {
  if (!confirm('Delete this project? This cannot be undone.')) return
  try { await apiFetch(`/api/projects/${id}`, { method: 'DELETE' }); load() }
  catch (e: any) { error.value = e.message }
}

async function action(id: number, act: string, body?: object) {
  if (actionLoading.value[id]) return
  actionLoading.value[id] = true; output.value = ''
  const project = projects.value.find((p) => p.id === id)
  const label = act === 'deploy' ? `deploy ${project?.name || id}` : `${act} ${project?.name || id}`
  const logId = actionLog.start(label, project?.name)
  actionLog.append(logId, `Starting ${act} on project ${project?.name || id}...`)
  try {
    await apiFetch(`/api/projects/${id}/${act}`, { method: 'POST', body: body ? JSON.stringify(body) : undefined })
    const poll = setInterval(async () => {
      try {
        const res: any = await apiFetch(`/api/projects/${id}/output`)
        const lines: string[] = res.lines || []
        if (lines.length) actionLog.appendLines(logId, lines)
        output.value = lines.join('\n')
        nextTick(() => { if (outputEl.value) outputEl.value.scrollTop = outputEl.value.scrollHeight })
        if (res.done) {
          clearInterval(poll); actionLoading.value[id] = false; deployTarget.value = null; deployRef.value = ''
          actionLog.append(logId, 'Completed successfully')
          actionLog.end(logId, 'success')
        }
      } catch {
        clearInterval(poll); actionLoading.value[id] = false
        actionLog.append(logId, 'Connection lost while polling')
        actionLog.end(logId, 'error')
      }
    }, 500)
  } catch (e: any) {
    output.value = e.message; actionLoading.value[id] = false; deployTarget.value = null; deployRef.value = ''
    actionLog.append(logId, `Error: ${e.message}`)
    actionLog.end(logId, 'error')
  }
}

function startDeploy(p: Project) {
  if (!p.repo_url) { action(p.id, 'deploy', {}); return }
  deployTarget.value = p.id; deployRef.value = ''
}
function confirmDeploy() {
  if (deployTarget.value && deployRef.value) action(deployTarget.value, 'deploy', { ref: deployRef.value })
}
async function toggleHistory(p: Project) {
  if (historyTarget.value === p.id) { historyTarget.value = null; return }
  historyTarget.value = p.id; historyName.value = p.name
  try { history.value = await apiFetch(`/api/projects/${p.id}/history`) } catch { history.value = [] }
}

const anyRunning = () => Object.values(actionLoading.value).some(Boolean)

load()

const columns: Column<Project>[] = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'path', label: 'Path', sortable: true },
  { key: 'repo_url', label: 'Repository', sortable: true, hideBelow: 'lg' },
  { key: 'created_at', label: 'Created', sortable: true, align: 'right', width: '140px', cellClass: 'num font-mono text-caption', headerClass: 'num' },
  { key: 'actions', label: '', align: 'right', width: '260px' },
]
</script>

<style scoped>
.history-row {
  background: none;
  border: none;
  color: #94A3B8;
  font-size: 12px;
  padding: 6px 8px;
  cursor: pointer;
  text-align: left;
  transition: background 0.2s;
}
.history-row:hover {
  background: rgba(247, 147, 26, 0.03);
}
.running-dot {
  width: 6px;
  height: 6px;
  background: #F7931A;
  border-radius: 50%;
}
</style>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <span class="page-badge">Console &middot; Compose</span>
        <h1 class="page-title">Projects</h1>
      </div>
      <button class="btn btn-primary" @click="openAdd">
        <v-icon size="16">mdi-plus</v-icon>
        New project
      </button>
    </div>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <!-- Form -->
    <div v-if="showForm" class="card mb-4">
      <div class="card-header">
        <span class="card-title">{{ editingId ? 'Edit project' : 'New project' }}</span>
        <button class="btn-ghost" @click="showForm = false"><v-icon size="16">mdi-close</v-icon></button>
      </div>
      <div class="card-body">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
          <div>
            <label class="label-upper">Name</label>
            <input v-model="form.name" placeholder="my-app" class="input-line" />
          </div>
          <div>
            <label class="label-upper">Path</label>
            <input v-model="form.path" placeholder="/srv/my-app" class="input-line" />
          </div>
          <div style="grid-column: span 2;">
            <label class="label-upper">Repository URL <span style="color: #71717A;">(optional)</span></label>
            <input v-model="form.repo_url" placeholder="https://github.com/user/repo.git" class="input-line" />
          </div>
          <div>
            <label class="label-upper">Git username</label>
            <input v-model="form.git_username" placeholder="git" class="input-line" />
          </div>
          <div>
            <label class="label-upper">Git token</label>
            <input v-model="form.git_token" type="password" placeholder="ghp_..." class="input-line" />
          </div>
          <div style="grid-column: span 2;" class="d-flex ga-3 pt-1">
            <button class="btn btn-primary" @click="save">{{ editingId ? 'Save changes' : 'Create project' }}</button>
            <button class="btn btn-secondary" @click="showForm = false">Cancel</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Deploy prompt -->
    <div v-if="deployTarget" class="card mb-4">
      <div class="card-header">
        <div class="d-flex align-center ga-2">
          <v-icon size="18" color="#F7931A">mdi-rocket</v-icon>
          <span class="card-title">Deploy &mdash; choose ref</span>
        </div>
        <button class="btn-ghost" @click="deployTarget = null"><v-icon size="16">mdi-close</v-icon></button>
      </div>
      <div class="card-body d-flex ga-3">
        <input v-model="deployRef" placeholder="main / v1.0.0" class="input-line" style="flex: 1;" @keyup.enter="confirmDeploy" />
        <button class="btn btn-primary" :disabled="!deployRef || actionLoading[deployTarget]" @click="confirmDeploy">
          Deploy
        </button>
      </div>
    </div>

    <!-- Output -->
    <div v-if="output || anyRunning()" class="card mb-4">
      <div class="card-header">
        <span class="card-title">Output</span>
        <span v-if="anyRunning()" class="rounded-pill px-2 font-mono d-inline-flex align-center ga-1" style="background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.3); color: #F7931A; font-size: 11px; line-height: 22px;">
          <span class="rounded-full d-inline-block running-dot"></span> Running
        </span>
      </div>
      <pre ref="outputEl" class="font-mono pa-4 overflow-auto scrollbar-thin" style="color: #94A3B8; background: rgba(0,0,0,0.3); font-size: 12px; line-height: 1.25; max-height: 384px; white-space: pre-wrap;">{{ output || '(waiting)' }}</pre>
    </div>

    <DataTable
      :data="projects"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search projects..."
      :page-size="25"
      :row-class="(row) => historyTarget === row.id ? 'is-selected' : undefined"
    >
      <template #cell-name="{ row }">
        <span class="font-weight-medium" style="color: white;">{{ row.name }}</span>
      </template>
      <template #cell-path="{ row }">
        <span class="font-mono" style="color: #94A3B8; font-size: 12px; max-width: 260px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ row.path }}</span>
      </template>
      <template #cell-repo_url="{ row }">
        <span class="font-mono" style="color: rgba(148, 163, 184, 0.5); font-size: 12px; max-width: 260px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ row.repo_url || '\u2014' }}</span>
      </template>
      <template #cell-created_at="{ row }">
        <span class="font-mono" style="color: #94A3B8; font-size: 12px;">{{ row.created_at?.split('T')[0] || '\u2014' }}</span>
      </template>
      <template #cell-actions="{ row }">
        <div class="d-inline-flex align-center ga-1">
          <button title="Up" class="btn-ghost" :disabled="actionLoading[row.id]" @click="action(row.id, 'up')"><v-icon size="14">mdi-play</v-icon></button>
          <button title="Down" class="btn-ghost" :disabled="actionLoading[row.id]" @click="action(row.id, 'down')"><v-icon size="14">mdi-stop-circle</v-icon></button>
          <button title="Restart" class="btn-ghost" :disabled="actionLoading[row.id]" @click="action(row.id, 'restart')"><v-icon size="14">mdi-refresh</v-icon></button>
          <button title="Deploy" class="btn-ghost" style="color: #F7931A;" :disabled="actionLoading[row.id]" @click="startDeploy(row)"><v-icon size="14">mdi-rocket</v-icon></button>
          <button title="History" class="btn-ghost" @click="toggleHistory(row)"><v-icon size="14">mdi-history</v-icon></button>
          <button title="Edit" class="btn-ghost" @click="openEdit(row)"><v-icon size="14">mdi-pencil</v-icon></button>
          <button title="Delete" class="btn-ghost btn-ghost--danger" @click="remove(row.id)"><v-icon size="14">mdi-delete</v-icon></button>
        </div>
      </template>
      <template #empty>
        {{ loading ? 'Loading projects...' : 'No projects yet \u2014 click "New project" above.' }}
      </template>
    </DataTable>

    <!-- History -->
    <div v-if="historyTarget !== null" class="card mt-6">
      <div class="card-header">
        <span class="card-title">History &middot; <span class="font-mono" style="color: #F7931A;">{{ historyName }}</span></span>
        <button class="btn-ghost" @click="historyTarget = null"><v-icon size="16">mdi-close</v-icon></button>
      </div>
      <div class="pa-4 d-flex flex-column ga-1">
        <div v-if="history.length" class="d-flex flex-column ga-1">
          <div v-for="h in history" :key="h.id">
            <button
              class="w-100 d-flex align-center justify-space-between font-mono px-2 rounded-lg history-row"
              :disabled="h.action !== 'project_deploy'"
              @click="h._open = !h._open"
            >
              <span class="font-mono">{{ h.action }}</span>
              <span class="tabular-nums" style="color: rgba(148, 163, 184, 0.5);">{{ h.created_at }}</span>
            </button>
            <pre v-if="h._open" class="font-mono ml-3 my-1" style="color: #94A3B8; font-size: 11px; border-left: 2px solid rgba(247, 147, 26, 0.3); padding-left: 12px; white-space: pre-wrap;">{{ h.detail }}</pre>
          </div>
        </div>
        <p v-else class="font-mono px-2" style="color: rgba(148, 163, 184, 0.5); font-size: 12px;">No history yet.</p>
      </div>
    </div>
  </div>
</template>
