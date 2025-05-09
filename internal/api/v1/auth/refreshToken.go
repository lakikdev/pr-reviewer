package auth

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// swagger:model RefreshTokenRequest
type RefreshTokenRequest struct {
	model.SessionData

	RefreshToken       string `json:"refreshToken"`
	WithPendingRequest *bool  `json:"withPendingRequest,omitempty"`
}

// swagger:model RefreshTokenResponse
type RefreshTokenResponse struct {
	authHelper.Tokens

	UserID model.UserID `json:"userID"`
	Roles  []model.Role `json:"roles,omitempty"`
}

/*
		swagger:route POST /refreshToken Auth refreshToken

		Refresh access tokens

		<b>Description</b><br>
		Returns access and refresh tokens

	    <br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters) <br>
			Err400xBind_002 - Error while verifying if request data has all required fields - (not all fields found) <br>

		<br><b>401: Unauthorized</b><br>
			Err401xToken_001 - Error verifying refresh token - (token invalid) <br>
			Err401xSession_001 - Error session not exist - (session not found) <br>
			Err401xSession_002 - Error session not exist - (session not found) <br>

		<br><b>404: Not Found</b><br>
			Err404xUser_001 - Error while getting user by ID - (user not found) <br>

		<br><b>500: Internal Server Error</b><br>
			Err500xTransaction_001 - Error while connecting to database - (server internal error) <br>
			Err500xTransaction_002 - Error while committing all changes - (server internal error) <br>
			Err500xUser_001 - Error while getting user by ID - (server internal error) <br>
			Err500xToken_001 - Error while trying to issue access and refresh tokens - (server internal error) <br>
			Err500xRole_001 - Error while getting user roles - (server internal error) <br>

		Parameters:
		+	name: request
			description: Refresh token request data
			in: body
			type: RefreshTokenRequest

	 	Responses:
	  		200: RefreshTokenResponse token refreshed
			default: genericError
*/
func (api *API) RefreshToken(c echo.Context) error {
	cc := common.Convert(c)
	ctx := c.Request().Context()

	var requestData RefreshTokenRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	//verify if request data has all required fields
	if err := requestData.Verify(); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_002", "Error while verifying if request data has all required fields - (not all fields found)")
	}

	//set User's IP address on session data
	if err := api.setIPAddress(c.Request(), &requestData.SessionData); err != nil {
		return err
	}

	cc.Logrus().WithField("deviceID", requestData.DeviceID)

	principal, err := authHelper.VerifyToken(requestData.RefreshToken)
	if err != nil {
		return utils.NewHttpError(err, http.StatusUnauthorized, "Err401xToken_001", "Error verifying refresh token - (token invalid)")
	}

	cc.Logrus().WithField("targetUserID", principal.UserID)

	//if token is valid we need to check if combination UserID - DeviceID - Refresh Token exists in database

	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	data := model.Session{
		UserID:   principal.UserID,
		DeviceID: requestData.DeviceID,
	}

	session, err := tx.Session().Get(data)
	if err != nil || session == nil {
		return utils.NewHttpError(err, http.StatusUnauthorized, "Err401xSession_001", "Error session not exist - (session not found)")
	}

	if err := session.CheckToken(requestData.RefreshToken); err != nil {
		return utils.NewHttpError(err, http.StatusUnauthorized, "Err401xSession_002", "Error session not exist - (session not found)")
	}

	//Check if user exists
	user, err := tx.User().GetByID(principal.UserID)
	if err != nil {
		if err == dbHelper.ErrItemNotFound {
			return utils.NewHttpError(err, http.StatusNotFound, "Err404xUser_001", "Error while getting user by ID - (user not found)")
		}

		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error while getting user by ID - (server internal error)")
	}

	tokens, err := api.issueTokens(tx, principal.UserID, requestData.SessionData)
	if err != nil || tokens == nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xToken_001", "Error while trying to issue access and refresh tokens - (server internal error)")
	}

	roles, err := api.Permissions.GetRoles(user.ID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRole_001", "Error while getting user roles - (server internal error)")
	}

	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &RefreshTokenResponse{
		Tokens: *tokens,
		UserID: principal.UserID,
		Roles:  roles,
	})
}

// Verify all required fields before create or update
func (r RefreshTokenRequest) Verify() error {
	if err := r.SessionData.Verify(); err != nil {
		return err
	}

	if len(r.RefreshToken) == 0 {
		return errors.New("refreshToken is required")
	}

	return nil
}
