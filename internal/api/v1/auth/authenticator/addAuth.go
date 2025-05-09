package authenticator

import (
	"encoding/json"
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

func (auth *Authenticator) addAuth(c echo.Context, authType IAuthType) (*model.User, error) {
	cc := common.Convert(c)
	principal := authHelper.GetPrincipal(cc.Request())
	requestData := authType.GetRequestData()

	// Get data from request
	if err := cc.BindAndReset(requestData); err != nil {
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

	//if userAuth is not null meaning it used in another account
	if userAuth != nil {
		return nil, utils.NewHttpError(nil, http.StatusConflict, "Err409xAuthenticator_001", "Error while adding auth account - (auth option already in use)")
	}

	//load user by ID
	user, err := tx.User().GetByID(principal.UserID)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_001", "Error while getting user - (server internal error)")
	}

	//create new data struct for user auth
	authData := authType.BuildAuthData(requestData)

	jsonData, err := json.Marshal(authData)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_002", "Error while creating a user - (server internal error)")
	}

	// Link auth data with the user
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
