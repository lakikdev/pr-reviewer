package authenticator

import (
	"encoding/json"
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

func (auth *Authenticator) signUp(c echo.Context, authType IAuthType) (*model.User, error) {

	requestData := authType.GetRequestData()

	// Get data from request
	if err := c.Bind(requestData); err != nil {
		return nil, utils.NewHttpError(err, http.StatusBadRequest, "Err400xAuthenticatorBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	if d, ok := requestData.(IRequestData); ok {
		if err := d.Verify(); err != nil {
			return nil, utils.NewHttpError(err, http.StatusBadRequest, "Err400xAuthenticatorBind_002", "Error while decoding request data to JSON - (data invalid)")
		}
	}

	tx, err := auth.DB.BeginTxx(c.Request().Context())
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	// Check if auth option already in use
	userAuth, err := authType.GetUserAuth(tx, requestData)
	if err != nil {
		return nil, err
	}

	if userAuth != nil {
		return nil, utils.NewHttpError(err, http.StatusConflict, "Err409xAuthenticatorUser_001", "Error creating user - (auth option already in use)")
	}

	user := &model.User{}

	//Set required fields on user object, for now used only by SequenceAuth
	//since it just set values on user object no error returned for now
	//TODO Add proper error handling if we will have some type cast error
	_ = authType.UpdateUser(user, requestData)

	// If not create empty user
	if err := tx.User().Create(user); err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_001", "Error while creating a user - (server internal error)")
	}

	//create new data struct for user auth
	authData := authType.BuildAuthData(requestData)

	jsonData, err := json.Marshal(authData)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_002", "Error while creating a user - (server internal error)")
	}

	// Link UDID with the user
	newUserAuth := &model.UserAuth{
		UserID: user.ID,
		Type:   authType.Name(),
		Data:   jsonData,
	}

	if err := tx.UserAuth().AddToUser(newUserAuth); err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_003", "Error while adding auth option for user - (server internal error)")
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return user, nil
}
