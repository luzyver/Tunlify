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
    <h1 class="text-2xl font-bold tracking-tight">TCP Access</h1>

    <div class="card max-w-xl">
      <div class="card-header">Generate Command</div>
      <div class="card-body space-y-4">
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Hostname</label>
          <input v-model="form.hostname" placeholder="tcp.example.com" class="input" />
        </div>
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Local URL</label>
          <input v-model="form.local_url" placeholder="localhost:9999" class="input" />
        </div>
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Mode</label>
          <div class="flex gap-4">
            <label v-for="m in ['foreground', 'nohup', 'systemd']" :key="m" class="flex items-center gap-2 text-sm cursor-pointer" :class="form.mode === m ? 'text-accent font-medium' : 'text-text-muted'">
              <input type="radio" v-model="form.mode" :value="m" class="accent-accent" /> {{ m }}
            </label>
          </div>
        </div>
        <button @click="generate" class="btn-primary">Generate</button>
      </div>
    </div>

    <div v-if="generatedCommand" class="card max-w-xl relative">
      <div class="card-header flex items-center justify-between">
        <span>Output</span>
        <button @click="copy(generatedCommand)" class="btn-secondary !px-2 !py-1 !text-xs">{{ copied ? 'Copied!' : 'Copy' }}</button>
      </div>
      <div class="card-body">
        <pre class="font-mono text-sm text-text bg-surface p-4 rounded overflow-x-auto whitespace-pre-wrap">{{ generatedCommand }}</pre>
      </div>
    </div>
  </div>
</template>
