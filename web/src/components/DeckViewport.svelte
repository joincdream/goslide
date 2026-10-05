<script>
  import { onMount } from 'svelte';
  import { deck } from '../stores/deck.svelte.js';

  let deckEl;

  function resizeDeck() {
    if (!deckEl) {
      deckEl = document.getElementById('goslide-deck');
    }
    if (!deckEl || deck.isOverviewMode) return;

    if (deck.isLock1080p) {
      deckEl.style.transform = 'scale(1)';
      deckEl.style.transformOrigin = 'center center';
    } else {
      const sidebarWidth = deck.isSidebarOpen ? 480 : 0;
      const availableWidth = window.innerWidth - sidebarWidth;
      // Safe margins: 32px each side horizontally (64px total), 24px top/bottom (48px total)
      const scale = Math.min((availableWidth - 64) / 1920, (window.innerHeight - 48) / 1080);
      deckEl.style.transform = `scale(${Math.max(0.1, scale)})`;
      deckEl.style.transformOrigin = 'center center';
    }
  }

  $effect(() => {
    // Reactively re-scale when sidebar toggles or mode changes
    const _open = deck.isSidebarOpen;
    const _lock = deck.isLock1080p;
    const _overview = deck.isOverviewMode;
    resizeDeck();
  });

  onMount(() => {
    resizeDeck();
    window.addEventListener('resize', resizeDeck);
    return () => window.removeEventListener('resize', resizeDeck);
  });
</script>

<!-- Floating Progress Indicator -->
{#if !deck.isOverviewMode}
  <div
    class="fixed bottom-5 px-3 py-1.5 rounded-full font-mono text-xs z-[3000] pointer-events-none select-none transition-all duration-200 {deck.isDrawMode ? 'bg-red-500/90 text-white shadow-lg shadow-red-500/40' : 'bg-slate-900/80 text-slate-400 backdrop-blur border border-slate-700/50'}"
    style="right: {deck.isSidebarOpen ? '504px' : '1.5rem'};"
  >
    {deck.currentIndex + 1} / {deck.totalSlides}
  </div>
{/if}
