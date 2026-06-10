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
const showPassword = ref(false)

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
    <div class="login-root">
      <div class="login-bg">
        <div class="bg-glow bg-glow--orange" />
        <div class="bg-glow bg-glow--blue" />
        <div class="bg-glow bg-glow--orange-sm" />
        <div class="bg-glow bg-glow--blue-sm" />
        <div class="bg-grid" />

        <svg class="network-svg" viewBox="0 0 800 600" preserveAspectRatio="xMidYMid slice">
          <defs>
            <linearGradient id="g-orange" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#F59E0B" stop-opacity="0.25" />
              <stop offset="100%" stop-color="#EA580C" stop-opacity="0.08" />
            </linearGradient>
            <linearGradient id="g-blue" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#3B82F6" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#1D4ED8" stop-opacity="0.05" />
            </linearGradient>
            <radialGradient id="center-glow" cx="50%" cy="50%" r="50%">
              <stop offset="0%" stop-color="#F59E0B" stop-opacity="0.3" />
              <stop offset="60%" stop-color="#3B82F6" stop-opacity="0.08" />
              <stop offset="100%" stop-color="transparent" stop-opacity="0" />
            </radialGradient>
          </defs>

          <circle cx="400" cy="300" r="160" fill="url(#center-glow)" />

          <line class="net-line net-line--orange" x1="400" y1="300" x2="140" y2="120" />
          <line class="net-line net-line--orange" x1="400" y1="300" x2="660" y2="100" />
          <line class="net-line net-line--orange" x1="400" y1="300" x2="100" y2="480" />
          <line class="net-line net-line--orange" x1="400" y1="300" x2="700" y2="500" />
          <line class="net-line net-line--orange" x1="140" y1="120" x2="660" y2="100" />
          <line class="net-line net-line--orange" x1="100" y1="480" x2="700" y2="500" />
          <line class="net-line net-line--orange" x1="140" y1="120" x2="100" y2="480" />
          <line class="net-line net-line--orange" x1="660" y1="100" x2="700" y2="500" />

          <line class="net-line net-line--blue" x1="400" y1="300" x2="260" y2="200" />
          <line class="net-line net-line--blue" x1="400" y1="300" x2="540" y2="220" />
          <line class="net-line net-line--blue" x1="400" y1="300" x2="290" y2="420" />
          <line class="net-line net-line--blue" x1="400" y1="300" x2="530" y2="400" />
          <line class="net-line net-line--blue" x1="260" y1="200" x2="540" y2="220" />
          <line class="net-line net-line--blue" x1="290" y1="420" x2="530" y2="400" />
          <line class="net-line net-line--blue" x1="260" y1="200" x2="290" y2="420" />
          <line class="net-line net-line--blue" x1="540" y1="220" x2="530" y2="400" />

          <ellipse class="tunnel-ring tunnel-ring--orange" cx="400" cy="300" rx="90" ry="55" />
          <ellipse class="tunnel-ring tunnel-ring--blue" cx="400" cy="300" rx="170" ry="105" />
          <ellipse class="tunnel-ring tunnel-ring--orange" cx="400" cy="300" rx="280" ry="175" />

          <circle class="net-node net-node--center" cx="400" cy="300" r="14" />
          <circle class="net-node net-node--orange" cx="140" cy="120" r="7" style="animation-delay: 0s" />
          <circle class="net-node net-node--orange" cx="660" cy="100" r="7" style="animation-delay: 0.8s" />
          <circle class="net-node net-node--orange" cx="100" cy="480" r="7" style="animation-delay: 1.6s" />
          <circle class="net-node net-node--orange" cx="700" cy="500" r="7" style="animation-delay: 2.4s" />
          <circle class="net-node net-node--blue" cx="260" cy="200" r="5" style="animation-delay: 0.4s" />
          <circle class="net-node net-node--blue" cx="540" cy="220" r="5" style="animation-delay: 1.2s" />
          <circle class="net-node net-node--blue" cx="290" cy="420" r="5" style="animation-delay: 2s" />
          <circle class="net-node net-node--blue" cx="530" cy="400" r="5" style="animation-delay: 2.8s" />
        </svg>
      </div>

      <div class="login-content">
        <div
          v-motion
          :initial="{ opacity: 0, y: 30, scale: 0.97 }"
          :visible="{ opacity: 1, y: 0, scale: 1, transition: { duration: 600, ease: 'easeOut' } }"
          class="form-card"
        >
          <div class="form-brand">
            <div class="brand-logo">
              <div class="logo-icon">
                <img src="/icon.png" alt="Tunlify" class="logo-img" />
              </div>
              <div class="logo-text-group">
                <span class="logo-wordmark">Tunlify</span>
                <span class="logo-tagline">Tunnel Management</span>
              </div>
            </div>
          </div>

          <div class="form-divider" />

          <div class="form-header">
            <h2 class="form-title">Welcome back</h2>
            <p class="form-subtitle">Sign in to your account to continue.</p>
          </div>

          <v-alert
            v-if="error"
            type="error"
            variant="tonal"
            closable
            @click:close="error = ''"
            class="mb-4"
            density="compact"
            title="Authentication failed"
          >
            {{ error }}
          </v-alert>

          <form @submit.prevent="handleLogin">
            <div class="form-fields">
              <v-text-field
                v-model="form.username"
                label="Username"
                prepend-inner-icon="mdi-account-circle-outline"
                required
                autocomplete="username"
                variant="outlined"
                class="login-field"
              />
              <v-text-field
                v-model="form.password"
                label="Password"
                prepend-inner-icon="mdi-lock-outline"
                :append-inner-icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="current-password"
                variant="outlined"
                class="login-field"
                @click:append-inner="showPassword = !showPassword"
              />
            </div>

            <v-btn
              type="submit"
              :loading="loading"
              block
              size="large"
              class="submit-btn mt-4"
            >
              <template v-if="!loading">
                Sign in
                <v-icon end>mdi-arrow-right</v-icon>
              </template>
            </v-btn>
          </form>

          <div class="form-footer">
            <span class="status-indicator">
              <span class="status-dot" />
              All systems operational
            </span>
            <span class="footer-sep">&middot;</span>
            <span>&copy; 2026</span>
          </div>
        </div>
      </div>
    </div>
  </v-app>
</template>

<style scoped>
.login-root {
  position: relative;
  width: 100%;
  min-height: 100vh;
  background: #0a0c12;
  font-family: 'DM Sans', sans-serif;
  overflow: hidden;
}

/* ─── Background ─── */

.login-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
}

.bg-glow {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  filter: blur(80px);
}

.bg-glow--orange {
  width: 700px;
  height: 700px;
  top: -25%;
  left: -15%;
  background: radial-gradient(circle, rgba(247, 147, 26, 0.25) 0%, rgba(234, 88, 12, 0.08) 40%, transparent 70%);
  animation: orb-float-a 20s ease-in-out infinite;
}

.bg-glow--blue {
  width: 600px;
  height: 600px;
  bottom: -20%;
  right: -10%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.2) 0%, rgba(37, 99, 235, 0.06) 40%, transparent 70%);
  animation: orb-float-b 24s ease-in-out infinite;
}

.bg-glow--orange-sm {
  width: 400px;
  height: 400px;
  top: 50%;
  left: 20%;
  background: radial-gradient(circle, rgba(247, 147, 26, 0.12) 0%, transparent 60%);
  animation: orb-float-a 18s ease-in-out infinite reverse;
}

.bg-glow--blue-sm {
  width: 350px;
  height: 350px;
  top: 10%;
  right: 25%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 60%);
  animation: orb-float-b 16s ease-in-out infinite reverse;
}

@keyframes orb-float-a {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(50px, -35px) scale(1.08); }
  66% { transform: translate(-30px, 25px) scale(0.93); }
}

@keyframes orb-float-b {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-45px, 30px) scale(1.06); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
  background-size: 64px 64px;
  mask-image: radial-gradient(ellipse at 50% 50%, black 30%, transparent 70%);
  -webkit-mask-image: radial-gradient(ellipse at 50% 50%, black 30%, transparent 70%);
}

/* ─── Network SVG ─── */

.network-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0.6;
  pointer-events: none;
}

.net-line {
  stroke-width: 1.5;
  stroke-dasharray: 6 8;
}

.net-line--orange {
  stroke: url(#g-orange);
  animation: flow 3s linear infinite;
}

.net-line--blue {
  stroke: url(#g-blue);
  stroke-dasharray: 4 10;
  animation: flow-reverse 4s linear infinite;
}

@keyframes flow {
  to { stroke-dashoffset: -28; }
}

@keyframes flow-reverse {
  to { stroke-dashoffset: 28; }
}

.tunnel-ring {
  fill: none;
  stroke-width: 1.5;
  transform-origin: 400px 300px;
}

.tunnel-ring--orange {
  stroke: rgba(247, 147, 26, 0.12);
  animation: ring-spin 30s linear infinite;
}

.tunnel-ring--blue {
  stroke: rgba(59, 130, 246, 0.08);
  animation: ring-spin 42s linear infinite reverse;
}

@keyframes ring-spin {
  to { transform: rotate(360deg); }
}

.net-node {
  stroke-width: 1.5;
  animation: node-pulse 3s ease-in-out infinite;
}

.net-node--center {
  fill: rgba(247, 147, 26, 0.6);
  stroke: rgba(247, 147, 26, 0.5);
  stroke-width: 2.5;
  filter: drop-shadow(0 0 16px rgba(247, 147, 26, 0.5));
  animation: node-pulse-center 3s ease-in-out infinite;
}

.net-node--orange {
  fill: rgba(247, 147, 26, 0.5);
  stroke: rgba(247, 147, 26, 0.25);
  filter: drop-shadow(0 0 6px rgba(247, 147, 26, 0.3));
}

.net-node--blue {
  fill: rgba(59, 130, 246, 0.35);
  stroke: rgba(59, 130, 246, 0.2);
  filter: drop-shadow(0 0 6px rgba(59, 130, 246, 0.25));
}

@keyframes node-pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

@keyframes node-pulse-center {
  0%, 100% { opacity: 0.7; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.2); }
}

/* ─── Content ─── */

.login-content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 32px;
}

/* ─── Form Card ─── */

.form-card {
  width: 100%;
  max-width: 400px;
  background: rgba(8, 10, 15, 0.7);
  backdrop-filter: blur(32px);
  -webkit-backdrop-filter: blur(32px);
  border: 1px solid rgba(247, 147, 26, 0.08);
  border-radius: 20px;
  padding: 36px;
  box-shadow:
    0 0 60px -20px rgba(247, 147, 26, 0.15),
    0 0 120px -40px rgba(59, 130, 246, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.04);
}

/* ─── Brand inside card ─── */

.brand-logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(247, 147, 26, 0.1), rgba(234, 88, 12, 0.05));
  border: 1px solid rgba(247, 147, 26, 0.12);
  flex-shrink: 0;
}

.logo-img {
  width: 26px;
  height: 26px;
  object-fit: contain;
}

.logo-text-group {
  display: flex;
  flex-direction: column;
}

.logo-wordmark {
  font-family: 'Archivo', sans-serif;
  font-weight: 800;
  font-size: 20px;
  color: #E4E8EE;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.logo-tagline {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #52525B;
}

/* ─── Divider ─── */

.form-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(247, 147, 26, 0.15), rgba(59, 130, 246, 0.12), transparent);
  margin: 24px 0;
}

/* ─── Form Header ─── */

.form-header {
  margin-bottom: 24px;
}

.form-title {
  font-family: 'Archivo', sans-serif;
  font-weight: 800;
  font-size: 24px;
  color: #E4E8EE;
  margin: 0 0 6px;
  letter-spacing: -0.02em;
}

.form-subtitle {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
  line-height: 1.5;
}

/* ─── Form Fields ─── */

.form-fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-card :deep(.v-field) {
  background: rgba(255, 255, 255, 0.02) !important;
}

.form-card :deep(.v-field--variant-outlined .v-field__outline) {
  color: rgba(255, 255, 255, 0.08);
}

.form-card :deep(.v-field:hover .v-field__outline) {
  color: rgba(255, 255, 255, 0.18);
}

.form-card :deep(.v-field--focused .v-field__outline) {
  color: #F59E0B !important;
  opacity: 0.8;
}

.form-card :deep(.v-field-label--floating) {
  color: #F59E0B !important;
}

.form-card :deep(.v-field .v-icon) {
  color: rgba(255, 255, 255, 0.18);
}

.form-card :deep(.v-field--focused .v-icon) {
  color: rgba(247, 147, 26, 0.7);
}

.form-card :deep(.v-field__input) {
  color: #E4E8EE;
  font-size: 14px;
  padding-top: 18px !important;
}

.form-card :deep(.v-field-label) {
  color: #4B5563;
  font-size: 14px;
  font-weight: 400;
}

/* ─── Submit Button ─── */

.submit-btn {
  height: 48px !important;
  font-family: 'Archivo', sans-serif !important;
  font-weight: 700 !important;
  font-size: 15px !important;
  text-transform: none !important;
  letter-spacing: 0.01em !important;
  border-radius: 12px !important;
  background: linear-gradient(135deg, #F59E0B, #EA580C) !important;
  border: none !important;
  color: white !important;
  transition: all 0.3s ease !important;
  box-shadow: 0 0 40px -8px rgba(247, 147, 26, 0.4) !important;
  overflow: hidden !important;
  position: relative !important;
}

.submit-btn:hover {
  transform: translateY(-1px) !important;
  box-shadow: 0 0 50px -4px rgba(247, 147, 26, 0.55) !important;
}

.submit-btn:active {
  transform: translateY(0) scale(0.99) !important;
}

.submit-btn :deep(.v-btn__overlay) {
  display: none !important;
}

/* ─── Footer ─── */

.form-footer {
  margin-top: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 12px;
  color: #3F3F46;
  font-weight: 400;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #22C55E;
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.4);
  animation: status-pulse 2.5s ease-in-out infinite;
}

@keyframes status-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.footer-sep {
  color: #3F3F46;
}

/* ─── Responsive ─── */

@media (max-width: 480px) {
  .login-content {
    padding: 16px;
  }

  .form-card {
    padding: 28px 24px;
    border-radius: 16px;
  }

  .form-title {
    font-size: 20px;
  }

  .network-svg {
    opacity: 0.35;
  }
}
</style>
