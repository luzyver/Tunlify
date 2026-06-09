<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useActionLogStore } from '../stores/actionLog'
const store = useActionLogStore()
const authStore = useAuthStore()

const barHeight = 32
const activeTab = ref<'tunnel' | 'action'>('tunnel')
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
    <div class="bottom-dock__bar">
      <v-icon size="14" class="dock-toggle" @click="expanded = !expanded">{{ expanded ? 'mdi-chevron-down' : 'mdi-chevron-up' }}</v-icon>
      <span class="dock-title">
        {{ store.running.length ? `${store.running.length} running` : 'Console' }}
      </span>
      <v-spacer />
      <button v-if="expanded" class="dock-clear" @click.stop="store.clearFinished()">Clear</button>
      <span class="dock-lines">{{ liveLogs.length }} lines</span>
    </div>

    <div v-if="expanded" class="bottom-dock__body" :style="{ height: dockHeight + 'px' }">
      <div class="bottom-dock__resize-handle" @mousedown="startDrag" />

      <div class="dock-tabs">
        <button
          class="dock-tab"
          :class="{ 'dock-tab--active': activeTab === 'tunnel' }"
          @click="activeTab = 'tunnel'"
        >Tunnel</button>
        <button
          class="dock-tab"
          :class="{ 'dock-tab--active': activeTab === 'action' }"
          @click="activeTab = 'action'"
        >Action</button>
      </div>

      <div v-show="activeTab === 'tunnel'" ref="liveLogEl" class="dock-content">
        <template v-if="liveLogs.length">
          <div v-for="(line, li) in liveLogs.slice(-100)" :key="li" class="dock-log-line">
            <template v-for="(seg, si) in colorizeLine(line)" :key="si">
              <span v-if="seg.color" :style="{ color: seg.color }">{{ seg.text }}</span>
              <span v-else>{{ seg.text }}</span>
            </template>
          </div>
        </template>
        <span class="dock-muted">(connecting...)</span>
      </div>

      <div v-show="activeTab === 'action'" class="dock-content">
        <template v-if="store.actions.length">
          <div v-for="a in store.actions" :key="a.id" class="dock-action-card">
            <div class="dock-action-header" @click="toggleSection(a.id)">
              <span
                class="dock-action-dot"
                :class="`dock-action-dot--${a.status}`"
              ></span>
              <span class="dock-action-label">{{ a.label }}</span>
              <span v-if="a.projectName" class="dock-action-project">{{ a.projectName }}</span>
              <span class="dock-action-status" :class="`dock-action-status--${a.status}`">
                {{ a.status === 'running' ? 'RUNNING' : a.status === 'success' ? 'DONE' : 'FAILED' }}
              </span>
              <v-spacer />
              <button class="dock-action-remove" @click.stop="store.remove(a.id)">&times;</button>
            </div>
            <div
              :ref="(el) => setLogEl(a.id, el as HTMLElement | null)"
              class="dock-action-body"
            >
              <template v-if="a.lines.length">
                <div v-for="(line, ai) in a.lines" :key="ai" class="dock-log-line">
                  <template v-for="(seg, si) in colorizeLine(line)" :key="si">
                    <span v-if="seg.color" :style="{ color: seg.color }">{{ seg.text }}</span>
                    <span v-else>{{ seg.text }}</span>
                  </template>
                </div>
              </template>
              <span class="dock-muted">(waiting for output...)</span>
            </div>
          </div>
        </template>
        <span class="dock-muted" style="font-size: 11px;">No actions yet</span>
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

/* Bar */
.bottom-dock__bar {
  display: flex;
  align-items: center;
  padding: 0 12px;
  height: 32px;
  font-size: 11px;
  color: #94A3B8;
  background: #0F1115;
  user-select: none;
  border-top: 1px solid rgba(30, 41, 59, 0.6);
  gap: 8px;
}
.dock-toggle {
  cursor: pointer;
  color: #71717A;
  transition: color 0.2s;
}
.dock-toggle:hover {
  color: #F7931A;
}
.dock-title {
  font-family: 'Sora', sans-serif;
  font-weight: 600;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #94A3B8;
}
.dock-clear {
  background: none;
  border: none;
  color: #71717A;
  cursor: pointer;
  padding: 2px 6px;
  font-size: 11px;
  font-family: 'Sora', sans-serif;
  font-weight: 500;
  transition: color 0.2s;
}
.dock-clear:hover {
  color: #F7931A;
}
.dock-lines {
  font-family: 'Sora', sans-serif;
  font-size: 10px;
  color: rgba(148, 163, 184, 0.4);
}

/* Resize handle */
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
  transition: background 0.2s;
}
.bottom-dock__resize-handle:hover {
  background: rgba(247, 147, 26, 0.3);
}
.bottom-dock__resize-handle:hover::after {
  background: rgba(247, 147, 26, 0.5);
}

/* Body */
.bottom-dock__body {
  overflow-y: auto;
  background: #0A0C10;
  display: flex;
  flex-direction: column;
}

/* Tabs */
.dock-tabs {
  display: flex;
  border-bottom: 1px solid rgba(30, 41, 59, 0.4);
  background: #0A0C10;
  padding: 0 8px;
}
.dock-tab {
  padding: 6px 12px;
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  border: none;
  background: none;
  cursor: pointer;
  color: rgba(148, 163, 184, 0.5);
  border-bottom: 2px solid transparent;
  font-family: 'Sora', sans-serif;
  font-weight: 600;
  transition: color 0.2s, border-color 0.2s;
}
.dock-tab:hover {
  color: #94A3B8;
}
.dock-tab--active {
  color: #F7931A !important;
  border-bottom-color: #F7931A;
}

/* Content */
.dock-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  background: #0A0C10;
}
.dock-log-line {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  line-height: 1.4;
  color: #94A3B8;
}
.dock-muted {
  color: rgba(148, 163, 184, 0.4);
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
}

/* Action cards */
.dock-action-card {
  margin-bottom: 4px;
  border-radius: 8px;
  background: rgba(0,0,0,0.3);
  border: 1px solid rgba(30, 41, 59, 0.4);
  overflow: hidden;
}
.dock-action-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  height: 28px;
  font-size: 11px;
  cursor: pointer;
  border-bottom: 1px solid rgba(30, 41, 59, 0.2);
  font-family: 'JetBrains Mono', monospace;
}
.dock-action-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dock-action-dot--running {
  background: #F7931A;
  box-shadow: 0 0 6px rgba(247,147,26,0.6);
}
.dock-action-dot--success {
  background: #FFD600;
}
.dock-action-dot--error {
  background: #EF4444;
}
.dock-action-label {
  color: white;
}
.dock-action-project {
  color: #F7931A;
}
.dock-action-status {
  font-size: 9px;
  letter-spacing: 0.06em;
}
.dock-action-status--running { color: #F7931A; }
.dock-action-status--success { color: #FFD600; }
.dock-action-status--error { color: #EF4444; }
.dock-action-remove {
  background: none;
  border: none;
  color: rgba(148, 163, 184, 0.4);
  cursor: pointer;
  padding: 2px;
  font-size: 13px;
  transition: color 0.2s;
}
.dock-action-remove:hover {
  color: #EF4444;
}
.dock-action-body {
  padding: 8px;
  font-family: 'JetBrains Mono', monospace;
  color: #94A3B8;
  font-size: 11px;
  line-height: 1.35;
  max-height: 100px;
  overflow-y: auto;
  white-space: pre-wrap;
}
</style>
