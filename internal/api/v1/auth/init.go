package auth

import (
	"pr-reviewer/internal/api/apiHelper"
	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/v1/auth/authenticator"
	"pr-reviewer/internal/api/v1/auth/authenticator/credentialsAuth"
	"pr-reviewer/internal/api/v1/auth/authenticator/udidAuth"
	"pr-reviewer/internal/common"

	"github.com/labstack/echo/v4"
)

// AuthAPI - provides REST for Auth
type API struct {
	Permissions    authHelper.Permissions
	Authenticator  authenticator.Authenticator
	DataRepository *common.DataRepository
}

func SetAPI(router *echo.Group, permissions authHelper.Permissions, dataRepository *common.DataRepository) {

	authenticator := authenticator.New(dataRepository.DB)
	authenticator.AddAuthType(udidAuth.New(dataRepository.DB))
	authenticator.AddAuthType(credentialsAuth.New(dataRepository.DB))

	api := API{
		Permissions:    permissions,
		Authenticator:  authenticator,
		DataRepository: dataRepository,
	}

	endpoints := []apiHelper.Endpoint{
		apiHelper.NewEndpoint("/signup", "POST", api.SignUp, authHelper.Any),
		apiHelper.NewEndpoint("/signin", "POST", api.SignIn, authHelper.Any),
		apiHelper.NewEndpoint("/addAuth", "POST", api.AddAuth, authHelper.Member),
		apiHelper.NewEndpoint("/refreshToken", "POST", api.RefreshToken, authHelper.Any),
	}

	for _, endpoint := range endpoints {
		router.Add(endpoint.Method, endpoint.Path, permissions.Wrap(endpoint.Func, endpoint.PermissionTypes...))
	}
}
