<script setup lang="ts">
import { ref } from 'vue'
import { usePomodoroStore } from '@/stores/pomodoro'
import Button from '@/components/ui/Button.vue'

const store = usePomodoroStore()
const workMinutes = ref(25)
const breakMinutes = ref(5)

function handleStart() {
  store.start(workMinutes.value, breakMinutes.value)
}
</script>

<template>
  <div class="flex flex-col items-center">
    <!-- Timer display -->
    <div class="relative mb-6 flex h-48 w-48 items-center justify-center">
      <svg class="absolute h-full w-full -rotate-90" viewBox="0 0 100 100">
        <circle cx="50" cy="50" r="45" fill="none" stroke-width="3" class="stroke-zinc-200 dark:stroke-zinc-800" />
        <circle
          cx="50" cy="50" r="45" fill="none" stroke-width="3"
          stroke-linecap="round"
          class="transition-all duration-1000"
          :class="store.state.phase === 'break' ? 'stroke-emerald-500' : 'stroke-zinc-900 dark:stroke-zinc-50'"
          :stroke-dasharray="283"
          :stroke-dashoffset="283 - (283 * store.progress / 100)"
        />
      </svg>
      <div class="text-center">
        <div class="text-3xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
          {{ store.state.active ? store.formattedTime : `${workMinutes}:00` }}
        </div>
        <div v-if="store.state.active" class="mt-1 text-xs font-medium uppercase tracking-wide"
          :class="store.state.phase === 'work' ? 'text-zinc-500' : 'text-emerald-600'">
          {{ store.state.phase }}
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div v-if="!store.state.active" class="mb-6 flex items-center gap-4">
      <label class="text-sm text-zinc-600 dark:text-zinc-400">
        Work
        <input v-model.number="workMinutes" type="number" min="1" max="120"
          class="ml-2 w-16 rounded-lg border border-zinc-200 bg-white px-2 py-1 text-sm dark:border-zinc-800 dark:bg-zinc-950" />
        min
      </label>
      <label class="text-sm text-zinc-600 dark:text-zinc-400">
        Break
        <input v-model.number="breakMinutes" type="number" min="1" max="30"
          class="ml-2 w-16 rounded-lg border border-zinc-200 bg-white px-2 py-1 text-sm dark:border-zinc-800 dark:bg-zinc-950" />
        min
      </label>
    </div>

    <div class="flex gap-3">
      <Button v-if="!store.state.active" @click="handleStart">Start Focus</Button>
      <Button v-else variant="destructive" @click="store.stop()">Stop</Button>
    </div>

    <p v-if="store.error" class="mt-3 text-sm text-red-600">{{ store.error }}</p>
  </div>
</template>
