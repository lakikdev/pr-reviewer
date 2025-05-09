resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "terraform_state" {
  project       = google_project.project.project_id
  name          = "${random_id.bucket_prefix.hex}-bucket-tfstate"
  force_destroy = true
  location      = var.gcp_gcs_location
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}