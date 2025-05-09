package auth

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// swagger:model SignUpRequest
type SignUpRequest struct {
	model.SessionData

	Type string `json:"type"`
}

// swagger:model SignUpResponse
type SignUpResponse struct {
	authHelper.Tokens

	UserID model.UserID `json:"userID"`
	Roles  []model.Role `json:"roles,omitempty"`
}

/*
		swagger:route POST /signup Auth authSignup

		Create user

		<b>Description</b><br>
		Create user using UDID, Credentials, Sequence Web3 Wallet

		If "type" was set to "udid" - endpoint will expect fields 'udid' to be provided
		If "type" was set to "credentials" - endpoint will expect fields 'email' and 'password' to be provided
		If "type" was set to "sequenceWallet" - endpoint will expect fields 'walletAddress' and 'proof' to be provided

	 	<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters) <br>
			Err400xBind_002 - Error while verifying if request data has all required fields - (not all fields found) <br>
			Err400xAuthenticatorBind_001 - Error while decoding request data to JSON - (could not decode parameters) <br>
			Err400xAuthenticatorBind_002 - Error while verifying if request data has all required fields - (not all fields found) <br>
			Err400xAuthenticator_001 - Error while verifying if request data has all required fields - (unknown auth type) <br>

		<br><b>403: Forbidden</b><br>
			Err403xCredentialsAuth_001 - Error checking auth data - (invalid credentials) <br>


		<br><b>409: Conflict</b><br>
			Err409xAuthenticatorUser_001 - Error creating user - (auth option already in use) <br>

		<br><b>500: Internal Server Error</b><br>
			Err500xIssueToken_001 - Error while trying to issue access and refresh tokens - (server internal error) <br>
			Err500xRoles_001 - Error while getting user roles - (server internal error) <br>
			Err500xTransaction_001 - Error while connecting to database - (server internal error) <br>
			Err500xTransaction_002 - Error while committing all changes - (server internal error) <br>
			Err500xAuthenticatorTransaction_001 - Error while connecting to database - (server internal error) <br>
			Err500xAuthenticatorTransaction_002 - Error while committing all changes - (server internal error) <br>
			Err500xAuthenticatorUser_001 - Error while creating a user - (server internal error) <br>
			Err500xAuthenticatorUser_002 - Error while creating a user - (server internal error) <br>
			Err500xAuthenticatorUser_003 - Error while adding auth option for user - (server internal error) <br>
			Err500xUDIDAuth_001 - Error getting user data - (server internal error) <br>
			Err500xUDIDAuth_002 - Error checking auth data - (server internal error) <br>
			Err500xUDIDAuth_003 - Error checking auth data - (server internal error) <br>
			Err500xCredentialsAuth_001 - Error getting user data - (server internal error) <br>
			Err500xCredentialsAuth_002 - Error checking auth data - (server internal error) <br>
			Err500xCredentialsAuth_003 - Error checking auth data - (server internal error) <br>
			Err500xCredentialsAuth_004 - Error while setting password - (server internal error) <br>
			Err500xSequenceAuth_001 - Error validating wallet proof - (server internal error) <br>
			Err500xSequenceAuth_002 - Error validating wallet proof - (proof is invalid) <br>
			Err500xSequenceAuth_003 - Error getting user data - (server internal error) <br>
			Err500xSequenceAuth_004 - Error checking auth data - (server internal error) <br>
			Err500xSequenceAuth_005 - Error checking auth data - (server internal error) <br>

		Parameters:
		+	name: request
			description: Sign Up request data
			in: body
			type: SignUpRequest

		Responses:
	  	201: SignUpResponse User is created successfully
		default: genericError
*/

func (api *API) SignUp(c echo.Context) error {
	cc := common.Convert(c)
	ctx := c.Request().Context()

	var requestData SignUpRequest
	if err := cc.BindAndReset(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	if err := requestData.Verify(); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_002", "Error while verifying if request data has all required fields - (not all fields found)")
	}

	//set User's IP address on session data
	if err := api.setIPAddress(c.Request(), &requestData.SessionData); err != nil {
		return err
	}

	cc.Logrus().WithField("type", requestData.Type)

	user, err := api.Authenticator.SignUp(c, requestData.Type)
	if err != nil {
		return err
	}

	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	tokens, err := api.issueTokens(tx, user.ID, requestData.SessionData)
	if err != nil || tokens == nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xIssueToken_001", "Error while trying to issue access and refresh tokens - (server internal error)")
	}

	roles, err := api.Permissions.GetRoles(user.ID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error while getting user roles - (server internal error)")
	}

	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	cc.Logrus().WithField("userID", user.ID)

	return utils.NewWriter().WriteJSON(c, http.StatusCreated, &SignUpResponse{
		Tokens: *tokens,
		UserID: user.ID,
		Roles:  roles,
	})
}

// Verify all required fields before create or update
func (r *SignUpRequest) Verify() error {
	if err := r.SessionData.Verify(); err != nil {
		return err
	}

	if len(r.Type) == 0 {
		return errors.New("type is required")
	}

	return nil
}
