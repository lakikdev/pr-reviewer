package config


import "github.com/namsral/flag"

var BitbucketBaseURL = flag.String("bitbucket-base-url", "https://bitbucket.org", "Bitbucket base URL.")
var BitbucketClientID = flag.String("bitbucket-client-id", "", "Bitbucket client ID.")
var BitbucketClientSecret = flag.String("bitbucket-client-secret", "", "Bitbucket client secret.")
var BitbucketUsername = flag.String("bitbucket-username", "", "Bitbucket username.")
var BitbucketPassword = flag.String("bitbucket-password", "", "Bitbucket password.")