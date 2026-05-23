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
  try {
    const data = await apiFetch<{ content: string }>('/api/config')
    parseYaml(data.content)
  } catch {}
}

function parseYaml(content: string) {
  const lines = content.split('\n')
  rules.value = []
  for (const line of lines) {
    const m = line.match(/^tunnel:\s*(.+)/)
    if (m) tunnelId.value = m[1].trim()
  }
  let inIngress = false
  let cur: Partial<IngressRule> = {}
  for (const line of lines) {
    if (line.match(/^ingress:/)) { inIngress = true; continue }
    if (!inIngress) continue
    const hm = line.match(/^\s+-\s+hostname:\s*(.+)/)
    const sm = line.match(/^\s+service:\s*(.+)/)
    const tm = line.match(/noTLSVerify:\s*true/)
    if (hm) {
      if (cur.hostname) rules.value.push({ hostname: cur.hostname, service: cur.service || '', noTLSVerify: cur.noTLSVerify || false })
      cur = { hostname: hm[1].trim() }
    } else if (sm && cur.hostname) {
      cur.service = sm[1].trim()
    } else if (tm) {
      cur.noTLSVerify = true
    }
  }
  if (cur.hostname) rules.value.push({ hostname: cur.hostname, service: cur.service || '', noTLSVerify: cur.noTLSVerify || false })
}

function addRule() {
  rules.value.push({ hostname: '', service: '', noTLSVerify: false })
}

function isLocked(rule: IngressRule) {
  return rule.service.includes('tunlify-frontend')
}

function removeRule(i: number) {
  if (isLocked(rules.value[i])) return
  rules.value.splice(i, 1)
}

function generateYaml(): string {
  let y = `tunnel: ${tunnelId.value}\ncredentials-file: /etc/cloudflared/credentials.json\n\ningress:\n`
  for (const r of rules.value) {
    if (!r.hostname || !r.service) continue
    y += `  - hostname: ${r.hostname}\n    service: ${r.service}\n`
    if (r.noTLSVerify) y += `    originRequest:\n      noTLSVerify: true\n`
  }
  return y + `  - service: http_status:404\n`
}

async function save() {
  saving.value = true
  message.value = null
  try {
    await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify({ content: generateYaml() }) })
    message.value = { type: 'success', text: 'Configuration saved' }
  } catch (e: any) {
    message.value = { type: 'error', text: e.message }
  } finally {
    saving.value = false
  }
}

async function saveAndRestart() {
  await save()
  if (message.value?.type !== 'success') return
  try {
    await apiFetch('/api/control/restart', { method: 'POST' })
    message.value = { type: 'success', text: 'Saved & tunnel restarted' }
  } catch (e: any) {
    message.value = { type: 'error', text: `Restart failed: ${e.message}` }
  }
}

loadConfig()
</script>

<template>
  <div class="space-y-6">
    <header class="flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow mb-2">Console · Ingress</p>
        <h1 class="text-2xl font-semibold tracking-tight text-text">Tunnel configuration</h1>
      </div>
      <button class="btn-secondary" @click="addRule">
        <Plus class="w-4 h-4" :stroke-width="1.75" />
        Add rule
      </button>
    </header>

    <div v-if="message" :class="message.type === 'success' ? 'alert-success' : 'alert-danger'">
      {{ message.text }}
    </div>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Tunnel</span>
        <span class="text-2xs text-text-dim">Read from cloudflared/config.yml</span>
      </div>
      <div class="card-body grid sm:grid-cols-2 gap-4">
        <div>
          <label class="field-label">Tunnel ID</label>
          <input v-model="tunnelId" class="input font-mono" placeholder="00000000-0000-0000-0000-000000000000" />
        </div>
        <div>
          <label class="field-label">Credentials</label>
          <input value="/etc/cloudflared/credentials.json" disabled class="input font-mono" />
        </div>
      </div>
    </section>

    <section class="card overflow-hidden">
      <div class="card-header">
        <span class="card-title">Ingress rules</span>
        <span class="text-2xs text-text-dim tabular-nums">{{ rules.length }} rules</span>
      </div>
      <table class="table-tight">
        <thead>
          <tr>
            <th class="w-[34%]">Hostname</th>
            <th class="w-[34%]">Service</th>
            <th>TLS</th>
            <th class="w-10"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(rule, i) in rules" :key="i">
            <td>
              <input
                v-model="rule.hostname"
                placeholder="app.example.com"
                class="input !py-1 !text-xs font-mono"
                :disabled="isLocked(rule)"
              />
            </td>
            <td>
              <input
                v-model="rule.service"
                placeholder="http://container:port"
                class="input !py-1 !text-xs font-mono"
                :disabled="isLocked(rule)"
              />
            </td>
            <td>
              <label class="inline-flex items-center gap-2 text-xs text-text-muted">
                <input type="checkbox" v-model="rule.noTLSVerify" class="accent-accent" :disabled="isLocked(rule)" />
                Skip verify
                <span v-if="isLocked(rule)" class="badge-neutral">Locked</span>
              </label>
            </td>
            <td>
              <button
                v-if="!isLocked(rule)"
                @click="removeRule(i)"
                class="btn-icon-ghost text-danger hover:bg-danger/5"
                title="Remove rule"
              >
                <Trash2 class="w-3.5 h-3.5" :stroke-width="1.75" />
              </button>
            </td>
          </tr>
          <tr v-if="!rules.length">
            <td colspan="4" class="text-center text-text-muted py-8 text-sm">
              No ingress rules yet — add one above.
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <p class="text-xs text-text-dim">
      Catch-all <code class="font-mono text-text-muted">http_status:404</code> is auto-appended on save.
    </p>

    <div class="flex items-center gap-3">
      <button @click="save" :disabled="saving" class="btn-secondary">
        {{ saving ? 'Saving…' : 'Save' }}
      </button>
      <button @click="saveAndRestart" :disabled="saving" class="btn-primary">
        {{ saving ? 'Saving…' : 'Save & restart tunnel' }}
      </button>
    </div>
  </div>
</template>
