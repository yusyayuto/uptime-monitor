variable "name_prefix" {
  type = string
}

variable "container_image" {
  type = string
}

variable "desired_count" {
  type    = number
  default = 1
}

variable "cpu" {
  type    = number
  default = 256
}

variable "memory" {
  type    = number
  default = 512
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "public_subnet_ids" {
  type = list(string)
}

variable "alb_security_group_id" {
  type = string
}

variable "ecs_security_group_id" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "certificate_arn" {
  type        = string
  description = "ACM certificate for HTTPS listener"
}

variable "test_listener_port" {
  type    = number
  default = 8443
}

variable "task_execution_role_arn" {
  type = string
}

variable "task_role_arn" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
