---
type: decision
title: ADR-001 Architecture Simplification & Anti-pattern Removal
description: Architectural Decision Record for eliminating over-engineering patterns based on YAGNI, KISS, and SoC
tags:
  - adr
  - decisions
  - simplification
  - yagni
  - kiss
timestamp: 2026-10-02T22:50:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# ADR-001: 아키텍처 단순화 및 오버엔지니어링 제거 결정

## 1. 배경 및 문제 제기 (Context)
초기 아키텍처(v1.2.0)는 가상의 요구사항과 불필요한 기능들을 수용하려다 다음과 같은 복잡도와 안티패턴을 초래했습니다:
1. Base64 이미지 번들링으로 인한 파일 크기 33% 팽창 및 CDP 브라우저 크래시.
2. 로컬 CLI 도구에 부적합한 Path Traversal 보안 모듈 탑재로 정상 파일 참조 차단.
3. 파서의 추측성 레이아웃 자동 추론(`inference.go`)으로 인한 결정론적 출력 훼손.
4. 브라우저 CSS Cascading 기능을 Go 백엔드에서 중복 구현(`cascader.go`).
5. 이미지 기반 PPTX로 확정되었음에도 미사용 AST 모델(`element.go`) 및 6개 파일 방치.
6. 단순 핫 리로드를 위해 외부 라이브러리(`gorilla/websocket`) 및 연결 관리 허브 도입.
7. 코드 구현 전 40여 개 파일로 잘게 쪼갠 조기 과분할(Over-modularization).

---

## 2. 결정 사항 (Decision Outcome)

YAGNI, KISS, SoC 원칙에 따라 다음과 같이 과감한 단순화 결정을 내리고 설계를 확정했습니다:

| 항목 | 기존 (AS-IS) | 결정 (TO-BE) | 근거 및 원칙 |
| :--- | :--- | :--- | :--- |
| **이미지 경로** | Base64 인라인 번들링 (`bundler.go`) | 마크다운 상대 경로 원본 보존 | **KISS**: 웹 표준 준수, 공유 시 zip/PDF/PPTX 활용 |
| **보안 모듈** | Path Traversal 검증기 (`path.go`) | 완전 삭제 | **YAGNI**: 로컬 CLI에서 상위 공통 에셋 참조 보장 |
| **레이아웃 판정** | 스마트 자동 추론 (`inference.go`) | 완전 삭제, 명시적 선언 100% 매핑 | **결정론적 출력**: 파서의 자의적 추측 원천 차단 |
| **테마 우선순위** | Go 백엔드 6단계 연산 (`cascader.go`) | 완전 삭제, 브라우저 CSS 엔진 위임 | **SoC**: 브라우저 고유 역할 침범 방지 |
| **PPTX 아키텍처** | `element.go` + 6개 파일 분할 | `element.go` 삭제, `pptx.go` 1개로 통폐합 | **YAGNI**: 캡처 기반 방식에 맞게 잔재 완전 청산 |
| **개발 서버 리로드**| `gorilla/websocket` + `hub.go` | Go 표준 `net/http` SSE (15줄) | **KISS/YAGNI**: 외부 의존성 제거, 단방향 스트리밍 |
| **패키지 분할** | 40여 개 파일 조기 파편화 | 15개 이내 핵심 파일로 통폐합 | **KISS**: Go 관용적 패키지 응집도 극대화 |

---

## 3. 결과 및 긍정적 효과 (Consequences)

- **인지 부하 및 복잡도 대폭 감소**: 파일 수가 40개에서 14개로 줄어 핵심 비즈니스 로직(파서, 렌더러, 익스포터)에만 온전히 집중 가능.
- **버그 및 크래시 원천 차단**: Base64 메모리 누수, Data URI 길이 초과, 웹소켓 연결 불안정, 레이아웃 오판정이 아키텍처 수준에서 박멸됨.
- **순수 단일 바이너리 강화**: `gorilla/websocket` 등 외부 의존성을 제거하여 CGO 0% 단일 실행 파일의 완성도 극대화.
