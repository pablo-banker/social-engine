<script lang="ts">
  import { enhance } from '$app/forms'
  import type { PageProps } from './$types'

  type ForgotPasswordForm = {
    sent?: boolean
    errors?: {
      general?: string
      email?: string
    }
    values?: {
      email?: string
    }
  }

  let { form }: PageProps = $props()

  let loading = $state(false)

  const typedForm = $derived(form as ForgotPasswordForm | undefined)
</script>

<div class="auth-card relative w-full max-w-[470px]">
  <div class="auth-card-inner">
    <div class="auth-shine"></div>
    <div class="auth-top-gradient"></div>
    <div class="auth-soft-highlight"></div>

    <div class="relative z-10">
      <div class="mb-8">
        <div class="mb-6 flex items-center gap-3">
          <div class="grid size-9 place-items-center rounded-xl border border-cyan-300/25 bg-cyan-300/10 shadow-lg shadow-cyan-400/10">
            <span class="size-2.5 rounded-full bg-cyan-300 shadow-[0_0_18px_rgba(103,232,249,0.9)]"></span>
          </div>

          <p class="text-xs font-bold uppercase tracking-[0.32em] text-cyan-300">
            Social Engine
          </p>
        </div>

        <h1 class="text-4xl font-bold tracking-tight text-white">
          Recupere sua senha.
        </h1>

        <p class="mt-4 max-w-sm text-base leading-7 text-zinc-400">
          Informe seu email e enviaremos as instruções para você voltar ao acesso.
        </p>
      </div>

      {#if typedForm?.sent}
        <div class="rounded-2xl border border-cyan-300/20 bg-cyan-300/10 p-5">
          <p class="text-sm font-bold text-cyan-200">
            Email enviado.
          </p>

          <p class="mt-2 text-sm leading-6 text-zinc-400">
            Verifique sua caixa de entrada para continuar a recuperação da senha.
          </p>
        </div>

        <a
          href="/login"
          class="mt-6 flex h-14 w-full items-center justify-center rounded-2xl bg-white font-bold text-zinc-950 transition hover:scale-[1.01] active:scale-[0.985]"
        >
          Voltar para login
        </a>
      {:else}
        <form
          method="POST"
          action="?/forgotPassword"
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

          {#if typedForm?.errors?.general}
            <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
              {typedForm.errors.general}
            </p>
          {/if}

          <button type="submit" disabled={loading} class="auth-button">
            <span class="relative z-10">
              {loading ? 'Enviando...' : 'Enviar instruções'}
            </span>
          </button>
        </form>

        <div class="mt-7 flex items-center justify-center text-sm">
          <a href="/login" class="font-bold text-cyan-300 transition hover:text-cyan-200">
            Voltar para login
          </a>
        </div>
      {/if}
    </div>
  </div>
</div>