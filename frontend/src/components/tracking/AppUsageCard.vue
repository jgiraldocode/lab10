<script setup lang="ts">
import Card from '@/components/ui/Card.vue'
import type { AppSummary } from '@/lib/api'

defineProps<{
  apps: AppSummary[]
}>()

function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  const m = Math.floor(seconds / 60)
  if (m < 60) return `${m}m`
  const h = Math.floor(m / 60)
  return `${h}h ${m % 60}m`
}
</script>

<template>
  <Card>
    <h3 class="mb-4 text-sm font-semibold text-zinc-900 dark:text-zinc-50">Top Apps</h3>
    <div v-if="apps.length === 0" class="text-sm text-zinc-500">No data</div>
    <ul class="space-y-3">
      <li v-for="app in apps" :key="app.bundle_id" class="flex items-center justify-between">
        <span class="text-sm text-zinc-700 dark:text-zinc-300">{{ app.app_name || app.bundle_id }}</span>
        <span class="text-sm font-medium text-zinc-900 dark:text-zinc-50">{{ formatDuration(app.total_seconds) }}</span>
      </li>
    </ul>
  </Card>
</template>
