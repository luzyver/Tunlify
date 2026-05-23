<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Menu, X } from 'lucide-vue-next'
import AppSidebar from './components/AppSidebar.vue'

const route = useRoute()
const isAppMode = computed(() => route.path !== '/login')
const mobileOpen = ref(false)
</script>

<template>
  <div v-if="!isAppMode" class="min-h-screen bg-bg text-text">
    <router-view />
  </div>

  <div v-else class="app-shell">
    <div
      class="md:hidden fixed inset-x-0 top-0 z-40 h-14
             bg-bg/95 backdrop-blur-sm border-b border-border
             flex items-center justify-between px-4"
    >
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="" class="w-6 h-6" />
        <span class="text-sm font-semibold text-text">Tunlify</span>
      </div>
      <button class="btn-icon-ghost" @click="mobileOpen = !mobileOpen" :title="mobileOpen ? 'Close menu' : 'Open menu'">
        <Menu v-if="!mobileOpen" class="w-4 h-4" :stroke-width="1.75" />
        <X v-else class="w-4 h-4" :stroke-width="1.75" />
      </button>
    </div>

    <div
      v-if="mobileOpen"
      class="md:hidden fixed inset-0 z-30 bg-text/20"
      @click="mobileOpen = false"
    ></div>

    <div
      class="fixed md:sticky inset-y-0 left-0 z-40 transition-transform duration-150"
      :class="mobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
    >
      <AppSidebar @navigate="mobileOpen = false" />
    </div>

    <main class="app-content">
      <div class="app-content-inner pt-20 md:pt-6">
        <router-view />
      </div>
    </main>
  </div>
</template>
