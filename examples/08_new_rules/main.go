package main

import (
	"fmt"

	validation "github.com/abdalzizmahmoud2511/go-validation"
)

type NewRulesUser struct {
	AlphaUnicode    string `validate:"required;alpha_unicode"`
	AlphaNumUnicode string `validate:"required;alphanum_unicode"`
	ASCII           string `validate:"required;ascii"`
	PrintASCII      string `validate:"required;print_ascii"`
	Lowercase       string `validate:"required;lowercase"`
	Uppercase       string `validate:"required;uppercase"`
	Contains        string `validate:"required;contains=hello"`
	ContainsAny     string `validate:"required;containsany=aeiou"`
	StartsWith      string `validate:"required;startswith=Mr"`
	EndsWith        string `validate:"required;endswith=Smith"`
	StartsNotWith   string `validate:"required;startsnotwith=Dr"`
	EndsNotWith     string `validate:"required;endsnotwith=Jones"`
	Gt              int    `validate:"required;gt=10"`
	Gte             int    `validate:"required;gte=10"`
	Lt              int    `validate:"required;lt=100"`
	Lte             int    `validate:"required;lte=100"`
	Eq              int    `validate:"required;eq=42"`
	Ne              int    `validate:"required;ne=0"`
	Base64          string `validate:"required;base64"`
	Hexadecimal     string `validate:"required;hexadecimal"`
	HexColor        string `validate:"required;hexcolor"`
	RGB             string `validate:"required;rgb"`
	RGBA            string `validate:"required;rgba"`
	HSL             string `validate:"required;hsl"`
	HSLA            string `validate:"required;hsla"`
	ISBN            string `validate:"required;isbn"`
	ISBN10          string `validate:"required;isbn10"`
	ISBN13          string `validate:"required;isbn13"`
	ISSN            string `validate:"required;issn"`
	JWT             string `validate:"required;jwt"`
	Semver          string `validate:"required;semver"`
	SSN             string `validate:"required;ssn"`
	Latitude        string `validate:"required;latitude"`
	Longitude       string `validate:"required;longitude"`
	NoneOf          string `validate:"required;noneof=admin root"`
	Unique          []int  `validate:"required;unique"`
}

func main() {
	fmt.Println("\n" + repeat("=", 61))
	fmt.Println("NEW RULES EXAMPLE")
	fmt.Println(repeat("=", 61))

	// Invalid user with new rules
	invalidUser := NewRulesUser{
		AlphaUnicode:    "John123",
		AlphaNumUnicode: "John!",
		ASCII:           "Hello",
		PrintASCII:      "Hello",
		Lowercase:       "Hello",
		Uppercase:       "hello",
		Contains:        "world",
		ContainsAny:     "bcdfg",
		StartsWith:      "Dr Smith",
		EndsWith:        "Jones",
		StartsNotWith:   "Dr Jones",
		EndsNotWith:     "Dr Jones",
		Gt:              5,
		Gte:             5,
		Lt:              150,
		Lte:             150,
		Eq:              10,
		Ne:              0,
		Base64:          "not-base64!",
		Hexadecimal:     "not-hex!",
		HexColor:        "not-hex",
		RGB:             "not-rgb",
		RGBA:            "not-rgba",
		HSL:             "not-hsl",
		HSLA:            "not-hsla",
		ISBN:            "not-isbn",
		ISBN10:          "not-isbn10",
		ISBN13:          "not-isbn13",
		ISSN:            "not-issn",
		JWT:             "not-jwt",
		Semver:          "not-semver",
		SSN:             "not-ssn",
		Latitude:        "91",
		Longitude:       "181",
		NoneOf:          "admin",
		Unique:          []int{1, 2, 2},
	}

	fmt.Println("\n--- English Errors ---")
	errs := validation.Validate("en", invalidUser)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("\n--- Arabic Errors ---")
	errs = validation.Validate("ar", invalidUser)
	for _, err := range errs {
		fmt.Printf("  - %s\n", err.Error())
	}

	// Valid user with new rules
	fmt.Println("\n--- Valid User ---")
	validUser := NewRulesUser{
		AlphaUnicode:    "John",
		AlphaNumUnicode: "John123",
		ASCII:           "Hello",
		PrintASCII:      "Hello",
		Lowercase:       "hello",
		Uppercase:       "HELLO",
		Contains:        "hello world",
		ContainsAny:     "hello",
		StartsWith:      "Mr Smith",
		EndsWith:        "Mr Smith",
		StartsNotWith:   "Mr Jones",
		EndsNotWith:     "Mr Smith",
		Gt:              15,
		Gte:             15,
		Lt:              50,
		Lte:             50,
		Eq:              42,
		Ne:              5,
		Base64:          "SGVsbG8gV29ybGQ=",
		Hexadecimal:     "1234567890ABCDEF",
		HexColor:        "#FF5733",
		RGB:             "rgb(255, 87, 51)",
		RGBA:            "rgba(255, 87, 51, 0.5)",
		HSL:             "hsl(12, 100%, 60%)",
		HSLA:            "hsla(12, 100%, 60%, 0.5)",
		ISBN:            "978-3-16-148410-0",
		ISBN10:          "0-306-40615-2",
		ISBN13:          "978-3-16-148410-0",
		ISSN:            "0378-5955",
		JWT:             "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
		Semver:          "v1.2.3",
		SSN:             "123-45-6789",
		Latitude:        "40.7128",
		Longitude:       "-74.0060",
		NoneOf:          "user",
		Unique:          []int{1, 2, 3},
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
