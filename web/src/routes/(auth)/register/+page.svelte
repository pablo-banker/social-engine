<script lang="ts">
  import { enhance } from '$app/forms'
  import type { PageProps } from './$types'

  type RegisterForm = {
    errors?: {
      general?: string
      firstName?: string
      lastName?: string
      email?: string
      password?: string
      confirmPassword?: string
    }
    values?: {
      firstName?: string
      lastName?: string
      email?: string
    }
  }

  let { form }: PageProps = $props()

  let step = $state(1)
  let loading = $state(false)

  let firstName = $state('')
  let lastName = $state('')
  let email = $state('')
  let password = $state('')
  let confirmPassword = $state('')

  let clientError = $state('')

  const typedForm = $derived(form as RegisterForm | undefined)

  $effect(() => {
    if (typedForm?.values?.firstName) firstName = typedForm.values.firstName
    if (typedForm?.values?.lastName) lastName = typedForm.values.lastName
    if (typedForm?.values?.email) email = typedForm.values.email
  })

  function nextStep() {
    clientError = ''

    if (step === 1 && (!firstName.trim() || !lastName.trim())) {
      clientError = 'Informe nome e sobrenome.'
      return
    }

    if (step === 2 && !email.trim()) {
      clientError = 'Informe seu email.'
      return
    }

    if (step < 3) {
      step += 1
    }
  }

  function previousStep() {
    clientError = ''

    if (step > 1) {
      step -= 1
    }
  }

  function validateBeforeSubmit(event: SubmitEvent) {
    clientError = ''

    if (step < 3) {
      event.preventDefault()
      nextStep()
      return
    }

    if (!password || !confirmPassword) {
      event.preventDefault()
      clientError = 'Informe e confirme sua senha.'
      return
    }

    if (password !== confirmPassword) {
      event.preventDefault()
      clientError = 'As senhas não conferem.'
    }
  }
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
          Crie sua conta.
        </h1>

        <p class="mt-4 max-w-sm text-base leading-7 text-zinc-400">
          Comece a publicar ideias, seguir pessoas e construir sua própria rede.
        </p>
      </div>

      <div class="mb-7 grid grid-cols-3 gap-2">
        {#each [1, 2, 3] as item}
          <div
            class={[
              'h-1.5 rounded-full transition',
              item <= step ? 'bg-cyan-300 shadow-[0_0_14px_rgba(103,232,249,0.7)]' : 'bg-white/10'
            ]}
          ></div>
        {/each}
      </div>

      <form
        method="POST"
        action="?/register"
        onsubmit={validateBeforeSubmit}
        use:enhance={() => {
          loading = true

          return async ({ update }) => {
            await update()
            loading = false
          }
        }}
        class="space-y-5"
      >
        <div class={step === 1 ? 'space-y-5' : 'hidden'}>
          <div>
            <label for="firstName" class="mb-2 block text-sm font-semibold text-zinc-300">
              Nome
            </label>

            <input
              id="firstName"
              name="firstName"
              bind:value={firstName}
              type="text"
              placeholder="Pablo"
              autocomplete="given-name"
              class="auth-input"
            />

            {#if typedForm?.errors?.firstName}
              <p class="mt-2 text-sm text-red-300">{typedForm.errors.firstName}</p>
            {/if}
          </div>

          <div>
            <label for="lastName" class="mb-2 block text-sm font-semibold text-zinc-300">
              Sobrenome
            </label>

            <input
              id="lastName"
              name="lastName"
              bind:value={lastName}
              type="text"
              placeholder="Banker"
              autocomplete="family-name"
              class="auth-input"
            />

            {#if typedForm?.errors?.lastName}
              <p class="mt-2 text-sm text-red-300">{typedForm.errors.lastName}</p>
            {/if}
          </div>
        </div>

        <div class={step === 2 ? 'space-y-5' : 'hidden'}>
          <div>
            <label for="email" class="mb-2 block text-sm font-semibold text-zinc-300">
              Email
            </label>

            <input
              id="email"
              name="email"
              bind:value={email}
              type="email"
              placeholder="seu@email.com"
              autocomplete="email"
              class="auth-input"
            />

            {#if typedForm?.errors?.email}
              <p class="mt-2 text-sm text-red-300">{typedForm.errors.email}</p>
            {/if}
          </div>
        </div>

        <div class={step === 3 ? 'space-y-5' : 'hidden'}>
          <div>
            <label for="password" class="mb-2 block text-sm font-semibold text-zinc-300">
              Senha
            </label>

            <input
              id="password"
              name="password"
              bind:value={password}
              type="password"
              placeholder="••••••••"
              autocomplete="new-password"
              class="auth-input"
            />

            {#if typedForm?.errors?.password}
              <p class="mt-2 text-sm text-red-300">{typedForm.errors.password}</p>
            {/if}
          </div>

          <div>
            <label for="confirmPassword" class="mb-2 block text-sm font-semibold text-zinc-300">
              Confirmar senha
            </label>

            <input
              id="confirmPassword"
              name="confirmPassword"
              bind:value={confirmPassword}
              type="password"
              placeholder="••••••••"
              autocomplete="new-password"
              class="auth-input"
            />

            {#if typedForm?.errors?.confirmPassword}
              <p class="mt-2 text-sm text-red-300">{typedForm.errors.confirmPassword}</p>
            {/if}
          </div>
        </div>

        {#if clientError}
          <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
            {clientError}
          </p>
        {/if}

        {#if typedForm?.errors?.general}
          <p class="rounded-2xl border border-red-400/20 bg-red-400/10 px-4 py-3 text-sm text-red-200">
            {typedForm.errors.general}
          </p>
        {/if}

        <div class="flex gap-3">
          {#if step > 1}
            <button
              type="button"
              onclick={previousStep}
              class="h-14 flex-1 rounded-2xl border border-white/10 font-bold text-zinc-300 transition hover:border-white/20 hover:bg-white/[0.04] hover:text-white active:scale-[0.985]"
            >
              Voltar
            </button>
          {/if}

          {#if step < 3}
            <button
              type="button"
              onclick={nextStep}
              class="auth-button flex-1"
            >
              <span class="relative z-10">
                Continuar
              </span>
            </button>
          {:else}
            <button type="submit" disabled={loading} class="auth-button flex-1">
              <span class="relative z-10">
                {loading ? 'Criando...' : 'Criar conta'}
              </span>
            </button>
          {/if}
        </div>
      </form>

      <div class="mt-7 flex items-center justify-center text-sm">
        <span class="text-zinc-500">
          Já tem uma conta?
        </span>

        <a href="/login" class="ml-2 font-bold text-cyan-300 transition hover:text-cyan-200">
          Entrar
        </a>
      </div>
    </div>
  </div>
</div>