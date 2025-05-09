package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"pr-reviewer/internal/config"

	"github.com/labstack/echo/v4"
)

func DirectToHttps(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		r := c.Request()
		if config.IsLocal() {
			return next(c)
		}

		if r.URL.Scheme == "https" || strings.HasPrefix(r.Proto, "HTTPS") || r.Header.Get("X-Forwarded-Proto") == "https" {
			fmt.Printf("Continue to API: %s/%s\n", r.Host, r.URL.Path)
			return next(c)
		} else {
			target := "https://" + r.Host + r.URL.Path
			fmt.Printf("redirect to: %s\n", target)

			http.Redirect(c.Response().Writer, r, target, http.StatusTemporaryRedirect)
			return next(c)
		}
	})
}
