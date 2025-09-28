variable "name_prefix" {
  type = string
}

variable "dynamodb_table_arns" {
  type = list(string)
}

variable "github_oidc_provider_arn" {
  type    = string
  default = ""
}

variable "github_repo" {
  type    = string
  default = ""
  description = "owner/repo formatted identifier"
}

variable "tags" {
  type    = map(string)
  default = {}
}
