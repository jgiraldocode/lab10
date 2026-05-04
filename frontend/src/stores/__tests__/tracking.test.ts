import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTrackingStore } from '../tracking'

vi.mock('@/lib/api', () => ({
  getSegments: vi.fn().mockResolvedValue([
    { id: 1, start_time: '2024-01-01T10:00:00Z', end_time: '2024-01-01T10:05:00Z', bundle_id: 'com.app.test', app_name: 'Test', tab_host: '', tab_title: '', seconds: 300, pomodoro_session_id: null },
    { id: 2, start_time: '2024-01-01T10:05:00Z', end_time: '2024-01-01T10:10:00Z', bundle_id: 'com.google.Chrome', app_name: 'Chrome', tab_host: 'github.com', tab_title: 'GitHub', seconds: 300, pomodoro_session_id: null },
  ]),
  getTimeline: vi.fn().mockResolvedValue([
    { start_time: '2024-01-01T10:00:00Z', end_time: '2024-01-01T10:05:00Z', app_name: 'Test', bundle_id: 'com.app.test', tab_host: '', seconds: 300, percent: 50 },
    { start_time: '2024-01-01T10:05:00Z', end_time: '2024-01-01T10:10:00Z', app_name: 'Chrome', bundle_id: 'com.google.Chrome', tab_host: 'github.com', seconds: 300, percent: 50 },
  ]),
  getApps: vi.fn().mockResolvedValue([
    { bundle_id: 'com.app.test', app_name: 'Test', total_seconds: 300 },
    { bundle_id: 'com.google.Chrome', app_name: 'Chrome', total_seconds: 300 },
  ]),
}))

describe('useTrackingStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with empty state', () => {
    const store = useTrackingStore()
    expect(store.segments).toEqual([])
    expect(store.timeline).toEqual([])
    expect(store.apps).toEqual([])
    expect(store.totalSeconds).toBe(0)
  })

  it('fetches segments', async () => {
    const store = useTrackingStore()
    await store.fetchSegments('2024-01-01')
    expect(store.segments).toHaveLength(2)
    expect(store.totalSeconds).toBe(600)
  })

  it('fetches all data', async () => {
    const store = useTrackingStore()
    await store.fetchAll('2024-01-01')
    expect(store.segments).toHaveLength(2)
    expect(store.timeline).toHaveLength(2)
    expect(store.apps).toHaveLength(2)
    expect(store.topApps).toHaveLength(2)
    expect(store.selectedDate).toBe('2024-01-01')
  })
})
