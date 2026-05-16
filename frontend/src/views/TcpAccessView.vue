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
    <h1 class="cyber-heading text-xl lg:text-2xl">TCP Access</h1>

    <div class="cyber-card space-y-4 max-w-xl">
      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">--hostname</label>
        <input v-model="form.hostname" placeholder="tcp.example.com" class="cyber-input" />
      </div>
      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">--url</label>
        <input v-model="form.local_url" placeholder="localhost:9999" class="cyber-input" />
      </div>
      <div class="flex gap-4">
        <label v-for="m in ['foreground', 'nohup', 'systemd']" :key="m" class="flex items-center gap-2 text-xs uppercase tracking-wider cursor-pointer" :class="form.mode === m ? 'text-neon' : 'text-muted'">
          <input type="radio" v-model="form.mode" :value="m" class="accent-neon" />
          {{ m }}
        </label>
      </div>
      <button @click="generate" class="cyber-btn">Generate</button>
    </div>

    <div v-if="generatedCommand" class="cyber-card relative max-w-xl">
      <p class="text-[10px] text-muted uppercase tracking-widest mb-2">// Output</p>
      <pre class="font-mono text-sm text-neon bg-void p-4 border border-border overflow-x-auto whitespace-pre-wrap">{{ generatedCommand }}</pre>
      <button @click="copy(generatedCommand)" class="cyber-btn-secondary !px-2 !py-1 !text-[10px] absolute top-4 right-4">
        {{ copied ? 'COPIED' : 'COPY' }}
      </button>
    </div>
  </div>
</template>
