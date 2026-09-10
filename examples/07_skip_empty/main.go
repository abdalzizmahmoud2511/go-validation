package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("SKIP EMPTY - default is skip empty fields")
	fmt.Println(repeat("=", 61))

	type User struct {
		Name    string `validate:"required;alpha"`
		Email   string `validate:"required;email"`
		Age     int    `validate:"required;min=18"`
		Address string `validate:"min=5"`
	}

	// User with some empty fields
	user := User{
		Name:    "John",
		Email:   "",
		Age:     0,
		Address: "",
	}

	fmt.Println("\n--- Default (skipEmpty=true) ---")
	fmt.Println("  Empty fields skipped, only required on non-empty:")
	errs := validation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("    - %s\n", err.Error())
	}

	fmt.Println("\n--- With skipEmpty=false ---")
	fmt.Println("  All fields validated, empty fields fail:")
	errs = validation.Validate("en", user, false)
	for _, err := range errs {
		fmt.Printf("    - %s\n", err.Error())
	}

	fmt.Println("\n--- Default, Name empty ---")
	fmt.Println("  Name is required even with skipEmpty:")
	userEmpty := User{
		Name:    "",
		Email:   "",
		Age:     0,
		Address: "abc",
	}
	errs = validation.Validate("en", userEmpty)
	for _, err := range errs {
		fmt.Printf("    - %s\n", err.Error())
	}

	// Map example
	fmt.Println("\n--- Map default (skipEmpty=true) ---")
	data := map[string]interface{}{
		"name":    "John",
		"email":   "",
		"age":     0,
		"address": "",
	}
	rules := map[string][]string{
		"name":    {"required", "alpha"},
		"email":   {"required", "email"},
		"age":     {"required", "min=18"},
		"address": {"min=5"},
	}
	errs = validation.ValidateMap("en", data, rules)
	if len(errs) == 0 {
		fmt.Println("    All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("    - %s\n", err.Error())
		}
	}

	// ValidateStructWithRules example
	fmt.Println("\n--- ValidateStructWithRules default ---")
	type Profile struct {
		Name  string
		Email string
		Age   int
	}
	profile := Profile{
		Name:  "John",
		Email: "",
		Age:   0,
	}
	structRules := map[string][]string{
		"Name":  {"required", "alpha"},
		"Email": {"required", "email"},
		"Age":   {"required", "min=18"},
	}
	errs = validation.ValidateStructWithRules("en", profile, structRules)
	if len(errs) == 0 {
		fmt.Println("    All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("    - %s\n", err.Error())
		}
	}
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
