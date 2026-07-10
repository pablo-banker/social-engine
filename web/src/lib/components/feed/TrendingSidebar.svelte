<script lang="ts">
  import type { TrendingTopic } from '$lib/types'

  type Props = {
    topics?: TrendingTopic[]
  }

  let { topics = [] }: Props = $props()
</script>

<aside class="hidden xl:block">
  <div class="sticky top-28 space-y-4">
    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl" style="--stagger: 1">
      <div class="flex items-center justify-between">
        <p class="text-sm font-bold uppercase tracking-[0.25em] text-cyan-300">
          Trending
        </p>

        <a href="/trending" class="text-xs text-zinc-500 transition hover:text-white">
          ver tudo
        </a>
      </div>

      {#if topics.length}
        <div class="mt-4 space-y-1">
          {#each topics.slice(0, 6) as topic (topic.tag)}
            <a
              href="/explore?tag={topic.tag}"
              class="group flex items-center justify-between rounded-2xl px-3 py-2 transition duration-300 hover:translate-x-1 hover:bg-white/[0.04]"
            >
              <span class="text-sm font-semibold text-zinc-200 transition group-hover:text-cyan-200">#{topic.tag}</span>
              <span class="text-xs text-zinc-500">
                {topic.posts}
                {topic.posts === 1 ? 'post' : 'posts'}
              </span>
            </a>
          {/each}
        </div>
      {:else}
        <p class="mt-4 text-sm text-zinc-500">
          Nada em alta ainda.
        </p>
      {/if}
    </div>

    <div class="rise rounded-3xl border border-white/10 bg-white/[0.035] p-5 backdrop-blur-xl" style="--stagger: 2">
      <p class="font-bold text-white">
        Social Engine MVP
      </p>

      <p class="mt-2 text-sm leading-6 text-zinc-400">
        Próximo passo: conectar o feed com o backend Go e persistir tudo no PostgreSQL.
      </p>
    </div>
  </div>
</aside>