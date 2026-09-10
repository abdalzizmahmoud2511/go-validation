package govalidation

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed locale/*.json
var localeFS embed.FS

// allLocales holds all pre-loaded locale messages at package init time.
// Key: language code ("en", "ar"), Value: map of rule -> message template.
var allLocales map[string]map[string]string

// customLocales holds locales loaded from an external path via SetLocalePath().
// These override the embedded locales.
var (
	customLocales map[string]map[string]string
	customMu      sync.RWMutex
	customPath    string
)

func init() {
	allLocales = make(map[string]map[string]string)
	loadAllLocales()
}

// loadAllLocales parses all embedded locale/*.json files at startup.
func loadAllLocales() {
	entries, err := localeFS.ReadDir("locale")
	if err != nil {
		allLocales["en"] = make(map[string]string)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		lang := strings.TrimSuffix(entry.Name(), ".json")

		data, err := localeFS.ReadFile("locale/" + entry.Name())
		if err != nil {
			continue
		}

		var msgs map[string]string
		if err := json.Unmarshal(data, &msgs); err != nil {
			continue
		}
		allLocales[lang] = msgs
	}

	if _, ok := allLocales["en"]; !ok {
		allLocales["en"] = make(map[string]string)
	}
}

// loadCustomLocales loads locale files from an external directory.
// Called when SetLocalePath() is used.
func loadCustomLocales(path string) {
	customMu.Lock()
	defer customMu.Unlock()

	customLocales = make(map[string]map[string]string)

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		lang := strings.TrimSuffix(entry.Name(), ".json")

		data, err := os.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			continue
		}

		var msgs map[string]string
		if err := json.Unmarshal(data, &msgs); err != nil {
			continue
		}
		customLocales[lang] = msgs
	}
}

// getLocaleMessages returns messages for a language.
// Checks custom path first, then embedded, then falls back to English.
func getLocaleMessages(lang string) map[string]string {
	// Check custom locales first
	customMu.RLock()
	if customLocales != nil {
		if msgs, ok := customLocales[lang]; ok {
			customMu.RUnlock()
			return msgs
		}
	}
	customMu.RUnlock()

	// Fall back to embedded locales
	if msgs, ok := allLocales[lang]; ok {
		return msgs
	}

	// Fall back to English embedded
	return allLocales["en"]
}

// getMessage retrieves a localized message by key.
// Uses pre-loaded locale data with no file I/O or mutex overhead.
func getMessage(lang, key string, vars map[string]string) string {
	msgs := getLocaleMessages(lang)

	// Try type-specific message (e.g., "min_string")
	if msgType, ok := vars["_type"]; ok && msgType != "" {
		typeKey := key + "_" + msgType
		if msg, ok := msgs[typeKey]; ok {
			return replaceVars(msg, vars)
		}
	}

	// Fallback to base key (e.g., "min")
	if msg, ok := msgs[key]; ok {
		return replaceVars(msg, vars)
	}

	// Fallback to English if not already English
	if lang != "en" {
		enMsgs := getLocaleMessages("en")
		if msgType, ok := vars["_type"]; ok && msgType != "" {
			typeKey := key + "_" + msgType
			if msg, ok := enMsgs[typeKey]; ok {
				return replaceVars(msg, vars)
			}
		}
		if msg, ok := enMsgs[key]; ok {
			return replaceVars(msg, vars)
		}
	}

	return "Validation failed for :field"
}

// replaceVars replaces :field, :arg, :min, :max placeholders in message.
func replaceVars(msg string, vars map[string]string) string {
	for k, v := range vars {
		msg = strings.ReplaceAll(msg, ":"+k, v)
	}
	return msg
}

// ClearLocaleCache clears the custom locale cache.
// Embedded locales are not affected (they are compiled into the binary).
func ClearLocaleCache() {
	customMu.Lock()
	defer customMu.Unlock()
	customLocales = nil
	customPath = ""
}

// GetAvailableLanguages returns a list of available language codes.
// Combines embedded and custom locale languages.
func GetAvailableLanguages() []string {
	seen := make(map[string]bool)
	var langs []string

	// Add embedded languages
	for lang := range allLocales {
		if !seen[lang] {
			seen[lang] = true
			langs = append(langs, lang)
		}
	}

	// Add custom languages
	customMu.RLock()
	if customLocales != nil {
		for lang := range customLocales {
			if !seen[lang] {
				seen[lang] = true
				langs = append(langs, lang)
			}
		}
	}
	customMu.RUnlock()

	return langs
}

// loadLocale is kept for backward compatibility but returns pre-loaded data.
func loadLocale(lang string) map[string]string {
	return getLocaleMessages(lang)
}
