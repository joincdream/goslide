# [GOS-36] Chrome 미설치 안내 가이드 에러 래핑 및 CLI 런타임 사용성 개선 계획서

> **티켓 번호**: [GOS-36](https://joincdream.atlassian.net/browse/GOS-36)  
> **마일스톤**: v1.0.0 Release Polish & Error UX  
> **마감일**: 2026-10-31  
> **상태**: 진행 중 (In Progress)  
> **담당 패키지**:  
> - [`internal/browser/`](file:///home/yundream/myjob/cloit/Goslide/internal/browser/) (크롬 탐색 및 플랫폼별 설치 가이드 에러)  
> - [`cmd/goslide/serve.go`](file:///home/yundream/myjob/cloit/Goslide/cmd/goslide/serve.go) (개발 서버 시작 배너 및 포트 충돌 에러 래핑)  
> - [`internal/i18n/`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/) (한/영 에러 및 도움말 메시지 카탈로그)  
> **연관 문서**:  
> - [`task/release/end-to-end-user-journey.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/end-to-end-user-journey.md) *(사용자 여정 및 갭 분석)*  
> - [`task/release/linux-user-test-checklist.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/linux-user-test-checklist.md) *(사전 환경 및 런타임 검증)*  

---

## 1. 개요 및 배경 (Context & Objective)

### 1.1 배경 및 문제의식
[`task/release/end-to-end-user-journey.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/end-to-end-user-journey.md)의 6단계(빌드) 및 4단계(개발 서버) 분석 결과, 신규 사용자가 런타임 환경에서 겪는 두 가지 핵심 UX 마찰이 확인되었습니다:

1. **Chrome/Chromium 미설치 환경에서의 불친절한 에러**:
   - Goslide의 PDF 및 PPTX 내보내기는 헤드리스 브라우저(`chromedp`)를 사용합니다.
   - 현재 시스템에 Chrome/Chromium이 없을 경우, [`internal/browser/browser.go`](file:///home/yundream/myjob/cloit/Goslide/internal/browser/browser.go#L83)는 단순 영문 에러(`could not locate Google Chrome or Chromium in PATH...`)를 반환합니다.
   - 사용자는 자신의 OS(Ubuntu, Fedora, Arch, macOS 등)에서 어떤 명령어로 브라우저를 설치해야 하는지 알지 못해 작업을 중단하거나 구글링을 강요받습니다.
2. **개발 서버(`goslide serve`)의 헤드리스/원격 터미널 사용성 부족**:
   - WSL2, 헤드리스 리눅스 서버, SSH 원격 개발 환경에서는 브라우저 자동 오픈(`--open`)이 실패하거나 불가능합니다.
   - 터미널 출력에 로컬 접속 URL(`http://localhost:8080`)과 프레젠테이션 핵심 단축키(`?`, `P`, `D`, `F`)가 직관적으로 노출되지 않고 단순 텍스트로 흘러가 사용자가 다음 행동을 파악하기 어렵습니다.
3. **8080 포트 점유(Port Collision) 시 원시 에러 노출**:
   - 다른 웹 서버나 프로세스가 8080 포트를 점유하고 있을 때 `failed to bind address localhost:8080: bind: address already in use`라는 시스템 소켓 에러가 그대로 출력됩니다.
   - `--port` 플래그로 포트를 손쉽게 변경할 수 있다는 가이드가 제공되어야 합니다.

### 1.2 개발 목표 (Objectives)
* **가이드형 에러 UX 제공 (Actionable Error Messages)**: 크롬이 없는 환경에서 에러 발생 시, OS를 자동 판별하여 터미널에 복사해 실행할 수 있는 패키지 설치 명령어(Ubuntu: `sudo apt install -y chromium-browser` 등)를 즉시 안내합니다.
* **개발 서버 시작 안내 강화**: 터미널 배너를 클린 엔지니어링 스타일로 정돈하고, 접속 URL 강조 및 브라우저에서 사용할 주요 단축키 팁을 함께 표시합니다.
* **포트 충돌 복구 가이드**: 포트 바인딩 실패 시 직관적인 `--port` 변경 팁을 제공합니다.
* **다국어(한/영) 100% 지원**: 에러 및 안내 메시지를 `internal/i18n` 시스템과 연동합니다.

---

## 2. 세부 사양 및 설계 (Detailed Specifications)

### 2.1 OS별 Chrome/Chromium 설치 가이드 에러 메시지 설계

[`internal/browser/browser.go`](file:///home/yundream/myjob/cloit/Goslide/internal/browser/browser.go)의 `FindChrome` 함수가 실패했을 때, 현재 OS(`runtime.GOOS`)에 맞춰 복사 가능한 설치 명령어를 조합하여 반환합니다:

#### 한국어 안내 예시
```text
PDF 및 PPTX 생성을 위해 Google Chrome 또는 Chromium 브라우저가 필요합니다.

설치 방법:
  • Ubuntu / Debian: sudo apt update && sudo apt install -y chromium-browser
  • Fedora / RHEL:   sudo dnf install -y chromium
  • Arch Linux:      sudo pacman -S chromium
  • macOS:           brew install --cask google-chrome
  • Windows:         Google Chrome을 설치하거나 Microsoft Edge 경로를 확인하세요.

임의 경로 지정:
  환경 변수 GOSLIDE_CHROME_BIN=<경로> 또는 --chrome-path 플래그를 설정하세요.
```

#### 영문 안내 예시
```text
Google Chrome or Chromium is required for PDF and PPTX export.

Installation:
  • Ubuntu / Debian: sudo apt update && sudo apt install -y chromium-browser
  • Fedora / RHEL:   sudo dnf install -y chromium
  • Arch Linux:      sudo pacman -S chromium
  • macOS:           brew install --cask google-chrome
  • Windows:         Install Google Chrome or verify Microsoft Edge path.

Custom Path:
  Set environment variable GOSLIDE_CHROME_BIN=<path> or use --chrome-path flag.
```

---

### 2.2 개발 서버 (`goslide serve`) 터미널 배너 개선

이모지를 절제한 단정하고 시인성 높은 클린 엔지니어링 포맷으로 개편합니다:

```text
┌─────────────────────────────────────────────────────────────┐
│  Goslide Live Server                                        │
│  URL:      http://localhost:8080                            │
│  Watching: demo.md (Theme: clean)                           │
│  Hotkeys:  [?] Help  [P] Presenter  [D] Pen  [F] Fullscreen │
└─────────────────────────────────────────────────────────────┘
Press Ctrl+C to stop
```

---

### 2.3 포트 바인딩 충돌 친절 래핑

`net.Listen` 에러 검출 시 `bind: address already in use` 또는 시스템 `EADDRINUSE` 소켓 에러를 감지하여 다음과 같이 친절하게 안내합니다:

```text
오류: 포트 8080이(가) 이미 다른 프로그램에서 사용 중입니다.
해결 방법: --port 옵션으로 다른 포트를 지정해 보세요 (예: goslide serve demo.md --port 8081)
```

---

## 3. 단계별 실행 계획 (Execution Phases)

```
[Phase 1: Chrome 에러 래핑] ───► [Phase 2: serve 배너 & 포트 에러] ───► [Phase 3: 단위 테스트 검증]
  - internal/browser/           - cmd/goslide/serve.go                 - browser_test.go
  - internal/i18n/ (카탈로그)   - i18n 연동                           - serve_test.go
```

### Phase 1. `internal/browser` 가이드형 에러 구현
1. `internal/browser/browser.go`에 `FormatChromeInstallGuide(lang string) string` 함수 구현.
2. `FindChrome` 실패 시 포맷팅된 가이드가 포함된 센티넬 에러(`model.ErrChromeNotFound`) 래핑 반환.
3. `internal/i18n`에 관련 메시지 키 등록 (`err.chrome_not_found`, `err.chrome_install_guide`).

### Phase 2. `cmd/goslide/serve.go` 배너 및 포트 에러 개선
1. 서버 구동 시 시인성 높은 프레임 박스 배너 출력.
2. `srv.Start(ctx)`에서 포트 충돌 에러(`address already in use`) 감지 시 직관적인 해결 가이드 출력.

### Phase 3. 단위 테스트 및 회귀 검증
1. `internal/browser/browser_test.go`: 가이드 문자열 포맷팅 및 에러 래핑 단위 테스트.
2. `cmd/goslide/serve_test.go`: 배너 출력 및 옵션 검증 테스트 보강.

---

## 4. 완료 조건 (Definition of Done - DoD)

1. [x] **티켓 상태 전이 완료**: Jira 티켓 [GOS-36](https://joincdream.atlassian.net/browse/GOS-36)가 `진행 중`으로 전환되었는가?
2. [x] **플랫폼별 Chrome 설치 가이드 구현**: `internal/browser`에서 브라우저 미발견 시 OS별 패키지 매니저 설치 명령어가 명확히 안내되는가?
3. [x] **개발 서버 배너 개선**: `goslide serve` 실행 시 접속 URL과 주요 단축키 힌트가 박스형 배너로 가독성 높게 표시되는가?
4. [x] **포트 충돌 가이드 래핑**: 포트 점유 에러 발생 시 `--port <새포트>` 해결 가이드가 안내되는가?
5. [x] **다국어(한/영) 지원**: 에러 및 배너 메시지가 `--lang` 및 시스템 로케일에 따라 올바르게 전환되는가?
6. [x] **단위 테스트 작성**: 관련 단위 테스트(`browser_test.go`, `serve_test.go`)가 작성되었는가?
