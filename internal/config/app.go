package config

import (
	"github.com/namsral/flag"
)

var ShortURLLength = flag.Int("short-url-length", 8, "Length of the short URL")
var ShortURLBase = flag.String("short-url-base", "http://localhost:8080/api/v1", "Prefix of the short URL")
var ShortURLPrefixHTTPS = flag.String("short-url-prefix-https", "https://localhost:8080", "Prefix of the short URL with HTTPS")

var GenerationCollisionRetry = flag.Int("generation-collision-retry", 5, "Number of retries to generate a new short URL if collision occurs")
