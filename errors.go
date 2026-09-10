package govalidation

import (
	"strconv"
	"strings"
)

// ValError represents a validation error with field, rule, and message.
type ValError struct {
	Field string
	Rule  string
	Msg   string
}

func (e ValError) Error() string {
	return e.Msg
}

func (e ValError) GetField() string {
	return e.Field
}

func (e ValError) GetRule() string {
	return e.Rule
}

func (e ValError) GetMessage() string {
	return e.Msg
}

// locErr creates a ValError with localized message.
// Optimized: no map allocation, direct string replacement.
func locErr(lang, rule, field, msgType string, args ...string) error {
	msg := resolveMsg(lang, rule, msgType, field, args)
	return ValError{field, rule, msg}
}

// resolveMsg looks up the locale message and replaces variables inline.
// Zero map allocations — direct string replacement.
func resolveMsg(lang, rule, msgType, field string, args []string) string {
	msgs, ok := allLocales[lang]
	if !ok {
		msgs = allLocales["en"]
	}

	// Try type-specific message (e.g., "min_string")
	if msgType != "" {
		typeKey := rule + "_" + msgType
		if msg, ok := msgs[typeKey]; ok {
			return fillTemplate(msg, field, args)
		}
	}

	// Fallback to base key (e.g., "min")
	if msg, ok := msgs[rule]; ok {
		return fillTemplate(msg, field, args)
	}

	// Fallback to English
	if lang != "en" {
		if enMsgs, ok := allLocales["en"]; ok {
			if msgType != "" {
				typeKey := rule + "_" + msgType
				if msg, ok := enMsgs[typeKey]; ok {
					return fillTemplate(msg, field, args)
				}
			}
			if msg, ok := enMsgs[rule]; ok {
				return fillTemplate(msg, field, args)
			}
		}
	}

	return "Validation failed for :field"
}

// fillTemplate replaces :field, :arg, :min, :max in message.
// Zero allocations for the common case (no placeholders).
func fillTemplate(msg, field string, args []string) string {
	// Fast path: check if msg contains any placeholders
	if !strings.Contains(msg, ":") {
		return msg
	}

	msg = strings.ReplaceAll(msg, ":field", field)

	if len(args) > 0 {
		msg = strings.ReplaceAll(msg, ":arg", args[0])
	}
	if len(args) > 1 {
		msg = strings.ReplaceAll(msg, ":min", args[0])
		msg = strings.ReplaceAll(msg, ":max", args[1])
	}

	return msg
}

// customErr returns a ValError with custom message if provided, otherwise localized message.
func customErr(lang, field, rule, customMsg, defaultMsg string) error {
	if customMsg != "" {
		msg := resolveLocaleMsg(lang, field, rule, customMsg)
		return ValError{field, rule, msg}
	}
	msg := resolveMsg(lang, rule, "", field, nil)
	if msg == "" {
		msg = defaultMsg
	}
	return ValError{field, rule, msg}
}

// resolveLocaleMsg resolves _locale:key from locale file, otherwise returns msg as-is.
func resolveLocaleMsg(lang, field, rule, msg string) string {
	if strings.HasPrefix(msg, "_locale:") {
		key := strings.TrimPrefix(msg, "_locale:")
		resolved := resolveMsg(lang, key, "", field, nil)
		if resolved != "" {
			return resolved
		}
	}
	return msg
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 10, 64)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
