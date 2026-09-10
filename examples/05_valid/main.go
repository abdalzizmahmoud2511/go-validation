package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("Valid() FUNCTION")
	fmt.Println(repeat("=", 61))

	type User struct {
		Name  string `validate:"required;alpha"`
		Email string `validate:"required;email"`
	}

	validUser := User{Name: "John", Email: "john@example.com"}
	invalidUser := User{Name: "", Email: "invalid"}

	fmt.Printf("\n  Valid user passes: %v\n", validation.Valid("en", validUser))
	fmt.Printf("  Invalid user passes: %v\n", validation.Valid("en", invalidUser))
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
