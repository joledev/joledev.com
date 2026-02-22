<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    direction?: 'up' | 'down' | 'left' | 'right';
    delay?: number;
    duration?: number;
    threshold?: number;
    class?: string;
  }

  let {
    direction = 'up',
    delay = 0,
    duration = 0.6,
    threshold = 0.1,
    class: className = '',
  }: Props = $props();

  let element: HTMLDivElement;
  let visible = $state(false);

  const transforms: Record<string, string> = {
    up: 'translateY(24px)',
    down: 'translateY(-24px)',
    left: 'translateX(24px)',
    right: 'translateX(-24px)',
  };

  onMount(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          visible = true;
          observer.unobserve(element);
        }
      },
      { threshold },
    );
    observer.observe(element);
    return () => observer.disconnect();
  });
</script>

<div
  bind:this={element}
  class={className}
  style="opacity: {visible ? 1 : 0}; transform: {visible
    ? 'none'
    : transforms[direction]}; transition: opacity {duration}s ease-out {delay}s, transform {duration}s ease-out {delay}s;"
>
  <slot />
</div>
