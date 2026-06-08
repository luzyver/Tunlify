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

<template>
  <div class="pa-6">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; Alerts</p>
      <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Webhook notifications</h1>
    </header>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <div class="rounded-2xl" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <div class="d-flex align-center ga-2">
          <v-icon size="18" color="#F7931A">mdi-bell</v-icon>
          <span class="font-heading font-semibold" style="color: white;">Configuration</span>
        </div>
        <span class="rounded-pill px-2 font-mono d-inline-flex align-center ga-1" :style="{ background: form.enabled ? 'rgba(255, 214, 0, 0.1)' : 'rgba(148, 163, 184, 0.1)', border: '1px solid ' + (form.enabled ? 'rgba(255, 214, 0, 0.3)' : 'rgba(148, 163, 184, 0.3)'), color: form.enabled ? '#FFD600' : '#94A3B8', fontSize: '11px', lineHeight: '22px' }">
          <span class="rounded-full d-inline-block" :style="{ width: '6px', height: '6px', background: form.enabled ? '#FFD600' : '#94A3B8' }"></span>
          {{ form.enabled ? 'Enabled' : 'Disabled' }}
        </span>
      </div>
      <div class="pa-6 d-flex flex-column ga-5">
        <label class="d-flex align-center ga-3 font-mono" style="color: white; cursor: pointer;">
          <input type="checkbox" v-model="form.enabled" style="accent-color: #F7931A; width: 16px; height: 16px;" />
          Enable notifications
        </label>

        <div>
          <p class="font-mono mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Provider</p>
          <div class="d-inline-flex rounded-lg overflow-hidden" style="border: 1px solid rgba(30, 41, 59, 0.6);">
            <button v-for="t in ['discord', 'telegram', 'slack']" :key="t"
              class="px-3 font-mono"
              :style="{ height: '32px', background: form.type === t ? 'rgba(247, 147, 26, 0.15)' : 'transparent', border: 'none', color: form.type === t ? '#F7931A' : '#94A3B8', fontSize: '11px', cursor: 'pointer', transition: 'all 0.2s' }"
              @click="form.type = t"
              @mouseenter="$event.target.style.background = form.type === t ? 'rgba(247, 147, 26, 0.15)' : 'rgba(255,255,255,0.05)'"
              @mouseleave="$event.target.style.background = form.type === t ? 'rgba(247, 147, 26, 0.15)' : 'transparent'"
            >{{ t }}</button>
          </div>
        </div>

        <div>
          <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Webhook URL</label>
          <input v-model="form.webhook_url" placeholder="https://..." style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
        </div>

        <div class="d-flex ga-3 pt-1">
          <button
            class="d-inline-flex align-center ga-2 rounded-pill px-5 font-mono"
            style="height: 40px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
            :disabled="saving"
            @click="save"
            @mouseenter="if(!saving) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)' }"
            @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
          >{{ saving ? 'Saving...' : 'Save settings' }}</button>
          <button
            class="rounded-pill px-4 font-mono"
            style="height: 40px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 12px; cursor: pointer; transition: all 0.3s;"
            :disabled="testing || !form.enabled"
            @click="test"
            @mouseenter="if(!testing && form.enabled) { $event.target.style.background = 'rgba(247, 147, 26, 0.2)'; $event.target.style.borderColor = '#F7931A' }"
            @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'; $event.target.style.borderColor = 'rgba(247, 147, 26, 0.2)'"
          >{{ testing ? 'Sending...' : 'Send test alert' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
