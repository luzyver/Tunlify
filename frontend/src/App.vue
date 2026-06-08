<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from './components/AppSidebar.vue'

const route = useRoute()
const isAppMode = computed(() => route.path !== '/login')
const rail = ref(false)
</script>

<template>
  <v-app v-if="isAppMode">
    <AppSidebar :rail="rail" @update:rail="rail = $event" />
    <v-main>
      <router-view />
    </v-main>
  </v-app>

  <router-view v-else />
</template>
