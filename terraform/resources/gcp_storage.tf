resource "google_storage_bucket" "gcs_bucket_public" {
  for_each = local.environments
  project       = var.gcloud_project_id
  name          = "${var.gcloud_project_id}-public-${each.key}"
  location      = var.gcp_gcs_location
  force_destroy = true
  storage_class = each.value.storage.class
}

resource "google_storage_bucket_iam_member" "member" {
  for_each = local.environments
  bucket = google_storage_bucket.gcs_bucket_public[each.key].name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_storage_bucket" "gcs_bucket_private" {
  for_each = local.environments
  project       = var.gcloud_project_id
  name          = "${var.gcloud_project_id}-private-${each.key}"
  location      = var.gcp_gcs_location
  force_destroy = true
  storage_class = each.value.storage.class
}