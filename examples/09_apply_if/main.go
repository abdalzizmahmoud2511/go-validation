package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

type ApplyIfUser struct {
	Name    string
	Country string
	ZipCode string
	State   string
	Email   string
	Age     int
	Role    string
}

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("ApplyIf() FUNCTION")
	fmt.Println(repeat("=", 61))

	// Example 1: Basic ApplyIf
	fmt.Println("\n--- Example 1: Apply 'required' to ZipCode if Country passes conditions ---")
	user1 := ApplyIfUser{
		Name:    "John",
		Country: "US",
		ZipCode: "12345",
	}

	rules1 := map[string][]string{
		"ZipCode": {
			validation.ApplyIf("required", map[string][]string{
				"Country": {"required", "alpha"},
			}),
		},
	}
	fmt.Printf("  Generated rule: %s\n", rules1["ZipCode"][0])

	errs := validation.ValidateWithRules("en", user1, rules1)
	if len(errs) == 0 {
		fmt.Println("  All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("  - %s\n", err.Error())
		}
	}

	// Example 2: Condition fails - rule not applied
	fmt.Println("\n--- Example 2: Condition fails - rule not applied ---")
	user2 := ApplyIfUser{
		Name:    "John",
		Country: "US",
		ZipCode: "",
	}

	errs = validation.ValidateWithRules("en", user2, rules1)
	fmt.Println("  Country=US, ZipCode='' (should fail):")
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Example 3: Multiple fields
	fmt.Println("\n--- Example 3: Multiple fields with different conditions ---")
	user3 := ApplyIfUser{
		Name:    "John",
		Country: "US",
		ZipCode: "12345",
		State:   "NY",
	}

	rules3 := map[string][]string{
		"ZipCode": {
			validation.ApplyIf("required", map[string][]string{
				"Country": {"required", "alpha"},
			}),
		},
		"State": {
			validation.ApplyIf("required", map[string][]string{
				"Country": {"required"},
			}),
		},
	}

	errs = validation.ValidateWithRules("en", user3, rules3)
	if len(errs) == 0 {
		fmt.Println("  All validations passed!")
	} else {
		for _, err := range errs {
			fmt.Printf("  - %s\n", err.Error())
		}
	}

	// Example 4: Struct tag usage
	fmt.Println("\n--- Example 4: Struct tag usage ---")

	type Address struct {
		Country string
		ZipCode string `validate:"required"`
	}

	type Profile struct {
		Name    string
		Address Address
	}

	rules4 := map[string][]string{
		"Address.ZipCode": {
			validation.ApplyIf("required", map[string][]string{
				"Address.Country": {"required", "alpha"},
			}),
		},
	}

	profile := Profile{
		Name: "John",
		Address: Address{
			Country: "US",
			ZipCode: "",
		},
	}

	errs = validation.ValidateStructWithRules("en", profile, rules4)
	fmt.Println("  Country=US, ZipCode='' (should fail):")
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Example 5: ApplyIfElse
	fmt.Println("\n--- Example 5: ApplyIfElse - apply 'required' or 'nullable' ---")
	user5 := ApplyIfUser{
		Name:    "John",
		Country: "US",
		ZipCode: "",
	}

	rules5 := map[string][]string{
		"ZipCode": {
			validation.ApplyIfElse("required", "nullable", map[string][]string{
				"Country": {"required", "alpha"},
			}),
		},
	}

	errs = validation.ValidateWithRules("en", user5, rules5)
	fmt.Println("  Country=US, ZipCode='' (should fail with required):")
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Example 6: Condition with numeric comparison
	fmt.Println("\n--- Example 6: Condition with numeric comparison ---")

	type Order struct {
		Age     int
		ZipCode string `validate:"required"`
	}

	rules6 := map[string][]string{
		"ZipCode": {
			validation.ApplyIf("required", map[string][]string{
				"Age": {"gte=18"},
			}),
		},
	}

	order1 := Order{Age: 25, ZipCode: ""}
	errs = validation.ValidateWithRules("en", order1, rules6)
	fmt.Println("  Age=25, ZipCode='' (should fail - Age >= 18):")
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	order2 := Order{Age: 15, ZipCode: ""}
	errs = validation.ValidateWithRules("en", order2, rules6)
	fmt.Println("  Age=15, ZipCode='' (should pass - Age < 18):")
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Example 7: Map validation with ApplyIf
	fmt.Println("\n--- Example 7: Map validation with ApplyIf ---")

	data := map[string]interface{}{
		"Country": "US",
		"ZipCode": "",
	}

	rules7 := map[string][]string{
		"ZipCode": {
			validation.ApplyIf("required", map[string][]string{
				"Country": {"required", "alpha"},
			}),
		},
	}

	errs = validation.ValidateMap("en", data, rules7)
	fmt.Println("  Country=US, ZipCode='' (should fail):")
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
