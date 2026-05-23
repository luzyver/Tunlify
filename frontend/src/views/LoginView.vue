<script setup lang="ts">
import { reactive, ref } from 'vue'
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
  <div class="min-h-screen bg-bg flex flex-col">

    <header class="max-w-marketing mx-auto w-full px-6 py-6 flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <img src="/icon.png" alt="" class="w-7 h-7 rounded-sm" />
        <span class="text-md font-semibold tracking-tight text-text">Tunlify</span>
      </div>
      <a
        href="https://github.com/luzyver/Tunlify"
        target="_blank"
        rel="noopener"
        class="text-sm text-text-muted hover:text-text transition-colors"
      >Source ↗</a>
    </header>

    <div class="flex-1 flex items-center justify-center px-6 py-12">
      <div class="w-full max-w-[420px]">

        <div class="mb-12 text-center">
          <p class="eyebrow mb-4">Tunlify · Console</p>
          <h1 class="text-display text-3xl text-text leading-tight">
            Sign in to your tunnel.
          </h1>
          <p class="mt-3 text-md text-text-muted">
            Manage Cloudflare tunnels and Compose projects.
          </p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div v-if="error" class="alert-danger">{{ error }}</div>

          <div>
            <label for="login-username" class="field-label">Username</label>
            <input
              id="login-username"
              v-model="form.username"
              type="text"
              required
              autofocus
              autocomplete="username"
              class="input-pill"
              placeholder="admin"
            />
          </div>

          <div>
            <label for="login-password" class="field-label">Password</label>
            <input
              id="login-password"
              v-model="form.password"
              type="password"
              required
              autocomplete="current-password"
              class="input-pill"
              placeholder="••••••••"
            />
          </div>

          <button type="submit" :disabled="loading" class="btn-pill w-full">
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </button>
        </form>

        <p class="mt-12 text-center text-xs text-text-dim">
          Self-hosted · MIT licensed
        </p>
      </div>
    </div>
  </div>
</template>
