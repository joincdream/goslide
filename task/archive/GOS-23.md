# [GOS-23] 팝아웃 발표자 콘솔(Presenter Console) UI 개편 및 시인성 최적화 계획서

> **티켓 번호**: [GOS-23](https://joincdream.atlassian.net/browse/GOS-23)  
> **마일스톤**: 발표자 콘솔(Presenter Console) 시인성 및 대본 가독성 고도화  
> **상태**: 완료 (Done)  
> **담당 패키지**: `web/`, `internal/theme/assets/js/`  
> **연관 소스 파일**:  
> - [`web/src/components/PresenterSidebar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterSidebar.svelte)  
> - [`web/src/stores/deck.svelte.js`](file:///home/yundream/myjob/cloit/Goslide/web/src/stores/deck.svelte.js)  
> - [`internal/theme/assets/js/goslide-core.js`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/js/goslide-core.js) (컴파일 번들 산출물)  

---

## 1. 개요 및 배경 (Overview & Problem Statement)

### 1.1 현상 및 문제점 분석

현재 팝아웃 발표자 콘솔(`'P'` 단축키 실행 화면)은 2×2 그리드 레이아웃을 사용하고 있어 실제 발표 현장에서 다음과 같은 심각한 시인성 및 인체공학적(Ergonomics) 한계를 지님:

```text
+---------------------------+---------------------------+
| Current Slide (1fr)       | Next Slide (1fr)          |  <- Equal size, no hierarchy
+---------------------------+---------------------------+
| Speaker Notes (100% width, 1fr)                       |  <- Excessive width, too shallow
+-------------------------------------------------------+
```

1. **대본(Script) 가독성 파괴 (Line Measure 과다)**:
   * 메모 영역이 가로 1100px 전체로 펼쳐져 있어 한 줄의 글자 수가 100자 이상으로 길어짐.
   * 발표 도중 시선이 좌우로 과도하게 움직여 다음 줄을 놓치는 시선 추적(Eye-tracking) 피로 발생.
2. **세로 높이 부족으로 인한 잦은 스크롤**:
   * 메모 창 세로 높이가 약 250px에 불과해 대본이 3~4줄만 되어도 스크롤이 발생함.
3. **슬라이드 위계(Visual Hierarchy) 부재**:
   * 라이브 청중 화면인 '현재 슬라이드'와 보조 화면인 '다음 슬라이드'가 1:1 같은 크기여서 주목도 분산.
4. **메모 폰트 조절 기능 부재**:
   * 발표 환경(노트북 바로 앞 vs 1.5m 강단 거치대)에 따라 글자 크기를 조절할 수 없음.

---

## 2. 상세 기술 설계 (Technical Design)

### 2.1 2-Column 전문 발표자 레이아웃 (TO-BE)

업계 표준 도구(Keynote, Google Slides)의 설계를 반영하여 **좌측 슬라이드 스택(상하) + 우측 세로형 대본 전용 뷰**로 전면 개편:

```text
+-----------------------------------+-----------------------------------+
| Left: Slides Column (~48%)        | Right: Speaker Notes Column (~52%)|
| +-------------------------------+ | +-------------------------------+ |
| | Current Slide (Live, ~62%)    | | | Speaker Notes Header          | |
| | (Audience Screen Focus)       | | | [A-] [18px] [A+] Font Size    | |
| |                               | | +-------------------------------+ |
| +-------------------------------+ | | Full-height Script View       | |
| +-------------------------------+ | | (Optimal 50-60 chars/line)    | |
| | Next Slide (Preview, ~38%)    | | | (12-18 lines visible at once)| |
| | (Upcoming Slide Context)      | | | (No frequent scrolling)       | |
| +-------------------------------+ | +-------------------------------+ |
+-----------------------------------+-----------------------------------+
| Bottom Toolbar: Pen | Laser | Clear | Colors | Widths | Slide Counter |
+-----------------------------------------------------------------------+
```

### 2.2 레이아웃 그리드 & 플렉스 사양

* **`main` 컨테이너**:
  * `display: grid; grid-template-columns: 48fr 52fr; gap: 14px; padding: 12px 16px;`
* **좌측 슬라이드 스택 (`.slides-col`)**:
  * `display: flex; flex-direction: column; gap: 12px; height: 100%; min-height: 0;`
  * **현재 슬라이드 카드 (`.cur-slide-card`)**: `flex: 62 1 0; min-height: 0;` (라이브 뱃지 추가)
  * **다음 슬라이드 카드 (`.next-slide-card`)**: `flex: 38 1 0; min-height: 0; opacity: 0.9;`
* **우측 대본 카드 (`.notes-card`)**:
  * `height: 100%; min-height: 0; display: flex; flex-direction: column;`
  * **메모 박스 (`.notes-box`)**: `flex: 1; padding: 16px 20px; font-size: 18px; line-height: 1.75;`
* **폰트 컨트롤러**:
  * `[A-]` / `[A+]` 버튼 클릭 시 `14px ~ 32px` 범위 조절.
  * 메인 창(사이드바)과 `BroadcastChannel`을 통한 폰트 크기 동기화.

---

## 3. 코드 수정 상세 명세 (Code Changes Specification)

### 3.1 [`web/src/components/PresenterSidebar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterSidebar.svelte)

#### 1) HTML 마크업 구조 개편
```html
<main>
  <!-- 좌측 컬럼: 슬라이드 스택 -->
  <div class="slides-col">
    <div class="card cur-slide-card">
      <div class="card-hdr">
        <span class="live-title"><span class="live-indicator"></span>현재 슬라이드 (Live)</span>
        <span id="p-cur-num">Slide 1</span>
      </div>
      <div class="frame-box" id="p-cur-frame"></div>
    </div>
    <div class="card next-slide-card">
      <div class="card-hdr">
        <span>다음 슬라이드 (Next)</span>
        <span id="p-next-num">Slide 2</span>
      </div>
      <div class="frame-box" id="p-next-frame"></div>
    </div>
  </div>

  <!-- 우측 컬럼: 풀스크린 발표자 대본 -->
  <div class="card notes-card">
    <div class="card-hdr">
      <span>📝 발표자 대본 (Notes)</span>
      <div class="font-controls">
        <button class="font-btn" id="p-font-dec" title="글자 축소">A-</button>
        <span class="font-val" id="p-font-val">18px</span>
        <button class="font-btn" id="p-font-inc" title="글자 확대">A+</button>
      </div>
    </div>
    <div class="notes-box" id="p-notes"></div>
  </div>
</main>
```

#### 2) CSS 스타일 고도화
* `main`: 2컬럼 그리드(`48fr 52fr`) 정의
* `.slides-col`, `.cur-slide-card`, `.next-slide-card`: 62:38 상하 위계 및 반응형 핏
* `.notes-box`: `font-size: 18px; line-height: 1.75; font-family: sans-serif;`
* `.font-controls`, `.font-btn`: 폰트 크기 조절 UI 스타일링
* `scaleFrames()`: 좌측 카드의 비대칭 높이에 맞춰 각각의 슬라이드 축소 배율 최적화

#### 3) 스크립트 이벤트 및 동기화
* `notesFontSize` 변수 관리 및 버튼 이벤트 리스너 바인딩
* `BroadcastChannel`의 `SYNC_INIT` 메시지에 `notesFontSize` 포함 및 송수신

---

## 4. 단계별 실행 결과 (Implementation Verification)

| 단계 | 작업 내용 | 검증 결과 |
| :---: | :--- | :--- |
| **Step 1** | `PresenterSidebar.svelte` 팝아웃 HTML/CSS/JS 2-Column 리팩토링 및 폰트 조절 구현 | 완료 (좌측 48% 상하 위계 슬라이드 + 우측 52% 세로 대본창) |
| **Step 2** | `deck.svelte.js`의 `getNextSlideHTML()` 및 `SYNC_INIT` 폰트 크기 동기화 | 완료 (`outerHTML`로 테마 보존 및 `notesFontSize` 동기화) |
| **Step 3** | Svelte 5 번들 빌드 및 Go 에셋 동기화 (`npm run build`) | 번들(`goslide-core.js`, `presenter.css`) 정상 생성 확인 |
| **Step 4** | 데모 슬라이드 빌드 및 팝아웃 콘솔 시각 검증 | `testdata/demo.html`에서 'P' 키 팝아웃 2-Column 레이아웃 및 폰트 조절 정상 동작 확인 |
| **Step 5** | Jira 티켓 코멘트 등록 및 상태 완료 전이 | 완료 |
