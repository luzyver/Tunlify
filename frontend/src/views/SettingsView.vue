<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../composables/useApi'

const { apiFetch } = useApi()
const authStore = useAuthStore()
const form = reactive({ current_password: '', new_password: '', confirm: '' })
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const loading = ref(false)

async function changePassword() {
  if (form.new_password !== form.confirm) {
    message.value = { type: 'error', text: 'New passwords do not match' }
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
  <div class="page">
    <header class="mb-4">
      <span class="page-badge">Console &middot; Settings</span>
      <h1 class="page-title" style="margin-top: 4px;">Account</h1>
    </header>

    <div class="card mb-4">
      <div class="card-header">
        <span class="card-title">Profile</span>
      </div>
      <div class="card-body">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
          <div>
            <p class="label-upper">Username</p>
            <p class="font-mono" style="color: white;">{{ authStore.username || 'admin' }}</p>
          </div>
          <div>
            <p class="label-upper">Role</p>
            <p class="font-mono" style="color: white;">Administrator</p>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <span class="card-title">Change password</span>
      </div>
      <form @submit.prevent="changePassword">
        <div class="card-body d-flex flex-column ga-5">
          <v-alert v-if="message" :type="message.type" variant="tonal">{{ message.text }}</v-alert>

          <div>
            <label class="label-upper">Current password</label>
            <input v-model="form.current_password" type="password" autocomplete="current-password" required class="input-line" />
          </div>

          <div>
            <label class="label-upper">New password</label>
            <input v-model="form.new_password" type="password" autocomplete="new-password" required minlength="8" class="input-line" />
          </div>

          <div>
            <label class="label-upper">Confirm new password</label>
            <input v-model="form.confirm" type="password" autocomplete="new-password" required minlength="8" class="input-line" />
          </div>

          <button type="submit" class="btn btn-primary align-self-start" :disabled="loading">
            {{ loading ? 'Updating...' : 'Update password' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
