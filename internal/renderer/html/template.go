package html

import (
	"html/template"
)

type documentTemplateData struct {
	Title       string
	ComposedCSS template.CSS
	CoreJS      template.JS
	Slides      []slideTemplateData
}

type slideTemplateData struct {
	Index       int
	IsFirst     bool
	Layout      string
	Classes     string
	BgColor     string
	BgImage     string
	BgDim       template.CSS
	Color       string
	Header      string
	Footer      string
	Paginate    bool
	HTMLContent template.HTML
	LeftHTML    template.HTML
	RightHTML   template.HTML
	Notes       string
}

const masterHTMLTemplate = `<!DOCTYPE html>
<html lang="ko">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
  <style id="goslide-theme-styles">
{{ .ComposedCSS }}
  </style>
</head>
<body>
  <div id="goslide-stage" class="goslide-stage">
    <div id="goslide-deck" class="goslide-deck">
      {{ range .Slides }}
      <section class="slide-card {{ .Layout }} layout-{{ .Layout }} {{ .Classes }}{{ if .IsFirst }} active{{ end }}{{ if .BgDim }} has-bg-dim{{ end }}"
               data-slide="{{ .Index }}"
               style="{{ if .BgColor }}background-color: {{ .BgColor }};{{ end }}{{ if .BgImage }}background-image: url('{{ .BgImage }}'); background-size: cover; background-position: center;{{ end }}{{ if .Color }}color: {{ .Color }};{{ end }}">
        {{ if .BgDim }}<div class="slide-bg-dim" style="background-color: {{ .BgDim }};"></div>{{ end }}
        {{ if .Header }}<div class="slide-header">{{ .Header }}</div>{{ end }}
        <div class="slide-body">
          {{ if and .LeftHTML .RightHTML }}
          <div class="two-cols">
            <div class="col-left">{{ .LeftHTML }}</div>
            <div class="col-right">{{ .RightHTML }}</div>
          </div>
          {{ else }}
          {{ .HTMLContent }}
          {{ end }}
        </div>
        <div class="slide-footer">
          <span>{{ .Footer }}</span>
          {{ if .Paginate }}<span>{{ .Index }}</span>{{ end }}
        </div>
        {{ if .Notes }}<aside class="slide-notes" style="display:none;">{{ .Notes }}</aside>{{ end }}
      </section>
      {{ end }}
    </div>
  </div>

  <!-- Interactive Presentation Tools -->
  <div id="goslide-laser" class="goslide-laser"></div>
  <div id="goslide-spotlight" class="goslide-spotlight"></div>
  <div id="goslide-blackout" class="goslide-blackout"></div>
  <div id="goslide-whiteout" class="goslide-whiteout"></div>

  <!-- In-Window Presenter Sidebar -->
  <aside id="goslide-presenter-sidebar" class="goslide-presenter-sidebar">
    <div class="sidebar-header">
      <div class="sidebar-title">🎙️ Presenter View</div>
      <div class="sidebar-actions">
        <button id="btn-sidebar-popout" class="sidebar-btn" title="별도 창으로 분리 (P)">↗ Pop out</button>
        <button id="btn-sidebar-close" class="sidebar-btn" title="사이드바 닫기 (N)">✕</button>
      </div>
    </div>
    <div class="sidebar-mode-selector">
      <button id="btn-mode-fit" class="mode-tab active" title="창 크기에 맞게 자동 리사이즈">↔ 화면 맞춤</button>
      <button id="btn-mode-1080p" class="mode-tab" title="1920×1080 고정 (스크린캐스트용)">🔒 1080p 고정</button>
    </div>
    <div class="sidebar-timer-bar">
      <span id="sidebar-timer-display" class="sidebar-timer">00:00:00</span>
      <button id="sidebar-timer-toggle" class="sidebar-btn-sm">시작</button>
      <button id="sidebar-timer-reset" class="sidebar-btn-sm">리셋</button>
      <span id="sidebar-clock" class="sidebar-clock"></span>
    </div>
    <div class="sidebar-section">
      <div class="sidebar-section-title">다음 슬라이드 (<span id="sidebar-next-index">Slide 2</span>)</div>
      <div id="sidebar-next-preview" class="sidebar-preview-box"></div>
    </div>
    <div class="sidebar-section sidebar-notes-section">
      <div class="sidebar-section-title">발표자 메모 (Notes)</div>
      <div id="sidebar-notes-content" class="sidebar-notes-body"></div>
    </div>
    <div class="sidebar-footer">
      <span id="sidebar-progress-text">1 / {{ len .Slides }}</span>
      <div class="sidebar-progress-bar">
        <div id="sidebar-progress-fill" class="sidebar-progress-fill"></div>
      </div>
    </div>
  </aside>

  <canvas id="goslide-canvas" class="goslide-canvas"></canvas>
  <div id="goslide-indicator" class="goslide-indicator">1 / {{ len .Slides }}</div>

  <script id="goslide-runtime-script">
{{ .CoreJS }}
  </script>
</body>
</html>
`

func parseMasterTemplate() (*template.Template, error) {
	return template.New("slideDocument").Parse(masterHTMLTemplate)
}
