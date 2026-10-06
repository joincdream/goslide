# [GOS-17] 단계적 빌드(Incremental Reveal) 및 코드 블록 라인 포커스(Line Highlighting) 개발 계획서

> **티켓 번호**: [GOS-17](https://joincdream.atlassian.net/browse/GOS-17)  
> **마일스톤**: Advanced Dynamics (DYN 6.5 ➔ 9.5 도약)  
> **마감일**: 2026-10-16  
> **상태**: 완료 (Done)  
> **담당 패키지**: `internal/parser/`, `internal/theme/assets/css/`, `web/src/stores/`, `web/src/components/`  
> **참조 문서**: [competitive-analysis-and-strategy.md](file:///home/yundream/myjob/cloit/Goslide/docs/competitive-analysis-and-strategy.md), [contracts-interfaces.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/contracts-interfaces.md), [hard-constraints.md](file:///home/yundream/myjob/cloit/Goslide/docs/okf/hard-constraints.md), [AGENTS.md](file:///home/yundream/myjob/cloit/Goslide/AGENTS.md)

---

## 1. 개요 및 가치 제안 (Overview & Value Proposition)

본 태스크의 목적은 기술 세미나 발표, 컨퍼런스 키노트, 그리고 온라인 강의 녹화 시 청중의 시선을 완벽히 통제하고 코드 설명의 명확성을 극대화하기 위해, Goslide의 유일한 취약점이었던 **DYN(Advanced Dynamics, 기존 6.5점)**의 2대 핵심 기능인 **"단계적 빌드(Incremental Reveal / Fragments)"**와 **"코드 블록 라인 포커스(Code Line Highlighting)"**를 구현하는 것입니다.

### 핵심 달성 목표
1. **청중 시선 제어 (단계적 빌드)**: 슬라이드 내의 불릿 포인트나 단락을 한 번에 노출하지 않고, 발표자의 발표 흐름(스페이스/화살표 키)에 맞추어 하나씩 부드럽게 등장(`<!-- pause -->`).
2. **코드 설명력 극대화 (라인 포커스)**: 긴 소스 코드 블록에서 설명하려는 특정 줄만 환하게 강조하고, 나머지 코드는 30% 투명도로 어둡게(dim) 처리하여 청중이 발표자의 설명에 즉각 집중하도록 유도(````go {2-4,7}````).
3. **정적 출력 결정론 보장 (PDF / PPTX)**: 브라우저 발표 시에는 단계적으로 노출되지만, 비즈니스 보고서로 내보내는 PDF 및 PPTX에서는 모든 내용이 완전히 펼쳐진 최종 상태(Final Reveal)로 인쇄되어 내용 누락 방지.
4. **리더보드 압도적 1위 등극**: Slidev의 유일한 무기였던 다이내믹스를 CGO 0% 무의존 Go 단일 바이너리로 흡수하여, 종합 점수 **57.5점 (100점 환산 95.8점)**으로 전 세계 마크다운 슬라이드 도구 1위 확립.

---

## 2. 아키텍처 및 세부 설계 명세 (Specification)

### 2.1 단계적 빌드 (Incremental Reveal / Fragments)

#### 1) 마크다운 작성 문법
슬라이드 내용 중 발표 도중 멈추고 다음 단계로 노출하고 싶은 지점에 `<!-- pause -->` 태그를 삽입합니다:

```markdown
## High-Performance Architecture

- Pure Go Single Static Binary
<!-- pause -->
- Zero CGO & Zero Node.js Dependencies
<!-- pause -->
- Millisecond Compilation (< 300ms)
```

#### 2) 파서 처리 메커니즘 (`internal/parser/`)
* Goldmark AST 트랜스포머 단계에서 `<!-- pause -->` 주석 블록을 감지합니다.
* `<!-- pause -->` 이후에 등장하는 첫 번째 블록 요소(리스트 항목 `<li>` 또는 단락 `<p>`)에 다음 속성을 주입합니다:
  ```html
  <li class="fragment" data-fragment-index="1">Zero CGO & Zero Node.js Dependencies</li>
  ```
* 동일 슬라이드 내 `<!-- pause -->`가 나타날 때마다 `data-fragment-index`가 1씩 자동 증가합니다.

#### 3) 프론트엔드 상태 머신 제어 (`web/src/stores/deck.svelte.js` & `KeyboardHandler.svelte`)
* 슬라이드별 총 프래그먼트 개수(`totalFragments`)와 현재 노출 인덱스(`currentFragmentIndex`)를 관리합니다.
* **다음 키 입력 (`ArrowRight`, `Space`, `Enter`)**:
  * `currentFragmentIndex < totalFragments`: 다음 슬라이드로 넘어가지 않고, 현재 슬라이드의 다음 프래그먼트에 `.visible` 클래스를 추가하고 `currentFragmentIndex++`.
  * `currentFragmentIndex == totalFragments`: 모든 프래그먼트가 이미 노출되었으므로 다음 슬라이드로 정상 이동.
* **이전 키 입력 (`ArrowLeft`, `Backspace`)**:
  * `currentFragmentIndex > 0`: 현재 프래그먼트의 `.visible` 클래스를 제거하고 `currentFragmentIndex--`.
  * `currentFragmentIndex == 0`: 이전 슬라이드로 이동 (이전 슬라이드의 모든 프래그먼트는 완료 상태로 표시).

#### 4) CSS 전환 효과 및 스타일링 (`internal/theme/assets/css/deck-content.css`)
* **기본값 (Default: 완전 은닉 ➔ 부드러운 등장)**:
  * 처음에는 `opacity: 0`으로 화면에서 완전히 숨겨져 청중의 스포일러를 차단하며, 키 입력 시 위로 6px 떠오르며 스르륵 나타납니다.
* **선택 옵션 (Dim Mode: 흐릿 ➔ 선명)**:
  * 슬라이드에 `<!-- _fragmentStyle: dim -->` 선언 시, 처음부터 25% 반투명으로 자리를 잡고 있다가 순차적으로 100% 선명해지는 딤 모드를 지원합니다.

```css
/* [기본 모드] 아예 안 보이다가 스르륵 등장 (Hidden ➔ Reveal) */
.fragment {
  opacity: 0;
  transform: translateY(6px);
  transition: opacity 0.25s ease-out, transform 0.25s ease-out;
  pointer-events: none;
}

.fragment.visible {
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

/* [선택 딤 모드] 흐릿했다가 자기 차례에 선명해짐 (Dim ➔ Focus) */
.dim-fragments .fragment {
  opacity: 0.25;
  transform: none;
}

.dim-fragments .fragment.visible {
  opacity: 1;
  font-weight: 500;
}
```

---

### 2.2 코드 블록 라인 포커스 (Code Line Highlighting)

#### 1) 마크다운 작성 문법
코드 펜스 언어 식별자 바로 뒤에 중괄호 `{...}`로 강조할 라인 번호 또는 범위를 지정합니다:

````markdown
```go {2,4-6}
package main

import "fmt" // 2번 라인 강조

func main() {
    server := NewServer() // 4번 라인 강조
    server.Configure()    // 5번 라인 강조
    server.Start()        // 6번 라인 강조
}
```
````

#### 2) 라인 파서 문법 규칙
* 단일 행: `{3}`
* 연속 범위: `{2-5}`
* 복합 행: `{1,4-6,9}`
* 잘못된 형식이나 범위를 벗어난 번호는 무시하고 일반 코드 블록으로 안전 폴백(Graceful Fallback).

#### 3) Chroma 구문 강조 연동 및 DOM 구조
* Chroma 렌더러가 코드를 줄 단위 `<span>`으로 감싸 출력할 때, 지정된 하이라이트 행에 `.hl` 클래스를 부여하고 컨테이너에 `.has-highlights`를 부여합니다:
  ```html
  <pre class="highlight-container has-highlights"><code class="language-go">
    <span class="line">package main</span>
    <span class="line hl">import "fmt"</span>
    <span class="line"></span>
    <span class="line hl">    server := NewServer()</span>
    <span class="line hl">    server.Configure()</span>
    <span class="line hl">    server.Start()</span>
  </code></pre>
  ```

#### 4) 시각적 딤(Dimming) 효과 스타일링
```css
/* 강조 라인이 존재하는 코드 블록의 기본 줄은 35%로 어둡게 처리 */
.highlight-container.has-highlights .line {
  opacity: 0.35;
  transition: opacity 0.2s ease, background-color 0.2s ease;
}

/* 포커스된 줄은 100% 선명도 및 은은한 하이라이트 배경색 부여 */
.highlight-container.has-highlights .line.hl {
  opacity: 1.0;
  font-weight: 600;
  background-color: rgba(56, 189, 248, 0.15); /* Sky blue tint */
  border-left: 3px solid #38bdf8;
  padding-left: 6px;
  margin-left: -9px;
}
```

---

### 2.3 엔터프라이즈 정적 출력 보장 (PDF / PPTX Determinism)

* **원칙**: PDF 인쇄 및 파워포인트(PPTX) 슬라이드는 정적 문서이므로, 단계별 중간 상태가 아닌 **"모든 불릿이 표시되고 코드가 온전히 보이는 완전체 상태"**여야 합니다.
* **구현**:
  * Headless Chrome PDF 렌더러(`internal/exporter/pdf/`) 및 PPTX 캡처 시, 인쇄 전용 CSS 미디어 쿼리(`@media print`) 또는 `.print-mode` 클래스를 통해 모든 `.fragment`를 무조건 `opacity: 1 !important; transform: none !important;`로 오버라이드합니다.

---

## 3. 단계별 개발 일정 및 세부 태스크 (Execution Plan)

| 단계 | 작업 내용 | 담당 파일 | 예상 산출물 |
| :---: | :--- | :--- | :--- |
| **Phase 1** | **마크다운 파서 확장**<br>• `<!-- pause -->` 태그 감지 및 `class="fragment"` 속성 주입 AST 확장<br>• 코드 블록 언어 태그 뒤 `{range}` 파서 구현 | `internal/parser/fragment.go`<br>`internal/parser/code_highlight.go`<br>`internal/parser/parser_test.go` | 파서 단위 테스트 통과 |
| **Phase 2** | **프론트엔드 Fragment 상태 머신 구현**<br>• `deck.svelte.js`에 프래그먼트 인덱스 및 스텝 네비게이션 액션 추가<br>• `KeyboardHandler.svelte` 키보드 이벤트 분기 연동 | `web/src/stores/deck.svelte.js`<br>`web/src/components/KeyboardHandler.svelte` | 단계적 노출 키 인터랙션 완성 |
| **Phase 3** | **시각적 스타일링 및 전환 효과 구현**<br>• `.fragment` 부드러운 페이드인 트랜지션 정의<br>• `.has-highlights` 및 `.line.hl` 포커스 딤 스타일링 적용 | `internal/theme/assets/css/deck-content.css`<br>`internal/theme/assets/css/theme-*.css` | 세련된 시각 효과 완성 |
| **Phase 4** | **발표자 콘솔 및 PDF/PPTX 출력 연동**<br>• 발표자 뷰(`PresenterSidebar`)에서 다음 단계 프래그먼트 프리뷰 동기화<br>• PDF/PPTX 내보내기 시 전량 노출 보장 CSS 오버라이드 | `web/src/components/PresenterSidebar.svelte`<br>`internal/exporter/pdf/print.go` | 발표자 뷰 & 익스포터 일치 |
| **Phase 5** | **E2E 통합 검증 & 빌드 테스트**<br>• `examples/example-dsl.md`에 `<!-- pause -->` 및 `{2-4}` 예제 추가<br>• `goslide serve` 실시간 검증 및 `goslide build -f pdf,pptx` 최종 검증 | `examples/example-dsl.md`<br>통합 테스트 스위트 | DoD 100% 검증 완료 |

---

## 4. 완료 기준 (Definition of Done - DoD)

1. **단계적 빌드 동작 검증**:
   * `<!-- pause -->`가 포함된 슬라이드에서 스페이스/화살표 키를 누르면 다음 슬라이드로 넘어가지 않고 항목이 순차적으로 하나씩 페이드인됨을 확인.
2. **코드 라인 포커스 동작 검증**:
   * ````go {2-4}```` 구문이 포함된 코드 블록에서 2~4번 줄만 선명하게 강조되고, 나머지 줄은 은은하게 반투명(dim) 처리됨을 육안 확인.
3. **엔터프라이즈 익스포트 무결성**:
   * `goslide build talk.md -f pdf,pptx` 실행 시 생성된 PDF와 PPTX에 모든 프래그먼트 항목과 코드가 누락 없이 100% 온전하게 인쇄됨을 확인.
4. **품질 검증 통과**:
   * `go test -v -race ./...` 전체 테스트 스위트 통과 (데이터 레이스 0건).
