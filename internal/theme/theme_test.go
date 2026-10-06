package theme_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yundream/goslide/internal/theme"
)

func TestManager_GetBaseCSS(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)
	css, err := mgr.GetBaseCSS()
	if err != nil {
		t.Fatalf("unexpected error getting base css: %v", err)
	}

	if !strings.Contains(css, ".slide-card") {
		t.Errorf("base css missing .slide-card definition")
	}
	if !strings.Contains(css, "word-break: keep-all") {
		t.Errorf("base css missing CJK word-break definition")
	}
}

func TestManager_GetThemeCSS(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)

	tests := []struct {
		name       string
		themeName  string
		wantErr    bool
		errType    error
		containStr string
	}{
		{
			name:       "fallback to default on empty",
			themeName:  "",
			wantErr:    false,
			containStr: "Goslide Default Theme",
		},
		{
			name:       "explicit default theme",
			themeName:  "default",
			wantErr:    false,
			containStr: "Goslide Default Theme",
		},
		{
			name:       "clean theme",
			themeName:  "clean",
			wantErr:    false,
			containStr: "Goslide Clean Theme",
		},
		{
			name:       "dark theme",
			themeName:  "dark",
			wantErr:    false,
			containStr: "Goslide Dark Theme",
		},
		{
			name:       "case insensitivity",
			themeName:  "DARK",
			wantErr:    false,
			containStr: "Goslide Dark Theme",
		},
		{
			name:      "unknown theme",
			themeName: "matrix-neon",
			wantErr:   true,
			errType:   theme.ErrThemeNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			css, err := mgr.GetThemeCSS(tt.themeName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetThemeCSS(%q) error = %v, wantErr %v", tt.themeName, err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error %v, got %v", tt.errType, err)
				}
				return
			}

			if !strings.Contains(css, tt.containStr) {
				t.Errorf("expected css to contain %q, but did not", tt.containStr)
			}
		})
	}
}

func TestManager_ComposeFullCSS_BaseAndTheme(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)
	composed, err := mgr.ComposeFullCSS("clean", "", "")
	if err != nil {
		t.Fatalf("unexpected error composing css: %v", err)
	}

	baseIdx := strings.Index(composed, "/* --- Base Styles --- */")
	themeIdx := strings.Index(composed, "/* --- Theme: clean --- */")

	if baseIdx == -1 || themeIdx == -1 {
		t.Fatalf("missing section headers in composed css")
	}
	if baseIdx >= themeIdx {
		t.Errorf("expected Base Styles to precede Theme Styles in cascading order")
	}
}

func TestManager_ComposeFullCSS_CascadingOrder(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)
	tmpDir := t.TempDir()
	customCSSPath := filepath.Join(tmpDir, "custom.css")
	customContent := ".slide-card { border: 5px solid red; }"
	if err := os.WriteFile(customCSSPath, []byte(customContent), 0600); err != nil {
		t.Fatalf("failed to write custom css: %v", err)
	}

	inlineCSS := "h1 { text-transform: uppercase; }"

	composed, err := mgr.ComposeFullCSS("dark", customCSSPath, inlineCSS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	baseIdx := strings.Index(composed, "/* --- Base Styles --- */")
	themeIdx := strings.Index(composed, "/* --- Theme: dark --- */")
	customIdx := strings.Index(composed, "/* --- Custom CSS File:")
	inlineIdx := strings.Index(composed, "/* --- Inline Styles --- */")

	if baseIdx == -1 || themeIdx == -1 || customIdx == -1 || inlineIdx == -1 {
		t.Fatalf("missing cascading sections: base=%d, theme=%d, custom=%d, inline=%d",
			baseIdx, themeIdx, customIdx, inlineIdx)
	}

	if !(baseIdx < themeIdx && themeIdx < customIdx && customIdx < inlineIdx) {
		t.Errorf("invalid cascading order: expected base < theme < custom < inline")
	}

	if !strings.Contains(composed, customContent) {
		t.Errorf("composed css missing custom file content")
	}
	if !strings.Contains(composed, inlineCSS) {
		t.Errorf("composed css missing inline content")
	}
}

func TestManager_ComposeFullCSS_Errors(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)

	t.Run("non-existent custom css file returns error", func(t *testing.T) {
		t.Parallel()
		_, err := mgr.ComposeFullCSS("default", "/non/existent/path.css", "")
		if err == nil {
			t.Errorf("expected error for non-existent custom css file, got nil")
		}
	})

	t.Run("invalid theme returns error", func(t *testing.T) {
		t.Parallel()
		_, err := mgr.ComposeFullCSS("invalid-theme-xyz", "", "")
		if !errors.Is(err, theme.ErrThemeNotFound) {
			t.Errorf("expected ErrThemeNotFound, got %v", err)
		}
	})
}

func TestManager_GetCoreJS(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)
	js, err := mgr.GetCoreJS()
	if err != nil {
		t.Fatalf("unexpected error getting core js: %v", err)
	}

	if !strings.Contains(js, "goslide-deck") {
		t.Errorf("core js missing goslide-deck identifier")
	}
	if !strings.Contains(js, "goslide-canvas") {
		t.Errorf("core js missing goslide-canvas identifier")
	}
	if !strings.Contains(js, "toggleDrawMode") {
		t.Errorf("core js missing toggleDrawMode function")
	}
}

func TestManager_GetStarterTemplate(t *testing.T) {
	t.Parallel()

	mgr := theme.NewManager(nil)

	t.Run("default template load", func(t *testing.T) {
		tpl, err := mgr.GetStarterTemplate("")
		if err != nil {
			t.Fatalf("unexpected error getting starter template: %v", err)
		}

		if !strings.Contains(tpl, "theme: \"clean\"") {
			t.Errorf("expected default template to contain theme: \"clean\", got:\n%s", tpl[:100])
		}
		if !strings.Contains(tpl, "<!-- _layout: cover -->") {
			t.Errorf("expected starter template to contain cover layout")
		}
		if !strings.Contains(tpl, "flowchart LR") {
			t.Errorf("expected starter template to contain mermaid diagram")
		}
	})

	t.Run("custom theme override", func(t *testing.T) {
		tpl, err := mgr.GetStarterTemplate("dark")
		if err != nil {
			t.Fatalf("unexpected error getting starter template: %v", err)
		}

		if !strings.Contains(tpl, "theme: \"dark\"") {
			t.Errorf("expected overridden template to contain theme: \"dark\", got:\n%s", tpl[:100])
		}
	})
}

