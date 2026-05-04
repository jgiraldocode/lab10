<script setup lang="ts">
defineProps<{
  value: number
  max?: number
  color?: string
}>()
</script>

<template>
  <div class="h-2 w-full overflow-hidden rounded-full bg-zinc-200 dark:bg-zinc-800">
    <div
      class="h-full rounded-full transition-all duration-300"
      :class="color || 'bg-zinc-900 dark:bg-zinc-50'"
      :style="{ width: `${Math.min((value / (max || 100)) * 100, 100)}%` }"
    />
  </div>
</template>
