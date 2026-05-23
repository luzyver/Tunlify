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
  loading.value = true; message.value = null
  try {
    await apiFetch('/api/settings/password', { method: 'POST', body: JSON.stringify(form) })
    message.value = { type: 'success', text: 'Password updated' }
    form.current_password = ''; form.new_password = ''; form.confirm = ''
  } catch (e: any) { message.value = { type: 'error', text: e.message } }
  finally { loading.value = false }
}
</script>

<template>
  <div class="space-y-6">
    <header class="border-b border-border pb-6">
      <p class="section-marker mb-2">ACCOUNT</p>
      <h1 class="editorial-h1 !text-3xl">Settings</h1>
    </header>

    <div class="card max-w-md space-y-4">
      <p class="text-sm font-medium text-text">Change Password</p>

      <div v-if="message" class="p-3 rounded-md text-sm border" :class="message.type === 'success' ? 'border-emerald-200 bg-emerald-50 text-emerald-700' : 'border-red-200 bg-red-50 text-red-700'">
        {{ message.text }}
      </div>

      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Current Password</label>
        <input v-model="form.current_password" type="password" class="input" />
      </div>
      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">New Password</label>
        <input v-model="form.new_password" type="password" class="input" />
      </div>
      <div class="space-y-1.5">
        <label class="text-xs text-text-muted font-medium">Confirm</label>
        <input v-model="form.confirm" type="password" class="input" />
      </div>
      <button @click="changePassword" :disabled="loading" class="btn-primary">
        {{ loading ? 'Updating…' : 'Update Password' }}
      </button>
    </div>
  </div>
</template>
