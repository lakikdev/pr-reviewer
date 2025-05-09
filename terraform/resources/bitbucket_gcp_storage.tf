//GCP STORAGE
resource "bitbucket_repository_variable" "storage_bucket_public" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_PUBLIC_STORAGE_BUCKET"
  value      = google_storage_bucket.gcs_bucket_public[each.key].name
  secured    = false
}


//GCP STORAGE
resource "bitbucket_repository_variable" "storage_bucket_private" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_PRIVATE_STORAGE_BUCKET"
  value      = google_storage_bucket.gcs_bucket_private[each.key].name
  secured    = false
}
