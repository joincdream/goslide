# Goslide 디렉토리 구조 정리 및 유지보수성 개선 계획서

> **문서 상태**: 계획 수립 (Planning)  
> **연관 티켓**: [GOS-25](https://joincdream.atlassian.net/browse/GOS-25)  
> **마일스톤**: 유지보수성 및 프로젝트 구조 최적화  
> **대상 디렉토리**: 루트(`.`), `task/`, `testdata/`, `examples/`, `docs/`, `internal/`  
> **핵심 가치**: **클린 워크스페이스(Clean Workspace)**, **테스트 격리(Test Isolation)**, **직관적 온보딩(Clear Onboarding)**, **문서 위계화(Modular Documentation)**

---

## 1. 개요 및 배경 (Context & Problem Statement)

Goslide는 Go 표준 패키지 레이아웃(`cmd/`, `internal/`, `pkg/`)과 단방향 파이프라인 아키텍처를 엄격하게 준수하며 고품질 코드를 유지해왔습니다. 그러나 지속적인 기능 추가 및 테스트 과정에서 산출물 방치, 디렉토리 간 역할 중복, 문서 누적으로 인한 관리 비용 증가가 관찰되었습니다.

본 문서는 프로젝트의 장기적 유지보수성 향상과 신규 기여자의 원활한 온보딩을 위해 데이터 분석에 기반한 **4단계 점진적 디렉토리 정리 계획**을 정의합니다.

### 1.1 데이터 기반 현황 및 문제점 분석

| 영역 | 현황 데이터 (As-Is) | 문제점 및 유지보수 위험 |
| :--- | :--- | :--- |
| **루트 디렉토리 (`.`)** | • `demo-standalone.html` (4.8MB, Untracked)<br>• `error.txt`, `.error.txt.swp` (Untracked) | 빌드 산출물과 디버깅 임시 파일이 프로젝트 최상위에 방치되어 저장소가 오염되고, 실수로 커밋될 위험이 있음 |
| **`testdata/` vs `examples/`** | • `testdata/demo.html` (Git Tracked)<br>• `testdata/demo.pdf`, `demo.pptx`, `demo-test.html`<br>• `testdata/demo.md`, `demo.ko.md`, 이미지 3종<br>• `examples/`에는 DSL 가이드 문서 2종만 존재 | • `testdata/`는 Go 관례상 테스트 전용 픽스처 공간이나, 사용자 데모 파일과 빌드 생성물이 혼재됨<br>• 사용자가 저장소를 클론했을 때 공식 예제 위치 탐색에 혼선 발생 |
| **`task/` 디렉토리** | • `GOS-2.md` ~ `GOS-24.md` (19개 티켓)<br>• `*-plan.md` 3개 등 총 22개 파일이 루트에 평면(flat) 나열 | 티켓이 지속적으로 추가됨에 따라 파일 탐색 비용이 급증하며, 이미 완료된 작업과 현재 활성 태스크 간의 구분이 어려움 |
| **`docs/` 디렉토리** | • `docs/architecture/` (01~09 모듈화 우수)<br>• 최상위 `docs/`에 7개 기획/보고서 분산<br>• `well_architected_assessment.md`와 `docs/release/well-architected-report.md` 중복성 | 모듈화된 아키텍처 문서와 초기 기획서/평가 보고서가 한곳에 섞여 있어 문서 관리 위계가 불명확함 |
| **`internal/` 패키지** | • `internal/browser` (Chrome 경로 탐색)<br>• `internal/i18n` (다국어 지원) 구현체 존재 | 실제 구현 코드는 우수하게 분리되어 있으나, [`docs/architecture/03-package-structure-and-c4.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/03-package-structure-and-c4.md)에 2개 패키지가 누락되어 문서 불일치 발생 |

---

## 2. 디렉토리 구조 Before & After 비교

```text
[Before: 현재 구조]
Goslide/
├── demo-standalone.html          # [오염] 4.8MB 빌드 산출물 방치
├── error.txt, .error.txt.swp     # [오염] 디버깅 임시 파일
├── task/                         # [비대] 22개 티켓/계획서 파일 평면 나열
│   ├── GOS-2.md ~ GOS-24.md
│   └── *-plan.md
├── testdata/                     # [혼재] 테스트 픽스처 + 사용자 데모 + 빌드 산출물
│   ├── demo.html (git tracked!)
│   ├── demo.pdf, demo.pptx, demo-test.html (untracked)
│   ├── demo.md, demo.ko.md, sample-img-*.png, hands-on-bg.jpeg
│   ├── golden/                   # 회귀 테스트 레퍼런스
│   └── slides/                   # 회귀 테스트 입력
├── examples/                     # [부실] DSL 가이드 2개만 존재
│   ├── example-dsl.md
│   └── example-dsl.ko.md
└── docs/                         # [분산] 최상위 디렉토리에 7개 파일 분산
    ├── architecture/             # 01~09 상세 설계서
    ├── okf/                      # 시맨틱 지식 라우터
    ├── release/                  # 릴리즈 검증 보고서
    └── *.md (7개)

                    ▼ ▼ ▼ ▼ ▼

[After: 권장 목표 구조]
Goslide/
├── task/
│   ├── archive/                  # 완료된 티켓 명세 보관 (GOS-2 ~ GOS-23)
│   ├── GOS-24.md                 # 현재 활성 태스크
│   └── *-plan.md                 # 중장기 실행 계획서
├── testdata/                     # [순수 테스트 전용] Go 툴체인 및 golden_test 전용
│   ├── golden/                   # 레퍼런스 HTML 픽스처
│   └── slides/                   # 마크다운 테스트 픽스처
├── examples/                     # [사용자 데모 및 가이드 전용]
│   ├── demo/                     # README 공식 데모 (demo.md, demo.ko.md, 이미지 에셋)
│   └── dsl/                      # Goslide DSL 사양 및 골든 템플릿 가이드
└── docs/                         # [위계화된 기술 문서 스위트]
    ├── architecture/             # 01~09 아키텍처 상세 사양서 (browser, i18n 동기화)
    ├── okf/                      # 시맨틱 지식 라우터
    ├── planning/                 # 제품 기획서, 로드맵, 기능 사양서
    └── release/                  # 릴리즈 체크리스트, 품질 보고서
```

---

## 3. 단계별 상세 실행 로드맵 (Action Items)

기존 기능 및 테스트에 부작용(Side Effects)을 최소화하기 위해 **위험도가 낮고 즉각적인 개선 효과가 큰 단계부터 순차적으로 진행**합니다.

### Phase 1: 루트 워크스페이스 청결화 및 산출물 격리 (즉시 실행 권장)
* **목표**: 저장소 오염 제거 및 빌드 산출물의 실수 커밋 방지
* **세부 작업**:
  1. **루트 임시 파일 삭제**:
     - `rm -f demo-standalone.html error.txt .error.txt.swp`
  2. **오추적된 산출물 캐시 제거**:
     - `git rm --cached testdata/demo.html`
     - `rm -f testdata/demo-test.html testdata/demo_out.html testdata/demo.pdf testdata/demo.pptx`
  3. **[`.gitignore`](file:///home/yundream/myjob/cloit/Goslide/.gitignore) 보강**:
     ```gitignore
     # Build and export outputs
     *.html
     !testdata/golden/*.html
     *.pdf
     *.pptx
     demo-standalone.html
     *-standalone.html
     error.txt
     ```
  4. **[`Makefile`](file:///home/yundream/myjob/cloit/Goslide/Makefile) `clean` 타겟 고도화**:
     ```makefile
     clean:
     	@echo "==> Cleaning build artifacts..."
     	@rm -rf $(BIN_DIR) dist *.out *.test coverage.* *.pdf *.pptx *_out.html *-standalone.html error.txt
     ```

### Phase 2: 예제(`examples/`)와 테스트 데이터(`testdata/`)의 명확한 분리
* **목표**: 사용자 대면 예제와 순수 Go 자동화 테스트 픽스처의 격리
* **세부 작업**:
  1. **디렉토리 생성 및 에셋 이동**:
     - `mkdir -p examples/demo examples/dsl`
     - `git mv testdata/demo.md testdata/demo.ko.md testdata/sample-img-01.png testdata/sample-img-02.jpg testdata/hands-on-bg.jpeg examples/demo/`
     - `git mv examples/example-dsl.md examples/example-dsl.ko.md examples/dsl/`
  2. **관련 문서 및 템플릿 경로 갱신**:
     - [`README.md`](file:///home/yundream/myjob/cloit/Goslide/README.md) 및 [`README.ko.md`](file:///home/yundream/myjob/cloit/Goslide/README.ko.md):  
       `goslide build testdata/demo.md` ➔ `goslide build examples/demo/demo.md`
     - [`internal/theme/assets/templates/demo.md`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/templates/demo.md):  
       이미지 참조 경로를 `examples/demo/...`로 동기화.
  3. **`testdata/` 순수화 확인**:
     - `testdata/` 내부에는 오직 `testdata/slides/*.md`와 `testdata/golden/*.html`만 남아 `go test ./...` 실행 시 완벽한 테스트 격리 달성.

### Phase 3: `task/` 티켓 아카이빙 체계 구축
* **목표**: 활성 작업 공간 시인성 확보 및 과거 히스토리 체계적 보관
* **세부 작업**:
  1. **아카이브 디렉토리 생성**:
     - `mkdir -p task/archive`
  2. **완료된 티켓 이동**:
     - `git mv task/GOS-2.md task/GOS-3.md ... task/GOS-23.md task/archive/`
  3. **활성 문서 유지**:
     - 현재 진행 중인 티켓(`task/GOS-24.md`)과 전체 계획서(`*-plan.md`)는 `task/` 최상위에 유지하여 작업 맥락 집중.

### Phase 4: 기술 문서군 위계화 및 아키텍처 사양 동기화
* **목표**: Diátaxis 및 arc42 표준에 맞춘 기술 문서 구조화
* **세부 작업**:
  1. **기획/평가 문서 카테고리화**:
     - `mkdir -p docs/planning docs/reports`
     - `git mv docs/functional_specification.md docs/project_plan.md docs/development_roadmap.md docs/competitive-analysis-and-strategy.md docs/planning/`
     - `git mv docs/architecture_review.md docs/well_architected_assessment.md docs/reports/`
  2. **내비게이션 링크 동기화**:
     - [`docs/okf/index.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/index.md) 및 [`docs/architecture/index.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/index.md)의 상대 링크 경로 일괄 갱신.
  3. **아키텍처 명세서 최신화**:
     - [`docs/architecture/03-package-structure-and-c4.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/03-package-structure-and-c4.md)의 다이어그램 및 패키지 책임 표에 `internal/browser`와 `internal/i18n` 공식 추가.

---

## 4. 검증 체크리스트 및 롤백 전략 (Verification & Safety)

| 검증 항목 | 수행 명령어 | 기대 결과 |
| :--- | :--- | :--- |
| **단위 테스트 회귀** | `go test -v -cover ./...` | 모든 단위 테스트 PASS (100% 통과) |
| **골든 파일 회귀** | `go test -v ./internal/renderer/html` | 바이트 단위 일치 및 회귀 없음 |
| **데모 빌드 검증** | `go run ./cmd/goslide build examples/demo/demo.md -o bin/demo.html` | 오류 없이 정상 빌드 및 리소스 로드 |
| **Git 상태 청결도** | `git status` | Untracked 임시 파일 없이 깨끗한 상태 유지 |

> **안전 가드레일 (Safety)**:
> 모든 파일 이동은 `git mv`를 사용하여 git 히스토리를 100% 보존하며, 각 단계는 독립된 커밋으로 나누어 문제 발생 시 언제든지 특정 단계만 롤백할 수 있도록 안전하게 진행합니다.
