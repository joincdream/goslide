/**
 * Goslide Deck Reactive State Store (Svelte 5 Runes)
 * Centralizes all presentation state, view modes, interactive tools, and sync.
 */

export const COLOR_PRESETS = [
  { id: 'red', hex: '#ef4444', label: '레드' },
  { id: 'blue', hex: '#3b82f6', label: '블루' },
  { id: 'green', hex: '#22c55e', label: '그린' },
  { id: 'yellow', hex: '#eab308', label: '옐로' }
];

export const WIDTH_PRESETS = {
  thin: { pen: 2, laser: 8, glow: 6, label: '얇게' },
  medium: { pen: 4, laser: 14, glow: 10, label: '보통' },
  thick: { pen: 8, laser: 22, glow: 16, label: '굵게' }
};

export class DeckStore {
  currentIndex = $state(0);
  totalSlides = $state(0);
  slides = $state([]);

  isSidebarOpen = $state(false);
  isLock1080p = $state(false);
  isOverviewMode = $state(false);

  // Tools (Mutually exclusive: Laser, Spotlight, Draw)
  isLaserActive = $state(false);
  laserPos = $state({ x: 0, y: 0 });
  isSpotlightActive = $state(false);
  spotlightPos = $state({ x: 0, y: 0 });

  isBlackout = $state(false);
  isWhiteout = $state(false);

  isDrawMode = $state(false);
  penColor = $state('#ef4444');
  penWidth = $state(4);

  // Active Tool Color & Width Presets (Drawing & Laser Pointer)
  activeColor = $state('#ef4444');
  activeWidthPreset = $state('medium'); // 'thin' | 'medium' | 'thick'

  get currentPenWidth() {
    return (WIDTH_PRESETS[this.activeWidthPreset] || WIDTH_PRESETS.medium).pen;
  }

  get currentLaserSize() {
    return (WIDTH_PRESETS[this.activeWidthPreset] || WIDTH_PRESETS.medium).laser;
  }

  get currentLaserGlow() {
    return (WIDTH_PRESETS[this.activeWidthPreset] || WIDTH_PRESETS.medium).glow;
  }

  // Presenter Timer & Clock
  timerSeconds = $state(0);
  isTimerRunning = $state(false);
  timerInterval = null;
  currentTime = $state('');

  // Speaker notes typography size (default: 18px)
  notesFontSize = $state(18);

  // Per-slide drawing cache (slideIndex -> ImageData)
  drawings = new Map();

  channel = null;

  constructor() {
    this.updateClock();
    setInterval(() => this.updateClock(), 1000);
    this.initSyncChannel();
  }

  updateClock() {
    this.currentTime = new Date().toTimeString().split(' ')[0];
  }

  init(slideElements) {
    this.slides = Array.from(slideElements);
    this.totalSlides = this.slides.length;

    // Read initial index from URL hash
    const hashIdx = this.getIndexFromHash();
    this.goToSlide(hashIdx, false);

    window.addEventListener('hashchange', () => {
      const idx = this.getIndexFromHash();
      if (idx !== this.currentIndex) {
        this.goToSlide(idx, false);
      }
    });
  }

  getIndexFromHash() {
    const raw = window.location.hash.replace('#', '');
    const num = parseInt(raw, 10);
    if (!isNaN(num) && num >= 1 && num <= this.totalSlides) {
      return num - 1;
    }
    return 0;
  }

  goToSlide(index, updateHash = true) {
    if (this.totalSlides === 0) return;
    const target = Math.max(0, Math.min(index, this.totalSlides - 1));
    this.currentIndex = target;

    // Toggle active class on DOM slide elements
    this.slides.forEach((slide, idx) => {
      slide.classList.toggle('active', idx === target);
    });

    if (updateHash) {
      window.location.hash = '#' + (target + 1);
    }

    this.applyAutofit(this.slides[target]);
    this.broadcastSlideChange();
  }

  nextSlide() {
    if (this.currentIndex < this.totalSlides - 1) {
      this.goToSlide(this.currentIndex + 1);
    }
  }

  prevSlide() {
    if (this.currentIndex > 0) {
      this.goToSlide(this.currentIndex - 1);
    }
  }

  toggleFullscreen() {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(() => {});
    } else {
      document.exitFullscreen().catch(() => {});
    }
  }

  // --- Interactive Tools ---
  toggleLaser(force) {
    this.isLaserActive = typeof force === 'boolean' ? force : !this.isLaserActive;
    if (this.isLaserActive) {
      if (this.isDrawMode) this.toggleDrawMode(false);
      if (this.isSpotlightActive) this.toggleSpotlight(false);
    }
    document.body.classList.toggle('laser-mode', this.isLaserActive);
    this.broadcastToolSettings();
  }

  toggleSpotlight(force) {
    this.isSpotlightActive = typeof force === 'boolean' ? force : !this.isSpotlightActive;
    if (this.isSpotlightActive) {
      if (this.isDrawMode) this.toggleDrawMode(false);
      if (this.isLaserActive) this.toggleLaser(false);
    }
    this.broadcastToolSettings();
  }

  toggleBlackout(force) {
    this.isBlackout = typeof force === 'boolean' ? force : !this.isBlackout;
    if (this.isBlackout && this.isWhiteout) this.toggleWhiteout(false);
  }

  toggleWhiteout(force) {
    this.isWhiteout = typeof force === 'boolean' ? force : !this.isWhiteout;
    if (this.isWhiteout && this.isBlackout) this.toggleBlackout(false);
  }

  toggleDrawMode(force) {
    this.isDrawMode = typeof force === 'boolean' ? force : !this.isDrawMode;
    if (this.isDrawMode) {
      if (this.isLaserActive) this.toggleLaser(false);
      if (this.isSpotlightActive) this.toggleSpotlight(false);
    }
    document.body.classList.toggle('drawing-mode', this.isDrawMode);
    this.broadcastToolSettings();
  }

  setPenColor(color) {
    this.penColor = color;
    this.activeColor = color;
    this.broadcastToolSettings();
  }

  adjustPenWidth(delta) {
    this.penWidth = Math.max(1.5, Math.min(16, this.penWidth + delta));
  }

  setPresetColor(color) {
    this.activeColor = color;
    this.penColor = color;
    this.broadcastToolSettings();
  }

  setPresetWidth(presetKey) {
    if (WIDTH_PRESETS[presetKey]) {
      this.activeWidthPreset = presetKey;
      this.penWidth = WIDTH_PRESETS[presetKey].pen;
      this.broadcastToolSettings();
    }
  }

  cycleWidthPreset(direction = 1) {
    const keys = ['thin', 'medium', 'thick'];
    const curIdx = keys.indexOf(this.activeWidthPreset);
    const nextIdx = Math.max(0, Math.min(keys.length - 1, curIdx + direction));
    this.setPresetWidth(keys[nextIdx]);
  }

  broadcastToolSettings() {
    if (!this.channel) return;
    this.channel.postMessage({
      type: 'TOOL_SETTINGS_SYNC',
      payload: {
        activeColor: this.activeColor,
        activeWidthPreset: this.activeWidthPreset,
        isLaserActive: this.isLaserActive,
        isDrawMode: this.isDrawMode
      }
    });
  }

  toggleOverview(force) {
    this.isOverviewMode = typeof force === 'boolean' ? force : !this.isOverviewMode;
    document.body.classList.toggle('overview-mode', this.isOverviewMode);
    window.dispatchEvent(new Event('resize'));
  }

  // --- Presenter Sidebar & Modes ---
  toggleSidebar(force) {
    this.isSidebarOpen = typeof force === 'boolean' ? force : !this.isSidebarOpen;
    document.body.classList.toggle('sidebar-open', this.isSidebarOpen);
    window.dispatchEvent(new Event('resize'));
  }

  setScreencastMode(mode) {
    this.isLock1080p = (mode === '1080p');
    document.body.classList.toggle('lock-1080p', this.isLock1080p);
    window.dispatchEvent(new Event('resize'));
  }

  // --- Timer Controls ---
  toggleTimer() {
    if (this.isTimerRunning) {
      clearInterval(this.timerInterval);
      this.isTimerRunning = false;
    } else {
      this.isTimerRunning = true;
      this.timerInterval = setInterval(() => {
        this.timerSeconds++;
      }, 1000);
    }
  }

  resetTimer() {
    if (this.timerInterval) {
      clearInterval(this.timerInterval);
      this.timerInterval = null;
    }
    this.isTimerRunning = false;
    this.timerSeconds = 0;
  }

  formatTimer() {
    const hrs = String(Math.floor(this.timerSeconds / 3600)).padStart(2, '0');
    const mins = String(Math.floor((this.timerSeconds % 3600) / 60)).padStart(2, '0');
    const secs = String(this.timerSeconds % 60).padStart(2, '0');
    return `${hrs}:${mins}:${secs}`;
  }

  // --- Speaker Notes Typography ---
  increaseNotesFontSize() {
    this.notesFontSize = Math.min(32, this.notesFontSize + 2);
  }

  decreaseNotesFontSize() {
    this.notesFontSize = Math.max(14, this.notesFontSize - 2);
  }

  resetNotesFontSize() {
    this.notesFontSize = 18;
  }

  // --- Autofit Helper ---
  applyAutofit(slide) {
    if (!slide) return;
    const isAutofit = slide.classList.contains('has-autofit') || slide.dataset.autofit === 'true';
    if (!isAutofit) return;

    const body = slide.querySelector('.slide-body');
    if (!body) return;

    body.style.transform = '';
    body.style.transformOrigin = '';
    body.style.width = '';

    const availHeight = body.clientHeight;
    const contentHeight = body.scrollHeight;

    if (contentHeight > availHeight && availHeight > 0) {
      const scale = Math.max(0.4, (availHeight / contentHeight) * 0.98);
      body.style.transform = `scale(${scale})`;
      const isCentered = slide.classList.contains('cover') ||
                         slide.classList.contains('section') ||
                         slide.classList.contains('layout-cover') ||
                         slide.classList.contains('layout-section');
      body.style.transformOrigin = isCentered ? 'center center' : 'top left';
      body.style.width = `${100 / scale}%`;
    }
  }

  // --- BroadcastChannel Sync ---
  initSyncChannel() {
    try {
      const channelName = 'goslide-sync-' + window.location.pathname;
      this.channel = new BroadcastChannel(channelName);
      this.channel.onmessage = (e) => {
        const msg = e.data;
        if (!msg) return;
        if (msg.type === 'REQUEST_INIT') {
          this.broadcastInit();
        } else if (msg.type === 'NAV_NEXT') {
          this.nextSlide();
        } else if (msg.type === 'NAV_PREV') {
          this.prevSlide();
        } else if (msg.type === 'TOOL_SETTINGS_SYNC') {
          const p = msg.payload;
          if (p.activeColor) {
            this.activeColor = p.activeColor;
            this.penColor = p.activeColor;
          }
          if (p.activeWidthPreset && WIDTH_PRESETS[p.activeWidthPreset]) {
            this.activeWidthPreset = p.activeWidthPreset;
            this.penWidth = WIDTH_PRESETS[p.activeWidthPreset].pen;
          }
          if (typeof p.isLaserActive === 'boolean' && p.isLaserActive !== this.isLaserActive) {
            this.toggleLaser(p.isLaserActive);
          }
          if (typeof p.isDrawMode === 'boolean' && p.isDrawMode !== this.isDrawMode) {
            this.toggleDrawMode(p.isDrawMode);
          }
        } else if (msg.type === 'TOOL_ACTION') {
          if (msg.payload?.action === 'clear-canvas') {
            window.dispatchEvent(new CustomEvent('goslide:clear-canvas'));
          } else if (msg.payload?.action === 'toggle-laser') {
            this.toggleLaser();
          } else if (msg.payload?.action === 'toggle-draw') {
            this.toggleDrawMode();
          }
        }
      };
    } catch (err) {
      console.warn('BroadcastChannel not supported');
    }
  }

  broadcastInit() {
    if (!this.channel) return;
    const slidesData = this.slides.map((slide, idx) => {
      const notesEl = slide.querySelector('.slide-notes');
      return {
        index: idx,
        notes: notesEl ? notesEl.textContent : '',
        html: slide.outerHTML
      };
    });
    this.channel.postMessage({
      type: 'SYNC_INIT',
      payload: {
        currentIndex: this.currentIndex,
        totalSlides: this.totalSlides,
        slidesData: slidesData,
        activeColor: this.activeColor,
        activeWidthPreset: this.activeWidthPreset,
        isLaserActive: this.isLaserActive,
        isDrawMode: this.isDrawMode
      }
    });
  }

  broadcastSlideChange() {
    if (!this.channel) return;
    this.channel.postMessage({
      type: 'SLIDE_CHANGE',
      payload: { index: this.currentIndex }
    });
  }

  getCurrentNotes() {
    const activeSlide = this.slides[this.currentIndex];
    const notesEl = activeSlide ? activeSlide.querySelector('.slide-notes') : null;
    return notesEl ? notesEl.textContent : '';
  }

  getNextSlideHTML() {
    if (this.currentIndex + 1 < this.totalSlides) {
      return this.slides[this.currentIndex + 1].innerHTML;
    }
    return '';
  }
}

export const deck = new DeckStore();
