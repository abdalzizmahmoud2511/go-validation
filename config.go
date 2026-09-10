package govalidation

import "sync"

// Config holds the validator configuration.
type Config struct {
	LocalePath string // Path to locale directory
	Language   string // Language code (e.g., "en", "ar")
}

var (
	defaultConfig = Config{
		LocalePath: "locale",
		Language:   "en",
	}
	configMu sync.RWMutex
)

// SetConfig sets the global validator configuration.
func SetConfig(cfg Config) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfig = cfg
}

// GetConfig returns the current global configuration.
func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfig
}

// SetLanguage sets the current language.
func SetLanguage(lang string) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfig.Language = lang
}

// GetLanguage returns the current language.
func GetLanguage() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfig.Language
}

// SetLocalePath sets the path to a custom locale directory and loads locales from it.
// Custom locales override embedded locales.
func SetLocalePath(path string) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfig.LocalePath = path
	// Load custom locales from the specified path
	loadCustomLocales(path)
}

// GetLocalePath returns the current locale path.
func GetLocalePath() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfig.LocalePath
}
