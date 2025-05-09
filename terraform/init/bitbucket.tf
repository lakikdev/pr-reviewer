provider "bitbucket" {
  username = var.bitbucket_username
  password = var.bitbucket_password
}

//BITBUCKET SETUP


resource "bitbucket_repository_variable" "bitbucket_repository_owner" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_bitbucket_repo_owner"
  value      = var.bitbucket_repo_owner
  secured    = false
}

resource "bitbucket_repository_variable" "bitbucket_repository_slug" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "TF_VAR_bitbucket_repo_slug"
  value      = var.bitbucket_repo_slug
  secured    = false
}

resource "bitbucket_repository_variable" "bitbucket_username" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "BITBUCKET_USERNAME"
  value      = var.bitbucket_username
  secured    = false
  lifecycle {
    ignore_changes = all
  }
}

resource "bitbucket_repository_variable" "bitbucket_password" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "BITBUCKET_PASSWORD"
  value      = var.bitbucket_password
  secured    = true
  lifecycle {
    ignore_changes = all
  }
}