<script lang="ts">
  import Quoter from './Quoter.svelte';

  interface Props {
    lang: 'es' | 'en';
    apiUrl?: string;
  }

  let { lang, apiUrl = '' }: Props = $props();

  let isOpen = $state(false);

  export function open() {
    isOpen = true;
    document.body.style.overflow = 'hidden';
  }

  function close() {
    isOpen = false;
    document.body.style.overflow = '';
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }

  // Expose open method via custom event
  $effect(() => {
    const handler = () => open();
    window.addEventListener('open-quoter-modal', handler);
    return () => window.removeEventListener('open-quoter-modal', handler);
  });
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="modal-backdrop" onclick={handleBackdropClick} onkeydown={() => {}}>
    <div class="modal-content" role="dialog" aria-modal="true" aria-label={lang === 'es' ? 'Cotizador' : 'Quoter'}>
      <button class="close-btn" onclick={close} type="button" aria-label={lang === 'es' ? 'Cerrar' : 'Close'}>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
      <Quoter {lang} {apiUrl} />
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    padding: 1rem;
    overflow-y: auto;
  }

  .modal-content {
    position: relative;
    width: 100%;
    max-width: 860px;
    max-height: 90vh;
    overflow-y: auto;
    animation: modalIn 0.2s ease-out;
  }

  @media (max-width: 640px) {
    .modal-backdrop {
      align-items: flex-start;
      padding: 0;
    }

    .modal-content {
      max-width: 100%;
      max-height: 100vh;
      min-height: 100vh;
      border-radius: 0;
    }
  }

  .close-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    z-index: 10;
    padding: 0.5rem;
    border: none;
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    border-radius: 0.5rem;
    transition: color 0.2s;
  }

  .close-btn:hover {
    color: var(--color-text-primary);
  }

  @keyframes modalIn {
    from { opacity: 0; transform: scale(0.97); }
    to { opacity: 1; transform: scale(1); }
  }
</style>
