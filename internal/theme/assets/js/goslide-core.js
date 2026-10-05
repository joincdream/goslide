/**
 * Goslide Core Presentation Runtime (Vanilla JS)
 * Handles slide navigation, presenter console, laser pointer, spotlight,
 * persistent canvas drawing, overview grid mode, and BroadcastChannel sync.
 */
(function() {
  'use strict';

  // Presentation State
  let slides = [];
  let currentIndex = 0;
  let isOverviewMode = false;
  let numberBuffer = '';
  let numberTimeout = null;

  // Drawing Canvas State (Persistent per-slide)
  let isDrawMode = false;
  let isDrawing = false;
  let lastX = 0;
  let lastY = 0;
  let currentColor = '#ef4444';
  let currentLineWidth = 3.5;
  const slideDrawings = new Map(); // slideIndex -> ImageData

  // Interactive Tools State
  let isLaserActive = false;
  let isSpotlightActive = false;
  let isBlackout = false;
  let isWhiteout = false;

  // Sidebar & Timer State
  let isSidebarOpen = false;
  let isLock1080p = false;
  let timerSeconds = 0;
  let timerInterval = null;
  let isTimerRunning = false;

  // BroadcastChannel for Pop-out Window Sync
  const channelName = 'goslide-sync-' + window.location.pathname;
  let channel = null;
  try {
    channel = new BroadcastChannel(channelName);
  } catch (err) {
    console.warn('BroadcastChannel not supported in this environment');
  }

  // DOM Elements
  let deck = null;
  let indicator = null;
  let canvas = null;
  let ctx = null;
  let laserEl = null;
  let spotlightEl = null;
  let blackoutEl = null;
  let whiteoutEl = null;
  let sidebarEl = null;

  // Sidebar Controls
  let sidebarTimerDisplay = null;
  let sidebarTimerToggleBtn = null;
  let sidebarTimerResetBtn = null;
  let sidebarClockEl = null;
  let sidebarNextPreview = null;
  let sidebarNextIndex = null;
  let sidebarNotesContent = null;
  let sidebarProgressText = null;
  let sidebarProgressFill = null;

  function init() {
    deck = document.getElementById('goslide-deck');
    slides = Array.from(document.querySelectorAll('.slide-card'));
    indicator = document.getElementById('goslide-indicator');
    canvas = document.getElementById('goslide-canvas');

    laserEl = document.getElementById('goslide-laser');
    spotlightEl = document.getElementById('goslide-spotlight');
    blackoutEl = document.getElementById('goslide-blackout');
    whiteoutEl = document.getElementById('goslide-whiteout');
    sidebarEl = document.getElementById('goslide-presenter-sidebar');

    if (canvas) {
      ctx = canvas.getContext('2d');
      setupCanvasAndViewport();
    }

    setupTools();
    setupSidebar();
    setupSyncChannel();
    setupKeyboard();

    if (slides.length > 0) {
      const initialIndex = getIndexFromHash();
      goToSlide(initialIndex, false);
    }

    window.addEventListener('hashchange', () => {
      const hashIndex = getIndexFromHash();
      if (hashIndex !== currentIndex) {
        goToSlide(hashIndex, false);
      }
    });

    // Overview slide click handler
    slides.forEach((slide, idx) => {
      slide.addEventListener('click', () => {
        if (isOverviewMode) {
          goToSlide(idx);
          toggleOverviewMode(false);
        }
      });
    });
  }

  function setupCanvasAndViewport() {
    function resize() {
      const dpr = window.devicePixelRatio || 1;

      // High-DPI canvas buffer
      if (isLock1080p) {
        canvas.width = Math.round(1920 * dpr);
        canvas.height = Math.round(1080 * dpr);
        canvas.style.width = '1920px';
        canvas.style.height = '1080px';
      } else {
        canvas.width = Math.round(window.innerWidth * dpr);
        canvas.height = Math.round(window.innerHeight * dpr);
        canvas.style.width = window.innerWidth + 'px';
        canvas.style.height = window.innerHeight + 'px';
      }

      ctx.scale(dpr, dpr);
      ctx.lineCap = 'round';
      ctx.lineJoin = 'round';
      ctx.strokeStyle = currentColor;
      ctx.lineWidth = currentLineWidth;

      // Restore drawing buffer if available
      restoreCurrentCanvas();

      // Responsive 16:9 Deck Scaling (Physical Fixed: 1920x1080)
      if (deck && !isOverviewMode) {
        if (isLock1080p) {
          deck.style.transform = 'scale(1)';
          deck.style.transformOrigin = 'center center';
        } else {
          const availableWidth = isSidebarOpen ? window.innerWidth - 400 : window.innerWidth;
          const scale = Math.min(availableWidth / 1920, window.innerHeight / 1080) * 0.96;
          deck.style.transform = 'scale(' + Math.max(0.1, scale) + ')';
          deck.style.transformOrigin = 'center center';
        }
      }

      if (slides[currentIndex]) {
        applyAutofit(slides[currentIndex]);
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
        saveCurrentCanvas();
      }
      isDrawing = false;
    });
  }

  function saveCurrentCanvas() {
    if (ctx && canvas) {
      slideDrawings.set(currentIndex, ctx.getImageData(0, 0, canvas.width, canvas.height));
    }
  }

  function restoreCurrentCanvas() {
    if (!ctx || !canvas) return;
    ctx.save();
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    if (slideDrawings.has(currentIndex)) {
      ctx.putImageData(slideDrawings.get(currentIndex), 0, 0);
    }
    ctx.restore();
    ctx.strokeStyle = currentColor;
    ctx.lineWidth = currentLineWidth;
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
  }

  function clearCanvas() {
    slideDrawings.delete(currentIndex);
    if (ctx && canvas) {
      ctx.save();
      ctx.setTransform(1, 0, 0, 1, 0, 0);
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.restore();
    }
  }

  function toggleDrawMode(force) {
    isDrawMode = typeof force === 'boolean' ? force : !isDrawMode;
    if (isDrawMode) {
      if (isLaserActive) toggleLaser(false);
      if (isSpotlightActive) toggleSpotlight(false);
    }
    if (canvas) {
      canvas.style.pointerEvents = isDrawMode ? 'auto' : 'none';
      canvas.style.cursor = isDrawMode ? 'crosshair' : 'default';
    }
    document.body.classList.toggle('drawing-mode', isDrawMode);
    if (indicator) {
      indicator.classList.toggle('pen-active', isDrawMode);
    }
  }

  function setPenColor(color) {
    currentColor = color;
    if (ctx) ctx.strokeStyle = currentColor;
  }

  function adjustPenWidth(delta) {
    currentLineWidth = Math.max(1.5, Math.min(16, currentLineWidth + delta));
    if (ctx) ctx.lineWidth = currentLineWidth;
  }

  function setupTools() {
    window.addEventListener('mousemove', (e) => {
      if (isLaserActive && laserEl) {
        laserEl.style.left = e.clientX + 'px';
        laserEl.style.top = e.clientY + 'px';
      }
      if (isSpotlightActive && spotlightEl) {
        spotlightEl.style.background = 'radial-gradient(circle 140px at ' + e.clientX + 'px ' + e.clientY + 'px, transparent 0%, rgba(0, 0, 0, 0.78) 100%)';
      }
    });

    if (blackoutEl) {
      blackoutEl.addEventListener('click', () => toggleBlackout(false));
    }
    if (whiteoutEl) {
      whiteoutEl.addEventListener('click', () => toggleWhiteout(false));
    }
  }

  function toggleLaser(force) {
    isLaserActive = typeof force === 'boolean' ? force : !isLaserActive;
    if (isLaserActive) {
      if (isDrawMode) toggleDrawMode(false);
      if (isSpotlightActive) toggleSpotlight(false);
    }
    if (laserEl) laserEl.classList.toggle('active', isLaserActive);
    document.body.classList.toggle('laser-mode', isLaserActive);
  }

  function toggleSpotlight(force) {
    isSpotlightActive = typeof force === 'boolean' ? force : !isSpotlightActive;
    if (isSpotlightActive) {
      if (isDrawMode) toggleDrawMode(false);
      if (isLaserActive) toggleLaser(false);
    }
    if (spotlightEl) spotlightEl.classList.toggle('active', isSpotlightActive);
  }

  function toggleBlackout(force) {
    isBlackout = typeof force === 'boolean' ? force : !isBlackout;
    if (blackoutEl) blackoutEl.classList.toggle('active', isBlackout);
    if (isBlackout && isWhiteout) toggleWhiteout(false);
  }

  function toggleWhiteout(force) {
    isWhiteout = typeof force === 'boolean' ? force : !isWhiteout;
    if (whiteoutEl) whiteoutEl.classList.toggle('active', isWhiteout);
    if (isWhiteout && isBlackout) toggleBlackout(false);
  }

  function toggleOverviewMode(force) {
    isOverviewMode = typeof force === 'boolean' ? force : !isOverviewMode;
    document.body.classList.toggle('overview-mode', isOverviewMode);
    window.dispatchEvent(new Event('resize'));
  }

  function setupSidebar() {
    if (!sidebarEl) return;

    sidebarTimerDisplay = document.getElementById('sidebar-timer-display');
    sidebarTimerToggleBtn = document.getElementById('sidebar-timer-toggle');
    sidebarTimerResetBtn = document.getElementById('sidebar-timer-reset');
    sidebarClockEl = document.getElementById('sidebar-clock');
    sidebarNextPreview = document.getElementById('sidebar-next-preview');
    sidebarNextIndex = document.getElementById('sidebar-next-index');
    sidebarNotesContent = document.getElementById('sidebar-notes-content');
    sidebarProgressText = document.getElementById('sidebar-progress-text');
    sidebarProgressFill = document.getElementById('sidebar-progress-fill');

    function updateClock() {
      if (sidebarClockEl) {
        sidebarClockEl.textContent = new Date().toTimeString().split(' ')[0];
      }
    }
    setInterval(updateClock, 1000);
    updateClock();

    function updateTimer() {
      timerSeconds++;
      const hrs = String(Math.floor(timerSeconds / 3600)).padStart(2, '0');
      const mins = String(Math.floor((timerSeconds % 3600) / 60)).padStart(2, '0');
      const secs = String(timerSeconds % 60).padStart(2, '0');
      if (sidebarTimerDisplay) {
        sidebarTimerDisplay.textContent = hrs + ':' + mins + ':' + secs;
      }
    }

    if (sidebarTimerToggleBtn) {
      sidebarTimerToggleBtn.addEventListener('click', () => {
        if (isTimerRunning) {
          clearInterval(timerInterval);
          isTimerRunning = false;
          sidebarTimerToggleBtn.textContent = '재개';
        } else {
          timerInterval = setInterval(updateTimer, 1000);
          isTimerRunning = true;
          sidebarTimerToggleBtn.textContent = '일시정지';
        }
      });
    }

    if (sidebarTimerResetBtn) {
      sidebarTimerResetBtn.addEventListener('click', () => {
        clearInterval(timerInterval);
        isTimerRunning = false;
        timerSeconds = 0;
        if (sidebarTimerDisplay) sidebarTimerDisplay.textContent = '00:00:00';
        if (sidebarTimerToggleBtn) sidebarTimerToggleBtn.textContent = '시작';
      });
    }

    const popoutBtn = document.getElementById('btn-sidebar-popout');
    if (popoutBtn) {
      popoutBtn.addEventListener('click', openPresenterPopout);
    }

    const modeFitBtn = document.getElementById('btn-mode-fit');
    const mode1080pBtn = document.getElementById('btn-mode-1080p');
    if (modeFitBtn) {
      modeFitBtn.addEventListener('click', () => setScreencastMode('fit'));
    }
    if (mode1080pBtn) {
      mode1080pBtn.addEventListener('click', () => setScreencastMode('1080p'));
    }

    try {
      const savedMode = localStorage.getItem('goslide-screencast-mode') || 'fit';
      setScreencastMode(savedMode, false);
    } catch (e) {}

    const closeBtn = document.getElementById('btn-sidebar-close');
    if (closeBtn) {
      closeBtn.addEventListener('click', () => togglePresenterSidebar(false));
    }
  }

  function setScreencastMode(mode, triggerResize = true) {
    isLock1080p = (mode === '1080p');
    document.body.classList.toggle('lock-1080p', isLock1080p);

    const modeFitBtn = document.getElementById('btn-mode-fit');
    const mode1080pBtn = document.getElementById('btn-mode-1080p');
    if (modeFitBtn) modeFitBtn.classList.toggle('active', !isLock1080p);
    if (mode1080pBtn) mode1080pBtn.classList.toggle('active', isLock1080p);

    try {
      localStorage.setItem('goslide-screencast-mode', mode);
    } catch (e) {}

    if (triggerResize) {
      window.dispatchEvent(new Event('resize'));
    }
  }

  function togglePresenterSidebar(force) {
    isSidebarOpen = typeof force === 'boolean' ? force : !isSidebarOpen;
    document.body.classList.toggle('sidebar-open', isSidebarOpen);
    updateSidebarContent();
    window.dispatchEvent(new Event('resize'));
  }

  function updateSidebarContent() {
    if (!isSidebarOpen || slides.length === 0) return;

    if (sidebarProgressText) {
      sidebarProgressText.textContent = (currentIndex + 1) + ' / ' + slides.length;
    }
    if (sidebarProgressFill) {
      sidebarProgressFill.style.width = (((currentIndex + 1) / slides.length) * 100) + '%';
    }

    // Update speaker notes
    const activeSlide = slides[currentIndex];
    const notesEl = activeSlide ? activeSlide.querySelector('.slide-notes') : null;
    if (sidebarNotesContent) {
      sidebarNotesContent.textContent = notesEl ? notesEl.textContent : '';
    }

    // Update next slide preview
    if (sidebarNextPreview && sidebarNextIndex) {
      if (currentIndex + 1 < slides.length) {
        sidebarNextIndex.textContent = 'Slide ' + (currentIndex + 2);
        const nextSlideClone = slides[currentIndex + 1].cloneNode(true);
        nextSlideClone.classList.add('active');
        sidebarNextPreview.innerHTML = '';
        sidebarNextPreview.appendChild(nextSlideClone);
      } else {
        sidebarNextIndex.textContent = '마지막 장';
        sidebarNextPreview.innerHTML = '<div style="color:#64748b;font-size:0.9rem;">다음 슬라이드가 없습니다.</div>';
      }
    }
  }

  function openPresenterPopout() {
    const popoutWin = window.open('', 'goslide-presenter-' + window.location.pathname, 'width=1100,height=750');
    if (!popoutWin) {
      alert('팝업 차단을 해제해주세요.');
      return;
    }

    const themeStylesEl = document.getElementById('goslide-theme-styles');
    const themeStylesCSS = themeStylesEl ? themeStylesEl.innerHTML : '';

    const presenterHTML = `<!DOCTYPE html>
<html lang="ko">
<head>
  <meta charset="utf-8">
  <title>Goslide Presenter Console</title>
  <style id="goslide-theme-styles">
${themeStylesCSS}
  </style>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { background: #0f172a; color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; height: 100vh; display: flex; flex-direction: column; overflow: hidden; user-select: none; }
    header { background: #1e293b; border-bottom: 1px solid #334155; padding: 10px 20px; display: flex; justify-content: space-between; align-items: center; flex-shrink: 0; }
    .timer-box { font-family: monospace; font-size: 1.2rem; font-weight: bold; }
    main { flex: 1; display: grid; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; gap: 12px; padding: 12px; overflow: hidden; min-height: 0; }
    .card { background: #1e293b; border: 1px solid #334155; border-radius: 8px; display: flex; flex-direction: column; overflow: hidden; min-height: 0; }
    .card-hdr { padding: 6px 12px; font-size: 0.8rem; color: #94a3b8; border-bottom: 1px solid #334155; font-weight: 600; display: flex; justify-content: space-between; flex-shrink: 0; }
    .frame-box { flex: 1; background: #000; overflow: hidden; display: flex; align-items: center; justify-content: center; position: relative; }
    .frame-box .slide-card {
      position: absolute !important;
      top: 50% !important;
      left: 50% !important;
      width: 1920px !important;
      height: 1080px !important;
      display: flex !important;
      opacity: 1 !important;
      box-shadow: none !important;
      border-radius: 0 !important;
      transform-origin: center center !important;
      pointer-events: none !important;
    }
    .notes-card { grid-column: 1 / -1; }
    .notes-box { flex: 1; padding: 14px; overflow-y: auto; font-size: 1.15rem; line-height: 1.6; user-select: text; white-space: pre-wrap; background: #1e293b; color: #e2e8f0; }
    .notes-box:empty::before { content: "작성된 발표자 메모가 없습니다."; color: #64748b; font-style: italic; }
    footer { background: #1e293b; border-top: 1px solid #334155; padding: 8px 20px; display: flex; justify-content: space-between; align-items: center; flex-shrink: 0; }
    .btn { background: #334155; color: #fff; border: 1px solid #475569; padding: 6px 16px; border-radius: 4px; cursor: pointer; font-weight: 600; }
    .btn:hover { background: #0284c7; }
  </style>
</head>
<body>
  <header>
    <div style="font-weight:700;color:#38bdf8;">🖥️ Goslide Presenter Console</div>
    <div class="timer-box" id="p-timer">00:00:00</div>
  </header>
  <main>
    <div class="card">
      <div class="card-hdr"><span>현재 슬라이드</span><span id="p-cur-num">Slide 1</span></div>
      <div class="frame-box" id="p-cur-frame"></div>
    </div>
    <div class="card">
      <div class="card-hdr"><span>다음 슬라이드</span><span id="p-next-num">Slide 2</span></div>
      <div class="frame-box" id="p-next-frame"></div>
    </div>
    <div class="card notes-card">
      <div class="card-hdr"><span>발표자 메모 (Notes)</span></div>
      <div class="notes-box" id="p-notes"></div>
    </div>
  </main>
  <footer>
    <div id="p-counter" style="font-weight:600;color:#94a3b8;">1 / 1</div>
    <div>
      <button class="btn" id="p-prev">이전 (←)</button>
      <button class="btn" id="p-next">다음 (→)</button>
    </div>
  </footer>
  <script>
    const ch = new BroadcastChannel('${channelName}');
    let curIdx = 0, total = 1, sData = [];
    let timerSec = 0;

    setInterval(() => {
      timerSec++;
      const hrs = String(Math.floor(timerSec / 3600)).padStart(2, '0');
      const mins = String(Math.floor((timerSec % 3600) / 60)).padStart(2, '0');
      const secs = String(timerSec % 60).padStart(2, '0');
      document.getElementById('p-timer').textContent = hrs + ':' + mins + ':' + secs;
    }, 1000);

    ch.onmessage = (e) => {
      const m = e.data;
      if (!m) return;
      if (m.type === 'SYNC_INIT') {
        curIdx = m.payload.currentIndex;
        total = m.payload.totalSlides;
        sData = m.payload.slidesData;
        render();
      } else if (m.type === 'SLIDE_CHANGE') {
        curIdx = m.payload.index;
        render();
      }
    };

    function scaleFrames() {
      const curBox = document.getElementById('p-cur-frame');
      const nextBox = document.getElementById('p-next-frame');
      [curBox, nextBox].forEach(box => {
        if (!box) return;
        const slide = box.querySelector('.slide-card');
        if (!slide) return;
        const scale = Math.min((box.clientWidth - 16) / 1920, (box.clientHeight - 16) / 1080);
        slide.style.transform = 'translate(-50%, -50%) scale(' + Math.max(0.1, scale) + ')';
      });
    }

    window.addEventListener('resize', scaleFrames);

    function render() {
      document.getElementById('p-counter').textContent = (curIdx + 1) + ' / ' + total;
      document.getElementById('p-cur-num').textContent = 'Slide ' + (curIdx + 1);
      if (sData[curIdx]) {
        document.getElementById('p-notes').textContent = sData[curIdx].notes;
        document.getElementById('p-cur-frame').innerHTML = sData[curIdx].html;
      }
      if (curIdx + 1 < total) {
        document.getElementById('p-next-num').textContent = 'Slide ' + (curIdx + 2);
        document.getElementById('p-next-frame').innerHTML = sData[curIdx + 1].html;
      } else {
        document.getElementById('p-next-num').textContent = '마지막 장';
        document.getElementById('p-next-frame').innerHTML = '<div style="color:#64748b;font-size:1.1rem;display:flex;height:100%;align-items:center;justify-content:center;">다음 슬라이드가 없습니다.</div>';
      }
      requestAnimationFrame(scaleFrames);
    }

    document.getElementById('p-prev').onclick = () => ch.postMessage({ type: 'NAV_PREV' });
    document.getElementById('p-next').onclick = () => ch.postMessage({ type: 'NAV_NEXT' });
    window.onkeydown = (e) => {
      if (['ArrowRight',' ','PageDown'].includes(e.key)) ch.postMessage({ type: 'NAV_NEXT' });
      if (['ArrowLeft','PageUp'].includes(e.key)) ch.postMessage({ type: 'NAV_PREV' });
    };
    ch.postMessage({ type: 'REQUEST_INIT' });
  <\/script>
</body>
</html>`;

    popoutWin.document.open();
    popoutWin.document.write(presenterHTML);
    popoutWin.document.close();
  }

  function setupSyncChannel() {
    if (!channel) return;

    channel.onmessage = (e) => {
      const msg = e.data;
      if (!msg) return;

      if (msg.type === 'REQUEST_INIT') {
        broadcastInit();
      } else if (msg.type === 'NAV_NEXT') {
        nextSlide();
      } else if (msg.type === 'NAV_PREV') {
        prevSlide();
      }
    };
  }

  function broadcastInit() {
    if (!channel) return;
    const slidesData = slides.map((slide, idx) => {
      const notesEl = slide.querySelector('.slide-notes');
      return {
        index: idx,
        notes: notesEl ? notesEl.textContent : '',
        html: slide.outerHTML
      };
    });

    channel.postMessage({
      type: 'SYNC_INIT',
      payload: {
        currentIndex: currentIndex,
        totalSlides: slides.length,
        slidesData: slidesData
      }
    });
  }

  function broadcastSlideChange() {
    if (!channel) return;
    channel.postMessage({
      type: 'SLIDE_CHANGE',
      payload: { index: currentIndex }
    });
  }

  // Autofit Safety Net: dynamically scale overflowing slides
  function applyAutofit(slide) {
    if (!slide) return;
    const isAutofit = slide.classList.contains('has-autofit') || slide.dataset.autofit === 'true';
    if (!isAutofit) return;

    const body = slide.querySelector('.slide-body');
    if (!body) return;

    // Reset inline styles to measure natural dimensions
    body.style.transform = '';
    body.style.transformOrigin = '';
    body.style.width = '';

    const availHeight = body.clientHeight;
    const contentHeight = body.scrollHeight;

    if (contentHeight > availHeight && availHeight > 0) {
      const scale = Math.max(0.4, (availHeight / contentHeight) * 0.98);
      body.style.transform = 'scale(' + scale + ')';
      const isCentered = slide.classList.contains('cover') ||
                         slide.classList.contains('section') ||
                         slide.classList.contains('layout-cover') ||
                         slide.classList.contains('layout-section');
      body.style.transformOrigin = isCentered ? 'center center' : 'top left';
      body.style.width = (100 / scale) + '%';
    }
  }

  function goToSlide(index, updateHash = true) {
    if (slides.length === 0) return;
    const bounded = Math.max(0, Math.min(index, slides.length - 1));

    slides.forEach((slide, idx) => {
      slide.classList.toggle('active', idx === bounded);
    });

    currentIndex = bounded;
    applyAutofit(slides[bounded]);
    restoreCurrentCanvas();

    if (indicator) {
      indicator.textContent = (currentIndex + 1) + ' / ' + slides.length;
    }

    if (updateHash) {
      window.location.hash = '#' + (currentIndex + 1);
    }

    updateSidebarContent();
    broadcastSlideChange();
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

      // When blackout/whiteout is active, any key dismissal
      if ((isBlackout || isWhiteout) && !['b', 'B', 'w', 'W'].includes(e.key)) {
        toggleBlackout(false);
        toggleWhiteout(false);
        return;
      }

      // Number jump buffer (0-9)
      if (e.key >= '0' && e.key <= '9' && !isDrawMode) {
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

      // Pen controls during drawing mode
      if (isDrawMode) {
        switch (e.key) {
          case '1': setPenColor('#ef4444'); return; // Red
          case '2': setPenColor('#3b82f6'); return; // Blue
          case '3': setPenColor('#22c55e'); return; // Green
          case '4': setPenColor('#eab308'); return; // Yellow
          case '5': setPenColor('#ffffff'); return; // White
          case '+':
          case '=': adjustPenWidth(1.5); return;
          case '-':
          case '_': adjustPenWidth(-1.5); return;
        }
      }

      // Navigation & action shortcuts
      switch (e.key) {
        case ' ':
        case 'ArrowRight':
        case 'PageDown':
        case 'j':
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
        case 'l':
        case 'L':
          e.preventDefault();
          toggleLaser();
          break;
        case 's':
        case 'S':
          e.preventDefault();
          toggleSpotlight();
          break;
        case 'b':
        case 'B':
          e.preventDefault();
          toggleBlackout();
          break;
        case 'w':
        case 'W':
          e.preventDefault();
          toggleWhiteout();
          break;
        case 'n':
        case 'N':
          e.preventDefault();
          togglePresenterSidebar();
          break;
        case 'p':
        case 'P':
          e.preventDefault();
          openPresenterPopout();
          break;
        case 'o':
        case 'O':
          e.preventDefault();
          toggleOverviewMode();
          break;
        case 'Escape':
          e.preventDefault();
          if (isOverviewMode) {
            toggleOverviewMode(false);
          } else if (isSidebarOpen) {
            togglePresenterSidebar(false);
          } else if (isLaserActive) {
            toggleLaser();
          } else if (isSpotlightActive) {
            toggleSpotlight();
          } else if (isBlackout || isWhiteout) {
            toggleBlackout(false);
            toggleWhiteout(false);
          }
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
