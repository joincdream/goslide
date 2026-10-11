<script>
  import { hmr } from '../stores/hmr.svelte.js';
</script>

{#if hmr.error}
  <!-- Fullscreen Error Overlay -->
  <div
    id="goslide-error-overlay"
    class="fixed top-5 left-1/2 -translate-x-1/2 w-[min(900px,92vw)] z-[999999] bg-zinc-900/95 backdrop-blur-md border border-red-500 rounded-xl shadow-2xl p-5 text-zinc-100 font-sans box-border"
  >
    <div class="flex items-center justify-between mb-3 pb-2.5 border-b border-red-500/30">
      <div class="flex items-center gap-2.5">
        <span class="text-xl">⚠️</span>
        <span class="font-bold text-base text-red-300">{hmr.error.title || 'Build Error'}</span>
        {#if hmr.error.file}
          <span class="text-xs bg-red-500/20 text-red-300 px-2 py-0.5 rounded font-mono">
            {hmr.error.file}
          </span>
        {/if}
      </div>
      <button
        onclick={() => hmr.clearError()}
        class="bg-transparent border-0 text-zinc-400 hover:text-zinc-200 text-2xl cursor-pointer leading-none p-1 rounded"
        title="Dismiss"
      >
        &times;
      </button>
    </div>

    <pre class="m-0 p-3.5 bg-black/60 rounded-md text-red-400 font-mono text-xs leading-relaxed overflow-x-auto max-h-[45vh] whitespace-pre-wrap break-words">{hmr.error.message || 'Unknown error occurred'}</pre>

    <div class="mt-3 text-xs text-zinc-400 flex items-center justify-between">
      <span>💡 Fix the syntax error in your editor and save to resume live preview.</span>
      <span class="text-[11px] text-zinc-500">Goslide Live Resilience Engine</span>
    </div>
  </div>
{/if}

{#if hmr.warnings && hmr.warnings.length > 0}
  <!-- Floating Diagnostic Toast -->
  <div
    id="goslide-warning-toast"
    class="fixed top-5 right-5 max-w-[440px] w-[calc(100vw-40px)] z-[999998] bg-zinc-900/95 backdrop-blur-md border border-amber-500 rounded-lg shadow-xl p-3.5 text-zinc-100 font-sans text-xs leading-normal box-border transition-all duration-300"
  >
    <div class="flex items-center justify-between mb-2">
      <div class="flex items-center gap-1.5 font-bold text-amber-400 text-sm">
        <span>⚠️</span> <span>DSL Warning</span>
      </div>
      <button
        onclick={() => hmr.clearWarnings()}
        class="bg-transparent border-0 text-zinc-400 hover:text-zinc-200 text-lg cursor-pointer leading-none"
        title="Dismiss"
      >
        &times;
      </button>
    </div>

    <div class="max-h-[220px] overflow-y-auto space-y-2">
      {#each hmr.warnings as d}
        <div class="pb-2 border-b border-amber-500/20 last:border-b-0 last:pb-0">
          <div class="font-semibold text-amber-200">
            Slide {d.slideIndex} (Line {d.line})
          </div>
          <div class="text-zinc-300 mt-0.5">
            {d.message}
          </div>
          {#if d.candidates && d.candidates.length > 0}
            <div class="mt-1 text-amber-300 text-[11px]">
              💡 Available: <code>{d.candidates.join(', ')}</code>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </div>
{/if}
