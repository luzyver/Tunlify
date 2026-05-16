<script setup lang="ts">
import { LayoutDashboard, ScrollText, Settings, Terminal, Shield, LogOut, Activity, BarChart3, Archive, Bell, Cog } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'

const emit = defineEmits<{ navigate: [] }>()
const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()

const nav = [
  { to: '/', label: 'Status', icon: LayoutDashboard },
  { to: '/health', label: 'Health', icon: Activity },
  { to: '/metrics', label: 'Metrics', icon: BarChart3 },
  { to: '/logs', label: 'Logs', icon: ScrollText },
  { to: '/config', label: 'Ingress', icon: Settings },
  { to: '/backups', label: 'Backups', icon: Archive },
  { to: '/tcp-access', label: 'TCP', icon: Terminal },
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
  <aside class="w-56 h-full bg-surface border-r border-border flex flex-col">
    <div class="p-5 border-b border-border">
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="Tunlify" class="w-7 h-7" />
        <h1 class="font-display text-lg text-neon text-glow tracking-widest">TUNLIFY</h1>
      </div>
      <p class="text-xs text-muted mt-1 tracking-wider">// TUNNEL CTRL</p>
    </div>

    <nav class="flex-1 p-3 space-y-1">
      <router-link
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        @click="emit('navigate')"
        class="flex items-center gap-3 px-3 py-2.5 text-sm text-muted uppercase tracking-wider hover:text-neon hover:bg-neon/5 transition-all duration-150"
        active-class="!text-neon bg-neon/5 border-l-2 border-neon shadow-neon"
      >
        <component :is="item.icon" class="w-4 h-4" :stroke-width="1.5" />
        {{ item.label }}
      </router-link>
    </nav>

    <div class="p-4 border-t border-border">
      <div class="flex items-center justify-between">
        <span class="text-xs text-muted tracking-wider">{{ authStore.username }}<span class="animate-blink text-neon">_</span></span>
        <button @click="handleLogout" class="text-muted hover:text-danger transition-colors">
          <LogOut class="w-4 h-4" :stroke-width="1.5" />
        </button>
      </div>
    </div>
  </aside>
</template>
