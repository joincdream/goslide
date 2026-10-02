# [GOS-2] M1-1: 프로젝트 베이스라인 및 개발 인프라 구축 계획서

> **티켓 번호**: [GOS-2](https://joincdream.atlassian.net/browse/GOS-2)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-1  
> **마감일**: 2026-10-06  
> **상태**: 진행 중 (In Progress)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [architecture_design.md](file:///home/yundream/myjob/cloit/Goslide/docs/architecture_design.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md)

---

## 1. 개요 및 목적

본 태스크의 목적은 Goslide의 프로덕션 수준 개발 환경을 구성하고, 14개 핵심 파일 린 아키텍처에 맞춘 표준 Go 디렉토리 레이아웃과 일관된 빌드, 린트, 테스트 및 지속적 통합(CI) 자동화 파이프라인을 구축하는 것입니다.

---

## 2. 세부 개발 태스크 (Work Breakdown Structure)

### 2.1 Go 모듈 초기화 (`go.mod`)
- 모듈명: `github.com/yundream/goslide`
- Go 버전: `go 1.22` (또는 상위 버전)
- 핵심 의존성 계획:
  - CLI: `github.com/spf13/cobra`
  - 마크다운 AST: `github.com/yuin/goldmark`
  - 구문 강조: `github.com/alecthomas/chroma/v2` (Pure Go)
  - Frontmatter 파싱: `gopkg.in/yaml.v3`
  - 파일 감시: `github.com/fsnotify/fsnotify`
  - Headless 브라우징: `github.com/chromedp/chromedp` (로드맵 2 연계)
- 원칙: CGO 라이브러리 일체 배제 (100% Pure Go 보장)

### 2.2 14개 파일 린 패키지 레이아웃 스캐폴딩
`docs/architecture_design.md` 및 `docs/okf/core-architecture.md`에 정의된 패키지 디렉토리를 생성합니다:

```
Goslide/
├── cmd/
│   └── goslide/          # CLI 진입점 (main.go, root.go, build.go, serve.go)
├── pkg/
│   └── goslide/          # 퍼블릭 라이브러리 API (goslide.go)
├── internal/
│   ├── model/            # 도메인 IR (deck.go, slide.go, directive.go, interfaces.go, errors.go)
│   ├── parser/           # 파서 및 지시어 (parser.go, frontmatter.go, splitter.go, directive.go, highlight.go)
│   ├── theme/            # 테마 에셋 (embed.go, theme.go 및 assets/css/, assets/js/)
│   ├── renderer/
│   │   └── html/         # HTML 렌더러 (renderer.go, template.go)
│   ├── exporter/
│   │   ├── pdf/          # PDF 익스포터 (exporter.go)
│   │   └── pptx/         # PPTX 익스포터 (exporter.go)
│   ├── server/           # Live Preview 서버 (server.go, sse.go)
│   ├── i18n/             # 국제화 사전 및 엔진 (i18n.go, locales/en.json, locales/ko.json)
│   └── testutil/         # 테스트 헬퍼 (golden.go)
└── testdata/             # 골든 파일 픽스처 및 샘플 마크다운
```

### 2.3 `Makefile` 작성
개발 생산성과 로컬/CI 환경의 일관성을 위한 표준 타깃 정의:
- `make build`: `CGO_ENABLED=0` 바이너리 컴파일 (`bin/goslide`)
- `make test`: 단위 테스트 실행
- `make test-race`: 동시성 데이터 레이스 감지 테스트 (`go test -race -v ./...`)
- `make lint`: `golangci-lint` 정적 분석 실행 (복잡도 검사 포함)
- `make complexity`: 순환 복잡도(`gocyclo`) 및 인지 복잡도(`gocognit`) 정밀 측정 및 상위 랭킹 리포트 출력
- `make fmt`: `gofmt` 및 `goimports` 코드 스타일 정렬
- `make clean`: 빌드 산출물(`bin/`, `dist/`, `*.out`) 정리
- `make help`: 타깃 도움말 출력

### 2.4 정적 분석 린터 룰셋 구성 (`.golangci.yml`)
Go 표준 코딩 관용구와 잠재적 버그, 그리고 코드 복잡도를 사전에 차단하기 위한 린터 활성화:
- 기본 검사 린터:
  - `govet` (Go 표준 검사)
  - `errcheck` (반환 에러 누락 검사)
  - `staticcheck` (정적 버그 패턴)
  - `unused` (미사용 코드)
  - `gosimple` (코드 단순화 권고)
  - `ineffassign` (비효율적 변수 할당)
  - `gofmt` / `goimports` (포맷팅 검증)
- **코드 복잡도 관리 린터**:
  - `gocyclo`: 순환 복잡도 (최대 15)
  - `gocognit`: 인지 복잡도 (최대 20)
  - `nestif`: if문 중첩 깊이 (최대 4단계)
  - `funlen`: 함수 길이 제한 (최대 80줄 / 40 Statements)
- 제약 조건:
  - `panic` 방지 및 명시적 에러 반환 규칙 강제

### 2.5 GitHub Actions CI 워크플로우 구성 (`.github/workflows/ci.yml`)
- 트리거: `main` 브랜치 PR 및 Push
- 실행 환경: `ubuntu-latest`
- Go 버전 매트릭스: `1.22.x`, `1.23.x`
- 파이프라인 단계:
  1. 저장소 체크아웃 (`actions/checkout@v4`)
  2. Go 환경 설정 및 모듈 캐싱 (`actions/setup-go@v5`)
  3. 린터 검증 (`golangci/golangci-lint-action@v6`)
  4. 테스트 및 레이스 검증 (`go test -v -race -coverprofile=coverage.txt ./...`)
  5. 정적 빌드 무결성 검증 (`CGO_ENABLED=0 go build -v ./...`)

### 2.6 국제화(i18n) 인프라 구축 (`internal/i18n`)
글로벌 개발자와 한국어 사용자를 동시에 지원하기 위한 내장형 경량 i18n 아키텍처 수립:
- **원칙 (KISS & Pure Go)**: 외부 무거운 프레임워크 없이 Go `embed.FS`와 JSON 사전으로 구성.
- **주요 구성 요소**:
  - `internal/i18n/i18n.go`: 
    - 시스템 로케일 자동 감지 (`LANG`, `LC_ALL`) 및 명시적 설정 지원.
    - 번역 헬퍼 `T(key string, args ...any) string` 및 `SetLocale(lang string)` 구현.
    - 누락된 키 요청 시 폴백(`en` 사전 또는 원본 키 반환) 메커니즘.
  - `internal/i18n/locales/en.json`: 기본 영문 카탈로그 (CLI 도움말, 상태 메시지, 슬라이드 UI 텍스트).
  - `internal/i18n/locales/ko.json`: 한국어 번역 카탈로그.
- **CLI 글로벌 플래그**: Cobra 루트 커맨드에 `--lang` 플래그(`auto`, `en`, `ko`) 바인딩.
- **슬라이드 CJK 타이포그래피 사전 고려**: 기본 CSS에 `word-break: keep-all;` 및 Pretendard 한글 폰트 스택 기본 적용.

---

## 3. 엔지니어링 가드레일 (Hard Constraints)

1. **CGO 의존성 0%**: `CGO_ENABLED=0` 컴파일이 기본이어야 하며 CGO 라이브러리 참조 금지.
2. **순환 참조 방지 계층 준수**:
   - `internal/model`은 외부 및 내부 타 패키지를 절대 임포트하지 않음.
   - `internal/i18n`은 Go 표준 라이브러리(`embed`, `encoding/json`, `os`, `strings`) 외 어떠한 내부/외부 의존성도 갖지 않는 최하위 유틸리티 패키지로 격리.
   - `internal/parser`는 오직 `model`에만 의존.
   - `internal/renderer`는 `model`, `theme`, `i18n`만 의존.
3. **글로벌 상태 배제**: 전역 변수 설정 금지, 생성자 주입 원칙 준수 (단, CLI 수준의 기본 번역 엔진은 불변 싱글턴 인스턴스 또는 DI 패턴 지원).

---

## 4. 완료 기준 (Definition of Done)

- [x] `go.mod` 및 `go.sum`이 올바르게 생성되고 모듈 의존성이 동기화됨.
- [x] 린 아키텍처에 맞는 패키지 디렉토리 구조(`internal/i18n` 포함)가 생성됨.
- [x] `internal/i18n` 패키지가 구현되고 `en.json`, `ko.json`이 `embed.FS`로 번들링되어 `i18n.T()` 단위 테스트가 통과함.
- [x] `Makefile`의 모든 기본 타깃(`build`, `test`, `lint`, `complexity`, `clean`)이 정상 작동함.
- [x] `.golangci.yml` 파일이 유효하며 `golangci-lint run` 및 `go vet ./...`에서 경고가 0건임.
- [x] `make complexity` 실행 시 순환/인지 복잡도 임계치를 초과하지 않음.
- [x] `.github/workflows/ci.yml`이 유효한 GitHub Actions 문법으로 작성됨.
- [x] `CGO_ENABLED=0 go build ./...` 명령이 오류 없이 통과함.
- [x] `./bin/goslide --lang=ko --help` 및 `./bin/goslide --lang=en --help` 실행 시 각 언어별 도움말이 정상 출력됨.

---

## 5. 실행 절차 (Step-by-Step Execution Plan)

1. **디렉토리 레이아웃 스캐폴딩**: `mkdir -p`를 통해 패키지 구조(`internal/i18n/locales` 포함) 생성.
2. **Go 모듈 초기화**: `go mod init github.com/yundream/goslide`.
3. **i18n 경량 엔진 구현**: `internal/i18n/i18n.go` 및 영문/한글 JSON 카탈로그 작성, 단위 테스트 작성.
4. **기본 더미 패키지 파일 및 CLI 뼈대 배치**: `cmd/goslide`에서 `--lang` 플래그 및 i18n 연동.
5. **Makefile 작성**: `Makefile` 작성 및 타깃 동작 테스트.
6. **린터 및 CI 워크플로우 파일 작성**: `.golangci.yml` 및 `.github/workflows/ci.yml` 작성.
7. **검증 수행**: `make lint`, `go vet ./...`, `go test ./...`, `CGO_ENABLED=0 go build ./...`, 다국어 CLI 출력 검증.
8. **Jira 및 Git 작업 보고**: GOS-2 진행 상황 기록.
