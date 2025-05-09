package database

import (
	"pr-reviewer/internal/database/apiKey"
	"pr-reviewer/internal/database/session"
	"pr-reviewer/internal/database/user"
	"pr-reviewer/internal/database/userAuth"
	"pr-reviewer/internal/database/userRole"

	"github.com/jmoiron/sqlx"
)

//Create transaction struct to be used in all database operations

type TxInterface interface {
	User() user.Interface
	UserRole() userRole.Interface
	UserAuth() userAuth.Interface
	Session() session.Interface
	APIKey() apiKey.Interface

	Commit() error
	Rollback() error
}

type Tx struct {
	tx *sqlx.Tx

	user     user.Interface
	userRole userRole.Interface
	userAuth userAuth.Interface
	session  session.Interface
	apiKey   apiKey.Interface
}

// commit transaction
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// rollback transaction
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) User() user.Interface {
	if t.user == nil {
		t.user = user.New(t.tx)
	}
	return t.user
}

func (t *Tx) UserRole() userRole.Interface {
	if t.userRole == nil {
		t.userRole = userRole.New(t.tx)
	}
	return t.userRole
}

func (t *Tx) UserAuth() userAuth.Interface {
	if t.userAuth == nil {
		t.userAuth = userAuth.New(t.tx)
	}
	return t.userAuth
}

func (t *Tx) Session() session.Interface {
	if t.session == nil {
		t.session = session.New(t.tx)
	}
	return t.session
}

func (t *Tx) APIKey() apiKey.Interface {
	if t.apiKey == nil {
		t.apiKey = apiKey.New(t.tx)
	}
	return t.apiKey
}
