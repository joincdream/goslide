<script>
  import { deck } from '../stores/deck.svelte.js';

  let mouseX = $state(0);
  let mouseY = $state(0);

  function handleMouseMove(e) {
    if (deck.isLaserActive || deck.isSpotlightActive) {
      mouseX = e.clientX;
      mouseY = e.clientY;
    }
  }
</script>

<svelte:window onmousemove={handleMouseMove} />

{#if deck.isLaserActive}
  <div
    class="fixed w-3.5 h-3.5 rounded-full pointer-events-none z-[9999] -translate-x-1/2 -translate-y-1/2 transition-[transform] duration-[40ms] ease-out"
    style="left: {mouseX}px; top: {mouseY}px; background-color: #ff0055; box-shadow: 0 0 8px 2px #ff0055, 0 0 16px 4px rgba(255, 0, 85, 0.6);"
  ></div>
{/if}

{#if deck.isSpotlightActive}
  <div
    class="fixed inset-0 pointer-events-none z-[9997]"
    style="background: radial-gradient(circle 140px at {mouseX}px {mouseY}px, transparent 0%, rgba(0, 0, 0, 0.78) 100%);"
  ></div>
{/if}
