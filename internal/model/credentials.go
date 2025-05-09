package model

import (
	"fmt"
)

// Credentials used in login API
type Credentials struct {
	SessionData

	Email    string `json:"email"`
	Password string `json:"password"`
}

// Principal is an authenticated entity
type Principal struct {
	UserID      UserID `json:"userID,omitempty"`
	Role        *Role  `json:"role,omitempty"`
	ValidAPIKey bool   `json:"validAPIKey,omitempty"`
}

// NilPrincipal is an uninitialized Principal
var NilPrincipal Principal

func (p Principal) String() string {
	if p.UserID != NilUserID {
		return fmt.Sprintf("UserID[%s]", p.UserID)
	}
	return "(none)"
}
