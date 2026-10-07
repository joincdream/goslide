# [GOS-5] M1-4: 테마 시스템 및 내장 정적 에셋 파이프라인 계획서

> **티켓 번호**: [GOS-5](https://joincdream.atlassian.net/browse/GOS-5)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-4  
> **마감일**: 2026-10-16  
> **완료일**: 2026-10-03  
> **상태**: 완료 (Done)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [core-architecture.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/core-architecture.md), [decisions-simplification.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 슬라이드의 시각적 완성도를 결정짓는 **내장 테마 시스템(Theme System)**과 Go `embed.FS` 기반 **정적 에셋 파이프라인**을 구축하는 것입니다.

임시 프리뷰 스타일을 넘어, 단일 실행 파일 내에 전문 프레젠테이션 테마 3종(`default`, `clean`, `dark`)과 공통 프레임 규격(`base.css`)을 내장하여 **외부 네트워크나 별도 CSS 파일 없이도 완벽한 슬라이드 디자인을 즉시 제공**합니다.

### 1.1 슬라이드 작성자(End-User)에게 제공하는 가치
1. **설정 없는 프로덕션 품질 디자인 (Zero Setup Aesthetic)**:
   - CSS를 한 줄도 모르는 개발자라도 마크다운만 작성하면 Pretendard 기반의 세련된 타이포그래피, 황금비 여백, 16:9 와이드 규격의 전문 발표 자료를 얻을 수 있습니다.
2. **상황별 테마 3종 기본 제공**:
   - `default`: 비즈니스 및 기술 발표를 위한 표준 모던 라이트 테마
   - `clean`: 군더더기 없는 미니멀리스트 화이트 테마 (Inter 폰트 스택)
   - `dark`: 어두운 발표장과 개발자 세미나에 최적화된 고대비 다크 테마 (One Dark 팔레트)
3. **완벽한 한글/CJK 타이포그래피**:
   - 한글 단어가 어색하게 줄바꿈되지 않도록 `word-break: keep-all;` 및 Pretendard 한글 폰트 스택을 기본 적용합니다.
4. **자유로운 스타일 확장성 (Custom CSS & Theme Override)**:
   - Frontmatter의 `custom_css` / `style` 지시어나 CLI의 `--theme-path` 플래그를 통해 기업 브랜드 컬러나 전용 폰트로 100% 오버라이드할 수 있습니다.

### 1.2 시스템 및 후속 단계(Renderer/Exporter)에게 제공하는 가치
1. **단일 바이너리(Single Binary) 자립성**:
   - 테마 CSS가 바이너리에 완전 내장(`embed.FS`)되어, 오프라인 환경이나 폐쇄망에서도 스타일이 깨지지 않는 완전한 독립형 HTML을 생성합니다.
2. **브라우저 네이티브 캐스케이딩 (ADR-001 준수)**:
   - 백엔드 Go 메모리에서 복잡한 6단계 스타일 우선순위 연산기를 돌리지 않고, 웹 표준 브라우저 CSS 엔진에 캐스케이딩을 위임하여 복잡도를 낮추고 렌더링 성능을 극대화합니다.

---

## 2. 결과물의 구체적 모습 (Concrete Manifestation)

사용자는 마크다운 Frontmatter에 `theme`를 선언하거나, CLI 플래그로 즉시 테마를 전환하여 브라우저에서 디자인 완성도를 검토할 수 있습니다.

### 2.1 CLI 명령을 통한 테마 전환 확인
```bash
# 1. 기본 테마로 빌드 (Default)
./bin/goslide build presentation.md -o talk-default.html

# 2. 미니멀 클린 테마로 빌드 (Clean)
./bin/goslide build presentation.md --theme=clean -o talk-clean.html

# 3. 고대비 다크 테마로 빌드 (Dark)
./bin/goslide build presentation.md --theme=dark -o talk-dark.html

# 4. 외부 기업 커스텀 CSS 적용
./bin/goslide build presentation.md --theme-path=./corporate.css -o talk-corp.html
```

### 2.2 테마별 디자인 스펙 비교

| 테마명 | 배경색 | 텍스트 기본색 | 헤딩 및 강조 컬러 | 주 사용 폰트 스택 | 분위기 및 권장 용도 |
| :--- | :---: | :---: | :---: | :--- | :--- |
| **`default`** | `#ffffff` | `#1e293b` | `#0f172a`, `#2563eb` | Pretendard, sans-serif | 기술 세미나, 사내 보고, 공식 발표 |
| **`clean`** | `#fafafa` | `#262626` | `#171717`, `#525252` | Inter, sans-serif | 미니멀 아키텍처 다이어그램, 간결한 요약 |
| **`dark`** | `#0f172a` | `#f1f5f9` | `#ffffff`, `#38bdf8` | Pretendard, Consolas | 개발자 컨퍼런스, 어두운 무대 발표 |

---

## 3. 핵심 설계 원칙 및 규칙 (Architectural Rules)

### 3.1 브라우저 네이티브 캐스케이딩 ([ADR-001](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md))
복잡한 Go단 캐스케이더 연산(`cascader.go`)을 배제하고, HTML `<head>` 내에 CSS 태그를 웹 표준 순서대로 배치하여 브라우저 엔진에 위임합니다:

```html
<!-- 1. 기본 뷰포트/16:9 규격 프레임워크 -->
<style id="goslide-base">
  /* base.css */
</style>

<!-- 2. 선택된 테마 스타일 -->
<style id="goslide-theme">
  /* default.css or clean.css or dark.css */
</style>

<!-- 3. CLI 외부 커스텀 CSS 파일 (--theme-path) -->
<style id="goslide-custom-file">
  /* custom file css */
</style>

<!-- 4. Frontmatter 인라인 커스텀 스타일 (custom_css / style) -->
<style id="goslide-inline-style">
  /* markdown inline custom css */
</style>
```
* **결과**: 뒤에 선언된 스타일이 자연스럽게 앞의 규칙을 덮어쓰므로(Cascading), Go 백엔드 코드 복잡도 0줄로 완벽한 오버라이드 보장.

### 3.2 단방향 의존성 및 패키지 격리
- `internal/theme` 패키지는 Go 표준 라이브러리(`embed`, `io/fs`, `os`) 외에 어떠한 외부 패키지나 상위 모듈도 참조하지 않는 최하위 리소스 계층으로 격리합니다.
- `embed.FS`는 OS 독립적인 슬래시(`/`) 가상 경로 체계를 사용하여 Windows, macOS, Linux에서 동일하게 컴파일 및 구동됩니다.

---

## 4. 세부 개발 태스크 (Work Breakdown Structure)

### 4.1 핵심 CSS 4종 작성 (`internal/theme/assets/css/`)
1. **`base.css`**:
   - 슬라이드 뷰포트 규격: 16:9 비율(기본 960x540 가상 박스, 반응형 스케일링 준비)
   - 타이포그래피 초기화 및 CJK `word-break: keep-all;`
   - 헤더, 푸터, 페이지 번호 인디케이터 절대/상대 위치 배치
   - 2단 그리드(`two-cols`) 및 표(`table`), 인라인 코드 기본 여백
2. **`default.css`**:
   - Pretendard 기반 기술 표준 라이트 팔레트
   - H1~H3 헤딩 계층, 인용구(Blockquote) 좌측 보더 강조
3. **`clean.css`**:
   - Inter 기반 초경량 여백 중심 디자인
   - 불필요한 장식선 제거 및 은은한 그레이 팔레트
4. **`dark.css`**:
   - 고대비 One Dark 팔레트 (배경 `#0f172a`, 텍스트 `#f8fafc`)
   - 야간 발표 가독성에 최적화된 링크/강조 컬러

### 4.2 내장 에셋 관리자 구현 (`internal/theme/embed.go` & `theme.go`)
- `//go:embed assets/css/*.css`를 통한 바이너리 내장
- `type ThemeManager struct`:
  - `NewThemeManager() *ThemeManager` 생성자
  - `GetThemeCSS(themeName string) (string, error)`: 테마명에 따른 CSS 반환 (미지원 테마는 `default`로 자동 폴백)
  - `ComposeFullCSS(themeName string, themePath string, inlineCSS string) (string, error)`: base + theme + external + inline CSS 단일 스트림 합성

### 4.3 CLI 테마 플래그 연동 (`cmd/goslide/build.go` & `root.go`)
- Cobra 글로벌/로컬 플래그 바인딩:
  - `--theme`: 테마명 지정 (`default`, `clean`, `dark` / 기본값: 마크다운 Frontmatter 설정값 또는 `default`)
  - `--theme-path`: 외부 커스텀 CSS 파일 경로 주입
- `goslide build` 실행 시 테마 에셋 관리자를 호출하여 실제 테마 CSS가 온전히 주입된 HTML 슬라이드 생성

### 4.4 단위 테스트 스위트 작성 (`internal/theme/theme_test.go`)
- 내장 테마 3종 정상 로딩 검증
- 존재하지 않는 테마 요청 시 `default` 폴백 검증
- 외부 커스텀 CSS 파일 로딩 및 합성 검증
- CSS 캐스케이딩 순서(Base -> Theme -> Custom -> Inline) 무결성 검증

---

## 5. 엔지니어링 가드레일 (Hard Constraints)

1. **CGO 배제 (100% Pure Go)**: Go 표준 `embed.FS`만 활용하여 단일 정적 바이너리 컴파일 보장.
2. **글로벌 상태 배제**: 전역 변수 대신 `NewThemeManager()` 인스턴스 생성자 패턴 준수.
3. **순환 참조 방지 계층 엄수**: `internal/theme`는 `internal/model`만 참조하거나 독립 유지.
4. **보안 및 예외 안전성**: 외부 CSS 로딩 시 파일 존재 여부 및 권한 오류에 대해 명시적 에러 반환.

---

## 6. 완료 기준 (Definition of Done)

- [x] `internal/theme/assets/css/` 디렉토리에 `base.css`, `default.css`, `clean.css`, `dark.css`가 작성됨.
- [x] `internal/theme/embed.go`에서 `embed.FS`로 CSS 에셋이 정상 번들링됨.
- [x] `internal/theme/theme.go`가 구현되어 테마별 CSS 로딩 및 폴백, 외부 CSS 합성이 정상 작동함.
- [x] `./bin/goslide build demo.md --theme=clean -o demo.html` 실행 시 실제 테마 CSS가 적용된 완성도 높은 슬라이드가 출력됨.
- [x] 테이블 기반 단위 테스트가 작성되고 `go test -v -race ./internal/theme/...`를 100% 통과함.
- [x] `make lint` 및 `make complexity` 기준을 무결하게 만족함.

---

## 7. 실행 절차 (Step-by-Step Execution Plan)

1. **CSS 에셋 작성**: `base.css`, `default.css`, `clean.css`, `dark.css` 구현
2. **에셋 내장 및 로더 구현**: `embed.go`, `theme.go` 작성
3. **단위 테스트 작성 및 검증**: `theme_test.go` 테이블 기반 테스트
4. **CLI 연동**: `cmd/goslide/build.go`에 `--theme`, `--theme-path` 플래그 및 테마 합성기 연결
5. **품질 검증**: `make lint`, `make complexity`, `make test-race`, `make build`
6. **실제 출력 확인**: 3종 테마 빌드 및 브라우저 시각 디자인 검토
