<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="flex min-h-screen bg-zinc-50 dark:bg-zinc-950">
    <Sidebar />
    <main class="flex-1 overflow-auto">
      <div class="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <slot />
      </div>
    </main>
  </div>
</template>
