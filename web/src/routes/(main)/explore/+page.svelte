<script lang="ts">
  import PostCard from '$lib/components/feed/PostCard.svelte'
  import SideNav from '$lib/components/nav/SideNav.svelte'
  import TrendingSidebar from '$lib/components/feed/TrendingSidebar.svelte'
  import { avatarBackground } from '$lib/appearance'
  import type { PageProps } from './$types'

  let { data }: PageProps = $props()

  const loggedIn = $derived(Boolean(data.session))
</script>

<svelte:head>
  <title>{data.tag ? `#${data.tag}` : 'Explore'} · Social Engine</title>
</svelte:head>

<main class="mx-auto grid max-w-[1800px] grid-cols-1 gap-8 px-6 py-8 lg:grid-cols-[280px_minmax(0,1fr)_360px] lg:px-10">
  <SideNav />

  <section class="min-w-0">
    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl">
      <p class="text-sm font-bold uppercase tracking-[0.3em] text-cyan-300">
        Explore
      </p>

      {#if data.tag}
        <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
          #{data.tag}
        </h1>

        <p class="mt-3 text-zinc-400">
          {data.posts.length}
          {data.posts.length === 1 ? 'post' : 'posts'} com essa tag.
          <a href="/explore" class="font-semibold text-cyan-300 transition hover:text-cyan-200">
            Limpar filtro
          </a>
        </p>
      {:else}
        <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
          Descubra pessoas e ideias.
        </h1>

        <p class="mt-3 max-w-2xl text-zinc-400">
          Encontre quem seguir e veja os posts em alta na comunidade.
        </p>
      {/if}
    </div>

    {#if !data.tag && data.users.length}
      <div class="rise mt-6 rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl" style="--stagger: 1">
        <p class="text-sm font-bold uppercase tracking-[0.25em] text-cyan-300">
          Quem seguir
        </p>

        <div class="mt-5 grid gap-3 sm:grid-cols-2">
          {#each data.users as user, i (user.id)}
            <a
              href="/{user.username}"
              class="rise group flex items-start gap-3 rounded-2xl border border-white/10 bg-white/[0.02] p-4 transition duration-300 hover:-translate-y-1 hover:border-cyan-300/20 hover:bg-white/[0.04] hover:shadow-[0_14px_40px_rgba(34,211,238,0.08)]"
              style="--stagger: {i + 2}"
            >
              <div
                class="size-11 shrink-0 rounded-full transition duration-300 group-hover:scale-105"
                style="background: {avatarBackground(user.avatarId)}"
              ></div>

              <div class="min-w-0">
                <p class="truncate font-bold text-white">
                  {user.name}
                </p>

                <p class="truncate text-sm text-zinc-500">
                  @{user.username}
                </p>

                <p class="mt-1 line-clamp-2 text-sm text-zinc-400">
                  {user.bio}
                </p>
              </div>
            </a>
          {/each}
        </div>
      </div>
    {/if}

    <div class="mt-6">
      {#if !data.tag}
        <h2 class="mb-4 px-1 text-sm font-bold uppercase tracking-[0.25em] text-zinc-500">
          Em alta
        </h2>
      {/if}

      <div class="space-y-4">
        {#each data.posts as post, i (post.id)}
          <div class="rise" style="--stagger: {i + 3}">
            <PostCard {post} {loggedIn} />
          </div>
        {:else}
          <div class="rounded-3xl border border-white/10 bg-white/[0.035] p-10 text-center backdrop-blur-xl">
            <p class="font-bold text-white">
              Nada encontrado.
            </p>

            <p class="mt-2 text-sm text-zinc-400">
              {#if data.tag}
                Nenhum post com <span class="text-zinc-200">#{data.tag}</span> ainda.
              {:else}
                Nenhum post por aqui ainda.
              {/if}
            </p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <TrendingSidebar topics={data.topics} />
</main>
