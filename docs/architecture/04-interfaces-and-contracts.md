---
title: 코어 인터페이스 및 도메인 에러 계약 (Core Interfaces & Domain Contracts)
type: contract
description: Core Go interfaces (Parser, Renderer, Exporter) and domain sentinel error handling specifications
tags:
  - interfaces
  - parser-contract
  - renderer-contract
  - exporter-contract
  - sentinel-errors
related_code:
  - internal/model/interfaces.go
  - internal/model/errors.go
timestamp: 2026-10-05T10:30:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# 4. 코어 인터페이스 및 도메인 에러 계약 (Interfaces & Contracts)

Goslide는 의존성 역전 원칙(DIP)과 인터페이스 분리 원칙(ISP)을 충실히 지키기 위해 비즈니스 로직과 포맷 변환기를 명확한 인터페이스 계약으로 분리합니다.

## 4.1 공통 인터페이스 명세 ([`internal/model/interfaces.go`](../../internal/model/interfaces.go))

```go
package model

import (
	"context"
	"io"
)

// Parser는 마크다운 소스를 분석하여 불변 Presentation Deck을 빌드합니다.
type Parser interface {
	Parse(ctx context.Context, r io.Reader) (*Deck, error)
}

// Renderer는 Deck을 특정 포맷(주로 HTML)의 스트림으로 렌더링합니다.
type Renderer interface {
	Render(ctx context.Context, deck *Deck, w io.Writer) error
}

// Exporter는 Deck을 외부 파일 포맷(PDF, PPTX)으로 생성하여 파일 시스템에 출력합니다.
type Exporter interface {
	Export(ctx context.Context, deck *Deck, outputPath string) error
}
```

### 인터페이스 설계 원칙:
1. **Context 전파 필수**: 모든 인터페이스 메서드는 첫 번째 인자로 `context.Context`를 전달받아 취소 및 타임아웃을 감지합니다.
2. **I/O 스트림 추상화**: `Renderer`는 물리 파일 대신 `io.Writer`를 받아 인메모리 버퍼 렌더링(PDF/PPTX 파이프라인 연계)과 파일 출력을 모두 지원합니다.
3. **불변 모델 전달**: `deck *Deck` 포인터는 읽기 전용으로 전달되며, 수신자가 내부 상태를 변경해서는 안 됩니다.

---

## 4.2 도메인 센티넬 에러 명세 ([`internal/model/errors.go`](../../internal/model/errors.go))

호출 측에서 `errors.Is()` 또는 `errors.As()`를 통해 명확히 분기할 수 있도록 계층형 에러 상수를 정의합니다:

```go
package model

import "errors"

var (
	// 파싱 관련 에러
	ErrSlideNotFound      = errors.New("slide not found")
	ErrInvalidFrontmatter = errors.New("invalid frontmatter syntax")
	ErrThemeNotFound      = errors.New("theme not found")

	// 익스포트 및 런타임 관련 에러
	ErrBrowserNotFound    = errors.New("headless browser executable not found")
	ErrExportFailed       = errors.New("export operation failed")
	ErrCanceled           = errors.New("operation canceled")
)
```

---

## 4.3 에러 처리 및 전파 규칙 (Error Handling Idioms)

1. **명시적 에러 래핑 (Go 1.13+ `fmt.Errorf("...: %w", err)`)**:
   - 하위 계층에서 발생한 구체적인 에러 원인을 손실 없이 상위로 전파하기 위해 `%w` 서식을 준수합니다.
   ```go
   // Good
   if err := p.parseFrontmatter(scanner); err != nil {
       return nil, fmt.Errorf("%w: %v", model.ErrInvalidFrontmatter, err)
   }
   ```
2. **Panic 금지 (No Panics)**:
   - 라이브러리 및 내부 패키지 로직에서 `panic`은 전면 금지됩니다. 모든 예외 상황은 명시적인 `error` 반환으로 처리합니다. (단, 컴파일 타임 정규식 `regexp.MustCompile` 등 정적 초기화 제외).
3. **컨텍스트 취소 우선 확인**:
   - 무거운 I/O나 루프 진입 전 반드시 `ctx.Err()`를 확인합니다.
   ```go
   if err := ctx.Err(); err != nil {
       return fmt.Errorf("%w: %v", model.ErrCanceled, err)
   }
   ```

---

## 4.4 관련 문서

* **[`01-principles-and-vision.md`](./01-principles-and-vision.md)**: 5대 원칙 (원칙 4: DIP/ISP, 원칙 5: Resource Safety)
* **[`05-rendering-and-export-pipelines.md`](./05-rendering-and-export-pipelines.md)**: `Renderer`와 `Exporter` 인터페이스 구현체 세부 사항
* **[`07-operations-build-and-cicd.md`](./07-operations-build-and-cicd.md)**: 센티넬 에러가 CLI 표준 종료 코드로 매핑되는 규칙
