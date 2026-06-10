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
  saving.value = true; message.value = null
  try {
    await apiFetch('/api/notifications', { method: 'PUT', body: JSON.stringify(form) })
    message.value = { type: 'success', text: 'Notification settings saved' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { saving.value = false }
}
async function test() {
  testing.value = true; message.value = null
  try {
    await apiFetch('/api/notifications/test', { method: 'POST' })
    message.value = { type: 'success', text: 'Test alert sent' }
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { testing.value = false }
}
load()
</script>

<style scoped>
.status-badge {
  font-size: 11px;
  line-height: 22px;
}
.status-badge--on {
  background: rgba(247, 147, 26, 0.1);
  border: 1px solid rgba(247, 147, 26, 0.2);
  color: #F7931A;
}
.status-badge--off {
  background: rgba(148, 163, 184, 0.1);
  border: 1px solid rgba(148, 163, 184, 0.3);
  color: #94A3B8;
}
.status-dot {
  width: 6px;
  height: 6px;
}
.status-dot--on { background: #F7931A; }
.status-dot--off { background: #94A3B8; }
.provider-btn {
  height: 32px;
  background: transparent;
  border: none;
  color: #94A3B8;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}
.provider-btn:hover {
  background: rgba(255,255,255,0.05);
}
.provider-btn--active {
  background: rgba(247, 147, 26, 0.15) !important;
  color: #F7931A !important;
}
</style>

<template>
  <div class="page">
    <header class="mb-6">
      <span class="page-badge">Console &middot; Alerts</span>
      <h1 class="page-title" style="margin-top: 4px;">Webhook notifications</h1>
    </header>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <div class="card">
      <div class="card-header">
        <div class="d-flex align-center ga-2">
          <v-icon size="18" color="#F7931A">mdi-bell</v-icon>
          <span class="card-title">Configuration</span>
        </div>
        <span class="rounded-pill px-2 font-mono d-inline-flex align-center ga-1 status-badge" :class="form.enabled ? 'status-badge--on' : 'status-badge--off'">
          <span class="rounded-full d-inline-block status-dot" :class="form.enabled ? 'status-dot--on' : 'status-dot--off'"></span>
          {{ form.enabled ? 'Enabled' : 'Disabled' }}
        </span>
      </div>
      <div class="card-body d-flex flex-column ga-5">
        <label class="d-flex align-center ga-3 font-mono" style="color: white; cursor: pointer;">
          <input type="checkbox" v-model="form.enabled" style="accent-color: #F7931A; width: 16px; height: 16px;" />
          Enable notifications
        </label>

        <div>
          <p class="label-upper">Provider</p>
          <div class="d-inline-flex rounded-lg overflow-hidden" style="border: 1px solid rgba(30, 41, 59, 0.6);">
            <button
              v-for="t in ['discord', 'telegram', 'slack']" :key="t"
              class="px-3 font-mono provider-btn"
              :class="{ 'provider-btn--active': form.type === t }"
              @click="form.type = t"
            >{{ t }}</button>
          </div>
        </div>

        <div>
          <label class="label-upper">Webhook URL</label>
          <input v-model="form.webhook_url" placeholder="https://..." class="input-line" />
        </div>

        <div class="d-flex ga-3 pt-1">
          <button class="btn btn-primary" :disabled="saving" @click="save">
            {{ saving ? 'Saving...' : 'Save settings' }}
          </button>
          <button class="btn btn-secondary" :disabled="testing || !form.enabled" @click="test">
            {{ testing ? 'Sending...' : 'Send test alert' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
