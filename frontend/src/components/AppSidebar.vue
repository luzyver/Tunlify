<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'

const emit = defineEmits<{ navigate: [] }>()
const props = withDefaults(defineProps<{ rail: boolean }>(), { rail: false })

const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()

const nav = [
  { to: '/', label: 'Status', icon: 'mdi-view-dashboard' },
  { to: '/projects', label: 'Projects', icon: 'mdi-docker' },
  { to: '/health', label: 'Health', icon: 'mdi-heart-pulse' },
  { to: '/metrics', label: 'Metrics', icon: 'mdi-chart-bar' },
  { to: '/logs', label: 'Logs', icon: 'mdi-text-box-search-outline' },
  { to: '/config', label: 'Ingress', icon: 'mdi-tune' },
  { to: '/backups', label: 'Backups', icon: 'mdi-archive' },
  { to: '/tcp-access', label: 'TCP Access', icon: 'mdi-console' },
  { to: '/notifications', label: 'Alerts', icon: 'mdi-bell-outline' },
  { to: '/audit', label: 'Audit', icon: 'mdi-shield-account' },
  { to: '/settings', label: 'Settings', icon: 'mdi-cog' },
]

async function handleLogout() {
  try { await apiFetch('/auth/logout', { method: 'POST' }) } catch {}
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <v-navigation-drawer
    :rail="props.rail"
    permanent
    :width="240"
    :rail-width="64"
    color="background"
    @update:rail="$emit('update:rail', $event)"
  >
    <template #prepend>
      <v-list-item
        class="px-4 py-2"
        :prepend-avatar="'/icon.png'"
        title="Tunlify"
        subtitle="Console"
        nav
      />
      <v-divider />
    </template>

    <v-list density="compact" nav>
      <v-list-item
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        :prepend-icon="item.icon"
        :title="item.label"
        :active="router.currentRoute.value.path === item.to"
        color="primary"
        variant="text"
        density="compact"
        class="mb-1"
        @click="emit('navigate')"
      />
    </v-list>

    <template #append>
      <v-divider />
      <v-list-item
        :prepend-avatar="`https://ui-avatars.com/api/?name=${authStore.username?.[0] || 'A'}&background=e2e6fb&color=5266eb&size=28`"
        :title="authStore.username || 'admin'"
        subtitle="Signed in"
        class="pt-2"
      >
        <template #append>
          <v-btn icon="mdi-logout" variant="text" size="small" @click="handleLogout" />
        </template>
      </v-list-item>
    </template>
  </v-navigation-drawer>
</template>
