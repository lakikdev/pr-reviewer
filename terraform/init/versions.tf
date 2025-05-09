terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.27.0"
    }

    bitbucket = {
      source = "terraform-providers/bitbucket"
    }

    random = {
      source = "hashicorp/random"
    }
  }
}