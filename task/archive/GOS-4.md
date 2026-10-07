# [GOS-4] M1-3: 지시어 파서, 명시적 레이아웃 엔진 및 조기 프리뷰 빌드 계획서

> **티켓 번호**: [GOS-4](https://joincdream.atlassian.net/browse/GOS-4)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-3  
> **마감일**: 2026-10-14  
> **상태**: 완료 (Completed)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [core-architecture.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/core-architecture.md), [decisions-simplification.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 슬라이드별 서식을 세밀하게 제어하는 **주석 지시어(Directive) 파서**와 **명시적 시맨틱 레이아웃 엔진**을 구축하고, 순수 Go 기반 **코드 구문 강조(Chroma)**를 연동하는 것입니다.  
특히 빅뱅 검증 리스크를 방지하기 위해 **CLI 빌드(`goslide build`) 뼈대를 조기에 연결**하여, 본 태스크 완료 즉시 사용자가 브라우저에서 직접 슬라이드 결과물을 확인할 수 있는 워킹 스켈레톤(Walking Skeleton)을 제공합니다.

### 1.1 슬라이드 작성자(End-User)에게 제공하는 가치
1. **슬라이드별 맞춤 스타일링 (Scoped Directives)**:
   - 마크다운 주석(`<!-- _backgroundColor: #1e293b -->`, `<!-- _class: lead -->`, `<!-- _color: white -->`)을 통해 특정 슬라이드만 어두운 배경이나 강조 스타일로 손쉽게 전환할 수 있습니다.
2. **2단 병렬 컬럼 비교 레이아웃 (`<!-- split -->`)**:
   - `<!-- _layout: two-cols -->` 선언과 `<!-- split -->` 구분자만으로 복잡한 HTML/CSS 없이 좌/우 2단 그리드(비교표, 이미지+설명 등)를 직관적으로 구성합니다.
3. **발표자 전용 메모 작성 (`<!-- note: ... -->`)**:
   - 청중에게 보여주지 않고 발표자 콘솔이나 PPTX 슬라이드 메모에만 기록할 발표 대본을 마크다운 주석으로 자유롭게 작성합니다.
4. **미려한 소스 코드 구문 강조 (Pure Go Chroma)**:
   - Go, Python, JavaScript, YAML, JSON, Bash 등 다양한 언어의 코드 블록이 순수 Go 환경에서 CGO 없이 미려한 컬러로 하이라이팅됩니다.
5. **조기 실행 및 브라우저 프리뷰 체감 (Walking Skeleton)**:
   - 3주 차 최종 출시(GOS-7)까지 기다릴 필요 없이, **본 작업 직후 `./bin/goslide build talk.md -o preview.html`을 터미널에서 직접 실행하여 브라우저에서 슬라이드 경계선과 2단 레이아웃을 즉시 확인**할 수 있습니다.

---

## 2. 핵심 설계 원칙 및 규칙 (Architectural Rules)

### 2.1 결정론적 레이아웃 정책 (ADR-001: Zero Guesswork)
[ADR-001](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md)에 따라, 파서는 슬라이드의 헤딩 개수나 줄 수를 자의적으로 추측하는 스마트 추론을 배제하고 **명시적 선언 100% 매핑 원칙**을 엄격히 준수합니다:
1. **1순위 (로컬 오버라이드)**: 슬라이드 내 `<!-- _layout: ... -->`가 선언된 경우 해당 레이아웃 적용.
2. **2순위 (상속 지시어)**: 이전 슬라이드에서 `<!-- layout: ... -->`로 상속된 값이 있는 경우 적용.
3. **3순위 (데크 기본값)**: Frontmatter에 선언된 전역 `layout` 적용.
4. **기본 Fallback**: 아무 선언도 없으면 무조건 표준 레이아웃인 **`default`** 적용.

### 2.2 지시어 스코프 규칙 (Marp 100% 호환)
- **언더스코어 접두사(`_`) 있음**: 해당 슬라이드 1장에만 적용되는 **일회성 로컬 지시어** (예: `<!-- _class: lead -->`, `<!-- _backgroundColor: #000 -->`).
- **언더스코어 접두사(`_`) 없음**: 해당 슬라이드부터 문서 끝까지 유지되는 **상속형 전역 지시어** (예: `<!-- header: "2장 아키텍처" -->`).

### 2.3 2단 컬럼 분할 (`<!-- split -->`)
- `Layout == "two-cols"`인 슬라이드 본문에서 `<!-- split -->` 주석이 감지되면, 앞부분은 `Slide.LeftHTML`, 뒷부분은 `Slide.RightHTML`로 분리하여 렌더링합니다.

---

## 3. 결과물의 구체적 모습 (Concrete Manifestation)

### 3.1 마크다운 원문 입력 예시 (`demo.md`)

```markdown
---
title: "Goslide 아키텍처 및 코드 데모"
theme: "default"
size: "16:9"
---

<!-- _layout: cover -->
<!-- _backgroundColor: #0f172a -->
<!-- _color: #f8fafc -->

# Goslide 프레젠테이션 엔진
### 고성능 순수 Go 마크다운 슬라이드 빌더

<!-- note:
이 슬라이드에서는 단일 바이너리와 고성능 철학을 먼저 강조할 것.
-->

---

<!-- _layout: two-cols -->
## 레거시 vs 신규 아키텍처 비교

### 기존 구조 (Node.js)
- 수백 MB의 무거운 런타임
- 느린 cold start

<!-- split -->

### 신규 구조 (Pure Go)
- 10MB 단일 바이너리
- 밀리초 단위 즉시 빌드

---

## 코드 구문 강조 (Go)

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Goslide!")
}
```
```

### 3.2 CLI 실행 및 생성되는 프리뷰 HTML 모습

사용자가 터미널에서 다음 명령을 실행합니다:
```bash
./bin/goslide build demo.md -o preview.html
```

생성된 `preview.html`을 더블클릭하여 브라우저에서 열면:
1. **1번 슬라이드**: 어두운 배경색(`#0f172a`)에 흰색 글씨, 표지 스타일로 중앙 정렬된 타이틀 표시.
2. **2번 슬라이드**: `<!-- split -->`을 기준으로 좌측(기존 구조)과 우측(신규 구조)이 깔끔하게 2단 그리드로 나뉘어 배치.
3. **3번 슬라이드**: Go 소스 코드가 Chroma를 통해 형형색색의 구문 강조(Syntax Highlighting) 스팬으로 렌더링.
4. **발표자 노트**: 슬라이드 화면에는 노출되지 않고 내부 데이터 모델에 안전 보존.

---

## 4. 세부 개발 태스크 (Work Breakdown Structure)

### 4.1 Pure Go 의존성 추가 (`go.mod`)
- `github.com/alecthomas/chroma/v2`: 순수 Go 구문 강조 라이브러리
- CGO 배제 (`CGO_ENABLED=0`) 준수

### 4.2 인라인 지시어 및 발표자 노트 토크나이저 (`internal/parser/directive.go`)
- 정규식 또는 줄 스캐너를 통해 `<!-- ... -->` 주석 토큰 추출:
  - 발표자 노트: `<!-- note: ... -->` 및 다중 행 주석 블록 $\rightarrow$ `Slide.Notes`에 적재 (HTML 렌더링 본문에서는 제거)
  - 지시어 매핑: `class`, `backgroundColor`, `backgroundImage`, `color`, `header`, `footer`, `paginate`, `layout`
  - 언더스코어(`_`) 접두사 유무에 따른 로컬 vs 상속 상태 관리자(`directiveState`) 구현

### 4.3 2단 컬럼 분할 처리 (`internal/parser/directive.go`)
- 슬라이드 본문에서 `<!-- split -->` 구분자 주석 감지:
  - `Layout == LayoutTwoCols`일 때 분할점 이전 내용을 `LeftHTML`, 이후 내용을 `RightHTML`로 각각 Goldmark 파싱하여 저장
  - `Slide.HTMLContent`는 좌/우 컬럼을 감싸는 기본 2단 컨테이너 구조로 합성

### 4.4 Chroma 순수 Go 구문 강조 연동 (`internal/parser/highlight.go`)
- Goldmark 커스텀 렌더러 또는 훅을 작성하여 Fenced Code Block을 감지:
  - 지정된 언어(`go`, `python`, `json`, `yaml`, `bash` 등)에 맞는 Chroma Lexer 탐색 (미지정 시 Fallback Lexer)
  - Pure Go 스타일(기본 다크 팔레트 등)로 HTML `<pre><code>` 내부에 인라인 스타일 스팬 태그 생성

### 4.5 파서 파이프라인 통합 (`internal/parser/parser.go`)
- `splitSlides()` 이후 각 슬라이드 청크에 대해:
  1. 지시어 추출 및 상속 상태 전이
  2. 발표자 노트 추출 및 본문에서 분리
  3. Chroma 코드 구문 강조와 함께 마크다운 $\rightarrow$ HTML 변환
  4. 2단 분할 처리(`LeftHTML`, `RightHTML`)
  5. `model.Slide`에 계산된 속성(`Directives`, `Layout`, `Notes`) 주입

### 4.6 조기 프리뷰 빌드 CLI 연동 (`cmd/goslide/build.go`)
- `cmd/goslide/build.go` 작성:
  - `goslide build <input.md> -o <output.html>` 플래그 파싱
  - 입력 파일 읽기 $\rightarrow$ `parser.Parse()` 호출
  - 미니멀 프리뷰 템플릿(슬라이드별 `<section>` 테두리 박스 및 2단 그리드 CSS 기본 내장)을 통해 브라우저에서 즉시 열람 가능한 `preview.html` 파일 생성

### 4.7 테이블 기반 단위 테스트 스위트 (`internal/parser/*_test.go`)
- `directive_test.go`: 로컬 지시어(`_`), 상속 지시어, 발표자 노트 추출, 다중 행 노트 검증 (15+ 케이스)
- `layout_test.go`: 2단 분할(`<!-- split -->`), 명시적 레이아웃 매핑 검증 (5+ 케이스)
- `highlight_test.go`: Go, Python, JSON 코드 블록 구문 강조 HTML 스팬 검증 (5+ 케이스)
- `cmd/goslide/build_test.go`: `goslide build` CLI 파일 생성 E2E 테스트

---

## 5. 엔지니어링 가드레일 (Hard Constraints)

1. **CGO 배제 (100% Pure Go)**: `alecthomas/chroma/v2` 활용, 외부 의존성 최소화.
2. **결정론적 출력 (Zero Guesswork)**: 파서가 콘텐츠를 자의적으로 추측하여 레이아웃을 바꾸지 않음.
3. **순환 참조 방지 계층 엄수**: `internal/parser`는 `internal/model`에만 의존.
4. **`panic()` 절대 금지 및 컨텍스트 취소 보장**: 모든 에러는 명시적 래핑 반환.

---

## 6. 완료 기준 (Definition of Done)

- [x] `go.mod`에 `github.com/alecthomas/chroma/v2`가 등록되고 빌드 통과.
- [x] `internal/parser/directive.go`가 구현되어 로컬/상속 지시어 및 `Notes`가 정상 추출됨.
- [x] `<!-- split -->`을 통한 좌/우 2단 컬럼 분할 및 HTML 생성이 검증됨.
- [x] `internal/parser/highlight.go`가 연동되어 주요 언어 코드 블록이 구문 강조 HTML로 변환됨.
- [x] `./bin/goslide build sample.md -o preview.html` 명령으로 사용자가 브라우저에서 직접 슬라이드를 확인할 수 있는 프리뷰 HTML 파일이 생성됨.
- [x] 단위 테스트가 모두 작성되고 `go test -v -race ./...`를 100% 통과함.
- [x] `make lint` 및 `make complexity` 기준을 무결하게 만족함.

---

## 7. 실행 절차 (Step-by-Step Execution Plan)

1. **의존성 추가**: `go get github.com/alecthomas/chroma/v2`
2. **지시어 파서 및 테스트**: `directive.go`, `directive_test.go` (스코프 및 노트 추출)
3. **Chroma 구문 강조 및 테스트**: `highlight.go`, `highlight_test.go`
4. **파서 파이프라인 통합**: `parser.go`에 지시어, 2단 분할, 하이라이팅 결합
5. **조기 프리뷰 CLI 구현**: `cmd/goslide/build.go` 작성 및 `preview.html` 생성 연동
6. **품질 검증**: `make lint`, `make complexity`, `make test-race`, `make build`
7. **사용자 확인**: 프리뷰 실행 명령어 안내 및 사용자 검토 후 승인
