variable "bitbucket_repo_owner" {
  type = string
  description = "Account ID of the Bitbucket repository owner (the organization.)"
}

variable "bitbucket_repo_slug" {
  type = string
  description = "Slug of the repository within the organization"
}

variable "bitbucket_username" {
  type = string
  description = "Username of the Bitbucket app"
}

variable "bitbucket_password" {
  type = string
  description = "Password of the Bitbucket app"
}