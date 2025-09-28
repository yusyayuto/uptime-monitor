variable "name_prefix" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "service_name" {
  type = string
}

variable "blue_target_group_name" {
  type = string
}

variable "green_target_group_name" {
  type = string
}

variable "prod_listener_arn" {
  type = string
}

variable "test_listener_arn" {
  type = string
}

variable "codedeploy_role_arn" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
