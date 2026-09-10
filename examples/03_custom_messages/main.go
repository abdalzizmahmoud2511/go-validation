package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("CUSTOM MESSAGES WITH Msg()")
	fmt.Println(repeat("=", 61))

	type User struct {
		Name  string `validate:"required;alpha"`
		Email string `validate:"required;email"`
		Age   int    `validate:"required;min=18"`
	}

	user := User{Name: "", Email: "invalid", Age: 15}

	// Custom messages from locale files
	rules := map[string][]string{
		"Name":  {validation.Msg("required", "custom_name_required"), validation.Msg("alpha", "custom_name_alpha")},
		"Email": {validation.Msg("email", "custom_email_invalid")},
		"Age":   {validation.Msg("min", "custom_age_min")},
	}

	fmt.Println("\n--- Arabic with Custom Locale Keys ---")
	errs := validation.ValidateWithRules("ar", user, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- English with Custom Locale Keys ---")
	errs = validation.ValidateWithRules("en", user, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
