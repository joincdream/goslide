# [GOS-11] M2-4: 고급 인터랙티브 발표자 도구 킷 개발 계획서

> **티켓 번호**: [GOS-11](https://joincdream.atlassian.net/browse/GOS-11)  
> **마일스톤**: 로드맵 2 (Release) / Milestone M2-4  
> **마감일**: 2026-11-06  
> **상태**: 완료 (Done)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [screencast-annotation.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/screencast-annotation.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 외부 무거운 프레임워크(React, Vue 등)나 별도의 중계 웹소켓 서버 없이, **순수 웹 표준 기술(Web Standard Vanilla JS, BroadcastChannel API, HTML5 Canvas, CSS Grid)**만을 활용하여 **프로덕션 레벨의 듀얼 스크린 발표자 콘솔 및 인터랙티브 프레젠테이션 도구 킷**을 개발하는 것입니다.

발표자는 청중 앞 오프라인 발표뿐만 아니라 **스크린캐스트 녹화 및 4K/울트라와이드 단일 대화면 환경**에서도 최적의 발표 경험을 누릴 수 있도록 **하이브리드 발표자 뷰(인-윈도우 사이드바 모드 + 독립 팝업 윈도우 모드)**를 지원합니다. 이를 통해 **현재 슬라이드, 다음 슬라이드 미리보기, 마크다운 발표자 메모, 발표 경과 타이머**를 실시간으로 확인하며 발표를 원활히 통제할 수 있습니다. 또한 시선을 사로잡는 **레이저 포인터(`L`), 스포트라이트(`S`), 화면 블라인드(`B`/`W`), 슬라이드별 영속 판서(`D`/`C`), 그리드 개요 모드(`O`/`ESC`)**를 제공합니다.

### 1.1 사용자에게 제공하는 핵심 가치
1. **스크린캐스트 & 4K 단일 모니터 최적화 (In-Window Sidebar 모드)**:
   - `N` 키로 **현재 창 우측에 슬라이드와 함께 발표자 노트·타이머 사이드바**를 즉시 토글할 수 있습니다.
   - OBS/스크린 레코더에서 화면 전체가 아닌 **슬라이드 영역(16:9)만 캡처 영역으로 지정**해 두면, 단일 모니터에서도 창 전환 스트레스 없이 완벽한 스크린캐스트 녹화가 가능합니다.
2. **무설치·무지연 듀얼 스크린 오프라인 발표 (Pop-out Window 모드)**:
   - `P` 키(또는 사이드바의 `[ ↗ 창 분리 ]` 버튼) 입력 즉시 새 창으로 발표자 전용 콘솔이 분리 오픈되며, `BroadcastChannel`을 통해 메인 프로젝터 창과 0ms에 수렴하는 무지연 양방향 동기화(슬라이드 전환, 레이저 포인터 좌표, 타이머)를 제공합니다.
3. **발표 몰입을 위한 시선 유도 도구 (Audience Engagement)**:
   - 레이저 포인터(`L`): 마우스 위치를 따라가는 고광택 붉은 점 애니메이션.
   - 스포트라이트(`S`): 마우스 주변 원형 영역만 밝히고 나머지 화면을 암전 마스킹 처리하여 핵심 지점에 시선 집중.
   - 화면 암전/백색화(`B`/`W`): Q&A나 화두 전환 시 즉시 화면을 닫아 청중이 발표자에게 집중하도록 유도.
4. **슬라이드별 필기 상태 보존 (Persistent In-Slide Annotation)**:
   - 기존의 일회성 전체 캔버스를 개선하여, 슬라이드별로 독립된 드로잉 버퍼(ImageData/Path)를 캐싱함으로써 앞뒤 슬라이드를 오가더라도 필기 내용이 사라지지 않고 유지됩니다.
5. **전체 슬라이드 탐색을 위한 그리드 개요 모드 (`O` 또는 `ESC`)**:
   - 수십 장의 슬라이드를 한눈에 조망할 수 있는 반응형 썸네일 그리드 뷰를 토글하여 원하는 슬라이드로 즉시 점프 가능.

---

## 2. 결과물의 구체적 모습 (Concrete Manifestation)

### 2.1 단축키 인터랙션 명세 (Keymap Matrix)
| 단축키 | 기능 명칭 | 설명 및 동작 방식 |
| :---: | :--- | :--- |
| **`N`** | **Presenter Sidebar** | **[인-윈도우]** 발표자 사이드바 토글 (내부 스위치로 `화면 맞춤(Fit)` vs `1080p 고정(Screencast)` 모드 즉시 전환 지원) |
| **`P`** | **Presenter Window** | **[듀얼 윈도우]** 메인 창과 동기화되는 독립 발표자 콘솔 팝업 윈도우 생성/토글 (프로젝터/오프라인 발표용) |
| **`L`** | **Laser Pointer** | 마우스 커서 위치에 실시간 붉은색 레이저 포인터 도트 토글 |
| **`S`** | **Spotlight** | 마우스 커서 주변 150px 원형 강조 빔 외 배경 딤(75% Black) 처리 |
| **`B`** | **Blackout Screen** | 화면 전체를 검은색으로 암전 (발표자에게 시선 유도, 아무 키나 누르면 해제) |
| **`W`** | **Whiteout Screen** | 화면 전체를 흰색으로 전환 |
| **`D`** | **Draw / Pen Mode** | 슬라이드 위 자유 필기 모드 토글 (마우스 드래그로 필기) |
| **`C`** | **Clear Annotation** | 현재 슬라이드의 필기 레이어 전체 지우기 |
| **`1` ~ `5`** | **Pen Color** | 펜 색상 변경 (1: 레드, 2: 블루, 3: 그린, 4: 옐로, 5: 화이트) |
| **`+` / `-`** | **Pen Width** | 펜 굵기 확대 / 축소 (2px ~ 12px) |
| **`O` / `ESC`** | **Grid Overview** | 전체 슬라이드 썸네일 그리드 뷰 모드 토글 |

> 💡 **도구 상호 배타성 (Mutual Exclusivity)**: 시선 집중 및 판서 도구(`L` 레이저 포인터, `D` 펜 판서, `S` 스포트라이트)는 상호 배타적으로 동작합니다. 즉, 하나의 도구를 활성화하면 기존에 켜져 있던 다른 도구는 자동으로 안전하게 꺼집니다.

### 2.2 발표자 도구 UI 레이아웃 구조

#### 1) [인-윈도우] 단일 창 스크린캐스트 레이아웃 (`N` 키)
```
┌───────────────────────────────────────────────┬────────────────────────────────┐
│  MAIN SLIDE VIEWPORT (16:9 녹화 캡처 영역)     │  PRESENTER SIDEBAR (비녹화 영역)│
│                                               │  [ ↗ Pop-out ]  [ 14:32 | 12:45]│
│                                               ├────────────────────────────────┤
│                                               │ [ NEXT SLIDE PREVIEW ]         │
│             # Slide Heading                   │ ┌────────────────────────────┐ │
│                                               │ │ Next Slide Preview...      │ │
│         • Content Point 1                     │ └────────────────────────────┘ │
│         • Content Point 2                     ├────────────────────────────────┤
│                                               │ [ PRESENTER NOTES ]            │
│                                               │ • 이 장표에서 트레이드오프 언급 │
│                                               │ • 2분 이내 진행 권장           │
│                                               ├────────────────────────────────┤
│                                               │ [PROGRESS] [██████░░] Slide 3/10│
└───────────────────────────────────────────────┴────────────────────────────────┘
```

#### 2) [듀얼 윈도우] 독립 OS 팝업 창 레이아웃 (`P` 키 또는 `[ ↗ Pop-out ]` 클릭)
```
┌────────────────────────────────────────────────────────────────────────┐
│  GOSLIDE PRESENTER CONSOLE                         [ 14:32:05 | 12:45 ] │
├───────────────────────────────────┬────────────────────────────────────┤
│ [ CURRENT SLIDE (Slide 3/12) ]    │ [ NEXT SLIDE PREVIEW ]             │
│ ┌───────────────────────────────┐ │ ┌────────────────────────────────┐ │
│ │                               │ │ │                                │ │
│ │       Microservices Core      │ │ │       API Gateway Pattern      │ │
│ │                               │ │ │                                │ │
│ └───────────────────────────────┘ │ └────────────────────────────────┘ │
├───────────────────────────────────┴────────────────────────────────────┤
│ [ PRESENTER NOTES ]                                                    │
│  • 이 장표에서는 모놀리스 분해 시 겪었던 트레이드오프를 강조할 것.      │
│  • 청중에게 네트워크 레이턴시 문제 겪어본 적 있는지 질문 던지기.        │
│  • 시간 배분: 최대 2분 소요 권장                                       │
├────────────────────────────────────────────────────────────────────────┤
│ [PROGRESS] [████████████░░░░░░░░░░] 25% | Elapsed: 08:30 / 30:00 (Pause)│
└────────────────────────────────────────────────────────────────────────┘
```

### 2.3 BroadcastChannel 동기화 프로토콜 명세
- **채널명**: `goslide-channel-${documentId || window.location.pathname}`
- **메시지 패킷 구조**:
```json
{
  "type": "SLIDE_CHANGE | LASER_MOVE | SPOTLIGHT_MOVE | TOOL_TOGGLE | TIMER_CONTROL",
  "sender": "main | presenter",
  "payload": {
    "currentIndex": 2,
    "cursor": { "x": 0.45, "y": 0.62 },
    "tool": "laser",
    "active": true,
    "timer": { "elapsedMs": 512000, "isRunning": true }
  },
  "timestamp": 1728000000000
}
```

---

## 3. 핵심 설계 원칙 및 규칙 (Architectural Rules)

### 3.1 제로 런타임 의존성 (Pure Web Standards)
- 외부 자바스크립트 라이브러리(Node/npm, CDN 등)를 전혀 임포트하지 않고, **100% 브라우저 네이티브 바닐라 자바스크립트**로 구현합니다.
- HTML 파일 하나만으로 발표자 콘솔과 모든 도구가 동작해야 하므로, 발표자 콘솔 템플릿과 스타일 역시 `internal/theme` 내에 Go `embed.FS`로 완결되도록 내장합니다.

### 3.2 브라우저 네이티브 IPC (`BroadcastChannel`)
- 동일 오리진 내 서로 다른 윈도우/탭 간 통신을 위해 `BroadcastChannel`을 활용합니다.
- 서버가 없는 정적 로컬 HTML 환경(`file://` 또는 `goslide serve`)에서도 메인 창과 발표자 창이 무설정으로 즉각 페어링됩니다.
- 한쪽 창이 닫히거나 새로고침되어도 연결 오류 없이 자동 재접속(Graceful Recovery)되도록 설계합니다.

### 3.3 정규화된 좌표계 기반 레이저/스포트라이트 전파
- 창 크기와 해상도가 다른 듀얼 모니터 환경을 고려하여, 레이저 포인터와 스포트라이트 좌표는 절대 픽셀(`clientX/Y`)이 아닌 **슬라이드 기준 0.0 ~ 1.0 비율(Normalized Ratio)**로 전파합니다.

### 3.4 슬라이드별 판서 캐시 아키텍처
- 캔버스 드로잉은 슬라이드 번호(`slideIndex`)를 키로 하는 메모리 맵(`Map<number, ImageData>`)에 저장됩니다.
- 슬라이드 전환 이벤트 발생 시:
  1. 현재 슬라이드의 캔버스 버퍼(`ctx.getImageData`)를 캐시에 저장.
  2. 캔버스를 클리어(`ctx.clearRect`).
  3. 목적지 슬라이드의 캐시가 존재하면 복원(`ctx.putImageData`), 없으면 빈 캔버스 유지.

---

## 4. 세부 개발 태스크 (Work Breakdown Structure)

### 4.1 발표자 메모 파싱 및 HTML DOM 주입 파이프라인
1. **`internal/parser/directive.go` / `model/slide.go`**:
   - 슬라이드 마크다운 내 `<!-- note: ... -->` 지시어 파싱 지원 확인 및 보강.
   - `model.Slide`에 `Notes string` 필드가 온전히 채워지도록 보장.
2. **`internal/renderer/html/template.go`**:
   - 개별 슬라이드 DOM(`<section class="slide-card">`)에 `<aside class="slide-notes">` 또는 `data-notes="..."` 속성 주입.
   - 마크다운 형태의 메모를 HTML로 렌더링하여 안전하게 포함.

### 4.2 발표자 뷰 UI 구현 (인-윈도우 사이드바 & 팝업 콘솔)
1. **인-윈도우 발표자 사이드바 (`#goslide-presenter-sidebar`)**:
   - `N` 키로 슬라이드 화면 우측에 슬라이드 뷰와 함께 배치되는 사이드바 토글.
   - 메인 슬라이드 뷰포트를 16:9 비율로 유지하면서 우측 영역에 메모/타이머/다음 슬라이드 렌더링.
   - 상단 `[ ↗ Pop-out 창 분리 ]` 버튼 제공.
2. **독립 팝업 템플릿 (`assets/templates/presenter.html`)**:
   - `P` 키 또는 `[ ↗ Pop-out ]` 클릭 시 새 창으로 열리는 독립 콘솔 레이아웃.
3. **공통 대시보드 컴포넌트**:
   - 경과 시간 스톱워치 (시작, 일시정지, 리셋).
   - 실시간 시계 (현재 로컬 시각).
   - 슬라이드 진행률 인디케이터.
   - 큰 글씨 발표자 메모 뷰어.

### 4.3 `goslide-core.js` 동기화 엔진 및 듀얼 스크린 통신
1. `BroadcastChannel` 초기화 및 리스너 등록.
2. 메인 창 $\leftrightarrow$ 발표자 창 간 `goToSlide` 양방향 동기화.
3. `P` 키 및 `[ ↗ Pop-out ]` 클릭 시 팝업 창 생성 (`window.open('', 'goslide-presenter', '...')`).
   - 팝업 차단 발생 시 콘솔 안내 및 UI 툴바 버튼 fallback 제공.
4. `N` 키 이벤트 바인딩 및 사이드바 토글/리사이즈 로직.

### 4.4 시선 유도 도구 킷 (Laser Pointer, Spotlight, Black/White Screen)
1. **레이저 포인터 (`L`)**:
   - `#goslide-laser` DOM 엘리먼트 생성 및 CSS 펄스 애니메이션 적용.
   - 마우스 이동에 따른 실시간 위치 업데이트 (`requestAnimationFrame`).
2. **스포트라이트 (`S`)**:
   - CSS Radial Gradient 마스크 또는 반투명 캔버스 오버레이를 활용하여 커서 주변 원형만 밝게 투과.
3. **화면 블라인드 (`B` / `W`)**:
   - 오버레이 레이어 투명도 전환으로 무지연 암전 및 백색화 구현.

### 4.5 슬라이드별 판서 캐싱 고도화
1. `Map<number, ImageData>` 기반 슬라이드 단위 필기 상태 격리.
2. `C` 키 입력 시 현재 슬라이드의 필기 데이터만 선택적 삭제.
3. 펜 색상(`1`~`5`) 및 두께(`+`/`-`) 단축키 및 시각적 피드백 인디케이터 추가.

### 4.6 그리드 개요 모드 (`O` 또는 `ESC`)
1. CSS Grid 기반 뷰포트 전환 클래스(`.overview-mode`) 정의.
2. 모든 슬라이드를 한 화면에 스케일 다운하여 타일 형태로 정렬.
3. 특정 슬라이드 클릭 시 즉시 해당 슬라이드로 이동하고 개요 모드 종료.

---

## 5. 엔지니어링 가드레일 (Hard Constraints)

1. **외부 번들러 배제**: Webpack, Vite 등의 복잡한 빌드 파이프라인 없이, 바닐라 JS 단일 파일(`goslide-core.js`) 내에 모듈화된 패턴으로 구현.
2. **CGO 배제 및 크로스 플랫폼**: 템플릿 내장(`embed.FS`)과 Go 코드는 100% Pure Go 유지.
3. **성능 가드레일**:
   - 마우스 이벤트(`mousemove`)는 `requestAnimationFrame`을 통해 쓰로틀링(Throttling)하여 60fps 부드러운 반응성 유지.
   - 캔버스 리사이즈 및 버퍼 캐싱 시 메모리 릭(Memory Leak) 방지.
4. **AGENTS.md 규칙 6번 준수 (Strict User-Directed Scope)**:
   - 사용자가 명시하지 않은 추가 작업(임의 테스트 실행, 린트 등)은 수행하지 않고 오직 요청된 작업 범위에 집중.

---

## 6. 완료 기준 (Definition of Done)

- [x] `N` 키로 스크린캐스트용 인-윈도우 발표자 사이드바가 정상 토글됨 (슬라이드 16:9 비율 유지).
- [x] `P` 키 및 `[ ↗ Pop-out ]` 클릭 시 독립된 발표자 콘솔 창이 분리되고, 양방향 슬라이드 전환이 즉각 동기화됨.
- [x] 발표자 도구(사이드바/팝업 콘솔)에 현재 슬라이드, 다음 슬라이드, 발표자 메모, 경과 타이머가 정상 출력됨.
- [x] `L` 키(레이저 포인터), `S` 키(스포트라이트), `B`/`W` 키(화면 암전)가 단축키로 토글됨.
- [x] `D` 키(판서) 모드에서 슬라이드별 필기 내용이 페이지 이동 후에도 캐시에서 복원됨.
- [x] `O` 또는 `ESC` 키로 전체 슬라이드 그리드 개요 화면이 전환되고 클릭으로 이동 가능함.
- [x] 모든 기능이 외부 라이브러리 없이 순수 바닐라 JS/HTML5로 동작함.

---

## 7. 단계별 실행 계획 (Step-by-Step Execution Plan)

1. **발표자 메모 렌더링 파이프라인 구축**: `internal/parser`, `internal/renderer/html`의 발표자 메모 데이터 바인딩.
2. **발표자 콘솔 템플릿 구현**: `internal/theme/assets/templates/presenter.html` 및 레이아웃/스타일 구성.
3. **`goslide-core.js` 코어 통신 계층 개발**: `BroadcastChannel` 기반 메인-발표자 간 이벤트 송수신 엔진.
4. **인터랙티브 발표 도구 추가**: 레이저 포인터(`L`), 스포트라이트(`S`), 화면 블라인드(`B`/`W`).
5. **인-슬라이드 판서 캐시 고도화**: 슬라이드별 ImageData 영속 캐싱 및 단축키 바인딩.
6. **그리드 개요 모드 구현**: CSS Grid 썸네일 탐색 및 점프 인터랙션 완성.
7. **최종 점검 및 사용자 보고**: 듀얼 스크린 인터랙션 검증 가이드 안내.
