package config

import "github.com/namsral/flag"

var GooglePublicStorageBucket = flag.String("public-storage-bucket", "", "Google Cloud Storage bucket for storing files.")
var GoogleBucketURL = flag.String("storage-bucket-url", "https://storage.googleapis.com/", "Google Cloud Storage bucket for storing files.")

var GoogleProjectID = flag.String("gcloud-project-id", "", "Google Cloud Project ID.")
var GoogleFunctionPrefix = flag.String("cloud-func-prefix", "", "Google Cloud Function prefix.")
