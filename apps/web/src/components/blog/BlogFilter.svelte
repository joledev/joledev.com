<script lang="ts">
  interface BlogPost {
    id: string;
    data: {
      title: string;
      description: string;
      pubDate: Date;
      tags: string[];
      category: string;
      lang: string;
      heroImage?: string;
      draft: boolean;
    };
  }

  interface Props {
    posts: BlogPost[];
    lang: 'es' | 'en';
    categories: { value: string; label: string }[];
    noPostsMessage: string;
    readMoreLabel: string;
    blogBasePath: string;
  }

  let { posts, lang, categories, noPostsMessage, readMoreLabel, blogBasePath }: Props = $props();

  let activeCategory = $state('all');

  let filteredPosts = $derived(
    activeCategory === 'all'
      ? posts
      : posts.filter((p) => p.data.category === activeCategory),
  );

  function formatDate(date: Date): string {
    return new Date(date).toLocaleDateString(lang === 'es' ? 'es-MX' : 'en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }

  const categoryLabels: Record<string, string> = {
    tutorial: lang === 'es' ? 'Tutorial' : 'Tutorial',
    'case-study': lang === 'es' ? 'Caso de estudio' : 'Case Study',
    opinion: lang === 'es' ? 'Opinión' : 'Opinion',
    project: lang === 'es' ? 'Proyecto' : 'Project',
  };
</script>

<!-- Category filters -->
<div class="filters">
  <button
    class="filter-pill"
    class:active={activeCategory === 'all'}
    onclick={() => (activeCategory = 'all')}
    type="button"
  >
    {lang === 'es' ? 'Todas' : 'All'}
  </button>
  {#each categories as cat}
    <button
      class="filter-pill"
      class:active={activeCategory === cat.value}
      onclick={() => (activeCategory = cat.value)}
      type="button"
    >
      {cat.label}
    </button>
  {/each}
</div>

<!-- Posts grid -->
{#if filteredPosts.length === 0}
  <p class="no-posts">{noPostsMessage}</p>
{:else}
  <div class="posts-grid">
    {#each filteredPosts as post (post.id)}
      <a href={`${blogBasePath}${post.id.replace(/^(es|en)\//, '')}/`} class="post-card">
        <div class="post-image"></div>
        <div class="post-body">
          <div class="post-meta">
            <time datetime={new Date(post.data.pubDate).toISOString()}>
              {formatDate(post.data.pubDate)}
            </time>
            <span class="category-badge">{categoryLabels[post.data.category] ?? post.data.category}</span>
          </div>
          <h3 class="post-title">{post.data.title}</h3>
          <p class="post-desc">{post.data.description}</p>
          <div class="post-tags">
            {#each post.data.tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        </div>
      </a>
    {/each}
  </div>
{/if}

<!-- TODO: Implement pagination when posts > 9 -->

<style>
  .filters {
    display: flex;
    gap: 0.5rem;
    overflow-x: auto;
    padding-bottom: 0.5rem;
    margin-bottom: 2rem;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
  }

  .filters::-webkit-scrollbar {
    display: none;
  }

  .filter-pill {
    flex-shrink: 0;
    padding: 0.5rem 1.25rem;
    border-radius: 9999px;
    border: 1px solid var(--color-border);
    background: transparent;
    color: var(--color-text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    font-family: var(--font-sans);
  }

  .filter-pill:hover {
    background: var(--color-accent-subtle);
    border-color: var(--color-accent-light);
    color: var(--color-accent-primary);
  }

  .filter-pill.active {
    background: var(--color-accent-primary);
    border-color: var(--color-accent-primary);
    color: #fff;
  }

  .no-posts {
    text-align: center;
    color: var(--color-text-muted);
    padding: 4rem 0;
    font-size: 1.125rem;
  }

  .posts-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: 1fr;
  }

  @media (min-width: 640px) {
    .posts-grid { grid-template-columns: repeat(2, 1fr); }
  }

  @media (min-width: 1024px) {
    .posts-grid { grid-template-columns: repeat(3, 1fr); }
  }

  .post-card {
    display: block;
    text-decoration: none;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    overflow: hidden;
    transition: border-color 0.2s, box-shadow 0.2s, transform 0.2s;
  }

  .post-card:hover {
    border-color: var(--color-accent-light);
    box-shadow: 0 8px 24px rgba(37, 99, 235, 0.06);
    transform: translateY(-2px);
  }

  .post-image {
    aspect-ratio: 16 / 9;
    background: linear-gradient(135deg, var(--color-accent-subtle), var(--color-bg-secondary));
  }

  .post-body {
    padding: 1.25rem;
  }

  .post-meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  .post-meta time {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }

  .category-badge {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding: 0.15rem 0.5rem;
    border-radius: 9999px;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
  }

  .post-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    line-height: 1.4;
    color: var(--color-text-primary);
    margin-bottom: 0.5rem;
    transition: color 0.2s;
  }

  .post-card:hover .post-title {
    color: var(--color-accent-primary);
  }

  .post-desc {
    font-size: 0.875rem;
    line-height: 1.6;
    color: var(--color-text-secondary);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    margin-bottom: 0.75rem;
  }

  .post-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
  }

  .tag {
    font-size: 0.75rem;
    padding: 0.15rem 0.5rem;
    border-radius: 0.375rem;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    font-weight: 500;
  }
</style>
