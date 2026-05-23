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
  <div class="flex h-screen max-w-app mx-auto">
    <!-- Mobile header -->
    <div v-if="showSidebar" class="fixed top-0 left-0 right-0 z-40 bg-bg border-b border-border px-4 py-3 flex items-center justify-between md:hidden">
      <div class="flex items-center gap-2">
        <img src="/icon.png" alt="Tunlify" class="w-6 h-6" />
        <span class="text-sm font-bold text-text">Tunlify</span>
      </div>
      <button @click="mobileOpen = !mobileOpen" class="btn-icon">
        <Menu v-if="!mobileOpen" class="w-4 h-4" />
        <X v-else class="w-4 h-4" />
      </button>
    </div>

    <!-- Mobile overlay -->
    <div v-if="showSidebar && mobileOpen" class="fixed inset-0 z-30 bg-black/20 md:hidden" @click="mobileOpen = false"></div>

    <!-- Sidebar -->
    <div v-if="showSidebar" :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'" class="fixed z-30 inset-y-0 left-0 transition-transform md:translate-x-0 md:static">
      <AppSidebar @navigate="mobileOpen = false" />
    </div>

    <!-- Main -->
    <main class="flex-1 overflow-auto pt-14 md:pt-0">
      <div class="p-6">
        <router-view />
      </div>
    </main>
  </div>
</template>
