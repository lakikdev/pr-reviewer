locals {
  environments = {
    development = var.dev_config
    test = var.test_config
    staging = var.staging_config
    production = var.prod_config
  }
}

variable "pipelines_ip_addresses" {
  type = list(string)
  description = "Bitbucket Pipelines IP addresses for remote connections to SQL instances upon deployment"
  default = [
    "34.199.54.113/32",
    "34.232.25.90/32",
    "34.232.119.183/32",
    "34.236.25.177/32",
    "35.171.175.212/32",
    "52.54.90.98/32",
    "52.202.195.162/32",
    "52.203.14.55/32",
    "52.204.96.37/32",
    "34.218.156.209/32",
    "34.218.168.212/32",
    "52.41.219.63/32",
    "35.155.178.254/32",
    "35.160.177.10/32",
    "34.216.18.129/32"
  ]
}