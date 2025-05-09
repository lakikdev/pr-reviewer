package model

import (
	"encoding/json"
	"time"
)

type AuthID string

type UserAuth struct {
	AuthID AuthID `json:"authID" db:"auth_id"`
	UserID UserID `json:"userID" db:"user_id"`
	Type   string `json:"type" db:"type"`

	Data json.RawMessage `json:"data,omitempty" db:"data"`

	Active  *bool      `json:"active,omitempty" db:"active"`
	AddedAt *time.Time `json:"addedAt" db:"added_at"`
}
