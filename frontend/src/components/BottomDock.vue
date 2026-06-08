<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useActionLogStore } from '../stores/actionLog'
const store = useActionLogStore()
const authStore = useAuthStore()

const barHeight = 32
const expanded = ref(false)
const dockHeight = ref(300)
const dragging = ref(false)

const liveLogs = ref<string[]>([])
const liveLogEl = ref<HTMLElement | null>(null)
let ws: WebSocket | null = null
let wsReconnectTimer: number | null = null
const historyLoaded = ref(false)

async function loadHistory() {
  if (!authStore.token || historyLoaded.value) return
  try {
    const res = await fetch(`/api/logs/history?limit=200`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (data.logs?.length) {
      liveLogs.value = data.logs.map((l: any) => decodeLog(typeof l === 'string' ? l : l.text || JSON.stringify(l)))
      historyLoaded.value = true
      nextTick(() => { if (liveLogEl.value) liveLogEl.value.scrollTop = liveLogEl.value.scrollHeight })
    }
  } catch {}
}

function connectWs() {
  if (!authStore.token) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${proto}//${location.host}/api/logs/ws?token=${encodeURIComponent(authStore.token)}`)
  ws.onmessage = (e) => {
    liveLogs.value.push(decodeLog(e.data))
    if (liveLogs.value.length > 2000) liveLogs.value.splice(0, liveLogs.value.length - 2000)
    nextTick(() => { if (liveLogEl.value) liveLogEl.value.scrollTop = liveLogEl.value.scrollHeight })
  }
  ws.onclose = () => {
    ws = null
    wsReconnectTimer = window.setTimeout(connectWs, 3000)
  }
}

onMounted(() => {
  if (expanded.value) { loadHistory(); connectWs() }
})

onUnmounted(() => {
  if (ws) { ws.onclose = null; ws.close(); ws = null }
  if (wsReconnectTimer !== null) clearTimeout(wsReconnectTimer)
})

watch(expanded, (v) => {
  if (v) { if (!historyLoaded.value) loadHistory(); if (!ws) connectWs() }
  if (!v && ws) { ws.onclose = null; ws.close(); ws = null }
  if (!v && wsReconnectTimer !== null) { clearTimeout(wsReconnectTimer); wsReconnectTimer = null }
})

const logEls = ref<Map<number, HTMLElement>>(new Map())

function setLogEl(id: number, el: HTMLElement | null) {
  if (el) logEls.value.set(id, el)
  else logEls.value.delete(id)
}

watch(
  () => store.actions.map((a) => a.lines.length),
  () => {
    nextTick(() => {
      for (const el of logEls.value.values()) el.scrollTop = el.scrollHeight
    })
  },
  { deep: true }
)

function toggleSection(id: number) {
  const a = store.actions.find((a) => a.id === id)
  if (a) a.lines = [...a.lines]
}

function decodeLog(raw: string): string {
  try {
    const decoded = atob(raw)
    try {
      const parsed = JSON.parse(decoded)
      return parsed.message || parsed.text || decoded
    } catch {
      return decoded
    }
  } catch {
    try {
      const parsed = JSON.parse(raw)
      return parsed.message || parsed.text || raw
    } catch {
      return raw
    }
  }
}

function colorizeLine(line: string): { text: string; color: string | null }[] {
  const rules: { re: RegExp; color: string }[] = [
    { re: /\b(ERR[OR]?|FATAL|CRITICAL|CRIT|PANIC)\b/g, color: '#EF4444' },
    { re: /\b(WARN[ING]?)\b/g, color: '#F59E0B' },
    { re: /\b(INFO?)\b/g, color: '#3B82F6' },
    { re: /\b(DEBUG|TRACE?)\b/g, color: '#8B5CF6' },
    { re: /\b(OK|SUCCESS|DONE)\b/g, color: '#10B981' },
    { re: /(https?:\/\/\S+)/g, color: '#60A5FA' },
    { re: /(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/g, color: '#34D399' },
    { re: /\b(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2})Z?\b/g, color: '#64748B' },
  ]
  const segments: { text: string; color: string | null }[] = [{ text: line, color: null }]
  for (const { re, color } of rules) {
    const newSegments: { text: string; color: string | null }[] = []
    for (const seg of segments) {
      if (seg.color !== null) {
        newSegments.push(seg)
        continue
      }
      let last = 0
      let m: RegExpExecArray | null
      const r = new RegExp(re.source, 'g')
      while ((m = r.exec(seg.text)) !== null) {
        if (m.index > last) newSegments.push({ text: seg.text.slice(last, m.index), color: null })
        newSegments.push({ text: m[1] || m[0], color })
        last = r.lastIndex
      }
      if (last < seg.text.length) newSegments.push({ text: seg.text.slice(last), color: null })
    }
    segments.length = 0
    segments.push(...newSegments)
  }
  return segments
}

function startDrag(e: MouseEvent) {
  dragging.value = true
  const startY = e.clientY
  const startH = dockHeight.value
  function onMove(ev: MouseEvent) {
    dockHeight.value = Math.max(120, startH + (startY - ev.clientY))
  }
  function onUp() {
    dragging.value = false
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}
</script>

<template>
  <div
    class="bottom-dock"
    :class="{ 'bottom-dock--expanded': expanded, 'bottom-dock--dragging': dragging }"
    :style="{ height: barHeight + (expanded ? dockHeight : 0) + 'px' }"
  >
    <div
      class="bottom-dock__bar d-flex align-center px-4 font-mono"
      style="height: 32px; cursor: pointer; font-size: 11px; color: #94A3B8; background: #0F1115; border-top: 1px solid rgba(30, 41, 59, 0.6); user-select: none;"
      @click="expanded = !expanded"
    >
      <v-icon size="14">{{ expanded ? 'mdi-chevron-down' : 'mdi-chevron-up' }}</v-icon>
      <span class="ml-1" style="text-transform: uppercase; letter-spacing: 0.08em;">
        {{ store.running.length ? `${store.running.length} running` : 'Console' }}
      </span>
      <v-spacer />
      <span class="font-mono" style="color: rgba(148, 163, 184, 0.5);">{{ liveLogs.length }} lines</span>
      <button
        class="ml-2 font-mono"
        style="background: none; border: none; color: #94A3B8; cursor: pointer; padding: 2px 6px; font-size: 11px;"
        @click.stop="store.clearFinished()"
      >Clear</button>
    </div>

    <div v-if="expanded" class="bottom-dock__body" :style="{ height: dockHeight + 'px' }">
      <div class="bottom-dock__resize-handle" @mousedown="startDrag" />

      <!-- Live tunnel logs -->
      <div class="pa-2" style="background: #0A0C10;">
        <div class="font-mono mb-1" style="color: rgba(148, 163, 184, 0.5); font-size: 10px; text-transform: uppercase; letter-spacing: 0.08em;">Live tunnel logs</div>
        <div
          ref="liveLogEl"
          class="font-mono pa-2 rounded-lg"
          style="color: #94A3B8; font-size: 11px; line-height: 1.35; max-height: 120px; overflow-y: auto; white-space: pre-wrap; background: rgba(0,0,0,0.3); border: 1px solid rgba(30, 41, 59, 0.4);"
        >
          <template v-if="liveLogs.length">
            <div v-for="(line, li) in liveLogs.slice(-100)" :key="li">
              <template v-for="(seg, si) in colorizeLine(line)" :key="si">
                <span v-if="seg.color" :style="{ color: seg.color }">{{ seg.text }}</span>
                <span v-else>{{ seg.text }}</span>
              </template>
            </div>
          </template>
          <span v-else style="color: rgba(148,163,184,0.4);">(connecting...)</span>
        </div>
      </div>

      <!-- Action logs -->
      <div class="pa-2" style="background: #0A0C10; border-top: 1px solid rgba(30, 41, 59, 0.4);">
        <div class="font-mono mb-1" style="color: rgba(148, 163, 184, 0.5); font-size: 10px; text-transform: uppercase; letter-spacing: 0.08em;">Actions</div>
        <div v-if="!store.actions.length" class="font-mono pa-2 text-center" style="color: rgba(148, 163, 184, 0.4); font-size: 12px;">
          No actions yet
        </div>
        <div v-for="a in store.actions" :key="a.id" class="mb-1 rounded-lg" style="background: rgba(0,0,0,0.3); border: 1px solid rgba(30, 41, 59, 0.4);">
          <div
            class="d-flex align-center ga-2 px-3 font-mono"
            style="height: 28px; font-size: 11px; cursor: pointer; border-bottom: 1px solid rgba(30, 41, 59, 0.2);"
            @click="toggleSection(a.id)"
          >
            <span
              class="rounded-full d-inline-block"
              style="width: 6px; height: 6px; flex-shrink: 0;"
              :style="{
                background: a.status === 'running' ? '#F7931A' : a.status === 'success' ? '#FFD600' : '#EF4444',
                boxShadow: a.status === 'running' ? '0 0 6px rgba(247,147,26,0.6)' : 'none',
              }"
            ></span>
            <span style="color: white;">{{ a.label }}</span>
            <span v-if="a.projectName" class="font-mono" style="color: #F7931A;">{{ a.projectName }}</span>
            <span class="font-mono" style="color: rgba(148, 163, 184, 0.4);">
              {{ a.status === 'running' ? 'RUNNING' : a.status === 'success' ? 'DONE' : 'FAILED' }}
            </span>
            <v-spacer />
            <button
              style="background: none; border: none; color: rgba(148, 163, 184, 0.4); cursor: pointer; padding: 2px; font-size: 11px;"
              @click.stop="store.remove(a.id)"
            >&times;</button>
          </div>
          <div
            :ref="(el) => setLogEl(a.id, el as HTMLElement | null)"
            class="font-mono pa-2"
            style="color: #94A3B8; font-size: 11px; line-height: 1.35; max-height: 80px; overflow-y: auto; white-space: pre-wrap;"
          >
            <template v-if="a.lines.length">
              <div v-for="(line, ai) in a.lines" :key="ai">
                <template v-for="(seg, si) in colorizeLine(line)" :key="si">
                  <span v-if="seg.color" :style="{ color: seg.color }">{{ seg.text }}</span>
                  <span v-else>{{ seg.text }}</span>
                </template>
              </div>
            </template>
            <span v-else style="color: rgba(148,163,184,0.4);">(waiting for output...)</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bottom-dock {
  position: fixed;
  left: 240px;
  right: 0;
  bottom: 0;
  z-index: 100;
  background: #0F1115;
  border-top: 1px solid rgba(30, 41, 59, 0.6);
}
.bottom-dock--dragging {
  user-select: none;
}
.bottom-dock__bar {
  border-top: none !important;
}
.bottom-dock__resize-handle {
  height: 8px;
  cursor: ns-resize;
  background: rgba(30, 41, 59, 0.6);
  transition: background 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0;
}
.bottom-dock__resize-handle::after {
  content: '';
  width: 32px;
  height: 3px;
  border-radius: 2px;
  background: rgba(148, 163, 184, 0.3);
}
.bottom-dock__resize-handle:hover {
  background: rgba(247, 147, 26, 0.3);
}
.bottom-dock__resize-handle:hover::after {
  background: rgba(247, 147, 26, 0.5);
}
.bottom-dock__body {
  overflow-y: auto;
  background: #0A0C10;
  display: flex;
  flex-direction: column;
}
</style>
