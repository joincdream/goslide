# [GOS-34] Frontend(웹 슬라이드/판서/프리젠테이션 뷰) i18n 다국어 지원 연동 계획서

> **티켓 번호**: [GOS-34](https://joincdream.atlassian.net/browse/GOS-34)  
> **마일스톤**: v1.0.0 Frontend i18n & Global UX Polish  
> **마감일**: 2026-11-20  
> **상태**: 진행 중 (In Progress)  
> **담당 패키지**:  
> - [`web/src/`](file:///home/yundream/myjob/cloit/Goslide/web/src/) (Svelte 5 런타임 UI 컴포넌트 & Micro i18n 스토어)  
> - [`internal/renderer/html/`](file:///home/yundream/myjob/cloit/Goslide/internal/renderer/html/) (HTML 마스터 템플릿 및 번역 주입)  
> - [`internal/i18n/`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/) (한/영 다국어 카탈로그)  
> **연관 문서**:  
> - [`task/GOS-31.md`](file:///home/yundream/myjob/cloit/Goslide/task/GOS-31.md) *(온보딩 번들 및 CLI i18n 완료)*  
> - [`docs/okf/screencast-annotation.md`](file:///home/yundream/myjob/cloit/Goslide/docs/okf/screencast-annotation.md) *(스크린캐스트 판서 규격)*  

---

## 1. 개요 및 배경 (Context & Objective)

### 1.1 배경 및 문제의식
GOS-31을 통해 CLI 명령어(`demo`, `prompt`, `theme`, `serve`, `build`, `init`) 및 도움말, 온보딩 가이드(`QUICKSTART`, `README`)의 다국어(한/영) 시스템은 성공적으로 구축되었습니다.  
그러나 **브라우저에서 실행되는 웹 프론트엔드 UI(Svelte 5 런타임)**는 여전히 한국어가 강하게 하드코딩되어 있어, 영문 환경(`--lang en`)으로 슬라이드를 빌드하거나 외국 사용자가 프레젠테이션을 진행할 때 심각한 언어 불일치와 사용성 저해를 유발합니다:

1. **판서 및 레이저 포인터 툴바 ([`PresenterToolbar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterToolbar.svelte))**:
   - `"펜"`, `"포인터"`, `"지우기"`, `"얇게"`, `"보통"`, `"굵게"` 버튼 레이블 및 단축키 툴팁 전체가 한국어로 하드코딩됨.
2. **듀얼 스크린 발표자 콘솔 ([`PresenterSidebar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterSidebar.svelte))**:
   - 팝아웃 창 렌더링 HTML 내 `"현재 슬라이드 (Live)"`, `"다음 슬라이드 (Next)"`, `"발표자 대본 / 메모 (Notes)"`, `"마지막 장"`, `"다음 슬라이드가 없습니다."`, `"이전 (←)"`, `"다음 (→)"`, 글자 확대/축소 툴팁 등이 한국어로 고정됨.
3. **슬라이드 개요 오버레이 ([`OverviewGrid.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/OverviewGrid.svelte))**:
   - `"슬라이드 개요 (Overview)"`, `"단축키 ESC 또는 O 발표 복귀"`, `"닫기 (ESC)"` 텍스트 고정.
4. **HTML 마스터 템플릿 ([`template.go`](file:///home/yundream/myjob/cloit/Goslide/internal/renderer/html/template.go#L37))**:
   - 마크다운 Frontmatter나 CLI `--lang` 설정과 무관하게 `<html lang="ko">`로 정적 하드코딩됨.

### 1.2 개발 목표 (Objectives)
* **완벽한 UI 언어 동기화**: CLI 플래그(`--lang ko|en`), Frontmatter(`lang: en`), 또는 시스템 로케일에 맞춰 브라우저 상의 모든 UI 문자열이 일관되게 다국어로 전환.
* **단일 바이너리 및 경량성 유지 (Zero Dependency)**: 외부 무거운 npm 라이브러리 없이, Svelte 5 Runes 기반 Micro i18n 스토어를 구축하여 런타임 오버헤드 0KB 유지.
* **팝아웃 윈도우 동기화 보장**: `BroadcastChannel`을 통해 메인 창과 분리된 발표자 콘솔 팝아웃 창에도 동일한 로케일 번역이 완벽히 전달되도록 설계.

---

## 2. 아키텍처 의사결정 (ADR: Svelte 5 Micro i18n 채택)

### 2.1 외부 라이브러리 대비 직접 구현(Micro i18n) 채택 근거

| 비교 항목 | 외부 라이브러리 (`svelte-i18n` 등) | **Svelte 5 Micro i18n (채택)** | 채택 사유 및 기술적 타당성 |
| :--- | :--- | :--- | :--- |
| **npm 추가 의존성** | 1~3개 패키지 추가 | **0개 (Zero Dependency)** | Goslide 프론트엔드는 `svelte` 외 외부 런타임 의존성이 없는 순수 아키텍처 유지 |
| **번들 오버헤드** | 약 15KB ~ 30KB 증가 | **0KB 수준 (순수 코드 ~35줄, 1KB 미만)** | 슬라이드 UI 번역 키는 20~30개 내외의 라벨/툴팁에 불과하여 무거운 파서 불필요 |
| **Svelte 5 호환성** | Svelte 3/4 레거시 Store 래퍼 필요 | **Svelte 5 Runes(`$state`) 네이티브** | Svelte 5의 신규 반응형 프리미티브와 100% 네이티브 통합 |
| **오프라인 단일 HTML** | 비동기 fetch 비활성화 및 우회 설정 필요 | **HTML 인라인 주입으로 즉시 동기 바인딩** | 인터넷이 없는 환경 및 `file:///` 로컬 파일 실행 시 CORS/fetch 보안 에러 원천 차단 |
| **카탈로그 단일화** | 프론트엔드용 JSON 별도 분리 관리 | Go 백엔드의 [`internal/i18n`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/)과 **100% 동일한 JSON 공유** | 백엔드와 프론트엔드의 번역 키 불일치 및 중복 유지보수 비용 제거 |

### 2.2 언어 결정 우선순위 (Locale Resolution Hierarchy)
슬라이드 렌더링 시 적용되는 로케일은 다음 우선순위에 따라 결정론적으로 확정됩니다:

1. **1순위 (CLI Flag)**: `goslide build --lang en` 또는 `goslide serve --lang ko`
2. **2순위 (Frontmatter)**: 마크다운 상단 메타데이터 `lang: en` 또는 `lang: ko`
3. **3순위 (OS Locale)**: 환경 변수(`LC_ALL`, `LC_MESSAGES`, `LANG`) 자동 감지 (`i18n.DetectLocale()`)
4. **4순위 (Fallback)**: 기본값 `"en"`

---

## 3. 시스템 아키텍처 및 데이터 흐름

```
[Markdown Source / CLI]
       │  (CLI --lang 또는 Frontmatter lang 또는 OS Locale)
       ▼
[Go Backend: internal/renderer/html]
  ├─ 1. 유효 Locale 결정 (i18n.Bundle)
  ├─ 2. 마스터 템플릿 <html lang="{{ .Lang }}"> 주입
  └─ 3. UI 카탈로그(ui.*) JSON 추출 후 <script> 인라인 주입
       │
       ▼
[Browser DOM Mount: Self-contained HTML]
  └─ <script id="goslide-i18n-data">
       window.__GOSLIDE_I18N__ = {
         "locale": "ko",
         "messages": { "ui.toolbar.pen": "펜", ... }
       };
     </script>
       │
       ▼
[Svelte 5: web/src/stores/i18n.svelte.js]
  └─ let messages = $state(window.__GOSLIDE_I18N__.messages)
  └─ export function t(key, fallback, ...args)
       │
       ├─────────────────────┼─────────────────────┐
       ▼                     ▼                     ▼
[PresenterToolbar]   [PresenterSidebar]     [OverviewGrid]
  • 펜/포인터/지우기     • 듀얼 모니터 콘솔       • 슬라이드 개요 그리드
  • 굵기/단축키 툴팁     • 대본/타이머/이동 버튼    • 단축키 안내/닫기 버튼
```

---

## 4. 세부 구현 명세 (Detailed Specifications)

### 4.1 i18n 카탈로그 키 확장 ([`internal/i18n/locales/`](file:///home/yundream/myjob/cloit/Goslide/internal/i18n/locales/))
`ko.json` 및 `en.json`에 프론트엔드 UI 전용 키 네임스페이스(`ui.*`)를 정의합니다:

```json
{
  "ui.toolbar.pen": "펜",
  "ui.toolbar.pen_tip": "펜 판서 모드 토글 (D)",
  "ui.toolbar.pointer": "포인터",
  "ui.toolbar.pointer_tip": "레이저 포인터 토글 (L)",
  "ui.toolbar.clear": "지우기",
  "ui.toolbar.clear_tip": "현재 슬라이드 판서 지우기 (C)",
  "ui.toolbar.colors_tip": "색상 선택 (단축키: 1~4)",
  "ui.toolbar.widths_tip": "굵기 선택 (단축키: - / +)",
  "ui.toolbar.width_thin": "얇게",
  "ui.toolbar.width_medium": "보통",
  "ui.toolbar.width_thick": "굵게",

  "ui.presenter.title": "Goslide 발표자 콘솔",
  "ui.presenter.live_slide": "현재 슬라이드 (Live)",
  "ui.presenter.next_slide": "다음 슬라이드 (Next)",
  "ui.presenter.notes": "발표자 대본 / 메모 (Notes)",
  "ui.presenter.no_notes": "작성된 발표자 메모가 없습니다.",
  "ui.presenter.last_slide": "마지막 장",
  "ui.presenter.no_next_slide": "다음 슬라이드가 없습니다.",
  "ui.presenter.prev": "이전 (←)",
  "ui.presenter.next": "다음 (→)",
  "ui.presenter.font_dec": "글자 축소",
  "ui.presenter.font_inc": "글자 확대",

  "ui.overview.title": "슬라이드 개요 (Overview)",
  "ui.overview.slides_count": "%d개 슬라이드",
  "ui.overview.back_hint": "단축키 ESC 또는 O 발표 복귀",
  "ui.overview.close": "닫기 (ESC)"
}
```

---

### 4.2 Go 백엔드 HTML 렌더러 연동 ([`internal/renderer/html/`](file:///home/yundream/myjob/cloit/Goslide/internal/renderer/html/))

1. **`internal/i18n/i18n.go` 헬퍼 메서드 추가**:
   - `GetCatalog(lang string) map[string]string`: 지정된 로케일의 전체 맵 반환 (영문 폴백 포함).
   - `GetCatalogJSON(lang string) template.JS`: 템플릿에 안전하게 인라인 삽입 가능한 JSON 문자열 생성.

2. **마스터 템플릿 필드 확장 ([`template.go`](file:///home/yundream/myjob/cloit/Goslide/internal/renderer/html/template.go))**:
   - `documentTemplateData` 구조체에 `Lang string`, `I18nJSON template.JS` 추가.
   - 마스터 HTML 상단 `<html lang="{{ .Lang }}">` 동적 속성 적용.
   - 바디 하단 `<script id="goslide-i18n-data">window.__GOSLIDE_I18N__ = {{ .I18nJSON }}; window.__GOSLIDE_LANG__ = "{{ .Lang }}";</script>` 주입.

3. **렌더러 로직 확장 ([`renderer.go`](file:///home/yundream/myjob/cloit/Goslide/internal/renderer/html/renderer.go))**:
   - `WithLang(lang string)` 옵션 함수 지원.
   - Frontmatter(`lang`) 또는 CLI 글로벌 플래그(`--lang`), OS 로케일 순으로 로케일 결정 후 템플릿 데이터에 바인딩.

---

### 4.3 Svelte 5 Micro i18n 스토어 설계 ([`web/src/stores/i18n.svelte.js`](file:///home/yundream/myjob/cloit/Goslide/web/src/stores/i18n.svelte.js))

외부 라이브러리 없이 Svelte 5 `$state`를 사용하는 초경량 모듈 구현:

```javascript
// web/src/stores/i18n.svelte.js
const defaultData = (typeof window !== 'undefined' && window.__GOSLIDE_I18N__) || {
  locale: 'en',
  messages: {}
};

let currentLocale = $state(defaultData.locale || 'en');
let catalog = $state(defaultData.messages || {});

/**
 * 번역 키에 해당하는 텍스트를 반환합니다.
 * @param {string} key 번역 키 (예: 'ui.toolbar.pen')
 * @param {string} fallback 기본 폴백 텍스트
 * @param {...any} args 포맷팅 인자 (%s, %d)
 * @returns {string}
 */
export function t(key, fallback = '', ...args) {
  let msg = catalog[key] || fallback || key;
  if (args.length > 0) {
    args.forEach(arg => {
      msg = msg.replace(/%[sd]/, String(arg));
    });
  }
  return msg;
}

export function getLocale() {
  return currentLocale;
}

export function setLocale(locale, newMessages) {
  currentLocale = locale;
  if (newMessages) {
    catalog = newMessages;
  }
}
```

---

### 4.4 Svelte UI 컴포넌트 마이그레이션 ([`web/src/components/`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/))

1. **[`PresenterToolbar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterToolbar.svelte)**:
   - 툴팁 및 버튼 라벨을 `t('ui.toolbar.pen', '펜')`, `t('ui.toolbar.pointer', '포인터')`, `t('ui.toolbar.clear', '지우기')` 등으로 교체.
   - 굵기 라벨 `t('ui.toolbar.width_thin', '얇게')` 등으로 교체.
2. **[`PresenterSidebar.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/PresenterSidebar.svelte)**:
   - 팝아웃 창 렌더링 HTML 생성부(`createPresenterWindow`)의 하드코딩 한글 문자열을 `t('ui.presenter.*')` 호출로 전면 치환.
3. **[`OverviewGrid.svelte`](file:///home/yundream/myjob/cloit/Goslide/web/src/components/OverviewGrid.svelte)**:
   - 개요 헤더 타이틀, 슬라이드 개수 표시(`t('ui.overview.slides_count', '%d개 슬라이드', slides.length)`), 닫기 버튼 안내 치환.

---

### 4.5 빌드 및 패키징 자동화
1. `cd web && npm run build` 실행하여 `web/dist/goslide-core.js` 번들 갱신.
2. `internal/theme/assets/js/goslide-core.js`로 복사 및 Go 바이너리(`//go:embed`)에 임베딩.
3. `make check` 전체 게이트웨이(테스트, 린트, 레이스 디텍터) 통과 확인.

---

## 5. 단계별 실행 계획 (Action Items)

| 단계 | 작업 내용 | 대상 파일 | 검증 기준 |
| :---: | :--- | :--- | :--- |
| **Phase 1** | 한/영 UI 다국어 사전 키 확장 | `internal/i18n/locales/{ko,en}.json` | 누락된 UI 키 전수 정의 |
| **Phase 2** | Go i18n 카탈로그 메서드 추가 및 HTML 렌더러 연동 | `internal/i18n/i18n.go`<br>`internal/renderer/html/{template,renderer}.go` | `html lang` 및 `__GOSLIDE_I18N__` 주입 단위 테스트 통과 |
| **Phase 3** | Svelte 5 Micro i18n 스토어 구현 및 컴포넌트 리팩토링 | `web/src/stores/i18n.svelte.js`<br>`web/src/components/*.svelte` | 하드코딩 한글 0건, `npm run build` 성공 |
| **Phase 4** | 프론트엔드 번들 갱신 및 바이너리 임베딩 | `internal/theme/assets/js/goslide-core.js` | 바이너리 빌드 및 임베딩 완료 |
| **Phase 5** | 통합 검증 및 시각적 회귀 점검 | `cmd/goslide/` 및 브라우저 검증 | `--lang en` 및 `--lang ko` 브라우저 UI 검증 및 `make check` 통과 |

---

## 6. 완료 기준 (Definition of Done)

- [x] `PresenterToolbar`의 펜/레이저/지우기 및 굵기 라벨/툴팁이 지정된 언어(ko/en)로 정확히 출력된다.
- [x] 팝아웃 `PresenterSidebar` 콘솔 창의 슬라이드 타이틀, 발표자 메모 기본 문구, 이전/다음 버튼이 언어 설정에 따라 정확히 출력된다.
- [x] `OverviewGrid`의 슬라이드 개수, 안내문, 닫기 버튼이 언어 설정에 따라 정확히 출력된다.
- [x] 빌드된 HTML의 `<html lang="...">` 속성이 지정된 언어(`ko` 또는 `en`)와 일치한다.
- [x] 외부 npm 런타임 의존성이 0개로 유지되고 번들 크기 증가가 없다.
- [x] `make check` (gofmt, golangci-lint, gocyclo, test-race)가 오류 없이 100% 통과한다.
