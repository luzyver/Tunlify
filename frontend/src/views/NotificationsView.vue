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
    message.value = { type: 'success', text: 'Test sent' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { testing.value = false }
}

load()
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold tracking-tight">Notifications</h1>

    <div class="card max-w-md">
      <div class="card-header">Webhook Alerts</div>
      <div class="card-body space-y-4">
        <div v-if="message" class="p-3 rounded text-sm" :class="message.type === 'success' ? 'bg-success/10 text-emerald-800' : 'bg-danger/10 text-danger'">
          {{ message.text }}
        </div>

        <label class="flex items-center gap-3 text-sm">
          <input type="checkbox" v-model="form.enabled" class="accent-accent w-4 h-4" />
          <span class="font-medium">Enable notifications</span>
        </label>

        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Provider</label>
          <div class="flex gap-4">
            <label v-for="t in ['discord', 'telegram', 'slack']" :key="t" class="flex items-center gap-2 text-sm cursor-pointer" :class="form.type === t ? 'text-accent font-medium' : 'text-text-muted'">
              <input type="radio" v-model="form.type" :value="t" class="accent-accent" /> {{ t }}
            </label>
          </div>
        </div>

        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Webhook URL</label>
          <input v-model="form.webhook_url" placeholder="https://..." class="input" />
        </div>

        <div class="flex gap-2">
          <button @click="save" class="btn-primary">Save</button>
          <button @click="test" :disabled="testing" class="btn-secondary">{{ testing ? 'Sending…' : 'Test' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
