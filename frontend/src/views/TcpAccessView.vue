<script setup lang="ts">
import { reactive, ref } from 'vue'
import { Check, Copy, Terminal } from 'lucide-vue-next'
import { useClipboard } from '@vueuse/core'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const { copy, copied } = useClipboard({ copiedDuring: 1500 })

const form = reactive({ hostname: '', local_url: 'localhost:', mode: 'foreground' })
const generated = ref('')
const error = ref('')
const loading = ref(false)

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
</script>

<template>
  <div class="space-y-6">
    <header>
      <p class="eyebrow mb-2">Console · TCP</p>
      <h1 class="text-2xl font-semibold tracking-tight text-text">TCP access command</h1>
    </header>

    <div v-if="error" class="alert-danger">{{ error }}</div>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title flex items-center gap-2">
          <Terminal class="w-4 h-4 text-text-muted" :stroke-width="1.75" />
          Configure
        </span>
      </div>
      <div class="card-body space-y-4">
        <div>
          <label class="field-label">Hostname</label>
          <input v-model="form.hostname" placeholder="tcp.example.com" class="input font-mono" />
        </div>
        <div>
          <label class="field-label">Local URL</label>
          <input v-model="form.local_url" placeholder="localhost:9999" class="input font-mono" />
        </div>
        <div>
          <label class="field-label">Mode</label>
          <div class="inline-flex border border-border rounded-md overflow-hidden bg-surface">
            <button
              v-for="m in ['foreground', 'nohup', 'systemd']"
              :key="m"
              type="button"
              class="px-3 h-8 text-xs font-medium transition-colors duration-100"
              :class="form.mode === m ? 'bg-accent-soft text-accent' : 'text-text-muted hover:bg-bg-alt'"
              @click="form.mode = m"
            >
              {{ m }}
            </button>
          </div>
        </div>
        <button @click="generate" :disabled="loading" class="btn-primary">
          {{ loading ? 'Generating…' : 'Generate command' }}
        </button>
      </div>
    </section>

    <section v-if="generated" class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Output</span>
        <button @click="copy(generated)" class="btn-secondary !py-1 !px-2.5 !text-xs">
          <Check v-if="copied" class="w-3.5 h-3.5 text-success" :stroke-width="1.75" />
          <Copy v-else class="w-3.5 h-3.5" :stroke-width="1.75" />
          {{ copied ? 'Copied' : 'Copy' }}
        </button>
      </div>
      <pre
        class="font-mono text-xs leading-5 text-text bg-bg-alt
               p-4 whitespace-pre-wrap overflow-x-auto"
      >{{ generated }}</pre>
    </section>
  </div>
</template>
