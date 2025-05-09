//GENERAL SETUP
resource "bitbucket_repository_variable" "gcloud_project_id" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_gcloud_project_id"
  value      = google_project.project.project_id
  secured    = false
}

resource "bitbucket_repository_variable" "gcp_region" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_gcp_region"
  value      = var.gcp_region
  secured    = false
}

resource "bitbucket_repository_variable" "gcp_gcs_location" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_gcp_gcs_location"
  value      = var.gcp_gcs_location
  secured    = false
}

resource "bitbucket_repository_variable" "gcp_gcr_location" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_gcp_gcr_location"
  value      = var.gcp_gcr_location
  secured    = false
}

resource "bitbucket_repository_variable" "terraform_state_bucket" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_TERRAFORM_STATE_BUCKET"
  value      = google_storage_bucket.terraform_state.name
  secured    = false
}