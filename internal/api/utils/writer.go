package utils

import (
	"time"

	"pr-reviewer/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type writer struct {
	Logger *logrus.Entry
}

func NewWriter() *writer {
	return &writer{}
}

func NewWriterWithLogger(logger *logrus.Entry) *writer {
	return &writer{
		Logger: logger,
	}
}

// GenericError - represent error structure for generic error (we need this to make all our error response the same)
// swagger:model genericError
type GenericError struct {
	ErrorCode    string      `json:"code"`
	ErrorMessage string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

// WriteError returns a JSON error message and HTTP status code.
func (wr *writer) WriteError(c echo.Context, httpError *HttpError, err error) error {
	if wr.Logger != nil {
		wr.Logger.WithError(err).WithField("errorCode", httpError.Code).Warn(httpError.Message)
	}

	response := GenericError{
		ErrorCode:    httpError.Code,
		ErrorMessage: httpError.Message,
	}

	if config.IsDevelop() || config.IsLocal() {
		if err != nil {
			response.Data = map[string]string{
				"innerErrorMsg": err.Error(),
			}
		}
	}

	setResponseHeaders(c)
	return c.JSON(httpError.StatusCode, response)

}

type GenericResponse struct {
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// WriteJSON returns a JSON data and HTTP status code
func (wr *writer) WriteJSON(c echo.Context, code int, data interface{}) error {
	setResponseHeaders(c)

	response := GenericResponse{
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	return c.JSON(code, response)
}

func setResponseHeaders(c echo.Context) {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000")
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("X-Frame-Options", "DENY")
}
