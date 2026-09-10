# go-validation

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
- **Embedded Locale Files** - no external files needed

## Installation

```bash
go get github.com/abdalzizmahmoud2511/go-validation
```

## Quick Start

### Struct Validation

```go
package main

import (
    "fmt"
    validation "github.com/abdalzizmahmoud2511/go-validation"
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

### Struct with External Rules

```go
type User struct {
    Name    string
    Email   string
    Address Address
}

type Address struct {
    Street string
    City   string
}

rules := map[string][]string{
    "Name":          {"required", "alpha"},
    "Email":         {"required", "email"},
    "Address.Street": {"required", "min=2"},
    "Address.City":   {"required"},
}

errs := validation.ValidateStructWithRules("en", user, rules)
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

// Apply different rules based on conditions
rules := map[string][]string{
    "ZipCode": {
        validation.ApplyIfElse("required", "nullable", map[string][]string{
            "Country": {"required", "alpha"},
        }),
    },
}
```

### Skip Empty Fields

```go
// Default: skip empty fields (only validate non-empty)
errs := validation.Validate("en", user)

// Validate all fields (empty fields fail)
errs = validation.Validate("en", user, false)
```

## Available Rules

### String Rules
| Rule | Description | Example |
|------|-------------|---------|
| `required` | Field is required | `validate:"required"` |
| `alpha` | Letters only | `validate:"alpha"` |
| `alpha_num` | Letters and numbers | `validate:"alpha_num"` |
| `alpha_dash` | Letters, numbers, dashes, underscores | `validate:"alpha_dash"` |
| `alpha_space` | Letters, numbers, dashes, underscores, spaces | `validate:"alpha_space"` |
| `alpha_unicode` | Unicode letters only | `validate:"alpha_unicode"` |
| `alphanum_unicode` | Unicode letters and numbers | `validate:"alphanum_unicode"` |
| `ascii` | ASCII characters only | `validate:"ascii"` |
| `print_ascii` | Printable ASCII only | `validate:"print_ascii"` |
| `lowercase` | All lowercase | `validate:"lowercase"` |
| `uppercase` | All uppercase | `validate:"uppercase"` |
| `contains` | Contains substring | `validate:"contains=hello"` |
| `containsany` | Contains any of chars | `validate:"containsany=aeiou"` |
| `containsrune` | Contains specific rune | `validate:"containsrune=@@" |
| `excludes` | Must not contain substring | `validate:"excludes=bad"` |
| `excludesall` | Must not contain any of chars | `validate:"excludesall=!@#"` |
| `excludesrune` | Must not contain specific rune | `validate:"excludesrune=@@" |
| `startswith` | Starts with | `validate:"startswith=Mr"` |
| `endswith` | Ends with | `validate:"endswith=Jr"` |
| `startsnotwith` | Must not start with | `validate:"startsnotwith=Dr"` |
| `endsnotwith` | Must not end with | `validate:"endsnotwith=Jr"` |
| `regex` | Regex pattern match | `validate:"regex=^[A-Z]+$"` |
| `number` | Numeric string | `validate:"number"` |

### Numeric Rules
| Rule | Description | Example |
|------|-------------|---------|
| `numeric` | Numeric value | `validate:"numeric"` |
| `min` | Minimum value/length | `validate:"min=5"` |
| `max` | Maximum value/length | `validate:"max=100"` |
| `between` | Value between min and max | `validate:"between=5,50"` |
| `len` | Exact length | `validate:"len=10"` |
| `digits` | Exactly N digits | `validate:"digits=5"` |
| `digits_between` | Digits between min and max | `validate:"digits_between=10,15"` |
| `gt` | Greater than | `validate:"gt=0"` |
| `gte` | Greater than or equal | `validate:"gte=18"` |
| `lt` | Less than | `validate:"lt=100"` |
| `lte` | Less than or equal | `validate:"lte=100"` |
| `eq` | Equal to | `validate:"eq=5"` |
| `ne` | Not equal to | `validate:"ne=0"` |

### Format Rules
| Rule | Description | Example |
|------|-------------|---------|
| `email` | Valid email | `validate:"email"` |
| `url` | Valid URL | `validate:"url"` |
| `uri` | Valid URI | `validate:"uri"` |
| `uuid` | Valid UUID | `validate:"uuid"` |
| `uuid3` | Valid UUID v3 | `validate:"uuid3"` |
| `uuid4` | Valid UUID v4 | `validate:"uuid4"` |
| `uuid5` | Valid UUID v5 | `validate:"uuid5"` |
| `date` | Valid date (yyyy-mm-dd) | `validate:"date"` |
| `date_dmy` | Valid date (dd-mm-yyyy) | `validate:"date_dmy"` |
| `datetime` | Valid datetime with format | `validate:"datetime=2006-01-02 15:04"` |
| `timezone` | Valid timezone | `validate:"timezone"` |
| `ip` | Valid IP | `validate:"ip"` |
| `ip_v4` | Valid IPv4 | `validate:"ip_v4"` |
| `ip_v6` | Valid IPv6 | `validate:"ip_v6"` |
| `mac_address` | Valid MAC | `validate:"mac_address"` |
| `latitude` | Valid latitude | `validate:"latitude"` |
| `longitude` | Valid longitude | `validate:"longitude"` |
| `lat` | Valid latitude | `validate:"lat"` |
| `lon` | Valid longitude | `validate:"lon"` |
| `coordinate` | Valid lat,lon | `validate:"coordinate"` |
| `hexcolor` | Valid hex color | `validate:"hexcolor"` |
| `css_color` | Valid CSS color | `validate:"css_color"` |
| `rgb` | Valid RGB | `validate:"rgb"` |
| `rgba` | Valid RGBA | `validate:"rgba"` |
| `hsl` | Valid HSL | `validate:"hsl"` |
| `hsla` | Valid HSLA | `validate:"hsla"` |
| `json` | Valid JSON | `validate:"json"` |
| `credit_card` | Valid credit card | `validate:"credit_card"` |
| `base64` | Valid base64 | `validate:"base64"` |
| `hexadecimal` | Valid hex | `validate:"hexadecimal"` |
| `isbn` | Valid ISBN | `validate:"isbn"` |
| `isbn10` | Valid ISBN-10 | `validate:"isbn10"` |
| `isbn13` | Valid ISBN-13 | `validate:"isbn13"` |
| `issn` | Valid ISSN | `validate:"issn"` |
| `jwt` | Valid JWT | `validate:"jwt"` |
| `semver` | Valid semver | `validate:"semver"` |
| `ssn` | Valid SSN | `validate:"ssn"` |

### Enum Rules
| Rule | Description | Example |
|------|-------------|---------|
| `oneof` | One of allowed values | `validate:"oneof=admin user"` |
| `in` | In allowed values | `validate:"in=active inactive"` |
| `not_in` | Not in values | `validate:"not_in=banned"` |
| `noneof` | Not in values (alias) | `validate:"noneof=banned"` |
| `bool` | Boolean value | `validate:"bool"` |
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

### Available Message Variables
- `:field` - The field name
- `:arg` - The rule argument
- `:min` - Minimum value (for between)
- `:max` - Maximum value (for between)

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

// Remove a custom rule
validation.RemoveCustomRule("my_rule")

// Check if a custom rule exists
if validation.HasCustomRule("my_rule") {
    // ...
}
```

## ValError Type

All validation errors return `ValError` which provides:

```go
type ValError struct {
    Field string
    Rule  string
    Msg   string
}

func (e ValError) Error() string     // Returns Msg
func (e ValError) GetField() string  // Returns Field
func (e ValError) GetRule() string   // Returns Rule
func (e ValError) GetMessage() string // Returns Msg
```

## Configuration

```go
// Set default language
validation.SetLanguage("ar")

// Set locale file path (for custom locale files)
validation.SetLocalePath("/path/to/locale")

// Get current config
cfg := validation.GetConfig()

// Clear locale cache (after adding new locale files)
validation.ClearLocaleCache()

// Get available languages
langs := validation.GetAvailableLanguages()
```

## Examples

```go
package main

import (
    "fmt"
    validation "github.com/abdalzizmahmoud2511/go-validation"
)

type User struct {
    Name    string  `validate:"required;alpha"`
    Email   string  `validate:"required;email"`
    Age     int     `validate:"required;min=18;max=100"`
    Website string  `validate:"url"`
    Role    string  `validate:"oneof=admin user guest"`
}

type Address struct {
    Street string `validate:"required;min=5"`
    City   string `validate:"required;alpha"`
    Zip    string `validate:"required;len=5"`
}

func main() {
    user := User{
        Name:    "John",
        Email:   "john@example.com",
        Age:     25,
        Website: "https://example.com",
        Role:    "admin",
    }

    // Validate with struct tags
    errs := validation.Validate("en", user)
    if len(errs) > 0 {
        for _, err := range errs {
            fmt.Println(err.Error())
        }
    }

    // Validate with external rules
    rules := map[string][]string{
        "Name":  {"required", "alpha"},
        "Email": {"required", "email"},
        "Age":   {"required", "min=18"},
    }
    errs = validation.ValidateStructWithRules("en", user, rules)
}
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
