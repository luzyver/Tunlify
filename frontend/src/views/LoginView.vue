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
    const data =await apiFetch<{ access_token: string }>('/auth/login', {
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
    <v-main>
      <div class="bg-orange-glow" style="position: fixed; inset: 0; z-index: 0;" />
      <div class="bg-blue-glow" style="position: fixed; inset: 0; z-index: 0;" />
      <div class="bg-grid" style="position: fixed; inset: 0; z-index: 0;" />
      <div style="position: relative; z-index: 1; min-height: 100vh; display: flex; flex-direction: column;">
        <header style="display: flex; align-items: center; justify-content: space-between; padding: 24px 32px; max-width: 1280px; width: 100%; margin: 0 auto;">
          <div class="d-flex align-center ga-3">
            <v-avatar size="36" color="surface" variant="outlined" style="border: 1px solid rgba(247, 147, 26, 0.3);">
              <img src="/icon.png?v=2" alt="" style="width: 22px; height: 22px;" />
            </v-avatar>
            <span class="font-heading font-semibold" style="color: white; font-size: 18px;">Tunlify</span>
          </div>
          <a
            href="https://github.com/luzyver/Tunlify"
            target="_blank"
            rel="noopener"
            class="font-mono text-decoration-none"
            style="color: #94A3B8; font-size: 13px; transition: color 0.2s;"
            @mouseenter="($event.target as HTMLElement).style.color = '#F7931A'"
            @mouseleave="($event.target as HTMLElement).style.color = '#94A3B8'"
          >Source &Nearr;</a>
        </header>

        <div style="flex: 1; display: flex; align-items: center; justify-content: center; padding: 48px 24px;">
          <div style="max-width: 420px; width: 100%;">
            <div class="text-center mb-12">
              <p class="eyebrow mb-4">Tunlify &middot; Console</p>
              <h1 class="font-heading" style="font-size: 32px; font-weight: 600; line-height: 1.2; color: white;">
                Sign in to your
                <span class="text-gradient-gold">tunnel</span>
              </h1>
              <p style="margin-top: 12px; color: #94A3B8; font-size: 15px;">
                Manage Cloudflare tunnels and Compose projects.
              </p>
            </div>

            <v-card flat color="surface" class="rounded-2xl" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
              <v-card-text class="pa-8">
                <form @submit.prevent="handleLogin">
                  <v-alert v-if="error" type="error" class="mb-4" variant="tonal" closable @click:close="error = ''">{{ error }}</v-alert>

                  <v-text-field
                    v-model="form.username"
                    label="Username"
                    type="text"
                    required
                    autofocus
                    autocomplete="username"
                    class="mb-2"
                    hide-details="auto"
                    variant="filled"
                    color="primary"
                    bg-color="#0A0C10"
                  />

                  <v-text-field
                    v-model="form.password"
                    label="Password"
                    type="password"
                    required
                    autocomplete="current-password"
                    class="mb-6"
                    hide-details="auto"
                    variant="filled"
                    color="primary"
                    bg-color="#0A0C10"
                  />

                  <v-btn
                    type="submit"
                    :loading="loading"
                    block
                    size="large"
                    color="primary"
                    variant="flat"
                    class="rounded-pill font-mono text-none"
                    style="height: 48px; font-size: 14px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; background: linear-gradient(to right, #EA580C, #F7931A); box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5);"
                  >
                    Sign in
                  </v-btn>
                </form>
              </v-card-text>
            </v-card>

            <p class="text-center font-mono" style="margin-top: 48px; color: rgba(148, 163, 184, 0.5); font-size: 11px;">
              Self-hosted &middot; MIT licensed
            </p>
          </div>
        </div>
      </div>
    </v-main>
  </v-app>
</template>
