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
    <v-main style="background: #030304;">
      <div class="bg-orange-glow" style="position: fixed; inset: 0; z-index: 0;" />
      <div class="bg-grid" style="position: fixed; inset: 0; z-index: 0;" />
      <div style="position: relative; z-index: 1; min-height: 100vh; display: flex; flex-direction: column;">
        <header style="display: flex; align-items: center; justify-content: space-between; padding: 24px 32px; max-width: 1280px; width: 100%; margin: 0 auto;">
          <div class="d-flex align-center ga-3">
            <div class="rounded-lg d-flex align-center justify-center" style="width: 36px; height: 36px; background: rgba(247, 147, 26, 0.15); border: 1px solid rgba(247, 147, 26, 0.3);">
              <img src="/icon.png" alt="" style="width: 22px; height: 22px;" />
            </div>
            <span class="font-heading font-semibold" style="color: white; font-size: 18px;">Tunlify</span>
          </div>
          <a
            href="https://github.com/luzyver/Tunlify"
            target="_blank"
            rel="noopener"
            class="font-mono text-decoration-none"
            style="color: #94A3B8; font-size: 13px; transition: color 0.2s;"
            @mouseenter="$event.target.style.color = '#F7931A'"
            @mouseleave="$event.target.style.color = '#94A3B8'"
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

            <div class="rounded-2xl p-8" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
              <form @submit.prevent="handleLogin">
                <v-alert v-if="error" type="error" class="mb-4" variant="tonal">{{ error }}</v-alert>

                <div class="mb-5">
                  <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 12px; text-transform: uppercase; letter-spacing: 0.08em;">Username</label>
                  <input
                    v-model="form.username"
                    type="text"
                    required
                    autofocus
                    autocomplete="username"
                    placeholder="admin"
                    style="
                      width: 100%; height: 48px; background: rgba(0,0,0,0.5); border: none;
                      border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white;
                      padding: 8px 16px; font-family: 'JetBrains Mono', monospace; font-size: 14px;
                      outline: none; transition: border-color 0.2s, box-shadow 0.2s;
                    "
                    @focus="$event.target.style.borderColor = '#F7931A'; $event.target.style.boxShadow = '0 10px 20px -10px rgba(247, 147, 26, 0.3)'"
                    @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'; $event.target.style.boxShadow = 'none'"
                  />
                </div>

                <div class="mb-6">
                  <label class="font-mono d-block mb-2" style="color: #94A3B8; font-size: 12px; text-transform: uppercase; letter-spacing: 0.08em;">Password</label>
                  <input
                    v-model="form.password"
                    type="password"
                    required
                    autocomplete="current-password"
                    placeholder="&bull;&bull;&bull;&bull;&bull;&bull;&bull;&bull;"
                    style="
                      width: 100%; height: 48px; background: rgba(0,0,0,0.5); border: none;
                      border-bottom: 2px solid rgba(30, 41, 59, 0.8); color: white;
                      padding: 8px 16px; font-family: 'JetBrains Mono', monospace; font-size: 14px;
                      outline: none; transition: border-color 0.2s, box-shadow 0.2s;
                    "
                    @focus="$event.target.style.borderColor = '#F7931A'; $event.target.style.boxShadow = '0 10px 20px -10px rgba(247, 147, 26, 0.3)'"
                    @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'; $event.target.style.boxShadow = 'none'"
                  />
                </div>

                <button
                  type="submit"
                  :disabled="loading"
                  class="w-100 rounded-pill font-mono"
                  style="
                    height: 48px; background: linear-gradient(to right, #EA580C, #F7931A);
                    border: none; color: white; font-size: 14px; font-weight: 600;
                    text-transform: uppercase; letter-spacing: 0.08em;
                    box-shadow: 0 0 20px -5px rgba(234, 88, 12, 0.5);
                    cursor: pointer; transition: all 0.3s;
                  "
                  @mouseenter="if (!loading) { $event.target.style.transform = 'scale(1.02)'; $event.target.style.boxShadow = '0 0 30px -5px rgba(247, 147, 26, 0.6)'; }"
                  @mouseleave="$event.target.style.transform = 'scale(1)'; $event.target.style.boxShadow = '0 0 20px -5px rgba(234, 88, 12, 0.5)'"
                >
                  {{ loading ? 'Signing in...' : 'Sign in' }}
                </button>
              </form>
            </div>

            <p class="text-center font-mono" style="margin-top: 48px; color: rgba(148, 163, 184, 0.5); font-size: 11px;">
              Self-hosted &middot; MIT licensed
            </p>
          </div>
        </div>
      </div>
    </v-main>
  </v-app>
</template>
