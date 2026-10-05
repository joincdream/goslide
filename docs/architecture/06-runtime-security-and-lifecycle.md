---
title: 런타임 보안 및 브라우저 프로세스 생명주기 (Runtime Security & Process Lifecycle)
type: operations
description: Headless browser process isolation, zombie process prevention, concurrency throttling, and security hardening
tags:
  - chromedp
  - process-isolation
  - zombie-prevention
  - concurrency
  - security
  - xss
related_code:
  - internal/exporter/pdf/
  - internal/exporter/pptx/
  - internal/parser/
timestamp: 2026-10-05T10:30:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# 6. 런타임 보안 및 브라우저 프로세스 생명주기 (Security & Lifecycle)

## 6.1 Headless 브라우저 프로세스 격리 및 회수 (Zombie Prevention)

`chromedp` 실행 시 발생할 수 있는 프로세스 누수(Zombie Process) 및 시스템 자원 고갈을 방지하기 위해 엄격한 라이프사이클 격리 정책을 구현합니다.

```go
func (e *PDFExporter) Export(ctx context.Context, deck *model.Deck, outPath string) (err error) {
    // 1. 타임아웃 컨텍스트 결합 (최대 60초 보장)
    ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
    defer cancelTimeout()

    // 2. Headless Exec Allocator 생성 (임시 사용자 디렉토리 사용)
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoSandbox,
        chromedp.Headless,
    )
    allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
    defer cancelAlloc() // 서브프로세스 트리 즉시 종료 보장

    // 3. 브라우저 탭 컨텍스트 생성
    taskCtx, cancelTask := chromedp.NewContext(allocCtx)
    defer cancelTask()

    // 4. 실행 및 에러 감지
    // ... chromedp.Run(taskCtx, ...)
    return nil
}
```

### 핵심 수명주기 3대 규칙:
1. **프로세스 트리 자동 회수**: `cancelAlloc()` 함수 호출을 반드시 `defer`로 등록하여, 고루틴 패닉이나 예기치 않은 오류 시에도 Chrome 서브프로세스가 부모 프로세스와 함께 즉시 SIGKILL되어 회수되도록 보장합니다.
2. **명시적 타임아웃**: 파일 I/O나 렌더링 지연으로 인해 무한 대기에 빠지지 않도록 최대 60초(또는 CLI 플래그 지정값) 타임아웃 컨텍스트를 강제합니다.
3. **동시성 스로틀링 (Concurrency Throttling)**: 대용량 슬라이드 캡처 시 CPU 및 메모리 스파이크를 방지하기 위해 최대 동시 브라우저 작업 수를 제한하는 워커 풀 또는 세마포어(`sync.Semaphore`) 패턴을 적용합니다.

---

## 6.2 보안 샌드박싱 및 입력 격리 (Security Hardening)

### Raw HTML 및 악성 스크립트 방어
마크다운 소스는 잠재적으로 신뢰할 수 없는 외부 입력일 수 있으므로 XSS(Cross-Site Scripting) 및 원격 코드 실행을 방어합니다:

* **기본값 (`--unsafe-html=false`)**:
  - Goldmark 파서의 `html.WithUnsafe()`를 기본적으로 비활성화합니다.
  - 마크다운 본문에 삽입된 `<script>`, `<iframe>`, `<object>` 등의 위험 태그는 파서 단계에서 안전하게 이스케이프(HTML Entity 치환) 처리됩니다.
* **명시적 허용 (`--unsafe-html=true`)**:
  - 사용자가 신뢰할 수 있는 소스에 대해 명시적으로 플래그를 지정했을 때만 원본 Raw HTML 태그 삽입을 허용합니다.
* **로컬 파일 접근 통제**:
  - Headless 브라우저 렌더링 시 외부 임의 파일 경로에 대한 접근을 제한하고, 프로젝트 작업 디렉토리 하위의 리소스만 참조하도록 격리합니다.

---

## 6.3 관련 문서

* **[`01-principles-and-vision.md`](./01-principles-and-vision.md)**: 원칙 5 (Context 기반 제어 및 리소스 누수 방지)
* **[`05-rendering-and-export-pipelines.md`](./05-rendering-and-export-pipelines.md)**: chromedp를 활용하는 PDF/PPTX 파이프라인
* **[`07-operations-build-and-cicd.md`](./07-operations-build-and-cicd.md)**: 브라우저 실패 시의 Exit Code 매핑 (`ExitExportFailed`)
