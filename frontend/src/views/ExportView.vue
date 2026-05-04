<script setup lang="ts">
import { ref } from 'vue'
import { exportCSV, exportJSON } from '@/lib/api'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

const date = ref(new Date().toISOString().slice(0, 10))

async function download(format: 'csv' | 'json') {
  const blob = format === 'csv' ? await exportCSV(date.value) : await exportJSON(date.value)
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `focustrack-${date.value}.${format}`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div class="space-y-6">
    <h2 class="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">Export Data</h2>

    <Card>
      <div class="space-y-4">
        <label class="block text-sm text-zinc-600 dark:text-zinc-400">
          Date
          <input v-model="date" type="date"
            class="ml-2 rounded-lg border border-zinc-200 bg-white px-3 py-1.5 text-sm dark:border-zinc-800 dark:bg-zinc-950" />
        </label>

        <div class="flex gap-3">
          <Button @click="download('csv')">Download CSV</Button>
          <Button variant="secondary" @click="download('json')">Download JSON</Button>
        </div>
      </div>
    </Card>
  </div>
</template>
