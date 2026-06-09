<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useApi } from '../composables/useApi'


const { apiFetch } = useApi()

const form = reactive({ hostname: '', local_url: 'localhost:', mode: 'foreground' })
const generated = ref('')
const error = ref('')
const loading = ref(false)
const copied = ref(false)

async function generate() {
  loading.value = true
  error.value = ''
  try {
    const data = await apiFetch<{ command: string }>('/api/tcp-access/generate', {
      method: 'POST', body: JSON.stringify(form),
    })
    generated.value = data.command
  } catch (e: any) { error.value = e.message }
  finally { loading.value = false }
}

function copyCommand() {
  navigator.clipboard.writeText(generated.value)
  copied.value = true
  setTimeout(() => copied.value = false, 1500)
}
</script>

<style scoped>
.mode-btn {
  height: 32px;
  background: transparent;
  border: none;
  color: #94A3B8;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}
.mode-btn:hover {
  background: rgba(255,255,255,0.05);
}
.mode-btn--active {
  background: rgba(247, 147, 26, 0.15) !important;
  color: #F7931A !important;
}
</style>

<template>
  <div class="page">
    <header class="mb-6">
      <span class="page-badge">Console &middot; TCP</span>
      <h1 class="page-title" style="margin-top: 4px;">TCP access command</h1>
    </header>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <div class="card mb-6">
      <div class="card-header">
        <div class="d-flex align-center ga-2">
          <v-icon size="18" color="#F7931A">mdi-console</v-icon>
          <span class="card-title">Configure</span>
        </div>
      </div>
      <div class="card-body d-flex flex-column ga-5">
        <div>
          <label class="label-upper">Hostname</label>
          <input v-model="form.hostname" placeholder="tcp.example.com" class="input-line" />
        </div>
        <div>
          <label class="label-upper">Local URL</label>
          <input v-model="form.local_url" placeholder="localhost:9999" class="input-line" />
        </div>
        <div>
          <p class="label-upper">Mode</p>
          <div class="d-inline-flex rounded-lg overflow-hidden" style="border: 1px solid rgba(30, 41, 59, 0.6);">
            <button
              v-for="m in ['foreground', 'nohup', 'systemd']"
              :key="m"
              class="px-3 font-mono mode-btn"
              :class="{ 'mode-btn--active': form.mode === m }"
              @click="form.mode = m"
            >{{ m }}</button>
          </div>
        </div>
        <button class="btn btn-primary align-self-start" :disabled="loading" @click="generate">
          {{ loading ? 'Generating...' : 'Generate command' }}
        </button>
      </div>
    </div>

    <div v-if="generated" class="card">
      <div class="card-header">
        <span class="card-title">Output</span>
        <button class="btn btn-secondary btn-sm" @click="copyCommand">
          <v-icon v-if="copied" size="12" color="#FFD600">mdi-check</v-icon>
          <v-icon v-else size="12">mdi-content-copy</v-icon>
          {{ copied ? 'Copied' : 'Copy' }}
        </button>
      </div>
      <pre class="font-mono pa-4 overflow-auto" style="color: white; background: rgba(0,0,0,0.3); font-size: 13px; line-height: 1.5; white-space: pre-wrap; max-height: 400px;">{{ generated }}</pre>
    </div>
  </div>
</template>
