package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

type ExampleAddress struct {
	Street  string
	City    string
	ZipCode string
	Country string
}

type ExampleChild struct {
	Name string
	Age  int
}

type ExampleUser struct {
	Name     string
	Email    string
	Age      int
	Role     string
	Address  ExampleAddress
	Children []ExampleChild
	Emails   []string
	Numbers  []int
}

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("ValidateStructWithRules() - STRUCT WITH MAP RULES")
	fmt.Println(repeat("=", 61))

	rules := map[string][]string{
		"Name":            {"required", "alpha"},
		"Email":           {"required", "email"},
		"Age":             {"required", "min=18"},
		"Role":            {"required", "oneof=admin user guest"},
		"Address.Street":  {"required", "min=2"},
		"Address.City":    {"required", "alpha_space"},
		"Address.ZipCode": {"required", "len=6"},
		"Address.Country": {"required", "len=2"},
		"Children.Name":   {"required", "alpha"},
		"Children.Age":    {"required", "min=1"},
		"Emails":          {"required", "email"},
		"Numbers":         {"required", "min=1"},
	}

	// Invalid struct
	fmt.Println("\n--- English Errors ---")
	invalidUser := ExampleUser{
		Name:  "John123",
		Email: "invalid",
		Age:   15,
		Role:  "superadmin",
		Address: ExampleAddress{
			Street:  "",
			City:    "",
			ZipCode: "123",
			Country: "USA",
		},
		Children: []ExampleChild{
			{Name: "Alice", Age: 5},
			{Name: "", Age: 0},
		},
		Emails:  []string{"valid@test.com", "invalid"},
		Numbers: []int{5, 0, 10},
	}

	errs := validation.ValidateStructWithRules("en", invalidUser, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- Arabic Errors ---")
	errs = validation.ValidateStructWithRules("ar", invalidUser, rules)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Valid struct
	fmt.Println("\n--- Valid User ---")
	validUser := ExampleUser{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
		Role:  "admin",
		Address: ExampleAddress{
			Street:  "123 Main St",
			City:    "New York",
			ZipCode: "100001",
			Country: "US",
		},
		Children: []ExampleChild{
			{Name: "Alice", Age: 5},
			{Name: "Bob", Age: 8},
		},
		Emails:  []string{"test@example.com", "valid@email.com"},
		Numbers: []int{1, 5, 10},
	}

	errs = validation.ValidateStructWithRules("en", validUser, rules)
	if len(errs) == 0 {
		fmt.Println("  All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("  - %s\n", err.Error())
		}
	}

	// With custom messages using Msg()
	fmt.Println("\n--- With Custom Messages (Arabic) ---")
	customRules := map[string][]string{
		"Name":           {validation.Msg("required", "custom_name_required"), validation.Msg("alpha", "custom_name_alpha")},
		"Email":          {validation.Msg("email", "custom_email_invalid")},
		"Age":            {validation.Msg("min", "custom_age_min")},
		"Address.Street": {validation.Msg("required", "custom_street_required")},
	}

	errs = validation.ValidateStructWithRules("ar", invalidUser, customRules)
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
