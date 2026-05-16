<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const form = reactive({ current_password: '', new_password: '', confirm: '' })
const message = ref<{ type: string; text: string } | null>(null)
const loading = ref(false)

async function changePassword() {
  if (form.new_password !== form.confirm) {
    message.value = { type: 'error', text: 'Passwords do not match' }
    return
  }
  loading.value = true
  message.value = null
  try {
    await apiFetch('/api/settings/password', { method: 'POST', body: JSON.stringify(form) })
    message.value = { type: 'success', text: 'PASSWORD UPDATED' }
    form.current_password = ''; form.new_password = ''; form.confirm = ''
  } catch (e: any) {
    message.value = { type: 'error', text: e.message }
  } finally { loading.value = false }
}
</script>

<template>
  <div class="space-y-6">
    <h1 class="cyber-heading text-xl lg:text-2xl">Settings</h1>

    <div class="cyber-card max-w-md space-y-4">
      <h2 class="text-xs text-muted uppercase tracking-widest">// Change Password</h2>

      <div v-if="message" class="border p-3 text-sm font-mono" :class="message.type === 'success' ? 'border-neon text-neon bg-neon/5' : 'border-danger text-danger bg-danger/5'">
        [{{ message.type === 'success' ? 'OK' : 'ERR' }}] {{ message.text }}
      </div>

      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">Current Password</label>
        <input v-model="form.current_password" type="password" class="cyber-input" />
      </div>
      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">New Password</label>
        <input v-model="form.new_password" type="password" class="cyber-input" />
      </div>
      <div class="space-y-1">
        <label class="text-[10px] text-muted uppercase tracking-widest">Confirm</label>
        <input v-model="form.confirm" type="password" class="cyber-input" />
      </div>
      <button @click="changePassword" :disabled="loading" class="cyber-btn">
        {{ loading ? 'UPDATING...' : 'UPDATE' }}
      </button>
    </div>
  </div>
</template>
