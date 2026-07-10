<script lang="ts">
  import './layout.css'
  import { onNavigate } from '$app/navigation'
  import type { Snippet } from 'svelte'

  let { children }: { children: Snippet } = $props()

  // Transição de página com a View Transitions API.
  // No-op em navegadores sem suporte ou com prefers-reduced-motion.
  onNavigate((navigation) => {
    if (!document.startViewTransition) return
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve()
        await navigation.complete
      })
    })
  })
</script>

<svelte:head>
  <link rel="icon" href="/favicon.png" />
</svelte:head>

{@render children()}