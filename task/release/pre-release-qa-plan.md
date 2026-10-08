# Goslide v1.0.0 정식 출시 전 품질 보증(QA) 및 E2E 테스트 계획서

> **문서 상태**: 계획 및 설계 (Planning & Specification)  
> **연관 티켓**: [GOS-18](https://joincdream.atlassian.net/browse/GOS-18), [GOS-12](https://joincdream.atlassian.net/browse/GOS-12)  
> **마일스톤**: 로드맵 2 (Release) / Milestone M2-5 (CLI 고도화 및 v1.0.0 공식 출시)  
> **대상 플랫폼**: Linux (amd64, arm64), macOS (Apple Silicon, Intel), Windows (amd64)  
> **최종 갱신일**: 2026-10-06

---

## 1. 개요 및 목적 (Overview & Goals)

본 문서는 Goslide의 공식 정식 버전인 **v1.0.0 출시**에 앞서, 소프트웨어의 안정성, 정확성, 크로스 플랫폼 호환성, 그리고 실전 발표 현장에서의 사용자 경험(UX)을 완벽히 검증하기 위한 **통합 품질 보증(QA) 명세서**입니다.

### 핵심 목표
1. **기능 전수 검증**: 6대 핵심 기능 도메인의 정상 작동 및 예외 케이스 처리 검증.
2. **테스트 사각지대 제거**: 기존 78개 단위 테스트 분석을 기반으로 미흡한 영역(Windows 파일 경로, Public API) 보강.
3. **사용자 여정 중심 시나리오 테스트**: 실제 발표자의 사용 흐름을 반영한 4대 E2E 시나리오 수립.
4. **브라우저 자동화 E2E 구현**: Go 내장 `chromedp`를 활용한 프론트엔드 UI 인터랙션 자동 검증 파이프라인 확립.
5. **3대 OS 실기기 스모크 테스트**: 실제 Microsoft PowerPoint 및 Apple Keynote 호환성 검증.
6. **문서화 및 릴리즈 에셋 무결성**: README.md(영/한), 공식 예제, 릴리즈 바이너리 다운로드 경로 및 CLI 가이드라인 최신화 검증.

---

## 2. Goslide 전체 기능 카탈로그 (Feature Catalog)

| 도메인 | 세부 기능 (Feature) | 설명 및 세부 요구사항 | 검증 방법 |
| :--- | :--- | :--- | :--- |
| **1. 파서 & DSL<br>(Parser)** | Frontmatter 추출 | YAML 메타데이터(`title`, `theme`, `paginate`, `class`) 추출 및 기본값 폴백 | Unit Test |
| | 슬라이드 분할 | `---` 구분자 기반 분할 및 코드 블록 내부 보호 | Unit Test |
| | 확장 지시어 | 전역/로컬 스코프 지시어 (`_layout`, `_bg`, `_color`, `_header` 등) | Unit Test |
| | 2단 레이아웃 | `<!-- split -->` 구분자 기반 좌우 2단 그리드 컬럼 분할 | Unit Test |
| | 구문 강조 (Chroma) | 100+ 프로그래밍 언어 하이라이팅, Mermaid 다이어그램 블록 통과 | Unit Test |
| | 라인별 하이라이트 | `{1,3-5}` 문법 기반 특정 코드 라인 강조 및 흐림(Dim) 처리 | Unit Test |
| | 스텝 애니메이션 | `* item`, `dim-fragments` 지시어 기반 순차 노출(Fragment) | Unit Test |
| | 미디어 & 콜아웃 | `![w:400]` 이미지 크기 조절, `[!NOTE]` 알림 박스, YouTube 임베드 | Unit Test |
| **2. 테마 & 렌더링<br>(Renderer)** | 테마 엔진 | 기본 테마 3종(`default`, `clean`, `dark`) CSS 합성 및 캐스케이딩 | Unit Test / Golden |
| | 바이너리 임베딩 | Go `embed.FS` 기반 프론트엔드 번들(Svelte 5, CSS, JS) 내장 | Build Check |
| | 단독 번들러 | `--standalone` 옵션 적용 시 로컬 이미지를 Base64 data:URI로 인라인 임베딩 | Unit Test / Golden |
| **3. 익스포터<br>(Exporter)** | HTML 렌더러 | 독립 실행형 인터랙티브 웹 슬라이드 생성 | Unit Test / Golden |
| | PDF 익스포터 | Headless Chrome 연동, 16:9 무여백 인쇄, 시스템 Chrome 자동 탐색 | Unit Test / E2E |
| | PPTX 익스포터 | 고해상도 뷰포트 캡처 ➔ OpenXML ZIP 패키징, 발표자 메모 연동 | Unit Test / App QA |
| **4. 개발 서버<br>(Live Server)** | 로컬 웹 서버 | Go `net/http` 기반 초경량 슬라이드 호스팅 | Unit Test |
| | 파일 감시 & 핫리로드 | `fsnotify` 기반 마크다운 변경 감지, SSE 브로드캐스팅, 슬라이드 위치 보존 | Unit Test / Scenario |
| | 경로 보안 서빙 | 디렉토리 경로 탐색(`../` Path Traversal) 차단 | Security Test |
| **5. 발표자 도구<br>(Presenter)** | 듀얼 발표자 콘솔 | `P` 키 팝업, `BroadcastChannel` 기반 실시간 양방향 동기화, 타이머 | Browser E2E |
| | 시선 집중 도구 | `L` 레이저 포인터, `S` 스포트라이트, `B`/`W` 화면 블라인드 | Browser E2E |
| | 드로잉 캔버스 | `D` 키 자유 필기 캔버스, 펜 색상 변경, 슬라이드별 필기 유지 | Browser E2E |
| **6. CLI & 플랫폼<br>(CLI & Runtime)** | 스타터 생성 (`init`) | 표준 스타터 슬라이드 생성 (`--theme`, `--force`) | Unit Test / CLI |
| | 다중 빌드 (`build`) | `-f all` 및 `-f html,pdf,pptx` 쉼표 구분 3종 일괄 컴파일 | Unit Test / CLI |
| | 다국어 지원 (i18n) | CLI 도움말 및 에러 메시지 영/한(`--lang=ko`) 자동 전환 | Unit Test |
| | 크로스 플랫폼 지원 | Linux, macOS, Windows 파일 경로 및 단일 바이너리 실행 | Cross-OS E2E |

---

## 3. 단위 테스트(Unit Test) 인벤토리 및 커버리지 분석

### 3.1 패키지별 단위 테스트 현황 요약
현재 Goslide에는 **총 23개 테스트 파일, 78개 단위 테스트 함수**가 작성되어 있습니다:

| 패키지 경로 | 테스트 파일 수 | 테스트 함수 수 | 커버하는 주요 기능 |
| :--- | :---: | :---: | :--- |
| `cmd/goslide/` | 4 | 15 | CLI 커맨드(`build`, `init`, `serve`), 플래그 파싱, 다중 빌드, 성능 벤치마크 |
| `internal/parser/` | 10 | 32 | Frontmatter, 슬라이드 분할, 지시어, 2단 레이아웃, 코드 하이라이트, 프래그먼트, CRLF |
| `internal/renderer/html/` | 3 | 10 | HTML 골든 파일 비교, 캔버스 런타임, 테마 주입, Base64 이미지 번들링 |
| `internal/exporter/pdf/` | 1 | 3 | 시스템 Chrome 탐색, Headless 인쇄 파이프라인, 에러 핸들링 |
| `internal/exporter/pptx/` | 1 | 6 | OpenXML XML 구조체 생성, 슬라이드/메모 XML 매핑, ZIP 패키징 |
| `internal/server/` | 2 | 7 | 로컬 웹 서버, SSE 이벤트 스트림, 파일 변경 디바운싱 감시 |
| `internal/theme/` | 1 | 7 | Base/Theme CSS 합성, 스타터 템플릿 로드, Svelte 코어 JS 임베딩 |
| `internal/i18n/` | 1 | 2 | 다국어 사전 키 매핑, CLI 메시지 번역 |

### 3.2 테스트 사각지대(Test Gaps) 및 보강 계획
* ⚠️ **GAP-1: Windows 파일 경로 전용 단위 테스트 부재**
  * *현황*: `print.go`, `capture.go`에 `filepath.ToSlash` 기반의 `file:///C:/...` 변환 로직이 추가되었으나, 테스트 코드에서 가상 Windows 드라이브 경로를 명시적으로 검증하지 않음.
  * *보강*: `TestWindowsFileURI_Conversion` 단위 테스트 추가.
* ⚠️ **GAP-2: Public Facade API 단위 테스트 부재**
  * *현황*: `pkg/goslide/goslide.go`에 대한 독립 테스트 파일 없음.
  * *보강*: `pkg/goslide/goslide_test.go` 작성 (`Build`, `Parse` 함수형 옵션 검증).
* ⚠️ **GAP-3: 디렉토리 탐색(Path Traversal) 보안 테스트 부재**
  * *현황*: `internal/server/server.go`에 `filepath.Rel` 방어 로직이 적용됨.
  * *보강*: `../../etc/passwd` 또는 `..\..\Windows` 요청 시 403 Forbidden 반환 검증 단위 테스트 추가.

---

## 4. 시나리오 테스트 명세서 (4대 사용자 여정 E2E)

### 시나리오 1: 신규 발표자 온보딩 및 빠른 슬라이드 작성 (The First-Time Presenter)
* **목적**: 사용자가 Goslide를 처음 설치하여 스타터 마크다운을 생성하고, 실시간 프리뷰를 보며 슬라이드를 완성하는 전체 흐름 검증.
* **테스트 절차**:
  1. 임의의 빈 디렉토리 생성 후 `goslide init presentation.md --theme=clean` 실행.
  2. `presentation.md` 파일이 생성되었는지 확인하고 Frontmatter(`theme: "clean"`) 무결성 점검.
  3. `goslide serve presentation.md --port=8080` 실행.
  4. 로컬 브라우저에서 `http://localhost:8080` 접속하여 커버 슬라이드 노출 확인.
  5. 에디터에서 `presentation.md` 파일 끝에 새로운 슬라이드(`--- \n # New Slide`)를 추가하고 저장.
  6. **합격 판정 기준 (Pass Criteria)**:
     * 새로고침 없이 300ms 이내에 브라우저 화면에 새 슬라이드가 자동 갱신(Hot Reload)되어야 함.
     * 보고 있던 슬라이드 위치(Slide Index)가 첫 페이지로 리셋되지 않아야 함.

---

### 시나리오 2: 복합 마크다운 다중 포맷 일괄 배포 (The Tech Speaker's Batch Build)
* **목적**: 코드 강조, 2단 레이아웃, 표, 발표자 메모 등 Goslide의 모든 확장 문법이 포함된 슬라이드를 단 한 번의 명령으로 3대 포맷(HTML, PDF, PPTX)으로 일괄 제작하는 파이프라인 검증.
* **테스트 대상 파일**: `examples/example-dsl.md`
* **테스트 명령**:
  ```bash
  goslide build examples/example-dsl.md -f all -o dist/ --standalone
  ```
* **합격 판정 기준 (Pass Criteria)**:
  1. **HTML (`dist/example-dsl.html`)**:
     * 단일 독립형 파일로 생성되어야 하며, 이미지 에셋이 Base64로 인라인 임베딩되어 로컬 파일 더블클릭만으로 정상 표시.
  2. **PDF (`dist/example-dsl.pdf`)**:
     * 매직 넘버(`%PDF-`) 유효성 확인.
     * 여백(Margin)이 없는 완벽한 16:9 비율 벡터 PDF 생성.
  3. **PPTX (`dist/example-dsl.pptx`)**:
     * OpenXML ZIP 포맷 규격 충족.
     * 각 슬라이드별 발표자 메모(`<!-- note: -->`)가 PPTX 하단 발표자 노트 영역에 정상 매핑.

---

### 시나리오 3: 실전 프레젠테이션 무대 및 인터랙션 시뮬레이션 (The Live Presentation Run)
* **목적**: 발표 현장에서 프로젝터와 노트북 듀얼 모니터 환경에서 발생하는 단축키 및 인터랙티브 툴킷 검증.
* **합격 판정 기준 (Pass Criteria)**:
  1. **전체 화면 ('F' 키)**: 브라우저가 화면 전체로 전환되며 슬라이드 비율(16:9)이 화면 중앙에 반응형 유지.
  2. **발표자 콘솔 ('P' 키)**:
     * 브라우저 팝업으로 발표자 창이 분리 기동.
     * 메인 창에서 슬라이드를 넘기면 발표자 창의 슬라이드 번호, 경과 시간 타이머, 발표자 메모가 지연 없이 실시간 동기화.
  3. **시선 집중 도구**:
     * `L` 키: 마우스 커서를 따라 발광 레드 레이저 포인터가 부드럽게 추종.
     * `S` 키: 마우스 커서 주변 원형 영역 외 화면이 어둡게 딤 처리(스포트라이트).
     * `B` 키: 화면 전체가 암전(Blackout), 재입력 시 즉시 복구.
  4. **화이트보드 캔버스 ('D' or 'W' 키)**:
     * 슬라이드 위에 마우스 드래그로 자유 필기 가능.
     * 슬라이드를 이전/다음으로 넘겼다가 다시 돌아와도 해당 슬라이드의 필기 데이터가 유지되어야 함.

---

### 시나리오 4: 완전 무의존성 오프라인 자립 환경 (Zero-Dependency Offline Test)
* **목적**: 인터넷이 연결되지 않은 발표장 환경에서 외부 CDN이나 네트워크 호출 없이 슬라이드가 완벽히 작동하는지 검증.
* **테스트 절차**:
  1. Wi-Fi 및 이더넷 연결 차단 (Airplane Mode).
  2. 생성된 `dist/example-dsl.html`을 브라우저에서 열람.
* **합격 판정 기준 (Pass Criteria)**:
  * 브라우저 개발자 도구(Console/Network)에서 외부 URL(CDN, Google Fonts, JS 라이브러리) 요청 실패 에러(ERR_INTERNET_DISCONNECTED)가 단 1건도 없어야 함.
  * 모든 CSS, JS, 웹폰트, 아이콘이 100% 정상 렌더링되어야 함.

---

## 5. Frontend 브라우저 자동화 E2E 테스트 아키텍처

Goslide는 별도의 무거운 외부 도구(Node.js, Playwright, Cypress) 없이, 이미 내장된 **Go 기반 `chromedp` 라이브러리를 활용하여 `go test` 파이프라인 안에서 실제 브라우저 E2E를 자동 실행**할 수 있습니다.

### 5.1 아키텍처 다이어그램
```
[go test -v ./internal/renderer/html/e2e_test.go]
                      │
                      ▼
        [chromedp.NewExecAllocator]
                      │ (시스템 Chrome 헤드리스 기동)
                      ▼
   [HTML 슬라이드 로드: file:///.../slide.html]
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
   [키보드 이벤트]  [DOM 상태 검증] [스크린샷/캔버스]
  - Space/Arrows - active 클래스 - 레이저 SVG 감지
  - 'P' 팝업     - 타이머 카운트 - 드로잉 픽셀 확인
```

### 5.2 브라우저 자동화 테스트 시나리오 명세 (`e2e_test.go`)

```go
// internal/renderer/html/e2e_test.go (예시 설계)
func TestE2E_BrowserSlideInteraction(t *testing.T) {
    // 1. 테스트용 슬라이드 HTML 렌더링
    // 2. Headless Chrome 컨텍스트 할당 (chromedp)
    // 3. 브라우저 내비게이션 검증:
    //    - 'Space' 키 입력 ➔ 2번 슬라이드로 전환되는지 DOM attribute 검증
    // 4. 단축키 검증:
    //    - 'L' 키 입력 ➔ <div class="laser-pointer"> 요소가 화면에 나타나는지 검증
    //    - 'B' 키 입력 ➔ <div class="blackout-mask"> 활성화 확인
    // 5. 드로잉 검증:
    //    - 'W' 키 입력 후 MouseDown ➔ MouseMove ➔ MouseUp 시뮬레이션
    //    - HTML5 Canvas getImageData() 호출로 투명하지 않은 픽셀 존재 확인
}
```

---

## 6. 운영체제별 실기기 E2E 스모크 테스트 체크리스트 (Real-Device QA)

릴리즈 후보(`v1.0.0-rc.1`) 발행 후 실제 하드웨어에서 수행하는 10분 수동 점검표:

| 운영체제 | 점검 항목 | 검증 방법 및 세부 내용 | 판정 (P/F) |
| :--- | :--- | :--- | :---: |
| **Linux** | 정적 바이너리 무결성 | `ldd bin/goslide` 실행 시 동적 링크 라이브러리가 없어야 함 (`statically linked`) | [ ] |
| | 헤드리스 크롬 자동 연동 | 시스템 크롬 감지 후 `goslide build -f pdf` 정상 출력 확인 | [ ] |
| **macOS** | Apple Silicon 네이티브 실행 | `file bin/goslide-darwin-arm64`가 `Mach-O 64-bit arm64` 확인 (Rosetta 불필요) | [ ] |
| | **Apple Keynote 호환성** | 생성된 PPTX를 **Keynote** 앱으로 열었을 때 슬라이드 깨짐 및 폰트 레이아웃 육안 점검 | [ ] |
| | macOS 전용 단축키 점검 | `Cmd + F`(전체화면), `Cmd + P`(인쇄창 방지 및 발표자 콘솔 연동) 충돌 여부 확인 | [ ] |
| **Windows**| **MS PowerPoint 365 호환성** | 생성된 PPTX를 **실제 Windows MS PowerPoint**에서 열람 시 "손상된 파일" 경고 없이 오픈되는지 확인 | [ ] |
| | PPTX 개별 편집 무결성 | PowerPoint에서 텍스트 및 도형이 비트맵 이미지가 아니라 **개별 텍스트 박스로 편집 가능한지** 확인 | [ ] |
| | Windows 파일 경로 처리 | `C:\Users\...` 경로 및 역슬래시(`\`) 환경에서 로컬 이미지 및 폰트 정상 렌더링 확인 | [ ] |
| | SmartScreen / 보안 점검 | `goslide.exe` 실행 시 보안 경고 발생 시 "추가 정보 ➔ 실행"으로 정상 기동되는지 확인 | [ ] |

---

## 7. 문서화 및 릴리즈 에셋 정비 계획 (Documentation & Release Readiness)

성공적인 v1.0.0 출시와 개발자 온보딩 경험(DX)을 극대화하기 위해 코드 외적인 문서와 릴리즈 에셋을 체계적으로 정비합니다.

### 7.1 README.md (영문) & README.ko.md (한글) 전면 리뉴얼
| 구분 | 점검 및 정비 항목 | 세부 개선 내용 | 상태 |
| :--- | :--- | :--- | :---: |
| **릴리즈 다운로드** | OS별 사전 빌드 바이너리 링크 | GitHub Releases 최신 다운로드 링크 및 뱃지 전면 배치, Linux/macOS/Windows 원클릭 다운로드 표 제공 | [ ] |
| **빠른 시작 워크플로우** | `goslide init` 기반 개편 | 빈 파일 수동 복사 대신 `goslide init presentation.md --theme clean` 명령어 중심의 3단계 가이드로 개편 | [ ] |
| **다중 포맷 일괄 빌드** | `-f all` 파이프라인 반영 | 포맷별 3회 분할 빌드 예시를 `goslide build presentation.md -f all -o dist/` 단일 명령어로 교체 | [ ] |
| **단축키 & 툴킷 테이블** | 최신 인터랙션 도구 반영 | `B`(화면 암전), `W`(백색 전환), `N`(사이드바), `?`(단축키 도움말), 판서 펜 툴바 안내 갱신 | [ ] |
| **지시어 문법 정합성** | 로컬 지시어 표기 통일 | 본문 예제의 지시어를 실제 파서 표준인 `<!-- _layout: cover -->` (언더스코어 포함)로 일괄 동기화 | [ ] |

### 7.2 공식 예제 및 DSL 사양서 동기화
- [ ] **`examples/dsl/example-dsl.md` & `examples/dsl/example-dsl.ko.md`**:
  - v1.0.0 스펙에 맞춰 라인 하이라이트(`{1,3-5}`), 스텝 애니메이션(`dim-fragments`), 미디어 임베드 등 신규 DSL 문법 반영 확인.
- [ ] **공식 데모 슬라이드 (`examples/demo/demo.md` & `examples/demo/demo.ko.md`)**:
  - `goslide build examples/demo/demo.md -f all` 실행 시 오류 없이 3종 포맷(HTML, PDF, PPTX)이 완벽히 생성되는지 회귀 검증.

### 7.3 GitHub Releases 배포 자동화 및 설치 스크립트 점검
- [ ] **OS별 아카이브 구조**:
  - `goslide_1.0.0_linux_amd64.tar.gz`, `goslide_1.0.0_darwin_arm64.tar.gz`, `goslide_1.0.0_windows_amd64.zip` 패키징 무결성.
- [ ] **무결성 검증 파일**:
  - 각 바이너리의 SHA-256 해시를 담은 `checksums.txt` 자동 생성 및 검증.
- [ ] **빠른 설치 스크립트 지원 검토**:
  - `curl -sSL https://raw.githubusercontent.com/joincdream/goslide/main/install.sh | sh` 형태의 원라이너 설치 지원 여부 검토.

---

## 8. 출시 판정 게이트 (Release Gate / Definition of Done)

공식 `v1.0.0` 태그를 발행하기 위한 최종 통과 기준:

1. **테스트 게이트**:
   * 단위 테스트 78건 + 신규 보강 단위 테스트 100% 통과 (`go test -v -race ./...`).
   * 브라우저 인터랙션 자동화 E2E 테스트 통과.
2. **크로스 플랫폼 게이트**:
   * 3대 OS 실기기 스모크 테스트에서 **P0(치명적 크래시) 및 P1(슬라이드 깨짐, PPTX 열람 불가) 결함 0건**.
3. **성능 게이트**:
   * 50장 복합 슬라이드 일괄 빌드 시 메모리 피크 100MB 이하, 빌드 소요 시간 3초 이내.
4. **배포 게이트**:
   * GitHub Actions 가상머신에서 Linux, macOS, Windows 5개 바이너리 및 아카이브, `checksums.txt` 정상 업로드 확인.
5. **문서화 및 온보딩 게이트 (Documentation & Onboarding Gate)**:
   * `README.md` 및 `README.ko.md`에 최신 CLI 워크플로우(`init`, `-f all`) 및 단축키 표 100% 반영 확인.
   * GitHub Releases 공식 다운로드 링크(`https://github.com/joincdream/goslide/releases/latest`) 유효성 확인.
   * README에 기재된 예제 코드 및 데모 파일 빌드 실행 시 에러 0건 확인.

