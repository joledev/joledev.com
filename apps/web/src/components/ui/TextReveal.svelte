<script lang="ts">
  import { onMount } from 'svelte';
  import { animate, stagger } from 'animejs';

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

  // Check reduced motion preference
  const prefersReducedMotion =
    typeof window !== 'undefined' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  // Split text into spans
  function splitText(text: string, splitMode: 'letters' | 'words'): string[] {
    if (splitMode === 'letters') {
      return text.split('').map((char) => (char === ' ' ? '\u00A0' : char));
    }
    return text.split(/\s+/);
  }

  const pieces = $derived(splitText(text, mode));

  onMount(() => {
    if (prefersReducedMotion) {
      // Show immediately with a simple fade
      hasAnimated = true;
      const spans = container.querySelectorAll<HTMLElement>('.tr-piece');
      spans.forEach((el) => {
        el.style.opacity = '1';
        el.style.transform = 'none';
        el.style.clipPath = 'none';
      });
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
    const targets = container.querySelectorAll('.tr-piece');
    if (!targets.length) return;

    switch (effect) {
      case 'fadeUp':
        animate(targets, {
          opacity: [0, 1],
          translateY: ['1.2em', '0em'],
          delay: stagger(staggerDelay, { start: delay }),
          duration,
          ease: 'outQuint',
        });
        break;

      case 'clipReveal':
        animate(targets, {
          opacity: [0, 1],
          clipPath: ['inset(0 100% 0 0)', 'inset(0 0% 0 0)'],
          delay: stagger(staggerDelay, { start: delay }),
          duration: duration * 1.2,
          ease: 'outQuint',
        });
        break;

      case 'wave':
        animate(targets, {
          opacity: [0, 1],
          translateY: ['1em', '0em'],
          delay: stagger(staggerDelay, { start: delay }),
          duration,
          ease: 'outElastic(1, .6)',
        });
        break;

      case 'typewriter':
        animate(targets, {
          opacity: [0, 1],
          delay: stagger(staggerDelay * 0.8, { start: delay }),
          duration: 100,
          ease: 'linear',
        });
        break;
    }
  }
</script>

{#if tag === 'h1'}
  <h1 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0;">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h1>
{:else if tag === 'h2'}
  <h2 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0;">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h2>
{:else if tag === 'h3'}
  <h3 bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0;">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </h3>
{:else if tag === 'p'}
  <p bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0;">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </p>
{:else}
  <span bind:this={container} class="tr-container {className}">
    {#each pieces as piece, i}
      <span class="tr-piece" style="opacity: 0;">{piece}</span>{#if mode === 'words' && i < pieces.length - 1}{' '}{/if}
    {/each}
  </span>
{/if}

<style>
  .tr-container {
    overflow: hidden;
  }

  .tr-piece {
    display: inline-block;
    will-change: transform, opacity, clip-path;
  }
</style>
