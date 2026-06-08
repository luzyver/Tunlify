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

<template>
  <div class="pa-6">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; TCP</p>
      <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">TCP access command</h1>
    </header>

    <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

    <div class="rounded-2xl mb-6" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center ga-3 px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <div class="rounded-lg p-2 d-flex" style="background: rgba(234, 88, 12, 0.15); border: 1px solid rgba(234, 88, 12, 0.3);">
          <v-icon size="18" color="#F7931A">mdi-console</v-icon>
        </div>
        <span class="font-heading font-semibold" style="color: white;">Configure</span>
      </div>
      <div class="pa-6 d-flex flex-column ga-5">
        <div>
          <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Hostname</label>
          <input v-model="form.hostname" placeholder="tcp.example.com" style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'; $event.target.style.boxShadow = '0 10px 20px -10px rgba(247, 147, 26, 0.3)'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'; $event.target.style.boxShadow = 'none'" />
        </div>
        <div>
          <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Local URL</label>
          <input v-model="form.local_url" placeholder="localhost:9999" style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
        </div>
        <div>
          <p class="font-mono mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Mode</p>
          <div class="d-inline-flex rounded-lg overflow-hidden" style="border: 1px solid rgba(30, 41, 59, 0.6);">
            <button v-for="m in ['foreground', 'nohup', 'systemd']" :key="m"
              class="px-3 font-mono"
              :style="{ height: '32px', background: form.mode === m ? 'rgba(247, 147, 26, 0.15)' : 'transparent', border: 'none', color: form.mode === m ? '#F7931A' : '#94A3B8', fontSize: '11px', cursor: 'pointer', transition: 'all 0.2s' }"
              @click="form.mode = m"
              @mouseenter="$event.target.style.background = form.mode === m ? 'rgba(247, 147, 26, 0.15)' : 'rgba(255,255,255,0.05)'"
              @mouseleave="$event.target.style.background = form.mode === m ? 'rgba(247, 147, 26, 0.15)' : 'transparent'"
            >{{ m }}</button>
          </div>
        </div>
        <button
          class="align-self-start d-inline-flex align-center ga-2 rounded-pill px-5 font-mono"
          style="height: 40px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
          :disabled="loading"
          @click="generate"
          @mouseenter="if(!loading) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)' }"
          @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
        >
          {{ loading ? 'Generating...' : 'Generate command' }}
        </button>
      </div>
    </div>

    <div v-if="generated" class="rounded-2xl" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Output</span>
        <button
          class="d-inline-flex align-center ga-1 rounded-pill px-3 font-mono"
          style="height: 28px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 11px; cursor: pointer; transition: all 0.2s;"
          @click="copyCommand"
          @mouseenter="$event.target.style.background = 'rgba(247, 147, 26, 0.2)'"
          @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'"
        >
          <v-icon v-if="copied" size="12" color="#FFD600">mdi-check</v-icon>
          <v-icon v-else size="12">mdi-content-copy</v-icon>
          {{ copied ? 'Copied' : 'Copy' }}
        </button>
      </div>
      <pre class="font-mono pa-4 overflow-auto" style="color: white; background: rgba(0,0,0,0.3); font-size: 13px; line-height: 1.5; white-space: pre-wrap; max-height: 400px;">{{ generated }}</pre>
    </div>
  </div>
</template>
