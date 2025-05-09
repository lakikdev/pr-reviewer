package authenticator

import (
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/model"
)

type AuthType string

type IAuthType interface {
	CreateIfNotExists() bool
	GetRequestData() interface{}
	GetUserAuth(tx database.TxInterface, requestData interface{}) (*model.UserAuth, error)
	ValidateAuth(userAuth *model.UserAuth, requestData interface{}) error
	BuildAuthData(requestData interface{}) interface{}
	UpdateUser(user *model.User, requestData interface{}) error
	Name() string
}

type IRequestData interface {
	Verify() error
}
