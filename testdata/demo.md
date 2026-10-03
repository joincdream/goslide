---
title: "Goslide 아키텍처 및 기능 프리뷰"
author: "Goslide Team"
theme: "clean"
size: "16:9"
paginate: true
header: "Goslide 기술 세미나"
footer: "© 2026 Cloit Corp."
---

<!-- _layout: cover -->
<!-- _backgroundColor: #0f172a -->
<!-- _color: #f8fafc -->

# Goslide 프레젠테이션 엔진
### 고성능 순수 Go 마크다운 슬라이드 빌더

<!-- note:
발표 시작 시 단일 바이너리와 빠른 빌드 속도를 강조할 것.
-->

---

<!-- _layout: two-cols -->
## 아키텍처 비교

### 기존 슬라이드 도구
- Node.js 무거운 런타임 의존성
- 대용량 파일 빌드 시 메모리 부담
- 복잡한 2단 레이아웃 CSS 직접 작성

<!-- split -->

### Goslide (Pure Go)
- CGO 0% 단일 바이너리
- 밀리초 단위 즉각 변환
- `<!-- split -->` 네이티브 2단 분할 지원

---

<!-- _backgroundColor: #1e1e2e -->
<!-- _color: #cdd6f4 -->

## 순수 Go 구문 강조 (Chroma)

```go
package main

import "fmt"

// Slide represents a single presentation slide
type Slide struct {
    Index       int    `json:"index"`
    HTMLContent string `json:"html_content"`
}

func main() {
    fmt.Println("Hello, Goslide with Pure Go Chroma!")
}
```

---

## GFM 기능 지원 현황

| 기능 | 지원 마일스톤 | 상태 |
| :--- | :---: | :---: |
| YAML Frontmatter | M1-2 (GOS-3) | **완료** |
| 주석 지시어 & 스코프 | M1-3 (GOS-4) | **완료** |
| 2단 컬럼 분할 (`<!-- split -->`) | M1-3 (GOS-4) | **완료** |
| Chroma 코드 구문 강조 | M1-3 (GOS-4) | **완료** |
| 조기 프리뷰 빌드 (`goslide build`) | M1-3 (GOS-4) | **완료** |

- [x] 순수 Go 컴파일 보장 (`CGO_ENABLED=0`)
- [x] 동시성 데이터 레이스 0건
- [ ] 브라우저 슬라이드 쇼 런타임 (M1-5 예정)
