package validation

import "unicode"

func ValidateName(input string) bool {
	if len(input) < 1 {
		return false
	}
	for _, c := range input {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
