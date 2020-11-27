variable "name" {
  description = "Prefix for resource names and tags"
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
