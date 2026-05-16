<script setup lang="ts">
import { ref, computed, nextTick, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const logs = ref<{ time: string; message: string }[]>([])
const filter = ref('all')
const search = ref('')
const autoScroll = ref(true)
const logContainer = ref<HTMLElement>()
const wsStatus = ref('CONNECTING')

const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/api/logs/ws?token=${authStore.token}`
let ws: WebSocket | null = null

function connect() {
  ws = new WebSocket(wsUrl)
  ws.onopen = () => { wsStatus.value = 'OPEN' }
  ws.onclose = () => { wsStatus.value = 'CLOSED'; setTimeout(connect, 3000) }
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
  if (msg.includes('ERR') || msg.includes('error')) return 'text-danger'
  if (msg.includes('WRN') || msg.includes('warn')) return 'text-magenta'
  return 'text-gray-300'
}
</script>

<template>
  <div class="space-y-4 h-full flex flex-col">
    <div class="flex items-center justify-between">
      <h1 class="cyber-heading text-xl lg:text-2xl">Log Stream</h1>
      <span class="cyber-badge" :class="wsStatus === 'OPEN' ? 'border-neon text-neon' : 'border-danger text-danger'">
        {{ wsStatus }}
      </span>
    </div>

    <div class="flex flex-col sm:flex-row gap-2 sm:items-center">
      <div class="flex border border-border overflow-hidden">
        <button v-for="f in ['all', 'error', 'warn', 'info']" :key="f" @click="filter = f"
          :class="filter === f ? 'bg-neon/10 text-neon' : 'text-muted hover:text-gray-200'"
          class="px-3 py-1.5 text-xs uppercase tracking-wider transition-colors">{{ f }}</button>
      </div>
      <div class="flex gap-2 sm:ml-auto items-center">
        <input v-model="search" placeholder="grep..." class="cyber-input !py-1.5 flex-1 sm:!w-48 !text-xs" />
        <label class="flex items-center gap-2 text-xs text-muted whitespace-nowrap">
          <input type="checkbox" v-model="autoScroll" class="accent-neon" /> tail
        </label>
      </div>
    </div>

    <div ref="logContainer" class="flex-1 cyber-card !p-0 overflow-auto min-h-[400px]">
      <div v-for="(log, i) in filteredLogs" :key="i" class="px-4 py-1 border-b border-border/30 hover:bg-neon/5 text-xs font-mono">
        <span class="text-muted mr-2">{{ log.time?.split('T')[1]?.slice(0, 8) }}</span>
        <span :class="getLogClass(log.message)">{{ log.message }}</span>
      </div>
      <div v-if="!filteredLogs.length" class="text-muted text-center py-12 text-sm tracking-wider">
        // AWAITING DATA<span class="animate-blink text-neon">_</span>
      </div>
    </div>
  </div>
</template>
