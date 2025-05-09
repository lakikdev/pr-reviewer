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

// swagger:model RemoveRoleRequest
type RemoveRoleRequest struct {
	Role   model.Role   `json:"role"`
	UserID model.UserID `param:"userID"`
}

// swagger:model RemoveRoleResponse
type RemoveRoleResponse struct {
	Removed          bool              `json:"removed"`
	RemovedFromCache bool              `json:"removedFromCache"`
	Roles            []*model.UserRole `json:"roles"`
}

/*
		swagger:route DELETE /users/userID/roles UserRoles userRoleRemove

		Remove role from user ----- [Admin]

		<b>Description</b><br>
		Remove role from user by ID. example role: ["admin"]

		<br><b>400: Bad Request</b><br>
			Err400x001 - Error while decoding request data to JSON - (could not decode parameters)

		<br><b>500: Internal Server Error</b><br>
			Err500x001 - Error while removing user's role

		Parameters:
		+	name: userID
			description: User ID
			in: path
			type: string
		+	name:        request
	   		description: Role grant.
	   		in:          body
	   		type:        RemoveRoleRequest

		Security:
			api_key: []

		Responses:
	  		200: RemoveRoleResponse Role is granted successfully
			default: genericError
*/
func (api *API) RemoveRole(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData RemoveRoleRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400x001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID":       requestData.UserID,
		"targetUserID": principal.UserID,
	})

	//begin transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	//remove role from database
	if err := tx.UserRole().Revoke(requestData.UserID, requestData.Role); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error while removing user's role - (server internal error)")

	}

	//get user's roles from database
	roles, err := tx.UserRole().ListByUser(requestData.UserID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_002", "Error while getting user's roles - (server internal error)")
	}

	//commit transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	deletedFromCache := api.Permissions.RemoveFromCache(requestData.UserID)

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &RemoveRoleResponse{
		Removed:          true,
		RemovedFromCache: deletedFromCache,
		Roles:            roles,
	})
}
