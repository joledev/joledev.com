<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    value: number;
    suffix?: string;
    prefix?: string;
    duration?: number;
    delay?: number;
    class?: string;
  }

  let {
    value,
    suffix = '',
    prefix = '',
    duration = 2000,
    delay = 0,
    class: className = '',
  }: Props = $props();

  let display = $state('0');
  let element: HTMLSpanElement;

  const prefersReducedMotion =
    typeof window !== 'undefined' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  onMount(() => {
    if (prefersReducedMotion) {
      display = String(value);
      return;
    }

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          observer.unobserve(element);
          runAnimation();
        }
      },
      { threshold: 0.3 },
    );
    observer.observe(element);
    return () => observer.disconnect();
  });

  function runAnimation() {
    const start = performance.now() + delay;
    const end = start + duration;

    function tick(now: number) {
      if (now < start) {
        requestAnimationFrame(tick);
        return;
      }
      const progress = Math.min((now - start) / duration, 1);
      // easeOutExpo
      const eased = progress === 1 ? 1 : 1 - Math.pow(2, -10 * progress);
      display = String(Math.round(eased * value));
      if (progress < 1) {
        requestAnimationFrame(tick);
      }
    }
    requestAnimationFrame(tick);
  }
</script>

<span bind:this={element} class="countup {className}">
  {prefix}{display}{suffix}
</span>

<style>
  .countup {
    font-variant-numeric: tabular-nums;
  }
</style>
