package api

import (
	"fmt"
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/api/middlewares"
	v1 "pr-reviewer/internal/api/v1"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/database"

	v1Analyze "pr-reviewer/internal/api/v1/analyze"
	v1APIKey "pr-reviewer/internal/api/v1/apiKey"
	v1Auth "pr-reviewer/internal/api/v1/auth"
	v1Role "pr-reviewer/internal/api/v1/role"
	v1User "pr-reviewer/internal/api/v1/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter provide a handler API service.
func NewRouter(db *database.DB) (*echo.Echo, error) {
	dataRepository := common.NewDataRepository(db)
	permissions := authHelper.NewPermissions(dataRepository)

	router := echo.New()

	router.GET("/version", v1.VersionHandler)
	apiRouter := router.Group("/api/v1")

	router.Use(middlewares.InitContext)
	router.Use(middlewares.Logger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Accept", "Accept-Language", "Content-Language", "Origin", "Authorization", "Content-Type"},
		AllowOrigins: []string{"*"},
	}))

	if config.IsDevelop() || config.IsLocal() {
		SetSwagger(router)
	} else { //All other environments
		router.Use(middlewares.LimitVisitor)
	}
	router.Use(middlewares.DirectToHttps)
	router.Use(authHelper.AuthorizationTokenMiddleware)
	router.Use(authHelper.APIKeyMiddleware)

	v1Auth.SetAPI(apiRouter, permissions, dataRepository)
	v1User.SetAPI(apiRouter, permissions, dataRepository)
	v1Role.SetAPI(apiRouter, permissions, dataRepository)
	v1Analyze.SetAPI(apiRouter, permissions, dataRepository)
	v1APIKey.SetAPI(apiRouter, permissions, dataRepository)

	return router, nil
}

func SetSwagger(router *echo.Echo) {
	fmt.Println("Swagger is enabled")
	router.File("/api/v1/swagger.json", fmt.Sprintf("%s/swagger.json", *config.DataDirectory))
	router.Static("/swagger", fmt.Sprintf("%s/swaggerui", *config.DataDirectory))
}
