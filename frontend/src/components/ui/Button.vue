<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'secondary' | 'ghost' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}>()
</script>

<template>
  <button
    :disabled="disabled"
    :class="[
      'inline-flex items-center justify-center rounded-lg font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-400 focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none',
      {
        'bg-zinc-900 text-white hover:bg-zinc-800 dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200': variant === 'primary' || !variant,
        'border border-zinc-200 bg-white text-zinc-900 hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-950 dark:text-zinc-50 dark:hover:bg-zinc-900': variant === 'secondary',
        'text-zinc-600 hover:text-zinc-900 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:text-zinc-50 dark:hover:bg-zinc-800': variant === 'ghost',
        'bg-red-600 text-white hover:bg-red-700': variant === 'destructive',
      },
      {
        'px-3 py-1.5 text-xs': size === 'sm',
        'px-4 py-2 text-sm': size === 'md' || !size,
        'px-6 py-3 text-base': size === 'lg',
      },
    ]"
  >
    <slot />
  </button>
</template>
