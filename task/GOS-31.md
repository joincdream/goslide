# [GOS-31] v1.0.0 릴리즈 패키징 아키텍처 수립 및 첫 사용자 경험(FTUX) 온보딩 번들 개발 계획서

> **티켓 번호**: [GOS-31](https://joincdream.atlassian.net/browse/GOS-31)  
> **마일스톤**: v1.0.0 정식 릴리즈 배포 및 온보딩 (Release Packaging & FTUX)  
> **마감일**: 2026-11-20  
> **상태**: 할 일 (To Do)  
> **담당 패키지**: [`cmd/goslide/`](../cmd/goslide/), [`.goreleaser.yaml`](../.goreleaser.yaml), [`themes/`](../themes/), [`examples/`](../examples/)  
> **연관 문서**:  
> - [`docs/README.md`](../docs/README.md) *(Directory Role Taxonomy & SSOT)*  
> - [`docs/dsl-guide.md`](../docs/dsl-guide.md) *(Slide Authoring & DSL Specification)*  
> - [`docs/theme-guide.md`](../docs/theme-guide.md) *(Theme Development Specification)*  
> - [`docs/okf/decisions-simplification.md`](../docs/okf/decisions-simplification.md) *(ADR-001: 단순성 및 단일 바이너리 원칙)*  
> - [`task/release/pre-release-qa-plan.md`](release/pre-release-qa-plan.md)  

---

## 1. 개요 및 배경 (Context & Problem Statement)

Goslide는 CGO 0%의 단일 정적 바이너리(Single Binary)와 초고속 렌더링 성능을 자랑합니다. 그러나 바이너리 파일 하나만 배포할 경우, 사용자가 최초로 다운로드받아 실행했을 때 다음과 같은 **첫 사용자 경험(FTUX: First-Time User Experience) 상의 진입 장벽**이 발생할 수 있습니다:

1. **외장 테마 탐색 및 적용의 번거로움**:
   - `themes/` 디렉토리에 구축된 공식 외장 테마(Corporate, Academic, Cyber-Dark)를 사용자가 직접 깃허브에서 복사/다운로드하지 않고도 CLI 또는 번들을 통해 즉시 프로젝트로 가져올 수 있어야 함.
2. **실전 예제(Side Examples)의 부재**:
   - 최소한의 템플릿 외에 실제 완성도 높은 데모 슬라이드(도표, Mermaid 다이어그램, 2단 분할 레이아웃 등)를 로컬에서 즉시 확인하고 발표해 볼 수 있는 스타터 덱 제공 필요.
3. **LLM 친화적 Slide DSL 가이드 제공**:
   - 사용자가 ChatGPT, Claude, Cursor 등에 Goslide 슬라이드 작성을 요청할 수 있도록, DSL 규격 요약 및 1줄 복사 프롬프트를 온/오프라인 환경에서 손쉽게 획득할 수 있어야 함.

본 과제는 Goslide를 처음 접하는 사용자가 바이너리를 다운로드하고 **3분 이내에 매력적인 첫 프레젠테이션을 완성하고 발표할 수 있도록 지원하는 패키징 및 온보딩 아키텍처**를 수립하고 개발합니다.

---

## 2. 핵심 검토 및 개발 과제 (Key Scope)

### 2.1 외장 테마(themes/) 패키징 및 스캐폴딩 방안
- **아카이브 번들링 vs CLI 온디맨드 추출**:
  - `tar.gz` / `zip` 배포 아카이브에 `themes/` 디렉토리를 포함할지 여부 결정.
  - 또는 CLI 명령어(예: `goslide theme init` 또는 `goslide init --extract-themes`)를 통해 내장 에셋이나 공식 템플릿을 현재 디렉토리로 스캐폴딩하는 방안 검토.
- **단일 바이너리 가치 훼손 방지**:
  - 외장 테마를 기본 탑재하되 바이너리 단독 실행 시에도 정상 폴백(fallback)되도록 무결성 유지.

### 2.2 스타터 예제(Side Examples) 및 `goslide init` 템플릿 고도화
- `goslide init [filename]` 실행 시 기본 템플릿 외에 목적별 템플릿 선택 지원:
  - `--template corporate` (비즈니스/제안서 스타터 덱)
  - `--template academic` (학술/강의 스타터 덱)
  - `--template cyber-dark` (테크 발표 스타터 덱)
- 릴리즈 패키지 내 `examples/` 디렉토리(골든 샘플 덱, 데모 덱) 포함 정책 정의.

### 2.3 LLM 연동용 Slide DSL 가이드라인 배포
- 사용자가 AI 도구에 프롬프트를 복사하여 붙여넣을 수 있는 최적의 전달 경로 수립:
  - CLI 명령어: `goslide dsl` 또는 `goslide prompt` 실행 시 LLM 프롬프트 가이드 표준 마크다운 출력.
  - 릴리즈 번들 내 `LLM_PROMPT.md` 또는 `DSL_GUIDE.md` 배치 여부 검토.
  - 공식 온라인 SSOT([`docs/dsl-guide.md`](../docs/dsl-guide.md))와의 실시간 동기화 방안.

### 2.4 GoReleaser 배포 패키지 아키텍처 수립
- [`.goreleaser.yaml`](../.goreleaser.yaml) 아카이브 설정 업데이트:
  - 바이너리 외에 라이선스(`LICENSE`), 안내문(`README.md`), 외장 테마(`themes/`), 예제(`examples/`) 파일의 아카이브 패키징 범위 확정.
  - OS별(Linux, macOS, Windows) 패키지 압축 포맷 및 디렉토리 구조 검증.

---

## 3. 단계별 실행 계획 (Implementation Steps)

| 단계 | 작업 내용 | 타겟 소스 및 산출물 |
| :---: | :--- | :--- |
| **Step 1** | **패키징 및 FTUX 배포 아키텍처 의사결정서 작성**<br>- 번들 포함 범위 확정 (themes, examples, LLM prompt)<br>- CLI 스캐폴딩 명령어 필요성 검토 | `docs/okf/decisions-packaging.md` (ADR) |
| **Step 2** | **CLI 온보딩 및 템플릿 생성 기능 고도화**<br>- `goslide init` 테마별 템플릿 확장<br>- LLM 프롬프트 지원 CLI 명령어 구현 (`goslide prompt` 등) | [`cmd/goslide/init.go`](../cmd/goslide/init.go)<br>[`cmd/goslide/`](../cmd/goslide/) |
| **Step 3** | **GoReleaser 아카이브 번들링 구성 및 테스트**<br>- `.goreleaser.yaml` 내 `files` 매핑 최적화<br>- 로컬 아카이브 생성 테스트 (`goreleaser release --snapshot --skip=publish --clean`) | [`.goreleaser.yaml`](../.goreleaser.yaml) |
| **Step 4** | **첫 사용자 온보딩 경험(FTUX) E2E 검증**<br>- 바이너리 압축 해제 후 3분 내 슬라이드 빌드 및 서빙 시나리오 검증<br>- Windows / Linux / macOS 아카이브 무결성 확인 | [`task/release/linux-user-test-checklist.md`](release/linux-user-test-checklist.md) |

---

## 4. 수용 기준 (Acceptance Criteria)

- [ ] **패키징 정책 확정**: 바이너리 아카이브(zip/tar.gz)에 포함될 디렉토리 및 파일 목록이 문서화되어야 함.
- [ ] **외장 테마 온보딩**: 사용자가 추가 작업 없이 공식 외장 테마(corporate, academic, cyber-dark)를 손쉽게 적용할 수 있어야 함.
- [ ] **LLM 프롬프트 연동**: 사용자가 Goslide용 슬라이드 작성을 LLM에 즉시 요청할 수 있는 프롬프트 에셋이 제공되어야 함.
- [ ] **GoReleaser 빌드 검증**: `goreleaser --snapshot` 실행 시 정해진 패키지 구조대로 아카이브가 정상 생성되어야 함.
- [ ] **단일 바이너리 원칙 준수**: 외부 에셋 번들 없이 바이너리 단독으로 실행하더라도 기존 기능이 100% 정상 작동해야 함.
