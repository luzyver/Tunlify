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
  <v-app>
    <v-main class="bg-background">
      <div class="d-flex flex-column" style="min-height: 100vh;">
        <header class="d-flex align-center justify-space-between pa-6" style="max-width: 1080px; width: 100%; margin: 0 auto;">
          <div class="d-flex align-center ga-2">
            <v-avatar size="28" rounded>
              <v-img src="/icon.png" />
            </v-avatar>
            <span class="text-h6 font-weight-semibold text-on-surface">Tunlify</span>
          </div>
          <a
            href="https://github.com/luzyver/Tunlify"
            target="_blank"
            rel="noopener"
            class="text-medium-emphasis text-decoration-none"
          >Source &Nearr;</a>
        </header>

        <div class="flex-grow-1 d-flex align-center justify-center pa-6">
          <div style="max-width: 420px; width: 100%;">
            <div class="text-center mb-12">
              <p class="eyebrow mb-4">Tunlify &middot; Console</p>
              <h1 class="text-display" style="font-size: 30px; line-height: 1.2; color: #2a2924;">
                Sign in to your tunnel.
              </h1>
              <p class="mt-3 text-body-1 text-medium-emphasis">
                Manage Cloudflare tunnels and Compose projects.
              </p>
            </div>

            <v-form @submit.prevent="handleLogin">
              <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

              <v-text-field
                v-model="form.username"
                label="Username"
                placeholder="admin"
                required
                autofocus
                autocomplete="username"
                class="mb-4"
              />

              <v-text-field
                v-model="form.password"
                label="Password"
                placeholder="&bull;&bull;&bull;&bull;&bull;&bull;&bull;&bull;"
                type="password"
                required
                autocomplete="current-password"
                class="mb-6"
              />

              <v-btn
                type="submit"
                color="primary"
                size="large"
                block
                :loading="loading"
                class="rounded-pill"
              >
                {{ loading ? 'Signing in...' : 'Sign in' }}
              </v-btn>
            </v-form>

            <p class="mt-12 text-center text-caption text-disabled">
              Self-hosted &middot; MIT licensed
            </p>
          </div>
        </div>
      </div>
    </v-main>
  </v-app>
</template>
