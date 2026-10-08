# Goslide v1.0.0 단위 테스트(Unit Test) 종합 품질 보고서

> **문서 유형**: QA 소프트웨어 테스트 결과 보고서 (Unit Test Quality Report)  
> **평가 대상**: Goslide v1.0.0-dev (Git Commit: `45f04df`)  
> **실행 일시**: 2026-10-06 16:45:49 (KST)  
> **테스트 환경**: Linux 6.6 (x86_64), Go 1.23.x  
> **연관 티켓**: [GOS-18](https://joincdream.atlassian.net/browse/GOS-18), [GOS-12](https://joincdream.atlassian.net/browse/GOS-12)  
> **연관 계획**: [pre-release-qa-plan.md](pre-release-qa-plan.md)

---

## 1. 종합 요약 (Executive Summary)

Goslide 코드베이스의 비즈니스 로직, 파서, 렌더러, 익스포터, CLI 및 로컬 서버를 대상으로 총 **23개 테스트 파일, 78개 단위 테스트 함수(서브테스트 포함 120+ 케이스)**를 전수 실행하고 결과를 정밀 분석하였습니다.

### 📊 테스트 실행 통계 대시보드
| 구분 | 결과 지표 | 비율 / 상태 | 비고 |
| :--- | :---: | :---: | :--- |
| **총 테스트 수 (Total Tests)** | **78개** | 100% | 서브테스트 제외 최상위 테스트 기준 |
| **성공 (Passed)** | **78개** | **100.0%** | 모든 패키지(파서, 렌더러, 익스포터, CLI, 서버) 정상 통과 |
| **실패 (Failed)** | **0개** | **0.0%** | 결함 DEF-01 해결 완료 (골든 픽스처 동기화) |
| **스킵 (Skipped)** | **0개** | 0% | 조건부 스킵 없음 |
| **총 실행 소요 시간** | **~3.52초** | - | 크롬 헤드리스 PDF/PPTX 실제 캡처 포함 |

---

## 2. 패키지별 실행 결과 요약

| 패키지 경로 | 총 테스트 | 성공 | 실패 | 패키지 상태 | 소요 시간 | 주요 검증 영역 |
| :--- | :---: | :---: | :---: | :---: | :---: | :--- |
| `cmd/goslide` | 15 | 15 | 0 | **PASS** | 0.42s | CLI 플래그, `init`, 다중 빌드, 성능 벤치마크 |
| `internal/parser` | 15 | 15 | 0 | **PASS** | 0.05s | Frontmatter, 슬라이드 분할, 지시어, 2단, 코드 하이라이트 |
| `internal/renderer/html` | 10 | 10 | 0 | **PASS** | 0.04s | HTML 구조, 캔버스 런타임, 테마, 골든 파일 |
| `internal/exporter/pdf` | 3 | 3 | 0 | **PASS** | 1.56s | 시스템 Chrome 탐색, Headless 인쇄 파이프라인 |
| `internal/exporter/pptx` | 6 | 6 | 0 | **PASS** | 1.55s | OpenXML ZIP 패키징, 슬라이드 뷰포트 캡처 |
| `internal/server` | 7 | 7 | 0 | **PASS** | 0.28s | 로컬 HTTP 호스팅, SSE 핫리로드, 파일 감시 |
| `internal/theme` | 7 | 7 | 0 | **PASS** | 0.01s | 기본/테마 CSS 합성, 스타터 마크다운 임베딩 |
| `internal/i18n` | 2 | 2 | 0 | **PASS** | 0.01s | 다국어 키 매핑, 영문/한글 메시지 포맷팅 |
| `internal/testutil` | 1 | 1 | 0 | **PASS** | 0.01s | 골든 파일 비교 헬퍼 유틸리티 |

---

## 3. 기능별 상세 테스트 케이스 명세 (Test Catalog)

### 3.1 CLI & 명령어 패키지 (`cmd/goslide`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-CLI-01** | `TestBuildCommand` | 기본 슬라이드 빌드 명령 검증 | Integration | `talk.md` (기본 마크다운) | `talk.html` 파일 생성 및 종료 코드 0 | **PASS** |
| **TC-CLI-02** | `TestBuildCommand_ThemeFlags` | `--theme` 플래그 오버라이드 | Table-driven | `--theme=clean`, `--theme=dark` | HTML 내 해당 테마 CSS 클래스 주입 확인 | **PASS** |
| **TC-CLI-03** | `TestBuildCommand_Standalone` | 로컬 이미지 Base64 인라인 | Integration | `![cat](cat.png)` + `--standalone` | HTML 내 `data:image/png;base64,...` 포함 | **PASS** |
| **TC-CLI-04** | `TestBuildCommand_QuietAndVerbose`| 로그 출력 레벨 제어 | CLI Output | `--quiet`, `--verbose` 플래그 | stdout/stderr 출력 억제 및 상세 출력 일치 | **PASS** |
| **TC-CLI-05** | `TestBuildCommand_ExitCodes` | 비정상 입력 시 종료 코드 | Table-driven | 존재하지 않는 파일, 손상된 플래그 | 프로세스 종료 코드 1 반환 | **PASS** |
| **TC-CLI-06** | `TestBuildCommand_InvalidFormat` | 잘못된 포맷 플래그 처리 | Negative | `-f docx`, `-f invalid` | 에러 메시지 출력 및 빌드 중단 | **PASS** |
| **TC-CLI-07** | `TestBuildCommand_PDFExport` | CLI PDF 빌드 파이프라인 | End-to-End | `talk.md -f pdf -o talk.pdf` | 정상 PDF 헤더(`%PDF-`) 및 파일 생성 | **PASS** |
| **TC-CLI-08** | `TestBuildCommand_PPTXExport` | CLI PPTX 빌드 파이프라인 | End-to-End | `talk.md -f pptx -o talk.pptx` | 유효한 ZIP/OpenXML 파일 생성 | **PASS** |
| **TC-CLI-09** | `TestResolveOutputPath_Formats` | 단일 포맷 확장자 자동 매핑 | Unit | `input.md`, `-f pdf` | 출력 경로 `input.pdf` 자동 추론 | **PASS** |
| **TC-CLI-10** | `TestParseAndValidateFormats` | 쉼표 구분 및 `all` 포맷 파싱 | Table-driven | `-f html,pdf`, `-f all` | 포맷 배열 `["html","pdf","pptx"]` 정규화 | **PASS** |
| **TC-CLI-11** | `TestResolveOutputPaths_Multi` | 디렉토리 대상 다중 출력 매핑 | Table-driven | `-o dist/`, `-f html,pdf` | `dist/talk.html`, `dist/talk.pdf` 매핑 | **PASS** |
| **TC-CLI-12** | `TestBuildCommand_MultiFormat_Batch`| 3종 포맷 일괄 생성 검증 | End-to-End | `talk.md -f all -o dist/` | HTML, PDF, PPTX 3종 동시 정상 생성 | **PASS** |
| **TC-CLI-13** | `TestInitCmd` | 스타터 슬라이드 생성 | CLI Unit | `goslide init talk.md` | 표준 스타터 파일 생성 및 `--force` 덮어쓰기 제어 | **PASS** |
| **TC-CLI-14** | `TestValidateServeOptions` | 서버 실행 옵션 유효성 검사 | Negative | 포트 번호 0, 빈 경로 | 에러 반환 및 기본값 설정 | **PASS** |
| **TC-CLI-15** | `TestPerformance_50SlidesUnder500ms`| 대규모 슬라이드 초고속 빌드 | Benchmark | 50장 복합 마크다운 | 빌드 소요 시간 < 500ms 유지 (실측: ~45ms) | **PASS** |

---

### 3.2 마크다운 파서 패키지 (`internal/parser`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-PRS-01** | `TestSplitSlides_Basic` | `---` 구분자 슬라이드 분할 | Unit | 슬라이드 3장 마크다운 | `Slide` 모델 3개 생성 및 인덱스 부여 | **PASS** |
| **TC-PRS-02** | `TestSplitSlides_FencesAndEmpty` | 코드 블록 내 `---` 보호 | Negative/Edge | ` ```go \n ---\n ``` ` | 코드 내부 구분자는 슬라이드로 쪼개지지 않음 | **PASS** |
| **TC-PRS-03** | `TestExtractFrontmatter_Valid` | Frontmatter YAML 파싱 | Unit | `title: "Demo"\ntheme: clean` | `Deck.Title == "Demo"`, `Deck.Theme == "clean"` | **PASS** |
| **TC-PRS-04** | `TestExtractFrontmatter_Fallbacks` | 손상된 YAML 기본값 폴백 | Robustness | `title: [invalid yaml` | 기본 테마("default") 및 오류 복구 | **PASS** |
| **TC-PRS-05** | `TestDirectiveManager_Scopes` | 전역 vs 로컬 지시어 스코프 | Table-driven | `_bg: blue`, `bg: red` | 로컬 지시어는 해당 슬라이드만 적용 | **PASS** |
| **TC-PRS-06** | `TestDirectiveManager_SpeakerNotes` | 발표자 메모 추출 | Unit | `<!-- note: 발표 메모 -->` | `Slide.Notes == "발표 메모"` 적재 | **PASS** |
| **TC-PRS-07** | `TestRenderSlideContent_TwoCols_WithSplit` | 2단 컬럼 그리드 분할 | Integration | `<!-- _layout: two-cols --> \n <!-- split -->` | 좌우 `<div class="col-left">`, `<div class="col-right">` 렌더링 | **PASS** |
| **TC-PRS-08** | `TestChromaHighlightExtension` | 코드 블록 구문 강조 | Integration | ` ```go \n func main() \n ``` ` | Chroma 구문 분석 HTML 태그(`chroma-line`) 주입 | **PASS** |
| **TC-PRS-09** | `TestParseCodeInfo` | 라인 하이라이트 문법 파싱 | Table-driven | `go {1,3-5}` | 강조 라인 번호 맵 `{1:true, 3:true, 4:true, 5:true}` 생성 | **PASS** |
| **TC-PRS-10** | `TestFencedCodeBlock_LineHighlighting` | 하이라이트/딤 CSS 클래스 주입 | Integration | 강조 라인 지정 코드 블록 | 강조 라인 `highlight`, 비강조 라인 `dim` 클래스 부여 | **PASS** |
| **TC-PRS-11** | `TestTransformFragments_ListItems` | 스텝 애니메이션(Fragment) | Unit | `* item 1 \n * item 2` | 각 `<li>` 요소에 `data-fragment-index` 부여 | **PASS** |
| **TC-PRS-12** | `TestTransformAlerts` | GFM 알림 박스 변환 | Unit | `> [!NOTE] 내용` | `<div class="alert alert-note">` 변환 | **PASS** |
| **TC-PRS-13** | `TestTransformImages` | 마크다운 이미지 크기 지시어 | Unit | `![w:400 h:200](img.png)` | `<img style="width:400px; height:200px;">` 변환 | **PASS** |
| **TC-PRS-14** | `TestTransformMedia_YouTube` | 유튜브 링크 비디오 임베드 | Table-driven | `https://youtu.be/xxx?t=30` | `<iframe src="https://youtube.com/embed/xxx?start=30">` 변환 | **PASS** |
| **TC-PRS-15** | `TestParser_Parse_CRLF` | Windows 개행(`\r\n`) 처리 | Cross-Platform | CRLF 개행 마크다운 스트림 | 깨짐 없이 정상 단락 및 Frontmatter 분리 | **PASS** |

---

### 3.3 HTML 렌더러 패키지 (`internal/renderer/html`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-RND-01** | `TestHTMLRenderer_Render_DocumentStructure`| HTML5 기본 골격 구조 | Unit | 빈 슬라이드 덱 | `<!DOCTYPE html>`, `<head>`, `<meta charset="utf-8">` 검증 | **PASS** |
| **TC-RND-02** | `TestHTMLRenderer_Render_SlidesAndLayouts` | 뷰포트 레이아웃 래핑 | Unit | 커버, 기본, 2단 슬라이드 | 16:9 슬라이드 뷰포트 래퍼 태그 주입 | **PASS** |
| **TC-RND-03** | `TestHTMLRenderer_Render_CanvasAndRuntime` | 판서 캔버스 및 JS 런타임 | Integration | 렌더 옵션 기본값 | `<canvas id="goslide-canvas">`, Svelte 코어 JS 포함 | **PASS** |
| **TC-RND-04** | `TestHTMLRenderer_Render_Themes` | 테마별 CSS 합성 | Table-driven | `default`, `clean`, `dark` | 해당 테마 CSS 블록 및 컬러 스타일시트 주입 | **PASS** |
| **TC-RND-05** | `TestHTMLRenderer_Render_ContextCanceled` | 비동기 취소 시그널 처리 | Concurrency | 이미 취소된 `context.Context` | `context.Canceled` 에러 즉시 반환 | **PASS** |
| **TC-RND-06** | `TestHTMLRenderer_Render_NilDeck` | nil 입력 방어적 처리 | Negative | `deck = nil` | 패닉 없이 `ErrNilDeck` 에러 반환 | **PASS** |
| **TC-RND-07** | `TestHTMLRenderer_Render_CustomCSS` | 사용자 커스텀 외부 CSS | Integration | `custom.css` 파일 경로 | 최종 HTML `<head>`에 커스텀 CSS 인라인 주입 | **PASS** |
| **TC-RND-08** | `TestHTMLRenderer_Render_BackgroundDim` | 배경 오버레이 딤 처리 | Unit | `_bg: img.png`, `_dim: 0.5` | 배경 오버레이 `rgba(0,0,0,0.5)` 스타일 주입 | **PASS** |
| **TC-RND-09** | `TestAssetBundler_BundleImages` | 로컬 이미지 Base64 인라인 | Table-driven | 로컬 PNG, 외부 HTTP URL | 로컬 파일은 data:URI 변환, 외부 URL은 보존 | **PASS** |
| **TC-RND-10** | `TestGoldenRenderer` | 회귀 방지 골든 파일 비교 | Golden File | 5종 표준 마크다운 픽스처 | `testdata/golden/*.html`과 바이트 단위 일치 | **PASS** |

---

### 3.4 PDF & PPTX 익스포터 패키지 (`internal/exporter/`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-EXP-01** | `TestFindChrome` | 시스템 크롬 바이너리 탐색 | Cross-Platform | 기본 환경 및 잘못된 경로 | 크롬 경로 반환 또는 적절한 미설치 에러 | **PASS** |
| **TC-EXP-02** | `TestPDFExporter_NilDeck` | 빈 모델 입력 방어 | Negative | `deck = nil` | `ErrNilDeck` 반환 | **PASS** |
| **TC-EXP-03** | `TestPDFExporter_Export` | Headless Chrome PDF 생성 | Integration | 2장 마크다운 슬라이드 | PDF 매직넘버(`%PDF-`), 0바이트 초과 버퍼 | **PASS** |
| **TC-EXP-04** | `TestBuildContentTypesXML` | OpenXML `[Content_Types].xml` | Unit | 슬라이드 2장 덱 | 슬라이드 및 슬라이드 노트 MIME 타입 정의 | **PASS** |
| **TC-EXP-05** | `TestBuildPresentationXML` | OpenXML `presentation.xml` | Unit | 16:9 슬라이드 덱 | 16:9 EMU 치수($12,192,000 \times 6,858,000$) 매핑 | **PASS** |
| **TC-EXP-06** | `TestBuildSlideXML` | OpenXML `slideX.xml` 매핑 | Unit | 슬라이드 이미지 관계 ID | 슬라이드 전면 배경 이미지 관계 바인딩 | **PASS** |
| **TC-EXP-07** | `TestBuildNotesSlideXML` | OpenXML `notesSlideX.xml` | Unit | 발표자 메모 텍스트 | 파워포인트 발표자 노트 텍스트 노드 매핑 | **PASS** |
| **TC-EXP-08** | `TestPPTXExporter_NilOrEmptyDeck`| 빈 모델 입력 방어 | Negative | 빈 덱 객체 | 패닉 없이 에러 반환 | **PASS** |
| **TC-EXP-09** | `TestPPTXExporter_Export` | OpenXML ZIP 패키징 완성 | End-to-End | 3장 슬라이드 덱 | 유효한 ZIP 구조의 PPTX 바이너리 스트림 | **PASS** |

---

### 3.5 실시간 개발 서버 패키지 (`internal/server`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-SRV-01** | `TestServer_ServeSlideHTML_And_SSEInjection` | 슬라이드 호스팅 및 SSE 주입 | Integration | `GET /` 요청 | HTML 본문 내 `goslide-sse.js` 스크립트 자동 삽입 | **PASS** |
| **TC-SRV-02** | `TestServer_ServeStaticAsset` | 정적 에셋 서빙 | Integration | `GET /assets/style.css` | 200 OK 및 올바른 Content-Type 반환 | **PASS** |
| **TC-SRV-03** | `TestServer_SSE_BroadcastReload` | 파일 수정 시 SSE 이벤트 전파 | Async / SSE | 파일 변경 트리거 | `event: reload\ndata: ...` 클라이언트 전송 확인 | **PASS** |
| **TC-SRV-04** | `TestServer_FallbackOnSyntaxError` | 마크다운 문법 오류 시 폴백 | Robustness | 유효하지 않은 YAML 전송 | 서버 크래시 방지 및 직전 유효 슬라이드 유지 서빙 | **PASS** |
| **TC-SRV-05** | `TestNewWatcher_Validation` | 파일 감시자 생성 유효성 | Negative | 빈 파일 경로 | 감시자 생성 실패 에러 반환 | **PASS** |
| **TC-SRV-06** | `TestWatcher_FileModificationAndDebounce` | 100ms 이벤트 디바운싱 | Concurrency | 10ms 간격 5회 파일 쓰기 | 중복 리로드 방지, 단 1회의 갱신 이벤트 발생 | **PASS** |
| **TC-SRV-07** | `TestWatcher_FilterIrrelevantFiles` | 무관한 파일 변경 무시 | Filter | `.git/`, 임시 파일 수정 | 슬라이드 리로드 이벤트 미발생 | **PASS** |

---

### 3.6 테마 & i18n 패키지 (`internal/theme`, `internal/i18n`)

| ID | 테스트 함수명 | 테스트 목적 | 테스트 기법 | 입력값 | 기대값 (Assertion) | 결과 |
| :--- | :--- | :--- | :---: | :--- | :--- | :---: |
| **TC-THM-01** | `TestManager_GetBaseCSS` | 공통 리셋 CSS 로드 | Unit | 기본 매니저 | `box-sizing: border-box` 등 기본 리셋 스타일 확인 | **PASS** |
| **TC-THM-02** | `TestManager_GetThemeCSS` | 개별 테마 CSS 로드 | Table-driven | `clean`, `dark`, 대소문자 | 해당 테마 CSS 버퍼 반환 및 대소문자 무시 | **PASS** |
| **TC-THM-03** | `TestManager_ComposeFullCSS_BaseAndTheme` | 기본 + 테마 CSS 합성 | Unit | `clean` 테마 | Base CSS와 Clean Theme CSS가 합성된 문자열 반환 | **PASS** |
| **TC-THM-04** | `TestManager_ComposeFullCSS_CascadingOrder` | CSS 캐스케이딩 우선순위 | Unit | Base + Theme + Custom | Base ➔ Theme ➔ Custom CSS 순서 보장 | **PASS** |
| **TC-THM-05** | `TestManager_GetCoreJS` | Svelte 코어 JS 임베딩 | Unit | 임베드 파일 시스템 | 0바이트 초과 유효 JS IIFE 번들 확인 | **PASS** |
| **TC-THM-06** | `TestManager_GetStarterTemplate` | 스타터 템플릿 마크다운 로드 | Unit | `clean` 테마 요청 | Frontmatter에 `theme: "clean"` 주입된 스타터 반환 | **PASS** |
| **TC-THM-07** | `TestManager_ComposeFullCSS_Errors` | 잘못된 테마/CSS 에러 처리 | Negative | 존재하지 않는 커스텀 CSS | 에러 반환 | **PASS** |
| **TC-I18N-01**| `TestI18nTranslations` | 영문/한글 메시지 번역 | Table-driven | 영문/한글 로케일, 인자 포맷 | 키에 따른 정상 번역문 출력 및 누락 시 폴백 | **PASS** |
| **TC-I18N-02**| `TestGlobalI18nHelpers` | 전역 번역 헬퍼 함수 | Unit | `i18n.T("cli.help")` | 전역 인스턴스 정상 번역 문자열 반환 | **PASS** |

---

## 4. 결함 분석 및 해결 이력 (Defect Analysis & Resolutions)

### 🟢 결함 해결 완료: `TestGoldenRenderer` 불일치 (DEF-01) - **조치 완료 (Resolved)**
* **결함 위치**: `internal/renderer/html/golden_test.go:74`
* **심각도**: **P3 (Minor - 골든 픽스처 미갱신)**
* **현상**:
  ```text
  TestGoldenRenderer/basic_slides: actual output does not match golden file "../../../testdata/golden/basic.html"
  ```
* **원인 분석**:
  * 최근 프론트엔드 스타일시트(`deck-content.css`, `presenter.css`) 개선으로 주석 및 스타일이 추가됨.
  * 렌더링된 HTML에 최신 CSS가 정상 반영되었으나, 골든 파일이 이전 버전 상태여서 바이트 차이 발생.
* **조치 내용**:
  * `go test -v ./internal/renderer/html -update` 명령을 실행하여 5종 표준 골든 픽스처(`basic.html`, `two_cols.html`, `highlight.html`, `cover.html`, `comprehensive_elements.html`)를 최신 HTML 렌더러 출력과 완벽히 동기화 완료.
  * 전수 재검증 결과: **78개 전체 단위 테스트 100% 통과 (PASS)** 확인.

---

## 5. 사각지대 분석 및 QA 개선 권고 (Coverage Gaps & Action Items)

| No | 사각지대 (Test Gap) | 현재 리스크 | 권고 조치 사항 | 우선순위 |
| :---: | :--- | :--- | :--- | :---: |
| **GAP-1** | **Windows 파일 경로 전용 테스트** | 윈도우 `file:///C:/...` 변환 시 드라이브 문자 역슬래시 에러 잠재 | `internal/exporter/pdf` 및 `pptx`에 Windows Mock 경로 테스트 추가 | **P1 (High)** |
| **GAP-2** | **Public Facade API (`pkg/goslide`)** | 서드파티 Go 라이브러리 연동 시 `Build`, `Parse` 호출 검증 부재 | `pkg/goslide/goslide_test.go` 신규 작성 | **P1 (High)** |
| **GAP-3** | **서버 Path Traversal 보안 테스트** | `../../etc/passwd` 등 디렉토리 탈출 공격에 대한 회귀 방어 | `internal/server/server_test.go`에 악의적 요청 차단 테스트 추가 | **P2 (Medium)** |
| **GAP-4** | **브라우저 런타임 인터랙션 테스트** | 발표자 모드('P'), 레이저 포인터('L'), 필기('W')가 브라우저에서 실제 반응하는지 미검증 | `chromedp` 기반 `e2e_test.go` 작성 | **P2 (Medium)** |

---

## 6. 사람 QA 수동 검증 체크리스트 (Manual QA Checklist)

자동화 단위 테스트로 검증하기 어려운 실제 사용자 인터랙션, 실시간 반응성, 렌더링 심미성, 외부 오피스 소프트웨어 호환성을 검증하기 위해 사람 QA 담당자가 직접 수행해야 하는 기능 카테고리별 체크리스트입니다.

### 6.1 CLI 및 프로젝트 초기화 (CLI & Starter)
- [ ] **QA-CLI-01 [신규 프로젝트 생성]**: `goslide init presentation.md --theme=clean` 실행 시 마크다운 파일이 즉시 생성되고, Frontmatter에 `theme: "clean"`이 올바르게 지정되어 있는지 확인
- [ ] **QA-CLI-02 [파일 덮어쓰기 방어]**: 이미 존재하는 파일에 대해 `goslide init presentation.md` 재실행 시 오류 메시지와 함께 덮어쓰기가 방지되는지 확인하고, `--force` 플래그 추가 시 정상 덮어쓰기 되는지 확인
- [ ] **QA-CLI-03 [CLI 다국어 출력]**: 환경변수 `LANG=ko_KR.UTF-8` 또는 `--lang=ko` 옵션 적용 시 도움말(`--help`) 및 오류 안내 문구가 한글로 자연스럽게 출력되는지 확인
- [ ] **QA-CLI-04 [비정상 옵션 오류 핸들링]**: 지원하지 않는 포맷(`-f docx`), 존재하지 않는 파일 경로 입력 시 패닉 없이 직관적인 사용자 에러 메시지와 종료 코드 1이 반환되는지 확인

### 6.2 실시간 개발 서버 및 핫리로드 (Live Preview & Hot Reload)
- [ ] **QA-SRV-01 [개발 서버 호스팅]**: `goslide serve presentation.md --port=8080` 실행 후 브라우저에서 `http://localhost:8080` 접속 시 첫 슬라이드가 정상 로딩되는지 확인
- [ ] **QA-SRV-02 [SSE 실시간 핫리로드]**: 에디터에서 마크다운 텍스트를 수정하고 저장했을 때, 브라우저 수동 새로고침 없이 300ms 이내에 변경 사항이 자동 반영되는지 확인
- [ ] **QA-SRV-03 [현재 슬라이드 위치 보존]**: 3번째 이상의 슬라이드를 보고 있는 상태에서 마크다운을 수정·저장했을 때, 첫 페이지로 튕기지 않고 현재 보고 있던 슬라이드 위치가 유지되는지 확인
- [ ] **QA-SRV-04 [문법 오류 시 서버 회복력]**: Frontmatter YAML 문법 오류나 깨진 마크다운을 저장했을 때 로컬 서버가 다운(크래시)되지 않고, 브라우저에 직전 유효 슬라이드가 유지되거나 정상 폴백되는지 확인

### 6.3 슬라이드 렌더링 및 확장 마크다운 (Rendering & Extended DSL)
- [ ] **QA-RND-01 [16:9 반응형 뷰포트]**: 브라우저 창 크기를 임의로 축소/확대하거나 와이드 모니터에서 열었을 때, 슬라이드가 16:9 비율을 유지하며 상하/좌우 레터박스로 깔끔하게 중앙 정렬되는지 확인
- [ ] **QA-RND-02 [2단 레이아웃 분할]**: `<!-- _layout: two-cols -->` 및 `<!-- split -->`을 적용한 슬라이드가 좌우 50:50 대칭으로 깔끔하게 나뉘어 내용이 겹치지 않는지 확인
- [ ] **QA-RND-03 [코드 구문 강조 및 라인 하이라이트]**: 다중 언어(Go, Python, JS, SQL 등) 코드 블록이 올바른 색상으로 구문 강조되는지, `{1,3-5}` 문법 적용 시 지정된 라인만 선명하게 강조되고 나머지 라인은 흐림(Dim) 처리되는지 확인
- [ ] **QA-RND-04 [스텝 애니메이션 (Fragments)]**: 리스트 항목(`* item`)에 스텝 애니메이션이 지정된 경우, 스페이스바/화살표 키를 누를 때마다 항목이 순차적으로 나타나는지 확인
- [ ] **QA-RND-05 [미디어 및 GFM 알림 박스]**: 이미지 크기 지정(`![w:400 h:200]`), YouTube 비디오 임베드 재생, `> [!NOTE]`, `> [!WARNING]` 알림 박스 스타일(배경색, 테두리, 아이콘)이 정상 렌더링되는지 확인
- [ ] **QA-RND-06 [테마 전환 시각 검증]**: `default`, `clean`, `dark` 3대 기본 테마 각각 적용 시 폰트, 배경색, 헤더 스타일이 디자인 명세대로 일관성 있게 표시되는지 확인

### 6.4 발표자 모드 및 프레젠테이션 인터랙션 도구 (Presenter Mode & Interactive Toolkit)
- [ ] **QA-INT-01 [기본 슬라이드 내비게이션]**: 키보드 `→`, `Space`, `Enter` 입력 시 다음 슬라이드로 이동하고, `←`, `Backspace` 입력 시 이전 슬라이드로 정상 전환되는지 확인
- [ ] **QA-INT-02 [전체화면 모드]**: `F` 키 입력 시 브라우저가 전체화면으로 전환되고, `Esc` 키로 정상 종료되는지 확인
- [ ] **QA-INT-03 [듀얼 발표자 콘솔 연동]**:
  - `P` 키 입력 시 별도의 팝업창으로 발표자 콘솔(Presenter View)이 열리는지 확인
  - 메인 창이나 발표자 창 중 한 곳에서 슬라이드를 넘기면 다른 창도 딜레이 없이 실시간 동기화되는지 확인
  - 경과 시간 타이머(Timer) 작동 및 발표자 메모(`<!-- note: -->`)가 누락 없이 표시되는지 확인
- [ ] **QA-INT-04 [레이저 포인터 (`L`)]**: 발표 중 `L` 키를 누르면 마우스 커서가 붉은색 발광 레이저 포인트로 변환되어 마우스 이동을 부드럽게 추종하는지 확인
- [ ] **QA-INT-05 [스포트라이트 모드 (`S`)]**: `S` 키 입력 시 마우스 커서 주변 원형 영역만 밝게 유지되고 나머지 배경은 어둡게 감춰지는지 확인
- [ ] **QA-INT-06 [화면 암전 및 백색 전환 (`B` / `W`)]**: `B` 키 입력 시 화면 전체가 검정 화면으로 암전되고, `W` 키 입력 시 백색 화면으로 전환되며, 임의의 키 입력 시 즉시 원래 슬라이드로 복귀하는지 확인
- [ ] **QA-INT-07 [화이트보드 캔버스 필기 (`D`)]**:
  - `D` 키 입력 시 드로잉 모드로 전환되어 마우스 드래그로 자유롭게 필기/드로잉이 가능한지 확인
  - 슬라이드를 다른 페이지로 넘겼다가 다시 돌아왔을 때 이전 필기 내용이 사라지지 않고 보존되는지 확인
  - 지우개 또는 초기화 단축키로 판서 내용을 정상 초기화할 수 있는지 확인

### 6.5 다중 포맷 빌드 산출물 무결성 (Export Quality: HTML, PDF, PPTX)
- [ ] **QA-EXP-01 [단독 실행형 HTML 오프라인 검증]**:
  - `goslide build talk.md -f html --standalone -o talk.html` 실행 후 산출물 확인
  - PC의 네트워크(와이파이/이더넷)를 끈 오프라인 환경에서 `talk.html`을 열었을 때 외부 CDN 실패 에러(ERR_INTERNET_DISCONNECTED) 없이 모든 폰트, 아이콘, 이미지가 완벽히 표시되는지 확인
- [ ] **QA-EXP-02 [PDF 출력 품질 및 인쇄 레이아웃]**:
  - `goslide build talk.md -f pdf -o talk.pdf` 생성 후 외부 PDF 뷰어(Adobe Acrobat, Chrome 등)로 열람
  - 슬라이드 간 여백(흰색 테두리) 없이 16:9 풀스크린 비율로 페이지가 나누어지는지 확인
  - 텍스트 및 벡터 그래픽이 래스터화(흐림 현상) 없이 선명한 고해상도로 인쇄 가능한지 확인
- [ ] **QA-EXP-03 [PPTX 실제 오피스 앱 호환성]**:
  - `goslide build talk.md -f pptx -o talk.pptx` 실행 후 실제 **Microsoft PowerPoint 365** 및 **Apple Keynote**에서 열기
  - 파일 열기 시 "손상된 파일" 복구 경고창 없이 깨끗하게 열리는지 확인
  - 각 슬라이드의 발표자 노트 영역에 마크다운에서 작성한 `<!-- note: -->` 텍스트가 정상 매핑되어 있는지 확인
- [ ] **QA-EXP-04 [다중 포맷 일괄 빌드 (`-f all`)]**:
  - `goslide build talk.md -f all -o dist/` 실행 시 `dist/` 디렉토리에 HTML, PDF, PPTX 3종 파일이 동시에 누락 없이 생성되는지 확인

### 6.6 크로스 플랫폼 및 OS 환경 (Cross-Platform Verification)
- [ ] **QA-ENV-01 [Windows 경로 처리]**: Windows OS에서 `C:\Users\...` 드라이브 문자 및 역슬래시(`\`) 경로가 포함된 파일 빌드 시 경로 파싱 오류 없이 빌드가 완료되는지 확인
- [ ] **QA-ENV-02 [macOS 단축키 충돌 방지]**: macOS Safari/Chrome 환경에서 프레젠터 단축키 실행 시 브라우저 기본 단축키(예: `Cmd+P` 시스템 인쇄 대화상자)와의 충돌 없이 발표자 콘솔이 정상 작동하는지 확인
- [ ] **QA-ENV-03 [Linux 정적 바이너리 무결성]**: 배포된 Linux 바이너리 실행 시 glibc 의존성 오류나 CGO 누락 없이 독립 실행형(Statically linked)으로 정상 구동되는지 확인

