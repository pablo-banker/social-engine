<script lang="ts">
  import { enhance } from '$app/forms'
  import type { PageProps } from './$types'

  type LoginForm = {
    errors?: {
      general?: string
      email?: string
      password?: string
    }
    values?: {
      email?: string
    }
  }

  let { form }: PageProps = $props()

  let loading = $state(false)

  const typedForm = $derived(form as LoginForm | undefined)
</script>

<div class="auth-card relative w-full max-w-[470px]">
  <div class="auth-card-inner">
    <div class="auth-shine"></div>
    <div class="auth-top-gradient"></div>
    <div class="auth-soft-highlight"></div>

    <div class="relative z-10">
      <div class="mb-9">
        <div class="mb-6 flex items-center gap-3">
          <div class="grid size-9 place-items-center rounded-xl border border-cyan-300/25 bg-cyan-300/10 shadow-lg shadow-cyan-400/10">
            <span class="size-2.5 rounded-full bg-cyan-300 shadow-[0_0_18px_rgba(103,232,249,0.9)]"></span>
          </div>

          <p class="text-xs font-bold uppercase tracking-[0.32em] text-cyan-300">
            Social Engine
          </p>
        </div>

        <h1 class="text-4xl font-bold tracking-tight text-white">
          Bem-vindo de volta.
        </h1>

        <p class="mt-4 max-w-sm text-base leading-7 text-zinc-400">
          Entre para continuar compartilhando ideias, conectando pessoas e construindo sua comunidade.
        </p>
      </div>

      <form
        method="POST"
        action="?/login"
        use:enhance={() => {
          loading = true

          return async ({ update }) => {
            await update()
            loading = false
          }
        }}
        class="space-y-5"
      >
        <div>
          <label for="email" class="mb-2 block text-sm font-semibold text-zinc-300">
            Email
          </label>

          <input
            id="email"
            name="email"
            type="email"
            placeholder="seu@email.com"
            autocomplete="email"
            value={typedForm?.values?.email ?? ''}
            class="auth-input"
          />

          {#if typedForm?.errors?.email}
            <p class="mt-2 text-sm text-red-300">{typedForm.errors.email}</p>
          {/if}
        </div>

        <div>
          <label for="password" class="mb-2 block text-sm font-semibold text-zinc-300">
            Senha
          </label>

          <input
            id="password"
            name="password"
            type="password"
            placeholder="••••••••"
            autocomplete="current-password"
            class="auth-input"
          />

          {#if typedForm?.errors?.password}
            <p class="mt-2 text-sm text-red-300">{typedForm.errors.password}</p>
          {/if}
        </div>

        {#if typedForm?.errors?.general}
          <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
            {typedForm.errors.general}
          </p>
        {/if}

        <button type="submit" disabled={loading} class="auth-button">
          <span class="relative z-10">
            {loading ? 'Entrando...' : 'Entrar'}
          </span>
        </button>
      </form>

      <div class="mt-7 flex items-center justify-between text-sm">
        <a href="/forgot-password" class="text-zinc-400 transition hover:text-white">
          Esqueci minha senha
        </a>

        <a href="/register" class="font-bold text-cyan-300 transition hover:text-cyan-200">
          Criar conta
        </a>
      </div>
    </div>
  </div>
</div>