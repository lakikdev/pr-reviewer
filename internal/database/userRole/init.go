package userRole

import (
	"pr-reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	tx *sqlx.Tx
}

func New(tx *sqlx.Tx) Interface {
	return &DB{
		tx: tx,
	}
}

type Interface interface {
	Grant(userID model.UserID, role model.Role) error
	Revoke(userID model.UserID, role model.Role) error
	ListByUser(userID model.UserID) ([]*model.UserRole, error)
}
