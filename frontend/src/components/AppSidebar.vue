<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'

const emit = defineEmits<{ navigate: [] }>()

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const { apiFetch } = useApi()

const nav = [
  { to: '/', label: 'Status', icon: 'mdi-view-dashboard' },
  { to: '/projects', label: 'Projects', icon: 'mdi-docker' },
  { to: '/health', label: 'Health', icon: 'mdi-heart-pulse' },
  { to: '/metrics', label: 'Metrics', icon: 'mdi-chart-bar' },
  { to: '/config', label: 'Ingress', icon: 'mdi-tune' },
  { to: '/backups', label: 'Backups', icon: 'mdi-archive' },
  { to: '/tcp-access', label: 'TCP Access', icon: 'mdi-console' },
  { to: '/notifications', label: 'Alerts', icon: 'mdi-bell' },
  { to: '/audit', label: 'Audit', icon: 'mdi-shield' },
  { to: '/settings', label: 'Settings', icon: 'mdi-cog' },
]

async function handleLogout() {
  try { await apiFetch('/auth/logout', { method: 'POST' }) } catch {}
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <v-navigation-drawer permanent :width="240" color="background">
    <template #prepend>
      <div class="d-flex align-center ga-3 px-4 py-3 border-bottom" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <div class="rounded-lg d-flex align-center justify-center" style="width: 32px; height: 32px; background: rgba(247, 147, 26, 0.15); border: 1px solid rgba(247, 147, 26, 0.3);">
          <img src="/icon.png?v=2" alt="" style="width: 20px; height: 20px;" />
        </div>
        <div>
          <div class="font-heading font-semibold text-body-1" style="color: white;">Tunlify</div>
          <div class="eyebrow-dim" style="font-size: 10px;">Console</div>
        </div>
      </div>
    </template>

    <v-list density="compact" nav class="px-2 pt-3">
      <v-list-item
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        :active="route.path === item.to || (item.to !== '/' && route.path.startsWith(item.to))"
        color="primary"
        variant="text"
        density="compact"
        class="mb-1 rounded-lg"
        :style="{
          background: (route.path === item.to || (item.to !== '/' && route.path.startsWith(item.to))) ? 'rgba(247, 147, 26, 0.1)' : 'transparent',
          borderLeft: (route.path === item.to || (item.to !== '/' && route.path.startsWith(item.to))) ? '2px solid #F7931A' : '2px solid transparent',
        }"
        @click="emit('navigate')"
      >
        <template #prepend>
          <v-icon size="18">{{ item.icon }}</v-icon>
        </template>
        <template #title>
          <span class="font-mono" style="font-size: 13px;">{{ item.label }}</span>
        </template>
      </v-list-item>
    </v-list>

    <template #append>
      <div class="border-top pa-4" style="border-top: 1px solid rgba(30, 41, 59, 0.6);">
        <div class="d-flex align-center ga-3">
          <div
            class="rounded-full d-flex align-center justify-center font-mono font-weight-bold"
            style="width: 32px; height: 32px; background: rgba(247, 147, 26, 0.15); border: 1px solid rgba(247, 147, 26, 0.3); color: #F7931A; font-size: 13px;"
          >
            {{ authStore.username?.[0]?.toUpperCase() || 'A' }}
          </div>
          <div class="flex-grow-1" style="min-width: 0;">
            <div class="font-mono text-body-2" style="color: white; line-height: 1.2;">{{ authStore.username || 'admin' }}</div>
            <div class="eyebrow-dim" style="font-size: 10px;">Signed in</div>
          </div>
          <v-btn icon variant="text" size="small" @click="handleLogout" style="color: #94A3B8;">
            <v-icon size="16">mdi-logout</v-icon>
          </v-btn>
        </div>
      </div>
    </template>
  </v-navigation-drawer>
</template>
