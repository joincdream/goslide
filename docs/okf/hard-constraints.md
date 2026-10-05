---
type: reference
title: Goslide Hard Engineering Constraints
description: Mandatory software engineering rules, circular dependency guards, and verification pointers
tags:
  - constraints
  - guidelines
  - verification
  - reference
timestamp: 2026-10-05T11:40:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Goslide 절대적 엔지니어링 제약 (Hard Constraints)

본 문서는 Goslide의 모든 코드 작성 및 리팩토링 시 예외 없이 준수해야 하는 **하드 가드레일(Hard Constraints)** 요약 및 단일 진실 공급원(SSOT) 포인터입니다.

---

## 1. 5대 하드 가드레일 (절대 위반 금지)

1. **CGO 의존성 0% (Pure Go)**: `CGO_ENABLED=0` 환경에서 단일 정적 바이너리로 컴파일되어야 합니다.
2. **글로벌 가변 상태 전면 금지**: 패키지 레벨 전역 변수 금지, 모든 컴포넌트는 `New...` 생성자 주입(DI) 필수.
3. **`panic()` 남용 금지**: 라이브러리 및 코어 로직 내 `panic` 금지, Go 1.13+ `%w` 에러 래핑 반환 필수.
4. **Context 생명주기 및 리소스 회수 보장**: 모든 I/O 및 브라우저 제어에 `context.Context` 전달, `defer close()` 보장.
5. **다형성 분기 체인 배제 (OCP)**: 레이아웃/지시어 등 3단 이상의 `if-else`/`switch` 체인 금지, Registry/Strategy 패턴 필수.

---

## 2. 단일 진실 공급원 (SSOT Navigation)

상세한 엔지니어링 가이드라인, 코딩 관용구, 테스트 규격 및 정적 분석 규칙은 루트 문서에서 관리됩니다:

* **엔지니어링 절대 원칙 & LLM 가드레일 (SSOT)**: [**`AGENTS.md`**](file:///home/yundream/myjob/cloit/Goslide/AGENTS.md)
* **빌드 및 린트 검증 규격**: [`Makefile`](file:///home/yundream/myjob/cloit/Goslide/Makefile) (`make test-race`, `make lint`, `make complexity`)
* **아키텍처 감사 루브릭**: [`.agents/skills/well-architected-review/references/rubrics.md`](file:///home/yundream/myjob/cloit/Goslide/.agents/skills/well-architected-review/references/rubrics.md)
