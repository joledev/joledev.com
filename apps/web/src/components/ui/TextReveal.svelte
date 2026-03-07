<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    text: string;
    tag?: 'h1' | 'h2' | 'h3' | 'p' | 'span';
    mode?: 'letters' | 'words';
    staggerDelay?: number;
    duration?: number;
    delay?: number;
    effect?: 'fadeUp' | 'clipReveal' | 'wave' | 'typewriter';
    class?: string;
    threshold?: number;
  }

  let {
    text,
    tag = 'span',
    mode = 'words',
    staggerDelay = 40,
    duration = 700,
    delay = 0,
    effect = 'fadeUp',
    class: className = '',
    threshold = 0.2,
  }: Props = $props();

  let container: HTMLElement;
  let hasAnimated = $state(false);

  const prefersReducedMotion =
    typeof window !== 'undefined' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  function splitText(text: string, splitMode: 'letters' | 'words'): string[] {
    if (splitMode === 'letters') {
      return text.split('').map((char) => (char === ' ' ? '\u00A0' : char));
    }
    return text.split(/\s+/);
  }

  const pieces = $derived(splitText(text, mode));

  function showImmediately() {
    const spans = container?.querySelectorAll<HTMLElement>('.tr-piece');
    spans?.forEach((el) => {
      el.style.opacity = '1';
      el.style.transform = 'none';
      el.style.clipPath = 'none';
    });
  }

  onMount(() => {
    if (prefersReducedMotion) {
      hasAnimated = true;
      showImmediately();
      return;
    }

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !hasAnimated) {
          hasAnimated = true;
          observer.unobserve(container);
          runAnimation();
        }
      },
      { threshold },
    );
    observer.observe(container);
    return () => observer.disconnect();
  });

  function runAnimation() {
    const targets = container.querySelectorAll<HTMLElement>('.tr-piece');
    if (!targets.length) return;

    // CSS transitions — works on all devices, no external dependency
    targets.forEach((el, i) => {
      const itemDelay = delay + i * staggerDelay;
      el.style.transition = `opacity ${duration}ms ease-out ${itemDelay}ms, transform ${duration}ms ease-out ${itemDelay}ms, clip-path ${duration}ms ease-out ${itemDelay}ms`;

      // Use requestAnimationFrame to ensure the transition triggers
      requestAnimationFrame(() => {
        el.style.opacity = '1';
        el.style.transform = 'none';
        if (effect === 'clipReveal') {
          el.style.clipPath = 'inset(0 0% 0 0)';
        }
      });
    });
  }
</script>

{#if tag === 'h1'}
  <h1 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0; transform: translateY(1em);">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h1>
{:else if tag === 'h2'}
  <h2 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0; transform: translateY(1em);">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h2>
{:else if tag === 'h3'}
  <h3 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0; transform: translateY(1em);">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h3>
{:else if tag === 'p'}
  <p bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0; transform: translateY(1em);">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </p>
{:else}
  <span bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0; transform: translateY(1em);">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </span>
{/if}

<style>
  .tr-container {
    overflow: hidden;
  }

  .tr-piece {
    display: inline-block;
  }
</style>
