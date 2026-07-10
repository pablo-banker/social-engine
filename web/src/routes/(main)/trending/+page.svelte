<script lang="ts">
  import PostCard from '$lib/components/feed/PostCard.svelte'
  import SideNav from '$lib/components/nav/SideNav.svelte'
  import type { PageProps } from './$types'

  let { data }: PageProps = $props()

  const loggedIn = $derived(Boolean(data.session))
</script>

<svelte:head>
  <title>Trending · Social Engine</title>
</svelte:head>

<main class="mx-auto grid max-w-[1800px] grid-cols-1 gap-8 px-6 py-8 lg:grid-cols-[280px_minmax(0,1fr)_360px] lg:px-10">
  <SideNav />

  <section class="min-w-0">
    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl">
      <p class="text-sm font-bold uppercase tracking-[0.3em] text-cyan-300">
        Trending
      </p>

      <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
        O que está em alta agora.
      </h1>

      <p class="mt-3 max-w-2xl text-zinc-400">
        Tópicos e posts mais populares da comunidade, calculados em tempo real.
      </p>
    </div>

    <div class="rise mt-6 rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl" style="--stagger: 1">
      <p class="text-sm font-bold uppercase tracking-[0.25em] text-cyan-300">
        Tópicos em alta
      </p>

      {#if data.topics.length}
        <ol class="mt-5 space-y-1">
          {#each data.topics as topic, i (topic.tag)}
            <a
              href="/explore?tag={topic.tag}"
              class="rise group flex items-center gap-4 rounded-2xl px-3 py-3 transition duration-300 hover:translate-x-1 hover:bg-white/[0.04]"
              style="--stagger: {i + 2}"
            >
              <span class="w-6 text-lg font-bold tabular-nums text-zinc-600 transition duration-300 group-hover:text-cyan-300">
                {i + 1}
              </span>

              <span class="flex-1 font-semibold text-zinc-100 transition group-hover:text-cyan-200">
                #{topic.tag}
              </span>

              <span class="text-sm text-zinc-500">
                {topic.posts}
                {topic.posts === 1 ? 'post' : 'posts'}
              </span>
            </a>
          {/each}
        </ol>
      {:else}
        <p class="mt-4 text-sm text-zinc-400">
          Nenhum tópico em alta ainda. Publique com #tags para começar.
        </p>
      {/if}
    </div>

    <div class="mt-6">
      <h2 class="mb-4 px-1 text-sm font-bold uppercase tracking-[0.25em] text-zinc-500">
        Posts em alta
      </h2>

      <div class="space-y-4">
        {#each data.posts as post, i (post.id)}
          <div class="rise" style="--stagger: {i + 3}">
            <PostCard {post} {loggedIn} />
          </div>
        {:else}
          <div class="rounded-3xl border border-white/10 bg-white/[0.035] p-10 text-center backdrop-blur-xl">
            <p class="text-sm text-zinc-400">
              Nenhum post por aqui ainda.
            </p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <aside class="hidden xl:block">
    <div class="rise sticky top-28 rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl" style="--stagger: 1">
      <p class="font-bold text-white">
        Como o trending funciona
      </p>

      <p class="mt-2 text-sm leading-6 text-zinc-400">
        Os tópicos vêm das <span class="text-zinc-200">#tags</span> usadas nos posts e os "posts em alta"
        são ordenados por curtidas. Tudo em memória por enquanto — depois virá do Go + PostgreSQL.
      </p>
    </div>
  </aside>
</main>
