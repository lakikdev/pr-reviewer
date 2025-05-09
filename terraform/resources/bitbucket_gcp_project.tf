resource "bitbucket_repository_variable" "gcloud_service_account_key" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "GCLOUD_SERVICE_ACCOUNT_KEY"
  value      = base64decode(google_service_account_key.bitbucket_pipelines.private_key)
  secured    = true
}
