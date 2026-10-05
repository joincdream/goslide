<script>
  import { onMount } from 'svelte';
  import { deck } from '../stores/deck.svelte.js';

  let canvasEl;
  let ctx;
  let isDrawing = false;
  let lastX = 0;
  let lastY = 0;

  function resizeCanvas() {
    if (!canvasEl) return;
    const dpr = window.devicePixelRatio || 1;
    if (deck.isLock1080p) {
      canvasEl.width = Math.round(1920 * dpr);
      canvasEl.height = Math.round(1080 * dpr);
      canvasEl.style.width = '1920px';
      canvasEl.style.height = '1080px';
    } else {
      canvasEl.width = Math.round(window.innerWidth * dpr);
      canvasEl.height = Math.round(window.innerHeight * dpr);
      canvasEl.style.width = window.innerWidth + 'px';
      canvasEl.style.height = window.innerHeight + 'px';
    }

    ctx = canvasEl.getContext('2d');
    ctx.scale(dpr, dpr);
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.strokeStyle = deck.penColor;
    ctx.lineWidth = deck.penWidth;
    restoreCanvas();
  }

  export function saveCanvas() {
    if (!ctx || !canvasEl) return;
    deck.drawings.set(deck.currentIndex, ctx.getImageData(0, 0, canvasEl.width, canvasEl.height));
  }

  export function restoreCanvas() {
    if (!ctx || !canvasEl) return;
    ctx.save();
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.clearRect(0, 0, canvasEl.width, canvasEl.height);
    if (deck.drawings.has(deck.currentIndex)) {
      ctx.putImageData(deck.drawings.get(deck.currentIndex), 0, 0);
    }
    ctx.restore();
    ctx.strokeStyle = deck.penColor;
    ctx.lineWidth = deck.penWidth;
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
  }

  export function clearCanvas() {
    deck.drawings.delete(deck.currentIndex);
    if (ctx && canvasEl) {
      ctx.save();
      ctx.setTransform(1, 0, 0, 1, 0, 0);
      ctx.clearRect(0, 0, canvasEl.width, canvasEl.height);
      ctx.restore();
    }
  }

  $effect(() => {
    // Re-restore canvas whenever current slide changes
    const _idx = deck.currentIndex;
    restoreCanvas();
  });

  $effect(() => {
    if (ctx) {
      ctx.strokeStyle = deck.penColor;
      ctx.lineWidth = deck.penWidth;
    }
  });

  onMount(() => {
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);
    return () => window.removeEventListener('resize', resizeCanvas);
  });

  function handleMouseDown(e) {
    if (!deck.isDrawMode || e.button !== 0) return;
    isDrawing = true;
    lastX = e.clientX;
    lastY = e.clientY;
    ctx.beginPath();
    ctx.moveTo(lastX, lastY);
  }

  function handleMouseMove(e) {
    if (!isDrawing || !deck.isDrawMode) return;
    const currX = e.clientX;
    const currY = e.clientY;
    const midX = (lastX + currX) / 2;
    const midY = (lastY + currY) / 2;
    ctx.quadraticCurveTo(lastX, lastY, midX, midY);
    ctx.stroke();
    lastX = currX;
    lastY = currY;
  }

  function handleMouseUp() {
    if (isDrawing && deck.isDrawMode) {
      ctx.lineTo(lastX, lastY);
      ctx.stroke();
      saveCanvas();
    }
    isDrawing = false;
  }
</script>

<svelte:window onmouseup={handleMouseUp} />

<canvas
  bind:this={canvasEl}
  id="goslide-canvas"
  class="fixed top-0 left-0 z-[2000] {deck.isDrawMode ? 'pointer-events-auto cursor-crosshair' : 'pointer-events-none'}"
  onmousedown={handleMouseDown}
  onmousemove={handleMouseMove}
></canvas>
