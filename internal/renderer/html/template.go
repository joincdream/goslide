package html

import (
	"html/template"
)

type documentTemplateData struct {
	Title       string
	Theme       string
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
	TitleHTML   template.HTML
	Footer      string
	Paginate    bool
	Autofit     bool
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
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
  <style id="goslide-theme-styles">
{{ .ComposedCSS }}
  </style>
</head>
<body>
  <div id="goslide-stage" class="goslide-stage">
    <div id="goslide-deck" class="goslide-deck">
      {{ range .Slides }}
      <section class="slide-card {{ .Layout }} layout-{{ .Layout }} {{ .Classes }}{{ if .IsFirst }} active{{ end }}{{ if .BgDim }} has-bg-dim{{ end }}{{ if .Autofit }} has-autofit{{ end }}"
               data-slide="{{ .Index }}"
               {{ if .Autofit }}data-autofit="true"{{ end }}
               style="{{ if .BgColor }}background-color: {{ .BgColor }};{{ end }}{{ if .BgImage }}background-image: url('{{ .BgImage }}'); background-size: cover; background-position: center;{{ end }}{{ if .Color }}color: {{ .Color }};{{ end }}">
        {{ if .BgDim }}<div class="slide-bg-dim" style="background-color: {{ .BgDim }};"></div>{{ end }}
        <div class="slide-header slide-tracker">{{ .Header }}</div>
        {{ if .TitleHTML }}
        <div class="slide-title-box">{{ .TitleHTML }}</div>
        {{ end }}
        <div class="slide-body slide-content-box">
          {{ if and .LeftHTML .RightHTML }}
          <div class="two-cols">
            <div class="col-left">{{ .LeftHTML }}</div>
            <div class="col-right">{{ .RightHTML }}</div>
          </div>
          {{ else }}
          {{ .HTMLContent }}
          {{ end }}
        </div>
        <div class="slide-footer">{{ if or .Footer .Paginate }}<span>{{ .Footer }}</span>{{ if .Paginate }}<span>{{ .Index }}</span>{{ end }}{{ end }}</div>
        {{ if .Notes }}<aside class="slide-notes" style="display:none;">{{ .Notes }}</aside>{{ end }}
      </section>
      {{ end }}
    </div>
  </div>

  <!-- KaTeX Math Rendering Support -->
  <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
  <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js"></script>

  <!-- Mermaid Diagram Rendering Support -->
  <script defer src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
  <script>
    window.__goslide_mermaid_promise = new Promise(function(resolve) {
      function renderMermaid() {
        if (typeof mermaid !== 'undefined') {
          mermaid.initialize({
            startOnLoad: false,
            theme: {{ if eq .Theme "dark" }}'dark'{{ else }}'default'{{ end }},
            securityLevel: 'loose'
          });
          var els = document.querySelectorAll('.mermaid');
          if (els.length > 0) {
            mermaid.run({ querySelector: '.mermaid' }).then(resolve).catch(function(e) {
              console.warn("Mermaid render error:", e);
              resolve();
            });
            return;
          }
        }
        resolve();
      }

      function init() {
        if (typeof mermaid !== 'undefined') {
          renderMermaid();
        } else {
          window.addEventListener('load', renderMermaid);
        }
      }

      if (document.readyState === 'loading') {
        document.addEventListener("DOMContentLoaded", init);
      } else {
        init();
      }
    });

    document.addEventListener("DOMContentLoaded", function() {
      function renderMath() {
        if (typeof renderMathInElement === 'function') {
          renderMathInElement(document.body, {
            delimiters: [
              {left: '$$', right: '$$', display: true},
              {left: '$', right: '$', display: false}
            ],
            throwOnError: false
          });
        }
      }
      if (typeof renderMathInElement === 'function') {
        renderMath();
      } else {
        window.addEventListener('load', renderMath);
      }
    });
  </script>

  <!-- Svelte 5 Application Mount Point -->
  <div id="goslide-app"></div>

  <script id="goslide-runtime-script">
{{ .CoreJS }}
  </script>
</body>
</html>
`

func parseMasterTemplate() (*template.Template, error) {
	return template.New("slideDocument").Parse(masterHTMLTemplate)
}
