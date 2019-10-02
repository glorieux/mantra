package strings

import "unicode"

// HasAlphaNumeric returns wether or not a string only contain alphanumeric
// characters
func HasAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
