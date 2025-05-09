package session

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
	Save(session model.Session) error
	Get(data model.Session) (*model.Session, error)
	GetLatest(userID model.UserID) (*model.Session, error)
	ClearAllForUser(userID model.UserID) error
}
