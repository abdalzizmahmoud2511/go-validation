package govalidation

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	alphaRe      = regexp.MustCompile(`^[a-zA-Z]+$`)
	alphaNumRe   = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	alphaDashRe  = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	alphaSpaceRe = regexp.MustCompile(`^[-a-zA-Z0-9_ ]+$`)
	asciiRe      = regexp.MustCompile(`^[\x00-\x7F]+$`)
	printASCIIRe = regexp.MustCompile(`^[\x20-\x7E]+$`)
	numberRe     = regexp.MustCompile(`^\d+(\.\d+)?$`)
)

func validateAlpha(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alpha", msg, "only supported on strings")
	}
	if !alphaRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "alpha", msg}
		}
		return locErr(lang, "alpha", name, "")
	}
	return nil
}

func validateAlphaNum(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alpha_num", msg, "only supported on strings")
	}
	if !alphaNumRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "alpha_num", msg}
		}
		return locErr(lang, "alpha_num", name, "")
	}
	return nil
}

func validateAlphaDash(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alpha_dash", msg, "only supported on strings")
	}
	if !alphaDashRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "alpha_dash", msg}
		}
		return locErr(lang, "alpha_dash", name, "")
	}
	return nil
}

func validateAlphaSpace(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alpha_space", msg, "only supported on strings")
	}
	if !alphaSpaceRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "alpha_space", msg}
		}
		return locErr(lang, "alpha_space", name, "")
	}
	return nil
}

func validateAlphaUnicode(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alpha_unicode", msg, "only supported on strings")
	}
	for _, r := range fv.String() {
		if !unicode.IsLetter(r) {
			if msg != "" {
				return ValError{name, "alpha_unicode", msg}
			}
			return locErr(lang, "alpha_unicode", name, "")
		}
	}
	return nil
}

func validateAlphaNumUnicode(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "alphanum_unicode", msg, "only supported on strings")
	}
	for _, r := range fv.String() {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if msg != "" {
				return ValError{name, "alphanum_unicode", msg}
			}
			return locErr(lang, "alphanum_unicode", name, "")
		}
	}
	return nil
}

func validateASCII(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "ascii", msg, "only supported on strings")
	}
	if !asciiRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "ascii", msg}
		}
		return locErr(lang, "ascii", name, "")
	}
	return nil
}

func validatePrintASCII(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "print_ascii", msg, "only supported on strings")
	}
	if !printASCIIRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "print_ascii", msg}
		}
		return locErr(lang, "print_ascii", name, "")
	}
	return nil
}

func validateLowercase(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "lowercase", msg, "only supported on strings")
	}
	if fv.String() != strings.ToLower(fv.String()) {
		if msg != "" {
			return ValError{name, "lowercase", msg}
		}
		return locErr(lang, "lowercase", name, "")
	}
	return nil
}

func validateUppercase(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uppercase", msg, "only supported on strings")
	}
	if fv.String() != strings.ToUpper(fv.String()) {
		if msg != "" {
			return ValError{name, "uppercase", msg}
		}
		return locErr(lang, "uppercase", name, "")
	}
	return nil
}

func validateContains(lang, name string, fv reflect.Value, substr, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "contains", msg, "only supported on strings")
	}
	if !strings.Contains(fv.String(), substr) {
		if msg != "" {
			return ValError{name, "contains", msg}
		}
		return locErr(lang, "contains", name, "")
	}
	return nil
}

func validateContainsAny(lang, name string, fv reflect.Value, chars, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "containsany", msg, "only supported on strings")
	}
	if !strings.ContainsAny(fv.String(), chars) {
		if msg != "" {
			return ValError{name, "containsany", msg}
		}
		return locErr(lang, "containsany", name, "")
	}
	return nil
}

func validateContainsRune(lang, name string, fv reflect.Value, runeStr, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "containsrune", msg, "only supported on strings")
	}
	r, _ := utf8.DecodeRuneInString(runeStr)
	if !strings.ContainsRune(fv.String(), r) {
		if msg != "" {
			return ValError{name, "containsrune", msg}
		}
		return locErr(lang, "containsrune", name, "")
	}
	return nil
}

func validateExcludes(lang, name string, fv reflect.Value, substr, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "excludes", msg, "only supported on strings")
	}
	if strings.Contains(fv.String(), substr) {
		if msg != "" {
			return ValError{name, "excludes", msg}
		}
		return locErr(lang, "excludes", name, "")
	}
	return nil
}

func validateExcludesAll(lang, name string, fv reflect.Value, chars, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "excludesall", msg, "only supported on strings")
	}
	if strings.ContainsAny(fv.String(), chars) {
		if msg != "" {
			return ValError{name, "excludesall", msg}
		}
		return locErr(lang, "excludesall", name, "")
	}
	return nil
}

func validateExcludesRune(lang, name string, fv reflect.Value, runeStr, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "excludesrune", msg, "only supported on strings")
	}
	r, _ := utf8.DecodeRuneInString(runeStr)
	if strings.ContainsRune(fv.String(), r) {
		if msg != "" {
			return ValError{name, "excludesrune", msg}
		}
		return locErr(lang, "excludesrune", name, "")
	}
	return nil
}

func validateStartsWith(lang, name string, fv reflect.Value, prefix, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "startswith", msg, "only supported on strings")
	}
	if !strings.HasPrefix(fv.String(), prefix) {
		if msg != "" {
			return ValError{name, "startswith", msg}
		}
		return locErr(lang, "startswith", name, "")
	}
	return nil
}

func validateEndsWith(lang, name string, fv reflect.Value, suffix, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "endswith", msg, "only supported on strings")
	}
	if !strings.HasSuffix(fv.String(), suffix) {
		if msg != "" {
			return ValError{name, "endswith", msg}
		}
		return locErr(lang, "endswith", name, "")
	}
	return nil
}

func validateStartsNotWith(lang, name string, fv reflect.Value, prefix, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "startsnotwith", msg, "only supported on strings")
	}
	if strings.HasPrefix(fv.String(), prefix) {
		if msg != "" {
			return ValError{name, "startsnotwith", msg}
		}
		return locErr(lang, "startsnotwith", name, "")
	}
	return nil
}

func validateEndsNotWith(lang, name string, fv reflect.Value, suffix, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "endsnotwith", msg, "only supported on strings")
	}
	if strings.HasSuffix(fv.String(), suffix) {
		if msg != "" {
			return ValError{name, "endsnotwith", msg}
		}
		return locErr(lang, "endsnotwith", name, "")
	}
	return nil
}

func validateRegex(lang, name string, fv reflect.Value, pattern, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "regex", msg, "only supported on strings")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return customErr(lang, name, "regex", msg, "invalid pattern: "+err.Error())
	}
	if !re.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "regex", msg}
		}
		return locErr(lang, "regex", name, "")
	}
	return nil
}

func validateNumber(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "number", msg, "only supported on strings")
	}
	if !numberRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "number", msg}
		}
		return locErr(lang, "number", name, "")
	}
	return nil
}
