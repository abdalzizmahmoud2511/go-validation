package govalidation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	localeCache = make(map[string]map[string]string)
	localeMu    sync.RWMutex
)

// loadLocale loads a locale file and caches it.
func loadLocale(lang string) map[string]string {
	localeMu.RLock()
	if msgs, ok := localeCache[lang]; ok {
		localeMu.RUnlock()
		return msgs
	}
	localeMu.RUnlock()

	localeMu.Lock()
	defer localeMu.Unlock()

	// Double check after acquiring write lock
	if msgs, ok := localeCache[lang]; ok {
		return msgs
	}

	path := filepath.Join(GetLocalePath(), lang+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Fallback to English if file not found
		if lang != "en" {
			return loadLocale("en")
		}
		return make(map[string]string)
	}

	var msgs map[string]string
	if err := json.Unmarshal(data, &msgs); err != nil {
		return make(map[string]string)
	}

	localeCache[lang] = msgs
	return msgs
}

// getMessage retrieves a localized message by key.
// Supports type suffix like "min_string" or "min_number".
func getMessage(lang, key string, vars map[string]string) string {
	msgs := loadLocale(lang)

	// Try to get type-specific message (e.g., "min_string")
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

	// Fallback to English
	if lang != "en" {
		return getMessage("en", key, vars)
	}

	// Default message if not found
	return "Validation failed for :field"
}

// replaceVars replaces :field, :arg, :min, :max placeholders in message.
func replaceVars(msg string, vars map[string]string) string {
	for k, v := range vars {
		msg = strings.ReplaceAll(msg, ":"+k, v)
	}
	return msg
}

// ClearLocaleCache clears the locale cache.
func ClearLocaleCache() {
	localeMu.Lock()
	defer localeMu.Unlock()
	localeCache = make(map[string]map[string]string)
}

// GetAvailableLanguages returns a list of available language codes.
func GetAvailableLanguages() []string {
	path := GetLocalePath()
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var langs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			lang := strings.TrimSuffix(entry.Name(), ".json")
			langs = append(langs, lang)
		}
	}
	return langs
}
