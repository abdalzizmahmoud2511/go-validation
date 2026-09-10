package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("MAP VALIDATION")
	fmt.Println(repeat("=", 61))

	// Invalid map
	data := map[string]interface{}{
		"name":     "John123",
		"username": "john!",
		"age":      15,
		"score":    150,
		"email":    "invalid",
		"website":  "not-a-url",
		"role":     "superadmin",
		"address": map[string]interface{}{
			"street":  "",
			"city":    "",
			"zip":     "123",
			"country": "USA",
		},
		"children": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 5},
			map[string]interface{}{"name": "", "age": 0},
		},
		"emails":  []interface{}{"valid@test.com", "invalid"},
		"numbers": []interface{}{5, 0, 10},
	}

	rules := map[string][]string{
		"name":            {"required", "alpha"},
		"username":        {"required", "alpha_num"},
		"age":             {"required", "min=18", "max=100"},
		"score":           {"required", "between=0,100"},
		"email":           {"required", "email"},
		"website":         {"required", "url"},
		"role":            {"required", "oneof=admin user guest"},
		"address.street":  {"required", "min=2"},
		"address.city":    {"required", "alpha_space"},
		"address.zip":     {"required", "len=6"},
		"address.country": {"required", "len=2"},
		"children.name":   {"required"},
		"children.age":    {"required", "min=1"},
		"emails":          {"required", "email"},
		"numbers":         {"required", "min=1"},
	}

	fmt.Println("\n--- English Errors ---")
	errs := validation.ValidateMap("en", data, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- Arabic Errors ---")
	errs = validation.ValidateMap("ar", data, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Valid map
	fmt.Println("\n--- Valid Map ---")
	validData := map[string]interface{}{
		"name":     "John",
		"username": "john123",
		"age":      25,
		"score":    85.5,
		"email":    "john@example.com",
		"website":  "https://example.com",
		"role":     "admin",
		"address": map[string]interface{}{
			"street":  "123 Main St",
			"city":    "New York",
			"zip":     "100001",
			"country": "US",
		},
		"children": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 5},
			map[string]interface{}{"name": "Bob", "age": 8},
		},
		"emails":  []interface{}{"test@example.com", "valid@email.com"},
		"numbers": []interface{}{1, 5, 10},
	}

	errs = validation.ValidateMap("en", validData, rules)
	if len(errs) == 0 {
		fmt.Println("  All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("  - %s\n", err.Error())
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
