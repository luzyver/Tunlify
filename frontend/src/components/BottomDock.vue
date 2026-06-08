<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useActionLogStore } from '../stores/actionLog'

const store = useActionLogStore()
const expanded = ref(false)
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

function toggle(id: number) {
  const a = store.actions.find((a) => a.id === id)
  if (a) a.lines = [...a.lines]
}
</script>

<template>
  <div
    class="bottom-dock"
    :class="{ 'bottom-dock--expanded': expanded }"
  >
    <div
      class="bottom-dock__bar d-flex align-center px-4 font-mono"
      style="height: 32px; cursor: pointer; font-size: 11px; color: #94A3B8; background: #0F1115; border-top: 1px solid rgba(30, 41, 59, 0.6); user-select: none;"
      @click="expanded = !expanded"
    >
      <v-icon size="14">mdi-chevron-right</v-icon>
      <span class="ml-1" style="text-transform: uppercase; letter-spacing: 0.08em;">
        {{ store.running.length }} running
      </span>
      <v-spacer />
      <span class="font-mono" style="color: rgba(148, 163, 184, 0.5);">{{ store.actions.length }} total</span>
      <button
        class="ml-2 font-mono"
        style="background: none; border: none; color: #94A3B8; cursor: pointer; padding: 2px 6px; font-size: 11px;"
        @click.stop="store.clearFinished()"
      >Clear</button>
    </div>

    <div v-if="expanded" class="bottom-dock__body pa-2" style="background: #0A0C10; border-top: 1px solid rgba(30, 41, 59, 0.4); max-height: 240px; overflow-y: auto;">
      <div v-if="!store.actions.length" class="font-mono pa-4 text-center" style="color: rgba(148, 163, 184, 0.4); font-size: 12px;">
        No actions yet
      </div>
      <div v-for="a in store.actions" :key="a.id" class="mb-1 rounded-lg" style="background: rgba(0,0,0,0.3); border: 1px solid rgba(30, 41, 59, 0.4);">
        <div
          class="d-flex align-center ga-2 px-3 font-mono"
          style="height: 28px; font-size: 11px; cursor: pointer; border-bottom: 1px solid rgba(30, 41, 59, 0.2);"
          @click="toggle(a.id)"
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
          style="color: #94A3B8; font-size: 11px; line-height: 1.35; max-height: 120px; overflow-y: auto; white-space: pre-wrap;"
        >{{ a.lines.length ? a.lines.join('\n') : '(waiting for output...)' }}</pre>
      </div>
    </div>
  </div>
</template>
