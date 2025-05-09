package userAuth

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
	GetByData(authType string, fieldName string, fieldValue interface{}) (*model.UserAuth, error)
	AddToUser(userAuth *model.UserAuth) error
}
