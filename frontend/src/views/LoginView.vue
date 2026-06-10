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
        <div class="bg-orb bg-orb--1" />
        <div class="bg-orb bg-orb--2" />
        <div class="bg-orb bg-orb--3" />
        <div class="bg-grid" />

        <svg class="network-svg" viewBox="0 0 800 600" preserveAspectRatio="xMidYMid slice">
          <line class="net-line" x1="400" y1="300" x2="140" y2="120" />
          <line class="net-line" x1="400" y1="300" x2="660" y2="100" />
          <line class="net-line" x1="400" y1="300" x2="100" y2="480" />
          <line class="net-line" x1="400" y1="300" x2="700" y2="500" />
          <line class="net-line" x1="140" y1="120" x2="660" y2="100" />
          <line class="net-line" x1="100" y1="480" x2="700" y2="500" />
          <line class="net-line" x1="140" y1="120" x2="100" y2="480" />
          <line class="net-line" x1="660" y1="100" x2="700" y2="500" />

          <line class="net-line net-line--s" x1="400" y1="300" x2="260" y2="200" />
          <line class="net-line net-line--s" x1="400" y1="300" x2="540" y2="220" />
          <line class="net-line net-line--s" x1="400" y1="300" x2="290" y2="420" />
          <line class="net-line net-line--s" x1="400" y1="300" x2="530" y2="400" />
          <line class="net-line net-line--s" x1="260" y1="200" x2="540" y2="220" />
          <line class="net-line net-line--s" x1="290" y1="420" x2="530" y2="400" />
          <line class="net-line net-line--s" x1="260" y1="200" x2="290" y2="420" />
          <line class="net-line net-line--s" x1="540" y1="220" x2="530" y2="400" />

          <ellipse class="tunnel-ring" cx="400" cy="300" rx="90" ry="55" />
          <ellipse class="tunnel-ring tunnel-ring--b" cx="400" cy="300" rx="170" ry="105" />
          <ellipse class="tunnel-ring tunnel-ring--c" cx="400" cy="300" rx="280" ry="175" />

          <circle class="net-node net-node--center" cx="400" cy="300" r="14" />
          <circle class="net-node" cx="140" cy="120" r="7" style="animation-delay: 0s" />
          <circle class="net-node" cx="660" cy="100" r="7" style="animation-delay: 0.8s" />
          <circle class="net-node" cx="100" cy="480" r="7" style="animation-delay: 1.6s" />
          <circle class="net-node" cx="700" cy="500" r="7" style="animation-delay: 2.4s" />
          <circle class="net-node net-node--s" cx="260" cy="200" r="5" style="animation-delay: 0.4s" />
          <circle class="net-node net-node--s" cx="540" cy="220" r="5" style="animation-delay: 1.2s" />
          <circle class="net-node net-node--s" cx="290" cy="420" r="5" style="animation-delay: 2s" />
          <circle class="net-node net-node--s" cx="530" cy="400" r="5" style="animation-delay: 2.8s" />
        </svg>
      </div>

      <div class="login-content">
        <div
          class="brand-section"
          v-motion="'brand'"
          :initial="{ opacity: 0, x: -60 }"
          :visible="{ opacity: 1, x: 0, transition: { duration: 700, ease: 'easeOut' } }"
        >
          <div class="brand-header">
            <div class="brand-logo">
              <div class="logo-icon">
                <svg viewBox="0 0 24 24" fill="none" width="22" height="22">
                  <path d="M12 2L2 7v10l10 5 10-5V7l-10-5z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" />
                  <path d="M2 7l10 5m0 0l10-5m-10 5v10" stroke="currentColor" stroke-width="1.5" />
                  <path d="M7 9.5l5-2.5 5 2.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                </svg>
              </div>
              <span class="logo-text">Tunlify</span>
            </div>
          </div>

          <div class="brand-body">
            <h1 class="brand-title">
              Tunnel infrastructure,<br />
              <span class="text-highlight">reimagined.</span>
            </h1>
            <p class="brand-desc">
              Manage Cloudflare tunnels and Docker Compose projects from a single, unified command center.
            </p>
          </div>

          <div class="brand-footer">
            <div class="status-badge">
              <span class="status-dot" />
              <span>All systems operational</span>
            </div>
            <div class="brand-metrics">
              <div class="metric">
                <span class="metric-value">99.9%</span>
                <span class="metric-label">Uptime</span>
              </div>
              <div class="metric">
                <span class="metric-value">&lt;15ms</span>
                <span class="metric-label">Latency</span>
              </div>
              <div class="metric">
                <span class="metric-value">Real-time</span>
                <span class="metric-label">Logs</span>
              </div>
            </div>
          </div>
        </div>

        <div
          class="form-section"
          v-motion="'form'"
          :initial="{ opacity: 0, x: 60 }"
          :visible="{ opacity: 1, x: 0, transition: { duration: 700, delay: 200, ease: 'easeOut' } }"
        >
          <div class="form-card">
            <div class="form-header">
              <div class="form-badge">
                <svg viewBox="0 0 16 16" fill="none" width="12" height="12">
                  <rect x="2" y="6" width="12" height="8" rx="1.5" stroke="currentColor" stroke-width="1.2" />
                  <path d="M5 6V4a3 3 0 016 0v2" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
                </svg>
                Secure Access
              </div>
              <h2 class="form-title">Welcome back</h2>
              <p class="form-subtitle">Sign in to manage your infrastructure.</p>
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
                  color="teal"
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
                  color="teal"
                  class="login-field"
                  @click:append-inner="showPassword = !showPassword"
                />
              </div>

              <v-btn
                type="submit"
                :loading="loading"
                block
                size="large"
                class="submit-btn mt-2"
              >
                <template v-if="!loading">
                  Sign in
                  <v-icon end>mdi-arrow-right</v-icon>
                </template>
              </v-btn>
            </form>

            <div class="form-footer">
              <span>&copy; 2026 Tunlify</span>
              <span class="footer-sep">&middot;</span>
              <span>All connections encrypted</span>
            </div>
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
  background: #07080B;
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

.bg-orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.bg-orb--1 {
  width: 640px;
  height: 640px;
  top: -20%;
  left: -10%;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.07) 0%, transparent 60%);
  animation: orb-float-a 22s ease-in-out infinite;
}

.bg-orb--2 {
  width: 520px;
  height: 520px;
  bottom: -15%;
  right: -5%;
  background: radial-gradient(circle, rgba(6, 182, 212, 0.05) 0%, transparent 60%);
  animation: orb-float-b 26s ease-in-out infinite;
}

.bg-orb--3 {
  width: 360px;
  height: 360px;
  top: 40%;
  left: 35%;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.04) 0%, transparent 60%);
  animation: orb-float-a 18s ease-in-out infinite reverse;
}

@keyframes orb-float-a {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(40px, -30px) scale(1.06); }
  66% { transform: translate(-25px, 20px) scale(0.95); }
}

@keyframes orb-float-b {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-35px, 25px) scale(1.04); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.008) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.008) 1px, transparent 1px);
  background-size: 64px 64px;
  mask-image: radial-gradient(ellipse at 50% 50%, black 25%, transparent 65%);
  -webkit-mask-image: radial-gradient(ellipse at 50% 50%, black 25%, transparent 65%);
}

/* ─── Network SVG ─── */

.network-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0.55;
  pointer-events: none;
}

.net-line {
  stroke: rgba(16, 185, 129, 0.18);
  stroke-width: 1.2;
  stroke-dasharray: 5 7;
  animation: flow 2.5s linear infinite;
}

.net-line--s {
  stroke: rgba(6, 182, 212, 0.1);
  stroke-dasharray: 3 9;
  animation: flow-reverse 3.5s linear infinite;
}

@keyframes flow {
  to { stroke-dashoffset: -24; }
}

@keyframes flow-reverse {
  to { stroke-dashoffset: 24; }
}

.tunnel-ring {
  fill: none;
  stroke: rgba(16, 185, 129, 0.07);
  stroke-width: 1;
}

.tunnel-ring {
  transform-origin: 400px 300px;
  animation: ring-spin 28s linear infinite;
}

.tunnel-ring--b {
  stroke: rgba(6, 182, 212, 0.05);
  animation: ring-spin 40s linear infinite reverse;
}

.tunnel-ring--c {
  stroke: rgba(16, 185, 129, 0.04);
  animation: ring-spin 52s linear infinite;
}

@keyframes ring-spin {
  to { transform: rotate(360deg); }
}

.net-node {
  fill: rgba(16, 185, 129, 0.4);
  stroke: rgba(16, 185, 129, 0.15);
  stroke-width: 1;
  animation: node-pulse 3s ease-in-out infinite;
}

.net-node--center {
  fill: rgba(16, 185, 129, 0.5);
  stroke: rgba(16, 185, 129, 0.4);
  stroke-width: 2;
  filter: drop-shadow(0 0 10px rgba(16, 185, 129, 0.35));
  animation: node-pulse-center 3s ease-in-out infinite;
}

.net-node--s {
  fill: rgba(6, 182, 212, 0.25);
  stroke: rgba(6, 182, 212, 0.12);
}

@keyframes node-pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}

@keyframes node-pulse-center {
  0%, 100% { opacity: 0.6; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.15); }
}

/* ─── Content Layout ─── */

.login-content {
  position: relative;
  z-index: 1;
  display: flex;
  min-height: 100vh;
}

/* ─── Brand Section ─── */

.brand-section {
  flex: 0 0 45%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 64px 56px 64px 72px;
}

.brand-header {
  margin-bottom: 40px;
}

.brand-logo {
  display: inline-flex;
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
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.15);
  color: #34D399;
}

.logo-text {
  font-family: 'Archivo', sans-serif;
  font-weight: 800;
  font-size: 22px;
  color: #E4E8EE;
  letter-spacing: -0.02em;
}

.brand-body {
  max-width: 460px;
}

.brand-title {
  font-family: 'Archivo', sans-serif;
  font-weight: 900;
  font-size: 44px;
  line-height: 1.1;
  color: #E4E8EE;
  margin: 0 0 20px;
  letter-spacing: -0.03em;
}

.text-highlight {
  background: linear-gradient(135deg, #10B981 0%, #14B8A6 50%, #06B6D4 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.brand-desc {
  font-size: 15px;
  line-height: 1.7;
  color: #6B7280;
  margin: 0;
  font-weight: 400;
  max-width: 360px;
}

.brand-footer {
  margin-top: 56px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 500;
  color: #6B7280;
  background: rgba(16, 185, 129, 0.05);
  border: 1px solid rgba(16, 185, 129, 0.1);
  margin-bottom: 28px;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #10B981;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
  animation: status-pulse 2.5s ease-in-out infinite;
}

@keyframes status-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.brand-metrics {
  display: flex;
  gap: 32px;
}

.metric {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-value {
  font-family: 'Archivo', sans-serif;
  font-weight: 700;
  font-size: 16px;
  color: #E4E8EE;
  letter-spacing: -0.01em;
}

.metric-label {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #4B5563;
}

/* ─── Form Section ─── */

.form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 64px 72px 64px 40px;
}

.form-card {
  width: 100%;
  max-width: 420px;
  background: rgba(15, 17, 21, 0.65);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 40px;
  box-shadow:
    0 0 60px -20px rgba(0, 0, 0, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.03);
}

.form-header {
  margin-bottom: 28px;
}

.form-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(16, 185, 129, 0.8);
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.12);
  margin-bottom: 20px;
}

.form-title {
  font-family: 'Archivo', sans-serif;
  font-weight: 800;
  font-size: 26px;
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

.form-fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* Vuetify field overrides */
.form-card :deep(.v-field) {
  background: rgba(255, 255, 255, 0.02) !important;
}

.form-card :deep(.v-field--variant-outlined .v-field__outline) {
  color: rgba(255, 255, 255, 0.08);
}

.form-card :deep(.v-field:hover .v-field__outline) {
  color: rgba(255, 255, 255, 0.15);
}

.form-card :deep(.v-field--focused .v-field__outline) {
  color: #10B981 !important;
  opacity: 0.6;
}

.form-card :deep(.v-field-label--floating) {
  color: #10B981 !important;
}

.form-card :deep(.v-field .v-icon) {
  color: rgba(255, 255, 255, 0.15);
}

.form-card :deep(.v-field--focused .v-icon) {
  color: rgba(16, 185, 129, 0.6);
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

/* Submit Button */
.submit-btn {
  height: 48px !important;
  font-family: 'Archivo', sans-serif !important;
  font-weight: 700 !important;
  font-size: 15px !important;
  text-transform: none !important;
  letter-spacing: 0.01em !important;
  border-radius: 12px !important;
  background: linear-gradient(135deg, #10B981, #059669) !important;
  border: none !important;
  color: white !important;
  transition: all 0.3s ease !important;
  box-shadow: 0 0 30px -8px rgba(16, 185, 129, 0.3) !important;
  overflow: hidden !important;
  position: relative !important;
}

.submit-btn:hover {
  transform: translateY(-1px) !important;
  box-shadow: 0 0 40px -6px rgba(16, 185, 129, 0.5) !important;
}

.submit-btn:active {
  transform: translateY(0) scale(0.99) !important;
}

/* Override Vuetify's default loading overlay */
.submit-btn :deep(.v-btn__overlay) {
  display: none !important;
}

/* Form Footer */
.form-footer {
  margin-top: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 12px;
  color: #3F3F46;
  font-weight: 400;
}

.footer-sep {
  color: #3F3F46;
}

/* ─── Responsive ─── */

@media (max-width: 1024px) {
  .brand-section {
    padding: 56px 40px 56px 48px;
  }

  .brand-title {
    font-size: 36px;
  }

  .form-section {
    padding: 56px 48px 56px 32px;
  }
}

@media (max-width: 900px) {
  .login-content {
    flex-direction: column;
  }

  .brand-section {
    flex: none;
    padding: 48px 32px 32px;
    min-height: 45vh;
    justify-content: flex-end;
  }

  .brand-body {
    max-width: 100%;
  }

  .brand-title {
    font-size: 32px;
  }

  .brand-desc {
    max-width: 100%;
  }

  .brand-footer {
    margin-top: 32px;
  }

  .brand-metrics {
    gap: 24px;
  }

  .form-section {
    padding: 24px 32px 48px;
  }

  .form-card {
    max-width: 100%;
    padding: 32px;
  }

  .network-svg {
    opacity: 0.3;
  }
}

@media (max-width: 480px) {
  .brand-section {
    padding: 32px 20px 24px;
    min-height: 35vh;
  }

  .brand-title {
    font-size: 26px;
  }

  .brand-desc {
    font-size: 13px;
  }

  .brand-metrics {
    gap: 16px;
  }

  .metric-value {
    font-size: 14px;
  }

  .form-section {
    padding: 16px 20px 40px;
  }

  .form-card {
    padding: 24px 20px;
    border-radius: 16px;
  }

  .form-title {
    font-size: 22px;
  }
}
</style>
