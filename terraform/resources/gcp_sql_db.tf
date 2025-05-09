resource "random_password" "sql_root_password" {
  for_each = local.environments
  length   = 16
  special  = true
  keepers = {
    ami_id = var.gcloud_project_id
  }
}

resource "google_sql_database_instance" "database_instance" {
  for_each = local.environments
  project          = var.gcloud_project_id
  name             = "postgresql-db-${each.key}"
  database_version = each.value.db.databaseVersion
  region           = var.gcp_region

  timeouts {
    create = "2h"
    update = "2h"
    delete = "20m"
  }

  settings {
    tier              = each.value.db.instanceTier
    availability_type = each.value.db.instanceAvailabilityType

    backup_configuration {
      enabled            = each.value.db.instanceBackupEnabled
      point_in_time_recovery_enabled = each.value.db.instanceBackupEnabled
      start_time         = "00:00"
      transaction_log_retention_days = "3"

    }

    ip_configuration {
      dynamic "authorized_networks" {
        for_each = var.pipelines_ip_addresses
        iterator = pipelines_ip_addresses

        content {
          name  = "Bitbucket Pipelines IP ${pipelines_ip_addresses.key}"
          value = pipelines_ip_addresses.value
        }
      }
    }
  }
}


resource "google_sql_user" "user" {
  for_each = local.environments
  project  = var.gcloud_project_id
  instance = google_sql_database_instance.database_instance[each.key].name
  name     = each.value.db.user.username
  password = each.value.db.user.generatePassword ? random_password.sql_root_password[each.key].result : each.value.db.user.defaultPassword
}



resource "google_sql_database" "sql_database" {
  for_each = local.environments
  project  = var.gcloud_project_id
  instance = google_sql_database_instance.database_instance[each.key].name
  name     = "main_db"
}


output "db_password" {
  value      = values(google_sql_user.user)[*].password
  sensitive = true
}