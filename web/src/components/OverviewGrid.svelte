<script>
  import { deck } from '../stores/deck.svelte.js';

  function handleSelect(index) {
    deck.goToSlide(index);
    deck.toggleOverview(false);
  }
</script>

{#if deck.isOverviewMode}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-[8000] bg-slate-950/95 overflow-y-auto p-10 select-none"
    onclick={() => deck.toggleOverview(false)}
  >
    <div class="max-w-7xl mx-auto grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      {#each deck.slides as slide, idx}
        <div
          class="aspect-video bg-slate-900 border-2 rounded-lg overflow-hidden relative cursor-pointer transition-all duration-200 hover:scale-105 hover:border-sky-400 {idx === deck.currentIndex ? 'border-sky-400 ring-4 ring-sky-400/30' : 'border-slate-700 opacity-80 hover:opacity-100'}"
          onclick={(e) => { e.stopPropagation(); handleSelect(idx); }}
        >
          <div class="absolute top-2 left-2 z-10 bg-slate-800/90 text-slate-300 text-xs px-2 py-0.5 rounded font-mono">
            {idx + 1}
          </div>
          <div class="w-[1920px] h-[1080px] absolute top-0 left-0 origin-top-left scale-[0.16] pointer-events-none p-8">
            {@html slide.innerHTML}
          </div>
        </div>
      {/each}
    </div>
  </div>
{/if}
