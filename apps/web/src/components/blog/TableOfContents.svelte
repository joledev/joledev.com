<script lang="ts">
  import { onMount } from 'svelte';

  interface Heading {
    depth: number;
    slug: string;
    text: string;
  }

  interface Props {
    headings: Heading[];
    label: string;
  }

  let { headings, label }: Props = $props();

  let activeSlug = $state('');
  let isOpen = $state(false);

  const tocHeadings = headings.filter((h) => h.depth === 2 || h.depth === 3);

  onMount(() => {
    const elements = tocHeadings
      .map((h) => document.getElementById(h.slug))
      .filter(Boolean) as HTMLElement[];

    if (elements.length === 0) return;

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            activeSlug = entry.target.id;
          }
        }
      },
      { rootMargin: '-80px 0px -70% 0px', threshold: 0 },
    );

    for (const el of elements) observer.observe(el);
    return () => observer.disconnect();
  });
</script>

{#if tocHeadings.length > 0}
  <!-- Mobile: collapsible -->
  <div class="toc-mobile">
    <button class="toc-toggle" onclick={() => (isOpen = !isOpen)} type="button">
      {label}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        class="chevron"
        class:open={isOpen}
      >
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </button>
    {#if isOpen}
      <nav class="toc-nav">
        {#each tocHeadings as heading}
          <a
            href={`#${heading.slug}`}
            class="toc-link"
            class:depth-3={heading.depth === 3}
            class:active={activeSlug === heading.slug}
            onclick={() => (isOpen = false)}
          >
            {heading.text}
          </a>
        {/each}
      </nav>
    {/if}
  </div>

  <!-- Desktop: sticky sidebar -->
  <nav class="toc-desktop">
    <p class="toc-title">{label}</p>
    {#each tocHeadings as heading}
      <a
        href={`#${heading.slug}`}
        class="toc-link"
        class:depth-3={heading.depth === 3}
        class:active={activeSlug === heading.slug}
      >
        {heading.text}
      </a>
    {/each}
  </nav>
{/if}

<style>
  .toc-mobile {
    display: block;
    margin-bottom: 2rem;
    border: 1px solid var(--color-border);
    border-radius: 0.75rem;
    overflow: hidden;
  }

  .toc-toggle {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 0.75rem 1rem;
    background: var(--color-bg-secondary);
    border: none;
    color: var(--color-text-primary);
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
  }

  .chevron {
    transition: transform 0.2s;
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .toc-nav {
    padding: 0.5rem;
  }

  .toc-desktop {
    display: none;
  }

  @media (min-width: 1024px) {
    .toc-mobile {
      display: none;
    }

    .toc-desktop {
      display: block;
      position: sticky;
      top: 5rem;
      max-height: calc(100vh - 6rem);
      overflow-y: auto;
      padding-left: 1.5rem;
      border-left: 1px solid var(--color-border);
    }
  }

  .toc-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
    margin-bottom: 0.75rem;
  }

  .toc-link {
    display: block;
    padding: 0.3rem 0.75rem;
    font-size: 0.8125rem;
    line-height: 1.5;
    color: var(--color-text-secondary);
    text-decoration: none;
    border-left: 2px solid transparent;
    transition: color 0.2s, border-color 0.2s;
  }

  .toc-link:hover {
    color: var(--color-accent-primary);
  }

  .toc-link.active {
    color: var(--color-accent-primary);
    border-left-color: var(--color-accent-primary);
    font-weight: 500;
  }

  .toc-link.depth-3 {
    padding-left: 1.5rem;
    font-size: 0.78rem;
  }
</style>
