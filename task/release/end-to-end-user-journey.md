# Goslide 엔드투엔드(End-to-End) 사용자 여정 및 릴리스 갭 분석

> **문서 목적**: Goslide의 신규 사용자가 제품을 최초 인지(방문)하는 순간부터 다운로드, 압축 해제, 설치, 첫 실행, 슬라이드 저작, AI 활용, 최종 빌드 및 발표/녹화까지 이어지는 전체 여정(User Journey)을 단계별로 분석하고, 각 접점에서 발생할 수 있는 마찰 지점(Friction Points) 및 누락된 요구사항(Gaps)을 선제적으로 발굴·개선하기 위함입니다.  
> **문서 위치**: `task/release/end-to-end-user-journey.md`  
> **대상 릴리스**: Goslide v1.0.0 (Linux, macOS, Windows)  
> **작성 일자**: 2026-10-10  

---

## 1. 전체 사용자 여정 요약 (Journey Map Overview)

Goslide의 사용자 경험은 아래의 7단계 생명주기를 거칩니다:

```
[1. 탐색 & 인지] ───► [2. 릴리스 다운로드] ───► [3. 압축 해제 & 설치] ───► [4. 3분 첫 체험 (Aha!)]
     (사이트/저장소)          (OS/Arch 선택)            (PATH/권한 설정)            (demo & serve)
                                                                                     │
[7. 발표 & 스크린캐스트] ◄── [6. 다중 포맷 빌드] ◄── [5. 저작 & AI 협업] ◄──────────┘
  (듀얼 콘솔/판서 녹화)        (HTML/PDF/PPTX)          (Live Reload/DSL)
```

---

## 2. 단계별 세부 사용자 여정 및 갭 분석 (Gap Analysis)

---

### 단계 1. 탐색 및 웹사이트/저장소 방문 (Discovery & Landing)

* **사용자 시나리오**:
  1. 개발자/강사/유튜버가 GitHub 저장소(`github.com/joincdream/goslide`)나 기술 블로그(`joinc.co.kr`)의 소개 글을 통해 유입됩니다.
  2. 첫 10초 이내에 "이 도구가 내게 필요한 도구인가?", "기존 Marp나 Slidev와 무엇이 다른가?"를 파악하려 합니다.
  3. README 상단의 설치 및 다운로드 링크를 찾습니다.

* **사용자 액션 & 터치포인트**:
  * [README.ko.md](file:///home/yundream/myjob/cloit/Goslide/README.ko.md), [README.md](file:///home/yundream/myjob/cloit/Goslide/README.md) 열람
  * [QUICKSTART.ko.md](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.ko.md) 링크 클릭
  * `[GitHub Releases]` 다운로드 링크 클릭

* **현재 상태 & 강점**:
  * 단일 바이너리, 순수 Go, CGO 의존성 0%라는 명확한 엔지니어링 가치 전달.
  * 3분 퀵스타트 문서와 다국어(한/영) 지원 완비.

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **시각적 데모(GIF / 스크린샷) 부재**:
     - *현상*: 텍스트 설명은 훌륭하나, 실제 동작 화면(실시간 핫리로드, 판서 오버레이, 듀얼 발표자 뷰)을 보여주는 애니메이션 GIF나 스크린샷이 README 상단에 없어 시각적 흡입력이 떨어짐.
  2. **시스템 요구사항(사전 준비물) 안내 위치**:
     - *현상*: 순수 Go 단일 바이너리이므로 Go 언어나 Node.js는 필요 없지만, **PDF/PPTX 익스포트를 위해서는 시스템에 Chrome/Chromium이 필요**합니다. 이 사전 조건이 다운로드 전 명확히 인지되지 않으면 추후 빌드 단계에서 당황할 수 있음.

---

### 단계 2. 릴리스 아티팩트 다운로드 (Download Release Artifact)

* **사용자 시나리오**:
  1. GitHub Releases 페이지(`releases/tag/v1.0.0`)에 진입합니다.
  2. 자신의 운영체제(Linux, macOS, Windows)와 CPU 아키텍처(x86_64, arm64)에 맞는 파일을 찾아 다운로드합니다.
  3. (보안 중시 사용자) `checksums.txt`를 통해 파일 무결성을 검증합니다.

* **제공 아티팩트 목록 (Makefile cross-build 기준)**:
  * Linux: `goslide_linux_amd64.tar.gz`, `goslide_linux_arm64.tar.gz`
  * macOS: `goslide_darwin_amd64.tar.gz`, `goslide_darwin_arm64.tar.gz`
  * Windows: `goslide_windows_amd64.zip`, `goslide_windows_arm64.zip`
  * 무결성 검증: `checksums.txt`

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **아키텍처 명칭 혼선 (특히 Mac 사용자)**:
     - *현상*: Apple Silicon(M1/M2/M3/M4) 사용자는 `arm64`를 받아야 하고, Intel Mac은 `amd64`를 받아야 합니다. Release Note 본문에 표 형태로 매핑 가이드(`Apple Silicon Mac ➔ goslide_darwin_arm64.tar.gz`)가 명시되어야 함.
  2. **원라이너(One-liner) 간편 설치 스크립트 부재**:
     - *현상*: 터미널 환경에서 작업하는 리눅스/맥 개발자는 웹 브라우저에서 다운로드하는 것보다 `curl -fsSL https://.../install.sh | bash` 형태를 선호합니다. 현재는 수동 다운로드 또는 `go install`만 제시되어 있음.
  3. **체크섬 검증 방법 안내 부재**:
     - *현상*: `checksums.txt`는 생성되지만, 이를 터미널에서 어떻게 확인하는지(`sha256sum -c checksums.txt`)에 대한 가이드가 README/Release Note에 없음.

---

### 단계 3. 압축 풀기 및 실행 환경(PATH) 설정 (Extraction & Installation)

* **사용자 시나리오**:
  1. 다운로드 폴더에서 압축을 풉니다.
  2. 아카이브 안의 내용물을 확인합니다.
  3. 터미널 어디서나 `goslide`를 실행할 수 있도록 시스템 PATH 경로로 복사합니다.
  4. `goslide --help` 또는 `goslide version`을 실행하여 정상 동작을 확인합니다.

* **아카이브 내부 구성 (Makefile package 기준)**:
  * 실행 바이너리: `goslide` (Windows: `goslide.exe`)
  * 안내 문서: `README.md`, `README.ko.md`, `QUICKSTART.md`, `QUICKSTART.ko.md`, `LICENSE`
  * 프롬프트 팩: `prompts/` (AI 슬라이드/테마 프롬프트 템플릿)

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **실행 권한 누락 위험 (`chmod +x`)**:
     - *현상*: 일부 리눅스/맥 환경이나 GUI 압축 해제 툴을 사용할 경우 실행 권한(`+x`)이 유실되어 `./goslide` 실행 시 `Permission denied` 에러가 발생할 수 있습니다. 퀵스타트에 `chmod +x goslide` 안내를 기본 포함해야 안전함.
  2. **PATH 복사 경로 가이드 누락**:
     - *현상*: 초보 사용자는 다운로드한 폴더에서만 바이너리가 실행되는 문제를 겪습니다. `sudo cp goslide /usr/local/bin/` 또는 `cp goslide ~/.local/bin/` 명령어 가이드가 필수적입니다.
  3. **macOS 게이트키퍼(Gatekeeper) 격리 차단**:
     - *현상*: Apple 공증(Notarization)이 없는 독립 바이너리는 macOS에서 더블클릭이나 터미널 실행 시 `"goslide는 확인되지 않은 개발자가 배포했기 때문에 열 수 없습니다"`라는 경고와 함께 실행이 차단됩니다.  
       ➔ **해결책**: `xattr -d com.apple.quarantine goslide` 해제 가이드 안내 필요.
  4. **버전 플래그 일관성**:
     - *현상*: 사용자가 설치 후 가장 먼저 치는 명령은 `goslide -v`, `goslide version`, `goslide --version`입니다. 모든 형태의 버전 조회가 자연스럽게 지원되는지 보장되어야 함.

---

### 단계 4. 최초 실행 및 "Aha!" 모멘트 (First Run & 3-Min Onboarding)

* **사용자 시나리오**:
  1. [QUICKSTART.ko.md](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.ko.md)를 따라 임의의 빈 디렉토리로 이동합니다.
  2. `goslide demo`를 실행하여 내장된 14장 쇼케이스 슬라이드와 테마를 언팩합니다.
  3. `goslide serve demo/demo.md`를 실행합니다.
  4. 기본 브라우저가 자동으로 열리며 세련된 디자인의 슬라이드가 눈앞에 나타나는 "Aha!" 모멘트를 경험합니다.

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **헤드리스/WSL2/SSH 원격 환경에서의 브라우저 자동 오픈 실패**:
     - *현상*: 원격 리눅스 서버나 WSL2 환경에서는 GUI 브라우저 자동 실행(`xdg-open` 등)이 실패하거나 에러 로그를 뱉을 수 있습니다.  
       ➔ **대응**: 브라우저 오픈 실패 시에도 터미널에 `웹 브라우저에서 http://localhost:8080 으로 접속하세요`라는 접속 URL이 눈에 띄게 강조되어야 함.
  2. **포트 충돌(Port Collision) 대처**:
     - *현상*: 이미 8080 포트를 다른 프로세스(Tomcat, 다른 개발 서버 등)가 점유 중일 때 `goslide serve`가 즉시 크래시되는지, 아니면 다음 가용 포트(8081)로 자동 폴백하거나 친절한 에러(`--port=8888 옵션으로 포트를 변경하세요`)를 안내하는지 점검 필요.
  3. **첫 화면 내 단축키 안내 힌트**:
     - *현상*: 브라우저가 열렸을 때 다음 장으로 넘기는 방법이나 판서 모드(`D`), 발표자 모드(`P`), 도움말(`?`)을 모르는 사용자를 위해, 화면 좌측 하단이나 우측 상단에 작게 `Press ? for shortcuts` 안내가 노출되는지 점검 필요.

---

### 단계 5. 나만의 슬라이드 저작 및 실시간 편집 (Authoring & Live Reload)

* **사용자 시나리오**:
  1. `goslide init my-presentation.md`로 자신만의 새 슬라이드 파일을 생성합니다.
  2. `goslide serve my-presentation.md`를 띄워두고 VS Code / Cursor 등의 텍스트 에디터로 마크다운을 편집합니다.
  3. 저장(`Ctrl+S`)할 때마다 브라우저가 300ms 이내에 자동 갱신(SSE Hot Reload)되는 것을 체감합니다.
  4. AI의 도움을 받기 위해 `goslide prompt`를 터미널에서 실행하여 시스템 프롬프트를 복사합니다.

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **상대 경로 로컬 이미지 에셋 서빙 (Asset Serving)**:
     - *현상*: 마크다운 본문에 `![Architecture](images/arch.png)`를 넣었을 때, `goslide serve` 서버가 마크다운 파일의 기준 디렉토리(CWD) 내 로컬 정적 파일(`images/`)을 올바르게 호스팅하는지 확인 필요.
  2. **문법 오류 시 화면 피드백 (Error Overlay)**:
     - *현상*: YAML Frontmatter 오타나 잘못된 지시어 입력 시 개발 서버 프로세스가 종료(Panic)되면 안 되며, 브라우저 상단에 에러 오버레이 카드가 떠서 몇 번째 줄에 오류가 있는지 즉시 안내해야 함.
  3. **긴 프롬프트 터미널 복사의 번거로움**:
     - *현상*: `goslide prompt` 실행 시 수십 줄의 마크다운 프롬프트가 출력되는데, 터미널 스크롤을 올려 수동으로 드래그 복사하는 과정이 번거로울 수 있음.  
       ➔ **개선안**: 아카이브에 `prompts/slide-prompt.md` 파일이 동봉되어 있음을 안내하거나, 파이프라인 연계(`goslide prompt | pbcopy` 등) 팁 제공.

---

### 단계 6. 다중 포맷 일괄 변환 (Multi-Format Build & Export)

* **사용자 시나리오**:
  1. 작성이 완료된 슬라이드를 배포하기 위해 빌드 명령을 실행합니다:
     ```bash
     goslide build my-presentation.md -f html,pdf,pptx
     ```
  2. `my-presentation.html`, `my-presentation.pdf`, `my-presentation.pptx` 3종 파일이 생성되었는지 확인합니다.
  3. 브라우저 없이 독립 실행되는 HTML, 인쇄용 고해상도 PDF, 파워포인트에서 편집 가능한 PPTX를 확인합니다.

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **시스템 Chrome/Chromium 부재 시 에러 메시지 친절도 (가장 중요)**:
     - *현상*: PDF와 PPTX 변환에는 `chromedp`가 내부적으로 사용됩니다. 만약 사용자의 시스템에 크롬이 없다면 `exec: "google-chrome": executable file not found in $PATH` 같은 원시 에러가 날 수 있습니다.  
       ➔ **개선안**: 크롬이 없을 경우 `"PDF 및 PPTX 생성을 위해 Google Chrome 또는 Chromium이 필요합니다. 설치 후 다시 시도해 주세요 (Ubuntu: sudo apt install -y chromium-browser)"`와 같이 친절한 가이드형 에러로 래핑되어 있는지 확인 필수.
  2. **출력 디렉토리 지정 옵션 (`-o / --output`)**:
     - *현상*: 기본값은 마크다운 파일과 동일한 위치에 생성되지만, CI/CD 배포 파이프라인이나 `dist/` 폴더로 모으고 싶은 사용자를 위한 `-o dist/` 플래그 지원 여부 확인.
  3. **한글 폰트 렌더링 (Tofu 현상)**:
     - *현상*: 리눅스 서버나 최소 설치 환경에서 PDF 변환 시 한글 폰트가 설치되어 있지 않으면 글자가 네모(`□□□`)로 깨질 수 있습니다. 배포 문서에 폰트 패키지(`fonts-nanum` 등) 설치 권고 명시 필요.

---

### 단계 7. 실전 발표 진행 및 스크린캐스트 강의 녹화 (Presentation Delivery)

* **사용자 시나리오**:
  1. **현장 발표 (발표자 모드)**:
     - 프로젝터에는 슬라이드 전체화면(`F` 키)을 띄우고, 노트북 화면에는 `P` 키를 눌러 듀얼스크린 발표자 콘솔(스피커 노트, 경과 시간, 다음 슬라이드 미리보기)을 엽니다.
     - 청중의 시선을 유도하기 위해 레이저 포인터(`L`)나 스포트라이트(`S`)를 조작합니다.
  2. **온라인 강의 녹화 (유튜브 스크린캐스트)**:
     - OBS나 화면 녹화 프로그램을 켜고, 슬라이드를 넘기며 설명합니다.
     - 강조가 필요한 부분에서 `D` 키를 눌러 투명 캔버스 판서 모드를 켜고 마우스/태블릿으로 동그라미나 밑줄을 칩니다.
     - 설명이 끝나면 `C` 키를 눌러 필기 내용을 깔끔하게 지우고 다음 슬라이드로 넘어갑니다.

* **발견된 잠재적 마찰 지점 (Friction Points) & 누락 사항 (Gaps)**:
  1. **판서 모드 상태에서의 슬라이드 넘김 간섭**:
     - *현상*: `D` 키로 판서 모드가 활성화되어 마우스 클릭이 드로잉으로 인식될 때, 키보드 화살표(`→`)나 `Space`로 슬라이드 이동이 정상적으로 동작하는지, 혹은 판서 모드임을 알 수 있는 시각적 커서/UI 표시가 충분한지 점검.
  2. **오프라인 환경(비행기, 외부 세미나장) 완결성**:
     - *현상*: 빌드된 `my-presentation.html`을 들고 인터넷이 안 되는 발표장에 갔을 때, 외부 CDN(KaTeX 폰트, Chroma 스타일, 자바스크립트 등)에 의존하고 있어 수식이나 스타일이 깨지지 않는지 완벽한 Self-contained 여부 보장 필요. (현재 Goslide는 `embed.FS` 기반으로 오프라인 완결성을 보장함)

---

## 3. 릴리스 전 체크리스트 및 우선순위 개선 과제 (Action Items)

사용자 여정 전수 조사 결과, v1.0.0 공식 릴리스 전 반드시 챙겨야 할 액션 아이템 매트릭스입니다:

| 구분 | 과제 내용 | 대상 컴포넌트 / 문서 | 중요도 | 상태 |
| :--- | :--- | :--- | :---: | :---: |
| **설치 편의성** | README/Release Note에 OS 및 CPU별 아티팩트 매핑 표 제공 (Apple Silicon `arm64` 안내) | `README.md`, `README.ko.md`, GitHub Release | P1 | 점검 필요 |
| **설치 편의성** | 압축 해제 후 `chmod +x goslide` 및 PATH 복사(`~/.local/bin`) 안내 명시 | `QUICKSTART.ko.md`, `README.ko.md` | P1 | 완료 (보강) |
| **보안/신뢰** | macOS 게이트키퍼 차단 해제 가이드(`xattr -d com.apple.quarantine`) 안내 | `QUICKSTART.md`, `QUICKSTART.ko.md` | P1 | 신규 추가 권장 |
| **빌드 안정성** | 시스템 크롬 미설치 시 PDF/PPTX 변환 에러 메시지의 친절한 설치 안내 래핑 확인 | `internal/exporter/pdf`, `pptx` | P1 | 코드 점검 |
| **시각적 경험** | README 상단에 판서 오버레이 및 듀얼 콘솔을 보여주는 대표 스크린샷 1장 추가 | `README.md`, `README.ko.md` | P2 | 권장 |
| **온보딩** | `goslide serve` 구동 시 터미널 출력에 눈에 띄는 접속 URL 및 단축키 도움말(`?`) 강조 | `cmd/goslide/serve.go` | P2 | 점검 완료 |
| **배포 자동화** | 향후 릴리스를 위한 원라이너 설치 스크립트(`install.sh`) 또는 Homebrew 패키징 로드맵 수립 | `scripts/install.sh` | P3 | v1.1 백로그 |

---

## 4. 결론 및 종합 의견

Goslide는 **"단일 정적 바이너리 & CGO 0% & 내장 테마(embed.FS)"**라는 독보적인 아키텍처 덕분에 타 마크다운 슬라이드 도구(Node.js, npm, Python 등의 복잡한 런타임 설치 필요)에 비해 **온보딩 허들이 압도적으로 낮습니다.**

사용자의 전체 여정 중 가장 주의해야 할 3대 마찰 지점은:
1. **다운로드 단계**: macOS 사용자(M시리즈 vs Intel)의 아키텍처 선택 및 게이트키퍼 경고 해제.
2. **설치 단계**: 압축 해제 후 시스템 PATH(`~/.local/bin`) 등록 안내.
3. **익스포트 단계**: PDF/PPTX 빌드 시 시스템 Chrome 브라우저 존재 여부 및 에러 안내.

본 여정 맵에서 도출된 가이드라인을 Release Note 및 온보딩 문서에 충실히 반영하면, 신규 사용자가 3분 이내에 성공적으로 안착(First Successful Presentation)할 수 있을 것입니다.
