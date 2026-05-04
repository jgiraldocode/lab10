<script setup lang="ts">
defineProps<{
  open: boolean
  title?: string
}>()

defineEmits<{
  close: []
}>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/50" @click="$emit('close')" />
        <div class="relative z-10 w-full max-w-md rounded-xl border border-zinc-200 bg-white p-6 shadow-lg dark:border-zinc-800 dark:bg-zinc-950">
          <h2 v-if="title" class="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">{{ title }}</h2>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
