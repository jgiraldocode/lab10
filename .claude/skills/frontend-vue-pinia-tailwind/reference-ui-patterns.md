# Patrones UI (Vue + Tailwind)

## Card estándar

```vue
<article
  class="rounded-xl border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-950"
>
  <slot />
</article>
```

## Botón primario (discreto)

```vue
<button
  type="button"
  class="inline-flex items-center justify-center rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-400 focus-visible:ring-offset-2 dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200"
>
  <slot />
</button>
```

## Store Pinia mínima con fetch

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchUser } from '@/lib/api'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      user.value = await fetchUser()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error'
    } finally {
      loading.value = false
    }
  }

  return { user, loading, error, load }
})
```

Estos fragmentos son guía; adaptar nombres y rutas al proyecto.
