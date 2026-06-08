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
  <div class="pa-6" style="max-width: 1280px;">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; Alerts</p>
      <h1 class="text-h4 font-weight-semibold text-on-surface">Webhook notifications</h1>
    </header>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <v-card>
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Configuration</span>
        <v-chip
          :color="form.enabled ? 'success' : 'default'"
          size="x-small"
          variant="tonal"
        >
          <template #prepend>
            <v-icon size="x-small">mdi-checkbox-blank-circle</v-icon>
          </template>
          {{ form.enabled ? 'Enabled' : 'Disabled' }}
        </v-chip>
      </v-card-title>
      <v-card-text class="d-flex flex-column ga-5">
        <v-checkbox
          v-model="form.enabled"
          label="Enable notifications"
          hide-details
          density="compact"
        />

        <div>
          <p class="text-caption font-weight-medium text-medium-emphasis mb-2">Provider</p>
          <v-chip-group v-model="form.type" mandatory color="primary" variant="tonal" density="compact">
            <v-chip value="discord" size="small">discord</v-chip>
            <v-chip value="telegram" size="small">telegram</v-chip>
            <v-chip value="slack" size="small">slack</v-chip>
          </v-chip-group>
        </div>

        <v-text-field
          v-model="form.webhook_url"
          label="Webhook URL"
          placeholder="https://..."
          variant="outlined"
          density="compact"
          hide-details
        />

        <div class="d-flex ga-3 pt-1">
          <v-btn color="primary" :loading="saving" @click="save">
            {{ saving ? 'Saving...' : 'Save settings' }}
          </v-btn>
          <v-btn variant="tonal" :loading="testing" :disabled="testing || !form.enabled" @click="test">
            {{ testing ? 'Sending...' : 'Send test alert' }}
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>
