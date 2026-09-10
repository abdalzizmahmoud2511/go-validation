package govalidation

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	emailRe       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	urlRe         = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	uuidRe        = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	ipRe          = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	ipV6Re        = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	macRe         = regexp.MustCompile(`^([0-9a-fA-F]{2}[:-]){5}[0-9a-fA-F]{2}$`)
	latRe         = regexp.MustCompile(`^-?([1-8]?\d(\.\d+)?|90(\.0+)?)$`)
	lonRe         = regexp.MustCompile(`^-?((1[0-7]\d|[1-9]?\d)(\.\d+)?|180(\.0+)?)$`)
	coordRe       = regexp.MustCompile(`^-?([1-8]?\d(\.\d+)?|90(\.0+)?)\s*,\s*-?((1[0-7]\d|[1-9]?\d)(\.\d+)?|180(\.0+)?)$`)
	hexColorRe    = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	cssColorRe    = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$|^[a-zA-Z]+$|^rgb\(|^rgba\(|^hsl\(|^hsla\(`)
	jsonRe        = regexp.MustCompile(`^\s*[\{\[]`)
	creditCardRe  = regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$|^5[1-5][0-9]{14}$|^3[47][0-9]{13}$|^3(?:0[0-5]|[68][0-9])[0-9]{11}$|^6(?:011|5[0-9]{2})[0-9]{12}$|^(?:2131|1800|35\d{3})\d{11}$`)
	base64Re      = regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	hexadecimalRe = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	rgbRe         = regexp.MustCompile(`^rgb\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)$`)
	rgbaRe        = regexp.MustCompile(`^rgba\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*(0(\.\d+)?|1(\.0+)?)\s*\)$`)
	hslRe         = regexp.MustCompile(`^hsl\(\s*\d{1,3}\s*,\s*\d{1,3}%\s*,\s*\d{1,3}%\s*\)$`)
	hslaRe        = regexp.MustCompile(`^hsla\(\s*\d{1,3}\s*,\s*\d{1,3}%\s*,\s*\d{1,3}%\s*,\s*(0(\.\d+)?|1(\.0+)?)\s*\)$`)
	semverRe      = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	ssnRe         = regexp.MustCompile(`^\d{9}$`)
)

func validateEmail(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "email", msg, "only supported on strings")
	}
	_, err := mail.ParseAddress(fv.String())
	if err != nil || !emailRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "email", msg}
		}
		return locErr(lang, "email", name, "")
	}
	return nil
}

func validateURL(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "url", msg, "only supported on strings")
	}
	u, err := url.Parse(fv.String())
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		if msg != "" {
			return ValError{name, "url", msg}
		}
		return locErr(lang, "url", name, "")
	}
	return nil
}

func validateURI(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uri", msg, "only supported on strings")
	}
	u, err := url.Parse(fv.String())
	if err != nil || u.Scheme == "" || u.Host == "" {
		if msg != "" {
			return ValError{name, "uri", msg}
		}
		return locErr(lang, "uri", name, "")
	}
	return nil
}

// isValidUUIDChar checks if byte is a valid hex or dash character for UUID
func isValidUUIDChar(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F') || b == '-'
}

func validateUUID(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uuid", msg, "only supported on strings")
	}
	s := fv.String()
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		if msg != "" {
			return ValError{name, "uuid", msg}
		}
		return locErr(lang, "uuid", name, "")
	}
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !isValidUUIDChar(s[i]) {
			if msg != "" {
				return ValError{name, "uuid", msg}
			}
			return locErr(lang, "uuid", name, "")
		}
	}
	return nil
}

func validateUUID3(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uuid3", msg, "only supported on strings")
	}
	s := fv.String()
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' || s[14] != '3' {
		if msg != "" {
			return ValError{name, "uuid3", msg}
		}
		return locErr(lang, "uuid3", name, "")
	}
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !isValidUUIDChar(s[i]) {
			if msg != "" {
				return ValError{name, "uuid3", msg}
			}
			return locErr(lang, "uuid3", name, "")
		}
	}
	return nil
}

func validateUUID4(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uuid4", msg, "only supported on strings")
	}
	s := fv.String()
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' || s[14] != '4' {
		if msg != "" {
			return ValError{name, "uuid4", msg}
		}
		return locErr(lang, "uuid4", name, "")
	}
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !isValidUUIDChar(s[i]) {
			if msg != "" {
				return ValError{name, "uuid4", msg}
			}
			return locErr(lang, "uuid4", name, "")
		}
	}
	return nil
}

func validateUUID5(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "uuid5", msg, "only supported on strings")
	}
	s := fv.String()
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' || s[14] != '5' {
		if msg != "" {
			return ValError{name, "uuid5", msg}
		}
		return locErr(lang, "uuid5", name, "")
	}
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !isValidUUIDChar(s[i]) {
			if msg != "" {
				return ValError{name, "uuid5", msg}
			}
			return locErr(lang, "uuid5", name, "")
		}
	}
	return nil
}

func validateDate(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "date", msg, "only supported on strings")
	}
	_, err := time.Parse("2006-01-02", fv.String())
	if err != nil {
		if msg != "" {
			return ValError{name, "date", msg}
		}
		return locErr(lang, "date", name, "")
	}
	return nil
}

func validateDateDDMMYY(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "date_dmy", msg, "only supported on strings")
	}
	_, err := time.Parse("02-01-2006", fv.String())
	if err != nil {
		if msg != "" {
			return ValError{name, "date_dmy", msg}
		}
		return locErr(lang, "date_dmy", name, "")
	}
	return nil
}

func validateDateTime(lang, name string, fv reflect.Value, format, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "datetime", msg, "only supported on strings")
	}
	_, err := time.Parse(format, fv.String())
	if err != nil {
		if msg != "" {
			return ValError{name, "datetime", msg}
		}
		return locErr(lang, "datetime", name, "")
	}
	return nil
}

func validateTimezone(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "timezone", msg, "only supported on strings")
	}
	_, err := time.LoadLocation(fv.String())
	if err != nil {
		if msg != "" {
			return ValError{name, "timezone", msg}
		}
		return locErr(lang, "timezone", name, "")
	}
	return nil
}

func validateIP(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "ip", msg, "only supported on strings")
	}
	if net.ParseIP(fv.String()) == nil {
		if msg != "" {
			return ValError{name, "ip", msg}
		}
		return locErr(lang, "ip", name, "")
	}
	return nil
}

func validateIPv4(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "ip_v4", msg, "only supported on strings")
	}
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() == nil {
		if msg != "" {
			return ValError{name, "ip_v4", msg}
		}
		return locErr(lang, "ip_v4", name, "")
	}
	return nil
}

func validateIPv6(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "ip_v6", msg, "only supported on strings")
	}
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() != nil {
		if msg != "" {
			return ValError{name, "ip_v6", msg}
		}
		return locErr(lang, "ip_v6", name, "")
	}
	return nil
}

func validateMAC(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "mac_address", msg, "only supported on strings")
	}
	_, err := net.ParseMAC(fv.String())
	if err != nil {
		if msg != "" {
			return ValError{name, "mac_address", msg}
		}
		return locErr(lang, "mac_address", name, "")
	}
	return nil
}

func validateLatitude(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "latitude", msg, "only supported on strings")
	}
	if !latRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "latitude", msg}
		}
		return locErr(lang, "latitude", name, "")
	}
	val, _ := strconv.ParseFloat(fv.String(), 64)
	if val < -90 || val > 90 {
		if msg != "" {
			return ValError{name, "latitude", msg}
		}
		return locErr(lang, "latitude", name, "")
	}
	return nil
}

func validateLongitude(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "longitude", msg, "only supported on strings")
	}
	if !lonRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "longitude", msg}
		}
		return locErr(lang, "longitude", name, "")
	}
	val, _ := strconv.ParseFloat(fv.String(), 64)
	if val < -180 || val > 180 {
		if msg != "" {
			return ValError{name, "longitude", msg}
		}
		return locErr(lang, "longitude", name, "")
	}
	return nil
}

func validateCoordinate(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "coordinate", msg, "only supported on strings")
	}
	if !coordRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "coordinate", msg}
		}
		return locErr(lang, "coordinate", name, "")
	}
	parts := strings.Split(fv.String(), ",")
	lat, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lon, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		if msg != "" {
			return ValError{name, "coordinate", msg}
		}
		return locErr(lang, "coordinate", name, "")
	}
	return nil
}

func validateHexColor(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "hexcolor", msg, "only supported on strings")
	}
	if !hexColorRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "hexcolor", msg}
		}
		return locErr(lang, "hexcolor", name, "")
	}
	return nil
}

func validateCSSColor(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "css_color", msg, "only supported on strings")
	}
	if !cssColorRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "css_color", msg}
		}
		return locErr(lang, "css_color", name, "")
	}
	return nil
}

func validateRGB(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "rgb", msg, "only supported on strings")
	}
	if !rgbRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "rgb", msg}
		}
		return locErr(lang, "rgb", name, "")
	}
	return nil
}

func validateRGBA(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "rgba", msg, "only supported on strings")
	}
	if !rgbaRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "rgba", msg}
		}
		return locErr(lang, "rgba", name, "")
	}
	return nil
}

func validateHSL(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "hsl", msg, "only supported on strings")
	}
	if !hslRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "hsl", msg}
		}
		return locErr(lang, "hsl", name, "")
	}
	return nil
}

func validateHSLA(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "hsla", msg, "only supported on strings")
	}
	if !hslaRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "hsla", msg}
		}
		return locErr(lang, "hsla", name, "")
	}
	return nil
}

func validateJSON(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "json", msg, "only supported on strings")
	}
	s := fv.String()
	// Fast path: skip leading whitespace and check first char
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	if i >= len(s) || (s[i] != '{' && s[i] != '[') {
		if msg != "" {
			return ValError{name, "json", msg}
		}
		return locErr(lang, "json", name, "")
	}
	var js json.RawMessage
	if json.Unmarshal([]byte(s), &js) != nil {
		if msg != "" {
			return ValError{name, "json", msg}
		}
		return locErr(lang, "json", name, "")
	}
	return nil
}

func validateCreditCard(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "credit_card", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), " ", "")
	s = strings.ReplaceAll(s, "-", "")
	// Fast reject: credit cards are 13-19 digits
	if len(s) < 13 || len(s) > 19 {
		if msg != "" {
			return ValError{name, "credit_card", msg}
		}
		return locErr(lang, "credit_card", name, "")
	}
	if !creditCardRe.MatchString(s) {
		if msg != "" {
			return ValError{name, "credit_card", msg}
		}
		return locErr(lang, "credit_card", name, "")
	}
	sum := 0
	alt := false
	for i := len(s) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(s[i]))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	if sum%10 != 0 {
		if msg != "" {
			return ValError{name, "credit_card", msg}
		}
		return locErr(lang, "credit_card", name, "")
	}
	return nil
}

func validateBase64(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "base64", msg, "only supported on strings")
	}
	if _, err := base64.StdEncoding.DecodeString(fv.String()); err != nil {
		if _, err2 := base64.RawStdEncoding.DecodeString(fv.String()); err2 != nil {
			if msg != "" {
				return ValError{name, "base64", msg}
			}
			return locErr(lang, "base64", name, "")
		}
	}
	return nil
}

func validateHexadecimal(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "hexadecimal", msg, "only supported on strings")
	}
	if !hexadecimalRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "hexadecimal", msg}
		}
		return locErr(lang, "hexadecimal", name, "")
	}
	return nil
}

func validateISBN(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "isbn", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 13 && len(s) != 10 {
		if msg != "" {
			return ValError{name, "isbn", msg}
		}
		return locErr(lang, "isbn", name, "")
	}
	if len(s) == 13 {
		return validateISBN13(lang, name, fv, msg)
	}
	return validateISBN10(lang, name, fv, msg)
}

func validateISBN10(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "isbn10", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 10 {
		if msg != "" {
			return ValError{name, "isbn10", msg}
		}
		return locErr(lang, "isbn10", name, "")
	}
	sum := 0
	for i := 0; i < 9; i++ {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			if msg != "" {
				return ValError{name, "isbn10", msg}
			}
			return locErr(lang, "isbn10", name, "")
		}
		sum += n * (10 - i)
	}
	last := s[9]
	if last == 'X' || last == 'x' {
		sum += 10
	} else {
		n, err := strconv.Atoi(string(last))
		if err != nil {
			if msg != "" {
				return ValError{name, "isbn10", msg}
			}
			return locErr(lang, "isbn10", name, "")
		}
		sum += n
	}
	if sum%11 != 0 {
		if msg != "" {
			return ValError{name, "isbn10", msg}
		}
		return locErr(lang, "isbn10", name, "")
	}
	return nil
}

func validateISBN13(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "isbn13", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 13 {
		if msg != "" {
			return ValError{name, "isbn13", msg}
		}
		return locErr(lang, "isbn13", name, "")
	}
	sum := 0
	for i := 0; i < 13; i++ {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			if msg != "" {
				return ValError{name, "isbn13", msg}
			}
			return locErr(lang, "isbn13", name, "")
		}
		if i%2 == 0 {
			sum += n
		} else {
			sum += n * 3
		}
	}
	if sum%10 != 0 {
		if msg != "" {
			return ValError{name, "isbn13", msg}
		}
		return locErr(lang, "isbn13", name, "")
	}
	return nil
}

func validateISSN(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "issn", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 8 {
		if msg != "" {
			return ValError{name, "issn", msg}
		}
		return locErr(lang, "issn", name, "")
	}
	sum := 0
	for i := 0; i < 7; i++ {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			if msg != "" {
				return ValError{name, "issn", msg}
			}
			return locErr(lang, "issn", name, "")
		}
		sum += n * (8 - i)
	}
	last := s[7]
	if last == 'X' || last == 'x' {
		sum += 10
	} else {
		n, err := strconv.Atoi(string(last))
		if err != nil {
			if msg != "" {
				return ValError{name, "issn", msg}
			}
			return locErr(lang, "issn", name, "")
		}
		sum += n
	}
	if sum%11 != 0 {
		if msg != "" {
			return ValError{name, "issn", msg}
		}
		return locErr(lang, "issn", name, "")
	}
	return nil
}

func validateJWT(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "jwt", msg, "only supported on strings")
	}
	parts := strings.Split(fv.String(), ".")
	if len(parts) != 3 {
		if msg != "" {
			return ValError{name, "jwt", msg}
		}
		return locErr(lang, "jwt", name, "")
	}
	for _, part := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			if msg != "" {
				return ValError{name, "jwt", msg}
			}
			return locErr(lang, "jwt", name, "")
		}
	}
	return nil
}

func validateSemver(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "semver", msg, "only supported on strings")
	}
	if !semverRe.MatchString(fv.String()) {
		if msg != "" {
			return ValError{name, "semver", msg}
		}
		return locErr(lang, "semver", name, "")
	}
	return nil
}

func validateSSN(lang, name string, fv reflect.Value, msg string) error {
	if fv.Kind() != reflect.String {
		return customErr(lang, name, "ssn", msg, "only supported on strings")
	}
	s := strings.ReplaceAll(fv.String(), "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if !ssnRe.MatchString(s) {
		if msg != "" {
			return ValError{name, "ssn", msg}
		}
		return locErr(lang, "ssn", name, "")
	}
	return nil
}

func validateIsbn(lang, name string, fv reflect.Value, msg string) error {
	return validateISBN(lang, name, fv, msg)
}

func validateIsbn10(lang, name string, fv reflect.Value, msg string) error {
	return validateISBN10(lang, name, fv, msg)
}

func validateIsbn13(lang, name string, fv reflect.Value, msg string) error {
	return validateISBN13(lang, name, fv, msg)
}

func validateLatitudeDecimal(lang, name string, fv reflect.Value, msg string) error {
	return validateLatitude(lang, name, fv, msg)
}

func validateLongitudeDecimal(lang, name string, fv reflect.Value, msg string) error {
	return validateLongitude(lang, name, fv, msg)
}
