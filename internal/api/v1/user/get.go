package user

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetUserRequest struct {
	UserID model.UserID `param:"userID"`
}

// swagger:model GetUserResponse
type GetUserResponse struct {
	User *model.User `json:"user"`
}

/*
		swagger:route GET /users/:userID User userGet

		Get user by id  ----- [Admin, MemberIsTarget]

		<b>Description</b><br>
		Returns user by id

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters)

	 	<br><b>404: Not Found</b><br>
			Err404xUser_001 - Error while getting user - (user not found)

		<br><b>500: Internal Server Error</b><br>
			Err500xUser_001 - Error while getting user - (internal server error)
			Err500xRoles_001 - Error while getting user roles - (internal server error)
			Err500xTransaction_001 - Error while connecting to database - (server internal error)
			Err500xTransaction_002 - Error while committing all changes - (server internal error)

		Parameters:
		+	name: 	     userID
			description: User ID
			in:          path
			type:        string

		Security:
			api_key: []

		Responses:
	  		200: GetUserResponse User data
			default: genericError
*/
func (api *API) Get(c echo.Context) error {
	cc := common.Convert(c)
	ctx := cc.Request().Context()

	var requestData GetUserRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID":       principal.UserID,
		"targetUserID": requestData.UserID,
	})

	//begin transaction
	tx, err := api.DataRepository.DB.BeginTxx(ctx)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	user, err := tx.User().GetByID(requestData.UserID)
	if err != nil {
		if err == dbHelper.ErrItemNotFound {
			return utils.NewHttpError(err, http.StatusNotFound, "Err404xUser_001", "Error while getting user - (user not found)")
		}

		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error while getting user - (internal server error)")
	}

	roles, err := api.Permissions.GetRoles(user.ID)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error while getting user roles - (internal server error)")
	}

	user.Roles = roles

	//commit transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &GetUserResponse{
		User: user,
	})
}
