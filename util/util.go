package util

import (
	"unicode"
)

var regexChars = `\.+*?()|[]{}^$`
var whiteSpaceChars = `\t\n\r `

// Returns true if the string contains any regex characters
func ContainsRegex(str string) bool {
	for _, ch := range str {
		for _, reg := range regexChars {
			if ch == reg {
				return true
			}
		}
	}
	return false

}

func IsWhiteSpace(str string) bool {
	for _, ch := range str {

		if !unicode.IsSpace(ch) {
			return false
		}

	}
	// Reached the end without seeing a non-blank char
	return true
}
