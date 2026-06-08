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
  <div class="pa-6">
    <div class="d-flex align-center justify-space-between ga-4 mb-6 flex-wrap" style="gap: 16px;">
      <div>
        <p class="eyebrow mb-2">Console &middot; Ingress</p>
        <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Tunnel configuration</h1>
      </div>
      <button
        class="d-inline-flex align-center ga-2 rounded-pill px-4 font-mono"
        style="height: 40px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
        @click="addRule"
        @mouseenter="$event.target.style.transform = 'scale(1.03)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)'"
        @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
      >
        <v-icon size="16">mdi-plus</v-icon>
        Add rule
      </button>
    </div>

    <v-alert v-if="message" :type="message.type" class="mb-4" variant="tonal">{{ message.text }}</v-alert>

    <div class="rounded-2xl mb-6" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Tunnel</span>
        <span class="font-mono" style="color: #94A3B8; font-size: 11px;">Read from cloudflared/config.yml</span>
      </div>
      <div class="p-6">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
          <div>
            <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Tunnel ID</label>
            <input v-model="tunnelId" placeholder="00000000-0000-0000-0000-000000000000" class="w-100 font-mono" style="background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
          </div>
          <div>
            <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Credentials</label>
            <input value="/etc/cloudflared/credentials.json" disabled class="w-100 font-mono" style="background: rgba(0,0,0,0.3); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.4); color: rgba(255,255,255,0.5); padding: 8px 12px; font-size: 13px; outline: none;" />
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-2xl mb-6" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="d-flex align-center justify-space-between px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Ingress rules</span>
        <span class="font-mono tabular-nums" style="color: #94A3B8; font-size: 11px;">{{ rules.length }} rules</span>
      </div>
      <div style="overflow-x: auto;">
        <table style="width: 100%; border-collapse: collapse;">
          <thead>
            <tr>
              <th class="eyebrow text-left px-4" style="height: 40px; border-bottom: 1px solid rgba(30, 41, 59, 0.6); width: 34%;">Hostname</th>
              <th class="eyebrow text-left px-4" style="height: 40px; border-bottom: 1px solid rgba(30, 41, 59, 0.6); width: 34%;">Service</th>
              <th class="eyebrow text-left px-4" style="height: 40px; border-bottom: 1px solid rgba(30, 41, 59, 0.6);">TLS</th>
              <th class="px-4" style="height: 40px; border-bottom: 1px solid rgba(30, 41, 59, 0.6); width: 40px;"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(rule, i) in rules" :key="i">
              <td class="px-2" style="border-bottom: 1px solid rgba(30, 41, 59, 0.3);">
                <input v-model="rule.hostname" placeholder="app.example.com" :disabled="isLocked(rule)" style="width: 100%; background: rgba(0,0,0,0.3); border: none; border-bottom: 1px solid rgba(30, 41, 59, 0.6); color: white; padding: 8px 10px; font-family: 'JetBrains Mono', monospace; font-size: 12px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.6)'" />
              </td>
              <td class="px-2" style="border-bottom: 1px solid rgba(30, 41, 59, 0.3);">
                <input v-model="rule.service" placeholder="http://container:port" :disabled="isLocked(rule)" style="width: 100%; background: rgba(0,0,0,0.3); border: none; border-bottom: 1px solid rgba(30, 41, 59, 0.6); color: white; padding: 8px 10px; font-family: 'JetBrains Mono', monospace; font-size: 12px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.6)'" />
              </td>
              <td style="border-bottom: 1px solid rgba(30, 41, 59, 0.3);">
                <div class="d-flex align-center ga-2">
                  <label class="d-flex align-center ga-2 font-mono" style="color: #94A3B8; font-size: 12px; cursor: pointer;">
                    <input type="checkbox" v-model="rule.noTLSVerify" :disabled="isLocked(rule)" style="accent-color: #F7931A;" />
                    Skip verify
                  </label>
                  <span v-if="isLocked(rule)" class="rounded-pill px-2 font-mono" style="background: rgba(148, 163, 184, 0.1); border: 1px solid rgba(148, 163, 184, 0.2); color: #94A3B8; font-size: 10px; line-height: 20px;">Locked</span>
                </div>
              </td>
              <td style="border-bottom: 1px solid rgba(30, 41, 59, 0.3);">
                <button v-if="!isLocked(rule)" @click="removeRule(i)" style="background: none; border: none; color: #EF4444; cursor: pointer; padding: 4px; transition: color 0.2s;" @mouseenter="$event.target.style.color = '#F7931A'" @mouseleave="$event.target.style.color = '#EF4444'">
                  <v-icon size="14">mdi-delete</v-icon>
                </button>
              </td>
            </tr>
            <tr v-if="!rules.length">
              <td colspan="4" class="text-center font-mono py-8" style="color: #94A3B8;">No ingress rules yet &mdash; add one above.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <p class="font-mono mb-4" style="color: #94A3B8; font-size: 11px;">
      Catch-all <code class="font-mono" style="color: #F7931A;">http_status:404</code> is auto-appended on save.
    </p>

    <div class="d-flex ga-3">
      <button class="rounded-pill px-5 font-mono" style="height: 40px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-size: 12px; font-weight: 600; cursor: pointer; transition: all 0.3s;" :disabled="saving" @click="save" @mouseenter="if(!saving) { $event.target.style.background = 'rgba(247, 147, 26, 0.2)'; $event.target.style.borderColor = '#F7931A' }" @mouseleave="$event.target.style.background = 'rgba(247, 147, 26, 0.1)'; $event.target.style.borderColor = 'rgba(247, 147, 26, 0.2)'">
        {{ saving ? 'Saving...' : 'Save' }}
      </button>
      <button
        class="rounded-pill px-5 font-mono"
        style="height: 40px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
        :disabled="saving"
        @click="saveAndRestart"
        @mouseenter="if(!saving) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)' }"
        @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
      >
        {{ saving ? 'Saving...' : 'Save & restart tunnel' }}
      </button>
    </div>
  </div>
</template>
