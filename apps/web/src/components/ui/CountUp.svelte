<script lang="ts">
  import { onMount } from 'svelte';
  import { animate } from 'animejs';

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
    const obj = { val: 0 };
    animate(obj, {
      val: [0, value],
      duration,
      delay,
      ease: 'outExpo',
      onUpdate: () => {
        display = String(Math.round(obj.val));
      },
    });
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
