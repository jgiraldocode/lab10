<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const date = ref(props.modelValue)

watch(date, (v) => emit('update:modelValue', v))
</script>

<template>
  <div class="flex items-center gap-4">
    <label class="text-sm text-zinc-600 dark:text-zinc-400">
      Date
      <input v-model="date" type="date"
        class="ml-2 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-sm dark:border-zinc-800 dark:bg-zinc-950" />
    </label>
  </div>
</template>
