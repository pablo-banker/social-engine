<script lang="ts">
  import { enhance } from '$app/forms'
  import { avatarBackground } from '$lib/appearance'
  import type { Post } from '$lib/types'

  type Props = {
    post: Post
    loggedIn?: boolean
    /** Torna o card inteiro clicável (abre os comentários). */
    linkToPost?: boolean
  }

  let { post, loggedIn = false, linkToPost = true }: Props = $props()

  // Dispara a animação de "estouro" do coração ao curtir.
  let burst = $state(false)
</script>

{#snippet heart(filled: boolean)}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    width="22"
    height="22"
    fill={filled ? 'currentColor' : 'none'}
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path
      d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"
    />
  </svg>
{/snippet}

{#snippet heartButton()}
  <span class="relative grid place-items-center">
    {#if burst}
      <span class="heart-ring" aria-hidden="true"></span>
    {/if}

    <span
      class={['grid place-items-center', burst && 'heart-pop']}
      onanimationend={() => (burst = false)}
    >
      {@render heart(post.likedByMe)}
    </span>
  </span>

  {#key post.likes}
    <span class="count-swap text-xs font-semibold tabular-nums">{post.likes}</span>
  {/key}
{/snippet}

<article
  class="group relative rounded-3xl border border-white/10 bg-white/[0.035] p-6 backdrop-blur-xl transition duration-300 hover:-translate-y-0.5 hover:border-cyan-300/20 hover:bg-white/[0.045] hover:shadow-[0_18px_50px_rgba(34,211,238,0.07)]"
>
  {#if linkToPost}
    <!-- Link "esticado": o card inteiro abre os comentários. Os elementos
         interativos abaixo ficam acima dele (z-10) e continuam clicáveis. -->
    <a
      href="/post/{post.id}"
      class="absolute inset-0 rounded-3xl"
      aria-label="Ver post de {post.author.name} e comentários"
    ></a>
  {/if}

  <!-- Curtir: canto superior direito, coração + contagem embaixo -->
  <div class="absolute right-4 top-4 z-10">
    {#if loggedIn}
      <form
        method="POST"
        action="/post/{post.id}?/toggleLike"
        use:enhance={() => {
          const willLike = !post.likedByMe

          return async ({ update }) => {
            await update()

            if (willLike) burst = true
          }
        }}
      >
        <button
          type="submit"
          aria-pressed={post.likedByMe}
          aria-label={post.likedByMe ? 'Descurtir' : 'Curtir'}
          class={[
            'flex flex-col items-center gap-1 transition active:scale-90',
            post.likedByMe ? 'text-rose-500' : 'text-zinc-500 hover:text-rose-400'
          ]}
        >
          {@render heartButton()}
        </button>
      </form>
    {:else}
      <a
        href="/login"
        aria-label="Entrar para curtir"
        class="flex flex-col items-center gap-1 text-zinc-500 transition hover:text-rose-400"
      >
        {@render heartButton()}
      </a>
    {/if}
  </div>

  <div class="flex items-center gap-3 pr-16">
    <a
      href="/{post.author.username}"
      style="background: {avatarBackground(post.author.avatarId)}"
      class="relative z-10 size-11 shrink-0 rounded-full ring-0 ring-cyan-300/40 transition duration-300 hover:scale-105 hover:ring-2"
      aria-label="Perfil de {post.author.name}"
    ></a>

    <div class="min-w-0">
      <a href="/{post.author.username}" class="relative z-10 font-bold text-white hover:underline">
        {post.author.name}
      </a>

      <p class="text-sm text-zinc-500">
        <a href="/{post.author.username}" class="relative z-10 transition hover:text-zinc-300">
          @{post.author.username}
        </a>
        · {post.createdAt}
      </p>
    </div>
  </div>

  <p class="mt-5 leading-7 whitespace-pre-wrap text-zinc-300">
    {post.content}
  </p>

  <div class="mt-5 flex items-center gap-2 text-sm text-zinc-500 transition duration-300 group-hover:text-zinc-400">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      width="16"
      height="16"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
      class="transition duration-300 group-hover:scale-110 group-hover:text-cyan-300"
    >
      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
    </svg>

    <span>{post.comments} comentários</span>
  </div>
</article>