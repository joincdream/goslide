---
type: index
title: Goslide Knowledge Bundle Index
description: Entry point and knowledge routing map for Goslide presentation deck engine
tags:
  - goslide
  - architecture
  - okf
  - index
timestamp: 2026-10-02T22:50:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Goslide Open Knowledge Bundle (OKF)

Goslide의 아키텍처, 데이터 모델, 인터페이스 계약, 제약 조건을 AI 에이전트와 엔지니어가 신속하게 탐색할 수 있도록 표준화한 **Open Knowledge Format (OKF v0.2)** 번들입니다.

---

## 1. 지시 및 지식 맵 (Knowledge Map)

| 문서 | Type | 역할 및 내용 |
| :--- | :---: | :--- |
| **[core-architecture.md](./core-architecture.md)** | `concept` | 단방향 파이프라인, 포맷별 렌더링/익스포트 흐름 및 실시간 SSE 프리뷰 아키텍처 |
| **[contracts-interfaces.md](./contracts-interfaces.md)** | `reference` | `internal/model`의 코어 인터페이스(`Parser`, `Renderer`, `Exporter`) 및 Go 불변 데이터 모델 |
| **[hard-constraints.md](./hard-constraints.md)** | `reference` | CGO 배제, panic 금지, 전역상태 금지, CLI 종료 코드 및 DoD 검증 체크리스트 |
| **[decisions-simplification.md](./decisions-simplification.md)** | `decision` | YAGNI/KISS/SoC 원칙에 따른 7대 오버엔지니어링 제거 결정 기록 (ADR) |
| **[screencast-annotation.md](./screencast-annotation.md)** | `concept` | 유튜브 교육 영상 스크린캐스트 녹화 특화 초경량 투명 캔버스 판서(Annotation) 레이어 사양 |
| **[layout-design.md](../layout-design.md)** | `design` | FHD 16:9 기준 In-Flow Flexbox 3단(헤더/본문/푸터) 완전 격리 레이아웃 설계 사양 |

---

## 2. 시스템 개요 (Executive Summary)

- **정의**: Go 기반 초경량, 고성능 마크다운 프레젠테이션 데크 빌더 (HTML, PDF, PPTX).
- **특화 가치**: 스크린캐스트 유튜브 교육 영상 제작에 최적화된 **초경량 투명 캔버스 판서 레이어** 내장.
- **실행 형태**: 외부 런타임 의존성이 완전히 배제된 **순수 Go 단일 실행 바이너리 (Single Binary)**.

---

## 3. 린(Lean) 패키지 레이아웃 (14개 핵심 파일)

```text
Goslide/
├── cmd/goslide/              # [main.go, build.go] CLI 진입점
├── pkg/goslide/              # [goslide.go] 외부 연동용 Facade Public API
├── internal/
│   ├── model/                # [deck.go, interfaces.go] 불변 도메인 모델 (Zero Dep)
│   ├── parser/               # [parser.go, highlight.go, parser_test.go] 파서 & 하이라이터
│   ├── theme/                # [embed.go, theme.go] embed.FS 테마 로더 (CSS/JS 내장)
│   ├── renderer/             # [html.go, html_test.go] HTML 슬라이드 렌더러
│   ├── exporter/             # [pdf.go, pptx.go, exporter_test.go] chromedp 기반 PDF/PPTX 익스포터
│   ├── server/               # [server.go, watcher.go] net/http 정적 서빙 + SSE 리로드
│   └── testutil/             # [golden.go] 골든 파일 회귀 검증 헬퍼
├── testdata/                 # 테스트 마크다운 픽스처 및 골든 파일
├── go.mod                    # Go 1.22+ 의존성 (Zero WebSocket, Zero CGO)
├── Makefile                  # build, test, lint 자동화
└── AGENTS.md                 # 엔지니어링 가이드라인
```
