<script setup lang="ts">
import { LayoutDashboard, ScrollText, Settings, Terminal, Shield, LogOut, Activity, BarChart3, Archive, Bell, Cog, Container } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'

const emit = defineEmits<{ navigate: [] }>()
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
  <aside class="w-[256px] h-full bg-bg border-r border-border flex flex-col">
    <div class="px-6 py-6 border-b border-border">
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="Tunlify" class="w-6 h-6" />
        <h1 class="text-base font-semibold text-text">Tunlify</h1>
      </div>
      <p class="text-xs text-text-dim mt-1">Tunnel Management</p>
    </div>

    <nav class="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
      <router-link
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        @click="emit('navigate')"
        class="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-text-muted hover:text-text hover:bg-surface transition-colors"
        active-class="!text-text !bg-surface font-medium"
      >
        <component :is="item.icon" class="w-4 h-4" :stroke-width="1.5" />
        {{ item.label }}
      </router-link>
    </nav>

    <div class="px-6 py-4 border-t border-border">
      <div class="flex items-center justify-between">
        <span class="text-sm text-text-muted">{{ authStore.username }}</span>
        <button @click="handleLogout" class="text-text-dim hover:text-text transition-colors" title="Logout">
          <LogOut class="w-4 h-4" :stroke-width="1.5" />
        </button>
      </div>
    </div>
  </aside>
</template>
