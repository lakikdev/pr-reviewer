package roles

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// swagger:model GetRolesResponse
type GetRolesRequest struct {
	UserID model.UserID `param:"userID"`
}

type GetRolesResponse struct {
	Roles []*model.UserRole `json:"roles"`
}

/*
		swagger:route GET /users/:userID/roles UserRoles userGetRolesByUser

		Get user's roles ----- [Admin]

		<b>Description</b><br>
		Get user's roles by ID. example role: ["admin"]

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters) <br>

		<br><b>500: Internal Server Error</b><br>
			Err500xRoles_001 - Error getting user's roles - (server internal error) <br>
			Err500xTransaction_001 - Error while connecting to database - (server internal error) <br>
			Err500xTransaction_002 - Error while committing all changes - (server internal error) <br>


		Parameters:
		+	name: userID
			description: User ID
			in: path
			type: string

		Security:
			api_key: []

		Responses:
	  		200: GetRolesResponse User's roles
			default: genericError
*/
func (api *API) GetRolesByUser(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData GetRolesRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"targetUserID": requestData.UserID,
		"UserID":       principal.UserID,
	})

	//begin db transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	roles, err := tx.UserRole().ListByUser(requestData.UserID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error getting user's roles - (server internal error)")
	}

	//commit db transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &GetRolesResponse{
		Roles: roles,
	})
}
