<script lang="ts">
  import PostCard from '$lib/components/feed/PostCard.svelte'
  import { avatarBackground, bannerGradient } from '$lib/appearance'
  import type { PageProps } from './$types'

  let { data }: PageProps = $props()

  const profile = $derived(data.profile)
  const isOwnProfile = $derived(data.session?.user.username === profile.username)

  // Follow é visual por enquanto (não persiste) — vira api.users.follow depois.
  let following = $state(false)
</script>

<svelte:head>
  <title>{profile.name} (@{profile.username}) · Social Engine</title>
</svelte:head>

<main class="mx-auto max-w-5xl px-6 py-8 lg:px-8">
  <div class="rise overflow-hidden rounded-3xl border border-white/10 bg-white/[0.035] backdrop-blur-xl">
    <div class="banner-live h-36" style="background: {bannerGradient(profile.bannerId)}"></div>

    <div class="px-6 pb-6">
      <div class="-mt-12 flex items-end justify-between">
        <div
          class="size-24 rounded-full border-4 border-black shadow-xl shadow-cyan-500/20 transition duration-300 hover:scale-105 hover:shadow-cyan-500/40"
          style="background: {avatarBackground(profile.avatarId)}"
        ></div>

        {#if isOwnProfile}
          <a
            href="/settings/profile"
            class="rounded-2xl border border-white/10 px-5 py-2.5 text-sm font-bold text-zinc-300 transition hover:border-white/20 hover:bg-white/[0.04] hover:text-white"
          >
            Editar perfil
          </a>
        {:else if data.session}
          <button
            onclick={() => (following = !following)}
            class={[
              'rounded-2xl px-6 py-2.5 text-sm font-bold transition active:scale-[0.98]',
              following
                ? 'border border-white/15 text-zinc-300 hover:border-red-400/30 hover:bg-red-400/10 hover:text-red-200'
                : 'bg-white text-zinc-950 hover:scale-[1.02]'
            ]}
          >
            {following ? 'Seguindo' : 'Seguir'}
          </button>
        {:else}
          <a
            href="/login"
            class="rounded-2xl bg-white px-6 py-2.5 text-sm font-bold text-zinc-950 transition hover:scale-[1.02] active:scale-[0.98]"
          >
            Seguir
          </a>
        {/if}
      </div>

      <h1 class="mt-4 text-2xl font-bold tracking-tight text-white">
        {profile.name}
      </h1>

      <p class="text-sm text-zinc-500">
        @{profile.username}
      </p>

      {#if profile.bio}
        <p class="mt-4 leading-7 text-zinc-300">
          {profile.bio}
        </p>
      {/if}

      <p class="mt-3 text-sm text-zinc-500">
        Entrou em {profile.joinedAt}
      </p>

      <div class="mt-5 flex flex-wrap gap-6 text-sm text-zinc-400">
        <span><span class="font-bold text-white">{profile.stats.posts}</span> posts</span>
        <span><span class="font-bold text-white">{profile.stats.followers}</span> seguidores</span>
        <span><span class="font-bold text-white">{profile.stats.following}</span> seguindo</span>
      </div>
    </div>
  </div>

  <div class="rise mt-6 border-b border-white/10" style="--stagger: 1">
    <span class="inline-block border-b-2 border-cyan-300 pb-3 text-sm font-bold text-white">
      Posts
    </span>
  </div>

  <div class="mt-6 space-y-4">
    {#each data.posts as post, i (post.id)}
      <div class="rise" style="--stagger: {i + 2}">
        <PostCard {post} loggedIn={Boolean(data.session)} />
      </div>
    {:else}
      <div class="rounded-3xl border border-white/10 bg-white/[0.035] p-10 text-center backdrop-blur-xl">
        <p class="font-bold text-white">
          Nada por aqui ainda.
        </p>

        <p class="mt-2 text-sm text-zinc-400">
          @{profile.username} ainda não publicou nada.
        </p>
      </div>
    {/each}
  </div>
</main>
