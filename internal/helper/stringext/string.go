package stringext

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

func New(text string) *string {
	return &text
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]+$`)
	return emailRegex.MatchString(e)
}

// NormalizeEmail standardizes the email address by converting it to lowercase
// and trimming any leading or trailing whitespace.
func NormalizeEmail(email string) string {
	// Convert the email to lowercase
	email = strings.ToLower(email)
	// Trim any leading or trailing whitespace
	email = strings.TrimSpace(email)
	return email
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate a random string of n length
func GenerateShortID(size int) (string, error) {
	b := make([]byte, size)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[num.Int64()]
	}
	return string(b), nil
}
