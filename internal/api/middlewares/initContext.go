package middlewares

import (
	"pr-reviewer/internal/common"

	"github.com/labstack/echo/v4"
)

func InitContext(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		cc := common.NewContext(c)
		cc.Logrus().WithField("requestID", cc.RequestID())
		return next(cc)
	})
}
