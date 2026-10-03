# [GOS-6] M1-5: HTML 렌더러 및 슬라이드 네비게이션 런타임 개발 계획서

> **티켓 번호**: [GOS-6](https://joincdream.atlassian.net/browse/GOS-6)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-5  
> **마감일**: 2026-10-21  
> **완료일**: 2026-10-03  
> **상태**: 완료 (Done)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [core-architecture.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/core-architecture.md), [screencast-annotation.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/screencast-annotation.md), [decisions-simplification.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 파싱된 Presentation IR(`model.Deck`, `model.Slide`)과 테마 스타일을 결합하여, 웹 브라우저에서 독립적으로 동작하는 **인터랙티브 단일 HTML 슬라이드 문서**를 생성하고 **키보드 네비게이션 및 초경량 스크린캐스트 캔버스 판서 런타임**을 구축하는 것입니다.

임시 프리뷰 렌더링을 넘어, `model.Renderer` 표준 인터페이스를 구현하는 정식 HTML 렌더러와 경량 자바스크립트 엔진(`goslide-core.js`)을 Go 바이너리에 내장(`embed.FS`)하여 **외부 프레임워크나 네트워크 의존 없이 오프라인에서도 즉시 발표 및 녹화가 가능한 완성형 슬라이드**를 제공합니다.

### 1.1 슬라이드 발표자 및 시청자에게 제공하는 가치
1. **단일 HTML 독립 구동 (Zero Dependency Self-Contained)**:
   - 생성된 `.html` 파일 하나만 있으면 브라우저 어디서나 프레젠테이션이 즉시 열리며, 내장 CSS 및 핵심 JS 엔진이 모두 인라인 번들링되어 오프라인 환경에서도 100% 동일하게 동작합니다.
2. **직관적인 슬라이드 네비게이션**:
   - `Space`, `Enter`, `→`, `PageDown` 키로 다음 슬라이드 이동.
   - `Backspace`, `←`, `PageUp` 키로 이전 슬라이드 이동.
   - 숫자 입력 후 `Enter`로 특정 슬라이드 즉시 점프.
   - `Home` / `End`로 처음/마지막 슬라이드 즉시 이동.
   - `F` 키를 이용한 브라우저 전체화면(Fullscreen) 모드 전환.
3. **유튜브 스크린캐스트 특화 1:1 뷰포트 투명 캔버스 판서 ([screencast-annotation.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/screencast-annotation.md))**:
   - **화면 전체 1:1 오버레이**: 캔버스를 슬라이드 내부가 아닌 브라우저 뷰포트 전체(`100vw × 100vh`)에 투명 오버레이로 배치하여 **좌표 보정 연산 0줄, 100% 직관적 판서 UX** 달성.
   - `D` 키로 판서 모드 토글 (커서 크로스헤어 전환 및 마우스 드래그 자유 필기).
   - `C` 키로 현재 슬라이드의 판서 즉시 초기화.
   - 슬라이드 전환 및 창 크기 변경 시 자동으로 판서가 깨끗이 클리어되어 정갈한 화면 녹화 보장.
4. **반응형 16:9 뷰포트 피팅**:
   - 브라우저 창 크기가 달라져도 16:9 슬라이드 비율(기본 960x540 가상 캔버스)을 유지하면서 화면 중앙에 최적 스케일링으로 맞추어 표시.

### 1.2 아키텍처 및 파이프라인 관점의 가치
1. **표준 인터페이스 계약 준수 ([contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md))**:
   - `internal/renderer/html`은 `model.Renderer`(`Render(ctx context.Context, deck *model.Deck, w io.Writer) error`) 계약을 충실히 구현하여 CLI, 웹 서버, 향후 익스포터 간 결합도를 최소화합니다.
2. **단순화 및 실용주의 준수 ([ADR-001](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md))**:
   - 무거운 외부 캔버스/슬라이드 라이브러리를 일체 배제하고 **단 50~80줄의 순수 바닐라 JS**로 모든 인터랙션을 구현합니다.
   - 과도한 Base64 이미지 변환으로 인한 파일 크기 팽창 및 CDP 크래시를 방지하기 위해 마크다운의 상대 경로를 표준 유지하며, `--standalone` 옵션 적용 시에만 로컬 이미지를 Data URI로 변환합니다.

---

## 2. 결과물의 구체적 모습 (Concrete Manifestation)

### 2.1 단일 HTML 슬라이드 DOM 구조
생성되는 HTML 파일은 다음과 같이 2계층 뷰포트 구조와 인라인 스크립트를 포함합니다:

```html
<!DOCTYPE html>
<html lang="ko">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Presentation Title</title>
  <style id="goslide-theme-styles">
    /* base.css + 테마 CSS + 인라인/커스텀 스타일이 합성된 CSS */
  </style>
  <style id="goslide-runtime-styles">
    /* 슬라이드 뷰포트 중앙 배치, active 슬라이드 토글, 1:1 고정 캔버스 스타일 */
  </style>
</head>
<body>
  <!-- Layer 1 (하단): 슬라이드 프레젠테이션 스테이지 -->
  <div id="goslide-stage" class="goslide-stage">
    <div id="goslide-deck" class="goslide-deck">
      <section class="slide-card layout-cover active" data-slide="1" style="...">
        <!-- Slide 1 내용 -->
        <div class="slide-body">...</div>
        <div class="slide-footer"><span>1 / 10</span></div>
      </section>
      <section class="slide-card layout-default" data-slide="2" style="...">
        <!-- Slide 2 내용 -->
      </section>
      ...
    </div>
  </div>

  <!-- Layer 2 (상단): 뷰포트 1:1 투명 캔버스 판서 오버레이 (좌표 오차 0%) -->
  <canvas id="goslide-annotation-canvas" class="goslide-canvas"></canvas>

  <!-- 플로팅 슬라이드 인디케이터 -->
  <div id="goslide-progress-indicator" class="goslide-indicator">1 / 10</div>

  <script id="goslide-runtime-script">
    /* internal/theme/assets/js/goslide-core.js 인라인 임베딩 (50~80줄 바닐라 JS) */
  </script>
</body>
</html>
```

### 2.2 브라우저 키보드 인터랙션 명세
- **슬라이드 탐색**:
  - `→`, `Space`, `Enter`, `PageDown`, `j`, `l`: 다음 슬라이드 표시.
  - `←`, `Backspace`, `PageUp`, `k`, `h`: 이전 슬라이드 표시.
  - `Home`: 첫 슬라이드 이동.
  - `End`: 마지막 슬라이드 이동.
  - 숫자(0~9) 버퍼링 후 `Enter`: 해당 번호의 슬라이드로 즉시 점프 (예: `5` 입력 후 `Enter` 시 5번 슬라이드 이동).
- **화면 및 판서 제어**:
  - `F`: 브라우저 전체화면 토글 (`requestFullscreen()` / `exitFullscreen()`).
  - `D`: 1:1 뷰포트 판서 모드 토글 (캔버스 활성화 `pointer-events: auto`, 커서 펜 모양).
  - `C`: 현재 화면의 캔버스 판서 즉시 삭제 (`clearRect`).
- **자동 정리 UX**:
  - 슬라이드 전환 시 판서 자동 클리어.
  - 창 크기 변경(`resize`) 시 판서 자동 클리어 및 캔버스 해상도 1:1 동기화.
- **URL 해시 동기화**:
  - 슬라이드 이동 시 `window.location.hash = '#1'`과 같이 브라우저 히스토리/해시 연동으로 새로고침 시 위치 복원.

---

## 3. 핵심 설계 원칙 및 규칙 (Architectural Rules)

### 3.1 `model.Renderer` 인터페이스 완전 구현
`internal/renderer/html/renderer.go`는 `internal/model/interfaces.go`의 계약을 엄격히 구현합니다:

```go
type HTMLRenderer struct {
    themeMgr   *theme.Manager
    themeName  string
    themePath  string
    standalone bool
}

func NewRenderer(opts ...Option) *HTMLRenderer
func (r *HTMLRenderer) Render(ctx context.Context, deck *model.Deck, w io.Writer) error
```

### 3.2 패키지 및 에셋 관리 (embed.FS)
- `goslide-core.js`는 `internal/theme/assets/js/` 및 `internal/theme/embed.go`를 통해 바이너리에 완전히 내장합니다.
- `internal/theme/embed.go`에서 `//go:embed assets/css/*.css assets/js/*.js` 형태로 JS 에셋까지 함께 임베딩하도록 확장합니다.
- `internal/theme` 패키지에서 `GetCoreJS() (string, error)` API를 제공하여 렌더러가 깔끔하게 스크립트를 주입할 수 있도록 합니다.

### 3.3 ADR-001 단순화 원칙 적용 (YAGNI/KISS)
- 무거운 외부 라이브러리(Fabric.js, Reveal.js 등)를 배제하고 **순수 바닐라 JS 50~80줄 내외**로 키보드 네비게이션과 1:1 뷰포트 캔버스 판서를 완성합니다.
- 캔버스를 슬라이드 내부에 종속시키지 않고 화면 전체 뷰포트에 1:1 매핑하여 복잡한 좌표 변환 수식을 원천 제거합니다.
- `internal/renderer/html/bundler.go`:
  - 마크다운 로컬 이미지 경로는 웹 표준 상대 경로를 기본 유지하되, `--standalone` 옵션 지정 시에만 Base64 Data URI로 안전하게 치환합니다.

---

## 4. 세부 개발 태스크 (Work Breakdown Structure)

### 4.1 에셋 확장 및 임베딩 (`internal/theme/`)
1. **`internal/theme/assets/js/goslide-core.js` 작성**:
   - 슬라이드 전환 컨트롤러 (`currentSlide`, `totalSlides`, CSS 클래스 `active` 토글).
   - 키보드 이벤트 리스너 (`keydown`): 방향키, 스페이스, PageUp/Down, 숫자 점프 버퍼, F(전체화면).
   - 1:1 뷰포트 투명 캔버스 판서 엔진:
     - `D` 키: 펜 모드 토글 (`pointer-events: auto` / `crosshair`).
     - 마우스 드래그 드로잉: `e.clientX`, `e.clientY` 1:1 직결.
     - `C` 키: `ctx.clearRect` 초기화.
     - 슬라이드 이동 및 `resize` 시 캔버스 자동 리셋.
   - URL Hash 연동 (`#1`, `#2`).
2. **`internal/theme/embed.go` & `theme.go` 확장**:
   - `//go:embed assets/css/*.css assets/js/*.js` 임베딩 범위 확장.
   - `GetCoreJS() (string, error)` 메서드 추가.

### 4.2 HTML 템플릿 및 데이터 모델 (`internal/renderer/html/`)
1. **`internal/renderer/html/template.go`**:
   - `html/template` 기반의 마스터 HTML 템플릿 정의 (`slideDocumentTemplate`).
   - 슬라이드 메타데이터(타이틀, 작성자, 뷰포트 비율 16:9/4:3, 헤더/푸터).
   - 합성된 CSS 및 인라인 JS 스크립트 슬롯.
   - 각 슬라이드 레이아웃별 컨테이너(`<section class="slide-card layout-...">`) 합성.
2. **`internal/renderer/html/renderer.go`**:
   - `HTMLRenderer` 구조체 및 옵션 패턴(`WithTheme(name)`, `WithCustomCSS(path)`, `WithStandalone(bool)`).
   - `model.Renderer` 인터페이스의 `Render(ctx, deck, w)` 구현.
   - Context 취소(`ctx.Err()`) 감지 및 리소스 안전 처리.
3. **`internal/renderer/html/bundler.go`**:
   - 단일 HTML 생성을 위한 보조 번들러 유틸리티.
   - `--standalone` 활성화 시 로컬 이미지 파일 Data URI 변환 지원.

### 4.3 CLI 연동 리팩토링 (`cmd/goslide/build.go`)
1. `cmd/goslide/build.go`에 임시로 들어가 있던 `renderPreviewHTML` 및 임시 템플릿 제거.
2. 정식 `html.NewRenderer()`를 생성하여 `p.Parse()`된 `deck`을 `renderer.Render()`로 스트리밍 출력하도록 연결.

### 4.4 단위 및 골든 테스트 작성 (`internal/renderer/html/`)
1. **`internal/renderer/html/renderer_test.go`**:
   - `model.Deck` 샘플 데이터를 사용한 HTML 렌더링 정상 검증.
   - 기본 테마, 클린 테마, 다크 테마 적용 검증.
   - 2단 레이아웃(`two-cols`), 커버 레이아웃(`cover`) 정상 DOM 출력 검증.
   - `goslide-core.js` 및 테마 CSS 인라인 포함 검증.
   - Context 취소 시 에러 반환 검증.

---

## 5. 엔지니어링 가드레일 (Hard Constraints)

1. **CGO 배제 (100% Pure Go)**: Go 표준 `html/template`, `embed.FS`, `io`만 사용하여 순수 Go로 구현.
2. **글로벌 상태 전면 금지**: 패키지 전역 변수 배제, `NewRenderer(opts...)` 생성자 패턴 준수.
3. **`panic()` 남용 금지**: 렌더링 실패나 템플릿 실행 오류는 `fmt.Errorf("...: %w", err)`로 명시적 반환.
4. **Context 생명주기 준수**: `Render(ctx, deck, w)` 시작 및 중간 과정에서 `ctx.Err()`를 확인하여 취소 신호 전파.
5. **순환 참조 방지**: `internal/renderer/html`은 `internal/model`과 `internal/theme`만 임포트하며, 상위 패키지나 파서 패키지를 절대 임포트하지 않음.

---

## 6. 완료 기준 (Definition of Done)

- [x] `internal/theme/assets/js/goslide-core.js`가 순수 바닐라 JS로 작성되고 `embed.FS`에 포함됨.
- [x] `internal/renderer/html/renderer.go`가 `model.Renderer` 인터페이스를 완벽히 구현함.
- [x] `internal/renderer/html/template.go`에서 Go `html/template` 기반의 완성형 슬라이드 HTML 합성기가 동작함.
- [x] `cmd/goslide/build.go`의 임시 템플릿 코드가 제거되고 `internal/renderer/html` 정식 렌더러로 교체됨.
- [x] 생성된 HTML을 브라우저에서 열었을 때 키보드 슬라이드 이동, 숫자 점프, 전체화면(F), 1:1 캔버스 판서(D/C)가 매끄럽게 동작함 (High-DPI 레티나 스케일링 및 2차 베지어 곡선 보정 적용).
- [x] 단위 테스트가 작성되고 `go test -v -race ./internal/renderer/html/...`를 100% 통과함.
- [x] `make lint` 및 `make complexity` 기준을 무결하게 만족함.

---

## 7. 단계별 실행 계획 (Step-by-Step Execution Plan)

1. **코어 JS 런타임 작성**: `internal/theme/assets/js/goslide-core.js` 구현 (키보드 네비게이션, 해시 라우팅, 1:1 뷰포트 캔버스 판서).
2. **에셋 파이프라인 확장**: `internal/theme/embed.go` 및 `theme.go`에 JS 에셋 임베딩 및 접근 메서드 추가.
3. **HTML 템플릿 및 렌더러 구현**: `internal/renderer/html/template.go`, `renderer.go`, `bundler.go` 구현.
4. **CLI 파이프라인 연동**: `cmd/goslide/build.go`를 `internal/renderer/html` 기반으로 리팩토링.
5. **단위 테스트 검증**: `renderer_test.go` 작성 및 레이스 컨디션 검증.
6. **실제 슬라이드 빌드 및 브라우저 인터랙션 검증**: 샘플 마크다운 빌드 후 브라우저에서 키 네비게이션 및 판서 기능 실기 검증.
