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
const expanded = ref<Set<string>>(new Set())

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
  estimateSize: () => 26,
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
  if (msg.includes('ERR') || msg.includes('error') || msg.includes('Error')) return 'text-error'
  if (msg.includes('WRN') || msg.includes('warn') || msg.includes('Warn')) return 'text-warning'
  return 'text-on-surface'
}

function formatTime(t?: string) {
  return t?.split('T')[1]?.slice(0, 8) || '\u2014'
}

function logKey(e: LogEntry) {
  return `${e.time}|${e.message}`
}

function isExpanded(e: LogEntry) {
  return expanded.value.has(logKey(e))
}

function toggle(e: LogEntry) {
  const k = logKey(e)
  const next = new Set(expanded.value)
  if (next.has(k)) next.delete(k)
  else next.add(k)
  expanded.value = next
}

function expandAll() {
  expanded.value = new Set(filtered.value.map(logKey))
}
function collapseAll() {
  expanded.value = new Set()
}

function clear() {
  logs.value = []
  expanded.value = new Set()
}

function setRowRef(el: Element | null) {
  if (el instanceof HTMLElement) rowVirtualizer.value.measureElement(el)
}
</script>

<template>
  <div class="pa-6 d-flex flex-column" style="max-width: 1280px; height: calc(100vh - 48px);">
    <div class="d-flex align-end justify-space-between ga-4 flex-shrink-0 mb-4">
      <div>
        <p class="eyebrow mb-2">Console &middot; Logs</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Cloudflared stream</h1>
      </div>
      <div class="d-flex align-center ga-2">
        <v-chip
          size="small"
          :color="wsStatus === 'connected' ? 'success' : wsStatus === 'connecting' ? 'warning' : 'error'"
          variant="tonal"
        >
          <template #prepend>
            <v-icon size="x-small">mdi-checkbox-blank-circle</v-icon>
          </template>
          {{ wsStatus }}
        </v-chip>
        <v-btn variant="tonal" @click="clear">Clear</v-btn>
      </div>
    </div>

    <div class="d-flex align-center ga-2 flex-shrink-0 mb-3 flex-wrap">
      <v-chip-group v-model="filter" mandatory color="primary" variant="tonal" density="compact">
        <v-chip value="all" size="small">all</v-chip>
        <v-chip value="error" size="small">error</v-chip>
        <v-chip value="warn" size="small">warn</v-chip>
        <v-chip value="info" size="small">info</v-chip>
      </v-chip-group>

      <v-text-field
        v-model="search"
        placeholder="Search..."
        prepend-inner-icon="mdi-magnify"
        clearable
        hide-details
        density="compact"
        variant="outlined"
        style="max-width: 260px;"
      />

      <v-spacer />

      <v-btn
        v-if="expanded.size"
        variant="text"
        size="small"
        @click="collapseAll"
      >Collapse all</v-btn>
      <v-btn
        v-else-if="filtered.length"
        variant="text"
        size="small"
        @click="expandAll"
      >Expand all</v-btn>

      <v-checkbox
        v-model="autoScroll"
        label="Tail"
        hide-details
        density="compact"
      />

      <span class="text-caption text-medium-emphasis font-mono tabular-nums">
        {{ filtered.length }} / {{ logs.length }}
      </span>
    </div>

    <v-card class="flex-grow-1 d-flex flex-column overflow-hidden">
      <div class="d-grid border-bottom bg-grey-lighten-4 eyebrow" style="grid-template-columns: 120px 1fr; border-bottom: 1px solid #ded9ca;">
        <div class="px-4" style="height: 36px; display: flex; align-items: center;">Time</div>
        <div class="px-4" style="height: 36px; display: flex; align-items: center;">Message</div>
      </div>

      <div ref="scrollEl" class="flex-grow-1 overflow-auto scrollbar-thin">
        <div
          v-if="filtered.length"
          class="position-relative"
          :style="{ height: rowVirtualizer.getTotalSize() + 'px' }"
        >
          <div
            v-for="vrow in rowVirtualizer.getVirtualItems()"
            :key="vrow.key"
            :ref="(el) => setRowRef(el as Element | null)"
            :data-index="vrow.index"
            class="position-absolute d-grid align-start border-bottom cursor-pointer transition-colors"
            :class="isExpanded(filtered[vrow.index]) ? 'bg-grey-lighten-3' : 'hover-bg-grey-lighten-4'"
            :style="{ transform: `translateY(${vrow.start}px)`, gridTemplateColumns: '120px 1fr', insetInline: 0, borderBottom: '1px solid #ded9ca' }"
            tabindex="0"
            @click="toggle(filtered[vrow.index])"
            @keydown.enter.prevent="toggle(filtered[vrow.index])"
            @keydown.space.prevent="toggle(filtered[vrow.index])"
          >
            <div class="px-4 py-1 font-mono text-caption text-medium-emphasis tabular-nums">{{ formatTime(filtered[vrow.index].time) }}</div>
            <div
              class="px-4 py-1 font-mono text-caption"
              :class="[
                levelClass(filtered[vrow.index].message),
                isExpanded(filtered[vrow.index]) ? 'text-pre-wrap word-break' : 'text-truncate'
              ]"
            >{{ filtered[vrow.index].message }}</div>
          </div>
        </div>
        <div v-else class="text-center text-medium-emphasis py-12 text-body-2">
          {{ logs.length ? 'No logs match the current filter' : 'Waiting for logs...' }}
        </div>
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.hover-bg-grey-lighten-4:hover {
  background: #f5f5f5;
}
.text-pre-wrap {
  white-space: pre-wrap;
  word-break: break-word;
}
.border-bottom {
  border-bottom: 1px solid #ded9ca;
}
</style>
