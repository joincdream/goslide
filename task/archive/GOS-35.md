# [GOS-35] 사용자 여정 기반 설치 가이드 및 OS별 온보딩 문서 종합 보강 계획서

> **티켓 번호**: [GOS-35](https://joincdream.atlassian.net/browse/GOS-35)  
> **마일스톤**: v1.0.0 Release Documentation & Onboarding Polish  
> **마감일**: 2026-10-31  
> **상태**: 진행 중 (In Progress)  
> **담당 문서**:  
> - [`QUICKSTART.ko.md`](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.ko.md) / [`QUICKSTART.md`](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.md)  
> - [`README.ko.md`](file:///home/yundream/myjob/cloit/Goslide/README.ko.md) / [`README.md`](file:///home/yundream/myjob/cloit/Goslide/README.md)  
> **연관 문서**:  
> - [`task/release/end-to-end-user-journey.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/end-to-end-user-journey.md) *(사용자 여정 및 갭 분석 보고서)*  
> - [`task/release/linux-user-test-checklist.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/linux-user-test-checklist.md) *(리눅스 사용자 테스트 체크리스트)*  
> - [`Makefile`](file:///home/yundream/myjob/cloit/Goslide/Makefile) *(크로스 플랫폼 패키징 사양)*  

---

## 1. 개요 및 배경 (Context & Objective)

### 1.1 배경 및 문제의식
[`task/release/end-to-end-user-journey.md`](file:///home/yundream/myjob/cloit/Goslide/task/release/end-to-end-user-journey.md)를 통해 신규 사용자의 전체 여정(사이트 방문 ➔ 다운로드 ➔ 설치 ➔ 첫 실행)을 전수 조사한 결과, Goslide의 최대 강점인 **"단일 바이너리 & CGO 0% 무의존성"**에도 불구하고, 첫 설치 단계에서 몇 가지 실질적인 마찰 지점(Friction Points)이 확인되었습니다:

1. **OS 및 CPU 아키텍처 선택 혼선**:
   - GitHub Releases에 6종의 아티팩트(`goslide_linux_amd64.tar.gz`, `goslide_darwin_arm64.tar.gz` 등)가 올라오지만, 일반 사용자(특히 Apple Silicon Mac 사용자)는 자신이 `darwin_arm64`를 받아야 하는지 `darwin_amd64`를 받아야 하는지 명확한 안내 표가 없어 혼란을 겪을 수 있습니다.
2. **압축 해제 후 실행 권한(`chmod +x`) 및 PATH 등록 가이드 부재**:
   - 아카이브를 해제한 뒤 다운로드 폴더에서만 `./goslide`로 실행하다가, 다른 디렉토리에서 `goslide: command not found`를 겪는 초보 사용자 마찰이 존재합니다.
   - `cp goslide ~/.local/bin/` 또는 `sudo cp goslide /usr/local/bin/` 등 시스템 어디서나 호출할 수 있는 원라이너(One-liner) 명령어가 누락되어 있습니다.
3. **macOS 게이트키퍼(Gatekeeper) 격리 차단 경고**:
   - 코드 사인이 포함되지 않은 오픈소스 바이너리는 macOS에서 실행 시 `"확인되지 않은 개발자가 배포했기 때문에 열 수 없습니다"`라는 보안 팝업이 뜨며 차단됩니다. 이를 해제하는 `xattr -d com.apple.quarantine goslide` 안내가 필수적입니다.
4. **체크섬 무결성 검증 안내 부재**:
   - `checksums.txt`가 릴리즈에 동봉되지만, 이를 터미널에서 검증하는 표준 명령어(`sha256sum -c`)가 안내되지 않았습니다.

### 1.2 개발 목표 (Objectives)
* **3분 온보딩 완결성 보장**: 신규 사용자가 외부 구글링 없이 온보딩 문서만 보고 3줄의 터미널 명령어로 다운로드부터 `goslide demo` 실행까지 막힘없이 도달하도록 가이드를 완비합니다.
* **한/영 문서 100% 동기화**: 국문(`*.ko.md`)과 영문(`*.md`) 문서 간 설명 및 명령어 사양이 완벽히 일치하도록 정비합니다.
* **불필요한 작업 분산 방지**: 여러 개로 쪼개질 수 있는 문서 보강 작업을 본 티켓(GOS-35) 하나로 통합 처리하여 업무 효율성을 극대화합니다.

---

## 2. 세부 사양 및 반영 내용 (Detailed Specifications)

### 2.1 OS 및 CPU 아키텍처 매핑 가이드 (Download Matrix)
`README` 및 `QUICKSTART`의 다운로드 섹션에 직관적인 아키텍처 매핑 표를 추가합니다:

| 운영체제 (OS) | CPU 아키텍처 / 대상 기기 | 권장 다운로드 아티팩트 |
| :--- | :--- | :--- |
| **macOS** | **Apple Silicon (M1 / M2 / M3 / M4)** | `goslide_darwin_arm64.tar.gz` |
| **macOS** | Intel 프로세서 | `goslide_darwin_amd64.tar.gz` |
| **Linux** | 64-bit Intel / AMD (x86_64) | `goslide_linux_amd64.tar.gz` |
| **Linux** | 64-bit ARM (aarch64 / 라즈베리파이 4, 5) | `goslide_linux_arm64.tar.gz` |
| **Windows** | 64-bit Intel / AMD (x64) | `goslide_windows_amd64.zip` |
| **Windows** | 64-bit ARM (Snapdragon X Elite 등) | `goslide_windows_arm64.zip` |

---

### 2.2 압축 해제 및 시스템 PATH 설치 원라이너 가이드

#### Linux & macOS
```bash
# 1. 아카이브 압축 해제
tar -xzf goslide_linux_amd64.tar.gz

# 2. 실행 권한 부여 (필요 시)
chmod +x goslide

# 3. 사용자 PATH 디렉토리로 이동 (권장: ~/.local/bin 또는 /usr/local/bin)
mkdir -p ~/.local/bin && mv goslide ~/.local/bin/

# 4. 설치 확인
goslide --help
```

#### Windows
- 다운로드한 `.zip` 파일 우클릭 ➔ **압축 풀기**
- 추출된 `goslide.exe`를 사용자 경로(예: `C:\Users\<사용자>\bin`)로 이동 후 시스템 환경 변수 `Path`에 해당 경로 추가

---

### 2.3 OS별 특이사항 및 보안 경고 대처법

#### macOS 미확인 개발자(Gatekeeper) 경고 해제
터미널에서 아래 명령을 실행하여 다운로드 파일의 격리 속성(Quarantine)을 제거합니다:
```bash
xattr -d com.apple.quarantine goslide
```

#### 체크섬(SHA-256) 파일 무결성 검증
다운로드한 아카이브의 위변조 방지를 위해 동봉된 `checksums.txt`로 검증합니다:
```bash
# Linux
sha256sum -c checksums.txt --ignore-missing

# macOS
shasum -a 256 -c checksums.txt --ignore-missing
```

---

## 3. 작업 대상 파일 및 단계별 실행 계획

```
[Phase 1: QUICKSTART 정비] ───► [Phase 2: README 정비] ───► [Phase 3: 링크 & 정합성 검증]
  - QUICKSTART.ko.md              - README.ko.md                - 마크다운 링크 유효성
  - QUICKSTART.md                 - README.md                   - 한/영 문구 일치 검증
```

### Phase 1. `QUICKSTART` 가이드 정비 (최우선 온보딩 문서)
1. [`QUICKSTART.ko.md`](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.ko.md):
   - 문서 서두에 "0단계: 바이너리 설치 및 PATH 설정" 섹션 신설.
   - 아키텍처 매핑 표, `chmod +x`, `mv ~/.local/bin/`, macOS `xattr` 안내 명시.
2. [`QUICKSTART.md`](file:///home/yundream/myjob/cloit/Goslide/QUICKSTART.md):
   - 영문으로 완벽히 동일한 구조 반영.

### Phase 2. `README` 설치 섹션 정비
1. [`README.ko.md`](file:///home/yundream/myjob/cloit/Goslide/README.ko.md):
   - `### 1. 설치 방법 ➔ #### 바이너리 다운로드 (권장)` 부분에 아키텍처 선택표 및 설치 원라이너 명령 추가.
2. [`README.md`](file:///home/yundream/myjob/cloit/Goslide/README.md):
   - 영문 동기화.

### Phase 3. 정합성 및 검증 (Verification)
- 문서 내 마크다운 링크 및 서식 깨짐 여부 정적 점검.
- 실제 터미널 명령어 구문 오류 여부 정적 확인.

---

## 4. 완료 조건 (Definition of Done - DoD)

1. [x] **티켓 상태 전이 완료**: Jira 티켓 [GOS-35](https://joincdream.atlassian.net/browse/GOS-35)가 `진행 중`으로 전환되었는가?
2. [x] **설치 매핑 표 반영**: Linux, macOS(Apple Silicon/Intel), Windows 6종 아티팩트 매핑 표가 `QUICKSTART` 및 `README`에 명시되었는가?
3. [x] **PATH 및 권한 원라이너 제공**: `chmod +x` 및 `mv ~/.local/bin/` 안내가 빠짐없이 기술되었는가?
4. [x] **macOS 게이트키퍼 해제 명령어 명시**: `xattr -d com.apple.quarantine goslide`가 포함되었는가?
5. [x] **다국어 100% 동기화**: `ko.md`와 `.md` 간 내용 불일치나 번역 누락이 없는가?
