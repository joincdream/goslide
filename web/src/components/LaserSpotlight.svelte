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
    class="fixed rounded-full pointer-events-none z-[9999] -translate-x-1/2 -translate-y-1/2 transition-[transform,width,height] duration-[40ms] ease-out"
    style="left: {mouseX}px; top: {mouseY}px; width: {deck.currentLaserSize}px; height: {deck.currentLaserSize}px; background-color: {deck.activeColor}; box-shadow: 0 0 {deck.currentLaserGlow}px 2px {deck.activeColor}, 0 0 {deck.currentLaserGlow * 2}px 4px {deck.activeColor}99;"
  ></div>
{/if}

{#if deck.isSpotlightActive}
  <div
    class="fixed inset-0 pointer-events-none z-[9997]"
    style="background: radial-gradient(circle 140px at {mouseX}px {mouseY}px, transparent 0%, rgba(0, 0, 0, 0.78) 100%);"
  ></div>
{/if}
