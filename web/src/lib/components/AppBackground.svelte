<!-- Fundo ambiente do app: versão discreta das auroras do auth.
     Fica atrás de todo o conteúdo (o conteúdo do layout usa z-10). -->
<div class="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
  <div class="glow glow-one"></div>
  <div class="glow glow-two"></div>
  <div class="noise"></div>
</div>

<style>
  .glow {
    position: absolute;
    border-radius: 9999px;
    filter: blur(150px);
    opacity: 0.07;
    animation: glow-float 16s ease-in-out infinite alternate;
  }

  .glow-one {
    width: 900px;
    height: 900px;
    left: 16%;
    top: -34%;
    background: #22d3ee;
  }

  .glow-two {
    width: 700px;
    height: 700px;
    right: -12%;
    bottom: -26%;
    background: #8b5cf6;
    animation-delay: -7s;
  }

  .noise {
    position: absolute;
    inset: 0;
    opacity: 0.03;
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='256' height='256' filter='url(%23n)' opacity='0.7'/%3E%3C/svg%3E");
  }

  @keyframes glow-float {
    from {
      transform: translate3d(-30px, -20px, 0) scale(1);
    }

    to {
      transform: translate3d(32px, 26px, 0) scale(1.08);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .glow {
      animation: none;
    }
  }
</style>