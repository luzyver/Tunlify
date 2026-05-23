<script setup lang="ts">
import {
  Activity,
  Archive,
  BarChart3,
  Bell,
  Cog,
  Container,
  LayoutDashboard,
  LogOut,
  ScrollText,
  Search,
  Settings,
  Shield,
  Terminal,
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'

const emit = defineEmits<{
  navigate: []
  'open-cmd': []
}>()

const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()

const nav = [
  { to: '/', label: 'Status', icon: LayoutDashboard },
  { to: '/projects', label: 'Projects', icon: Container },
  { to: '/health', label: 'Health', icon: Activity },
  { to: '/metrics', label: 'Metrics', icon: BarChart3 },
  { to: '/logs', label: 'Logs', icon: ScrollText },
  { to: '/config', label: 'Ingress', icon: Settings },
  { to: '/backups', label: 'Backups', icon: Archive },
  { to: '/tcp-access', label: 'TCP Access', icon: Terminal },
  { to: '/notifications', label: 'Alerts', icon: Bell },
  { to: '/audit', label: 'Audit', icon: Shield },
  { to: '/settings', label: 'Settings', icon: Cog },
]

async function handleLogout() {
  try { await apiFetch('/auth/logout', { method: 'POST' }) } catch {}
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <aside
    class="w-sidebar shrink-0 h-screen sticky top-0
           bg-bg border-r border-border
           flex flex-col"
  >

    <div class="h-14 flex items-center gap-2.5 px-4 border-b border-border shrink-0">
      <img src="/icon.png" alt="" class="w-6 h-6 rounded-sm" />
      <span class="text-sm font-semibold tracking-tight text-text">Tunlify</span>
    </div>

    <button
      type="button"
      class="mx-3 mt-3 mb-2 flex items-center gap-2.5 h-9 px-3
             bg-surface border border-border rounded-md
             text-text-dim text-sm
             transition-colors duration-100
             hover:border-border-strong hover:text-text-muted"
      @click="$emit('open-cmd')"
    >
      <Search class="w-4 h-4 shrink-0" :stroke-width="1.75" />
      <span class="flex-1 text-left">Search…</span>
      <span class="kbd">⌘K</span>
    </button>

    <nav class="flex-1 px-2 pb-3 space-y-px overflow-y-auto scrollbar-thin">
      <router-link
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        @click="emit('navigate')"
        class="flex items-center gap-2.5 h-8 px-3
               rounded-md text-sm text-text-muted
               transition-colors duration-100
               hover:bg-bg-alt hover:text-text"
        active-class="!text-accent !bg-accent-soft font-medium"
      >
        <component :is="item.icon" class="w-4 h-4 shrink-0" :stroke-width="1.75" />
        <span class="truncate">{{ item.label }}</span>
      </router-link>
    </nav>

    <div class="px-3 py-3 border-t border-border flex items-center gap-2 shrink-0">
      <div
        class="w-7 h-7 rounded-full bg-accent-soft text-accent
               flex items-center justify-center
               text-xs font-semibold uppercase shrink-0"
      >
        {{ authStore.username?.[0] || 'A' }}
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-sm font-medium text-text truncate leading-tight">
          {{ authStore.username || 'admin' }}
        </p>
        <p class="text-2xs text-text-dim truncate leading-tight mt-0.5">Signed in</p>
      </div>
      <button @click="handleLogout" class="btn-icon-ghost" title="Sign out">
        <LogOut class="w-3.5 h-3.5" :stroke-width="1.75" />
      </button>
    </div>
  </aside>
</template>
