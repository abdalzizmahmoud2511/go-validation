# Go-validation

A powerful Go validation library with support for structs, maps, nested fields, localization, and conditional rules.


## Features

- **70+ Built-in Rules** - required, email, url, uuid, min, max, between, regex, and more
- **Struct & Map Validation** - validate both structs and maps with the same API
- **Nested Struct Support** - dot notation for deep fields (`Address.City`)
- **Array/Slice Validation** - validate each element in arrays
- **Custom Rules** - register your own validation rules
- **Custom Error Messages** - override messages per field
- **Localization (i18n)** - English and Arabic out of the box
- **Conditional Validation** - `ApplyIf()` and `ApplyIfElse()` for conditional rules
- **Skip Empty Fields** - configurable behavior (default: skip empty)
- **Zero Dependencies** - only stdlib

## Installation

```bash
go get github.com/yourusername/goven
```

## Quick Start

### Struct Validation

```go
package main

import (
    "fmt"
    "github.com/yourusername/goven/pkg/validation"
)

type User struct {
    Name  string `validate:"required;alpha"`
    Email string `validate:"required;email"`
    Age   int    `validate:"required;min=18;max=100"`
}

func main() {
    user := User{Name: "John123", Email: "invalid", Age: 15}
    
    errs := validation.Validate("en", user)
    for _, err := range errs {
        fmt.Println(err.Error())
    }
    // Output:
    // The Name field must contain only letters
    // The Email field must be a valid email address
    // The Age field must be at least 18
}
```

### Map Validation

```go
data := map[string]interface{}{
    "name":  "John",
    "email": "invalid",
}

rules := map[string][]string{
    "name":  {"required", "alpha"},
    "email": {"required", "email"},
}

errs := validation.ValidateMap("en", data, rules)
```

### Custom Messages

```go
rules := map[string][]string{
    "Name":  {validation.Msg("required", "custom_name_required")},
    "Email": {validation.Msg("email", "custom_email_invalid")},
}
```

### Conditional Rules

```go
// Apply 'required' to ZipCode if Country is valid
rules := map[string][]string{
    "ZipCode": {
        validation.ApplyIf("required", map[string][]string{
            "Country": {"required", "alpha"},
        }),
    },
}
```

## Available Rules

| Rule | Description | Example |
|------|-------------|---------|
| `required` | Field is required | `validate:"required"` |
| `alpha` | Letters only | `validate:"alpha"` |
| `alpha_num` | Letters and numbers | `validate:"alpha_num"` |
| `alpha_dash` | Letters, numbers, dashes, underscores | `validate:"alpha_dash"` |
| `email` | Valid email | `validate:"email"` |
| `url` | Valid URL | `validate:"url"` |
| `uuid` | Valid UUID | `validate:"uuid"` |
| `min` | Minimum value/length | `validate:"min=5"` |
| `max` | Maximum value/length | `validate:"max=100"` |
| `between` | Value between min and max | `validate:"between=5,50"` |
| `len` | Exact length | `validate:"len=10"` |
| `regex` | Regex pattern match | `validate:"regex=^[A-Z]+$"` |
| `oneof` | One of allowed values | `validate:"oneof=admin user"` |
| `in` | In allowed values | `validate:"in=active inactive"` |
| `not_in` | Not in values | `validate:"not_in=banned"` |
| `gt` | Greater than | `validate:"gt=0"` |
| `gte` | Greater than or equal | `validate:"gte=18"` |
| `lt` | Less than | `validate:"lt=100"` |
| `lte` | Less than or equal | `validate:"lte=100"` |
| `digits` | Exactly N digits | `validate:"digits=5"` |
| `digits_between` | Digits between min and max | `validate:"digits_between=10,15"` |
| `date` | Valid date (yyyy-mm-dd) | `validate:"date"` |
| `date_dmy` | Valid date (dd-mm-yyyy) | `validate:"date_dmy"` |
| `ip` | Valid IP | `validate:"ip"` |
| `ip_v4` | Valid IPv4 | `validate:"ip_v4"` |
| `ip_v6` | Valid IPv6 | `validate:"ip_v6"` |
| `mac_address` | Valid MAC | `validate:"mac_address"` |
| `lat` | Valid latitude | `validate:"lat"` |
| `lon` | Valid longitude | `validate:"lon"` |
| `coordinate` | Valid lat,lon | `validate:"coordinate"` |
| `css_color` | Valid CSS color | `validate:"css_color"` |
| `json` | Valid JSON | `validate:"json"` |
| `credit_card` | Valid credit card | `validate:"credit_card"` |
| `bool` | Boolean value | `validate:"bool"` |
| `contains` | Contains substring | `validate:"contains=hello"` |
| `containsany` | Contains any of chars | `validate:"containsany=aeiou"` |
| `startswith` | Starts with | `validate:"startswith=Mr"` |
| `endswith` | Ends with | `validate:"endswith=Jr"` |
| `lowercase` | All lowercase | `validate:"lowercase"` |
| `uppercase` | All uppercase | `validate:"uppercase"` |
| `base64` | Valid base64 | `validate:"base64"` |
| `hexadecimal` | Valid hex | `validate:"hexadecimal"` |
| `hexcolor` | Valid hex color | `validate:"hexcolor"` |
| `rgb` | Valid RGB | `validate:"rgb"` |
| `rgba` | Valid RGBA | `validate:"rgba"` |
| `hsl` | Valid HSL | `validate:"hsl"` |
| `hsla` | Valid HSLA | `validate:"hsla"` |
| `isbn` | Valid ISBN | `validate:"isbn"` |
| `isbn10` | Valid ISBN-10 | `validate:"isbn10"` |
| `isbn13` | Valid ISBN-13 | `validate:"isbn13"` |
| `issn` | Valid ISSN | `validate:"issn"` |
| `jwt` | Valid JWT | `validate:"jwt"` |
| `semver` | Valid semver | `validate:"semver"` |
| `ssn` | Valid SSN | `validate:"ssn"` |
| `unique` | Unique values in array | `validate:"unique"` |
| `min_word` | Minimum word count | `validate:"min_word=5"` |
| `max_word` | Maximum word count | `validate:"max_word=10"` |
| `len_word` | Exact word count | `validate:"len_word=5"` |

## Localization

### Supported Languages
- English (`en`)
- Arabic (`ar`)

### Adding New Language

```go
// Add locale file: locale/fr.json
{
    "required": "Le champ :field est obligatoire",
    "email": "Le champ :field doit être un email valide",
    "min_string": "Le champ :field doit avoir au moins :arg caractères"
}

// Use it
errs := validation.Validate("fr", user)
```

## Custom Rules

```go
func init() {
    validation.AddCustomRule("my_rule", func(lang, field, rule, message string, value interface{}) error {
        // Your validation logic
        if !isValid(value) {
            if message != "" {
                return fmt.Errorf(message)
            }
            return fmt.Errorf("validation failed for %s", field)
        }
        return nil
    })
}
```

## Skip Empty Fields

```go
// Default: skip empty fields (only validate non-empty)
errs := validation.Validate("en", user)

// Validate all fields (empty fields fail)
errs = validation.Validate("en", user, false)
```

## Examples

Run the examples:

```bash
go run examples/01_struct/main.go
go run examples/09_apply_if/main.go
```

## License

MIT License
