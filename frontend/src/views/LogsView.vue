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

function levelColor(msg: string) {
  if (msg.includes('ERR') || msg.includes('error') || msg.includes('Error')) return '#EF4444'
  if (msg.includes('WRN') || msg.includes('warn') || msg.includes('Warn')) return '#F7931A'
  return 'white'
}

function formatTime(t?: string) {
  return t?.split('T')[1]?.slice(0, 8) || '\u2014'
}

function logKey(e: LogEntry) { return `${e.time}|${e.message}` }
function isExpanded(e: LogEntry) { return expanded.value.has(logKey(e)) }

function toggle(e: LogEntry) {
  const k = logKey(e)
  const next = new Set(expanded.value)
  if (next.has(k)) next.delete(k); else next.add(k)
  expanded.value = next
}

function expandAll() { expanded.value = new Set(filtered.value.map(logKey)) }
function collapseAll() { expanded.value = new Set() }
function clear() { logs.value = []; expanded.value = new Set() }

function setRowRef(el: Element | null) {
  if (el instanceof HTMLElement) rowVirtualizer.value.measureElement(el)
}

function handleMouseLeave(e: MouseEvent, entry: LogEntry) {
  if (!isExpanded(entry)) {
    (e.currentTarget as HTMLElement).style.background = 'transparent'
  }
}
</script>

<template>
  <div class="pa-6 d-flex flex-column" style="max-width: 1280px; height: calc(100vh - 48px);">
    <div class="d-flex align-center justify-space-between ga-4 flex-shrink-0 mb-4 flex-wrap" style="gap: 12px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Logs</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Cloudflared stream</h1>
      </div>
      <div class="d-flex align-center ga-2">
        <span class="rounded-pill px-3 font-mono d-inline-flex align-center ga-1" :style="{ height: '28px', background: wsStatus === 'connected' ? 'rgba(255, 214, 0, 0.1)' : wsStatus === 'connecting' ? 'rgba(247, 147, 26, 0.1)' : 'rgba(239, 68, 68, 0.1)', border: '1px solid ' + (wsStatus === 'connected' ? 'rgba(255, 214, 0, 0.3)' : wsStatus === 'connecting' ? 'rgba(247, 147, 26, 0.3)' : 'rgba(239, 68, 68, 0.3)'), color: wsStatus === 'connected' ? '#FFD600' : wsStatus === 'connecting' ? '#F7931A' : '#EF4444', fontSize: '11px' }">
          <span class="rounded-full d-inline-block" :style="{ width: '6px', height: '6px', background: wsStatus === 'connected' ? '#FFD600' : wsStatus === 'connecting' ? '#F7931A' : '#EF4444' }"></span>
          {{ wsStatus }}
        </span>
        <button
          class="rounded-pill px-3 font-mono"
          style="height: 28px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 11px; cursor: pointer; transition: all 0.2s;"
          @click="clear"
          @mouseenter="$event.target.style.background = 'rgba(247, 147, 26, 0.2)'"
          @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'"
        >Clear</button>
      </div>
    </div>

    <div class="d-flex align-center ga-2 flex-shrink-0 mb-3 flex-wrap" style="gap: 8px;">
      <div class="d-inline-flex rounded-lg overflow-hidden" style="border: 1px solid rgba(30, 41, 59, 0.6);">
        <button v-for="f in ['all', 'error', 'warn', 'info'] as const" :key="f"
          class="px-3 font-mono"
          :style="{ height: '32px', background: filter === f ? 'rgba(247, 147, 26, 0.15)' : 'transparent', border: 'none', color: filter === f ? '#F7931A' : '#94A3B8', fontSize: '11px', cursor: 'pointer', transition: 'all 0.2s' }"
          @click="filter = f"
          @mouseenter="$event.target.style.background = filter === f ? 'rgba(247, 147, 26, 0.15)' : 'rgba(255,255,255,0.05)'"
          @mouseleave="$event.target.style.background = filter === f ? 'rgba(247, 147, 26, 0.15)' : 'transparent'"
        >{{ f }}</button>
      </div>

      <div class="position-relative" style="max-width: 260px;">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94A3B8" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="position-absolute" style="left: 10px; top: 50%; transform: translateY(-50%); z-index: 1;"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
        <input v-model="search" placeholder="Search..." style="width: 100%; background: rgba(0,0,0,0.5); border: 1px solid rgba(30, 41, 59, 0.8); border-radius: 8px; color: white; padding: 6px 10px 6px 32px; font-family: 'JetBrains Mono', monospace; font-size: 12px; height: 32px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
      </div>

      <v-spacer />

      <button v-if="expanded.size" class="font-mono text-decoration-none" style="background: none; border: none; color: #94A3B8; font-size: 11px; cursor: pointer; transition: color 0.2s;" @click="collapseAll" @mouseenter="$event.target.style.color = '#F7931A'" @mouseleave="$event.target.style.color = '#94A3B8'">Collapse all</button>
      <button v-else-if="filtered.length" class="font-mono text-decoration-none" style="background: none; border: none; color: #94A3B8; font-size: 11px; cursor: pointer; transition: color 0.2s;" @click="expandAll" @mouseenter="$event.target.style.color = '#F7931A'" @mouseleave="$event.target.style.color = '#94A3B8'">Expand all</button>

      <label class="d-flex align-center ga-1 font-mono" style="color: #94A3B8; font-size: 11px; cursor: pointer;">
        <input type="checkbox" v-model="autoScroll" style="accent-color: #F7931A;" /> Tail
      </label>

      <span class="font-mono tabular-nums" style="color: #94A3B8; font-size: 11px;">{{ filtered.length }} / {{ logs.length }}</span>
    </div>

    <div class="rounded-xl flex-grow-1 d-flex flex-column overflow-hidden" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-grid eyebrow" style="grid-template-columns: 120px 1fr; border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <div class="px-4" style="height: 36px; display: flex; align-items: center;">Time</div>
        <div class="px-4" style="height: 36px; display: flex; align-items: center;">Message</div>
      </div>

      <div ref="scrollEl" class="flex-grow-1 overflow-auto scrollbar-thin">
        <div v-if="filtered.length" class="position-relative" :style="{ height: rowVirtualizer.getTotalSize() + 'px' }">
          <div
            v-for="vrow in rowVirtualizer.getVirtualItems()" :key="vrow.key"
            :ref="(el) => setRowRef(el as Element | null)"
            :data-index="vrow.index"
            class="position-absolute d-grid align-start"
            :style="{
              transform: `translateY(${vrow.start}px)`,
              gridTemplateColumns: '120px 1fr',
              insetInline: 0,
              borderBottom: '1px solid rgba(30, 41, 59, 0.2)',
              background: isExpanded(filtered[vrow.index]) ? 'rgba(247, 147, 26, 0.03)' : 'transparent',
            }"
            tabindex="0"
            @click="toggle(filtered[vrow.index])"
            @keydown.enter.prevent="toggle(filtered[vrow.index])"
            @keydown.space.prevent="toggle(filtered[vrow.index])"
            @mouseenter="$event.currentTarget.style.background = 'rgba(247, 147, 26, 0.03)'"
            @mouseleave="handleMouseLeave($event, filtered[vrow.index])"
          >
            <div class="px-4 py-1 font-mono" style="color: #94A3B8; font-size: 12px; font-variant-numeric: tabular-nums;">{{ formatTime(filtered[vrow.index].time) }}</div>
            <div
              class="px-4 py-1 font-mono" style="font-size: 12px; line-height: 1.5; cursor: pointer;"
              :style="{
                color: levelColor(filtered[vrow.index].message),
                whiteSpace: isExpanded(filtered[vrow.index]) ? 'pre-wrap' : 'nowrap',
                overflow: isExpanded(filtered[vrow.index]) ? 'visible' : 'hidden',
                textOverflow: isExpanded(filtered[vrow.index]) ? 'clip' : 'ellipsis',
                wordBreak: isExpanded(filtered[vrow.index]) ? 'break-word' : 'normal',
              }"
            >{{ filtered[vrow.index].message }}</div>
          </div>
        </div>
        <div v-else class="text-center font-mono py-12" style="color: #94A3B8; font-size: 13px;">
          {{ logs.length ? 'No logs match the current filter' : 'Waiting for logs...' }}
        </div>
      </div>
    </div>
  </div>
</template>
