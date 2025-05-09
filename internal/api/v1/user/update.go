package user

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/helper"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// swagger:model UpdateUserRequest
type UpdateUserRequest struct {
	model.User

	Password *string `json:"password"`
}

// swagger:model UpdateUserResponse
type UpdateUserResponse struct {
	*model.User
}

/*
		swagger:route PATCH /users/:userID User userUpdate

		Update user  ----- [Admin, MemberIsTarget]

		<b>Description</b><br>
		Updates user by ID and returns updated user

	    <br><b>400: Bad Request</b><br>
			Err400x001 - Error while decoding request data to JSON - (could not decode parameters)
			Err400x002 - Error while verifying if request data has all required fields - (not all fields found)
			Err400x003 - User not authorized to edit password - (user not authorized)

		<br><b>404: Not Found</b><br>
			Err404x001 - Error while getting user by email - (user not found)
			Err404x002 - Error while updating user - (user not found)

		<br><b>500: Internal Server Error</b><br>
			Err500x001 - Error while getting user by email - (internal server error)
			Err500x002 - Error while hashing and setting new password - (internal server error)
			Err500x003 - Error while updating metadata - (internal server error)

		Parameters:
		+ 	name:        userID
		   	description: User ID
		   	in:          path
		   	type:        string
		+	name: request
			description: User model
			in: body
			type: UpdateUserRequest

		Security:
			api_key: []

	 	Responses:
	  		200: UpdateUserResponse User is updated successfully, user model
			default: genericError
*/
func (api *API) Update(c echo.Context) error {
	cc := common.Convert(c)

	userID := model.UserID(cc.Param("userID"))
	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID":       principal.UserID,
		"targetUserID": userID,
	})

	//begin transaction
	tx, err := api.DataRepository.DB.BeginTxx(cc.Request().Context())
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	loadedUser, err := tx.User().GetByID(userID)
	if err != nil || loadedUser == nil {
		if err == dbHelper.ErrItemNotFound {
			return utils.NewHttpError(err, http.StatusNotFound, "Err404xUser_001", "Error while getting user by email - (user not found)")
		}

		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error while getting user by email - (internal server error)")
	}

	var request UpdateUserRequest
	request.User, err = helper.DeepCopy(*loadedUser)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error while update user - (server internal error)")
	}
	if err := c.Bind(&request); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	if err := request.Verify(); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_002", "Error while verifying if request data has all required fields - (not all fields found)")

	}

	if err := tx.User().Update(&request.User); err != nil {
		if err == dbHelper.ErrItemNotFound {
			return utils.NewHttpError(err, http.StatusNotFound, "Err404xUser_002", "Error while updating user - (user not found)")
		}

		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_003", "Error while updating metadata - (internal server error)")
	}

	//commit transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &UpdateUserResponse{
		&request.User,
	})
}
