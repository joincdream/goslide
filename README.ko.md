<div align="center">

# Goslide 📽️

**순수 Go(Pure Go)로 구현된 초고속 단일 바이너리 마크다운 슬라이드 데크 빌더**

외부 런타임 의존성 없이 마크다운 문서를 인터랙티브 HTML, 벡터 PDF, 파워포인트(PPTX) 파일로 즉시 변환합니다.  
개발자 기술 발표, IT 교육 강의 및 유튜브 스크린캐스트 동영상 강의 녹화에 최적화되어 있습니다.

[English](README.md) | [한국어](README.ko.md)

</div>

---

## ✨ 핵심 기능

- **🚀 단일 바이너리 & 순수 Go (Pure Go)**: CGO 의존성 0%. `embed.FS`를 통해 테마와 런타임이 바이너리 안에 내장되어 단일 실행 파일 하나로 크로스 플랫폼(`Linux`, `macOS`, `Windows`) 배포 가능.
- **🎨 다중 포맷 일괄 변환 (Multi-Format Export)**:
  - **인터랙티브 HTML**: 외부 네트워크 없이 오프라인 구동되는 반응형 웹 슬라이드.
  - **고해상도 벡터 PDF**: Headless Chrome 기반 16:9 무마진 벡터 PDF 완벽 출력.
  - **고화질 PPTX**: HTML 시각 효과를 100% 보존하고 발표자 노트를 유지하는 파워포인트 생성.
- **🖍️ 유튜브 스크린캐스트 판서 오버레이**: 강의 녹화 중 슬라이드 위 자유 필기를 지원하는 경량 투명 캔버스 내장 (`D` 키로 펜 켜기/끄기, `C` 키로 전체 지우기).
- **🎯 결정론적 명시적 레이아웃 (Deterministic Layout)**: 어설픈 AI 추측을 배제하고 지시어(`cover`, `two-cols`, `default`)에 따라 100% 예측 가능한 레이아웃 렌더링.
- **⚡ 실시간 Live Preview 개발 서버**: Go 표준 라이브러리 기반 SSE (Server-Sent Events) 핫 리로드 지원 (`goslide serve`). 편집 중 슬라이드 위치 유지.
- **🖥️ 듀얼 스크린 발표자 콘솔**: 듀얼 모니터 발표자 전용 뷰(`P` 키), 레이저 포인터(`L` 키), 스포트라이트 집중 모드(`S` 키).
- **💻 풍부한 기술 콘텐츠 지원**: Chroma 기반 순수 Go 구문 강조, KaTeX 인라인/블록 수식(`$...$`, `$$...$$`), GitHub Flavored Markdown (GFM) 표 완벽 렌더링.

---

## 🚀 빠른 시작 (Quick Start)

### 설치 방법

```bash
# go install 사용 (Go 1.22+)
go install github.com/joincdream/goslide/cmd/goslide@latest

# 또는 소스코드 직접 빌드
git clone https://github.com/joincdream/goslide.git
cd goslide
make build
```

### 슬라이드 마크다운 작성 예시

`presentation.md` 파일을 작성합니다:

````markdown
---
title: "현대 클라우드 아키텍처"
theme: default
paginate: true
header: "Goslide 기술 발표"
footer: "© 2026 Cloud IT"
---

# 현대 클라우드 아키텍처
### 대규모 분산 시스템 설계와 구현

<!-- layout: cover -->

---

## 📌 마이크로서비스 파이프라인

- gRPC 기반 고성능 마이크로서비스
- NATS를 활용한 실시간 이벤트 스트리밍
- 쿠버네티스 오케스트레이션 자동화

<!-- layout: two-cols -->

<!-- split -->

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Goslide!")
}
```

<!-- note: 분리된 아키텍처의 장점과 성능 이점을 설명할 것 -->
````

### 🤖 AI / LLM으로 슬라이드 자동 생성

ChatGPT, Claude 등의 LLM 프롬프트에 아래 문장을 그대로 입력하면 Goslide 규격에 맞는 슬라이드를 즉시 생성할 수 있습니다:

> "https://raw.githubusercontent.com/joincdream/goslide/main/examples/dsl/example-dsl.ko.md 
> 문서를 참고해서 [원하는 발표 주제]에 대한 슬라이드를 마크다운으로 작성해 줘."

👉 지원하는 레이아웃, 배경 처리, 수식 및 세부 지시어는 [Goslide DSL 규격 및 골든 예제 가이드](examples/dsl/example-dsl.ko.md)에서 확인하실 수 있습니다.

### 슬라이드 빌드

저장소에 포함된 공식 데모 슬라이드([`examples/demo/demo.ko.md`](examples/demo/demo.ko.md))를 바로 빌드할 수 있습니다:

```bash
# 독립형 인터랙티브 HTML 생성
goslide build examples/demo/demo.ko.md -o index.html

# 16:9 무마진 벡터 PDF 생성
goslide build examples/demo/demo.ko.md -f pdf -o demo.pdf

# 파워포인트(PPTX) 파일 생성
goslide build examples/demo/demo.ko.md -f pptx -o demo.pptx
```

### 로컬 라이브 프리뷰 서버 구동

```bash
# SSE 기반 실시간 변경 감지 개발 서버 실행
goslide serve examples/demo/demo.ko.md --port 8080
```

웹 브라우저에서 `http://localhost:8080`을 열면 `examples/demo/demo.ko.md` 수정 저장 시 보고 있던 슬라이드 위치를 유지하며 300ms 이내에 즉시 새로고침됩니다.

---

## ⌨️ 단축키 안내 (프레젠테이션 & 스크린캐스트)

| 단축키 | 동작 설명 | 활용 상황 |
| :---: | :--- | :--- |
| `→` / `Space` / `PageDown` | 다음 슬라이드로 이동 | 슬라이드 진행 |
| `←` / `PageUp` | 이전 슬라이드로 이동 | 슬라이드 진행 |
| `Home` / `End` | 첫 번째 / 마지막 슬라이드로 이동 | 슬라이드 진행 |
| `[슬라이드 번호]` + `Enter` | 지정한 슬라이드로 즉시 점프 | 슬라이드 진행 |
| `F` | 전체 화면 모드 토글 | 프레젠테이션 |
| **`D`** | **투명 캔버스 판서 펜 켜기 / 끄기** | **유튜브 녹화 / 실시간 판서** |
| **`C`** | **현재 슬라이드 필기 내용 전체 지우기** | **유튜브 녹화 / 실시간 판서** |
| `P` | 듀얼 스크린 발표자 콘솔 분리 창 열기 | 발표자 모드 |
| `L` | 디지털 레이저 포인터 켜기 / 끄기 | 발표자 모드 |
| `S` | 스포트라이트 마우스 집중 영역 토글 | 발표자 모드 |

---

## 📐 지시어 및 레이아웃 사용법

### Frontmatter (전역 지시어)

문서 최상단 YAML 블록에 프레젠테이션 메타데이터를 선언합니다:

```yaml
---
title: 프레젠테이션 제목
theme: default        # 테마: default, clean, dark
paginate: true        # 페이지 번호 표시 (예: 3 / 24)
header: 머리말 텍스트
footer: 꼬리말 텍스트
size: 16:9            # 비율: 16:9 (기본값) 또는 4:3
---
```

### 슬라이드 지시어 (Scoped Directives)

HTML 주석을 사용하여 개별 슬라이드의 레이아웃과 스타일을 제어합니다:

- `<!-- layout: cover -->`: 타이틀 및 간지용 중앙 정렬 레이아웃.
- `<!-- layout: two-cols -->` (또는 `<!-- split -->`): 좌/우 2단 그리드 컬럼 레이아웃.
- `<!-- backgroundColor: #1e1e1e -->`: 특정 슬라이드의 배경색 오버라이드.
- `<!-- note: 발표자 메모 본문 -->`: 발표자 콘솔 및 PPTX 슬라이드 메모에 기록되는 발표자 대본.

---

## 🏗️ 시스템 아키텍처

Goslide는 결합도가 낮은 단방향 파이프라인 아키텍처를 준수합니다:

```
[Markdown Source]
       │
       ▼ (1. Parse)
[불변 Slide IR] (Deck, Slide, Directives)
       │
       ▼ (2. Transform & Resolve)
[Themed Presentation Model]
       │
       ├───────────────────┼───────────────────┐
       ▼ (3a. Render)      ▼ (3b. Export)      ▼ (3c. Export)
  HTML Renderer       PDF Exporter        PPTX Exporter
  (html/template)       (chromedp)        (OPC zip builder)
```

상세 스펙 및 엔지니어링 가이드라인:
- [코어 아키텍처 명세서](docs/okf/core-architecture.md)
- [인터페이스 및 도메인 모델 계약](docs/okf/contracts-interfaces.md)
- [스크린캐스트 판서 오버레이 규격](docs/okf/screencast-annotation.md)
- [엔지니어링 제약사항 및 DoD 체크리스트](docs/okf/hard-constraints.md)

---

## 🛠️ 개발 및 기여 가이드

```bash
# 단위 테스트 및 동시성 Race Condition 검증
make test-race

# 정적 분석 린트 수행
make lint

# 대용량 슬라이드 변환 성능 벤치마크
make bench
```

---

## 📄 라이선스

본 프로젝트는 [Apache 2.0 라이선스](LICENSE)를 따릅니다.
