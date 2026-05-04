const BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

// Types
export interface SampleEvent {
  id: number
  timestamp: string
  bundle_id: string
  app_name: string
  window_title: string
  browser_family: string
  tab_title: string
  tab_url_host: string
  source: string
  confidence: string
}

export interface TimeSegment {
  id: number
  start_time: string
  end_time: string
  bundle_id: string
  app_name: string
  tab_host: string
  tab_title: string
  seconds: number
  pomodoro_session_id?: number
}

export interface TimelineEntry {
  start_time: string
  end_time: string
  app_name: string
  bundle_id: string
  tab_host: string
  seconds: number
  percent: number
}

export interface AppSummary {
  bundle_id: string
  app_name: string
  total_seconds: number
}

export interface PomodoroSession {
  id: number
  started_at: string
  ended_at?: string
  work_minutes: number
  break_minutes: number
  status: string
  current_phase: string
}

export interface PomodoroState {
  active: boolean
  session_id: number
  phase: string
  remaining_seconds: number
  total_seconds: number
  violation?: FocusViolation
}

export interface FocusViolation {
  id: number
  session_id: number
  timestamp: string
  bundle_id: string
  app_name: string
  action: string
}

export interface AllowlistEntry {
  id: number
  bundle_id: string
  app_name: string
  is_default: boolean
  session_id?: number
}

// Tracking
export const getEvents = (date: string) =>
  request<SampleEvent[]>(`/events?date=${date}`)

export const getSegments = (date: string) =>
  request<TimeSegment[]>(`/segments?date=${date}`)

export const getTimeline = (date: string) =>
  request<TimelineEntry[]>(`/timeline?date=${date}`)

export const getApps = (date: string) =>
  request<AppSummary[]>(`/apps?date=${date}`)

// Pomodoro
export const startPomodoro = (workMinutes: number, breakMinutes: number) =>
  request<PomodoroSession>('/pomodoro/start', {
    method: 'POST',
    body: JSON.stringify({ work_minutes: workMinutes, break_minutes: breakMinutes }),
  })

export const stopPomodoro = () =>
  request<void>('/pomodoro/stop', { method: 'POST' })

export const getPomodoroState = () =>
  request<PomodoroState>('/pomodoro/state')

export const getPomodoroSessions = () =>
  request<PomodoroSession[]>('/pomodoro/sessions')

export const handleViolation = (action: string) =>
  request<void>('/pomodoro/violations', {
    method: 'POST',
    body: JSON.stringify({ action }),
  })

// Allowlist
export const getAllowlist = (sessionId?: number) =>
  request<AllowlistEntry[]>(`/allowlist${sessionId ? `?session_id=${sessionId}` : ''}`)

export const addToAllowlist = (entry: { bundle_id: string; app_name: string; is_default: boolean; session_id?: number }) =>
  request<AllowlistEntry>('/allowlist', {
    method: 'POST',
    body: JSON.stringify(entry),
  })

export const removeFromAllowlist = (id: number) =>
  request<void>(`/allowlist/${id}`, { method: 'DELETE' })

// Export
export const exportCSV = (date: string) =>
  fetch(`${BASE}/export/csv?date=${date}`).then(r => r.blob())

export const exportJSON = (date: string) =>
  fetch(`${BASE}/export/json?date=${date}`).then(r => r.blob())
