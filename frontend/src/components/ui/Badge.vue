<script setup lang="ts">
defineProps<{
  variant?: 'default' | 'success' | 'warning' | 'error'
}>()
</script>

<template>
  <span
    :class="[
      'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
      {
        'bg-zinc-100 text-zinc-800 dark:bg-zinc-800 dark:text-zinc-200': !variant || variant === 'default',
        'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': variant === 'success',
        'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': variant === 'warning',
        'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': variant === 'error',
      }
    ]"
  >
    <slot />
  </span>
</template>
