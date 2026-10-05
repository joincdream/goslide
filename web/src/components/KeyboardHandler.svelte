<script>
  import { deck } from '../stores/deck.svelte.js';

  let numberBuffer = '';
  let numberTimeout = null;

  function handleKeydown(e) {
    const activeTag = document.activeElement ? document.activeElement.tagName : '';
    if (['INPUT', 'TEXTAREA', 'SELECT'].includes(activeTag)) return;

    // When blackout/whiteout is active, any key dismissal
    if ((deck.isBlackout || deck.isWhiteout) && !['b', 'B', 'w', 'W'].includes(e.key)) {
      deck.toggleBlackout(false);
      deck.toggleWhiteout(false);
      return;
    }

    // Number jump buffer (0-9)
    if (e.key >= '0' && e.key <= '9' && !deck.isDrawMode) {
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
      if (target >= 1 && target <= deck.totalSlides) {
        deck.goToSlide(target - 1);
      }
      return;
    }

    // Pen controls during drawing mode
    if (deck.isDrawMode) {
      switch (e.key) {
        case '1': deck.setPenColor('#ef4444'); return;
        case '2': deck.setPenColor('#3b82f6'); return;
        case '3': deck.setPenColor('#22c55e'); return;
        case '4': deck.setPenColor('#eab308'); return;
        case '5': deck.setPenColor('#ffffff'); return;
        case '+':
        case '=': deck.adjustPenWidth(1.5); return;
        case '-':
        case '_': deck.adjustPenWidth(-1.5); return;
      }
    }

    // Action & navigation shortcuts
    switch (e.key) {
      case ' ':
      case 'ArrowRight':
      case 'PageDown':
      case 'j':
        e.preventDefault();
        deck.nextSlide();
        break;
      case 'ArrowLeft':
      case 'Backspace':
      case 'PageUp':
      case 'k':
      case 'h':
        e.preventDefault();
        deck.prevSlide();
        break;
      case 'Home':
        e.preventDefault();
        deck.goToSlide(0);
        break;
      case 'End':
        e.preventDefault();
        deck.goToSlide(deck.totalSlides - 1);
        break;
      case 'f':
      case 'F':
        e.preventDefault();
        deck.toggleFullscreen();
        break;
      case 'd':
      case 'D':
        e.preventDefault();
        deck.toggleDrawMode();
        break;
      case 'c':
      case 'C':
        e.preventDefault();
        // Handled via event or direct canvas clear
        window.dispatchEvent(new CustomEvent('goslide:clear-canvas'));
        break;
      case 'l':
      case 'L':
        e.preventDefault();
        deck.toggleLaser();
        break;
      case 's':
      case 'S':
        e.preventDefault();
        deck.toggleSpotlight();
        break;
      case 'b':
      case 'B':
        e.preventDefault();
        deck.toggleBlackout();
        break;
      case 'w':
      case 'W':
        e.preventDefault();
        deck.toggleWhiteout();
        break;
      case 'n':
      case 'N':
        e.preventDefault();
        deck.toggleSidebar();
        break;
      case 'p':
      case 'P':
        e.preventDefault();
        window.dispatchEvent(new CustomEvent('goslide:open-popout'));
        break;
      case 'o':
      case 'O':
        e.preventDefault();
        deck.toggleOverview();
        break;
      case 'Escape':
        e.preventDefault();
        if (deck.isOverviewMode) {
          deck.toggleOverview(false);
        } else if (deck.isSidebarOpen) {
          deck.toggleSidebar(false);
        } else if (deck.isLaserActive) {
          deck.toggleLaser(false);
        } else if (deck.isSpotlightActive) {
          deck.toggleSpotlight(false);
        } else if (deck.isBlackout || deck.isWhiteout) {
          deck.toggleBlackout(false);
          deck.toggleWhiteout(false);
        }
        break;
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />
