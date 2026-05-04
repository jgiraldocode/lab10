<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useTrackingStore } from '@/stores/tracking'
import StatsOverview from '@/components/stats/StatsOverview.vue'
import FilterBar from '@/components/stats/FilterBar.vue'
import AppUsageCard from '@/components/tracking/AppUsageCard.vue'
import DailyTimeline from '@/components/tracking/DailyTimeline.vue'
import Card from '@/components/ui/Card.vue'

const tracking = useTrackingStore()

onMounted(() => tracking.fetchAll())

watch(() => tracking.selectedDate, (d) => tracking.fetchAll(d))

function onDateChange(date: string) {
  tracking.selectedDate = date
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">Statistics</h2>
      <FilterBar :model-value="tracking.selectedDate" @update:model-value="onDateChange" />
    </div>

    <StatsOverview
      :total-seconds="tracking.totalSeconds"
      :app-count="tracking.apps.length"
      :segment-count="tracking.segments.length"
    />

    <div class="grid gap-6 lg:grid-cols-3">
      <div class="lg:col-span-2">
        <Card>
          <h3 class="mb-4 text-sm font-semibold text-zinc-900 dark:text-zinc-50">Activity</h3>
          <DailyTimeline :entries="tracking.timeline" />
        </Card>
      </div>
      <AppUsageCard :apps="tracking.apps" />
    </div>
  </div>
</template>
