package v1

import (
	"net/http"
	"os"

	"pr-reviewer/internal/api/utils"

	"github.com/labstack/echo/v4"
)

//API for returning version
//When up starts, we set version and than use it if necessary .

// represents the up version.
//
//swagger:model ServerVersion
type ServerVersion struct {
	CommitHash string `json:"commit"`
	AppVersion string `json:"version"`
}

// Marshaled JSON
var versionJSON ServerVersion

func init() {
	versionJSON = ServerVersion{
		CommitHash: os.Getenv("BITBUCKET_COMMIT_SHORT"),
		AppVersion: os.Getenv("APP_VERSION"),
	}
}

// swagger:route GET /version Version version
//
// # App Version
//
// Responses:
//
//	200: ServerVersion Return current Git Header version used

// swagger:route GET /api/v1/version Version versionV1
//
// # App Version
//
// Responses:
//
//	200: ServerVersion Return current Git Header version used
func VersionHandler(c echo.Context) error {
	return utils.NewWriter().WriteJSON(c, http.StatusOK, versionJSON)
}
