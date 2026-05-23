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
  <div class="min-h-screen flex items-center justify-center bg-bg-alt p-6">
    <div class="w-full max-w-sm">
      <div class="card p-6">
        <div class="text-center mb-6">
          <img src="/icon.png" alt="Tunlify" class="w-10 h-10 mx-auto mb-3" />
          <h1 class="text-xl font-bold text-text">Sign in to Tunlify</h1>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div v-if="error" class="bg-danger/10 border border-danger/20 text-danger text-sm p-3 rounded">
            {{ error }}
          </div>

          <div class="space-y-1">
            <label class="text-sm font-medium text-text">Username</label>
            <input v-model="form.username" type="text" required class="input" placeholder="admin" autofocus />
          </div>

          <div class="space-y-1">
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
