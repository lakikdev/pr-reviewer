variable "staging_config" {
    default = {
      basic_auth = {
        isActive = false
        username = "lakik"
        password = "pass_2025"
      }
      cloud_run = {
        cpu = 1
        memory = "1Gi"
        maxScale = "100"
      }
      db = {
        isActive = true
        user = {
          username = "root"
          generatePassword = true
          defaultPassword = "password"
        }
        databaseVersion = "POSTGRES_16"
        instanceName = "postgresql-db-staging"
        instanceTier = "db-f1-micro"
        instanceBackupEnabled = true
        instanceAvailabilityType = "ZONAL"
      }
      storage = {
        class = "STANDARD"
      }
    }
}