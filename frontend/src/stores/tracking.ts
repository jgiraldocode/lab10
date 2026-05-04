import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getSegments, getTimeline, getApps } from '@/lib/api'
import type { TimeSegment, TimelineEntry, AppSummary } from '@/lib/api'

export const useTrackingStore = defineStore('tracking', () => {
  const segments = ref<TimeSegment[]>([])
  const timeline = ref<TimelineEntry[]>([])
  const apps = ref<AppSummary[]>([])
  const selectedDate = ref(new Date().toISOString().slice(0, 10))
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalSeconds = computed(() =>
    segments.value.reduce((sum, s) => sum + s.seconds, 0)
  )

  const topApps = computed(() => apps.value.slice(0, 10))

  async function fetchSegments(date?: string) {
    const d = date || selectedDate.value
    loading.value = true
    error.value = null
    try {
      segments.value = (await getSegments(d)) || []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error'
    } finally {
      loading.value = false
    }
  }

  async function fetchTimeline(date?: string) {
    const d = date || selectedDate.value
    try {
      timeline.value = (await getTimeline(d)) || []
    } catch {
      // silent
    }
  }

  async function fetchApps(date?: string) {
    const d = date || selectedDate.value
    try {
      apps.value = (await getApps(d)) || []
    } catch {
      // silent
    }
  }

  async function fetchAll(date?: string) {
    const d = date || selectedDate.value
    selectedDate.value = d
    await Promise.all([fetchSegments(d), fetchTimeline(d), fetchApps(d)])
  }

  return { segments, timeline, apps, selectedDate, loading, error, totalSeconds, topApps, fetchSegments, fetchTimeline, fetchApps, fetchAll }
})
