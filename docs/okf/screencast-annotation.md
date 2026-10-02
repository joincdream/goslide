---
type: concept
title: Screencast Annotation & Presentation Runtime
description: Specifications for lightweight transparent canvas drawing overlay optimized for YouTube educational screencasts
tags:
  - screencast
  - annotation
  - canvas
  - runtime
timestamp: 2026-10-02T22:50:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# 스크린캐스트 특화 판서(Annotation) 레이어 사양

Goslide는 일반적인 마크다운 슬라이드 빌더와 달리 **"스크린캐스트 기반 유튜브 교육 영상 제작"**을 핵심 유즈케이스로 설정하고 이에 최적화된 **초경량 투명 캔버스 판서 엔진**을 내장합니다.

---

## 1. 2계층 뷰포트 구조 (Two-Layer Viewport)

외부 복잡한 그래픽 라이브러리 없이 순수 HTML5/CSS 표준으로 구현됩니다:

```text
[브라우저 뷰포트]
  ┌────────────────────────────────────────────────────────┐
  │  Layer 1 (하단): 슬라이드 HTML 콘텐츠 (텍스트, 코드, 이미지) │
  │  Layer 2 (상단): 투명 <canvas> 레이어 (전체화면 오버레이)   │
  └────────────────────────────────────────────────────────┘
```

1. **평상시 (슬라이드 네비게이션 모드)**:
   - `<canvas>`의 스타일을 `pointer-events: none;`으로 설정.
   - 마우스 클릭, 텍스트 드래그, 링크 클릭, 방향키 슬라이드 넘김이 일반 웹 브라우저처럼 100% 정상 동작.
2. **판서 모드 활성화 (단축키 `D`)**:
   - `<canvas>`의 스타일을 `pointer-events: auto;`로 전환하고 마우스 커서를 펜 형태(`cursor: crosshair;`)로 변경.
   - 마우스 드래그 이벤트(`mousedown` $\rightarrow$ `mousemove` $\rightarrow$ `mouseup`) 단 3개로 Canvas 2D 컨텍스트에 즉각적인 선(`ctx.stroke()`) 드로잉 수행.
3. **판서 즉시 삭제 (단축키 `C`)**:
   - `ctx.clearRect(0, 0, width, height)` 호출로 현재 슬라이드의 판서를 즉각 초기화.
4. **슬라이드 전환 연동**:
   - 슬라이드가 이전/다음으로 이동(`Space`, `ArrowRight` 등)할 때 캔버스를 자동으로 깨끗이 리셋하여 다음 설명 녹화에 방해되지 않도록 보장.

---

## 2. 유튜브 스크린캐스트 최적화 키 바인딩

| 단축키 | 동작 | 역할 및 설명 |
| :---: | :--- | :--- |
| **`D`** | 펜 토글 (Draw Mode) | 일반 네비게이션 모드 $\leftrightarrow$ 투명 캔버스 판서 모드 전환 |
| **`C`** | 전체 지우기 (Clear) | 현재 슬라이드에 그려진 모든 판서 즉시 삭제 |
| **`Space` / `→` / `J`** | 다음 슬라이드 | 다음 슬라이드로 이동 (판서 자동 클리어) |
| **`Backspace` / `←` / `K`** | 이전 슬라이드 | 이전 슬라이드로 이동 (판서 자동 클리어) |
| **`F`** | 전체화면 (Fullscreen) | 브라우저 전체화면 토글 (OBS/화면 캡처 최적화) |

---

## 3. 부가 장식 도구 배제 원칙 (YAGNI)

스크린캐스트 교육 영상 녹화 환경에서는 아래 도구들을 일체 구현하지 않습니다:
- **레이저 포인터**: 마우스 커서 자체가 녹화 화면에 선명히 노출되므로 불필요.
- **스포트라이트 마스크**: 영상 시청자의 시야 흐름을 방해하고 편집 시 부자연스러움.
- **화면 블라인드 (B/W)**: OBS Studio 등 화면 녹화 소프트웨어의 씬 전환 기능으로 완벽 대체 가능.
- **듀얼스크린 콘솔**: 단일 모니터 화면 녹화 시 별도 창 통신 오버헤드만 발생.

---

## 4. 구현 규모 계약

- **런타임 파일**: `internal/theme/assets/js/goslide-core.js`
- **코드 크기**: 외부 프레임워크(Fabric.js 등) 없이 **순수 바닐라 JS 단 50줄 내외**로 완성.
- **성능**: 60fps 무지연 드로잉 및 제로 메모리 누수 보장.
