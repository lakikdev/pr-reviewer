package analyze

import (
	"context"
	"net/http"
	"time"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// swagger:model AnalyzeRequest
type AnalyzeRequest struct {
	Owner string `json:"owner"`
	Slug  string `json:"slug"`
	ID    string `json:"id"`
}

// swagger:model AnalyzeResponse
type AnalyzeResponse struct {
	Started bool `json:"started"`
}

/*
		swagger:route POST /analyze Analyze analyzeStart

		Analyze PR ----- [Admin]

		<b>Description</b><br>
		Starts analyzing the PR specified in request body

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters)

		<br><b>500: Internal Server Error</b><br>
			Err500xTransaction_001 - Error while connecting to database - (server internal error)
			Err500xTransaction_002 - Error while committing all changes - (server internal error)

		Parameters:
		+	name:        request
	   		description: PR details.
	   		in:          body
	   		type:        AnalyzeRequest

		Security:
			api_key: []

		Responses:
	  		201: AnalyzeResponse PR analysis response
			default: genericError
*/
func (api *API) Analyze(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData AnalyzeRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"owner":        requestData.Owner,
		"slug":         requestData.Slug,
		"id":           requestData.ID,
		"targetUserID": principal.UserID,
	})

	//begin db transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	//start background go routine to analyze PR with timeout to clean context
	//we don't use context now during analyze but we will in the future to call DB
	//so just pass it to the function for future use
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer cancel()
		log := logrus.WithFields(logrus.Fields{
			"owner": requestData.Owner,
			"slug":  requestData.Slug,
			"id":    requestData.ID,
		})

		commentsAdded, err := api.PRReviewer.Analyze(requestData.Owner, requestData.Slug, requestData.ID)
		if err != nil {
			log.Errorf("Error analyzing PR: %v", err)
			return
		}

		//log task completion
		log.Infof("PR analysis completed. Comments added: %d", commentsAdded)
	}(ctx, cancel)

	//commit db transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusCreated, &AnalyzeResponse{
		Started: true,
	})
}
