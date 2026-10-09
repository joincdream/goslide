package i18n_test

import (
	"testing"

	"github.com/yundream/goslide/internal/i18n"
)

func TestI18nTranslations(t *testing.T) {
	b, err := i18n.NewBundle()
	if err != nil {
		t.Fatalf("unexpected error creating bundle: %v", err)
	}

	tests := []struct {
		name     string
		locale   string
		key      string
		args     []any
		expected string
	}{
		{
			name:     "English simple translation",
			locale:   "en",
			key:      "presenter.elapsed",
			args:     nil,
			expected: "Elapsed Time",
		},
		{
			name:     "English with format args",
			locale:   "en",
			key:      "slide.page_indicator",
			args:     []any{3, 24},
			expected: "Page 3 of 24",
		},
		{
			name:     "Korean simple translation",
			locale:   "ko",
			key:      "presenter.elapsed",
			args:     nil,
			expected: "경과 시간",
		},
		{
			name:     "Korean with format args",
			locale:   "ko",
			key:      "slide.page_indicator",
			args:     []any{3, 24},
			expected: "3 / 24 페이지",
		},
		{
			name:     "Fallback to key when missing entirely",
			locale:   "ko",
			key:      "unknown.key",
			args:     nil,
			expected: "unknown.key",
		},
		{
			name:     "Fallback with args when missing key",
			locale:   "en",
			key:      "unknown.format %s",
			args:     []any{"test"},
			expected: "unknown.format test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b.SetLocale(tt.locale)
			got := b.T(tt.key, tt.args...)
			if got != tt.expected {
				t.Errorf("T(%q) with locale %q = %q; want %q", tt.key, tt.locale, got, tt.expected)
			}
		})
	}
}

func TestGlobalI18nHelpers(t *testing.T) {
	i18n.SetLocale("en")
	if got := i18n.T("presenter.elapsed"); got != "Elapsed Time" {
		t.Errorf("global T() for en = %q; want %q", got, "Elapsed Time")
	}

	i18n.SetLocale("ko")
	if got := i18n.T("presenter.elapsed"); got != "경과 시간" {
		t.Errorf("global T() for ko = %q; want %q", got, "경과 시간")
	}
}

func TestGetCatalog(t *testing.T) {
	koCat := i18n.GetCatalog("ko")
	if koCat["ui.toolbar.pen"] != "펜" {
		t.Errorf("ko catalog ui.toolbar.pen = %q; want %q", koCat["ui.toolbar.pen"], "펜")
	}

	enCat := i18n.GetCatalog("en")
	if enCat["ui.toolbar.pen"] != "Pen" {
		t.Errorf("en catalog ui.toolbar.pen = %q; want %q", enCat["ui.toolbar.pen"], "Pen")
	}

	koJSON := i18n.GetCatalogJSON("ko")
	if len(koJSON) < 10 {
		t.Errorf("ko JSON string too short: %q", koJSON)
	}
}
