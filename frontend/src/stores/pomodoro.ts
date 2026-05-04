import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  startPomodoro,
  stopPomodoro,
  getPomodoroState,
  getPomodoroSessions,
  handleViolation,
  getAllowlist,
  addToAllowlist,
  removeFromAllowlist,
} from '@/lib/api'
import type { PomodoroState, PomodoroSession, AllowlistEntry } from '@/lib/api'

export const usePomodoroStore = defineStore('pomodoro', () => {
  const state = ref<PomodoroState>({
    active: false,
    session_id: 0,
    phase: '',
    remaining_seconds: 0,
    total_seconds: 0,
  })
  const sessions = ref<PomodoroSession[]>([])
  const allowlist = ref<AllowlistEntry[]>([])
  const polling = ref(false)
  const error = ref<string | null>(null)

  let pollInterval: ReturnType<typeof setInterval> | null = null

  const hasViolation = computed(() => !!state.value.violation)

  const formattedTime = computed(() => {
    const secs = state.value.remaining_seconds
    const m = Math.floor(secs / 60)
    const s = secs % 60
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  })

  const progress = computed(() => {
    if (state.value.total_seconds === 0) return 0
    return ((state.value.total_seconds - state.value.remaining_seconds) / state.value.total_seconds) * 100
  })

  async function start(workMinutes = 25, breakMinutes = 5) {
    error.value = null
    try {
      await startPomodoro(workMinutes, breakMinutes)
      startPolling()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error'
    }
  }

  async function stop() {
    try {
      await stopPomodoro()
      stopPolling()
      state.value = { active: false, session_id: 0, phase: '', remaining_seconds: 0, total_seconds: 0 }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error'
    }
  }

  async function fetchState() {
    try {
      state.value = await getPomodoroState()
      if (!state.value.active && polling.value) {
        stopPolling()
      }
    } catch {
      // silent
    }
  }

  function startPolling() {
    if (pollInterval) return
    polling.value = true
    fetchState()
    pollInterval = setInterval(fetchState, 1000)
  }

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval)
      pollInterval = null
    }
    polling.value = false
  }

  async function fetchSessions() {
    try {
      sessions.value = (await getPomodoroSessions()) || []
    } catch {
      // silent
    }
  }

  async function fetchAllowlist(sessionId?: number) {
    try {
      allowlist.value = (await getAllowlist(sessionId)) || []
    } catch {
      // silent
    }
  }

  async function addApp(bundleId: string, appName: string, isDefault: boolean, sessionId?: number) {
    await addToAllowlist({ bundle_id: bundleId, app_name: appName, is_default: isDefault, session_id: sessionId })
    await fetchAllowlist(sessionId)
  }

  async function removeApp(id: number) {
    await removeFromAllowlist(id)
    const sid = state.value.active ? state.value.session_id : undefined
    await fetchAllowlist(sid)
  }

  async function dismissViolation() {
    await handleViolation('returned')
    if (state.value.violation) {
      state.value = { ...state.value, violation: undefined }
    }
  }

  async function addViolatingAppToAllowlist() {
    const v = state.value.violation
    if (!v) return
    await addToAllowlist({
      bundle_id: v.bundle_id,
      app_name: v.app_name,
      is_default: false,
      session_id: state.value.session_id,
    })
    await handleViolation('added_to_allowlist')
    state.value = { ...state.value, violation: undefined }
    await fetchAllowlist(state.value.session_id)
  }

  return {
    state, sessions, allowlist, polling, error,
    hasViolation, formattedTime, progress,
    start, stop, fetchState, startPolling, stopPolling,
    fetchSessions, fetchAllowlist, addApp, removeApp,
    dismissViolation, addViolatingAppToAllowlist,
  }
})
