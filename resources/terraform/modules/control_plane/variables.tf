variable "name" {
  description = "Prefix for resource names"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet ids to join to"
  type = list(string)
}

variable "k8s_version" {
  description = "Kubernetes version to install"
  type        = string
}