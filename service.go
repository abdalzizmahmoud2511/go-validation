package main

import (
	"fmt"
	"reflect"
	"strings"
)

// RuleFunc is the signature for custom validation rules.
type RuleFunc func(lang, field, rule, message string, value interface{}) error

var customRules = make(map[string]RuleFunc)

// AddCustomRule registers a custom validation rule.
func AddCustomRule(ruleName string, fn RuleFunc) {
	if _, exists := customRules[ruleName]; exists {
		panic(fmt.Sprintf("govalidator: rule %q is already defined", ruleName))
	}
	customRules[ruleName] = fn
}

// RemoveCustomRule removes a previously registered custom rule.
func RemoveCustomRule(ruleName string) {
	delete(customRules, ruleName)
}

// HasCustomRule returns true if a custom rule with the given name exists.
func HasCustomRule(ruleName string) bool {
	_, exists := customRules[ruleName]
	return exists
}

// Validate validates a struct using validate tags and returns all errors.
// If skipEmpty is not provided or true, fields with zero values are skipped (except required).
func Validate(lang string, v interface{}, skipEmpty ...bool) []error {
	skip := len(skipEmpty) == 0 || skipEmpty[0]
	return validateStruct(lang, v, "", skip)
}

// Valid returns true if the struct passes all validations.
func Valid(lang string, v interface{}, skipEmpty ...bool) bool {
	return len(Validate(lang, v, skipEmpty...)) == 0
}

// ValidateMap validates a map (nested or flat) using rules with dot notation.
// If skipEmpty is not provided or true, fields with zero values are skipped (except required).
func ValidateMap(lang string, data map[string]interface{}, rules map[string][]string, skipEmpty ...bool) []error {
	skip := len(skipEmpty) == 0 || skipEmpty[0]
	var errs []error
	for field, fieldRules := range rules {
		errs = append(errs, validateMapField(lang, data, field, fieldRules, skip)...)
	}
	return errs
}

// validateMapField validates a single field (possibly in an array) from a map.
func validateMapField(lang string, data map[string]interface{}, field string, fieldRules []string, skipEmpty bool) []error {
	// Check if field contains an array with sub-field (e.g., "children.age")
	arrayField, subField := splitArrayField(field)
	if arrayField != "" && subField != "" {
		val, exists := data[arrayField]
		if !exists {
			var errs []error
			for _, rule := range fieldRules {
				key, msg := parseRuleWithMsg(rule)
				if key == "required" {
					if msg != "" {
						errs = append(errs, ValError{field, "required", msg})
					} else {
						errs = append(errs, locErr(lang, "required", field, ""))
					}
				}
			}
			return errs
		}
		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			return validateMapArray(lang, arrayField, subField, fieldRules, rv, skipEmpty)
		}
	}

	val, exists := getNestedValue(data, field)
	if !exists {
		var errs []error
		for _, rule := range fieldRules {
			key, msg := parseRuleWithMsg(rule)
			if key == "required" {
				if msg != "" {
					errs = append(errs, ValError{field, "required", msg})
				} else {
					errs = append(errs, locErr(lang, "required", field, ""))
				}
			}
		}
		return errs
	}

	fv := reflect.ValueOf(val)

	// If skipEmpty is true and value is empty, skip non-required rules
	if skipEmpty && isZero(fv) {
		hasRequired := false
		hasApplyIf := false
		for _, rule := range fieldRules {
			key, _ := parseRuleWithMsg(rule)
			if key == "required" {
				hasRequired = true
			}
			if strings.HasPrefix(rule, "applyIf|") || strings.HasPrefix(rule, "applyIfElse|") {
				hasApplyIf = true
			}
		}
		if !hasRequired && !hasApplyIf {
			return nil
		}
	}

	// Handle arrays of primitive values (e.g., emails = ["a@b.com", "invalid"])
	if (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) && len(fieldRules) > 0 {
		return validateMapArrayValues(lang, field, fieldRules, fv, skipEmpty)
	}

	var errs []error
	for _, rule := range fieldRules {
		if err := applyRule(lang, field, fv, rule, data); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// validateMapArrayValues validates each element in an array of primitive values.
func validateMapArrayValues(lang string, field string, fieldRules []string, rv reflect.Value, skipEmpty bool) []error {
	var errs []error
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		elemName := fmt.Sprintf("%s[%d]", field, i)
		// If skipEmpty and element is zero, skip
		if skipEmpty && isZero(elem) {
			continue
		}
		for _, rule := range fieldRules {
			if err := applyRule(lang, elemName, elem, rule); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// splitArrayField splits "children.age" into ("children", "age").
// Returns ("", "") if the field is not an array field.
func splitArrayField(field string) (arrayField, subField string) {
	parts := strings.SplitN(field, ".", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// validateMapArray validates each element in an array/slice for a map field.
func validateMapArray(lang string, arrayField, subField string, fieldRules []string, rv reflect.Value, skipEmpty bool) []error {
	var errs []error
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		elemFieldName := fmt.Sprintf("%s[%d].%s", arrayField, i, subField)

		if elem.Kind() == reflect.Map {
			subVal, subExists := getMapValue(elem.Interface(), subField)
			if !subExists {
				for _, rule := range fieldRules {
					key, msg := parseRuleWithMsg(rule)
					if key == "required" {
						if msg != "" {
							errs = append(errs, ValError{elemFieldName, "required", msg})
						} else {
							errs = append(errs, locErr(lang, "required", elemFieldName, ""))
						}
					}
				}
				continue
			}
			subFv := reflect.ValueOf(subVal)
			// If skipEmpty and value is zero, skip non-required rules
			if skipEmpty && isZero(subFv) {
				continue
			}
			for _, rule := range fieldRules {
				if err := applyRule(lang, elemFieldName, subFv, rule); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return errs
}

// getMapValue retrieves a value from a map (interface{}) by key.
func getMapValue(data interface{}, key string) (interface{}, bool) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}
	val, ok := m[key]
	return val, ok
}

// getNestedValue retrieves a value from a map using dot notation.
// "address.street" → data["address"]["street"]
func getNestedValue(data map[string]interface{}, key string) (interface{}, bool) {
	parts := strings.Split(key, ".")
	current := data

	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			return nil, false
		}
		if i == len(parts)-1 {
			return val, true
		}
		next, ok := val.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current = next
	}
	return nil, false
}

// ValidateWithRules validates a struct using a map of rules with dot notation.
// If skipEmpty is not provided or true, fields with zero values are skipped (except required).
func ValidateWithRules(lang string, v interface{}, rules map[string][]string, skipEmpty ...bool) []error {
	skip := len(skipEmpty) == 0 || skipEmpty[0]
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []error{fmt.Errorf("value is nil")}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return []error{fmt.Errorf("expected struct, got %s", rv.Kind())}
	}

	flat := flattenStruct(rv, "")
	var errs []error
	for field, fieldRules := range rules {
		val, ok := flat[field]
		if !ok {
			for _, rule := range fieldRules {
				key, msg := parseRuleWithMsg(rule)
				if key == "required" {
					if msg != "" {
						errs = append(errs, ValError{field, "required", msg})
					} else {
						errs = append(errs, locErr(lang, "required", field, ""))
					}
				}
			}
			continue
		}
		fv := reflect.ValueOf(val)
		// If skipEmpty and value is zero, skip non-required rules
		if skip && isZero(fv) {
			hasRequired := false
			hasApplyIf := false
			for _, rule := range fieldRules {
				key, _ := parseRuleWithMsg(rule)
				if key == "required" {
					hasRequired = true
				}
				if strings.HasPrefix(rule, "applyIf|") || strings.HasPrefix(rule, "applyIfElse|") {
					hasApplyIf = true
				}
			}
			if !hasRequired && !hasApplyIf {
				continue
			}
		}
		for _, rule := range fieldRules {
			if err := applyRule(lang, field, fv, rule, v); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// ValidateStructWithRules validates a struct using a map of rules with dot notation.
// Supports nested structs, arrays of structs, and arrays of primitives.
//
// Rules use dot notation for nested fields:
//
//	rules := map[string][]string{
//	    "Name":              {"required", "alpha"},
//	    "Address.Street":    {"required", "min=2"},
//	    "Address.City":      {"required"},
//	    "Children.Name":     {"required"},
//	    "Children.Age":      {"required", "min=1"},
//	    "Emails":            {"required", "email"},
//	}
// If skipEmpty is not provided or true, fields with zero values are skipped (except required).
func ValidateStructWithRules(lang string, v interface{}, rules map[string][]string, skipEmpty ...bool) []error {
	skip := len(skipEmpty) == 0 || skipEmpty[0]
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []error{fmt.Errorf("value is nil")}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return []error{fmt.Errorf("expected struct, got %s", rv.Kind())}
	}

	var errs []error
	for field, fieldRules := range rules {
		errs = append(errs, validateStructFieldRules(lang, rv, field, fieldRules, skip, rv.Interface())...)
	}
	return errs
}

// validateStructFieldRules validates a struct field (possibly nested/array) with rules.
func validateStructFieldRules(lang string, rv reflect.Value, field string, fieldRules []string, skipEmpty bool, parentData ...interface{}) []error {
	// Try to find the field value
	fv, exists := getStructFieldValue(rv, field)
	if !exists {
		var errs []error
		for _, rule := range fieldRules {
			key, msg := parseRuleWithMsg(rule)
			if key == "required" {
				if msg != "" {
					errs = append(errs, ValError{field, "required", msg})
				} else {
					errs = append(errs, locErr(lang, "required", field, ""))
				}
			}
		}
		return errs
	}

	// Handle slices/arrays
	if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
		return validateStructArrayField(lang, fv, field, fieldRules, skipEmpty)
	}

	// If skipEmpty and value is zero, skip non-required rules
	if skipEmpty && isZero(fv) {
		hasRequired := false
		hasApplyIf := false
		for _, rule := range fieldRules {
			key, _ := parseRuleWithMsg(rule)
			if key == "required" {
				hasRequired = true
			}
			if strings.HasPrefix(rule, "applyIf|") || strings.HasPrefix(rule, "applyIfElse|") {
				hasApplyIf = true
			}
		}
		if !hasRequired && !hasApplyIf {
			return nil
		}
	}

	// Single value
	var errs []error
	for _, rule := range fieldRules {
		if err := applyRule(lang, field, fv, rule, parentData...); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// getStructFieldValue retrieves a field value from a struct using dot notation.
// "Address.Street" → struct.Address.Street
// "Children.Name" → struct.Children (slice), returns the slice
func getStructFieldValue(rv reflect.Value, field string) (reflect.Value, bool) {
	parts := strings.SplitN(field, ".", 2)
	fieldName := parts[0]

	// Find field by name
	var fv reflect.Value
	found := false
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Name == fieldName {
			fv = rv.Field(i)
			found = true
			break
		}
	}

	if !found {
		return reflect.Value{}, false
	}

	// Handle pointers
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return reflect.Value{}, false
		}
		fv = fv.Elem()
	}

	// If there are more parts, recurse into the value
	if len(parts) > 1 {
		// If it's a slice/array, return it as-is for array handling
		if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
			return fv, true
		}
		// If it's a struct, recurse
		if fv.Kind() == reflect.Struct {
			return getStructFieldValue(fv, parts[1])
		}
		return reflect.Value{}, false
	}

	return fv, true
}

// validateStructArrayField validates each element in a struct array field.
func validateStructArrayField(lang string, sv reflect.Value, field string, fieldRules []string, skipEmpty bool) []error {
	var errs []error
	for i := 0; i < sv.Len(); i++ {
		elem := sv.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		elemName := fmt.Sprintf("%s[%d]", field, i)

		// If element is a struct, check sub-rules (e.g., "Children.Name")
		if elem.Kind() == reflect.Struct {
			for _, rule := range fieldRules {
				ruleField, _ := parseRule(rule)
				// Check if this rule is for a sub-field (e.g., "Children.Name")
				if strings.Contains(ruleField, ".") {
					subParts := strings.SplitN(ruleField, ".", 2)
					if subParts[0] == field {
						subField := subParts[1]
						subFv, subExists := getStructFieldValue(elem, subField)
						if !subExists {
							key, msg := parseRuleWithMsg(rule)
							if key == "required" {
								if msg != "" {
									errs = append(errs, ValError{elemName+"."+subField, "required", msg})
								} else {
									errs = append(errs, locErr(lang, "required", elemName+"."+subField, ""))
								}
							}
							continue
						}
						// If skipEmpty and value is zero, skip non-required rules
						if skipEmpty && isZero(subFv) {
							continue
						}
						for _, r := range fieldRules {
							rField, _ := parseRule(r)
							if rField == ruleField {
								if err := applyRule(lang, elemName+"."+subField, subFv, r); err != nil {
									errs = append(errs, err)
								}
							}
						}
					}
				}
			}
		} else {
			// If skipEmpty and element is zero, skip
			if skipEmpty && isZero(elem) {
				continue
			}
			// Primitive element, apply rules directly
			for _, rule := range fieldRules {
				// Only apply rules without dots (e.g., "Emails" not "Children.Name")
				ruleField, _ := parseRule(rule)
				if !strings.Contains(ruleField, ".") {
					if err := applyRule(lang, elemName, elem, rule); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}
	}
	return errs
}

// validateStruct recursively validates a struct using validate tags.
// Supports slices of structs (e.g., []Children).
func validateStruct(lang string, v interface{}, prefix string, skipEmpty bool) []error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []error{fmt.Errorf("value is nil")}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return []error{fmt.Errorf("expected struct, got %s", rv.Kind())}
	}

	rt := rv.Type()
	var errs []error
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		rules := field.Tag.Get("validate")
		fieldName := field.Name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		fv := rv.Field(i)

		// Handle slices/arrays of structs
		if (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) && rules != "" {
			errs = append(errs, validateStructSlice(lang, fv, fieldName, rules, skipEmpty)...)
			continue
		}

		// Recurse into nested structs without tags
		if isNestedStruct(fv) && rules == "" {
			if fv.CanAddr() {
				errs = append(errs, validateStruct(lang, fv.Addr().Interface(), fieldName, skipEmpty)...)
			} else {
				errs = append(errs, validateStruct(lang, fv.Interface(), fieldName, skipEmpty)...)
			}
			continue
		}

		if rules == "" {
			continue
		}

		// Apply rules to this field
		for _, rule := range strings.Split(rules, ";") {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}
			ruleKey, _, _ := parseRuleFull(rule)
			// If skipEmpty and value is zero, skip non-required rules
			if skipEmpty && isZero(fv) && ruleKey != "required" {
				continue
			}
			if err := applyRule(lang, fieldName, fv, rule, rv.Interface()); err != nil {
				errs = append(errs, err)
			}
		}

		// Recurse into nested structs that have rules
		if isNestedStruct(fv) && hasRule(rules, "required") {
			if fv.CanAddr() {
				errs = append(errs, validateStruct(lang, fv.Addr().Interface(), fieldName, skipEmpty)...)
			} else {
				errs = append(errs, validateStruct(lang, fv.Interface(), fieldName, skipEmpty)...)
			}
		}
	}
	return errs
}

// validateStructSlice validates each element in a slice of structs.
func validateStructSlice(lang string, sv reflect.Value, prefix string, rules string, skipEmpty bool) []error {
	var errs []error
	for i := 0; i < sv.Len(); i++ {
		elem := sv.Index(i)
		elemName := fmt.Sprintf("%s[%d]", prefix, i)

		// If skipEmpty and element is zero, skip
		if skipEmpty && isZero(elem) {
			continue
		}

		// Apply rules to the slice element itself (e.g., required)
		for _, rule := range strings.Split(rules, ";") {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}
			if err := applyRule(lang, elemName, elem, rule); err != nil {
				errs = append(errs, err)
			}
		}

		// Recurse into struct elements to validate their fields
		if elem.Kind() == reflect.Struct {
			errs = append(errs, validateStruct(lang, elem.Interface(), elemName, skipEmpty)...)
		}
	}
	return errs
}

// applyRule applies a single rule to a field value.
// If the value is a slice/array and the rule is not "required", the rule is applied to each element.
func applyRule(lang, name string, fv reflect.Value, rule string, data ...interface{}) error {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return nil
	}

	// Handle applyIf rule
	if strings.HasPrefix(rule, "applyIf|") {
		return applyApplyIfRule(lang, name, fv, rule, data...)
	}

	// Handle applyIfElse rule
	if strings.HasPrefix(rule, "applyIfElse|") {
		return applyApplyIfElseRule(lang, name, fv, rule, data...)
	}

	key, arg, msg := parseRuleFull(rule)
	msg = resolveLocaleMsg(lang, name, key, msg)

	// Handle slices/arrays: apply rule to each element (except "required")
	if (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) && key != "required" {
		var errs []error
		for i := 0; i < fv.Len(); i++ {
			elem := fv.Index(i)
			if elem.Kind() == reflect.Interface {
				elem = elem.Elem()
			}
			elemName := fmt.Sprintf("%s[%d]", name, i)
			if err := applyRule(lang, elemName, elem, rule); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errs[0] // Return first error
		}
		return nil
	}

	switch key {
	case "required":
		if isZero(fv) {
			if msg != "" {
				return ValError{name, key, msg}
			}
			return locErr(lang, "required", name, "")
		}

	// String rules
	case "alpha":
		return validateAlpha(lang, name, fv, msg)
	case "alpha_num":
		return validateAlphaNum(lang, name, fv, msg)
	case "alpha_dash":
		return validateAlphaDash(lang, name, fv, msg)
	case "alpha_space":
		return validateAlphaSpace(lang, name, fv, msg)
	case "alpha_unicode":
		return validateAlphaUnicode(lang, name, fv, msg)
	case "alphanum_unicode":
		return validateAlphaNumUnicode(lang, name, fv, msg)
	case "ascii":
		return validateASCII(lang, name, fv, msg)
	case "print_ascii":
		return validatePrintASCII(lang, name, fv, msg)
	case "lowercase":
		return validateLowercase(lang, name, fv, msg)
	case "uppercase":
		return validateUppercase(lang, name, fv, msg)
	case "contains":
		return validateContains(lang, name, fv, arg, msg)
	case "containsany":
		return validateContainsAny(lang, name, fv, arg, msg)
	case "containsrune":
		return validateContainsRune(lang, name, fv, arg, msg)
	case "excludes":
		return validateExcludes(lang, name, fv, arg, msg)
	case "excludesall":
		return validateExcludesAll(lang, name, fv, arg, msg)
	case "excludesrune":
		return validateExcludesRune(lang, name, fv, arg, msg)
	case "startswith":
		return validateStartsWith(lang, name, fv, arg, msg)
	case "endswith":
		return validateEndsWith(lang, name, fv, arg, msg)
	case "startsnotwith":
		return validateStartsNotWith(lang, name, fv, arg, msg)
	case "endsnotwith":
		return validateEndsNotWith(lang, name, fv, arg, msg)
	case "regex":
		return validateRegex(lang, name, fv, arg, msg)
	case "number":
		return validateNumber(lang, name, fv, msg)

	// Numeric rules
	case "numeric":
		return validateNumeric(lang, name, fv, msg)
	case "min":
		return validateMin(lang, name, fv, arg, msg, "")
	case "max":
		return validateMax(lang, name, fv, arg, msg, "")
	case "len":
		return validateLen(lang, name, fv, arg, msg, "")
	case "between":
		return validateBetween(lang, name, fv, arg, msg, "")
	case "digits":
		return validateDigits(lang, name, fv, arg, msg)
	case "digits_between":
		return validateDigitsBetween(lang, name, fv, arg, msg)
	case "gt":
		return validateGt(lang, name, fv, arg, msg, "")
	case "gte":
		return validateGte(lang, name, fv, arg, msg, "")
	case "lt":
		return validateLt(lang, name, fv, arg, msg, "")
	case "lte":
		return validateLte(lang, name, fv, arg, msg, "")
	case "eq":
		return validateEq(lang, name, fv, arg, msg, "")
	case "ne":
		return validateNe(lang, name, fv, arg, msg, "")

	// Format rules
	case "email":
		return validateEmail(lang, name, fv, msg)
	case "url":
		return validateURL(lang, name, fv, msg)
	case "uri":
		return validateURI(lang, name, fv, msg)
	case "uuid":
		return validateUUID(lang, name, fv, msg)
	case "uuid3":
		return validateUUID3(lang, name, fv, msg)
	case "uuid4":
		return validateUUID4(lang, name, fv, msg)
	case "uuid5":
		return validateUUID5(lang, name, fv, msg)
	case "date":
		return validateDate(lang, name, fv, msg)
	case "date_dmy":
		return validateDateDDMMYY(lang, name, fv, msg)
	case "datetime":
		return validateDateTime(lang, name, fv, arg, msg)
	case "timezone":
		return validateTimezone(lang, name, fv, msg)
	case "ip":
		return validateIP(lang, name, fv, msg)
	case "ip_v4":
		return validateIPv4(lang, name, fv, msg)
	case "ip_v6":
		return validateIPv6(lang, name, fv, msg)
	case "mac_address":
		return validateMAC(lang, name, fv, msg)
	case "latitude":
		return validateLatitude(lang, name, fv, msg)
	case "longitude":
		return validateLongitude(lang, name, fv, msg)
	case "lat":
		return validateLatitudeDecimal(lang, name, fv, msg)
	case "lon":
		return validateLongitudeDecimal(lang, name, fv, msg)
	case "coordinate":
		return validateCoordinate(lang, name, fv, msg)
	case "hexcolor":
		return validateHexColor(lang, name, fv, msg)
	case "css_color":
		return validateCSSColor(lang, name, fv, msg)
	case "rgb":
		return validateRGB(lang, name, fv, msg)
	case "rgba":
		return validateRGBA(lang, name, fv, msg)
	case "hsl":
		return validateHSL(lang, name, fv, msg)
	case "hsla":
		return validateHSLA(lang, name, fv, msg)
	case "json":
		return validateJSON(lang, name, fv, msg)
	case "credit_card":
		return validateCreditCard(lang, name, fv, msg)
	case "base64":
		return validateBase64(lang, name, fv, msg)
	case "hexadecimal":
		return validateHexadecimal(lang, name, fv, msg)
	case "isbn":
		return validateISBN(lang, name, fv, msg)
	case "isbn10":
		return validateISBN10(lang, name, fv, msg)
	case "isbn13":
		return validateISBN13(lang, name, fv, msg)
	case "issn":
		return validateISSN(lang, name, fv, msg)
	case "jwt":
		return validateJWT(lang, name, fv, msg)
	case "semver":
		return validateSemver(lang, name, fv, msg)
	case "ssn":
		return validateSSN(lang, name, fv, msg)

	// Enum rules
	case "oneof":
		return validateOneOf(lang, name, fv, arg, msg)
	case "in":
		return validateOneOf(lang, name, fv, arg, msg)
	case "not_in":
		return validateNotIn(lang, name, fv, arg, msg)
	case "noneof":
		return validateNoneOf(lang, name, fv, arg, msg)
	case "bool":
		return validateBool(lang, name, fv, msg)
	case "unique":
		return validateUnique(lang, name, fv, msg)
	case "max_word":
		return validateMaxWord(lang, name, fv, arg, msg)
	case "min_word":
		return validateMinWord(lang, name, fv, arg, msg)
	case "len_word":
		return validateLenWord(lang, name, fv, arg, msg)
	}

	// Check custom rules
	if fn, ok := customRules[key]; ok {
		return fn(lang, name, rule, msg, fv.Interface())
	}
	for k, fn := range customRules {
		if strings.HasPrefix(rule, k+"=") || strings.HasPrefix(rule, k+":") {
			return fn(lang, name, rule, msg, fv.Interface())
		}
	}

	return nil
}

// applyApplyIfRule handles the "applyIf|rule|conditions" rule
func applyApplyIfRule(lang, name string, fv reflect.Value, rule string, data ...interface{}) error {
	resolvedRule, conditions, ok := parseApplyIfRule(rule)
	if !ok || len(conditions) == 0 {
		return nil
	}

	// Need parent data to evaluate conditions
	var parent interface{}
	if len(data) > 0 {
		parent = data[0]
	}
	if parent == nil {
		return nil
	}

	// Evaluate all conditions
	for field, fieldRules := range conditions {
		if !checkConditions(parent, field, fieldRules) {
			return nil // Condition not met, skip rule
		}
	}

	// All conditions passed, apply the resolved rule
	return applyRule(lang, name, fv, resolvedRule, data...)
}

// applyApplyIfElseRule handles the "applyIfElse|thenRule|elseRule|conditions" rule
func applyApplyIfElseRule(lang, name string, fv reflect.Value, rule string, data ...interface{}) error {
	thenRule, elseRule, conditions, ok := parseApplyIfElseRule(rule)
	if !ok {
		return nil
	}

	// Need parent data to evaluate conditions
	var parent interface{}
	if len(data) > 0 {
		parent = data[0]
	}
	if parent == nil {
		return applyRule(lang, name, fv, thenRule, data...)
	}

	// Evaluate all conditions
	allPassed := true
	for field, fieldRules := range conditions {
		if !checkConditions(parent, field, fieldRules) {
			allPassed = false
			break
		}
	}

	// Apply appropriate rule
	if allPassed {
		return applyRule(lang, name, fv, thenRule, data...)
	}
	return applyRule(lang, name, fv, elseRule, data...)
}

// checkConditions checks if all conditions pass for a field against the parent data
func checkConditions(parent interface{}, field string, rules []string) bool {
	rv := reflect.ValueOf(parent)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}

	var fv reflect.Value

	switch rv.Kind() {
	case reflect.Struct:
		fv = getCheckFieldValue(rv, field)
	case reflect.Map:
		fv = getCheckMapValue(rv, field)
	default:
		return false
	}

	if !fv.IsValid() {
		return false
	}

	// Check all rules against the field
	for _, rule := range rules {
		if !checkRuleAgainstField(fv, rule) {
			return false
		}
	}

	return true
}

// getCheckMapValue retrieves a value from a map using dot notation
func getCheckMapValue(rv reflect.Value, field string) reflect.Value {
	parts := strings.SplitN(field, ".", 2)
	currentField := parts[0]

	fv := rv.MapIndex(reflect.ValueOf(currentField))
	if !fv.IsValid() {
		return reflect.Value{}
	}
	if fv.Kind() == reflect.Interface {
		fv = fv.Elem()
	}

	// Handle nested maps
	if len(parts) > 1 && fv.Kind() == reflect.Map {
		return getCheckMapValue(fv, parts[1])
	}

	return fv
}

func getCheckFieldValue(rv reflect.Value, fieldName string) reflect.Value {
	// Handle dot notation for nested structs
	parts := strings.SplitN(fieldName, ".", 2)
	currentField := parts[0]

	for i := 0; i < rv.NumField(); i++ {
		if strings.EqualFold(rv.Type().Field(i).Name, currentField) {
			fv := rv.Field(i)
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					return reflect.Value{}
				}
				fv = fv.Elem()
			}
			if len(parts) > 1 && fv.Kind() == reflect.Struct {
				return getCheckFieldValue(fv, parts[1])
			}
			return fv
		}
	}
	return reflect.Value{}
}

func checkRuleAgainstField(fv reflect.Value, rule string) bool {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return true
	}

	// Handle pointers
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return rule == "required"
		}
		fv = fv.Elem()
	}

	// Parse rule
	parts := strings.SplitN(rule, "=", 2)
	key := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	// Also check for : separator
	if !strings.Contains(rule, "=") {
		parts = strings.SplitN(rule, ":", 2)
		key = parts[0]
		if len(parts) > 1 {
			arg = parts[1]
		}
	}

	switch key {
	case "required":
		return !isZero(fv)
	case "alpha":
		if fv.Kind() != reflect.String {
			return false
		}
		return alphaRe.MatchString(fv.String())
	case "alpha_num":
		if fv.Kind() != reflect.String {
			return false
		}
		return alphaNumRe.MatchString(fv.String())
	case "email":
		if fv.Kind() != reflect.String {
			return false
		}
		return emailRe.MatchString(fv.String())
	case "min":
		return checkNumericComparison(fv, arg, ">=")
	case "max":
		return checkNumericComparison(fv, arg, "<=")
	case "len":
		return checkLengthComparison(fv, arg, "==")
	case "eq":
		return checkStringComparison(fv, arg, "==")
	case "ne":
		return checkStringComparison(fv, arg, "!=")
	case "gt":
		return checkNumericComparison(fv, arg, ">")
	case "lt":
		return checkNumericComparison(fv, arg, "<")
	case "gte":
		return checkNumericComparison(fv, arg, ">=")
	case "lte":
		return checkNumericComparison(fv, arg, "<=")
	case "oneof", "in":
		validValues := strings.Split(arg, ",")
		s := fvToString(fv)
		for _, v := range validValues {
			if s == strings.TrimSpace(v) {
				return true
			}
		}
		return false
	case "not_in":
		validValues := strings.Split(arg, ",")
		s := fvToString(fv)
		for _, v := range validValues {
			if s == strings.TrimSpace(v) {
				return false
			}
		}
		return true
	case "empty":
		return isZero(fv)
	case "not_empty":
		return !isZero(fv)
	}

	return true
}

func checkNumericComparison(fv reflect.Value, arg string, op string) bool {
	argFloat, err := parseFloat(arg)
	if err != nil {
		return false
	}

	var valFloat float64
	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valFloat = float64(fv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valFloat = float64(fv.Uint())
	case reflect.Float32, reflect.Float64:
		valFloat = fv.Float()
	default:
		return false
	}

	switch op {
	case ">":
		return valFloat > argFloat
	case ">=":
		return valFloat >= argFloat
	case "<":
		return valFloat < argFloat
	case "<=":
		return valFloat <= argFloat
	case "==":
		return valFloat == argFloat
	case "!=":
		return valFloat != argFloat
	}
	return false
}

func checkLengthComparison(fv reflect.Value, arg string, op string) bool {
	argFloat, err := parseFloat(arg)
	if err != nil {
		return false
	}

	var length int
	switch fv.Kind() {
	case reflect.String:
		length = fv.Len()
	case reflect.Slice, reflect.Map, reflect.Array:
		length = fv.Len()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// For numbers, length is the number itself
		length = int(fv.Int())
	default:
		return false
	}

	lengthFloat := float64(length)
	switch op {
	case "==":
		return lengthFloat == argFloat
	case "!=":
		return lengthFloat != argFloat
	case ">":
		return lengthFloat > argFloat
	case ">=":
		return lengthFloat >= argFloat
	case "<":
		return lengthFloat < argFloat
	case "<=":
		return lengthFloat <= argFloat
	}
	return false
}

func checkStringComparison(fv reflect.Value, arg string, op string) bool {
	if fv.Kind() != reflect.String {
		return false
	}
	s := fv.String()
	switch op {
	case "==":
		return s == arg
	case "!=":
		return s != arg
	}
	return false
}
