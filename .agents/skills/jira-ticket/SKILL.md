---
name: jira-ticket
description: >-
  Manage Jira tickets using the Jira CLI. Use this skill when viewing, searching,
  creating, updating, commenting on, or transitioning Jira issues.
---

# Jira Ticket Management Skill

이 스킬은 프로젝트 디렉토리의 [.jira-profile](../../../.jira-profile) 설정을 기반으로 Jira CLI(`jira`)를 사용하여 티켓을 조회, 생성, 수정, 상태 전이 및 코멘트 작업을 수행할 때 적용합니다.

Jira CLI는 작업 디렉토리의 `.jira-profile`을 자동으로 감지하므로 별도의 `--profile` 옵션을 지정할 필요 없이 명령어를 바로 사용합니다.

---

## 1. 티켓 조회 및 맥락 분석

AI 에이전트 작업 시 프롬프트 주입 및 컨텍스트 파악에 최적화된 `--md` 옵션을 적극 활용합니다.

- **전체 티켓 목록 조회**:
  ```bash
  jira list --md
  ```
- **조건부 티켓 목록 조회 (JQL)**:
  ```bash
  jira list "status = '진행 중'" --md
  jira list "assignee = currentUser() AND status != '완료'" --md
  ```
- **단일 티켓 전체 맥락 파악**:
  ```bash
  jira get <KEY> --md
  ```
  > `jira get <KEY> --md`를 실행하면 이슈 설명(Description), 하위작업(Subtasks), 코멘트(Comments)가 하나의 마크다운 문서로 출력되어 추가 질의 없이 티켓의 전체 작업 문맥을 즉시 파악할 수 있습니다.

---

## 2. 작업 수명주기 및 상태 전이 (Workflow Transitions)

상태를 변경하기 전 반드시 전이 가능한 상태명을 먼저 확인합니다.

1. **전이 가능한 상태 목록 확인**:
   ```bash
   jira transitions <KEY>
   ```
2. **작업 시작 시 상태 변경**:
   ```bash
   jira move <KEY> "진행 중"
   ```
3. **작업 완료 시 코멘트 등록 및 상태 종료**:
   ```bash
   jira comment <KEY> "작업 완료 요약 내용 및 검증 결과"
   jira move <KEY> "완료"
   ```

---

## 3. 티켓 생성 및 하위 작업(Subtask)

1. **팀 표준 라벨 확인**:
   오타 및 비표준 라벨 방지를 위해 생성 전 표준 라벨을 확인합니다.
   ```bash
   jira labels
   ```
2. **신규 이슈 생성**:
   ```bash
   jira create "<요약/제목>" -d "<상세 설명 (Markdown 지원)>" -l "<라벨1,라벨2>"
   ```
   - 주요 옵션:
     - `-t, --type`: 이슈 유형 (작업, 스토리, 에픽, Bug 등 / 기본: '작업')
     - `-l, --labels`: 쉼표로 구분된 라벨 목록
     - `--due`: 마감일 (YYYY-MM-DD)
3. **하위 작업(Subtask) 생성**:
   ```bash
   jira create "<하위 작업 제목>" -t Subtask --parent <KEY> -d "<작업 내용>"
   ```

---

## 4. 티켓 수정 및 코멘트

- **티켓 정보 수정 (라벨, 마감일 등)**:
  ```bash
  jira edit <KEY> -l "AI,Development" --due 2026-10-31
  ```
- **진행 상황 코멘트 추가**:
  ```bash
  jira comment <KEY> "<진행 상황 또는 논의 내용>"
  ```

---

## 5. 에이전트 작업 모범 사례 (Best Practices)

1. **상태 변경 전 검증**: 상태 전이 명령어(`jira move`)를 실행하기 전 반드시 `jira transitions <KEY>` 결과에 대상 상태가 존재하는지 확인합니다.
2. **전체 맥락 우선 파악**: 작업을 시작하기 전 `jira get <KEY> --md`로 기존 논의 내용 및 요구사항을 완전히 파악합니다.
3. **작업 결과 기록**: 코드 수정 및 테스트 검증이 끝난 후에는 `jira comment`로 구체적인 변경 사항과 테스트 결과를 기록하고 이슈를 완료 처리합니다.
