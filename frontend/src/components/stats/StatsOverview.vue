<script setup lang="ts">
import Card from '@/components/ui/Card.vue'

defineProps<{
  totalSeconds: number
  appCount: number
  segmentCount: number
}>()

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-3">
    <Card>
      <p class="text-xs font-medium text-zinc-500">Total Active Time</p>
      <p class="mt-1 text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">{{ formatDuration(totalSeconds) }}</p>
    </Card>
    <Card>
      <p class="text-xs font-medium text-zinc-500">Apps Used</p>
      <p class="mt-1 text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">{{ appCount }}</p>
    </Card>
    <Card>
      <p class="text-xs font-medium text-zinc-500">Context Switches</p>
      <p class="mt-1 text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">{{ segmentCount }}</p>
    </Card>
  </div>
</template>
