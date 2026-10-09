package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

// Bundle manages internationalization message catalogs.
type Bundle struct {
	mu       sync.RWMutex
	locale   string
	messages map[string]map[string]string
}

var (
	defaultBundle *Bundle
	once          sync.Once
	initErr       error
)

// NewBundle creates and initializes a new i18n Bundle from embedded locales.
func NewBundle() (*Bundle, error) {
	b := &Bundle{
		locale:   "en",
		messages: make(map[string]map[string]string),
	}

	for _, lang := range []string{"en", "ko"} {
		data, err := localesFS.ReadFile(fmt.Sprintf("locales/%s.json", lang))
		if err != nil {
			return nil, fmt.Errorf("failed to read locale file for %s: %w", lang, err)
		}

		var catalog map[string]string
		if err := json.Unmarshal(data, &catalog); err != nil {
			return nil, fmt.Errorf("failed to parse locale json for %s: %w", lang, err)
		}
		b.messages[lang] = catalog
	}

	return b, nil
}

// GetDefaultBundle returns the singleton i18n Bundle initialized with system locale.
func GetDefaultBundle() (*Bundle, error) {
	once.Do(func() {
		defaultBundle, initErr = NewBundle()
		if initErr == nil {
			defaultBundle.SetLocale(DetectLocale())
		}
	})
	return defaultBundle, initErr
}

// SetLocale changes the active locale for the bundle.
func (b *Bundle) SetLocale(lang string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "auto" || lang == "" {
		lang = DetectLocale()
	}

	if _, ok := b.messages[lang]; ok {
		b.locale = lang
	} else {
		b.locale = "en"
	}
}

// GetLocale returns the currently active locale.
func (b *Bundle) GetLocale() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.locale
}

// T translates a key using the active locale, falling back to English if missing.
func (b *Bundle) T(key string, args ...any) string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	msg := key
	if catalog, ok := b.messages[b.locale]; ok {
		if val, found := catalog[key]; found {
			msg = val
		}
	}

	// Fallback to English if not found
	if msg == key && b.locale != "en" {
		if enCatalog, ok := b.messages["en"]; ok {
			if val, found := enCatalog[key]; found {
				msg = val
			}
		}
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// DetectLocale inspects system environment variables (LC_ALL, LC_MESSAGES, LANG).
func DetectLocale() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		val := strings.ToLower(os.Getenv(env))
		if strings.HasPrefix(val, "ko") {
			return "ko"
		}
		if strings.HasPrefix(val, "en") {
			return "en"
		}
	}
	return "en"
}

// Global convenience functions for CLI and application use.

// Lookup translates a message key without variadic format arguments.
func Lookup(key string) string {
	b, err := GetDefaultBundle()
	if err != nil {
		return key
	}
	return b.T(key)
}

// T translates a message key using the default bundle.
func T(key string, args ...any) string {
	b, err := GetDefaultBundle()
	if err != nil {
		if len(args) > 0 {
			return fmt.Sprintf(key, args...)
		}
		return key
	}
	return b.T(key, args...)
}

// SetLocale sets the locale on the default bundle.
func SetLocale(lang string) {
	b, err := GetDefaultBundle()
	if err == nil {
		b.SetLocale(lang)
	}
}

// GetCatalog returns a copy of message catalog for the given language, falling back to English.
func (b *Bundle) GetCatalog(lang string) map[string]string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "auto" || lang == "" {
		lang = DetectLocale()
	}

	result := make(map[string]string)
	// Base with English
	if enCatalog, ok := b.messages["en"]; ok {
		for k, v := range enCatalog {
			result[k] = v
		}
	}
	// Overlay with target lang if different
	if lang != "en" {
		if targetCatalog, ok := b.messages[lang]; ok {
			for k, v := range targetCatalog {
				result[k] = v
			}
		}
	}
	return result
}

// GetCatalogJSON returns the JSON serialized catalog for the given language.
func (b *Bundle) GetCatalogJSON(lang string) string {
	cat := b.GetCatalog(lang)
	bytes, err := json.Marshal(cat)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// GetCatalog returns a copy of message catalog using the default bundle.
func GetCatalog(lang string) map[string]string {
	b, err := GetDefaultBundle()
	if err != nil {
		return make(map[string]string)
	}
	return b.GetCatalog(lang)
}

// GetCatalogJSON returns the JSON serialized catalog using the default bundle.
func GetCatalogJSON(lang string) string {
	b, err := GetDefaultBundle()
	if err != nil {
		return "{}"
	}
	return b.GetCatalogJSON(lang)
}
