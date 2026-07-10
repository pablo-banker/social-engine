<script lang="ts">
  import { page } from '$app/state'
</script>

<main class="mx-auto flex max-w-3xl flex-col items-center px-6 py-24 text-center">
  <p class="text-sm font-bold uppercase tracking-[0.3em] text-cyan-300">
    Erro {page.status}
  </p>

  <h1 class="mt-4 text-4xl font-bold tracking-tight text-white">
    {page.error?.message ?? 'Algo deu errado.'}
  </h1>

  <p class="mt-4 max-w-md text-zinc-400">
    A página que você procura não existe ou foi movida. Volte ao feed para continuar explorando.
  </p>

  <a
    href="/"
    class="mt-8 flex h-12 items-center justify-center rounded-2xl bg-white px-6 font-bold text-zinc-950 transition hover:scale-[1.02] active:scale-[0.98]"
  >
    Voltar ao feed
  </a>
</main>
