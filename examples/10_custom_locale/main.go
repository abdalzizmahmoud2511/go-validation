package main

import (
	"fmt"
	"strings"

	govalidation "github.com/abdalzizmahmoud2511/go-validation"
)

type User struct {
	Name  string `validate:"required;alpha"`
	Email string `validate:"required;email"`
	Age   int    `validate:"required;min=18;max=120"`
	Role  string `validate:"oneof=admin user guest"`
}

func main() {
	user := User{
		Name:  "John123",
		Email: "not-an-email",
		Age:   15,
		Role:  "superadmin",
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("CUSTOM LOCALE PATH EXAMPLE")
	fmt.Println(strings.Repeat("=", 60))

	// Default: uses embedded locales (package/locale/*.json)
	fmt.Println("\n--- 1. Embedded English (default) ---")
	fmt.Println("  Path: embedded in binary")
	errs := govalidation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Set custom locale path to load from external directory
	// This loads locales from ./locales/ directory (not from package)
	fmt.Println("\n--- 2. Custom Path: ./locales/ ---")
	fmt.Println("  Path: ./locales/")
	govalidation.SetLocalePath("examples/10_custom_locale/locales")

	// English from custom path (ALL CAPS style)
	fmt.Println("\n  English (custom):")
	errs = govalidation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// German from custom path (new language not in package)
	fmt.Println("\n  German (custom):")
	errs = govalidation.Validate("de", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Arabic still works from embedded (not overridden)
	fmt.Println("\n  Arabic (embedded fallback):")
	errs = govalidation.Validate("ar", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Show all available languages (embedded + custom)
	fmt.Println("\n--- 3. Available Languages ---")
	langs := govalidation.GetAvailableLanguages()
	fmt.Printf("  %v\n", langs)

	// Clear custom path and return to embedded
	fmt.Println("\n--- 4. After ClearLocaleCache() ---")
	govalidation.ClearLocaleCache()
	fmt.Println("  Back to embedded English:")
	errs = govalidation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}
}
