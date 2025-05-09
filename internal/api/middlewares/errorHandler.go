package middlewares

import (
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/modules/log"

	"github.com/labstack/echo/v4"
)

// GenericError - represent error structure for generic error (we need this to make all our error response the same)
// swagger:model genericError
type GenericError struct {
	ErrorCode    string      `json:"code"`
	ErrorMessage string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	cc, ok := c.(*common.Context)
	if !ok {
		c.Logger().Error("Server Internal Error - Error unknown")
		_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "error unknown"})
	}

	if c.Response().Committed {
		return
	}

	he, ok := err.(*utils.HttpError)
	if !ok {
		he = utils.NewHttpError(err, http.StatusInternalServerError, "500x9999", "Server Internal Error")
	}

	errorResponse := GenericError{
		ErrorCode:    he.Code,
		ErrorMessage: he.Message,
	}

	if config.IsDevelop() || config.IsLocal() {
		if he.InternalError != nil {
			errorResponse.Data = map[string]interface{}{
				"innerErrorMsg": he.InternalError.Error(),
			}
		}
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000")
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("X-Frame-Options", "DENY")

	_ = c.JSON(he.StatusCode, errorResponse)

	if err != nil {
		cc.Logrus().WithFields(log.Fields{
			"file":         he.CallerFileName,
			"func":         he.CallerFuncName,
			"errorCode":    he.Code,
			"errorMessage": he.Message,
			"errorData":    he.InternalError,
		}).Error(he.Message)
	}
}
