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
  <div class="flex h-screen bg-bg">
    <!-- Mobile header -->
    <div v-if="showSidebar" class="fixed top-0 left-0 right-0 z-40 bg-bg border-b border-border px-4 py-3 flex items-center justify-between lg:hidden">
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="Tunlify" class="w-5 h-5" />
        <span class="text-sm font-medium text-text">Tunlify</span>
      </div>
      <button @click="mobileOpen = !mobileOpen" class="text-text-muted hover:text-text">
        <Menu v-if="!mobileOpen" class="w-5 h-5" />
        <X v-else class="w-5 h-5" />
      </button>
    </div>

    <!-- Mobile overlay -->
    <div v-if="showSidebar && mobileOpen" class="fixed inset-0 z-30 bg-black/20 lg:hidden" @click="mobileOpen = false"></div>

    <!-- Sidebar -->
    <div v-if="showSidebar" :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'" class="fixed z-30 inset-y-0 left-0 transition-transform lg:translate-x-0 lg:static">
      <AppSidebar @navigate="mobileOpen = false" />
    </div>

    <!-- Main content -->
    <main class="flex-1 overflow-auto pt-14 lg:pt-0">
      <div class="max-w-content mx-auto px-6 py-8">
        <router-view />
      </div>
    </main>
  </div>
</template>
