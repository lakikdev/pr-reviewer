package model

import (
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DeviceID string

var NilDeviceID DeviceID

// Session is represent user's session
type Session struct {
	UserID           UserID     `db:"user_id"`
	DeviceID         DeviceID   `db:"device_id"`
	IPAddress        *string    `db:"ip_address"`
	OSVersion        *string    `db:"os_version"`
	AppVersion       *string    `db:"app_version"`
	DeviceType       *string    `db:"device_type"`
	RefreshToken     string     `db:"refresh_token"`
	RefreshTokenHash *[]byte    `db:"refresh_token_hash"`
	ExpiresAt        int64      `db:"expires_at"`
	LastUpdated      *time.Time `db:"last_updated"`
}

// SessionData used to represent data sent in json body with requests
type SessionData struct {
	DeviceID   DeviceID `json:"deviceID,omitempty"`
	OSVersion  *string  `json:"osVersion,omitempty"`
	AppVersion *string  `json:"appVersion,omitempty"`
	DeviceType *string  `json:"deviceType,omitempty"`
	IPAddress  *string  `json:"-"`
}

// Verify all required fields before create or update
func (u *SessionData) Verify() error {
	if len(u.DeviceID) == 0 {
		return errors.New("DeviceID is required")
	}

	return nil
}

// Set Password updates a user's password
func (u *Session) SetToken(token string) error {
	hash, err := HashToken(token)
	if err != nil {
		return err
	}
	u.RefreshTokenHash = &hash
	return nil
}

// CheckPassword verifies user's password
func (u *Session) CheckToken(token string) error {
	if u.RefreshTokenHash != nil && len(*u.RefreshTokenHash) == 0 {
		return errors.New("token not set")
	}
	//hash token before compared since Token in DB was hashed before storage
	hash := sha256.Sum256([]byte(token))
	return bcrypt.CompareHashAndPassword(*u.RefreshTokenHash, hash[:])
}

// HashPassword hashes a user's raw password
func HashToken(token string) ([]byte, error) {
	//hash before use bcrypt to avoid 72 bytes limit
	hash := sha256.Sum256([]byte(token))
	return bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
}

type OSType string

var (
	IOS     OSType = "iOS"
	Android OSType = "Android"
	OtherOS OSType = "Other"
)

// Get OS type extracted from OS Version
func (u *Session) GetOS() OSType {
	if u.OSVersion != nil && len(*u.OSVersion) > 0 {
		switch s := strings.Split(*u.OSVersion, " "); s[0] {
		case "iOS", "iPadOS":
			return IOS
		case "Android":
			return Android
		}
	}

	return OtherOS
}
