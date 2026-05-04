<script setup lang="ts">
import { onMounted } from 'vue'
import { usePomodoroStore } from '@/stores/pomodoro'
import PomodoroTimer from '@/components/pomodoro/PomodoroTimer.vue'
import AllowlistManager from '@/components/pomodoro/AllowlistManager.vue'
import FocusWarning from '@/components/pomodoro/FocusWarning.vue'
import SessionHistory from '@/components/pomodoro/SessionHistory.vue'
import Card from '@/components/ui/Card.vue'

const store = usePomodoroStore()

onMounted(async () => {
  await store.fetchState()
  if (store.state.active) {
    store.startPolling()
  }
})
</script>

<template>
  <div class="space-y-6">
    <h2 class="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">Pomodoro</h2>

    <Card>
      <PomodoroTimer />
    </Card>

    <div class="grid gap-6 lg:grid-cols-2">
      <AllowlistManager />
      <SessionHistory />
    </div>

    <FocusWarning />
  </div>
</template>
