package govalidation

import (
	"reflect"
	"strconv"
	"strings"
)

func validateNumeric(lang, name string, fv reflect.Value, msg string) error {
	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return nil
	case reflect.String:
		_, err := strconv.ParseFloat(fv.String(), 64)
		if err != nil {
			if msg != "" {
				return ValError{name, "numeric", msg}
			}
			return locErr(lang, "numeric", name, "")
		}
		return nil
	default:
		if msg != "" {
			return ValError{name, "numeric", msg}
		}
		return locErr(lang, "numeric", name, "")
	}
}

func validateMin(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "min", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) < parsedFloat {
			return sizeErr(lang, name, "min", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) < parsedFloat {
			return sizeErr(lang, name, "min", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() < parsedFloat {
			return sizeErr(lang, name, "min", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) < parsedFloat {
			return sizeErr(lang, name, "min", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) < parsedFloat {
			return sizeErr(lang, name, "min", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "min", msg, "unsupported type")
	}
	return nil
}

func validateMax(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "max", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) > parsedFloat {
			return sizeErr(lang, name, "max", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) > parsedFloat {
			return sizeErr(lang, name, "max", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() > parsedFloat {
			return sizeErr(lang, name, "max", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) > parsedFloat {
			return sizeErr(lang, name, "max", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) > parsedFloat {
			return sizeErr(lang, name, "max", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "max", msg, "unsupported type")
	}
	return nil
}

func validateLen(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "len", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) != parsedFloat {
			return sizeErr(lang, name, "len", msgType, msg, "", param, "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) != parsedFloat {
			return sizeErr(lang, name, "len", msgType, msg, "", param, "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() != parsedFloat {
			return sizeErr(lang, name, "len", msgType, msg, "", param, "")
		}
	case reflect.String:
		if float64(fv.Len()) != parsedFloat {
			return sizeErr(lang, name, "len", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) != parsedFloat {
			return sizeErr(lang, name, "len", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "len", msg, "unsupported type")
	}
	return nil
}

func validateBetween(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parts := splitRange(param)
	if len(parts) != 2 {
		return customErr(lang, name, "between", msg, "invalid parameters: "+param)
	}
	minStr, maxStr := parts[0], parts[1]
	min, err1 := strconv.ParseFloat(minStr, 64)
	max, err2 := strconv.ParseFloat(maxStr, 64)
	if err1 != nil || err2 != nil {
		return customErr(lang, name, "between", msg, "invalid parameters: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := float64(fv.Int())
		if val < min || val > max {
			return sizeErr(lang, name, "between", msgType, msg, minStr, maxStr, "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := float64(fv.Uint())
		if val < min || val > max {
			return sizeErr(lang, name, "between", msgType, msg, minStr, maxStr, "")
		}
	case reflect.Float32, reflect.Float64:
		val := fv.Float()
		if val < min || val > max {
			return sizeErr(lang, name, "between", msgType, msg, minStr, maxStr, "")
		}
	case reflect.String:
		length := float64(fv.Len())
		if length < min || length > max {
			return sizeErr(lang, name, "between", msgType, msg, "", minStr, maxStr)
		}
	case reflect.Slice, reflect.Map:
		length := float64(fv.Len())
		if length < min || length > max {
			return sizeErr(lang, name, "between", msgType, msg, "", minStr, maxStr)
		}
	default:
		return customErr(lang, name, "between", msg, "unsupported type")
	}
	return nil
}

func validateDigits(lang, name string, fv reflect.Value, param, msg string) error {
	parsedInt, err := strconv.Atoi(param)
	if err != nil {
		return customErr(lang, name, "digits", msg, "invalid parameter: "+param)
	}

	s := strings.TrimSpace(fvToString(fv))

	dotIndex := strings.Index(s, ".")
	if dotIndex != -1 {
		decimalPart := s[dotIndex+1:]
		intPart := s[:dotIndex]
		if len(decimalPart) != parsedInt || len(intPart) == 0 {
			if msg != "" {
				return ValError{name, "digits", msg}
			}
			return locErr(lang, "digits", name, "")
		}
	} else {
		if len(s) != parsedInt {
			if msg != "" {
				return ValError{name, "digits", msg}
			}
			return locErr(lang, "digits", name, "")
		}
	}
	return nil
}

func validateDigitsBetween(lang, name string, fv reflect.Value, param, msg string) error {
	parts := splitRange(param)
	if len(parts) != 2 {
		return customErr(lang, name, "digits_between", msg, "invalid parameters: "+param)
	}
	minStr, maxStr := parts[0], parts[1]
	min, err1 := strconv.Atoi(minStr)
	max, err2 := strconv.Atoi(maxStr)
	if err1 != nil || err2 != nil {
		return customErr(lang, name, "digits_between", msg, "invalid parameters: "+param)
	}

	s := strings.TrimSpace(fvToString(fv))
	dotParts := strings.Split(s, ".")
	count := 0
	if len(dotParts) == 2 {
		count = len(dotParts[1])
	} else {
		count = len(dotParts[0])
	}

	if count < min || count > max {
		if msg != "" {
			return ValError{name, "digits_between", msg}
		}
		return locErr(lang, "digits_between", name, "")
	}
	return nil
}

func validateGt(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "gt", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) <= parsedFloat {
			return sizeErr(lang, name, "gt", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) <= parsedFloat {
			return sizeErr(lang, name, "gt", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() <= parsedFloat {
			return sizeErr(lang, name, "gt", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) <= parsedFloat {
			return sizeErr(lang, name, "gt", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) <= parsedFloat {
			return sizeErr(lang, name, "gt", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "gt", msg, "unsupported type")
	}
	return nil
}

func validateGte(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "gte", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) < parsedFloat {
			return sizeErr(lang, name, "gte", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) < parsedFloat {
			return sizeErr(lang, name, "gte", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() < parsedFloat {
			return sizeErr(lang, name, "gte", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) < parsedFloat {
			return sizeErr(lang, name, "gte", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) < parsedFloat {
			return sizeErr(lang, name, "gte", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "gte", msg, "unsupported type")
	}
	return nil
}

func validateLt(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "lt", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) >= parsedFloat {
			return sizeErr(lang, name, "lt", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) >= parsedFloat {
			return sizeErr(lang, name, "lt", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() >= parsedFloat {
			return sizeErr(lang, name, "lt", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) >= parsedFloat {
			return sizeErr(lang, name, "lt", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) >= parsedFloat {
			return sizeErr(lang, name, "lt", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "lt", msg, "unsupported type")
	}
	return nil
}

func validateLte(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "lte", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) > parsedFloat {
			return sizeErr(lang, name, "lte", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) > parsedFloat {
			return sizeErr(lang, name, "lte", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() > parsedFloat {
			return sizeErr(lang, name, "lte", msgType, msg, param, "", "")
		}
	case reflect.String:
		if float64(fv.Len()) > parsedFloat {
			return sizeErr(lang, name, "lte", msgType, msg, "", param, "")
		}
	case reflect.Slice, reflect.Map:
		if float64(fv.Len()) > parsedFloat {
			return sizeErr(lang, name, "lte", msgType, msg, "", param, "")
		}
	default:
		return customErr(lang, name, "lte", msg, "unsupported type")
	}
	return nil
}

func validateEq(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "eq", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) != parsedFloat {
			return sizeErr(lang, name, "eq", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) != parsedFloat {
			return sizeErr(lang, name, "eq", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() != parsedFloat {
			return sizeErr(lang, name, "eq", msgType, msg, param, "", "")
		}
	case reflect.String:
		if fv.String() != param {
			return sizeErr(lang, name, "eq", msgType, msg, param, "", "")
		}
	default:
		return customErr(lang, name, "eq", msg, "unsupported type")
	}
	return nil
}

func validateNe(lang, name string, fv reflect.Value, param, msg, msgType string) error {
	parsedFloat, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return customErr(lang, name, "ne", msg, "invalid parameter: "+param)
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(fv.Int()) == parsedFloat {
			return sizeErr(lang, name, "ne", msgType, msg, param, "", "")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(fv.Uint()) == parsedFloat {
			return sizeErr(lang, name, "ne", msgType, msg, param, "", "")
		}
	case reflect.Float32, reflect.Float64:
		if fv.Float() == parsedFloat {
			return sizeErr(lang, name, "ne", msgType, msg, param, "", "")
		}
	case reflect.String:
		if fv.String() == param {
			return sizeErr(lang, name, "ne", msgType, msg, param, "", "")
		}
	default:
		return customErr(lang, name, "ne", msg, "unsupported type")
	}
	return nil
}

func sizeErr(lang, field, rule, msgType, customMsg, val, min, max string) error {
	if customMsg != "" {
		return ValError{field, rule, customMsg}
	}

	msgKey := rule
	if msgType == "string" {
		msgKey = rule + "_string"
	} else if msgType == "number" {
		msgKey = rule + "_number"
	}

	if min == "" && max == "" {
		return locErr(lang, msgKey, field, val)
	} else if min == "" {
		return locErr(lang, msgKey, field, max)
	} else if max == "" {
		return locErr(lang, msgKey, field, min)
	}
	return locErr(lang, msgKey, field, min, max)
}
