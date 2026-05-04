<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePomodoroStore } from '@/stores/pomodoro'
import { useTrackingStore } from '@/stores/tracking'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const pomodoroStore = usePomodoroStore()
const trackingStore = useTrackingStore()

const selectedApp = ref('')
const saveAsDefault = ref(true)

onMounted(async () => {
  const sid = pomodoroStore.state.active ? pomodoroStore.state.session_id : undefined
  await pomodoroStore.fetchAllowlist(sid)
  await trackingStore.fetchApps()
})

async function handleAdd() {
  const app = trackingStore.apps.find(a => a.bundle_id === selectedApp.value)
  if (!app) return

  const sid = pomodoroStore.state.active ? pomodoroStore.state.session_id : undefined
  await pomodoroStore.addApp(app.bundle_id, app.app_name, saveAsDefault.value, sid)
  selectedApp.value = ''
}
</script>

<template>
  <Card>
    <h3 class="mb-4 text-sm font-semibold text-zinc-900 dark:text-zinc-50">Allowed Apps</h3>

    <!-- Add form -->
    <div class="mb-4 flex items-end gap-2">
      <div class="flex-1">
        <label class="mb-1 block text-xs text-zinc-500">App</label>
        <select v-model="selectedApp"
          class="w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-800 dark:bg-zinc-950">
          <option value="">Select an app...</option>
          <option v-for="app in trackingStore.apps" :key="app.bundle_id" :value="app.bundle_id">
            {{ app.app_name || app.bundle_id }}
          </option>
        </select>
      </div>
      <label class="flex items-center gap-1.5 text-xs text-zinc-600 dark:text-zinc-400">
        <input v-model="saveAsDefault" type="checkbox" class="rounded" />
        Default
      </label>
      <Button size="sm" :disabled="!selectedApp" @click="handleAdd">Add</Button>
    </div>

    <!-- List -->
    <ul class="space-y-2">
      <li v-for="entry in pomodoroStore.allowlist" :key="entry.id"
        class="flex items-center justify-between rounded-lg px-3 py-2 text-sm hover:bg-zinc-50 dark:hover:bg-zinc-900">
        <div>
          <span class="text-zinc-900 dark:text-zinc-50">{{ entry.app_name }}</span>
          <span v-if="entry.is_default" class="ml-2 text-xs text-zinc-500">(default)</span>
        </div>
        <button @click="pomodoroStore.removeApp(entry.id)"
          class="text-xs text-zinc-400 hover:text-red-500 transition-colors">
          Remove
        </button>
      </li>
    </ul>

    <p v-if="pomodoroStore.allowlist.length === 0" class="text-sm text-zinc-500">
      No apps in allowlist. Add apps that are permitted during focus sessions.
    </p>
  </Card>
</template>
