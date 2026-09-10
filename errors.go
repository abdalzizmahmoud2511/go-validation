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
func locErr(lang, rule, field, msgType string, args ...string) error {
	vars := map[string]string{
		"field": field,
	}

	if len(args) > 0 {
		vars["arg"] = args[0]
	}
	if len(args) > 1 {
		vars["min"] = args[0]
		vars["max"] = args[1]
	}
	if msgType != "" {
		vars["_type"] = msgType
	}

	msg := getMessage(lang, rule, vars)
	return ValError{field, rule, msg}
}

// customErr returns a ValError with custom message if provided, otherwise localized message.
func customErr(lang, field, rule, customMsg, defaultMsg string) error {
	if customMsg != "" {
		msg := resolveLocaleMsg(lang, field, rule, customMsg)
		return ValError{field, rule, msg}
	}
	vars := map[string]string{
		"field": field,
	}
	msg := getMessage(lang, rule, vars)
	if msg == "" {
		msg = defaultMsg
	}
	return ValError{field, rule, msg}
}

// resolveLocaleMsg resolves _locale:key from locale file, otherwise returns msg as-is.
func resolveLocaleMsg(lang, field, rule, msg string) string {
	if strings.HasPrefix(msg, "_locale:") {
		key := strings.TrimPrefix(msg, "_locale:")
		vars := map[string]string{"field": field}
		resolved := getMessage(lang, key, vars)
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
