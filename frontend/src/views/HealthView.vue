<script setup lang="ts">
import { ref } from 'vue'
import { RefreshCw } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const results = ref<any[]>([])
const loading = ref(false)

async function check() {
  loading.value = true
  try { results.value = await apiFetch('/api/health') } catch {}
  finally { loading.value = false }
}
check()
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="cyber-heading text-xl lg:text-2xl">Health Check</h1>
      <button @click="check" :disabled="loading" class="cyber-btn-secondary flex items-center gap-2">
        <RefreshCw class="w-4 h-4" :stroke-width="1.5" :class="loading && 'animate-spin'" /> Refresh
      </button>
    </div>

    <div class="space-y-2">
      <div v-for="r in results" :key="r.hostname" class="cyber-card flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-200 font-mono">{{ r.hostname }}</p>
          <p class="text-xs text-muted mt-0.5">{{ r.service }}</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs text-muted font-mono">{{ r.latency }}</span>
          <span class="cyber-badge" :class="r.status === 'up' ? 'border-neon text-neon' : 'border-danger text-danger'">
            {{ r.status }}
          </span>
        </div>
      </div>
      <div v-if="!results.length && !loading" class="cyber-card text-center text-muted py-8 tracking-wider text-sm">
        // NO SERVICES
      </div>
    </div>
  </div>
</template>
