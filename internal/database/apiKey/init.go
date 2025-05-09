package apiKey

import (
	"errors"

	"pr-reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

var ErrUserExists = errors.New("user with that auth_id already exists")

// ErrForgotPasswordBackoff is returned when requesting a password reset too frequently.
var ErrForgotPasswordBackoff = errors.New("last password reset request was sent too recently, please wait")

type DB struct {
	tx *sqlx.Tx
}

func New(tx *sqlx.Tx) Interface {
	return &DB{
		tx: tx,
	}
}

type Interface interface {
	Create(apiKey *model.APIKey) error
	Exists(hashedKey string) (bool, error)
}
