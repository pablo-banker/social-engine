<script lang="ts">
  import type { Snippet } from 'svelte'
  import AuthBackground from '$lib/components/auth/AuthBackground.svelte'

  import {
    SiGo,
    SiSvelte,
    SiPostgresql
  } from '@icons-pack/svelte-simple-icons'

  type Props = {
    children: Snippet
  }

  let { children }: Props = $props()
</script>

<section class="relative min-h-screen overflow-hidden">
  <AuthBackground />

  <div class="relative z-10 grid min-h-screen grid-cols-1 items-center px-6 py-10 lg:grid-cols-[560px_1fr] lg:px-20 xl:px-28">
    {@render children()}

    <div class="hidden lg:block">
      <div class="max-w-2xl">
        <p class="mb-5 text-sm font-bold uppercase tracking-[0.35em] text-cyan-300/80">
          Open social platform
        </p>

        <h2 class="text-6xl font-bold tracking-tight text-white xl:text-7xl">
          Create.
          <br />
          Share.
          <br />
          Connect.
        </h2>

        <p class="mt-7 max-w-xl text-lg leading-8 text-zinc-400">
          Uma engine social modular construída para posts, perfis, conexões,
          comentários, likes e comunidades modernas.
        </p>

        <div class="mt-10 grid max-w-xl grid-cols-3 gap-4">
          <div class="tech-card tech-card-go">
            <div class="mb-5 flex items-center justify-between">
              <span class="tech-icon">
                <SiGo size={44} color="#00ADD8" />
              </span>

              <span class="status-dot bg-cyan-400 shadow-[0_0_18px_rgba(34,211,238,0.9)]"></span>
            </div>

            <p class="text-2xl font-bold text-white">Go</p>
            <p class="mt-1 text-sm text-zinc-500">Backend</p>
          </div>

          <div class="tech-card tech-card-svelte">
            <div class="mb-5 flex items-center justify-between">
              <span class="tech-icon">
                <SiSvelte size={40} color="#FF3E00" />
              </span>

              <span class="status-dot bg-orange-400 shadow-[0_0_18px_rgba(251,146,60,0.9)]"></span>
            </div>

            <p class="text-2xl font-bold text-white">Svelte</p>
            <p class="mt-1 text-sm text-zinc-500">Frontend</p>
          </div>

          <div class="tech-card tech-card-postgres">
            <div class="mb-5 flex items-center justify-between">
              <span class="tech-icon">
                <SiPostgresql size={40} color="#4169E1" />
              </span>

              <span class="status-dot bg-sky-400 shadow-[0_0_18px_rgba(56,189,248,0.9)]"></span>
            </div>

            <p class="text-2xl font-bold text-white">PostgreSQL</p>
            <p class="mt-1 text-sm text-zinc-500">Database</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<style>
  :global(.auth-card) {
    animation: card-enter 800ms ease both;
  }

  :global(.auth-card-inner) {
    position: relative;
    overflow: hidden;
    border-radius: 2rem;
    padding: 2rem;
    border: 1px solid rgba(255, 255, 255, 0.08);
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.045), rgba(255, 255, 255, 0.015)),
      rgba(9, 9, 11, 0.82);

    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.06),
      inset 0 -1px 0 rgba(255, 255, 255, 0.02),
      0 25px 60px rgba(0, 0, 0, 0.5);

    backdrop-filter: blur(28px);
  }

  @media (min-width: 768px) {
    :global(.auth-card-inner) {
      padding: 2.5rem;
    }
  }

  :global(.auth-card-inner::before) {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    padding: 1px;
    background: linear-gradient(
      180deg,
      rgba(34, 211, 238, 0.18),
      rgba(255, 255, 255, 0.05),
      rgba(139, 92, 246, 0.12)
    );
    mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask-composite: exclude;
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    pointer-events: none;
    z-index: 3;
  }

  :global(.auth-shine) {
    position: absolute;
    inset: -120%;
    z-index: 2;
    background: linear-gradient(
      115deg,
      transparent 42%,
      rgba(255, 255, 255, 0.04),
      transparent 58%
    );
    transform: translateX(-100%) rotate(18deg);
    animation: shine-pass 10s ease-in-out infinite;
    pointer-events: none;
  }

  :global(.auth-top-gradient) {
    position: absolute;
    inset: 0;
    z-index: 1;
    background: radial-gradient(
      circle at 50% 0%,
      rgba(34, 211, 238, 0.13),
      transparent 38%
    );
    pointer-events: none;
  }

  :global(.auth-soft-highlight) {
    position: absolute;
    inset: 0;
    z-index: 1;
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.07),
      transparent 35%,
      transparent 70%,
      rgba(139, 92, 246, 0.06)
    );
    pointer-events: none;
  }

  :global(.auth-input) {
    height: 3.5rem;
    width: 100%;
    border-radius: 1rem;
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(255, 255, 255, 0.035);
    padding-inline: 1.25rem;
    color: white;
    outline: none;
    transition: 300ms;
  }

  :global(.auth-input::placeholder) {
    color: rgb(82, 82, 91);
  }

  :global(.auth-input:hover) {
    border-color: rgba(103, 232, 249, 0.3);
  }

  :global(.auth-input:focus) {
    border-color: rgba(103, 232, 249, 0.7);
    background: rgba(255, 255, 255, 0.055);
    box-shadow: 0 0 0 4px rgba(103, 232, 249, 0.1);
  }

  :global(.auth-button) {
    position: relative;
    height: 3.5rem;
    width: 100%;
    overflow: hidden;
    border-radius: 1rem;
    background: white;
    font-weight: 700;
    color: rgb(9, 9, 11);
    transition: 300ms;
  }

  :global(.auth-button:hover) {
    transform: scale(1.01);
    box-shadow: 0 25px 50px rgba(34, 211, 238, 0.1);
  }

  :global(.auth-button:active) {
    transform: scale(0.985);
  }

  :global(.auth-button:disabled) {
    cursor: not-allowed;
    opacity: 0.7;
  }

  :global(.auth-button::before) {
    content: '';
    position: absolute;
    inset: 0;
    transform: translateX(-120%);
    background: linear-gradient(
      90deg,
      transparent,
      rgba(34, 211, 238, 0.45),
      transparent
    );
    transition: transform 850ms ease;
  }

  :global(.auth-button:hover::before) {
    transform: translateX(120%);
  }

  .tech-card {
    position: relative;
    overflow: hidden;
    border-radius: 1.5rem;
    border: 1px solid rgba(255, 255, 255, 0.1);
    background:
      radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.055), transparent 45%),
      rgba(255, 255, 255, 0.025);
    padding: 1.25rem;
    backdrop-filter: blur(20px);
    transition:
      transform 260ms ease,
      border-color 260ms ease,
      background 260ms ease,
      box-shadow 260ms ease;
  }

  .tech-card::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(
      320px circle at 50% 0%,
      rgba(255, 255, 255, 0.08),
      transparent 60%
    );
    opacity: 0;
    transition: opacity 260ms ease;
    pointer-events: none;
  }

  .tech-card:hover {
    transform: translateY(-6px);
    background:
      radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.075), transparent 45%),
      rgba(255, 255, 255, 0.04);
  }

  .tech-card:hover::before {
    opacity: 1;
  }

  .tech-card-go:hover {
    border-color: rgba(34, 211, 238, 0.35);
    box-shadow: 0 18px 50px rgba(34, 211, 238, 0.08);
  }

  .tech-card-svelte:hover {
    border-color: rgba(249, 115, 22, 0.35);
    box-shadow: 0 18px 50px rgba(249, 115, 22, 0.08);
  }

  .tech-card-postgres:hover {
    border-color: rgba(56, 189, 248, 0.35);
    box-shadow: 0 18px 50px rgba(56, 189, 248, 0.08);
  }

  .tech-icon {
    display: inline-flex;
    position: relative;
    z-index: 1;
    opacity: 0.9;
    filter: drop-shadow(0 0 14px rgba(255, 255, 255, 0.12));
    transition:
      transform 260ms ease,
      opacity 260ms ease,
      filter 260ms ease;
  }

  .tech-card:hover .tech-icon {
    opacity: 1;
    transform: translateY(-2px) scale(1.08);
  }

  .tech-card-go:hover .tech-icon {
    filter: drop-shadow(0 0 18px rgba(34, 211, 238, 0.55));
  }

  .tech-card-svelte:hover .tech-icon {
    filter: drop-shadow(0 0 18px rgba(255, 62, 0, 0.55));
  }

  .tech-card-postgres:hover .tech-icon {
    filter: drop-shadow(0 0 18px rgba(65, 105, 225, 0.55));
  }

  .status-dot {
    position: relative;
    z-index: 1;
    width: 0.55rem;
    height: 0.55rem;
    border-radius: 9999px;
    animation: status-pulse 2.8s ease-in-out infinite;
  }

  @keyframes card-enter {
    from {
      opacity: 0;
      filter: blur(18px);
      transform: translateY(24px) scale(0.96);
    }

    to {
      opacity: 1;
      filter: blur(0);
      transform: translateY(0) scale(1);
    }
  }

  @keyframes shine-pass {
    0%, 48% {
      transform: translateX(-100%) rotate(18deg);
      opacity: 0;
    }

    58% {
      opacity: 1;
    }

    78%, 100% {
      transform: translateX(100%) rotate(18deg);
      opacity: 0;
    }
  }

  @keyframes status-pulse {
    0%, 100% {
      opacity: 0.45;
      transform: scale(0.85);
    }

    50% {
      opacity: 1;
      transform: scale(1.15);
    }
  }
</style>