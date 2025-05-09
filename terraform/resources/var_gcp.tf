
variable "gcloud_project_id" {
  type = string
  description = "Project ID of the Google Project we created in the previous step"
}

variable "gcp_region" {
  type = string
  description = "Google Cloud Region that the project should be created within"
  default = "us-central1"
}

variable "gcp_gcs_location" {
  type = string
  description = "Google Cloud Storage location that should be used (e.g. for Strapi uploads)"
  default = "us"
}

variable "gcp_gcr_location" {
  type = string
  description = "Google Container Registry location that should be used for Docker images"
  default = "us"
}

variable "google_credentials" {
  type = string
  description = "Credentioals to access GCP"
}
