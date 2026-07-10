<script lang="ts">
  import { page } from '$app/state'
  import { avatarBackground } from '$lib/appearance'
  import AppBackground from '$lib/components/AppBackground.svelte'
  import type { Snippet } from 'svelte'
  import type { LayoutProps } from './$types'

  let { children, data }: LayoutProps & { children: Snippet } = $props()

  const navItems = [
    { href: '/', label: 'Feed', isActive: (path: string) => path === '/' },
    { href: '/explore', label: 'Explore', isActive: (path: string) => path.startsWith('/explore') },
    { href: '/trending', label: 'Trending', isActive: (path: string) => path.startsWith('/trending') }
  ]
</script>

<section class="min-h-screen bg-black text-white">
  <AppBackground />

  <header
    class="sticky top-0 z-50 border-b border-white/10 bg-black/70 backdrop-blur-2xl"
    style="view-transition-name: header"
  >
    <div class="mx-auto flex max-w-[1800px] items-center justify-between px-6 py-4 lg:px-10">
      <a href="/" class="group flex items-center gap-3">
        <img
          src="/favicon.png"
          alt="Social Engine"
          class="size-8 transition duration-300 group-hover:scale-110 group-hover:-rotate-6"
        />

        <div>
          <p class="font-bold text-white">
            Social Engine
          </p>

          <p class="text-xs text-zinc-500 transition group-hover:text-cyan-300/80">
            Open social platform
          </p>
        </div>
      </a>

      <nav class="hidden items-center gap-6 text-sm font-medium md:flex">
        {#each navItems as item (item.href)}
          {@const active = item.isActive(page.url.pathname)}

          <a
            href={item.href}
            aria-current={active ? 'page' : undefined}
            class={[
              'relative py-1 transition after:absolute after:inset-x-0 after:-bottom-0.5 after:h-px after:origin-left after:bg-gradient-to-r after:from-cyan-300 after:to-violet-400 after:transition-transform after:duration-300 hover:after:scale-x-100',
              active
                ? 'font-semibold text-white after:scale-x-100'
                : 'text-zinc-400 after:scale-x-0 hover:text-white'
            ]}
          >
            {item.label}
          </a>
        {/each}
      </nav>

      <div class="flex items-center gap-3">
        {#if data.session}
          <a
            href="/{data.session.user.username}"
            aria-label="Meu perfil"
            class="hidden items-center gap-2.5 rounded-2xl border border-white/10 py-1.5 pl-1.5 pr-3.5 transition duration-300 hover:-translate-y-0.5 hover:border-cyan-300/25 hover:bg-white/[0.04] hover:shadow-[0_10px_30px_rgba(34,211,238,0.08)] sm:flex"
          >
            <span class="size-7 rounded-full" style="background: {avatarBackground(data.session.user.avatarId)}"></span>
            <span class="text-sm font-bold text-white">{data.session.user.name}</span>
          </a>

          <form method="POST" action="/logout">
            <button class="rounded-2xl border border-white/10 px-4 py-2.5 text-sm font-bold text-zinc-300 transition duration-300 hover:-translate-y-0.5 hover:border-red-400/30 hover:bg-red-400/10 hover:text-red-200">
              Sair
            </button>
          </form>
        {:else}
          <a
            href="/login"
            class="rounded-2xl border border-white/10 px-4 py-2.5 text-sm font-bold text-zinc-300 transition duration-300 hover:-translate-y-0.5 hover:border-white/20 hover:bg-white/[0.04] hover:text-white"
          >
            Login
          </a>

          <a
            href="/register"
            class="btn-shine rounded-2xl bg-white px-4 py-2.5 text-sm font-bold text-zinc-950 transition hover:scale-[1.02] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.98]"
          >
            <span class="relative z-10">Criar conta</span>
          </a>
        {/if}
      </div>
    </div>
  </header>

  <div class="relative z-10">
    {@render children()}
  </div>
</section>