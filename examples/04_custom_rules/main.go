package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func init() {
	validation.AddCustomRule("my_custom", func(lang, field, rule, message string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("my_custom rule only supported on strings")
		}
		if len(str) < 3 {
			if message != "" {
				return fmt.Errorf(message)
			}
			return fmt.Errorf("validation failed for %s", field)
		}
		return nil
	})
}

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("CUSTOM RULES")
	fmt.Println(repeat("=", 61))

	type User struct {
		Bio string `validate:"required;max_word=5"`
		Age int    `validate:"required;min=18"`
	}

	user := User{Bio: "This has way too many words in it", Age: 15}

	fmt.Println("\n--- English ---")
	errs := validation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- Arabic ---")
	errs = validation.Validate("ar", user)
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
