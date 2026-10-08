# Goslide External Themes

Goslide 사용자가 직접 디자인 코드를 확인하고, 복사하여 기업 및 개인 브랜드에 맞게 자유롭게 커스터마이징할 수 있는 **공식 외장 테마 템플릿(External Theme Templates)** 디렉토리입니다.

---

## 📁 디렉토리 구성

```text
themes/
├── README.md           # 외장 테마 가이드 및 커스텀 매뉴얼
├── corporate.css       # [공식 외장 테마] 신뢰성 높은 비즈니스/IR/제안서 테마
├── academic.css        # (준비 중) 학술 발표, 연구 논문, 세미나용 테마
└── cyber-dark.css      # (준비 중) 개발자 컨퍼런스, 테크 밋업용 다크 테마
```

---

## 🚀 외장 테마 사용법

바이너리에 내장된 기본 테마 대신 사용자가 작성한 외장 CSS 파일을 적용할 때는 `--theme-path` 플래그를 사용합니다.

### 1. 개발 및 실시간 미리보기 (`serve`)
```bash
goslide serve presentation.md --theme-path themes/corporate.css
```

### 2. 정적 HTML / PDF / PPTX 내보내기 (`build`)
```bash
# HTML 생성
goslide build presentation.md --theme-path themes/corporate.css -o dist/index.html

# PDF / PPTX 일괄 내보내기
goslide build presentation.md --theme-path themes/corporate.css -f all -o dist/
```

---

## 🎨 커스텀 테마 만들기 (3단계)

### 1단계: 기존 외장 테마 복사
`themes/corporate.css`를 원하는 이름으로 복사합니다.
```bash
cp themes/corporate.css themes/my-company.css
```

### 2단계: 브랜드 컬러 및 폰트 변경 (`:root`)
CSS 파일 상단의 `:root` 변수를 회사 브랜드 가이드라인에 맞게 수정합니다.
```css
:root {
  /* 메인 브랜드 컬러 */
  --corporate-primary: #1e3a8a;      /* 브랜드 메인 색상 */
  --corporate-accent: #0284c7;       /* 강조/포인트 색상 */
  
  /* 배경 및 서피스 */
  --corporate-bg-card: #ffffff;      /* 슬라이드 캔버스 배경 */
  --corporate-bg-stage: #f8fafc;     /* 프레젠테이션 전체 배경 */

  /* 폰트 지정 */
  --corporate-font-family: "Pretendard", -apple-system, BlinkMacSystemFont, sans-serif;
}
```

### 3단계: 슬라이드 레이아웃별 오버라이드
Goslide 표준 레이아웃 클래스를 타겟팅하여 각 슬라이드 형태별 디자인을 자유롭게 확장할 수 있습니다:
- `.slide-card.layout-cover`: 타이틀 및 표지 슬라이드
- `.slide-card.layout-section`: 중간 간지(챕터 디바이더) 슬라이드
- `.slide-card.layout-two-cols`: 2단 그리드 분할 슬라이드
- `.slide-card.layout-lead`: 키 메시지 강조 슬라이드
- `.slide-card`: 기본 본문 슬라이드

---

## 📖 고급 테마 제작 & LLM 프롬프트 가이드

Goslide의 4-Tier DOM 트리 계층 구조와 레이아웃별 상세 CSS 셀렉터 명세, 그리고 **LLM(ChatGPT, Claude 등)에 전달하여 새 테마 CSS를 1분 만에 자동 생성할 수 있는 프롬프트 템플릿**은 공식 가이드를 참조하세요:

👉 [Goslide Theme Development Guide (docs/theme-guide.md)](file:///home/yundream/myjob/cloit/Goslide/docs/theme-guide.md)
