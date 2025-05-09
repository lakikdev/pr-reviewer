terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }

    bitbucket = {
      source = "terraform-providers/bitbucket"
    }

    random = {
      source = "hashicorp/random"
    }
  }
}