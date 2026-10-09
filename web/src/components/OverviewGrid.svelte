<script>
  import { deck } from '../stores/deck.svelte.js';
  import { t } from '../stores/i18n.svelte.js';

  function handleSelect(index) {
    deck.goToSlide(index);
    deck.toggleOverview(false);
  }

  function autoscale(node) {
    function update() {
      const parentWidth = node.parentElement ? node.parentElement.clientWidth : 0;
      if (parentWidth > 0) {
        const scale = parentWidth / 1920;
        node.style.transform = `scale(${scale})`;
      }
    }

    update();
    const observer = new ResizeObserver(update);
    if (node.parentElement) {
      observer.observe(node.parentElement);
    }

    return {
      destroy() {
        observer.disconnect();
      }
    };
  }
</script>

{#if deck.isOverviewMode}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-[8000] bg-slate-950/95 overflow-y-auto p-6 md:p-10 select-none flex flex-col"
    onclick={() => deck.toggleOverview(false)}
  >
    <!-- Overview Header Bar -->
    <div
      class="max-w-7xl w-full mx-auto mb-6 flex justify-between items-center text-slate-300 border-b border-slate-800 pb-3"
      onclick={(e) => e.stopPropagation()}
    >
      <div class="flex items-center gap-2">
        <span class="text-xl">🗂️</span>
        <span class="font-bold text-slate-100 text-lg">{t('ui.overview.title', '슬라이드 개요 (Overview)')}</span>
        <span class="text-xs bg-slate-800 text-slate-400 px-2.5 py-0.5 rounded-full font-mono font-medium">
          {t('ui.overview.slides_count', '%d Slides', deck.totalSlides)}
        </span>
      </div>
      <div class="flex items-center gap-3 text-xs text-slate-400">
        <span>{t('ui.overview.back_hint', '단축키 ESC 또는 O 발표 복귀')}</span>
        <button
          onclick={() => deck.toggleOverview(false)}
          class="p-1 hover:bg-slate-800 hover:text-white rounded-md transition-colors font-bold text-sm"
          title={t('ui.overview.close', '닫기 (ESC)')}
        >
          ✕
        </button>
      </div>
    </div>

    <!-- Overview Grid Tiles -->
    <div class="max-w-7xl w-full mx-auto grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      {#each deck.slides as slide, idx}
        <div
          class="aspect-video bg-slate-900 border-2 rounded-lg overflow-hidden relative cursor-pointer transition-all duration-200 hover:scale-105 hover:border-sky-400 {idx === deck.currentIndex ? 'border-sky-400 ring-4 ring-sky-400/30' : 'border-slate-700 opacity-85 hover:opacity-100'}"
          onclick={(e) => { e.stopPropagation(); handleSelect(idx); }}
        >
          <!-- Slide Number Badge -->
          <div class="absolute top-2 left-2 z-20 bg-slate-900/90 text-slate-200 text-xs px-2 py-0.5 rounded font-mono shadow-md border border-slate-700/50">
            {idx + 1}
          </div>

          <!-- Scaled Slide Container -->
          <div
            use:autoscale
            class="overview-slot w-[1920px] h-[1080px] absolute top-0 left-0 origin-top-left pointer-events-none"
          >
            {@html slide.outerHTML}
          </div>
        </div>
      {/each}
    </div>
  </div>
{/if}

<style>
  /* 
   * Force all slide-cards inside overview slots to be visible, fully opaque, 
   * and correctly positioned regardless of active status 
   */
  .overview-slot :global(.slide-card) {
    position: absolute !important;
    top: 0 !important;
    left: 0 !important;
    width: 1920px !important;
    height: 1080px !important;
    display: flex !important;
    visibility: visible !important;
    opacity: 1 !important;
    box-shadow: none !important;
    border-radius: 0 !important;
    pointer-events: none !important;
    transform-origin: top left !important;
  }
</style>
