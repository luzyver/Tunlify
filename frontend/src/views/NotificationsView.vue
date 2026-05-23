<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const form = reactive({ enabled: false, webhook_url: '', type: 'discord' })
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const testing = ref(false)
const saving = ref(false)

async function load() {
  try { Object.assign(form, await apiFetch('/api/notifications')) } catch {}
}

async function save() {
  saving.value = true
  message.value = null
  try {
    await apiFetch('/api/notifications', { method: 'PUT', body: JSON.stringify(form) })
    message.value = { type: 'success', text: 'Notification settings saved' }
  } catch (e: any) {
    message.value = { type: 'error', text: e.message }
  } finally {
    saving.value = false
  }
}

async function test() {
  testing.value = true
  message.value = null
  try {
    await apiFetch('/api/notifications/test', { method: 'POST' })
    message.value = { type: 'success', text: 'Test alert sent' }
  } catch (e: any) {
    message.value = { type: 'error', text: e.message }
  } finally {
    testing.value = false
  }
}

load()
</script>

<template>
  <div class="space-y-6">
    <header>
      <p class="eyebrow mb-2">Console · Alerts</p>
      <h1 class="text-2xl font-semibold tracking-tight text-text">Webhook notifications</h1>
    </header>

    <div v-if="message" :class="message.type === 'success' ? 'alert-success' : 'alert-danger'">
      {{ message.text }}
    </div>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Configuration</span>
        <span class="badge" :class="form.enabled ? 'badge-success' : 'badge-neutral'">
          <span class="dot" :class="form.enabled ? 'bg-success' : 'bg-text-dim'"></span>
          {{ form.enabled ? 'Enabled' : 'Disabled' }}
        </span>
      </div>
      <div class="card-body space-y-5">
        <label class="flex items-center gap-3 text-sm cursor-pointer select-none">
          <input type="checkbox" v-model="form.enabled" class="accent-accent w-4 h-4" />
          <span class="font-medium text-text">Enable notifications</span>
        </label>

        <div>
          <label class="field-label">Provider</label>
          <div class="inline-flex border border-border rounded-md overflow-hidden bg-surface">
            <button
              v-for="t in ['discord', 'telegram', 'slack']"
              :key="t"
              type="button"
              class="px-3 h-8 text-xs font-medium transition-colors duration-100"
              :class="form.type === t ? 'bg-accent-soft text-accent' : 'text-text-muted hover:bg-bg-alt'"
              @click="form.type = t"
            >
              {{ t }}
            </button>
          </div>
        </div>

        <div>
          <label class="field-label">Webhook URL</label>
          <input v-model="form.webhook_url" placeholder="https://…" class="input font-mono" />
        </div>

        <div class="flex items-center gap-3 pt-1">
          <button @click="save" :disabled="saving" class="btn-primary">
            {{ saving ? 'Saving…' : 'Save settings' }}
          </button>
          <button @click="test" :disabled="testing || !form.enabled" class="btn-secondary">
            {{ testing ? 'Sending…' : 'Send test alert' }}
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
