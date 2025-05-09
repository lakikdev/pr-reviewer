package utils

import (
	"fmt"

	"pr-reviewer/internal/modules/log"
)

type HttpError struct {
	CallerFileName string
	CallerFuncName string

	Code          string
	Message       string
	StatusCode    int
	InternalError error
}

func NewHttpError(internalError error, statusCode int, code string, message string) *HttpError {
	return &HttpError{
		CallerFileName: log.GetCallerFileName(3),
		CallerFuncName: log.GetCallerFuncName(3),
		StatusCode:     statusCode,
		Code:           code,
		Message:        message,
		InternalError:  internalError,
	}
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("StatusCode: %d - Code: %s - Message: %s", he.StatusCode, he.Code, he.Message)
}
