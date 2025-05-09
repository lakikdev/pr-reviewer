package apiHelper

import (
	"pr-reviewer/internal/api/authHelper"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Path            string
	Method          string
	Func            echo.HandlerFunc
	PermissionTypes []authHelper.PermissionType
}

func NewEndpoint(path string, method string, handlerFunc echo.HandlerFunc, permissionTypes ...authHelper.PermissionType) Endpoint {
	return Endpoint{path, method, handlerFunc, permissionTypes}
}
