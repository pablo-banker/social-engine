<div class="absolute inset-0 overflow-hidden bg-black">
  <div class="aurora aurora-one"></div>
  <div class="aurora aurora-two"></div>
  <div class="aurora aurora-three"></div>

  <svg class="connections" viewBox="0 0 1200 800" preserveAspectRatio="none">
    <path d="M40 520 L180 430 L310 500 L460 360 L610 420" />
    <path d="M690 210 L800 280 L940 230 L1100 340" />
    <path d="M760 520 L890 450 L1020 520 L1160 410" />

    <circle cx="40" cy="520" r="4" />
    <circle cx="180" cy="430" r="4" />
    <circle cx="310" cy="500" r="4" />
    <circle cx="460" cy="360" r="4" />
    <circle cx="610" cy="420" r="4" />

    <circle cx="690" cy="210" r="4" />
    <circle cx="800" cy="280" r="4" />
    <circle cx="940" cy="230" r="4" />
    <circle cx="1100" cy="340" r="4" />

    <circle cx="760" cy="520" r="4" />
    <circle cx="890" cy="450" r="4" />
    <circle cx="1020" cy="520" r="4" />
    <circle cx="1160" cy="410" r="4" />
  </svg>

  <div class="grid"></div>
  <div class="vignette"></div>
  <div class="noise"></div>
</div>

<style>
  .aurora {
    position: absolute;
    border-radius: 9999px;
    filter: blur(130px);
    opacity: 0.18;
    animation: aurora-float 10s ease-in-out infinite alternate;
  }

  .aurora-one {
    width: 760px;
    height: 760px;
    left: 32%;
    top: -35%;
    background: #22d3ee;
  }

  .aurora-two {
    width: 560px;
    height: 560px;
    right: -12%;
    bottom: 4%;
    background: #8b5cf6;
    animation-delay: -4s;
  }

  .aurora-three {
    width: 440px;
    height: 440px;
    left: -8%;
    bottom: -18%;
    background: #06b6d4;
    animation-delay: -7s;
  }

  .connections {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    opacity: 0.42;
  }

  .connections path {
    fill: none;
    stroke: rgba(34, 211, 238, 0.32);
    stroke-width: 1;
    stroke-dasharray: 8 14;
    animation: dash-flow 8s linear infinite;
  }

  .connections path:nth-child(2) {
    stroke: rgba(139, 92, 246, 0.32);
    animation-duration: 10s;
  }

  .connections path:nth-child(3) {
    stroke: rgba(34, 211, 238, 0.24);
    animation-duration: 12s;
  }

  .connections circle {
    fill: #22d3ee;
    filter: drop-shadow(0 0 12px rgba(34, 211, 238, 0.9));
    animation: node-pulse 2.6s ease-in-out infinite;
  }

  .connections circle:nth-of-type(3n) {
    fill: #8b5cf6;
    filter: drop-shadow(0 0 12px rgba(139, 92, 246, 0.9));
    animation-delay: -1s;
  }

  .connections circle:nth-of-type(4n) {
    animation-delay: -1.7s;
  }

  .grid {
    position: absolute;
    inset: 0;
    opacity: 0.055;
    background-image:
      linear-gradient(to right, #ffffff 1px, transparent 1px),
      linear-gradient(to bottom, #ffffff 1px, transparent 1px);
    background-size: 64px 64px;
    animation: grid-drift 18s linear infinite;
    mask-image: radial-gradient(circle at center, black 0%, transparent 78%);
  }

  .vignette {
    position: absolute;
    inset: 0;
    background:
      radial-gradient(circle at 55% 35%, transparent 0%, rgba(0, 0, 0, 0.68) 48%, #000 88%),
      linear-gradient(to bottom, transparent, #000);
  }

  .noise {
    position: absolute;
    inset: 0;
    opacity: 0.035;
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='256' height='256' filter='url(%23n)' opacity='0.7'/%3E%3C/svg%3E");
  }

  @keyframes aurora-float {
    from {
      transform: translate3d(-28px, -18px, 0) scale(1);
    }
    to {
      transform: translate3d(32px, 24px, 0) scale(1.1);
    }
  }

  @keyframes dash-flow {
    from {
      stroke-dashoffset: 0;
    }
    to {
      stroke-dashoffset: -220;
    }
  }

  @keyframes node-pulse {
    0%, 100% {
      opacity: 0.28;
      transform: scale(0.75);
    }
    50% {
      opacity: 1;
      transform: scale(1.35);
    }
  }

  @keyframes grid-drift {
    from {
      background-position: 0 0;
    }
    to {
      background-position: 64px 64px;
    }
  }
</style>