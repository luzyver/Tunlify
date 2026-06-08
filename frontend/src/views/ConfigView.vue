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
  <div class="pa-6" style="max-width: 1280px;">
    <div class="d-flex align-end justify-space-between ga-4 mb-6">
      <div>
        <p class="eyebrow mb-2">Console &middot; Ingress</p>
        <h1 class="text-h4 font-weight-semibold text-on-surface">Tunnel configuration</h1>
      </div>
      <v-btn color="primary" @click="addRule">
        <v-icon start>mdi-plus</v-icon>
        Add rule
      </v-btn>
    </div>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <v-card class="mb-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Tunnel</span>
        <span class="text-caption text-medium-emphasis">Read from cloudflared/config.yml</span>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field v-model="tunnelId" label="Tunnel ID" placeholder="00000000-0000-0000-0000-000000000000" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field model-value="/etc/cloudflared/credentials.json" label="Credentials" disabled variant="outlined" density="compact" hide-details />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-card class="mb-6">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>Ingress rules</span>
        <span class="text-caption text-medium-emphasis tabular-nums">{{ rules.length }} rules</span>
      </v-card-title>
      <v-table density="compact">
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
              <v-text-field
                v-model="rule.hostname"
                placeholder="app.example.com"
                density="compact"
                variant="outlined"
                hide-details
                :disabled="isLocked(rule)"
              />
            </td>
            <td>
              <v-text-field
                v-model="rule.service"
                placeholder="http://container:port"
                density="compact"
                variant="outlined"
                hide-details
                :disabled="isLocked(rule)"
              />
            </td>
            <td>
              <div class="d-flex align-center ga-2">
                <v-checkbox
                  v-model="rule.noTLSVerify"
                  label="Skip verify"
                  hide-details
                  density="compact"
                  :disabled="isLocked(rule)"
                />
                <v-chip v-if="isLocked(rule)" size="x-small" variant="outlined">Locked</v-chip>
              </div>
            </td>
            <td>
              <v-btn
                v-if="!isLocked(rule)"
                icon="mdi-delete"
                color="error"
                variant="text"
                size="small"
                @click="removeRule(i)"
              />
            </td>
          </tr>
          <tr v-if="!rules.length">
            <td colspan="4" class="text-center text-medium-emphasis py-8">
              No ingress rules yet &mdash; add one above.
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <p class="text-caption text-medium-emphasis mb-4">
      Catch-all <code class="font-mono">http_status:404</code> is auto-appended on save.
    </p>

    <div class="d-flex ga-3">
      <v-btn variant="tonal" :loading="saving" @click="save">
        {{ saving ? 'Saving...' : 'Save' }}
      </v-btn>
      <v-btn color="primary" :loading="saving" @click="saveAndRestart">
        {{ saving ? 'Saving...' : 'Save & restart tunnel' }}
      </v-btn>
    </div>
  </div>
</template>
