provider "bitbucket" {
  username = var.bitbucket_username
  password = var.bitbucket_password
}

//BITBUCKET SETUP

resource "bitbucket_repository_variable" "bitbucket_environments" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_ENVIRONMENT"
  value      = each.key
  secured    = false
}

resource "bitbucket_repository_variable" "bitbucket_DOCKER_BUILDKIT" {
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "DOCKER_BUILDKIT"
  value      = 0
  secured    = false
}