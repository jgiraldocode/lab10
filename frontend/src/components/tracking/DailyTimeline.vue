<script setup lang="ts">
import { computed } from 'vue'
import type { TimelineEntry } from '@/lib/api'

const props = defineProps<{
  entries: TimelineEntry[]
}>()

const colors = [
  'bg-blue-500', 'bg-emerald-500', 'bg-amber-500', 'bg-purple-500',
  'bg-rose-500', 'bg-cyan-500', 'bg-indigo-500', 'bg-orange-500',
  'bg-teal-500', 'bg-pink-500', 'bg-lime-500', 'bg-fuchsia-500',
]

function hashColor(bundleId: string): string {
  let hash = 0
  for (let i = 0; i < bundleId.length; i++) {
    hash = ((hash << 5) - hash) + bundleId.charCodeAt(i)
    hash |= 0
  }
  return colors[Math.abs(hash) % colors.length]
}

function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  if (m < 60) return `${m}m ${s}s`
  const h = Math.floor(m / 60)
  return `${h}h ${m % 60}m`
}

function formatTime(iso: string): string {
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const legend = computed(() => {
  const map = new Map<string, { name: string; color: string }>()
  for (const e of props.entries) {
    if (!map.has(e.bundle_id)) {
      map.set(e.bundle_id, { name: e.app_name, color: hashColor(e.bundle_id) })
    }
  }
  return Array.from(map.values())
})
</script>

<template>
  <div>
    <div v-if="entries.length === 0" class="py-12 text-center text-sm text-zinc-500">
      No activity recorded yet
    </div>
    <div v-else>
      <!-- Timeline bar -->
      <div class="mb-4 flex h-8 overflow-hidden rounded-lg border border-zinc-200 dark:border-zinc-800">
        <div
          v-for="(entry, i) in entries"
          :key="i"
          :class="['relative transition-all', hashColor(entry.bundle_id)]"
          :style="{ width: `${Math.max(entry.percent, 0.5)}%` }"
          :title="`${entry.app_name}${entry.tab_host ? ' (' + entry.tab_host + ')' : ''} - ${formatDuration(entry.seconds)}`"
        />
      </div>

      <!-- Legend -->
      <div class="mb-6 flex flex-wrap gap-3">
        <div v-for="item in legend" :key="item.name" class="flex items-center gap-1.5 text-xs text-zinc-600 dark:text-zinc-400">
          <div :class="['h-2.5 w-2.5 rounded-full', item.color]" />
          {{ item.name }}
        </div>
      </div>

      <!-- Segment list -->
      <div class="space-y-1">
        <div
          v-for="(entry, i) in entries"
          :key="i"
          class="flex items-center justify-between rounded-lg px-3 py-2 text-sm hover:bg-zinc-50 dark:hover:bg-zinc-900"
        >
          <div class="flex items-center gap-3">
            <div :class="['h-2.5 w-2.5 rounded-full', hashColor(entry.bundle_id)]" />
            <span class="font-medium text-zinc-900 dark:text-zinc-50">{{ entry.app_name }}</span>
            <span v-if="entry.tab_host" class="text-zinc-500">{{ entry.tab_host }}</span>
          </div>
          <div class="flex items-center gap-4 text-xs text-zinc-500">
            <span>{{ formatTime(entry.start_time) }} - {{ formatTime(entry.end_time) }}</span>
            <span class="font-medium text-zinc-700 dark:text-zinc-300">{{ formatDuration(entry.seconds) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
