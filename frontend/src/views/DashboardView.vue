<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useTrackingStore } from '@/stores/tracking'
import { usePomodoroStore } from '@/stores/pomodoro'
import DailyTimeline from '@/components/tracking/DailyTimeline.vue'
import AppUsageCard from '@/components/tracking/AppUsageCard.vue'
import StatsOverview from '@/components/stats/StatsOverview.vue'
import FilterBar from '@/components/stats/FilterBar.vue'
import Card from '@/components/ui/Card.vue'

const tracking = useTrackingStore()
const pomodoro = usePomodoroStore()

let refreshInterval: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  await tracking.fetchAll()
  await pomodoro.fetchState()
  if (pomodoro.state.active) {
    pomodoro.startPolling()
  }
  refreshInterval = setInterval(() => tracking.fetchAll(), 5000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})

function onDateChange(date: string) {
  tracking.fetchAll(date)
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">Dashboard</h2>
      <FilterBar :model-value="tracking.selectedDate" @update:model-value="onDateChange" />
    </div>

    <StatsOverview
      :total-seconds="tracking.totalSeconds"
      :app-count="tracking.apps.length"
      :segment-count="tracking.segments.length"
    />

    <!-- Mini Pomodoro status -->
    <Card v-if="pomodoro.state.active">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="h-2.5 w-2.5 animate-pulse rounded-full"
            :class="pomodoro.state.phase === 'work' ? 'bg-zinc-900 dark:bg-zinc-50' : 'bg-emerald-500'" />
          <span class="text-sm font-medium text-zinc-900 dark:text-zinc-50">
            Pomodoro {{ pomodoro.state.phase }}
          </span>
        </div>
        <span class="text-lg font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
          {{ pomodoro.formattedTime }}
        </span>
      </div>
    </Card>

    <div class="grid gap-6 lg:grid-cols-3">
      <div class="lg:col-span-2">
        <Card>
          <h3 class="mb-4 text-sm font-semibold text-zinc-900 dark:text-zinc-50">Timeline</h3>
          <DailyTimeline :entries="tracking.timeline" />
        </Card>
      </div>
      <AppUsageCard :apps="tracking.topApps" />
    </div>

    <p v-if="tracking.loading" class="text-center text-sm text-zinc-500">Loading...</p>
    <p v-if="tracking.error" class="text-center text-sm text-red-600">{{ tracking.error }}</p>
  </div>
</template>
