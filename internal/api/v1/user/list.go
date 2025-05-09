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

// swagger:model ListUserRequest
type ListUserRequest struct {
	model.ListDataParameters
}

// swagger:model ListUserResponse
type ListUserResponse struct {
	Items []*model.User `json:"items"`
	Total *int64        `json:"total"`
}

/*
		swagger:route POST /users/list User userList

		List and search users  ----- [Admin]

		<b>Description</b><br>
		Returns list of users;

		<br><b>400: Bad Request</b><br>
			Err400xBind_001 - Error while decoding request data to JSON - (could not decode parameters) <br>
			Err400xBind_002 - Error while verifying request data - (could not verify parameters) <br>

		<br><b>500: Internal Server Error</b><br>
			Err500xUser_001 - Error getting users list - (internal server error) <br>
			Err500xRoles_001 - Error while getting user roles - (internal server error) <br>
			Err500xTransaction_001 - Error while connecting to database - (server internal error) <br>
			Err500xTransaction_002 - Error while committing all changes - (server internal error) <br>


		Parameters:
		+ 	name:        request
	   		description: Search query specifications, page limit and offset
	  		in:          body
	  		type:        ListUserRequest

		Security:
			api_key: []

		Responses:
	  		200: ListUserResponse List of users returned by the search
			default: genericError
*/
func (api *API) List(c echo.Context) error {
	cc := common.Convert(c)

	principal := authHelper.GetPrincipal(cc.Request())

	cc.Logrus().WithFields(logrus.Fields{
		"UserID": principal.UserID,
	})

	var requestData ListUserRequest
	if err := c.Bind(&requestData); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_001", "Error while decoding request data to JSON - (could not decode parameters)")
	}

	if err := requestData.ListDataParameters.Verify(); err != nil {
		return utils.NewHttpError(err, http.StatusBadRequest, "Err400xBind_002", "Error while verifying request data - (could not verify parameters)")

	}

	//Begin transaction
	tx, err := api.DataRepository.DB.BeginTxx(cc.Request().Context())
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_001", "Error while connecting to database - (server internal error)")
	}
	defer func() { _ = tx.Rollback() }()

	items, total, err := tx.User().List(requestData.ListDataParameters)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xUser_001", "Error getting users list - (internal server error)")
	}

	// Load roles for each user.
	for _, user := range items {
		roles, err := api.Permissions.GetRoles(user.ID)
		if err != nil {
			return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xRoles_001", "Error while getting user roles - (internal server error)")

		}
		user.Roles = roles
	}

	if items == nil {
		items = make([]*model.User, 0)
	}

	cc.Logrus().WithField("total", total)

	//Commit transaction
	if err := tx.Commit(); err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500xTransaction_002", "Error while committing all changes - (server internal error)")
	}

	return utils.NewWriter().WriteJSON(c, http.StatusOK, &ListUserResponse{
		Items: items,
		Total: total,
	})
}
