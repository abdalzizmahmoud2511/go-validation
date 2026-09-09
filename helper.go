package main

import (
	"reflect"
	"strings"
)

// isZero returns true if the value is zero/empty.
func isZero(fv reflect.Value) bool {
	switch fv.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return fv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return fv.IsNil()
	case reflect.Struct:
		return fv.IsZero()
	default:
		return fv.IsZero()
	}
}

// isNestedStruct returns true if the value is a struct (or pointer to struct).
func isNestedStruct(fv reflect.Value) bool {
	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}
	return fv.Kind() == reflect.Struct
}

// hasRule checks if a semicolon-separated rule string contains a specific rule.
func hasRule(rules, target string) bool {
	for _, r := range strings.Split(rules, ";") {
		key := strings.TrimSpace(r)
		if i := strings.Index(key, "="); i >= 0 {
			key = key[:i]
		}
		if i := strings.Index(key, "|"); i >= 0 {
			key = key[:i]
		}
		if key == target {
			return true
		}
	}
	return false
}

// flattenStruct converts a nested struct into a flat map with dot notation keys.
func flattenStruct(rv reflect.Value, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		fv := rv.Field(i)
		key := field.Name
		if prefix != "" {
			key = prefix + "." + key
		}

		inner := fv
		if inner.Kind() == reflect.Ptr {
			if inner.IsNil() {
				result[key] = nil
				continue
			}
			inner = inner.Elem()
		}

		if inner.Kind() == reflect.Struct {
			nested := flattenStruct(inner, key)
			for k, v := range nested {
				result[k] = v
			}
		} else {
			result[key] = fv.Interface()
		}
	}
	return result
}

// parseRule splits "rule=arg" into key and arg.
func parseRule(rule string) (key, arg string) {
	rule = strings.TrimSpace(rule)
	if i := strings.Index(rule, "="); i >= 0 {
		return rule[:i], rule[i+1:]
	}
	return rule, ""
}

// parseRuleWithMsg splits "rule|message" into key and message.
func parseRuleWithMsg(rule string) (key, msg string) {
	rule = strings.TrimSpace(rule)
	if i := strings.Index(rule, "|"); i >= 0 {
		return rule[:i], rule[i+1:]
	}
	return rule, ""
}

// parseRuleFull splits "rule=arg|message" into key, arg, and message.
func parseRuleFull(rule string) (key, arg, msg string) {
	rule = strings.TrimSpace(rule)

	// Extract message first (after |)
	if i := strings.Index(rule, "|"); i >= 0 {
		msg = rule[i+1:]
		rule = rule[:i]
	}

	// Extract key and arg (after =)
	if i := strings.Index(rule, "="); i >= 0 {
		return rule[:i], rule[i+1:], msg
	}
	return rule, "", msg
}

func fvToString(fv reflect.Value) string {
	switch fv.Kind() {
	case reflect.String:
		return fv.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intToString(fv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintToString(fv.Uint())
	case reflect.Float32, reflect.Float64:
		return floatToString(fv.Float())
	default:
		return ""
	}
}

func intToString(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func uintToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func floatToString(f float64) string {
	return strings.TrimRight(strings.TrimRight(
		strings.Replace(intToString(int64(f))+".0", ".0", ".", 1)+intToString(int64((f-float64(int64(f)))*1000000)),
		"0"), ".")
}

func splitRange(arg string) []string {
	return splitComma(arg)
}

// Msg returns a rule string with a custom message from locale.
// Usage: {"Name": {Msg("required", "custom_required"), Msg("min", "custom_min")}}
// The key is resolved from the locale file at validation time.
func Msg(rule, localeKey string) string {
	return rule + "|_locale:" + localeKey
}

func splitComma(s string) []string {
	parts := make([]string, 0)
	for _, p := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
