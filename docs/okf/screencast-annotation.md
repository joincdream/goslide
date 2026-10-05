---
type: concept
title: Screencast Annotation & Presentation Runtime
description: Specifications and component pointers for transparent canvas drawing overlay and presenter runtime
tags:
  - screencast
  - annotation
  - canvas
  - runtime
timestamp: 2026-10-05T11:40:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# 스크린캐스트 특화 판서 및 런타임 사양 (Screencast & Runtime)

Goslide는 **"스크린캐스트 기반 교육 영상 제작"**에 최적화된 **1:1 뷰포트 투명 캔버스 판서 엔진**과 **발표자 콘솔(Presenter View)**을 내장합니다.

---

## 1. 2계층 뷰포트 및 제어 원칙

```text
[뷰포트 1920x1080]
  Layer 1 (하단): 슬라이드 HTML 콘텐츠 (CSS Grid/Flexbox 마스터 구조)
  Layer 2 (상단): 투명 <canvas id="goslide-canvas"> (1:1 풀스크린 드로잉 오버레이)
```

* **포인터 이벤트 분리**: 평상시 `pointer-events: none` 유지, `D` 키 입력 시 `pointer-events: auto`로 전환하여 펜 드로잉 활성화.
* **자동 클리어**: 슬라이드 전환(`Space`, `ArrowRight` 등) 시 캔버스 버퍼를 즉각 초기화하여 다음 슬라이드 녹화 방해 차단.
* **발표자 뷰**: `P` 키로 분리 창 모드, `N` 키로 인-윈도우 480px 사이드바 모드 전환.

---

## 2. 구현 소스 포인터 (Navigation)

| 컴포넌트 | 소스 위치 | 역할 |
| :--- | :--- | :--- |
| **캔버스 드로잉 UI** | [`web/src/components/DrawingCanvas.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/DrawingCanvas.svelte) | Svelte 5 기반 펜/지우개 드로잉 엔진 |
| **발표자 사이드바 UI** | [`web/src/components/PresenterSidebar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterSidebar.svelte) | 480px 사이드바, 타이머, 16:9 프리뷰, 노트 |
| **캔버스 & 발표자 CSS** | [`internal/theme/assets/css/deck-canvas.css`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/css/deck-canvas.css) | 캔버스 오버레이, 인디케이터, 4단 마스터 구조 |
| **내장 번들 에셋** | [`internal/theme/assets/js/goslide-core.js`](file:///home/yundream/myjob/cloit/Goslide/internal/theme/assets/js/goslide-core.js) | Svelte 5 컴파일된 IIFE 프로덕션 번들 |
