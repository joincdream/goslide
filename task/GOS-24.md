# [GOS-24] 마크다운 문법 오류 시 실시간 SSE 에러 오버레이(Error Overlay) 표시 및 서버 회복력 개선 계획서

> **티켓 번호**: [GOS-24](https://joincdream.atlassian.net/browse/GOS-24)  
> **마일스톤**: 개발 서버(`goslide serve`) 회복력(Resilience) 및 개발자 경험(DX) 고도화  
> **상태**: 구현 완료 (Implemented)  
> **담당 패키지**: `internal/server/`, `docs/release/`  
> **연관 소스 파일**:  
> - [`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go)  
> - [`internal/server/server_test.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server_test.go)  
> - [`docs/release/linux-user-test-checklist.md`](file:///home/yundream/myjob/cloit/Goslide/docs/release/linux-user-test-checklist.md) (체크리스트 3.5 항목)  

---

## 1. 개요 및 배경 (Overview & Problem Statement)

### 1.1 현상 및 문제점 분석

현재 `goslide serve` 구동 중 마크다운 또는 Frontmatter(YAML)의 문법 오류가 발생하면, 서버는 다음과 같이 동작합니다:

1. **서버 터미널**: `[goslide] Warning: rebuild failed on file change: ...` 경고 로그만 출력.
2. **브라우저 통신**: SSE로 일괄 `reload` 신호를 브로드캐스트.
3. **브라우저 화면 (Silent Fallback)**:
   - 브라우저가 새로고침되지만 서버가 에러로 인해 직전 캐시(`lastHTML`)를 조용히 응답함.
   - **결과**: 화면에 아무런 시각적 변화나 에러 안내가 없어, 사용자는 "왜 저장이 반영되지 않지?", "핫 리로드가 멈췄나?" 또는 "서버가 죽었나?" 하고 혼란을 겪음 (Silent Failure로 인한 심각한 DX 저하).

### 1.2 개선 목표

1. **Fail-Safe & Non-Crashing**: 문법 오류가 발생해도 `goslide serve` 프로세스는 절대 종료되거나 크래시되지 않음.
2. **Real-time Visual Feedback (Error Overlay)**: Vite, Next.js 등 모던 웹 개발 서버와 같이, 문법 오류 발생 즉시 브라우저 화면 상단/중앙에 눈에 띄는 **에러 오버레이(Error Overlay) 카드**를 띄워 파일명과 구체적인 오류 원인(줄 번호 등)을 직관적으로 안내.
3. **Seamless Self-Healing**: 사용자가 에디터에서 오타를 수정하고 저장하면, 핫 리로드가 동작하면서 에러 오버레이가 자동으로 걷히고 새 슬라이드로 즉시 정상 복원.
4. **리눅스 사용자 테스트 체크리스트 (3.5 항목) 충족**: `docs/release/linux-user-test-checklist.md` 기준을 완벽하게 만족.

---

## 2. 상세 기술 설계 (Technical Design)

### 2.1 SSE 메시지 프로토콜 구조화 (JSON Schema)

기존 단순 문자열(`reload`) 전송에서 구조화된 JSON 페이로드 포맷으로 확장합니다:

```json
{
  "type": "reload" | "error",
  "title": "오류 제목 (예: Markdown Parse Error)",
  "message": "상세 오류 내용 (예: yaml: line 4: did not find expected key)",
  "file": "파일명 (예: talk.md)"
}
```

#### Go 데이터 모델 ([`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go))
```go
type EventType string

const (
	EventReload EventType = "reload"
	EventError  EventType = "error"
)

type SSEMessage struct {
	Type    EventType `json:"type"`
	Title   string    `json:"title,omitempty"`
	Message string    `json:"message,omitempty"`
	File    string    `json:"file,omitempty"`
}
```

### 2.2 브라우저 Error Overlay UI 설계

* **위치 및 스타일**:
  * 뷰포트 최상단(`z-index: 99999`)에 고정된 플로팅 알림 모달/배너.
  * 배경: 다크 반투명 블러(`backdrop-filter: blur(8px)`) 및 붉은색 테두리(`border: 2px solid #ef4444`).
  * 텍스트: 선명한 흰색 폰트, 모노스페이스 코드 블록(`background: rgba(0,0,0,0.4); font-family: monospace;`)으로 파싱 에러 스택/메시지 출력.
  * 슬라이드 원본을 완전히 지우지 않고 배경에 유지하여 작성 맥락을 보존.
* **생명주기(Lifecycle)**:
  * `{"type":"error", ...}` 수신 ➔ DOM에 `#goslide-error-overlay` 삽입 또는 텍스트 갱신.
  * `{"type":"reload"}` 수신 ➔ `#goslide-error-overlay` 제거 후 `location.reload()`.

```text
+-------------------------------------------------------------------------+
|  ⚠️ [Goslide Syntax Error] talk.md                                   [×] |
|  ---------------------------------------------------------------------  |
|  yaml: unmarshal errors:                                                |
|    line 4: cannot unmarshal !!seq into string                           |
|                                                                         |
|  Fix the syntax error and save the file to continue live reload.        |
+-------------------------------------------------------------------------+
|                                                                         |
|                     (기존 슬라이드 배경 유지)                            |
|                                                                         |
+-------------------------------------------------------------------------+
```

### 2.3 서버 에러 감지 및 브로드캐스트 파이프라인

```text
[파일 저장 감지 (fsnotify Watcher)]
               │
               ▼
   [s.renderMarkdown(ctx) 실행]
               │
      ┌────────┴────────┐
   (성공)              (실패 - err != nil)
      │                 │
      ▼                 ▼
[s.BroadcastMessage] [s.BroadcastMessage]
type: "reload"       type: "error"
                     title: "Markdown Syntax Error"
                     message: err.Error()
                     file: s.cfg.MarkdownPath
      │                 │
      └────────┬────────┘
               ▼
    [브라우저 EventSource 수신]
         ┌─────┴─────┐
         ▼           ▼
   (location.reload) (showErrorOverlay)
```

---

## 3. 코드 수정 상세 명세 (Code Changes Specification)

### 3.1 [`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go)

1. **타입 정의 및 브로드캐스트 메서드 확장**:
   - `SSEMessage` 구조체 선언.
   - `BroadcastMessage(msg SSEMessage)` 메서드 추가 (JSON 직렬화 후 `clientChan` 전송).
   - 기존 `BroadcastReload()`는 `BroadcastMessage(SSEMessage{Type: EventReload})`를 호출하도록 리팩토링.
   - `BroadcastError(title, msg, file string)` 헬퍼 메서드 추가.

2. **Watcher 이벤트 루프 수정**:
   ```go
   case <-s.watcher.Events():
       start := time.Now()
       _, renderErr := s.renderMarkdown(ctx)
       if renderErr != nil {
           log.Printf("[goslide] Warning: rebuild failed on file change: %v", renderErr)
           s.BroadcastError("Markdown Syntax Error", renderErr.Error(), filepath.Base(s.absFile))
       } else {
           log.Printf("[goslide] Rebuild completed in %v. Broadcasting reload...", time.Since(start))
           s.BroadcastReload()
       }
   ```

3. **클라이언트 주입 스크립트(`sseScript`) 강화**:
   - EventSource의 `onmessage`에서 `JSON.parse` 파싱.
   - `error` 타입 수신 시 `#goslide-error-overlay` 엘리먼트를 동적으로 생성/주입.
   - `reload` 타입 수신 시 에러 엘리먼트 제거 후 `location.reload()`.
   - 레거시 단순 텍스트 `reload` 수신 시에도 하위 호환 처리.

### 3.2 [`internal/server/server_test.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server_test.go)

1. `TestServer_SSE_BroadcastReload`:
   - JSON 직렬화된 `{"type":"reload"}` 수신 검증.
2. `TestServer_SSE_BroadcastError` (신규 추가):
   - `s.BroadcastError(...)` 호출 시 SSE 스트림을 통해 `{"type":"error", ...}` JSON 페이로드가 정상 브로드캐스트되는지 검증.
3. `TestServer_Watcher_ErrorHandling` (신규 또는 보강):
   - 문법이 깨진 파일 저장 시 서버가 다운되지 않고 `error` 이벤트를 정상 송신하는지 검증.

### 3.3 [`docs/release/linux-user-test-checklist.md`](file:///home/yundream/myjob/cloit/Goslide/docs/release/linux-user-test-checklist.md)

- 3.5 항목 기대 동작 문구 갱신:
  - 기존: "터미널에 에러 로그가 출력되지만 서버가 다운되지 않고, 브라우저에는 직전 유효 슬라이드가 안정적으로 유지되어야 함"
  - 변경: "서버가 다운되지 않고 터미널에 에러 로그가 출력되며, **브라우저 화면 상단에 에러 오버레이(Error Overlay)가 표시되어 문법 오류 원인을 즉시 인지할 수 있어야 함**"

---

## 4. 단계별 실행 계획 (Implementation Steps)

| 단계 | 작업 내용 | 검증 기준 |
| :---: | :--- | :--- |
| **Step 1** | `internal/server/server.go`에 `SSEMessage` 구조체 및 `BroadcastError()` / `BroadcastReload()` 구현 | 컴파일 통과 및 단위 테스트 작성 기반 마련 |
| **Step 2** | `sseScript`에 Error Overlay HTML/CSS DOM 렌더러 및 리로드 핸들러 추가 | 구문 오류 메시지 DOM 렌더링 및 복구 스크립트 무결성 |
| **Step 3** | `server.go` Watcher 루프에 `BroadcastError` 연동 | 파일 저장 시 에러 브로드캐스트 동작 확인 |
| **Step 4** | `internal/server/server_test.go` 단위 테스트 추가 및 `go test -v -race ./internal/server/...` 검증 | 100% 테스트 통과 |
| **Step 5** | `docs/release/linux-user-test-checklist.md` 3.5 항목 갱신 | 체크리스트 기준 최신화 |
| **Step 6** | Jira 티켓([GOS-24](https://joincdream.atlassian.net/browse/GOS-24)) 상태 업데이트 및 보고 | 완료 |
