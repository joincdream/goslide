# [GOS-30] 리스트 내 !-- pause -- 사용 시 단일 ul 그룹 유지 및 class="fragment" 정규화 계획서

> **티켓 번호**: [GOS-30](https://joincdream.atlassian.net/browse/GOS-30)  
> **마일스톤**: v1.0.0 정식 릴리즈 안정화 (Parser & Incremental Reveal Polish)  
> **마감일**: 2026-11-15  
> **상태**: 완료 (Done)  
> **담당 패키지**: [`internal/parser/`](../internal/parser/)  
> **연관 문서**:  
> - [`docs/dsl-guide.md`](../docs/dsl-guide.md) *(Section 4: Incremental Reveal & Animation)*  
> - [`docs/theme-guide.md`](../docs/theme-guide.md) *(Section 4: Component Selectors Specification)*  
> - [`docs/okf/decisions-simplification.md`](../docs/okf/decisions-simplification.md) *(ADR-001: 단순성 및 무상태 변환 원칙)*  
> **연관 소스 파일**:  
> - [`internal/parser/fragment.go`](../internal/parser/fragment.go)  
> - [`internal/parser/fragment_test.go`](../internal/parser/fragment_test.go)  
> - [`internal/parser/postprocess.go`](../internal/parser/postprocess.go)  

---

## 1. 개요 및 배경 (Context & Problem Statement)

Goslide는 발표자가 슬라이드 내용을 순차적으로 노출할 수 있도록 마크다운 내 `<!-- pause -->` 주석 지시어를 지원합니다.

그러나 순서 없는 목록(`*` 또는 `-`)의 항목 사이에 `<!-- pause -->`를 삽입할 경우, 기반 마크다운 파서인 CommonMark(Goldmark) 표준 특성상 리스트가 종료된 것으로 판단하여 태그를 닫고 새로 여는 현상이 발생합니다.

### 1.1 현재 생성되는 문제의 HTML 구조
```html
<!-- 현재 파서의 출력: 3개의 분리된 ul 생성 -->
<ul>
  <li><strong>Deterministic Verification</strong> — ...</li>
</ul>
<ul class="fragment" data-fragment-index="1">
  <li><strong>Enterprise Private Serving</strong> — ...</li>
</ul>
<ul class="fragment" data-fragment-index="2">
  <li><strong>Real-Time Observability</strong> — ...</li>
</ul>
```

### 1.2 시각적/구조적 결함
1. **중앙 정렬 레이아웃에서의 계단식 들여쓰기 왜곡**:  
   `section` 레이아웃 등 부모가 `text-align: center`이고 자식이 `ul { display: inline-block; }`인 경우, 분할된 각 `ul`의 너비(width)가 텍스트 길이에 맞춰 수축(Shrink-to-fit)됩니다. 이로 인해 3개의 `ul`이 개별 중앙 정렬되면서 **왼쪽 불릿 기호의 시작 위치가 지그재그(계단식 들여쓰기)로 어긋나는 치명적 시각 버그**가 발생합니다.
2. **시맨틱 무결성 훼손**:  
   논리적으로 1개의 목록 그룹이어야 하는 항목들이 브라우저 및 접근성 도구에서 3개의 별개 목록으로 인식됩니다.
3. **상하 마진 누적**:  
   `ul` 간 기본 마진(`margin-bottom`)이 중첩되어 항목 사이 간격이 비정상적으로 벌어집니다.

---

## 2. 목표 아키텍처 및 정규화 설계 (Target Architecture)

Reveal.js, Slidev, Marp 등 업계 표준 웹 프레젠테이션 엔진과 일치하도록, 점진적 노출 효과는 **`<ul>` 컨테이너가 아닌 개별 `<li>` 태그에 직접 부여**되어야 합니다.

### 2.1 목표 정규화 HTML (Target Semantic State)
```html
<!-- ✅ 목표: 단일 ul 컨테이너 내 li 단위 fragment 적용 -->
<ul>
  <li><strong>Deterministic Verification</strong> — ...</li>
  <li class="fragment" data-fragment-index="1"><strong>Enterprise Private Serving</strong> — ...</li>
  <li class="fragment" data-fragment-index="2"><strong>Real-Time Observability</strong> — ...</li>
</ul>
```

### 2.2 무상태(Stateless) 리스트 병합 원칙 (ADR-001 준수)
복잡한 상태 머신(State Machine)이나 런타임 추적 없이, 파서 후처리 단계([`internal/parser/fragment.go`](../internal/parser/fragment.go))에서 순수 함수(Pure Function) 기반으로 닫힌 `</ul>`와 바로 이어지는 `<ul class="fragment">`의 껍데기만 안전하게 벗겨내는 **무상태 DOM/HTML 정규화**를 수행합니다.

```
[ 파서 원시 HTML ]
  <ul><li>Item 1</li></ul>
  <ul class="fragment" data-fragment-index="1"><li>Item 2</li></ul>
           │
           ▼ (Stateless List Merger Regex / Parser)
[ 정규화된 최종 HTML ]
  <ul>
    <li>Item 1</li>
    <li class="fragment" data-fragment-index="1">Item 2</li>
  </ul>
```

---

## 3. 단계별 상세 구현 계획 (Implementation Steps)

|     단계     | 작업 내용                                                                                                                                                                                                                                                                                                                       | 타겟 소스 파일                                                                  |
| :--------: | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------ |
| **Step 1** | **인접 리스트 병합 및 `<li>` fragment 이전 로직 구현**<br>- `internal/parser/fragment.go`에 `normalizeListFragments(html string) string` 함수 구현<br>- `</ul>\s*<ul([^>]*)>\s*<li` 패턴 감지 시 앞선 `ul`로 병합<br>- 부모 `ul`에 있던 `class="...fragment..."` 및 `data-fragment-index="..."` 속성을 자식 `<li`에 안전하게 병합·주입<br>- `ol`(순서 있는 목록)에 대해서도 동일 병합 규칙 적용 | [`internal/parser/fragment.go`](../internal/parser/fragment.go)           |
| **Step 2** | **파이프라인 연동**<br>- `postProcessHTML(html)`의 `transformFragments(html)` 호출 직후 `normalizeListFragments(html)` 파이프라인 연결                                                                                                                                                                                                         | [`internal/parser/postprocess.go`](../internal/parser/postprocess.go)     |
| **Step 3** | **단위 테스트(Unit Tests) 및 엣지 케이스 검증**<br>- `fragment_test.go`에 테이블 기반 테스트 추가:<br>  1. `ul` 내 3단 pause 항목 병합 및 li fragment 속성 검증<br>  2. `ol` 순서 목록의 pause 병합 및 번호 연속성 유지 검증<br>  3. 리스트 외 일반 단락(`<p>`), 헤딩, 코드 블록 간 pause 동작 보존 검증<br>  4. 중첩 서브리스트(`li > ul`)가 포함된 복합 리스트 내 pause 회귀 검증                                       | [`internal/parser/fragment_test.go`](../internal/parser/fragment_test.go) |
| **Step 4** | **실제 데모 슬라이드 검증**<br>- `examples/demo/demo.md` 16번 슬라이드(`<!-- _layout: section -->`) 렌더링 확인<br>- 불릿 시작선 정렬(Left Alignment) 및 스텝별 Reveal 동작 확인                                                                                                                                                                               | [`examples/demo/demo.md`](../examples/demo/demo.md)                       |
| **Step 5** | **사양 문서 동기화**<br>- `docs/dsl-guide.md`의 Incremental Reveal 항목에 리스트 아이템 단위 pause 표준 동작 명시                                                                                                                                                                                                                                    | [`docs/dsl-guide.md`](../docs/dsl-guide.md)                               |

---

## 4. 수용 기준 및 검증 체크리스트 (Acceptance Criteria)

- [x] **단일 `ul`/`ol` 보장**: 마크다운 리스트 사이에 `<!-- pause -->`를 사용했을 때 결과 HTML이 단 1개의 `<ul>` 또는 `<ol>`로 묶여야 함.
- [x] **`<li class="fragment">` 속성 이전**: 각 `li` 요소에 `class="fragment"` 및 순차적인 `data-fragment-index="1"`, `"2"` 속성이 정확히 주입되어야 함.
- [x] **시각적 정렬 일치**: `section`, `cover`, `lead`, `two-cols`를 포함한 모든 레이아웃에서 불릿 기호의 왼쪽 시작선이 오차 없이 일직선으로 정렬되어야 함.
- [x] **상하 마진 정상화**: 분할된 `ul` 사이의 불필요한 마진이 사라지고 리스트 아이템 기본 간격(`li { margin-bottom: 0.5rem; }`)이 일관되게 유지되어야 함.
- [x] **하위 호환성 유지**: 단락(`<p>`), 인용구(`<blockquote>`), 코드 블록(`pre`) 간의 기존 `<!-- pause -->` 동작이 100% 정상 작동해야 함.
- [x] **단위 테스트 통과**: `go test -v -race ./internal/parser` 실행 시 모든 테스트 케이스를 통과해야 함.
