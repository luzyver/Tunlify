<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { useApi } from '../composables/useApi'
import DataTable, { type Column } from '../components/DataTable.vue'

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
const historyName = ref('')

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
  const confirmed = await confirm('Delete this project? This cannot be undone.')
  if (!confirmed) return
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

async function toggleHistory(p: Project) {
  if (historyTarget.value === p.id) { historyTarget.value = null; return }
  historyTarget.value = p.id
  historyName.value = p.name
  try { history.value = await apiFetch(`/api/projects/${p.id}/history`) }
  catch { history.value = [] }
}

const anyRunning = () => Object.values(actionLoading.value).some(Boolean)

load()

const columns: Column<Project>[] = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'path', label: 'Path', sortable: true },
  { key: 'repo_url', label: 'Repository', sortable: true, hideBelow: 'lg' },
  { key: 'created_at', label: 'Created', sortable: true, align: 'right', width: '140px', cellClass: 'num text-caption text-medium-emphasis', headerClass: 'num' },
  { key: 'actions', label: '', align: 'right', width: '260px' },
]
</script>

<template>
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Compose</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Projects</h1>
      </div>
      <v-btn color="primary" @click="openAdd">
        <v-icon start>mdi-plus</v-icon>
        New project
      </v-btn>
    </div>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <v-card v-if="showForm" class="mb-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>{{ editingId ? 'Edit project' : 'New project' }}</span>
        <v-btn icon="mdi-close" variant="text" size="small" @click="showForm = false" />
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field v-model="form.name" label="Name" placeholder="my-app" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field v-model="form.path" label="Path" placeholder="/srv/my-app" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="form.repo_url" label="Repository URL (optional)" placeholder="https://github.com/user/repo.git" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field v-model="form.git_username" label="Git username" placeholder="git" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field v-model="form.git_token" label="Git token" placeholder="ghp_..." type="password" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12">
            <div class="d-flex ga-3 pt-1">
              <v-btn color="primary" @click="save">{{ editingId ? 'Save changes' : 'Create project' }}</v-btn>
              <v-btn variant="tonal" @click="showForm = false">Cancel</v-btn>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-card v-if="deployTarget" class="mb-6">
      <v-card-title class="d-flex align-center ga-2">
        <v-icon color="primary">mdi-rocket-launch</v-icon>
        Deploy &mdash; choose ref
        <v-spacer />
        <v-btn icon="mdi-close" variant="text" size="small" @click="deployTarget = null" />
      </v-card-title>
      <v-card-text class="d-flex ga-3">
        <v-text-field
          v-model="deployRef"
          placeholder="main / v1.0.0"
          variant="outlined"
          density="compact"
          hide-details
          @keyup.enter="confirmDeploy"
        />
        <v-btn
          color="primary"
          :disabled="!deployRef || actionLoading[deployTarget]"
          @click="confirmDeploy"
        >
          Deploy
        </v-btn>
      </v-card-text>
    </v-card>

    <v-card v-if="output || anyRunning()" class="mb-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Output</span>
        <v-chip v-if="anyRunning()" color="primary" size="x-small" variant="tonal">
          <template #prepend>
            <v-icon size="x-small">mdi-checkbox-blank-circle</v-icon>
          </template>
          Running
        </v-chip>
      </v-card-title>
      <pre
        ref="outputEl"
        class="font-mono text-caption text-medium-emphasis bg-grey-lighten-3 pa-4 overflow-auto scrollbar-thin"
        style="max-height: 384px; white-space: pre-wrap; line-height: 1.25;"
      >{{ output || '(waiting)' }}</pre>
    </v-card>

    <DataTable
      :data="projects"
      :columns="columns"
      :searchable="true"
      search-placeholder="Search projects..."
      :page-size="25"
      :row-class="(row) => historyTarget === row.id ? 'is-selected' : undefined"
    >
      <template #cell-name="{ row }">
        <span class="font-weight-medium text-on-surface">{{ row.name }}</span>
      </template>
      <template #cell-path="{ row }">
        <span class="font-mono text-caption text-medium-emphasis text-truncate" style="max-width: 260px;">{{ row.path }}</span>
      </template>
      <template #cell-repo_url="{ row }">
        <span class="font-mono text-caption text-disabled text-truncate" style="max-width: 260px;">{{ row.repo_url || '\u2014' }}</span>
      </template>
      <template #cell-created_at="{ row }">
        {{ row.created_at?.split('T')[0] || '\u2014' }}
      </template>
      <template #cell-actions="{ row }">
        <div class="d-inline-flex align-center ga-1">
          <v-btn icon="mdi-play" variant="text" size="small" :disabled="actionLoading[row.id]" @click="action(row.id, 'up')" />
          <v-btn icon="mdi-stop" variant="text" size="small" :disabled="actionLoading[row.id]" @click="action(row.id, 'down')" />
          <v-btn icon="mdi-restart" variant="text" size="small" :disabled="actionLoading[row.id]" @click="action(row.id, 'restart')" />
          <v-btn icon="mdi-rocket-launch" color="primary" variant="text" size="small" :disabled="actionLoading[row.id]" @click="startDeploy(row)" />
          <v-btn icon="mdi-history" variant="text" size="small" @click="toggleHistory(row)" />
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEdit(row)" />
          <v-btn icon="mdi-delete" color="error" variant="text" size="small" @click="remove(row.id)" />
        </div>
      </template>
      <template #empty>
        {{ loading ? 'Loading projects...' : 'No projects yet \u2014 click "New project" above.' }}
      </template>
    </DataTable>

    <v-card v-if="historyTarget !== null" class="mt-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>History &middot; <span class="font-mono text-medium-emphasis">{{ historyName }}</span></span>
        <v-btn icon="mdi-close" variant="text" size="small" @click="historyTarget = null" />
      </v-card-title>
      <v-card-text class="d-flex flex-column ga-1">
        <div v-if="history.length" class="d-flex flex-column ga-1">
          <div v-for="h in history" :key="h.id">
            <v-btn
              variant="text"
              block
              class="d-flex align-center justify-space-between text-medium-emphasis"
              :disabled="h.action !== 'project_deploy'"
              @click="h._open = !h._open"
            >
              <span class="font-mono">{{ h.action }}</span>
              <span class="tabular-nums text-disabled">{{ h.created_at }}</span>
            </v-btn>
            <pre
              v-if="h._open"
              class="font-mono text-caption text-medium-emphasis ml-3 my-1"
              style="border-left: 2px solid #ded9ca; padding-left: 12px; white-space: pre-wrap;"
            >{{ h.detail }}</pre>
          </div>
        </div>
        <p v-else class="text-caption text-disabled">No history yet.</p>
      </v-card-text>
    </v-card>
  </div>
</template>
