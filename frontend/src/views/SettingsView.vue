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
  loading.value = true
  message.value = null
  try {
    await apiFetch('/api/settings/password', {
      method: 'POST',
      body: JSON.stringify(form),
    })
    message.value = { type: 'success', text: 'Password updated' }
    form.current_password = ''
    form.new_password = ''
    form.confirm = ''
  } catch (e: any) {
    message.value = { type: 'error', text: e.message }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <header>
      <p class="eyebrow mb-2">Console · Settings</p>
      <h1 class="text-2xl font-semibold tracking-tight text-text">Account</h1>
    </header>

    <section class="card overflow-hidden max-w-[640px]">
      <div class="card-header">
        <span class="card-title">Profile</span>
      </div>
      <div class="card-body grid sm:grid-cols-2 gap-4">
        <div>
          <span class="field-label">Username</span>
          <p class="font-mono text-text">{{ authStore.username || 'admin' }}</p>
        </div>
        <div>
          <span class="field-label">Role</span>
          <p class="text-text">Administrator</p>
        </div>
      </div>
    </section>

    <section class="card overflow-hidden max-w-[640px]">
      <div class="card-header">
        <span class="card-title">Change password</span>
      </div>
      <form @submit.prevent="changePassword" class="card-body space-y-4">
        <div v-if="message" :class="message.type === 'success' ? 'alert-success' : 'alert-danger'">
          {{ message.text }}
        </div>

        <div>
          <label class="field-label">Current password</label>
          <input v-model="form.current_password" type="password" autocomplete="current-password" class="input" required />
        </div>
        <div>
          <label class="field-label">New password</label>
          <input v-model="form.new_password" type="password" autocomplete="new-password" class="input" required minlength="8" />
        </div>
        <div>
          <label class="field-label">Confirm new password</label>
          <input v-model="form.confirm" type="password" autocomplete="new-password" class="input" required minlength="8" />
        </div>

        <button type="submit" :disabled="loading" class="btn-primary">
          {{ loading ? 'Updating…' : 'Update password' }}
        </button>
      </form>
    </section>
  </div>
</template>
