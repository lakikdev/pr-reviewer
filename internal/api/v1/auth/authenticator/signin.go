package authenticator

import (
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

func (auth *Authenticator) signIn(c echo.Context, authType IAuthType) (*model.User, error) {
	cc := common.Convert(c)
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

	//If auth type not exists and CreateIfNotExists is true we send all info to signUp function to handle creation
	//and stop current function
	if userAuth == nil && authType.CreateIfNotExists() {
		return auth.signUp(c, authType)
	} else if userAuth == nil {
		return nil, utils.NewHttpError(err, http.StatusNotFound, "Err404xAuthenticatorData_001", "Error checking auth data - (invalid credentials)")
	}

	if err := authType.ValidateAuth(userAuth, requestData); err != nil {
		return nil, err
	}

	user, err := tx.User().GetByID(userAuth.UserID)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorUser_001", "Error while getting user - (server internal error)")
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAuthenticatorTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return user, nil

}
