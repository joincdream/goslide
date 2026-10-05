<script>
  import { deck } from '../stores/deck.svelte.js';

  function openPopout() {
    const popoutWin = window.open('', 'goslide-presenter-' + window.location.pathname, 'width=1100,height=750');
    if (!popoutWin) {
      alert('팝업 차단을 해제해주세요.');
      return;
    }

    const themeStylesEl = document.getElementById('goslide-theme-styles');
    const themeStylesCSS = themeStylesEl ? themeStylesEl.innerHTML : '';

    const channelName = 'goslide-sync-' + window.location.pathname;

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
</script>

<aside
  class="fixed top-0 right-0 w-[480px] max-w-[90vw] h-screen bg-slate-900 text-slate-100 border-l border-slate-700 z-[9990] flex flex-col shadow-2xl transition-transform duration-250 ease-out select-none {deck.isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}"
>
  <!-- Header -->
  <div class="px-4 py-2.5 bg-slate-800 border-b border-slate-700 flex justify-between items-center flex-shrink-0">
    <div class="text-slate-300 flex items-center gap-1.5">
      <span class="text-base" title="Presenter View">🎙️</span>
    </div>
    <div class="flex items-center gap-1.5">
      <button
        onclick={openPopout}
        class="bg-slate-700 hover:bg-sky-600 text-slate-100 p-1.5 rounded-md transition-colors flex items-center justify-center"
        title="새 창으로 분리 (P)"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg>
      </button>
      <button
        onclick={() => deck.toggleSidebar(false)}
        class="bg-slate-700 hover:bg-slate-600 text-slate-200 text-xs px-2 py-1.5 rounded-md font-bold transition-colors"
        title="사이드바 닫기 (N)"
      >
        ✕
      </button>
    </div>
  </div>

  <!-- Mode Selector -->
  <div class="flex bg-slate-950 p-1.5 gap-1.5 border-b border-slate-700 flex-shrink-0">
    <button
      onclick={() => deck.setScreencastMode('fit')}
      class="flex-1 py-1.5 px-3 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 {!deck.isLock1080p ? 'bg-slate-700 text-sky-400 shadow-md' : 'text-slate-400 hover:text-slate-200'}"
      title="가용 화면에 맞게 자동 리사이즈"
    >
      ↔ 화면 맞춤
    </button>
    <button
      onclick={() => deck.setScreencastMode('1080p')}
      class="flex-1 py-1.5 px-3 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 {deck.isLock1080p ? 'bg-slate-700 text-sky-400 shadow-md' : 'text-slate-400 hover:text-slate-200'}"
      title="1920×1080 고정 (스크린캐스트 녹화용)"
    >
      🔒 1080p 고정
    </button>
  </div>

  <!-- Timer Bar (Single Clean Stopwatch) -->
  <div class="px-4 py-3 bg-slate-800/90 border-b border-slate-700 flex flex-col items-center gap-2 flex-shrink-0">
    <span class="font-mono text-3xl font-black text-white tracking-wider tabular-nums whitespace-nowrap drop-shadow text-center">
      {deck.formatTimer()}
    </span>
    <div class="flex items-center justify-center gap-2 w-full">
      <button
        type="button"
        onclick={() => deck.toggleTimer()}
        class="flex-1 max-w-[130px] bg-slate-700 hover:bg-sky-400 hover:text-slate-950 text-slate-100 font-bold py-1.5 text-xs rounded transition-colors text-center"
      >
        {deck.isTimerRunning ? '일시정지' : (deck.timerSeconds > 0 ? '재개' : '시작')}
      </button>
      <button
        type="button"
        onclick={() => deck.resetTimer()}
        class="flex-1 max-w-[90px] bg-slate-700 hover:bg-rose-500 hover:text-white text-slate-300 font-bold py-1.5 text-xs rounded transition-colors text-center"
      >
        리셋
      </button>
    </div>
  </div>

  <!-- Next Slide Preview (448px x 252px 16:9 Card) -->
  <div class="p-4 border-b border-slate-700 flex-shrink-0">
    <div class="text-xs font-bold text-slate-300 uppercase tracking-wider mb-2 flex justify-between items-center">
      <span>Next</span>
      <span class="text-sky-400 font-mono">
        {deck.currentIndex + 1 < deck.totalSlides ? `Slide ${deck.currentIndex + 2}` : 'END'}
      </span>
    </div>
    <div class="w-full h-[252px] bg-black border border-slate-700 rounded-lg overflow-hidden relative flex items-center justify-center shadow-lg">
      {#if deck.currentIndex + 1 < deck.totalSlides}
        <div class="w-[1920px] h-[1080px] absolute top-0 left-0 origin-top-left scale-[0.2333] pointer-events-none p-10 bg-slate-900">
          {@html deck.getNextSlideHTML()}
        </div>
      {:else}
        <div class="text-slate-500 text-sm font-medium">다음 슬라이드가 없습니다.</div>
      {/if}
    </div>
  </div>

  <!-- Speaker Notes with Font Size Adjuster -->
  <div class="flex-1 p-4 flex flex-col min-h-0">
    <div class="flex justify-between items-center mb-2 flex-shrink-0">
      <span class="text-xs font-bold text-slate-300 uppercase tracking-wider">
        Notes
      </span>
      <div class="flex items-center gap-1 bg-slate-800 rounded px-1.5 py-0.5 border border-slate-700">
        <button
          onclick={() => deck.decreaseNotesFontSize()}
          class="px-1.5 py-0.5 text-xs font-mono font-bold text-slate-300 hover:text-white"
          title="메모 글자 축소"
        >
          A-
        </button>
        <span class="text-[11px] font-mono text-sky-400 px-1 font-semibold">{deck.notesFontSize}px</span>
        <button
          onclick={() => deck.increaseNotesFontSize()}
          class="px-1.5 py-0.5 text-xs font-mono font-bold text-slate-300 hover:text-white"
          title="메모 글자 확대"
        >
          A+
        </button>
      </div>
    </div>
    <div
      class="flex-1 overflow-y-auto text-slate-100 bg-slate-950/70 p-4 rounded-lg border border-slate-700 whitespace-pre-wrap select-text tracking-wide shadow-inner"
      style="font-size: {deck.notesFontSize}px; line-height: 1.75;"
    >
      {#if deck.getCurrentNotes()}
        {deck.getCurrentNotes()}
      {:else}
        <span class="text-slate-500 italic">작성된 발표자 메모가 없습니다.</span>
      {/if}
    </div>
  </div>

  <!-- Footer Navigation & Progress -->
  <div class="px-4 py-3 bg-slate-800 border-t border-slate-700 flex flex-col gap-2 flex-shrink-0">
    <div class="flex justify-between items-center">
      <span class="text-sm font-bold text-slate-200">
        Slide {deck.currentIndex + 1} <span class="text-slate-400 font-normal">/ {deck.totalSlides}</span>
      </span>
      <div class="flex items-center gap-1.5">
        <button
          onclick={() => deck.prevSlide()}
          disabled={deck.currentIndex === 0}
          class="bg-slate-700 hover:bg-slate-600 disabled:opacity-40 disabled:hover:bg-slate-700 text-slate-200 text-xs px-2.5 py-1 rounded font-semibold transition-colors"
          title="이전 슬라이드 (PageUp / ←)"
        >
          ← 이전
        </button>
        <button
          onclick={() => deck.nextSlide()}
          disabled={deck.currentIndex >= deck.totalSlides - 1}
          class="bg-sky-600 hover:bg-sky-500 disabled:opacity-40 disabled:hover:bg-sky-600 text-white text-xs px-2.5 py-1 rounded font-semibold transition-colors"
          title="다음 슬라이드 (PageDown / → / Space)"
        >
          다음 →
        </button>
      </div>
    </div>
    <div class="w-full h-1.5 bg-slate-700 rounded-full overflow-hidden">
      <div
        class="h-full bg-sky-400 transition-[width] duration-200"
        style="width: {deck.totalSlides > 0 ? ((deck.currentIndex + 1) / deck.totalSlides) * 100 : 0}%;"
      ></div>
    </div>
  </div>
</aside>
