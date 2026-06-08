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
  <div class="pa-6" style="max-width: 1280px;">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; Settings</p>
      <h1 class="text-h4 font-weight-semibold text-on-surface">Account</h1>
    </header>

    <v-card class="mb-6">
      <v-card-title>Profile</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <p class="text-caption font-weight-medium text-medium-emphasis mb-1">Username</p>
            <p class="font-mono text-on-surface">{{ authStore.username || 'admin' }}</p>
          </v-col>
          <v-col cols="12" sm="6">
            <p class="text-caption font-weight-medium text-medium-emphasis mb-1">Role</p>
            <p class="text-on-surface">Administrator</p>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-card>
      <v-card-title>Change password</v-card-title>
      <v-form @submit.prevent="changePassword">
        <v-card-text class="d-flex flex-column ga-4">
          <v-alert v-if="message" :type="message.type" variant="tonal">{{ message.text }}</v-alert>

          <v-text-field
            v-model="form.current_password"
            label="Current password"
            type="password"
            autocomplete="current-password"
            variant="outlined"
            density="compact"
            hide-details
            required
          />

          <v-text-field
            v-model="form.new_password"
            label="New password"
            type="password"
            autocomplete="new-password"
            variant="outlined"
            density="compact"
            hide-details
            required
            minlength="8"
          />

          <v-text-field
            v-model="form.confirm"
            label="Confirm new password"
            type="password"
            autocomplete="new-password"
            variant="outlined"
            density="compact"
            hide-details
            required
            minlength="8"
          />

          <v-btn type="submit" color="primary" :loading="loading" class="align-self-start">
            {{ loading ? 'Updating...' : 'Update password' }}
          </v-btn>
        </v-card-text>
      </v-form>
    </v-card>
  </div>
</template>
