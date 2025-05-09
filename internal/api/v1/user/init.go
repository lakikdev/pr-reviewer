package user

import (
	"pr-reviewer/internal/api/apiHelper"
	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/common"

	"github.com/labstack/echo/v4"
)

type API struct {
	Permissions    authHelper.Permissions
	DataRepository *common.DataRepository
}

func SetAPI(router *echo.Group, permissions authHelper.Permissions, dataRepository *common.DataRepository) {

	api := API{
		Permissions:    permissions,
		DataRepository: dataRepository,
	}

	endpoints := []apiHelper.Endpoint{
		apiHelper.NewEndpoint("/users/:userID", "GET", api.Get, authHelper.Admin, authHelper.MemberIsTarget),
		apiHelper.NewEndpoint("/users/:userID", "PATCH", api.Update, authHelper.Admin, authHelper.MemberIsTarget),
		apiHelper.NewEndpoint("/users/list", "POST", api.List, authHelper.Admin),
		apiHelper.NewEndpoint("/users/:userID", "DELETE", api.Delete, authHelper.Admin, authHelper.MemberIsTarget),
	}

	for _, endpoint := range endpoints {
		router.Add(endpoint.Method, endpoint.Path, permissions.Wrap(endpoint.Func, endpoint.PermissionTypes...))
	}

}
