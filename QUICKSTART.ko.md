<div align="center">

# Goslide 3분 퀵스타트 가이드

[English](QUICKSTART.md) | [한국어](QUICKSTART.ko.md)

</div>

---

**Goslide**에 오신 것을 환영합니다! Goslide는 순수 Go 기반의 단일 바이너리 마크다운 슬라이드 데크 빌더입니다.  
아래 단계를 따라 하시면 **3분 이내**에 첫 프레젠테이션을 완성하고 PDF/PPTX로 변환할 수 있습니다.

---

## 0단계. 바이너리 다운로드 및 설치 (30초)

### 가장 빠른 원라이너 설치 (Linux & macOS)
터미널에서 자동 설치 스크립트를 실행합니다:
```bash
curl -fsSL https://raw.githubusercontent.com/joincdream/goslide/main/install.sh | bash
```

이미 설치되어 있다면 언제든 최신 버전으로 업데이트: `goslide update`

---

### 수동 설치 (모든 플랫폼)
[GitHub Releases](https://github.com/joincdream/goslide/releases)에서 운영체제에 맞는 아티팩트를 다운로드합니다:

| 운영체제 (OS) | CPU 아키텍처 / 대상 기기 | 권장 다운로드 아티팩트 |
| :--- | :--- | :--- |
| **macOS** | **Apple Silicon (M1 / M2 / M3 / M4)** | `goslide_darwin_arm64.tar.gz` |
| **macOS** | Intel 프로세서 | `goslide_darwin_amd64.tar.gz` |
| **Linux** | 64-bit Intel / AMD (x86_64) | `goslide_linux_amd64.tar.gz` |
| **Linux** | 64-bit ARM (aarch64 / 라즈베리파이 4, 5) | `goslide_linux_arm64.tar.gz` |
| **Windows** | 64-bit Intel / AMD (x64) | `goslide_windows_amd64.zip` |
| **Windows** | 64-bit ARM (Snapdragon X Elite 등) | `goslide_windows_arm64.zip` |

**Linux & macOS 수동 설치:**
```bash
# 1. 아카이브 압축 해제 (예: Linux amd64 기준)
tar -xzf goslide_linux_amd64.tar.gz

# 2. 실행 권한 부여 및 시스템 PATH 경로로 복사
chmod +x goslide
mkdir -p ~/.local/bin && cp goslide ~/.local/bin/

# 3. 설치 확인
goslide --help
```

> 참고 (macOS 사용자): 미확인 개발자 보안 알림으로 실행이 차단되는 경우 터미널에서 다음 명령을 실행하세요:  
> `xattr -d com.apple.quarantine goslide`

**Windows 설치:**
- 다운로드한 `.zip` 압축을 해제한 후 `goslide.exe`를 원하는 경로(예: `C:\bin`)로 이동하고 시스템 환경 변수 `Path`에 해당 경로를 추가합니다.

> [!NOTE]
> 인터랙티브 HTML 슬라이드 및 실시간 프리뷰(`goslide serve`)는 **외부 런타임 의존성이 0%**입니다. 고해상도 벡터 **PDF** 및 파워포인트 **PPTX** 변환 시에는 시스템에 **Google Chrome 또는 Chromium** 브라우저가 필요합니다.

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

## 유용한 추가 명령어 및 팁

- `goslide theme list`: 사용 가능한 내장 테마 목록 조회
- `goslide theme export clean my-style.css`: 내장 테마 CSS를 로컬로 추출하여 커스텀
- `goslide --help`: 전체 CLI 옵션 및 도움말 확인

### 다운로드 파일 무결성 검증 (SHA-256)
```bash
# Linux
sha256sum -c checksums.txt --ignore-missing

# macOS
shasum -a 256 -c checksums.txt --ignore-missing
```
