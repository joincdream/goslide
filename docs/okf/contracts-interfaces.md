---
type: reference
title: Goslide Contracts & Domain Interfaces
description: Specification of core interfaces, domain IR models, sentinel errors, and exit codes
tags:
  - interfaces
  - models
  - contracts
  - reference
timestamp: 2026-10-02T22:50:00Z
trust:
  level: authoritative
  owner: architecture-engineering
---

# Goslide 도메인 인터페이스 및 모델 명세

`internal/model` 패키지는 최하위 계층으로서 Go 표준 라이브러리 외에 외부 의존성을 일체 갖지 않는(Zero Dependency) 순수 데이터 모델 및 인터페이스 계약입니다.

---

## 1. 코어 인터페이스 (`internal/model/interfaces.go`)

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

---

## 2. 불변 데이터 모델 (`internal/model/deck.go`)

```go
package model

import "time"

type SizeRatio string
const (
	Ratio16x9 SizeRatio = "16:9"
	Ratio4x3  SizeRatio = "4:3"
)

type LayoutType string
const (
	LayoutDefault LayoutType = "default"
	LayoutCover   LayoutType = "cover"
	LayoutSection LayoutType = "section"
	LayoutTwoCols LayoutType = "two-cols"
	LayoutBlank   LayoutType = "blank"
)

// Deck은 단일 프레젠테이션 전체를 포괄하는 최상위 불변 모델입니다.
type Deck struct {
	Title       string           `json:"title"`
	Author      string           `json:"author"`
	CreatedAt   time.Time        `json:"created_at"`
	GlobalAttrs GlobalDirectives `json:"global_attributes"`
	CustomCSS   string           `json:"custom_css,omitempty"`
	Slides      []*Slide         `json:"slides"`
}

// GlobalDirectives는 데크 전체의 전역 속성 지시자입니다.
type GlobalDirectives struct {
	Theme    string     `json:"theme"`
	Layout   LayoutType `json:"layout"`
	Size     SizeRatio  `json:"size"`
	Paginate bool       `json:"paginate"`
	Header   string     `json:"header,omitempty"`
	Footer   string     `json:"footer,omitempty"`
}

// SlideDirectives는 개별 슬라이드의 유효 스타일 지시자입니다.
type SlideDirectives struct {
	Layout          LayoutType `json:"layout,omitempty"`
	Class           []string   `json:"class,omitempty"`
	BackgroundColor string     `json:"background_color,omitempty"`
	BackgroundImage string     `json:"background_image,omitempty"`
	Color           string     `json:"color,omitempty"`
	Header          string     `json:"header,omitempty"`
	Footer          string     `json:"footer,omitempty"`
	Paginate        bool       `json:"paginate"`
}

// Slide는 단일 화면 슬라이드 단위 모델입니다.
type Slide struct {
	Index       int             `json:"index"`
	Layout      LayoutType      `json:"layout"`
	Directives  SlideDirectives `json:"directives"`
	RawContent  string          `json:"-"`
	HTMLContent string          `json:"html_content"`
	Notes       string          `json:"notes,omitempty"`
	LeftHTML    string          `json:"left_html,omitempty"`
	RightHTML   string          `json:"right_html,omitempty"`
}
```

---

## 3. 도메인 센티넬 에러 (`internal/model/interfaces.go`)

호출 측에서 `errors.Is()` 또는 `errors.As()`로 명확히 분기할 수 있도록 정의합니다:

```go
var (
	ErrSlideNotFound      = errors.New("slide not found")
	ErrInvalidFrontmatter = errors.New("invalid frontmatter syntax")
	ErrThemeNotFound      = errors.New("theme not found")
	ErrBrowserNotFound    = errors.New("headless browser executable not found")
	ErrExportFailed       = errors.New("export operation failed")
	ErrCanceled           = errors.New("operation canceled")
)
```

---

## 4. CLI 표준 종료 코드 (Exit Codes)

| 코드 | 식별자 | 발생 시나리오 |
| :---: | :--- | :--- |
| **`0`** | `ExitSuccess` | 정상 변환 및 커맨드 수행 완료 |
| **`1`** | `ExitGeneralError` | 분류되지 않은 일반 런타임 오류 또는 패닉 |
| **`2`** | `ExitInvalidUsage` | 잘못된 CLI 플래그 조합, 필수 인자 누락 |
| **`3`** | `ExitFileNotFound` | 입력 마크다운 파일 또는 참조 에셋 미존재 |
| **`4`** | `ExitParseError` | YAML Frontmatter 파싱 실패 또는 마크다운 문법 오류 |
| **`5`** | `ExitExportFailed` | `chromedp` 구동 실패, PDF 인쇄 실패, PPTX 패키징 오류 |
