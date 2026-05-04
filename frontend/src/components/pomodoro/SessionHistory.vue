<script setup lang="ts">
import { onMounted } from 'vue'
import { usePomodoroStore } from '@/stores/pomodoro'
import Card from '@/components/ui/Card.vue'
import Badge from '@/components/ui/Badge.vue'

const store = usePomodoroStore()

onMounted(() => store.fetchSessions())

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <Card>
    <h3 class="mb-4 text-sm font-semibold text-zinc-900 dark:text-zinc-50">Session History</h3>
    <div v-if="!store.sessions.length" class="text-sm text-zinc-500">No sessions yet</div>
    <ul class="space-y-2">
      <li v-for="s in store.sessions" :key="s.id"
        class="flex items-center justify-between rounded-lg px-3 py-2 text-sm hover:bg-zinc-50 dark:hover:bg-zinc-900">
        <div>
          <span class="text-zinc-700 dark:text-zinc-300">{{ formatDate(s.started_at) }}</span>
          <span class="ml-2 text-zinc-500">{{ s.work_minutes }}m work / {{ s.break_minutes }}m break</span>
        </div>
        <Badge
          :variant="s.status === 'completed' ? 'success' : s.status === 'active' ? 'warning' : 'error'">
          {{ s.status }}
        </Badge>
      </li>
    </ul>
  </Card>
</template>
