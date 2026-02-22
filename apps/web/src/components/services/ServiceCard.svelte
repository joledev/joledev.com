<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    title: string;
    description: string;
    icon: string;
    image?: string;
    delay?: number;
  }

  let { title, description, icon, image = '', delay = 0 }: Props = $props();

  let element: HTMLDivElement;
  let visible = $state(false);
  let hovered = $state(false);
  let imgLoaded = $state(false);

  onMount(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          visible = true;
          observer.unobserve(element);
        }
      },
      { threshold: 0.1 },
    );
    observer.observe(element);

    // Preload image
    if (image) {
      const img = new Image();
      img.onload = () => { imgLoaded = true; };
      img.src = image;
    }

    return () => observer.disconnect();
  });
</script>

<div
  bind:this={element}
  class="card"
  class:hovered
  style="opacity: {visible ? 1 : 0}; transform: {visible
    ? 'none'
    : 'translateY(24px)'}; transition: opacity 0.5s ease-out {delay}s, transform 0.5s ease-out {delay}s, border-color 0.3s, box-shadow 0.3s;"
  role="article"
  onmouseenter={() => (hovered = true)}
  onmouseleave={() => (hovered = false)}
>
  {#if image && imgLoaded}
    <div class="card-bg" style="background-image: url({image});"></div>
  {/if}
  <div class="card-content">
    <div class="icon-wrapper">
      {@html icon}
    </div>
    <h3 class="card-title">{title}</h3>
    <p class="card-description">{description}</p>
  </div>
</div>

<style>
  .card {
    position: relative;
    padding: 1.5rem;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-glass);
    backdrop-filter: blur(8px);
    overflow: hidden;
    cursor: default;
    transition:
      border-color 0.3s ease,
      box-shadow 0.3s ease;
  }

  .card:hover {
    transform: translateY(-4px) !important;
    border-color: var(--color-accent-light);
    box-shadow: 0 8px 32px rgba(37, 99, 235, 0.12);
  }

  /* Background image layer */
  .card-bg {
    position: absolute;
    inset: 0;
    background-size: cover;
    background-position: center;
    opacity: 0;
    transition: opacity 0.4s ease;
  }

  .card-bg::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to bottom,
      rgba(0, 0, 0, 0.45) 0%,
      rgba(0, 0, 0, 0.7) 100%
    );
  }

  .card.hovered .card-bg {
    opacity: 1;
  }

  /* Content layer */
  .card-content {
    position: relative;
    z-index: 1;
    transition: color 0.3s ease;
  }

  .icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 0.75rem;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    margin-bottom: 1rem;
    transition: background 0.3s ease, color 0.3s ease;
  }

  .card.hovered .icon-wrapper {
    background: rgba(37, 99, 235, 0.25);
    color: #93c5fd;
  }

  .icon-wrapper :global(svg) {
    width: 1.25rem;
    height: 1.25rem;
  }

  .card-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    margin-bottom: 0.5rem;
    color: var(--color-text-primary);
    transition: color 0.3s ease;
  }

  .card.hovered .card-title {
    color: #ffffff;
  }

  .card-description {
    font-size: 0.875rem;
    line-height: 1.6;
    color: var(--color-text-secondary);
    transition: color 0.3s ease;
  }

  .card.hovered .card-description {
    color: rgba(255, 255, 255, 0.85);
  }
</style>
