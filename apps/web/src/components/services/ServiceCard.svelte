<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    title: string;
    description: string;
    icon: string;
    delay?: number;
  }

  let { title, description, icon, delay = 0 }: Props = $props();

  let element: HTMLDivElement;
  let visible = $state(false);

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

    return () => observer.disconnect();
  });
</script>

<div
  bind:this={element}
  class="card"
  style="opacity: {visible ? 1 : 0}; transform: {visible
    ? 'none'
    : 'translateY(24px)'}; transition: opacity 0.5s ease-out {delay}s, transform 0.5s ease-out {delay}s, border-color 0.3s, box-shadow 0.3s;"
  role="article"
>
  <div class="icon-wrapper">
    {@html icon}
  </div>
  <h3 class="card-title">{title}</h3>
  <p class="card-description">{description}</p>
</div>

<style>
  .card {
    padding: 2rem 1.5rem;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-glass);
    backdrop-filter: blur(8px);
    cursor: default;
    transition: transform 200ms ease, border-color 200ms ease, box-shadow 200ms ease;
  }

  .card:hover {
    transform: translateY(-4px) !important;
    border-color: var(--color-accent-primary);
    box-shadow: 0 8px 24px rgba(37, 99, 235, 0.15);
  }

  .icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
    border-radius: 0.75rem;
    background: var(--color-secondary-subtle);
    color: var(--color-secondary-primary);
    margin-bottom: 1.25rem;
    transition: background 0.3s ease, color 0.3s ease;
  }

  .card:hover .icon-wrapper {
    background: var(--color-secondary-primary);
    color: #fff;
  }

  .icon-wrapper :global(svg) {
    width: 1.5rem;
    height: 1.5rem;
  }

  .card-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    margin-bottom: 0.5rem;
    color: var(--color-text-primary);
  }

  .card-description {
    font-size: 0.875rem;
    line-height: 1.6;
    color: var(--color-text-secondary);
  }
</style>
