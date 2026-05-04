import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'dashboard', component: DashboardView },
    { path: '/pomodoro', name: 'pomodoro', component: () => import('@/views/PomodoroView.vue') },
    { path: '/stats', name: 'stats', component: () => import('@/views/StatsView.vue') },
    { path: '/export', name: 'export', component: () => import('@/views/ExportView.vue') },
  ],
})

export default router
