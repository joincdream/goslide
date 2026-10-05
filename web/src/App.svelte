<script>
  import { onMount } from 'svelte';
  import { deck } from './stores/deck.svelte.js';
  import DeckViewport from './components/DeckViewport.svelte';
  import DrawingCanvas from './components/DrawingCanvas.svelte';
  import LaserSpotlight from './components/LaserSpotlight.svelte';
  import ScreenMask from './components/ScreenMask.svelte';
  import PresenterSidebar from './components/PresenterSidebar.svelte';
  import PresenterToolbar from './components/PresenterToolbar.svelte';
  import OverviewGrid from './components/OverviewGrid.svelte';
  import KeyboardHandler from './components/KeyboardHandler.svelte';

  let canvasComponent;
  let sidebarComponent;

  onMount(() => {
    const slideElements = document.querySelectorAll('.slide-card');
    deck.init(slideElements);

    // Overview slide click listeners
    slideElements.forEach((slide, idx) => {
      slide.addEventListener('click', () => {
        if (deck.isOverviewMode) {
          deck.goToSlide(idx);
          deck.toggleOverview(false);
        }
      });
    });

    const handleClearCanvas = () => {
      if (canvasComponent) canvasComponent.clearCanvas();
    };

    const handleOpenPopout = () => {
      if (sidebarComponent) sidebarComponent.openPopout();
    };

    window.addEventListener('goslide:clear-canvas', handleClearCanvas);
    window.addEventListener('goslide:open-popout', handleOpenPopout);

    return () => {
      window.removeEventListener('goslide:clear-canvas', handleClearCanvas);
      window.removeEventListener('goslide:open-popout', handleOpenPopout);
    };
  });
</script>

<KeyboardHandler />
<DeckViewport />
<DrawingCanvas bind:this={canvasComponent} />
<LaserSpotlight />
<ScreenMask />
<PresenterSidebar bind:this={sidebarComponent} />
<PresenterToolbar />
<OverviewGrid />
