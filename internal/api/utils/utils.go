package utils

import (
	"crypto/rand"

	"pr-reviewer/internal/modules/swearFilter"
)

const (
	shortIDCharacters      = "ABCDEFGHIJKLMNOPRSTUVWXYZ"
	referralCodeCharacters = "ABCDEFGHIJKLMNOPRSTUVWXYZ123456789"
)

var (
	shortIDCharacterCount      = len(shortIDCharacters)
	referralCodeCharacterCount = len(referralCodeCharacters)
	shortIDBlackList           = map[string]struct{}{
		"FUCK": {},
		"POOP": {},
		"SHIT": {},
	}
)

// GenerateRandomID generates a random ID
func GenerateRandomID(len int) string {
	randBytes := make([]byte, shortIDCharacterCount)
	_, _ = rand.Read(randBytes)

	for {
		randomShortID := make([]byte, len)
		for i := range randomShortID {
			randomShortID[i] = shortIDCharacters[int(randBytes[i])%shortIDCharacterCount]
		}
		shortID := string(randomShortID)
		if _, blacklisted := shortIDBlackList[shortID]; blacklisted {
			continue
		}
		return shortID
	}
}

// GenerateReferralCode generates a random referral code
func GenerateReferralCode(len int) string {
	randBytes := make([]byte, referralCodeCharacterCount)

	for {
		_, _ = rand.Read(randBytes)
		randomCode := make([]byte, len)
		for i := range randomCode {
			randomCode[i] = referralCodeCharacters[int(randBytes[i])%referralCodeCharacterCount]
		}
		code := string(randomCode)
		if clean := swearFilter.Check(code); !clean {
			continue
		}
		return code
	}
}
