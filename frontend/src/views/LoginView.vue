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
    error.value = e.message || 'ACCESS DENIED'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-void p-6">
    <div class="w-full max-w-sm">
      <div class="cyber-card">
        <div class="text-center mb-8">
          <img src="/icon.png" alt="Tunlify" class="w-12 h-12 mx-auto mb-3" />
          <h1 class="font-display text-2xl text-neon text-glow tracking-widest animate-glitch">TUNLIFY</h1>
          <p class="text-xs text-muted mt-2 tracking-wider">// AUTHENTICATE TO PROCEED</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <div v-if="error" class="border border-danger bg-danger/10 text-danger text-sm p-3 font-mono">
            <span class="text-danger/60">[ERR]</span> {{ error }}
          </div>

          <div class="space-y-1">
            <label class="text-xs text-muted uppercase tracking-wider">User ID</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-neon/50 text-sm">›</span>
              <input v-model="form.username" type="text" required class="cyber-input !pl-7" autofocus />
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-xs text-muted uppercase tracking-wider">Passkey</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-neon/50 text-sm">›</span>
              <input v-model="form.password" type="password" required class="cyber-input !pl-7" />
            </div>
          </div>

          <button type="submit" :disabled="loading" class="cyber-btn w-full">
            {{ loading ? 'CONNECTING...' : 'INITIATE' }}
          </button>
        </form>

        <div class="mt-4 text-xs text-muted text-center tracking-wider">
          SYS.READY<span class="animate-blink text-neon">_</span>
        </div>
      </div>
    </div>
    <div class="scanlines"></div>
  </div>
</template>
