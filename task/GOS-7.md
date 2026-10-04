# [GOS-7] M1-6: CLI 빌드 통합 및 MVP 게이트웨이 품질 검증 계획서

> **티켓 번호**: [GOS-7](https://joincdream.atlassian.net/browse/GOS-7)  
> **마일스톤**: 로드맵 1 (MVP) / Milestone M1-6 (로드맵 1 최종 관문)  
> **마감일**: 2026-10-23  
> **상태**: 진행 중 (In Progress)  
> **담당자**: Goslide Core Team  
> **참조 문서**: [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [core-architecture.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/core-architecture.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md), [decisions-simplification.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 지금까지 개발된 코어 IR(`model`), 마크다운/지시어 파서(`parser`), 테마 에셋 엔진(`theme`), HTML 렌더러 및 스크린캐스트 판서 런타임(`renderer/html`)을 **Cobra 기반 CLI 단일 진입점으로 온전히 통합**하고, **골든 파일 회귀 테스트 및 50장 대용량 성능 벤치마크**를 통해 **로드맵 1(MVP)의 최종 완료를 선언(Gateway Review)**하는 것입니다.

이로써 Goslide는 외부 런타임(Node.js 등)이나 CGO 없이, 누구나 터미널에서 단일 명령어(`goslide build`)로 전문 프레젠테이션 슬라이드를 밀리초 단위로 생성할 수 있는 **완전한 상용 수준의 MVP(Minimum Viable Product)** 단계에 도달합니다.

### 1.1 사용자 및 개발자에게 제공하는 가치
1. **일관되고 신뢰할 수 있는 CLI 경험**:
   - `goslide build <input.md> -o <out.html> [flags]` 명령을 통해 직관적이고 표준화된 인터페이스 제공.
   - `--theme`, `--theme-path`, `--standalone`, `--quiet`, `--verbose` 등 프로덕션 레벨의 플래그 바인딩.
   - 표준 Unix/POSIX 종료 코드(Exit Codes: `0` 성공, `2` 인자 오류, `3` 파일 없음, `4` 파싱 오류 등)를 지원하여 CI/CD 및 자동화 스크립트 연동성 보장.
2. **결정론적 출력 무결성 보장 (Golden File Regression Testing)**:
   - 다양한 마크다운 샘플(기본, 2단, 코드 구문 강조, 테마별)에 대해 골든 파일(`testdata/golden/*.html`)을 유지하여, 향후 어떠한 리팩토링이나 기능 추가 시에도 HTML 레이아웃이 깨지지 않음을 자동 검증.
   - `-update` 플래그를 통한 손쉬운 골든 파일 갱신 지원.
3. **극도의 고성능 보장 (Benchmark Guardrail)**:
   - 50장 대용량 슬라이드 변환 시간이 **500ms(0.5초) 미만**임을 자동화 벤치마크 테스트로 검증하여 "Go 기반 초경량/초고속 슬라이드 빌더"의 핵심 가치를 증명.
4. **로드맵 1 MVP 공식 완료 게이트웨이**:
   - 로드맵 1의 8대 완료 기준(DoD)을 전수 점검하여, 로드맵 2(PDF/PPTX 익스포터, 라이브 프리뷰 개발 서버, 발표자 뷰)로 안전하게 진입할 수 있는 기반 확립.

---

## 2. 결과물의 구체적 모습 (Concrete Manifestation)

### 2.1 CLI 빌드 인터랙션 명세
```bash
# 1. 기본 빌드 (입력 파일명 기반 HTML 자동 생성)
goslide build presentation.md
# 출력: ✓ Successfully built 12 slides to presentation.html [theme: clean] (Exit Code: 0)

# 2. 출력 파일 및 테마 명시 지정
goslide build presentation.md -o output/talk.html --theme=dark

# 3. 독립형 에셋 인라인 번들링 모드 (--standalone)
goslide build presentation.md --standalone -o dist/bundle.html

# 4. 외부 커스텀 기업 CSS 오버라이드
goslide build presentation.md --theme-path=./corporate.css

# 5. 비정상 상황 시 명확한 에러 메시지와 표준 종료 코드
goslide build not-found.md
# 출력: Error: input file "not-found.md" not found (Exit Code: 3)
```

### 2.2 CLI 표준 종료 코드 ([contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md))
| 종료 코드 | 상수 식별자 | 발생 조건 |
| :---: | :--- | :--- |
| **`0`** | `ExitSuccess` | 슬라이드 파싱 및 빌드 성공 완료 |
| **`1`** | `ExitGeneralError` | 분류되지 않은 런타임 오류 |
| **`2`** | `ExitInvalidUsage` | 필수 인자(`<input.md>`) 누락, 유효하지 않은 플래그 |
| **`3`** | `ExitFileNotFound` | 입력 마크다운 파일 미존재 |
| **`4`** | `ExitParseError` | YAML Frontmatter 파싱 실패 또는 마크다운 문법 오류 |

---

## 3. 핵심 설계 원칙 및 규칙 (Architectural Rules)

### 3.1 CLI 얇은 래퍼 원칙 (Thin CLI Wrapper)
- `cmd/goslide` 패키지는 순수하게 CLI 인자 파싱, 플래그 바인딩, i18n 메시지 출력 및 OS 종료 코드 제어만을 담당합니다.
- 실제 변환 비즈니스 로직은 `internal/parser`, `internal/theme`, `internal/renderer/html`의 조합으로만 오케스트레이션합니다.

### 3.2 골든 파일 회귀 검증 메커니즘
- `internal/testutil/golden.go`에 골든 파일 비교 헬퍼 작성:
  ```go
  func AssertGolden(t *testing.T, actual []byte, goldenPath string, update bool)
  ```
- 테스트 실행 시 `-update` 플래그(`go test ./... -update`)를 주면 실제 결과물로 골든 파일을 자동 갱신하고, 플래그가 없으면 골든 파일과 바이트 단위로 정확히 일치하는지 비교 검증합니다.
- HTML 내부의 타임스탬프 등 비결정론적 요소가 포함되지 않도록 결정론적 출력을 보장합니다.

### 3.3 대용량 벤치마크 테스트 규격
- `cmd/goslide/benchmark_test.go`에 `BenchmarkBuild_50Slides(b *testing.B)` 구현.
- 50장의 슬라이드(커버, 2단, 복잡한 코드블록, 텍스트)를 메모리에서 합성하거나 픽스처로 준비하여 파싱 $\rightarrow$ 스타일 합성 $\rightarrow$ HTML 렌더링 전체 파이프라인의 속도를 측정.
- **성능 게이트웨이 통과 기준**: 50장 기준 500ms 미만 (단일 슬라이드당 10ms 이하).

---

## 4. 세부 개발 태스크 (Work Breakdown Structure)

### 4.1 CLI 진입점 완성 및 종료 코드 표준화 (`cmd/goslide/`)
1. **`cmd/goslide/build.go`**:
   - 플래그 추가: `--quiet` (`-q`), `--verbose` (`-v`), `--standalone`, `--theme`, `--theme-path`.
   - 에러 발생 시 도메인 센티넬 에러(`os.IsNotExist`, `model.ErrInvalidFrontmatter` 등)를 식별하여 규격화된 종료 코드 반환 처리.
2. **`cmd/goslide/root.go` & `main.go`**:
   - `buildCmd` 헬프 및 사용법 설명 정비.
   - 비정상 종료 시 `os.Exit(code)`로 POSIX 표준 종료 코드 전파.

### 4.2 골든 파일 테스트 프레임워크 구축 (`internal/testutil/` & `testdata/`)
1. **`internal/testutil/golden.go`**:
   - `AssertGolden(t *testing.T, actual []byte, goldenPath string)` 구현.
   - `flag.Bool("update", false, "update golden files")` 연동.
2. **골든 파일 픽스처 마크다운 작성 (`testdata/slides/`)**:
   - `basic.md`: 일반 텍스트, H1~H3, 목록, 헤더/푸터.
   - `two_cols.md`: `<!-- layout: two-cols -->` 및 `<!-- split -->`.
   - `highlight.md`: Go, Python, JSON 코드 블록 Chroma 구문 강조.
   - `cover.md`: `<!-- layout: cover -->` 타이틀 슬라이드.
3. **골든 파일 생성 및 회귀 테스트 (`internal/renderer/html/golden_test.go`)**:
   - 4종 픽스처에 대해 `AssertGolden` 검증 수행.
   - `testdata/golden/*.html` 레퍼런스 파일 생성.

### 4.3 성능 벤치마크 테스트 작성 (`cmd/goslide/benchmark_test.go`)
1. 50장 슬라이드로 구성된 대형 마크다운 데크 생성 헬퍼.
2. 파서 + 렌더러 파이프라인 벤치마크 함수 `BenchmarkBuild_50Slides`.
3. 50장 기준 변환 시간이 500ms(0.5초) 미만인지 확인하는 어설션 테스트 `TestPerformance_50SlidesUnder500ms` 추가.

### 4.4 Makefile 타깃 확장
1. `make golden-update`: 골든 파일 일괄 갱신 타깃 (`go test ./... -update`).
2. `make bench`: 벤치마크 실행 타깃 (`go test -bench=. -benchmem ./cmd/goslide/...`).

### 4.5 로드맵 1 MVP 게이트웨이 검증 (Gate 1 DoD Review)
- 로드맵 1 8대 DoD 체크리스트 전수 점검 및 검증 보고서 작성.

---

## 5. 엔지니어링 가드레일 (Hard Constraints)

1. **CGO 배제 (100% Pure Go)**: `CGO_ENABLED=0` 환경에서 단일 바이너리 컴파일 보장.
2. **결정론적 출력 (Deterministic Output)**: 동일한 마크다운 입력에 대해 항상 100% 동일한 HTML 바이트 생성.
3. **메모리 및 성능 가드레일**: 50장 슬라이드 빌드 메모리 피크 100MB 이하, 처리 시간 500ms 미만.
4. **글로벌 상태 금지 및 Context 준수**: CLI 오케스트레이션에서도 패키지 전역 변수 공유 배제.

---

## 6. 완료 기준 (Definition of Done)

- [ ] `cmd/goslide/build.go`에서 `--theme`, `--theme-path`, `--standalone`, `--quiet` 플래그가 온전히 작동함.
- [ ] 입력 파일 미존재 시 `ExitFileNotFound(3)`, 파싱 에러 시 `ExitParseError(4)` 등 표준 종료 코드가 반환됨.
- [ ] `internal/testutil/golden.go` 골든 파일 헬퍼가 구현되고 `-update` 플래그를 지원함.
- [ ] `testdata/slides/` 픽스처 및 `testdata/golden/` 레퍼런스 HTML이 구축되고 회귀 테스트를 100% 통과함.
- [ ] 50장 대용량 슬라이드 빌드 성능이 500ms 미만을 달성함 (`benchmark_test.go`).
- [ ] `make test`, `make test-race`, `make lint`, `make complexity` 검사를 무결하게 통과함.
- [ ] 로드맵 1 게이트웨이 리뷰 체크리스트 8개 항목이 모두 충족됨.

---

## 7. 단계별 실행 계획 (Step-by-Step Execution Plan)

1. **골든 파일 헬퍼 구현**: `internal/testutil/golden.go` 작성.
2. **골든 픽스처 및 회귀 테스트 작성**: `testdata/slides/*.md` 작성, `internal/renderer/html/golden_test.go` 작성 및 레퍼런스 생성.
3. **대용량 성능 벤치마크 구현**: `cmd/goslide/benchmark_test.go` 작성 및 500ms 이내 검증.
4. **CLI 플래그 및 표준 종료 코드 고도화**: `cmd/goslide/build.go`, `root.go` 리팩토링 및 `build_test.go` 테스트 추가.
5. **Makefile 타깃 정비**: `golden-update`, `bench` 추가.
6. **품질 검증**: `make test`, `make test-race`, `make lint`, `make complexity`, `make bench`.
7. **로드맵 1 게이트웨이 리뷰 및 완료 보고**: MVP 완료 승인 및 다음 로드맵 2 안내.
