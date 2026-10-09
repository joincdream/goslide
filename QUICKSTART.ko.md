<div align="center">

# 🚀 Goslide 3분 퀵스타트 가이드

[English](QUICKSTART.md) | [한국어](QUICKSTART.ko.md)

</div>

---

**Goslide**에 오신 것을 환영합니다! Goslide는 순수 Go 기반의 단일 바이너리 마크다운 슬라이드 데크 빌더입니다.  
아래 5단계를 따라 하시면 **3분 이내**에 첫 프레젠테이션을 완성하고 PDF/PPTX로 변환할 수 있습니다.

---

## 1단계. 데모 환경 준비 (30초)

터미널에서 아래 명령을 실행하여 실전 데모 슬라이드(`demo/`)와 프로젝트 공통 테마(`themes/`)를 준비합니다:

```bash
# 데모 마크다운, 이미지 에셋 및 프로젝트 공통 테마 언팩
goslide demo
```

생성되는 파일 구조:
- `demo/`: 14장 규모의 쇼케이스 슬라이드(`demo.md`) 및 이미지 에셋
- `themes/`: 프로젝트 전체에서 공용으로 사용할 수 있는 커스텀 테마 스타일시트 (`clean.css`, `dark.css`, `academic.css`)

---

## 2단계. 실시간 라이브 프리뷰 및 편집 체험 (60초)

```bash
# 로컬 개발 서버 시작 (브라우저 자동 오픈)
goslide serve demo/demo.md
```

- **실시간 편집 (Hot Reload)**: 에디터에서 `demo/demo.md`를 열고 내용을 수정한 후 저장(`Ctrl+S`)해 보세요. 브라우저가 깜빡임 없이 즉시 갱신됩니다!
- **테마 변경**: `demo/demo.md` 상단 Frontmatter의 `theme: "clean"`을 `"dark"` 또는 `"academic"`으로 바꿔보거나, 루트의 `themes/` 폴더 내 CSS를 직접 수정해 보세요.
- **발표자 단축키**:
  - `P`: 분리형 듀얼스크린 발표자 뷰 (스피커 노트, 타이머, 다음 슬라이드 미리보기)
  - `D`: 투명 캔버스 판서 그리기 모드 (`C` 키로 판서 지우기)
  - `L`: 레이저 포인터 모드
  - `S`: 스포트라이트 집중 모드
  - `B` / `W`: 화면 가리기 (블랙아웃 / 화이트아웃)
  - `?`: 전체 단축키 도움말

---

## 3단계. 다양한 포맷으로 일괄 변환 (30초)

마크다운 파일 하나로 웹 슬라이드, 벡터 PDF, 파워포인트 파일을 한 번에 빌드합니다:

```bash
# HTML, PDF, PPTX 3종 동시 생성
goslide build demo/demo.md -f html,pdf,pptx
```

생성 결과물:
- `demo/demo.html`: 브라우저에서 오프라인으로 열 수 있는 반응형 웹 슬라이드
- `demo/demo.pdf`: 16:9 무마진 인쇄용 고해상도 벡터 PDF
- `demo/demo.pptx`: 발표자 노트가 완벽히 보존되는 Office PowerPoint 파일

---

## 4단계. 나만의 첫 슬라이드 만들기 (60초)

자신만의 새 슬라이드를 작성하려면 `init` 명령어를 실행합니다:

```bash
# 새 프레젠테이션 템플릿 생성
goslide init my-slide.md

# 실시간 미리보기 시작
goslide serve my-slide.md
```

---

## 5단계. AI(ChatGPT, Claude)를 활용한 슬라이드 자동 생성

AI 에이전트에게 슬라이드 작성을 맡기고 싶다면, Goslide 전용 시스템 프롬프트를 복사하여 전달하세요:

```bash
# 터미널에 AI 시스템 프롬프트 출력
goslide prompt
```

출력된 프롬프트를 ChatGPT 또는 Claude에 붙여넣고 원하는 주제(예: *"Go 언어 동시성 프로그래밍에 대한 7장짜리 세미나 슬라이드 만들어줘"*)를 입력하면, 완벽한 Goslide 마크다운 코드가 생성됩니다.

---

## 💡 유용한 추가 명령어

- `goslide theme list`: 사용 가능한 내장 테마 목록 조회
- `goslide theme export clean my-style.css`: 내장 테마 CSS를 로컬로 추출하여 커스텀
- `goslide --help`: 전체 CLI 옵션 및 도움말 확인
