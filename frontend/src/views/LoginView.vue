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
  <div class="login-root">
    <div class="orb orb-1" />
    <div class="orb orb-2" />
    <div class="orb orb-3" />
    <div class="dot-grid" />

    <header class="login-header">
      <div class="login-brand">
        <div class="brand-icon">
          <img src="/icon.png" alt="Tunlify" class="brand-img" />
        </div>
        <span class="brand-name">Tunlify</span>
      </div>
      <a
        href="https://github.com/luzyver/Tunlify"
        target="_blank"
        rel="noopener"
        class="source-link"
      >Source &rarr;</a>
    </header>

    <main class="login-main">
      <div class="login-card">
        <div class="login-card-inner">
          <div class="login-card-header">
            <span class="eyebrow">Tunlify &middot; Console</span>
            <h1 class="login-title">
              Sign in to your
              <span class="text-gradient-gold">tunnel</span>
            </h1>
            <p class="login-subtitle">
              Manage Cloudflare tunnels and Compose projects.
            </p>
          </div>

          <form @submit.prevent="handleLogin" class="login-form">
            <v-alert v-if="error" type="error" class="mb-4" variant="tonal" closable @click:close="error = ''">{{ error }}</v-alert>

            <v-text-field
              v-model="form.username"
              label="Username"
              type="text"
              required
              autofocus
              autocomplete="username"
              prepend-inner-icon="mdi-account-outline"
              variant="filled"
              color="primary"
              bg-color="#080A0E"
              hide-details="auto"
              class="login-input"
            />

            <v-text-field
              v-model="form.password"
              label="Password"
              type="password"
              required
              autocomplete="current-password"
              prepend-inner-icon="mdi-lock-outline"
              variant="filled"
              color="primary"
              bg-color="#080A0E"
              hide-details="auto"
              class="login-input"
            />

            <v-btn
              type="submit"
              :loading="loading"
              block
              size="large"
              color="primary"
              variant="flat"
              class="login-btn"
              style="background: linear-gradient(135deg, #EA580C, #F7931A); box-shadow: 0 0 24px -6px rgba(234, 88, 12, 0.4);"
            >
              Sign in
            </v-btn>
          </form>
        </div>
      </div>

      <p class="login-footer">
        &copy; 2026 Tunlify
      </p>
    </main>
  </div>
</template>

<style scoped>
.login-root {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #0d0400 0%, #030304 40%, #00040d 100%);
  overflow: hidden;
}

.orb {
  position: fixed;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  will-change: transform;
  z-index: 0;
}

.orb-1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(247, 147, 26, 0.25) 0%, transparent 70%);
  top: -200px;
  left: -200px;
  animation: float-orb 20s ease-in-out infinite;
}

.orb-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(37, 99, 235, 0.2) 0%, transparent 70%);
  bottom: -150px;
  right: -150px;
  animation: float-orb 25s ease-in-out infinite reverse;
}

.orb-3 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(234, 88, 12, 0.15) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: float-orb 30s ease-in-out infinite 5s;
}

@keyframes float-orb {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(60px, -40px) scale(1.05); }
  50% { transform: translate(-30px, 60px) scale(0.95); }
  75% { transform: translate(40px, 20px) scale(1.02); }
}

.dot-grid {
  position: fixed;
  inset: 0;
  z-index: 0;
  background-size: 40px 40px;
  background-image:
    radial-gradient(circle, rgba(148, 163, 184, 0.08) 1px, transparent 1px);
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
  -webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
  pointer-events: none;
}

.login-header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 32px;
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
}

.login-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  border: 1px solid rgba(247, 147, 26, 0.25);
  background: rgba(247, 147, 26, 0.06);
}

.brand-img {
  width: 20px;
  height: 20px;
}

.brand-name {
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 600;
  font-size: 18px;
  color: white;
}

.source-link {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: #94A3B8;
  text-decoration: none;
  transition: color 0.2s;
}

.source-link:hover {
  color: #F7931A;
}

.login-main {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 20px;
  background: rgba(15, 17, 21, 0.6);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(30, 41, 59, 0.5);
  box-shadow: 0 0 60px -15px rgba(247, 147, 26, 0.08);
  transition: box-shadow 0.4s, border-color 0.4s;
  animation: card-in 0.6s ease-out both;
}

.login-card:focus-within {
  border-color: rgba(247, 147, 26, 0.3);
  box-shadow: 0 0 60px -10px rgba(247, 147, 26, 0.15);
}

@keyframes card-in {
  from {
    opacity: 0;
    transform: translateY(16px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-card-inner {
  padding: 40px;
}

.login-card-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 32px;
  font-weight: 600;
  line-height: 1.2;
  color: white;
  margin-top: 12px;
}

.login-subtitle {
  margin-top: 10px;
  color: #94A3B8;
  font-size: 15px;
  line-height: 1.5;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.login-input :deep(.v-field) {
  border-radius: 12px;
  transition: box-shadow 0.25s;
}

.login-input :deep(.v-field--focused) {
  box-shadow: 0 0 0 1px rgba(247, 147, 26, 0.3), 0 0 20px -8px rgba(247, 147, 26, 0.2);
}

.login-input :deep(.v-field__prepend-inner) {
  color: rgba(148, 163, 184, 0.5);
  margin-top: 0;
}

.login-input :deep(.v-field--focused .v-field__prepend-inner) {
  color: #F7931A;
}

.login-btn {
  height: 48px;
  border-radius: 999px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: none;
  position: relative;
  overflow: hidden;
  transition: box-shadow 0.3s, transform 0.2s;
}

.login-btn:hover {
  box-shadow: 0 0 32px -4px rgba(234, 88, 12, 0.6) !important;
  transform: scale(1.01);
}

.login-btn:active {
  transform: scale(0.98);
}

.login-footer {
  margin-top: 48px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: rgba(148, 163, 184, 0.4);
  text-align: center;
}
</style>