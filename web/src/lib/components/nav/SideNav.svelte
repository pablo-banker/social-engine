<script lang="ts">
  import { page } from '$app/state'

  const items = [
    { href: '/', label: 'Feed', isActive: (path: string) => path === '/' },
    { href: '/explore', label: 'Explore', isActive: (path: string) => path.startsWith('/explore') },
    { href: '/trending', label: 'Trending', isActive: (path: string) => path.startsWith('/trending') }
  ]
</script>

<aside class="hidden lg:block">
  <div class="rise sticky top-28 rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl">
    <p class="text-sm font-bold uppercase tracking-[0.25em] text-cyan-300">
      Menu
    </p>

    <nav class="mt-6 space-y-2">
      {#each items as item (item.href)}
        {@const active = item.isActive(page.url.pathname)}

        <a
          href={item.href}
          aria-current={active ? 'page' : undefined}
          class={[
            'relative block rounded-2xl px-4 py-3 text-sm transition duration-300',
            active
              ? 'bg-cyan-300/10 font-bold text-cyan-200'
              : 'font-medium text-zinc-400 hover:translate-x-1 hover:bg-white/[0.04] hover:text-white'
          ]}
        >
          <span
            class={[
              'absolute left-0 top-1/2 h-5 w-0.5 -translate-y-1/2 rounded-full bg-cyan-300 shadow-[0_0_12px_rgba(103,232,249,0.9)] transition-all duration-300',
              active ? 'opacity-100' : 'opacity-0'
            ]}
            aria-hidden="true"
          ></span>

          {item.label}
        </a>
      {/each}
    </nav>
  </div>
</aside>