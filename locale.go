package govalidation

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed locale/*.json
var localeFS embed.FS

// allLocales holds all pre-loaded locale messages at package init time.
// Key: language code ("en", "ar"), Value: map of rule -> message template.
var allLocales map[string]map[string]string

func init() {
	allLocales = make(map[string]map[string]string)
	loadAllLocales()
}

// loadAllLocales parses all embedded locale/*.json files at startup.
func loadAllLocales() {
	entries, err := localeFS.ReadDir("locale")
	if err != nil {
		// Fallback: empty locales
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

	// Ensure English always exists as fallback
	if _, ok := allLocales["en"]; !ok {
		allLocales["en"] = make(map[string]string)
	}
}

// getMessage retrieves a localized message by key.
// Uses pre-loaded locale data with no file I/O or mutex overhead.
func getMessage(lang, key string, vars map[string]string) string {
	msgs, ok := allLocales[lang]
	if !ok {
		// Fallback to English
		msgs = allLocales["en"]
	}

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
		if enMsgs, ok := allLocales["en"]; ok {
			// Try type-specific in English
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

// ClearLocaleCache is kept for API compatibility but is now a no-op
// since locales are embedded and loaded at init time.
func ClearLocaleCache() {}

// GetAvailableLanguages returns a list of available language codes.
func GetAvailableLanguages() []string {
	langs := make([]string, 0, len(allLocales))
	for lang := range allLocales {
		langs = append(langs, lang)
	}
	return langs
}

// loadLocale is kept for backward compatibility but returns pre-loaded data.
func loadLocale(lang string) map[string]string {
	if msgs, ok := allLocales[lang]; ok {
		return msgs
	}
	return allLocales["en"]
}
