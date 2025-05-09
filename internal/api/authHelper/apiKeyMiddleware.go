package authHelper

import (
	"context"
	"net/http"

	"pr-reviewer/internal/api/utils"

	"github.com/labstack/echo/v4"
)

type apiKeyContextKeyType struct{}

var apiKeyContextKey apiKeyContextKeyType

func APIKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		r := c.Request()

		newReq, err := CheckAPIKey(r)
		if err != nil {
			return utils.NewHttpError(err, http.StatusUnauthorized, "Err401x9999", "Error while checking api key")
		}

		c.SetRequest(newReq)

		return next(c)
	})
}

func CheckAPIKey(r *http.Request) (*http.Request, error) {
	//extract apiKey from header
	apiKey := GetAPIKey(r)

	//allow to continue without token
	if apiKey == "" {
		return r, nil
	}

	//hash apiKey
	hashedAPIKey, err := HashAPIKey(apiKey)
	if err != nil {
		return r, err
	}

	return r.WithContext(WithAPIKeyContext(r.Context(), hashedAPIKey)), nil
}

// Set principal in context to get it in API
func WithAPIKeyContext(ctx context.Context, hashedAPIKey string) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, hashedAPIKey)
}

func GetAPIKeyFromContext(r *http.Request) string {
	if hashedKey, ok := r.Context().Value(apiKeyContextKey).(string); ok {
		return hashedKey
	}
	return ""
}

func GetAPIKey(r *http.Request) string {
	apiKey := r.Header.Get("x-api-key")
	if apiKey == "" {
		return ""
	}

	return apiKey
}
