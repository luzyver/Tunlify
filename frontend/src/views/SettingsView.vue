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
  <div class="pa-6" style="max-width: 1280px;">
    <header class="mb-6">
      <p class="eyebrow mb-2">Console &middot; Settings</p>
      <h1 class="font-heading" style="font-size: 28px; font-weight: 600; color: white;">Account</h1>
    </header>

    <div class="rounded-2xl mb-6" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Profile</span>
      </div>
      <div class="pa-6">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
          <div>
            <p class="font-mono mb-1" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Username</p>
            <p class="font-mono" style="color: white;">{{ authStore.username || 'admin' }}</p>
          </div>
          <div>
            <p class="font-mono mb-1" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Role</p>
            <p class="font-mono" style="color: white;">Administrator</p>
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-2xl" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <div class="px-6 py-4" style="border-bottom: 1px solid rgba(30, 41, 59, 0.6);">
        <span class="font-heading font-semibold" style="color: white;">Change password</span>
      </div>
      <form @submit.prevent="changePassword">
        <div class="pa-6 d-flex flex-column ga-5">
          <v-alert v-if="message" :type="message.type" variant="tonal">{{ message.text }}</v-alert>

          <div>
            <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Current password</label>
            <input v-model="form.current_password" type="password" autocomplete="current-password" required style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
          </div>

          <div>
            <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">New password</label>
            <input v-model="form.new_password" type="password" autocomplete="new-password" required minlength="8" style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
          </div>

          <div>
            <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em;">Confirm new password</label>
            <input v-model="form.confirm" type="password" autocomplete="new-password" required minlength="8" style="width: 100%; background: rgba(0,0,0,0.5); border: none; border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white; padding: 8px 12px; font-family: 'JetBrains Mono', monospace; font-size: 13px; outline: none; transition: border-color 0.2s;" @focus="$event.target.style.borderColor = '#F7931A'" @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'" />
          </div>

          <button
            type="submit"
            class="align-self-start d-inline-flex align-center ga-2 rounded-pill px-5 font-mono"
            style="height: 40px; background: linear-gradient(to right, #EA580C, #F7931A); border: none; color: white; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5); cursor: pointer; transition: all 0.3s;"
            :disabled="loading"
            @mouseenter="if(!loading) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)' }"
            @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
          >{{ loading ? 'Updating...' : 'Update password' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>
