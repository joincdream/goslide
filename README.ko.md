<div align="center">

# Goslide 📽️

**순수 Go(Pure Go)로 구현된 초고속 단일 바이너리 마크다운 슬라이드 데크 빌더**

외부 런타임 의존성 없이 마크다운 문서를 인터랙티브 HTML, 벡터 PDF, 파워포인트(PPTX) 파일로 즉시 변환합니다.  
개발자 기술 발표, IT 교육 강의 및 유튜브 스크린캐스트 동영상 강의 녹화에 최적화되어 있습니다.

[English](README.md) | [한국어](README.ko.md) · [3분 퀵스타트](QUICKSTART.ko.md)

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

> 💡 3분 안에 첫 슬라이드를 완성하는 전체 가이드는 **[3분 퀵스타트 가이드](QUICKSTART.ko.md)**를 참고하세요.

### 1. 설치 방법

#### 바이너리 다운로드 (권장)
[GitHub Releases](https://github.com/joincdream/goslide/releases)에서 운영체제(Linux, macOS, Windows)에 맞는 압축 파일을 다운로드하여 압축을 풀면 바로 실행할 수 있습니다.

#### Go 도구로 설치 (Go 1.22+)
```bash
go install github.com/joincdream/goslide/cmd/goslide@latest
```

#### 소스코드 직접 빌드
```bash
git clone https://github.com/joincdream/goslide.git
cd goslide
make build
```

---

### 2. 3분 실전 체험 (Demo & Live Reload)

Goslide는 바이너리 자체에 완성도 높은 14장 규모의 쇼케이스 프레젠테이션과 커스텀 테마를 내장하고 있습니다:

```bash
# 1. 데모 슬라이드(demo/) 및 프로젝트 공통 테마(themes/) 언팩
goslide demo

# 2. 실시간 라이브 프리뷰 개발 서버 시작 (브라우저 자동 실행)
goslide serve demo/demo.md

# 3. HTML, 16:9 벡터 PDF, 파워포인트 PPTX 일괄 동시 생성
goslide build demo/demo.md -f html,pdf,pptx

# 4. 나만의 새 슬라이드 생성
goslide init my-slide.md
```

---

### 3. AI / LLM으로 슬라이드 및 테마 자동 생성

ChatGPT, Claude, Cursor 등의 AI 에이전트에게 슬라이드 작성이나 커스텀 테마 디자인을 맡기고 싶다면:

#### 1) 슬라이드 마크다운 생성
```bash
# 슬라이드 생성용 AI 시스템 프롬프트 출력 (공식 DSL 규격 URL 포함)
goslide prompt
```
출력된 프롬프트를 복사하여 AI 도구에 붙여넣고 원하는 발표 주제(예: *"클라우드 네이티브 마이크로서비스 설계 8장 작성해 줘"*)를 입력하면, 완벽한 Goslide 마크다운 코드가 생성됩니다.  
👉 [공식 Goslide Markdown DSL 규격서 (docs/dsl-guide.md)](docs/dsl-guide.md)

#### 2) 커스텀 테마 CSS 생성
```bash
# 테마 CSS 생성용 AI 시스템 프롬프트 출력 (공식 4-Tier DOM 규격 URL 포함)
goslide prompt theme
```
출력된 프롬프트를 AI 도구에 붙여넣고 원하는 브랜드 색상과 스타일을 요청하면, 1920×1080 뷰포트에 완벽히 최적화된 독립형 커스텀 테마 CSS가 생성됩니다.  
👉 [공식 Goslide 테마 CSS 아키텍처 규격서 (docs/theme-guide.md)](docs/theme-guide.md)

---

## ⌨️ 단축키 안내 (프레젠테이션 & 스크린캐스트)

| 단축키 | 동작 설명 | 활용 상황 |
| :---: | :--- | :--- |
| `→` / `Space` / `PageDown` / `J` | 다음 슬라이드로 이동 | 슬라이드 진행 |
| `←` / `PageUp` / `K` / `H` | 이전 슬라이드로 이동 | 슬라이드 진행 |
| `Home` / `End` | 첫 번째 / 마지막 슬라이드로 이동 | 슬라이드 진행 |
| `[슬라이드 번호]` + `Enter` | 지정한 슬라이드로 즉시 점프 | 슬라이드 진행 |
| `F` | 전체 화면 모드 토글 | 프레젠테이션 |
| **`D`** | **투명 캔버스 판서 펜 켜기 / 끄기** | **유튜브 녹화 / 실시간 판서** |
| **`C`** | **현재 슬라이드 필기 내용 전체 지우기** | **유튜브 녹화 / 실시간 판서** |
| `L` | 디지털 레이저 포인터 켜기 / 끄기 | 발표자 모드 |
| `S` | 스포트라이트 마우스 집중 영역 토글 | 발표자 모드 |
| `P` | 듀얼 스크린 발표자 콘솔 분리 창 열기 | 발표자 모드 |
| `B` / `W` | 슬라이드 블랙아웃(암전) / 화이트아웃(백전) | 청중 집중 유도 |
| `1` ~ `4` | 판서 펜 / 레이저 색상 변경 (빨강, 파랑, 초록, 노랑) | 판서 / 포인터 |
| `+` / `-` | 판서 펜 / 레이저 굵기 변경 (얇게, 보통, 굵게) | 판서 / 포인터 |
| `?` | 전체 단축키 도움말 팝업 열기 | 도움말 |

---

## 📐 지시어 및 레이아웃 사용법

### Frontmatter (전역 지시어)

문서 최상단 YAML 블록에 프레젠테이션 메타데이터를 선언합니다:

```yaml
---
title: 프레젠테이션 제목
theme: clean          # 테마: clean, dark, academic, cyber-dark, default
size: 16:9            # 비율: 16:9 (기본값) 또는 4:3
paginate: true        # 페이지 번호 표시 (예: 3 / 24)
header: 머리말 텍스트
footer: 꼬리말 텍스트
autofit: true         # 텍스트 오버플로우 방지 자동 축소
---
```

### 슬라이드 지시어 (Scoped Directives)

HTML 주석을 사용하여 개별 슬라이드의 레이아웃과 스타일을 제어합니다:

- `<!-- layout: cover -->`: 타이틀 및 간지용 중앙 정렬 레이아웃.
- `<!-- layout: two-cols -->` (또는 `<!-- split -->`): 좌/우 2단 그리드 컬럼 레이아웃.
- `<!-- layout: section -->`: 대주제 구분용 섹션 레이아웃.
- `<!-- layout: blank -->`: 여백 및 패딩이 없는 전체 캔버스 레이아웃 (대형 이미지/다이어그램용).
- `<!-- backgroundColor: #1e1e1e -->`: 특정 슬라이드의 배경색 오버라이드.
- `<!-- backgroundImage: url(bg.jpg) -->`: 슬라이드 배경 이미지 적용.
- `<!-- note: 발표자 메모 본문 -->`: 발표자 콘솔 및 PPTX 슬라이드 메모에 기록되는 발표자 대본.

---

## 🎨 테마 관리 및 커스텀 테마 제작

Goslide는 100% CSS 기반으로 동작하므로 나만의 브랜드 스타일 테마를 손쉽게 제작하고 적용할 수 있습니다.

### 1. 내장 테마 목록 조회 및 추출
```bash
# 사용 가능한 내장 테마 목록 조회 (clean, dark, academic, cyber-dark, default)
goslide theme list

# 내장 테마 CSS를 로컬 파일로 추출
goslide theme export clean themes/my-company.css
```

### 2. 커스텀 테마 자동 인식
프로젝트의 `themes/` 폴더에 CSS 파일(예: `themes/my-company.css`)로 저장하면, 마크다운 Frontmatter에서 이름만으로 즉시 적용됩니다:
```yaml
---
title: "기업 IR 자료"
theme: my-company      # themes/my-company.css 자동 탐색
---
```

### 3. 주요 CSS 셀렉터 구조
| 셀렉터 | 설명 | 기본 높이 / 속성 |
| :--- | :--- | :--- |
| `.slide-card` | 슬라이드 캔버스 컨테이너 | 배경색, 텍스트 기본색, 테두리, 그림자 |
| `.slide-card .slide-header` | 상단 고정 머리말 텍스트 | 40px |
| `.slide-card .slide-title-box` | 타이틀 및 서브타이틀 영역 | 120px |
| `.slide-card .slide-body` | 슬라이드 본문 플렉스박스 | 840px ~ 960px |
| `.slide-card .slide-footer` | 하단 고정 꼬리말 및 페이지 번호 | 80px |
| `.slide-card.layout-cover` | 표지 슬라이드 전용 스타일 | 중앙 정렬 H1 |
| `.slide-card.layout-two-cols` | 2단 컬럼 분할 그리드 | `.col-left`, `.col-right` |

👉 더 자세한 테마 DOM 아키텍처 및 하드 가이드라인은 **[Goslide 테마 CSS 아키텍처 규격서 (docs/theme-guide.md)](docs/theme-guide.md)**를 참조하세요.

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
# 전체 품질 검증 (포맷, 린트, 복잡도, 동시성 레이스 테스트)
make check

# 단위 테스트 및 동시성 Race Condition 검증
make test-race

# 정적 분석 린트 수행
make lint

# 대용량 슬라이드 변환 성능 벤치마크
make bench
```

---

## 👤 제작자 및 문의 (Author & Contact)

**yundream (조인씨 / Joinc)**

- 🌐 **웹사이트**: [https://www.joinc.co.kr](https://www.joinc.co.kr)
- 💼 **LinkedIn**: [linkedin.com/in/yundream](https://www.linkedin.com/in/yundream/)
- ✉️ **이메일**: [yundream@gmail.com](mailto:yundream@gmail.com)
- 🐙 **GitHub**: [@joincdream](https://github.com/joincdream)

---

## 📄 라이선스

본 프로젝트는 [Apache 2.0 라이선스](LICENSE)를 따릅니다.
