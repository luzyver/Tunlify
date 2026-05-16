<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const form = reactive({ enabled: false, webhook_url: '', type: 'discord' })
const message = ref<{ type: string; text: string } | null>(null)
const testing = ref(false)

async function load() {
  try { Object.assign(form, await apiFetch('/api/notifications')) } catch {}
}

async function save() {
  message.value = null
  try {
    await apiFetch('/api/notifications', { method: 'PUT', body: JSON.stringify(form) })
    message.value = { type: 'success', text: 'SAVED' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
}

async function test() {
  testing.value = true
  try {
    await apiFetch('/api/notifications/test', { method: 'POST' })
    message.value = { type: 'success', text: 'TEST SENT' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { testing.value = false }
}

load()
</script>

<template>
  <div class="space-y-6">
    <h1 class="cyber-heading text-xl lg:text-2xl">Notifications</h1>

    <div class="cyber-card max-w-md space-y-4">
      <div v-if="message" class="border p-3 text-sm font-mono" :class="message.type === 'success' ? 'border-neon text-neon bg-neon/5' : 'border-danger text-danger bg-danger/5'">
        [{{ message.type === 'success' ? 'OK' : 'ERR' }}] {{ message.text }}
      </div>

      <label class="flex items-center gap-3 text-sm">
        <input type="checkbox" v-model="form.enabled" class="accent-neon w-4 h-4" />
        <span class="text-gray-200 uppercase tracking-wider text-xs">Enable Notifications</span>
      </label>

      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">Type</label>
        <div class="flex gap-4">
          <label v-for="t in ['discord', 'telegram', 'slack']" :key="t" class="flex items-center gap-2 text-xs uppercase tracking-wider" :class="form.type === t ? 'text-neon' : 'text-muted'">
            <input type="radio" v-model="form.type" :value="t" class="accent-neon" /> {{ t }}
          </label>
        </div>
      </div>

      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">Webhook URL</label>
        <input v-model="form.webhook_url" placeholder="https://..." class="cyber-input" />
      </div>

      <div class="flex gap-2">
        <button @click="save" class="cyber-btn">Save</button>
        <button @click="test" :disabled="testing" class="cyber-btn-secondary">
          {{ testing ? 'SENDING...' : 'TEST' }}
        </button>
      </div>
    </div>
  </div>
</template>
