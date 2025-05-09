resource "google_service_account" "bitbucket_pipelines" {
  project      = var.gcloud_project_id
  account_id   = "bitbucket-pipelines"
  display_name = "Bitbucket Pipelines Service Account"
}

resource "google_project_iam_member" "bitbucket_pipelines_roles" {
  for_each = toset(var.bitbucket_pipelines_service_account_iam_roles)

  project = var.gcloud_project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.bitbucket_pipelines.email}"
}


resource "google_service_account_key" "bitbucket_pipelines" {
  service_account_id = google_service_account.bitbucket_pipelines.name
}

variable "bitbucket_pipelines_service_account_iam_roles" {
  type        = list(string)
  description = "List of IAM roles to assign to the Bitbucket Pipelines service account."
  default     = [
    "roles/compute.storageAdmin",
    "roles/cloudbuild.builds.editor",
    "roles/run.admin",
    "roles/storage.admin",
    "roles/iam.serviceAccountUser",
    "roles/viewer",
    "roles/artifactregistry.admin",
    "roles/artifactregistry.createOnPushRepoAdmin"
  ]
}
