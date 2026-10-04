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
- 렌더러와 익스포터는 명확한 인터페이스를 통해 분리되어야 하며 상호 의존하지 않습니다.

```go
// internal/model/interfaces.go 예시
type Renderer interface {
    Render(ctx context.Context, deck *Deck, w io.Writer) error
}

type Exporter interface {
    Export(ctx context.Context, deck *Deck, outputPath string) error
}
```

---

## 3. LLM 에이전트 행동 가드레일 (Hard Rules)

### 🔴 금지 사항 (Strict Don'ts)
1. **존재하지 않는 가상 라이브러리/함수 환각(Hallucination) 금지**: 검증되지 않은 패키지를 `go.mod`에 임의로 추가하지 마십시오. 필요한 경우 Go 표준 라이브러리 및 널리 검증된 오픈소스를 우선 제안하십시오.
2. **CGO 의존성 남발 금지**: CGO가 필요한 라이브러리를 추가하기 전 반드시 순수 Go(Pure Go) 대안을 먼저 검토해야 합니다.
3. **`panic()` 남용 금지**: 라이브러리 및 내부 패키지 로직에서 `panic`은 허용되지 않습니다. 모든 오류는 명시적인 `error` 반환으로 처리합니다. (단, 프로그램 시작 시 정적 정규식 컴파일 `regexp.MustCompile` 등 제외)
4. **글로벌 상태(Global State) 사용 금지**: 패키지 레벨 전역 변수(패키지 설정, 캐시 등)는 동시성 버그를 유발하므로 금지합니다. 의존성은 구조체 생성자(`New...`)를 통한 주입(Dependency Injection) 방식을 취합니다.
5. **Context 무시 금지**: 파일 I/O, 네트워크(외부 이미지 다운로드 등), 무거운 변환 작업(Headless 브라우저 프로세스 등)에는 항상 `context.Context`를 전달하고 취소(Cancellation) 신호를 처리해야 합니다.
6. **사용자 미지시 작업 임의 수행 금지 (Strict User-Directed Scope)**: 사용자가 명시적으로 지시하지 않은 추가 작업(테스트 실행, 린트/복잡도 검사, 빌드, 부가 스크립트 실행 등)은 절대로 임의로 먼저 수행하지 마십시오. 오직 사용자가 지시한 작업(예: 티켓 상태 전이, Git 커밋, 특정 파일 수정 등)의 범위 내에서만 엄격히 수행합니다.

### 🟢 권장 사항 (Must-Haves)
1. **명시적 에러 래핑**: Go 1.13+ 에러 래핑 규칙(`fmt.Errorf("...: %w", err)`)을 준수하여 에러 원인을 추적할 수 있도록 합니다.
2. **에러 타입/센티넬 에러 정의**: 호출 측에서 `errors.Is` 또는 `errors.As`로 분기할 수 있도록 명확한 에러 상수를 정의합니다.
3. **리소스 안전성**: 파일, 네트워크 커넥션, 서브프로세스는 반드시 `defer close()` 또는 리소스 정리 로직을 보장합니다.

---

## 4. 모듈별 구현 세부 가이드라인

### 4.1 Parser (Markdown & Directives)
- **기반 파서**: 검증되고 확장성이 뛰어난 Go Markdown 파서(예: `github.com/yuin/goldmark`)의 AST 확장 메커니즘을 사용합니다.
- **슬라이드 분할 규칙**:
  - `---` (수평선)을 기준으로 개별 슬라이드로 분리합니다.
  - 마크다운 Frontmatter(`---`로 시작하는 상단 YAML 블록)와 슬라이드 분할자를 정확히 구분해야 합니다.
- **Directives 지원**:
  - Global Directives: `theme`, `paginate`, `header`, `footer`, `size` (16:9, 4:3 등).
  - Scoped Directives (슬라이드 단위): `<!-- _class: lead -->`, `<!-- backgroundColor: #f0f0f0 -->`.
- **불변성(Immutability)**: 파싱된 `Deck` 및 `Slide` 구조체는 파싱 완료 후 읽기 전용으로 취급되어야 합니다.

### 4.2 HTML Renderer
- **독립형 번들링**:
  - Go의 `embed.FS`를 사용하여 기본 테마 CSS 및 필수 JS 런타임을 바이너리에 내장합니다.
  - `--standalone` 옵션 적용 시 외부 리소스(이미지 등)를 Data URI(Base64)로 인라인 임베딩할 수 있는 옵션을 제공합니다.
- **보안(XSS 방지)**:
  - 마크다운 파싱 시 기본적으로 위험한 스크립트 실행을 방지하도록 샌드박싱 처리를 고려합니다. (Raw HTML 허용 여부는 플래그로 격리)

### 4.3 PDF Exporter
- **Headless 브라우저 연동 방식**:
  - `chromedp/chromedp`를 활용하여 생성된 HTML 슬라이드를 로컬에서 헤드리스로 렌더링 후 `Page.printToPDF` API를 호출하는 방식을 1차 전략으로 취합니다.
  - 브라우저 인스턴스 라이프사이클을 안전하게 관리하고, 타임아웃 발생 시 좀비 프로세스가 남지 않도록 보장합니다.

### 4.4 PPTX Exporter
- **오피스 오픈 XML(OpenXML) 사양 준수**:
  - 슬라이드 레이아웃, 마스터 슬라이드, 텍스트 상자(Rich Text), 표, 코드 블록, 이미지를 PPTX 형태(`p:sp`, `a:p`, `a:r` 등)로 매핑합니다.
  - 슬라이드 크기(16:9 기준 12192000 x 6858000 EMU) 등 정밀한 치수 단위를 구조화합니다.

---

## 5. 코딩 표준 및 Go 관용구

### 5.1 네이밍 및 가시성
- Go 표준 규칙(Effective Go)을 따릅니다.
- 축약어는 일관된 대소문자를 유지합니다 (`HTMLRenderer`, `PDFExporter`, `PPTXWriter`).
- 패키지명은 간결한 소문자 단수 명사를 사용합니다 (`parser`, `theme`, `model`).

### 5.2 에러 핸들링 패턴
```go
// Good: 컨텍스트와 함께 에러 래핑 및 센티넬 에러 지원
var ErrSlideNotFound = errors.New("slide not found")

func (p *Parser) Parse(ctx context.Context, r io.Reader) (*model.Deck, error) {
    if err := ctx.Err(); err != nil {
        return nil, fmt.Errorf("parse canceled: %w", err)
    }
    // ...
    if err != nil {
        return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
    }
    return deck, nil
}

// Bad: 에러 무시, 메시지 손실, raw string 에러
func Parse(r io.Reader) *model.Deck {
    deck, _ := doSomething(r) // 절대 금지
    return deck
}
```

### 5.3 동시성 및 리소스 누수 방지
- 고루틴(Goroutine)을 생성할 때는 반드시 언제, 어떻게 종료되는지 생명주기를 정의해야 합니다.
- `sync.WaitGroup` 또는 `errgroup.Group`을 사용하여 고루틴 완료 및 에러 전파를 확실히 처리합니다.

---

## 6. 테스트 및 검증 규격 (Testing Guardrails)

에이전트가 코드를 작성하거나 변경할 때, **테스트 코드가 수반되지 않은 기능 추가는 완료된 것으로 간주하지 않습니다.**

### 6.1 테이블 기반 테스트 (Table-Driven Tests)
파서, 유틸리티, 디렉티브 처리기는 테이블 기반 테스트 패턴을 작성합니다:

```go
func TestDirectiveParser(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected model.Directives
        wantErr  bool
    }{
        {
            name:  "valid theme directive",
            input: "<!-- theme: default -->",
            expected: model.Directives{Theme: "default"},
            wantErr: false,
        },
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 실행 및 검증
        })
    }
}
```

### 6.2 골든 파일 테스트 (Golden File Testing)
- HTML 및 복잡한 구조체 생성 결과는 `testdata/` 디렉토리에 `.golden` 파일을 두고 비교 검증합니다.
- `-update` 플래그를 통해 의도된 변경 시 골든 파일을 갱신할 수 있는 패턴을 적용합니다.

### 6.3 린트 및 정적 분석 준수
모든 코드는 아래 도구의 검사를 통과해야 합니다:
```bash
go vet ./...
test -z "$(gofmt -l .)"
golangci-lint run
```

---

## 7. LLM 에이전트 작업 절차 (Step-by-Step Execution Workflow)

LLM 에이전트는 작업을 시작할 때 다음 단계를 순차적으로 수행합니다:

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

