variable "gcp_org_id" {
  type = string
  description = "Google Cloud Organization ID for the project to be created within"
}

variable "gcp_billing_account" {
  type = string
  description = "Google Cloud Billing Account ID to be associated with the created project"
}

variable "gcp_project_name" {
  type = string
  description = "Name of the Google Cloud Platform Project to be created."
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

variable "google_project_services" {
  type = list(string)
  description = "Google Project Services to be enabled"
  default = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "containerregistry.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "stackdriver.googleapis.com",
    "clouderrorreporting.googleapis.com",
    "vpcaccess.googleapis.com",
    "redis.googleapis.com",
    "artifactregistry.googleapis.com"
  ]
}

variable "google_credentials" {
  type = string
  description = "Credentioals to access GCP"
}
