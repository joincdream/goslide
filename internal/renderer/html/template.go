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
	Color       string
	Header      string
	Footer      string
	Paginate    bool
	HTMLContent template.HTML
	LeftHTML    template.HTML
	RightHTML   template.HTML
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
      <section class="slide-card {{ .Layout }} layout-{{ .Layout }} {{ .Classes }}{{ if .IsFirst }} active{{ end }}"
               data-slide="{{ .Index }}"
               style="{{ if .BgColor }}background-color: {{ .BgColor }};{{ end }}{{ if .BgImage }}background-image: url('{{ .BgImage }}'); background-size: cover; background-position: center;{{ end }}{{ if .Color }}color: {{ .Color }};{{ end }}">
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
      </section>
      {{ end }}
    </div>
  </div>

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
