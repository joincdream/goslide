# [GOS-12] CLI 편의 기능 및 GoReleaser 크로스 배포 자동화 (v1.0.0 릴리즈) 개발 계획서

> **티켓 번호**: [GOS-12](https://joincdream.atlassian.net/browse/GOS-12)  
> **마일스톤**: 로드맵 2 (Release) / Milestone M2-5 (CLI 고도화, GoReleaser 배포 및 공식 릴리즈)  
> **마감일**: 2026-11-13  
> **상태**: 진행 중 (In Progress)  
> **담당 패키지**: `cmd/goslide/`, `pkg/goslide/`, `.goreleaser.yaml`, `.github/workflows/release.yml`  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md), [AGENTS.md](file:///home/yundream/myjob/cloit/Goslide/AGENTS.md)

---

## 1. 개요 및 목적 (Overview & Goals)

본 태스크의 목적은 Goslide를 독립 실행 가능한 완결된 단일 정적 바이너리로 패키징하고, 서드파티 개발자와 일반 사용자가 손쉽게 프레젠테이션을 시작·빌드·배포할 수 있도록 **CLI 편의 기능(`goslide init`), 퍼블릭 라이브러리 API(`pkg/goslide`), 다중 포맷 일괄 빌드(`-f html,pdf,pptx`), 그리고 5대 주요 OS/Arch 자동 릴리즈 파이프라인(GoReleaser + GitHub Actions)**을 구축하여 공식 **`v1.0.0`** 정식 버전을 성공적으로 출시하는 것입니다.

### 핵심 달성 목표
1. **신속한 슬라이드 작성 진입 (`goslide init`)**: 빈 파일 대신 Goslide 표준 DSL 문법(커버, 2단 레이아웃, 코드 하이라이트)이 포함된 마크다운 스타터 생성.
2. **Go 생태계 라이브러리 연동 (`pkg/goslide`)**: 다른 Go 애플리케이션에서 `goslide`를 임포트하여 마크다운 파싱 및 슬라이드 빌드를 호출할 수 있는 Public Facade API 확립.
3. **단일 명령 다중 포맷 출력**: `goslide build talk.md -f html,pdf,pptx -o dist/`로 3대 산출물(HTML, PDF, PPTX)을 단 한 번의 명령으로 일괄 생성.
4. **글로벌 크로스 플랫폼 패키징 (GoReleaser)**: Linux(amd64/arm64), macOS(Intel/Apple Silicon), Windows(amd64) 5종 아키텍처에 대해 CGO 0% 무의존 단일 바이너리 자동 빌드 및 체크섬 생성.
5. **v1.0.0 출시 품질 보증 (DoD)**: 50장 이상 대규모 슬라이드 일괄 빌드 시 메모리 피크 100MB 이하 검증 및 동시성 데이터 레이스 0건 달성.

---

## 2. 세부 설계 명세 (Specification)

### 2.1 `goslide init` 템플릿 생성 커맨드 (`cmd/goslide/init.go`)

사용자가 문법 가이드를 일일이 찾아보지 않고도 즉시 작동하는 슬라이드를 생성할 수 있도록 지원합니다.

```bash
# 기본 clean 테마의 presentation.md 생성
goslide init

# 특정 파일명 및 테마 지정 생성
goslide init talk.md --theme=dark

# 기존 파일이 있을 경우 덮어쓰기 방지 (강제 시 --force)
goslide init presentation.md --force
```

#### 옵션 플래그
* `-t, --theme string`: 스타터 템플릿에 적용할 기본 테마 (기본값: `clean`)
* `-f, --force`: 이미 파일이 존재하는 경우 경고 없이 덮어쓰기
* `-h, --help`: init 커맨드 도움말 출력

#### 템플릿 임베딩 아키텍처 (Single Source of Truth)
Go 소스 코드 내부에 마크다운 문자열을 하드코딩하지 않고, 실제 검증 가능한 독립 마크다운 파일(`starter.md`)로 관리하여 Go 바이너리에 임베딩합니다:
* **독립 템플릿 파일**: [`internal/theme/assets/templates/starter.md`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/templates/starter.md)
* **바이너리 임베딩 선언**: `internal/theme/embed.go` (`//go:embed assets/css/*.css assets/js/*.js assets/templates/*.md`)
* **템플릿 제공 메서드**: `theme.Manager.GetStarterTemplate(theme string)`가 템플릿을 로드하고 지정된 테마명을 Frontmatter에 주입하여 반환.
* **공통 활용 보장**: 터미널 CLI `goslide init`과 웹 브라우저 Deck Hub의 "신규 데크 생성(`POST /api/v1/files`)"이 동일한 임베딩 템플릿을 공유하여 완벽한 일관성 유지.

#### 생성되는 스타터 마크다운 구성 (`starter.md`)
```markdown
---
title: "Presentation Title"
theme: "clean"
paginate: true
---

<!-- _layout: cover -->
# Presentation Title

Sub-heading or Presenter Introduction

---

<!-- _layout: two-cols -->
## 2-Column Layout Demonstration

Left column points:
* High-Performance Pure Go
* Zero-Node & Zero-CGO Runtime

<!-- split -->

Right column points:
* 16:9 Borderless Vector PDF
* Editable PowerPoint Export

---

## Code Highlighting & Presentation Shortcuts

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Goslide!")
}
```

<!-- note:
이곳은 발표자 노트 영역입니다.
'P' 키를 누르면 발표자 콘솔에서 확인할 수 있습니다.
-->
```

---

### 2.2 서드파티 연동용 Public Go API (`pkg/goslide/goslide.go`)

외부 Go 프로젝트에서 CLI 바이너리를 통하지 않고 직접 Go 코드로 슬라이드를 렌더링하고 파싱할 수 있는 안정적인 공개 API를 제공합니다:

```go
package goslide

import (
	"context"
	"io"

	"github.com/yundream/goslide/internal/model"
)

// Version은 현재 Goslide 라이브러리의 릴리즈 버전 문자열입니다.
const Version = "v1.0.0"

// BuildOption은 빌드 파이프라인 설정을 위한 함수형 옵션입니다.
type BuildOption func(*buildOptions)

// Build는 지정된 마크다운 파일을 파싱하여 요청된 포맷(HTML, PDF, PPTX)으로 컴파일합니다.
func Build(ctx context.Context, inputPath, outputPath string, opts ...BuildOption) error

// Parse는 io.Reader로부터 마크다운 스트림을 읽어 불변 Deck 도메인 모델을 생성합니다.
func Parse(ctx context.Context, r io.Reader) (*model.Deck, error)
```

#### 제공 옵션 함수
* `WithFormat(fmt string)`: 출력 포맷 지정 (`html`, `pdf`, `pptx`)
* `WithTheme(theme string)`: 적용 테마 지정 (`clean`, `dark`, `default`)
* `WithCustomCSS(cssPath string)`: 커스텀 외부 CSS 스타일시트 경로 지정
* `WithStandalone(standalone bool)`: 로컬 이미지 에셋 Base64 인라인 임베딩 여부

---

### 2.3 다중 포맷 일괄 빌드 확장 (`cmd/goslide/build.go`)

현재 단일 포맷만 허용하던 `-f, --format` 플래그를 확장하여, 쉼표(`,`)로 구분된 다중 포맷 인자를 처리할 수 있도록 리팩토링합니다.

#### 동작 사양
* 입력: `goslide build talk.md -f html,pdf,pptx -o dist/`
* 처리 절차:
  1. 마크다운 파일 파싱은 **단 1회만 수행**하여 메모리에 불변 `model.Deck` IR 적재 (중복 파싱 방지).
  2. `-o` 대상이 디렉토리인 경우 `dist/talk.html`, `dist/talk.pdf`, `dist/talk.pptx`로 자동 확장자 분기.
  3. `-o` 대상이 지정되지 않은 경우 원본 파일 기준 디렉토리에 각 포맷별 파일 생성.
  4. 각 포맷 렌더러/익스포터를 순차적으로 안전하게 실행하고 최종 생성 결과 요약 출력.

---

### 2.4 GoReleaser 크로스 컴파일 및 배포 자동화

#### 1) `.goreleaser.yaml` 설정 사양
* **프로젝트명**: `goslide`
* **빌드 바이너리 매트릭스**:
  * `linux/amd64`, `linux/arm64`
  * `darwin/amd64`, `darwin/arm64` (Universal 바이너리 병합 지원)
  * `windows/amd64`
* **빌드 플래그**:
  * `CGO_ENABLED=0` (엄격한 Pure Go 정적 링크)
  * `-ldflags="-s -w -X github.com/yundream/goslide/pkg/goslide.Version={{.Version}}"`
* **아카이브 및 체크섬**:
  * `.tar.gz` (Linux, macOS), `.zip` (Windows)
  * `checksums.txt` SHA-256 자동 계산 및 동봉

#### 2) `.github/workflows/release.yml` GitHub Actions 워크플로우
* **트리거**: `git push origin v*.*.*` 태그 푸시 시 자동 실행.
* **작업 단계**:
  1. `actions/checkout@v4` (전체 태그 페치)
  2. `actions/setup-go@v5` (Go 1.23+ 설정)
  3. `actions/setup-node@v4` & `npm run build` (Svelte 5 웹 에셋 번들 생성 및 임베드)
  4. `goreleaser/goreleaser-action@v5` 실행 및 GitHub Releases 자동 배포

---

## 3. 단계별 개발 일정 및 세부 태스크 (Execution Plan)

| 단계 | 작업 내용 | 담당 파일 | 예상 산출물 |
| :---: | :--- | :--- | :--- |
| **Phase 1** | **스타터 덱 임베딩 & `goslide init` 커맨드 구현**<br>• `assets/templates/starter.md` 독립 템플릿 파일 생성 및 `embed.FS` 바인딩<br>• `theme.Manager.GetStarterTemplate(theme)` 메서드 구현<br>• `initCmd` 선언, 파일 덮어쓰기 방지 및 `--force` 옵션 테스트 | `internal/theme/assets/templates/starter.md`<br>`internal/theme/embed.go`<br>`internal/theme/theme.go`<br>`cmd/goslide/init.go`<br>`cmd/goslide/init_test.go` | `goslide init` 및 임베딩 템플릿 완성 |
| **Phase 2** | **Public Facade API 확립**<br>• 서드파티 연동용 `goslide.Build`, `goslide.Parse` 구현<br>• 함수형 옵션 패턴(`With...`) 정의<br>• 라이브러리 임포트 통합 단위 테스트 | `pkg/goslide/goslide.go`<br>`pkg/goslide/goslide_test.go` | `pkg/goslide` 공개 API |
| **Phase 3** | **다중 포맷 일괄 빌드 파이프라인 연동**<br>• `-f html,pdf,pptx` 쉼표 구분 파싱 로직 추가<br>• `Deck` 1회 파싱 후 각 익스포터 순차 파이프라인 체이닝<br>• 디렉토리 대상 출력 경로 자동 매핑 | `cmd/goslide/build.go`<br>`cmd/goslide/build_test.go` | 일괄 빌드 기능 |
| **Phase 4** | **GoReleaser 및 CI/CD 워크플로우 구축**<br>• 5개 플랫폼 타겟 `.goreleaser.yaml` 작성<br>• `.github/workflows/release.yml` 생성<br>• `goreleaser check` 정적 문법 유효성 검증 | `.goreleaser.yaml`<br>`.github/workflows/release.yml` | 크로스 배포 자동화 |
| **Phase 5** | **종합 회귀 테스트 & v1.0.0 릴리즈**<br>• 50장 슬라이드 일괄 빌드 메모리/성능 검증 (< 100MB)<br>• `go test -race ./...` 동시성 안전성 검증<br>• v1.0.0 릴리즈 노트 작성 및 태그 발행 준비 | 전체 패키지<br>`docs/release-notes-v1.0.0.md` | v1.0.0 정식 릴리즈 완료 |

---

## 4. 완료 기준 (Definition of Done - DoD)

1. **CLI init 정상 동작**:
   * 임의의 빈 폴더에서 `goslide init` 실행 시 깨지지 않는 표준 스타터 슬라이드 생성 확인.
2. **다중 포맷 일괄 빌드 검증**:
   * `goslide build examples/example-dsl.md -f html,pdf,pptx -o /tmp/dist/` 실행 시 HTML, PDF, PPTX 3종 파일이 정상 생성되는지 확인.
3. **Public API 테스트 통과**:
   * `go test -v ./pkg/goslide/...` 단위 테스트 100% 통과.
4. **GoReleaser 유효성 통과**:
   * 로컬에서 `goreleaser check` 실행 시 에러 없이 통과.
5. **동시성 및 메모리 안정성**:
   * `go test -race ./...` 실행 시 레이스 컨디션 0건.
   * 50장 슬라이드 일괄 빌드 시 메모리 피크 100MB 이하 유지.
