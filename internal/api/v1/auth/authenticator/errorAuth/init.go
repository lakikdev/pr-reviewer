package errorAuth

import (
	"errors"

	"pr-reviewer/internal/api/v1/auth/authenticator"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/model"
)

// errorAuth is used to catch errors on AuthType creations.
// if during creation and adding Auth Type to authenticator error was thrown we are creating errorAuth instead
// and in AddAuthType function catching that Auth type
type ErrorAuth struct {
	error
}

func New(message string) authenticator.IAuthType {
	return &ErrorAuth{errors.New(message)}
}

func (auth *ErrorAuth) Name() string { return "unknown" }
func (auth *ErrorAuth) CreateIfNotExists() bool {
	return true
}
func (auth *ErrorAuth) GetRequestData() interface{} {
	return nil
}

func (auth *ErrorAuth) GetUserAuth(tx database.TxInterface, requestData interface{}) (*model.UserAuth, error) {
	return nil, nil
}

func (auth *ErrorAuth) ValidateAuth(userAuth *model.UserAuth, requestData interface{}) error {
	return nil
}

func (auth *ErrorAuth) BuildAuthData(requestData interface{}) interface{} {
	return nil
}

func (auth *ErrorAuth) UpdateUser(user *model.User, requestData interface{}) error {
	return nil
}
