import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePomodoroStore } from '../pomodoro'

const mockState = {
  active: true,
  session_id: 1,
  phase: 'work',
  remaining_seconds: 1480,
  total_seconds: 1500,
  violation: null as any,
}

vi.mock('@/lib/api', () => ({
  startPomodoro: vi.fn().mockResolvedValue({ id: 1, status: 'active' }),
  stopPomodoro: vi.fn().mockResolvedValue(undefined),
  getPomodoroState: vi.fn().mockImplementation(() => Promise.resolve({ ...mockState })),
  getPomodoroSessions: vi.fn().mockResolvedValue([
    { id: 1, started_at: '2024-01-01T10:00:00Z', work_minutes: 25, break_minutes: 5, status: 'completed', current_phase: 'work' },
  ]),
  handleViolation: vi.fn().mockResolvedValue(undefined),
  getAllowlist: vi.fn().mockResolvedValue([
    { id: 1, bundle_id: 'com.apple.Terminal', app_name: 'Terminal', is_default: true },
  ]),
  addToAllowlist: vi.fn().mockResolvedValue({ id: 2, bundle_id: 'com.google.Chrome', app_name: 'Chrome', is_default: false }),
  removeFromAllowlist: vi.fn().mockResolvedValue(undefined),
}))

describe('usePomodoroStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initializes inactive', () => {
    const store = usePomodoroStore()
    expect(store.state.active).toBe(false)
    expect(store.hasViolation).toBe(false)
    expect(store.formattedTime).toBe('00:00')
  })

  it('starts a session and begins polling', async () => {
    const store = usePomodoroStore()
    await store.start(25, 5)
    expect(store.polling).toBe(true)
    store.stopPolling()
  })

  it('formats time correctly', async () => {
    const store = usePomodoroStore()
    await store.fetchState()
    expect(store.formattedTime).toBe('24:40')
  })

  it('computes progress', async () => {
    const store = usePomodoroStore()
    await store.fetchState()
    const expected = ((1500 - 1480) / 1500) * 100
    expect(store.progress).toBeCloseTo(expected, 1)
  })

  it('fetches sessions', async () => {
    const store = usePomodoroStore()
    await store.fetchSessions()
    expect(store.sessions).toHaveLength(1)
    expect(store.sessions[0].status).toBe('completed')
  })

  it('fetches allowlist', async () => {
    const store = usePomodoroStore()
    await store.fetchAllowlist()
    expect(store.allowlist).toHaveLength(1)
    expect(store.allowlist[0].app_name).toBe('Terminal')
  })

  it('detects violations', async () => {
    const store = usePomodoroStore()
    mockState.violation = { id: 1, session_id: 1, bundle_id: 'com.spotify.client', app_name: 'Spotify', timestamp: '2024-01-01T10:05:00Z', action: '' }
    await store.fetchState()
    expect(store.hasViolation).toBe(true)
    expect(store.state.violation?.app_name).toBe('Spotify')
    mockState.violation = null
  })

  it('stops session', async () => {
    const store = usePomodoroStore()
    await store.stop()
    expect(store.state.active).toBe(false)
  })
})
