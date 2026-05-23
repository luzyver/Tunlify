<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useClipboard } from '@vueuse/core'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const { copy, copied } = useClipboard()

const form = reactive({ hostname: '', local_url: 'localhost:', mode: 'foreground' })
const generatedCommand = ref('')

async function generate() {
  const data = await apiFetch<{ command: string }>('/api/tcp-access/generate', {
    method: 'POST',
    body: JSON.stringify(form),
  })
  generatedCommand.value = data.command
}
</script>

<template>
  <div class="space-y-6">
    <header class="border-b border-border pb-6">
      <p class="section-marker mb-2">TOOLS</p>
      <h1 class="editorial-h1 !text-3xl">TCP Access</h1>
    </header>

    <div class="card space-y-4 max-w-xl">
      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Hostname</label>
        <input v-model="form.hostname" placeholder="tcp.example.com" class="input" />
      </div>
      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Local URL</label>
        <input v-model="form.local_url" placeholder="localhost:9999" class="input" />
      </div>
      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Mode</label>
        <div class="flex gap-4">
          <label v-for="m in ['foreground', 'nohup', 'systemd']" :key="m" class="flex items-center gap-2 text-sm cursor-pointer" :class="form.mode === m ? 'text-text font-medium' : 'text-text-muted'">
            <input type="radio" v-model="form.mode" :value="m" class="accent-accent" /> {{ m }}
          </label>
        </div>
      </div>
      <button @click="generate" class="btn-primary">Generate Command</button>
    </div>

    <div v-if="generatedCommand" class="card relative max-w-xl">
      <p class="text-xs text-text-muted font-medium mb-2">Output</p>
      <pre class="font-mono text-sm text-text bg-surface p-4 rounded-md overflow-x-auto whitespace-pre-wrap">{{ generatedCommand }}</pre>
      <button @click="copy(generatedCommand)" class="btn-secondary !px-2 !py-1 !text-xs absolute top-4 right-4">
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
  </div>
</template>
