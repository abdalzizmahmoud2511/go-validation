package govalidation

import (
	"fmt"
	"strings"
)

// ApplyIf creates a conditional rule that applies a rule only if conditions pass.
// Similar to Msg() but for conditional validation.
//
// Parameters:
//   - rule: the rule to apply if conditions pass (e.g., "required", "min=5")
//   - conditions: map[fieldName][]rulesToCheck - if ALL rules pass for the field, apply the rule
//
// Returns: a rule string that can be used in struct tags or rules maps
//
// Usage:
//
//	// In rules map:
//	rules := map[string][]string{
//	    "ZipCode": {ApplyIf("required", map[string][]string{
//	        "Country": {"required", "alpha"},
//	    })},
//	}
//	// This generates: "applyIf|required|Country,required|Country,alpha"
//
//	// In struct tags (not directly, but the generated string works):
//	// ZipCode string `validate:"applyIf|required|Country,required|Country,alpha"`
func ApplyIf(rule string, conditions map[string][]string) string {
	if len(conditions) == 0 {
		return rule
	}

	// Build condition parts: field,rules;field,rules
	var conditionParts []string
	for field, fieldRules := range conditions {
		conditionParts = append(conditionParts, field+":"+strings.Join(fieldRules, ","))
	}

	return fmt.Sprintf("applyIf|%s|%s", rule, strings.Join(conditionParts, ";"))
}

// ApplyIfElse creates a conditional rule with then/else branches.
// If conditions pass, applies thenRule; otherwise applies elseRule.
//
// Usage:
//
//	rules := map[string][]string{
//	    "ZipCode": {ApplyIfElse("required", "nullable", map[string][]string{
//	        "Country": {"required", "alpha"},
//	    })},
//	}
func ApplyIfElse(thenRule, elseRule string, conditions map[string][]string) string {
	if len(conditions) == 0 {
		return thenRule
	}

	var conditionParts []string
	for field, fieldRules := range conditions {
		conditionParts = append(conditionParts, field+":"+strings.Join(fieldRules, ","))
	}

	return fmt.Sprintf("applyIfElse|%s|%s|%s", thenRule, elseRule, strings.Join(conditionParts, ";"))
}

// parseApplyIfRule parses a rule string like "applyIf|required|Country,required|Country,alpha"
// Returns: rule to apply, conditions map
func parseApplyIfRule(ruleStr string) (string, map[string][]string, bool) {
	if !strings.HasPrefix(ruleStr, "applyIf|") && !strings.HasPrefix(ruleStr, "applyIfElse|") {
		return "", nil, false
	}

	// Split by |
	parts := strings.SplitN(ruleStr, "|", 3)
	if len(parts) < 3 {
		return "", nil, false
	}

	// parts[0] = "applyIf" or "applyIfElse"
	// parts[1] = rule to apply
	// parts[2] = conditions string
	rule := parts[1]
	conditionsStr := parts[2]

	// Parse conditions: "field:rule1,rule2;field2:rule3"
	conditions := make(map[string][]string)
	fieldGroups := strings.Split(conditionsStr, ";")

	for _, group := range fieldGroups {
		fieldParts := strings.SplitN(group, ":", 2)
		if len(fieldParts) != 2 {
			continue
		}
		field := fieldParts[0]
		rules := strings.Split(fieldParts[1], ",")
		conditions[field] = rules
	}

	return rule, conditions, true
}

// parseApplyIfElseRule parses a rule string like "applyIfElse|required|nullable|Country,required"
// Returns: thenRule, elseRule, conditions map
func parseApplyIfElseRule(ruleStr string) (string, string, map[string][]string, bool) {
	if !strings.HasPrefix(ruleStr, "applyIfElse|") {
		return "", "", nil, false
	}

	// Split by |
	parts := strings.SplitN(ruleStr, "|", 4)
	if len(parts) < 4 {
		return "", "", nil, false
	}

	// parts[0] = "applyIfElse"
	// parts[1] = thenRule
	// parts[2] = elseRule
	// parts[3] = conditions string
	thenRule := parts[1]
	elseRule := parts[2]
	conditionsStr := parts[3]

	// Parse conditions: "field:rule1,rule2;field2:rule3"
	conditions := make(map[string][]string)
	fieldGroups := strings.Split(conditionsStr, ";")

	for _, group := range fieldGroups {
		fieldParts := strings.SplitN(group, ":", 2)
		if len(fieldParts) != 2 {
			continue
		}
		field := fieldParts[0]
		rules := strings.Split(fieldParts[1], ",")
		conditions[field] = rules
	}

	return thenRule, elseRule, conditions, true
}
