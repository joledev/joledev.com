<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    photos: string[];
    interval?: number;
  }

  let { photos, interval = 4000 }: Props = $props();

  let current = $state(0);
  let timer: ReturnType<typeof setInterval>;

  onMount(() => {
    timer = setInterval(() => {
      current = (current + 1) % photos.length;
    }, interval);

    return () => clearInterval(timer);
  });
</script>

<div class="photo-carousel">
  {#each photos as photo, i}
    <img
      src={photo}
      alt="Joel Lopez - AWS User Group Ensenada"
      class="photo-img"
      class:active={i === current}
      loading={i === 0 ? 'eager' : 'lazy'}
      decoding="async"
      width="600"
      height="600"
    />
  {/each}

  <!-- Dots -->
  <div class="photo-dots">
    {#each photos as _, i}
      <button
        class="photo-dot"
        class:active={i === current}
        onclick={() => { current = i; clearInterval(timer); timer = setInterval(() => { current = (current + 1) % photos.length; }, interval); }}
        aria-label="Foto {i + 1}"
      ></button>
    {/each}
  </div>
</div>

<style>
  .photo-carousel {
    position: relative;
    aspect-ratio: 1;
    overflow: hidden;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-secondary);
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  }

  .photo-img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    opacity: 0;
    transition: opacity 0.8s ease-in-out;
  }

  .photo-img.active {
    opacity: 1;
  }

  .photo-dots {
    position: absolute;
    bottom: 0.75rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 0.5rem;
    z-index: 2;
  }

  .photo-dot {
    width: 0.5rem;
    height: 0.5rem;
    border-radius: 50%;
    border: none;
    background: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: background 0.3s, transform 0.3s;
    padding: 0;
  }

  .photo-dot.active {
    background: white;
    transform: scale(1.3);
  }

  .photo-dot:hover {
    background: rgba(255, 255, 255, 0.8);
  }
</style>
