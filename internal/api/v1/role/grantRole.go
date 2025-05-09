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

// swagger:model GrantRoleRequest
type GrantRoleRequest struct {
	Role   model.Role   `json:"role"`
	UserID model.UserID `param:"userID"`
}

// swagger:model GrantRoleResponse
type GrantRoleResponse struct {
	Granted bool              `json:"granted"`
	Roles   []*model.UserRole `json:"roles"`
}

/*
		swagger:route POST /users/:userID/roles UserRoles userRoleGrant

		Grant role to user ----- [Admin]

		<b>Description</b><br>
		Grant role to user by ID. example role: ["admin"]

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters)

		<br><b>500: Internal Server Error</b><br>
			Err500xRoles_001 - Error granting role - (server internal error)
			Err500xRoles_002 - Error getting user's roles - (server internal error)
			Err500xTransaction_001 - Error while connecting to database - (server internal error)
			Err500xTransaction_002 - Error while committing all changes - (server internal error)

		Parameters:
		+	name: userID
			description: User ID
			in: path
			type: string
		+	name:        request
	   		description: Role grant.
	   		in:          body
	   		type:        GrantRoleRequest

		Security:
			api_key: []

		Responses:
	  		201: GrantRoleResponse List of user's roles
			default: genericError
*/
func (api *API) GrantRole(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData GrantRoleRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID":       requestData.UserID,
		"targetUserID": principal.UserID,
	})

	//begin db transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	//store role in database
	if err := tx.UserRole().Grant(requestData.UserID, requestData.Role); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error granting role - (server internal error)")
	}

	cc.Logrus().WithField("role", requestData.Role)

	//clear user's roles from cache
	_ = api.Permissions.RemoveFromCache(requestData.UserID)

	//get user's roles from database
	roles, err := tx.UserRole().ListByUser(requestData.UserID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_002", "Error getting user's roles - (server internal error)")
	}

	//commit db transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusCreated, &GrantRoleResponse{
		Granted: true,
		Roles:   roles,
	})
}
