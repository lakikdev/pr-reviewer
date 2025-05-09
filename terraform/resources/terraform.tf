terraform {
   backend "gcs" {
    bucket  = "terraform_state_bucket"
    prefix  = "terraform/state"
    credentials = "service-account-key.json"
  }
}
