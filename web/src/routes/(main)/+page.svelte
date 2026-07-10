<script lang="ts">
  import { enhance } from '$app/forms'
  import PostCard from '$lib/components/feed/PostCard.svelte'
  import SideNav from '$lib/components/nav/SideNav.svelte'
  import TrendingSidebar from '$lib/components/feed/TrendingSidebar.svelte'
  import type { PageProps } from './$types'

  let { data, form }: PageProps = $props()

  const maxLength = 500

  let loading = $state(false)
  let content = $state('')

  const remaining = $derived(maxLength - content.length)
  const canPublish = $derived(content.trim().length > 0 && remaining >= 0)
</script>

<main class="mx-auto grid max-w-[1800px] grid-cols-1 gap-8 px-6 py-8 lg:grid-cols-[280px_minmax(0,1fr)_360px] lg:px-10">
  <SideNav />

  <section class="min-w-0">
    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl">
      <p class="text-sm font-bold uppercase tracking-[0.3em] text-cyan-300">
        Public feed
      </p>

      {#if data.session}
        <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
          Bem-vindo de volta, {data.session.user.name.split(' ')[0]}.
        </h1>

        <p class="mt-3 max-w-2xl text-zinc-400">
          Compartilhe o que você está construindo e acompanhe o que a comunidade está publicando.
        </p>
      {:else}
        <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
          Veja o que a comunidade está publicando.
        </h1>

        <p class="mt-3 max-w-2xl text-zinc-400">
          O feed é público. Para postar, comentar, curtir ou seguir pessoas, você precisa estar logado.
        </p>
      {/if}
    </div>

    <div class="rise mt-6 rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl" style="--stagger: 1">
      {#if data.session}
        <form
          method="POST"
          action="?/createPost"
          use:enhance={() => {
            loading = true

            return async ({ result, update }) => {
              await update()
              loading = false

              if (result.type === 'success') {
                content = ''
              }
            }
          }}
          class="space-y-4"
        >
          <textarea
            name="content"
            bind:value={content}
            rows="4"
            maxlength={maxLength}
            placeholder="O que você está construindo hoje?"
            class="w-full resize-none rounded-2xl border border-white/10 bg-white/[0.035] p-4 text-white outline-none transition placeholder:text-zinc-600 hover:border-cyan-300/30 focus:border-cyan-300/70 focus:ring-4 focus:ring-cyan-300/10"
          ></textarea>

          {#if form?.error}
            <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
              {form.error}
            </p>
          {/if}

          <div class="flex items-center justify-between">
            <span class="text-sm text-zinc-500">
              {remaining} caracteres restantes
            </span>

            <button
              type="submit"
              disabled={loading || !canPublish}
              class="btn-shine h-12 rounded-2xl bg-white px-5 font-bold text-zinc-950 transition hover:scale-[1.01] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.985] disabled:cursor-not-allowed disabled:opacity-70"
            >
              <span class="relative z-10">{loading ? 'Publicando...' : 'Publicar'}</span>
            </button>
          </div>
        </form>
      {:else}
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <p class="font-bold text-white">
              Quer participar da conversa?
            </p>

            <p class="mt-1 text-sm text-zinc-400">
              Crie uma conta para publicar, comentar, curtir e seguir pessoas.
            </p>
          </div>

          <a
            href="/login"
            class="btn-shine flex h-12 items-center justify-center rounded-2xl bg-white px-5 font-bold text-zinc-950 transition hover:scale-[1.01] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.985]"
          >
            <span class="relative z-10">Fazer login</span>
          </a>
        </div>
      {/if}
    </div>

    <div class="mt-6 space-y-4">
      {#each data.posts as post, i (post.id)}
        <div class="rise" style="--stagger: {i + 2}">
          <PostCard {post} loggedIn={Boolean(data.session)} />
        </div>
      {/each}
    </div>
  </section>

  <TrendingSidebar topics={data.topics} />
</main>