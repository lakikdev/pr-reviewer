resource "google_cloud_run_service" "cloud_run" {
  for_each = local.environments
  project  = var.gcloud_project_id
  name     = "${var.gcloud_project_id}-${each.key}"
  location = var.gcp_region

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        resources {
          limits = {
            cpu    = each.value.cloud_run.cpu
            memory = each.value.cloud_run.memory
          }
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = each.value.cloud_run.maxScale
        "run.googleapis.com/cloudsql-instances" = each.value.db.isActive ? "${var.gcloud_project_id}:${var.gcp_region}:${google_sql_database_instance.database_instance[each.key].name}" : "null"
      }
    }
  }

  autogenerate_revision_name = true

  # lifecycle {
  #   ignore_changes = all
  # }
}

# Create public access
data "google_iam_policy" "noauth-cloud_run" {
  for_each = local.environments
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

# Enable public access on Cloud Run service
resource "google_cloud_run_service_iam_policy" "noauth-cloud_run" {
  for_each = local.environments
  location    = google_cloud_run_service.cloud_run[each.key].location
  project     = google_cloud_run_service.cloud_run[each.key].project
  service     = google_cloud_run_service.cloud_run[each.key].name
  policy_data = data.google_iam_policy.noauth-cloud_run[each.key].policy_data
}

