<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../composables/useApi'

const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()

const form = reactive({ username: '', password: '' })
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    const data = await apiFetch<{ access_token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(form),
    })
    authStore.setAuth(data.access_token, form.username)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || 'Invalid credentials'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-bg p-6">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <img src="/icon.png" alt="Tunlify" class="w-10 h-10 mx-auto mb-4" />
        <h1 class="font-display text-3xl font-semibold text-text">Tunlify</h1>
        <p class="text-sm text-text-muted mt-2">Sign in to your dashboard</p>
      </div>

      <div class="card p-6 space-y-6">
        <div v-if="error" class="border border-red-200 bg-red-50 text-red-700 text-sm p-3 rounded-md">
          {{ error }}
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-text">Username</label>
            <input v-model="form.username" type="text" required class="input" placeholder="admin" autofocus />
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-text">Password</label>
            <input v-model="form.password" type="password" required class="input" placeholder="••••••••" />
          </div>

          <button type="submit" :disabled="loading" class="btn-primary w-full">
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
