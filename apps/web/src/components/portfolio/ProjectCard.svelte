<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    title: string;
    description: string;
    tags: string[];
    slug: string;
    lang: string;
    heroImage?: string;
    category?: string;
    delay?: number;
  }

  let { title, description, tags, slug, lang, heroImage, category, delay = 0 }: Props = $props();

  let element: HTMLDivElement;
  let visible = $state(false);

  const projectsPath = lang === 'es' ? 'proyectos' : 'projects';

  const categoryLabels: Record<string, Record<string, string>> = {
    web: { es: 'Web', en: 'Web' },
    system: { es: 'Sistema', en: 'System' },
    integration: { es: 'Integración', en: 'Integration' },
    mobile: { es: 'Móvil', en: 'Mobile' },
    game: { es: 'Videojuego', en: 'Game' },
    devops: { es: 'DevOps', en: 'DevOps' },
  };

  let categoryLabel = $derived(category ? (categoryLabels[category]?.[lang] || category) : '');

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

<a
  href={`/${lang}/${projectsPath}/${slug.replace(/^(es|en)\//, '')}/`}
  bind:this={element}
  class="card"
  style="opacity: {visible ? 1 : 0}; transform: {visible
    ? 'none'
    : 'translateY(24px)'}; transition: opacity 0.5s ease-out {delay}s, transform 0.5s ease-out {delay}s;"
>
  <!-- Image -->
  <div class="image-wrapper">
    {#if heroImage}
      <img src={heroImage} alt={title} class="image" loading="lazy" />
    {:else}
      <div class="image-placeholder">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
      </div>
    {/if}
    <div class="image-overlay">
      <p class="overlay-text">{description}</p>
    </div>
    {#if categoryLabel}
      <span class="category-badge">{categoryLabel}</span>
    {/if}
  </div>

  <!-- Info -->
  <div class="card-body">
    <h3 class="card-title">{title}</h3>
    <div class="tags">
      {#each tags.slice(0, 4) as tag}
        <span class="tag">{tag}</span>
      {/each}
      {#if tags.length > 4}
        <span class="tag tag-more">+{tags.length - 4}</span>
      {/if}
    </div>
  </div>
</a>

<style>
  .card {
    display: block;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    overflow: hidden;
    transition: border-color 0.2s, box-shadow 0.2s, transform 0.2s;
    text-decoration: none;
    color: inherit;
  }

  .card:hover {
    border-color: var(--color-accent-light);
    box-shadow: 0 8px 32px rgba(37, 99, 235, 0.08);
    transform: translateY(-2px);
  }

  .image-wrapper {
    position: relative;
    aspect-ratio: 16 / 10;
    overflow: hidden;
  }

  .image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .card:hover .image {
    transform: scale(1.03);
  }

  .image-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--color-accent-subtle), var(--color-bg-secondary));
    color: var(--color-accent-primary);
    transition: transform 0.3s ease;
  }

  .card:hover .image-placeholder {
    transform: scale(1.03);
  }

  .image-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: flex-end;
    padding: 1.25rem;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.7) 0%, transparent 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .card:hover .image-overlay {
    opacity: 1;
  }

  .overlay-text {
    color: #f1f5f9;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  .category-badge {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding: 0.25rem 0.625rem;
    border-radius: 9999px;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    color: #f1f5f9;
  }

  .card-body {
    padding: 1.25rem;
  }

  .card-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1rem;
    color: var(--color-text-primary);
    margin-bottom: 0.75rem;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
  }

  .tag {
    font-size: 0.75rem;
    padding: 0.2rem 0.5rem;
    border-radius: 0.375rem;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    font-weight: 500;
  }

  .tag-more {
    background: var(--color-bg-secondary);
    color: var(--color-text-muted);
  }
</style>
