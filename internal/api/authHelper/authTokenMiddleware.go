package authHelper

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

type principalContextKeyType struct{}

var principalContextKey principalContextKeyType

func AuthorizationTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		r := c.Request()

		newReq, err := CheckToken(r)
		if err != nil {
			return utils.NewHttpError(err, http.StatusUnauthorized, "Err401x9999", "Error while checking access token")
		}

		c.SetRequest(newReq)

		return next(c)
	})
}

func CheckToken(r *http.Request) (*http.Request, error) {
	//extract token from header
	token, err := GetToken(r)
	if err != nil {
		return r, err
	}

	//allow to continue without token
	if token == "" {
		return r, nil
	}

	principal, err := VerifyToken(token)
	if err != nil {
		return r, err
	}

	return r.WithContext(WithPrincipalContext(r.Context(), *principal)), nil
}

// Set principal in context to get it in API
func WithPrincipalContext(ctx context.Context, principal model.Principal) context.Context {
	return context.WithValue(ctx, principalContextKey, principal)
}

func GetToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", nil
	}

	tokenParts := strings.SplitN(token, " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" || len(tokenParts[1]) == 0 {
		//cms passes basic token on /signin API so we skip check if it do that
		if strings.ToLower(tokenParts[0]) == "basic" {
			return "", nil
		}
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return tokenParts[1], nil
}

func GetPrincipal(r *http.Request) model.Principal {
	if principal, ok := r.Context().Value(principalContextKey).(model.Principal); ok {
		return principal
	}
	return model.NilPrincipal
}
