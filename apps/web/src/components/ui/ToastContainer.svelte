<script lang="ts">
  import { getToasts, removeToast } from '../../lib/toast.svelte';

  let toasts = $derived(getToasts());
</script>

{#if toasts.length > 0}
  <div class="toast-container" role="status" aria-live="polite">
    {#each toasts as t (t.id)}
      <div class="toast toast-{t.type}">
        <span class="toast-icon">
          {#if t.type === 'success'}
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
          {:else if t.type === 'error'}
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
          {:else}
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
          {/if}
        </span>
        <span class="toast-message">{t.message}</span>
        <button class="toast-close" onclick={() => removeToast(t.id)} aria-label="Close">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    {/each}
  </div>
{/if}

<style>
  .toast-container {
    position: fixed;
    bottom: 1.5rem;
    right: 1.5rem;
    z-index: 9999;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-width: 380px;
    pointer-events: none;
  }

  @media (max-width: 639px) {
    .toast-container {
      left: 1rem;
      right: 1rem;
      bottom: 1rem;
      max-width: none;
    }
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    background: var(--color-glass);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border: 1px solid var(--color-glass-border);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
    pointer-events: auto;
    animation: toast-in 300ms ease-out;
  }

  .toast-success { border-left: 3px solid var(--color-success); }
  .toast-error { border-left: 3px solid var(--color-error); }
  .toast-info { border-left: 3px solid var(--color-accent-primary); }

  .toast-icon {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }
  .toast-success .toast-icon { color: var(--color-success); }
  .toast-error .toast-icon { color: var(--color-error); }
  .toast-info .toast-icon { color: var(--color-accent-primary); }

  .toast-message {
    flex: 1;
    font-size: 0.875rem;
    line-height: 1.4;
    color: var(--color-text-primary);
  }

  .toast-close {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.25rem;
    border: none;
    background: none;
    color: var(--color-text-muted);
    cursor: pointer;
    border-radius: 0.25rem;
    transition: color 150ms ease;
  }
  .toast-close:hover {
    color: var(--color-text-primary);
  }

  @keyframes toast-in {
    from {
      opacity: 0;
      transform: translateY(12px) scale(0.95);
    }
    to {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .toast {
      animation: none;
    }
  }
</style>
