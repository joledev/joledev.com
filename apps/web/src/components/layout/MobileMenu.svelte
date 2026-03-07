<script lang="ts">
  import { onMount } from 'svelte';

  interface Props {
    lang: 'es' | 'en';
    navItems: { label: string; href: string }[];
  }

  let { lang, navItems }: Props = $props();
  let isOpen = $state(false);

  const icons: Record<string, string> = {
    Servicios: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z',
    Services: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z',
    Proyectos: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z',
    Projects: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z',
    Blog: 'M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z',
    Cotizar: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z',
    'Get a Quote': 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z',
    Agendar: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
    Schedule: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
  };

  function toggle() {
    isOpen = !isOpen;
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
  }

  function close() {
    isOpen = false;
    document.body.style.overflow = '';
  }

  onMount(() => {
    return () => { document.body.style.overflow = ''; };
  });
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

<!-- Mobile overlay + panel -->
{#if isOpen}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="overlay" onclick={close} onkeydown={close}></div>
  <nav class="mobile-nav">
    <div class="nav-header">
      <span class="nav-brand">JoleDev</span>
      <button class="close-btn" onclick={close} aria-label="Close menu" type="button">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <div class="nav-links">
      {#each navItems as item, i}
        <a
          href={item.href}
          class="mobile-link"
          onclick={close}
          style="animation-delay: {(i + 1) * 50}ms"
        >
          <span class="link-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d={icons[item.label] || icons['Blog']}/>
            </svg>
          </span>
          <span class="link-text">{item.label}</span>
          <svg class="link-arrow" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
        </a>
      {/each}
    </div>

    <div class="nav-footer">
      <span class="footer-text">contacto@joledev.com</span>
    </div>
  </nav>
{/if}

<style>
  .hamburger {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 5px;
    padding: 0.75rem;
    min-width: 44px;
    min-height: 44px;
    background: transparent;
    border: none;
    cursor: pointer;
  }

  .hamburger-line {
    display: block;
    width: 22px;
    height: 2px;
    background-color: var(--color-text-primary);
    border-radius: 2px;
    transition: transform 0.3s ease, opacity 0.3s ease;
    transform-origin: center;
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
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    z-index: 40;
    animation: fadeIn 0.2s ease;
  }

  .mobile-nav {
    position: fixed;
    top: 0;
    right: 0;
    width: min(320px, 85vw);
    height: 100dvh;
    background: var(--color-bg-primary);
    border-left: 1px solid var(--color-border);
    display: flex;
    flex-direction: column;
    z-index: 50;
    animation: slideIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    box-shadow: -8px 0 32px rgba(0, 0, 0, 0.2);
  }

  .nav-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--color-border);
  }

  .nav-brand {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    color: var(--color-accent-primary);
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 0.5rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-secondary);
    color: var(--color-text-secondary);
    cursor: pointer;
    transition: all 0.2s;
  }

  .close-btn:hover {
    color: var(--color-text-primary);
    border-color: var(--color-text-muted);
  }

  .nav-links {
    flex: 1;
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    overflow-y: auto;
  }

  .mobile-link {
    display: flex;
    align-items: center;
    gap: 0.875rem;
    padding: 0.875rem 1rem;
    color: var(--color-text-primary);
    text-decoration: none;
    border-radius: 0.75rem;
    transition: background-color 0.15s, transform 0.15s;
    animation: linkSlide 0.35s cubic-bezier(0.16, 1, 0.3, 1) both;
  }

  .mobile-link:hover {
    background-color: var(--color-bg-secondary);
  }

  .mobile-link:active {
    transform: scale(0.98);
    background-color: var(--color-accent-subtle);
  }

  .link-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 0.625rem;
    background: var(--color-bg-secondary);
    color: var(--color-accent-primary);
    flex-shrink: 0;
    border: 1px solid var(--color-border);
  }

  .link-text {
    flex: 1;
    font-size: 0.9375rem;
    font-weight: 500;
  }

  .link-arrow {
    color: var(--color-text-muted);
    opacity: 0.5;
    transition: opacity 0.15s, transform 0.15s;
  }

  .mobile-link:hover .link-arrow {
    opacity: 1;
    transform: translateX(2px);
  }

  .nav-footer {
    padding: 1rem 1.25rem;
    border-top: 1px solid var(--color-border);
    text-align: center;
  }

  .footer-text {
    font-size: 0.75rem;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes slideIn {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }

  @keyframes linkSlide {
    from {
      opacity: 0;
      transform: translateX(16px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .mobile-nav,
    .overlay,
    .mobile-link {
      animation: none;
    }
  }
</style>
