<script setup lang="ts">
import { ref } from 'vue'
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
  <div class="page">
    <div class="page-header">
      <div>
        <span class="page-badge">Console &middot; Ingress</span>
        <h1 class="page-title">Tunnel configuration</h1>
      </div>
      <button class="btn btn-primary" @click="addRule">
        <v-icon size="16">mdi-plus</v-icon>
        Add rule
      </button>
    </div>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <div class="card mb-6">
      <div class="card-header">
        <span class="card-title">Tunnel</span>
        <span class="stat-label" style="margin-bottom: 0;">Read from cloudflared/config.yml</span>
      </div>
      <div class="card-body">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
          <div>
            <label class="label-upper">Tunnel ID</label>
            <input v-model="tunnelId" placeholder="00000000-0000-0000-0000-000000000000" class="input-line" />
          </div>
          <div>
            <label class="label-upper">Credentials</label>
            <input value="/etc/cloudflared/credentials.json" disabled class="input-line" />
          </div>
        </div>
      </div>
    </div>

    <div class="card mb-6">
      <div class="card-header">
        <span class="card-title">Ingress rules</span>
        <span class="stat-label" style="margin-bottom: 0; color: #71717A;">{{ rules.length }} rules</span>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width: 34%;">Hostname</th>
              <th style="width: 34%;">Service</th>
              <th>TLS</th>
              <th style="width: 40px;"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(rule, i) in rules" :key="i">
              <td>
                <input v-model="rule.hostname" placeholder="app.example.com" :disabled="isLocked(rule)" class="input-line" style="border-bottom-width: 1px; background: rgba(0,0,0,0.3);" />
              </td>
              <td>
                <input v-model="rule.service" placeholder="http://container:port" :disabled="isLocked(rule)" class="input-line" style="border-bottom-width: 1px; background: rgba(0,0,0,0.3);" />
              </td>
              <td>
                <div class="d-flex align-center ga-2">
                  <label class="d-flex align-center ga-2 font-mono" style="color: #94A3B8; font-size: 12px; cursor: pointer;">
                    <input type="checkbox" v-model="rule.noTLSVerify" :disabled="isLocked(rule)" style="accent-color: #F7931A;" />
                    Skip verify
                  </label>
                  <span v-if="isLocked(rule)" class="rounded-pill px-2 font-mono" style="background: rgba(148, 163, 184, 0.1); border: 1px solid rgba(148, 163, 184, 0.2); color: #94A3B8; font-size: 10px; line-height: 20px;">Locked</span>
                </div>
              </td>
              <td>
                <button v-if="!isLocked(rule)" class="btn-ghost btn-ghost--danger" @click="removeRule(i)">
                  <v-icon size="14">mdi-delete</v-icon>
                </button>
              </td>
            </tr>
            <tr v-if="!rules.length">
              <td colspan="4" class="empty">No ingress rules yet &mdash; add one above.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <p class="stat-label mb-4" style="color: #71717A;">
      Catch-all <code class="font-mono" style="color: #F7931A;">http_status:404</code> is auto-appended on save.
    </p>

    <div class="d-flex ga-3">
      <button class="btn btn-secondary" :disabled="saving" @click="save">
        {{ saving ? 'Saving...' : 'Save' }}
      </button>
      <button class="btn btn-primary" :disabled="saving" @click="saveAndRestart">
        {{ saving ? 'Saving...' : 'Save & restart tunnel' }}
      </button>
    </div>
  </div>
</template>
