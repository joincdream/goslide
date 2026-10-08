package model

// DiagnosticSeverity defines the visual and execution impact of a diagnostic.
type DiagnosticSeverity string

const (
	SeverityError   DiagnosticSeverity = "error"   // 구문 미종료 등 본문 유실 위험
	SeverityWarning DiagnosticSeverity = "warning" // 미지원 레이아웃/지시어 오타 (기본값 폴백)
	SeverityInfo    DiagnosticSeverity = "info"    // 권장 문법 가이드
)

// Diagnostic represents an actionable diagnostic feedback item.
type Diagnostic struct {
	Severity   DiagnosticSeverity `json:"severity"`   // "error" | "warning" | "info"
	SlideIndex int                `json:"slideIndex"` // 1-based 슬라이드 번호
	Line       int                `json:"line"`       // 슬라이드 내 상대 줄 번호
	Rule       string             `json:"rule"`       // "layout.unknown", "directive.missing_underscore" 등
	Message    string             `json:"message"`    // 명확한 팩트 메시지 (예: "Unknown layout 'two-col'")
	RawSnippet string             `json:"rawSnippet"` // 문제가 발생한 원본 텍스트 ("<!-- _layout: two-col -->")
	Candidates []string           `json:"candidates"` // 선택 가능한 유효 옵션 목록 (VS Code QuickFix/Completion 연계용)
}
