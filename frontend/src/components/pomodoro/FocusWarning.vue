<script setup lang="ts">
import { usePomodoroStore } from '@/stores/pomodoro'
import Modal from '@/components/ui/Modal.vue'
import Button from '@/components/ui/Button.vue'

const store = usePomodoroStore()
</script>

<template>
  <Modal :open="store.hasViolation" title="Focus Lost" @close="store.dismissViolation()">
    <template v-if="store.state.violation">
      <p class="mb-4 text-sm text-zinc-600 dark:text-zinc-400">
        You switched to <strong class="text-zinc-900 dark:text-zinc-50">{{ store.state.violation.app_name }}</strong>,
        which is not in your allowed apps for this focus session.
      </p>
      <div class="flex gap-3">
        <Button @click="store.dismissViolation()">Return to Focus</Button>
        <Button variant="secondary" @click="store.addViolatingAppToAllowlist()">
          Add to Session
        </Button>
      </div>
    </template>
  </Modal>
</template>
