<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from './components/AppSidebar.vue'

const route = useRoute()
const isAppMode = computed(() => route.path !== '/login')
</script>

<template>
  <v-app v-if="isAppMode">
    <div class="bg-orange-glow pointer-events-none" style="position: fixed; inset: 0; z-index: 0;" />
    <div class="bg-grid pointer-events-none" style="position: fixed; inset: 0; z-index: 0;" />
    <AppSidebar />
    <v-main style="position: relative; z-index: 1;">
      <router-view />
    </v-main>
  </v-app>

  <router-view v-else />
</template>
