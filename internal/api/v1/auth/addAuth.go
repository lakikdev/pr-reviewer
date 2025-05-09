package auth

import (
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// swagger:model SignInRequest
type AddAuthRequest struct {
	model.SessionData

	Type string `json:"type"`
}

// swagger:model SignInResponse
type AddAuthResponse struct {
	Added bool `json:"added"`
}

/*
		swagger:route POST /migrate Auth authSignin

		Login user

		<b>Description</b><br>
		Login user using UDID, Credentials, Sequence Web3 Wallet

		If "type" was set to "udid" - endpoint will expect fields 'udid' to be provided
		If "type" was set to "credentials" - endpoint will expect fields 'email' and 'password' to be provided
		If "type" was set to "sequenceWallet" - endpoint will expect fields 'walletAddress' and 'proof' to be provided
		If "type" was set to "sequenceWaasWallet" - endpoint will expect fields 'token' to be provided

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
			Err500xAuthenticatorUser_001 - Error while getting user - (server internal error) <br>
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
			Err500xSequenceWaaSAuth_001 - Error validating token - (server internal error) <br>
			Err500xSequenceWaaSAuth_002 - Error validating token - (token is invalid) <br>
			Err500xSequenceWaaSAuth_003 - Error getting user data - (server internal error) <br>
			Err500xSequenceWaaSAuth_004 - Error checking auth data - (server internal error) <br>
			Err500xSequenceWaaSAuth_005	- Error checking auth data - (server internal error) <br>

		Parameters:
		+	name: request
			description: Sign In request data
			in: body
			type: SignInRequest

	 	Responses:
	  		200: SignInResponse User is logged in successfully, returns tokens with userId
			default: genericError
*/
func (api *API) AddAuth(c echo.Context) error {
	cc := common.Convert(c)

	var requestData AddAuthRequest
	if err := cc.BindAndReset(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	if err := requestData.Verify(); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_002", "Error while verifying if request data has all required fields - (not all fields found)")
	}

	cc.Logrus().WithField("type", requestData.Type)

	_, err := api.Authenticator.AddAuth(c, requestData.Type)
	if err != nil {
		return err
	}

	return utils.NewWriter().WriteJSON(c, http.StatusCreated, &AddAuthResponse{
		Added: true,
	})
}

// Verify all required fields before create or update
func (r AddAuthRequest) Verify() error {
	if err := r.SessionData.Verify(); err != nil {
		return err
	}

	if len(r.Type) == 0 {
		return errors.New("type is required")
	}

	return nil
}
