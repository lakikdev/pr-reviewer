package user

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DeleteUserRequest struct {
	UserID model.UserID `param:"userID"`
}

// swagger:model DeleteUserResponse
type DeleteUserResponse struct {
	Deleted bool `json:"deleted"`
}

/*
		swagger:route DELETE /users/:userID User userDelete

		Delete user ----- [Admin, MemberIsTarget]

		<b>Description</b><br>
		Deletes user by ID

		<br><b>500: Internal Server Error</b><br>
			Err500xUser_001 - Error deleting user - (internal server error)
			Err500xTransaction_001 - Error while connecting to database - (server internal error)
			Err500xTransaction_002 - Error while committing all changes - (server internal error)

		Parameters:
		+	name: id
			description: User ID
			in: path
			type: string

		Security:
			api_key: []

		Responses:
	  		200: DeleteUserResponse User is deleted successfully
			default: genericError
*/
func (api *API) Delete(c echo.Context) error {
	cc := common.Convert(c)

	var requestData DeleteUserRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID":       principal.UserID,
		"targetUserID": requestData.UserID,
	})

	//begin transaction
	tx, err := api.DataRepository.DB.BeginTxx(cc.Request().Context())
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	if err := tx.User().Delete(requestData.UserID); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error deleting user - (internal server error)")
	}

	//clear all sessions for user
	if err := tx.Session().ClearAllForUser(requestData.UserID); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xSession_001", "Error deleting user sessions - (internal server error)")
	}

	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &DeleteUserResponse{
		Deleted: true,
	})
}
