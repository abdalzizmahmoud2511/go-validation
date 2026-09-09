package main

import (
	"reflect"
	"strconv"
	"strings"
)

func validateOneOf(lang, name string, fv reflect.Value, param, msg string) error {
	validValues := strings.Split(param, ",")
	found := false

	switch fv.Kind() {
	case reflect.String:
		for _, v := range validValues {
			if fv.String() == v {
				found = true
				break
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, v := range validValues {
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				if fv.Int() == intVal {
					found = true
					break
				}
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for _, v := range validValues {
			if uintVal, err := strconv.ParseUint(v, 10, 64); err == nil {
				if fv.Uint() == uintVal {
					found = true
					break
				}
			}
		}
	case reflect.Float32, reflect.Float64:
		for _, v := range validValues {
			if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
				if fv.Float() == floatVal {
					found = true
					break
				}
			}
		}
	case reflect.Bool:
		for _, v := range validValues {
			if boolVal, err := strconv.ParseBool(v); err == nil {
				if fv.Bool() == boolVal {
					found = true
					break
				}
			}
		}
	}

	if !found {
		if msg != "" {
			return ValError{name, "oneof", msg}
		}
		return locErr(lang, "oneof", name, "")
	}
	return nil
}

func validateIn(lang, name string, fv reflect.Value, param, msg string) error {
	return validateOneOf(lang, name, fv, param, msg)
}

func validateNotIn(lang, name string, fv reflect.Value, param, msg string) error {
	validValues := strings.Split(param, ",")
	found := false

	switch fv.Kind() {
	case reflect.String:
		for _, v := range validValues {
			if fv.String() == v {
				found = true
				break
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, v := range validValues {
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				if fv.Int() == intVal {
					found = true
					break
				}
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for _, v := range validValues {
			if uintVal, err := strconv.ParseUint(v, 10, 64); err == nil {
				if fv.Uint() == uintVal {
					found = true
					break
				}
			}
		}
	case reflect.Float32, reflect.Float64:
		for _, v := range validValues {
			if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
				if fv.Float() == floatVal {
					found = true
					break
				}
			}
		}
	case reflect.Bool:
		for _, v := range validValues {
			if boolVal, err := strconv.ParseBool(v); err == nil {
				if fv.Bool() == boolVal {
					found = true
					break
				}
			}
		}
	}

	if found {
		if msg != "" {
			return ValError{name, "not_in", msg}
		}
		return locErr(lang, "not_in", name, "")
	}
	return nil
}

func validateNoneOf(lang, name string, fv reflect.Value, param, msg string) error {
	return validateNotIn(lang, name, fv, param, msg)
}

func validateBool(lang, name string, fv reflect.Value, msg string) error {
	switch fv.Kind() {
	case reflect.Bool:
		return nil
	case reflect.String:
		_, err := strconv.ParseBool(fv.String())
		if err != nil {
			if msg != "" {
				return ValError{name, "bool", msg}
			}
			return locErr(lang, "bool", name, "")
		}
		return nil
	default:
		if msg != "" {
			return ValError{name, "bool", msg}
		}
		return locErr(lang, "bool", name, "")
	}
}

func validateUnique(lang, name string, fv reflect.Value, msg string) error {
	switch fv.Kind() {
	case reflect.Slice, reflect.Array:
		seen := make(map[interface{}]bool)
		for i := 0; i < fv.Len(); i++ {
			val := fv.Index(i).Interface()
			if seen[val] {
				if msg != "" {
					return ValError{name, "unique", msg}
				}
				return locErr(lang, "unique", name, "")
			}
			seen[val] = true
		}
		return nil
	default:
		return nil
	}
}

func validateMaxWord(lang, name string, fv reflect.Value, param, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "max_word", msg, "only supported on strings")
	}
	max, err := strconv.Atoi(param)
	if err != nil {
		return customErr(lang, name, "max_word", msg, "invalid parameter: "+param)
	}
	words := strings.Fields(fv.String())
	if len(words) > max {
		if msg != "" {
			return ValError{name, "max_word", msg}
		}
		return locErr(lang, "max_word", name, "")
	}
	return nil
}

func validateMinWord(lang, name string, fv reflect.Value, param, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "min_word", msg, "only supported on strings")
	}
	min, err := strconv.Atoi(param)
	if err != nil {
		return customErr(lang, name, "min_word", msg, "invalid parameter: "+param)
	}
	words := strings.Fields(fv.String())
	if len(words) < min {
		if msg != "" {
			return ValError{name, "min_word", msg}
		}
		return locErr(lang, "min_word", name, "")
	}
	return nil
}

func validateLenWord(lang, name string, fv reflect.Value, param, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "len_word", msg, "only supported on strings")
	}
	length, err := strconv.Atoi(param)
	if err != nil {
		return customErr(lang, name, "len_word", msg, "invalid parameter: "+param)
	}
	words := strings.Fields(fv.String())
	if len(words) != length {
		if msg != "" {
			return ValError{name, "len_word", msg}
		}
		return locErr(lang, "len_word", name, "")
	}
	return nil
}
