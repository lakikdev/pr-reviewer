
//Cloud run variables
resource "bitbucket_repository_variable" "basic_auth_username" {
  for_each   = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_BASIC_AUTH_USERNAME"
  value      = each.value.basic_auth.username
  secured    = false
}

resource "bitbucket_repository_variable" "basic_auth_password" {
  for_each   = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_BASIC_AUTH_PASSWORD"
  value      = each.value.basic_auth.isActive ? each.value.basic_auth.password : "disabled"
  secured    = each.value.basic_auth.isActive
}
