package swearFilter

import (
	"strings"
	"unicode"
)

// Check if a string is free of any profanity
// returns true if no profanity is found, false otherwise
func Check(text string) bool {
	// Convert input to lowercase since our blacklist is all lowercase
	text = strings.ToLower(text)

	// Tokenize the input text
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	// Check for single-word profanities
	for _, word := range words {
		if containsProfanity(word) {
			return false
		}
	}

	// Check for multi-word profanities
	for i := 0; i < len(words); i++ {
		for j := i + 1; j <= len(words); j++ {
			phrase := strings.Join(words[i:j], " ")
			if containsProfanity(phrase) {
				return false
			}
		}
	}

	return true
}

// Helper function to check if a phrase contains any profane words or phrases
func containsProfanity(phrase string) bool {
	for _, swear := range swearBlackList {
		if phrase == swear {
			return true
		}
	}
	return false
}
