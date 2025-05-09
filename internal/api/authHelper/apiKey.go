package authHelper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateAPIKey() (raw string, hashed string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", "", err
	}
	rawKey := base64.URLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(rawKey))
	return rawKey, base64.URLEncoding.EncodeToString(hash[:]), nil
}

func HashAPIKey(raw string) (hashed string, err error) {
	hash := sha256.Sum256([]byte(raw))
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}
