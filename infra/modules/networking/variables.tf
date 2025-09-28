variable "name_prefix" {
  type        = string
  description = "Prefix used for resource names"
}

variable "cidr_block" {
  type        = string
  description = "CIDR block for the VPC"
}

variable "azs" {
  type        = list(string)
  description = "Availability zones to spread subnets across"
}

variable "tags" {
  type        = map(string)
  description = "Common resource tags"
  default     = {}
}
