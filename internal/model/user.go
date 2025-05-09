package model

import (
	"strings"
	"time"
)

// UserID is identifier for User
type UserID string

// NilUserID is an empty UserID
var NilUserID UserID

// User is structure represent User object
type User struct {
	ID UserID `json:"userID,omitempty" db:"user_id"`

	DisplayName       *string `json:"displayName" db:"display_name"`
	Email             *string `json:"email" db:"email"`
	ProfileColorIndex *int    `json:"profileColorIndex,omitempty" db:"profile_color_index"`

	SequenceWallet        *string `json:"sequenceWallet,omitempty" db:"sequence_wallet"`
	StripeCustomerID      *string `json:"stripeCustomerID,omitempty" db:"stripe_customer_id"`
	StripePaymentMethodID *string `json:"stripePaymentMethodID,omitempty" db:"stripe_payment_method_id"`

	CreatedAt *time.Time `json:"createdAt" db:"created_at"`

	Roles     []Role `json:"roles,omitempty"`
	IsCreator bool   `json:"isCreator,omitempty" db:"is_creator"`
}

// Verify all required fields before create or update
func (u *User) Verify() error {

	if u.DisplayName != nil {
		//trim spaces
		*u.DisplayName = strings.TrimSpace(*u.DisplayName)
	}

	return nil
}
