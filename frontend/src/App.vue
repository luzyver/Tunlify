<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import AppSidebar from './components/AppSidebar.vue'
import BottomDock from './components/BottomDock.vue'

const route = useRoute()
const authStore = useAuthStore()
const isAppMode = computed(() => authStore.isAuthenticated && route.path !== '/login')
</script>

<template>
  <v-app v-if="isAppMode">
    <div class="bg-orange-glow pointer-events-none" style="position: fixed; inset: 0; z-index: 0;" />
    <div class="bg-blue-glow pointer-events-none" style="position: fixed; inset: 0; z-index: 0;" />
    <div class="bg-grid pointer-events-none" style="position: fixed; inset: 0; z-index: 0;" />
    <AppSidebar />
    <v-main style="position: relative; z-index: 1; overflow: hidden;">
      <div style="height: 100%; overflow: auto;">
        <router-view />
      </div>
    </v-main>
    <BottomDock />
  </v-app>

  <router-view v-else />
</template>
