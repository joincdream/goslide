<script>
  import { deck, COLOR_PRESETS, WIDTH_PRESETS } from '../stores/deck.svelte.js';
  import { t } from '../stores/i18n.svelte.js';

  function handleClear() {
    window.dispatchEvent(new CustomEvent('goslide:clear-canvas'));
  }
</script>

{#if deck.isSidebarOpen && !deck.isOverviewMode}
  <div
    class="fixed bottom-6 z-[9980] flex items-center gap-3 px-4 py-2 bg-slate-900/90 backdrop-blur-md border border-slate-700/80 shadow-2xl rounded-2xl select-none text-slate-200 transition-all duration-200"
    style="left: calc((100vw - 480px) / 2); transform: translateX(-50%);"
  >
    <!-- Tool Toggle Buttons -->
    <div class="flex items-center gap-1.5">
      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={() => deck.toggleDrawMode()}
        class="px-2.5 py-1.5 text-xs font-semibold rounded-lg flex items-center gap-1.5 transition-all {deck.isDrawMode ? 'bg-sky-600 text-white shadow-md shadow-sky-600/30' : 'bg-slate-800 text-slate-300 hover:bg-slate-700 hover:text-white'}"
        title={t('ui.toolbar.pen_tip', '펜 판서 모드 토글 (D)')}
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>
        <span>{t('ui.toolbar.pen', '펜')}</span>
      </button>

      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={() => deck.toggleLaser()}
        class="px-2.5 py-1.5 text-xs font-semibold rounded-lg flex items-center gap-1.5 transition-all {deck.isLaserActive ? 'bg-rose-600 text-white shadow-md shadow-rose-600/30' : 'bg-slate-800 text-slate-300 hover:bg-slate-700 hover:text-white'}"
        title={t('ui.toolbar.pointer_tip', '레이저 포인터 토글 (L)')}
      >
        <span class="inline-block w-2.5 h-2.5 rounded-full bg-current animate-pulse"></span>
        <span>{t('ui.toolbar.pointer', '포인터')}</span>
      </button>

      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={handleClear}
        class="px-2 py-1.5 text-xs font-semibold rounded-lg bg-slate-800 text-slate-400 hover:bg-slate-700 hover:text-rose-400 transition-all flex items-center gap-1"
        title={t('ui.toolbar.clear_tip', '현재 슬라이드 판서 지우기 (C)')}
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
        <span>{t('ui.toolbar.clear', '지우기')}</span>
      </button>
    </div>

    <!-- Divider -->
    <div class="w-[1px] h-6 bg-slate-700/80"></div>

    <!-- 4 Color Presets -->
    <div class="flex items-center gap-2" title={t('ui.toolbar.colors_tip', '색상 선택 (단축키: 1~4)')}>
      {#each COLOR_PRESETS as preset, idx}
        <button
          onmousedown={(e) => e.preventDefault()}
          onclick={() => deck.setPresetColor(preset.hex)}
          class="relative w-6 h-6 rounded-full transition-all flex items-center justify-center {deck.activeColor.toLowerCase() === preset.hex.toLowerCase() ? 'ring-2 ring-white scale-110 shadow-md shadow-black/40' : 'opacity-70 hover:opacity-100 hover:scale-105'}"
          style="background-color: {preset.hex};"
          title="{t('ui.toolbar.color_' + preset.id, preset.label)} [{idx + 1}]"
        >
          {#if deck.activeColor.toLowerCase() === preset.hex.toLowerCase()}
            <span class="w-1.5 h-1.5 rounded-full bg-white shadow-sm"></span>
          {/if}
        </button>
      {/each}
    </div>

    <!-- Divider -->
    <div class="w-[1px] h-6 bg-slate-700/80"></div>

    <!-- 3 Width Presets -->
    <div class="flex items-center bg-slate-800/90 rounded-lg p-0.5 border border-slate-700/60" title={t('ui.toolbar.widths_tip', '굵기 선택 (단축키: - / +)')}>
      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={() => deck.setPresetWidth('thin')}
        class="px-2.5 py-1 text-xs rounded-md font-medium transition-all flex items-center gap-1 {deck.activeWidthPreset === 'thin' ? 'bg-slate-700 text-sky-400 font-bold shadow-sm' : 'text-slate-400 hover:text-slate-200'}"
        title={t('ui.toolbar.width_thin_tip', '얇게 (펜 2px / 레이저 8px)')}
      >
        <span class="inline-block w-2.5 h-[2px] bg-current rounded-full"></span>
        <span>{t('ui.toolbar.width_thin', '얇게')}</span>
      </button>
      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={() => deck.setPresetWidth('medium')}
        class="px-2.5 py-1 text-xs rounded-md font-medium transition-all flex items-center gap-1 {deck.activeWidthPreset === 'medium' ? 'bg-slate-700 text-sky-400 font-bold shadow-sm' : 'text-slate-400 hover:text-slate-200'}"
        title={t('ui.toolbar.width_medium_tip', '보통 (펜 4px / 레이저 14px)')}
      >
        <span class="inline-block w-2.5 h-[3px] bg-current rounded-full"></span>
        <span>{t('ui.toolbar.width_medium', '보통')}</span>
      </button>
      <button
        onmousedown={(e) => e.preventDefault()}
        onclick={() => deck.setPresetWidth('thick')}
        class="px-2.5 py-1 text-xs rounded-md font-medium transition-all flex items-center gap-1 {deck.activeWidthPreset === 'thick' ? 'bg-slate-700 text-sky-400 font-bold shadow-sm' : 'text-slate-400 hover:text-slate-200'}"
        title={t('ui.toolbar.width_thick_tip', '굵게 (펜 8px / 레이저 22px)')}
      >
        <span class="inline-block w-2.5 h-[5px] bg-current rounded-full"></span>
        <span>{t('ui.toolbar.width_thick', '굵게')}</span>
      </button>
    </div>
  </div>
{/if}
