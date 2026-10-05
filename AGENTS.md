# AGENTS.md — Goslide Engineering Guidelines & LLM Guardrails

이 문서는 **Goslide**(Go 기반의 Markdown 슬라이드 생성 도구: HTML, PDF, PPTX 변환) 프로젝트에서 LLM 및 AI 코딩 에이전트가 코드를 작성, 리팩토링, 검증할 때 반드시 준수해야 하는 **소프트웨어 엔지니어링 원칙과 가드레일**을 정의합니다.

모든 에이전트는 이 문서에 명시된 규칙을 절대적인 제약 조건(Hard Constraints)으로 받아들여야 하며, 이를 위반하는 코드를 생성해서는 안 됩니다.

---

## 1. 프로젝트 비전 및 핵심 철학

- **목표**: [Marp](https://marp.app/)와 유사한 Markdown 기반 슬라이드 데크 빌더를 순수 Go(또는 최소한의 외부 런타임 의존성)로 구현합니다.
- **입력**: Frontmatter(메타데이터, 테마 지시자), 슬라이드 분할자(`---`), 확장 마크다운 지시어(Directives: 배경, 레이아웃, 클래스 등).
- **출력**: 
  - **HTML**: 웹 브라우저에서 실행 가능한 독립형(Self-contained) 인터랙티브 슬라이드.
  - **PDF**: 인쇄 및 발표 품질의 벡터/고해상도 PDF.
  - **PPTX**: 프레젠테이션 편집이 가능한 Office Open XML 기반 PowerPoint 파일.
- **핵심 가치**:
  1. **단일 바이너리(Single Binary) 지향**: 크로스 플랫폼 빌드가 쉬워야 하며 불필요한 CGO 의존을 지양합니다.
  2. **결정론적 출력(Deterministic Output)**: 동일한 마크다운과 설정에 대해 항상 동일한 결과물을 생성해야 합니다.
  3. **고성능 & 낮은 리소스 소비**: Go의 동시성을 안전하게 활용하여 수백 장의 슬라이드도 밀리초 단위로 파싱/변환합니다.

---

## 2. 시스템 아키텍처 원칙

### 2.1 단방향 파이프라인 (Pipeline Architecture)
모든 변환 작업은 결합도가 낮은 파이프라인 단계를 거칩니다:

```
[Markdown Source]
       │
       ▼ (1. Parse)
[AST & Presentation Model] (Immutable Slide IR)
       │
       ▼ (2. Transform & Resolve)
[Themed Presentation Model] (Styles, Layouts, Directives Resolved)
       │
       ├───────────────────┼───────────────────┐
       ▼ (3a. Render)      ▼ (3b. Export)      ▼ (3c. Export)
     HTML                 PDF                 PPTX
```

### 2.2 패키지 및 모듈 경계 (Go Standard Layout)
에이전트는 다음 패키지 레이아웃을 엄격히 준수하며, 순환 참조(Circular Dependency)를 발생시켜서는 안 됩니다:

```
Goslide/
├── cmd/
│   └── goslide/          # CLI 진입점 (플래그 파싱, 사용자 인터랙션, 종료 코드 제어)
├── pkg/
│   └── goslide/          # 외부에서 라이브러리로 임포트 가능한 공개 API
└── internal/             # 내부 비즈니스 로직 (외부 노출 금지)
    ├── model/            # 슬라이드, 데크, 지시어 등 불변 IR(Intermediate Representation)
    ├── parser/           # Markdown 파서 및 지시어(Directive) 추출기
    ├── theme/            # 테마, CSS 파싱, 폰트 및 에셋 관리 (embed.FS 활용)
    ├── renderer/         # HTML 렌더러 구현체
    │   └── html/
    ├── exporter/         # 외부 포맷 익스포터
    │   ├── pdf/          # PDF 생성 엔진 (예: Headless 브라우저 드라이버 or 네이티브)
    │   └── pptx/         # PPTX(OpenXML) 빌더
    └── testutil/         # 테스트 헬퍼 및 골든 파일(Golden File) 픽스처
```

### 2.3 인터페이스 격리 (DIP / ISP)
- 렌더러와 익스포터는 상호 의존하지 않으며, `internal/model/interfaces.go`에 정의된 인터페이스(`Renderer`, `Exporter`)를 통해 명확히 분리됩니다.

---

## 3. LLM 에이전트 행동 가드레일 (Hard Rules)

### 🔴 금지 사항 (Strict Don'ts)
1. **존재하지 않는 가상 라이브러리/함수 환각(Hallucination) 금지**: 검증되지 않은 패키지를 `go.mod`에 임의로 추가하지 마십시오. 필요한 경우 Go 표준 라이브러리 및 널리 검증된 오픈소스를 우선 제안하십시오.
2. **CGO 의존성 남발 금지**: CGO가 필요한 라이브러리를 추가하기 전 반드시 순수 Go(Pure Go) 대안을 먼저 검토해야 합니다.
3. **`panic()` 남용 금지**: 라이브러리 및 내부 패키지 로직에서 `panic`은 허용되지 않습니다. 모든 오류는 명시적인 `error` 반환으로 처리합니다. (단, 프로그램 시작 시 정적 정규식 컴파일 `regexp.MustCompile` 등 제외)
4. **글로벌 상태(Global State) 사용 금지**: 패키지 레벨 전역 변수(패키지 설정, 캐시 등)는 동시성 버그를 유발하므로 금지합니다. 의존성은 구조체 생성자(`New...`)를 통한 주입(Dependency Injection) 방식을 취합니다.
5. **Context 무시 금지**: 파일 I/O, 네트워크(외부 이미지 다운로드 등), 무거운 변환 작업(Headless 브라우저 프로세스 등)에는 항상 `context.Context`를 전달하고 취소(Cancellation) 신호를 처리해야 합니다.
6. **사용자 미지시 작업 임의 수행 금지 (Strict User-Directed Scope)**: 사용자가 명시적으로 지시하지 않은 추가 작업(테스트 실행, 린트/복잡도 검사, 빌드, 부가 스크립트 실행 등)은 절대로 임의로 먼저 수행하지 마십시오. 오직 사용자가 지시한 작업(예: 티켓 상태 전이, Git 커밋, 특정 파일 수정 등)의 범위 내에서만 엄격히 수행합니다.
7. **OKF 라우터 우회 및 광범위 파일 탐색 금지 (No Blind Grep/Scan)**: 파일 조사 전 반드시 `docs/okf/index.md`의 라우팅 맵만 참조해야 하며, 프로젝트 전체를 대상으로 하는 광범위한 `grep`, 무작위 디렉터리 순회, 임의의 파일 열람을 엄격히 금지합니다.
8. **과잉 조사(Over-investigation) 및 불필요한 연쇄 탐색 금지 (Strict Pinpoint Action)**: 문제 원인이 특정되었거나 수정 대상이 명확한 경우, 연관성이 떨어지는 주변 파일들을 '혹시나' 하는 목적으로 연쇄 조회하는 행위를 절대 금지합니다. 오직 단일 목적 파일 1곳만 핀포인트로 접근하여 최소 단위로 수정합니다.
9. **Jira Task 단위 스코프 엄격 준수 및 반복 회귀 테스트 금지 (Strict Jira Task Scope & No Redundant Testing)**:
   - 모든 작업은 Jira Task 단위로 엄격히 분할되어 있으므로, 작업 범위는 오직 사용자가 지정한 단일 컴포넌트/파일 1곳으로 제한합니다.
   - 전체 코드베이스/문서 스캔 및 무작위 탐색을 엄격히 금지합니다.
   - 이미 앞선 단계에서 검증된 전체 회귀 테스트(`make test`, `make test-race`, `make golden-update`, 헤드리스 브라우저 스크린샷 덤프 등)를 매 작업 단위마다 습관적으로 재실행하여 토큰과 시간을 낭비하는 행위를 절대 금지합니다.
   - 테스트 및 검증 도구는 오직 사용자가 명시적으로 "테스트 실행해주세요" 또는 "검증해주세요"라고 지시한 경우에만 실행합니다.
10. **단순 작업의 최소주의 원칙 및 자의적 연쇄 검증 절대 금지 (KISS: Keep It Simple & No Self-Verification Spiral)**:
   - 단순한 스타일(CSS), 마크다운, 설정, 텍스트 수정 작업은 **[단일 파일 수정 ➔ 즉시 완료 보고]**의 1단계 직행으로 즉시 종료해야 합니다.
   - '제대로 반영되었는지 확인하겠다'는 명목으로 빌드(`go build`, `goslide build`), 빌드 결과물 열람/검색(`grep`, `cat`), 추가 테스트를 자의적으로 연쇄 실행하는 오버엔지니어링(Self-Verification Spiral)을 절대 금지합니다.
   - 빌드나 테스트, 실행 검증은 오직 사용자가 "빌드해주세요" 또는 "검증해주세요"라고 직접 지시한 경우에 한해서만 수행합니다.

### 🟢 권장 사항 (Must-Haves)
1. **명시적 에러 래핑**: Go 1.13+ 에러 래핑 규칙(`fmt.Errorf("...: %w", err)`)을 준수하여 에러 원인을 추적할 수 있도록 합니다.
2. **에러 타입/센티넬 에러 정의**: 호출 측에서 `errors.Is` 또는 `errors.As`로 분기할 수 있도록 명확한 에러 상수를 정의합니다.
3. **리소스 안전성**: 파일, 네트워크 커넥션, 서브프로세스는 반드시 `defer close()` 또는 리소스 정리 로직을 보장합니다.

## 4. 코딩 표준 및 Go 관용구

### 4.1 네이밍 및 가시성
- Go 표준 규칙(Effective Go)을 따릅니다.
- 축약어는 일관된 대소문자를 유지합니다 (`HTMLRenderer`, `PDFExporter`, `PPTXWriter`).
- 패키지명은 간결한 소문자 단수 명사를 사용합니다 (`parser`, `theme`, `model`).

### 4.2 에러 핸들링 원칙
- 모든 에러는 컨텍스트와 함께 명시적으로 래핑(`fmt.Errorf("...: %w", err)`)합니다.
- 호출 측에서 `errors.Is` 또는 `errors.As`로 분기할 수 있도록 명확한 센티넬 에러(`var Err... = errors.New(...)`)를 정의합니다.
- 에러 무시(`_ = ...`)나 raw string 에러 생성을 엄격히 금지합니다.

### 4.3 동시성 및 리소스 누수 방지
- 고루틴(Goroutine) 생성 시에는 반드시 생명주기와 종료 시점을 명확히 정의합니다.
- `sync.WaitGroup` 또는 `errgroup.Group`을 사용하여 고루틴 완료 및 에러 전파를 안전하게 처리합니다.
- 파일, 네트워크, 브라우저 프로세스 등 외부 리소스는 `defer Close()`로 누수를 방지합니다.

---

## 5. 테스트 및 검증 규격 (Testing Guardrails)

- **테스트 동반 필수**: 신규 기능 추가 시 단위 테스트 코드가 수반되지 않은 변경은 완료된 것으로 간주하지 않습니다.
- **테이블 기반 테스트 (Table-Driven Tests)**: 파서, 지시어, 유틸리티 등 다양한 입력 케이스는 테이블 주도 테스트 패턴으로 간결하게 구성합니다.
- **골든 파일 테스트 (Golden File Testing)**: HTML 및 복잡한 구조체 생성 결과는 `testdata/` 디렉토리의 `.golden` 파일과 비교 검증합니다 (`-update` 플래그로 의도된 변경 갱신).
- **린트 및 정적 분석 준수**: 모든 코드는 `go vet`, `gofmt -l`, `golangci-lint` 검사를 통과해야 합니다.

---

## 6. LLM 에이전트 작업 절차 (Step-by-Step Execution Workflow)

### 6.1 단순 작업 (Simple Task: CSS, 마크다운, 문서/오탈자, 설정 수정 등)
1. **타겟 파일 핀포인트 수정 (Single-Target Edit)**: 사용자가 지정한 파일 1곳만 최소 단위로 수정.
2. **즉시 완료 보고 (Immediate Report)**: 자의적인 빌드, 결과 확인 덤프, 테스트 실행 없이 즉시 수정 결과를 보고하고 작업을 종료(Early Exit).

### 6.2 신규 기능 및 복합 구현 작업 (Feature Implementation)
LLM 에이전트는 복합 작업 시 다음 단계를 순차적으로 수행합니다:

1. **사전 분석 및 OKF 지식 검색 (Pre-check & OKF Retrieval)**:
   - 구현하려는 태스크와 관련된 Google OKF 스펙을 필수로 확인하고 맥락을 동기화합니다:
     - 진입점 및 인덱스: `docs/okf/index.md`
     - 아키텍처 및 렌더링 파이프라인: `docs/okf/core-architecture.md`
     - 인터페이스 계약 및 도메인 모델: `docs/okf/contracts-interfaces.md`
     - 코딩 제약 및 검증 체크리스트: `docs/okf/hard-constraints.md`
     - 스크린캐스트 판서 오버레이: `docs/okf/screencast-annotation.md`
     - 아키텍처 단순화 의사결정: `docs/okf/decisions-simplification.md`
   - 관련 모듈의 인터페이스 및 기존 데이터 모델 확인.
   - 불필요한 의존성 추가 방지 및 영향도 평가.
2. **인터페이스 우선 설계 (Interface First)**:
   - 비즈니스 로직 작성 전 타입 정의 및 인터페이스 선언 (`internal/model/interfaces.go` 준수).
3. **구현 및 단위 테스트 동시 작성**:
   - 최소한의 변경으로 요구사항 충족.
   - 예외 케이스(빈 입력, 손상된 마크다운, 타임아웃 등) 테스트 추가.
4. **검증 (Verification - 사용자가 명시적으로 요청한 경우에만 수행)**:
   - 사용자가 테스트나 품질 검증을 요청한 경우에 한해 `go test -v -race ./...`, 린트 등을 실행합니다.
   - 단순 커밋, 티켓 상태 전이, 파일 열람/문서 수정 등 검증을 요청하지 않은 작업 시에는 테스트나 분석 도구를 임의로 실행하지 않습니다.
   - 리소스 해제(`Close()`) 누락 여부 등은 코드 작성 시 정적으로 점검합니다.
5. **결과 보고**:
   - 변경 사항, 설계 결정 배경, 검증 결과를 사용자에게 명확히 보고.

