---
type: reference
title: Goslide Domain Contracts & Interfaces
description: Thin contract reference and entry point pointers for Goslide domain models and interfaces
tags:
  - interfaces
  - models
  - contracts
  - reference
timestamp: 2026-10-05T11:40:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Goslide 도메인 계약 및 인터페이스 (Contracts & Interfaces)

`internal/model` 패키지는 최하위 도메인 계층으로서 Go 표준 라이브러리 외에 외부 의존성을 일체 갖지 않는(Zero Dependency) 순수 데이터 모델 및 인터페이스 계약입니다.

---

## 1. 핵심 계약 규칙 (Core Invariants)

1. **단일 메서드 인터페이스 (Single-Method ISP)**:
   * 모든 코어 인터페이스(`Parser`, `Renderer`, `Exporter`)는 결합도를 최소화하기 위해 1개의 메서드만 정의합니다.
2. **불변 IR 모델 (Immutability)**:
   * 파싱 완료된 `Deck`, `Slide`는 읽기 전용으로 소비되며, 렌더링/익스포트 계층에서 직접 수정할 수 없습니다.
3. **AST 오염 방지 (Zero AST Dependency)**:
   * 도메인 모델은 마크다운 파서 AST(`goldmark/ast`)에 오염되지 않은 순수 Go 구조체로 유지됩니다.
4. **마스터 레이아웃 레지스트리 (Layout OCP)**:
   * 슬라이드 레이아웃은 `MasterLayoutRegistry` 및 `LayoutSpec`을 통해 다형성 전략 패턴으로 디스패치됩니다.

---

## 2. 코드 및 상세 문서 포인터 (Navigation)

| 대상 | 위치 | 역할 |
| :--- | :--- | :--- |
| **코어 인터페이스** | [`internal/model/interfaces.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/interfaces.go) | `Parser`, `Renderer`, `Exporter` 인터페이스 정의 |
| **불변 슬라이드 IR** | [`internal/model/deck.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/deck.go) | `Deck`, `Slide`, `SlideDirectives` 데이터 구조체 |
| **마스터 레이아웃 규격** | [`internal/model/layout.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/layout.go) | `MasterLayoutRegistry`, `LayoutSpec`, 카테고리 5종 |
| **도메인 센티넬 에러** | [`internal/model/errors.go`](file:///home/yundream/myjob/cloit/Goslide/internal/model/errors.go) | `ErrThemeNotFound`, `ErrInvalidMarkdown` 등 |
| **상세 아키텍처 사양** | [`docs/architecture/04-interfaces-and-contracts.md`](file:///home/yundream/myjob/cloit/Goslide/docs/architecture/04-interfaces-and-contracts.md) | 입출력 명세 및 포맷별 매핑 상세 가이드 |
