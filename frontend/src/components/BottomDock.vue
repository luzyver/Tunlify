<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useActionLogStore } from '../stores/actionLog'

const store = useActionLogStore()
const authStore = useAuthStore()
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
      liveLogs.value = data.logs.map((l: any) => typeof l === 'string' ? l : l.text || JSON.stringify(l))
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
    liveLogs.value.push(e.data)
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

    <div v-if="expanded" class="bottom-dock__resize-handle" @mousedown="startDrag" />

    <div
      v-if="expanded"
      class="bottom-dock__body"
      :style="{ maxHeight: dockHeight + 'px' }"
    >
      <!-- Live tunnel logs -->
      <div class="pa-2" style="background: #0A0C10;">
        <div class="font-mono mb-1" style="color: rgba(148, 163, 184, 0.5); font-size: 10px; text-transform: uppercase; letter-spacing: 0.08em;">Live tunnel logs</div>
        <pre
          ref="liveLogEl"
          class="font-mono pa-2 rounded-lg"
          style="color: #94A3B8; font-size: 11px; line-height: 1.35; max-height: 120px; overflow-y: auto; white-space: pre-wrap; background: rgba(0,0,0,0.3); border: 1px solid rgba(30, 41, 59, 0.4);"
        >{{ liveLogs.length ? liveLogs.slice(-100).join('\n') : '(connecting...)' }}</pre>
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
          <pre
            :ref="(el) => setLogEl(a.id, el as HTMLElement | null)"
            class="font-mono pa-2"
            style="color: #94A3B8; font-size: 11px; line-height: 1.35; max-height: 80px; overflow-y: auto; white-space: pre-wrap;"
          >{{ a.lines.length ? a.lines.join('\n') : '(waiting for output...)' }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bottom-dock--dragging {
  user-select: none;
}
.bottom-dock__resize-handle {
  height: 4px;
  cursor: ns-resize;
  background: rgba(30, 41, 59, 0.6);
  transition: background 0.2s;
}
.bottom-dock__resize-handle:hover {
  background: rgba(247, 147, 26, 0.4);
}
.bottom-dock__body {
  overflow-y: auto;
  background: #0A0C10;
  border-top: 1px solid rgba(30, 41, 59, 0.4);
}
</style>
