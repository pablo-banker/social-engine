<script lang="ts">
  import { untrack } from 'svelte'
  import { enhance } from '$app/forms'
  import { AVATARS, BANNERS, avatarBackground, bannerGradient } from '$lib/appearance'
  import type { PageProps } from './$types'

  let { data, form }: PageProps = $props()

  const maxBio = 160

  // Valores iniciais: da sessão (load) ou repopulados após erro de validação.
  const initial = untrack(() => form?.values ?? data.settings)

  let name = $state(initial.name)
  let bio = $state(initial.bio)
  let avatarId = $state(initial.avatarId)
  let bannerId = $state(initial.bannerId)
  let loading = $state(false)

  const bioRemaining = $derived(maxBio - bio.length)
</script>

<svelte:head>
  <title>Editar perfil · Social Engine</title>
</svelte:head>

<main class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
  <a
    href="/{data.session?.user.username}"
    class="inline-flex items-center gap-2 text-sm font-medium text-zinc-400 transition hover:text-white"
  >
    <span aria-hidden="true">←</span>
    Voltar ao perfil
  </a>

  <h1 class="mt-4 text-3xl font-bold tracking-tight text-white">
    Editar perfil
  </h1>

  <!-- Preview ao vivo -->
  <div class="rise mt-6 overflow-hidden rounded-3xl border border-white/10 bg-white/[0.02]">
    <div class="banner-live h-28" style="background: {bannerGradient(bannerId)}"></div>

    <div class="px-6 pb-5">
      <div
        class="-mt-9 size-16 rounded-full border-4 border-black shadow-lg shadow-black/40"
        style="background: {avatarBackground(avatarId)}"
      ></div>

      <p class="mt-3 text-lg font-bold text-white">
        {name || 'Seu nome'}
      </p>

      <p class="text-sm text-zinc-500">
        @{data.session?.user.username}
      </p>
    </div>
  </div>

  <form
    method="POST"
    action="?/save"
    use:enhance={() => {
      loading = true

      return async ({ update }) => {
        await update()
        loading = false
      }
    }}
    class="mt-6 space-y-6"
  >
    <input type="hidden" name="avatarId" value={avatarId} />
    <input type="hidden" name="bannerId" value={bannerId} />

    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl" style="--stagger: 1">
      <div>
        <label for="name" class="mb-2 block text-sm font-semibold text-zinc-300">
          Nome
        </label>

        <input
          id="name"
          name="name"
          bind:value={name}
          type="text"
          maxlength="50"
          placeholder="Seu nome"
          class="h-12 w-full rounded-2xl border border-white/10 bg-white/[0.035] px-4 text-white outline-none transition placeholder:text-zinc-600 hover:border-cyan-300/30 focus:border-cyan-300/70 focus:ring-4 focus:ring-cyan-300/10"
        />
      </div>

      <div class="mt-5">
        <label for="bio" class="mb-2 block text-sm font-semibold text-zinc-300">
          Bio
        </label>

        <textarea
          id="bio"
          name="bio"
          bind:value={bio}
          rows="3"
          maxlength={maxBio}
          placeholder="Conte um pouco sobre você..."
          class="w-full resize-none rounded-2xl border border-white/10 bg-white/[0.035] p-4 text-white outline-none transition placeholder:text-zinc-600 hover:border-cyan-300/30 focus:border-cyan-300/70 focus:ring-4 focus:ring-cyan-300/10"
        ></textarea>

        <p class="mt-2 text-right text-xs text-zinc-500">
          {bioRemaining} caracteres restantes
        </p>
      </div>
    </div>

    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl" style="--stagger: 2">
      <p class="text-sm font-semibold text-zinc-300">
        Avatar
      </p>

      <div class="mt-4 grid grid-cols-5 gap-3 sm:grid-cols-10">
        {#each AVATARS as avatar (avatar.id)}
          <button
            type="button"
            onclick={() => (avatarId = avatar.id)}
            aria-label="Escolher avatar"
            aria-pressed={avatarId === avatar.id}
            style="background: {avatarBackground(avatar.id)}"
            class={[
              'aspect-square rounded-full transition duration-300 hover:scale-110 active:scale-95',
              avatarId === avatar.id
                ? 'ring-2 ring-cyan-300 ring-offset-2 ring-offset-zinc-950'
                : 'ring-1 ring-white/10'
            ]}
          ></button>
        {/each}
      </div>
    </div>

    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl" style="--stagger: 3">
      <p class="text-sm font-semibold text-zinc-300">
        Banner
      </p>

      <div class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-5">
        {#each BANNERS as banner (banner.id)}
          <button
            type="button"
            onclick={() => (bannerId = banner.id)}
            aria-label="Escolher banner"
            aria-pressed={bannerId === banner.id}
            style="background: {banner.gradient}"
            class={[
              'h-14 rounded-2xl transition duration-300 hover:scale-[1.04] active:scale-[0.97]',
              bannerId === banner.id
                ? 'ring-2 ring-cyan-300 ring-offset-2 ring-offset-zinc-950'
                : 'ring-1 ring-white/10'
            ]}
          ></button>
        {/each}
      </div>
    </div>

    {#if form?.error}
      <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
        {form.error}
      </p>
    {/if}

    <div class="flex items-center justify-end gap-3">
      <a
        href="/{data.session?.user.username}"
        class="flex h-12 items-center rounded-2xl border border-white/10 px-5 text-sm font-bold text-zinc-300 transition hover:border-white/20 hover:bg-white/[0.04] hover:text-white"
      >
        Cancelar
      </a>

      <button
        type="submit"
        disabled={loading}
        class="btn-shine h-12 rounded-2xl bg-white px-6 font-bold text-zinc-950 transition hover:scale-[1.01] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.985] disabled:cursor-not-allowed disabled:opacity-70"
      >
        <span class="relative z-10">{loading ? 'Salvando...' : 'Salvar'}</span>
      </button>
    </div>
  </form>
</main>
