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
    message.value = { type: 'success', text: 'Saved' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
}

async function test() {
  testing.value = true
  try {
    await apiFetch('/api/notifications/test', { method: 'POST' })
    message.value = { type: 'success', text: 'Test notification sent' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { testing.value = false }
}

load()
</script>

<template>
  <div class="space-y-6">
    <header class="border-b border-border pb-6">
      <p class="section-marker mb-2">ALERTS</p>
      <h1 class="editorial-h1 !text-3xl">Notifications</h1>
    </header>

    <div class="card max-w-md space-y-4">
      <div v-if="message" class="p-3 rounded-md text-sm border" :class="message.type === 'success' ? 'border-emerald-200 bg-emerald-50 text-emerald-700' : 'border-red-200 bg-red-50 text-red-700'">
        {{ message.text }}
      </div>

      <label class="flex items-center gap-3 text-sm">
        <input type="checkbox" v-model="form.enabled" class="accent-accent w-4 h-4" />
        <span class="text-text font-medium">Enable notifications</span>
      </label>

      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Provider</label>
        <div class="flex gap-4">
          <label v-for="t in ['discord', 'telegram', 'slack']" :key="t" class="flex items-center gap-2 text-sm cursor-pointer" :class="form.type === t ? 'text-text font-medium' : 'text-text-muted'">
            <input type="radio" v-model="form.type" :value="t" class="accent-accent" /> {{ t }}
          </label>
        </div>
      </div>

      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Webhook URL</label>
        <input v-model="form.webhook_url" placeholder="https://..." class="input" />
      </div>

      <div class="flex gap-3">
        <button @click="save" class="btn-primary">Save</button>
        <button @click="test" :disabled="testing" class="btn-secondary">
          {{ testing ? 'Sending…' : 'Send Test' }}
        </button>
      </div>
    </div>
  </div>
</template>
