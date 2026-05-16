<script setup lang="ts">
import { ref } from 'vue'
import { Plus, Trash2 } from 'lucide-vue-next'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()

interface IngressRule { hostname: string; service: string; noTLSVerify: boolean }

const rules = ref<IngressRule[]>([])
const tunnelId = ref('')
const saving = ref(false)
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)

async function loadConfig() {
  try { const data = await apiFetch<{ content: string }>('/api/config'); parseYaml(data.content) } catch {}
}

function parseYaml(content: string) {
  const lines = content.split('\n')
  rules.value = []
  for (const line of lines) { const m = line.match(/^tunnel:\s*(.+)/); if (m) tunnelId.value = m[1].trim() }
  let inIngress = false, cur: Partial<IngressRule> = {}
  for (const line of lines) {
    if (line.match(/^ingress:/)) { inIngress = true; continue }
    if (!inIngress) continue
    const hm = line.match(/^\s+-\s+hostname:\s*(.+)/), sm = line.match(/^\s+service:\s*(.+)/), tm = line.match(/noTLSVerify:\s*true/)
    if (hm) { if (cur.hostname) rules.value.push({ hostname: cur.hostname, service: cur.service || '', noTLSVerify: cur.noTLSVerify || false }); cur = { hostname: hm[1].trim() } }
    else if (sm && cur.hostname) cur.service = sm[1].trim()
    else if (tm) cur.noTLSVerify = true
  }
  if (cur.hostname) rules.value.push({ hostname: cur.hostname, service: cur.service || '', noTLSVerify: cur.noTLSVerify || false })
}

function addRule() { rules.value.push({ hostname: '', service: '', noTLSVerify: false }) }
function removeRule(i: number) { if (isTunlifyRule(rules.value[i])) return; rules.value.splice(i, 1) }
function isTunlifyRule(rule: IngressRule) { return rule.service.includes('tunlify-frontend') }

function generateYaml(): string {
  let y = `tunnel: ${tunnelId.value}\ncredentials-file: /etc/cloudflared/credentials.json\n\ningress:\n`
  for (const r of rules.value) { if (!r.hostname || !r.service) continue; y += `  - hostname: ${r.hostname}\n    service: ${r.service}\n`; if (r.noTLSVerify) y += `    originRequest:\n      noTLSVerify: true\n` }
  return y + `  - service: http_status:404\n`
}

async function save() {
  saving.value = true; message.value = null
  try { await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify({ content: generateYaml() }) }); message.value = { type: 'success', text: 'CONFIG WRITTEN' } }
  catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { saving.value = false }
}

async function saveAndRestart() {
  await save()
  if (message.value?.type === 'success') {
    try { await apiFetch('/api/control/restart', { method: 'POST' }); message.value = { type: 'success', text: 'CONFIG SAVED & TUNNEL RESTARTED' } }
    catch (e: any) { message.value = { type: 'error', text: `Restart failed: ${e.message}` } }
  }
}

loadConfig()
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <h1 class="cyber-heading text-xl lg:text-2xl">Ingress Config</h1>
      <button @click="addRule" class="cyber-btn flex items-center gap-2 w-full sm:w-auto justify-center">
        <Plus class="w-4 h-4" :stroke-width="1.5" /> Add Rule
      </button>
    </div>

    <div v-if="message" class="border p-3 text-sm font-mono" :class="message.type === 'success' ? 'border-neon text-neon bg-neon/5' : 'border-danger text-danger bg-danger/5'">
      [{{ message.type === 'success' ? 'OK' : 'ERR' }}] {{ message.text }}
    </div>

    <div class="space-y-3">
      <div v-for="(rule, i) in rules" :key="i" class="cyber-card flex gap-4 items-start">
        <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3">
          <div>
            <label class="text-[10px] text-muted uppercase tracking-widest block mb-1">Hostname</label>
            <input v-model="rule.hostname" placeholder="app.example.com" class="cyber-input" :disabled="isTunlifyRule(rule)" :class="isTunlifyRule(rule) && 'opacity-50 cursor-not-allowed'" />
          </div>
          <div>
            <label class="text-[10px] text-muted uppercase tracking-widest block mb-1">Service</label>
            <input v-model="rule.service" placeholder="http://container:port" class="cyber-input" :disabled="isTunlifyRule(rule)" :class="isTunlifyRule(rule) && 'opacity-50 cursor-not-allowed'" />
          </div>
          <label class="flex items-center gap-2 text-xs text-muted col-span-full">
            <input type="checkbox" v-model="rule.noTLSVerify" class="accent-neon" :disabled="isTunlifyRule(rule)" />
            <span class="uppercase tracking-wider">no-tls-verify</span>
            <span v-if="isTunlifyRule(rule)" class="cyber-badge border-magenta text-magenta ml-2">LOCKED</span>
          </label>
        </div>
        <button v-if="!isTunlifyRule(rule)" @click="removeRule(i)" class="cyber-btn-danger !px-2 !py-2 mt-5">
          <Trash2 class="w-4 h-4" :stroke-width="1.5" />
        </button>
      </div>

      <div v-if="!rules.length" class="cyber-card text-center text-muted py-8 tracking-wider text-sm">
        // NO RULES CONFIGURED
      </div>
    </div>

    <p class="text-xs text-muted tracking-wider border-t border-border pt-4">
      // catch-all: <span class="text-cyan">http_status:404</span> (auto-appended)
    </p>

    <div class="flex flex-col sm:flex-row gap-3">
      <button @click="save" :disabled="saving" class="cyber-btn-secondary w-full sm:w-auto">Save</button>
      <button @click="saveAndRestart" :disabled="saving" class="cyber-btn w-full sm:w-auto">
        {{ saving ? 'WRITING...' : 'SAVE & RESTART' }}
      </button>
    </div>
  </div>
</template>
