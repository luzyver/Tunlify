<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  Activity,
  Archive,
  BarChart3,
  Bell,
  Cog,
  Container,
  LayoutDashboard,
  LogOut,
  RotateCcw,
  ScrollText,
  Search,
  Settings,
  Shield,
  Terminal,
} from 'lucide-vue-next'
import { useApi } from '../composables/useApi'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const { apiFetch } = useApi()

const open = ref(false)
const query = ref('')
const cursor = ref(0)
const inputEl = ref<HTMLInputElement | null>(null)

interface CmdItem {
  id: string
  label: string
  icon: any
  group: 'navigate' | 'action'
  hint?: string
  to?: string
  run?: () => void | Promise<void>
}

const items: CmdItem[] = [
  { id: 'nav-status', label: 'Go to Status', icon: LayoutDashboard, group: 'navigate', hint: 'G S', to: '/' },
  { id: 'nav-projects', label: 'Go to Projects', icon: Container, group: 'navigate', hint: 'G P', to: '/projects' },
  { id: 'nav-health', label: 'Go to Health', icon: Activity, group: 'navigate', hint: 'G H', to: '/health' },
  { id: 'nav-metrics', label: 'Go to Metrics', icon: BarChart3, group: 'navigate', hint: 'G M', to: '/metrics' },
  { id: 'nav-logs', label: 'Go to Logs', icon: ScrollText, group: 'navigate', hint: 'G L', to: '/logs' },
  { id: 'nav-config', label: 'Go to Ingress', icon: Settings, group: 'navigate', hint: 'G I', to: '/config' },
  { id: 'nav-backups', label: 'Go to Backups', icon: Archive, group: 'navigate', hint: 'G B', to: '/backups' },
  { id: 'nav-tcp', label: 'Go to TCP Access', icon: Terminal, group: 'navigate', hint: 'G T', to: '/tcp-access' },
  { id: 'nav-alerts', label: 'Go to Alerts', icon: Bell, group: 'navigate', hint: 'G A', to: '/notifications' },
  { id: 'nav-audit', label: 'Go to Audit', icon: Shield, group: 'navigate', hint: 'G U', to: '/audit' },
  { id: 'nav-settings', label: 'Go to Settings', icon: Cog, group: 'navigate', hint: 'G ,', to: '/settings' },
  {
    id: 'act-restart',
    label: 'Restart cloudflared tunnel',
    icon: RotateCcw,
    group: 'action',
    hint: 'R',
    run: async () => {
      try {
        await apiFetch('/api/control/restart', { method: 'POST' })
      } catch {
      }
    },
  },
  {
    id: 'act-logout',
    label: 'Sign out',
    icon: LogOut,
    group: 'action',
    hint: 'Q',
    run: async () => {
      try {
        await apiFetch('/auth/logout', { method: 'POST' })
      } catch {}
      authStore.logout()
      router.push('/login')
    },
  },
]

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return items
  return items.filter((i) => i.label.toLowerCase().includes(q))
})

const groups = computed(() => {
  const navigate = filtered.value.filter((i) => i.group === 'navigate')
  const action = filtered.value.filter((i) => i.group === 'action')
  return { navigate, action }
})

function show() {
  open.value = true
  query.value = ''
  cursor.value = 0
  nextTick(() => inputEl.value?.focus())
}
function hide() {
  open.value = false
}
function toggle() {
  open.value ? hide() : show()
}

async function runItem(item: CmdItem) {
  hide()
  if (item.to) router.push(item.to)
  else if (item.run) await item.run()
}

function onKey(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    toggle()
    return
  }
  if (!open.value) return
  if (e.key === 'Escape') {
    e.preventDefault()
    hide()
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    cursor.value = Math.min(cursor.value + 1, filtered.value.length - 1)
    return
  }
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    cursor.value = Math.max(cursor.value - 1, 0)
    return
  }
  if (e.key === 'Enter') {
    e.preventDefault()
    const item = filtered.value[cursor.value]
    if (item) runItem(item)
  }
}

watch(query, () => {
  cursor.value = 0
})

onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))

defineExpose({ open: show })
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="cmdk-overlay" @click.self="hide" role="dialog" aria-modal="true" aria-label="Command palette">
      <div class="cmdk-panel">
        <div class="cmdk-input-row">
          <Search class="w-4 h-4 text-text-muted shrink-0" :stroke-width="1.75" />
          <input
            ref="inputEl"
            v-model="query"
            class="cmdk-input"
            placeholder="Type a command or search…"
            spellcheck="false"
            autocomplete="off"
          />
          <span class="kbd">ESC</span>
        </div>

        <div class="cmdk-list scrollbar-thin">
          <template v-if="groups.navigate.length">
            <div class="cmdk-section">Navigate</div>
            <div
              v-for="(item, i) in groups.navigate"
              :key="item.id"
              class="cmdk-row"
              :class="{ 'is-active': cursor === filtered.indexOf(item) }"
              @mouseenter="cursor = filtered.indexOf(item)"
              @click="runItem(item)"
            >
              <component :is="item.icon" class="cmdk-row-icon" :stroke-width="1.75" />
              <span class="cmdk-row-label">{{ item.label }}</span>
              <span v-if="item.hint" class="cmdk-row-hint">{{ item.hint }}</span>
            </div>
</template>

          <template v-if="groups.action.length">
            <div class="cmdk-section">Actions</div>
            <div
              v-for="item in groups.action"
              :key="item.id"
              class="cmdk-row"
              :class="{ 'is-active': cursor === filtered.indexOf(item) }"
              @mouseenter="cursor = filtered.indexOf(item)"
              @click="runItem(item)"
            >
              <component :is="item.icon" class="cmdk-row-icon" :stroke-width="1.75" />
              <span class="cmdk-row-label">{{ item.label }}</span>
              <span v-if="item.hint" class="cmdk-row-hint">{{ item.hint }}</span>
            </div>
</template>

          <div v-if="!filtered.length" class="cmdk-empty">No results for "{{ query }}"</div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
