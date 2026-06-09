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
  <div class="login-root">
    <div class="brand-panel">
      <div class="brand-scene">
        <div class="glow" />
        <div class="ring ring--a" />
        <div class="ring ring--b" />
        <div class="ring ring--c" />
        <div class="ring ring--d" />
        <div class="dot dot--1" />
        <div class="dot dot--2" />
        <div class="dot dot--3" />
        <div class="dot dot--4" />
        <div class="dot dot--5" />
        <div class="dot dot--6" />
        <div class="grid-mesh" />
      </div>
      <div class="brand-body">
        <div class="brand-logo">
          <div class="brand-icon-frame">
            <img src="/icon.png" alt="Tunlify" class="brand-icon" />
          </div>
          <span class="brand-wordmark">Tunlify</span>
        </div>
        <div class="brand-text">
          <h1 class="brand-headline">
            Tunnel infrastructure,<br />
            <span class="text-gradient">simplified.</span>
          </h1>
          <p class="brand-description">
            Manage Cloudflare tunnels and Docker Compose projects from a single, modern dashboard.
          </p>
        </div>
      </div>
      <div class="brand-footer">
        <span class="status-indicator">
          <span class="status-dot" />
          <span>All systems operational</span>
        </span>
      </div>
    </div>

    <div class="form-panel">
      <div class="form-body">
        <div class="form-header">
          <span class="form-badge">Secure area</span>
          <h2 class="form-title">Welcome back</h2>
          <p class="form-subtitle">Sign in to your account to continue.</p>
        </div>

        <form @submit.prevent="handleLogin" class="form-element" novalidate>
          <Transition name="alert">
            <div v-if="error" class="form-alert" role="alert">
              <svg class="alert-icon" viewBox="0 0 20 20" fill="none" aria-hidden="true">
                <circle cx="10" cy="10" r="9" stroke="currentColor" stroke-width="1.5" />
                <path d="M10 6v4M10 13v0" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
              <span class="alert-text">{{ error }}</span>
              <button class="alert-close" @click="error = ''" type="button" aria-label="Dismiss">
                <svg viewBox="0 0 20 20" fill="none" width="14" height="14" aria-hidden="true">
                  <path d="M5 5l10 10M15 5l-10 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                </svg>
              </button>
            </div>
          </Transition>

          <div class="field-group">
            <div class="field-input-wrap">
              <input
                id="username"
                v-model="form.username"
                type="text"
                required
                placeholder=" "
                autocomplete="username"
                class="field-input"
              />
              <label for="username" class="field-label">
                <svg class="field-icon" viewBox="0 0 20 20" fill="none" width="16" height="16" aria-hidden="true">
                  <circle cx="10" cy="7" r="3.5" stroke="currentColor" stroke-width="1.5" />
                  <path d="M3 18c0-3.5 3.5-6 7-6s7 2.5 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                </svg>
                Username
              </label>
            </div>
          </div>

          <div class="field-group">
            <div class="field-input-wrap">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                placeholder=" "
                autocomplete="current-password"
                class="field-input"
              />
              <label for="password" class="field-label">
                <svg class="field-icon" viewBox="0 0 20 20" fill="none" width="16" height="16" aria-hidden="true">
                  <rect x="3" y="9" width="14" height="8" rx="2" stroke="currentColor" stroke-width="1.5" />
                  <path d="M6 9V6a4 4 0 118 0v3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                </svg>
                Password
              </label>
              <button
                class="field-toggle"
                type="button"
                @click="showPassword = !showPassword"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
              >
                <svg v-if="!showPassword" viewBox="0 0 20 20" fill="none" width="18" height="18" aria-hidden="true">
                  <path d="M2 10s2.5-5.5 8-5.5 8 5.5 8 5.5-2.5 5.5-8 5.5-8-5.5-8-5.5z" stroke="currentColor" stroke-width="1.5" />
                  <circle cx="10" cy="10" r="2.5" stroke="currentColor" stroke-width="1.5" />
                </svg>
                <svg v-else viewBox="0 0 20 20" fill="none" width="18" height="18" aria-hidden="true">
                  <path d="M4 4l12 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                  <path d="M2 10s2.5-5.5 8-5.5a8.3 8.3 0 013.5.8M16.5 7.5C17.8 8.8 18 10 18 10s-2.5 5.5-8 5.5a8.3 8.3 0 01-3.5-.8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                  <circle cx="10" cy="10" r="2.5" stroke="currentColor" stroke-width="1.5" />
                </svg>
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="submit-btn"
            :class="{ 'submit-btn--loading': loading }"
          >
            <span v-if="!loading" class="submit-text">Sign in</span>
            <span v-else class="submit-loader">
              <span class="spinner" />
              <span>Signing in...</span>
            </span>
          </button>
        </form>

        <p class="form-footer">
          &copy; 2026 Tunlify
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-root {
  display: flex;
  min-height: 100vh;
  background: #070809;
  font-family: 'Sora', sans-serif;
}

/* ─── Brand Panel ─── */

.brand-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 64px;
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(ellipse 80% 60% at 30% 50%, rgba(245, 158, 11, 0.04) 0%, transparent 70%),
    radial-gradient(ellipse 60% 80% at 70% 30%, rgba(234, 88, 12, 0.03) 0%, transparent 60%),
    linear-gradient(160deg, #070809 0%, #0a0b0e 50%, #090a0c 100%);
}

.brand-scene {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.glow {
  position: absolute;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(245, 158, 11, 0.08) 0%, transparent 60%);
  animation: login-pulse 6s ease-in-out infinite;
}

@keyframes login-pulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.1); opacity: 1; }
}

.ring {
  position: absolute;
  border-radius: 50%;
  will-change: transform;
}

.ring--a {
  width: 420px;
  height: 420px;
  border: 1px solid rgba(245, 158, 11, 0.06);
  animation: login-spin 40s linear infinite;
}

.ring--b {
  width: 320px;
  height: 320px;
  background: conic-gradient(from 0deg, transparent 20%, rgba(245, 158, 11, 0.08) 50%, transparent 80%);
  mask: radial-gradient(circle, transparent 47%, black 48%, black 50%, transparent 51%);
  -webkit-mask: radial-gradient(circle, transparent 47%, black 48%, black 50%, transparent 51%);
  animation: login-spin-reverse 30s linear infinite;
}

.ring--c {
  width: 220px;
  height: 220px;
  border: 2px dashed rgba(245, 158, 11, 0.08);
  animation: login-spin 25s linear infinite;
}

.ring--d {
  width: 120px;
  height: 120px;
  background: conic-gradient(from 0deg, transparent, rgba(234, 88, 12, 0.12), transparent, rgba(245, 158, 11, 0.06), transparent);
  mask: radial-gradient(circle, transparent 40%, black 41%, black 44%, transparent 45%);
  -webkit-mask: radial-gradient(circle, transparent 40%, black 41%, black 44%, transparent 45%);
  animation: login-spin-reverse 20s linear infinite;
}

@keyframes login-spin {
  to { transform: rotate(360deg); }
}

@keyframes login-spin-reverse {
  to { transform: rotate(-360deg); }
}

.dot {
  position: absolute;
  border-radius: 50%;
  background: rgba(245, 158, 11, 0.5);
  box-shadow: 0 0 6px rgba(245, 158, 11, 0.2);
}

.dot--1 { width: 3px; height: 3px; top: 18%; left: 25%; animation: login-float 8s ease-in-out infinite; }
.dot--2 { width: 4px; height: 4px; top: 65%; left: 15%; animation: login-float 11s ease-in-out infinite 1s; }
.dot--3 { width: 2px; height: 2px; top: 30%; right: 20%; animation: login-float 9s ease-in-out infinite 2s; }
.dot--4 { width: 4px; height: 4px; bottom: 25%; right: 30%; animation: login-float 12s ease-in-out infinite 0.5s; }
.dot--5 { width: 3px; height: 3px; top: 55%; left: 60%; animation: login-float 10s ease-in-out infinite 3s; }
.dot--6 { width: 2px; height: 2px; top: 15%; right: 35%; animation: login-float 7s ease-in-out infinite 1.5s; }

@keyframes login-float {
  0%, 100% { transform: translateY(0) scale(1); opacity: 0.4; }
  50% { transform: translateY(-20px) scale(1.3); opacity: 1; }
}

.grid-mesh {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.012) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.012) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse at 40% 50%, black 15%, transparent 65%);
  -webkit-mask-image: radial-gradient(ellipse at 40% 50%, black 15%, transparent 65%);
}

.brand-body {
  position: relative;
  z-index: 1;
  max-width: 480px;
  animation: login-fade-up 0.8s ease-out both;
}

.brand-logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 32px;
}

.brand-icon-frame {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  border: 1px solid rgba(245, 158, 11, 0.2);
  background: rgba(245, 158, 11, 0.05);
}

.brand-icon {
  width: 22px;
  height: 22px;
}

.brand-wordmark {
  font-weight: 700;
  font-size: 20px;
  color: #F1F5F9;
  letter-spacing: -0.01em;
}

.brand-headline {
  font-weight: 700;
  font-size: 38px;
  line-height: 1.15;
  color: #F1F5F9;
  margin: 0 0 16px;
  letter-spacing: -0.02em;
}

.text-gradient {
  background: linear-gradient(135deg, #F59E0B, #EA580C);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.brand-description {
  font-size: 15px;
  line-height: 1.7;
  color: #71717A;
  margin: 0;
  font-weight: 400;
}

.brand-footer {
  position: relative;
  z-index: 1;
  margin-top: 48px;
  animation: login-fade-up 0.8s ease-out 0.15s both;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #52525B;
  font-weight: 500;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #22C55E;
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.4);
  animation: login-pulse 3s ease-in-out infinite;
}

/* ─── Form Panel ─── */

.form-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 64px;
  position: relative;
  background:
    radial-gradient(ellipse 100% 80% at 50% 60%, rgba(245, 158, 11, 0.02) 0%, transparent 70%),
    linear-gradient(180deg, #0b0c10 0%, #0d0e12 100%);
}

.form-body {
  width: 100%;
  max-width: 400px;
  animation: login-fade-up 0.8s ease-out 0.1s both;
}

.form-header {
  margin-bottom: 32px;
}

.form-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: rgba(245, 158, 11, 0.8);
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.12);
  margin-bottom: 20px;
}

.form-title {
  font-weight: 700;
  font-size: 28px;
  color: #F1F5F9;
  margin: 0 0 8px;
  letter-spacing: -0.02em;
}

.form-subtitle {
  font-size: 14px;
  color: #71717A;
  margin: 0;
  line-height: 1.5;
  font-weight: 400;
}

/* ─── Alert ─── */

.form-alert {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.15);
  color: #FCA5A5;
  font-size: 13px;
  margin-bottom: 20px;
  font-weight: 500;
}

.alert-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.alert-text {
  flex: 1;
}

.alert-close {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  color: rgba(252, 165, 165, 0.5);
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s, color 0.15s;
  padding: 0;
}

.alert-close:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #FCA5A5;
}

.alert-enter-active {
  transition: all 0.25s ease-out;
}

.alert-leave-active {
  transition: all 0.2s ease-in;
}

.alert-enter-from,
.alert-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ─── Form Fields ─── */

.form-element {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field-group {
  display: flex;
  flex-direction: column;
}

.field-input-wrap {
  position: relative;
}

.field-input {
  width: 100%;
  padding: 22px 16px 8px 44px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 12px;
  color: #F1F5F9;
  font-size: 15px;
  font-family: 'Sora', sans-serif;
  font-weight: 400;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
  box-sizing: border-box;
}

.field-input:hover {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
}

.field-input:focus {
  border-color: rgba(245, 158, 11, 0.35);
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.06), 0 0 24px -10px rgba(245, 158, 11, 0.15);
  background: rgba(255, 255, 255, 0.05);
}

.field-input::placeholder {
  color: transparent;
}

.field-label {
  position: absolute;
  left: 44px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
  color: #52525B;
  font-weight: 400;
  pointer-events: none;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
}

.field-icon {
  color: rgba(255, 255, 255, 0.15);
  flex-shrink: 0;
  transition: color 0.2s;
}

.field-input:focus ~ .field-label,
.field-input:not(:placeholder-shown) ~ .field-label {
  top: 12px;
  transform: translateY(0);
  font-size: 11px;
  color: rgba(245, 158, 11, 0.7);
  font-weight: 500;
}

.field-input:focus ~ .field-label .field-icon,
.field-input:not(:placeholder-shown) ~ .field-label .field-icon {
  color: rgba(245, 158, 11, 0.5);
}

.field-toggle {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: #52525B;
  cursor: pointer;
  border-radius: 8px;
  transition: color 0.2s, background 0.15s;
  padding: 0;
}

.field-toggle:hover {
  color: #A1A1AA;
  background: rgba(255, 255, 255, 0.05);
}

/* ─── Submit Button ─── */

.submit-btn {
  width: 100%;
  height: 48px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #F59E0B, #EA580C);
  color: white;
  font-family: 'Sora', sans-serif;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: box-shadow 0.3s, transform 0.15s, opacity 0.2s;
  position: relative;
  overflow: hidden;
  margin-top: 8px;
}

.submit-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #FBBF24, #EA580C);
  opacity: 0;
  transition: opacity 0.3s;
  border-radius: 12px;
}

.submit-btn:hover:not(:disabled) {
  box-shadow: 0 0 32px -6px rgba(245, 158, 11, 0.4);
  transform: translateY(-1px);
}

.submit-btn:hover:not(:disabled)::before {
  opacity: 1;
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0) scale(0.99);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.submit-text,
.submit-loader {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: login-spin 0.6s linear infinite;
}

/* ─── Footer ─── */

.form-footer {
  margin-top: 48px;
  font-size: 12px;
  color: #3F3F46;
  text-align: center;
  font-weight: 400;
  letter-spacing: 0.02em;
}

/* ─── Animations ─── */

@keyframes login-fade-up {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ─── Responsive ─── */

@media (max-width: 900px) {
  .login-root {
    flex-direction: column;
  }

  .brand-panel {
    min-height: 35vh;
    padding: 40px 32px 32px;
    justify-content: flex-end;
  }

  .brand-scene {
    opacity: 0.5;
  }

  .brand-headline {
    font-size: 28px;
  }

  .brand-body {
    max-width: 100%;
  }

  .brand-footer {
    margin-top: 24px;
  }

  .form-panel {
    padding: 32px;
  }

  .form-body {
    max-width: 100%;
  }
}

@media (max-width: 480px) {
  .brand-panel {
    min-height: 30vh;
    padding: 32px 20px 24px;
  }

  .brand-headline {
    font-size: 24px;
  }

  .brand-logo {
    margin-bottom: 20px;
  }

  .form-panel {
    padding: 24px 20px;
  }

  .form-title {
    font-size: 24px;
  }

  .field-input {
    padding-left: 40px;
    padding-top: 20px;
    padding-bottom: 8px;
    font-size: 14px;
  }

  .field-label {
    left: 40px;
    font-size: 13px;
  }
}
</style>
