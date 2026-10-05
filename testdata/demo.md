---
title: "현대 대형 언어 모델(LLM) 아키텍처와 엔지니어링 실무"
author: "AI Core Architecture Lab"
theme: "clean"
size: "16:9"
paginate: true
header: "2026 AI Tech Seminar: Deep Dive into LLMs"
footer: "© 2026 Cloit Tech Architecture. All rights reserved."
---

<!-- _layout: cover -->

# 현대 대형 언어 모델(LLM) 아키텍처
### 트랜스포머 코어 메커니즘부터 에이전틱 AI 파이프라인까지

**발표자**: 클로잇 AI 기술연구소  
**일시**: 2026년 10월 기술 세미나  
**단축키 가이드**: <kbd>P</kbd> 발표자 뷰 | <kbd>N</kbd> 사이드바 | <kbd>D</kbd> 판서 모드 | <kbd>?</kbd> 전체 단축키

<!-- note:
[발표 도입부]
- 참석자분들께 환영 인사를 건네고, 오늘 세미나가 단순 개념 소개를 넘어 실제 프로덕션 레벨의 엔지니어링 아키텍처를 다룬다는 점을 강조합니다.
- 'P' 키를 눌러 발표자 콘솔을 열거나, 'N' 키로 인-윈도우 사이드바를 활성화하여 타이머를 시작하십시오.
-->

---

## 세미나 아젠다 & 진행 로드맵

오늘 다룰 4가지 핵심 엔지니어링 주제 및 세션 진행 현황입니다:

- [x] **트랜스포머 코어**: Self-Attention과 인과적 언어 모델링(Causal LM)
- [ ] **최신 모델 비교**: 파라미터 스케일링 법칙과 벤치마크 평가
- [ ] **도구 연동과 에이전트**: Tool Calling 및 구조화된 JSON 출력
- [ ] **프로덕션 서빙**: Go 기반 고성능 스트리밍 프록시 아키텍처

> [!NOTE] 세션 구성 안내
> 각 챕터는 이론 5분, 실무 코드 분석 5분으로 구성되어 있습니다. 발표 도중 언제든 질문해 주세요.

<!-- note:
- 목차를 간단히 브리핑하고 각 챕터별 소요 시간(각 10분 내외)을 안내합니다.
- 청중의 사전 지식(트랜스포머 경험 여부)을 가볍게 확인하며 아이스브레이킹을 진행합니다.
-->

---

<!--
_class: lead
-->

# The Bitter Lesson

> *"General methods that leverage computation are ultimately the most effective by a large margin."*
>
> — **Rich Sutton** (AI Pioneer)

<!-- note:
- 키노트 원포인트 강조 슬라이드입니다.
- AI 연구 역사에서 인간의 수작업 규칙보다 컴퓨팅 파워를 활용한 범용 학습이 결국 승리했다는 핵심 철학을 공유합니다.
-->

---

<!-- _layout: two-cols -->
## 패러다임의 전환: RNN vs Transformer

### 기존 시퀀스 모델 (RNN / LSTM)
- ~~순차 처리(Sequential)~~: 이전 은닉 상태에 묶여 병렬화 불가
- ~~학습 속도 급감~~: 긴 문장에서 $O(N)$ 시간 복잡도 누적
- ~~장기 의존성 소실~~: 문맥이 길어지면 앞쪽 정보 유실

<!-- split -->

### 트랜스포머 (Transformer Decoder)
- **완전 병렬화(Parallelized)**: GPU 텐서 코어 동시 연산
- **Self-Attention ($O(1)$ 경로)**: 거리 무관 토큰 간 직접 상관도 계산
- **스케일링 법칙(Scaling Law)**: 연산량 증가에 비례해 지능 폭발

<!-- note:
- 'L' 키를 눌러 레이저 포인터를 켜고, 우측의 '완전 병렬화'와 'Self-Attention' 항목을 강조해 주세요.
- 왜 GPU 클러스터에서 트랜스포머가 표준이 되었는지를 시간 복잡도 관점에서 설명합니다.
-->

---

## Scaled Dot-Product Attention 수식과 원리

어텐션 연산은 **Query**, **Key**, **Value** 세 벡터 간의 유사도 가중 합으로 계산됩니다:

$$\text{Attention}(Q, K, V) = \text{softmax}\left(\frac{QK^T}{\sqrt{d_k}}\right)V$$

* **Query ($Q$)**: 현재 주목하고자 하는 토큰의 질의 벡터
* **Key ($K$)**: 문맥 내 모든 토큰들이 가진 고유 속성 인덱스
* **Value ($V$)**: 실제로 조합하여 추출할 정보 벡터
* **Scaling Factor ($\sqrt{d_k}$)**: 차원 증가 시 Softmax 기울기 소실(Vanishing Gradient) 방지

<!-- note:
- 수식의 각 항이 데이터베이스의 질의(Query), 인덱스 키(Key), 실제 레코드(Value)와 유사한 개념임을 비유로 설명합니다.
- 'D' 키를 눌러 화면에 자유 필기(Annotation)를 켜고, 분모에 원을 그려 스케일링 계수의 역할을 직관적으로 짚어주세요.
-->

---

<!--
_class: lead
_backgroundImage: ./hands-on-bg.jpeg
_backgroundDim: 0.55
_color: #ffffff
-->

# HANDS ON
### 실무 엔지니어링: 도구 연동부터 서빙 파이프라인까지

<!-- note:
[세션 전환: Hands-on 실무 단계]
- 트랜스포머의 핵심 이론 파트를 마치고, 실제 구현 및 엔지니어링 실습 단계로 전환함을 안내합니다.
- 파이썬 Tool Calling과 Go 스트리밍 프록시 구현체를 차례로 살펴봅니다.
-->

---

<!--
_class: compact
-->

## Python: Tool Calling (함수 호출) 구현 예제

LLM이 정형화된 JSON 도구 스키마를 판별하여 외부 API를 실행하는 표준 패턴입니다:

```python
import json
from openai import OpenAI

client = OpenAI()
tools = [{
    "type": "function",
    "function": {
        "name": "query_database",
        "description": "기업 재무 데이터베이스에서 분기별 실적 조회",
        "parameters": {
            "type": "object",
            "properties": {
                "ticker": {"type": "string", "description": "종목 코드 (예: AAPL)"},
                "quarter": {"type": "string", "enum": ["Q1", "Q2", "Q3", "Q4"]}
            },
            "required": ["ticker", "quarter"]
        }
    }
}]
response = client.chat.completions.create(model="gpt-4o", messages=[{"role": "user", "content": "AAPL 3분기 실적 분석"}], tools=tools)
```

> [!TIP] 프롬프트 엔지니어링 팁
> 스키마의 `description`을 명확하고 구체적으로 서술할수록 도구 선택 정확도(Tool Selection Accuracy)가 비약적으로 상승합니다.

<!-- note:
- 코드 블록의 Chroma 구문 강조 품질을 확인합니다.
- LLM은 함수를 직접 실행하는 것이 아니라, 실행해야 할 '함수 이름과 인자 JSON'을 생성하는 것임을 청중에게 명확히 짚어줍니다.
-->

---

## 2026 주요 오픈/상용 LLM 스펙 비교

현대 엔터프라이즈 환경에서 도입을 고려하는 대표 모델들의 사양 비교입니다:

| 모델 명칭 | 제공사 | 파라미터 | 콘텍스트 | 라이선스 | 주 활용 분야 |
| :--- | :---: | ---: | ---: | :---: | :--- |
| **GPT-4o** | OpenAI | 비공개 (MoE) | 128k | 상용 API | 고난도 추론, 멀티모달 |
| **Claude 3.5 Sonnet** | Anthropic | 비공개 | 200k | 상용 API | 코딩 에이전트, 문서 분석 |
| **Llama 3.3** | Meta | 70B Dense | 128k | Community | 사내 온프레미스 구축 |
| **Qwen 2.5 Coder** | Alibaba | 32B Dense | 128k | Apache 2.0 | 오픈소스 코드 생성 특화 |
| **Gemma 2** | Google | 27B Dense | 8k | Terms | 엣지 디바이스 및 저전력 서빙 |

> [!NOTE] 엔터프라이즈 도입 권고
> 보안 규제가 엄격한 금융/의료 도메인은 Apache 2.0 또는 Community 라이선스 기반의 사내 자체 프라이빗 호스팅을 권장합니다.

<!-- note:
- 오픈소스 모델(Llama, Qwen)의 약진으로 사내 온프레미스(On-premise) 프라이빗 서빙 요구가 급증했음을 언급합니다.
- 표의 좌/중앙/우측 정렬 상태를 짚으며 가독성을 확인합니다.
-->

---

<!-- _layout: two-cols -->
## RAG vs 에이전틱(Agentic) 파이프라인 비교

### 단순 검색 증강 (Naïve RAG)
- **단방향 파이프라인**: 질의 $\rightarrow$ 임베딩 $\rightarrow$ 검색 $\rightarrow$ 답변
- **단일 턴(Single Turn)**: 검색 실패 시 잘못된 환각 생성
- **정적 지식 한정**: 외부 DB 쓰기나 조건부 재검색 불가

<!-- split -->

### 에이전틱 워크플로우 (Agentic AI)
- **반복적 추론(ReAct Loop)**: Thought $\rightarrow$ Action $\rightarrow$ Observation
- **자율 도구 사용**: SQL 실행, 코드 실행, 웹 브라우징
- **자기 교정(Self-Correction)**: 결과 검증 실패 시 스스로 쿼리 재시도

<!-- note:
- 'S' 키를 눌러 스포트라이트를 켜고, 우측의 '자기 교정(Self-Correction)' 부분에 마우스를 올려 주목도를 높여보세요.
- 단순 RAG에서 에이전틱 RAG로 진화하는 업계 트렌드를 설명합니다.
-->

---

<!--
_class: compact
-->

## Go 기반 초경량 LLM 스트리밍 서빙 프록시

Go의 고루틴과 표준 `net/http` SSE를 결합한 무할당(Zero-Alloc) 토큰 스트리밍 아키텍처:

```go
package main

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "net/http"
)

// StreamTokenProxy handles real-time token streaming with zero memory allocations
func StreamTokenProxy(ctx context.Context, w http.ResponseWriter, upstreamBody io.Reader) error {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    flusher, ok := w.(http.Flusher)
    if !ok {
        return fmt.Errorf("streaming unsupported by client")
    }

    scanner := bufio.NewScanner(upstreamBody)
    for scanner.Scan() {
        tokenLine := scanner.Text()
        fmt.Fprintf(w, "data: %s\n\n", tokenLine)
        flusher.Flush()
    }
    return scanner.Err()
}
```

> [!IMPORTANT] 무할당(Zero-Alloc) 고성능 원칙
> Go 런타임의 가비지 컬렉터(GC) 압력을 최소화하기 위해 스트리밍 버퍼 재사용과 채널 버퍼링을 적극 적용합니다.

<!-- note:
- Goslide의 Live Reload 기능(`goslide serve`)에서도 동일하게 활용된 Go 표준 SSE 메커니즘임을 강조합니다.
-->

---

<!-- _layout: section -->

# 엔지니어링 결론 및 로드맵
### "단순 호출을 넘어 자율 에이전트와 도메인 최적화로"

* **결정론적 로직 연계**: LLM 환각을 방어하는 스키마 검증 레이어 확립
* **사내 프라이빗 호스팅**: 오픈 가중치 모델을 활용한 TCO 절감 및 데이터 주권 보호
* **관측 가능성(Observability)**: 토큰 지연시간(TTFT) 및 분당 요청수(RPM) 실시간 계측

<!-- note:
- 테마 반전(Invert) 섹션 간지 슬라이드입니다.
- 결론을 브리핑하며 다음 질의응답 세션으로 자연스럽게 유도합니다.
-->

---

<!-- _layout: cover -->
<!-- _autofit: true -->

# 감사합니다 (Q & A)
### 질문과 자유 토론을 환영합니다

* **세미나 슬라이드 레포지토리**: `github.com/yundream/goslide`
* **빌더 엔진**: Goslide Pure Go Presentation Builder (Single Binary)
* **프레젠테이션 도구**: 1080p 스크린캐스트, 실시간 판서, 레이저 포인터 지원

- [x] 프레젠테이션 핵심 전달 완료
- [x] 라이브 프리뷰 및 실시간 핫 리로드 시연 완료
- [ ] 질의응답 및 테크 피드백 세션 진행

<!-- note:
[발표 마무리]
- 청중의 질문을 유도하고 세미나를 마무리합니다.
- '?' 키를 눌러 단축키 치트시트를 띄우고 참여자들과 질의응답을 나눕니다.
-->
