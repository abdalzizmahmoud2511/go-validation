package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

func main() {
	fmt.Println("=" + repeat("=", 60))
	fmt.Println("STRUCT VALIDATION")
	fmt.Println("=" + repeat("=", 60))

	type Address struct {
		Street  string `validate:"required;min=2"`
		City    string `validate:"required;alpha_space"`
		ZipCode string `validate:"required;len=6"`
		Country string `validate:"required;len=2"`
	}

	type Child struct {
		Name string `validate:"required;alpha"`
		Age  int    `validate:"required;min=1;max=17"`
	}

	type User struct {
		Name        string   `validate:"required;alpha"`
		Username    string   `validate:"required;alpha_num"`
		Slug        string   `validate:"required;alpha_dash"`
		Bio         string   `validate:"max_word=10"`
		Pattern     string   `validate:"regex=^[A-Z]{3}$"`
		Age         int      `validate:"required;min=18;max=100"`
		Score       float64  `validate:"required;between=0,100"`
		Zip         int      `validate:"required;digits=5"`
		Phone       string   `validate:"required;digits_between=10,15"`
		Email       string   `validate:"required;email"`
		Website     string   `validate:"required;url"`
		GUID        string   `validate:"required;uuid"`
		BirthDate   string   `validate:"required;date"`
		ExpireDate  string   `validate:"required;date_dmy"`
		IPv4        string   `validate:"required;ip_v4"`
		IPv6        string   `validate:"required;ip_v6"`
		MAC         string   `validate:"required;mac_address"`
		Latitude    string   `validate:"required;lat"`
		Longitude   string   `validate:"required;lon"`
		Coordinate  string   `validate:"required;coordinate"`
		Color       string   `validate:"required;css_color"`
		MetaData    string   `validate:"required;json"`
		CardNumber  string   `validate:"required;credit_card"`
		Role        string   `validate:"required;oneof=admin user guest"`
		Status      string   `validate:"required;in=active inactive banned"`
		CountryCode string   `validate:"required;not_in=XX YY ZZ"`
		IsActive    bool     `validate:"required;bool"`
		Address     Address  `validate:"required"`
		Children    []Child  `validate:"required"`
		Tags        []string `validate:"required"`
		Emails      []string `validate:"required;email"`
		Numbers     []int    `validate:"required;min=1"`
	}

	// Invalid user
	user := User{
		Name:        "John123",
		Username:    "john!",
		Slug:        "john doe",
		Bio:         "a b c d e f g h i j k",
		Pattern:     "abc",
		Age:         15,
		Score:       150,
		Zip:         1234,
		Phone:       "123",
		Email:       "invalid",
		Website:     "not-a-url",
		GUID:        "not-uuid",
		BirthDate:   "2000/01/01",
		ExpireDate:  "01-13-2025",
		IPv4:        "999.999.999.999",
		IPv6:        "not-ipv6",
		MAC:         "not-mac",
		Latitude:    "91",
		Longitude:   "181",
		Coordinate:  "invalid",
		Color:       "not-color",
		MetaData:    "not-json",
		CardNumber:  "1234",
		Role:        "superadmin",
		Status:      "deleted",
		CountryCode: "XX",
		IsActive:    true,
		Address:     Address{},
		Children:    []Child{},
		Tags:        []string{},
		Emails:      []string{"valid@test.com", "invalid"},
		Numbers:     []int{5, 0, 10},
	}

	fmt.Println("\n--- English Errors ---")
	errs := validation.Validate("en", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- Arabic Errors ---")
	errs = validation.Validate("ar", user)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Valid user
	fmt.Println("\n--- Valid User ---")
	validUser := User{
		Name:        "John",
		Username:    "john123",
		Slug:        "john-doe_123",
		Bio:         "Hello World",
		Pattern:     "ABC",
		Age:         25,
		Score:       85.5,
		Zip:         12345,
		Phone:       "1234567890",
		Email:       "john@example.com",
		Website:     "https://example.com",
		GUID:        "550e8400-e29b-41d4-a716-446655440000",
		BirthDate:   "1999-01-15",
		ExpireDate:  "15-06-2025",
		IPv4:        "192.168.1.1",
		IPv6:        "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		MAC:         "00:1A:2B:3C:4D:5E",
		Latitude:    "40.7128",
		Longitude:   "-74.0060",
		Coordinate:  "40.7128,-74.006",
		Color:       "rgb(255,87,51)",
		MetaData:    `{"key":"value"}`,
		CardNumber:  "4111111111111111",
		Role:        "admin",
		Status:      "active",
		CountryCode: "US",
		IsActive:    true,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			ZipCode: "100001",
			Country: "US",
		},
		Children: []Child{
			{Name: "Alice", Age: 5},
			{Name: "Bob", Age: 8},
		},
		Tags:    []string{"go", "validation"},
		Emails:  []string{"test@example.com", "valid@email.com"},
		Numbers: []int{1, 5, 10},
	}

	errs = validation.Validate("en", validUser)
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
