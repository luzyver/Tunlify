<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu, X } from 'lucide-vue-next'
import AppSidebar from './components/AppSidebar.vue'

const route = useRoute()
const showSidebar = computed(() => route.path !== '/login')
const mobileOpen = ref(false)
</script>

<template>
  <div class="flex h-screen">
    <div v-if="showSidebar" class="fixed top-0 left-0 right-0 z-40 bg-surface border-b border-border px-4 py-3 flex items-center justify-between lg:hidden">
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="Tunlify" class="w-5 h-5" />
        <span class="font-display text-sm text-neon text-glow tracking-wider">TUNLIFY</span>
      </div>
      <button @click="mobileOpen = !mobileOpen" class="text-neon">
        <Menu v-if="!mobileOpen" class="w-5 h-5" />
        <X v-else class="w-5 h-5" />
      </button>
    </div>

    <div v-if="showSidebar && mobileOpen" class="fixed inset-0 z-30 bg-black/70 lg:hidden" @click="mobileOpen = false"></div>

    <div v-if="showSidebar" :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'" class="fixed z-30 inset-y-0 left-0 transition-transform lg:translate-x-0 lg:static">
      <AppSidebar @navigate="mobileOpen = false" />
    </div>

    <main class="flex-1 overflow-auto p-4 pt-16 lg:p-6 lg:pt-6">
      <router-view />
    </main>
    <div class="scanlines"></div>
  </div>
</template>
