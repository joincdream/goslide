---
type: index
title: Goslide Open Knowledge Router (OKF)
description: Intent-based knowledge routing map and fast navigation index for Goslide
tags:
  - goslide
  - okf
  - index
  - router
  - navigation
timestamp: 2026-10-05T11:40:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Goslide Open Knowledge Router (OKF)

본 문서는 AI 에이전트와 엔지니어가 불필요한 컨텍스트 토큰 소모 없이, **작업 의도(User Intent)에 부합하는 최소한의 문서와 소스 코드만 $O(1)$로 탐색(Selective Loading)**하도록 안내하는 **시맨틱 라우팅 맵**입니다.

> [!IMPORTANT] 에이전트 탐색 가드레일
> 모든 문서를 한꺼번에 읽지 마십시오. 아래 라우팅 테이블에서 현재 사용자의 지시사항(Task Domain)과 일치하는 행의 **최소 지식 문서 1개와 소스 코드**만 열람하십시오.

---

## 1. 의도 기반 시맨틱 라우팅 테이블 (Semantic Intent Router)

| 작업 의도 및 키워드 (User Intent / Task Domain) | 로드할 최소 지식 문서 (On-Demand Knowledge) | 핵심 소스 코드 진입점 (Source Pointers) |
| :--- | :--- | :--- |
| **슬라이드 레이아웃 / 본문 구조 / 2단 컬럼 / 타이틀 박스**<br>`layout`, `two-cols`, `lead`, `compact`, `slide-title-box` | [`docs/planning/layout-design.md`](file:///home/yundream/myjob/cloit/Goslide/docs/planning/layout-design.md)<br>[`docs/architecture/05-rendering-and-export-pipelines.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/05-rendering-and-export-pipelines.md) | [`internal/model/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/layout.go)<br>[`internal/parser/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/layout.go) |
| **마크다운 파싱 / Frontmatter / 지시어 / Chroma 하이라이트**<br>`parse`, `frontmatter`, `directive`, `highlight`, `delimiter` | [`docs/architecture/02-slide-spec-and-domain-model.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/02-slide-spec-and-domain-model.md) | [`internal/parser/parser.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/parser.go)<br>[`internal/parser/directive.go`](file:///home/yundream/myjob/cloit/Goslide/internal/parser/directive.go) |
| **테마 / 모듈러 CSS / 폰트 / 이미지 번들러**<br>`theme`, `css`, `deck-canvas.css`, `deck-content.css`, `bundler` | [`docs/architecture/05-rendering-and-export-pipelines.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/05-rendering-and-export-pipelines.md) | [`internal/theme/theme.go`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/theme.go)<br>[`internal/theme/assets/css/`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/css/) |
| **발표자 콘솔 / 캔버스 판서 / 단축키 / Svelte 프론트엔드**<br>`presenter`, `drawing`, `canvas`, `annotation`, `svelte`, `shortcut` | [`docs/okf/screencast-annotation.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/screencast-annotation.md) | [`web/src/App.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/App.svelte)<br>[`web/src/components/`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/) |
| **라이브 서버 / SSE 핫리로드 / 파일 감시기**<br>`serve`, `sse`, `watcher`, `reload`, `live-reload` | [`docs/architecture/06-runtime-security-and-lifecycle.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/06-runtime-security-and-lifecycle.md) | [`internal/server/server.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/server.go)<br>[`internal/server/watcher.go`](file:///home/yundream/myjob/cloit/Goslide/internal/server/watcher.go) |
| **도메인 모델 / 인터페이스 계약 / 센티넬 에러**<br>`interface`, `contract`, `Deck`, `Slide`, `errors` | [`docs/okf/contracts-interfaces.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md) | [`internal/model/interfaces.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/interfaces.go)<br>[`internal/model/deck.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/deck.go) |
| **아키텍처 단순화 결정 (ADR-001) / YAGNI 원칙**<br>`adr`, `simplification`, `yagni`, `kiss`, `anti-pattern` | [`docs/okf/decisions-simplification.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md) | [ADR-001](file:///home/yundream/myjob/cloit/Goslide/docs/okf/decisions-simplification.md) |
| **웹 에디터 / Vim 모드 / 스튜디오 레이아웃 / 자동 저장**<br>`editor`, `vim`, `studio`, `codemirror`, `auto-save` | [`docs/architecture/08-embedded-editor-and-vim-mode.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/08-embedded-editor-and-vim-mode.md) | [`internal/server/`](file:///home/yundream/myjob/cloit/Goslide/internal/server/)<br>[`web/src/`](file:///home/yundream/myjob/cloit/Goslide/web/src/) |
| **테마 관리 / 커스텀 CSS 공유 / 구글 슬라이드 스타일 스위처**<br>`theme-manager`, `custom-css`, `sharing`, `theme-drawer` | [`docs/architecture/09-theme-management-and-sharing.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/09-theme-management-and-sharing.md) | [`internal/theme/`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/)<br>[`web/src/`](file:///home/yundream/myjob/cloit/Goslide/web/src/) |
| **코딩 제약 / 정적 분석 / CGO 0% / panic 금지**<br>`constraints`, `rules`, `lint`, `test-race`, `complexity` | [**`AGENTS.md`**](file:///home/yundream/myjob/cloit/Goslide/AGENTS.md)<br>[`docs/okf/hard-constraints.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md) | [`Makefile`](file:///home/yundream/myjob/cloit/Goslide/Makefile) |

---

## 2. 전체 상세 기술 사양서 스위트 (Deep-Dive Specifications)

시스템 전반의 심층 설계 사양이 필요할 때는 모듈화된 아키텍처 문서를 참조합니다:

* **시스템 아키텍처 인덱스**: [**`docs/architecture/index.md`**](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/index.md) (01~09 모듈 전체 맵)
* **품질 감사 보고서**: [`docs/reports/well_architected_assessment.md`](file:///home/yundream/myjob/cloit/Goslide/docs/reports/well_architected_assessment.md)
* **경쟁 도구 분석 및 제품 전략 보고서**: [`docs/planning/competitive-analysis-and-strategy.md`](file:///home/yundream/myjob/cloit/Goslide/docs/planning/competitive-analysis-and-strategy.md)
