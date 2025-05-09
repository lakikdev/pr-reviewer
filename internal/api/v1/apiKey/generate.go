package apiKey

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// swagger:model GenerateRequest
type GenerateRequest struct {
	Name string `json:"name"`
}

// swagger:model GenerateResponse
type GenerateResponse struct {
	Key string `json:"key"`
}

/*
		swagger:route POST /api-key/generate APIKey apiKeyGenerate

		Generate API KEY ----- [API Key]

		<b>Description</b><br>
		Generate API KEY

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters)

		<br><b>500: Internal Server Error</b><br>
			Err500xAPIkey_001 - Error generating API Key - (server internal error)
			Err500xAPIkey_002 - Error storing API Key - (server internal error)
			Err500xTransaction_001 - Error while connecting to database - (server internal error)
			Err500xTransaction_002 - Error while committing all changes - (server internal error)

		Parameters:
		+	name:        request
	   		description: API Key Data.
	   		in:          body
	   		type:        GenerateRequest

		Security:
			api_key: []

		Responses:
	  		201: GenerateResponse API key response
			default: genericError
*/
func (api *API) Generate(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData GenerateRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"targetUserID": principal.UserID,
	})

	//begin db transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	//generate new API Key
	rawKey, hashed, err := authHelper.GenerateAPIKey()
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAPIkey_001", "Error generating API Key - (server internal error)")
	}

	apiKey := &model.APIKey{
		Name:    &requestData.Name,
		KeyHash: &hashed,
	}

	//store role in database
	if err := tx.APIKey().Create(apiKey); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xAPIkey_002", "Error storing API Key - (server internal error)")
	}

	//commit db transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusCreated, &GenerateResponse{
		Key: rawKey,
	})
}
