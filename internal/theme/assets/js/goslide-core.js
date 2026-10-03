/**
 * Goslide Core Presentation Runtime (Vanilla JS)
 * Handles slide navigation, keyboard shortcuts, fullscreen, and 1:1 transparent canvas drawing.
 */
(function() {
  'use strict';

  // State
  let slides = [];
  let currentIndex = 0;
  let isDrawMode = false;
  let isDrawing = false;
  let lastX = 0;
  let lastY = 0;
  let numberBuffer = '';
  let numberTimeout = null;

  // DOM Elements
  let deck = null;
  let indicator = null;
  let canvas = null;
  let ctx = null;

  function init() {
    deck = document.getElementById('goslide-deck');
    slides = Array.from(document.querySelectorAll('.slide-card'));
    indicator = document.getElementById('goslide-indicator');
    canvas = document.getElementById('goslide-canvas');

    if (canvas) {
      ctx = canvas.getContext('2d');
      setupCanvasAndViewport();
    }

    if (slides.length > 0) {
      const initialIndex = getIndexFromHash();
      goToSlide(initialIndex, false);
    }

    setupKeyboard();
    window.addEventListener('hashchange', () => {
      const hashIndex = getIndexFromHash();
      if (hashIndex !== currentIndex) {
        goToSlide(hashIndex, false);
      }
    });
  }

  function setupCanvasAndViewport() {
    function resize() {
      const dpr = window.devicePixelRatio || 1;

      // 1. High-DPI physical buffer scaling to eliminate aliasing/stair-stepping
      canvas.width = Math.round(window.innerWidth * dpr);
      canvas.height = Math.round(window.innerHeight * dpr);
      canvas.style.width = window.innerWidth + 'px';
      canvas.style.height = window.innerHeight + 'px';

      // 2. Re-apply scale and stroke styles (buffer resizing resets context state)
      ctx.scale(dpr, dpr);
      ctx.lineCap = 'round';
      ctx.lineJoin = 'round';
      ctx.strokeStyle = '#ef4444';
      ctx.lineWidth = 3.5;

      // 3. Responsive 16:9 Deck Scaling
      if (deck) {
        const scale = Math.min(window.innerWidth / 960, window.innerHeight / 540) * 0.96;
        deck.style.transform = 'scale(' + scale + ')';
        deck.style.transformOrigin = 'center center';
      }
    }

    window.addEventListener('resize', resize);
    resize();

    canvas.addEventListener('mousedown', (e) => {
      if (!isDrawMode || e.button !== 0) return;
      isDrawing = true;
      lastX = e.clientX;
      lastY = e.clientY;
      ctx.beginPath();
      ctx.moveTo(lastX, lastY);
    });

    canvas.addEventListener('mousemove', (e) => {
      if (!isDrawing || !isDrawMode) return;
      const currX = e.clientX;
      const currY = e.clientY;

      // Midpoint quadratic bezier interpolation for smooth handwriting stroke
      const midX = (lastX + currX) / 2;
      const midY = (lastY + currY) / 2;

      ctx.quadraticCurveTo(lastX, lastY, midX, midY);
      ctx.stroke();

      lastX = currX;
      lastY = currY;
    });

    window.addEventListener('mouseup', () => {
      if (isDrawing && isDrawMode) {
        ctx.lineTo(lastX, lastY);
        ctx.stroke();
      }
      isDrawing = false;
    });
  }

  function clearCanvas() {
    if (ctx && canvas) {
      ctx.save();
      ctx.setTransform(1, 0, 0, 1, 0, 0);
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.restore();
    }
  }

  function toggleDrawMode() {
    isDrawMode = !isDrawMode;
    if (canvas) {
      canvas.style.pointerEvents = isDrawMode ? 'auto' : 'none';
      canvas.style.cursor = isDrawMode ? 'crosshair' : 'default';
    }
    document.body.classList.toggle('drawing-mode', isDrawMode);
    if (indicator) {
      indicator.classList.toggle('pen-active', isDrawMode);
    }
  }

  function goToSlide(index, updateHash = true) {
    if (slides.length === 0) return;
    const bounded = Math.max(0, Math.min(index, slides.length - 1));

    slides.forEach((slide, idx) => {
      slide.classList.toggle('active', idx === bounded);
    });

    currentIndex = bounded;
    clearCanvas();

    if (indicator) {
      indicator.textContent = (currentIndex + 1) + ' / ' + slides.length;
    }

    if (updateHash) {
      window.location.hash = '#' + (currentIndex + 1);
    }
  }

  function nextSlide() {
    if (currentIndex < slides.length - 1) {
      goToSlide(currentIndex + 1);
    }
  }

  function prevSlide() {
    if (currentIndex > 0) {
      goToSlide(currentIndex - 1);
    }
  }

  function getIndexFromHash() {
    const raw = window.location.hash.replace('#', '');
    const num = parseInt(raw, 10);
    if (!isNaN(num) && num >= 1 && num <= slides.length) {
      return num - 1;
    }
    return 0;
  }

  function toggleFullscreen() {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(() => {});
    } else {
      document.exitFullscreen().catch(() => {});
    }
  }

  function setupKeyboard() {
    window.addEventListener('keydown', (e) => {
      const activeTag = document.activeElement ? document.activeElement.tagName : '';
      if (['INPUT', 'TEXTAREA', 'SELECT'].includes(activeTag)) return;

      // Number jump buffer (0-9)
      if (e.key >= '0' && e.key <= '9') {
        numberBuffer += e.key;
        clearTimeout(numberTimeout);
        numberTimeout = setTimeout(() => { numberBuffer = ''; }, 1500);
        return;
      }

      // Enter key jump
      if (e.key === 'Enter' && numberBuffer !== '') {
        e.preventDefault();
        const target = parseInt(numberBuffer, 10);
        numberBuffer = '';
        clearTimeout(numberTimeout);
        if (target >= 1 && target <= slides.length) {
          goToSlide(target - 1);
        }
        return;
      }

      // Navigation & action shortcuts
      switch (e.key) {
        case ' ':
        case 'ArrowRight':
        case 'PageDown':
        case 'j':
        case 'l':
          e.preventDefault();
          nextSlide();
          break;
        case 'ArrowLeft':
        case 'Backspace':
        case 'PageUp':
        case 'k':
        case 'h':
          e.preventDefault();
          prevSlide();
          break;
        case 'Home':
          e.preventDefault();
          goToSlide(0);
          break;
        case 'End':
          e.preventDefault();
          goToSlide(slides.length - 1);
          break;
        case 'f':
        case 'F':
          e.preventDefault();
          toggleFullscreen();
          break;
        case 'd':
        case 'D':
          e.preventDefault();
          toggleDrawMode();
          break;
        case 'c':
        case 'C':
          e.preventDefault();
          clearCanvas();
          break;
      }
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
