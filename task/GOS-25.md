# [GOS-25] 프로젝트 디렉토리 구조 정리 및 유지보수성 개선

> **티켓 번호**: [GOS-25](https://joincdream.atlassian.net/browse/GOS-25)  
> **마일스톤**: 프로젝트 유지보수성 및 개발 환경 최적화  
> **상태**: 구현 완료 (Implemented)  
> **담당 패키지/영역**: 루트(`.`), `task/`, `testdata/`, `examples/`, `docs/`, `internal/`  
> **상세 계획서**: [`task/directory-restructuring-plan.md`](file:///home/yundream/myjob/cloit/Goslide/task/directory-restructuring-plan.md)  
> **핵심 가치**: **클린 워크스페이스**, **테스트 격리**, **직관적 온보딩**, **문서 위계화**

---

## 1. 개요 및 배경 (Overview & Problem Statement)

Goslide는 Go 표준 패키지 레이아웃(`cmd/`, `internal/`, `pkg/`)과 단방향 파이프라인 아키텍처를 충실히 구현하고 있으나, 기능 구현 및 테스트 과정에서 산출물 방치, 디렉토리 간 역할 중복(`testdata` vs `examples`), 티켓 문서의 평면적 누적으로 인한 관리 비용 증가가 확인되었습니다.

장기적인 유지보수성 확보와 신규 기여자의 명확한 온보딩을 위해 점진적인 디렉토리 구조 개편을 추진했습니다.

---

## 2. 4단계 점진적 실행 계획 (4-Phase Action Plan)

| 단계 | 목표 | 주요 작업 내용 |
| :--- | :--- | :--- |
| **Phase 1: 루트 청결화 및 산출물 격리** | 저장소 오염 제거 및 실수 커밋 방지 | • 루트의 `demo-standalone.html`, `error.txt`, `.error.txt.swp` 삭제<br>• git tracked 산출물(`testdata/demo.html`) 캐시 제거 (`git rm`)<br>• [`.gitignore`](file:///home/yundream/myjob/cloit/Goslide/.gitignore) 산출물 패턴 보강<br>• [`Makefile`](file:///home/yundream/myjob/cloit/Goslide/Makefile) `clean` 타겟 확장 |
| **Phase 2: 예제와 테스트 데이터 분리** | 공식 데모와 순수 테스트 픽스처 격리 | • `testdata/demo.md`, `demo.ko.md`, 이미지 에셋을 `examples/demo/`로 이관<br>• `examples/example-dsl*.md`를 `examples/dsl/`로 정리<br>• `README.md`, `README.ko.md`의 가이드 명령어 경로 갱신<br>• `testdata/`는 `golden/`, `slides/` 전용으로 순수화 |
| **Phase 3: task/ 티켓 아카이빙** | 활성 작업 시인성 확보 | • `task/archive/` 디렉토리 생성<br>• 완료된 과거 티켓 17건 아카이빙<br>• 현재 활성 티켓(`GOS-12`, `18`, `24`, `25`) 및 계획서만 루트 유지 |
| **Phase 4: 기술 문서군 위계화** | 문서 표준(arc42/Diátaxis) 정렬 | • 기획/평가 문서를 `docs/planning/`, `docs/reports/`로 분류<br>• `docs/okf/index.md` 및 `docs/architecture/index.md` 링크 동기화<br>• `docs/architecture/03-package-structure-and-c4.md`에 `browser`, `i18n` 공식 반영 |

---

## 3. 검증 및 수용 기준 (Acceptance Criteria)

- [x] **Phase 1 완료**: 루트 임시 산출물 정리, `.gitignore` 및 `Makefile clean` 보강
- [x] **Phase 2 완료**: `examples/demo/` 및 `examples/dsl/` 분리, `testdata/` 순수화, README 경로 갱신
- [x] **Phase 3 완료**: `task/archive/` 생성 및 완료 티켓 17건 아카이빙 완료
- [x] **Phase 4 완료**: `docs/planning/` 및 `docs/reports/` 위계화, 아키텍처 문서 최신화
- [x] `git status` 상에 빌드 산출물 및 임시 파일이 나타나지 않는 클린 상태 달성
