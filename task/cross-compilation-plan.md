# Goslide 크로스 컴파일 및 멀티 플랫폼 배포 파이프라인 개발 계획서

> **문서 상태**: 계획 수립 (Planning)  
> **연관 티켓**: [GOS-12](file:///home/yundream/myjob/cloit/Goslide/task/GOS-12.md) (Phase 4: GoReleaser 크로스 컴파일 및 배포 자동화)  
> **마일스톤**: 로드맵 2 (Release) / M2-5 (공식 v1.0.0 릴리즈)  
> **담당 패키지**: `pkg/goslide/`, `cmd/goslide/`, `Makefile`, `.goreleaser.yaml`, `.github/workflows/release.yml`  
> **핵심 가치**: **CGO 0% Pure Go**, **단일 바이너리(Single Binary)**, **5대 메이저 OS/Arch 동시 지원**, **SemVer 기반 버전 관리**

---

## 1. 개요 및 목적 (Overview & Goals)

본 태스크는 Goslide를 Linux, macOS, Windows 등 개발자와 발표자가 사용하는 모든 주요 데스크톱 및 서버 환경에서 **별도의 의존성(Node.js, C 컴파일러, 파이썬 등) 없이 단일 바이너리만 다운로드하여 즉시 실행**할 수 있도록, **로컬 크로스 컴파일 메커니즘**과 **GoReleaser + GitHub Actions 기반 글로벌 자동 릴리즈 파이프라인**을 구축하는 것을 목표로 합니다.

특히 사용자의 설치 환경과 사용 목적에 맞추어 **단일 독립 실행파일(Raw Standalone Binary)**과 **문서/예제가 동봉된 압축 아카이브(Compressed Archive)**를 동시에 배포하는 **듀얼 패키징(Dual Packaging) 전략**을 채택하고, 컴파일 타임 메타데이터 주입 기반의 **시맨틱 버저닝(Semantic Versioning)** 체계를 정립합니다.

### 핵심 지원 플랫폼 매트릭스 (Target Matrix)
| 운영체제 (OS)             | 아키텍처 (Arch) | 단일 바이너리 산출물 (Raw)           | 배포 아카이브 (Archive)             | 주요 타겟 환경                                     |     |
| :-------------------- | :---------- | :-------------------------- | :---------------------------- | :------------------------------------------- | --- |
| **Linux**             | `amd64`     | `goslide-linux-amd64`       | `goslide_linux_amd64.tar.gz`  | 일반 x86_64 서버, Ubuntu/Debian/RHEL 데스크톱, WSL2  |     |
| **Linux**             | `arm64`     | `goslide-linux-arm64`       | `goslide_linux_arm64.tar.gz`  | AWS Graviton, Raspberry Pi 4/5, Linux ARM 서버 |     |
| **macOS**             | `arm64`     | `goslide-darwin-arm64`      | `goslide_darwin_arm64.tar.gz` | Apple Silicon (M1/M2/M3/M4) Mac              |     |
| **macOS**             | `amd64`     | `goslide-darwin-amd64`      | `goslide_darwin_amd64.tar.gz` | Intel 기반 레거시 Mac                             |     |
| **macOS (Universal)** | `universal` | `goslide-darwin-universal`  | `goslide_darwin_all.tar.gz`   | `lipo`로 통합된 macOS 단일 범용 바이너리                 |     |
| **Windows**           | `amd64`     | `goslide-windows-amd64.exe` | `goslide_windows_amd64.zip`   | 64비트 Windows 10/11 데스크톱 및 워크스테이션             |     |

---

## 2. 세부 설계 사양 (Technical Specification)

### 2.1 사전 웹 번들 컴파일 및 정적 임베딩 파이프라인
Goslide의 크로스 컴파일은 Go 소스 코드뿐만 아니라 프론트엔드 에셋(Svelte 5 + Tailwind)이 사전 빌드되어 `internal/theme/assets/`에 포함되어 있어야 합니다.

```
[web/ (Svelte 5)] ➔ npm run build ➔ internal/theme/assets/js/goslide-core.js
                                  ➔ internal/theme/assets/css/presenter.css
                                          │
                                          ▼
                      [Go embed.FS (internal/theme/embed.go)]
                                          │
                     ┌────────────────────┼────────────────────┐
                     ▼                    ▼                    ▼
             [Linux Binary]         [macOS Binary]       [Windows .exe]
             (CGO_ENABLED=0)        (CGO_ENABLED=0)      (CGO_ENABLED=0)
```

* **보장 규칙**: 크로스 컴파일 실행 전 반드시 `make build-web`이 선행되어 최신 프론트엔드 번들이 바이너리에 번들링되어야 합니다.

---

### 2.2 버전 관리 체계 및 릴리즈 라이프사이클 (Version Management & Lifecycle)

Goslide는 소스 코드 내 하드코딩된 버전을 지양하고, Git Tag를 단일 진실 공급원(Single Source of Truth)으로 삼아 빌드 시점에 링커 플래그(`-ldflags`)를 통해 메타데이터를 주입합니다.

#### 1) 시맨틱 버저닝 (SemVer 2.0.0) 규칙
버전 번호는 `vMAJOR.MINOR.PATCH` 형식을 엄격히 따릅니다:
* **MAJOR (`v2.0.0`)**: CLI 플래그, 마크다운 DSL 지시어 문법의 하위 호환성을 깨는 변경 (Breaking Changes)
* **MINOR (`v1.1.0`)**: 새로운 렌더러/익스포터 추가, 신규 지시어 지원, 하위 호환되는 기능 확장
* **PATCH (`v1.0.1`)**: 하위 호환되는 버그 수정, 성능 최적화, 보안 패치
* **Pre-release (`v1.0.0-rc.1`, `v1.0.0-beta.1`)**: 공식 릴리즈 전 커뮤니티 및 검증 단계 배포

#### 2) 메타데이터 변수 선언 (`pkg/goslide/goslide.go`)
`const` 대신 `var`로 선언하여 컴파일 타임 주입(`-X`)을 허용하며, 진단용 헬퍼 함수를 제공합니다:
```go
package goslide

import (
	"fmt"
	"runtime"
)

var (
	// Version은 릴리즈 태그 (예: v1.0.0, 로컬 빌드 시 v1.0.0-dev)
	Version = "v1.0.0-dev"
	// Commit은 Git 커밋 SHA 해시
	Commit = "none"
	// Date는 빌드된 RFC3339 일시
	Date = "unknown"
)

// FullVersionString은 CLI -v 또는 시스템 진단 시 출력할 포맷팅된 버전 문자열을 반환합니다.
func FullVersionString() string {
	return fmt.Sprintf("Goslide %s (commit: %s, built at: %s, %s/%s)",
		Version, Commit, Date, runtime.GOOS, runtime.GOARCH)
}
```

#### 3) Makefile 및 링커 플래그 연동
로컬 빌드(`make build`, `make cross-build`) 시 Git 메타데이터를 자동 추출하여 주입합니다:
```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0-dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS = -s -w \
	-X github.com/yundream/goslide/pkg/goslide.Version=$(VERSION) \
	-X github.com/yundream/goslide/pkg/goslide.Commit=$(COMMIT) \
	-X github.com/yundream/goslide/pkg/goslide.Date=$(DATE)
```

#### 4) CLI 버전 출력 사양 (`cmd/goslide/root.go`)
사용자가 `-v` 또는 `--version` 플래그를 전달할 때 빌드 메타데이터를 가독성 높게 출력합니다:
```bash
$ goslide --version
Goslide v1.0.0 (commit: 0c07bd9, built at: 2026-10-06T15:43:00Z, linux/amd64)
```

#### 5) 버전 태그 릴리즈 런북 (Release Runbook)
신규 버전을 배포하는 표준 워크플로우:
1. **사전 검증**: `make test` 및 `make test-race` 통과 확인.
2. **웹 번들 검증**: `make build-web` 최신 프론트엔드 에셋 컴파일 확인.
3. **Git 태그 생성**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0: First stable release with HTML/PDF/PPTX support"
   ```
4. **원격 태그 푸시**:
   ```bash
   git push origin v1.0.0
   ```
5. **CI/CD 자동 파이프라인 트리거**: GitHub Actions가 태그 이벤트를 감지하여 5대 플랫폼 바이너리 빌드, `checksums.txt` 생성, 릴리즈 노트 자동 추출 후 GitHub Releases에 에셋 자동 업로드.

---

### 2.3 배포 패키징 전략: 단일 바이너리 + 압축 아카이브 듀얼 배포 (Dual Distribution Policy)

사용자의 다양한 설치 환경과 편의성을 극대화하기 위해 **단일 바이너리(Standalone Executable)**와 **압축 아카이브(Archive)** 두 가지 산출물을 모두 릴리즈 에셋으로 동시 제공합니다.

#### 1) 단일 바이너리 (Raw Standalone Binary) 배포의 필요성 및 활용
* **형태**: 압축 없는 독립 실행 파일 (`goslide-linux-amd64`, `goslide-windows-amd64.exe`, `goslide-darwin-universal`)
* **장점 및 사용처**:
  * **초간편 단일 명령어 설치**: `curl` 또는 `wget`으로 즉시 다운로드하여 `PATH`에 등록 가능:
    ```bash
    curl -Lo /usr/local/bin/goslide https://github.com/yundream/goslide/releases/download/v1.0.0/goslide-linux-amd64
    chmod +x /usr/local/bin/goslide
    ```
  * **CI/CD 및 컨테이너 최적화**: Dockerfile 및 CI 파이프라인에서 압축 해제 도구(`tar`, `unzip`) 설치 없이 곧바로 바이너리를 레이어에 추가 가능.

#### 2) 압축 아카이브 (Archive: `.tar.gz` / `.zip`) 배포의 필요성 및 활용
* **형태**: 바이너리와 메타데이터/예제가 함께 묶인 압축 파일 (`goslide_1.0.0_linux_amd64.tar.gz`, `goslide_1.0.0_windows_amd64.zip`)
* **장점 및 사용처**:
  * **필수 문서 및 기본 예제 동봉**: `README.md`, `LICENSE`, `examples/example-dsl.md`가 함께 압축되어 첫 사용자 온보딩 경험 제공.
  * **Unix 실행 권한(`chmod +x`) 유지**: `tar.gz` 아카이빙은 파일의 실행 권한 비트를 보존하므로 사용자가 압축 해제 후 별도의 `chmod +x`를 수행할 필요가 없음.
  * **패키지 관리자 공식 배포 표준**: Homebrew (`brew install yundream/tap/goslide`), Windows Scoop, Arch AUR 등 공식 패키지 매니저는 압축 아카이브 URL과 SHA-256 체크섬을 배포 규격으로 요구함.
  * **Windows 다운로드 보안 경고 완화**: 원시 `.exe` 직접 다운로드 시 발생하는 웹 브라우저/SmartScreen 보안 경고를 완화하고 다운로드 파일 크기를 최적화.

---

### 2.4 로컬 개발용 크로스 빌드 명령 (`Makefile`)

개발자가 로컬 머신에서 직접 모든 타겟 플랫폼 바이너리를 한 번에 빌드하고 검증할 수 있도록 `Makefile`을 확장합니다:

```makefile
## cross-build: Build static binaries for all 5 target OS/Architectures
cross-build: build-web
	@echo "==> Cross-compiling Goslide for all target platforms (Version: $(VERSION))..."
	@mkdir -p $(BIN_DIR)
	# Linux amd64 & arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/goslide
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/goslide
	# macOS amd64 & arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/goslide
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/goslide
	# Windows amd64 (.exe)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/goslide
	@echo "==> Generating SHA256 checksums..."
	@cd $(BIN_DIR) && sha256sum $(BINARY_NAME)-* > checksums.txt
	@echo "🎉 Cross-compilation completed successfully!"
```

---

### 2.5 GoReleaser 듀얼 패키징 자동화 사양 (`.goreleaser.yaml`)

GoReleaser v2 표준을 사용하여 태그 푸시 시 **압축 아카이브**와 **단일 독립 바이너리**를 동시에 빌드하고 릴리즈 에셋으로 업로드합니다:

```yaml
version: 2

project_name: goslide

before:
  hooks:
    - cd web && npm ci && npm run build
    - go mod tidy

builds:
  - id: goslide
    main: ./cmd/goslide
    binary: goslide
    env:
      - CGO_ENABLED=0
    flags:
      - -trimpath
    ldflags:
      - -s -w
      - -X github.com/yundream/goslide/pkg/goslide.Version={{.Version}}
      - -X github.com/yundream/goslide/pkg/goslide.Commit={{.Commit}}
      - -X github.com/yundream/goslide/pkg/goslide.Date={{.Date}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64

universal_binaries:
  - id: goslide-darwin-universal
    ids:
      - goslide
    name_template: "goslide"
    replace: false

archives:
  # 1. 문서 및 예제가 동봉된 압축 아카이브 배포
  - id: default-archive
    builds:
      - goslide
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    files:
      - README.md
      - LICENSE*
      - examples/example-dsl.md
      - examples/example-dsl.ko.md

  # 2. curl/wget 원클릭 다운로드를 위한 단일 무압축 바이너리 배포
  - id: raw-binary
    builds:
      - goslide
    name_template: "{{ .Binary }}-{{ .Os }}-{{ .Arch }}"
    format: binary

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

---

### 2.6 GitHub Actions 릴리즈 파이프라인 (`.github/workflows/release.yml`)

`git push origin v*.*.*` 태그 푸시 시 트리거되는 CI/CD 배포 워크플로우:

```yaml
name: Release Goslide

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: web/package-lock.json

      - name: Build Web Frontend Assets
        run: |
          cd web
          npm ci
          npm run build

      - name: Run Unit Tests
        run: go test -v -race ./...

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 3. 단계별 실행 계획 (Execution Phases)

| 단계 | 작업 내용 | 담당 파일 | 예상 산출물 |
| :---: | :--- | :--- | :--- |
| **Phase 1** | **버전 관리 변수 및 CLI 메타데이터 주입 구현**<br>• `pkg/goslide/goslide.go`에 `var Version, Commit, Date` 및 `FullVersionString()` 정의<br>• `cmd/goslide/root.go`에 `--version` 상세 출력 로직 반영<br>• `Makefile`에 `LDFLAGS` 및 `cross-build` 타겟 구현 및 로컬 정적 바이너리 검증 | `pkg/goslide/goslide.go`<br>`cmd/goslide/root.go`<br>`Makefile` | 버전 정보 주입 및 로컬 5대 플랫폼 일괄 빌드 기능 |
| **Phase 2** | **GoReleaser 듀얼 패키징 설정 파일 구현**<br>• `.goreleaser.yaml` 작성<br>• 단일 바이너리(`raw-binary`) + 압축 아카이브(`default-archive`) 동시 생성 정의<br>• macOS Universal 바이너리 설정 및 `goreleaser check` 유효성 검증 | `.goreleaser.yaml` | GoReleaser 릴리즈 스펙 정의 |
| **Phase 3** | **GitHub Actions 배포 워크플로우 구성**<br>• `.github/workflows/release.yml` 작성<br>• Node 20 & Go 1.23 빌드 캐싱 및 자동 릴리즈 액션 바인딩<br>• Dry-run 로컬 테스트 (`goreleaser release --snapshot --clean`) | `.github/workflows/release.yml` | GitHub 자동 릴리즈 파이프라인 |
| **Phase 4** | **크로스 플랫폼 실행 및 바이너리 무결성 검증**<br>• Windows `.exe` 실행 파일 헤더 및 바이너리 무결성 검사<br>• Linux x86_64 / ARM64 아키텍처 바이너리 심볼 확인<br>• macOS Universal (`x86_64` + `arm64`) 바이너리 번들 무결성 확인 | 전체 빌드 산출물 (`bin/`) | 전 플랫폼 바이너리 검증 완료 |

---

## 4. 완료 기준 (Definition of Done - DoD)

1. **버전 메타데이터 주입 확인**:
   * 로컬 빌드 및 크로스 빌드 바이너리에서 `goslide --version` 실행 시 실제 Git 태그, 커밋 해시, 빌드 일시, OS/Arch가 정확히 출력됨.
2. **듀얼 패키징 산출물 무결성**:
   * 단일 실행파일(`goslide-linux-amd64` 등)과 압축 아카이브(`.tar.gz`, `.zip`)가 누락 없이 생성되고 올바른 체크섬(`checksums.txt`)이 생성됨.
3. **CGO 의존성 0% 보증**:
   * 생성된 바이너리들이 동적 공유 라이브러리 의존 없이 순수 정적 링크(`statically linked`) 상태임을 확인.
4. **GoReleaser 유효성 통과**:
   * `goreleaser check` 명령 실행 시 스키마 에러 0건.
5. **Windows 실행 파일 적합성**:
   * `goslide-windows-amd64.exe`가 정상 PE32+ 실행 파일 형식으로 생성됨을 확인.
6. **CI/CD 워크플로우 통과**:
   * `.github/workflows/release.yml` 문법 검증 및 태그 트리거 정상 바인딩 확인.
