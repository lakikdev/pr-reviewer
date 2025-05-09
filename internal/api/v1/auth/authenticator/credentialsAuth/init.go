package credentialsAuth

import (
	"encoding/json"
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/api/v1/auth/authenticator"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/helper/stringext"
	"pr-reviewer/internal/model"
)

type CredentialsAuth struct {
	name string
	DB   *database.DB
}

func New(db *database.DB) authenticator.IAuthType {
	return &CredentialsAuth{
		name: "credentials",
		DB:   db,
	}
}

func (auth *CredentialsAuth) Name() string {
	return auth.name
}

func (auth *CredentialsAuth) CreateIfNotExists() bool {
	return false
}

func (auth *CredentialsAuth) GetRequestData() interface{} {
	return &CredentialsAuthTypeRequest{}
}

func (auth *CredentialsAuth) GetUserAuth(tx database.TxInterface, requestData interface{}) (*model.UserAuth, error) {
	data := requestData.(*CredentialsAuthTypeRequest)
	// Check if Email already in use
	userAuth, err := tx.UserAuth().GetByData(auth.Name(), "email", data.Email)
	if err != nil && err != dbHelper.ErrItemNotFound {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xCredentialsAuth_001", "Error getting user data - (server internal error)")
	}

	return userAuth, nil
}

func (auth *CredentialsAuth) ValidateAuth(userAuth *model.UserAuth, requestData interface{}) error {
	data := requestData.(*CredentialsAuthTypeRequest)

	var userAuthData CredentialsAuthTypeJSON
	err := json.Unmarshal(userAuth.Data, &userAuthData)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xCredentialsAuth_002", "Error checking auth data - (server internal error)")
	}

	//It should be the same all the time since we selected by Email,
	//Exists as example for other Auth Types and in case we will add additional validation conditions
	if data.Email != nil && userAuthData.Email != nil && *data.Email != *userAuthData.Email {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xCredentialsAuth_003", "Error checking auth data - (server internal error)")
	}

	//Checking stored hashed password if it match provided password
	if err := userAuthData.CheckPassword(*data.Password); err != nil {
		return utils.NewHttpError(err, http.StatusForbidden, "Err403xCredentialsAuth_001", "Error checking auth data - (invalid credentials)")
	}

	return nil
}

func (auth *CredentialsAuth) BuildAuthData(requestData interface{}) interface{} {
	data := requestData.(*CredentialsAuthTypeRequest)

	credentialsData := CredentialsAuthTypeJSON{
		Email: data.Email,
	}

	if err := credentialsData.SetPassword(*data.Password); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xCredentialsAuth_004", "Error while setting password - (server internal error)")
	}

	return &credentialsData
}

func (auth *CredentialsAuth) UpdateUser(user *model.User, requestData interface{}) error {
	data := requestData.(*CredentialsAuthTypeRequest)
	user.Email = stringext.New(*data.Email)
	return nil
}
