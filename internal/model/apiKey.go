package model

import (
	"time"
)

type APIKeyID string

type APIKey struct {
	ID      APIKeyID `db:"id"`
	Name    *string  `db:"name"`
	KeyHash *string  `db:"key_hash"`

	ExpiresAt *time.Time `db:"expires_at"`
	Active    *bool      `db:"active"`

	CreatedAt *time.Time `db:"created_at"`
}
