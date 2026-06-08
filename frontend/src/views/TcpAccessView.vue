<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()

const form = reactive({ hostname: '', local_url: 'localhost:', mode: 'foreground' })
const generated = ref('')
const error = ref('')
const loading = ref(false)
const snackbar = ref(false)

async function generate() {
  loading.value = true
  error.value = ''
  try {
    const data = await apiFetch<{ command: string }>('/api/tcp-access/generate', {
      method: 'POST',
      body: JSON.stringify(form),
    })
    generated.value = data.command
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function copyCommand() {
  navigator.clipboard.writeText(generated.value)
  snackbar.value = true
}
</script>

<template>
  <div class="pa-6" style="max-width: 1280px;">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; TCP</p>
      <h1 class="text-h4 font-weight-semibold text-on-surface">TCP access command</h1>
    </header>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <v-card class="mb-6">
      <v-card-title class="d-flex align-center ga-2">
        <v-icon>mdi-console</v-icon>
        Configure
      </v-card-title>
      <v-card-text class="d-flex flex-column ga-4">
        <v-text-field
          v-model="form.hostname"
          label="Hostname"
          placeholder="tcp.example.com"
          variant="outlined"
          density="compact"
          hide-details
        />
        <v-text-field
          v-model="form.local_url"
          label="Local URL"
          placeholder="localhost:9999"
          variant="outlined"
          density="compact"
          hide-details
        />
        <div>
          <p class="text-caption font-weight-medium text-medium-emphasis mb-2">Mode</p>
          <v-chip-group v-model="form.mode" mandatory color="primary" variant="tonal" density="compact">
            <v-chip value="foreground" size="small">foreground</v-chip>
            <v-chip value="nohup" size="small">nohup</v-chip>
            <v-chip value="systemd" size="small">systemd</v-chip>
          </v-chip-group>
        </div>
        <v-btn color="primary" :loading="loading" @click="generate" class="align-self-start">
          {{ loading ? 'Generating...' : 'Generate command' }}
        </v-btn>
      </v-card-text>
    </v-card>

    <v-card v-if="generated">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Output</span>
        <v-btn variant="tonal" size="small" @click="copyCommand">
          <v-icon start size="x-small">mdi-content-copy</v-icon>
          Copy
        </v-btn>
      </v-card-title>
      <pre
        class="font-mono text-body-2 text-on-surface bg-grey-lighten-3 pa-4 overflow-auto"
        style="white-space: pre-wrap; line-height: 1.5;"
      >{{ generated }}</pre>
    </v-card>

    <v-snackbar v-model="snackbar" timeout="1500" color="success">Copied</v-snackbar>
  </div>
</template>
