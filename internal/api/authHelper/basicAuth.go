package authHelper

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func BasicAuth(configUsername, configPassword string) echo.MiddlewareFunc {
	return echoMiddleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(configUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(configPassword)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
