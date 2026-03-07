<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    strength?: number;
    class?: string;
  }

  let { strength = 0.3, class: className = '' }: Props = $props();

  let wrapper: HTMLDivElement;
  let isMobile = $state(true);

  onMount(() => {
    isMobile = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
  });

  function handleMouseMove(e: MouseEvent) {
    if (isMobile) return;
    const rect = wrapper.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const deltaX = (e.clientX - centerX) * strength;
    const deltaY = (e.clientY - centerY) * strength;

    const maxOffset = 15;
    const clampedX = Math.max(-maxOffset, Math.min(maxOffset, deltaX));
    const clampedY = Math.max(-maxOffset, Math.min(maxOffset, deltaY));

    wrapper.style.transform = `translate(${clampedX}px, ${clampedY}px)`;
  }

  function handleMouseLeave() {
    if (isMobile) return;
    wrapper.style.transform = '';
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
    transition: transform 0.3s ease-out;
  }
</style>
