variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "name" {
  description = "Prefix for resource names"
  type        = string
}

variable "worker_groups" {
  description = "Worker groups definition list"
  type        = list(object({
    name                 = string
    instance_type        = string
    asg_desired_capacity = number
    asg_min_size         = number
    asg_max_size         = number
  }))
}

variable "subnet_ids" {
  description = "Subnet ids to join to"
  type = list(string)
}
