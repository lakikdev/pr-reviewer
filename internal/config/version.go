package config

import "github.com/namsral/flag"

//Version is the build version of this binary
//Will be changed on build
var CommitHash = flag.String("bitbucket-commit-short", "vHEAD", "Bitbucket commit short hash")
var AppVersion = flag.String("app-version", "0-0-1", "App Version")
