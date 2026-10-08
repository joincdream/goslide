# [GOS-29] v1.0.0 릴리즈 카테고리별 테마 패키지 확충 및 샘플 템플릿 제작 계획서

> **티켓 번호**: [GOS-29](https://joincdream.atlassian.net/browse/GOS-29)  
> **마일스톤**: v1.0.0 정식 릴리즈 준비 (Official Themes & Starter Templates)  
> **마감일**: 2026-11-30  
> **상태**: 완료 (Done)  
> **담당 패키지**: [`internal/theme/`](../internal/theme/), [`cmd/goslide/`](../cmd/goslide/)  
> **연관 문서**:  
> - [`docs/planning/competitive-analysis-and-strategy.md`](../docs/planning/competitive-analysis-and-strategy.md)  
> - [`docs/okf/decisions-simplification.md`](../docs/okf/decisions-simplification.md) *(ADR-001: 단순성 및 단일 바이너리 원칙)*  
> **연관 소스 파일**:  
> - [`internal/theme/theme.go`](../internal/theme/theme.go)  
> - [`internal/theme/assets/css/`](../internal/theme/assets/css/)  
> - [`cmd/goslide/init.go`](../cmd/goslide/init.go)  

---

## 1. 개요 및 배경 (Context & Problem Statement)

Goslide의 기술적 완성도(CGO 0% 단일 정적 바이너리, < 50ms 초고속 빌드, 16:9 무마진 PDF 및 PPTX 메모 완벽 보존, 발표자 런타임 툴킷)는 이미 글로벌 최상위 수준(98점)에 도달했습니다.

그러나 현재 내장된 테마는 `default`, `clean`, `dark` 3종으로, 기능 검증용으로는 우수하나 사용자가 실제 비즈니스 보고, 학술 강의, 테크 컨퍼런스 등 목적에 맞추어 즉시 실무에 투입하기에는 시각적 개성과 다양성이 부족합니다.

본 태스크는 v1.0.0 정식 릴리즈 시 사용자가 별도의 CSS 코딩 없이 마크다운 Frontmatter(`theme: corporate`) 또는 CLI 옵션(`--theme`)만으로 완성도 높은 프레젠테이션을 제작할 수 있도록 **3대 핵심 카테고리 공식 테마 패키지**를 확충하고 스타터 템플릿을 연동합니다.

---

## 2. 카테고리별 테마 디자인 명세 (Theme Specifications)

단일 바이너리 철학을 지키기 위해 외부 웹 폰트 무단 호출을 지양하고, 시스템 고품질 폰트 스택 및 내장 CSS 변수를 기반으로 제작합니다.

### 2.1 Corporate (비즈니스 / 엔터프라이즈 보고)
- **대상**: 임원 보고, 주간 업무 보고, 고객사 제안서, IR 피칭
- **배경색 / 메인 톤**: 신뢰감을 주는 딥 네이비 (`#0f172a` / 서브: `#1e293b`), 화이트 캔버스
- **타이포그래피**: 산세리프(Inter, Pretendard, Apple SD 산돌고딕, 맑은 고딕), 명확한 헤딩 위계
- **스타일 특징**:
  - 카드 박스 및 리드 문단 테두리 악센트
  - 정갈한 헤더 트래커 라인 및 슬라이드 번호(Pagination) 배치
  - 표(Table) 헤더 딥블루 음영 및 짝수행 얼룩무늬(Zebra striping)

### 2.2 Academic (학술 / 교육 / 온라인 강좌)
- **대상**: 대학 강의, 연구 세미나, 기술 교육 영상(YouTuber), 핸드북
- **배경색 / 메인 톤**: 눈의 피로를 줄이는 따뜻한 페이퍼/아이보리 (`#fbfbfa` / 다크 텍스트: `#292524`)
- **타이포그래피**: 가독성 높은 모던 세리프/산세리프 조화
- **스타일 특징**:
  - KaTeX 수식 블록 중앙 정렬 및 부드러운 하이라이트 박스
  - 블록 인용(Blockquote) 좌측 두꺼운 세리프 바 및 이탤릭 강조
  - 리스트 항목 간 여백 최적화

### 2.3 Cyber-Dark (테크 컨퍼런스 / 개발자 밋업)
- **대상**: 개발자 밋업, 기술 아키텍처 발표, 시스템 인프라 설명회
- **배경색 / 메인 톤**: 짙은 미드나이트 블랙 (`#090d16` / `#0b0f19`)
- **타이포그래피**: 코딩 친화적 모노스페이스 헤딩 + 시안(`cyan-400`), 일렉트릭 퍼플(`purple-400`) 네온 악센트
- **스타일 특징**:
  - Chroma 코드 블록의 `{1,3-5}` 라인별 포커스 강조 효과 극대화
  - 네온 글로우(Glow) 효과의 구분선 및 태그 뱃지
  - 다크 테마 기반 Mermaid 다이어그램 렌더링 최적화

---

## 3. 단계별 구현 계획 (Implementation Steps)

| 단계 | 작업 내용 | 타겟 소스 파일 |
| :---: | :--- | :--- |
| **Step 1** | **신규 3종 테마 CSS 제작**<br>- `corporate.css`: 비즈니스 딥네이비 테마<br>- `academic.css`: 교육/강의용 아이보리 페이퍼 테마<br>- `cyber-dark.css`: 테크 컨퍼런스용 네온 다크 테마 | [`internal/theme/assets/css/corporate.css`](../internal/theme/assets/css/corporate.css)<br>[`internal/theme/assets/css/academic.css`](../internal/theme/assets/css/academic.css)<br>[`internal/theme/assets/css/cyber-dark.css`](../internal/theme/assets/css/cyber-dark.css) |
| **Step 2** | **내장 테마 레지스트리 등록 및 유효성 검증**<br>- `theme.go`의 지원 테마 목록(`defaultThemes`)에 신규 3종 추가<br>- `GetThemeCSS()` 매핑 및 단위 테스트 갱신 | [`internal/theme/theme.go`](../internal/theme/theme.go)<br>[`internal/theme/theme_test.go`](../internal/theme/theme_test.go) |
| **Step 3** | **외부 커스텀용 독립 CSS 번들 배포**<br>- 사용자가 다운로드하여 수정할 수 있도록 루트 `themes/` 디렉토리에 CSS 원본 배치 | `themes/*.css` |
| **Step 4** | **`goslide init` 및 예제 슬라이드 템플릿 연동**<br>- `goslide init talk.md --theme corporate` 등 옵션 지원<br>- `examples/` 디렉토리에 테마별 데모 마크다운 슬라이드 제작 | [`cmd/goslide/init.go`](../cmd/goslide/init.go)<br>`examples/themes/*.md` |
| **Step 5** | **포맷별 무결성 검증**<br>- HTML, 16:9 무마진 PDF, 1080p PPTX 익스포트 시 서식/레이아웃 일치성 검증 | 전체 테스트 스위트 |

---

## 4. 수용 기준 (Acceptance Criteria)

- [x] **테마 등록 및 유효성**: Frontmatter `theme: corporate`, `theme: academic`, `theme: cyber-dark` 선언 시 올바른 CSS가 로드되어야 함.
- [x] **CLI 플래그 지원**: `goslide build talk.md -t academic` 및 `goslide serve talk.md -t corporate` 정상 동작.
- [x] **`goslide init` 연동**: `goslide init sample.md --theme cyber-dark` 실행 시 해당 테마가 지정된 스타터 마크다운 파일이 생성되어야 함.
- [x] **멀티 포맷 출력 무결성**: 신규 테마 3종에 대해 `-f all` 실행 시 HTML, PDF, PPTX 모두 폰트/레이아웃 깨짐 없이 렌더링되어야 함.
- [x] **바이너리 무의존성 유지**: 외부 CDN 폰트 로드 실패 시에도 시스템 로컬 폰트로 안전하게 폴백되어 오프라인 환경에서 100% 정상 작동해야 함.
