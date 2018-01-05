// Package texts provides string truncation.
package texts

import (
	"unicode"
	"unicode/utf8"
)

// Truncate shortens text until the length of the shortened text and appended
// suffix is equal to or less than maxLength. If exact is true, text is cut
// mid-word.
func Truncate(text string, maxLength int, suffix string, exact bool) string {
	if utf8.RuneCountInString(text) <= maxLength {
		return text
	}

	suffixLength := utf8.RuneCountInString(suffix)

	if suffixLength > maxLength {
		return ""
	}

	truncatedText := make([]rune, 0, maxLength)

	if exact {
		maxLengthWithoutSuffix := maxLength - suffixLength
		for _, character := range text {
			if len(truncatedText) == maxLengthWithoutSuffix {
				break
			}
			truncatedText = append(truncatedText, character)
		}
	} else {
		word := []rune{}
		for _, character := range text {
			word = append(word, character)

			if unicode.IsSpace(character) {
				if len(truncatedText)+len(word)+suffixLength > maxLength {
					break
				}
				truncatedText = append(truncatedText, word...)
				word = []rune{}
			}
		}
	}

	truncatedText = append(truncatedText, []rune(suffix)...)
	return string(truncatedText)
}
