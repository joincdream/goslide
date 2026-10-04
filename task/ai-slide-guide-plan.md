# [TASK] AI 슬라이드 생성을 위한 가이드 및 DSL 스펙 문서 구축 계획서

> **문서 상태**: 계획 수립 (Draft)  
> **작성일**: 2026-10-04  
> **담당자**: Goslide Core Team  
> **목적**: LLM(ChatGPT, Claude, Gemini, 에이전트 등)이 Goslide 전용 마크다운 슬라이드를 정확하고 일관되게 생성할 수 있도록 단일 진실 공급원(SSOT) 참조 문서(`example-dsl.md`)를 구축하고 배포하는 세부 계획  
> **참조 문서**: [functional_specification.md](file:///home/yundream/myjob/cloit/Goslide/docs/functional_specification.md), [development_roadmap.md](file:///home/yundream/myjob/cloit/Goslide/docs/development_roadmap.md), [README.ko.md](file:///home/yundream/myjob/cloit/Goslide/README.ko.md)

---

## 1. 개요 및 추진 배경 (Overview & Motivation)

### 1.1 배경 및 문제점
* **LLM의 슬라이드 문법 무지**: LLM은 일반 마크다운(블로그/문서)에는 익숙하지만, Goslide의 고유 지시어(`_layout`, `_backgroundImage`, `<!-- split -->`, `<!-- note: -->`) 및 16:9 슬라이드 화면 비율 제약을 알지 못합니다.
* **과도한 설명 프롬프트의 비효율**: 매번 사용자가 긴 자연어 프롬프트를 작성하면 입력 토큰 비용이 낭비되고 사용성이 저하됩니다.
* **사용자 친화적 원클릭 참조 요구**: 사용자가 복잡한 설정 없이 **GitHub Raw URL 하나만 LLM에 전달**하여 최적의 마크다운 슬라이드를 얻을 수 있는 표준화된 가이드/스펙 문서가 필요합니다.

### 1.2 핵심 목표
1. **단일 진실 공급원(SSOT) 문서 구축**: `examples/example-dsl.md` 파일 하나에 **핵심 DSL 스펙 + 토큰 최적화된 Few-Shot 예제**를 압축 통합.
2. **제로-프릭션(Zero-Friction) 사용자 경험**:  
   `"https://raw.githubusercontent.com/.../example-dsl.md 를 참고해서 [주제] 슬라이드 만들어줘"` 한 줄로 완벽한 슬라이드 출력 보장.
3. **결정론적 렌더링 품질 확보**: 뷰포트(16:9) 높이 예산 및 컴포넌트 격리 규칙을 포함하여 생성 즉시 `goslide build` 및 `serve`에서 깨짐 없이 렌더링되도록 강제.

---

## 2. 핵심 산출물 및 배포 구조 (Deliverables)

```
Goslide/
├── examples/
│   └── example-dsl.md            # [핵심 산출물 1] AI 참조용 통합 DSL 스펙 & 골든 데모 문서
├── docs/
│   └── ai-slide-generation.md    # [핵심 산출물 2] 공식 문서용 상세 AI 활용 가이드
├── README.md / README.ko.md      # [핵심 산출물 3] 사용자용 원클릭 프롬프트 안내 섹션 추가
└── task/
    └── ai-slide-guide-plan.md    # [본 계획서]
```

---

## 3. `example-dsl.md` 문서 세부 구성 명세 (Specification)

문서는 LLM의 컨텍스트 윈도우와 토큰 비용을 최소화(1.5KB~2.5KB 내외)하면서도 패턴 매칭을 유도할 수 있도록 3개 섹션으로 구성합니다.

### 3.1 섹션 1: 엔진 불변 규칙 및 포맷팅 제약 (Format Constraints)
* **슬라이드 분할**: 반드시 단독 줄의 `---` 사용 (코드 블록 내부 `---` 사용 엄격 금지).
* **16:9 뷰포트 높이 예산 (Vertical Budget)**:
  * 슬라이드 1장당 최대 10~12줄 이내.
  * 서술형 줄글 금지, 키워드 중심의 단문 불릿(`- **키워드**: 설명`) 원칙.
* **발표자 노트 격리**: 본문 텍스트 슬림화를 위해 부연 설명 및 스크립트는 `<!-- note: ... -->`로 완전 격리.

### 3.2 섹션 2: DSL 지시어 및 컴포넌트 카탈로그 (Cheat Sheet)
| 범주 | 지시어 / 태그 | 설명 |
| :--- | :--- | :--- |
| **전역 설정** | `title`, `author`, `theme`(clean/default/dark), `size`("16:9"), `paginate` | YAML Frontmatter 블록 |
| **레이아웃** | `<!-- _layout: cover -->` | 표지 슬라이드 (중앙 정렬, 타이틀/부제목/작성자) |
| | `<!-- _layout: section -->` | 중간 챕터 전환 간지 (대형 폰트 집중) |
| | `<!-- _layout: two-cols -->` | 2단 비교 레이아웃 (**반드시 `<!-- split -->` 구분자 필수**) |
| | `<!-- _layout: blank -->` | 여백 없는 풀스크린 다이어그램 / 이미지 |
| | `<!-- _layout: default -->` | 일반 본문 슬라이드 |
| **배경/스타일** | `<!-- _backgroundImage: url('...') -->` | 배경 이미지 삽입 |
| | `<!-- _backgroundDim: 0.5 -->` | 텍스트 가독성을 위한 배경 어둡게 처리 (0.1 ~ 0.9) |
| | `<!-- _backgroundColor: #1e293b -->` | 슬라이드 배경색 지정 |
| | `<!-- _color: white -->` | 텍스트 기본 색상 지정 |
| | `<!-- _class: lead -->` / `invert` | 강조 또는 색상 반전 클래스 |
| **로컬 오버라이드** | `_` 접두사 규칙 | `_`가 붙은 지시어는 해당 슬라이드 1회만 적용됨 |
| **발표자 노트** | `<!-- note: ... -->` | 슬라이드 화면에 노출되지 않는 발표자 전용 스크립트 |

### 3.3 섹션 3: 골든 샘플 데크 (Few-shot Examples)
LLM이 형식을 그대로 복제할 수 있도록, 지원하는 모든 레이아웃과 지시어가 1장씩 포함된 4~5장 규모의 실제 마크다운 데크를 문서 하단에 제공.

---

## 4. 단계별 실행 계획 (Action Items & Milestones)

| 단계 | 작업 내용 | 세부 액션 |
| :---: | :--- | :--- |
| **Phase 1** | **`examples/example-dsl.md` 초안 작성** | • Goslide 파서(`internal/parser`) 스펙과 100% 일치하는 지시어 및 규칙 정리<br>• 표지, 2단 분할(`<!-- split -->`), 배경 이미지, 발표자 노트가 포함된 완벽한 골든 데모 슬라이드 작성<br>• 토큰 효율을 위한 군더더기 텍스트 정제 |
| **Phase 2** | **다중 LLM 교차 검증 (Quality Verification)** | • ChatGPT(4o), Claude(3.5 Sonnet), Gemini를 대상으로 Raw URL 링크 기반 프롬프트 테스트 수행<br>• 생성된 마크다운을 `goslide build` 및 `goslide serve`로 실제 렌더링하여 파싱 오류 및 16:9 넘침 점검 |
| **Phase 3** | **사용자 문서화 및 안내 추가** | • `README.md` 및 `README.ko.md`에 "LLM으로 슬라이드 생성하기" 가이드 섹션 추가<br>• 사용자가 복사해서 쓸 수 있는 표준 프롬프트 템플릿(Raw URL 포함) 게시 |
| **Phase 4 (선택)** | **CLI 편의 기능 연동** | • `goslide init` 실행 시 `example-dsl.md`를 로컬에 생성해 주는 옵션 연동<br>• (필요 시) `.agents/skills/goslide` 스킬 파일 추가 배포 검토 |

---

## 5. 완료 기준 (Definition of Done)

1. `examples/example-dsl.md` 파일이 레포지토리에 커밋되어 Raw URL로 외부 접근이 가능해야 함.
2. LLM에게 해당 URL을 전달했을 때:
   - `---` 슬라이드 구분이 완벽히 동작해야 함.
   - `two-cols` 슬라이드에 `<!-- split -->`이 누락되지 않아야 함.
   - 표지와 간지에 적절한 `_layout` 지시어가 선언되어야 함.
   - 본문 텍스트가 16:9 화면을 초과하지 않고 발표자 노트(`<!-- note: -->`)로 적절히 분리되어야 함.
3. 생성된 마크다운이 `goslide build` 명령어에서 오류 없이 빌드되어야 함.
4. `README.md` 및 `README.ko.md`에 사용자용 원클릭 프롬프트 가이드가 반영되어야 함.
