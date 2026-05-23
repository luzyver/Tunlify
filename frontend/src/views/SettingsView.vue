<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const form = reactive({ current_password: '', new_password: '', confirm: '' })
const message = ref<{ type: string; text: string } | null>(null)
const loading = ref(false)

async function changePassword() {
  if (form.new_password !== form.confirm) { message.value = { type: 'error', text: 'Passwords do not match' }; return }
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
    <h1 class="text-2xl font-bold tracking-tight">Settings</h1>

    <div class="card max-w-md">
      <div class="card-header">Change Password</div>
      <div class="card-body space-y-4">
        <div v-if="message" class="p-3 rounded text-sm" :class="message.type === 'success' ? 'bg-success/10 text-emerald-800' : 'bg-danger/10 text-danger'">
          {{ message.text }}
        </div>
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Current Password</label>
          <input v-model="form.current_password" type="password" class="input" />
        </div>
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">New Password</label>
          <input v-model="form.new_password" type="password" class="input" />
        </div>
        <div class="space-y-1">
          <label class="text-xs font-medium text-text-muted">Confirm</label>
          <input v-model="form.confirm" type="password" class="input" />
        </div>
        <button @click="changePassword" :disabled="loading" class="btn-primary">
          {{ loading ? 'Updating…' : 'Update' }}
        </button>
      </div>
    </div>
  </div>
</template>
