<script lang="ts">
  import { onMount } from 'svelte';
  import { animate } from 'animejs';

  interface Props {
    strength?: number;
    class?: string;
  }

  let { strength = 0.3, class: className = '' }: Props = $props();

  let wrapper: HTMLDivElement;
  let isMobile = $state(true);

  onMount(() => {
    // Disable on touch devices
    isMobile = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
  });

  function handleMouseMove(e: MouseEvent) {
    if (isMobile) return;
    const rect = wrapper.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const deltaX = (e.clientX - centerX) * strength;
    const deltaY = (e.clientY - centerY) * strength;

    // Clamp to max offset
    const maxOffset = 15;
    const clampedX = Math.max(-maxOffset, Math.min(maxOffset, deltaX));
    const clampedY = Math.max(-maxOffset, Math.min(maxOffset, deltaY));

    animate(wrapper, {
      translateX: clampedX,
      translateY: clampedY,
      duration: 300,
      ease: 'outQuad',
    });
  }

  function handleMouseLeave() {
    if (isMobile) return;
    animate(wrapper, {
      translateX: 0,
      translateY: 0,
      duration: 600,
      ease: 'outElastic(1, .5)',
    });
  }
</script>

<div
  bind:this={wrapper}
  class="magnetic-wrapper {className}"
  onmousemove={handleMouseMove}
  onmouseleave={handleMouseLeave}
  role="presentation"
>
  <slot />
</div>

<style>
  .magnetic-wrapper {
    display: inline-block;
    will-change: transform;
  }
</style>
