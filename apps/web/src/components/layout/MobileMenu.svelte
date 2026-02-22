<script lang="ts">
  interface Props {
    lang: 'es' | 'en';
    navItems: { label: string; href: string }[];
  }

  let { lang, navItems }: Props = $props();
  let isOpen = $state(false);

  function toggle() {
    isOpen = !isOpen;
  }

  function close() {
    isOpen = false;
  }
</script>

<!-- Hamburger button -->
<button
  class="hamburger"
  onclick={toggle}
  aria-label={isOpen ? 'Close menu' : 'Open menu'}
  aria-expanded={isOpen}
  type="button"
>
  <span class="hamburger-line" class:open={isOpen}></span>
  <span class="hamburger-line" class:open={isOpen}></span>
  <span class="hamburger-line" class:open={isOpen}></span>
</button>

<!-- Mobile overlay -->
{#if isOpen}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="overlay" onclick={close} onkeydown={close}></div>
  <nav class="mobile-nav">
    {#each navItems as item}
      <a href={item.href} class="mobile-link" onclick={close}>
        {item.label}
      </a>
    {/each}
  </nav>
{/if}

<style>
  .hamburger {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 0.5rem;
    background: transparent;
    border: none;
    cursor: pointer;
  }

  .hamburger-line {
    display: block;
    width: 24px;
    height: 2px;
    background-color: var(--color-text-primary);
    transition: transform 0.3s, opacity 0.3s;
  }

  .hamburger-line.open:nth-child(1) {
    transform: rotate(45deg) translate(5px, 5px);
  }

  .hamburger-line.open:nth-child(2) {
    opacity: 0;
  }

  .hamburger-line.open:nth-child(3) {
    transform: rotate(-45deg) translate(5px, -5px);
  }

  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 40;
  }

  .mobile-nav {
    position: fixed;
    top: 0;
    right: 0;
    width: 280px;
    height: 100vh;
    background: var(--color-bg-elevated);
    border-left: 1px solid var(--color-border);
    padding: 5rem 1.5rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    z-index: 50;
  }

  .mobile-link {
    display: block;
    padding: 0.75rem 1rem;
    color: var(--color-text-primary);
    text-decoration: none;
    border-radius: 0.5rem;
    transition: background-color 0.2s;
  }

  .mobile-link:hover {
    background-color: var(--color-bg-secondary);
  }
</style>
