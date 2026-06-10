import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

// reportError surfaces a swallowed error to the console with consistent context,
// so failed background fetches are no longer silently lost.
export function reportError(context: string, err: unknown): void {
  const message = err instanceof Error ? err.message : String(err)
  console.error(`[Tunlify] ${context}: ${message}`)
}

export const useApi = () => {
  const authStore = useAuthStore()
  const router = useRouter()

  const apiFetch = async <T>(path: string, opts: RequestInit = {}): Promise<T> => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(opts.headers as Record<string, string>),
    }

    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`
    }

    const res = await fetch(path, { ...opts, headers })

    if (res.status === 401) {
      authStore.logout()
      router.push('/login')
      throw new Error('Unauthorized')
    }

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }))
      throw new Error(err.error || 'Request failed')
    }

    return res.json()
  }

  return { apiFetch }
}
