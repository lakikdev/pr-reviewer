package config

import (
	"github.com/namsral/flag"
)

// DataDirectory is the path used for loading templates/database migrations
var DataDirectory = flag.String("data-directory", "", "Path for loading templates and migration scripts.")
var Environment = flag.String("environment", "local", "Path for loading templates and migration scripts.")

var BasicAuthUsername = flag.String("basic-auth-username", "test", "Username to access all endpoints closed by basic auth")
var BasicAuthPassword = flag.String("basic-auth-password", "test", "Password to access all endpoints closed by basic auth")

var JWTSecretKey = flag.String("jwt-secret-key", "pass-secret-from-env", "Key used to sign JWT tokens")

func GetEnvironment() string {
	return *Environment
}

func IsDevelop() bool {
	return *Environment == "development"
}

func IsLocal() bool {
	return *Environment == "local"
}
