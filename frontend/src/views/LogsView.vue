<script setup lang="ts">
import { ref, computed, nextTick, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const logs = ref<{ time: string; message: string }[]>([])
const filter = ref('all')
const search = ref('')
const autoScroll = ref(true)
const logContainer = ref<HTMLElement>()
const wsStatus = ref('connecting')

const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/api/logs/ws?token=${authStore.token}`
let ws: WebSocket | null = null

function connect() {
  ws = new WebSocket(wsUrl)
  ws.onopen = () => { wsStatus.value = 'connected' }
  ws.onclose = () => { wsStatus.value = 'disconnected'; setTimeout(connect, 3000) }
  ws.onmessage = (e) => {
    try {
      const entry = JSON.parse(e.data)
      logs.value.push(entry)
      if (logs.value.length > 2000) logs.value.shift()
      if (autoScroll.value) nextTick(() => logContainer.value?.scrollTo(0, logContainer.value.scrollHeight))
    } catch {}
  }
}
connect()
onUnmounted(() => ws?.close())

const filteredLogs = computed(() => {
  let result = logs.value
  if (filter.value !== 'all') result = result.filter(l => l.message.toLowerCase().includes(filter.value))
  if (search.value) { const s = search.value.toLowerCase(); result = result.filter(l => l.message.toLowerCase().includes(s)) }
  return result
})

function getLogClass(msg: string) {
  if (msg.includes('ERR') || msg.includes('error')) return 'text-red-600'
  if (msg.includes('WRN') || msg.includes('warn')) return 'text-amber-600'
  return 'text-text'
}
</script>

<template>
  <div class="space-y-6 h-full flex flex-col">
    <header class="border-b border-border pb-6">
      <div class="flex items-center justify-between">
        <div>
          <p class="section-marker mb-2">LOGS</p>
          <h1 class="editorial-h1 !text-3xl">Log Stream</h1>
        </div>
        <span class="badge" :class="wsStatus === 'connected' ? 'badge-success' : 'badge-danger'">
          {{ wsStatus }}
        </span>
      </div>
    </header>

    <div class="flex flex-col sm:flex-row gap-3 sm:items-center">
      <div class="flex border border-border rounded-md overflow-hidden">
        <button v-for="f in ['all', 'error', 'warn', 'info']" :key="f" @click="filter = f"
          :class="filter === f ? 'bg-surface text-text font-medium' : 'text-text-muted hover:text-text'"
          class="px-3 py-1.5 text-xs transition-colors">{{ f }}</button>
      </div>
      <div class="flex gap-3 sm:ml-auto items-center">
        <input v-model="search" placeholder="Search logs…" class="input !py-1.5 !text-xs sm:!w-48" />
        <label class="flex items-center gap-2 text-xs text-text-muted whitespace-nowrap">
          <input type="checkbox" v-model="autoScroll" class="accent-accent" /> Auto-scroll
        </label>
      </div>
    </div>

    <div ref="logContainer" class="flex-1 card !p-0 overflow-auto min-h-[400px]">
      <div v-for="(log, i) in filteredLogs" :key="i" class="px-4 py-1.5 border-b border-border/50 hover:bg-surface text-xs font-mono">
        <span class="text-text-dim mr-3">{{ log.time?.split('T')[1]?.slice(0, 8) }}</span>
        <span :class="getLogClass(log.message)">{{ log.message }}</span>
      </div>
      <div v-if="!filteredLogs.length" class="text-text-muted text-center py-12 text-sm">
        Waiting for log data…
      </div>
    </div>
  </div>
</template>
