<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { useVirtualizer } from '@tanstack/vue-virtual'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

interface LogEntry { time: string; message: string }

const logs = ref<LogEntry[]>([])
const filter = ref<'all' | 'error' | 'warn' | 'info'>('all')
const search = ref('')
const autoScroll = ref(true)
const wsStatus = ref<'connecting' | 'connected' | 'disconnected'>('connecting')
const scrollEl = ref<HTMLElement | null>(null)

const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/api/logs/ws?token=${authStore.token}`
let ws: WebSocket | null = null

function connect() {
  ws = new WebSocket(wsUrl)
  ws.onopen = () => { wsStatus.value = 'connected' }
  ws.onclose = () => {
    wsStatus.value = 'disconnected'
    setTimeout(connect, 3000)
  }
  ws.onmessage = (e) => {
    try {
      const entry = JSON.parse(e.data) as LogEntry
      logs.value.push(entry)
      if (logs.value.length > 5000) logs.value.splice(0, logs.value.length - 5000)
    } catch {}
  }
}
connect()
onUnmounted(() => ws?.close())

const filtered = computed(() => {
  let result = logs.value
  if (filter.value !== 'all') result = result.filter((l) => l.message.toLowerCase().includes(filter.value))
  if (search.value) {
    const s = search.value.toLowerCase()
    result = result.filter((l) => l.message.toLowerCase().includes(s))
  }
  return result
})

const rowVirtualizer = useVirtualizer({
  get count() { return filtered.value.length },
  getScrollElement: () => scrollEl.value,
  estimateSize: () => 24,
  overscan: 12,
})

watch(
  () => filtered.value.length,
  () => {
    if (!autoScroll.value) return
    nextTick(() => {
      const el = scrollEl.value
      if (!el) return
      el.scrollTop = el.scrollHeight
    })
  }
)

function levelClass(msg: string) {
  if (msg.includes('ERR') || msg.includes('error') || msg.includes('Error')) return 'text-danger'
  if (msg.includes('WRN') || msg.includes('warn') || msg.includes('Warn')) return 'text-warning'
  return 'text-text'
}

function formatTime(t?: string) {
  return t?.split('T')[1]?.slice(0, 8) || '—'
}

function clear() { logs.value = [] }
</script>

<template>
  <div class="space-y-6 h-[calc(100vh-3rem)] flex flex-col">
    <header class="flex items-end justify-between gap-4 shrink-0">
      <div>
        <p class="eyebrow mb-2">Console · Logs</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Cloudflared stream</h1>
      </div>
      <div class="flex items-center gap-2">
        <span class="badge" :class="wsStatus === 'connected' ? 'badge-success' : wsStatus === 'connecting' ? 'badge-warning' : 'badge-danger'">
          <span class="dot" :class="wsStatus === 'connected' ? 'bg-success' : wsStatus === 'connecting' ? 'bg-warning' : 'bg-danger'"></span>
          {{ wsStatus }}
        </span>
        <button @click="clear" class="btn-secondary">Clear</button>
      </div>
    </header>

    <div class="flex flex-wrap items-center gap-2 shrink-0">
      <div class="inline-flex border border-border rounded-md overflow-hidden bg-surface">
        <button
          v-for="f in ['all', 'error', 'warn', 'info'] as const"
          :key="f"
          type="button"
          class="px-3 h-8 text-xs font-medium transition-colors duration-100"
          :class="filter === f ? 'bg-accent-soft text-accent' : 'text-text-muted hover:bg-bg-alt'"
          @click="filter = f"
        >
          {{ f }}
        </button>
      </div>
      <input
        v-model="search"
        placeholder="Search…"
        class="input !py-1.5 max-w-[260px]"
      />
      <label class="ml-auto flex items-center gap-2 text-xs text-text-muted">
        <input type="checkbox" v-model="autoScroll" class="accent-accent" /> Tail
      </label>
      <span class="text-2xs text-text-dim font-mono tabular-nums">
        {{ filtered.length }} / {{ logs.length }}
      </span>
    </div>

    <div class="card flex-1 overflow-hidden flex flex-col min-h-0">
      <div class="grid border-b border-border bg-bg-alt eyebrow shrink-0" style="grid-template-columns: 120px 1fr;">
        <div class="px-4 h-9 flex items-center">Time</div>
        <div class="px-4 h-9 flex items-center">Message</div>
      </div>

      <div ref="scrollEl" class="flex-1 overflow-auto scrollbar-thin">
        <div
          v-if="filtered.length"
          class="relative w-full"
          :style="{ height: rowVirtualizer.getTotalSize() + 'px' }"
        >
          <div
            v-for="vrow in rowVirtualizer.getVirtualItems()"
            :key="vrow.key"
            class="absolute inset-x-0 grid items-center border-b border-border hover:bg-surface-alt"
            :style="{
              transform: `translateY(${vrow.start}px)`,
              height: vrow.size + 'px',
              gridTemplateColumns: '120px 1fr',
            }"
          >
            <div class="px-4 num text-xs text-text-dim">{{ formatTime(filtered[vrow.index].time) }}</div>
            <div class="px-4 font-mono text-xs leading-5 truncate" :class="levelClass(filtered[vrow.index].message)">
              {{ filtered[vrow.index].message }}
            </div>
          </div>
        </div>
        <div v-else class="text-center text-text-muted py-12 text-sm">
          {{ logs.length ? 'No logs match the current filter' : 'Waiting for logs…' }}
        </div>
      </div>
    </div>
  </div>
</template>
