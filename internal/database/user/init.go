package user

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
	Create(user *model.User) error
	Update(user *model.User) error
	UpdatePassword(user *model.User) error
	GetByID(userID model.UserID) (*model.User, error)
	List(param model.ListDataParameters) (users []*model.User, total *int64, err error)
	Delete(userID model.UserID) error
	PermanentDelete(userID model.UserID) (bool, error)
}
