import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('../views/LoginView.vue') },
    { path: '/', component: () => import('../views/DashboardView.vue') },
    { path: '/logs', component: () => import('../views/LogsView.vue') },
    { path: '/config', component: () => import('../views/ConfigView.vue') },
    { path: '/health', component: () => import('../views/HealthView.vue') },
    { path: '/metrics', component: () => import('../views/MetricsView.vue') },
    { path: '/backups', component: () => import('../views/BackupsView.vue') },
    { path: '/tcp-access', component: () => import('../views/TcpAccessView.vue') },
    { path: '/audit', component: () => import('../views/AuditView.vue') },
    { path: '/notifications', component: () => import('../views/NotificationsView.vue') },
    { path: '/projects', component: () => import('../views/ProjectsView.vue') },
    { path: '/settings', component: () => import('../views/SettingsView.vue') },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.path !== '/login' && !auth.isAuthenticated) return '/login'
})

export default router
