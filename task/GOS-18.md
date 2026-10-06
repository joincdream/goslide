# [GOS-18] v1.0.0 정식 출시 릴리즈 및 크로스 플랫폼 E2E 검증 계획서

> **티켓 번호**: [GOS-18](https://joincdream.atlassian.net/browse/GOS-18)  
> **마일스톤**: 로드맵 2 (Release) / Milestone M2-5 (CLI 고도화, GoReleaser 배포 및 공식 릴리즈)  
> **마감일**: 2026-11-13  
> **상태**: 진행 중 (In Progress)  
> **담당 패키지**: `task/`, `docs/`, `cmd/goslide/`, `.goreleaser.yaml`, `.github/workflows/`  
> **연관 문서**: [unit-test-report.md](file:///home/yundream/myjob/cloit/Goslide/docs/release/unit-test-report.md), [pre-release-qa-plan.md](file:///home/yundream/myjob/cloit/Goslide/docs/release/pre-release-qa-plan.md), [cross-compilation-plan.md](file:///home/yundream/myjob/cloit/Goslide/task/cross-compilation-plan.md), [GOS-12.md](file:///home/yundream/myjob/cloit/Goslide/task/GOS-12.md), [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md)

---

## 1. 개요 및 목적 (Overview & Goals)

Goslide는 Go 언어의 크로스 컴파일 기능 덕분에 단일 코드베이스에서 Linux, macOS, Windows 실행파일을 모두 생성할 수 있습니다. 그러나 슬라이드 제작 도구의 특성상 **실제 데스크톱 환경에서의 E2E 사용자 경험(UX), OS별 기본 폰트 렌더링, 프레젠테이션 단축키 충돌 여부, 그리고 Microsoft PowerPoint 및 Apple Keynote에서의 PPTX 호환성**은 CI 가상머신만으로는 검증할 수 없으며, **실제 기기에서의 육안 및 기능 점검**이 필수적입니다.

본 태스크의 목적은 공식 **`v1.0.0` 정식 출시**에 앞서, **릴리즈 후보(RC - Release Candidate)**를 배포하고 3대 주요 OS 환경에서 체계적인 **실기기 E2E 스모크 테스트**를 거쳐 무결성을 확보한 뒤 최종 출시하는 릴리즈 계획(Release Plan)을 정의하고 실행하는 것입니다.

---

## 2. 릴리즈 라이프사이클 단계 (Release Lifecycle)

```
[Phase 1: 기능 완성 & CI 자동화 통과]
  - GOS-12 (init 커맨드, 다중 빌드, GoReleaser 파이프라인) 완료
  - GitHub Actions 멀티 OS 매트릭스 (Linux, macOS, Windows) 자동 빌드 성공
         │
         ▼
[Phase 2: 릴리즈 후보(RC) 태그 발행: v1.0.0-rc.1]
  - git tag -a v1.0.0-rc.1 -m "Release Candidate 1" && git push origin v1.0.0-rc.1
  - GitHub Releases에 5대 타겟 바이너리 및 아카이브 자동 업로드 확인
         │
         ▼
[Phase 3: 3대 운영체제 실기기 E2E 검증 (Visual & Functional QA)]
  - Linux (Ubuntu/WSL2)
  - macOS (Apple Silicon / Intel)
  - Windows (Windows 10/11)
  - 발견된 버그가 있을 경우 핫픽스 후 v1.0.0-rc.2 발행
         │
         ▼
[Phase 4: 공식 v1.0.0 정식 릴리즈]
  - git tag -a v1.0.0 -m "Official Release v1.0.0" && git push origin v1.0.0
  - GitHub Releases 정식 공개 및 릴리즈 노트 배포
```

---

## 3. 크로스 플랫폼 E2E 검증 매트릭스 및 시나리오 (Verification Scenarios)

### 3.1 공통 CLI 실행 및 산출물 생성 시나리오
모든 타겟 OS에서 동일한 표준 슬라이드(`examples/example-dsl.md`)를 기반으로 테스트합니다:

1. **단일 바이너리 직접 실행**:
   * 터미널에서 `goslide --version` 실행 시 버전, 커밋 해시, 빌드 일시, OS/Arch 정확히 출력되는지 확인.
2. **스타터 슬라이드 생성 (`goslide init`)**:
   * 임의의 빈 폴더에서 `goslide init my-talk.md` 실행 ➔ 표준 마크다운 생성 확인.
3. **다중 포맷 일괄 빌드 (`goslide build`)**:
   * `goslide build my-talk.md -f html,pdf,pptx -o dist/`
   * 산출물 3종(`dist/my-talk.html`, `dist/my-talk.pdf`, `dist/my-talk.pptx`) 생성 확인.

---

### 3.2 운영체제별 특화 검증 시나리오

#### 🐧 1) Linux 환경 (Ubuntu / Debian / WSL2)
* **단일 바이너리 무의존성 검증**:
  * `ldd bin/goslide` 실행 시 `not a dynamic executable` (Statically linked) 확인.
* **헤드리스 크롬 연동**:
  * 시스템에 설치된 Google Chrome 또는 Chromium을 자동 감지하여 PDF가 정상 벡터로 렌더링되는지 확인.
* **웹 서버 모드**:
  * `goslide serve my-talk.md -p 8080` 실행 후 브라우저 핫리로드 동작 확인.

#### 🍏 2) macOS 환경 (Apple Silicon M-시리즈 / Intel)
* **Apple Silicon 네이티브 실행**:
  * `goslide-darwin-arm64`가 Rosetta 에뮬레이션 없이 네이티브 Mach-O 64-bit ARM 바이너리로 구동되는지 확인.
* **Apple Keynote 호환성**:
  * 생성된 `my-talk.pptx`를 Mac의 **Apple Keynote** 앱에서 열었을 때:
    * 16:9 슬라이드 비율 유지 여부
    * 한글 폰트(San Francisco, Apple SD 산돌고딕) 텍스트 줄바꿈 자연스러움 여부
    * 배경 그라데이션 및 코드 블록 가독성 점검
* **macOS 단축키 충돌 여부**:
  * 브라우저 프레젠테이션 모드에서 `Cmd + F`(전체화면), `P`(발표자 콘솔), `B`(블랙아웃) 동작 확인.

#### 🪟 3) Windows 환경 (Windows 10 / 11 64-bit)
* **실행 파일(`goslide.exe`) 및 보안 경고 점검**:
  * Windows Defender / SmartScreen 실행 시 차단 여부 및 경고 문구 확인 (코드 사이닝 미적용 시 "추가 정보 ➔ 실행" 가이드 제공 여부).
* **Microsoft PowerPoint (MS Office 365) 정식 호환성**:
  * 생성된 `my-talk.pptx`를 실제 **PowerPoint 정품 데스크톱 앱**에서 열었을 때:
    * "파일이 손상되었습니다" 경고창 없이 즉시 열리는지 확인.
    * 텍스트 박스, 제목, 본문, 코드 영역이 개별 도형(Shape)으로 정상 분리되어 직접 편집(Edit) 가능한지 확인.
* **Windows 파일 경로 및 URL 정규화 검증**:
  * `C:\Users\...` 드라이브 문자 경로 및 역슬래시(`\`)가 포함된 마크다운에서 배경 이미지/로컬 에셋이 깨지지 않고 렌더링되는지 확인.
* **Chrome DevTools 연동**:
  * `file:///C:/...` 형식으로 헤드리스 크롬이 에러 없이 인쇄/캡처를 완료하는지 확인.

---

### 3.3 프레젠테이션 UX & 인터랙션 검증 (인간 QA)

브라우저에서 `my-talk.html`을 열고 발표자가 실제로 슬라이드를 넘기며 테스트합니다:
1. **듀얼 모니터 & 발표자 콘솔 (Presenter Console)**:
   * 'P' 키 입력 시 발표자 콘솔 팝업이 브라우저 팝업 차단에 걸리지 않고 원활히 열리는지 확인.
   * 메인 화면과 발표자 콘솔 간 슬라이드 싱크(타이머, 현재/다음 슬라이드, 발표자 노트) 실시간 연동 확인.
2. **판서 및 레이저 포인터**:
   * 'L' 키 (레이저 포인터), 'W' 키 (화이트보드 드로잉), 'B' 키 (블랙아웃) 입력 시 마우스 궤적 반응성 점검.
3. **오프라인 동작 (Self-contained)**:
   * 인터넷 연결을 끊은 오프라인 상태에서도 HTML 슬라이드가 CSS/JS 깨짐 없이 100% 정상 작동하는지 확인.

---

## 4. 최종 정식 릴리즈 완료 기준 (Definition of Done - DoD)

1. **RC 단계 통과**:
   * `v1.0.0-rc.1` 배포 후 3대 OS 실기기 E2E 점검에서 블로커(Blocker)급 결함 0건.
2. **PowerPoint 및 Keynote 상호 운용성 확인**:
   * 실제 Windows PowerPoint 및 Mac Keynote에서 PPTX가 정상 로드 및 편집 가능함을 육안으로 확인.
3. **바이너리 배포 무결성**:
   * GitHub Releases 페이지에 Linux, macOS, Windows 5종 바이너리와 압축 아카이브, `checksums.txt`가 정상 업로드됨을 확인.
4. **문서 및 공지 완료**:
   * `docs/release-notes-v1.0.0.md` 릴리즈 노트 작성 완료.
   * `README.md`의 설치 가이드(`curl`, `brew`, `download`) 최신화.

---

## 5. 단계별 실행 일정

| 단계 | 작업 내용 | 담당자 / 환경 | 상태 |
| :---: | :--- | :--- | :---: |
| **Step 1** | **GOS-12 크로스 컴파일 파이프라인 완성** (Makefile, .goreleaser.yaml, CI) | 개발팀 | 진행 예정 |
| **Step 2** | **v1.0.0-rc.1 태그 발행** 및 GitHub Actions 자동 빌드 산출물 확인 | 배포 담당 | 대기 |
| **Step 3** | **Linux / macOS / Windows 3대 OS 실기기 스모크 E2E 테스트 수행** | QA / 개발팀 | 대기 |
| **Step 4** | **결함 수정 (필요 시 v1.0.0-rc.2) 및 최종 DoD 검증** | 개발팀 | 대기 |
| **Step 5** | **공식 v1.0.0 태그 발행** 및 릴리즈 노트 공지 | 배포 담당 | 대기 |
