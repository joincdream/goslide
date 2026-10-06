---
title: Serve 모드 내장 웹 에디터 및 Vim 모드 아키텍처 설계서 (Embedded Editor & Vim Mode)
type: architecture
description: Architectural design for built-in web markdown editor with native Vim keybindings, split-view live rendering, and silent auto-save in goslide serve
tags:
  - goslide
  - architecture
  - editor
  - vim-mode
  - codemirror
  - serve
timestamp: 2026-10-06T00:05:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Serve 모드 내장 웹 에디터 및 Vim 모드 아키텍처 설계서 (Embedded Editor & Vim Mode)

본 문서는 `goslide serve` 로컬 개발 서버 환경에서 사용자가 별도의 데스크톱 에디터(VS Code 등)에 의존하지 않고, **웹 브라우저 단독으로 구글 슬라이드처럼 마크다운을 편집하고 16:9 와이드스크린 슬라이드를 실시간 확인·제어할 수 있는 "Goslide Studio(내장 웹 에디터)" 및 네이티브 Vim 모드**의 아키텍처와 엔지니어링 명세를 정의합니다.

---

## 1. 개요 및 비전 (Overview & Vision)

Goslide의 핵심 철학은 **"외부 의존성 제로(CGO 0%, Node 0%)의 단일 정적 바이너리"**입니다. 

Marp처럼 특정 에디터(VS Code)의 확장에 종속되는 방식은 비-VS Code 사용자(Vim, Emacs, JetBrains, 메모장, 태블릿 사용자)를 배제하게 됩니다. 반면 `goslide serve` 내부에 **초경량 웹 마크다운 에디터와 네이티브 Vim 키바인딩 모드**를 내장하면:

1. **에디터 무관 100% 독립성**: 터미널에서 바이너리 하나만 실행하면 브라우저에 완결된 슬라이드 제작 스튜디오가 열립니다.
2. **Vim 개발자의 거부감 0%**: Obsidian이 성공한 핵심 요인인 `@replit/codemirror-vim`을 탑재하여, 브라우저 안에서도 손가락이 기억하는 `:w`, `dd`, `ci"`, `yy`, `p` 단축키를 네이티브로 제공합니다.
3. **디스크 I/O 병목이 없는 메모리 직결 렌더링**: 에디터 타이핑이 디스크 저장-감시 루프를 거치지 않고 브라우저 메모리 내에서 0ms로 슬라이드 뷰에 즉각 반영됩니다.
4. **하이브리드 유연성**: 내장 에디터를 써도 되고, 기존처럼 터미널 외부 Vim이나 Cursor에서 파일을 고쳐도 기존 `watcher.go`가 감지하여 양방향으로 동기화됩니다.

---

## 2. 시스템 아키텍처 및 듀얼 모드 흐름 (Dual-Engine Architecture)

```mermaid
flowchart TD
    subgraph Browser["브라우저 : Goslide Studio (Svelte 5 Runes + CodeMirror 6)"]
        UI["Studio Layout\n[썸네일 바 | 마크다운 에디터 | 16:9 슬라이드 프리뷰]"]
        CM["CodeMirror 6 Editor Engine\n(Markdown 모드 + @replit/codemirror-vim 플러그인)"]
        VimSwitch["Vim 모드 토글 스위치\n[ Normal Mode ↔ Vim Keymap ]"]
        MemSync["In-Memory Reactive Bridge\n(타이핑 ➔ Svelte Store 0ms 즉각 슬라이드 갱신)"]
        AutoSaver["Debounced Auto-Saver\n(500ms 디바운스 or :w 명령 감지)"]

        UI --> CM
        VimSwitch --> CM
        CM --> MemSync
        CM --> AutoSaver
    end

    subgraph GoServer["Go 로컬 백엔드 (internal/server)"]
        HTTP["HTTP API Router\n(/api/v1/content, /api/v1/themes)"]
        AtomicWriter["Atomic File Writer\n(임시 파일 생성 ➔ 원자적 os.Rename)"]
        FSWatcher["fsnotify Watcher & 100ms Debouncer\n(internal/server/watcher.go)"]
        SSEHub["SSE Event Hub\n(/events 브로드캐스트)"]

        HTTP --> AtomicWriter
        FSWatcher --> SSEHub
    end

    subgraph LocalDisk["로컬 파일 시스템 (Local OS)"]
        SlideFile["talk.md (슬라이드 원본 소스)"]
        ExternalEditor["외부 에디터 (NeoVim, VS Code, Cursor)"]
    end

    MemSync --> UI
    AutoSaver -->|"POST /api/v1/content"| HTTP
    AtomicWriter -->|"안전한 디스크 영속화"| SlideFile
    ExternalEditor -->|"외부에서 편집 및 저장 (:w)"| SlideFile
    SlideFile -->|"OS 파일 이벤트 감지"| FSWatcher
    SSEHub -->|"SSE reload 신호"| UI
```

### 2.1 사용자 E2E 시퀀스 (Deck Hub & Starter Onboarding Sequence)

사용자가 `goslide serve`를 실행하여 기존 데크(Deck)를 탐색하거나 새 데크를 생성해 편집으로 진입하는 전체 흐름입니다:

```mermaid
sequenceDiagram
    autonumber
    actor User as 사용자 (브라우저)
    participant UI as Web Frontend (Svelte 5)
    participant API as Go Server (internal/server)
    participant FS as 로컬 파일 시스템 (Local Disk)

    Note over User, FS: 1. 초기 진입 (인자 없이 goslide serve 구동)
    User->>UI: http://localhost:8080 접속
    UI->>API: GET /api/v1/files (현재 폴더의 *.md Deck 목록 요청)
    API->>FS: 디렉토리 스캔
    FS-->>API: ["talk.md", "arch.md"] 반환
    API-->>UI: Deck 파일 목록 JSON 반환
    UI-->>User: [Deck Hub 화면]<br>• Slot 1: ➕ "+ New Deck" 카드<br>• 기존 Deck 카드 그리드 목록

    Note over User, FS: 2. 신규 Deck 생성 & 테마 선택
    User->>UI: "+ New Deck" 카드 클릭
    UI->>API: GET /api/v1/themes (사용 가능한 테마 목록)
    API-->>UI: ["clean", "dark", "default"]
    UI-->>User: [Theme Selection 모달 표시]<br>• 테마 카드 캐러셀<br>• [Skip: Clean] 빠른 생성 버튼<br>• Filename 입력창 (예: presentation.md)

    User->>UI: 테마 선택 ("clean") 및 "Create" 클릭
    
    Note over UI, FS: 3. 온보딩 템플릿(Starter Deck) 주입 및 파일 생성
    UI->>API: POST /api/v1/files<br>{ filename: "presentation.md", theme: "clean", template: "starter" }
    API->>FS: starter 템플릿 마크다운 파일 원자적 생성
    FS-->>API: 파일 생성 완료
    API-->>UI: 201 Created { file: "presentation.md" }

    Note over User, FS: 4. 3-Pane Studio 전환 및 즉시 렌더링
    UI->>UI: Studio 모드로 자동 라우팅 (/studio?file=presentation.md)
    UI->>API: GET /api/v1/content?file=presentation.md
    API-->>UI: 온보딩 마크다운 원문 반환
    UI-->>User: [3-Pane Studio 화면 오픈]<br>• 좌측: Slide Navigator (3 Slides 썸네일)<br>• 중앙: CodeMirror(Vim 모드) 에디터<br>• 우측: 16:9 슬라이드 프리뷰 실시간 렌더링
```

---

## 3. 세부 설계 명세 (Specification)

### 3.0 덱 허브 & 온보딩 스타터 (Deck Hub & Starter Template)

`goslide serve` 실행 시 특정 마크다운 파일을 지정하지 않은 경우, 사용자가 마주하게 되는 **덱 허브(Deck Hub Dashboard)** 및 **새 덱 생성 모달(New Deck Modal)**의 UI 레이아웃 설계입니다.

> [!NOTE] 도메인 용어 체계 (Ubiquitous Language)
> * **`Deck` (프레젠테이션 파일 단위)**: 마크다운 파일 1개 전체 (`talk.md`). 생성 버튼은 `+ New Deck`, 목록은 `Decks`로 통일합니다.
> * **`Slide` (낱장 페이지 단위)**: 데크 내부에서 `---` 구분자로 나뉘는 개별 슬라이드 (`12 Slides`, `Slide 1 / 12`).

#### 3.0.1 덱 허브 대시보드 레이아웃 (Deck Hub Layout)

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│ GOSLIDE HUB            [ Search decks... (Ctrl+K) ]          Dir: ~/workspace/slides   │
├────────────────────────────────────────────────────────────────────────────────────────┤
│ Decks                                                                                  │
│ ┌────────────────────────┐ ┌────────────────────────┐ ┌────────────────────────┐       │
│ │ ....................   │ │ ┌────────────────────┐ │ │ ┌────────────────────┐ │       │
│ │ :                  :   │ │ │ 16:9 Mini Preview  │ │ │ │ 16:9 Mini Preview  │ │       │
│ │ :        +         :   │ │ └────────────────────┘ │ │ └────────────────────┘ │       │
│ │ :     New Deck     :   │ │ talk.md                │ │ architecture.md        │       │
│ │ :                  :   │ │ Microservices Core     │ │ Pipeline Decoupling    │       │
│ │ :  (Choose Theme)  :   │ │ 12 slides - 10m ago    │ │ 24 slides - 2h ago     │       │
│ │ :..................:   │ │ [ Edit ]   [ Present ] │ │ [ Edit ]   [ Present ] │       │
│ └────────────────────────┘ └────────────────────────┘ └────────────────────────┘       │
│                                                                                        │
│ Themes (Click to Start New Deck)                                                       │
│ ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐ │
│ │ ┌──────────────┐ │  │ ┌──────────────┐ │  │ ┌──────────────┐ │  │ ┌──────────────┐ │ │
│ │ │ Aa Bb Cc     │ │  │ │ Aa Bb Cc     │ │  │ │ Aa Bb Cc     │ │  │ │ Aa Bb Cc     │ │ │
│ │ └──────────────┘ │  │ └──────────────┘ │  │ └──────────────┘ │  │ └──────────────┘ │ │
│ │ Clean (Default)  │  │ Dark (Developer) │  │ Dracula (High)   │  │ Minimal (Simple) │ │
│ │ [ Start Deck ]   │  │ [ Start Deck ]   │  │ [ Start Deck ]   │  │ [ Start Deck ]   │ │
│ └──────────────────┘  └──────────────────┘  └──────────────────┘  └──────────────────┘ │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

1. **상단 글로벌 헤더**:
   * 작업 디렉토리 경로 표기 및 실시간 덱 검색창(`Ctrl + K` / `Cmd + K` 단축키 바인딩).
2. **Slot 1: `+ New Deck` 카드**:
   * 그리드의 첫 번째 카드에 점선 테두리 점보 카드로 배치. 클릭 시 즉시 테마 선택 모달 호출.
3. **덱 카드 그리드 (Decks Grid)**:
   * 디렉토리 내 `*.md` 파일 목록을 16:9 미니 썸네일 카드로 렌더링.
   * `[ Edit ]`: 3-Pane Studio 에디터로 진입하여 편집 시작.
   * `[ Present ]`: 에디터를 거치지 않고 전체화면 발표 뷰로 직행.
4. **하단 테마 갤러리 (Themes Section)**:
   * 시스템 내장 및 `themes/` 폴더에 설치된 테마 목록을 폰트/컬러 팔레트 미리보기 카드로 상시 노출.
   * 원하는 테마 카드의 `[ Start Deck ]`을 클릭하거나 카드를 더블클릭하면, **별도 모달을 거치지 않고 해당 테마가 적용된 새 데크가 즉시 생성되어 3-Pane Studio로 직행** (원클릭 테마 기반 데크 생성: 1-Click Theme-to-Deck).

---

#### 3.0.2 신규 덱 생성 모달 레이아웃 (New Deck Modal Layout)

```
┌────────────────────────────────────────────────────────┐
│  + Create New Deck                                     │
├────────────────────────────────────────────────────────┤
│  Filename                                              │
│  [ presentation.md                                   ] │
│                                                        │
│  Select Theme                                          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐    │
│  │ [Selected]   │ │              │ │              │    │
│  │   Clean      │ │    Dark      │ │   Dracula    │    │
│  │  (Default)   │ │ (Engineering)│ │ (High-Dark)  │    │
│  └──────────────┘ └──────────────┘ └──────────────┘    │
├────────────────────────────────────────────────────────┤
│             [ Cancel ]   [ Skip: Clean ]   [  Create ] │
└────────────────────────────────────────────────────────┘
```

1. **테마 선택 캐러셀 (Theme Selection)**:
   * 시스템 내장 및 `themes/` 디렉토리에 존재하는 커스텀 테마 목록을 시각적 카드로 제공.
   * `Skip: Clean` 버튼을 누르면 기본 `clean` 테마가 자동 적용되어 초고속 생성 보장.
2. **온보딩 스타터 덱 주입 (Starter Deck)**:
   * 생성되는 파일은 빈 문서가 아닌, Goslide 핵심 기능을 즉시 학습할 수 있는 **3장의 인터랙티브 온보딩 슬라이드**로 구성됩니다:
     * **Slide 1 (Cover)**: 프론트매터 메타데이터(`title`, `theme`) 시연 및 커버 슬라이드.
     * **Slide 2 (Feature)**: 2단 컬럼(`layout: two-cols`), 코드 블록 구문 강조, 불릿 포인트 예제.
     * **Slide 3 (Tips & Shortcuts)**: `:w` 저장 단축키, `F` 키 전체화면, `P` 키 발표자 뷰 등 핵심 사용 팁.

---

### 3.1 UI 레이아웃 구조 (3-Pane Studio Layout)

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│ GOSLIDE STUDIO       [Theme: Clean v] [Vim: ON/OFF]    [ Fullscreen ]    [ Export v ] │
├─────────────────────────┬─────────────────────────────────┬────────────────────────────┤
│ [ Slide Navigator ]     │ [ Markdown Editor (CodeMirror) ]│ [ 16:9 Slide Preview ]     │
│                         │                                 │                            │
│ ┌─────────────────────┐ │ 1: ---                          │ ┌────────────────────────┐ │
│ │ 1. Cover Slide      │ │ 2: theme: clean                 │ │                        │ │
│ └─────────────────────┘ │ 3: ---                          │ │   Microservices Core   │ │
│ ┌─────────────────────┐ │ 4: # Microservices Core         │ │                        │ │
│ │ 2. Architecture     │ │ 5:                              │ └────────────────────────┘ │
│ └─────────────────────┘ │ 6: - Pipeline Decoupling        │                            │
│ ┌─────────────────────┐ │ 7: - Zero-Node Portability      │                            │
│ │ 3. Deployment       │ │ 8:                              │ [ < Prev ]   Slide 1 / 12  │
│ └─────────────────────┘ │ 9: <!-- note: Speaker Notes --> │ [ Next > ]                 │
└─────────────────────────┴─────────────────────────────────┴────────────────────────────┘
```

1. **좌측 썸네일 바 (Slide Navigator)**:
   * 슬라이드별 미니 썸네일 카드 나열.
   * 클릭 시 에디터 커서가 해당 슬라이드의 `---` 줄로 점프하며, 우측 프리뷰도 해당 슬라이드로 이동.
2. **중앙 에디터 (Markdown + Vim)**:
   * CodeMirror 6 기반 마크다운 구문 강조 에디터.
   * 우측 상단 토글로 Vim 모드 즉시 ON/OFF 전환.
3. **우측 슬라이드 프리뷰**:
   * 에디터에서 타이핑하는 내용이 메모리 상에서 실시간 렌더링.
   * 편집 집중도를 위해 판서/레이저 도구는 에디터 프리뷰 영역에서 숨기며, 순수 슬라이드 네비게이션 컨트롤(`Prev / Next / Index`)만 배치합니다.
   * 펜 판서 및 레이저 포인터는 상단의 `[ Fullscreen ]` 버튼 또는 단축키(`F`)로 **전체화면 발표 모드**에 진입했을 때만 활성화됩니다.

---

### 3.2 CodeMirror 6 및 네이티브 Vim 모드 연동

Obsidian과 동일한 검증된 스택인 `@replit/codemirror-vim`을 채택합니다:

```javascript
// web/src/components/studio/EditorPane.svelte
import { onMount } from 'svelte';
import { EditorView, basicSetup } from 'codemirror';
import { markdown } from '@codemirror/lang-markdown';
import { vim, Vim } from '@replit/codemirror-vim';
import { Compartment } from '@codemirror/state';

let editorContainer;
let view;
let isVimMode = $state(true);
const vimCompartment = new Compartment();

// :w 및 :write 커맨드 인터셉트
Vim.defineEx('write', 'w', () => {
  saveContentImmediately();
});

function initEditor(initialText) {
  view = new EditorView({
    doc: initialText,
    extensions: [
      basicSetup,
      markdown(),
      vimCompartment.of(isVimMode ? vim() : []),
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          const newContent = update.state.doc.toString();
          // 1. 메모리 직결 프리뷰 갱신 (0ms 지연)
          deck.updateContentFromEditor(newContent);
          // 2. 백그라운드 500ms 디바운스 자동 저장
          scheduleAutoSave(newContent);
        }
      })
    ],
    parent: editorContainer
  });
}

export function toggleVimMode(enabled) {
  isVimMode = enabled;
  view.dispatch({
    effects: vimCompartment.reconfigure(isVimMode ? vim() : [])
  });
}
```

---

### 3.3 백엔드 REST API 명세 (`internal/server/api.go`)

Goslide의 내장 HTTP 서버에 다음 경량 엔드포인트를 추가합니다:

| 엔드포인트 | 메서드 | 설명 | 요청/응답 페이로드 |
| :--- | :---: | :--- | :--- |
| `/api/v1/files` | `GET` | 현재 디렉토리 내 마크다운 슬라이드 목록 반환 | JSON: `[{ name: "talk.md", title: "Microservices", updatedAt: 1728172800 }]` |
| `/api/v1/files` | `POST` | 템플릿 기반 신규 슬라이드 파일 생성 | JSON: `{ filename: "demo.md", theme: "clean", template: "starter" }` |
| `/api/v1/content` | `GET` | 슬라이드 마크다운 원문 스트림 반환 (`?file=talk.md`) | `text/markdown` 원문 텍스트 |
| `/api/v1/content` | `POST` | 에디터에서 수정한 마크다운 원자적 저장 (`?file=talk.md`) | Body: `text/markdown`, Resp: `{"status":"ok","bytes":1420}` |
| `/api/v1/themes` | `GET` | 내장 및 로컬 커스텀 테마 목록 반환 | JSON: `["clean", "dark", "default"]` |
| `/api/v1/status` | `GET` | 서버 상태, 파일 경로, 테마 목록 정보 | JSON: `{ filePath: "talk.md", theme: "clean" }` |

#### 원자적 파일 쓰기 보장 (Atomic File Write)
편집 도중 시스템 충돌이나 전원 차단으로 원본 파일이 깨지는 것을 방지하기 위해 임시 파일 생성 후 Rename 패턴을 적용합니다:

```go
// internal/server/api.go
func (s *Server) handleSaveContent(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read payload", http.StatusBadRequest)
        return
    }

    // 1. 임시 파일 작성
    tempFile := s.targetFilePath + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10)
    if err := os.WriteFile(tempFile, body, 0644); err != nil {
        http.Error(w, "failed to write temp file", http.StatusInternalServerError)
        return
    }

    // 2. 파일 감시기가 자신의 저장 이벤트를 무한 루프로 처리하지 않도록 플래그 억제
    s.suppressWatcherOnce = true

    // 3. 원자적 파일 교체 (Atomic Rename)
    if err := os.Rename(tempFile, s.targetFilePath); err != nil {
        _ = os.Remove(tempFile)
        http.Error(w, "failed to atomic rename file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "size": len(body)})
}
```

---

## 4. 슬라이드 단위 커서 추적 및 동기화 (Cursor-to-Slide Mapping)

1. **커서 이동 ➔ 슬라이드 포커스 점프**:
   * 에디터에서 커서가 이동할 때(`selectionChange`), 현재 커서 위치의 상위에 있는 가장 가까운 `---` 구분자의 개수를 세어 현재 `slideIndex`를 산출합니다.
   * 우측 슬라이드 프리뷰가 해당 슬라이드로 부드럽게 자동 스크롤(Auto-focus)됩니다.
2. **슬라이드 클릭 ➔ 에디터 라인 점프**:
   * 우측 프리뷰 또는 좌측 썸네일에서 슬라이드를 클릭하면, 에디터의 커서가 해당 슬라이드의 시작 라인(`---` 다음 줄)으로 즉시 점프합니다.

---

## 5. 단계별 실행 계획 (5-Phase Execution Plan)

| 단계 | 작업 내용 | 대상 파일 |
| :---: | :--- | :--- |
| **Phase 1** | **백엔드 REST API 구현**<br>• 파일 목록/생성 API (`GET/POST /api/v1/files`)<br>• `GET/POST /api/v1/content` 엔드포인트 구현<br>• 원자적 파일 쓰기 및 watcher 연쇄 감지 억제 플래그 추가 | `internal/server/server.go`<br>`internal/server/api.go` |
| **Phase 2** | **CodeMirror 6 및 Vim 모드 컴포넌트 개발**<br>• `@replit/codemirror-vim` 연동 및 `:w` 커맨드 인터셉트<br>• Vim/Normal 모드 토글 상태 관리 | `web/src/components/studio/EditorPane.svelte`<br>`web/package.json` |
| **Phase 3** | **Studio 3-Pane 레이아웃 및 슬라이드 허브 구축**<br>• 슬라이드 허브(파일 목록 카드 + 새 슬라이드 생성 모달)<br>• 썸네일 네비게이터, 에디터, 16:9 슬라이드 프리뷰 통합<br>• 반응형 스플릿 바(Splitter) 리사이즈 지원 | `web/src/components/studio/SlideHub.svelte`<br>`web/src/components/studio/StudioLayout.svelte`<br>`web/src/stores/studio.svelte.js` |
| **Phase 4** | **메모리 직결 반응형 렌더링 & 커서 동기화**<br>• 디스크 안 거치는 인-메모리 실시간 파싱 및 뷰포트 바인딩<br>• 에디터 커서 라인 $\leftrightarrow$ 슬라이드 인덱스 양방향 동기화 | `web/src/stores/deck.svelte.js` |
| **Phase 5** | **번들 빌드 및 E2E 기능 검증**<br>• `npm run build`로 바이너리 임베드 및 `goslide serve` 구동 검증 | `internal/theme/assets/js/goslide-core.js` |
