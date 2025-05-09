package udidAuth

import (
	"encoding/json"
	"net/http"
	"strings"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/api/v1/auth/authenticator"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/helper/stringext"
	"pr-reviewer/internal/model"
)

type UDIDAuth struct {
	name string
	DB   *database.DB
}

func New(db *database.DB) authenticator.IAuthType {
	return &UDIDAuth{
		name: "udid",
		DB:   db,
	}
}

func (auth *UDIDAuth) Name() string {
	return auth.name
}

func (auth *UDIDAuth) CreateIfNotExists() bool {
	return true
}

func (auth *UDIDAuth) GetRequestData() interface{} {
	return &UDIDAuthTypeRequest{}
}

func (auth *UDIDAuth) GetUserAuth(tx database.TxInterface, requestData interface{}) (*model.UserAuth, error) {
	data := requestData.(*UDIDAuthTypeRequest)

	// Check if UDID already in use
	userAuth, err := tx.UserAuth().GetByData(auth.Name(), "udid", data.UDID)
	if err != nil && err != dbHelper.ErrItemNotFound {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUDIDAuth_001", "Error getting user data - (server internal error)")
	}

	return userAuth, nil
}

func (auth *UDIDAuth) ValidateAuth(userAuth *model.UserAuth, requestData interface{}) error {
	data := requestData.(*UDIDAuthTypeRequest)

	var userAuthData UDIDAuthTypeJSON
	err := json.Unmarshal(userAuth.Data, &userAuthData)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUDIDAuth_002", "Error checking auth data - (server internal error)")
	}

	//It should be the same all the time since we selected by UDID,
	//Exists as example for other Auth Types and in case we will add additional validation conditions
	if data.UDID != nil && userAuthData.UDID != nil && *data.UDID != *userAuthData.UDID {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUDIDAuth_003", "Error checking auth data - (server internal error)")
	}

	return nil
}

func (auth *UDIDAuth) BuildAuthData(requestData interface{}) interface{} {
	data := requestData.(*UDIDAuthTypeRequest)
	return &UDIDAuthTypeJSON{
		UDID: stringext.New(strings.ToLower(*data.UDID)),
	}
}

func (auth *UDIDAuth) UpdateUser(user *model.User, requestData interface{}) error {
	return nil
}
