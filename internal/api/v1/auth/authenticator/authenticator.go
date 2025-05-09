package authenticator

import (
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Authenticator struct {
	DB       *database.DB
	authList map[string]IAuthType
}

func New(db *database.DB) Authenticator {
	return Authenticator{db, map[string]IAuthType{}}
}

func (auth *Authenticator) AddAuthType(authType IAuthType) {
	if err, ok := authType.(error); ok {
		logrus.WithError(err).Warn("auth type is disabled")
	}

	auth.authList[authType.Name()] = authType
}

func (a *Authenticator) SignUp(c echo.Context, authName string) (*model.User, error) {
	if authType, ok := a.authList[authName]; ok {
		return a.signUp(c, authType)
	}

	return nil, utils.NewHttpError(nil, http.StatusBadRequest, "Err400xAuthenticator_001", "Error while verifying if request data has all required fields - (unknown auth type)")
}

func (a *Authenticator) SignIn(c echo.Context, authName string) (*model.User, error) {
	if authType, ok := a.authList[authName]; ok {
		return a.signIn(c, authType)
	}

	return nil, utils.NewHttpError(nil, http.StatusBadRequest, "Err400xAuthenticator_001", "Error while verifying if request data has all required fields - (unknown auth type)")
}

func (a *Authenticator) AddAuth(c echo.Context, authName string) (*model.User, error) {
	if authType, ok := a.authList[authName]; ok {
		return a.addAuth(c, authType)
	}

	return nil, utils.NewHttpError(nil, http.StatusBadRequest, "Err400xAuthenticator_001", "Error while verifying if request data has all required fields - (unknown auth type)")
}
