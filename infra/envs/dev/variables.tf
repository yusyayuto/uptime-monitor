variable "environment" {
  type    = string
  default = "dev"
}

variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "availability_zones" {
  type    = list(string)
  default = ["ap-northeast-1a", "ap-northeast-1c"]
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "container_image" {
  type    = string
  default = "123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/uptime-monitor:latest"
}

variable "desired_count" {
  type    = number
  default = 2
}

variable "task_cpu" {
  type    = number
  default = 512
}

variable "task_memory" {
  type    = number
  default = 1024
}

variable "certificate_arn" {
  type = string
}

variable "test_listener_port" {
  type    = number
  default = 8443
}

variable "alert_emails" {
  type    = list(string)
  default = []
}

variable "github_oidc_provider_arn" {
  type    = string
  default = ""
}

variable "github_repo" {
  type    = string
  default = ""
}

variable "additional_tags" {
  type    = map(string)
  default = {}
}
