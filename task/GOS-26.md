# [GOS-26] Goslide 확장 지시어(Directive) 문법 오류 및 오타 진단(Diagnostic/Lint) 시스템 구축 계획서

> **티켓 번호**: [GOS-26](https://joincdream.atlassian.net/browse/GOS-26)  
> **마일스톤**: 마크다운 파서 신뢰성 및 작성자 경험(Authoring DX) 극대화 (Parser Diagnostics & Linting)  
> **마감일**: 2026-11-25  
> **상태**: 할 일 (To Do)  
> **담당 패키지**: `internal/parser/`, `internal/model/`, `internal/server/`, `internal/i18n/`  
> **설계 철학**:  
> - **Fail-Fast & No Silent Fallback**: 알 수 없는 레이아웃이나 오타를 조용히 기본값으로 덮지 않고 즉각 명확한 진단(Diagnostic) 피드백 제공  
> - **Actionable Feedback ("Did you mean?")**: 단순 에러 출력을 넘어 오타 가능성이 높은 표준 키워드를 추천하여 작성 시간 단축  
> - **i18n First**: 모든 경고 및 진단 메시지의 다국어(ko/en) 카탈로그 바인딩  
> **연관 소스 파일**:  
> - [`internal/parser/directive.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive.go)  
> - [`internal/parser/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/layout.go)  
> - [`internal/parser/parser.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/parser.go)  
> - [`internal/model/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/layout.go)  
> - [`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go)  
> - [`internal/i18n/locales/ko.json`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/ko.json)  
> - [`internal/i18n/locales/en.json`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/en.json)  
> - [`task/GOS-19.md`](file:///home/yundream/myjob/cloit/Goslide/task/GOS-19.md)  
> - [`task/GOS-24.md`](file:///home/yundream/myjob/cloit/Goslide/task/GOS-24.md)  

---

## 1. 개요 및 배경 (Context & Problem Statement)

### 1.1 사각지대 분석 (The Invisible Failure Problem)

표준 마크다운 문법 오류는 브라우저 화면에 마크다운 기호가 그대로 노출되어 작성자가 즉시 눈으로 인지할 수 있습니다. 반면, Goslide의 확장 DSL(지시어, 마커)은 **HTML 주석(`<!-- ... -->`) 형태**로 작성됩니다.

이로 인해 다음과 같은 치명적인 **"보이지 않는 실패(Invisible Failure)"** 사각지대가 발생합니다:

1. **HTML 주석의 시각적 은닉성**:
   - 주석은 브라우저에서 렌더링되지 않고 완전히 숨겨집니다.
2. **조용한 실패(Silent Failure) 안티패턴**:
   - `_layout: two-col`처럼 오타가 발생하면, 파서는 아무런 경고나 에러 없이 조용히 `default` 1단 레이아웃으로 폴백합니다.
3. **디버깅 피로도 급증**:
   - 사용자는 화면에도 오타가 안 보이고, 터미널에도 아무런 경고가 없는데, 의도했던 2단 분할이나 배경색이 먹히지 않아 엉뚱한 곳(CSS, 브라우저 캐시 등)을 헤매며 긴 시간을 소모합니다.

---

### 1.2 대표적인 오류 케이스 분석

| 분류 | 작성자 오타 / 오류 입력 예시 | 현재 동작 (AS-IS) | 문제점 |
| :--- | :--- | :--- | :--- |
| **구분자 문법 오류** | `<!-- _layout: two-cols --->` | 뒤의 `-`가 값으로 흡수되어 `two-cols -`로 파싱 | 조용히 `default`로 폴백, 주석은 삭제됨 |
| **구분자 미종료** | `<!-- _layout: two-cols -- >`<br>`<!-- _color: red` | 다음 `-->`를 만날 때까지 본문 전체를 주석으로 흡수 | **슬라이드 본문 여러 줄이 통째로 증발(Swallowing)** |
| **레이아웃 키워드 오타** | `<!-- _layout: two-col -->`<br>`<!-- _layout: twocol -->` | 미등록 레이아웃으로 판정 | 조용히 `default` 1단으로 폴백 (2단 분할 깨짐) |
| **특수 마커 오타** | `<!-- splti -->`<br>`<!-- puase -->` | 일반 무의미 주석으로 간주되어 마크다운에서 제거 | 분할 또는 단계적 빌드가 조용히 무시됨 |
| **지시어 속성명 오타** | `<!-- _bakground: red -->`<br>`<!-- _pagiante: true -->` | 알려지지 않은 지시어 키로 간주되어 무시 | 스타일이 적용되지 않음 |

---

### 1.3 개선 목표

1. **Fail-Fast & 진단 피드백 (Actionable Diagnostics)**:
   - 오타나 문법 오류 발생 시 조용히 넘어가지 않고 터미널과 개발 서버에 명확한 Warning/Error 피드백 제공.
2. **유사도 기반 스마트 키워드 추천 ("Did you mean?")**:
   - Levenshtein 거리 알고리즘을 활용하여 `two-col` ➔ `Did you mean "two-cols"?`, `splti` ➔ `Did you mean "<!-- split -->"?` 추천 제공.
3. **비정상 주석 구분자 차단**:
   - `--->`, `<!-- ... -- >` 등 비표준 주석 닫기 구문을 정밀 검출하여 본문 증발(Comment Swallowing) 방지.

---

## 2. 상세 기술 설계 (Technical Specifications)

### 2.1 진단 도메인 모델 ([`internal/model/diagnostic.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/diagnostic.go) 신규 설계)

파싱 과정에서 수집된 진단 정보를 구조화된 모델로 관리합니다:

```go
package model

type DiagnosticSeverity string

const (
	SeverityError   DiagnosticSeverity = "error"   // 빌드 중단이 필요한 치명적 오류 (YAML Frontmatter 등)
	SeverityWarning DiagnosticSeverity = "warning" // 폴백되지만 의도와 다를 수 있는 오타/문법 경고 (DSL 오타 등)
	SeverityInfo    DiagnosticSeverity = "info"    // 권장사항 안내
)

type Diagnostic struct {
	Severity   DiagnosticSeverity `json:"severity"`
	SlideIndex int                // 1-based 슬라이드 번호
	Line       int                // 슬라이드 내 상대 줄 번호 (또는 문서 전체 줄 번호)
	Rule       string             // 규칙 코드 (예: "directive-unknown-layout", "malformed-comment")
	MessageKey string             // i18n 메시지 키
	MessageArgs []any             // i18n 포맷 인자
	Suggestion string             // 교정 추천 제안 (예: "two-cols")
	RawSnippet string             // 문제가 된 원문 조각 (예: "<!-- _layout: two-col -->")
}
```

---

### 2.2 주석 토크나이저 & 구분자 정밀 검증 ([`internal/parser/directive.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive.go))

정규식의 non-greedy 매칭에 의존하지 않고, 변형된 닫는 태그(`--->`, `-- >`)를 명시적으로 탐지합니다.

```go
// 1. 변형된 닫는 주석 감지 정규식
var malformedClosingRegex = regexp.MustCompile(`<!--(.*?)(--+>|--\s+>)`)

// 2. 주석 끝의 불필요한 대시(-) 정리 및 경고 수집
func cleanCommentDelimiter(rawComment string) (cleaned string, hasDelimiterWarning bool) {
	trimmed := strings.TrimSpace(rawComment)
	if strings.HasSuffix(trimmed, "-") {
		return strings.TrimRight(trimmed, "- \t"), true
	}
	return trimmed, false
}
```

---

### 2.3 레이아웃 및 지시어 유사도 알고리즘 (Levenshtein Distance)

사용자의 오타와 가장 가까운 표준 키워드를 $O(M \times N)$으로 계산하여 추천합니다:

```go
// FindBestMatch finds the closest valid keyword within maxDistance.
func FindBestMatch(input string, validCandidates []string, maxDistance int) string {
	bestCandidate := ""
	minDist := maxDistance + 1

	for _, candidate := range validCandidates {
		dist := LevenshteinDistance(strings.ToLower(input), strings.ToLower(candidate))
		if dist < minDist {
			minDist = dist
			bestCandidate = candidate
		}
	}

	if minDist <= maxDistance {
		return bestCandidate
	}
	return ""
}
```

#### 진단 예시:
* `_layout: two-col` 입력 시 ➔ `FindBestMatch("two-col", ["default", "cover", "section", "two-cols", "lead", "blank"], 2)` ➔ `"two-cols"` 도출!
* `<!-- splti -->` 입력 시 ➔ `FindBestMatch("splti", ["split", "pause", "note"], 2)` ➔ `"split"` 도출!

---

### 2.4 i18n 진단 메시지 카탈로그 ([`internal/i18n/locales/`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/))

#### 1) `locales/ko.json`
```json
{
  "diagnostic.warning.unknown_layout": "슬라이드 %d: 알 수 없는 레이아웃 \"%s\"입니다. 기본값(\"default\")으로 대체되었습니다.",
  "diagnostic.suggestion": "추천: \"%s\"을(를) 의도하셨나요?",
  "diagnostic.warning.malformed_delimiter": "슬라이드 %d: 비정상적인 주석 닫기 구분자(\"%s\")가 감지되었습니다. 표준 \"-->\"를 사용해 주세요.",
  "diagnostic.warning.unknown_directive": "슬라이드 %d: 알 수 없는 지시어 키 \"%s\"입니다. 무시되었습니다.",
  "diagnostic.warning.unknown_marker": "슬라이드 %d: 정의되지 않은 마커 \"<!-- %s -->\"입니다. 일반 주석으로 취급되었습니다."
}
```

#### 2) `locales/en.json`
```json
{
  "diagnostic.warning.unknown_layout": "Slide %d: Unknown layout \"%s\". Falling back to \"default\".",
  "diagnostic.suggestion": "Did you mean \"%s\"?",
  "diagnostic.warning.malformed_delimiter": "Slide %d: Malformed comment closing delimiter \"%s\" detected. Please use standard \"-->\".",
  "diagnostic.warning.unknown_directive": "Slide %d: Unknown directive key \"%s\". Ignored.",
  "diagnostic.warning.unknown_marker": "Slide %d: Undefined marker \"<!-- %s -->\". Treated as raw comment."
}
```

---

### 2.5 개발 서버(`goslide serve`) 및 CLI 진단 출력 연동

마크다운 빌드 및 HMR 감지 시 `deck.Diagnostics`를 검사하여 터미널에 시각적으로 눈에 띄는 Warning 박스를 출력합니다:

```text
[goslide] ⚠️  Slide 2 Warning: Unknown layout "two-col"
            💡 Did you mean "two-cols"?
            Fallback to "default" layout applied.
```

브라우저에서도 필요 시 경고 토스트(Warning Banner)를 띄워 작성자가 에디터에서 즉시 수정할 수 있도록 유도합니다.

---

## 3. 단계별 구현 계획 (Implementation Steps)

| 단계 | 작업 내용 | 타겟 소스 파일 |
| :---: | :--- | :--- |
| **Step 1** | **진단 도메인 모델 및 Levenshtein 유사도 유틸리티 구현**<br>- `model.Diagnostic` 구조체 정의<br>- 문자열 유사도 알고리즘 및 단위 테스트 | [`internal/model/diagnostic.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/diagnostic.go)<br>[`internal/parser/levenshtein.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/levenshtein.go)<br>[`internal/parser/levenshtein_test.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/levenshtein_test.go) |
| **Step 2** | **i18n 진단 메시지 카탈로그 등록**<br>- `ko.json` 및 `en.json`에 `diagnostic.*` 메시지 키 추가 | [`internal/i18n/locales/ko.json`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/ko.json)<br>[`internal/i18n/locales/en.json`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/en.json) |
| **Step 3** | **주석 구분자 유효성 검증기 구현**<br>- `--->`, 공백 오타, 미종료 주석 감지<br>- 슬라이드 본문 증발 방지 및 경고 수집 단위 테스트 | [`internal/parser/directive.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive.go)<br>[`internal/parser/directive_test.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive_test.go) |
| **Step 4** | **지시어 키/값 및 특수 마커 유효성 검사기 연동**<br>- 레이아웃 오타(`two-col`), 마커 오타(`splti`), 속성 오타 감지 및 "Did you mean?" 제안 바인딩 | [`internal/parser/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/layout.go)<br>[`internal/parser/directive.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive.go) |
| **Step 5** | **개발 서버 및 CLI 진단 출력 연동**<br>- `goslide serve` 파일 변경 시 진단 Warning 터미널 출력<br>- `goslide build` 시 진단 요약 리포트 출력 | [`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go)<br>[`cmd/goslide/`](file:///home/yundream/myjob/cloit/Goslide/cmd/goslide/) |
| **Step 6** | **통합 검증 및 회귀 테스트**<br>- 오타 시나리오(two-col, splti, --->)별 단위 테스트 및 E2E 진단 검증 | 전체 테스트 스위트 |

---

## 4. 수용 기준 및 검증 계획 (Acceptance Criteria & Verification)

### 4.1 수용 기준 (Acceptance Criteria)

- [ ] `<!-- _layout: two-col -->` 입력 시 `Warning: Unknown layout "two-col"`과 함께 `Did you mean "two-cols"?` 추천 메시지가 터미널에 출력되어야 함.
- [ ] `<!-- splti -->` 입력 시 알려지지 않은 마커 경고와 함께 `Did you mean "<!-- split -->"?` 추천이 출력되어야 함.
- [ ] `<!-- _layout: two-cols --->`처럼 비정상적인 주석 닫기 구분자(`--->`) 입력 시 구문 경고가 출력되고, 값이 오염(`two-cols -`)되지 않고 정상적으로 정제되어야 함.
- [ ] 진단 메시지가 `internal/i18n`을 통해 한국어/영어 로케일로 정확히 번역되어야 함.
- [ ] 모든 진단 기능 추가 후 기존의 정상적인 마크다운 슬라이드 렌더링에 일체 사이드이펙트나 성능 저하가 없어야 함.
