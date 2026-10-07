# [GOS-3] M1-2: 도메인 IR 및 코어 마크다운 파서 개발 계획서

> **티켓 번호**: [GOS-3](https://joincdream.atlassian.net/browse/GOS-3)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-2  
> **마감일**: 2026-10-09  
> **상태**: 완료 (Completed)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [architecture/index.md](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/index.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크는 마크다운 원문을 읽어 **Goslide의 표준 불변 도메인 모델(Presentation Deck IR)**로 변환하는 코어 파서 파이프라인을 구축하는 작업입니다.

단순한 기능 구현을 넘어, 본 작업이 완료되었을 때 슬라이드 작성자와 개발자에게 제공되는 실질적 가치는 다음과 같습니다:

### 1.1 슬라이드 작성자(End-User)에게 제공하는 가치
1. **기술 문서 작성의 안정성 (Code-Safe Splitting)**:
   - 기존 마크다운 슬라이드 도구의 고질적인 문제였던 **코드 블록(YAML, CI 스크립트, git diff 등) 내부의 `---` 기호로 인해 슬라이드가 의도치 않게 잘려나가는 문제를 원천 차단**합니다.
   - 엔지니어는 코드가 포함된 마크다운을 아무런 탈출(Escape) 트릭 없이 자연스럽게 작성할 수 있습니다.
2. **무설정(Zero-Config) 기본 사용성**:
   - 상단 Frontmatter를 작성하지 않아도 기본 테마(`default`), 16:9 비율, 표준 레이아웃이 자동으로 적용되어, 평범한 마크다운 파일도 즉시 슬라이드로 변환됩니다.
3. **표준 GFM(GitHub Flavored Markdown) 표현력 제공**:
   - 복잡한 표(Table), 작업 체크리스트(`- [x]`), 취소선(`~~text~~`), 자동 URL 링크가 깨짐 없이 브라우저 슬라이드에 렌더링됩니다.

### 1.2 렌더러/익스포터 개발자(System/Developer)에게 제공하는 가치
1. **단일 진실 원천(Single Source of Truth) 확보**:
   - 이후 단계인 HTML 렌더러(M1-5), 테마 엔진(M1-4), PDF/PPTX 익스포터(M2-1, M2-2)는 복잡한 원시 문자열을 다시 해석할 필요 없이, 완전히 정제된 불변 구조체(`model.Deck`, `model.Slide`)만 참조하면 됩니다.
2. **데이터 레이스 없는 불변성(Immutability)**:
   - 파싱이 완료된 `Deck` 객체는 읽기 전용으로 안전하게 공유되어 멀티스레드 렌더링 및 동시 변환 시에도 데이터 경합이 발생하지 않습니다.

---

## 2. 슬라이드 구조체 아키텍처 및 컴포넌트 결합 환경 (Slide Structural Architecture)

본 태스크의 핵심 결과물은 **슬라이드 단위의 구조체([`model.Deck`](file:///home/yundream/myjob/cloit/Goslide/internal/model/deck.go#L25) / [`model.Slide`](file:///home/yundream/myjob/cloit/Goslide/internal/model/deck.go#L57))**를 확립하고, 이후 렌더링 파이프라인의 컴포넌트들이 결합될 수 있는 환경을 제공하는 것입니다.

### 2.1 슬라이드 구조체와 컴포넌트 결합 다이어그램

```
       [GOS-3에서 구축하는 슬라이드 구조체 (model.Slide)]
       ┌────────────────────────────────────────────────────────┐
       │ Index / Layout / HTMLContent / Directives / Notes      │
       └────────────────────────────────────────────────────────┘
          ▲               ▲                    ▲             ▲
          │               │                    │             │
    [테마 컴포넌트]   [레이아웃 컴포넌트]    [JS 런타임]   [익스포터]
       (GOS-5)         (GOS-4)              (GOS-6)       (GOS-8/9)
     16:9 박스 규격,   커버 중앙 정렬,       키보드 넘김,   슬라이드별
     폰트, 색상 CSS     좌/우 2단 컬럼 분할    투명 캔버스 판서  PDF/PPTX 캡처
```

### 2.2 도메인 모델 상세 구조 명세 (`internal/model/deck.go`)

```go
package model

import "time"

// SizeRatio defines the presentation aspect ratio.
type SizeRatio string
const (
	Ratio16x9 SizeRatio = "16:9"
	Ratio4x3  SizeRatio = "4:3"
)

// LayoutType defines the structural layout archetype of a slide.
type LayoutType string
const (
	LayoutDefault LayoutType = "default"  // 기본 본문 슬라이드
	LayoutCover   LayoutType = "cover"    // 표지/중앙 정렬 타이틀 슬라이드
	LayoutSection LayoutType = "section"  // 중간 챕터 구분 슬라이드
	LayoutTwoCols LayoutType = "two-cols" // 좌/우 2단 분할 슬라이드
	LayoutBlank   LayoutType = "blank"    // 여백 없는 백지 슬라이드
)

// Deck represents the entire presentation container.
type Deck struct {
	Title       string           `json:"title" yaml:"title"`
	Author      string           `json:"author" yaml:"author"`
	CreatedAt   time.Time        `json:"created_at" yaml:"created_at"`
	GlobalAttrs GlobalDirectives `json:"global_attributes" yaml:",inline"`
	CustomCSS   string           `json:"custom_css,omitempty" yaml:"custom_css"`
	Slides      []*Slide         `json:"slides" yaml:"-"`
}

// GlobalDirectives holds deck-level settings extracted from Frontmatter.
type GlobalDirectives struct {
	Theme    string     `json:"theme" yaml:"theme"`
	Layout   LayoutType `json:"layout" yaml:"layout"`
	Size     SizeRatio  `json:"size" yaml:"size"`
	Paginate bool       `json:"paginate" yaml:"paginate"`
	Header   string     `json:"header,omitempty" yaml:"header"`
	Footer   string     `json:"footer,omitempty" yaml:"footer"`
}

// SlideDirectives holds scoped directives applied to a single slide.
type SlideDirectives struct {
	Layout          LayoutType `json:"layout,omitempty"`
	Class           []string   `json:"class,omitempty"`
	BackgroundColor string     `json:"background_color,omitempty"`
	BackgroundImage string     `json:"background_image,omitempty"`
	Color           string     `json:"color,omitempty"`
	Header          string     `json:"header,omitempty"`
	Footer          string     `json:"footer,omitempty"`
	Paginate        bool       `json:"paginate"`
}

// Slide represents an atomic presentation slide slot for component binding.
type Slide struct {
	Index       int             `json:"index"`                  // 1-based 슬라이드 고유 순번
	Layout      LayoutType      `json:"layout"`                 // 슬라이드 레이아웃 템플릿
	Directives  SlideDirectives `json:"directives"`             // 슬라이드 개별 스타일 속성
	RawContent  string          `json:"-"`                      // 파싱 전 원시 마크다운
	HTMLContent string          `json:"html_content"`           // 파싱된 본문 HTML 조각
	Notes       string          `json:"notes,omitempty"`        // 발표자 전용 메모
	LeftHTML    string          `json:"left_html,omitempty"`    // 2단 레이아웃 좌측 본문
	RightHTML   string          `json:"right_html,omitempty"`   // 2단 레이아웃 우측 본문
}
```

---

## 3. 핵심 설계 결정 사항 (Key Architectural Decisions)

엔지니어링 모호성과 추측을 방지하기 위해 확정된 3대 설계 결정입니다:

### 3.1 결정 1: 슬라이드 필드 기본 초기화 및 GOS-4 위임
- **GOS-3 책임 범위**: `model.Slide` 구조체 생성 시 `Layout`은 `LayoutDefault`, `Notes`는 `""`(빈 문자열), `Directives`는 빈 구조체, `LeftHTML`/`RightHTML`은 `""`로 **안전하게 기본 초기화**만 수행합니다.
- **GOS-4 위임**: 마크다운 인라인 HTML 주석(`<!-- layout: cover -->`, `<!-- note: ... -->`)의 실제 추출과 2단 분할(`LeftHTML/RightHTML`) 로직은 다음 마일스톤인 **GOS-4(M1-3)**에서 전담합니다. GOS-3에서는 주석 파싱을 수행하지 않고 Goldmark의 일반 파싱 파이프라인에 그대로 넘깁니다.

### 3.2 결정 2: 슬라이드 분할 규칙 (Strict 3 Dashes & Empty Skip)
- **엄격한 3개 대시 규칙 (`^---[ \t]*$`)**: 오직 정확히 대시 3개로만 구성된 줄만 슬라이드 분할자로 인식합니다.
  - `----` (4개 이상 대시), `***`, `___`는 슬라이드 분할자가 아닌 마크다운 본문의 구분선(`<hr>`)으로 온전히 보존하여 슬라이드 내부 콘텐츠 표현력을 해치지 않습니다.
- **빈 슬라이드 자동 스킵**: 슬라이드 분할자 사이나 문서 끝에 공백만 존재하는 빈 슬라이드는 슬라이드 목록(`Deck.Slides`)에서 자동 제외하여 불필요한 빈 화면 렌더링을 차단합니다.
- **중첩 Fenced Code Block 완벽 보존**:
  - 열린 백틱(```) 또는 물결표(~~~)의 개수($N \ge 3$)를 정확히 추적하여, 동일 문자 $N$개 이상의 닫는 펜스를 만날 때까지 내부의 `---`는 슬라이드 분할자로 처리하지 않습니다.

### 3.3 결정 3: Frontmatter 미지원 값의 Graceful Fallback
- **YAML 문법 오류**: YAML 파싱 자체가 불가능한 문법 오류(들여쓰기 깨짐 등)인 경우에만 명시적으로 `model.ErrInvalidFrontmatter` 에러를 반환합니다.
- **미지원 테마 및 해상도 기본값 폴백**:
  - `theme`: 미입력 또는 알 수 없는 테마가 입력된 경우 기본 테마인 `"default"`로 자동 폴백.
  - `size`: `"16:9"`, `"4:3"` 이외의 값이 입력된 경우 기본 해상도인 `model.Ratio16x9`(`"16:9"`)로 자동 폴백.
  - `layout`: 미지원 레이아웃은 `model.LayoutDefault`(`"default"`)로 폴백.
- **슬라이드 본문 누락 방어**: Frontmatter만 존재하고 슬라이드 본문이 0개인 경우, 빈 1장의 default 슬라이드를 생성하여 후속 렌더러가 정상 구동되도록 보장합니다.

---

## 4. 결과물의 구체적 모습 (Concrete Manifestation)

사용자가 입력한 마크다운이 파서를 거쳐 어떤 중간 도메인 모델과 HTML 조각으로 변환되는지 구체적인 데이터 흐름으로 나타냅니다.

### 4.1 마크다운 원문 입력 예시 (`presentation.md`)

```markdown
---
title: "Goslide 아키텍처 개요"
author: "윤상배"
theme: "default"
size: "16:9"
paginate: true
header: "Goslide 내부 기술 세미나"
footer: "© 2026 Cloit Corp."
---

# Goslide 프레젠테이션 엔진

Go 언어로 구현된 고성능 슬라이드 데크 빌더입니다.

---

# 개발 가이드 및 코드 예시

아래의 YAML 설정 파일 예제는 내부 `---`를 포함하지만, 정상적으로 한 장의 슬라이드로 보존됩니다:

```yaml
---
server:
  port: 8080
---
```

| 모듈 | 상태 | 지원 포맷 |
| :--- | :---: | :--- |
| 파서 | **완료** | Frontmatter, GFM, 수평선 분할 |
| 렌더러 | 진행 중 | HTML, PDF, PPTX |
```

### 4.2 파싱 후 불변 도메인 모델(`*model.Deck`) 인메모리 모습

파서가 생성하여 후속 렌더러에 전달하는 Go 인메모리 도메인 객체의 모습입니다:

```go
&model.Deck{
    Title:     "Goslide 아키텍처 개요",
    Author:    "윤상배",
    CreatedAt: time.Now(),
    GlobalAttrs: model.GlobalDirectives{
        Theme:    "default",
        Layout:   model.LayoutDefault, // "default"
        Size:     model.Ratio16x9,      // "16:9"
        Paginate: true,
        Header:   "Goslide 내부 기술 세미나",
        Footer:   "© 2026 Cloit Corp.",
    },
    Slides: []*model.Slide{
        // [슬라이드 1] 타이틀 페이지
        {
            Index:      1,
            Layout:     model.LayoutDefault,
            Directives: model.SlideDirectives{},
            RawContent: "# Goslide 프레젠테이션 엔진\n\nGo 언어로 구현된 고성능...",
            HTMLContent: "<h1>Goslide 프레젠테이션 엔진</h1>\n<p>Go 언어로 구현된 고성능 슬라이드 데크 빌더입니다.</p>\n",
        },
        // [슬라이드 2] 코드 블록 및 GFM 표가 포함된 본문 슬라이드
        {
            Index:      2,
            Layout:     model.LayoutDefault,
            Directives: model.SlideDirectives{},
            RawContent: "# 개발 가이드 및 코드 예시\n\n아래의 YAML 설정 파일...",
            HTMLContent: "<h1>개발 가이드 및 코드 예시</h1>\n" +
                         "<p>아래의 YAML 설정 파일 예제는 내부 <code>---</code>를 포함하지만, 정상적으로 한 장의 슬라이드로 보존됩니다:</p>\n" +
                         "<pre><code class=\"language-yaml\">---\nserver:\n  port: 8080\n---\n</code></pre>\n" +
                         "<table>\n<thead>\n<tr><th>모듈</th><th align=\"center\">상태</th><th>지원 포맷</th></tr>\n</thead>\n" +
                         "<tbody>\n<tr><td>파서</td><td align=\"center\"><strong>완료</strong></td><td>Frontmatter, GFM, 수평선 분할</td></tr>\n" +
                         "<tr><td>렌더러</td><td align=\"center\">진행 중</td><td>HTML, PDF, PPTX</td></tr>\n</tbody>\n</table>\n",
        },
    },
}
```

---

## 5. 사용성 및 사용자 경험 (Usability & Developer Experience)

### 5.1 슬라이드 작성자 사용성 (Author Experience)
1. **대시 3개(`---`) 기반 명확한 분할**:
   - `---` 앞뒤 빈 줄 유무와 상관없이 명확히 슬라이드를 분할합니다.
   - 단, `----`나 `***`는 본문 구분선으로 남아 슬라이드가 쪼개지지 않으므로 안심하고 서식을 지정할 수 있습니다.
2. **명확한 문법 에러 피드백 (Actionable Error Reporting)**:
   - Frontmatter YAML 구문 오류 시 줄 번호와 원인이 포함된 에러를 제공합니다:
     ```
     error: invalid frontmatter syntax in presentation.md: yaml: line 4: mapping values are not allowed in this context
     ```

### 5.2 Go 개발자 API 사용성 (Library Developer Experience)
```go
// 1. 파서 인스턴스 생성
p := parser.NewParser()

// 2. 표준 io.Reader로부터 Deck 파싱
deck, err := p.Parse(ctx, strings.NewReader(markdownContent))
if err != nil {
    if errors.Is(err, model.ErrInvalidFrontmatter) {
        // Frontmatter 구문 오류 전용 핸들링
    }
    return err
}

// 3. 정제된 불변 Deck 바로 활용
fmt.Printf("총 슬라이드 수: %d, 테마: %s, 비율: %s\n", len(deck.Slides), deck.GlobalAttrs.Theme, deck.GlobalAttrs.Size)
```

---

## 6. 세부 개발 태스크 (Work Breakdown Structure)

### 6.1 순수 Go 의존성 추가 (`go.mod`)
- `gopkg.in/yaml.v3`: Frontmatter YAML 메타데이터 언마샬링
- `github.com/yuin/goldmark`: 고성능 Pure Go 마크다운 AST 파서 및 GFM 확장
- CGO 의존성 0% 유지 (`CGO_ENABLED=0`)

### 6.2 도메인 모델 태그 보완 (`internal/model/deck.go`)
- `Deck` 및 `GlobalDirectives` 구조체에 YAML 태그(`yaml:"title"`, `yaml:"theme"` 등) 보강

### 6.3 Frontmatter 추출기 구현 (`internal/parser/frontmatter.go`)
- 입력 시작부 `---` 감지 및 닫는 `---`까지의 YAML 블록 추출
- Frontmatter 미존재 시 기본 전역 속성(`theme: default`, `size: 16:9`, `layout: default`, `paginate: false`) 주입
- `yaml.Unmarshal` 수행:
  - 구문 오류 시 `model.ErrInvalidFrontmatter` 반환
  - 파싱 후 `Theme`이 비어있거나 미지원이면 `"default"` 폴백
  - `Size`가 `"16:9"`, `"4:3"`이 아니면 `Ratio16x9` 폴백
  - `Layout`이 미지원이면 `LayoutDefault` 폴백

### 6.4 슬라이드 분할기 구현 (`internal/parser/splitter.go`)
- 개행 문자 일괄 정규화 (`\r\n` $\rightarrow$ `\n`)
- 줄(Line) 단위 정밀 스캐너:
  - Fenced Code Block 상태 머신: 열린 펜스 문자(``` 또는 ~~~)와 개수($N \ge 3$) 추적
  - 코드 블록 밖에서만 `^---[ \t]*$` 패턴을 분할자로 판별
  - `----` 및 `***`는 일반 텍스트로 보존
- 공백만 있는 빈 슬라이드 청크 스킵
- 슬라이드가 아예 없는 경우 1장의 빈 슬라이드 기본 생성
- 1부터 시작하는 순차 인덱스(`Index: 1, 2, ...`) 부여

### 6.5 Goldmark GFM 파서 파이프라인 (`internal/parser/parser.go`)
- `model.Parser` 인터페이스 구현체 `Parser` 구조체:
  - `goldmark.New` 시 `extension.GFM`, `parser.WithAutoHeadingID()`, `html.WithUnsafe()` 활성화
- `Parse(ctx context.Context, r io.Reader) (*model.Deck, error)`:
  - `ctx.Err()` 취소 신호 사전 확인
  - Frontmatter 추출 $\rightarrow$ 슬라이드 분할 $\rightarrow$ 각 슬라이드별 마크다운 $\rightarrow$ HTML 본문 변환
  - `Slide.Layout = model.LayoutDefault`, `Slide.Notes = ""`, `Slide.Directives = model.SlideDirectives{}` 기본 초기화
  - 완성된 불변 `*model.Deck` 반환

### 6.6 단위 테스트 스위트 작성 (`internal/parser/*_test.go`)
- **테이블 기반 테스트 25개 이상 작성**:
  1. `frontmatter_test.go`:
     - 정상 Frontmatter 파싱
     - Frontmatter 생략 시 기본값(`default`, `16:9`) 적용 검증
     - 미지원 `size: "32:9"` 입력 시 `16:9` 기본값 자동 폴백 검증
     - 미지원 `theme: "unknown"` 입력 시 `"default"` 자동 폴백 검증
     - 잘못된 YAML 문법 시 `ErrInvalidFrontmatter` 반환 검증
  2. `splitter_test.go`:
     - 정확히 `---` 3개로 분할되는지 검증
     - `----` (대시 4개) 및 `***`가 슬라이드로 쪼개지지 않고 본문으로 남는지 검증
     - 4개 이상 중첩 백틱(```` ````) 내부의 `---` 보존 검증
     - 물결표(~~~) 코드 블록 내부의 `---` 보존 검증
     - 공백만 있는 빈 슬라이드 자동 스킵 검증
     - 슬라이드 본문이 없는 경우 1장 기본 생성 검증
  3. `parser_test.go`:
     - 전체 파이프라인 통합 변환 (GFM 테이블, 체크리스트, 코드 블록 HTML 검증)
     - `Slide.Layout`, `Slide.Notes` 기본 초기화 값 검증
     - `context.WithCancel` 취소 시 `ErrCanceled` 즉시 반환 검증

---

## 7. 완료 기준 (Definition of Done)

- [x] `go.mod`에 `gopkg.in/yaml.v3` 및 `github.com/yuin/goldmark`가 등록되고 빌드 성공.
- [x] `internal/model/deck.go`에 YAML 직렬화 태그가 보강됨.
- [x] `internal/parser/frontmatter.go`가 구현되고 미지원 테마/해상도 기본값 폴백 검증 완료.
- [x] `internal/parser/splitter.go`가 구현되고 대시 3개 엄격 분할, 코드 블록 중첩 보존, 빈 슬라이드 스킵이 검증됨.
- [x] `internal/parser/parser.go`가 `model.Parser`를 구현하고 GFM 문법을 HTML로 변환함.
- [x] 25개 이상의 테이블 기반 단위 테스트가 작성되고 `go test -v -race ./...` 통과.
- [x] `make lint` 및 `make complexity` 실행 시 경고 0건 및 복잡도 임계치 만족.

---

## 8. 실행 절차 (Step-by-Step Execution Plan)

1. **모듈 의존성 설정**: `go get gopkg.in/yaml.v3 github.com/yuin/goldmark`
2. **도메인 모델 태그 보완**: `internal/model/deck.go` 수정
3. **Frontmatter 파서 구현 및 TDD 테스트**: `frontmatter.go` 및 `frontmatter_test.go` (폴백 로직 포함)
4. **슬라이드 분할기 구현 및 TDD 테스트**: `splitter.go` 및 `splitter_test.go` (대시 3개 엄격 검사, 중첩 펜스 추적, 빈 슬라이드 스킵)
5. **Goldmark 코어 파서 구현 및 통합 테스트**: `parser.go` 및 `parser_test.go`
6. **품질 검증**: `make lint`, `make complexity`, `make test-race`
7. **Jira 업데이트**: 완료 코멘트 작성 및 커밋
