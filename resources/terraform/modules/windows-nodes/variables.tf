variable "name" {
  description = "Prefix for resource names and tags"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet ids to join to"
  type = list(string)
}

variable "vpc_id" {
  description = "Subnet ids to join to"
  type = list(string)
}