//DATABSE SETUP
resource "bitbucket_repository_variable" "db_instance_name" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_INSTANCE_NAME"
  value      = google_sql_database_instance.database_instance[each.key].connection_name
  secured    = false
}

resource "bitbucket_repository_variable" "db_host" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_HOST"
  value      = google_sql_database_instance.database_instance[each.key].ip_address.0.ip_address
  secured    = false
}

resource "bitbucket_repository_variable" "db_username" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_USERNAME"
  value      = google_sql_user.user[each.key].name
  secured    = false
}

resource "bitbucket_repository_variable" "db_password" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_PASSWORD"
  value      = google_sql_user.user[each.key].password
  secured    = true
}

resource "bitbucket_repository_variable" "db_name" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_NAME"
  value      = google_sql_database.sql_database[each.key].name
  secured    = false
}

resource "bitbucket_repository_variable" "db_socket_path" {
  for_each = local.environments
  repository = "${var.bitbucket_repo_owner}/${var.bitbucket_repo_slug}"
  key        = "${each.key}_DB_SOCKET_PATH"
  value      = "/cloudsql/${google_sql_database_instance.database_instance[each.key].connection_name}"
  secured    = false
}
