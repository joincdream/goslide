# [GOS-26] Goslide 확장 지시어(Directive) 린(Lean) 진단 및 후보군 가이드 시스템 구축 계획서

> **티켓 번호**: [GOS-26](https://joincdream.atlassian.net/browse/GOS-26)  
> **마일스톤**: 마크다운 파서 신뢰성 및 작성자 경험(Authoring DX) 극대화 (Parser Diagnostics & Linting)  
> **마감일**: 2026-11-25  
> **상태**: 완료 (Done)  
> **담당 패키지**: [`internal/parser/`](../internal/parser/), [`internal/model/`](../internal/model/), [`internal/server/`](../internal/server/)  
> **아키텍처 상세 설계서**:  
> 👉 **[`docs/architecture/10-dsl-diagnostic-and-guiding-engine.md`](../docs/architecture/10-dsl-diagnostic-and-guiding-engine.md)**  
> *(상세 컴포넌트 구조, 슬라이드 블록 검사 파이프라인, LSP/VS Code 연계 메커니즘 등 세부 구현 명세는 위 아키텍처 문서를 단일 진실 공급원(SSOT)으로 참조합니다.)*  
> **연관 소스 파일**:  
> - [`internal/model/diagnostic.go`](../internal/model/diagnostic.go) *(신규)*  
> - [`internal/parser/directive.go`](../internal/parser/directive.go)  
> - [`internal/parser/parser.go`](../internal/parser/parser.go)  
> - [`internal/server/server.go`](../internal/server/server.go)  

---

## 1. 개요 및 배경 (Context & Problem Statement)

Goslide의 확장 DSL(지시어, 레이아웃, 마커)은 HTML 주석(`<!-- ... -->`) 형태로 작성됩니다. 문법 오류나 오타가 발생하면 브라우저에 에러가 노출되지 않고 조용히 기본값으로 폴백(Silent Fallback)되어 작성자가 원인을 찾는 데 많은 시간을 허비하는 문제가 있었습니다.

본 태스크는 과도한 추론 엔진(Levenshtein 거리 계산, 복잡한 i18n 템플릿 등)으로 인한 코드 복잡도를 배제하고, **"정확한 위치 지목(Location) + 정형화된 후보군 리스트(Candidate Listing)"**를 제공하는 린(Lean) 진단 시스템을 구축합니다.

### 핵심 설계 방향
1. **Fact & Location First**: 추측을 배제하고 정확한 위치(`SlideIndex`, `Line`)와 원문(`RawSnippet`), 그리고 선택 가능한 **유효 후보군(`Candidates []string`)**을 제공하여 최종 판단은 개발자가 직접 수행.
2. **Slide-Block Scope Inspection**: 변경된 단어 단위의 복잡한 증분 파싱 대신, 슬라이드 경계(`---`) 안의 **"해당 슬라이드 블록 전체"**를 단위로 검토하여 지시어 간 정합성(`_layout: two-cols` $\leftrightarrow$ `split` 마커)과 라인 오프셋을 안정적으로 보장.
3. **LSP & VS Code 친화적 모델**: 에러 메시지를 텍스트로 뭉뚱그리지 않고 정규화된 `Candidates` 슬라이스로 반환하여, 향후 VS Code의 **LSP `CodeAction` (💡 Quick Fix)** 및 **`CompletionItem` (자동 완성)**에 무개조로 1:1 직결.
4. **False-Positive Zero**: 마크다운 코드 블록 내부의 주석 및 단순 메모용 주석(`<!-- memo: ... -->`)을 철저히 격리.

---

## 2. 도메인 모델 요약 ([`internal/model/diagnostic.go`](../internal/model/diagnostic.go))

*상세 아키텍처 및 파이프라인 다이어그램은 [`docs/architecture/10-dsl-diagnostic-and-guiding-engine.md`](../docs/architecture/10-dsl-diagnostic-and-guiding-engine.md)를 참조하십시오.*

```go
package model

type DiagnosticSeverity string

const (
	SeverityError   DiagnosticSeverity = "error"   // 구문 미종료 등 본문 유실 위험
	SeverityWarning DiagnosticSeverity = "warning" // 미지원 레이아웃/지시어 오타 (기본값 폴백)
	SeverityInfo    DiagnosticSeverity = "info"    // 권장 문법 가이드
)

type Diagnostic struct {
	Severity   DiagnosticSeverity `json:"severity"`   // "error" | "warning" | "info"
	SlideIndex int                `json:"slideIndex"` // 1-based 슬라이드 번호
	Line       int                `json:"line"`       // 슬라이드 내 상대 줄 번호
	Rule       string             `json:"rule"`       // "layout.unknown", "directive.missing_underscore" 등
	Message    string             `json:"message"`    // 명확한 팩트 메시지 (예: "Unknown layout 'two-col'")
	RawSnippet string             `json:"rawSnippet"` // 문제가 발생한 원본 텍스트 ("<!-- _layout: two-col -->")
	Candidates []string           `json:"candidates"` // 선택 가능한 유효 옵션 목록 (VS Code QuickFix/Completion 연계용)
}
```

---

## 3. 단계별 구현 계획 (Implementation Steps)

| 단계 | 작업 내용 | 타겟 소스 파일 |
| :---: | :--- | :--- |
| **Step 1** | **진단 도메인 모델 정의**<br>- `model.Diagnostic` 구조체 및 `DiagnosticSeverity` 선언 | [`internal/model/diagnostic.go`](../internal/model/diagnostic.go) |
| **Step 2** | **표준 카탈로그 및 검사기(Catalog Matcher) 구현**<br>- `LayoutCatalog` (표준 레이아웃 슬라이스 관리)<br>- `DirectiveCatalog` (표준 전역/로컬 지시어 관리)<br>- 접두사 검사 기반 언더스코어 누락 진단 (`layout` $\rightarrow$ `_layout`) | [`internal/parser/catalog.go`](../internal/parser/catalog.go)<br>[`internal/parser/diagnostic.go`](../internal/parser/diagnostic.go) |
| **Step 3** | **주석 구분자 정제 및 코드 블록 AST 격리**<br>- `--->` 자동 정제 및 미종료 주석 본문 삼킴(Swallowing) 방지<br>- Goldmark AST 순회 시 코드 블록 내부 주석 스킵 | [`internal/parser/directive.go`](../internal/parser/directive.go) |
| **Step 4** | **슬라이드 블록 단위 진단 수집 파이프라인 결합**<br>- `splitSlides()` 후 각 슬라이드 블록별로 진단 수집<br>- 슬라이드 내 상대 라인 번호(`Line`) 및 `SlideIndex` 바인딩 | [`internal/parser/parser.go`](../internal/parser/parser.go) |
| **Step 5** | **CLI 및 개발 서버(`goslide serve`) 출력 연동**<br>- 터미널 경고 시 원문(`RawSnippet`)과 `Available: [...]` 후보군 출력<br>- SSE 이벤트 파이프라인을 통한 웹 브라우저 경고 전달 | [`internal/server/server.go`](../internal/server/server.go)<br>[`cmd/goslide/`](../cmd/goslide/) |
| **Step 6** | **단위 테스트 및 무결성 검증**<br>- 슬라이드 블록 검사, 후보군 리스트 정합성, 코드블록 오탐 방지 테스트 | [`internal/parser/diagnostic_test.go`](../internal/parser/diagnostic_test.go) |

---

## 4. 수용 기준 (Acceptance Criteria)

- [x] **후보군 리스트 반환**: 미지원 레이아웃(예: `_layout: two-col`) 입력 시 `Candidates` 슬라이스에 `[default, cover, two-cols, split, section, blank]`가 정상 반환되어야 함.
- [x] **언더스코어 누락 검출**: 슬라이드 로컬 본문에서 `<!-- layout: two-cols -->` 작성 시 `directive.missing_underscore` 경고와 함께 `Candidates: ["_layout"]`가 제시되어야 함.
- [x] **슬라이드 블록 단위 정확성**: 슬라이드 경계(`---`)를 기준으로 상대 라인 번호(`Line`)와 슬라이드 인덱스(`SlideIndex`)가 정확히 일치해야 함.
- [x] **코드 블록 오탐 방지 (False-Positive Zero)**: 마크다운 코드 블록(```html 등) 내부에 작성된 예시 주석에 대해서는 어떠한 진단 경고도 발생하지 않아야 함.
- [x] **주석 구분자 정제**: `--->` 입력 시 값이 오염되지 않고 정상 정제되어 슬라이드가 렌더링되어야 함.
- [x] **성능 및 빌드 무결성**: 단일 바이너리(CGO 0%), 50장 슬라이드 기준 파싱/진단 지연 2ms 이하, 기존 렌더링 결과에 일체의 부작용이 없어야 함.
