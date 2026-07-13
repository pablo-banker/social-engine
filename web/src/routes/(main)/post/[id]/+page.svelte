<script lang="ts">
  import { enhance } from '$app/forms'
  import PostCard from '$lib/components/feed/PostCard.svelte'
  import { avatarBackground } from '$lib/appearance'
  import { relativeTime } from '$lib/format'
  import type { PageProps } from './$types'

  let { data, form }: PageProps = $props()

  const maxLength = 300

  let loading = $state(false)
  let content = $state('')

  const remaining = $derived(maxLength - content.length)
  const canSend = $derived(content.trim().length > 0 && remaining >= 0)
</script>

<svelte:head>
  <title>Post de {data.post.author.name} · Social Engine</title>
</svelte:head>

<main class="mx-auto max-w-5xl px-6 py-8 lg:px-8">
  <a href="/" class="rise group inline-flex items-center gap-2 text-sm font-medium text-zinc-400 transition hover:text-white">
    <span aria-hidden="true" class="transition-transform duration-300 group-hover:-translate-x-1">←</span>
    Voltar ao feed
  </a>

  <div class="rise mt-4" style="--stagger: 1">
    <PostCard post={data.post} loggedIn={Boolean(data.session)} linkToPost={false} />
  </div>

  <div class="rise mt-6 rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl" style="--stagger: 2">
    {#if data.session}
      <form
        method="POST"
        action="?/addComment"
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
          rows="3"
          maxlength={maxLength}
          placeholder="Escreva um comentário..."
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
            disabled={loading || !canSend}
            class="btn-shine h-11 rounded-2xl bg-white px-5 font-bold text-zinc-950 transition hover:scale-[1.01] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.985] disabled:cursor-not-allowed disabled:opacity-70"
          >
            <span class="relative z-10">{loading ? 'Enviando...' : 'Comentar'}</span>
          </button>
        </div>
      </form>
    {:else}
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="font-bold text-white">
            Entre para comentar
          </p>

          <p class="mt-1 text-sm text-zinc-400">
            Você pode ler os comentários, mas precisa estar logado para responder.
          </p>
        </div>

        <a
          href="/login"
          class="btn-shine flex h-11 items-center justify-center rounded-2xl bg-white px-5 font-bold text-zinc-950 transition hover:scale-[1.01] hover:shadow-[0_15px_40px_rgba(34,211,238,0.12)] active:scale-[0.985]"
        >
          <span class="relative z-10">Fazer login</span>
        </a>
      </div>
    {/if}
  </div>

  <div class="mt-6 space-y-3">
    {#each data.comments as comment, i (comment.id)}
      <article
        class="rise rounded-2xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl transition duration-300 hover:border-white/15 hover:bg-white/[0.045]"
        style="--stagger: {i + 3}"
      >
        <div class="flex items-center gap-3">
          <a
            href="/{comment.author.username}"
            style="background: {avatarBackground(comment.author.avatarId)}"
            class="size-9 rounded-full transition hover:scale-105"
            aria-label="Perfil de {comment.author.name}"
          ></a>

          <div class="min-w-0">
            <a href="/{comment.author.username}" class="text-sm font-bold text-white hover:underline">
              {comment.author.name}
            </a>

            <p class="text-xs text-zinc-500">
              @{comment.author.username} · {relativeTime(comment.createdAt)}
            </p>
          </div>
        </div>

        <p class="mt-3 text-sm leading-6 whitespace-pre-wrap text-zinc-300">
          {comment.content}
        </p>
      </article>
    {:else}
      <div class="rounded-2xl border border-white/10 bg-white/[0.035] p-8 text-center backdrop-blur-xl">
        <p class="text-sm text-zinc-400">
          Nenhum comentário ainda. Seja o primeiro a responder.
        </p>
      </div>
    {/each}
  </div>
</main>
